package service

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
)

// TestWorkspaceIntegration_FullLifecycle 测试完整的工作空间生命周期
func TestWorkspaceIntegration_FullLifecycle(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 1. 创建所有者
	owner := &entity.User{
		Username:     "integration-owner-" + time.Now().Format("20060102150405"),
		Email:        "integration-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)
	t.Logf("✓ 步骤1: 创建所有者成功 (ID: %d)", owner.ID)

	// 2. 创建工作空间
	workspace := &entity.Workspace{
		OwnerID:     owner.ID,
		Name:        "Integration Test Workspace",
		Description: "This is a test workspace for integration testing",
		Type:        "team",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	t.Logf("✓ 步骤2: 创建工作空间成功 (ID: %d, UUID: %s)", workspace.ID, workspace.UUID)

	// 3. 添加多个成员
	members := make([]*entity.User, 3)
	for i := 0; i < 3; i++ {
		member := &entity.User{
			Username:     fmt.Sprintf("integration-member%d-%s", i, time.Now().Format("20060102150405")),
			Email:        fmt.Sprintf("integration-member%d-%s@example.com", i, time.Now().Format("20060102150405")),
			PasswordHash: "test-hash",
			Status:       "active",
		}
		if err := customerDao.Create(member); err != nil {
			t.Fatalf("创建成员%d失败: %v", i, err)
		}
		defer customerDao.Delete(member.ID)
		members[i] = member

		// 添加成员到工作空间
		role := "member"
		if i == 0 {
			role = "admin"
		}
		if err := service.AddMember(workspace.ID, member.ID, role); err != nil {
			t.Fatalf("添加成员%d失败: %v", i, err)
		}
		t.Logf("✓ 步骤3.%d: 添加成员成功 (ID: %d, Role: %s)", i+1, member.ID, role)
	}

	// 4. 查询成员列表
	memberList, err := service.ListMembers(workspace.ID)
	if err != nil {
		t.Fatalf("查询成员列表失败: %v", err)
	}
	if len(memberList) != 3 {
		t.Errorf("成员数量不正确，期望3个，实际%d个", len(memberList))
	}
	t.Logf("✓ 步骤4: 查询成员列表成功 (共%d个成员)", len(memberList))

	// 5. 更新工作空间信息
	workspace.Name = "Updated Integration Test Workspace"
	workspace.Description = "Updated description"
	if err := service.UpdateWorkspace(workspace); err != nil {
		t.Fatalf("更新工作空间失败: %v", err)
	}
	t.Logf("✓ 步骤5: 更新工作空间信息成功")

	// 验证更新
	updatedWs, err := service.GetWorkspace(workspace.ID)
	if err != nil {
		t.Fatalf("获取更新后的工作空间失败: %v", err)
	}
	if updatedWs.Name != "Updated Integration Test Workspace" {
		t.Errorf("工作空间名称未更新")
	}

	// 6. 移除一个成员
	if err := service.RemoveMember(workspace.ID, members[2].ID); err != nil {
		t.Fatalf("移除成员失败: %v", err)
	}
	t.Logf("✓ 步骤6: 移除成员成功")

	// 验证成员数量
	memberList, err = service.ListMembers(workspace.ID)
	if err != nil {
		t.Fatalf("查询成员列表失败: %v", err)
	}
	if len(memberList) != 2 {
		t.Errorf("移除成员后数量不正确，期望2个，实际%d个", len(memberList))
	}

	// 7. 删除工作空间
	if err := service.DeleteWorkspace(workspace.ID); err != nil {
		t.Fatalf("删除工作空间失败: %v", err)
	}
	t.Logf("✓ 步骤7: 删除工作空间成功")

	// 8. 验证所有数据已清理
	_, err = service.GetWorkspace(workspace.ID)
	if err == nil {
		t.Error("删除后仍能查询到工作空间")
	}
	t.Logf("✓ 步骤8: 验证数据清理成功")

	t.Log("✅ 完整生命周期测试通过")
}

// TestWorkspaceIntegration_MemberManagement 测试成员管理流程
func TestWorkspaceIntegration_MemberManagement(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 1. 创建工作空间和所有者
	owner := &entity.User{
		Username:     "member-mgmt-owner-" + time.Now().Format("20060102150405"),
		Email:        "member-mgmt-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Member Management Test Workspace",
		Type:    "team",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)
	t.Logf("✓ 步骤1: 创建工作空间成功")

	// 2. 添加成员A（member角色）
	memberA := &entity.User{
		Username:     "member-a-" + time.Now().Format("20060102150405"),
		Email:        "member-a-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(memberA); err != nil {
		t.Fatalf("创建成员A失败: %v", err)
	}
	defer customerDao.Delete(memberA.ID)

	if err := service.AddMember(workspace.ID, memberA.ID, "member"); err != nil {
		t.Fatalf("添加成员A失败: %v", err)
	}
	t.Logf("✓ 步骤2: 添加成员A（member角色）成功")

	// 3. 添加成员B（admin角色）
	memberB := &entity.User{
		Username:     "member-b-" + time.Now().Format("20060102150405"),
		Email:        "member-b-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(memberB); err != nil {
		t.Fatalf("创建成员B失败: %v", err)
	}
	defer customerDao.Delete(memberB.ID)

	if err := service.AddMember(workspace.ID, memberB.ID, "admin"); err != nil {
		t.Fatalf("添加成员B失败: %v", err)
	}
	t.Logf("✓ 步骤3: 添加成员B（admin角色）成功")

	// 4. 验证权限检查
	hasPermission, err := service.CheckPermission(workspace.ID, owner.ID)
	if err != nil || !hasPermission {
		t.Errorf("所有者权限检查失败")
	}
	t.Logf("✓ 步骤4: 所有者权限检查通过")

	hasPermission, err = service.CheckPermission(workspace.ID, memberA.ID)
	if err != nil || !hasPermission {
		t.Errorf("成员A权限检查失败")
	}
	t.Logf("✓ 步骤4: 成员A权限检查通过")

	hasPermission, err = service.CheckPermission(workspace.ID, memberB.ID)
	if err != nil || !hasPermission {
		t.Errorf("成员B权限检查失败")
	}
	t.Logf("✓ 步骤4: 成员B权限检查通过")

	// 5. 移除成员A
	if err := service.RemoveMember(workspace.ID, memberA.ID); err != nil {
		t.Fatalf("移除成员A失败: %v", err)
	}
	t.Logf("✓ 步骤5: 移除成员A成功")

	// 验证成员A已无权限
	hasPermission, err = service.CheckPermission(workspace.ID, memberA.ID)
	if err != nil || hasPermission {
		t.Errorf("移除后成员A仍有权限")
	}
	t.Logf("✓ 步骤5: 验证成员A已无权限")

	t.Log("✅ 成员管理流程测试通过")
}

// TestWorkspaceIntegration_PermissionCheck 测试权限检查流程
func TestWorkspaceIntegration_PermissionCheck(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 1. 创建工作空间和所有者
	owner := &entity.User{
		Username:     "perm-owner-" + time.Now().Format("20060102150405"),
		Email:        "perm-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	workspace := &entity.Workspace{
		OwnerID: owner.ID,
		Name:    "Permission Test Workspace",
		Type:    "team",
		Status:  "active",
	}
	if err := service.CreateWorkspace(workspace); err != nil {
		t.Fatalf("创建工作空间失败: %v", err)
	}
	defer service.DeleteWorkspace(workspace.ID)
	t.Logf("✓ 步骤1: 创建工作空间成功")

	// 2. 测试所有者权限
	hasPermission, err := service.CheckPermission(workspace.ID, owner.ID)
	if err != nil {
		t.Fatalf("检查所有者权限失败: %v", err)
	}
	if !hasPermission {
		t.Error("所有者应该有权限")
	}
	t.Logf("✓ 步骤2: 所有者权限检查通过")

	// 3. 添加成员并测试成员权限
	member := &entity.User{
		Username:     "perm-member-" + time.Now().Format("20060102150405"),
		Email:        "perm-member-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(member); err != nil {
		t.Fatalf("创建成员失败: %v", err)
	}
	defer customerDao.Delete(member.ID)

	if err := service.AddMember(workspace.ID, member.ID, "member"); err != nil {
		t.Fatalf("添加成员失败: %v", err)
	}

	hasPermission, err = service.CheckPermission(workspace.ID, member.ID)
	if err != nil {
		t.Fatalf("检查成员权限失败: %v", err)
	}
	if !hasPermission {
		t.Error("成员应该有权限")
	}
	t.Logf("✓ 步骤3: 成员权限检查通过")

	// 4. 测试非成员权限（应该拒绝）
	nonMember := &entity.User{
		Username:     "perm-nonmember-" + time.Now().Format("20060102150405"),
		Email:        "perm-nonmember-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(nonMember); err != nil {
		t.Fatalf("创建非成员失败: %v", err)
	}
	defer customerDao.Delete(nonMember.ID)

	hasPermission, err = service.CheckPermission(workspace.ID, nonMember.ID)
	if err != nil {
		t.Fatalf("检查非成员权限失败: %v", err)
	}
	if hasPermission {
		t.Error("非成员不应该有权限")
	}
	t.Logf("✓ 步骤4: 非成员权限检查通过（正确拒绝）")

	// 5. 测试archived工作空间权限
	workspace.Status = "archived"
	if err := service.UpdateWorkspace(workspace); err != nil {
		t.Fatalf("更新工作空间状态失败: %v", err)
	}

	hasPermission, err = service.CheckPermission(workspace.ID, owner.ID)
	if err == nil {
		t.Error("archived工作空间应该返回错误")
	}
	if hasPermission {
		t.Error("archived工作空间不应该有权限")
	}
	t.Logf("✓ 步骤5: archived工作空间权限检查通过（正确拒绝）")

	t.Log("✅ 权限检查流程测试通过")
}

// TestWorkspaceIntegration_ConcurrentCreation 测试并发创建工作空间
func TestWorkspaceIntegration_ConcurrentCreation(t *testing.T) {
	setupTestDB(t)

	db := database.GetDB()
	if err := db.AutoMigrate(&entity.Workspace{}, &entity.WorkspaceMember{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	service := NewWorkspaceService()
	customerDao := dao.NewUserDao()

	// 创建测试用户
	owner := &entity.User{
		Username:     "concurrent-owner-" + time.Now().Format("20060102150405"),
		Email:        "concurrent-owner-" + time.Now().Format("20060102150405") + "@example.com",
		PasswordHash: "test-hash",
		Status:       "active",
	}
	if err := customerDao.Create(owner); err != nil {
		t.Fatalf("创建所有者失败: %v", err)
	}
	defer customerDao.Delete(owner.ID)

	// 并发创建多个工作空间
	concurrency := 5
	var wg sync.WaitGroup
	errors := make(chan error, concurrency)
	workspaceIDs := make(chan uint, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			workspace := &entity.Workspace{
				OwnerID: owner.ID,
				Name:    fmt.Sprintf("Concurrent Workspace %d", index),
				Type:    "personal",
			}

			if err := service.CreateWorkspace(workspace); err != nil {
				errors <- err
				return
			}

			workspaceIDs <- workspace.ID
		}(i)
	}

	wg.Wait()
	close(errors)
	close(workspaceIDs)

	// 检查错误
	errorCount := 0
	for err := range errors {
		t.Errorf("并发创建失败: %v", err)
		errorCount++
	}

	// 收集创建的工作空间ID
	createdIDs := []uint{}
	for id := range workspaceIDs {
		createdIDs = append(createdIDs, id)
	}

	t.Logf("✓ 并发创建完成: 成功 %d 个, 失败 %d 个", len(createdIDs), errorCount)

	// 验证数据一致性
	for _, id := range createdIDs {
		ws, err := service.GetWorkspace(id)
		if err != nil {
			t.Errorf("获取工作空间 %d 失败: %v", id, err)
		} else {
			t.Logf("✓ 验证工作空间 %d: %s", id, ws.Name)
		}

		// 清理
		service.DeleteWorkspace(id)
	}

	t.Log("✅ 并发创建工作空间测试通过")
}
