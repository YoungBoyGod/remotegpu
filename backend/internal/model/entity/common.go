package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity contains common columns for all tables
type BaseEntity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UUIDEntity for tables using UUID as an external identifier but bigserial as PK
type UUIDEntity struct {
	BaseEntity
	UUID string `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`
}
