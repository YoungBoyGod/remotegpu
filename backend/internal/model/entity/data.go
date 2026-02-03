package entity

import (
	"time"
)

// Image 镜像实体，表示系统或自定义 Docker 镜像
type Image struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"type:varchar(256);not null;unique" json:"name"`
	DisplayName string `gorm:"type:varchar(256)" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`

	Category    string `gorm:"type:varchar(64)" json:"category"`
	Framework   string `gorm:"type:varchar(64)" json:"framework"`
	CUDAVersion string `gorm:"type:varchar(32)" json:"cuda_version"`

	RegistryURL string `gorm:"type:varchar(512)" json:"registry_url"`
	IsOfficial  bool   `gorm:"default:false" json:"is_official"`
	CustomerID  *uint  `json:"customer_id,omitempty"` // If private image

	Status    string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Dataset 数据集实体，表示数据文件的集合
type Dataset struct {
	BaseEntity
	UUID string `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex;not null" json:"uuid"`

	CustomerID  uint   `gorm:"not null" json:"customer_id"`
	WorkspaceID *uint  `json:"workspace_id,omitempty"`

	Name        string `gorm:"type:varchar(256);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	StoragePath string `gorm:"type:varchar(512);not null" json:"storage_path"`
	StorageType string `gorm:"type:varchar(32);default:'minio'" json:"storage_type"`
	TotalSize   int64  `gorm:"default:0" json:"total_size"`
	FileCount   int    `gorm:"default:0" json:"file_count"`

	Status     string `gorm:"type:varchar(20);default:'uploading'" json:"status"` // uploading, ready, error
	Visibility string `gorm:"type:varchar(20);default:'private'" json:"visibility"`

	// Relations
	DatasetMounts []DatasetMount `gorm:"foreignKey:DatasetID" json:"mounts,omitempty"`
}

// DatasetMount 数据集挂载实体，表示挂载到特定主机上的数据集
type DatasetMount struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	DatasetID uint   `gorm:"not null" json:"dataset_id"`
	HostID    string `gorm:"type:varchar(64);not null" json:"host_id"`
	MountPath string `gorm:"type:varchar(256);not null" json:"mount_path"`
	ReadOnly  bool   `gorm:"default:true" json:"read_only"`

	Status    string    `gorm:"type:varchar(20);default:'mounting'" json:"status"` // mounting, mounted, error, unmounted
	CreatedAt time.Time `json:"created_at"`
}
