package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// InvoiceDao 账单数据访问
type InvoiceDao struct {
	db *gorm.DB
}

func NewInvoiceDao(db *gorm.DB) *InvoiceDao {
	return &InvoiceDao{db: db}
}

func (d *InvoiceDao) Create(ctx context.Context, invoice *entity.Invoice) error {
	return d.db.WithContext(ctx).Create(invoice).Error
}

func (d *InvoiceDao) FindByID(ctx context.Context, id uint) (*entity.Invoice, error) {
	var invoice entity.Invoice
	if err := d.db.WithContext(ctx).First(&invoice, id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (d *InvoiceDao) FindByInvoiceNo(ctx context.Context, invoiceNo string) (*entity.Invoice, error) {
	var invoice entity.Invoice
	err := d.db.WithContext(ctx).Where("invoice_no = ?", invoiceNo).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

// ListByCustomer 按客户查询账单（分页）
func (d *InvoiceDao) ListByCustomer(ctx context.Context, customerID uint, page, pageSize int, status string) ([]entity.Invoice, int64, error) {
	var invoices []entity.Invoice
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.Invoice{}).Where("customer_id = ?", customerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&invoices).Error
	return invoices, total, err
}

// ListAll 查询所有账单（管理员，分页）
func (d *InvoiceDao) ListAll(ctx context.Context, page, pageSize int, status string) ([]entity.Invoice, int64, error) {
	var invoices []entity.Invoice
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.Invoice{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&invoices).Error
	return invoices, total, err
}

// UpdateStatus 更新账单状态
func (d *InvoiceDao) UpdateStatus(ctx context.Context, id uint, status string) error {
	return d.db.WithContext(ctx).Model(&entity.Invoice{}).Where("id = ?", id).Update("status", status).Error
}

// MarkPaid 标记账单已支付
func (d *InvoiceDao) MarkPaid(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Model(&entity.Invoice{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  "paid",
			"paid_at": gorm.Expr("NOW()"),
		}).Error
}
