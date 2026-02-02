package entity

import "time"

// SyncTask 同步任务实体
type SyncTask struct {
	ID              uint       `gorm:"primarykey" json:"id"`
	DatasetID       uint       `gorm:"not null;index:idx_sync_tasks_dataset" json:"dataset_id"`
	SourceID        uint       `gorm:"not null;index:idx_sync_tasks_source;comment:源存储ID" json:"source_id"`
	TargetID        uint       `gorm:"not null;index:idx_sync_tasks_target;comment:目标存储ID" json:"target_id"`
	Type            string     `gorm:"size:20;default:'full';comment:类型 full/incremental" json:"type"`
	Status          string     `gorm:"size:20;default:'pending';index:idx_sync_tasks_status;comment:状态 pending/running/success/failed" json:"status"`
	Progress        int        `gorm:"default:0;comment:进度(0-100)" json:"progress"`
	TransferredSize int64      `gorm:"default:0;comment:已传输大小" json:"transferred_size"`
	TotalSize       int64      `gorm:"default:0;comment:总大小" json:"total_size"`
	Speed           int64      `gorm:"default:0;comment:传输速度(字节/秒)" json:"speed"`
	ErrorMessage    string     `gorm:"type:text" json:"error_message"`
	RetryCount      int        `gorm:"default:0;comment:重试次数" json:"retry_count"`
	StartedAt       *time.Time `gorm:"column:started_at" json:"started_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at" json:"completed_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;index:idx_sync_tasks_created_at,sort:desc" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// 关联关系
	Dataset *Dataset       `gorm:"foreignKey:DatasetID" json:"dataset,omitempty"`
	Source  *StorageSource `gorm:"foreignKey:SourceID" json:"source,omitempty"`
	Target  *StorageSource `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName 指定表名
func (SyncTask) TableName() string {
	return "sync_tasks"
}
