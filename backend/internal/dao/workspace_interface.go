package dao

import "github.com/YoungBoyGod/remotegpu/internal/model/entity"

// WorkspaceDaoInterface 工作空间DAO接口
type WorkspaceDaoInterface interface {
	Create(workspace *entity.Workspace) error
	GetByID(id uint) (*entity.Workspace, error)
	GetByUUID(uuid string) (*entity.Workspace, error)
	Update(workspace *entity.Workspace) error
	Delete(id uint) error
	GetByOwnerID(ownerID uint) ([]*entity.Workspace, error)
	List(page, pageSize int) ([]*entity.Workspace, int64, error)
	GetByStatus(status string) ([]*entity.Workspace, error)
}

// WorkspaceMemberDaoInterface 工作空间成员DAO接口
type WorkspaceMemberDaoInterface interface {
	Create(member *entity.WorkspaceMember) error
	GetByID(id uint) (*entity.WorkspaceMember, error)
	Update(member *entity.WorkspaceMember) error
	Delete(id uint) error
	GetByWorkspaceID(workspaceID uint) ([]*entity.WorkspaceMember, error)
	GetByUserID(userID uint) ([]*entity.WorkspaceMember, error)
	GetByWorkspaceAndUser(workspaceID, userID uint) (*entity.WorkspaceMember, error)
}
