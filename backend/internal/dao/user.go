package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// UserDaoInterface 定义CustomerDao的接口
type UserDaoInterface interface {
	Create(customer *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(customer *entity.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*entity.User, int64, error)
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{
		db: database.GetDB(),
	}
}

// Create 创建客户
func (d *UserDao) Create(customer *entity.User) error {
	return d.db.Create(customer).Error
}

// GetByID 根据ID获取客户
func (d *UserDao) GetByID(id uint) (*entity.User, error) {
	var customer entity.User
	err := d.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetByUsername 根据用户名获取客户
func (d *UserDao) GetByUsername(username string) (*entity.User, error) {
	var customer entity.User
	err := d.db.Where("username = ?", username).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetByEmail 根据邮箱获取客户
func (d *UserDao) GetByEmail(email string) (*entity.User, error) {
	var customer entity.User
	err := d.db.Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Update 更新客户
func (d *UserDao) Update(customer *entity.User) error {
	return d.db.Save(customer).Error
}

// Delete 删除客户
func (d *UserDao) Delete(id uint) error {
	return d.db.Delete(&entity.User{}, id).Error
}

// List 获取客户列表
func (d *UserDao) List(page, pageSize int) ([]*entity.User, int64, error) {
	var customers []*entity.User
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Offset(offset).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
