package billing

import (
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"gorm.io/gorm"
)

// BillingService 计费服务
type BillingService struct {
	billingDao *dao.BillingRecordDao
	invoiceDao *dao.InvoiceDao
	db         *gorm.DB
}

func NewBillingService(db *gorm.DB) *BillingService {
	return &BillingService{
		billingDao: dao.NewBillingRecordDao(db),
		invoiceDao: dao.NewInvoiceDao(db),
		db:         db,
	}
}
