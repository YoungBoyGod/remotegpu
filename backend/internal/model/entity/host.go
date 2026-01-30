package entity

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Host 主机实体
type Host struct {
	ID             string         `gorm:"primaryKey;size:64" json:"id"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	Hostname       string         `gorm:"size:256" json:"hostname"`
	IPAddress      string         `gorm:"column:ip_address;size:64;not null" json:"ip_address"`
	PublicIP       string         `gorm:"column:public_ip;size:64" json:"public_ip"`
	OSType         string         `gorm:"column:os_type;size:20;not null" json:"os_type"`
	OSVersion      string         `gorm:"column:os_version;size:64" json:"os_version"`
	Arch           string         `gorm:"size:20;default:'x86_64'" json:"arch"`
	DeploymentMode string         `gorm:"column:deployment_mode;size:20;not null" json:"deployment_mode"`
	K8sNodeName    string         `gorm:"column:k8s_node_name;size:128" json:"k8s_node_name"`
	Status         string         `gorm:"size:20;default:'offline'" json:"status"`
	HealthStatus   string         `gorm:"column:health_status;size:20;default:'unknown'" json:"health_status"`
	TotalCPU       int            `gorm:"column:total_cpu;not null" json:"total_cpu"`
	TotalMemory    int64          `gorm:"column:total_memory;not null" json:"total_memory"`
	TotalDisk      int64          `gorm:"column:total_disk" json:"total_disk"`
	TotalGPU       int            `gorm:"column:total_gpu;default:0" json:"total_gpu"`
	UsedCPU        int            `gorm:"column:used_cpu;default:0" json:"used_cpu"`
	UsedMemory     int64          `gorm:"column:used_memory;default:0" json:"used_memory"`
	UsedDisk       int64          `gorm:"column:used_disk;default:0" json:"used_disk"`
	UsedGPU        int            `gorm:"column:used_gpu;default:0" json:"used_gpu"`
	SSHPort        int            `gorm:"column:ssh_port;default:22" json:"ssh_port"`
	WinRMPort      *int           `gorm:"column:winrm_port" json:"winrm_port"`
	AgentPort      int            `gorm:"column:agent_port;default:8080" json:"agent_port"`
	Labels         datatypes.JSON `gorm:"type:jsonb" json:"labels"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags"`
	LastHeartbeat  *time.Time     `gorm:"column:last_heartbeat" json:"last_heartbeat"`
	RegisteredAt   time.Time      `gorm:"column:registered_at;default:now()" json:"registered_at"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Host) TableName() string {
	return "hosts"
}
