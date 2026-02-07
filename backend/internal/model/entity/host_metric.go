package entity

import "time"

// HostMetric 机器监控指标
// @author Claude
// @description 存储机器的CPU、内存、磁盘等监控数据
// @modified 2026-02-06
type HostMetric struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	HostID string `gorm:"column:host_id;type:varchar(255);not null;index" json:"host_id"`

	// CPU 指标
	CPUUsagePercent *float64 `gorm:"column:cpu_usage_percent;type:decimal(5,2)" json:"cpu_usage_percent,omitempty"`
	CPUCoresUsed    *float64 `gorm:"column:cpu_cores_used;type:decimal(10,2)" json:"cpu_cores_used,omitempty"`

	// 内存指标
	MemoryTotalGB      *int64   `gorm:"column:memory_total_gb" json:"memory_total_gb,omitempty"`
	MemoryUsedGB       *int64   `gorm:"column:memory_used_gb" json:"memory_used_gb,omitempty"`
	MemoryUsagePercent *float64 `gorm:"column:memory_usage_percent;type:decimal(5,2)" json:"memory_usage_percent,omitempty"`

	// 磁盘指标
	DiskTotalGB      *int64   `gorm:"column:disk_total_gb" json:"disk_total_gb,omitempty"`
	DiskUsedGB       *int64   `gorm:"column:disk_used_gb" json:"disk_used_gb,omitempty"`
	DiskUsagePercent *float64 `gorm:"column:disk_usage_percent;type:decimal(5,2)" json:"disk_usage_percent,omitempty"`

	// 网络指标
	NetworkRxBytes *int64 `gorm:"column:network_rx_bytes" json:"network_rx_bytes,omitempty"`
	NetworkTxBytes *int64 `gorm:"column:network_tx_bytes" json:"network_tx_bytes,omitempty"`

	// 采集时间
	CollectedAt time.Time `gorm:"column:collected_at;not null;index" json:"collected_at"`
}

// TableName 指定表名
func (HostMetric) TableName() string {
	return "host_metrics"
}
