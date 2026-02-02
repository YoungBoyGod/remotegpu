package entity

import "time"

// StorageSource 存储源实体
type StorageSource struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:128;not null;uniqueIndex:idx_storage_sources_name" json:"name"`
	Type        string    `gorm:"size:32;not null;comment:类型 oss/minio/s3/rustfs" json:"type"`
	Endpoint    string    `gorm:"size:256;not null" json:"endpoint"`
	AccessKey   string    `gorm:"size:256;not null" json:"-"`
	SecretKey   string    `gorm:"size:256;not null" json:"-"`
	Bucket      string    `gorm:"size:128;not null" json:"bucket"`
	Region      string    `gorm:"size:64" json:"region"`
	IsPublic    bool      `gorm:"default:false;comment:是否公有云" json:"is_public"`
	Priority    int       `gorm:"default:0;comment:优先级(数字越小优先级越高)" json:"priority"`
	Status      string    `gorm:"size:20;default:'active';index:idx_storage_sources_status;comment:状态 active/inactive/error" json:"status"`
	HealthCheck *time.Time `gorm:"column:health_check" json:"health_check"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (StorageSource) TableName() string {
	return "storage_sources"
}
