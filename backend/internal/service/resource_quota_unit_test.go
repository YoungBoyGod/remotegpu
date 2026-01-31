package service

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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


// TestResourceQuotaService_CheckQuota_Success 测试CheckQuota成功场景
func TestResourceQuotaService_CheckQuota_Success(t *testing.T) {
	t.Run("Check quota with sufficient resources", func(t *testing.T) {
		// Setup mock DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		// Mock GetQuota
		quota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		// Mock GetUsedResources query
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 8192, 1, 200, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}

		ok, err := service.CheckQuota(1, nil, request)
		assert.NoError(t, err)
		assert.True(t, ok)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Check quota with exact available resources", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 8, 16384, 2, 500, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}

		ok, err := service.CheckQuota(1, nil, request)
		assert.NoError(t, err)
		assert.True(t, ok)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
// TestResourceQuotaService_CheckQuota_Errors 测试CheckQuota错误场景
func TestResourceQuotaService_CheckQuota_Errors(t *testing.T) {
	t.Run("Quota not found", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		mockDao.On("GetByUserID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}

		ok, err := service.CheckQuota(999, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "未设置资源配额")
		mockDao.AssertExpectations(t)
	})

	t.Run("CPU quota exceeded", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 12, 8192, 1, 200, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		request := &ResourceRequest{
			CPU:     8,
			Memory:  8192,
			GPU:     1,
			Storage: 200,
		}

		ok, err := service.CheckQuota(1, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)

		quotaErr, isQuotaErr := err.(*QuotaExceededError)
		assert.True(t, isQuotaErr)
		assert.Equal(t, "CPU", quotaErr.Resource)
		assert.Equal(t, int64(8), quotaErr.Requested)
		assert.Equal(t, int64(4), quotaErr.Available)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_CheckQuota_MemoryGPUStorageExceeded 测试内存/GPU/存储配额超限
func TestResourceQuotaService_CheckQuota_MemoryGPUStorageExceeded(t *testing.T) {
	t.Run("Memory quota exceeded", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 24576, 1, 200, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		request := &ResourceRequest{
			CPU:     4,
			Memory:  16384,
			GPU:     1,
			Storage: 200,
		}

		ok, err := service.CheckQuota(1, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)

		quotaErr, isQuotaErr := err.(*QuotaExceededError)
		assert.True(t, isQuotaErr)
		assert.Equal(t, "Memory", quotaErr.Resource)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GPU quota exceeded", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 8192, 3, 200, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		request := &ResourceRequest{
			CPU:     4,
			Memory:  8192,
			GPU:     2,
			Storage: 200,
		}

		ok, err := service.CheckQuota(1, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)

		quotaErr, isQuotaErr := err.(*QuotaExceededError)
		assert.True(t, isQuotaErr)
		assert.Equal(t, "GPU", quotaErr.Resource)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Storage quota exceeded", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 8192, 1, 700, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		request := &ResourceRequest{
			CPU:     4,
			Memory:  8192,
			GPU:     1,
			Storage: 500,
		}

		ok, err := service.CheckQuota(1, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)

		quotaErr, isQuotaErr := err.(*QuotaExceededError)
		assert.True(t, isQuotaErr)
		assert.Equal(t, "Storage", quotaErr.Resource)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_GetUsedResources_Success 测试GetUsedResources成功场景
func TestResourceQuotaService_GetUsedResources_Success(t *testing.T) {
	t.Run("Get used resources with multiple environments", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		storage1 := int64(200)
		storage2 := int64(300)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 8192, 1, storage1, "running").
			AddRow(2, 1, 1, 8, 16384, 2, storage2, "creating")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		used, err := service.GetUsedResources(1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, used)
		assert.Equal(t, 12, used.CPU)
		assert.Equal(t, int64(24576), used.Memory)
		assert.Equal(t, 3, used.GPU)
		assert.Equal(t, int64(500), used.Storage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get used resources with no environments", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		used, err := service.GetUsedResources(1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, used)
		assert.Equal(t, 0, used.CPU)
		assert.Equal(t, int64(0), used.Memory)
		assert.Equal(t, 0, used.GPU)
		assert.Equal(t, int64(0), used.Storage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get used resources for specific workspace", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		workspaceID := uint(10)
		storage := int64(200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 10, 4, 8192, 1, storage, "running")

		mock.ExpectQuery("SELECT \\* FROM `environments` WHERE.*customer_id.*status.*workspace_id").
			WithArgs(1, "running", "creating", 10).
			WillReturnRows(rows)

		used, err := service.GetUsedResources(1, &workspaceID)
		assert.NoError(t, err)
		assert.NotNil(t, used)
		assert.Equal(t, 4, used.CPU)
		assert.Equal(t, int64(8192), used.Memory)
		assert.Equal(t, 1, used.GPU)
		assert.Equal(t, int64(200), used.Storage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_GetUsedResources_Errors 测试GetUsedResources错误场景
func TestResourceQuotaService_GetUsedResources_Errors(t *testing.T) {
	t.Run("Database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnError(errors.New("database error"))

		used, err := service.GetUsedResources(1, nil)
		assert.Error(t, err)
		assert.Nil(t, used)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_GetAvailableQuota_Success 测试GetAvailableQuota成功场景
func TestResourceQuotaService_GetAvailableQuota_Success(t *testing.T) {
	t.Run("Get available quota with some resources used", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		storage := int64(200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 8192, 1, storage, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		available, err := service.GetAvailableQuota(1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, available)
		assert.Equal(t, 12, available.CPU)
		assert.Equal(t, int64(24576), available.Memory)
		assert.Equal(t, 3, available.GPU)
		assert.Equal(t, int64(800), available.Storage)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get available quota with no resources used", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		available, err := service.GetAvailableQuota(1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, available)
		assert.Equal(t, 16, available.CPU)
		assert.Equal(t, int64(32768), available.Memory)
		assert.Equal(t, 4, available.GPU)
		assert.Equal(t, int64(1000), available.Storage)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get available quota with over-used resources returns zero", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		storage := int64(1200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 20, 40960, 5, storage, "running")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		available, err := service.GetAvailableQuota(1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, available)
		assert.Equal(t, 0, available.CPU)
		assert.Equal(t, int64(0), available.Memory)
		assert.Equal(t, 0, available.GPU)
		assert.Equal(t, int64(0), available.Storage)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_GetAvailableQuota_Errors 测试GetAvailableQuota错误场景
func TestResourceQuotaService_GetAvailableQuota_Errors(t *testing.T) {
	t.Run("Quota not found", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		mockDao.On("GetByUserID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		available, err := service.GetAvailableQuota(999, nil)
		assert.Error(t, err)
		assert.Nil(t, available)
		mockDao.AssertExpectations(t)
	})

	t.Run("Database error on GetUsedResources", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		quota := &entity.ResourceQuota{
			ID:      1,
			UserID:  1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		mockDao.On("GetByUserID", uint(1)).Return(quota, nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnError(errors.New("database error"))

		available, err := service.GetAvailableQuota(1, nil)
		assert.Error(t, err)
		assert.Nil(t, available)
		mockDao.AssertExpectations(t)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_GetQuotaInTx_Success 测试GetQuotaInTx成功场景
func TestResourceQuotaService_GetQuotaInTx_Success(t *testing.T) {
	t.Run("Get quota in transaction with user level", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		mock.ExpectBegin()
		tx := gormDB.Begin()

		rows := sqlmock.NewRows([]string{"id", "user_id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage"}).
			AddRow(1, 1, 1, nil, 16, 32768, 4, 1000)

		mock.ExpectQuery("SELECT \\* FROM `resource_quotas` WHERE customer_id = \\? AND workspace_id IS NULL.*FOR UPDATE").
			WithArgs(1, 1).
			WillReturnRows(rows)

		quota, err := service.GetQuotaInTx(tx, 1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, quota)
		assert.Equal(t, uint(1), quota.ID)
		assert.Equal(t, uint(1), quota.UserID)
		assert.Equal(t, 16, quota.CPU)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get quota in transaction with workspace level", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		mock.ExpectBegin()
		tx := gormDB.Begin()

		workspaceID := uint(10)
		rows := sqlmock.NewRows([]string{"id", "user_id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage"}).
			AddRow(2, 1, 1, 10, 8, 16384, 2, 500)

		mock.ExpectQuery("SELECT \\* FROM `resource_quotas` WHERE customer_id = \\? AND workspace_id = \\?.*FOR UPDATE").
			WithArgs(1, 10, 1).
			WillReturnRows(rows)

		quota, err := service.GetQuotaInTx(tx, 1, &workspaceID)
		assert.NoError(t, err)
		assert.NotNil(t, quota)
		assert.Equal(t, uint(2), quota.ID)
		assert.Equal(t, uint(1), quota.UserID)
		assert.Equal(t, 8, quota.CPU)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_GetQuotaInTx_Errors 测试GetQuotaInTx错误场景
func TestResourceQuotaService_GetQuotaInTx_Errors(t *testing.T) {
	t.Run("Quota not found in transaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		mock.ExpectBegin()
		tx := gormDB.Begin()

		mock.ExpectQuery("SELECT \\* FROM `resource_quotas` WHERE customer_id = \\? AND workspace_id IS NULL.*FOR UPDATE").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		quota, err := service.GetQuotaInTx(tx, 999, nil)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaService_UpdateQuota_WithUsedResources 测试UpdateQuota与已使用资源的场景
func TestResourceQuotaService_UpdateQuota_WithUsedResources(t *testing.T) {
	t.Run("Update quota successfully when new quota is above used resources", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

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
			ID:      1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)

		// Mock GetUsedResources query
		storage := int64(200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 4, 8192, 1, storage, "running")

		sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		mockDao.On("Update", mock.MatchedBy(func(q *entity.ResourceQuota) bool {
			return q.ID == 1 && q.CPU == 16 && q.Memory == 32768 && q.GPU == 4 && q.Storage == 1000
		})).Return(nil)

		err = service.UpdateQuota(newQuota)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Update quota fails when CPU below used resources", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		existingQuota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		newQuota := &entity.ResourceQuota{
			ID:      1,
			CPU:     4,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)

		storage := int64(200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 8, 16384, 2, storage, "running")

		sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		err = service.UpdateQuota(newQuota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPU配额不能小于已使用量")
		mockDao.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Update quota fails when Memory below used resources", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		existingQuota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		newQuota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  8192,
			GPU:     4,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)

		storage := int64(200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 8, 16384, 2, storage, "running")

		sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		err = service.UpdateQuota(newQuota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "内存配额不能小于已使用量")
		mockDao.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Update quota fails when GPU below used resources", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		existingQuota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		newQuota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  32768,
			GPU:     1,
			Storage: 1000,
		}

		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)

		storage := int64(200)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 8, 16384, 2, storage, "running")

		sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		err = service.UpdateQuota(newQuota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GPU配额不能小于已使用量")
		mockDao.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Update quota fails when Storage below used resources", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockDao := new(MockResourceQuotaDao)
		service := NewResourceQuotaServiceWithDeps(mockDao, gormDB)

		existingQuota := &entity.ResourceQuota{
			ID:          1,
			UserID:      1,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}

		newQuota := &entity.ResourceQuota{
			ID:      1,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 100,
		}

		mockDao.On("GetByID", uint(1)).Return(existingQuota, nil)

		storage := int64(500)
		rows := sqlmock.NewRows([]string{"id", "customer_id", "workspace_id", "cpu", "memory", "gpu", "storage", "status"}).
			AddRow(1, 1, 1, 8, 16384, 2, storage, "running")

		sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `environments` WHERE customer_id = ? AND status IN (?,?)")).
			WithArgs(1, "running", "creating").
			WillReturnRows(rows)

		err = service.UpdateQuota(newQuota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "存储配额不能小于已使用量")
		mockDao.AssertExpectations(t)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}
