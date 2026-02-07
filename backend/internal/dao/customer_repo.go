package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type CustomerDao struct {
	*BaseDao[entity.Customer]
}

func NewCustomerDao(db *gorm.DB) *CustomerDao {
	return &CustomerDao{
		BaseDao: NewBaseDao[entity.Customer](db),
	}
}

func (d *CustomerDao) FindByUsername(ctx context.Context, username string) (*entity.Customer, error) {
	var customer entity.Customer
	if err := d.db.WithContext(ctx).Where("username = ?", username).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (d *CustomerDao) List(ctx context.Context, page, pageSize int) ([]entity.Customer, int64, error) {
	var customers []entity.Customer
	var total int64

	db := d.db.WithContext(ctx).Model(&entity.Customer{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (d *CustomerDao) UpdateStatus(ctx context.Context, id uint, status string) error {
	return d.db.WithContext(ctx).Model(&entity.Customer{}).Where("id = ?", id).Update("status", status).Error
}

// CountActive 统计活跃客户数量
// @description 统计状态为 active 的客户数量
// @modified 2026-02-04
func (d *CustomerDao) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Customer{}).
		Where("status = ?", "active").Count(&count).Error
	return count, err
}

// UpdateFields 更新客户指定字段
func (d *CustomerDao) UpdateFields(ctx context.Context, id uint, fields map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&entity.Customer{}).Where("id = ?", id).Updates(fields).Error
}

// Count 统计客户总数
// @modified 2026-02-04
func (d *CustomerDao) Count(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Customer{}).Count(&count).Error
	return count, err
}

// UpdateQuota 更新客户配额
func (d *CustomerDao) UpdateQuota(ctx context.Context, id uint, quotaGPU int, quotaStorage int64) error {
	return d.db.WithContext(ctx).Model(&entity.Customer{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"quota_gpu":     quotaGPU,
			"quota_storage": quotaStorage,
		}).Error
}
