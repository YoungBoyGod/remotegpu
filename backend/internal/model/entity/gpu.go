package entity

import "time"

// GPU GPU设备实体
type GPU struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	HostID            string     `gorm:"column:host_id;size:64;not null;index" json:"host_id"`
	GPUIndex          int        `gorm:"column:gpu_index;not null" json:"gpu_index"`
	UUID              string     `gorm:"size:128;uniqueIndex" json:"uuid"`
	Name              string     `gorm:"size:128" json:"name"`
	Brand             string     `gorm:"size:64" json:"brand"`
	Architecture      string     `gorm:"size:64" json:"architecture"`
	MemoryTotal       int64      `gorm:"column:memory_total" json:"memory_total"`
	CUDACores         int        `gorm:"column:cuda_cores" json:"cuda_cores"`
	ComputeCapability string     `gorm:"column:compute_capability;size:32" json:"compute_capability"`
	Status            string     `gorm:"size:20;default:'available'" json:"status"`
	HealthStatus      string     `gorm:"column:health_status;size:20;default:'healthy'" json:"health_status"`
	AllocatedTo       string     `gorm:"column:allocated_to;size:64" json:"allocated_to"`
	AllocatedAt       *time.Time `gorm:"column:allocated_at" json:"allocated_at"`
	PowerLimit        int        `gorm:"column:power_limit" json:"power_limit"`
	TemperatureLimit  int        `gorm:"column:temperature_limit" json:"temperature_limit"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (GPU) TableName() string {
	return "gpus"
}
