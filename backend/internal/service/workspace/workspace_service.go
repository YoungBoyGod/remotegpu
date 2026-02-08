package workspace

import (
	"context"
	"errors"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// WorkspaceService 工作空间业务逻辑层
type WorkspaceService struct {
	db        *gorm.DB
	wsDao     *dao.WorkspaceDao
	memberDao *dao.WorkspaceMemberDao
}

func NewWorkspaceService(db *gorm.DB) *WorkspaceService {
	return &WorkspaceService{
		db:        db,
		wsDao:     dao.NewWorkspaceDao(db),
		memberDao: dao.NewWorkspaceMemberDao(db),
	}
}

// Create 创建工作空间，同时插入 owner 成员记录
func (s *WorkspaceService) Create(ctx context.Context, ownerID uint, name, description string) (*entity.Workspace, error) {
	ws := &entity.Workspace{
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
		Type:        "personal",
		Status:      "active",
		MemberCount: 1,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ws).Error; err != nil {
			return err
		}
		member := &entity.WorkspaceMember{
			WorkspaceID: ws.ID,
			CustomerID:  ownerID,
			Role:        "owner",
			Status:      "active",
			JoinedAt:    time.Now(),
		}
		return tx.Create(member).Error
	})
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// List 获取用户可见的工作空间列表（拥有的 + 作为成员加入的）
func (s *WorkspaceService) List(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Workspace, int64, error) {
	wsIDs, err := s.memberDao.ListWorkspaceIDsByCustomer(ctx, customerID)
	if err != nil {
		return nil, 0, err
	}
	if len(wsIDs) == 0 {
		return []entity.Workspace{}, 0, nil
	}

	var total int64
	db := s.db.WithContext(ctx).Model(&entity.Workspace{}).Where("id IN ?", wsIDs)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var workspaces []entity.Workspace
	err = db.Preload("Owner").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&workspaces).Error
	return workspaces, total, err
}

// GetByID 获取工作空间详情
func (s *WorkspaceService) GetByID(ctx context.Context, id uint) (*entity.Workspace, error) {
	return s.wsDao.FindByID(ctx, id)
}

// Update 更新工作空间（仅 owner/admin 可操作）
func (s *WorkspaceService) Update(ctx context.Context, wsID, customerID uint, fields map[string]interface{}) error {
	if err := s.requireRole(ctx, wsID, customerID, "owner", "admin"); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Model(&entity.Workspace{}).Where("id = ?", wsID).Updates(fields).Error
}

// Delete 删除工作空间（仅 owner 可操作）
func (s *WorkspaceService) Delete(ctx context.Context, wsID, customerID uint) error {
	if err := s.requireRole(ctx, wsID, customerID, "owner"); err != nil {
		return err
	}

	// 检查是否有运行中的环境
	var count int64
	s.db.WithContext(ctx).Model(&entity.Environment{}).
		Where("workspace_id = ? AND status IN ?", wsID, []string{"running", "creating"}).
		Count(&count)
	if count > 0 {
		return errors.New("工作空间下有运行中的环境，请先停止")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("workspace_id = ?", wsID).Delete(&entity.WorkspaceMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Workspace{}, wsID).Error
	})
}

// AddMember 添加成员（仅 owner/admin 可操作）
func (s *WorkspaceService) AddMember(ctx context.Context, wsID, operatorID, targetID uint, role string) error {
	if err := s.requireRole(ctx, wsID, operatorID, "owner", "admin"); err != nil {
		return err
	}

	// admin 不能添加 owner 角色
	operatorMember, _ := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, operatorID)
	if operatorMember != nil && operatorMember.Role == "admin" && role == "owner" {
		return errors.New("admin 不能添加 owner 角色")
	}

	// 检查是否已存在
	existing, _ := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, targetID)
	if existing != nil {
		return errors.New("该用户已是工作空间成员")
	}

	if role == "" {
		role = "member"
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		member := &entity.WorkspaceMember{
			WorkspaceID: wsID,
			CustomerID:  targetID,
			Role:        role,
			Status:      "active",
			JoinedAt:    time.Now(),
		}
		if err := tx.Create(member).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Workspace{}).Where("id = ?", wsID).
			UpdateColumn("member_count", gorm.Expr("member_count + 1")).Error
	})
	return err
}

// RemoveMember 移除成员（仅 owner/admin 可操作）
func (s *WorkspaceService) RemoveMember(ctx context.Context, wsID, operatorID, targetID uint) error {
	if err := s.requireRole(ctx, wsID, operatorID, "owner", "admin"); err != nil {
		return err
	}

	target, err := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, targetID)
	if err != nil {
		return errors.New("成员不存在")
	}
	if target.Role == "owner" {
		return errors.New("不能移除 owner")
	}

	operator, _ := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, operatorID)
	if operator != nil && operator.Role == "admin" && target.Role == "admin" {
		return errors.New("admin 不能移除其他 admin")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("workspace_id = ? AND customer_id = ?", wsID, targetID).
			Delete(&entity.WorkspaceMember{}).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Workspace{}).Where("id = ?", wsID).
			UpdateColumn("member_count", gorm.Expr("member_count - 1")).Error
	})
}

// ListMembers 获取工作空间成员列表
func (s *WorkspaceService) ListMembers(ctx context.Context, wsID, customerID uint) ([]entity.WorkspaceMember, error) {
	// 校验当前用户是否为成员
	_, err := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, customerID)
	if err != nil {
		return nil, errors.New("无权查看该工作空间成员")
	}
	return s.memberDao.ListByWorkspace(ctx, wsID)
}

// IsMember 检查用户是否为工作空间成员
func (s *WorkspaceService) IsMember(ctx context.Context, wsID, customerID uint) bool {
	m, err := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, customerID)
	return err == nil && m != nil
}

// requireRole 校验用户在工作空间中的角色
func (s *WorkspaceService) requireRole(ctx context.Context, wsID, customerID uint, roles ...string) error {
	member, err := s.memberDao.FindByWorkspaceAndCustomer(ctx, wsID, customerID)
	if err != nil {
		return errors.New("无权操作该工作空间")
	}
	for _, r := range roles {
		if member.Role == r {
			return nil
		}
	}
	return errors.New("权限不足")
}
