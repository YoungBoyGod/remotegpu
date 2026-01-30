package service

import (
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupResourceQuotaServiceTest(t *testing.T) {
	cfg := database.Config{
		Host:     "192.168.10.210",
		Port:     5432,
		User:     "remotegpu_user",
		Password: "remotegpu_password",
		DBName:   "remotegpu",
	}
	if err := database.InitDB(cfg); err != nil {
		t.Skipf("跳过测试，无法连接数据库: %v", err)
	}

	db := database.GetDB()
	// 自动迁移表结构
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.ResourceQuota{}); err != nil {
		t.Fatalf("自动迁移表结构失败: %v", err)
	}
}

func TestResourceQuotaService_SetAndGetQuota(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-quota-service-" + uuid.New().String()[:8],
		Email:        "quota-service-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Quota Service User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 测试设置用户级别配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)
	t.Log("设置用户级别配额成功")

	// 测试获取配额
	retrieved, err := service.GetQuota(customer.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 16, retrieved.CPU)
	assert.Equal(t, int64(32768), retrieved.Memory)
	assert.Equal(t, 4, retrieved.GPU)
	assert.Equal(t, int64(1000), retrieved.Storage)
	t.Log("获取配额成功")

	// 清理
	service.DeleteQuota(retrieved.ID)
}

func TestResourceQuotaService_UpdateQuota(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-update-quota-" + uuid.New().String()[:8],
		Email:        "update-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Update Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 设置初始配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     500,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)

	// 获取配额ID
	retrieved, _ := service.GetQuota(customer.ID, nil)

	// 测试更新配额
	retrieved.CPU = 16
	retrieved.Memory = 32768
	err = service.UpdateQuota(retrieved)
	assert.NoError(t, err)
	t.Log("更新配额成功")

	// 验证更新
	updated, err := service.GetQuota(customer.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 16, updated.CPU)
	assert.Equal(t, int64(32768), updated.Memory)

	// 清理
	service.DeleteQuota(retrieved.ID)
}

func TestResourceQuotaService_CheckQuota(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-check-quota-" + uuid.New().String()[:8],
		Email:        "check-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Check Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 设置配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)
	defer func() {
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	}()

	// 测试配额足够的情况
	request := &ResourceRequest{
		CPU:     8,
		Memory:  16384,
		GPU:     2,
		Storage: 500,
	}
	ok, err := service.CheckQuota(customer.ID, nil, request)
	assert.NoError(t, err)
	assert.True(t, ok)
	t.Log("配额检查通过（足够）")

	// 测试配额不足的情况
	requestExceed := &ResourceRequest{
		CPU:     32,
		Memory:  16384,
		GPU:     2,
		Storage: 500,
	}
	ok, err = service.CheckQuota(customer.ID, nil, requestExceed)
	assert.Error(t, err)
	assert.False(t, ok)
	t.Log("配额检查失败（不足）")
}

func TestResourceQuotaService_GetAvailableQuota(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-available-quota-" + uuid.New().String()[:8],
		Email:        "available-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Available Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 设置配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)
	defer func() {
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	}()

	// 测试获取可用配额
	available, err := service.GetAvailableQuota(customer.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 16, available.CPU)
	assert.Equal(t, int64(32768), available.Memory)
	assert.Equal(t, 4, available.GPU)
	assert.Equal(t, int64(1000), available.Storage)
	t.Log("获取可用配额成功")
}

// TestResourceQuotaService_CheckQuota_AllResourceTypes 测试所有资源类型的配额检查
func TestResourceQuotaService_CheckQuota_AllResourceTypes(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-all-resources-" + uuid.New().String()[:8],
		Email:        "all-resources-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test All Resources User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 设置配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)
	defer func() {
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	}()

	// 测试 CPU 超限
	t.Run("CPU exceeded", func(t *testing.T) {
		request := &ResourceRequest{
			CPU:     32, // 超过配额 16
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "CPU")
		t.Log("CPU 配额检查失败（超限）")
	})

	// 测试 Memory 超限
	t.Run("Memory exceeded", func(t *testing.T) {
		request := &ResourceRequest{
			CPU:     8,
			Memory:  65536, // 超过配额 32768
			GPU:     2,
			Storage: 500,
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "Memory")
		t.Log("Memory 配额检查失败（超限）")
	})

	// 测试 GPU 超限
	t.Run("GPU exceeded", func(t *testing.T) {
		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     8, // 超过配额 4
			Storage: 500,
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "GPU")
		t.Log("GPU 配额检查失败（超限）")
	})

	// 测试 Storage 超限
	t.Run("Storage exceeded", func(t *testing.T) {
		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 2000, // 超过配额 1000
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "Storage")
		t.Log("Storage 配额检查失败（超限）")
	})
}

// TestResourceQuotaService_CheckQuotaInTx 测试事务中的配额检查（并发安全）
func TestResourceQuotaService_CheckQuotaInTx(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()
	db := database.GetDB()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-quota-tx-" + uuid.New().String()[:8],
		Email:        "quota-tx-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Quota Tx User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 设置配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)
	defer func() {
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	}()

	// 测试在事务中检查配额（足够）
	t.Run("Check quota in transaction - sufficient", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}
		ok, err := service.CheckQuotaInTx(tx, customer.ID, nil, request)
		assert.NoError(t, err)
		assert.True(t, ok)
		t.Log("事务中配额检查通过（足够）")
	})

	// 测试在事务中检查配额（不足）
	t.Run("Check quota in transaction - insufficient", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		request := &ResourceRequest{
			CPU:     32, // 超过配额
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}
		ok, err := service.CheckQuotaInTx(tx, customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "CPU")
		t.Log("事务中配额检查失败（不足）")
	})

	// 测试传入 nil 事务应该报错
	t.Run("Check quota with nil transaction", func(t *testing.T) {
		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}
		ok, err := service.CheckQuotaInTx(nil, customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), "事务不能为空")
		t.Log("nil 事务检查正确报错")
	})
}

// TestQuotaExceededError_NegativeAvailable 测试负数可用配额的错误信息
func TestQuotaExceededError_NegativeAvailable(t *testing.T) {
	// 测试正常的配额不足错误信息
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
		t.Logf("正常错误信息: %s", errMsg)
	})

	// 测试负数可用配额的错误信息（已超额使用）
	t.Run("Negative available", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "Memory",
			Requested: 16384,
			Available: -4096, // 负数表示已超额使用
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "Memory 配额不足")
		assert.Contains(t, errMsg, "需要 16384")
		assert.Contains(t, errMsg, "可用 0")
		assert.Contains(t, errMsg, "已超额使用 4096")
		t.Logf("负数可用配额错误信息: %s", errMsg)
	})

	// 测试零可用配额的错误信息
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
		t.Logf("零可用配额错误信息: %s", errMsg)
	})
}
