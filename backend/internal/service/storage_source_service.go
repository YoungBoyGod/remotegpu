package service

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/adapter"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

// StorageSourceService 存储源管理服务
type StorageSourceService struct {
	dao *dao.StorageSourceDao
}

// NewStorageSourceService 创建存储源服务
func NewStorageSourceService() *StorageSourceService {
	return &StorageSourceService{
		dao: dao.NewStorageSourceDao(),
	}
}

// CreateStorageSource 创建存储源
func (s *StorageSourceService) CreateStorageSource(source *entity.StorageSource) error {
	// 验证存储源配置
	if err := s.validateStorageSource(source); err != nil {
		return fmt.Errorf("存储源配置验证失败: %w", err)
	}

	// 测试连接
	if err := s.testConnection(source); err != nil {
		return fmt.Errorf("存储源连接测试失败: %w", err)
	}

	// 创建存储源
	source.Status = "active"
	now := time.Now()
	source.HealthCheck = &now

	return s.dao.Create(source)
}

// GetStorageSource 获取存储源
func (s *StorageSourceService) GetStorageSource(id uint) (*entity.StorageSource, error) {
	return s.dao.GetByID(id)
}

// ListStorageSources 获取存储源列表
func (s *StorageSourceService) ListStorageSources(status string) ([]*entity.StorageSource, error) {
	return s.dao.List(status)
}

// UpdateStorageSource 更新存储源
func (s *StorageSourceService) UpdateStorageSource(source *entity.StorageSource) error {
	// 验证存储源配置
	if err := s.validateStorageSource(source); err != nil {
		return fmt.Errorf("存储源配置验证失败: %w", err)
	}

	return s.dao.Update(source)
}

// DeleteStorageSource 删除存储源
func (s *StorageSourceService) DeleteStorageSource(id uint) error {
	// TODO: 检查是否有数据集副本使用此存储源
	return s.dao.Delete(id)
}

// HealthCheckAll 对所有存储源进行健康检查
func (s *StorageSourceService) HealthCheckAll() error {
	sources, err := s.dao.List("active")
	if err != nil {
		return err
	}

	for _, source := range sources {
		if err := s.healthCheck(source); err != nil {
			// 更新状态为error
			_ = s.dao.UpdateStatus(source.ID, "error")
		} else {
			// 更新健康检查时间
			now := time.Now()
			source.HealthCheck = &now
			_ = s.dao.Update(source)
		}
	}

	return nil
}

// validateStorageSource 验证存储源配置
func (s *StorageSourceService) validateStorageSource(source *entity.StorageSource) error {
	if source.Name == "" {
		return fmt.Errorf("存储源名称不能为空")
	}
	if source.Type == "" {
		return fmt.Errorf("存储源类型不能为空")
	}
	if source.Endpoint == "" {
		return fmt.Errorf("存储源端点不能为空")
	}
	if source.Bucket == "" {
		return fmt.Errorf("存储桶名称不能为空")
	}
	return nil
}

// testConnection 测试存储源连接
func (s *StorageSourceService) testConnection(source *entity.StorageSource) error {
	adapter, err := s.createAdapter(source)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return adapter.HealthCheck(ctx)
}

// healthCheck 健康检查
func (s *StorageSourceService) healthCheck(source *entity.StorageSource) error {
	adapter, err := s.createAdapter(source)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return adapter.HealthCheck(ctx)
}

// createAdapter 创建存储适配器
func (s *StorageSourceService) createAdapter(source *entity.StorageSource) (adapter.StorageAdapter, error) {
	switch source.Type {
	case "minio":
		return adapter.NewMinIOAdapter(source.Endpoint, source.AccessKey, source.SecretKey, source.Bucket, false)
	case "oss":
		return adapter.NewOSSAdapter(source.Endpoint, source.AccessKey, source.SecretKey, source.Bucket)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", source.Type)
	}
}
