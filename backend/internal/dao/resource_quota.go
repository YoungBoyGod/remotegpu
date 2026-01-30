package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

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

// GetByCustomerID 根据客户ID获取资源配额（用户级别配额，workspace_id为空）
func (d *ResourceQuotaDao) GetByCustomerID(customerID uint) (*entity.ResourceQuota, error) {
	var quota entity.ResourceQuota
	err := d.db.Where("customer_id = ? AND workspace_id IS NULL", customerID).First(&quota).Error
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetByWorkspaceID 根据工作空间ID获取资源配额
func (d *ResourceQuotaDao) GetByWorkspaceID(workspaceID uint) (*entity.ResourceQuota, error) {
	var quota entity.ResourceQuota
	err := d.db.Where("workspace_id = ?", workspaceID).First(&quota).Error
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetByCustomerAndWorkspace 根据客户ID和工作空间ID获取资源配额
func (d *ResourceQuotaDao) GetByCustomerAndWorkspace(customerID, workspaceID uint) (*entity.ResourceQuota, error) {
	var quota entity.ResourceQuota
	err := d.db.Where("customer_id = ? AND workspace_id = ?", customerID, workspaceID).First(&quota).Error
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
