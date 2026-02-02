package dao

import (
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// BaseDao 通用DAO基类，提供基础CRUD操作
// T 为实体类型，ID 为主键类型
type BaseDao[T any, ID comparable] struct {
	db *gorm.DB
}

// NewBaseDao 创建BaseDao实例
func NewBaseDao[T any, ID comparable]() *BaseDao[T, ID] {
	return &BaseDao[T, ID]{
		db: database.GetDB(),
	}
}

// NewBaseDaoWithDB 使用指定的数据库连接创建BaseDao实例
func NewBaseDaoWithDB[T any, ID comparable](db *gorm.DB) *BaseDao[T, ID] {
	return &BaseDao[T, ID]{
		db: db,
	}
}

// Create 创建实体
func (d *BaseDao[T, ID]) Create(entity *T) error {
	return d.db.Create(entity).Error
}

// GetByID 根据ID获取实体
func (d *BaseDao[T, ID]) GetByID(id ID) (*T, error) {
	var entity T
	err := d.db.Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 更新实体
func (d *BaseDao[T, ID]) Update(entity *T) error {
	return d.db.Save(entity).Error
}

// Delete 删除实体
func (d *BaseDao[T, ID]) Delete(id ID) error {
	var entity T
	return d.db.Delete(&entity, id).Error
}

// List 分页查询实体列表
func (d *BaseDao[T, ID]) List(page, pageSize int) ([]*T, int64, error) {
	var entities []*T
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := d.db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := d.db.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// ListAll 查询所有实体
func (d *BaseDao[T, ID]) ListAll() ([]*T, error) {
	var entities []*T
	if err := d.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count 统计实体数量
func (d *BaseDao[T, ID]) Count() (int64, error) {
	var count int64
	if err := d.db.Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists 检查实体是否存在
func (d *BaseDao[T, ID]) Exists(id ID) (bool, error) {
	var count int64
	err := d.db.Model(new(T)).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetDB 获取数据库连接（用于复杂查询）
func (d *BaseDao[T, ID]) GetDB() *gorm.DB {
	return d.db
}

// Transaction 执行事务
func (d *BaseDao[T, ID]) Transaction(fn func(tx *gorm.DB) error) error {
	return d.db.Transaction(fn)
}
