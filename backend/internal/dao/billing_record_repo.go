package dao

import (
	"context"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// BillingRecordDao 计费记录数据访问
type BillingRecordDao struct {
	db *gorm.DB
}

func NewBillingRecordDao(db *gorm.DB) *BillingRecordDao {
	return &BillingRecordDao{db: db}
}

func (d *BillingRecordDao) Create(ctx context.Context, record *entity.BillingRecord) error {
	return d.db.WithContext(ctx).Create(record).Error
}

func (d *BillingRecordDao) BatchCreate(ctx context.Context, records []entity.BillingRecord) error {
	if len(records) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Create(&records).Error
}

// ListByCustomer 按客户查询计费记录（分页）
func (d *BillingRecordDao) ListByCustomer(ctx context.Context, customerID uint, page, pageSize int) ([]entity.BillingRecord, int64, error) {
	var records []entity.BillingRecord
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.BillingRecord{}).Where("customer_id = ?", customerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// SumByCustomerAndPeriod 汇总指定客户在某时间段内的计费金额
func (d *BillingRecordDao) SumByCustomerAndPeriod(ctx context.Context, customerID uint, start, end time.Time) (float64, error) {
	var total float64
	err := d.db.WithContext(ctx).
		Model(&entity.BillingRecord{}).
		Where("customer_id = ? AND start_time >= ? AND end_time <= ?", customerID, start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// ListByPeriod 查询指定时间段内的所有计费记录
func (d *BillingRecordDao) ListByPeriod(ctx context.Context, start, end time.Time) ([]entity.BillingRecord, error) {
	var records []entity.BillingRecord
	err := d.db.WithContext(ctx).
		Where("start_time >= ? AND end_time <= ?", start, end).
		Order("created_at DESC").
		Find(&records).Error
	return records, err
}
