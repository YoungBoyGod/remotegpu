package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// WorkspaceDao 工作空间数据访问对象
type WorkspaceDao struct {
	db *gorm.DB
}

// NewWorkspaceDao 创建工作空间 DAO
func NewWorkspaceDao() *WorkspaceDao {
	return &WorkspaceDao{
		db: database.GetDB(),
	}
}

// Create 创建工作空间
func (d *WorkspaceDao) Create(workspace *entity.Workspace) error {
	return d.db.Create(workspace).Error
}

// GetByID 根据ID获取工作空间
func (d *WorkspaceDao) GetByID(id uint) (*entity.Workspace, error) {
	var workspace entity.Workspace
	err := d.db.Where("id = ?", id).First(&workspace).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

// GetByUUID 根据UUID获取工作空间
func (d *WorkspaceDao) GetByUUID(uuid string) (*entity.Workspace, error) {
	var workspace entity.Workspace
	err := d.db.Where("uuid = ?", uuid).First(&workspace).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

// Update 更新工作空间
func (d *WorkspaceDao) Update(workspace *entity.Workspace) error {
	return d.db.Save(workspace).Error
}

// Delete 删除工作空间（软删除）
func (d *WorkspaceDao) Delete(id uint) error {
	return d.db.Delete(&entity.Workspace{}, id).Error
}

// GetByOwnerID 根据所有者ID获取工作空间列表
func (d *WorkspaceDao) GetByOwnerID(ownerID uint) ([]*entity.Workspace, error) {
	var workspaces []*entity.Workspace
	err := d.db.Where("owner_id = ?", ownerID).Find(&workspaces).Error
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

// List 获取工作空间列表（分页）
func (d *WorkspaceDao) List(page, pageSize int) ([]*entity.Workspace, int64, error) {
	var workspaces []*entity.Workspace
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&entity.Workspace{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Offset(offset).Limit(pageSize).Find(&workspaces).Error; err != nil {
		return nil, 0, err
	}

	return workspaces, total, nil
}

// GetByStatus 根据状态获取工作空间列表
func (d *WorkspaceDao) GetByStatus(status string) ([]*entity.Workspace, error) {
	var workspaces []*entity.Workspace
	err := d.db.Where("status = ?", status).Find(&workspaces).Error
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

// WorkspaceMemberDao 工作空间成员数据访问对象
type WorkspaceMemberDao struct {
	db *gorm.DB
}

// NewWorkspaceMemberDao 创建工作空间成员 DAO
func NewWorkspaceMemberDao() *WorkspaceMemberDao {
	return &WorkspaceMemberDao{
		db: database.GetDB(),
	}
}

// Create 创建工作空间成员
func (d *WorkspaceMemberDao) Create(member *entity.WorkspaceMember) error {
	return d.db.Create(member).Error
}

// GetByID 根据ID获取工作空间成员
func (d *WorkspaceMemberDao) GetByID(id uint) (*entity.WorkspaceMember, error) {
	var member entity.WorkspaceMember
	err := d.db.Where("id = ?", id).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// Update 更新工作空间成员
func (d *WorkspaceMemberDao) Update(member *entity.WorkspaceMember) error {
	return d.db.Save(member).Error
}

// Delete 删除工作空间成员
func (d *WorkspaceMemberDao) Delete(id uint) error {
	return d.db.Delete(&entity.WorkspaceMember{}, id).Error
}

// GetByWorkspaceID 根据工作空间ID获取成员列表
func (d *WorkspaceMemberDao) GetByWorkspaceID(workspaceID uint) ([]*entity.WorkspaceMember, error) {
	var members []*entity.WorkspaceMember
	err := d.db.Where("workspace_id = ?", workspaceID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetByUserID 根据用户ID获取成员列表
func (d *WorkspaceMemberDao) GetByUserID(userID uint) ([]*entity.WorkspaceMember, error) {
	var members []*entity.WorkspaceMember
	err := d.db.Where("user_id = ?", userID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetByWorkspaceAndUser 根据工作空间ID和用户ID获取成员
func (d *WorkspaceMemberDao) GetByWorkspaceAndUser(workspaceID, userID uint) (*entity.WorkspaceMember, error) {
	var member entity.WorkspaceMember
	err := d.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}
