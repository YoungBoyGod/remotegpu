package adapter

import (
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSAdapter 阿里云OSS存储适配器
type OSSAdapter struct {
	client     *oss.Client
	bucket     *oss.Bucket
	bucketName string
}

// NewOSSAdapter 创建阿里云OSS适配器
func NewOSSAdapter(endpoint, accessKeyID, accessKeySecret, bucketName string) (*OSSAdapter, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &OSSAdapter{
		client:     client,
		bucket:     bucket,
		bucketName: bucketName,
	}, nil
}

// Upload 上传文件
func (a *OSSAdapter) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64) error {
	return a.bucket.PutObject(objectKey, reader)
}

// Download 下载文件
func (a *OSSAdapter) Download(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	return a.bucket.GetObject(objectKey)
}

// Delete 删除文件
func (a *OSSAdapter) Delete(ctx context.Context, objectKey string) error {
	return a.bucket.DeleteObject(objectKey)
}

// Exists 检查文件是否存在
func (a *OSSAdapter) Exists(ctx context.Context, objectKey string) (bool, error) {
	return a.bucket.IsObjectExist(objectKey)
}

// GetSize 获取文件大小
func (a *OSSAdapter) GetSize(ctx context.Context, objectKey string) (int64, error) {
	meta, err := a.bucket.GetObjectMeta(objectKey)
	if err != nil {
		return 0, err
	}

	// 从元数据中获取Content-Length
	contentLength := meta.Get("Content-Length")
	if contentLength == "" {
		return 0, nil
	}

	var size int64
	_, err = io.ReadFull(io.LimitReader(nil, 0), nil)
	return size, err
}

// List 列出指定前缀下的所有对象
func (a *OSSAdapter) List(ctx context.Context, prefix string) ([]string, error) {
	var objects []string
	marker := ""

	for {
		result, err := a.bucket.ListObjects(oss.Prefix(prefix), oss.Marker(marker))
		if err != nil {
			return nil, err
		}

		for _, object := range result.Objects {
			objects = append(objects, object.Key)
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}

	return objects, nil
}

// Copy 复制文件
func (a *OSSAdapter) Copy(ctx context.Context, srcKey, dstKey string) error {
	_, err := a.bucket.CopyObject(srcKey, dstKey)
	return err
}

// GetPresignedURL 生成预签名URL
func (a *OSSAdapter) GetPresignedURL(ctx context.Context, objectKey string, expireSeconds int) (string, error) {
	return a.bucket.SignURL(objectKey, oss.HTTPGet, int64(expireSeconds))
}

// HealthCheck 健康检查
func (a *OSSAdapter) HealthCheck(ctx context.Context) error {
	_, err := a.client.GetBucketInfo(a.bucketName)
	return err
}
