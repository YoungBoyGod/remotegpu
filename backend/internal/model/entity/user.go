package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UUID          uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;default:uuid_generate_v4()" json:"uuid"`
	Username      string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Email         string         `gorm:"uniqueIndex;size:128;not null" json:"email"`
	PasswordHash  string         `gorm:"column:password_hash;size:256;not null" json:"-"`
	DisplayName   string         `gorm:"size:128" json:"display_name"`
	AvatarURL     string         `gorm:"column:avatar_url;size:512" json:"avatar_url"`
	Phone         string         `gorm:"size:32" json:"phone"`
	FullName      string         `gorm:"size:256" json:"full_name"`
	Company       string         `gorm:"size:256" json:"company"`
	UserType      string         `gorm:"size:32;default:'external';comment:用户类型 admin/internal/external" json:"user_type"`
	AccountType   string         `gorm:"size:32;default:'individual';comment:账户类型 individual/enterprise" json:"account_type"`
	Status        string         `gorm:"size:32;default:'active';comment:状态 active/suspended/deleted" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	PhoneVerified bool           `gorm:"default:false" json:"phone_verified"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
