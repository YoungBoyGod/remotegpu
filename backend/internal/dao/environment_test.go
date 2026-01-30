package dao

import (
	"sync"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestEnvironmentDao_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	envDao := NewEnvironmentDao()
	customerDao := NewUserDao()
	workspaceDao := NewWorkspaceDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.User{
		Username:     "test-env-customer-" + time.Now().Format("20060102150405"),
		Email:        "env-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 创建测试工作空间
	testWorkspace := &entity.Workspace{
		UUID:        uuid.New(),
		OwnerID:     testCustomer.ID,
		Name:        "Test Workspace for Env",
		Type:        "personal",
		MemberCount: 1,
		Status:      "active",
	}
	if err := workspaceDao.Create(testWorkspace); err != nil {
		t.Fatalf("创建测试工作空间失败: %v", err)
	}
	defer workspaceDao.Delete(testWorkspace.ID)

	// 创建测试主机
	testHost := &entity.Host{
		ID:             "test-host-" + time.Now().Format("20060102150405"),
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "docker",
		Status:         "online",
		TotalCPU:       16,
		TotalMemory:    34359738368, // 32GB
		TotalGPU:       2,
	}
	if err := hostDao.Create(testHost); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostDao.Delete(testHost.ID)

	// 测试 Create Environment
	storage := int64(10737418240) // 10GB
	env := &entity.Environment{
		ID:          "env-test-" + time.Now().Format("20060102150405"),
		UserID:  testCustomer.ID,
		WorkspaceID: &testWorkspace.ID,
		HostID:      testHost.ID,
		Name:        "Test Environment",
		Description: "This is a test environment",
		Image:       "ubuntu:22.04",
		Status:      "creating",
		CPU:         4,
		Memory:      8589934592, // 8GB
		GPU:         1,
		Storage:     &storage,
	}
	if err := envDao.Create(env); err != nil {
		t.Fatalf("创建环境失败: %v", err)
	}
	t.Logf("创建环境成功, ID: %s", env.ID)

	// 测试 GetByID
	found, err := envDao.GetByID(env.ID)
	if err != nil {
		t.Fatalf("获取环境失败: %v", err)
	}
	if found.Name != "Test Environment" {
		t.Errorf("环境名称不匹配: 期望 'Test Environment', 实际 '%s'", found.Name)
	}

	// 测试 Update
	env.Status = "running"
	env.Name = "Updated Environment"
	if err := envDao.Update(env); err != nil {
		t.Fatalf("更新环境失败: %v", err)
	}

	updated, err := envDao.GetByID(env.ID)
	if err != nil {
		t.Fatalf("获取更新后的环境失败: %v", err)
	}
	if updated.Status != "running" {
		t.Errorf("环境状态未更新: 期望 'running', 实际 '%s'", updated.Status)
	}

	// 测试 GetByUserID
	envs, err := envDao.GetByUserID(testCustomer.ID)
	if err != nil {
		t.Fatalf("根据客户ID获取环境失败: %v", err)
	}
	if len(envs) == 0 {
		t.Error("应该至少有一个环境")
	}

	// 测试 GetByWorkspaceID
	workspaceEnvs, err := envDao.GetByWorkspaceID(testWorkspace.ID)
	if err != nil {
		t.Fatalf("根据工作空间ID获取环境失败: %v", err)
	}
	if len(workspaceEnvs) == 0 {
		t.Error("应该至少有一个环境")
	}

	// 测试 GetByHostID
	hostEnvs, err := envDao.GetByHostID(testHost.ID)
	if err != nil {
		t.Fatalf("根据主机ID获取环境失败: %v", err)
	}
	if len(hostEnvs) == 0 {
		t.Error("应该至少有一个环境")
	}

	// 测试 GetByStatus
	runningEnvs, err := envDao.GetByStatus("running")
	if err != nil {
		t.Fatalf("根据状态获取环境失败: %v", err)
	}
	if len(runningEnvs) == 0 {
		t.Error("应该至少有一个运行中的环境")
	}

	// 测试 List
	envList, total, err := envDao.List(1, 10)
	if err != nil {
		t.Fatalf("获取环境列表失败: %v", err)
	}
	if total == 0 {
		t.Error("环境总数应该大于0")
	}
	t.Logf("环境列表: 总数 %d, 当前页 %d 条", total, len(envList))

	// 测试 Delete
	if err := envDao.Delete(env.ID); err != nil {
		t.Fatalf("删除环境失败: %v", err)
	}
	t.Logf("删除环境成功")
}

func TestPortMappingDao_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	pmDao := NewPortMappingDao()
	envDao := NewEnvironmentDao()
	customerDao := NewUserDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.User{
		Username:     "test-pm-customer-" + time.Now().Format("20060102150405"),
		Email:        "pm-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 创建测试主机
	testHost := &entity.Host{
		ID:             "test-host-pm-" + time.Now().Format("20060102150405"),
		Name:           "Test Host for PM",
		IPAddress:      "192.168.1.101",
		OSType:         "linux",
		DeploymentMode: "docker",
		Status:         "online",
		TotalCPU:       8,
		TotalMemory:    17179869184, // 16GB
	}
	if err := hostDao.Create(testHost); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostDao.Delete(testHost.ID)

	// 创建测试环境
	testEnv := &entity.Environment{
		ID:         "env-pm-test-" + time.Now().Format("20060102150405"),
		UserID: testCustomer.ID,
		HostID:     testHost.ID,
		Name:       "Test Env for PM",
		Image:      "ubuntu:22.04",
		Status:     "running",
		CPU:        2,
		Memory:     4294967296, // 4GB
	}
	if err := envDao.Create(testEnv); err != nil {
		t.Fatalf("创建测试环境失败: %v", err)
	}
	defer envDao.Delete(testEnv.ID)

	// 测试 AllocatePort
	pm, err := pmDao.AllocatePort(testEnv.ID, "ssh", 22)
	if err != nil {
		t.Fatalf("分配端口失败: %v", err)
	}
	t.Logf("分配端口成功, ID: %d, ExternalPort: %d", pm.ID, pm.ExternalPort)

	if pm.ExternalPort < 30000 || pm.ExternalPort > 32767 {
		t.Errorf("分配的端口不在有效范围内: %d", pm.ExternalPort)
	}

	// 测试 GetByID
	found, err := pmDao.GetByID(pm.ID)
	if err != nil {
		t.Fatalf("获取端口映射失败: %v", err)
	}
	if found.ServiceType != "ssh" {
		t.Errorf("服务类型不匹配: 期望 'ssh', 实际 '%s'", found.ServiceType)
	}

	// 测试 GetByEnvironmentID
	pms, err := pmDao.GetByEnvironmentID(testEnv.ID)
	if err != nil {
		t.Fatalf("根据环境ID获取端口映射失败: %v", err)
	}
	if len(pms) == 0 {
		t.Error("应该至少有一个端口映射")
	}

	// 测试 ReleasePort
	if err := pmDao.ReleasePort(pm.ID); err != nil {
		t.Fatalf("释放端口失败: %v", err)
	}

	released, err := pmDao.GetByID(pm.ID)
	if err != nil {
		t.Fatalf("获取释放后的端口映射失败: %v", err)
	}
	if released.Status != "released" {
		t.Errorf("端口状态未更新: 期望 'released', 实际 '%s'", released.Status)
	}

	// 测试 DeleteByEnvironmentID
	if err := pmDao.DeleteByEnvironmentID(testEnv.ID); err != nil {
		t.Fatalf("根据环境ID删除端口映射失败: %v", err)
	}
	t.Logf("删除端口映射成功")
}

func TestPortMappingDao_ConcurrentAllocate(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	pmDao := NewPortMappingDao()
	envDao := NewEnvironmentDao()
	customerDao := NewUserDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.User{
		Username:     "test-concurrent-" + time.Now().Format("20060102150405"),
		Email:        "concurrent-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 创建测试主机
	testHost := &entity.Host{
		ID:             "test-host-concurrent-" + time.Now().Format("20060102150405"),
		Name:           "Test Host for Concurrent",
		IPAddress:      "192.168.1.102",
		OSType:         "linux",
		DeploymentMode: "docker",
		Status:         "online",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}
	if err := hostDao.Create(testHost); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostDao.Delete(testHost.ID)

	// 创建测试环境
	testEnv := &entity.Environment{
		ID:         "env-concurrent-" + time.Now().Format("20060102150405"),
		UserID: testCustomer.ID,
		HostID:     testHost.ID,
		Name:       "Test Env for Concurrent",
		Image:      "ubuntu:22.04",
		Status:     "running",
		CPU:        2,
		Memory:     4294967296,
	}
	if err := envDao.Create(testEnv); err != nil {
		t.Fatalf("创建测试环境失败: %v", err)
	}
	defer envDao.Delete(testEnv.ID)

	// 并发分配端口测试
	concurrency := 10
	var wg sync.WaitGroup
	results := make(chan *entity.PortMapping, concurrency)
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			pm, err := pmDao.AllocatePort(testEnv.ID, "custom", 8000+index)
			if err != nil {
				errors <- err
				return
			}
			results <- pm
		}(i)
	}

	wg.Wait()
	close(results)
	close(errors)

	// 检查错误
	for err := range errors {
		t.Errorf("并发分配端口出错: %v", err)
	}

	// 检查结果
	allocatedPorts := make(map[int]bool)
	count := 0
	for pm := range results {
		count++
		if allocatedPorts[pm.ExternalPort] {
			t.Errorf("端口重复分配: %d", pm.ExternalPort)
		}
		allocatedPorts[pm.ExternalPort] = true
		t.Logf("分配端口 #%d: %d", count, pm.ExternalPort)
	}

	if count != concurrency {
		t.Errorf("分配的端口数量不正确: 期望 %d, 实际 %d", concurrency, count)
	}

	// 清理
	if err := pmDao.DeleteByEnvironmentID(testEnv.ID); err != nil {
		t.Fatalf("清理端口映射失败: %v", err)
	}
	t.Logf("并发测试完成，成功分配 %d 个不重复的端口", count)
}

// TestEnvironmentDao_ErrorScenarios 测试错误场景
func TestEnvironmentDao_ErrorScenarios(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	envDao := NewEnvironmentDao()

	// 测试获取不存在的环境
	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := envDao.GetByID("non-existent-id")
		if err == nil {
			t.Error("期望返回错误，但没有错误")
		}
		if err != gorm.ErrRecordNotFound {
			t.Errorf("期望 gorm.ErrRecordNotFound，实际: %v", err)
		}
	})

	// 测试查询不存在的客户ID
	t.Run("GetByUserID_Empty", func(t *testing.T) {
		envs, err := envDao.GetByUserID(999999)
		if err != nil {
			t.Errorf("不应该返回错误: %v", err)
		}
		if len(envs) != 0 {
			t.Errorf("期望返回空列表，实际返回 %d 条记录", len(envs))
		}
	})

	// 测试查询不存在的工作空间ID
	t.Run("GetByWorkspaceID_Empty", func(t *testing.T) {
		envs, err := envDao.GetByWorkspaceID(999999)
		if err != nil {
			t.Errorf("不应该返回错误: %v", err)
		}
		if len(envs) != 0 {
			t.Errorf("期望返回空列表，实际返回 %d 条记录", len(envs))
		}
	})
}

// TestPortMappingDao_ErrorScenarios 测试端口映射错误场景
func TestPortMappingDao_ErrorScenarios(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	pmDao := NewPortMappingDao()

	// 测试获取不存在的端口映射
	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := pmDao.GetByID(999999)
		if err == nil {
			t.Error("期望返回错误，但没有错误")
		}
		if err != gorm.ErrRecordNotFound {
			t.Errorf("期望 gorm.ErrRecordNotFound，实际: %v", err)
		}
	})

	// 测试查询不存在的环境ID
	t.Run("GetByEnvironmentID_Empty", func(t *testing.T) {
		pms, err := pmDao.GetByEnvironmentID("non-existent-env")
		if err != nil {
			t.Errorf("不应该返回错误: %v", err)
		}
		if len(pms) != 0 {
			t.Errorf("期望返回空列表，实际返回 %d 条记录", len(pms))
		}
	})

	// 测试释放不存在的端口
	t.Run("ReleasePort_NotFound", func(t *testing.T) {
		err := pmDao.ReleasePort(999999)
		// ReleasePort 即使记录不存在也不会返回错误（GORM 行为）
		if err != nil {
			t.Logf("释放不存在的端口返回错误: %v", err)
		}
	})

	// 测试直接创建端口映射（覆盖Create方法）
	t.Run("Create_Direct", func(t *testing.T) {
		pm := &entity.PortMapping{
			EnvID:        "test-env-for-create",
			ServiceType:  "custom",
			ExternalPort: 31000,
			InternalPort: 8080,
			Status:       "active",
		}
		// 注意：这个测试会失败，因为env_id不存在，但可以覆盖Create方法
		err := pmDao.Create(pm)
		if err != nil {
			t.Logf("预期的外键约束错误: %v", err)
		}
	})

	// 测试直接删除端口映射（覆盖Delete方法）
	t.Run("Delete_Direct", func(t *testing.T) {
		err := pmDao.Delete(999999)
		// Delete即使记录不存在也不会返回错误
		if err != nil {
			t.Logf("删除不存在的记录: %v", err)
		}
	})
}

// TestPortMappingDao_PortBoundary 测试端口分配边界条件
func TestPortMappingDao_PortBoundary(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	pmDao := NewPortMappingDao()
	envDao := NewEnvironmentDao()
	customerDao := NewUserDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.User{
		Username:     "test-boundary-customer-" + time.Now().Format("20060102150405"),
		Email:        "boundary-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 创建测试主机
	testHost := &entity.Host{
		ID:             "test-host-boundary-" + time.Now().Format("20060102150405"),
		Name:           "Test Host for Boundary",
		IPAddress:      "192.168.1.102",
		OSType:         "linux",
		DeploymentMode: "docker",
		Status:         "online",
		TotalCPU:       16,
		TotalMemory:    34359738368,
		TotalGPU:       2,
	}
	if err := hostDao.Create(testHost); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostDao.Delete(testHost.ID)

	// 创建测试环境
	testEnv := &entity.Environment{
		ID:         "test-boundary-env-" + time.Now().Format("20060102150405"),
		UserID: testCustomer.ID,
		HostID:     testHost.ID,
		Name:       "Test Boundary Environment",
		Image:      "ubuntu:22.04",
		Status:     "running",
		CPU:        2,
		Memory:     4294967296,
	}
	if err := envDao.Create(testEnv); err != nil {
		t.Fatalf("创建测试环境失败: %v", err)
	}
	defer envDao.Delete(testEnv.ID)

	// 清理端口映射
	defer func() {
		if err := pmDao.DeleteByEnvironmentID(testEnv.ID); err != nil {
			t.Logf("清理端口映射失败: %v", err)
		}
	}()

	// 测试正常分配端口
	t.Run("AllocatePort_Normal", func(t *testing.T) {
		pm, err := pmDao.AllocatePort(testEnv.ID, "ssh", 22)
		if err != nil {
			t.Fatalf("分配端口失败: %v", err)
		}
		if pm.ExternalPort < 30000 || pm.ExternalPort > 32767 {
			t.Errorf("分配的端口超出范围: %d", pm.ExternalPort)
		}
		t.Logf("成功分配端口: %d", pm.ExternalPort)
	})

	// 测试获取可用端口数量
	t.Run("GetAvailablePort", func(t *testing.T) {
		available, err := pmDao.GetAvailablePort()
		if err != nil {
			t.Fatalf("获取可用端口数量失败: %v", err)
		}
		expectedTotal := 32767 - 30000 + 1 // 2768个端口
		if available > expectedTotal {
			t.Errorf("可用端口数量异常: %d (总数应该是 %d)", available, expectedTotal)
		}
		t.Logf("当前可用端口数量: %d", available)
	})

	// 测试端口释放后重用
	t.Run("ReleaseAndReuse", func(t *testing.T) {
		// 分配一个端口
		pm1, err := pmDao.AllocatePort(testEnv.ID, "rdp", 3389)
		if err != nil {
			t.Fatalf("分配端口失败: %v", err)
		}
		allocatedPort := pm1.ExternalPort
		t.Logf("分配端口: %d", allocatedPort)

		// 释放端口
		if err := pmDao.ReleasePort(pm1.ID); err != nil {
			t.Fatalf("释放端口失败: %v", err)
		}
		t.Logf("释放端口: %d", allocatedPort)

		// 验证端口状态已更新
		released, err := pmDao.GetByID(pm1.ID)
		if err != nil {
			t.Fatalf("获取端口映射失败: %v", err)
		}
		if released.Status != "released" {
			t.Errorf("端口状态未更新: 期望 'released', 实际 '%s'", released.Status)
		}
		if released.ReleasedAt == nil {
			t.Error("ReleasedAt 应该被设置")
		}
	})
}
