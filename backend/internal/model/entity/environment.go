package entity

import (
	"time"
)

// PortMapping 端口映射实体
type PortMapping struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	EnvID        string     `gorm:"column:env_id;size:64;not null;index:idx_port_mappings_env" json:"env_id"`
	ServiceType  string     `gorm:"column:service_type;size:32;not null;comment:服务类型 ssh/rdp/jupyter/custom" json:"service_type"`
	ExternalPort int        `gorm:"column:external_port;not null;uniqueIndex:idx_port_mappings_external_port" json:"external_port"`
	InternalPort int        `gorm:"column:internal_port;not null" json:"internal_port"`
	Status       string     `gorm:"size:20;default:'active';index:idx_port_mappings_status;comment:状态 active/released" json:"status"`
	AllocatedAt  time.Time  `gorm:"column:allocated_at;default:now()" json:"allocated_at"`
	ReleasedAt   *time.Time `gorm:"column:released_at" json:"released_at"`

	// 关联关系
	Environment *Environment `gorm:"foreignKey:EnvID" json:"environment,omitempty"`
}

// TableName 指定表名
func (PortMapping) TableName() string {
	return "port_mappings"
}

// Environment 开发环境实体
type Environment struct {
	ID           string     `gorm:"primaryKey;size:64" json:"id"`
	UserID       uint       `gorm:"column:user_id;not null;index:idx_environments_user" json:"user_id"`
	WorkspaceID  *uint      `gorm:"column:workspace_id;index:idx_environments_workspace" json:"workspace_id"`
	HostID       string     `gorm:"column:host_id;size:64;not null;index:idx_environments_host" json:"host_id"`
	Name         string     `gorm:"size:128;not null" json:"name"`
	Description  string     `gorm:"type:text" json:"description"`
	Image        string     `gorm:"size:256;not null" json:"image"`
	Status       string     `gorm:"size:20;default:'creating';index:idx_environments_status;comment:状态 creating/running/stopped/error/deleting" json:"status"`
	CPU          int        `gorm:"not null" json:"cpu"`
	Memory       int64      `gorm:"not null" json:"memory"`
	GPU          int        `gorm:"default:0" json:"gpu"`
	Storage      *int64     `gorm:"" json:"storage"`
	SSHPort      *int       `gorm:"column:ssh_port" json:"ssh_port"`
	RDPPort      *int       `gorm:"column:rdp_port" json:"rdp_port"`
	JupyterPort  *int       `gorm:"column:jupyter_port" json:"jupyter_port"`
	ContainerID  string     `gorm:"column:container_id;size:128" json:"container_id"`
	PodName      string     `gorm:"column:pod_name;size:128" json:"pod_name"`
	CreatedAt    time.Time  `gorm:"column:created_at;index:idx_environments_created_at,sort:desc" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	StartedAt    *time.Time `gorm:"column:started_at" json:"started_at"`
	StoppedAt    *time.Time `gorm:"column:stopped_at" json:"stopped_at"`

	// 关联关系
	User         *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Workspace    *Workspace     `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Host         *Host          `gorm:"foreignKey:HostID" json:"host,omitempty"`
	PortMappings []*PortMapping `gorm:"foreignKey:EnvID" json:"port_mappings,omitempty"`
}

// TableName 指定表名
func (Environment) TableName() string {
	return "environments"
}
