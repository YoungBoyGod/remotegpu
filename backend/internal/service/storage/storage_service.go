package storage

import (
	"context"
	"time"

	pkgStorage "github.com/YoungBoyGod/remotegpu/pkg/storage"
)

type StorageService struct {
	storageMgr *pkgStorage.Manager
}

func NewStorageService(mgr *pkgStorage.Manager) *StorageService {
	return &StorageService{
		storageMgr: mgr,
	}
}

// BackendDetail 存储后端详情
type BackendDetail struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
}

// StorageStats 存储使用统计
type StorageStats struct {
	BackendName string `json:"backend_name"`
	FileCount   int    `json:"file_count"`
	TotalSize   int64  `json:"total_size"`
}

// ListBackends 获取存储池列表
func (s *StorageService) ListBackends() []BackendDetail {
	infos := s.storageMgr.List()
	result := make([]BackendDetail, 0, len(infos))
	for _, info := range infos {
		result = append(result, BackendDetail{
			Name:      info.Name,
			Type:      info.Type,
			IsDefault: info.IsDefault,
		})
	}
	return result
}

// GetBackendStats 获取指定存储后端的使用统计
func (s *StorageService) GetBackendStats(ctx context.Context, backendName string) (*StorageStats, error) {
	backend, err := s.storageMgr.Get(backendName)
	if err != nil {
		return nil, err
	}

	files, err := backend.List(ctx, "")
	if err != nil {
		return nil, err
	}

	var totalSize int64
	for _, f := range files {
		if !f.IsDir {
			totalSize += f.Size
		}
	}

	return &StorageStats{
		BackendName: backendName,
		FileCount:   len(files),
		TotalSize:   totalSize,
	}, nil
}

// ListFiles 列出指定存储后端下的文件
func (s *StorageService) ListFiles(ctx context.Context, backendName, prefix string) ([]pkgStorage.FileInfo, error) {
	backend, err := s.storageMgr.Get(backendName)
	if err != nil {
		return nil, err
	}
	return backend.List(ctx, prefix)
}

// DeleteFile 删除指定存储后端下的文件
func (s *StorageService) DeleteFile(ctx context.Context, backendName, path string) error {
	backend, err := s.storageMgr.Get(backendName)
	if err != nil {
		return err
	}
	return backend.Delete(ctx, path)
}

// GetDownloadURL 获取文件下载预签名 URL
func (s *StorageService) GetDownloadURL(ctx context.Context, backendName, path string) (string, error) {
	backend, err := s.storageMgr.Get(backendName)
	if err != nil {
		return "", err
	}
	return backend.GetURL(ctx, path, 15*time.Minute)
}

// InitMultipart 初始化分片上传（占位，待对接具体实现）
func (s *StorageService) InitMultipart(ctx context.Context, bucket, objectName string) (string, error) {
	return "dummy_upload_id", nil
}
