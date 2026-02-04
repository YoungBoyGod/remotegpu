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

func (d *MachineDao) FindByIPAddress(ctx context.Context, ip string) (*entity.Host, error) {
	var host entity.Host
	if err := d.db.WithContext(ctx).Where("ip_address = ?", ip).First(&host).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

func (d *MachineDao) FindByHostname(ctx context.Context, hostname string) (*entity.Host, error) {
	var host entity.Host
	if err := d.db.WithContext(ctx).Where("hostname = ?", hostname).First(&host).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

type HostKey struct {
	IPAddress string
	Hostname  string
}

func (d *MachineDao) FindExistingKeys(ctx context.Context, ips []string, hostnames []string) (map[HostKey]entity.Host, error) {
	if len(ips) == 0 && len(hostnames) == 0 {
		return map[HostKey]entity.Host{}, nil
	}

	var hosts []entity.Host
	query := d.db.WithContext(ctx).Model(&entity.Host{})
	if len(ips) > 0 {
		query = query.Or("ip_address IN ?", ips)
	}
	if len(hostnames) > 0 {
		query = query.Or("hostname IN ?", hostnames)
	}
	if err := query.Find(&hosts).Error; err != nil {
		return nil, err
	}

	results := make(map[HostKey]entity.Host, len(hosts))
	for _, host := range hosts {
		results[HostKey{IPAddress: host.IPAddress, Hostname: host.Hostname}] = host
	}
	return results, nil
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

// CountByStatus 按状态统计机器数量
// @description 统计指定状态的机器数量，用于仪表盘展示
// @param status 机器状态 (idle/allocated/maintenance/offline)
// @return 数量和错误
// @modified 2026-02-04
func (d *MachineDao) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Host{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetStatusStats 获取各状态机器统计
// @description 一次查询获取所有状态的机器数量统计
// @return map[状态]数量
// @modified 2026-02-04
func (d *MachineDao) GetStatusStats(ctx context.Context) (map[string]int64, error) {
	type StatusCount struct {
		Status string
		Count  int64
	}
	var results []StatusCount

	err := d.db.WithContext(ctx).Model(&entity.Host{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}
	return stats, nil
}
