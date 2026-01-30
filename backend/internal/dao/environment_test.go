package dao

import (
	"sync"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
)

func TestEnvironmentDao_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Environment{}, &entity.PortMapping{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	envDao := NewEnvironmentDao()
	customerDao := NewCustomerDao()
	workspaceDao := NewWorkspaceDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.Customer{
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
		CustomerID:  testCustomer.ID,
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

	// 测试 GetByCustomerID
	envs, err := envDao.GetByCustomerID(testCustomer.ID)
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
	customerDao := NewCustomerDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.Customer{
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
		CustomerID: testCustomer.ID,
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
	customerDao := NewCustomerDao()
	hostDao := NewHostDao()

	// 创建测试用户
	testCustomer := &entity.Customer{
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
		CustomerID: testCustomer.ID,
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
