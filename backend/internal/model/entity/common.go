package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity 基础实体，包含所有表的通用字段
type BaseEntity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UUIDEntity UUID 实体，用于使用 UUID 作为外部标识符但使用 bigserial 作为主键的表
type UUIDEntity struct {
	BaseEntity
	UUID string `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`
}
