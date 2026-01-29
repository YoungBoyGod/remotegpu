package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage S3兼容存储（支持AWS S3、RustFS、MinIO等）
type S3Storage struct {
	name     string
	stype    string // rustfs 或 s3
	client   *minio.Client
	bucket   string
	endpoint string
}

// NewS3Storage 创建S3存储实例
func NewS3Storage(name, stype, endpoint, accessKey, secretKey, bucket, region string) (*S3Storage, error) {
	useSSL := false
	if len(endpoint) > 5 && endpoint[:5] == "https" {
		useSSL = true
	}

	// 移除协议前缀
	host := endpoint
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		host = endpoint[7:]
	} else if len(endpoint) > 8 && endpoint[:8] == "https://" {
		host = endpoint[8:]
	}

	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("创建S3客户端失败: %w", err)
	}

	return &S3Storage{
		name:     name,
		stype:    stype,
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
	}, nil
}

func (s *S3Storage) Name() string { return s.name }
func (s *S3Storage) Type() string { return s.stype }

func (s *S3Storage) Upload(ctx context.Context, path string, reader io.Reader, size int64, opts *UploadOptions) error {
	putOpts := minio.PutObjectOptions{}
	if opts != nil {
		putOpts.ContentType = opts.ContentType
		putOpts.UserMetadata = opts.Metadata
	}

	_, err := s.client.PutObject(ctx, s.bucket, path, reader, size, putOpts)
	if err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}
	return nil
}

func (s *S3Storage) Download(ctx context.Context, path string) (io.ReadCloser, *FileInfo, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("获取文件失败: %w", err)
	}

	info, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	fileInfo := &FileInfo{
		Name:         info.Key,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
		ETag:         info.ETag,
	}

	return obj, fileInfo, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func (s *S3Storage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *S3Storage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	var files []FileInfo
	opts := minio.ListObjectsOptions{Prefix: prefix, Recursive: false}

	for obj := range s.client.ListObjects(ctx, s.bucket, opts) {
		if obj.Err != nil {
			return nil, fmt.Errorf("列出文件失败: %w", obj.Err)
		}
		files = append(files, FileInfo{
			Name:         obj.Key,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
			LastModified: obj.LastModified,
			ETag:         obj.ETag,
		})
	}
	return files, nil
}

func (s *S3Storage) GetURL(ctx context.Context, path string, expires time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.bucket, path, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}
	return url.String(), nil
}

func (s *S3Storage) Stat(ctx context.Context, path string) (*FileInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}
	return &FileInfo{
		Name:         info.Key,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
		ETag:         info.ETag,
	}, nil
}

func (s *S3Storage) Copy(ctx context.Context, srcPath, dstPath string) error {
	src := minio.CopySrcOptions{Bucket: s.bucket, Object: srcPath}
	dst := minio.CopyDestOptions{Bucket: s.bucket, Object: dstPath}

	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

func (s *S3Storage) Move(ctx context.Context, srcPath, dstPath string) error {
	if err := s.Copy(ctx, srcPath, dstPath); err != nil {
		return err
	}
	return s.Delete(ctx, srcPath)
}
