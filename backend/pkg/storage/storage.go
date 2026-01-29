package storage

import (
	"context"
	"io"
	"time"
)

// FileInfo 文件信息
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag,omitempty"`
	IsDir        bool      `json:"is_dir"`
}

// UploadOptions 上传选项
type UploadOptions struct {
	ContentType string
	Metadata    map[string]string
}

// Storage 存储接口
type Storage interface {
	// Name 返回存储后端名称
	Name() string

	// Type 返回存储类型
	Type() string

	// Upload 上传文件
	Upload(ctx context.Context, path string, reader io.Reader, size int64, opts *UploadOptions) error

	// Download 下载文件
	Download(ctx context.Context, path string) (io.ReadCloser, *FileInfo, error)

	// Delete 删除文件
	Delete(ctx context.Context, path string) error

	// Exists 检查文件是否存在
	Exists(ctx context.Context, path string) (bool, error)

	// List 列出目录下的文件
	List(ctx context.Context, prefix string) ([]FileInfo, error)

	// GetURL 获取文件访问URL（可选支持预签名URL）
	GetURL(ctx context.Context, path string, expires time.Duration) (string, error)

	// Stat 获取文件信息
	Stat(ctx context.Context, path string) (*FileInfo, error)

	// Copy 复制文件
	Copy(ctx context.Context, srcPath, dstPath string) error

	// Move 移动文件
	Move(ctx context.Context, srcPath, dstPath string) error
}
