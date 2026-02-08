package entity

import (
	"time"
)

// Environment 开发环境实体，表示用户在主机上创建的容器化开发环境
type Environment struct {
	ID          string `gorm:"primarykey;type:varchar(64)" json:"id"`
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	WorkspaceID *uint  `gorm:"index" json:"workspace_id"`
	HostID      string `gorm:"type:varchar(64);not null;index" json:"host_id"`
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Image       string `gorm:"type:varchar(256);not null" json:"image"`

	// 状态: creating, running, stopped, error, deleting
	Status string `gorm:"type:varchar(20);default:'creating';index" json:"status"`

	// 资源配置
	CPU     int   `gorm:"not null" json:"cpu"`
	Memory  int64 `gorm:"not null" json:"memory"`
	GPU     int   `gorm:"default:0" json:"gpu"`
	Storage int64 `json:"storage"`

	// 端口
	SSHPort     *int `json:"ssh_port"`
	RDPPort     *int `json:"rdp_port"`
	JupyterPort *int `json:"jupyter_port"`

	// 容器信息
	ContainerID string `gorm:"type:varchar(128)" json:"container_id"`
	PodName     string `gorm:"type:varchar(128)" json:"pod_name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	StartedAt *time.Time `json:"started_at"`
	StoppedAt *time.Time `json:"stopped_at"`

	// Relations
	Host         *Host          `gorm:"foreignKey:HostID" json:"host,omitempty"`
	PortMappings []PortMapping  `gorm:"foreignKey:EnvID" json:"port_mappings,omitempty"`
}

func (Environment) TableName() string {
	return "environments"
}

// PortMapping 端口映射实体，表示环境的端口映射关系
type PortMapping struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	EnvID        string    `gorm:"type:varchar(64);not null;index" json:"env_id"`
	ServiceType  string    `gorm:"type:varchar(32);not null" json:"service_type"` // ssh, rdp, jupyter, custom
	ExternalPort int       `gorm:"not null;uniqueIndex" json:"external_port"`
	InternalPort int       `gorm:"not null" json:"internal_port"`
	Status       string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, released
	AllocatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"allocated_at"`
	ReleasedAt   *time.Time `json:"released_at"`
}

func (PortMapping) TableName() string {
	return "port_mappings"
}
