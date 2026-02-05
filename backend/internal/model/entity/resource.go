package entity

import (
	"time"
)

// Host 主机实体，表示物理机器或节点
type Host struct {
	ID       string `gorm:"primarykey;type:varchar(64)" json:"id"` // Machine ID, e.g., "node-01"
	Name     string `gorm:"type:varchar(128);not null" json:"name"`
	Hostname string `gorm:"type:varchar(256)" json:"hostname"`
	Region   string `gorm:"type:varchar(64);default:'default'" json:"region"`

	// Network
	IPAddress   string `gorm:"type:varchar(64);not null" json:"ip_address"` // Internal IP
	PublicIP    string `gorm:"type:varchar(64)" json:"public_ip"`           // External IP
	SSHPort     int    `gorm:"default:22" json:"ssh_port"`
	AgentPort   int    `gorm:"default:8080" json:"agent_port"`
	SSHUsername string `gorm:"type:varchar(128)" json:"ssh_username"`
	SSHPassword string `gorm:"type:text" json:"-"`
	SSHKey      string `gorm:"type:text" json:"-"`

	// Specs
	OSType        string `gorm:"type:varchar(20);default:'linux'" json:"os_type"`
	OSVersion     string `gorm:"type:varchar(64)" json:"os_version"`
	CPUInfo       string `gorm:"type:varchar(256)" json:"cpu_info"`
	TotalCPU      int    `gorm:"not null" json:"total_cpu"`
	TotalMemoryGB int64  `gorm:"not null" json:"total_memory_gb"`
	TotalDiskGB   int64  `json:"total_disk_gb"`

	// Status
	Status         string `gorm:"type:varchar(20);default:'offline'" json:"status"` // offline, idle, allocated, maintenance
	HealthStatus   string `gorm:"type:varchar(20);default:'unknown'" json:"health_status"`
	DeploymentMode string `gorm:"type:varchar(20);default:'traditional'" json:"deployment_mode"`
	NeedsCollect   bool   `gorm:"default:false" json:"needs_collect"`

	LastHeartbeat *time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relations
	GPUs        []GPU        `gorm:"foreignKey:HostID" json:"gpus,omitempty"`
	Allocations []Allocation `gorm:"foreignKey:HostID" json:"allocations,omitempty"`
}

// GPU GPU 设备实体，表示主机上的 GPU 设备
type GPU struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	HostID string `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_host_gpu_index" json:"host_id"`
	Index  int    `gorm:"not null;uniqueIndex:idx_host_gpu_index" json:"index"`
	UUID   string `gorm:"type:varchar(128);unique" json:"uuid"`

	Name          string `gorm:"type:varchar(128);not null" json:"name"`
	MemoryTotalMB int    `gorm:"not null" json:"memory_total_mb"`
	Brand         string `gorm:"type:varchar(64)" json:"brand"`

	// Status
	Status       string `gorm:"type:varchar(20);default:'available'" json:"status"` // available, allocated, error
	HealthStatus string `gorm:"type:varchar(20);default:'healthy'" json:"health_status"`
	AllocatedTo  string `gorm:"type:varchar(64)" json:"allocated_to"` // Can reference allocation_id

	UpdatedAt time.Time `json:"updated_at"`
}
