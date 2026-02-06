package system_config

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type SystemConfigService struct {
	configDao *dao.SystemConfigDao
}

func NewSystemConfigService(db *gorm.DB) *SystemConfigService {
	return &SystemConfigService{
		configDao: dao.NewSystemConfigDao(db),
	}
}

// GetAllConfigs 获取所有配置项
func (s *SystemConfigService) GetAllConfigs(ctx context.Context) ([]entity.SystemConfig, error) {
	return s.configDao.GetAll(ctx)
}

// UpdateConfigs 批量更新配置值
func (s *SystemConfigService) UpdateConfigs(ctx context.Context, updates map[string]string) error {
	return s.configDao.BatchUpdate(ctx, updates)
}
