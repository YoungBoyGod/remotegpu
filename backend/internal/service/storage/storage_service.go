package storage

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/pkg/storage"
)

type StorageService struct {
	storageMgr *storage.Manager
}

func NewStorageService(mgr *storage.Manager) *StorageService {
	return &StorageService{
		storageMgr: mgr,
	}
}

func (s *StorageService) InitMultipart(ctx context.Context, bucket, objectName string) (string, error) {
	// 封装 minio 或本地存储逻辑
	// return s.storageMgr.InitMultipartUpload(ctx, bucket, objectName)
	return "dummy_upload_id", nil
}
