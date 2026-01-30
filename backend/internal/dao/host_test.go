package dao

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupHostMockDB 创建Mock数据库连接
func setupHostMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
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

// TestHostDao_Create 测试创建主机
func TestHostDao_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		host := &entity.Host{
			ID:             "host-001",
			Name:           "Test Host",
			IPAddress:      "192.168.1.100",
			OSType:         "linux",
			DeploymentMode: "traditional",
			TotalCPU:       8,
			TotalMemory:    17179869184,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "hosts"`).
			WillReturnRows(sqlmock.NewRows([]string{"registered_at"}).AddRow(time.Now()))
		mock.ExpectCommit()

		err := dao.Create(host)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		host := &entity.Host{
			ID:             "host-002",
			Name:           "Test Host 2",
			IPAddress:      "192.168.1.101",
			OSType:         "linux",
			DeploymentMode: "traditional",
			TotalCPU:       4,
			TotalMemory:    8589934592,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "hosts"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Create(host)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_GetByID 测试根据ID获取主机
func TestHostDao_GetByID(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		hostID := "host-001"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "hostname", "ip_address", "public_ip",
			"os_type", "os_version", "arch", "deployment_mode", "k8s_node_name",
			"status", "health_status", "total_cpu", "total_memory", "total_disk",
			"total_gpu", "used_cpu", "used_memory", "used_disk", "used_gpu",
			"ssh_port", "winrm_port", "agent_port", "labels", "tags",
			"last_heartbeat", "registered_at", "created_at", "updated_at",
		}).AddRow(
			hostID, "Test Host", "testhost", "192.168.1.100", "1.2.3.4",
			"linux", "Ubuntu 22.04", "x86_64", "traditional", "",
			"online", "healthy", 8, 17179869184, 1099511627776,
			2, 0, 0, 0, 0,
			22, nil, 8080, []byte(`{}`), pq.StringArray{},
			&now, now, now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE id = .+`).
			WillReturnRows(rows)

		host, err := dao.GetByID(hostID)
		assert.NoError(t, err)
		assert.NotNil(t, host)
		assert.Equal(t, hostID, host.ID)
		assert.Equal(t, "Test Host", host.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE id = .+`).
			WillReturnError(gorm.ErrRecordNotFound)

		host, err := dao.GetByID("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, host)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE id = .+`).
			WillReturnError(sql.ErrConnDone)

		host, err := dao.GetByID("host-001")
		assert.Error(t, err)
		assert.Nil(t, host)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_GetByIPAddress 测试根据IP地址获取主机
func TestHostDao_GetByIPAddress(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		ipAddress := "192.168.1.100"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "hostname", "ip_address", "public_ip",
			"os_type", "os_version", "arch", "deployment_mode", "k8s_node_name",
			"status", "health_status", "total_cpu", "total_memory", "total_disk",
			"total_gpu", "used_cpu", "used_memory", "used_disk", "used_gpu",
			"ssh_port", "winrm_port", "agent_port", "labels", "tags",
			"last_heartbeat", "registered_at", "created_at", "updated_at",
		}).AddRow(
			"host-001", "Test Host", "testhost", ipAddress, "1.2.3.4",
			"linux", "Ubuntu 22.04", "x86_64", "traditional", "",
			"online", "healthy", 8, 17179869184, 1099511627776,
			2, 0, 0, 0, 0,
			22, nil, 8080, []byte(`{}`), pq.StringArray{},
			&now, now, now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE ip_address = .+`).
			WithArgs(ipAddress).
			WillReturnRows(rows)

		host, err := dao.GetByIPAddress(ipAddress)
		assert.NoError(t, err)
		assert.NotNil(t, host)
		assert.Equal(t, ipAddress, host.IPAddress)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		ipAddress := "192.168.1.200"

		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE ip_address = .+`).
			WithArgs(ipAddress).
			WillReturnError(gorm.ErrRecordNotFound)

		host, err := dao.GetByIPAddress(ipAddress)
		assert.Error(t, err)
		assert.Nil(t, host)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_Update 测试更新主机
func TestHostDao_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		host := &entity.Host{
			ID:             "host-001",
			Name:           "Updated Host",
			Hostname:       "updatedhost",
			IPAddress:      "192.168.1.100",
			OSType:         "linux",
			DeploymentMode: "traditional",
			TotalCPU:       16,
			TotalMemory:    34359738368,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "hosts"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Update(host)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		host := &entity.Host{
			ID:   "host-001",
			Name: "Updated Host",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "hosts"`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Update(host)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_Delete 测试删除主机
func TestHostDao_Delete(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "hosts" WHERE id = $1`).
			WithArgs(hostID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Delete(hostID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "hosts" WHERE id = $1`).
			WithArgs(hostID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Delete(hostID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_List 测试分页获取主机列表
func TestHostDao_List(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()

		// Mock Count query
		mock.ExpectQuery(`SELECT count(*) FROM "hosts"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock List query
		rows := sqlmock.NewRows([]string{
			"id", "name", "hostname", "ip_address", "public_ip",
			"os_type", "os_version", "arch", "deployment_mode", "k8s_node_name",
			"status", "health_status", "total_cpu", "total_memory", "total_disk",
			"total_gpu", "used_cpu", "used_memory", "used_disk", "used_gpu",
			"ssh_port", "winrm_port", "agent_port", "labels", "tags",
			"last_heartbeat", "registered_at", "created_at", "updated_at",
		}).
			AddRow("host-001", "Host 1", "host1", "192.168.1.100", "1.2.3.4",
				"linux", "Ubuntu 22.04", "x86_64", "traditional", "",
				"online", "healthy", 8, 17179869184, 1099511627776,
				2, 0, 0, 0, 0,
				22, nil, 8080, []byte(`{}`), pq.StringArray{},
				&now, now, now, now).
			AddRow("host-002", "Host 2", "host2", "192.168.1.101", "1.2.3.5",
				"linux", "Ubuntu 22.04", "x86_64", "traditional", "",
				"online", "healthy", 16, 34359738368, 2199023255552,
				4, 0, 0, 0, 0,
				22, nil, 8080, []byte(`{}`), pq.StringArray{},
				&now, now, now, now)

		mock.ExpectQuery(`SELECT * FROM "hosts" LIMIT $1 OFFSET $2`).
			WithArgs(10, 0).
			WillReturnRows(rows)

		hosts, total, err := dao.List(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, hosts, 2)
		assert.Equal(t, "host-001", hosts[0].ID)
		assert.Equal(t, "host-002", hosts[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count(*) FROM "hosts"`).
			WillReturnError(sql.ErrConnDone)

		hosts, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, hosts)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListError", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count(*) FROM "hosts"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`SELECT * FROM "hosts" LIMIT $1 OFFSET $2`).
			WithArgs(10, 0).
			WillReturnError(sql.ErrConnDone)

		hosts, total, err := dao.List(1, 10)
		assert.Error(t, err)
		assert.Nil(t, hosts)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_ListByStatus 测试根据状态获取主机列表
func TestHostDao_ListByStatus(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		status := "online"

		rows := sqlmock.NewRows([]string{
			"id", "name", "hostname", "ip_address", "public_ip",
			"os_type", "os_version", "arch", "deployment_mode", "k8s_node_name",
			"status", "health_status", "total_cpu", "total_memory", "total_disk",
			"total_gpu", "used_cpu", "used_memory", "used_disk", "used_gpu",
			"ssh_port", "winrm_port", "agent_port", "labels", "tags",
			"last_heartbeat", "registered_at", "created_at", "updated_at",
		}).AddRow(
			"host-001", "Host 1", "host1", "192.168.1.100", "1.2.3.4",
			"linux", "Ubuntu 22.04", "x86_64", "traditional", "",
			status, "healthy", 8, 17179869184, 1099511627776,
			2, 0, 0, 0, 0,
			22, nil, 8080, []byte(`{}`), pq.StringArray{},
			&now, now, now, now,
		)

		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE status = .+`).
			WithArgs(status).
			WillReturnRows(rows)

		hosts, err := dao.ListByStatus(status)
		assert.NoError(t, err)
		assert.Len(t, hosts, 1)
		assert.Equal(t, status, hosts[0].Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		status := "online"

		mock.ExpectQuery(`SELECT .* FROM "hosts" WHERE status = .+`).
			WithArgs(status).
			WillReturnError(sql.ErrConnDone)

		hosts, err := dao.ListByStatus(status)
		assert.Error(t, err)
		assert.Nil(t, hosts)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_UpdateStatus 测试更新主机状态
func TestHostDao_UpdateStatus(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		hostID := "host-001"
		status := "maintenance"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "hosts" SET "status"=$1,"updated_at"=$2 WHERE id = $3`).
			WithArgs(status, sqlmock.AnyArg(), hostID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.UpdateStatus(hostID, status)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		hostID := "host-001"
		status := "maintenance"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "hosts" SET "status"=$1,"updated_at"=$2 WHERE id = $3`).
			WithArgs(status, sqlmock.AnyArg(), hostID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.UpdateStatus(hostID, status)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestHostDao_UpdateHeartbeat 测试更新心跳时间
func TestHostDao_UpdateHeartbeat(t *testing.T) {
	gormDB, mock, sqlDB := setupHostMockDB(t)
	defer sqlDB.Close()

	dao := &HostDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "hosts" SET "last_heartbeat"=NOW(),"updated_at"=$1 WHERE id = $2`).
			WithArgs(sqlmock.AnyArg(), hostID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.UpdateHeartbeat(hostID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		hostID := "host-001"

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "hosts" SET "last_heartbeat"=NOW(),"updated_at"=$1 WHERE id = $2`).
			WithArgs(sqlmock.AnyArg(), hostID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.UpdateHeartbeat(hostID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
