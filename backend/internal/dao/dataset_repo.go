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

// UpdateFields 更新数据集指定字段
func (d *DatasetDao) UpdateFields(ctx context.Context, id uint, fields map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&entity.Dataset{}).Where("id = ?", id).Updates(fields).Error
}

// FindByID 根据ID查询数据集
// @author Claude
// @description 根据数据集ID查询详情，用于权限校验
// @modified 2026-02-04
func (d *DatasetDao) FindByID(ctx context.Context, id uint) (*entity.Dataset, error) {
	var dataset entity.Dataset
	if err := d.db.WithContext(ctx).First(&dataset, id).Error; err != nil {
		return nil, err
	}
	return &dataset, nil
}