package dao

import (
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
)

func TestWorkspaceDao_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	workspaceDao := NewWorkspaceDao()
	customerDao := NewCustomerDao()

	// 先创建测试用户
	testCustomer := &entity.Customer{
		Username:     "test-workspace-owner-" + time.Now().Format("20060102150405"),
		Email:        "workspace-test-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 测试 Create Workspace
	workspace := &entity.Workspace{
		UUID:        uuid.New(),
		OwnerID:     testCustomer.ID,
		Name:        "Test Workspace",
		Description: "This is a test workspace",
		Type:        "personal",
		MemberCount: 1,
		Status:      "active",
	}
	if err := workspaceDao.Create(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	t.Logf("创建工作空间成功, ID: %d", workspace.ID)

	// 测试 GetByID
	found, err := workspaceDao.GetByID(workspace.ID)
	if err != nil {
		t.Fatalf("获取工作空间失败: %v", err)
	}
	if found.Name != "Test Workspace" {
		t.Errorf("工作空间名称不匹配: 期望 'Test Workspace', 实际 '%s'", found.Name)
	}

	// 测试 GetByUUID
	foundByUUID, err := workspaceDao.GetByUUID(workspace.UUID.String())
	if err != nil {
		t.Fatalf("根据UUID获取工作空间失败: %v", err)
	}
	if foundByUUID.ID != workspace.ID {
		t.Errorf("工作空间ID不匹配: 期望 %d, 实际 %d", workspace.ID, foundByUUID.ID)
	}

	// 测试 Update
	workspace.Name = "Updated Workspace"
	workspace.Description = "Updated description"
	if err := workspaceDao.Update(workspace); err != nil {
		t.Fatalf("更新工作空间失败: %v", err)
	}

	updated, err := workspaceDao.GetByID(workspace.ID)
	if err != nil {
		t.Fatalf("获取更新后的工作空间失败: %v", err)
	}
	if updated.Name != "Updated Workspace" {
		t.Errorf("工作空间名称未更新: 期望 'Updated Workspace', 实际 '%s'", updated.Name)
	}

	// 测试 GetByOwnerID
	workspaces, err := workspaceDao.GetByOwnerID(testCustomer.ID)
	if err != nil {
		t.Fatalf("根据所有者ID获取工作空间失败: %v", err)
	}
	if len(workspaces) == 0 {
		t.Error("应该至少有一个工作空间")
	}

	// 测试 GetByStatus
	activeWorkspaces, err := workspaceDao.GetByStatus("active")
	if err != nil {
		t.Fatalf("根据状态获取工作空间失败: %v", err)
	}
	if len(activeWorkspaces) == 0 {
		t.Error("应该至少有一个活跃的工作空间")
	}

	// 测试 List
	workspaceList, total, err := workspaceDao.List(1, 10)
	if err != nil {
		t.Fatalf("获取工作空间列表失败: %v", err)
	}
	if total == 0 {
		t.Error("工作空间总数应该大于0")
	}
	t.Logf("工作空间列表: 总数 %d, 当前页 %d 条", total, len(workspaceList))

	// 测试 Delete
	if err := workspaceDao.Delete(workspace.ID); err != nil {
		t.Fatalf("删除工作空间失败: %v", err)
	}
	t.Logf("删除工作空间成功")
}

func TestWorkspaceMemberDao_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	workspaceDao := NewWorkspaceDao()
	memberDao := NewWorkspaceMemberDao()
	customerDao := NewCustomerDao()

	// 创建测试用户（工作空间所有者）
	owner := &entity.Customer{
		Username:     "test-member-owner-" + time.Now().Format("20060102150405"),
		Email:        "member-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	// 创建测试工作空间
	workspace := &entity.Workspace{
		UUID:        uuid.New(),
		OwnerID:     owner.ID,
		Name:        "Test Workspace for Members",
		Type:        "team",
		MemberCount: 1,
		Status:      "active",
	}
	if err := workspaceDao.Create(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer workspaceDao.Delete(workspace.ID)

	// 创建测试成员用户
	memberUser := &entity.Customer{
		Username:     "test-member-user-" + time.Now().Format("20060102150405"),
		Email:        "member-user-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(memberUser); err != nil {
		t.Fatalf("创建成员用户失败: %v", err)
	}
	defer customerDao.Delete(memberUser.ID)

	// 测试 Create Member
	member := &entity.WorkspaceMember{
		WorkspaceID: workspace.ID,
		CustomerID:  memberUser.ID,
		Role:        "member",
		Status:      "active",
	}
	if err := memberDao.Create(member); err != nil {
		t.Fatalf("创建工作空间成员失败: %v", err)
	}
	t.Logf("创建工作空间成员成功, ID: %d", member.ID)

	// 测试 GetByID
	found, err := memberDao.GetByID(member.ID)
	if err != nil {
		t.Fatalf("获取工作空间成员失败: %v", err)
	}
	if found.Role != "member" {
		t.Errorf("成员角色不匹配: 期望 'member', 实际 '%s'", found.Role)
	}

	// 测试 Update
	member.Role = "admin"
	if err := memberDao.Update(member); err != nil {
		t.Fatalf("更新工作空间成员失败: %v", err)
	}

	updated, err := memberDao.GetByID(member.ID)
	if err != nil {
		t.Fatalf("获取更新后的成员失败: %v", err)
	}
	if updated.Role != "admin" {
		t.Errorf("成员角色未更新: 期望 'admin', 实际 '%s'", updated.Role)
	}

	// 测试 GetByWorkspaceID
	members, err := memberDao.GetByWorkspaceID(workspace.ID)
	if err != nil {
		t.Fatalf("根据工作空间ID获取成员失败: %v", err)
	}
	if len(members) == 0 {
		t.Error("应该至少有一个成员")
	}

	// 测试 GetByCustomerID
	customerMembers, err := memberDao.GetByCustomerID(memberUser.ID)
	if err != nil {
		t.Fatalf("根据客户ID获取成员失败: %v", err)
	}
	if len(customerMembers) == 0 {
		t.Error("应该至少有一个成员记录")
	}

	// 测试 GetByWorkspaceAndCustomer
	foundMember, err := memberDao.GetByWorkspaceAndCustomer(workspace.ID, memberUser.ID)
	if err != nil {
		t.Fatalf("根据工作空间和客户ID获取成员失败: %v", err)
	}
	if foundMember.ID != member.ID {
		t.Errorf("成员ID不匹配: 期望 %d, 实际 %d", member.ID, foundMember.ID)
	}

	// 测试 Delete
	if err := memberDao.Delete(member.ID); err != nil {
		t.Fatalf("删除工作空间成员失败: %v", err)
	}
	t.Logf("删除工作空间成员成功")
}
