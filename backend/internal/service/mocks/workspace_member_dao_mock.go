package mocks

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/mock"
)

// MockWorkspaceMemberDao 是 WorkspaceMemberDaoInterface 的 mock 实现
type MockWorkspaceMemberDao struct {
	mock.Mock
}

// 确保 MockWorkspaceMemberDao 实现了 WorkspaceMemberDaoInterface 接口
var _ interface {
	Create(member *entity.WorkspaceMember) error
	GetByID(id uint) (*entity.WorkspaceMember, error)
	Update(member *entity.WorkspaceMember) error
	Delete(id uint) error
	GetByWorkspaceID(workspaceID uint) ([]*entity.WorkspaceMember, error)
	GetByUserID(userID uint) ([]*entity.WorkspaceMember, error)
	GetByWorkspaceAndUser(workspaceID, userID uint) (*entity.WorkspaceMember, error)
} = (*MockWorkspaceMemberDao)(nil)

// Create mocks the Create method
func (m *MockWorkspaceMemberDao) Create(member *entity.WorkspaceMember) error {
	args := m.Called(member)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockWorkspaceMemberDao) GetByID(id uint) (*entity.WorkspaceMember, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WorkspaceMember), args.Error(1)
}

// Update mocks the Update method
func (m *MockWorkspaceMemberDao) Update(member *entity.WorkspaceMember) error {
	args := m.Called(member)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockWorkspaceMemberDao) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetByWorkspaceID mocks the GetByWorkspaceID method
func (m *MockWorkspaceMemberDao) GetByWorkspaceID(workspaceID uint) ([]*entity.WorkspaceMember, error) {
	args := m.Called(workspaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.WorkspaceMember), args.Error(1)
}

// GetByUserID mocks the GetByUserID method
func (m *MockWorkspaceMemberDao) GetByUserID(userID uint) ([]*entity.WorkspaceMember, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.WorkspaceMember), args.Error(1)
}

// GetByWorkspaceAndUser mocks the GetByWorkspaceAndUser method
func (m *MockWorkspaceMemberDao) GetByWorkspaceAndUser(workspaceID, userID uint) (*entity.WorkspaceMember, error) {
	args := m.Called(workspaceID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WorkspaceMember), args.Error(1)
}
