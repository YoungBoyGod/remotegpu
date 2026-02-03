package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type MachineDao struct {
	db *gorm.DB
}

func NewMachineDao(db *gorm.DB) *MachineDao {
	return &MachineDao{db: db}
}

func (d *MachineDao) Create(ctx context.Context, host *entity.Host) error {
	return d.db.WithContext(ctx).Create(host).Error
}

func (d *MachineDao) FindByID(ctx context.Context, id string) (*entity.Host, error) {
	var host entity.Host
	if err := d.db.WithContext(ctx).Preload("GPUs").First(&host, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

func (d *MachineDao) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entity.Host, int64, error) {
	var hosts []entity.Host
	var total int64

	db := d.db.WithContext(ctx).Model(&entity.Host{})

	if status, ok := filters["status"]; ok && status != "" {
		db = db.Where("status = ?", status)
	}
	if region, ok := filters["region"]; ok && region != "" {
		db = db.Where("region = ?", region)
	}
	if gpuModel, ok := filters["gpu_model"]; ok && gpuModel != "" {
		db = db.Joins("JOIN gpus ON gpus.host_id = hosts.id").Where("gpus.name LIKE ?", "%"+gpuModel.(string)+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("GPUs").Offset((page - 1) * pageSize).Limit(pageSize).Find(&hosts).Error; err != nil {
		return nil, 0, err
	}

	return hosts, total, nil
}

func (d *MachineDao) UpdateStatus(ctx context.Context, id string, status string) error {
	return d.db.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", id).Update("status", status).Error
}

func (d *MachineDao) Count(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Host{}).Count(&count).Error
	return count, err
}
