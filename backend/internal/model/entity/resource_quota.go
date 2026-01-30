package entity

import (
	"time"
)

// ResourceQuota 资源配额实体
type ResourceQuota struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CustomerID  uint      `gorm:"not null;uniqueIndex:idx_customer_workspace" json:"customer_id"`
	WorkspaceID *uint     `gorm:"uniqueIndex:idx_customer_workspace" json:"workspace_id"`
	CPU         int       `gorm:"not null;default:0;comment:CPU核心数配额" json:"cpu"`
	Memory      int64     `gorm:"not null;default:0;comment:内存配额(MB)" json:"memory"`
	GPU         int       `gorm:"not null;default:0;comment:GPU数量配额" json:"gpu"`
	Storage     int64     `gorm:"not null;default:0;comment:存储配额(GB)" json:"storage"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Customer  *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

// TableName 指定表名
func (ResourceQuota) TableName() string {
	return "resource_quotas"
}
