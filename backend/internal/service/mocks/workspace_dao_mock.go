package mocks

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/mock"
)

// MockWorkspaceDao 是 WorkspaceDaoInterface 的 mock 实现
type MockWorkspaceDao struct {
	mock.Mock
}

// 确保 MockWorkspaceDao 实现了 WorkspaceDaoInterface 接口
var _ interface {
	Create(workspace *entity.Workspace) error
	GetByID(id uint) (*entity.Workspace, error)
	GetByUUID(uuid string) (*entity.Workspace, error)
	Update(workspace *entity.Workspace) error
	Delete(id uint) error
	GetByOwnerID(ownerID uint) ([]*entity.Workspace, error)
	List(page, pageSize int) ([]*entity.Workspace, int64, error)
	GetByStatus(status string) ([]*entity.Workspace, error)
} = (*MockWorkspaceDao)(nil)

// Create mocks the Create method
func (m *MockWorkspaceDao) Create(workspace *entity.Workspace) error {
	args := m.Called(workspace)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockWorkspaceDao) GetByID(id uint) (*entity.Workspace, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Workspace), args.Error(1)
}

// GetByUUID mocks the GetByUUID method
func (m *MockWorkspaceDao) GetByUUID(uuid string) (*entity.Workspace, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Workspace), args.Error(1)
}

// Update mocks the Update method
func (m *MockWorkspaceDao) Update(workspace *entity.Workspace) error {
	args := m.Called(workspace)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockWorkspaceDao) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetByOwnerID mocks the GetByOwnerID method
func (m *MockWorkspaceDao) GetByOwnerID(ownerID uint) ([]*entity.Workspace, error) {
	args := m.Called(ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Workspace), args.Error(1)
}

// List mocks the List method
func (m *MockWorkspaceDao) List(page, pageSize int) ([]*entity.Workspace, int64, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Workspace), args.Get(1).(int64), args.Error(2)
}

// GetByStatus mocks the GetByStatus method
func (m *MockWorkspaceDao) GetByStatus(status string) ([]*entity.Workspace, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Workspace), args.Error(1)
}
