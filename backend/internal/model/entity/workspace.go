package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Workspace 工作空间实体
type Workspace struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UUID        uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;default:uuid_generate_v4()" json:"uuid"`
	OwnerID     uint           `gorm:"column:owner_id;not null;index:idx_workspaces_owner" json:"owner_id"`
	Name        string         `gorm:"size:128;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Type        string         `gorm:"size:32;default:'personal';index:idx_workspaces_type;comment:类型 personal/team/enterprise" json:"type"`
	MemberCount int            `gorm:"default:1" json:"member_count"`
	Status      string         `gorm:"size:32;default:'active';index:idx_workspaces_status;comment:状态 active/archived" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Owner   *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []*WorkspaceMember `gorm:"foreignKey:WorkspaceID" json:"members,omitempty"`
}

// TableName 指定表名
func (Workspace) TableName() string {
	return "workspaces"
}

// WorkspaceMember 工作空间成员实体
type WorkspaceMember struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	WorkspaceID uint      `gorm:"not null;uniqueIndex:idx_workspace_user" json:"workspace_id"`
	UserID      uint      `gorm:"not null;uniqueIndex:idx_workspace_user;index:idx_workspace_members_user" json:"user_id"`
	Role        string    `gorm:"size:32;default:'member';index:idx_workspace_members_role;comment:角色 owner/admin/member/viewer" json:"role"`
	Status      string    `gorm:"size:32;default:'active';comment:状态 active/invited/suspended" json:"status"`
	JoinedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联关系
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (WorkspaceMember) TableName() string {
	return "workspace_members"
}
