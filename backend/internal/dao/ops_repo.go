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

// ListAlerts 分页查询告警列表
// @author Claude
// @description 支持分页和筛选的告警列表查询
// @param severity 告警级别筛选（可选）
// @param acknowledged 是否已确认筛选（可选）
// @modified 2026-02-04
func (d *OpsDao) ListAlerts(ctx context.Context, page, pageSize int, severity string, acknowledged *bool) ([]entity.ActiveAlert, int64, error) {
	var alerts []entity.ActiveAlert
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.ActiveAlert{})

	// 筛选条件
	if severity != "" {
		query = query.Joins("JOIN alert_rules ON alert_rules.id = active_alerts.rule_id").
			Where("alert_rules.severity = ?", severity)
	}
	if acknowledged != nil {
		query = query.Where("acknowledged = ?", *acknowledged)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Rule").
		Order("triggered_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&alerts).Error

	return alerts, total, err
}

// AcknowledgeAlert 确认告警
// @author Claude
// @description 将告警标记为已确认
// @modified 2026-02-04
func (d *OpsDao) AcknowledgeAlert(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).
		Model(&entity.ActiveAlert{}).
		Where("id = ?", id).
		Update("acknowledged", true).Error
}

func (d *OpsDao) CreateAlert(ctx context.Context, alert *entity.ActiveAlert) error {
	return d.db.WithContext(ctx).Create(alert).Error
}