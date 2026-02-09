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

// --- AlertRule CRUD ---

// CreateAlertRule 创建告警规则
func (d *OpsDao) CreateAlertRule(ctx context.Context, rule *entity.AlertRule) error {
	return d.db.WithContext(ctx).Create(rule).Error
}

// FindAlertRuleByID 根据 ID 查询告警规则
func (d *OpsDao) FindAlertRuleByID(ctx context.Context, id uint) (*entity.AlertRule, error) {
	var rule entity.AlertRule
	err := d.db.WithContext(ctx).First(&rule, "id = ?", id).Error
	return &rule, err
}

// UpdateAlertRule 更新告警规则
func (d *OpsDao) UpdateAlertRule(ctx context.Context, id uint, fields map[string]any) error {
	return d.db.WithContext(ctx).Model(&entity.AlertRule{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteAlertRule 删除告警规则
func (d *OpsDao) DeleteAlertRule(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&entity.AlertRule{}, "id = ?", id).Error
}

// ToggleAlertRule 切换告警规则启用状态
func (d *OpsDao) ToggleAlertRule(ctx context.Context, id uint) (*entity.AlertRule, error) {
	var rule entity.AlertRule
	err := d.db.WithContext(ctx).First(&rule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	rule.Enabled = !rule.Enabled
	if err := d.db.WithContext(ctx).Model(&rule).Update("enabled", rule.Enabled).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// ListAlertRules 分页查询告警规则
func (d *OpsDao) ListAlertRules(ctx context.Context, page, pageSize int, severity string, enabled *bool) ([]entity.AlertRule, int64, error) {
	var rules []entity.AlertRule
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.AlertRule{})
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at desc").
		Offset(offset).Limit(pageSize).
		Find(&rules).Error

	return rules, total, err
}