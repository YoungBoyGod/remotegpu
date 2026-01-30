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
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
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
		CustomerID:  customer.ID,
		WorkspaceID: &workspace.ID,
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     500,
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
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     500,
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
