package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage 本地文件存储
type LocalStorage struct {
	name     string
	basePath string
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(name, basePath string) (*LocalStorage, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %w", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	return &LocalStorage{
		name:     name,
		basePath: absPath,
	}, nil
}

func (l *LocalStorage) Name() string { return l.name }
func (l *LocalStorage) Type() string { return "local" }

func (l *LocalStorage) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}

func (l *LocalStorage) Upload(ctx context.Context, path string, reader io.Reader, size int64, opts *UploadOptions) error {
	fullPath := l.fullPath(path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

func (l *LocalStorage) Download(ctx context.Context, path string) (io.ReadCloser, *FileInfo, error) {
	fullPath := l.fullPath(path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("打开文件失败: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	fileInfo := &FileInfo{
		Name:         filepath.Base(path),
		Size:         info.Size(),
		LastModified: info.ModTime(),
		IsDir:        info.IsDir(),
	}

	return file, fileInfo, nil
}

func (l *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := l.fullPath(path)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func (l *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := l.fullPath(path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (l *LocalStorage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	fullPath := l.fullPath(prefix)
	var files []FileInfo

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:         entry.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
			IsDir:        entry.IsDir(),
		})
	}

	return files, nil
}

func (l *LocalStorage) GetURL(ctx context.Context, path string, expires time.Duration) (string, error) {
	return "file://" + l.fullPath(path), nil
}

func (l *LocalStorage) Stat(ctx context.Context, path string) (*FileInfo, error) {
	fullPath := l.fullPath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}
	return &FileInfo{
		Name:         filepath.Base(path),
		Size:         info.Size(),
		LastModified: info.ModTime(),
		IsDir:        info.IsDir(),
	}, nil
}

func (l *LocalStorage) Copy(ctx context.Context, srcPath, dstPath string) error {
	src, err := os.Open(l.fullPath(srcPath))
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer src.Close()

	dstFullPath := l.fullPath(dstPath)
	if err := os.MkdirAll(filepath.Dir(dstFullPath), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	dst, err := os.Create(dstFullPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

func (l *LocalStorage) Move(ctx context.Context, srcPath, dstPath string) error {
	srcFullPath := l.fullPath(srcPath)
	dstFullPath := l.fullPath(dstPath)

	if err := os.MkdirAll(filepath.Dir(dstFullPath), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	if err := os.Rename(srcFullPath, dstFullPath); err != nil {
		return fmt.Errorf("移动文件失败: %w", err)
	}
	return nil
}
