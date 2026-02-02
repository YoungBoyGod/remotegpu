package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// SyncTaskDao 同步任务数据访问对象
type SyncTaskDao struct {
	db *gorm.DB
}

// NewSyncTaskDao 创建同步任务DAO
func NewSyncTaskDao() *SyncTaskDao {
	return &SyncTaskDao{}
}

// NewSyncTaskDaoWithDB 创建同步任务DAO(带数据库连接)
func NewSyncTaskDaoWithDB(db *gorm.DB) *SyncTaskDao {
	return &SyncTaskDao{db: db}
}

// Create 创建同步任务
func (d *SyncTaskDao) Create(task *entity.SyncTask) error {
	return d.getDB().Create(task).Error
}

// GetByID 根据ID获取同步任务
func (d *SyncTaskDao) GetByID(id uint) (*entity.SyncTask, error) {
	var task entity.SyncTask
	err := d.getDB().Preload("Dataset").Preload("Source").Preload("Target").
		Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// ListByDataset 获取数据集的同步任务列表
func (d *SyncTaskDao) ListByDataset(datasetID uint, status string) ([]*entity.SyncTask, error) {
	var tasks []*entity.SyncTask
	query := d.getDB().Where("dataset_id = ?", datasetID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

// ListByStatus 根据状态获取同步任务列表
func (d *SyncTaskDao) ListByStatus(status string, limit int) ([]*entity.SyncTask, error) {
	var tasks []*entity.SyncTask
	query := d.getDB().Where("status = ?", status)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("created_at ASC").Find(&tasks).Error
	return tasks, err
}

// Update 更新同步任务
func (d *SyncTaskDao) Update(task *entity.SyncTask) error {
	return d.getDB().Save(task).Error
}

// UpdateStatus 更新任务状态
func (d *SyncTaskDao) UpdateStatus(id uint, status string) error {
	return d.getDB().Model(&entity.SyncTask{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateProgress 更新任务进度
func (d *SyncTaskDao) UpdateProgress(id uint, progress int, transferredSize int64, speed int64) error {
	return d.getDB().Model(&entity.SyncTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":         progress,
			"transferred_size": transferredSize,
			"speed":            speed,
		}).Error
}

// Delete 删除同步任务
func (d *SyncTaskDao) Delete(id uint) error {
	return d.getDB().Delete(&entity.SyncTask{}, id).Error
}

// getDB 获取数据库连接
func (d *SyncTaskDao) getDB() *gorm.DB {
	if d.db != nil {
		return d.db
	}
	return database.GetDB()
}
