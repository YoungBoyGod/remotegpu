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
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
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
		UserID:  customer.ID,
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
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
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
		UserID:  customer.ID,
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
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
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
		UserID:  customer.ID,
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
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
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
		UserID:  customer.ID,
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
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
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
		UserID:  customer.ID,
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
	customerDao := dao.NewUserDao()
	db := database.GetDB()

	// 创建测试客户
	customer := &entity.User{
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
		UserID:  customer.ID,
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

// TestResourceQuotaService_SetQuota_BoundaryConditions 测试SetQuota边界条件
func TestResourceQuotaService_SetQuota_BoundaryConditions(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-setquota-boundary-" + uuid.New().String()[:8],
		Email:        "setquota-boundary-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test SetQuota Boundary User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 测试负数配额值
	t.Run("Negative quota values", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         -1,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}
		err := service.SetQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
		t.Log("负数配额值正确返回错误")
	})

	// 测试零配额值（应该允许）
	t.Run("Zero quota values", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         0,
			Memory:      0,
			GPU:         0,
			Storage:     0,
		}
		err := service.SetQuota(quota)
		assert.NoError(t, err)
		t.Log("零配额值设置成功")

		// 清理
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	})

	// 测试重复设置配额（应该更新而不是创建新记录）
	t.Run("Duplicate set quota", func(t *testing.T) {
		// 第一次设置
		quota1 := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}
		err := service.SetQuota(quota1)
		assert.NoError(t, err)

		// 第二次设置（应该更新）
		quota2 := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}
		err = service.SetQuota(quota2)
		assert.NoError(t, err)

		// 验证只有一条记录，且值已更新
		retrieved, err := service.GetQuota(customer.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 16, retrieved.CPU)
		assert.Equal(t, int64(32768), retrieved.Memory)
		t.Log("重复设置配额正确更新而不是创建新记录")

		// 清理
		service.DeleteQuota(retrieved.ID)
	})
}

// TestResourceQuotaService_DeleteQuota 测试删除配额
func TestResourceQuotaService_DeleteQuota(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-delete-quota-" + uuid.New().String()[:8],
		Email:        "delete-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Delete Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 测试删除不存在的配额
	t.Run("Delete non-existent quota", func(t *testing.T) {
		err := service.DeleteQuota(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额不存在")
		t.Log("删除不存在的配额正确返回错误")
	})

	// 测试删除存在的配额
	t.Run("Delete existing quota", func(t *testing.T) {
		// 先创建配额
		quota := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}
		err := service.SetQuota(quota)
		assert.NoError(t, err)

		// 获取配额ID
		retrieved, err := service.GetQuota(customer.ID, nil)
		assert.NoError(t, err)

		// 删除配额
		err = service.DeleteQuota(retrieved.ID)
		assert.NoError(t, err)
		t.Log("删除存在的配额成功")

		// 验证配额已被删除
		_, err = service.GetQuota(customer.ID, nil)
		assert.Error(t, err)
		t.Log("验证配额已被删除")
	})
}

// TestResourceQuotaService_GetUsedResources 测试获取已使用资源
func TestResourceQuotaService_GetUsedResources(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()
	workspaceDao := dao.NewWorkspaceDao()
	db := database.GetDB()

	// 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-used-resources-" + uuid.New().String()[:8],
		Email:        "used-resources-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Used Resources User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 创建测试工作空间
	workspace := &entity.Workspace{
		Name:        "test-workspace-" + uuid.New().String()[:8],
		Description: "Test Workspace",
		OwnerID:     customer.ID,
		Status:      "active",
	}
	err = workspaceDao.Create(workspace)
	assert.NoError(t, err)
	defer workspaceDao.Delete(workspace.ID)

	// 创建测试主机
	host := &entity.Host{
		ID:             "test-host-" + uuid.New().String()[:8],
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "k8s",
		Status:         "active",
		TotalCPU:       16,
		TotalMemory:    32000,
		TotalGPU:       4,
	}
	err = db.Create(host).Error
	assert.NoError(t, err)
	defer db.Where("id = ?", host.ID).Delete(&entity.Host{})

	// 测试场景1：没有环境时，资源使用为0
	t.Run("No environments", func(t *testing.T) {
		used, err := service.GetUsedResources(customer.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, used.CPU)
		assert.Equal(t, int64(0), used.Memory)
		assert.Equal(t, 0, used.GPU)
		assert.Equal(t, int64(0), used.Storage)
		t.Log("没有环境时资源使用为0")
	})

	// 创建测试环境
	storage1 := int64(1000)
	env1 := &entity.Environment{
		ID:          "env-running-" + uuid.New().String()[:8],
		Name:        "env-running",
		UserID:  customer.ID,
		WorkspaceID: &workspace.ID,
		HostID:      host.ID,
		Status:      "running",
		CPU:         4,
		Memory:      8192,
		GPU:         1,
		Storage:     &storage1,
		Image:       "test-image:latest",
	}
	err = db.Create(env1).Error
	assert.NoError(t, err)
	defer db.Where("id = ?", env1.ID).Delete(&entity.Environment{})

	storage2 := int64(2000)
	env2 := &entity.Environment{
		ID:          "env-creating-" + uuid.New().String()[:8],
		Name:        "env-creating",
		UserID:  customer.ID,
		WorkspaceID: &workspace.ID,
		HostID:      host.ID,
		Status:      "creating",
		CPU:         2,
		Memory:      4096,
		GPU:         1,
		Storage:     &storage2,
		Image:       "test-image:latest",
	}
	err = db.Create(env2).Error
	assert.NoError(t, err)
	defer db.Where("id = ?", env2.ID).Delete(&entity.Environment{})

	env3 := &entity.Environment{
		ID:          "env-stopped-" + uuid.New().String()[:8],
		Name:        "env-stopped",
		UserID:  customer.ID,
		WorkspaceID: &workspace.ID,
		HostID:      host.ID,
		Status:      "stopped",
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     nil, // 测试 nil Storage
		Image:       "test-image:latest",
	}
	err = db.Create(env3).Error
	assert.NoError(t, err)
	defer db.Where("id = ?", env3.ID).Delete(&entity.Environment{})

	// 测试场景2：统计运行中和创建中的环境（不包括停止的）
	t.Run("Running and creating environments", func(t *testing.T) {
		used, err := service.GetUsedResources(customer.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 6, used.CPU) // 4 + 2
		assert.Equal(t, int64(12288), used.Memory) // 8192 + 4096
		assert.Equal(t, 2, used.GPU) // 1 + 1
		assert.Equal(t, int64(3000), used.Storage) // 1000 + 2000
		t.Log("正确统计运行中和创建中的环境资源")
	})

	// 测试场景3：工作空间级别资源统计
	t.Run("Workspace level resources", func(t *testing.T) {
		used, err := service.GetUsedResources(customer.ID, &workspace.ID)
		assert.NoError(t, err)
		assert.Equal(t, 6, used.CPU)
		assert.Equal(t, int64(12288), used.Memory)
		assert.Equal(t, 2, used.GPU)
		assert.Equal(t, int64(3000), used.Storage)
		t.Log("工作空间级别资源统计正确")
	})
}

// TestResourceQuotaService_UpdateQuota_ErrorCases 测试UpdateQuota的错误场景
func TestResourceQuotaService_UpdateQuota_ErrorCases(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-update-error-" + uuid.New().String()[:8],
		Email:        "update-error-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Update Error User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 测试场景1：更新不存在的配额
	t.Run("Update non-existent quota", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			ID:      99999,
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		err := service.UpdateQuota(quota)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额不存在")
		t.Log("更新不存在的配额正确返回错误")
	})

	// 测试场景2：更新为负数配额
	t.Run("Update with negative values", func(t *testing.T) {
		// 先创建配额
		quota := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}
		err := service.SetQuota(quota)
		assert.NoError(t, err)

		// 获取配额ID
		retrieved, _ := service.GetQuota(customer.ID, nil)

		// 尝试更新为负数
		retrieved.CPU = -1
		err = service.UpdateQuota(retrieved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "配额值不能为负数")
		t.Log("更新为负数配额正确返回错误")

		// 清理
		service.DeleteQuota(retrieved.ID)
	})

	// 测试场景3：更新配额低于已使用资源
	t.Run("Update quota below used resources", func(t *testing.T) {
		workspaceDao := dao.NewWorkspaceDao()
		db := database.GetDB()

		// 创建工作空间
		workspace := &entity.Workspace{
			Name:        "test-workspace-" + uuid.New().String()[:8],
			Description: "Test Workspace",
			OwnerID:     customer.ID,
			Status:      "active",
		}
		err := workspaceDao.Create(workspace)
		assert.NoError(t, err)
		defer workspaceDao.Delete(workspace.ID)

		// 创建主机
		host := &entity.Host{
			ID:             "test-host-" + uuid.New().String()[:8],
			Name:           "Test Host",
			IPAddress:      "192.168.1.100",
			OSType:         "linux",
			DeploymentMode: "k8s",
			Status:         "active",
			TotalCPU:       16,
			TotalMemory:    32000,
			TotalGPU:       4,
		}
		err = db.Create(host).Error
		assert.NoError(t, err)
		defer db.Where("id = ?", host.ID).Delete(&entity.Host{})

		// 创建配额
		quota := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     2000,
		}
		err = service.SetQuota(quota)
		assert.NoError(t, err)

		// 创建运行中的环境（使用部分资源）
		storage := int64(1000)
		env := &entity.Environment{
			ID:          "env-test-" + uuid.New().String()[:8],
			Name:        "env-test",
			UserID:  customer.ID,
			WorkspaceID: &workspace.ID,
			HostID:      host.ID,
			Status:      "running",
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     &storage,
			Image:       "test-image:latest",
		}
		err = db.Create(env).Error
		assert.NoError(t, err)
		defer db.Where("id = ?", env.ID).Delete(&entity.Environment{})

		// 获取配额ID
		retrieved, _ := service.GetQuota(customer.ID, nil)

		// 尝试更新CPU配额低于已使用量
		retrieved.CPU = 4 // 已使用8，新配额4
		err = service.UpdateQuota(retrieved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPU配额不能小于已使用量")
		t.Log("更新CPU配额低于已使用量正确返回错误")

		// 恢复配额
		retrieved.CPU = 16

		// 尝试更新Memory配额低于已使用量
		retrieved.Memory = 8192 // 已使用16384，新配额8192
		err = service.UpdateQuota(retrieved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "内存配额不能小于已使用量")
		t.Log("更新Memory配额低于已使用量正确返回错误")

		// 清理
		service.DeleteQuota(retrieved.ID)
	})
}

// TestResourceQuotaService_GetAvailableQuota_EdgeCases 测试GetAvailableQuota的边界情况
func TestResourceQuotaService_GetAvailableQuota_EdgeCases(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()

	// 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-available-edge-" + uuid.New().String()[:8],
		Email:        "available-edge-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Available Edge User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 测试场景1：配额不存在
	t.Run("Quota not found", func(t *testing.T) {
		_, err := service.GetAvailableQuota(customer.ID, nil)
		assert.Error(t, err)
		t.Log("配额不存在时正确返回错误")
	})

	// 测试场景2：已使用资源超过配额（负数可用配额应设为0）
	t.Run("Negative available resources", func(t *testing.T) {
		workspaceDao := dao.NewWorkspaceDao()
		db := database.GetDB()

		// 创建工作空间
		workspace := &entity.Workspace{
			Name:        "test-workspace-" + uuid.New().String()[:8],
			Description: "Test Workspace",
			OwnerID:     customer.ID,
			Status:      "active",
		}
		err := workspaceDao.Create(workspace)
		assert.NoError(t, err)
		defer workspaceDao.Delete(workspace.ID)

		// 创建主机
		host := &entity.Host{
			ID:             "test-host-" + uuid.New().String()[:8],
			Name:           "Test Host",
			IPAddress:      "192.168.1.100",
			OSType:         "linux",
			DeploymentMode: "k8s",
			Status:         "active",
			TotalCPU:       32,
			TotalMemory:    64000,
			TotalGPU:       8,
		}
		err = db.Create(host).Error
		assert.NoError(t, err)
		defer db.Where("id = ?", host.ID).Delete(&entity.Host{})

		// 创建较小的配额
		quota := &entity.ResourceQuota{
			UserID:  customer.ID,
			WorkspaceID: nil,
			CPU:         4,
			Memory:      8192,
			GPU:         1,
			Storage:     500,
		}
		err = service.SetQuota(quota)
		assert.NoError(t, err)

		// 创建使用超过配额的环境
		storage := int64(1000)
		env := &entity.Environment{
			ID:          "env-exceed-" + uuid.New().String()[:8],
			Name:        "env-exceed",
			UserID:  customer.ID,
			WorkspaceID: &workspace.ID,
			HostID:      host.ID,
			Status:      "running",
			CPU:         8,     // 超过配额4
			Memory:      16384, // 超过配额8192
			GPU:         2,     // 超过配额1
			Storage:     &storage,
			Image:       "test-image:latest",
		}
		err = db.Create(env).Error
		assert.NoError(t, err)
		defer db.Where("id = ?", env.ID).Delete(&entity.Environment{})

		// 获取可用配额
		available, err := service.GetAvailableQuota(customer.ID, nil)
		assert.NoError(t, err)

		// 验证负数可用配额被设为0
		assert.Equal(t, 0, available.CPU, "负数CPU可用配额应设为0")
		assert.Equal(t, int64(0), available.Memory, "负数Memory可用配额应设为0")
		assert.Equal(t, 0, available.GPU, "负数GPU可用配额应设为0")
		t.Log("负数可用配额正确设为0")

		// 清理
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	})
}

// TestResourceQuotaService_GetQuota_ErrorCases 测试GetQuota的错误场景
func TestResourceQuotaService_GetQuota_ErrorCases(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()
	workspaceDao := dao.NewWorkspaceDao()

	// 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-getquota-error-" + uuid.New().String()[:8],
		Email:        "getquota-error-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test GetQuota Error User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 创建工作空间
	workspace := &entity.Workspace{
		Name:        "test-workspace-" + uuid.New().String()[:8],
		Description: "Test Workspace",
		OwnerID:     customer.ID,
		Status:      "active",
	}
	err = workspaceDao.Create(workspace)
	assert.NoError(t, err)
	defer workspaceDao.Delete(workspace.ID)

	// 测试场景1：用户级别配额不存在
	t.Run("User level quota not found", func(t *testing.T) {
		_, err := service.GetQuota(customer.ID, nil)
		assert.Error(t, err)
		t.Log("用户级别配额不存在时正确返回错误")
	})

	// 测试场景2：工作空间级别配额不存在
	t.Run("Workspace level quota not found", func(t *testing.T) {
		_, err := service.GetQuota(customer.ID, &workspace.ID)
		assert.Error(t, err)
		t.Log("工作空间级别配额不存在时正确返回错误")
	})
}

// TestResourceQuotaService_ErrorScenarios 测试错误场景以提高覆盖率
func TestResourceQuotaService_ErrorScenarios(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()

	// 测试 QuotaExceededError.Error() - Available >= 0 的情况
	t.Run("QuotaExceededError_PositiveAvailable", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "CPU",
			Requested: 100,
			Available: 50,
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "CPU")
		assert.Contains(t, errMsg, "100")
		assert.Contains(t, errMsg, "50")
		assert.NotContains(t, errMsg, "已超额使用")
		t.Logf("QuotaExceededError (Available >= 0): %s", errMsg)
	})

	// 测试 QuotaExceededError.Error() - Available < 0 的情况
	t.Run("QuotaExceededError_NegativeAvailable", func(t *testing.T) {
		err := &QuotaExceededError{
			Resource:  "Memory",
			Requested: 100,
			Available: -20,
		}
		errMsg := err.Error()
		assert.Contains(t, errMsg, "Memory")
		assert.Contains(t, errMsg, "100")
		assert.Contains(t, errMsg, "已超额使用")
		assert.Contains(t, errMsg, "20")
		t.Logf("QuotaExceededError (Available < 0): %s", errMsg)
	})

	// 测试 GetQuota - 不存在的用户级别配额
	t.Run("GetQuota_UserQuotaNotFound", func(t *testing.T) {
		_, err := service.GetQuota(99999, nil)
		assert.Error(t, err)
		t.Logf("GetQuota 用户配额不存在: %v", err)
	})

	// 测试 GetQuota - 不存在的工作空间级别配额
	t.Run("GetQuota_WorkspaceQuotaNotFound", func(t *testing.T) {
		workspaceID := uint(99999)
		_, err := service.GetQuota(99999, &workspaceID)
		assert.Error(t, err)
		t.Logf("GetQuota 工作空间配额不存在: %v", err)
	})

	// 测试 GetQuotaInTx - nil 事务
	t.Run("GetQuotaInTx_NilTransaction", func(t *testing.T) {
		_, err := service.GetQuotaInTx(nil, 1, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "事务不能为空")
		t.Logf("GetQuotaInTx nil事务: %v", err)
	})

	t.Log("错误场景测试完成")
}

// TestResourceQuotaService_GetUsedResources_Bug3Fix 测试Bug #3修复：用户级别配额应该统计所有工作空间的环境
// Bug描述：当workspaceID为nil时，应该统计用户在所有工作空间的环境，而不是只统计workspace_id IS NULL的环境
func TestResourceQuotaService_GetUsedResources_Bug3Fix(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewUserDao()
	workspaceDao := dao.NewWorkspaceDao()
	envDao := dao.NewEnvironmentDao()

	// 1. 创建测试客户
	customer := &entity.User{
		UUID:         uuid.New(),
		Username:     "test-bug3-" + uuid.New().String()[:8],
		Email:        "bug3-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Bug3 Test User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	t.Logf("创建测试客户: ID=%d, Username=%s", customer.ID, customer.Username)

	// 2. 创建两个工作空间
	workspace1 := &entity.Workspace{
		OwnerID:     customer.ID,
		Name:        "workspace-1-" + uuid.New().String()[:8],
		Description: "Test Workspace 1",
		Status:      "active",
	}
	err = workspaceDao.Create(workspace1)
	assert.NoError(t, err)
	defer workspaceDao.Delete(workspace1.ID)

	workspace2 := &entity.Workspace{
		OwnerID:     customer.ID,
		Name:        "workspace-2-" + uuid.New().String()[:8],
		Description: "Test Workspace 2",
		Status:      "active",
	}
	err = workspaceDao.Create(workspace2)
	assert.NoError(t, err)
	defer workspaceDao.Delete(workspace2.ID)

	t.Logf("创建工作空间: workspace1.ID=%d, workspace2.ID=%d", workspace1.ID, workspace2.ID)

	// 3. 创建测试主机
	db := database.GetDB()
	host1 := &entity.Host{
		ID:             "host-1",
		Name:           "Test Host 1",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "k8s",
		Status:         "active",
		TotalCPU:       32,
		TotalMemory:    64000,
		TotalGPU:       8,
	}
	err = db.Create(host1).Error
	assert.NoError(t, err)
	defer db.Where("id = ?", host1.ID).Delete(&entity.Host{})

	host2 := &entity.Host{
		ID:             "host-2",
		Name:           "Test Host 2",
		IPAddress:      "192.168.1.101",
		OSType:         "linux",
		DeploymentMode: "k8s",
		Status:         "active",
		TotalCPU:       32,
		TotalMemory:    64000,
		TotalGPU:       8,
	}
	err = db.Create(host2).Error
	assert.NoError(t, err)
	defer db.Where("id = ?", host2.ID).Delete(&entity.Host{})

	t.Logf("创建测试主机: host1.ID=%s, host2.ID=%s", host1.ID, host2.ID)

	// 4. 在workspace1中创建2个running状态的环境
	storage1 := int64(10 * 1024 * 1024 * 1024) // 10GB
	env1 := &entity.Environment{
		ID:          "env-ws1-1-" + uuid.New().String()[:8],
		UserID:  customer.ID,
		WorkspaceID: &workspace1.ID,
		HostID:      "host-1",
		Name:        "env-ws1-1",
		Image:       "ubuntu:20.04",
		Status:      "running",
		CPU:         2,
		Memory:      1024,
		GPU:         1,
		Storage:     &storage1,
	}
	err = envDao.Create(env1)
	assert.NoError(t, err)
	defer envDao.Delete(env1.ID)

	env2 := &entity.Environment{
		ID:          "env-ws1-2-" + uuid.New().String()[:8],
		UserID:  customer.ID,
		WorkspaceID: &workspace1.ID,
		HostID:      "host-1",
		Name:        "env-ws1-2",
		Image:       "ubuntu:20.04",
		Status:      "running",
		CPU:         1,
		Memory:      512,
		GPU:         0,
		Storage:     &storage1,
	}
	err = envDao.Create(env2)
	assert.NoError(t, err)
	defer envDao.Delete(env2.ID)

	t.Logf("在workspace1中创建2个环境: env1(CPU=2,Mem=1024,GPU=1), env2(CPU=1,Mem=512,GPU=0)")

	// 5. 在workspace2中创建1个running状态的环境
	storage2 := int64(5 * 1024 * 1024 * 1024) // 5GB
	env3 := &entity.Environment{
		ID:          "env-ws2-1-" + uuid.New().String()[:8],
		UserID:  customer.ID,
		WorkspaceID: &workspace2.ID,
		HostID:      "host-2",
		Name:        "env-ws2-1",
		Image:       "ubuntu:20.04",
		Status:      "running",
		CPU:         4,
		Memory:      2048,
		GPU:         2,
		Storage:     &storage2,
	}
	err = envDao.Create(env3)
	assert.NoError(t, err)
	defer envDao.Delete(env3.ID)

	t.Logf("在workspace2中创建1个环境: env3(CPU=4,Mem=2048,GPU=2)")

	// 6. 创建1个creating状态的环境（应该也被统计）
	env4 := &entity.Environment{
		ID:          "env-ws1-3-" + uuid.New().String()[:8],
		UserID:  customer.ID,
		WorkspaceID: &workspace1.ID,
		HostID:      "host-1",
		Name:        "env-ws1-3",
		Image:       "ubuntu:20.04",
		Status:      "creating",
		CPU:         1,
		Memory:      256,
		GPU:         0,
		Storage:     nil,
	}
	err = envDao.Create(env4)
	assert.NoError(t, err)
	defer envDao.Delete(env4.ID)

	t.Logf("在workspace1中创建1个creating状态的环境: env4(CPU=1,Mem=256,GPU=0)")

	// 测试1: 用户级别配额（workspaceID=nil）应该统计所有工作空间的环境
	t.Run("UserLevel_ShouldCountAllWorkspaces", func(t *testing.T) {
		used, err := service.GetUsedResources(customer.ID, nil)
		assert.NoError(t, err)
		assert.NotNil(t, used)

		// 期望值：workspace1的3个环境 + workspace2的1个环境
		// CPU: 2 + 1 + 4 + 1 = 8
		// Memory: 1024 + 512 + 2048 + 256 = 3840
		// GPU: 1 + 0 + 2 + 0 = 3
		// Storage: 10GB + 10GB + 5GB + 0 = 25GB
		expectedCPU := 8
		expectedMemory := int64(3840)
		expectedGPU := 3
		expectedStorage := int64(25 * 1024 * 1024 * 1024)

		assert.Equal(t, expectedCPU, used.CPU, "用户级别配额应该统计所有工作空间的CPU")
		assert.Equal(t, expectedMemory, used.Memory, "用户级别配额应该统计所有工作空间的Memory")
		assert.Equal(t, expectedGPU, used.GPU, "用户级别配额应该统计所有工作空间的GPU")
		assert.Equal(t, expectedStorage, used.Storage, "用户级别配额应该统计所有工作空间的Storage")

		t.Logf("✅ Bug #3修复验证通过：用户级别配额统计了所有工作空间的环境")
		t.Logf("   已使用资源: CPU=%d, Memory=%d, GPU=%d, Storage=%dGB",
			used.CPU, used.Memory, used.GPU, used.Storage/(1024*1024*1024))
	})

	// 测试2: 工作空间级别配额应该只统计指定工作空间的环境
	t.Run("WorkspaceLevel_ShouldCountOnlySpecifiedWorkspace", func(t *testing.T) {
		// 测试workspace1
		used1, err := service.GetUsedResources(customer.ID, &workspace1.ID)
		assert.NoError(t, err)
		assert.NotNil(t, used1)

		// 期望值：workspace1的3个环境
		// CPU: 2 + 1 + 1 = 4
		// Memory: 1024 + 512 + 256 = 1792
		// GPU: 1 + 0 + 0 = 1
		// Storage: 10GB + 10GB + 0 = 20GB
		assert.Equal(t, 4, used1.CPU, "workspace1应该只统计自己的CPU")
		assert.Equal(t, int64(1792), used1.Memory, "workspace1应该只统计自己的Memory")
		assert.Equal(t, 1, used1.GPU, "workspace1应该只统计自己的GPU")
		assert.Equal(t, int64(20*1024*1024*1024), used1.Storage, "workspace1应该只统计自己的Storage")

		t.Logf("✅ workspace1资源统计正确: CPU=%d, Memory=%d, GPU=%d, Storage=%dGB",
			used1.CPU, used1.Memory, used1.GPU, used1.Storage/(1024*1024*1024))

		// 测试workspace2
		used2, err := service.GetUsedResources(customer.ID, &workspace2.ID)
		assert.NoError(t, err)
		assert.NotNil(t, used2)

		// 期望值：workspace2的1个环境
		// CPU: 4
		// Memory: 2048
		// GPU: 2
		// Storage: 5GB
		assert.Equal(t, 4, used2.CPU, "workspace2应该只统计自己的CPU")
		assert.Equal(t, int64(2048), used2.Memory, "workspace2应该只统计自己的Memory")
		assert.Equal(t, 2, used2.GPU, "workspace2应该只统计自己的GPU")
		assert.Equal(t, int64(5*1024*1024*1024), used2.Storage, "workspace2应该只统计自己的Storage")

		t.Logf("✅ workspace2资源统计正确: CPU=%d, Memory=%d, GPU=%d, Storage=%dGB",
			used2.CPU, used2.Memory, used2.GPU, used2.Storage/(1024*1024*1024))
	})

	t.Log("✅ Bug #3修复验证完成：GetUsedResources逻辑正确")
}
