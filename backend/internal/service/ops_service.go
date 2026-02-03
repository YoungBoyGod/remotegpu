package service

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type OpsService struct {
	opsRepo *dao.OpsRepo
}

func NewOpsService(db *gorm.DB) *OpsService {
	return &OpsService{
		opsRepo: dao.NewOpsRepo(db),
	}
}

func (s *OpsService) ListActiveAlerts(ctx context.Context) ([]entity.ActiveAlert, error) {
	return s.opsRepo.GetActiveAlerts(ctx)
}

func (s *OpsService) LogAudit(ctx context.Context, log *entity.AuditLog) error {
	return s.opsRepo.CreateAuditLog(ctx, log)
}
