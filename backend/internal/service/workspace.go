package service

import (
	"errors"
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkspaceService 工作空间服务
type WorkspaceService struct {
	workspaceDao *dao.WorkspaceDao
	memberDao    *dao.WorkspaceMemberDao
}

// NewWorkspaceService 创建工作空间服务实例
func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceDao: dao.NewWorkspaceDao(),
		memberDao:    dao.NewWorkspaceMemberDao(),
	}
}

// CreateWorkspace 创建工作空间
func (s *WorkspaceService) CreateWorkspace(workspace *entity.Workspace) error {
	// 生成 UUID
	if workspace.UUID == uuid.Nil {
		workspace.UUID = uuid.New()
	}

	// 设置默认值
	if workspace.Type == "" {
		workspace.Type = "personal"
	}
	if workspace.Status == "" {
		workspace.Status = "active"
	}
	if workspace.MemberCount == 0 {
		workspace.MemberCount = 1
	}

	return s.workspaceDao.Create(workspace)
}

// GetWorkspace 根据ID获取工作空间
func (s *WorkspaceService) GetWorkspace(id uint) (*entity.Workspace, error) {
	return s.workspaceDao.GetByID(id)
}

// UpdateWorkspace 更新工作空间
func (s *WorkspaceService) UpdateWorkspace(workspace *entity.Workspace) error {
	// 检查工作空间是否存在
	existing, err := s.workspaceDao.GetByID(workspace.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("工作空间不存在")
		}
		return err
	}

	// 不允许修改所有者
	if workspace.OwnerID != existing.OwnerID {
		return fmt.Errorf("不允许修改工作空间所有者")
	}

	return s.workspaceDao.Update(workspace)
}

// DeleteWorkspace 删除工作空间
func (s *WorkspaceService) DeleteWorkspace(id uint) error {
	// 先删除所有成员
	members, err := s.memberDao.GetByWorkspaceID(id)
	if err != nil {
		return err
	}

	for _, member := range members {
		if err := s.memberDao.Delete(member.ID); err != nil {
			return err
		}
	}

	return s.workspaceDao.Delete(id)
}

// ListWorkspaces 获取工作空间列表
func (s *WorkspaceService) ListWorkspaces(ownerID uint, page, pageSize int) ([]*entity.Workspace, int64, error) {
	if ownerID > 0 {
		workspaces, err := s.workspaceDao.GetByOwnerID(ownerID)
		if err != nil {
			return nil, 0, err
		}
		return workspaces, int64(len(workspaces)), nil
	}
	return s.workspaceDao.List(page, pageSize)
}

// AddMember 添加工作空间成员
func (s *WorkspaceService) AddMember(workspaceID, customerID uint, role string) error {
	// 检查工作空间是否存在
	workspace, err := s.workspaceDao.GetByID(workspaceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("工作空间不存在")
		}
		return err
	}

	// 检查成员是否已存在
	existing, err := s.memberDao.GetByWorkspaceAndCustomer(workspaceID, customerID)
	if err == nil && existing != nil {
		return fmt.Errorf("成员已存在")
	}

	// 设置默认角色
	if role == "" {
		role = "member"
	}

	// 创建成员
	member := &entity.WorkspaceMember{
		WorkspaceID: workspaceID,
		CustomerID:  customerID,
		Role:        role,
		Status:      "active",
	}

	if err := s.memberDao.Create(member); err != nil {
		return err
	}

	// 更新成员数量
	workspace.MemberCount++
	return s.workspaceDao.Update(workspace)
}

// RemoveMember 移除工作空间成员
func (s *WorkspaceService) RemoveMember(workspaceID, customerID uint) error {
	// 检查工作空间是否存在
	workspace, err := s.workspaceDao.GetByID(workspaceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("工作空间不存在")
		}
		return err
	}

	// 检查成员是否存在
	member, err := s.memberDao.GetByWorkspaceAndCustomer(workspaceID, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("成员不存在")
		}
		return err
	}

	// 不允许移除所有者
	if workspace.OwnerID == customerID {
		return fmt.Errorf("不允许移除工作空间所有者")
	}

	// 删除成员
	if err := s.memberDao.Delete(member.ID); err != nil {
		return err
	}

	// 更新成员数量
	if workspace.MemberCount > 0 {
		workspace.MemberCount--
	}
	return s.workspaceDao.Update(workspace)
}

// ListMembers 获取工作空间成员列表
func (s *WorkspaceService) ListMembers(workspaceID uint) ([]*entity.WorkspaceMember, error) {
	return s.memberDao.GetByWorkspaceID(workspaceID)
}

// CheckPermission 检查用户是否有工作空间权限
func (s *WorkspaceService) CheckPermission(workspaceID, customerID uint) (bool, error) {
	// 检查工作空间是否存在
	workspace, err := s.workspaceDao.GetByID(workspaceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("工作空间不存在")
		}
		return false, err
	}

	// 检查是否是所有者
	if workspace.OwnerID == customerID {
		return true, nil
	}

	// 检查是否是成员
	member, err := s.memberDao.GetByWorkspaceAndCustomer(workspaceID, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	// 检查成员状态
	return member.Status == "active", nil
}
