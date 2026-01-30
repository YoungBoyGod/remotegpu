package service

import (
	"errors"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)


// TestHostService_Create_Success 测试创建主机成功
func TestHostService_Create_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)

	host := &entity.Host{
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       16,
		TotalMemory:    32000,
	}

	mockHostDao.On("Create", mock.AnythingOfType("*entity.Host")).Return(nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Create(host)

	assert.NoError(t, err)
	assert.NotEmpty(t, host.ID)
	assert.Contains(t, host.ID, "host-")
	assert.False(t, host.RegisteredAt.IsZero())
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Create_WithExistingID 测试创建主机时已有ID
func TestHostService_Create_WithExistingID(t *testing.T) {
	mockHostDao := new(MockHostDao)

	host := &entity.Host{
		ID:             "host-existing",
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       16,
		TotalMemory:    32000,
	}

	mockHostDao.On("Create", host).Return(nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Create(host)

	assert.NoError(t, err)
	assert.Equal(t, "host-existing", host.ID)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Create_DatabaseError 测试创建主机时数据库错误
func TestHostService_Create_DatabaseError(t *testing.T) {
	mockHostDao := new(MockHostDao)

	host := &entity.Host{
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       16,
		TotalMemory:    32000,
	}

	mockHostDao.On("Create", mock.AnythingOfType("*entity.Host")).Return(errors.New("database error"))

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Create(host)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockHostDao.AssertExpectations(t)
}

// TestHostService_GetByID_Success 测试获取主机成功
func TestHostService_GetByID_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)

	expectedHost := &entity.Host{
		ID:        "host-123",
		Name:      "Test Host",
		IPAddress: "192.168.1.100",
		Status:    "active",
	}

	mockHostDao.On("GetByID", "host-123").Return(expectedHost, nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	host, err := service.GetByID("host-123")

	assert.NoError(t, err)
	assert.NotNil(t, host)
	assert.Equal(t, "host-123", host.ID)
	assert.Equal(t, "Test Host", host.Name)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_GetByID_NotFound 测试获取不存在的主机
func TestHostService_GetByID_NotFound(t *testing.T) {
	mockHostDao := new(MockHostDao)

	mockHostDao.On("GetByID", "host-999").Return(nil, gorm.ErrRecordNotFound)

	service := &HostService{
		hostDao: mockHostDao,
	}

	host, err := service.GetByID("host-999")

	assert.Error(t, err)
	assert.Nil(t, host)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Update_Success 测试更新主机成功
func TestHostService_Update_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)

	host := &entity.Host{
		ID:        "host-123",
		Name:      "Updated Host",
		IPAddress: "192.168.1.100",
		Status:    "active",
	}

	mockHostDao.On("Update", host).Return(nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Update(host)

	assert.NoError(t, err)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Update_DatabaseError 测试更新主机时数据库错误
func TestHostService_Update_DatabaseError(t *testing.T) {
	mockHostDao := new(MockHostDao)

	host := &entity.Host{
		ID:        "host-123",
		Name:      "Updated Host",
		IPAddress: "192.168.1.100",
	}

	mockHostDao.On("Update", host).Return(errors.New("database error"))

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Update(host)

	assert.Error(t, err)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Delete_Success 测试删除主机成功
func TestHostService_Delete_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("DeleteByHostID", "host-123").Return(nil)
	mockHostDao.On("Delete", "host-123").Return(nil)

	service := &HostService{
		hostDao: mockHostDao,
		gpuDao:  mockGPUDao,
	}

	err := service.Delete("host-123")

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Delete_GPUDeleteError 测试删除主机时GPU删除失败
func TestHostService_Delete_GPUDeleteError(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("DeleteByHostID", "host-123").Return(errors.New("gpu delete error"))

	service := &HostService{
		hostDao: mockHostDao,
		gpuDao:  mockGPUDao,
	}

	err := service.Delete("host-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gpu delete error")
	mockGPUDao.AssertExpectations(t)
}

// TestHostService_Delete_HostDeleteError 测试删除主机失败
func TestHostService_Delete_HostDeleteError(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("DeleteByHostID", "host-123").Return(nil)
	mockHostDao.On("Delete", "host-123").Return(errors.New("host delete error"))

	service := &HostService{
		hostDao: mockHostDao,
		gpuDao:  mockGPUDao,
	}

	err := service.Delete("host-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "host delete error")
	mockGPUDao.AssertExpectations(t)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_List_Success 测试获取主机列表成功
func TestHostService_List_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)

	expectedHosts := []*entity.Host{
		{ID: "host-1", Name: "Host 1", Status: "active"},
		{ID: "host-2", Name: "Host 2", Status: "active"},
	}

	mockHostDao.On("List", 1, 10).Return(expectedHosts, int64(2), nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	hosts, total, err := service.List(1, 10)

	assert.NoError(t, err)
	assert.Len(t, hosts, 2)
	assert.Equal(t, int64(2), total)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_List_DatabaseError 测试获取主机列表时数据库错误
func TestHostService_List_DatabaseError(t *testing.T) {
	mockHostDao := new(MockHostDao)

	mockHostDao.On("List", 1, 10).Return(nil, int64(0), errors.New("database error"))

	service := &HostService{
		hostDao: mockHostDao,
	}

	hosts, total, err := service.List(1, 10)

	assert.Error(t, err)
	assert.Nil(t, hosts)
	assert.Equal(t, int64(0), total)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_UpdateStatus_Success 测试更新主机状态成功
func TestHostService_UpdateStatus_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)

	mockHostDao.On("UpdateStatus", "host-123", "active").Return(nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.UpdateStatus("host-123", "active")

	assert.NoError(t, err)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_UpdateStatus_DatabaseError 测试更新主机状态时数据库错误
func TestHostService_UpdateStatus_DatabaseError(t *testing.T) {
	mockHostDao := new(MockHostDao)

	mockHostDao.On("UpdateStatus", "host-123", "active").Return(errors.New("database error"))

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.UpdateStatus("host-123", "active")

	assert.Error(t, err)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Heartbeat_Success 测试心跳更新成功
func TestHostService_Heartbeat_Success(t *testing.T) {
	mockHostDao := new(MockHostDao)

	existingHost := &entity.Host{
		ID:     "host-123",
		Name:   "Test Host",
		Status: "active",
	}

	mockHostDao.On("GetByID", "host-123").Return(existingHost, nil)
	mockHostDao.On("UpdateHeartbeat", "host-123").Return(nil)

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Heartbeat("host-123")

	assert.NoError(t, err)
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Heartbeat_HostNotFound 测试心跳更新时主机不存在
func TestHostService_Heartbeat_HostNotFound(t *testing.T) {
	mockHostDao := new(MockHostDao)

	mockHostDao.On("GetByID", "host-999").Return(nil, gorm.ErrRecordNotFound)

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Heartbeat("host-999")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "主机不存在")
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Heartbeat_GetByIDError 测试心跳更新时查询主机失败
func TestHostService_Heartbeat_GetByIDError(t *testing.T) {
	mockHostDao := new(MockHostDao)

	mockHostDao.On("GetByID", "host-123").Return(nil, errors.New("database error"))

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Heartbeat("host-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockHostDao.AssertExpectations(t)
}

// TestHostService_Heartbeat_UpdateError 测试心跳更新失败
func TestHostService_Heartbeat_UpdateError(t *testing.T) {
	mockHostDao := new(MockHostDao)

	existingHost := &entity.Host{
		ID:     "host-123",
		Name:   "Test Host",
		Status: "active",
	}

	mockHostDao.On("GetByID", "host-123").Return(existingHost, nil)
	mockHostDao.On("UpdateHeartbeat", "host-123").Return(errors.New("update error"))

	service := &HostService{
		hostDao: mockHostDao,
	}

	err := service.Heartbeat("host-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	mockHostDao.AssertExpectations(t)
}
