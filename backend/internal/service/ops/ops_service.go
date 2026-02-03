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

func (s *OpsService) LogAudit(ctx context.Context, log *entity.AuditLog) error {
	return s.opsDao.CreateAuditLog(ctx, log)
}