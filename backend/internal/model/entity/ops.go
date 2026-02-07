package entity

import (
	"time"

	"gorm.io/datatypes"
)

// AuditLog 审计日志实体，记录敏感操作
type AuditLog struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	CustomerID *uint  `json:"customer_id,omitempty"`
	Username   string `gorm:"type:varchar(128)" json:"username"`
	IPAddress  string `gorm:"type:varchar(64)" json:"ip_address"`
	Method     string `gorm:"type:varchar(10)" json:"method"`
	Path       string `gorm:"type:varchar(512)" json:"path"`

	Action       string         `gorm:"type:varchar(128);not null" json:"action"`
	ResourceType string         `gorm:"type:varchar(64)" json:"resource_type"`
	ResourceID   string         `gorm:"type:varchar(128)" json:"resource_id"`
	Detail       datatypes.JSON `gorm:"type:jsonb" json:"detail"`

	StatusCode int       `json:"status_code"`
	CreatedAt  time.Time `json:"created_at"`
}

// AlertRule 告警规则实体，定义触发告警的条件
type AlertRule struct {
	ID          uint    `gorm:"primarykey" json:"id"`
	Name        string  `gorm:"type:varchar(128);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	MetricType  string  `gorm:"type:varchar(64);not null" json:"metric_type"`
	Threshold   float64 `gorm:"not null" json:"threshold"`
	Condition   string  `gorm:"column:comparison;type:varchar(10);not null" json:"condition"`
	Duration    int     `gorm:"default:60" json:"duration"`
	Severity    string  `gorm:"type:varchar(20);default:'warning'" json:"severity"`

	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ActiveAlert 活动告警实体，表示当前触发的告警
type ActiveAlert struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	RuleID uint   `gorm:"not null" json:"rule_id"`
	HostID string `gorm:"type:varchar(64);not null" json:"host_id"`

	Value   float64 `gorm:"not null" json:"value"`
	Message string  `gorm:"type:text" json:"message"`

	TriggeredAt  time.Time `json:"triggered_at"`
	Acknowledged bool      `gorm:"default:false" json:"acknowledged"`

	// Relations
	Rule AlertRule `gorm:"foreignKey:RuleID" json:"rule"`
}
