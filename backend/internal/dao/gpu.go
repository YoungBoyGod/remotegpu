package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// GPUDao GPU数据访问对象
type GPUDao struct {
	db *gorm.DB
}

// NewGPUDao 创建GPU DAO实例
func NewGPUDao() *GPUDao {
	return &GPUDao{
		db: database.GetDB(),
	}
}

// Create 创建GPU
func (d *GPUDao) Create(gpu *entity.GPU) error {
	return d.db.Create(gpu).Error
}

// GetByID 根据ID获取GPU
func (d *GPUDao) GetByID(id uint) (*entity.GPU, error) {
	var gpu entity.GPU
	err := d.db.Where("id = ?", id).First(&gpu).Error
	if err != nil {
		return nil, err
	}
	return &gpu, nil
}

// GetByHostID 根据主机ID获取GPU列表
func (d *GPUDao) GetByHostID(hostID string) ([]*entity.GPU, error) {
	var gpus []*entity.GPU
	err := d.db.Where("host_id = ?", hostID).Find(&gpus).Error
	return gpus, err
}

// Update 更新GPU
func (d *GPUDao) Update(gpu *entity.GPU) error {
	return d.db.Save(gpu).Error
}

// Delete 删除GPU
func (d *GPUDao) Delete(id uint) error {
	return d.db.Delete(&entity.GPU{}, id).Error
}

// DeleteByHostID 根据主机ID删除所有GPU
func (d *GPUDao) DeleteByHostID(hostID string) error {
	return d.db.Where("host_id = ?", hostID).Delete(&entity.GPU{}).Error
}

// UpdateStatus 更新GPU状态
func (d *GPUDao) UpdateStatus(id uint, status string) error {
	return d.db.Model(&entity.GPU{}).Where("id = ?", id).Update("status", status).Error
}

// List 分页获取GPU列表
func (d *GPUDao) List(page, pageSize int) ([]*entity.GPU, int64, error) {
	var gpus []*entity.GPU
	var total int64

	d.db.Model(&entity.GPU{}).Count(&total)

	offset := (page - 1) * pageSize
	err := d.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&gpus).Error
	return gpus, total, err
}

// GetByStatus 根据状态获取GPU列表
func (d *GPUDao) GetByStatus(status string) ([]*entity.GPU, error) {
	var gpus []*entity.GPU
	err := d.db.Where("status = ?", status).Find(&gpus).Error
	return gpus, err
}

// Allocate 分配GPU
func (d *GPUDao) Allocate(id uint, allocatedTo string) error {
	return d.db.Model(&entity.GPU{}).Where("id = ?", id).Updates(map[string]any{
		"status":       "allocated",
		"allocated_to": allocatedTo,
		"allocated_at": gorm.Expr("NOW()"),
	}).Error
}

// Release 释放GPU
func (d *GPUDao) Release(id uint) error {
	return d.db.Model(&entity.GPU{}).Where("id = ?", id).Updates(map[string]any{
		"status":       "available",
		"allocated_to": "",
		"allocated_at": nil,
	}).Error
}
