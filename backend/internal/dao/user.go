package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

type CustomerDao struct {
	db *gorm.DB
}

func NewCustomerDao() *CustomerDao {
	return &CustomerDao{
		db: database.GetDB(),
	}
}

// Create 创建客户
func (d *CustomerDao) Create(customer *entity.Customer) error {
	return d.db.Create(customer).Error
}

// GetByID 根据ID获取客户
func (d *CustomerDao) GetByID(id uint) (*entity.Customer, error) {
	var customer entity.Customer
	err := d.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetByUsername 根据用户名获取客户
func (d *CustomerDao) GetByUsername(username string) (*entity.Customer, error) {
	var customer entity.Customer
	err := d.db.Where("username = ?", username).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetByEmail 根据邮箱获取客户
func (d *CustomerDao) GetByEmail(email string) (*entity.Customer, error) {
	var customer entity.Customer
	err := d.db.Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Update 更新客户
func (d *CustomerDao) Update(customer *entity.Customer) error {
	return d.db.Save(customer).Error
}

// Delete 删除客户
func (d *CustomerDao) Delete(id uint) error {
	return d.db.Delete(&entity.Customer{}, id).Error
}

// List 获取客户列表
func (d *CustomerDao) List(page, pageSize int) ([]*entity.Customer, int64, error) {
	var customers []*entity.Customer
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&entity.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Offset(offset).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
