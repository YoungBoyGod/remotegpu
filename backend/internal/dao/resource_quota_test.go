package dao

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupResourceQuotaMockDB 创建Mock数据库连接
func setupResourceQuotaMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
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

// TestResourceQuotaDao_Create 测试创建资源配额
func TestResourceQuotaDao_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success_UserLevel", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			UserID:           1,
			WorkspaceID:      nil,
			QuotaLevel:       "pro",
			CPU:              16,
			Memory:           32768,
			GPU:              4,
			Storage:          1000,
			EnvironmentQuota: 10,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "resource_quotas"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := dao.Create(quota)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success_WorkspaceLevel", func(t *testing.T) {
		workspaceID := uint(1)
		quota := &entity.ResourceQuota{
			UserID:           1,
			WorkspaceID:      &workspaceID,
			QuotaLevel:       "basic",
			CPU:              8,
			Memory:           16384,
			GPU:              2,
			Storage:          500,
			EnvironmentQuota: 5,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "resource_quotas"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectCommit()

		err := dao.Create(quota)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			UserID:     1,
			QuotaLevel: "pro",
			CPU:        16,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "resource_quotas"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Create(quota)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetByID 测试根据ID获取资源配额
func TestResourceQuotaDao_GetByID(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		quotaID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			quotaID, 1, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(quotaID, 1).
			WillReturnRows(rows)

		quota, err := dao.GetByID(quotaID)
		assert.NoError(t, err)
		assert.NotNil(t, quota)
		assert.Equal(t, quotaID, quota.ID)
		assert.Equal(t, 16, quota.CPU)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		quotaID := uint(999)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(quotaID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		quota, err := dao.GetByID(quotaID)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		quotaID := uint(1)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(quotaID, 1).
			WillReturnError(sql.ErrConnDone)

		quota, err := dao.GetByID(quotaID)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetByUserID 测试根据用户ID获取资源配额
func TestResourceQuotaDao_GetByUserID(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		quota, err := dao.GetByUserID(userID)
		assert.NoError(t, err)
		assert.NotNil(t, quota)
		assert.Equal(t, userID, quota.UserID)
		assert.Nil(t, quota.WorkspaceID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		userID := uint(999)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		quota, err := dao.GetByUserID(userID)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		userID := uint(1)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnError(sql.ErrConnDone)

		quota, err := dao.GetByUserID(userID)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetByWorkspaceID 测试根据工作空间ID获取资源配额
func TestResourceQuotaDao_GetByWorkspaceID(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		workspaceID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, 1, workspaceID, "basic", 8, 16384,
			2, 500, 5, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(workspaceID, 1).
			WillReturnRows(rows)

		quota, err := dao.GetByWorkspaceID(workspaceID)
		assert.NoError(t, err)
		assert.NotNil(t, quota)
		assert.Equal(t, workspaceID, *quota.WorkspaceID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		workspaceID := uint(999)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(workspaceID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		quota, err := dao.GetByWorkspaceID(workspaceID)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetByUserAndWorkspace 测试根据用户和工作空间ID获取资源配额
func TestResourceQuotaDao_GetByUserAndWorkspace(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		workspaceID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, workspaceID, "basic", 8, 16384,
			2, 500, 5, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, workspaceID, 1).
			WillReturnRows(rows)

		quota, err := dao.GetByUserAndWorkspace(userID, workspaceID)
		assert.NoError(t, err)
		assert.NotNil(t, quota)
		assert.Equal(t, userID, quota.UserID)
		assert.Equal(t, workspaceID, *quota.WorkspaceID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		userID := uint(999)
		workspaceID := uint(999)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, workspaceID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		quota, err := dao.GetByUserAndWorkspace(userID, workspaceID)
		assert.Error(t, err)
		assert.Nil(t, quota)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_Update 测试更新资源配额
func TestResourceQuotaDao_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			ID:               1,
			UserID:           1,
			WorkspaceID:      nil,
			QuotaLevel:       "pro",
			CPU:              32,
			Memory:           65536,
			GPU:              8,
			Storage:          2000,
			EnvironmentQuota: 20,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "resource_quotas"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Update(quota)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			ID:     1,
			UserID: 1,
			CPU:    32,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "resource_quotas"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Update(quota)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_Delete 测试删除资源配额
func TestResourceQuotaDao_Delete(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		quotaID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "resource_quotas"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Delete(quotaID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		quotaID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "resource_quotas"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Delete(quotaID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_List 测试获取资源配额列表
func TestResourceQuotaDao_List(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()

		// Mock Count查询
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "resource_quotas"`).
			WillReturnRows(countRows)

		// Mock List查询
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(1, 1, nil, "pro", 16, 32768, 4, 1000, 10, now, now, nil).
			AddRow(2, 2, nil, "basic", 8, 16384, 2, 500, 5, now, now, nil)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnRows(rows)

		quotas, total, err := dao.List(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, quotas, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PageLessThanOne", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "resource_quotas"`).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(1, 1, nil, "pro", 16, 32768, 4, 1000, 10, time.Now(), time.Now(), nil)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnRows(rows)

		quotas, total, err := dao.List(0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, quotas, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PageSizeLessThanOne", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "resource_quotas"`).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(1, 1, nil, "pro", 16, 32768, 4, 1000, 10, time.Now(), time.Now(), nil)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnRows(rows)

		quotas, total, err := dao.List(1, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, quotas, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PageSizeGreaterThan100", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "resource_quotas"`).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(1, 1, nil, "pro", 16, 32768, 4, 1000, 10, time.Now(), time.Now(), nil)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnRows(rows)

		quotas, total, err := dao.List(1, 200)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, quotas, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "resource_quotas"`).
			WillReturnError(sql.ErrConnDone)

		quotas, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, quotas)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "resource_quotas"`).
			WillReturnRows(countRows)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnError(sql.ErrConnDone)

		quotas, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, quotas)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetByQuotaLevel 测试根据配额级别获取资源配额列表
func TestResourceQuotaDao_GetByQuotaLevel(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		level := "pro"

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(1, 1, nil, level, 16, 32768, 4, 1000, 10, now, now, nil).
			AddRow(2, 2, nil, level, 16, 32768, 4, 1000, 10, now, now, nil)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnRows(rows)

		quotas, err := dao.GetByQuotaLevel(level)
		assert.NoError(t, err)
		assert.Len(t, quotas, 2)
		assert.Equal(t, level, quotas[0].QuotaLevel)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		level := "enterprise"

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		})

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnRows(rows)

		quotas, err := dao.GetByQuotaLevel(level)
		assert.NoError(t, err)
		assert.Len(t, quotas, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		level := "pro"

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WillReturnError(sql.ErrConnDone)

		quotas, err := dao.GetByQuotaLevel(level)
		assert.Error(t, err)
		assert.Nil(t, quotas)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetUsedResources 测试获取已使用的资源
func TestResourceQuotaDao_GetUsedResources(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success_WithResources", func(t *testing.T) {
		userID := uint(1)

		rows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(rows)

		used, err := dao.GetUsedResources(userID)
		assert.NoError(t, err)
		assert.NotNil(t, used)
		assert.Equal(t, 8, used.CPU)
		assert.Equal(t, int64(16384), used.Memory)
		assert.Equal(t, 2, used.GPU)
		assert.Equal(t, int64(500), used.Storage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success_NoResources", func(t *testing.T) {
		userID := uint(2)

		rows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(0, 0, 0, 0)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(rows)

		used, err := dao.GetUsedResources(userID)
		assert.NoError(t, err)
		assert.NotNil(t, used)
		assert.Equal(t, 0, used.CPU)
		assert.Equal(t, int64(0), used.Memory)
		assert.Equal(t, 0, used.GPU)
		assert.Equal(t, int64(0), used.Storage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		userID := uint(1)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnError(sql.ErrConnDone)

		used, err := dao.GetUsedResources(userID)
		assert.Error(t, err)
		assert.Nil(t, used)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_GetAvailableQuota 测试获取可用配额
func TestResourceQuotaDao_GetAvailableQuota(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		available, err := dao.GetAvailableQuota(userID)
		assert.NoError(t, err)
		assert.NotNil(t, available)
		assert.Equal(t, 8, available.CPU)
		assert.Equal(t, int64(16384), available.Memory)
		assert.Equal(t, 2, available.GPU)
		assert.Equal(t, int64(500), available.Storage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QuotaNotFound", func(t *testing.T) {
		userID := uint(999)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		available, err := dao.GetAvailableQuota(userID)
		assert.Error(t, err)
		assert.Nil(t, available)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUsedResourcesError", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources error
		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnError(sql.ErrConnDone)

		available, err := dao.GetAvailableQuota(userID)
		assert.Error(t, err)
		assert.Nil(t, available)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestResourceQuotaDao_CheckQuota 测试检查配额
func TestResourceQuotaDao_CheckQuota(t *testing.T) {
	gormDB, mock, sqlDB := setupResourceQuotaMockDB(t)
	defer sqlDB.Close()

	dao := &ResourceQuotaDao{db: gormDB}

	t.Run("Success_QuotaSufficient", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		// 请求: CPU=4, Memory=8192, GPU=1, Storage=200
		// 已用: CPU=8, Memory=16384, GPU=2, Storage=500
		// 总配额: CPU=16, Memory=32768, GPU=4, Storage=1000
		// 已用+请求: CPU=12, Memory=24576, GPU=3, Storage=700 < 总配额
		ok, err := dao.CheckQuota(userID, 4, 8192, 1, 200)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CPUQuotaExceeded", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		// 请求: CPU=10, Memory=8192, GPU=1, Storage=200
		// 已用+请求: CPU=18 > 16 (总配额)
		ok, err := dao.CheckQuota(userID, 10, 8192, 1, 200)
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MemoryQuotaExceeded", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		// 请求: CPU=4, Memory=20000, GPU=1, Storage=200
		// 已用+请求: Memory=36384 > 32768 (总配额)
		ok, err := dao.CheckQuota(userID, 4, 20000, 1, 200)
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GPUQuotaExceeded", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		// 请求: CPU=4, Memory=8192, GPU=3, Storage=200
		// 已用+请求: GPU=5 > 4 (总配额)
		ok, err := dao.CheckQuota(userID, 4, 8192, 3, 200)
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("StorageQuotaExceeded", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		// 请求: CPU=4, Memory=8192, GPU=1, Storage=600
		// 已用+请求: Storage=1100 > 1000 (总配额)
		ok, err := dao.CheckQuota(userID, 4, 8192, 1, 600)
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ExactQuotaMatch", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources
		usedRows := sqlmock.NewRows([]string{"cpu", "memory", "gpu", "storage"}).
			AddRow(8, 16384, 2, 500)

		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnRows(usedRows)

		// 请求: CPU=8, Memory=16384, GPU=2, Storage=500
		// 已用+请求: CPU=16, Memory=32768, GPU=4, Storage=1000 (刚好等于总配额)
		ok, err := dao.CheckQuota(userID, 8, 16384, 2, 500)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QuotaNotFound", func(t *testing.T) {
		userID := uint(999)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		ok, err := dao.CheckQuota(userID, 4, 8192, 1, 200)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUsedResourcesError", func(t *testing.T) {
		userID := uint(1)
		now := time.Now()

		// Mock GetByUserID
		quotaRows := sqlmock.NewRows([]string{
			"id", "user_id", "workspace_id", "quota_level", "cpu", "memory",
			"gpu", "storage", "environment_quota", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, userID, nil, "pro", 16, 32768,
			4, 1000, 10, now, now, nil,
		)

		mock.ExpectQuery(`SELECT \* FROM "resource_quotas"`).
			WithArgs(userID, 1).
			WillReturnRows(quotaRows)

		// Mock GetUsedResources error
		mock.ExpectQuery(`SELECT COALESCE\(SUM\(cpu\), 0\) as cpu`).
			WillReturnError(sql.ErrConnDone)

		ok, err := dao.CheckQuota(userID, 4, 8192, 1, 200)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
