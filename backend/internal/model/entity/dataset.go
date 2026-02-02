package entity

import (
	"time"

	"github.com/lib/pq"
)

// Dataset 数据集实体
type Dataset struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:128;not null;index:idx_datasets_name" json:"name"`
	UserID      uint           `gorm:"not null;index:idx_datasets_user" json:"user_id"`
	Size        int64          `gorm:"default:0;comment:大小(字节)" json:"size"`
	FileCount   int            `gorm:"default:0;comment:文件数量" json:"file_count"`
	Description string         `gorm:"type:text" json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Version     string         `gorm:"size:32;default:'1.0'" json:"version"`
	Status      string         `gorm:"size:20;default:'uploading';index:idx_datasets_status;comment:状态 uploading/ready/syncing/error" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at;index:idx_datasets_created_at,sort:desc" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`

	// 关联关系
	User     *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Replicas []*DatasetReplica `gorm:"foreignKey:DatasetID" json:"replicas,omitempty"`
}

// TableName 指定表名
func (Dataset) TableName() string {
	return "datasets"
}
