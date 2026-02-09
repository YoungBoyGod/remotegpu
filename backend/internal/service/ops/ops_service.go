package ops

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type OpsService struct {
	opsDao *dao.OpsDao
}

func NewOpsService(db *gorm.DB) *OpsService {
	return &OpsService{
		opsDao: dao.NewOpsDao(db),
	}
}

func (s *OpsService) ListActiveAlerts(ctx context.Context) ([]entity.ActiveAlert, error) {
	return s.opsDao.GetActiveAlerts(ctx)
}

// ListAlerts 分页查询告警列表
// @author Claude
// @description 支持分页和筛选的告警列表查询
// @modified 2026-02-04
func (s *OpsService) ListAlerts(ctx context.Context, page, pageSize int, severity string, acknowledged *bool) ([]entity.ActiveAlert, int64, error) {
	return s.opsDao.ListAlerts(ctx, page, pageSize, severity, acknowledged)
}

// AcknowledgeAlert 确认告警
// @author Claude
// @description 将告警标记为已确认
// @modified 2026-02-04
func (s *OpsService) AcknowledgeAlert(ctx context.Context, id uint) error {
	return s.opsDao.AcknowledgeAlert(ctx, id)
}

func (s *OpsService) LogAudit(ctx context.Context, log *entity.AuditLog) error {
	return s.opsDao.CreateAuditLog(ctx, log)
}

// --- 告警规则 CRUD ---

// CreateAlertRule 创建告警规则
func (s *OpsService) CreateAlertRule(ctx context.Context, rule *entity.AlertRule) error {
	return s.opsDao.CreateAlertRule(ctx, rule)
}

// GetAlertRule 获取告警规则详情
func (s *OpsService) GetAlertRule(ctx context.Context, id uint) (*entity.AlertRule, error) {
	return s.opsDao.FindAlertRuleByID(ctx, id)
}

// UpdateAlertRule 更新告警规则
func (s *OpsService) UpdateAlertRule(ctx context.Context, id uint, fields map[string]any) error {
	// 确认规则存在
	if _, err := s.opsDao.FindAlertRuleByID(ctx, id); err != nil {
		return err
	}
	return s.opsDao.UpdateAlertRule(ctx, id, fields)
}

// DeleteAlertRule 删除告警规则
func (s *OpsService) DeleteAlertRule(ctx context.Context, id uint) error {
	return s.opsDao.DeleteAlertRule(ctx, id)
}

// ToggleAlertRule 切换告警规则启用状态
func (s *OpsService) ToggleAlertRule(ctx context.Context, id uint) (*entity.AlertRule, error) {
	return s.opsDao.ToggleAlertRule(ctx, id)
}

// ListAlertRules 分页查询告警规则
func (s *OpsService) ListAlertRules(ctx context.Context, page, pageSize int, severity string, enabled *bool) ([]entity.AlertRule, int64, error) {
	return s.opsDao.ListAlertRules(ctx, page, pageSize, severity, enabled)
}