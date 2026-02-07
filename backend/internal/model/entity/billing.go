package entity

import (
	"time"
)

// BillingRecord 计费记录
type BillingRecord struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CustomerID   uint      `gorm:"not null" json:"customer_id"`
	EnvID        string    `gorm:"type:varchar(64)" json:"env_id"`
	ResourceType string    `gorm:"type:varchar(32);not null" json:"resource_type"`
	Quantity     float64   `gorm:"not null" json:"quantity"`
	UnitPrice    float64   `gorm:"type:decimal(10,4);not null" json:"unit_price"`
	Amount       float64   `gorm:"type:decimal(10,4);not null" json:"amount"`
	Currency     string    `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	StartTime    time.Time `gorm:"not null" json:"start_time"`
	EndTime      time.Time `gorm:"not null" json:"end_time"`
	CreatedAt    time.Time `json:"created_at"`
}

func (BillingRecord) TableName() string {
	return "billing_records"
}

// Invoice 账单
type Invoice struct {
	ID                 uint       `gorm:"primarykey" json:"id"`
	InvoiceNo          string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"invoice_no"`
	CustomerID         uint       `gorm:"not null" json:"customer_id"`
	BillingPeriodStart time.Time  `gorm:"not null" json:"billing_period_start"`
	BillingPeriodEnd   time.Time  `gorm:"not null" json:"billing_period_end"`
	TotalAmount        float64    `gorm:"type:decimal(10,4);not null" json:"total_amount"`
	Currency           string     `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	Status             string     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PaidAt             *time.Time `json:"paid_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Invoice) TableName() string {
	return "invoices"
}
