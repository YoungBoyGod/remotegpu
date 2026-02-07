package dao

import (
	"context"
	"time"

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
	if err := d.db.WithContext(ctx).
		Preload("GPUs").
		Preload("Allocations").
		Preload("Allocations.Customer").
		First(&host, "id = ?", id).Error; err != nil {
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
	if ds, ok := filters["device_status"]; ok && ds != "" {
		db = db.Where("device_status = ?", ds)
	}
	if as, ok := filters["allocation_status"]; ok && as != "" {
		db = db.Where("allocation_status = ?", as)
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

// UpdateDeviceStatus 更新设备在线状态
func (d *MachineDao) UpdateDeviceStatus(ctx context.Context, id string, deviceStatus string) error {
	return d.db.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", id).Update("device_status", deviceStatus).Error
}

// UpdateAllocationStatus 更新分配状态
func (d *MachineDao) UpdateAllocationStatus(ctx context.Context, id string, allocationStatus string) error {
	return d.db.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", id).Update("allocation_status", allocationStatus).Error
}

func (d *MachineDao) UpdateCollectFields(ctx context.Context, host *entity.Host) error {
	return d.db.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", host.ID).Updates(map[string]interface{}{
		"hostname":          host.Hostname,
		"name":              host.Name,
		"cpu_info":          host.CPUInfo,
		"total_cpu":         host.TotalCPU,
		"total_memory_gb":   host.TotalMemoryGB,
		"total_disk_gb":     host.TotalDiskGB,
		"needs_collect":     host.NeedsCollect,
		"status":            host.Status,
		"device_status":     host.DeviceStatus,
		"allocation_status": host.AllocationStatus,
		"health_status":     host.HealthStatus,
	}).Error
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
// @description 一次查询获取设备状态和分配状态的机器数量统计
// @return map[状态]数量
func (d *MachineDao) GetStatusStats(ctx context.Context) (map[string]int64, error) {
	type StatusCount struct {
		Status string
		Count  int64
	}

	stats := make(map[string]int64)

	// 按设备状态统计
	var deviceResults []StatusCount
	err := d.db.WithContext(ctx).Model(&entity.Host{}).
		Select("device_status as status, count(*) as count").
		Group("device_status").
		Scan(&deviceResults).Error
	if err != nil {
		return nil, err
	}
	for _, r := range deviceResults {
		stats["device:"+r.Status] = r.Count
	}

	// 按分配状态统计
	var allocResults []StatusCount
	err = d.db.WithContext(ctx).Model(&entity.Host{}).
		Select("allocation_status as status, count(*) as count").
		Group("allocation_status").
		Scan(&allocResults).Error
	if err != nil {
		return nil, err
	}
	for _, r := range allocResults {
		stats["allocation:"+r.Status] = r.Count
	}

	// 兼容旧字段
	var legacyResults []StatusCount
	err = d.db.WithContext(ctx).Model(&entity.Host{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&legacyResults).Error
	if err != nil {
		return nil, err
	}
	for _, r := range legacyResults {
		stats[r.Status] = r.Count
	}

	return stats, nil
}

func (d *MachineDao) ListNeedCollect(ctx context.Context, limit int) ([]entity.Host, error) {
	var hosts []entity.Host
	query := d.db.WithContext(ctx).Model(&entity.Host{}).Where("needs_collect = ?", true).Order("created_at asc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil
}

// UpdateHeartbeat 更新机器心跳时间和设备状态
func (d *MachineDao) UpdateHeartbeat(ctx context.Context, id string, deviceStatus string) error {
	now := time.Now()
	return d.db.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_heartbeat": now,
			"device_status":  deviceStatus,
		}).Error
}

// UpdateFields 更新机器指定字段
func (d *MachineDao) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除机器
func (d *MachineDao) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&entity.Host{}, "id = ?", id).Error
}

// ListAll 获取所有机器（不分页）
func (d *MachineDao) ListAll(ctx context.Context) ([]entity.Host, error) {
	var hosts []entity.Host
	err := d.db.WithContext(ctx).
		Preload("GPUs").
		Order("created_at desc").
		Find(&hosts).Error
	return hosts, err
}

// BatchUpdateAllocationStatus 批量更新分配状态
func (d *MachineDao) BatchUpdateAllocationStatus(ctx context.Context, ids []string, allocationStatus string) (int64, error) {
	result := d.db.WithContext(ctx).Model(&entity.Host{}).
		Where("id IN ?", ids).
		Update("allocation_status", allocationStatus)
	return result.RowsAffected, result.Error
}

// ListOnline 获取所有在线的机器
func (d *MachineDao) ListOnline(ctx context.Context) ([]entity.Host, error) {
	var hosts []entity.Host
	err := d.db.WithContext(ctx).
		Where("device_status = ?", "online").
		Find(&hosts).Error
	return hosts, err
}
