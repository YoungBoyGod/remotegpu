package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type OpsRepo struct {
	db *gorm.DB
}

func NewOpsRepo(db *gorm.DB) *OpsRepo {
	return &OpsRepo{db: db}
}

// Audit Logs
func (d *OpsRepo) CreateAuditLog(ctx context.Context, log *entity.AuditLog) error {
	return d.db.WithContext(ctx).Create(log).Error
}

// Alerts
func (d *OpsRepo) GetActiveAlerts(ctx context.Context) ([]entity.ActiveAlert, error) {
	var alerts []entity.ActiveAlert
	err := d.db.WithContext(ctx).
		Preload("Rule").
		Where("acknowledged = ?", false).
		Order("triggered_at desc").
		Find(&alerts).Error
	return alerts, err
}

func (d *OpsRepo) CreateAlert(ctx context.Context, alert *entity.ActiveAlert) error {
	return d.db.WithContext(ctx).Create(alert).Error
}
