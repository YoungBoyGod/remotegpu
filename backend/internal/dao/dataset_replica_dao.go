package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// DatasetReplicaDao 数据集副本数据访问对象
type DatasetReplicaDao struct {
	db *gorm.DB
}

// NewDatasetReplicaDao 创建数据集副本DAO
func NewDatasetReplicaDao() *DatasetReplicaDao {
	return &DatasetReplicaDao{}
}

// NewDatasetReplicaDaoWithDB 创建数据集副本DAO(带数据库连接)
func NewDatasetReplicaDaoWithDB(db *gorm.DB) *DatasetReplicaDao {
	return &DatasetReplicaDao{db: db}
}

// Create 创建数据集副本
func (d *DatasetReplicaDao) Create(replica *entity.DatasetReplica) error {
	return d.getDB().Create(replica).Error
}

// GetByID 根据ID获取数据集副本
func (d *DatasetReplicaDao) GetByID(id uint) (*entity.DatasetReplica, error) {
	var replica entity.DatasetReplica
	err := d.getDB().Preload("StorageSource").Where("id = ?", id).First(&replica).Error
	if err != nil {
		return nil, err
	}
	return &replica, nil
}

// ListByDataset 获取数据集的所有副本
func (d *DatasetReplicaDao) ListByDataset(datasetID uint) ([]*entity.DatasetReplica, error) {
	var replicas []*entity.DatasetReplica
	err := d.getDB().Preload("StorageSource").
		Where("dataset_id = ?", datasetID).
		Order("is_primary DESC, created_at DESC").
		Find(&replicas).Error
	return replicas, err
}

// ListByStorage 获取存储源上的所有副本
func (d *DatasetReplicaDao) ListByStorage(storageSourceID uint) ([]*entity.DatasetReplica, error) {
	var replicas []*entity.DatasetReplica
	err := d.getDB().Where("storage_source_id = ?", storageSourceID).
		Order("created_at DESC").
		Find(&replicas).Error
	return replicas, err
}

// GetPrimaryReplica 获取数据集的主副本
func (d *DatasetReplicaDao) GetPrimaryReplica(datasetID uint) (*entity.DatasetReplica, error) {
	var replica entity.DatasetReplica
	err := d.getDB().Preload("StorageSource").
		Where("dataset_id = ? AND is_primary = ?", datasetID, true).
		First(&replica).Error
	if err != nil {
		return nil, err
	}
	return &replica, nil
}

// Update 更新数据集副本
func (d *DatasetReplicaDao) Update(replica *entity.DatasetReplica) error {
	return d.getDB().Save(replica).Error
}

// UpdateStatus 更新副本状态
func (d *DatasetReplicaDao) UpdateStatus(id uint, status string) error {
	return d.getDB().Model(&entity.DatasetReplica{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除数据集副本
func (d *DatasetReplicaDao) Delete(id uint) error {
	return d.getDB().Delete(&entity.DatasetReplica{}, id).Error
}

// getDB 获取数据库连接
func (d *DatasetReplicaDao) getDB() *gorm.DB {
	if d.db != nil {
		return d.db
	}
	return database.GetDB()
}
