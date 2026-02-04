package audit

import (
	"context"
	"encoding/json"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditService struct {
	auditDao *dao.AuditDao
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{
		auditDao: dao.NewAuditDao(db),
	}
}

// LogAction 记录审计日志
func (s *AuditService) LogAction(ctx context.Context, log *entity.AuditLog) error {
	return s.auditDao.Create(ctx, log)
}

// ListLogs 查询审计日志
func (s *AuditService) ListLogs(ctx context.Context, params dao.AuditListParams) ([]entity.AuditLog, int64, error) {
	return s.auditDao.List(ctx, params)
}

// CreateLog 创建审计日志的便捷方法
func (s *AuditService) CreateLog(
	ctx context.Context,
	customerID *uint,
	username, ipAddress, method, path string,
	action, resourceType, resourceID string,
	detail map[string]interface{},
	statusCode int,
) error {
	// CodeX 2026-02-04: marshal detail explicitly to catch errors.
	var detailJSON datatypes.JSON
	if detail != nil {
		payload, err := json.Marshal(detail)
		if err != nil {
			return err
		}
		detailJSON = datatypes.JSON(payload)
	}

	log := &entity.AuditLog{
		CustomerID:   customerID,
		Username:     username,
		IPAddress:    ipAddress,
		Method:       method,
		Path:         path,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detailJSON,
		StatusCode:   statusCode,
	}

	return s.auditDao.Create(ctx, log)
}
