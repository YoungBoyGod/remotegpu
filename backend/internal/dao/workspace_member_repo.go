package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// WorkspaceMemberDao 工作空间成员数据访问层
type WorkspaceMemberDao struct {
	db *gorm.DB
}

func NewWorkspaceMemberDao(db *gorm.DB) *WorkspaceMemberDao {
	return &WorkspaceMemberDao{db: db}
}

// Create 创建成员记录
func (d *WorkspaceMemberDao) Create(ctx context.Context, member *entity.WorkspaceMember) error {
	return d.db.WithContext(ctx).Create(member).Error
}

// Delete 删除成员记录
func (d *WorkspaceMemberDao) Delete(ctx context.Context, workspaceID, customerID uint) error {
	return d.db.WithContext(ctx).
		Where("workspace_id = ? AND customer_id = ?", workspaceID, customerID).
		Delete(&entity.WorkspaceMember{}).Error
}

// FindByWorkspaceAndCustomer 查询指定工作空间的指定成员
func (d *WorkspaceMemberDao) FindByWorkspaceAndCustomer(ctx context.Context, workspaceID, customerID uint) (*entity.WorkspaceMember, error) {
	var member entity.WorkspaceMember
	err := d.db.WithContext(ctx).
		Where("workspace_id = ? AND customer_id = ?", workspaceID, customerID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// ListByWorkspace 查询工作空间的所有成员（带客户信息）
func (d *WorkspaceMemberDao) ListByWorkspace(ctx context.Context, workspaceID uint) ([]entity.WorkspaceMember, error) {
	var members []entity.WorkspaceMember
	err := d.db.WithContext(ctx).
		Preload("Customer").
		Where("workspace_id = ?", workspaceID).
		Order("joined_at ASC").
		Find(&members).Error
	return members, err
}

// ListWorkspaceIDsByCustomer 查询客户加入的所有工作空间 ID
func (d *WorkspaceMemberDao) ListWorkspaceIDsByCustomer(ctx context.Context, customerID uint) ([]uint, error) {
	var ids []uint
	err := d.db.WithContext(ctx).
		Model(&entity.WorkspaceMember{}).
		Where("customer_id = ? AND status = ?", customerID, "active").
		Pluck("workspace_id", &ids).Error
	return ids, err
}
