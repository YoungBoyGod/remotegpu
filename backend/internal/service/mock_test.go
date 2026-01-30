package service

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/mock"
)

// MockGPUDao 模拟 GPUDao，实现 GPUDaoInterface
type MockGPUDao struct {
	mock.Mock
}

func (m *MockGPUDao) Create(gpu *entity.GPU) error {
	args := m.Called(gpu)
	return args.Error(0)
}

func (m *MockGPUDao) GetByID(id uint) (*entity.GPU, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.GPU), args.Error(1)
}

func (m *MockGPUDao) GetByHostID(hostID string) ([]*entity.GPU, error) {
	args := m.Called(hostID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.GPU), args.Error(1)
}

func (m *MockGPUDao) Update(gpu *entity.GPU) error {
	args := m.Called(gpu)
	return args.Error(0)
}

func (m *MockGPUDao) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockGPUDao) DeleteByHostID(hostID string) error {
	args := m.Called(hostID)
	return args.Error(0)
}

func (m *MockGPUDao) UpdateStatus(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockGPUDao) List(page, pageSize int) ([]*entity.GPU, int64, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.GPU), args.Get(1).(int64), args.Error(2)
}

func (m *MockGPUDao) GetByStatus(status string) ([]*entity.GPU, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.GPU), args.Error(1)
}

func (m *MockGPUDao) Allocate(id uint, allocatedTo string) error {
	args := m.Called(id, allocatedTo)
	return args.Error(0)
}

func (m *MockGPUDao) Release(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockResourceQuotaDao 模拟 ResourceQuotaDao
type MockResourceQuotaDao struct {
	mock.Mock
}

func (m *MockResourceQuotaDao) Create(quota *entity.ResourceQuota) error {
	args := m.Called(quota)
	return args.Error(0)
}

func (m *MockResourceQuotaDao) GetByID(id uint) (*entity.ResourceQuota, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ResourceQuota), args.Error(1)
}

func (m *MockResourceQuotaDao) GetByUserID(userID uint) (*entity.ResourceQuota, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ResourceQuota), args.Error(1)
}

func (m *MockResourceQuotaDao) GetByWorkspaceID(workspaceID uint) (*entity.ResourceQuota, error) {
	args := m.Called(workspaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ResourceQuota), args.Error(1)
}

func (m *MockResourceQuotaDao) GetByUserAndWorkspace(userID, workspaceID uint) (*entity.ResourceQuota, error) {
	args := m.Called(userID, workspaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ResourceQuota), args.Error(1)
}

func (m *MockResourceQuotaDao) Update(quota *entity.ResourceQuota) error {
	args := m.Called(quota)
	return args.Error(0)
}

func (m *MockResourceQuotaDao) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// 确保MockResourceQuotaDao实现了ResourceQuotaDaoInterface接口
var _ ResourceQuotaDaoInterface = (*MockResourceQuotaDao)(nil)
