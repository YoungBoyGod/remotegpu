package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type DatasetDao struct {
	*BaseDao[entity.Dataset]
}

func NewDatasetDao(db *gorm.DB) *DatasetDao {
	return &DatasetDao{
		BaseDao: NewBaseDao[entity.Dataset](db),
	}
}

func (d *DatasetDao) ListByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Dataset, int64, error) {
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