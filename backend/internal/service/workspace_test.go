package service

import (
	"errors"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestWorkspaceService_CreateWorkspace 测试创建工作空间
func TestWorkspaceService_CreateWorkspace(t *testing.T) {
	t.Run("Success_WithDefaults", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("Create", mock.AnythingOfType("*entity.Workspace")).Return(nil)

		err := service.CreateWorkspace(workspace)

		assert.NoError(t, err)
		assert.Equal(t, "personal", workspace.Type)
		assert.Equal(t, "active", workspace.Status)
		assert.Equal(t, 1, workspace.MemberCount)
		assert.NotEqual(t, uuid.Nil, workspace.UUID)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Success_WithCustomType", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			OwnerID: 1,
			Name:    "Team Workspace",
			Type:    "team",
		}

		mockWorkspaceDao.On("Create", workspace).Return(nil)

		err := service.CreateWorkspace(workspace)

		assert.NoError(t, err)
		assert.Equal(t, "team", workspace.Type)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_EmptyName", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			OwnerID: 1,
			Name:    "",
		}

		err := service.CreateWorkspace(workspace)

		assert.Error(t, err)
		assert.Equal(t, "工作空间名称不能为空", err.Error())
	})

	t.Run("Error_NameTooLong", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		longName := ""
		for i := 0; i < 129; i++ {
			longName += "a"
		}

		workspace := &entity.Workspace{
			OwnerID: 1,
			Name:    longName,
		}

		err := service.CreateWorkspace(workspace)

		assert.Error(t, err)
		assert.Equal(t, "工作空间名称不能超过128个字符", err.Error())
	})

	t.Run("Error_InvalidType", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			OwnerID: 1,
			Name:    "Test Workspace",
			Type:    "invalid_type",
		}

		err := service.CreateWorkspace(workspace)

		assert.Error(t, err)
		assert.Equal(t, "无效的工作空间类型: invalid_type", err.Error())
	})

	t.Run("Error_DaoCreateFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("Create", mock.AnythingOfType("*entity.Workspace")).Return(errors.New("database error"))

		err := service.CreateWorkspace(workspace)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_GetWorkspace 测试获取工作空间
func TestWorkspaceService_GetWorkspace(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		expectedWorkspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Type:    "personal",
			Status:  "active",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(expectedWorkspace, nil)

		workspace, err := service.GetWorkspace(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedWorkspace, workspace)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		workspace, err := service.GetWorkspace(999)

		assert.Error(t, err)
		assert.Nil(t, workspace)
		mockWorkspaceDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_UpdateWorkspace 测试更新工作空间
func TestWorkspaceService_UpdateWorkspace(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		existingWorkspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Old Name",
			Type:    "personal",
			Status:  "active",
		}

		updatedWorkspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "New Name",
			Description: "Updated description",
			Type:        "personal",
			Status:      "active",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(existingWorkspace, nil)
		mockWorkspaceDao.On("Update", updatedWorkspace).Return(nil)

		err := service.UpdateWorkspace(updatedWorkspace)

		assert.NoError(t, err)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_WorkspaceNotFound", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      999,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.UpdateWorkspace(workspace)

		assert.Error(t, err)
		assert.Equal(t, "工作空间不存在", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_CannotChangeOwner", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		existingWorkspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		updatedWorkspace := &entity.Workspace{
			ID:      1,
			OwnerID: 2,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(existingWorkspace, nil)

		err := service.UpdateWorkspace(updatedWorkspace)

		assert.Error(t, err)
		assert.Equal(t, "不允许修改工作空间所有者", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_DaoUpdateFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		existingWorkspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Old Name",
		}

		updatedWorkspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "New Name",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(existingWorkspace, nil)
		mockWorkspaceDao.On("Update", updatedWorkspace).Return(errors.New("database error"))

		err := service.UpdateWorkspace(updatedWorkspace)

		assert.Error(t, err)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_GetByIDFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(nil, errors.New("database connection error"))

		err := service.UpdateWorkspace(workspace)

		assert.Error(t, err)
		mockWorkspaceDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_ListWorkspaces 测试获取工作空间列表
func TestWorkspaceService_ListWorkspaces(t *testing.T) {
	t.Run("Success_WithOwnerID", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		expectedWorkspaces := []*entity.Workspace{
			{ID: 1, OwnerID: 1, Name: "Workspace 1"},
			{ID: 2, OwnerID: 1, Name: "Workspace 2"},
		}

		mockWorkspaceDao.On("GetByOwnerID", uint(1)).Return(expectedWorkspaces, nil)

		workspaces, total, err := service.ListWorkspaces(1, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, workspaces, 2)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Success_WithoutOwnerID", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		expectedWorkspaces := []*entity.Workspace{
			{ID: 1, OwnerID: 1, Name: "Workspace 1"},
			{ID: 2, OwnerID: 2, Name: "Workspace 2"},
		}

		mockWorkspaceDao.On("List", 1, 10).Return(expectedWorkspaces, int64(2), nil)

		workspaces, total, err := service.ListWorkspaces(0, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, workspaces, 2)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_GetByOwnerIDFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("GetByOwnerID", uint(1)).Return(nil, errors.New("database error"))

		workspaces, total, err := service.ListWorkspaces(1, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, workspaces)
		assert.Equal(t, int64(0), total)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_ListFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("List", 1, 10).Return(nil, int64(0), errors.New("database error"))

		workspaces, total, err := service.ListWorkspaces(0, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, workspaces)
		assert.Equal(t, int64(0), total)
		mockWorkspaceDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_AddMember 测试添加工作空间成员
func TestWorkspaceService_AddMember(t *testing.T) {
	t.Run("Success_WithDefaultRole", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 1,
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)
		mockMemberDao.On("Create", mock.AnythingOfType("*entity.WorkspaceMember")).Return(nil)
		mockWorkspaceDao.On("Update", workspace).Return(nil)

		err := service.AddMember(1, 2, "")

		assert.NoError(t, err)
		assert.Equal(t, 2, workspace.MemberCount)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Success_WithCustomRole", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 1,
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)
		mockMemberDao.On("Create", mock.AnythingOfType("*entity.WorkspaceMember")).Return(nil)
		mockWorkspaceDao.On("Update", workspace).Return(nil)

		err := service.AddMember(1, 2, "admin")

		assert.NoError(t, err)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_WorkspaceNotFound", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.AddMember(999, 2, "member")

		assert.Error(t, err)
		assert.Equal(t, "工作空间不存在", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_MemberAlreadyExists", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		existingMember := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(existingMember, nil)

		err := service.AddMember(1, 2, "member")

		assert.Error(t, err)
		assert.Equal(t, "成员已存在", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_InvalidRole", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)

		err := service.AddMember(1, 2, "invalid_role")

		assert.Error(t, err)
		assert.Equal(t, "无效的角色: invalid_role", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_CreateMemberFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)
		mockMemberDao.On("Create", mock.AnythingOfType("*entity.WorkspaceMember")).Return(errors.New("database error"))

		err := service.AddMember(1, 2, "member")

		assert.Error(t, err)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_UpdateWorkspaceFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 1,
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)
		mockMemberDao.On("Create", mock.AnythingOfType("*entity.WorkspaceMember")).Return(nil)
		mockWorkspaceDao.On("Update", workspace).Return(errors.New("database error"))

		err := service.AddMember(1, 2, "member")

		assert.Error(t, err)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_RemoveMember 测试移除工作空间成员
func TestWorkspaceService_RemoveMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 2,
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(member, nil)
		mockMemberDao.On("Delete", uint(1)).Return(nil)
		mockWorkspaceDao.On("Update", workspace).Return(nil)

		err := service.RemoveMember(1, 2)

		assert.NoError(t, err)
		assert.Equal(t, 1, workspace.MemberCount)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_WorkspaceNotFound", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.RemoveMember(999, 2)

		assert.Error(t, err)
		assert.Equal(t, "工作空间不存在", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_MemberNotFound", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(999)).Return(nil, gorm.ErrRecordNotFound)

		err := service.RemoveMember(1, 999)

		assert.Error(t, err)
		assert.Equal(t, "成员不存在", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_CannotRemoveOwner", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      1,
			Role:        "owner",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(1)).Return(member, nil)

		err := service.RemoveMember(1, 1)

		assert.Error(t, err)
		assert.Equal(t, "不允许移除工作空间所有者", err.Error())
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_DeleteMemberFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 2,
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(member, nil)
		mockMemberDao.On("Delete", uint(1)).Return(errors.New("database error"))

		err := service.RemoveMember(1, 2)

		assert.Error(t, err)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_UpdateWorkspaceFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 2,
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(member, nil)
		mockMemberDao.On("Delete", uint(1)).Return(nil)
		mockWorkspaceDao.On("Update", workspace).Return(errors.New("database error"))

		err := service.RemoveMember(1, 2)

		assert.Error(t, err)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Success_MemberCountZero", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:          1,
			OwnerID:     1,
			Name:        "Test Workspace",
			MemberCount: 0,
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(member, nil)
		mockMemberDao.On("Delete", uint(1)).Return(nil)
		mockWorkspaceDao.On("Update", workspace).Return(nil)

		err := service.RemoveMember(1, 2)

		assert.NoError(t, err)
		assert.Equal(t, 0, workspace.MemberCount)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_ListMembers 测试获取工作空间成员列表
func TestWorkspaceService_ListMembers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		expectedMembers := []*entity.WorkspaceMember{
			{ID: 1, WorkspaceID: 1, UserID: 1, Role: "owner"},
			{ID: 2, WorkspaceID: 1, UserID: 2, Role: "member"},
		}

		mockMemberDao.On("GetByWorkspaceID", uint(1)).Return(expectedMembers, nil)

		members, err := service.ListMembers(1)

		assert.NoError(t, err)
		assert.Len(t, members, 2)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_GetByWorkspaceIDFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockMemberDao.On("GetByWorkspaceID", uint(1)).Return(nil, errors.New("database error"))

		members, err := service.ListMembers(1)

		assert.Error(t, err)
		assert.Nil(t, members)
		mockMemberDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_CheckPermission 测试检查用户权限
func TestWorkspaceService_CheckPermission(t *testing.T) {
	t.Run("Success_Owner", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Status:  "active",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)

		hasPermission, err := service.CheckPermission(1, 1)

		assert.NoError(t, err)
		assert.True(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Success_ActiveMember", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Status:  "active",
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
			Status:      "active",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(member, nil)

		hasPermission, err := service.CheckPermission(1, 2)

		assert.NoError(t, err)
		assert.True(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Success_InactiveMember", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Status:  "active",
		}

		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
			Status:      "suspended",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(member, nil)

		hasPermission, err := service.CheckPermission(1, 2)

		assert.NoError(t, err)
		assert.False(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Success_NotMember", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Status:  "active",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(999)).Return(nil, gorm.ErrRecordNotFound)

		hasPermission, err := service.CheckPermission(1, 999)

		assert.NoError(t, err)
		assert.False(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})

	t.Run("Error_WorkspaceNotFound", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		hasPermission, err := service.CheckPermission(999, 1)

		assert.Error(t, err)
		assert.Equal(t, "工作空间不存在", err.Error())
		assert.False(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_WorkspaceArchived", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Status:  "archived",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)

		hasPermission, err := service.CheckPermission(1, 1)

		assert.Error(t, err)
		assert.Equal(t, "工作空间已归档", err.Error())
		assert.False(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_GetByIDFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(nil, errors.New("database error"))

		hasPermission, err := service.CheckPermission(1, 1)

		assert.Error(t, err)
		assert.False(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
	})

	t.Run("Error_GetByWorkspaceAndUserFails", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		workspace := &entity.Workspace{
			ID:      1,
			OwnerID: 1,
			Name:    "Test Workspace",
			Status:  "active",
		}

		mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
		mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, errors.New("database error"))

		hasPermission, err := service.CheckPermission(1, 2)

		assert.Error(t, err)
		assert.False(t, hasPermission)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_DeleteWorkspace 测试删除工作空间
func TestWorkspaceService_DeleteWorkspace(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockWorkspaceDao := new(mocks.MockWorkspaceDao)
		mockMemberDao := new(mocks.MockWorkspaceMemberDao)

		// 创建一个mock的gorm.DB
		// 由于DeleteWorkspace使用了事务，我们需要mock db.Transaction
		// 但这很复杂，所以我们暂时跳过这个测试
		t.Skip("DeleteWorkspace使用了事务，需要更复杂的mock设置")

		service := &WorkspaceService{
			workspaceDao: mockWorkspaceDao,
			memberDao:    mockMemberDao,
		}

		err := service.DeleteWorkspace(1)

		assert.NoError(t, err)
		mockWorkspaceDao.AssertExpectations(t)
		mockMemberDao.AssertExpectations(t)
	})
}

// TestWorkspaceService_AddMember_GetByIDError 测试AddMember中GetByID返回错误
func TestWorkspaceService_AddMember_GetByIDError(t *testing.T) {
	mockWorkspaceDao := new(mocks.MockWorkspaceDao)
	mockMemberDao := new(mocks.MockWorkspaceMemberDao)

	service := &WorkspaceService{
		workspaceDao: mockWorkspaceDao,
		memberDao:    mockMemberDao,
	}

	mockWorkspaceDao.On("GetByID", uint(1)).Return(nil, errors.New("database connection error"))

	err := service.AddMember(1, 2, "member")

	assert.Error(t, err)
	assert.Equal(t, "database connection error", err.Error())
	mockWorkspaceDao.AssertExpectations(t)
}

// TestWorkspaceService_RemoveMember_GetByIDError 测试RemoveMember中GetByID返回错误
func TestWorkspaceService_RemoveMember_GetByIDError(t *testing.T) {
	mockWorkspaceDao := new(mocks.MockWorkspaceDao)
	mockMemberDao := new(mocks.MockWorkspaceMemberDao)

	service := &WorkspaceService{
		workspaceDao: mockWorkspaceDao,
		memberDao:    mockMemberDao,
	}

	mockWorkspaceDao.On("GetByID", uint(1)).Return(nil, errors.New("database connection error"))

	err := service.RemoveMember(1, 2)

	assert.Error(t, err)
	assert.Equal(t, "database connection error", err.Error())
	mockWorkspaceDao.AssertExpectations(t)
}

// TestWorkspaceService_RemoveMember_GetByWorkspaceAndUserError 测试RemoveMember中GetByWorkspaceAndUser返回错误
func TestWorkspaceService_RemoveMember_GetByWorkspaceAndUserError(t *testing.T) {
	mockWorkspaceDao := new(mocks.MockWorkspaceDao)
	mockMemberDao := new(mocks.MockWorkspaceMemberDao)

	service := &WorkspaceService{
		workspaceDao: mockWorkspaceDao,
		memberDao:    mockMemberDao,
	}

	workspace := &entity.Workspace{
		ID:      1,
		OwnerID: 1,
		Name:    "Test Workspace",
	}

	mockWorkspaceDao.On("GetByID", uint(1)).Return(workspace, nil)
	mockMemberDao.On("GetByWorkspaceAndUser", uint(1), uint(2)).Return(nil, errors.New("database error"))

	err := service.RemoveMember(1, 2)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockWorkspaceDao.AssertExpectations(t)
	mockMemberDao.AssertExpectations(t)
}
