package adapter

import (
	"context"
	"io"
)

// StorageAdapter 存储适配器接口
type StorageAdapter interface {
	// Upload 上传文件
	// objectKey: 对象键(路径)
	// reader: 文件内容读取器
	// size: 文件大小
	Upload(ctx context.Context, objectKey string, reader io.Reader, size int64) error

	// Download 下载文件
	// objectKey: 对象键(路径)
	// 返回文件内容读取器
	Download(ctx context.Context, objectKey string) (io.ReadCloser, error)

	// Delete 删除文件
	// objectKey: 对象键(路径)
	Delete(ctx context.Context, objectKey string) error

	// Exists 检查文件是否存在
	// objectKey: 对象键(路径)
	Exists(ctx context.Context, objectKey string) (bool, error)

	// GetSize 获取文件大小
	// objectKey: 对象键(路径)
	GetSize(ctx context.Context, objectKey string) (int64, error)

	// List 列出指定前缀下的所有对象
	// prefix: 对象键前缀
	List(ctx context.Context, prefix string) ([]string, error)

	// Copy 复制文件
	// srcKey: 源对象键
	// dstKey: 目标对象键
	Copy(ctx context.Context, srcKey, dstKey string) error

	// GetPresignedURL 生成预签名URL
	// objectKey: 对象键(路径)
	// expireSeconds: 过期时间(秒)
	GetPresignedURL(ctx context.Context, objectKey string, expireSeconds int) (string, error)

	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error
}
