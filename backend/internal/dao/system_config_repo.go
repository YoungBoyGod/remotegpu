package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type SystemConfigDao struct {
	db *gorm.DB
}

func NewSystemConfigDao(db *gorm.DB) *SystemConfigDao {
	return &SystemConfigDao{db: db}
}

// GetAll 获取所有配置项
func (d *SystemConfigDao) GetAll(ctx context.Context) ([]entity.SystemConfig, error) {
	var configs []entity.SystemConfig
	if err := d.db.WithContext(ctx).Order("id asc").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetByKey 按 key 查询配置
func (d *SystemConfigDao) GetByKey(ctx context.Context, key string) (*entity.SystemConfig, error) {
	var config entity.SystemConfig
	if err := d.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetByGroup 按分组查询配置
func (d *SystemConfigDao) GetByGroup(ctx context.Context, group string) ([]entity.SystemConfig, error) {
	var configs []entity.SystemConfig
	if err := d.db.WithContext(ctx).Where("config_group = ?", group).Order("id asc").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// ListGroups 获取所有配置分组名称
func (d *SystemConfigDao) ListGroups(ctx context.Context) ([]string, error) {
	var groups []string
	if err := d.db.WithContext(ctx).Model(&entity.SystemConfig{}).
		Distinct("config_group").
		Order("config_group asc").
		Pluck("config_group", &groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Create 创建配置项
func (d *SystemConfigDao) Create(ctx context.Context, config *entity.SystemConfig) error {
	return d.db.WithContext(ctx).Create(config).Error
}

// Update 更新单条配置项
func (d *SystemConfigDao) Update(ctx context.Context, config *entity.SystemConfig) error {
	return d.db.WithContext(ctx).Save(config).Error
}

// Delete 删除配置项
func (d *SystemConfigDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&entity.SystemConfig{}, id).Error
}

// GetByID 按 ID 查询配置
func (d *SystemConfigDao) GetByID(ctx context.Context, id uint) (*entity.SystemConfig, error) {
	var config entity.SystemConfig
	if err := d.db.WithContext(ctx).First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// BatchUpdate 批量更新配置值
func (d *SystemConfigDao) BatchUpdate(ctx context.Context, updates map[string]string) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range updates {
			result := tx.Model(&entity.SystemConfig{}).
				Where("config_key = ?", key).
				Update("config_value", value)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}
