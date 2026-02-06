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
