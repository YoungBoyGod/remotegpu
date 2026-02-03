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

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (d *CustomerDao) UpdateStatus(ctx context.Context, id uint, status string) error {
	return d.db.WithContext(ctx).Model(&entity.Customer{}).Where("id = ?", id).Update("status", status).Error
}
