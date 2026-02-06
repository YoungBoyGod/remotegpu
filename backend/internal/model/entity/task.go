package entity

import (
	"time"

	"gorm.io/datatypes"
)

// Task 任务实体
type Task struct {
	ID         string `gorm:"primarykey;type:varchar(64)" json:"id"`
	CustomerID uint   `gorm:"not null" json:"customer_id"`
	Name       string `gorm:"type:varchar(256);not null" json:"name"`
	Type       string `gorm:"type:varchar(32);not null;default:'shell'" json:"type"`

	// 执行信息
	Command string         `gorm:"type:text;not null" json:"command"`
	Args    datatypes.JSON `gorm:"type:jsonb" json:"args"`
	WorkDir string         `gorm:"type:varchar(500)" json:"workdir"`
	EnvVars datatypes.JSON `gorm:"type:jsonb" json:"env_vars"`
	Timeout int            `gorm:"default:3600" json:"timeout"`

	// 优先级和重试
	Priority   int `gorm:"default:5" json:"priority"`
	RetryCount int `gorm:"default:0" json:"retry_count"`
	RetryDelay int `gorm:"default:60" json:"retry_delay"`
	MaxRetries int `gorm:"default:3" json:"max_retries"`

	// 状态
	Status   string `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ExitCode int    `json:"exit_code"`
	ErrorMsg string `gorm:"type:text" json:"error_msg"`

	// 关联
	MachineID string `gorm:"type:varchar(64)" json:"machine_id"`
	GroupID   string `gorm:"type:varchar(64)" json:"group_id"`
	ParentID  string `gorm:"type:varchar(64)" json:"parent_id"`

	// 调度与租约
	AssignedAgentID string     `gorm:"type:varchar(64)" json:"assigned_agent_id"`
	LeaseExpiresAt  *time.Time `json:"lease_expires_at"`
	AttemptID       string     `gorm:"type:varchar(64)" json:"attempt_id"`

	// 时间戳
	CreatedAt  time.Time  `json:"created_at"`
	AssignedAt *time.Time `json:"assigned_at"`
	StartedAt  *time.Time `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at"`

	// 兼容旧字段
	HostID    string `gorm:"type:varchar(64)" json:"host_id"`
	ProcessID int    `gorm:"type:int;default:0" json:"process_id,omitempty"`
	ImageID   *uint  `json:"image_id,omitempty"`

	// Relations
	Image *Image `gorm:"foreignKey:ImageID" json:"image,omitempty"`
}
