package entity

import (
	"time"
)

// ProxyNode 代理节点实体，表示一个 Proxy 服务实例
type ProxyNode struct {
	ID             string     `gorm:"primarykey;type:varchar(64)" json:"id"`
	Name           string     `gorm:"type:varchar(128);not null" json:"name"`
	Host           string     `gorm:"type:varchar(256);not null" json:"host"`
	APIPort        int        `gorm:"default:9090" json:"api_port"`
	HTTPPort       int        `gorm:"default:9091" json:"http_port"`
	RangeStart     int        `gorm:"default:20000" json:"range_start"`
	RangeEnd       int        `gorm:"default:60000" json:"range_end"`
	Version        string     `gorm:"type:varchar(32)" json:"version"`
	Status         string     `gorm:"type:varchar(20);default:'offline'" json:"status"`
	ActiveMappings int        `gorm:"default:0" json:"active_mappings"`
	UsedPorts      int        `gorm:"default:0" json:"used_ports"`
	LastHeartbeat  *time.Time `json:"last_heartbeat"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (ProxyNode) TableName() string {
	return "proxy_nodes"
}
