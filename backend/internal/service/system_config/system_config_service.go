package system_config

import (
	"context"
	"errors"
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"gorm.io/gorm"
)

var (
	ErrConfigNotFound    = errors.New("配置项不存在")
	ErrConfigKeyConflict = errors.New("配置键已存在")
)

type SystemConfigService struct {
	configDao    *dao.SystemConfigDao
	auditService *audit.AuditService
}

func NewSystemConfigService(db *gorm.DB) *SystemConfigService {
	return &SystemConfigService{
		configDao: dao.NewSystemConfigDao(db),
	}
}

// SetAuditService 注入审计服务（在 router 初始化时调用）
func (s *SystemConfigService) SetAuditService(auditSvc *audit.AuditService) {
	s.auditService = auditSvc
}

// GetAllConfigs 获取所有配置项
func (s *SystemConfigService) GetAllConfigs(ctx context.Context) ([]entity.SystemConfig, error) {
	return s.configDao.GetAll(ctx)
}

// GetConfigsByGroup 按分组获取配置项
func (s *SystemConfigService) GetConfigsByGroup(ctx context.Context, group string) ([]entity.SystemConfig, error) {
	return s.configDao.GetByGroup(ctx, group)
}

// ListGroups 获取所有配置分组
func (s *SystemConfigService) ListGroups(ctx context.Context) ([]string, error) {
	return s.configDao.ListGroups(ctx)
}

// GetConfig 获取单条配置
func (s *SystemConfigService) GetConfig(ctx context.Context, id uint) (*entity.SystemConfig, error) {
	config, err := s.configDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	return config, nil
}

// UpdateConfigs 批量更新配置值（带审计）
func (s *SystemConfigService) UpdateConfigs(ctx context.Context, updates map[string]string, operator string) error {
	// 记录变更前的值用于审计
	if s.auditService != nil {
		oldValues := make(map[string]string)
		for key := range updates {
			if old, err := s.configDao.GetByKey(ctx, key); err == nil {
				oldValues[key] = old.ConfigValue
			}
		}
		defer func() {
			for key, newVal := range updates {
				detail := map[string]interface{}{
					"config_key": key,
					"old_value":  oldValues[key],
					"new_value":  newVal,
				}
				_ = s.auditService.CreateLog(ctx, nil, operator, "", "PUT", "/admin/settings/configs",
					"update_config", "system_config", key, detail, 200)
			}
		}()
	}
	return s.configDao.BatchUpdate(ctx, updates)
}

// CreateConfig 创建配置项（带审计）
func (s *SystemConfigService) CreateConfig(ctx context.Context, config *entity.SystemConfig, operator string) error {
	// 检查 key 是否已存在
	if _, err := s.configDao.GetByKey(ctx, config.ConfigKey); err == nil {
		return ErrConfigKeyConflict
	}
	if config.ConfigGroup == "" {
		config.ConfigGroup = "general"
	}
	if config.ConfigType == "" {
		config.ConfigType = "string"
	}

	if err := s.configDao.Create(ctx, config); err != nil {
		return err
	}

	s.logAudit(ctx, operator, "create_config", config.ConfigKey, map[string]interface{}{
		"config_key":   config.ConfigKey,
		"config_value": config.ConfigValue,
		"config_group": config.ConfigGroup,
	})
	return nil
}

// UpdateConfig 更新单条配置项（带审计）
func (s *SystemConfigService) UpdateConfig(ctx context.Context, id uint, fields map[string]interface{}, operator string) error {
	old, err := s.configDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConfigNotFound
		}
		return err
	}

	oldValue := old.ConfigValue
	if v, ok := fields["config_value"]; ok {
		if str, ok := v.(string); ok {
			old.ConfigValue = str
		}
	}
	if v, ok := fields["description"]; ok {
		if str, ok := v.(string); ok {
			old.Description = str
		}
	}
	if v, ok := fields["config_group"]; ok {
		if str, ok := v.(string); ok {
			old.ConfigGroup = str
		}
	}
	if v, ok := fields["config_type"]; ok {
		if str, ok := v.(string); ok {
			old.ConfigType = str
		}
	}
	if v, ok := fields["is_public"]; ok {
		if b, ok := v.(bool); ok {
			old.IsPublic = b
		}
	}

	if err := s.configDao.Update(ctx, old); err != nil {
		return err
	}

	s.logAudit(ctx, operator, "update_config", old.ConfigKey, map[string]interface{}{
		"config_key": old.ConfigKey,
		"old_value":  oldValue,
		"new_value":  old.ConfigValue,
	})
	return nil
}

// DeleteConfig 删除配置项（带审计）
func (s *SystemConfigService) DeleteConfig(ctx context.Context, id uint, operator string) error {
	config, err := s.configDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConfigNotFound
		}
		return err
	}

	if err := s.configDao.Delete(ctx, id); err != nil {
		return err
	}

	s.logAudit(ctx, operator, "delete_config", config.ConfigKey, map[string]interface{}{
		"config_key":   config.ConfigKey,
		"config_value": config.ConfigValue,
	})
	return nil
}

// logAudit 记录审计日志的辅助方法
func (s *SystemConfigService) logAudit(ctx context.Context, operator, action, resourceID string, detail map[string]interface{}) {
	if s.auditService == nil {
		return
	}
	_ = s.auditService.CreateLog(ctx, nil, operator, "", "POST",
		fmt.Sprintf("/admin/settings/configs/%s", resourceID),
		action, "system_config", resourceID, detail, 200)
}
