package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// WorkspaceDao 工作空间数据访问层
type WorkspaceDao struct {
	*BaseDao[entity.Workspace]
}

func NewWorkspaceDao(db *gorm.DB) *WorkspaceDao {
	return &WorkspaceDao{
		BaseDao: NewBaseDao[entity.Workspace](db),
	}
}

// FindByID 根据ID查询工作空间（带 Owner 预加载）
func (d *WorkspaceDao) FindByID(ctx context.Context, id uint) (*entity.Workspace, error) {
	var ws entity.Workspace
	if err := d.db.WithContext(ctx).Preload("Owner").First(&ws, id).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}

// FindByUUID 根据UUID查询工作空间
func (d *WorkspaceDao) FindByUUID(ctx context.Context, uuid string) (*entity.Workspace, error) {
	var ws entity.Workspace
	if err := d.db.WithContext(ctx).Where("uuid = ?", uuid).First(&ws).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}
