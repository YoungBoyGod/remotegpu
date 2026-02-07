package entity

import "time"

// Document 文档实体
type Document struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Title          string    `gorm:"type:varchar(256);not null" json:"title"`
	Category       string    `gorm:"type:varchar(64);not null;default:'general'" json:"category"`
	FileName       string    `gorm:"type:varchar(256);not null" json:"file_name"`
	FilePath       string    `gorm:"type:varchar(512);not null" json:"file_path"`
	FileSize       int64     `gorm:"not null;default:0" json:"file_size"`
	ContentType    string    `gorm:"type:varchar(128);not null;default:''" json:"content_type"`
	StorageBackend string    `gorm:"type:varchar(64);not null;default:''" json:"storage_backend"`
	UploadedBy     uint      `gorm:"index" json:"uploaded_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	Uploader *Customer `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

func (Document) TableName() string {
	return "documents"
}
