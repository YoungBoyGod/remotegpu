package document

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/storage"
	"gorm.io/gorm"
)

type DocumentService struct {
	documentDao *dao.DocumentDao
	storageMgr  *storage.Manager
	db          *gorm.DB
}

func NewDocumentService(db *gorm.DB, mgr *storage.Manager) *DocumentService {
	return &DocumentService{
		documentDao: dao.NewDocumentDao(db),
		storageMgr:  mgr,
		db:          db,
	}
}

// ListDocuments 获取文档列表
func (s *DocumentService) ListDocuments(ctx context.Context, page, pageSize int, category, keyword string) ([]entity.Document, int64, error) {
	return s.documentDao.List(ctx, page, pageSize, category, keyword)
}

// GetDocument 获取文档详情
func (s *DocumentService) GetDocument(ctx context.Context, id uint) (*entity.Document, error) {
	return s.documentDao.FindByID(ctx, id)
}

// CreateDocument 创建文档记录
func (s *DocumentService) CreateDocument(ctx context.Context, doc *entity.Document) error {
	return s.documentDao.Create(ctx, doc)
}

// UpdateDocument 更新文档信息
func (s *DocumentService) UpdateDocument(ctx context.Context, id uint, fields map[string]any) error {
	// 先确认文档存在
	if _, err := s.documentDao.FindByID(ctx, id); err != nil {
		return fmt.Errorf("文档不存在: %w", err)
	}
	return s.documentDao.Update(ctx, id, fields)
}

// DeleteDocument 删除文档（同时删除存储文件）
func (s *DocumentService) DeleteDocument(ctx context.Context, id uint) error {
	doc, err := s.documentDao.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除存储文件
	if s.storageMgr != nil && doc.FilePath != "" {
		backend, err := s.storageMgr.Get(doc.StorageBackend)
		if err == nil {
			_ = backend.Delete(ctx, doc.FilePath)
		}
	}

	return s.documentDao.Delete(ctx, id)
}

// ListCategories 获取所有分类
func (s *DocumentService) ListCategories(ctx context.Context) ([]string, error) {
	return s.documentDao.ListCategories(ctx)
}

// GetDownloadURL 获取文档下载 URL
func (s *DocumentService) GetDownloadURL(ctx context.Context, id uint) (string, error) {
	doc, err := s.documentDao.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	if s.storageMgr == nil {
		return "", fmt.Errorf("storage not available")
	}

	backend, err := s.storageMgr.Get(doc.StorageBackend)
	if err != nil {
		return "", err
	}

	return backend.GetURL(ctx, doc.FilePath, 30*time.Minute)
}

// BuildStoragePath 构建存储路径
func BuildStoragePath(category, fileName string) string {
	return filepath.Join("documents", category, fileName)
}
