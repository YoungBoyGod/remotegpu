package service

import (
	"errors"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestResourceQuotaService_SetQuota_Unit 测试SetQuota方法（单元测试）
func TestResourceQuotaService_SetQuota_Unit(t *testing.T) {
	t.Run("Create new user level quota", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		mockDao.On("GetByUserID", uint(1)).Return(nil, gorm.ErrRecordNotFound)
		mockDao.On("Create", quota).Return(nil)

		err := service.SetQuota(quota)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("Update existing user level quota", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		existingQuota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}

		newQuota := &entity.ResourceQuota{
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		mockDao.On("GetByUserID", uint(1)).Return(existingQuota, nil)
		mockDao.On("Update", mock.AnythingOfType("*entity.ResourceQuota")).Return(nil)

		err := service.SetQuota(newQuota)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), newQuota.ID)
		mockDao.AssertExpectations(t)
	})

	t.Run("Create workspace level quota", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		workspaceID := uint(10)
		quota := &entity.ResourceQuota{
			UserID:      1,
			WorkspaceID: &workspaceID,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}

		mockDao.On("GetByUserAndWorkspace", uint(1), uint(10)).Return(nil, gorm.ErrRecordNotFound)
		mockDao.On("Create", quota).Return(nil)

		err := service.SetQuota(quota)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("Negative CPU value", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			UserID:      1,
			WorkspaceID: nil,
			CPU:         -1,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}

		err := service.SetQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
	})

	t.Run("Negative Memory value", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      -1,
			GPU:         2,
			Storage:     500,
		}

		err := service.SetQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
	})

	t.Run("Database error on GetByUserID", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		mockDao.On("GetByUserID", uint(1)).Return(nil, errors.New("database error"))

		err := service.SetQuota(quota)
		assert.Error(t, err)
		mockDao.AssertExpectations(t)
	})
}

// TestResourceQuotaService_GetQuota_Unit 测试GetQuota方法（单元测试）
func TestResourceQuotaService_GetQuota_Unit(t *testing.T) {
	t.Run("Get user level quota successfully", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		expectedQuota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		mockDao.On("GetByUserID", uint(1)).Return(expectedQuota, nil)

		quota, err := service.GetQuota(1, nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuota, quota)
		mockDao.AssertExpectations(t)
	})

	t.Run("Get workspace level quota successfully", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		workspaceID := uint(10)
		expectedQuota := &entity.ResourceQuota{
			ID:          2,
			UserID:      1,
			WorkspaceID: &workspaceID,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}

		mockDao.On("GetByUserAndWorkspace", uint(1), uint(10)).Return(expectedQuota, nil)

		quota, err := service.GetQuota(1, &workspaceID)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuota, quota)
		mockDao.AssertExpectations(t)
	})

	t.Run("User level quota not found", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		mockDao.On("GetByUserID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		quota, err := service.GetQuota(999, nil)
		assert.Error(t, err)
		assert.Nil(t, quota)
		mockDao.AssertExpectations(t)
	})
}

// TestResourceQuotaService_GetQuotaByID_Unit 测试GetQuotaByID方法
func TestResourceQuotaService_GetQuotaByID_Unit(t *testing.T) {
	t.Run("Get quota by ID successfully", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		expectedQuota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(expectedQuota, nil)

		quota, err := service.GetQuotaByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuota, quota)
		mockDao.AssertExpectations(t)
	})

	t.Run("Quota not found by ID", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		mockDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		quota, err := service.GetQuotaByID(999)
		assert.Error(t, err)
		assert.Nil(t, quota)
		mockDao.AssertExpectations(t)
	})
}

// TestResourceQuotaService_DeleteQuota_Unit 测试DeleteQuota方法
func TestResourceQuotaService_DeleteQuota_Unit(t *testing.T) {
	t.Run("Delete quota successfully", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		existingQuota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)
		mockDao.On("Delete", uint(1)).Return(nil)

		err := service.DeleteQuota(1)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("Delete non-existent quota", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		mockDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.DeleteQuota(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额不存在")
		mockDao.AssertExpectations(t)
	})

	t.Run("Delete with database error", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		existingQuota := &entity.ResourceQuota{ID: 1}
		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)
		mockDao.On("Delete", uint(1)).Return(errors.New("database error"))

		err := service.DeleteQuota(1)
		assert.Error(t, err)
		mockDao.AssertExpectations(t)
	})
}

// TestResourceQuotaService_UpdateQuota_Unit 测试UpdateQuota方法
func TestResourceQuotaService_UpdateQuota_Unit(t *testing.T) {
	t.Run("Update with negative CPU", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			ID:      1,
			CPU:     -1,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
	})

	t.Run("Update with negative Memory", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  -1,
			GPU:     4,
			Storage: 1000,
		}

		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
	})

	t.Run("Update with negative GPU", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  32768,
			GPU:     -1,
			Storage: 1000,
		}

		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
	})

	t.Run("Update with negative Storage", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: -1,
		}

		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
	})

	t.Run("Update non-existent quota", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			ID:      999,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额不存在")
		mockDao.AssertExpectations(t)
	})

	t.Run("GetByID returns database error", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(nil, errors.New("database error"))

		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		mockDao.AssertExpectations(t)
	})
}

// TestQuotaExceededError_Unit 测试QuotaExceededError错误信息
func TestQuotaExceededError_Unit(t *testing.T) {
	t.Run("Positive available", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "CPU",
			Requested: 16,
			Available: 8,
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "CPU 配额不足")
		assert.Contains(t, errMsg, "需要 16")
		assert.Contains(t, errMsg, "可用 8")
		assert.NotContains(t, errMsg, "已超额使用")
	})

	t.Run("Negative available", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "Memory",
			Requested: 16384,
			Available: -4096,
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "Memory 配额不足")
		assert.Contains(t, errMsg, "需要 16384")
		assert.Contains(t, errMsg, "可用 0")
		assert.Contains(t, errMsg, "已超额使用 4096")
	})

	t.Run("Zero available", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "GPU",
			Requested: 4,
			Available: 0,
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "GPU 配额不足")
		assert.Contains(t, errMsg, "需要 4")
		assert.Contains(t, errMsg, "可用 0")
		assert.NotContains(t, errMsg, "已超额使用")
	})

	t.Run("Storage quota exceeded", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "Storage",
			Requested: 2000,
			Available: 500,
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "Storage 配额不足")
		assert.Contains(t, errMsg, "需要 2000")
		assert.Contains(t, errMsg, "可用 500")
	})
}

// TestResourceQuotaService_GetQuotaInTx_Unit 测试GetQuotaInTx方法
func TestResourceQuotaService_GetQuotaInTx_Unit(t *testing.T) {
	t.Run("Nil transaction", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		quota, err := service.GetQuotaInTx(nil, 1, nil)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.Contains(t, err.Error(), "事务不能为空")
	})
}

// TestResourceQuotaService_CheckQuotaInTx_Unit 测试CheckQuotaInTx方法
func TestResourceQuotaService_CheckQuotaInTx_Unit(t *testing.T) {
	t.Run("Nil transaction", func(t *testing.T) {
		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDao(mockDao)

		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}

		ok, err := service.CheckQuotaInTx(nil, 1, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "事务不能为空")
	})
}
