package entity

import "time"

// SystemConfig 系统配置表
type SystemConfig struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ConfigKey   string    `gorm:"column:config_key;type:varchar(128);uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"column:config_value;type:text;not null" json:"config_value"`
	ConfigType  string    `gorm:"column:config_type;type:varchar(32);not null;default:string" json:"config_type"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	IsPublic    bool      `gorm:"column:is_public;default:false" json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}
