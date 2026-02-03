package service

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
	// Wrapper around minio or local storage logic
	// return s.storageMgr.InitMultipartUpload(ctx, bucket, objectName)
	return "dummy_upload_id", nil
}
