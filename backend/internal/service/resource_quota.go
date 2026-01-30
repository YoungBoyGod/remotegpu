package service

import (
	"errors"
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// ResourceQuotaService 资源配额服务
type ResourceQuotaService struct {
	quotaDao *dao.ResourceQuotaDao
}

// NewResourceQuotaService 创建资源配额服务实例
func NewResourceQuotaService() *ResourceQuotaService {
	return &ResourceQuotaService{
		quotaDao: dao.NewResourceQuotaDao(),
	}
}

// ResourceRequest 资源请求结构
type ResourceRequest struct {
	CPU     int   `json:"cpu"`
	Memory  int64 `json:"memory"`
	GPU     int   `json:"gpu"`
	Storage int64 `json:"storage"`
}

// UsedResources 已使用资源结构
type UsedResources struct {
	CPU     int   `json:"cpu"`
	Memory  int64 `json:"memory"`
	GPU     int   `json:"gpu"`
	Storage int64 `json:"storage"`
}

// QuotaExceededError 配额超限错误
type QuotaExceededError struct {
	Resource  string
	Requested int64
	Available int64
}

func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("%s 配额不足: 需要 %d, 可用 %d", e.Resource, e.Requested, e.Available)
}

// SetQuota 设置资源配额
func (s *ResourceQuotaService) SetQuota(quota *entity.ResourceQuota) error {
	// 验证配额值
	if quota.CPU < 0 || quota.Memory < 0 || quota.GPU < 0 || quota.Storage < 0 {
		return fmt.Errorf("配额值不能为负数")
	}

	// 检查是否已存在配额
	var existing *entity.ResourceQuota
	var err error

	if quota.WorkspaceID != nil {
		existing, err = s.quotaDao.GetByCustomerAndWorkspace(quota.CustomerID, *quota.WorkspaceID)
	} else {
		existing, err = s.quotaDao.GetByCustomerID(quota.CustomerID)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 如果已存在，更新；否则创建
	if existing != nil {
		quota.ID = existing.ID
		return s.quotaDao.Update(quota)
	}

	return s.quotaDao.Create(quota)
}

// GetQuota 获取资源配额
func (s *ResourceQuotaService) GetQuota(customerID uint, workspaceID *uint) (*entity.ResourceQuota, error) {
	if workspaceID != nil {
		return s.quotaDao.GetByCustomerAndWorkspace(customerID, *workspaceID)
	}
	return s.quotaDao.GetByCustomerID(customerID)
}

// UpdateQuota 更新资源配额
func (s *ResourceQuotaService) UpdateQuota(quota *entity.ResourceQuota) error {
	// 验证配额值
	if quota.CPU < 0 || quota.Memory < 0 || quota.GPU < 0 || quota.Storage < 0 {
		return fmt.Errorf("配额值不能为负数")
	}

	// 检查配额是否存在
	existing, err := s.quotaDao.GetByID(quota.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("配额不存在")
		}
		return err
	}

	// 更新字段
	existing.CPU = quota.CPU
	existing.Memory = quota.Memory
	existing.GPU = quota.GPU
	existing.Storage = quota.Storage

	return s.quotaDao.Update(existing)
}

// DeleteQuota 删除资源配额
func (s *ResourceQuotaService) DeleteQuota(id uint) error {
	// 检查配额是否存在
	_, err := s.quotaDao.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("配额不存在")
		}
		return err
	}

	return s.quotaDao.Delete(id)
}

// CheckQuota 检查资源配额是否足够
func (s *ResourceQuotaService) CheckQuota(customerID uint, workspaceID *uint, request *ResourceRequest) (bool, error) {
	// 1. 获取配额
	quota, err := s.GetQuota(customerID, workspaceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("未设置资源配额")
		}
		return false, err
	}

	// 2. 获取已使用资源
	used, err := s.GetUsedResources(customerID, workspaceID)
	if err != nil {
		return false, err
	}

	// 3. 计算可用配额
	availableCPU := quota.CPU - used.CPU
	availableMemory := quota.Memory - used.Memory
	availableGPU := quota.GPU - used.GPU
	availableStorage := quota.Storage - used.Storage

	// 4. 检查是否足够
	if request.CPU > availableCPU {
		return false, &QuotaExceededError{
			Resource:  "CPU",
			Requested: int64(request.CPU),
			Available: int64(availableCPU),
		}
	}
	if request.Memory > availableMemory {
		return false, &QuotaExceededError{
			Resource:  "Memory",
			Requested: request.Memory,
			Available: availableMemory,
		}
	}
	if request.GPU > availableGPU {
		return false, &QuotaExceededError{
			Resource:  "GPU",
			Requested: int64(request.GPU),
			Available: int64(availableGPU),
		}
	}
	if request.Storage > availableStorage {
		return false, &QuotaExceededError{
			Resource:  "Storage",
			Requested: request.Storage,
			Available: availableStorage,
		}
	}

	return true, nil
}

// GetUsedResources 获取已使用的资源
// TODO: 需要等待 Environment 实体实现后才能完成此功能（A5任务）
func (s *ResourceQuotaService) GetUsedResources(customerID uint, workspaceID *uint) (*UsedResources, error) {
	// 暂时返回空的已使用资源，等待 Environment 实体实现
	// 实现逻辑：
	// 1. 查询所有运行中的环境（status = 'running'）
	// 2. 根据 customerID 和 workspaceID 过滤
	// 3. 统计 CPU、Memory、GPU、Storage 使用量

	return &UsedResources{
		CPU:     0,
		Memory:  0,
		GPU:     0,
		Storage: 0,
	}, nil
}

// GetAvailableQuota 获取可用配额
func (s *ResourceQuotaService) GetAvailableQuota(customerID uint, workspaceID *uint) (*entity.ResourceQuota, error) {
	// 1. 获取总配额
	quota, err := s.GetQuota(customerID, workspaceID)
	if err != nil {
		return nil, err
	}

	// 2. 获取已使用资源
	used, err := s.GetUsedResources(customerID, workspaceID)
	if err != nil {
		return nil, err
	}

	// 3. 计算可用配额
	available := &entity.ResourceQuota{
		CustomerID:  customerID,
		WorkspaceID: workspaceID,
		CPU:         quota.CPU - used.CPU,
		Memory:      quota.Memory - used.Memory,
		GPU:         quota.GPU - used.GPU,
		Storage:     quota.Storage - used.Storage,
	}

	return available, nil
}
