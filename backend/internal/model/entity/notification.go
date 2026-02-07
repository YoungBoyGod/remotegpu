package entity

import "time"

// Notification 通知实体
type Notification struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CustomerID uint       `gorm:"not null" json:"customer_id"`
	Title      string     `gorm:"type:varchar(256);not null" json:"title"`
	Content    string     `gorm:"type:text" json:"content"`
	Type       string     `gorm:"type:varchar(32);not null" json:"type"`
	Level      string     `gorm:"type:varchar(20);default:'info'" json:"level"`
	IsRead     bool       `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
