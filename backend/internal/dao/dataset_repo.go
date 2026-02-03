package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type DatasetRepo struct {
	*BaseDao[entity.Dataset]
}

func NewDatasetRepo(db *gorm.DB) *DatasetRepo {
	return &DatasetRepo{
		BaseDao: NewBaseDao[entity.Dataset](db),
	}
}

func (d *DatasetRepo) ListByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Dataset, int64, error) {
	var datasets []entity.Dataset
	var total int64

	db := d.db.WithContext(ctx).Model(&entity.Dataset{}).Where("customer_id = ?", customerID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&datasets).Error; err != nil {
		return nil, 0, err
	}

	return datasets, total, nil
}
