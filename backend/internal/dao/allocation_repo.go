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
