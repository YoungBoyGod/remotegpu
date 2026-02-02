package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// StorageSourceDao 存储源数据访问对象
type StorageSourceDao struct {
	db *gorm.DB
}

// NewStorageSourceDao 创建存储源DAO
func NewStorageSourceDao() *StorageSourceDao {
	return &StorageSourceDao{}
}

// NewStorageSourceDaoWithDB 创建存储源DAO(带数据库连接)
func NewStorageSourceDaoWithDB(db *gorm.DB) *StorageSourceDao {
	return &StorageSourceDao{db: db}
}

// Create 创建存储源
func (d *StorageSourceDao) Create(source *entity.StorageSource) error {
	return d.getDB().Create(source).Error
}

// GetByID 根据ID获取存储源
func (d *StorageSourceDao) GetByID(id uint) (*entity.StorageSource, error) {
	var source entity.StorageSource
	err := d.getDB().Where("id = ?", id).First(&source).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

// GetByName 根据名称获取存储源
func (d *StorageSourceDao) GetByName(name string) (*entity.StorageSource, error) {
	var source entity.StorageSource
	err := d.getDB().Where("name = ?", name).First(&source).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

// List 获取存储源列表
func (d *StorageSourceDao) List(status string) ([]*entity.StorageSource, error) {
	var sources []*entity.StorageSource
	query := d.getDB()

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("priority ASC, created_at DESC").Find(&sources).Error
	return sources, err
}

// ListByType 根据类型获取存储源列表
func (d *StorageSourceDao) ListByType(sourceType string, status string) ([]*entity.StorageSource, error) {
	var sources []*entity.StorageSource
	query := d.getDB().Where("type = ?", sourceType)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("priority ASC, created_at DESC").Find(&sources).Error
	return sources, err
}

// ListByPublic 根据是否公有云获取存储源列表
func (d *StorageSourceDao) ListByPublic(isPublic bool, status string) ([]*entity.StorageSource, error) {
	var sources []*entity.StorageSource
	query := d.getDB().Where("is_public = ?", isPublic)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("priority ASC, created_at DESC").Find(&sources).Error
	return sources, err
}

// Update 更新存储源
func (d *StorageSourceDao) Update(source *entity.StorageSource) error {
	return d.getDB().Save(source).Error
}

// UpdateStatus 更新存储源状态
func (d *StorageSourceDao) UpdateStatus(id uint, status string) error {
	return d.getDB().Model(&entity.StorageSource{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除存储源
func (d *StorageSourceDao) Delete(id uint) error {
	return d.getDB().Delete(&entity.StorageSource{}, id).Error
}

// getDB 获取数据库连接
func (d *StorageSourceDao) getDB() *gorm.DB {
	if d.db != nil {
		return d.db
	}
	// 如果没有注入DB,使用全局DB
	return database.GetDB()
}
