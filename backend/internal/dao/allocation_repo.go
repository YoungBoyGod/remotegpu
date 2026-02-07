package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type AllocationDao struct {
	db *gorm.DB
}

func NewAllocationDao(db *gorm.DB) *AllocationDao {
	return &AllocationDao{db: db}
}

func (d *AllocationDao) Create(ctx context.Context, allocation *entity.Allocation) error {
	return d.db.WithContext(ctx).Create(allocation).Error
}

func (d *AllocationDao) FindByID(ctx context.Context, id string) (*entity.Allocation, error) {
	var allocation entity.Allocation
	if err := d.db.WithContext(ctx).Preload("Customer").Preload("Host").First(&allocation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &allocation, nil
}

func (d *AllocationDao) UpdateStatus(ctx context.Context, id string, status string) error {
	return d.db.WithContext(ctx).Model(&entity.Allocation{}).Where("id = ?", id).Update("status", status).Error
}

func (d *AllocationDao) FindRecent(ctx context.Context, limit int) ([]entity.Allocation, error) {
	var allocations []entity.Allocation
	err := d.db.WithContext(ctx).
		Preload("Customer").
		Preload("Host").
		Order("created_at desc").
		Limit(limit).
		Find(&allocations).Error
	return allocations, err
}

func (d *AllocationDao) FindActiveByHostID(ctx context.Context, hostID string) (*entity.Allocation, error) {
	var allocation entity.Allocation
	err := d.db.WithContext(ctx).
		Where("host_id = ? AND status = ?", hostID, "active").
		First(&allocation).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

func (d *AllocationDao) FindActiveByHostAndCustomer(ctx context.Context, hostID string, customerID uint) (*entity.Allocation, error) {
	var allocation entity.Allocation
	err := d.db.WithContext(ctx).
		Where("host_id = ? AND customer_id = ? AND status = ?", hostID, customerID, "active").
		First(&allocation).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

// FindAllActiveByCustomerID 查询客户所有活跃分配（带 Host 预加载，不分页）
func (d *AllocationDao) FindAllActiveByCustomerID(ctx context.Context, customerID uint) ([]entity.Allocation, error) {
	var allocations []entity.Allocation
	err := d.db.WithContext(ctx).
		Where("customer_id = ? AND status = ?", customerID, "active").
		Preload("Host").
		Order("created_at desc").
		Find(&allocations).Error
	return allocations, err
}

// FindActiveByCustomerID 根据客户ID查询活跃分配记录
// @author Claude
// @description 查询指定客户的所有活跃分配，用于客户端机器列表过滤
// @param customerID 客户ID
// @param page 页码
// @param pageSize 每页数量
// @return 分配列表、总数、错误
// @modified 2026-02-04
func (d *AllocationDao) FindActiveByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Allocation, int64, error) {
	var allocations []entity.Allocation
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.Allocation{}).
		Where("customer_id = ? AND status = ?", customerID, "active")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Host").
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&allocations).Error

	return allocations, total, err
}
