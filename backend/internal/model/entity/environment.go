package entity

import (
	"time"
)

// PortMapping 端口映射实体
type PortMapping struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	EnvID        string     `gorm:"column:env_id;size:64;not null;index:idx_port_mappings_env" json:"env_id"`
	ServiceType  string     `gorm:"column:service_type;size:32;not null;comment:服务类型 ssh/rdp/jupyter/vnc/tensorboard/vscode/custom" json:"service_type"`
	InternalPort int        `gorm:"column:internal_port;not null" json:"internal_port"`
	ExternalPort int        `gorm:"column:external_port;not null;uniqueIndex:idx_port_mappings_external_port" json:"external_port"`
	PublicPort   *int       `gorm:"column:public_port;comment:公网端口(防火墙映射后)" json:"public_port,omitempty"`
	Protocol     string     `gorm:"column:protocol;size:10;default:'tcp';comment:协议 tcp/udp" json:"protocol"`
	Description  string     `gorm:"column:description;size:256" json:"description"`

	// 访问地址
	InternalAccessURL string `gorm:"column:internal_access_url;size:512;comment:内网访问地址" json:"internal_access_url,omitempty"`
	PublicAccessURL   string `gorm:"column:public_access_url;size:512;comment:公网访问地址" json:"public_access_url,omitempty"`
	PublicDomain      string `gorm:"column:public_domain;size:256;comment:公网域名" json:"public_domain,omitempty"`

	// 状态
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
	AccessInfo   []byte     `gorm:"column:access_info;type:jsonb" json:"access_info,omitempty"`

	// 部署和配置相关字段
	DeploymentMode   string `gorm:"column:deployment_mode;size:32;default:'k8s_pod'" json:"deployment_mode"`
	EnvironmentType  string `gorm:"column:environment_type;size:32;default:'ide'" json:"environment_type"`
	GPUMode          string `gorm:"column:gpu_mode;size:32;default:'exclusive'" json:"gpu_mode"`
	StorageType      string `gorm:"column:storage_type;size:32;default:'local'" json:"storage_type"`
	StorageConfig    []byte `gorm:"column:storage_config;type:jsonb" json:"storage_config,omitempty"`
	LifecyclePolicy  string `gorm:"column:lifecycle_policy;size:32;default:'persistent'" json:"lifecycle_policy"`
	NetworkConfig    []byte `gorm:"column:network_config;type:jsonb" json:"network_config,omitempty"`
	UseJumpserver    bool   `gorm:"column:use_jumpserver;default:false" json:"use_jumpserver"`
	UseGuacamole     bool   `gorm:"column:use_guacamole;default:false" json:"use_guacamole"`
	GuacamoleConnID  string `gorm:"column:guacamole_conn_id;size:64" json:"guacamole_conn_id,omitempty"`
	VNCPort          *int   `gorm:"column:vnc_port" json:"vnc_port"`
	VNCPassword      string `gorm:"column:vnc_password;size:128" json:"vnc_password,omitempty"`
	AdditionalPorts  []byte `gorm:"column:additional_ports;type:jsonb" json:"additional_ports,omitempty"`

	// 关联关系
	User         *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Host         *Host          `gorm:"foreignKey:HostID" json:"host,omitempty"`
	PortMappings []*PortMapping `gorm:"foreignKey:EnvID" json:"port_mappings,omitempty"`
}

// TableName 指定表名
func (Environment) TableName() string {
	return "environments"
}
