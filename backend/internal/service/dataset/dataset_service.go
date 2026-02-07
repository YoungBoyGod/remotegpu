package dataset

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

var (
	ErrDatasetNotReady    = errors.New("数据集未就绪，无法挂载")
	ErrMountAlreadyExists = errors.New("该数据集已挂载到此机器")
	ErrMountNotFound      = errors.New("挂载记录不存在")
	ErrInvalidMountPath   = errors.New("挂载路径不合法")
)

type DatasetService struct {
	datasetDao      *dao.DatasetDao
	datasetMountDao *dao.DatasetMountDao
}

func NewDatasetService(db *gorm.DB) *DatasetService {
	return &DatasetService{
		datasetDao:      dao.NewDatasetDao(db),
		datasetMountDao: dao.NewDatasetMountDao(db),
	}
}

func (s *DatasetService) ListDatasets(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Dataset, int64, error) {
	return s.datasetDao.ListByCustomerID(ctx, customerID, page, pageSize)
}

// GetDataset 根据ID获取数据集
// @author Claude
// @description 获取数据集详情，用于权限校验
// @modified 2026-02-04
func (s *DatasetService) GetDataset(ctx context.Context, id uint) (*entity.Dataset, error) {
	return s.datasetDao.FindByID(ctx, id)
}

// CompleteUpload 完成分片上传，更新数据集状态为 ready
func (s *DatasetService) CompleteUpload(ctx context.Context, id uint, name string, size int64) error {
	fields := map[string]interface{}{
		"status": "ready",
	}
	if name != "" {
		fields["name"] = name
	}
	if size > 0 {
		fields["total_size"] = size
	}
	return s.datasetDao.UpdateFields(ctx, id, fields)
}
// @author Claude
// @description 验证数据集是否属于指定用户
// @modified 2026-02-04
func (s *DatasetService) ValidateOwnership(ctx context.Context, datasetID uint, customerID uint) error {
	dataset, err := s.datasetDao.FindByID(ctx, datasetID)
	if err != nil {
		return err
	}
	if dataset.CustomerID != customerID {
		return entity.ErrUnauthorized
	}
	return nil
}

// MountDataset 创建数据集挂载记录
func (s *DatasetService) MountDataset(ctx context.Context, datasetID uint, hostID, mountPath string, readOnly bool) (*entity.DatasetMount, error) {
	// 校验挂载路径合法性
	if err := validateMountPath(mountPath); err != nil {
		return nil, err
	}

	// 校验数据集状态
	dataset, err := s.datasetDao.FindByID(ctx, datasetID)
	if err != nil {
		return nil, fmt.Errorf("数据集不存在: %w", err)
	}
	if dataset.Status != "ready" {
		return nil, ErrDatasetNotReady
	}

	// 检查是否已存在活跃挂载
	existing, err := s.datasetMountDao.FindActiveMount(ctx, datasetID, hostID)
	if err == nil && existing != nil {
		return nil, ErrMountAlreadyExists
	}

	mount := &entity.DatasetMount{
		DatasetID: datasetID,
		HostID:    hostID,
		MountPath: filepath.Clean(mountPath),
		ReadOnly:  readOnly,
		Status:    "mounting",
	}
	if err := s.datasetMountDao.Create(ctx, mount); err != nil {
		return nil, err
	}
	return mount, nil
}

// UnmountDataset 卸载数据集
func (s *DatasetService) UnmountDataset(ctx context.Context, mountID uint) error {
	mount, err := s.datasetMountDao.FindByID(ctx, mountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMountNotFound
		}
		return err
	}
	if mount.Status == "unmounted" {
		return nil
	}
	return s.datasetMountDao.UpdateStatus(ctx, mountID, "unmounted", "")
}

// UpdateMountStatus 更新挂载状态（供 Agent 回调使用）
func (s *DatasetService) UpdateMountStatus(ctx context.Context, mountID uint, status, errMsg string) error {
	_, err := s.datasetMountDao.FindByID(ctx, mountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMountNotFound
		}
		return err
	}
	return s.datasetMountDao.UpdateStatus(ctx, mountID, status, errMsg)
}

// ListMountsByDataset 获取数据集的所有挂载
func (s *DatasetService) ListMountsByDataset(ctx context.Context, datasetID uint) ([]entity.DatasetMount, error) {
	return s.datasetMountDao.ListByDatasetID(ctx, datasetID)
}

// ListMountsByHost 获取机器上的所有挂载
func (s *DatasetService) ListMountsByHost(ctx context.Context, hostID string) ([]entity.DatasetMount, error) {
	return s.datasetMountDao.ListByHostID(ctx, hostID)
}

// CleanupMountsByHost 清理机器上的所有挂载（回收机器时调用）
func (s *DatasetService) CleanupMountsByHost(ctx context.Context, hostID string) error {
	mounts, err := s.datasetMountDao.ListByHostID(ctx, hostID)
	if err != nil {
		return err
	}
	for _, m := range mounts {
		if err := s.datasetMountDao.UpdateStatus(ctx, m.ID, "unmounted", "机器回收清理"); err != nil {
			return fmt.Errorf("清理挂载 %d 失败: %w", m.ID, err)
		}
	}
	return nil
}

// validateMountPath 校验挂载路径合法性
func validateMountPath(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return ErrInvalidMountPath
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("%w: 必须为绝对路径", ErrInvalidMountPath)
	}
	// 禁止挂载到系统关键目录
	forbidden := []string{"/", "/bin", "/sbin", "/usr", "/etc", "/proc", "/sys", "/dev", "/boot", "/root"}
	cleaned := filepath.Clean(path)
	for _, f := range forbidden {
		if cleaned == f {
			return fmt.Errorf("%w: 不允许挂载到系统目录 %s", ErrInvalidMountPath, f)
		}
	}
	// 禁止路径穿越
	if strings.Contains(path, "..") {
		return fmt.Errorf("%w: 路径不允许包含 ..", ErrInvalidMountPath)
	}
	return nil
}
