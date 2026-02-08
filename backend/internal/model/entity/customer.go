package entity

import (
	"time"
)

// Customer 客户实体，表示系统中的用户/租户
type Customer struct {
	BaseEntity
	UUID string `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`

	Username     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(256);not null" json:"-"`

	// Profile
	DisplayName string `gorm:"type:varchar(128)" json:"display_name"`
	FullName    string `gorm:"type:varchar(256)" json:"full_name"`
	CompanyCode string `gorm:"type:varchar(64)" json:"company_code"`
	Company     string `gorm:"type:varchar(256)" json:"company"`
	Phone       string `gorm:"type:varchar(32)" json:"phone"`
	AvatarURL   string `gorm:"type:varchar(512)" json:"avatar_url"`

	// Role & Type
	Role        string `gorm:"type:varchar(32);default:'customer_owner';index" json:"role"` // admin, customer_owner, customer_member
	UserType    string `gorm:"type:varchar(32);default:'external'" json:"user_type"`        // admin, external
	AccountType string `gorm:"type:varchar(32);default:'individual'" json:"account_type"`   // individual, enterprise

	// Status
	Status             string `gorm:"type:varchar(32);default:'active'" json:"status"` // active, suspended, deleted
	EmailVerified      bool   `gorm:"default:false" json:"email_verified"`
	PhoneVerified      bool   `gorm:"default:false" json:"phone_verified"`
	MustChangePassword bool   `gorm:"default:false" json:"must_change_password"`

	// 配额限制
	QuotaGPU     int   `gorm:"default:0" json:"quota_gpu"`
	QuotaStorage int64 `gorm:"default:0" json:"quota_storage"`

	LastLoginAt *time.Time `json:"last_login_at"`

	// Relations
	SSHKeys     []SSHKey     `gorm:"foreignKey:CustomerID" json:"ssh_keys,omitempty"`
	Allocations []Allocation `gorm:"foreignKey:CustomerID" json:"allocations,omitempty"`
	Workspaces  []Workspace  `gorm:"foreignKey:OwnerID" json:"workspaces,omitempty"`
}

// TableName 覆盖 User 使用的表名为 `customers`
func (Customer) TableName() string {
	return "customers"
}

// SSHKey SSH密钥实体，表示客户添加的 SSH 公钥
type SSHKey struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CustomerID  uint      `gorm:"not null;index" json:"customer_id"`
	Name        string    `gorm:"type:varchar(64);not null" json:"name"`
	Fingerprint string    `gorm:"type:varchar(128)" json:"fingerprint"`
	PublicKey   string    `gorm:"type:text;not null" json:"public_key"`
	CreatedAt   time.Time `json:"created_at"`
}

