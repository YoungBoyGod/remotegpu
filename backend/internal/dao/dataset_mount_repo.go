package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type DatasetMountDao struct {
	db *gorm.DB
}

func NewDatasetMountDao(db *gorm.DB) *DatasetMountDao {
	return &DatasetMountDao{db: db}
}

// Create 创建挂载记录
func (d *DatasetMountDao) Create(ctx context.Context, mount *entity.DatasetMount) error {
	return d.db.WithContext(ctx).Create(mount).Error
}

// FindByID 根据 ID 查询挂载记录
func (d *DatasetMountDao) FindByID(ctx context.Context, id uint) (*entity.DatasetMount, error) {
	var mount entity.DatasetMount
	if err := d.db.WithContext(ctx).First(&mount, id).Error; err != nil {
		return nil, err
	}
	return &mount, nil
}

// ListByDatasetID 查询数据集的所有挂载
func (d *DatasetMountDao) ListByDatasetID(ctx context.Context, datasetID uint) ([]entity.DatasetMount, error) {
	var mounts []entity.DatasetMount
	err := d.db.WithContext(ctx).
		Where("dataset_id = ? AND status != ?", datasetID, "unmounted").
		Order("created_at desc").Find(&mounts).Error
	return mounts, err
}

// ListByHostID 查询机器上的所有挂载
func (d *DatasetMountDao) ListByHostID(ctx context.Context, hostID string) ([]entity.DatasetMount, error) {
	var mounts []entity.DatasetMount
	err := d.db.WithContext(ctx).
		Where("host_id = ? AND status != ?", hostID, "unmounted").
		Order("created_at desc").Find(&mounts).Error
	return mounts, err
}

// FindActiveMount 查询指定数据集在指定机器上的活跃挂载
func (d *DatasetMountDao) FindActiveMount(ctx context.Context, datasetID uint, hostID string) (*entity.DatasetMount, error) {
	var mount entity.DatasetMount
	err := d.db.WithContext(ctx).
		Where("dataset_id = ? AND host_id = ? AND status NOT IN ?", datasetID, hostID, []string{"unmounted", "error"}).
		First(&mount).Error
	if err != nil {
		return nil, err
	}
	return &mount, nil
}

// UpdateStatus 更新挂载状态
func (d *DatasetMountDao) UpdateStatus(ctx context.Context, id uint, status, errMsg string) error {
	fields := map[string]interface{}{"status": status}
	if errMsg != "" {
		fields["error_message"] = errMsg
	}
	return d.db.WithContext(ctx).Model(&entity.DatasetMount{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除挂载记录
func (d *DatasetMountDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&entity.DatasetMount{}, id).Error
}

// CountByHostID 统计机器上的活跃挂载数
func (d *DatasetMountDao) CountByHostID(ctx context.Context, hostID string) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.DatasetMount{}).
		Where("host_id = ? AND status NOT IN ?", hostID, []string{"unmounted", "error"}).
		Count(&count).Error
	return count, err
}
