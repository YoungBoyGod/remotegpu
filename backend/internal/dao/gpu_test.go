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

// setupGPUMockDB 创建Mock数据库连接
func setupGPUMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
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

// TestGPUDao_Create 测试创建GPU
func TestGPUDao_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpu := &entity.GPU{
			HostID:            "host-001",
			GPUIndex:          0,
			UUID:              "GPU-12345678",
			Name:              "Tesla V100",
			Brand:             "NVIDIA",
			Architecture:      "Volta",
			MemoryTotal:       34359738368,
			CUDACores:         5120,
			ComputeCapability: "7.0",
			Status:            "available",
			HealthStatus:      "healthy",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "gpus"`).
			WithArgs(
				gpu.HostID,
				gpu.GPUIndex,
				gpu.UUID,
				gpu.Name,
				gpu.Brand,
				gpu.Architecture,
				gpu.MemoryTotal,
				gpu.CUDACores,
				gpu.ComputeCapability,
				gpu.Status,
				gpu.HealthStatus,
				gpu.AllocatedTo,
				gpu.AllocatedAt,
				gpu.PowerLimit,
				gpu.TemperatureLimit,
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := dao.Create(gpu)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpu := &entity.GPU{
			HostID:   "host-001",
			GPUIndex: 0,
			Name:     "Tesla V100",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "gpus"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Create(gpu)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_GetByID 测试根据ID获取GPU
func TestGPUDao_GetByID(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpuID := uint(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "host_id", "gpu_index", "uuid", "name", "brand",
			"architecture", "memory_total", "cuda_cores", "compute_capability",
			"status", "health_status", "allocated_to", "allocated_at",
			"power_limit", "temperature_limit", "created_at", "updated_at",
		}).AddRow(
			gpuID, "host-001", 0, "GPU-12345678", "Tesla V100", "NVIDIA",
			"Volta", 34359738368, 5120, "7.0",
			"available", "healthy", "", nil,
			250, 85, now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE id = .+`).
			WithArgs(gpuID).
			WillReturnRows(rows)

		gpu, err := dao.GetByID(gpuID)
		assert.NoError(t, err)
		assert.NotNil(t, gpu)
		assert.Equal(t, gpuID, gpu.ID)
		assert.Equal(t, "Tesla V100", gpu.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		gpuID := uint(999)

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE id = .+`).
			WithArgs(gpuID).
			WillReturnError(gorm.ErrRecordNotFound)

		gpu, err := dao.GetByID(gpuID)
		assert.Error(t, err)
		assert.Nil(t, gpu)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpuID := uint(1)

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE id = .+`).
			WithArgs(gpuID).
			WillReturnError(sql.ErrConnDone)

		gpu, err := dao.GetByID(gpuID)
		assert.Error(t, err)
		assert.Nil(t, gpu)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_GetByHostID 测试根据主机ID获取GPU列表
func TestGPUDao_GetByHostID(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		hostID := "host-001"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "host_id", "gpu_index", "uuid", "name", "brand",
			"architecture", "memory_total", "cuda_cores", "compute_capability",
			"status", "health_status", "allocated_to", "allocated_at",
			"power_limit", "temperature_limit", "created_at", "updated_at",
		}).
			AddRow(1, hostID, 0, "GPU-001", "Tesla V100", "NVIDIA",
				"Volta", 34359738368, 5120, "7.0",
				"available", "healthy", "", nil,
				250, 85, now, now).
			AddRow(2, hostID, 1, "GPU-002", "Tesla V100", "NVIDIA",
				"Volta", 34359738368, 5120, "7.0",
				"available", "healthy", "", nil,
				250, 85, now, now)

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE host_id = .+`).
			WithArgs(hostID).
			WillReturnRows(rows)

		gpus, err := dao.GetByHostID(hostID)
		assert.NoError(t, err)
		assert.Len(t, gpus, 2)
		assert.Equal(t, hostID, gpus[0].HostID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE host_id = .+`).
			WithArgs(hostID).
			WillReturnError(sql.ErrConnDone)

		gpus, err := dao.GetByHostID(hostID)
		assert.Error(t, err)
		assert.Nil(t, gpus)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_Update 测试更新GPU
func TestGPUDao_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpu := &entity.GPU{
			ID:       1,
			HostID:   "host-001",
			GPUIndex: 0,
			Name:     "Updated GPU",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Update(gpu)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpu := &entity.GPU{
			ID:   1,
			Name: "Updated GPU",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Update(gpu)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_Delete 测试删除GPU
func TestGPUDao_Delete(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpuID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "gpus" WHERE "gpus"."id" = $1`).
			WithArgs(gpuID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Delete(gpuID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpuID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "gpus" WHERE "gpus"."id" = $1`).
			WithArgs(gpuID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Delete(gpuID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_DeleteByHostID 测试根据主机ID删除所有GPU
func TestGPUDao_DeleteByHostID(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "gpus" WHERE host_id = $1`).
			WithArgs(hostID).
			WillReturnResult(sqlmock.NewResult(1, 2))
		mock.ExpectCommit()

		err := dao.DeleteByHostID(hostID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "gpus" WHERE host_id = $1`).
			WithArgs(hostID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.DeleteByHostID(hostID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_UpdateStatus 测试更新GPU状态
func TestGPUDao_UpdateStatus(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpuID := uint(1)
		status := "maintenance"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus" SET "status"=$1,"updated_at"=$2 WHERE id = $3`).
			WithArgs(status, sqlmock.AnyArg(), gpuID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.UpdateStatus(gpuID, status)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpuID := uint(1)
		status := "maintenance"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus" SET "status"=$1,"updated_at"=$2 WHERE id = $3`).
			WithArgs(status, sqlmock.AnyArg(), gpuID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.UpdateStatus(gpuID, status)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_List 测试分页获取GPU列表
func TestGPUDao_List(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()

		// Mock Count query
		mock.ExpectQuery(`SELECT count(*) FROM "gpus"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock List query
		rows := sqlmock.NewRows([]string{
			"id", "host_id", "gpu_index", "uuid", "name", "brand",
			"architecture", "memory_total", "cuda_cores", "compute_capability",
			"status", "health_status", "allocated_to", "allocated_at",
			"power_limit", "temperature_limit", "created_at", "updated_at",
		}).
			AddRow(2, "host-002", 0, "GPU-002", "RTX 4090", "NVIDIA",
				"Ada Lovelace", 25769803776, 16384, "8.9",
				"available", "healthy", "", nil,
				450, 90, now, now).
			AddRow(1, "host-001", 0, "GPU-001", "Tesla V100", "NVIDIA",
				"Volta", 34359738368, 5120, "7.0",
				"available", "healthy", "", nil,
				250, 85, now, now)

		mock.ExpectQuery(`SELECT * FROM "gpus" ORDER BY id DESC LIMIT $1 OFFSET $2`).
			WithArgs(10, 0).
			WillReturnRows(rows)

		gpus, total, err := dao.List(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, gpus, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count(*) FROM "gpus"`).
			WillReturnError(sql.ErrConnDone)

		gpus, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, gpus)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_GetByStatus 测试根据状态获取GPU列表
func TestGPUDao_GetByStatus(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		status := "available"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "host_id", "gpu_index", "uuid", "name", "brand",
			"architecture", "memory_total", "cuda_cores", "compute_capability",
			"status", "health_status", "allocated_to", "allocated_at",
			"power_limit", "temperature_limit", "created_at", "updated_at",
		}).AddRow(
			1, "host-001", 0, "GPU-001", "Tesla V100", "NVIDIA",
			"Volta", 34359738368, 5120, "7.0",
			status, "healthy", "", nil,
			250, 85, now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE status = .+`).
			WithArgs(status).
			WillReturnRows(rows)

		gpus, err := dao.GetByStatus(status)
		assert.NoError(t, err)
		assert.Len(t, gpus, 1)
		assert.Equal(t, status, gpus[0].Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		status := "available"

		mock.ExpectQuery(`SELECT .* FROM "gpus" WHERE status = .+`).
			WithArgs(status).
			WillReturnError(sql.ErrConnDone)

		gpus, err := dao.GetByStatus(status)
		assert.Error(t, err)
		assert.Nil(t, gpus)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_Allocate 测试分配GPU
func TestGPUDao_Allocate(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpuID := uint(1)
		allocatedTo := "env-001"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus" SET "allocated_at"=NOW(),"allocated_to"=$1,"status"=$2,"updated_at"=$3 WHERE id = $4`).
			WithArgs(allocatedTo, "allocated", sqlmock.AnyArg(), gpuID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Allocate(gpuID, allocatedTo)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpuID := uint(1)
		allocatedTo := "env-001"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus" SET "allocated_at"=NOW(),"allocated_to"=$1,"status"=$2,"updated_at"=$3 WHERE id = $4`).
			WithArgs(allocatedTo, "allocated", sqlmock.AnyArg(), gpuID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Allocate(gpuID, allocatedTo)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGPUDao_Release 测试释放GPU
func TestGPUDao_Release(t *testing.T) {
	gormDB, mock, sqlDB := setupGPUMockDB(t)
	defer sqlDB.Close()

	dao := &GPUDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		gpuID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus" SET "allocated_at"=$1,"allocated_to"=$2,"status"=$3,"updated_at"=$4 WHERE id = $5`).
			WithArgs(nil, "", "available", sqlmock.AnyArg(), gpuID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Release(gpuID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		gpuID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "gpus" SET "allocated_at"=$1,"allocated_to"=$2,"status"=$3,"updated_at"=$4 WHERE id = $5`).
			WithArgs(nil, "", "available", sqlmock.AnyArg(), gpuID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Release(gpuID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
