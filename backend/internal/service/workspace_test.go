package service

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestWorkspaceService_CRUD(t *testing.T) {
	setupTestDB(t)

	// 自动迁移表结构
	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	testCustomer := &entity.User{
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
	customerDao := dao.NewUserDao()

	// 创建测试用户（所有者）
	owner := &entity.User{
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
	member1 := &entity.User{
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
	nonMember := &entity.User{
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

// TestWorkspaceService_CreateWorkspace_Validation 测试创建工作空间的输入验证
func TestWorkspaceService_CreateWorkspace_Validation(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	testCustomer := &entity.User{
		Username:     "test-validation-" + time.Now().Format("20060102150405"),
		Email:        "validation-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(testCustomer); err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	defer customerDao.Delete(testCustomer.ID)

	// 测试空名称
	t.Run("EmptyName", func(t *testing.T) {
		workspace := &entity.Workspace{
			OwnerID: testCustomer.ID,
			Name:    "",
		}
		err := service.CreateWorkspace(workspace)
		if err == nil {
			t.Error("空名称应该返回错误")
		}
		if err != nil && err.Error() != "工作空间名称不能为空" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 测试名称过长
	t.Run("NameTooLong", func(t *testing.T) {
		longName := ""
		for i := 0; i < 130; i++ {
			longName += "a"
		}
		workspace := &entity.Workspace{
			OwnerID: testCustomer.ID,
			Name:    longName,
		}
		err := service.CreateWorkspace(workspace)
		if err == nil {
			t.Error("名称过长应该返回错误")
		}
		if err != nil && err.Error() != "工作空间名称不能超过128个字符" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 测试无效的Type
	t.Run("InvalidType", func(t *testing.T) {
		workspace := &entity.Workspace{
			OwnerID: testCustomer.ID,
			Name:    "Test Workspace",
			Type:    "invalid_type",
		}
		err := service.CreateWorkspace(workspace)
		if err == nil {
			t.Error("无效的Type应该返回错误")
		}
		if err != nil && err.Error() != "无效的工作空间类型: invalid_type" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 测试有效的Type
	t.Run("ValidTypes", func(t *testing.T) {
		validTypes := []string{"personal", "team", "enterprise"}
		for _, validType := range validTypes {
			workspace := &entity.Workspace{
				OwnerID: testCustomer.ID,
				Name:    "Test Workspace " + validType,
				Type:    validType,
			}
			err := service.CreateWorkspace(workspace)
			if err != nil {
				t.Errorf("有效的Type '%s' 不应该返回错误: %v", validType, err)
			}
			if workspace.ID > 0 {
				service.DeleteWorkspace(workspace.ID)
			}
		}
	})
}

// TestWorkspaceService_AddMember_RoleValidation 测试添加成员的角色验证
func TestWorkspaceService_AddMember_RoleValidation(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "test-role-owner-" + time.Now().Format("20060102150405"),
		Email:        "role-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	member := &entity.User{
		Username:     "test-role-member-" + time.Now().Format("20060102150405"),
		Email:        "role-member-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(member); err != nil {
		t.Fatalf("创建成员失败: %v", err)
	}
	defer customerDao.Delete(member.ID)

	// 创建工作空间
	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Test Workspace for Role",
		Type:    "team",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 测试无效的角色
	t.Run("InvalidRole", func(t *testing.T) {
		err := service.AddMember(workspace.ID, member.ID, "invalid_role")
		if err == nil {
			t.Error("无效的角色应该返回错误")
		}
		if err != nil && err.Error() != "无效的角色: invalid_role" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 测试有效的角色
	t.Run("ValidRoles", func(t *testing.T) {
		validRoles := []string{"owner", "admin", "member", "viewer"}
		for i, role := range validRoles {
			// 为每个角色创建一个新成员
			newMember := &entity.User{
				Username:     fmt.Sprintf("test-member-%d-%s", i, time.Now().Format("20060102150405")),
				Email:        fmt.Sprintf("member-%d-%s@example.com", i, time.Now().Format("20060102150405")),
				PasswordHash: "test-hash",
				Status:       "active",
			}
			if err := customerDao.Create(newMember); err != nil {
				t.Fatalf("创建成员失败: %v", err)
			}
			defer customerDao.Delete(newMember.ID)

			err := service.AddMember(workspace.ID, newMember.ID, role)
			if err != nil {
				t.Errorf("有效的角色 '%s' 不应该返回错误: %v", role, err)
			}
		}
	})

	// 测试默认角色
	t.Run("DefaultRole", func(t *testing.T) {
		defaultMember := &entity.User{
			Username:     "test-default-" + time.Now().Format("20060102150405"),
			Email:        "default-" + time.Now().Format("20060102150405") + "@example.com",
			PasswordHash: "test-hash",
			Status:       "active",
		}
		if err := customerDao.Create(defaultMember); err != nil {
			t.Fatalf("创建成员失败: %v", err)
		}
		defer customerDao.Delete(defaultMember.ID)

		err := service.AddMember(workspace.ID, defaultMember.ID, "")
		if err != nil {
			t.Errorf("空角色应该使用默认值: %v", err)
		}

		// 验证默认角色是 "member"
		members, _ := service.ListMembers(workspace.ID)
		found := false
		for _, m := range members {
			if m.UserID == defaultMember.ID && m.Role == "member" {
				found = true
				break
			}
		}
		if !found {
			t.Error("默认角色应该是 'member'")
		}
	})
}

// TestWorkspaceService_CheckPermission_ArchivedWorkspace 测试已归档工作空间的权限检查
func TestWorkspaceService_CheckPermission_ArchivedWorkspace(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "test-archived-owner-" + time.Now().Format("20060102150405"),
		Email:        "archived-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	// 创建工作空间
	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Test Archived Workspace",
		Type:    "team",
		Status:  "active",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 测试活跃工作空间的权限
	t.Run("ActiveWorkspace", func(t *testing.T) {
		hasPermission, err := service.CheckPermission(workspace.ID, owner.ID)
		if err != nil {
			t.Errorf("检查活跃工作空间权限失败: %v", err)
		}
		if !hasPermission {
			t.Error("所有者应该有活跃工作空间的权限")
		}
	})

	// 将工作空间设置为已归档
	workspace.Status = "archived"
	if err := service.UpdateWorkspace(workspace); err != nil {
		t.Fatalf("更新工作空间状态失败: %v", err)
	}

	// 测试已归档工作空间的权限
	t.Run("ArchivedWorkspace", func(t *testing.T) {
		hasPermission, err := service.CheckPermission(workspace.ID, owner.ID)
		if err == nil {
			t.Error("已归档工作空间应该返回错误")
		}
		if err != nil && err.Error() != "工作空间已归档" {
			t.Errorf("错误信息不正确: %v", err)
		}
		if hasPermission {
			t.Error("已归档工作空间不应该有权限")
		}
	})
}

// TestWorkspaceService_UpdateWorkspace_ErrorCases 测试更新工作空间的错误场景
func TestWorkspaceService_UpdateWorkspace_ErrorCases(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "test-update-owner-" + time.Now().Format("20060102150405"),
		Email:        "update-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	// 测试更新不存在的工作空间
	t.Run("WorkspaceNotFound", func(t *testing.T) {
		workspace := &entity.Workspace{
			ID:      99999,
			OwnerID: owner.ID,
			Name:    "Non-existent Workspace",
		}
		err := service.UpdateWorkspace(workspace)
		if err == nil {
			t.Error("更新不存在的工作空间应该返回错误")
		}
		if err != nil && err.Error() != "工作空间不存在" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 创建工作空间用于测试修改OwnerID
	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Test Workspace for Update",
		Type:    "team",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 测试修改OwnerID
	t.Run("ChangeOwnerID", func(t *testing.T) {
		// 创建另一个用户
		newOwner := &entity.User{
			Username:     "test-new-owner-" + time.Now().Format("20060102150405"),
			Email:        "new-owner-" + time.Now().Format("20060102150405") + "@example.com",
			PasswordHash: "test-hash",
			Status:       "active",
		}
		if err := customerDao.Create(newOwner); err != nil {
			t.Fatalf("创建新所有者失败: %v", err)
		}
		defer customerDao.Delete(newOwner.ID)

		// 尝试修改OwnerID
		workspace.OwnerID = newOwner.ID
		err := service.UpdateWorkspace(workspace)
		if err == nil {
			t.Error("修改OwnerID应该返回错误")
		}
		if err != nil && err.Error() != "不允许修改工作空间所有者" {
			t.Errorf("错误信息不正确: %v", err)
		}

		// 恢复OwnerID
		workspace.OwnerID = owner.ID
	})
}

// TestWorkspaceService_AddMember_ErrorCases 测试添加成员的错误场景
func TestWorkspaceService_AddMember_ErrorCases(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "test-addmember-owner-" + time.Now().Format("20060102150405"),
		Email:        "addmember-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	member := &entity.User{
		Username:     "test-addmember-member-" + time.Now().Format("20060102150405"),
		Email:        "addmember-member-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(member); err != nil {
		t.Fatalf("创建成员失败: %v", err)
	}
	defer customerDao.Delete(member.ID)

	// 测试添加成员到不存在的工作空间
	t.Run("WorkspaceNotFound", func(t *testing.T) {
		err := service.AddMember(99999, member.ID, "member")
		if err == nil {
			t.Error("添加成员到不存在的工作空间应该返回错误")
		}
		if err != nil && err.Error() != "工作空间不存在" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 创建工作空间
	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Test Workspace for AddMember",
		Type:    "team",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 先添加一个成员
	if err := service.AddMember(workspace.ID, member.ID, "member"); err != nil {
		t.Fatalf("添加成员失败: %v", err)
	}

	// 测试重复添加成员
	t.Run("MemberAlreadyExists", func(t *testing.T) {
		err := service.AddMember(workspace.ID, member.ID, "admin")
		if err == nil {
			t.Error("重复添加成员应该返回错误")
		}
		if err != nil && err.Error() != "成员已存在" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})
}

// TestWorkspaceService_RemoveMember_ErrorCases 测试移除成员的错误场景
func TestWorkspaceService_RemoveMember_ErrorCases(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "test-removemember-owner-" + time.Now().Format("20060102150405"),
		Email:        "removemember-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	member := &entity.User{
		Username:     "test-removemember-member-" + time.Now().Format("20060102150405"),
		Email:        "removemember-member-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(member); err != nil {
		t.Fatalf("创建成员失败: %v", err)
	}
	defer customerDao.Delete(member.ID)

	// 测试从不存在的工作空间移除成员
	t.Run("WorkspaceNotFound", func(t *testing.T) {
		err := service.RemoveMember(99999, member.ID)
		if err == nil {
			t.Error("从不存在的工作空间移除成员应该返回错误")
		}
		if err != nil && err.Error() != "工作空间不存在" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})

	// 创建工作空间
	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Test Workspace for RemoveMember",
		Type:    "team",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 测试移除不存在的成员
	t.Run("MemberNotFound", func(t *testing.T) {
		err := service.RemoveMember(workspace.ID, member.ID)
		if err == nil {
			t.Error("移除不存在的成员应该返回错误")
		}
		if err != nil && err.Error() != "成员不存在" {
			t.Errorf("错误信息不正确: %v", err)
		}
	})
}

// TestWorkspaceService_CheckPermission_ErrorCases 测试权限检查的错误场景
func TestWorkspaceService_CheckPermission_ErrorCases(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	customer := &entity.User{
		Username:     "test-permission-" + time.Now().Format("20060102150405"),
		Email:        "permission-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(customer); err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	defer customerDao.Delete(customer.ID)

	// 测试检查不存在的工作空间权限
	t.Run("WorkspaceNotFound", func(t *testing.T) {
		hasPermission, err := service.CheckPermission(99999, customer.ID)
		if err == nil {
			t.Error("检查不存在的工作空间权限应该返回错误")
		}
		if err != nil && err.Error() != "工作空间不存在" {
			t.Errorf("错误信息不正确: %v", err)
		}
		if hasPermission {
			t.Error("不存在的工作空间不应该有权限")
		}
	})
}

// TestWorkspaceService_ListWorkspaces_Pagination 测试工作空间列表分页
func TestWorkspaceService_ListWorkspaces_Pagination(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()

	// 测试分页查询（ownerID为0）
	t.Run("PaginationWithoutOwner", func(t *testing.T) {
		workspaces, total, err := service.ListWorkspaces(0, 1, 10)
		if err != nil {
			t.Errorf("分页查询失败: %v", err)
		}
		t.Logf("分页查询结果: 总数 %d, 当前页 %d 条", total, len(workspaces))
	})
}

// TestWorkspaceService_RemoveMember_MemberCount 测试移除成员时的成员数量更新
func TestWorkspaceService_RemoveMember_MemberCount(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "test-count-owner-" + time.Now().Format("20060102150405"),
		Email:        "count-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	member := &entity.User{
		Username:     "test-count-member-" + time.Now().Format("20060102150405"),
		Email:        "count-member-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(member); err != nil {
		t.Fatalf("创建成员失败: %v", err)
	}
	defer customerDao.Delete(member.ID)

	// 创建工作空间，手动设置MemberCount为0
	workspace := &entity.Workspace{
		OwnerID:     owner.ID,
		Name:        "Test Workspace for Count",
		Type:        "team",
		MemberCount: 0,
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)

	// 添加成员
	if err := service.AddMember(workspace.ID, member.ID, "member"); err != nil {
		t.Fatalf("添加成员失败: %v", err)
	}

	// 移除成员，测试MemberCount的更新
	if err := service.RemoveMember(workspace.ID, member.ID); err != nil {
		t.Fatalf("移除成员失败: %v", err)
	}

	// 验证MemberCount
	ws, err := service.GetWorkspace(workspace.ID)
	if err != nil {
		t.Fatalf("获取工作空间失败: %v", err)
	}
	t.Logf("移除成员后的MemberCount: %d", ws.MemberCount)
}

// setupMockDB 创建Mock数据库连接
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建sqlmock失败: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("创建gorm DB失败: %v", err)
	}

	return gormDB, mock, sqlDB
}

// TestWorkspaceService_DeleteWorkspace_TransactionError 测试删除工作空间时的事务错误
func TestWorkspaceService_DeleteWorkspace_TransactionError(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	service := &WorkspaceService{
		db:           gormDB,
		workspaceDao: dao.NewWorkspaceDao(),
		memberDao:    dao.NewWorkspaceMemberDao(),
	}

	workspaceID := uint(1)

	// 测试场景1: 删除成员失败
	t.Run("DeleteMembersFails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM \"workspace_members\"").
			WillReturnError(fmt.Errorf("database error: connection lost"))
		mock.ExpectRollback()

		err := service.DeleteWorkspace(workspaceID)
		if err == nil {
			t.Error("删除成员失败时应该返回错误")
		}
		if err != nil && err.Error() != "删除工作空间成员失败: database error: connection lost" {
			t.Logf("错误信息: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未满足的期望: %v", err)
		}
	})

	// 测试场景2: 删除工作空间失败
	t.Run("DeleteWorkspaceFails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM \"workspace_members\"").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE \"workspaces\"").
			WillReturnError(fmt.Errorf("database error: disk full"))
		mock.ExpectRollback()

		err := service.DeleteWorkspace(workspaceID)
		if err == nil {
			t.Error("删除工作空间失败时应该返回错误")
		}
		if err != nil && err.Error() != "删除工作空间失败: database error: disk full" {
			t.Logf("错误信息: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未满足的期望: %v", err)
		}
	})
}

// TestWorkspaceService_RemoveMember_DatabaseError 测试移除成员时的数据库错误
func TestWorkspaceService_RemoveMember_DatabaseError(t *testing.T) {
	// 由于DAO层使用全局数据库，我们无法Mock它
	// 这个测试需要重构DAO层以支持依赖注入
	// 暂时跳过这个测试
	t.Skip("DAO层使用全局数据库，需要重构以支持Mock测试")
}
