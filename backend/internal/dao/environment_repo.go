package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// EnvironmentDao 开发环境数据访问层
type EnvironmentDao struct {
	db *gorm.DB
}

func NewEnvironmentDao(db *gorm.DB) *EnvironmentDao {
	return &EnvironmentDao{db: db}
}

// Create 创建环境
func (d *EnvironmentDao) Create(ctx context.Context, env *entity.Environment) error {
	return d.db.WithContext(ctx).Create(env).Error
}

// FindByID 根据ID查询环境（带 Host 和 PortMappings 预加载）
func (d *EnvironmentDao) FindByID(ctx context.Context, id string) (*entity.Environment, error) {
	var env entity.Environment
	err := d.db.WithContext(ctx).
		Preload("Host").
		Preload("PortMappings").
		First(&env, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}
