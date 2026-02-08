package entity

import (
	"time"
)

// Workspace 工作空间实体，表示资源的逻辑分组（团队）
type Workspace struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UUID        string    `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`
	OwnerID     uint      `gorm:"not null" json:"owner_id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"type:varchar(32);default:'personal'" json:"type"`   // personal, team, enterprise
	MemberCount int       `gorm:"default:1" json:"member_count"`
	Status      string    `gorm:"type:varchar(32);default:'active'" json:"status"`   // active, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Owner   *Customer          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []WorkspaceMember  `gorm:"foreignKey:WorkspaceID" json:"members,omitempty"`
}

func (Workspace) TableName() string {
	return "workspaces"
}

// WorkspaceMember 工作空间成员实体
type WorkspaceMember struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	WorkspaceID uint      `gorm:"not null;uniqueIndex:idx_ws_customer" json:"workspace_id"`
	CustomerID  uint      `gorm:"not null;uniqueIndex:idx_ws_customer" json:"customer_id"`
	Role        string    `gorm:"type:varchar(32);default:'member'" json:"role"`   // owner, admin, member, viewer
	Status      string    `gorm:"type:varchar(32);default:'active'" json:"status"` // active, invited, suspended
	JoinedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Customer  *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}

func (WorkspaceMember) TableName() string {
	return "workspace_members"
}
