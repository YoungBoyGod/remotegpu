package entity

import (
	"time"
)

// Customer represents a user/tenant in the system
type Customer struct {
	BaseEntity
	UUID string `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`

	Username     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(256);not null" json:"-"`

	// Profile
	DisplayName string `gorm:"type:varchar(128)" json:"display_name"`
	FullName    string `gorm:"type:varchar(256)" json:"full_name"`
	Company     string `gorm:"type:varchar(256)" json:"company"`
	Phone       string `gorm:"type:varchar(32)" json:"phone"`
	AvatarURL   string `gorm:"type:varchar(512)" json:"avatar_url"`

	// Role & Type
	Role        string `gorm:"type:varchar(32);default:'customer_owner';index" json:"role"` // admin, customer_owner, customer_member
	UserType    string `gorm:"type:varchar(32);default:'external'" json:"user_type"`        // admin, external
	AccountType string `gorm:"type:varchar(32);default:'individual'" json:"account_type"`   // individual, enterprise

	// Status
	Status        string `gorm:"type:varchar(32);default:'active'" json:"status"` // active, suspended, deleted
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	PhoneVerified bool   `gorm:"default:false" json:"phone_verified"`

	// Billing
	Balance  float64 `gorm:"type:decimal(10,4);default:0.00" json:"balance"`
	Currency string  `gorm:"type:varchar(10);default:'CNY'" json:"currency"`

	LastLoginAt *time.Time `json:"last_login_at"`

	// Relations
	SSHKeys     []SSHKey     `gorm:"foreignKey:CustomerID" json:"ssh_keys,omitempty"`
	Allocations []Allocation `gorm:"foreignKey:CustomerID" json:"allocations,omitempty"`
	Workspaces  []Workspace  `gorm:"foreignKey:OwnerID" json:"workspaces,omitempty"`
}

// TableName overrides the table name used by User to `customers`
func (Customer) TableName() string {
	return "customers"
}

// SSHKey represents an SSH public key added by a customer
type SSHKey struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CustomerID uint      `gorm:"not null;index" json:"customer_id"`
	Name       string    `gorm:"type:varchar(64);not null" json:"name"`
	Fingerprint string    `gorm:"type:varchar(128)" json:"fingerprint"`
	PublicKey  string    `gorm:"type:text;not null" json:"public_key"`
	CreatedAt  time.Time `json:"created_at"`
}

// Workspace represents a logical grouping of resources (Team)
type Workspace struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UUID      string    `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`
	OwnerID   uint      `gorm:"not null" json:"owner_id"`
	Name      string    `gorm:"type:varchar(128);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Type      string    `gorm:"type:varchar(32);default:'personal'" json:"type"` // personal, team
	Status    string    `gorm:"type:varchar(32);default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
