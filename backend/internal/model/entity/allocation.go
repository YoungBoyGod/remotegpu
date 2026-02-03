package entity

import (
	"time"
)

// Allocation 资源分配实体，表示分配给客户的资源（租约）
type Allocation struct {
	ID          string `gorm:"primarykey;type:varchar(64)" json:"id"`
	CustomerID  uint   `gorm:"not null;index" json:"customer_id"`
	HostID      string `gorm:"type:varchar(64);not null;index" json:"host_id"`
	WorkspaceID *uint  `json:"workspace_id,omitempty"`

	// Time
	StartTime     time.Time  `gorm:"not null" json:"start_time"`
	EndTime       time.Time  `gorm:"not null" json:"end_time"`
	ActualEndTime *time.Time `json:"actual_end_time,omitempty"`

	// Status
	Status string `gorm:"type:varchar(32);default:'active';index" json:"status"` // active, expired, reclaimed, pending

	// Metadata
	Remark    string    `gorm:"type:text" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Customer Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Host     Host      `gorm:"foreignKey:HostID" json:"host,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}
