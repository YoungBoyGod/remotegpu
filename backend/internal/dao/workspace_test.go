package dao

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupWorkspaceMockDB 创建Mock数据库连接
func setupWorkspaceMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
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

// TestWorkspaceDao_Create 测试创建工作空间
func TestWorkspaceDao_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspace := &entity.Workspace{
			UUID:        uuid.New(),
			OwnerID:     1,
			Name:        "Test Workspace",
			Description: "Test Description",
			Type:        "personal",
			MemberCount: 1,
			Status:      "active",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "workspaces"`).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "id"}).AddRow(workspace.UUID, 1))
		mock.ExpectCommit()

		err := dao.Create(workspace)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		workspace := &entity.Workspace{
			UUID:    uuid.New(),
			OwnerID: 1,
			Name:    "Test Workspace",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "workspaces"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Create(workspace)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_GetByID 测试根据ID获取工作空间
func TestWorkspaceDao_GetByID(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspaceID := uint(1)
		testUUID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			workspaceID, testUUID, 1, "Test Workspace", "Test Description", "personal",
			1, "active", now, now, nil,
		)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE id = .+`).
			WillReturnRows(rows)

		workspace, err := dao.GetByID(workspaceID)
		assert.NoError(t, err)
		assert.NotNil(t, workspace)
		assert.Equal(t, workspaceID, workspace.ID)
		assert.Equal(t, "Test Workspace", workspace.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		workspaceID := uint(999)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE id = .+`).
			WillReturnError(gorm.ErrRecordNotFound)

		workspace, err := dao.GetByID(workspaceID)
		assert.Error(t, err)
		assert.Nil(t, workspace)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		workspaceID := uint(1)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE id = .+`).
			WillReturnError(sql.ErrConnDone)

		workspace, err := dao.GetByID(workspaceID)
		assert.Error(t, err)
		assert.Nil(t, workspace)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_GetByUUID 测试根据UUID获取工作空间
func TestWorkspaceDao_GetByUUID(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		testUUID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, testUUID, 1, "Test Workspace", "Test Description", "personal",
			1, "active", now, now, nil,
		)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE uuid = .+`).
			WillReturnRows(rows)

		workspace, err := dao.GetByUUID(testUUID.String())
		assert.NoError(t, err)
		assert.NotNil(t, workspace)
		assert.Equal(t, testUUID, workspace.UUID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		testUUID := uuid.New()

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE uuid = .+`).
			WillReturnError(gorm.ErrRecordNotFound)

		workspace, err := dao.GetByUUID(testUUID.String())
		assert.Error(t, err)
		assert.Nil(t, workspace)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		testUUID := uuid.New()

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE uuid = .+`).
			WillReturnError(sql.ErrConnDone)

		workspace, err := dao.GetByUUID(testUUID.String())
		assert.Error(t, err)
		assert.Nil(t, workspace)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_Update 测试更新工作空间
func TestWorkspaceDao_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspace := &entity.Workspace{
			ID:          1,
			UUID:        uuid.New(),
			OwnerID:     1,
			Name:        "Updated Workspace",
			Description: "Updated Description",
			Type:        "team",
			MemberCount: 5,
			Status:      "active",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "workspaces"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Update(workspace)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		workspace := &entity.Workspace{
			ID:   1,
			Name: "Updated Workspace",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "workspaces"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Update(workspace)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_Delete 测试删除工作空间（软删除）
func TestWorkspaceDao_Delete(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspaceID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "workspaces" SET "deleted_at"=.+ WHERE`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Delete(workspaceID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		workspaceID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "workspaces" SET "deleted_at"=.+ WHERE`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Delete(workspaceID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_GetByOwnerID 测试根据所有者ID获取工作空间列表
func TestWorkspaceDao_GetByOwnerID(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		ownerID := uint(1)
		now := time.Now()
		testUUID1 := uuid.New()
		testUUID2 := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(1, testUUID1, ownerID, "Workspace 1", "Desc 1", "personal",
				1, "active", now, now, nil).
			AddRow(2, testUUID2, ownerID, "Workspace 2", "Desc 2", "team",
				5, "active", now, now, nil)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE owner_id = .+`).
			WillReturnRows(rows)

		workspaces, err := dao.GetByOwnerID(ownerID)
		assert.NoError(t, err)
		assert.Len(t, workspaces, 2)
		assert.Equal(t, ownerID, workspaces[0].OwnerID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		ownerID := uint(999)

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		})

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE owner_id = .+`).
			WillReturnRows(rows)

		workspaces, err := dao.GetByOwnerID(ownerID)
		assert.NoError(t, err)
		assert.Len(t, workspaces, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		ownerID := uint(1)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE owner_id = .+`).
			WillReturnError(sql.ErrConnDone)

		workspaces, err := dao.GetByOwnerID(ownerID)
		assert.Error(t, err)
		assert.Nil(t, workspaces)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_List 测试分页获取工作空间列表
func TestWorkspaceDao_List(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		testUUID1 := uuid.New()
		testUUID2 := uuid.New()

		mock.ExpectQuery(`SELECT count\(\*\) FROM "workspaces"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(1, testUUID1, 1, "Workspace 1", "Desc 1", "personal",
				1, "active", now, now, nil).
			AddRow(2, testUUID2, 2, "Workspace 2", "Desc 2", "team",
				5, "active", now, now, nil)

		mock.ExpectQuery(`SELECT \* FROM "workspaces"`).
			WillReturnRows(rows)

		workspaces, total, err := dao.List(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, workspaces, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "workspaces"`).
			WillReturnError(sql.ErrConnDone)

		workspaces, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, workspaces)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "workspaces"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`SELECT \* FROM "workspaces"`).
			WillReturnError(sql.ErrConnDone)

		workspaces, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, workspaces)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceDao_GetByStatus 测试根据状态获取工作空间列表
func TestWorkspaceDao_GetByStatus(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		status := "active"
		now := time.Now()
		testUUID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, testUUID, 1, "Active Workspace", "Desc", "personal",
			1, status, now, now, nil,
		)

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE status = .+`).
			WillReturnRows(rows)

		workspaces, err := dao.GetByStatus(status)
		assert.NoError(t, err)
		assert.Len(t, workspaces, 1)
		assert.Equal(t, status, workspaces[0].Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		status := "archived"

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "owner_id", "name", "description", "type",
			"member_count", "status", "created_at", "updated_at", "deleted_at",
		})

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE status = .+`).
			WillReturnRows(rows)

		workspaces, err := dao.GetByStatus(status)
		assert.NoError(t, err)
		assert.Len(t, workspaces, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		status := "active"

		mock.ExpectQuery(`SELECT .* FROM "workspaces" WHERE status = .+`).
			WillReturnError(sql.ErrConnDone)

		workspaces, err := dao.GetByStatus(status)
		assert.Error(t, err)
		assert.Nil(t, workspaces)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_Create 测试创建工作空间成员
func TestWorkspaceMemberDao_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		member := &entity.WorkspaceMember{
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
			Status:      "active",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "workspace_members"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := dao.Create(member)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		member := &entity.WorkspaceMember{
			WorkspaceID: 1,
			UserID:      2,
			Role:        "member",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "workspace_members"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Create(member)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_GetByID 测试根据ID获取工作空间成员
func TestWorkspaceMemberDao_GetByID(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		memberID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "workspace_id", "user_id", "role", "status", "joined_at", "created_at",
		}).AddRow(
			memberID, 1, 2, "member", "active", now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE id = .+`).
			WillReturnRows(rows)

		member, err := dao.GetByID(memberID)
		assert.NoError(t, err)
		assert.NotNil(t, member)
		assert.Equal(t, memberID, member.ID)
		assert.Equal(t, "member", member.Role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		memberID := uint(999)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE id = .+`).
			WillReturnError(gorm.ErrRecordNotFound)

		member, err := dao.GetByID(memberID)
		assert.Error(t, err)
		assert.Nil(t, member)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		memberID := uint(1)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE id = .+`).
			WillReturnError(sql.ErrConnDone)

		member, err := dao.GetByID(memberID)
		assert.Error(t, err)
		assert.Nil(t, member)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_Update 测试更新工作空间成员
func TestWorkspaceMemberDao_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		member := &entity.WorkspaceMember{
			ID:          1,
			WorkspaceID: 1,
			UserID:      2,
			Role:        "admin",
			Status:      "active",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "workspace_members"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Update(member)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		member := &entity.WorkspaceMember{
			ID:   1,
			Role: "admin",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "workspace_members"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Update(member)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_Delete 测试删除工作空间成员
func TestWorkspaceMemberDao_Delete(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		memberID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "workspace_members" WHERE`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Delete(memberID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		memberID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "workspace_members" WHERE`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Delete(memberID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_GetByWorkspaceID 测试根据工作空间ID获取成员列表
func TestWorkspaceMemberDao_GetByWorkspaceID(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspaceID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "workspace_id", "user_id", "role", "status", "joined_at", "created_at",
		}).
			AddRow(1, workspaceID, 2, "admin", "active", now, now).
			AddRow(2, workspaceID, 3, "member", "active", now, now)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE workspace_id = .+`).
			WillReturnRows(rows)

		members, err := dao.GetByWorkspaceID(workspaceID)
		assert.NoError(t, err)
		assert.Len(t, members, 2)
		assert.Equal(t, workspaceID, members[0].WorkspaceID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		workspaceID := uint(999)

		rows := sqlmock.NewRows([]string{
			"id", "workspace_id", "user_id", "role", "status", "joined_at", "created_at",
		})

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE workspace_id = .+`).
			WillReturnRows(rows)

		members, err := dao.GetByWorkspaceID(workspaceID)
		assert.NoError(t, err)
		assert.Len(t, members, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		workspaceID := uint(1)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE workspace_id = .+`).
			WillReturnError(sql.ErrConnDone)

		members, err := dao.GetByWorkspaceID(workspaceID)
		assert.Error(t, err)
		assert.Nil(t, members)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_GetByUserID 测试根据用户ID获取成员列表
func TestWorkspaceMemberDao_GetByUserID(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		userID := uint(2)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "workspace_id", "user_id", "role", "status", "joined_at", "created_at",
		}).
			AddRow(1, 1, userID, "admin", "active", now, now).
			AddRow(2, 2, userID, "member", "active", now, now)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE user_id = .+`).
			WillReturnRows(rows)

		members, err := dao.GetByUserID(userID)
		assert.NoError(t, err)
		assert.Len(t, members, 2)
		assert.Equal(t, userID, members[0].UserID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		userID := uint(999)

		rows := sqlmock.NewRows([]string{
			"id", "workspace_id", "user_id", "role", "status", "joined_at", "created_at",
		})

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE user_id = .+`).
			WillReturnRows(rows)

		members, err := dao.GetByUserID(userID)
		assert.NoError(t, err)
		assert.Len(t, members, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		userID := uint(2)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE user_id = .+`).
			WillReturnError(sql.ErrConnDone)

		members, err := dao.GetByUserID(userID)
		assert.Error(t, err)
		assert.Nil(t, members)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestWorkspaceMemberDao_GetByWorkspaceAndUser 测试根据工作空间ID和用户ID获取成员
func TestWorkspaceMemberDao_GetByWorkspaceAndUser(t *testing.T) {
	gormDB, mock, sqlDB := setupWorkspaceMockDB(t)
	defer sqlDB.Close()

	dao := &WorkspaceMemberDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspaceID := uint(1)
		userID := uint(2)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "workspace_id", "user_id", "role", "status", "joined_at", "created_at",
		}).AddRow(
			1, workspaceID, userID, "admin", "active", now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE workspace_id = .+ AND user_id = .+`).
			WillReturnRows(rows)

		member, err := dao.GetByWorkspaceAndUser(workspaceID, userID)
		assert.NoError(t, err)
		assert.NotNil(t, member)
		assert.Equal(t, workspaceID, member.WorkspaceID)
		assert.Equal(t, userID, member.UserID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		workspaceID := uint(1)
		userID := uint(999)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE workspace_id = .+ AND user_id = .+`).
			WillReturnError(gorm.ErrRecordNotFound)

		member, err := dao.GetByWorkspaceAndUser(workspaceID, userID)
		assert.Error(t, err)
		assert.Nil(t, member)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		workspaceID := uint(1)
		userID := uint(2)

		mock.ExpectQuery(`SELECT .* FROM "workspace_members" WHERE workspace_id = .+ AND user_id = .+`).
			WillReturnError(sql.ErrConnDone)

		member, err := dao.GetByWorkspaceAndUser(workspaceID, userID)
		assert.Error(t, err)
		assert.Nil(t, member)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
