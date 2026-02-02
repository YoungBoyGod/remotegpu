package adapter

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOAdapter MinIO存储适配器
type MinIOAdapter struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOAdapter 创建MinIO适配器
func NewMinIOAdapter(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIOAdapter, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOAdapter{
		client:     client,
		bucketName: bucket,
	}, nil
}

// Upload 上传文件
func (a *MinIOAdapter) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64) error {
	_, err := a.client.PutObject(ctx, a.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

// Download 下载文件
func (a *MinIOAdapter) Download(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	return a.client.GetObject(ctx, a.bucketName, objectKey, minio.GetObjectOptions{})
}

// Delete 删除文件
func (a *MinIOAdapter) Delete(ctx context.Context, objectKey string) error {
	return a.client.RemoveObject(ctx, a.bucketName, objectKey, minio.RemoveObjectOptions{})
}

// Exists 检查文件是否存在
func (a *MinIOAdapter) Exists(ctx context.Context, objectKey string) (bool, error) {
	_, err := a.client.StatObject(ctx, a.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		// 检查是否是对象不存在的错误
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetSize 获取文件大小
func (a *MinIOAdapter) GetSize(ctx context.Context, objectKey string) (int64, error) {
	stat, err := a.client.StatObject(ctx, a.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return 0, err
	}
	return stat.Size, nil
}

// List 列出指定前缀下的所有对象
func (a *MinIOAdapter) List(ctx context.Context, prefix string) ([]string, error) {
	var objects []string
	objectCh := a.client.ListObjects(ctx, a.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object.Key)
	}

	return objects, nil
}

// Copy 复制文件
func (a *MinIOAdapter) Copy(ctx context.Context, srcKey, dstKey string) error {
	src := minio.CopySrcOptions{
		Bucket: a.bucketName,
		Object: srcKey,
	}
	dst := minio.CopyDestOptions{
		Bucket: a.bucketName,
		Object: dstKey,
	}
	_, err := a.client.CopyObject(ctx, dst, src)
	return err
}

// GetPresignedURL 生成预签名URL
func (a *MinIOAdapter) GetPresignedURL(ctx context.Context, objectKey string, expireSeconds int) (string, error) {
	url, err := a.client.PresignedGetObject(ctx, a.bucketName, objectKey,
		time.Duration(expireSeconds)*time.Second, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// HealthCheck 健康检查
func (a *MinIOAdapter) HealthCheck(ctx context.Context) error {
	exists, err := a.client.BucketExists(ctx, a.bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return minio.ErrorResponse{Code: "NoSuchBucket"}
	}
	return nil
}
