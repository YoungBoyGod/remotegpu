package service

import (
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
)

func TestWorkspaceService_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewCustomerDao()

	// 创建测试用户
	testCustomer := &entity.Customer{
		Username:     "test-ws-service-" + time.Now().Format("20060102150405"),
		Email:        "ws-service-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 测试 CreateWorkspace
	workspace := &entity.Workspace{
		OwnerID:     testCustomer.ID,
		Name:        "Test Service Workspace",
		Description: "Test workspace for service",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	t.Logf("创建工作空间成功, ID: %d, UUID: %s", workspace.ID, workspace.UUID)

	// 验证默认值
	if workspace.Type != "personal" {
		t.Errorf("默认类型应该是 'personal', 实际: '%s'", workspace.Type)
	}
	if workspace.Status != "active" {
		t.Errorf("默认状态应该是 'active', 实际: '%s'", workspace.Status)
	}
	if workspace.MemberCount != 1 {
		t.Errorf("默认成员数应该是 1, 实际: %d", workspace.MemberCount)
	}

	// 测试 GetWorkspace
	found, err := service.GetWorkspace(workspace.ID)
	if err != nil {
		t.Fatalf("获取工作空间失败: %v", err)
	}
	if found.Name != "Test Service Workspace" {
		t.Errorf("工作空间名称不匹配")
	}

	// 测试 UpdateWorkspace
	workspace.Name = "Updated Service Workspace"
	workspace.Description = "Updated description"
	if err := service.UpdateWorkspace(workspace); err != nil {
		t.Fatalf("更新工作空间失败: %v", err)
	}

	updated, err := service.GetWorkspace(workspace.ID)
	if err != nil {
		t.Fatalf("获取更新后的工作空间失败: %v", err)
	}
	if updated.Name != "Updated Service Workspace" {
		t.Errorf("工作空间名称未更新")
	}

	// 测试 ListWorkspaces
	workspaces, total, err := service.ListWorkspaces(testCustomer.ID, 1, 10)
	if err != nil {
		t.Fatalf("获取工作空间列表失败: %v", err)
	}
	if len(workspaces) == 0 {
		t.Error("应该至少有一个工作空间")
	}
	t.Logf("工作空间列表: 总数 %d, 当前页 %d 条", total, len(workspaces))

	// 测试 DeleteWorkspace
	if err := service.DeleteWorkspace(workspace.ID); err != nil {
		t.Fatalf("删除工作空间失败: %v", err)
	}
	t.Logf("删除工作空间成功")
}

func TestWorkspaceService_MemberManagement(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewCustomerDao()

	// 创建测试用户（所有者）
	owner := &entity.Customer{
		Username:     "test-ws-owner-" + time.Now().Format("20060102150405"),
		Email:        "ws-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	// 创建工作空间
	workspace := &entity.Workspace{
		UUID:    uuid.New(),
		OwnerID: owner.ID,
		Name:    "Test Workspace for Members",
		Type:    "team",
		Status:  "active",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 创建测试成员用户
	member1 := &entity.Customer{
		Username:     "test-member1-" + time.Now().Format("20060102150405"),
		Email:        "member1-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(member1); err != nil {
		t.Fatalf("创建成员用户失败: %v", err)
	}
	defer customerDao.Delete(member1.ID)

	// 测试 AddMember
	if err := service.AddMember(workspace.ID, member1.ID, "member"); err != nil {
		t.Fatalf("添加成员失败: %v", err)
	}
	t.Logf("添加成员成功")

	// 验证成员数量更新
	ws, err := service.GetWorkspace(workspace.ID)
	if err != nil {
		t.Fatalf("获取工作空间失败: %v", err)
	}
	if ws.MemberCount != 2 {
		t.Errorf("成员数量应该是 2, 实际: %d", ws.MemberCount)
	}

	// 测试重复添加成员
	if err := service.AddMember(workspace.ID, member1.ID, "member"); err == nil {
		t.Error("重复添加成员应该失败")
	}

	// 测试 ListMembers
	members, err := service.ListMembers(workspace.ID)
	if err != nil {
		t.Fatalf("获取成员列表失败: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("成员数量应该是 1, 实际: %d", len(members))
	}

	// 测试 CheckPermission - 所有者
	hasPermission, err := service.CheckPermission(workspace.ID, owner.ID)
	if err != nil {
		t.Fatalf("检查所有者权限失败: %v", err)
	}
	if !hasPermission {
		t.Error("所有者应该有权限")
	}

	// 测试 CheckPermission - 成员
	hasPermission, err = service.CheckPermission(workspace.ID, member1.ID)
	if err != nil {
		t.Fatalf("检查成员权限失败: %v", err)
	}
	if !hasPermission {
		t.Error("成员应该有权限")
	}

	// 创建另一个用户（非成员）
	nonMember := &entity.Customer{
		Username:     "test-nonmember-" + time.Now().Format("20060102150405"),
		Email:        "nonmember-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(nonMember); err != nil {
		t.Fatalf("创建非成员用户失败: %v", err)
	}
	defer customerDao.Delete(nonMember.ID)

	// 测试 CheckPermission - 非成员
	hasPermission, err = service.CheckPermission(workspace.ID, nonMember.ID)
	if err != nil {
		t.Fatalf("检查非成员权限失败: %v", err)
	}
	if hasPermission {
		t.Error("非成员不应该有权限")
	}

	// 测试 RemoveMember
	if err := service.RemoveMember(workspace.ID, member1.ID); err != nil {
		t.Fatalf("移除成员失败: %v", err)
	}
	t.Logf("移除成员成功")

	// 验证成员数量更新
	ws, err = service.GetWorkspace(workspace.ID)
	if err != nil {
		t.Fatalf("获取工作空间失败: %v", err)
	}
	if ws.MemberCount != 1 {
		t.Errorf("成员数量应该是 1, 实际: %d", ws.MemberCount)
	}

	// 测试移除所有者（应该失败）
	if err := service.RemoveMember(workspace.ID, owner.ID); err == nil {
		t.Error("移除所有者应该失败")
	}
}
