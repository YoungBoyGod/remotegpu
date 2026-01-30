package service

import (
	"errors"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)


// TestGPUService_Create_Success 测试创建GPU成功
func TestGPUService_Create_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	gpu := &entity.GPU{
		HostID:      "host-123",
		GPUIndex:    0,
		Name:        "RTX 4090",
		Brand:       "NVIDIA",
		MemoryTotal: 24000,
	}

	mockGPUDao.On("Create", gpu).Return(nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Create(gpu)

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Create_DatabaseError 测试创建GPU时数据库错误
func TestGPUService_Create_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	gpu := &entity.GPU{
		HostID:      "host-123",
		GPUIndex:    0,
		Name:        "RTX 4090",
		Brand:       "NVIDIA",
		MemoryTotal: 24000,
	}

	mockGPUDao.On("Create", gpu).Return(errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Create(gpu)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_GetByID_Success 测试获取GPU成功
func TestGPUService_GetByID_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	expectedGPU := &entity.GPU{
		ID:          1,
		HostID:      "host-123",
		GPUIndex:    0,
		Name:        "RTX 4090",
		Brand:       "NVIDIA",
		MemoryTotal: 24000,
		Status:      "available",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(expectedGPU, nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpu, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, gpu)
	assert.Equal(t, uint(1), gpu.ID)
	assert.Equal(t, "RTX 4090", gpu.Name)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_GetByID_NotFound 测试获取不存在的GPU
func TestGPUService_GetByID_NotFound(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpu, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, gpu)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_GetByHostID_Success 测试根据主机ID获取GPU列表成功
func TestGPUService_GetByHostID_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	expectedGPUs := []*entity.GPU{
		{ID: 1, HostID: "host-123", GPUIndex: 0, Name: "RTX 4090"},
		{ID: 2, HostID: "host-123", GPUIndex: 1, Name: "RTX 4090"},
	}

	mockGPUDao.On("GetByHostID", "host-123").Return(expectedGPUs, nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpus, err := service.GetByHostID("host-123")

	assert.NoError(t, err)
	assert.Len(t, gpus, 2)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_GetByHostID_DatabaseError 测试根据主机ID获取GPU列表时数据库错误
func TestGPUService_GetByHostID_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("GetByHostID", "host-123").Return(nil, errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpus, err := service.GetByHostID("host-123")

	assert.Error(t, err)
	assert.Nil(t, gpus)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Update_Success 测试更新GPU成功
func TestGPUService_Update_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	gpu := &entity.GPU{
		ID:          1,
		HostID:      "host-123",
		GPUIndex:    0,
		Name:        "RTX 4090 Updated",
		Brand:       "NVIDIA",
		MemoryTotal: 24000,
	}

	mockGPUDao.On("Update", gpu).Return(nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Update(gpu)

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Update_DatabaseError 测试更新GPU时数据库错误
func TestGPUService_Update_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	gpu := &entity.GPU{
		ID:          1,
		HostID:      "host-123",
		GPUIndex:    0,
		Name:        "RTX 4090",
		Brand:       "NVIDIA",
		MemoryTotal: 24000,
	}

	mockGPUDao.On("Update", gpu).Return(errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Update(gpu)

	assert.Error(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Delete_Success 测试删除GPU成功
func TestGPUService_Delete_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("Delete", uint(1)).Return(nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Delete(1)

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Delete_DatabaseError 测试删除GPU时数据库错误
func TestGPUService_Delete_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("Delete", uint(1)).Return(errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Delete(1)

	assert.Error(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_UpdateStatus_Success 测试更新GPU状态成功
func TestGPUService_UpdateStatus_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("UpdateStatus", uint(1), "maintenance").Return(nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.UpdateStatus(1, "maintenance")

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_UpdateStatus_DatabaseError 测试更新GPU状态时数据库错误
func TestGPUService_UpdateStatus_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("UpdateStatus", uint(1), "maintenance").Return(errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.UpdateStatus(1, "maintenance")

	assert.Error(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_List_Success 测试获取GPU列表成功
func TestGPUService_List_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	expectedGPUs := []*entity.GPU{
		{ID: 1, HostID: "host-1", Name: "RTX 4090", Status: "available"},
		{ID: 2, HostID: "host-2", Name: "RTX 3090", Status: "available"},
	}

	mockGPUDao.On("List", 1, 10).Return(expectedGPUs, int64(2), nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpus, total, err := service.List(1, 10)

	assert.NoError(t, err)
	assert.Len(t, gpus, 2)
	assert.Equal(t, int64(2), total)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_List_DatabaseError 测试获取GPU列表时数据库错误
func TestGPUService_List_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("List", 1, 10).Return(nil, int64(0), errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpus, total, err := service.List(1, 10)

	assert.Error(t, err)
	assert.Nil(t, gpus)
	assert.Equal(t, int64(0), total)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_GetByStatus_Success 测试根据状态获取GPU列表成功
func TestGPUService_GetByStatus_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	expectedGPUs := []*entity.GPU{
		{ID: 1, HostID: "host-1", Name: "RTX 4090", Status: "available"},
		{ID: 2, HostID: "host-2", Name: "RTX 3090", Status: "available"},
	}

	mockGPUDao.On("GetByStatus", "available").Return(expectedGPUs, nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpus, err := service.GetByStatus("available")

	assert.NoError(t, err)
	assert.Len(t, gpus, 2)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_GetByStatus_DatabaseError 测试根据状态获取GPU列表时数据库错误
func TestGPUService_GetByStatus_DatabaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("GetByStatus", "available").Return(nil, errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	gpus, err := service.GetByStatus("available")

	assert.Error(t, err)
	assert.Nil(t, gpus)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Allocate_Success 测试分配GPU成功
func TestGPUService_Allocate_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	availableGPU := &entity.GPU{
		ID:     1,
		HostID: "host-123",
		Name:   "RTX 4090",
		Status: "available",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(availableGPU, nil)
	mockGPUDao.On("Allocate", uint(1), "env-123").Return(nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Allocate(1, "env-123")

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Allocate_GPUNotAvailable 测试分配不可用的GPU
func TestGPUService_Allocate_GPUNotAvailable(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	allocatedGPU := &entity.GPU{
		ID:     1,
		HostID: "host-123",
		Name:   "RTX 4090",
		Status: "allocated",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(allocatedGPU, nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Allocate(1, "env-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GPU不可用")
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Allocate_GetByIDError 测试分配GPU时查询失败
func TestGPUService_Allocate_GetByIDError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("GetByID", uint(1)).Return(nil, errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Allocate(1, "env-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Allocate_AllocateError 测试分配GPU时分配操作失败
func TestGPUService_Allocate_AllocateError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	availableGPU := &entity.GPU{
		ID:     1,
		HostID: "host-123",
		Name:   "RTX 4090",
		Status: "available",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(availableGPU, nil)
	mockGPUDao.On("Allocate", uint(1), "env-123").Return(errors.New("allocate error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Allocate(1, "env-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "allocate error")
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Release_Success 测试释放GPU成功
func TestGPUService_Release_Success(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	allocatedGPU := &entity.GPU{
		ID:     1,
		HostID: "host-123",
		Name:   "RTX 4090",
		Status: "allocated",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(allocatedGPU, nil)
	mockGPUDao.On("Release", uint(1)).Return(nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Release(1)

	assert.NoError(t, err)
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Release_GPUNotAllocated 测试释放未分配的GPU
func TestGPUService_Release_GPUNotAllocated(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	availableGPU := &entity.GPU{
		ID:     1,
		HostID: "host-123",
		Name:   "RTX 4090",
		Status: "available",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(availableGPU, nil)

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Release(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GPU未分配")
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Release_GetByIDError 测试释放GPU时查询失败
func TestGPUService_Release_GetByIDError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	mockGPUDao.On("GetByID", uint(1)).Return(nil, errors.New("database error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Release(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockGPUDao.AssertExpectations(t)
}

// TestGPUService_Release_ReleaseError 测试释放GPU时释放操作失败
func TestGPUService_Release_ReleaseError(t *testing.T) {
	mockGPUDao := new(MockGPUDao)

	allocatedGPU := &entity.GPU{
		ID:     1,
		HostID: "host-123",
		Name:   "RTX 4090",
		Status: "allocated",
	}

	mockGPUDao.On("GetByID", uint(1)).Return(allocatedGPU, nil)
	mockGPUDao.On("Release", uint(1)).Return(errors.New("release error"))

	service := &GPUService{
		gpuDao: mockGPUDao,
	}

	err := service.Release(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "release error")
	mockGPUDao.AssertExpectations(t)
}
