package entity

import (
	"time"

	"gorm.io/datatypes"
)

// Task represents a training or inference job
type Task struct {
	ID         string `gorm:"primarykey;type:varchar(64)" json:"id"`
	CustomerID uint   `gorm:"not null" json:"customer_id"`
	HostID     string `gorm:"type:varchar(64)" json:"host_id"`

	Name string `gorm:"type:varchar(256);not null" json:"name"`
	Type string `gorm:"type:varchar(32);not null" json:"type"` // training, inference

	ImageID  *uint          `json:"image_id,omitempty"`
	Command  string         `gorm:"type:text;not null" json:"command"`
	EnvVars  datatypes.JSON `gorm:"type:jsonb" json:"env_vars"`

	Status   string `gorm:"type:varchar(20);default:'queued'" json:"status"`
	ExitCode int    `json:"exit_code"`
	ErrorMsg string `gorm:"type:text" json:"error_msg"`

	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`

	// Relations
	Image *Image `gorm:"foreignKey:ImageID" json:"image,omitempty"`
}
