package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// HostDao 主机数据访问对象
type HostDao struct {
	db *gorm.DB
}

// NewHostDao 创建主机DAO实例
func NewHostDao() *HostDao {
	return &HostDao{
		db: database.GetDB(),
	}
}

// Create 创建主机
func (d *HostDao) Create(host *entity.Host) error {
	return d.db.Create(host).Error
}

// GetByID 根据ID获取主机
func (d *HostDao) GetByID(id string) (*entity.Host, error) {
	var host entity.Host
	err := d.db.Where("id = ?", id).First(&host).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

// GetByIPAddress 根据IP地址获取主机
func (d *HostDao) GetByIPAddress(ip string) (*entity.Host, error) {
	var host entity.Host
	err := d.db.Where("ip_address = ?", ip).First(&host).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

// Update 更新主机
func (d *HostDao) Update(host *entity.Host) error {
	return d.db.Save(host).Error
}

// Delete 删除主机
func (d *HostDao) Delete(id string) error {
	return d.db.Delete(&entity.Host{}, "id = ?", id).Error
}

// List 获取主机列表（分页）
func (d *HostDao) List(page, pageSize int) ([]*entity.Host, int64, error) {
	var hosts []*entity.Host
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&entity.Host{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Offset(offset).Limit(pageSize).Find(&hosts).Error; err != nil {
		return nil, 0, err
	}

	return hosts, total, nil
}

// ListByStatus 根据状态获取主机列表
func (d *HostDao) ListByStatus(status string) ([]*entity.Host, error) {
	var hosts []*entity.Host
	err := d.db.Where("status = ?", status).Find(&hosts).Error
	return hosts, err
}

// UpdateStatus 更新主机状态
func (d *HostDao) UpdateStatus(id, status string) error {
	return d.db.Model(&entity.Host{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateHeartbeat 更新心跳时间
func (d *HostDao) UpdateHeartbeat(id string) error {
	return d.db.Model(&entity.Host{}).Where("id = ?", id).Update("last_heartbeat", gorm.Expr("NOW()")).Error
}
