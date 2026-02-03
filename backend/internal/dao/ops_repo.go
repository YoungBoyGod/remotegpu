package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type OpsDao struct {
	db *gorm.DB
}

func NewOpsDao(db *gorm.DB) *OpsDao {
	return &OpsDao{db: db}
}

// Audit Logs
func (d *OpsDao) CreateAuditLog(ctx context.Context, log *entity.AuditLog) error {
	return d.db.WithContext(ctx).Create(log).Error
}

// Alerts
func (d *OpsDao) GetActiveAlerts(ctx context.Context) ([]entity.ActiveAlert, error) {
	var alerts []entity.ActiveAlert
	err := d.db.WithContext(ctx).
		Preload("Rule").
		Where("acknowledged = ?", false).
		Order("triggered_at desc").
		Find(&alerts).Error
	return alerts, err
}

func (d *OpsDao) CreateAlert(ctx context.Context, alert *entity.ActiveAlert) error {
	return d.db.WithContext(ctx).Create(alert).Error
}