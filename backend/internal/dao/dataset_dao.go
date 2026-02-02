package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// DatasetDao 数据集数据访问对象
type DatasetDao struct {
	db *gorm.DB
}

// NewDatasetDao 创建数据集DAO
func NewDatasetDao() *DatasetDao {
	return &DatasetDao{}
}

// NewDatasetDaoWithDB 创建数据集DAO(带数据库连接)
func NewDatasetDaoWithDB(db *gorm.DB) *DatasetDao {
	return &DatasetDao{db: db}
}

// Create 创建数据集
func (d *DatasetDao) Create(dataset *entity.Dataset) error {
	return d.getDB().Create(dataset).Error
}

// GetByID 根据ID获取数据集
func (d *DatasetDao) GetByID(id uint) (*entity.Dataset, error) {
	var dataset entity.Dataset
	err := d.getDB().Where("id = ?", id).First(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

// GetByIDWithReplicas 根据ID获取数据集(包含副本信息)
func (d *DatasetDao) GetByIDWithReplicas(id uint) (*entity.Dataset, error) {
	var dataset entity.Dataset
	err := d.getDB().Preload("Replicas").Preload("Replicas.StorageSource").
		Where("id = ?", id).First(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

// ListByUser 获取用户的数据集列表
func (d *DatasetDao) ListByUser(userID uint, status string) ([]*entity.Dataset, error) {
	var datasets []*entity.Dataset
	query := d.getDB().Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&datasets).Error
	return datasets, err
}

// ListByWorkspace 获取工作空间的数据集列表
func (d *DatasetDao) ListByWorkspace(workspaceID uint, status string) ([]*entity.Dataset, error) {
	var datasets []*entity.Dataset
	query := d.getDB().Where("workspace_id = ?", workspaceID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&datasets).Error
	return datasets, err
}

// Update 更新数据集
func (d *DatasetDao) Update(dataset *entity.Dataset) error {
	return d.getDB().Save(dataset).Error
}

// UpdateStatus 更新数据集状态
func (d *DatasetDao) UpdateStatus(id uint, status string) error {
	return d.getDB().Model(&entity.Dataset{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除数据集
func (d *DatasetDao) Delete(id uint) error {
	return d.getDB().Delete(&entity.Dataset{}, id).Error
}

// getDB 获取数据库连接
func (d *DatasetDao) getDB() *gorm.DB {
	if d.db != nil {
		return d.db
	}
	return database.GetDB()
}
