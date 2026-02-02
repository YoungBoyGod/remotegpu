package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// UsedResources 已使用的资源统计
type UsedResources struct {
	CPU     int   `json:"cpu"`
	Memory  int64 `json:"memory"`
	GPU     int   `json:"gpu"`
	Storage int64 `json:"storage"`
}

// AvailableQuota 可用配额
type AvailableQuota struct {
	CPU     int   `json:"cpu"`
	Memory  int64 `json:"memory"`
	GPU     int   `json:"gpu"`
	Storage int64 `json:"storage"`
}

// ResourceQuotaDao 资源配额数据访问对象
type ResourceQuotaDao struct {
	db *gorm.DB
}

// NewResourceQuotaDao 创建资源配额 DAO
func NewResourceQuotaDao() *ResourceQuotaDao {
	return &ResourceQuotaDao{
		db: database.GetDB(),
	}
}

// Create 创建资源配额
func (d *ResourceQuotaDao) Create(quota *entity.ResourceQuota) error {
	return d.db.Create(quota).Error
}

// GetByID 根据ID获取资源配额
func (d *ResourceQuotaDao) GetByID(id uint) (*entity.ResourceQuota, error) {
	var quota entity.ResourceQuota
	err := d.db.Where("id = ?", id).First(&quota).Error
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetByUserID 根据用户ID获取资源配额（用户级别配额，workspace_id为空）
func (d *ResourceQuotaDao) GetByUserID(userID uint) (*entity.ResourceQuota, error) {
	var quota entity.ResourceQuota
	err := d.db.Where("user_id = ? AND workspace_id IS NULL", userID).First(&quota).Error
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// Update 更新资源配额
func (d *ResourceQuotaDao) Update(quota *entity.ResourceQuota) error {
	return d.db.Save(quota).Error
}

// Delete 删除资源配额
func (d *ResourceQuotaDao) Delete(id uint) error {
	return d.db.Delete(&entity.ResourceQuota{}, id).Error
}

// List 获取资源配额列表（分页）
func (d *ResourceQuotaDao) List(page, pageSize int) ([]*entity.ResourceQuota, int64, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小，防止内存溢出
	}

	var quotas []*entity.ResourceQuota
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&entity.ResourceQuota{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Offset(offset).Limit(pageSize).Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	return quotas, total, nil
}

// GetByQuotaLevel 根据配额级别获取资源配额列表
func (d *ResourceQuotaDao) GetByQuotaLevel(level string) ([]*entity.ResourceQuota, error) {
	var quotas []*entity.ResourceQuota
	err := d.db.Where("quota_level = ?", level).Find(&quotas).Error
	if err != nil {
		return nil, err
	}
	return quotas, nil
}

// GetUsedResources 统计用户所有运行中环境的资源使用情况
func (d *ResourceQuotaDao) GetUsedResources(userID uint) (*UsedResources, error) {
	var result UsedResources

	// 查询用户所有运行中的环境，统计资源使用
	err := d.db.Model(&entity.Environment{}).
		Select("COALESCE(SUM(cpu), 0) as cpu, COALESCE(SUM(memory), 0) as memory, COALESCE(SUM(gpu), 0) as gpu, COALESCE(SUM(storage), 0) as storage").
		Where("user_id = ? AND status = ?", userID, "running").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAvailableQuota 计算用户的可用配额
func (d *ResourceQuotaDao) GetAvailableQuota(userID uint) (*AvailableQuota, error) {
	// 获取用户配额
	quota, err := d.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 获取已使用的资源
	used, err := d.GetUsedResources(userID)
	if err != nil {
		return nil, err
	}

	// 计算可用配额 = 总配额 - 已用配额
	available := &AvailableQuota{
		CPU:     quota.CPU - used.CPU,
		Memory:  quota.Memory - used.Memory,
		GPU:     quota.GPU - used.GPU,
		Storage: quota.Storage - used.Storage,
	}

	return available, nil
}

// CheckQuota 检查用户请求的资源是否超过配额限制
func (d *ResourceQuotaDao) CheckQuota(userID uint, cpu int, memory int64, gpu int, storage int64) (bool, error) {
	// 获取用户配额
	quota, err := d.GetByUserID(userID)
	if err != nil {
		return false, err
	}

	// 获取已使用的资源
	used, err := d.GetUsedResources(userID)
	if err != nil {
		return false, err
	}

	// 检查CPU配额
	if used.CPU+cpu > quota.CPU {
		return false, nil
	}

	// 检查内存配额
	if used.Memory+memory > quota.Memory {
		return false, nil
	}

	// 检查GPU配额
	if used.GPU+gpu > quota.GPU {
		return false, nil
	}

	// 检查存储配额
	if used.Storage+storage > quota.Storage {
		return false, nil
	}

	return true, nil
}
