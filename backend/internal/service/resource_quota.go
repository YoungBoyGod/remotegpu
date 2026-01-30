package service

import (
	"errors"
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	if e.Available < 0 {
		return fmt.Sprintf("%s 配额不足: 需要 %d, 可用 0 (已超额使用 %d)", e.Resource, e.Requested, -e.Available)
	}
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

// GetQuotaInTx 在事务中获取资源配额（使用悲观锁 FOR UPDATE）
// 用于需要并发安全的场景，如环境创建时的配额检查
func (s *ResourceQuotaService) GetQuotaInTx(tx *gorm.DB, customerID uint, workspaceID *uint) (*entity.ResourceQuota, error) {
	if tx == nil {
		return nil, fmt.Errorf("事务不能为空")
	}

	var quota entity.ResourceQuota
	query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("customer_id = ?", customerID)

	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}

	if err := query.First(&quota).Error; err != nil {
		return nil, err
	}

	return &quota, nil
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

	// 获取已使用资源
	used, err := s.GetUsedResources(existing.CustomerID, existing.WorkspaceID)
	if err != nil {
		return fmt.Errorf("获取已使用资源失败: %w", err)
	}

	// 验证新配额是否小于已使用资源
	if quota.CPU < used.CPU {
		return fmt.Errorf("CPU配额不能小于已使用量: 已使用%d，新配额%d", used.CPU, quota.CPU)
	}
	if quota.Memory < used.Memory {
		return fmt.Errorf("内存配额不能小于已使用量: 已使用%dMB，新配额%dMB", used.Memory, quota.Memory)
	}
	if quota.GPU < used.GPU {
		return fmt.Errorf("GPU配额不能小于已使用量: 已使用%d，新配额%d", used.GPU, quota.GPU)
	}
	if quota.Storage < used.Storage {
		return fmt.Errorf("存储配额不能小于已使用量: 已使用%dGB，新配额%dGB", used.Storage, quota.Storage)
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
// 注意：此方法不保证并发安全。如需在环境创建等场景中使用，请使用 CheckQuotaInTx
func (s *ResourceQuotaService) CheckQuota(customerID uint, workspaceID *uint, request *ResourceRequest) (bool, error) {
	return s.checkQuota(nil, customerID, workspaceID, request)
}

// CheckQuotaInTx 在事务中检查资源配额是否足够（使用悲观锁，保证并发安全）
// 用于环境创建等需要原子性操作的场景
// 使用方式：
//   tx := db.Begin()
//   ok, err := service.CheckQuotaInTx(tx, customerID, workspaceID, request)
//   if ok {
//       // 创建环境...
//       tx.Commit()
//   } else {
//       tx.Rollback()
//   }
func (s *ResourceQuotaService) CheckQuotaInTx(tx *gorm.DB, customerID uint, workspaceID *uint, request *ResourceRequest) (bool, error) {
	if tx == nil {
		return false, fmt.Errorf("事务不能为空")
	}
	return s.checkQuota(tx, customerID, workspaceID, request)
}

// checkQuota 内部实现：检查资源配额是否足够
// tx 为 nil 时使用普通查询，不为 nil 时在事务中执行并使用悲观锁
func (s *ResourceQuotaService) checkQuota(tx *gorm.DB, customerID uint, workspaceID *uint, request *ResourceRequest) (bool, error) {
	// 1. 获取配额
	var quota *entity.ResourceQuota
	var err error

	if tx != nil {
		// 在事务中使用悲观锁获取配额
		quota, err = s.GetQuotaInTx(tx, customerID, workspaceID)
	} else {
		// 普通查询
		quota, err = s.GetQuota(customerID, workspaceID)
	}

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
//
// 资源统计范围说明：
// - workspaceID != nil: 统计指定工作空间内的所有运行中和创建中环境的资源使用
// - workspaceID == nil: 统计用户在所有工作空间的所有运行中和创建中环境的资源使用总和
//
// 注意：用户级别配额会统计该用户在所有工作空间中的资源使用，确保用户总资源不超限。
func (s *ResourceQuotaService) GetUsedResources(customerID uint, workspaceID *uint) (*UsedResources, error) {
	db := database.GetDB()

	// 查询所有运行中和创建中的环境（创建中的环境也占用资源）
	var environments []*entity.Environment
	query := db.Where("customer_id = ? AND status IN ?", customerID, []string{"running", "creating"})

	// 根据 workspaceID 过滤环境范围
	if workspaceID != nil {
		// 工作空间级别：只统计该工作空间的环境
		query = query.Where("workspace_id = ?", *workspaceID)
	}
	// 用户级别（workspaceID == nil）：统计该用户所有环境，不添加workspace_id过滤条件

	if err := query.Find(&environments).Error; err != nil {
		return nil, err
	}

	// 统计资源使用量
	used := &UsedResources{
		CPU:     0,
		Memory:  0,
		GPU:     0,
		Storage: 0,
	}

	for _, env := range environments {
		used.CPU += env.CPU
		used.Memory += env.Memory
		used.GPU += env.GPU
		if env.Storage != nil {
			used.Storage += *env.Storage
		}
	}

	return used, nil
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

	// 3. 计算可用配额（确保不为负数）
	availableCPU := quota.CPU - used.CPU
	if availableCPU < 0 {
		availableCPU = 0
	}
	availableMemory := quota.Memory - used.Memory
	if availableMemory < 0 {
		availableMemory = 0
	}
	availableGPU := quota.GPU - used.GPU
	if availableGPU < 0 {
		availableGPU = 0
	}
	availableStorage := quota.Storage - used.Storage
	if availableStorage < 0 {
		availableStorage = 0
	}

	available := &entity.ResourceQuota{
		CustomerID:  customerID,
		WorkspaceID: workspaceID,
		CPU:         availableCPU,
		Memory:      availableMemory,
		GPU:         availableGPU,
		Storage:     availableStorage,
	}

	return available, nil
}
