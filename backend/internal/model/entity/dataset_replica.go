package entity

import "time"

// DatasetReplica 数据集副本实体
type DatasetReplica struct {
	ID              uint       `gorm:"primarykey" json:"id"`
	DatasetID       uint       `gorm:"not null;index:idx_dataset_replicas_dataset" json:"dataset_id"`
	StorageSourceID uint       `gorm:"not null;index:idx_dataset_replicas_storage" json:"storage_source_id"`
	Path            string     `gorm:"size:512;not null;comment:存储路径" json:"path"`
	Size            int64      `gorm:"default:0;comment:副本大小" json:"size"`
	Status          string     `gorm:"size:20;default:'syncing';index:idx_dataset_replicas_status;comment:状态 syncing/ready/error" json:"status"`
	IsPrimary       bool       `gorm:"default:false;comment:是否主副本" json:"is_primary"`
	LastSyncAt      *time.Time `gorm:"column:last_sync_at" json:"last_sync_at"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// 关联关系
	Dataset       *Dataset       `gorm:"foreignKey:DatasetID" json:"dataset,omitempty"`
	StorageSource *StorageSource `gorm:"foreignKey:StorageSourceID" json:"storage_source,omitempty"`
}

// TableName 指定表名
func (DatasetReplica) TableName() string {
	return "dataset_replicas"
}
