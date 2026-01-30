package dao

import (
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
)

func setupResourceQuotaTest(t *testing.T) {
	setupTestDB(t)
	db := database.GetDB()
	// 自动迁移表结构
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.ResourceQuota{}); err != nil {
		t.Fatalf("自动迁移表结构失败: %v", err)
	}
}

func TestResourceQuotaDao_CRUD(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()
	workspaceDao := NewWorkspaceDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-quota-user-" + time.Now().Format("20060102150405"),
		Email:        "quota-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	if err != nil {
		t.Fatalf("创建测试客户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 创建测试工作空间
	workspace := &entity.Workspace{
		UUID:        uuid.New(),
		OwnerID:     customer.ID,
		Name:        "Test Quota Workspace",
		Description: "Test workspace for quota",
		Type:        "personal",
		Status:      "active",
	}
	err = workspaceDao.Create(workspace)
	if err != nil {
		t.Fatalf("创建测试工作空间失败: %v", err)
	}
	defer workspaceDao.Delete(workspace.ID)

	// 测试 Create - 用户级别配额
	userQuota := &entity.ResourceQuota{
		CustomerID:       customer.ID,
		WorkspaceID:      nil,
		QuotaLevel:       "pro",
		CPU:              16,
		Memory:           32768,
		GPU:              4,
		Storage:          1000,
		EnvironmentQuota: 10,
	}
	err = dao.Create(userQuota)
	if err != nil {
		t.Fatalf("创建用户级别配额失败: %v", err)
	}
	t.Log("创建用户级别配额成功")

	// 测试 GetByID
	found, err := dao.GetByID(userQuota.ID)
	if err != nil {
		t.Fatalf("获取配额失败: %v", err)
	}
	if found.CPU != 16 {
		t.Fatalf("CPU配额不匹配: got %d, want 16", found.CPU)
	}
	t.Log("获取配额成功")

	// 测试 GetByCustomerID
	customerQuota, err := dao.GetByCustomerID(customer.ID)
	if err != nil {
		t.Fatalf("根据客户ID获取配额失败: %v", err)
	}
	if customerQuota.ID != userQuota.ID {
		t.Fatalf("配额ID不匹配")
	}
	t.Log("根据客户ID获取配额成功")

	// 测试 Update
	found.CPU = 32
	found.Memory = 65536
	err = dao.Update(found)
	if err != nil {
		t.Fatalf("更新配额失败: %v", err)
	}
	updated, _ := dao.GetByID(found.ID)
	if updated.CPU != 32 {
		t.Fatalf("更新后CPU配额不匹配: got %d, want 32", updated.CPU)
	}
	t.Log("更新配额成功")

	// 测试 Delete
	err = dao.Delete(userQuota.ID)
	if err != nil {
		t.Fatalf("删除配额失败: %v", err)
	}
	t.Log("删除配额成功")

	t.Log("ResourceQuota DAO CRUD 测试通过")
}

func TestResourceQuotaDao_WorkspaceQuota(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()
	workspaceDao := NewWorkspaceDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-ws-quota-user-" + time.Now().Format("20060102150405"),
		Email:        "ws-quota-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test WS Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	if err != nil {
		t.Fatalf("创建测试客户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 创建测试工作空间
	workspace := &entity.Workspace{
		UUID:        uuid.New(),
		OwnerID:     customer.ID,
		Name:        "Test WS Quota Workspace",
		Description: "Test workspace for quota",
		Type:        "personal",
		Status:      "active",
	}
	err = workspaceDao.Create(workspace)
	if err != nil {
		t.Fatalf("创建测试工作空间失败: %v", err)
	}
	defer workspaceDao.Delete(workspace.ID)

	// 创建工作空间级别配额
	wsQuota := &entity.ResourceQuota{
		CustomerID:       customer.ID,
		WorkspaceID:      &workspace.ID,
		QuotaLevel:       "basic",
		CPU:              8,
		Memory:           16384,
		GPU:              2,
		Storage:          500,
		EnvironmentQuota: 5,
	}
	err = dao.Create(wsQuota)
	if err != nil {
		t.Fatalf("创建工作空间级别配额失败: %v", err)
	}
	defer dao.Delete(wsQuota.ID)

	// 测试 GetByWorkspaceID
	found, err := dao.GetByWorkspaceID(workspace.ID)
	if err != nil {
		t.Fatalf("根据工作空间ID获取配额失败: %v", err)
	}
	if found.CPU != 8 {
		t.Fatalf("CPU配额不匹配: got %d, want 8", found.CPU)
	}
	t.Log("根据工作空间ID获取配额成功")

	// 测试 GetByCustomerAndWorkspace
	found2, err := dao.GetByCustomerAndWorkspace(customer.ID, workspace.ID)
	if err != nil {
		t.Fatalf("根据客户和工作空间ID获取配额失败: %v", err)
	}
	if found2.ID != wsQuota.ID {
		t.Fatalf("配额ID不匹配")
	}
	t.Log("根据客户和工作空间ID获取配额成功")

	t.Log("工作空间配额测试通过")
}

func TestResourceQuotaDao_List(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-list-quota-user-" + time.Now().Format("20060102150405"),
		Email:        "list-quota-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test List Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	if err != nil {
		t.Fatalf("创建测试客户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 创建多个配额记录
	quota1 := &entity.ResourceQuota{
		CustomerID:       customer.ID,
		WorkspaceID:      nil,
		QuotaLevel:       "free",
		CPU:              8,
		Memory:           16384,
		GPU:              2,
		Storage:          500,
		EnvironmentQuota: 3,
	}
	err = dao.Create(quota1)
	if err != nil {
		t.Fatalf("创建配额1失败: %v", err)
	}
	defer dao.Delete(quota1.ID)

	// 测试 List
	quotas, total, err := dao.List(1, 10)
	if err != nil {
		t.Fatalf("获取配额列表失败: %v", err)
	}
	if total < 1 {
		t.Fatalf("配额总数不正确: got %d, want >= 1", total)
	}
	if len(quotas) < 1 {
		t.Fatalf("配额列表为空")
	}
	t.Logf("获取配额列表成功，总数: %d", total)

	t.Log("配额列表测试通过")
}

// TestResourceQuotaDao_ErrorScenarios 测试错误场景
func TestResourceQuotaDao_ErrorScenarios(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()

	// 测试 GetByID - 不存在的ID
	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := dao.GetByID(99999)
		if err == nil {
			t.Error("获取不存在的配额应该返回错误")
		}
		t.Logf("GetByID 不存在的ID测试通过: %v", err)
	})

	// 测试 GetByCustomerID - 不存在的客户ID
	t.Run("GetByCustomerID_NotFound", func(t *testing.T) {
		_, err := dao.GetByCustomerID(99999)
		if err == nil {
			t.Error("获取不存在客户的配额应该返回错误")
		}
		t.Logf("GetByCustomerID 不存在的客户ID测试通过: %v", err)
	})

	// 测试 GetByWorkspaceID - 不存在的工作空间ID
	t.Run("GetByWorkspaceID_NotFound", func(t *testing.T) {
		_, err := dao.GetByWorkspaceID(99999)
		if err == nil {
			t.Error("获取不存在工作空间的配额应该返回错误")
		}
		t.Logf("GetByWorkspaceID 不存在的工作空间ID测试通过: %v", err)
	})

	// 测试 GetByCustomerAndWorkspace - 不存在的组合
	t.Run("GetByCustomerAndWorkspace_NotFound", func(t *testing.T) {
		_, err := dao.GetByCustomerAndWorkspace(99999, 99999)
		if err == nil {
			t.Error("获取不存在的客户和工作空间组合配额应该返回错误")
		}
		t.Logf("GetByCustomerAndWorkspace 不存在的组合测试通过: %v", err)
	})

	t.Log("错误场景测试通过")
}

// TestResourceQuotaDao_ListEdgeCases 测试List方法的边界情况
func TestResourceQuotaDao_ListEdgeCases(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-list-edge-" + time.Now().Format("20060102150405"),
		Email:        "list-edge-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test List Edge User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	if err != nil {
		t.Fatalf("创建测试客户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 创建测试配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     500,
	}
	err = dao.Create(quota)
	if err != nil {
		t.Fatalf("创建配额失败: %v", err)
	}
	defer dao.Delete(quota.ID)

	// 测试 page < 1（应该默认为1）
	t.Run("List_PageLessThanOne", func(t *testing.T) {
		quotas, total, err := dao.List(0, 10)
		if err != nil {
			t.Errorf("page < 1 应该使用默认值: %v", err)
		}
		if total < 1 {
			t.Errorf("应该至少有一条记录")
		}
		if len(quotas) < 1 {
			t.Errorf("应该返回至少一条记录")
		}
		t.Logf("page < 1 测试通过，返回 %d 条记录", len(quotas))
	})

	// 测试 pageSize < 1（应该默认为10）
	t.Run("List_PageSizeLessThanOne", func(t *testing.T) {
		quotas, total, err := dao.List(1, 0)
		if err != nil {
			t.Errorf("pageSize < 1 应该使用默认值: %v", err)
		}
		if total < 1 {
			t.Errorf("应该至少有一条记录")
		}
		t.Logf("pageSize < 1 测试通过，返回 %d 条记录", len(quotas))
	})

	// 测试 pageSize > 100（应该限制为100）
	t.Run("List_PageSizeGreaterThan100", func(t *testing.T) {
		quotas, total, err := dao.List(1, 200)
		if err != nil {
			t.Errorf("pageSize > 100 应该限制为100: %v", err)
		}
		if total < 1 {
			t.Errorf("应该至少有一条记录")
		}
		// 验证返回的记录数不超过100
		if len(quotas) > 100 {
			t.Errorf("返回的记录数不应该超过100，实际: %d", len(quotas))
		}
		t.Logf("pageSize > 100 测试通过，返回 %d 条记录", len(quotas))
	})

	// 测试负数参数
	t.Run("List_NegativeParameters", func(t *testing.T) {
		quotas, total, err := dao.List(-1, -10)
		if err != nil {
			t.Errorf("负数参数应该使用默认值: %v", err)
		}
		if total < 1 {
			t.Errorf("应该至少有一条记录")
		}
		t.Logf("负数参数测试通过，返回 %d 条记录", len(quotas))
	})

	t.Log("List边界情况测试通过")
}

func TestResourceQuotaDao_GetByQuotaLevel(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()

	// 创建测试客户
	customer1 := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-level-user1-" + time.Now().Format("20060102150405"),
		Email:        "level-test1-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Level User 1",
		Status:       "active",
	}
	err := customerDao.Create(customer1)
	if err != nil {
		t.Fatalf("创建测试客户1失败: %v", err)
	}
	defer customerDao.Delete(customer1.ID)

	customer2 := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-level-user2-" + time.Now().Format("20060102150405"),
		Email:        "level-test2-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Level User 2",
		Status:       "active",
	}
	err = customerDao.Create(customer2)
	if err != nil {
		t.Fatalf("创建测试客户2失败: %v", err)
	}
	defer customerDao.Delete(customer2.ID)

	// 创建不同级别的配额
	quotaPro := &entity.ResourceQuota{
		CustomerID:       customer1.ID,
		WorkspaceID:      nil,
		QuotaLevel:       "pro",
		CPU:              16,
		Memory:           32768,
		GPU:              4,
		Storage:          1000,
		EnvironmentQuota: 10,
	}
	err = dao.Create(quotaPro)
	if err != nil {
		t.Fatalf("创建pro配额失败: %v", err)
	}
	defer dao.Delete(quotaPro.ID)

	quotaBasic := &entity.ResourceQuota{
		CustomerID:       customer2.ID,
		WorkspaceID:      nil,
		QuotaLevel:       "basic",
		CPU:              8,
		Memory:           16384,
		GPU:              2,
		Storage:          500,
		EnvironmentQuota: 5,
	}
	err = dao.Create(quotaBasic)
	if err != nil {
		t.Fatalf("创建basic配额失败: %v", err)
	}
	defer dao.Delete(quotaBasic.ID)

	// 测试查询pro级别
	proQuotas, err := dao.GetByQuotaLevel("pro")
	if err != nil {
		t.Fatalf("查询pro级别配额失败: %v", err)
	}
	if len(proQuotas) < 1 {
		t.Fatalf("应该至少有1个pro级别配额")
	}
	t.Logf("查询到 %d 个pro级别配额", len(proQuotas))

	// 测试查询basic级别
	basicQuotas, err := dao.GetByQuotaLevel("basic")
	if err != nil {
		t.Fatalf("查询basic级别配额失败: %v", err)
	}
	if len(basicQuotas) < 1 {
		t.Fatalf("应该至少有1个basic级别配额")
	}
	t.Logf("查询到 %d 个basic级别配额", len(basicQuotas))

	t.Log("GetByQuotaLevel测试通过")
}

func TestResourceQuotaDao_CheckQuota(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()

	// 自动迁移Host和Environment表
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Host{}, &entity.Environment{}); err != nil {
		t.Fatalf("自动迁移表失败: %v", err)
	}

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-check-user-" + time.Now().Format("20060102150405"),
		Email:        "check-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Check User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	if err != nil {
		t.Fatalf("创建测试客户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 创建配额：CPU=16, Memory=32768, GPU=4, Storage=1000
	quota := &entity.ResourceQuota{
		CustomerID:       customer.ID,
		WorkspaceID:      nil,
		QuotaLevel:       "pro",
		CPU:              16,
		Memory:           32768,
		GPU:              4,
		Storage:          1000,
		EnvironmentQuota: 10,
	}
	err = dao.Create(quota)
	if err != nil {
		t.Fatalf("创建配额失败: %v", err)
	}
	defer dao.Delete(quota.ID)

	// 创建测试主机
	host := &entity.Host{
		ID:             "test-host-" + uuid.New().String(),
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "docker",
		TotalCPU:       32,
		TotalMemory:    65536,
		TotalGPU:       8,
	}
	err = db.Create(host).Error
	if err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer db.Delete(host)

	// 创建一个运行中的环境：CPU=8, Memory=16384, GPU=2, Storage=500
	env := &entity.Environment{
		ID:          "test-env-" + uuid.New().String(),
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		HostID:      host.ID,
		Name:        "Test Environment",
		Image:       "ubuntu:20.04",
		Status:      "running",
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     func() *int64 { s := int64(500); return &s }(),
	}
	err = db.Create(env).Error
	if err != nil {
		t.Fatalf("创建测试环境失败: %v", err)
	}
	defer db.Delete(env)

	// 测试1：配额充足的情况（请求：CPU=4, Memory=8192, GPU=1, Storage=200）
	// 已用：CPU=8, Memory=16384, GPU=2, Storage=500
	// 总配额：CPU=16, Memory=32768, GPU=4, Storage=1000
	// 已用+请求：CPU=12, Memory=24576, GPU=3, Storage=700 < 总配额
	ok, err := dao.CheckQuota(customer.ID, 4, 8192, 1, 200)
	if err != nil {
		t.Fatalf("CheckQuota失败: %v", err)
	}
	if !ok {
		t.Fatalf("配额应该充足，但返回false")
	}
	t.Log("配额充足测试通过")

	// 测试2：CPU配额不足（请求：CPU=10, Memory=8192, GPU=1, Storage=200）
	// 已用+请求：CPU=18 > 16（总配额）
	ok, err = dao.CheckQuota(customer.ID, 10, 8192, 1, 200)
	if err != nil {
		t.Fatalf("CheckQuota失败: %v", err)
	}
	if ok {
		t.Fatalf("CPU配额应该不足，但返回true")
	}
	t.Log("CPU配额不足测试通过")

	// 测试3：内存配额不足（请求：CPU=4, Memory=20000, GPU=1, Storage=200）
	// 已用+请求：Memory=36384 > 32768（总配额）
	ok, err = dao.CheckQuota(customer.ID, 4, 20000, 1, 200)
	if err != nil {
		t.Fatalf("CheckQuota失败: %v", err)
	}
	if ok {
		t.Fatalf("内存配额应该不足，但返回true")
	}
	t.Log("内存配额不足测试通过")

	// 测试4：GPU配额不足（请求：CPU=4, Memory=8192, GPU=3, Storage=200）
	// 已用+请求：GPU=5 > 4（总配额）
	ok, err = dao.CheckQuota(customer.ID, 4, 8192, 3, 200)
	if err != nil {
		t.Fatalf("CheckQuota失败: %v", err)
	}
	if ok {
		t.Fatalf("GPU配额应该不足，但返回true")
	}
	t.Log("GPU配额不足测试通过")

	// 测试5：存储配额不足（请求：CPU=4, Memory=8192, GPU=1, Storage=600）
	// 已用+请求：Storage=1100 > 1000（总配额）
	ok, err = dao.CheckQuota(customer.ID, 4, 8192, 1, 600)
	if err != nil {
		t.Fatalf("CheckQuota失败: %v", err)
	}
	if ok {
		t.Fatalf("存储配额应该不足，但返回true")
	}
	t.Log("存储配额不足测试通过")

	// 测试6：边界情况（请求刚好用完所有配额）
	// 请求：CPU=8, Memory=16384, GPU=2, Storage=500
	// 已用+请求：CPU=16, Memory=32768, GPU=4, Storage=1000（刚好等于总配额）
	ok, err = dao.CheckQuota(customer.ID, 8, 16384, 2, 500)
	if err != nil {
		t.Fatalf("CheckQuota失败: %v", err)
	}
	if !ok {
		t.Fatalf("边界情况应该允许，但返回false")
	}
	t.Log("边界情况测试通过")

	t.Log("CheckQuota测试通过")
}

func TestResourceQuotaDao_GetAvailableQuota(t *testing.T) {
	setupResourceQuotaTest(t)

	dao := NewResourceQuotaDao()
	customerDao := NewCustomerDao()

	// 自动迁移Host和Environment表
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Host{}, &entity.Environment{}); err != nil {
		t.Fatalf("自动迁移表失败: %v", err)
	}

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-available-user-" + time.Now().Format("20060102150405"),
		Email:        "available-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Available User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	if err != nil {
		t.Fatalf("创建测试客户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 创建配额：CPU=16, Memory=32768, GPU=4, Storage=1000
	quota := &entity.ResourceQuota{
		CustomerID:       customer.ID,
		WorkspaceID:      nil,
		QuotaLevel:       "pro",
		CPU:              16,
		Memory:           32768,
		GPU:              4,
		Storage:          1000,
		EnvironmentQuota: 10,
	}
	err = dao.Create(quota)
	if err != nil {
		t.Fatalf("创建配额失败: %v", err)
	}
	defer dao.Delete(quota.ID)

	// 创建测试主机
	host := &entity.Host{
		ID:             "test-host-" + uuid.New().String(),
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "docker",
		TotalCPU:       32,
		TotalMemory:    65536,
		TotalGPU:       8,
	}
	err = db.Create(host).Error
	if err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer db.Delete(host)

	// 创建一个运行中的环境：CPU=8, Memory=16384, GPU=2, Storage=500
	env := &entity.Environment{
		ID:          "test-env-" + uuid.New().String(),
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		HostID:      host.ID,
		Name:        "Test Environment",
		Image:       "ubuntu:20.04",
		Status:      "running",
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     func() *int64 { s := int64(500); return &s }(),
	}
	err = db.Create(env).Error
	if err != nil {
		t.Fatalf("创建测试环境失败: %v", err)
	}
	defer db.Delete(env)

	// 测试GetAvailableQuota
	available, err := dao.GetAvailableQuota(customer.ID)
	if err != nil {
		t.Fatalf("GetAvailableQuota失败: %v", err)
	}

	// 验证可用配额 = 总配额 - 已用配额
	// 总配额：CPU=16, Memory=32768, GPU=4, Storage=1000
	// 已用：CPU=8, Memory=16384, GPU=2, Storage=500
	// 可用：CPU=8, Memory=16384, GPU=2, Storage=500
	if available.CPU != 8 {
		t.Errorf("可用CPU配额不正确: got %d, want 8", available.CPU)
	}
	if available.Memory != 16384 {
		t.Errorf("可用内存配额不正确: got %d, want 16384", available.Memory)
	}
	if available.GPU != 2 {
		t.Errorf("可用GPU配额不正确: got %d, want 2", available.GPU)
	}
	if available.Storage != 500 {
		t.Errorf("可用存储配额不正确: got %d, want 500", available.Storage)
	}

	t.Log("GetAvailableQuota测试通过")
}
