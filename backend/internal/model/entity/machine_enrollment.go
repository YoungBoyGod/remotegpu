package entity

import "time"

// MachineEnrollment 用户机器添加任务
type MachineEnrollment struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	CustomerID  uint   `gorm:"index" json:"customer_id"`
	Name        string `gorm:"type:varchar(128)" json:"name"`
	Hostname    string `gorm:"type:varchar(256)" json:"hostname"`
	Region      string `gorm:"type:varchar(64)" json:"region"`
	Address     string `gorm:"type:varchar(256);not null" json:"address"`
	SSHPort     int    `gorm:"default:22" json:"ssh_port"`
	SSHUsername string `gorm:"type:varchar(128)" json:"ssh_username"`
	SSHPassword string `gorm:"type:text" json:"-"`
	SSHKey      string `gorm:"type:text" json:"-"`

	Status       string `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"`
	HostID       string `gorm:"type:varchar(64)" json:"host_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
