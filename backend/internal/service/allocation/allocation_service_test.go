package allocation

import (
	"context"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAllocationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		uuid TEXT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL DEFAULT '',
		display_name TEXT,
		full_name TEXT,
		company_code TEXT,
		company TEXT,
		phone TEXT,
		avatar_url TEXT,
		role TEXT DEFAULT 'customer_owner',
		user_type TEXT DEFAULT 'external',
		account_type TEXT DEFAULT 'individual',
		status TEXT DEFAULT 'active',
		email_verified INTEGER DEFAULT 0,
		phone_verified INTEGER DEFAULT 0,
		must_change_password INTEGER DEFAULT 0,
		quota_gpu INTEGER DEFAULT 0,
		quota_storage INTEGER DEFAULT 0,
		balance REAL DEFAULT 0,
		currency TEXT DEFAULT 'CNY',
		credit_limit REAL DEFAULT 0,
		billing_plan_id INTEGER,
		last_login_at DATETIME
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE hosts (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(128) NOT NULL DEFAULT '',
		hostname VARCHAR(256),
		region VARCHAR(64) DEFAULT 'default',
		ip_address VARCHAR(64) NOT NULL DEFAULT '',
		public_ip VARCHAR(64),
		ssh_host VARCHAR(255),
		ssh_port INTEGER DEFAULT 22,
		agent_port INTEGER DEFAULT 8080,
		ssh_username VARCHAR(128) DEFAULT 'root',
		ssh_password TEXT,
		ssh_key TEXT,
		jupyter_url VARCHAR(255),
		jupyter_token VARCHAR(255),
		vnc_url VARCHAR(255),
		vnc_password VARCHAR(255),
		os_type VARCHAR(20) DEFAULT 'linux',
		os_version VARCHAR(64),
		cpu_info VARCHAR(256),
		total_cpu INTEGER NOT NULL DEFAULT 0,
		total_memory_gb INTEGER NOT NULL DEFAULT 0,
		total_disk_gb INTEGER DEFAULT 0,
		status VARCHAR(20) DEFAULT 'offline',
		device_status VARCHAR(20) DEFAULT 'offline',
		allocation_status VARCHAR(20) DEFAULT 'idle',
		health_status VARCHAR(20) DEFAULT 'unknown',
		deployment_mode VARCHAR(20) DEFAULT 'traditional',
		needs_collect INTEGER DEFAULT 0,
		last_heartbeat DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE allocations (
		id VARCHAR(64) PRIMARY KEY,
		customer_id INTEGER NOT NULL,
		host_id VARCHAR(64) NOT NULL,
		workspace_id INTEGER,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		actual_end_time DATETIME,
		status VARCHAR(32) DEFAULT 'active',
		remark TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE gpus (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_id VARCHAR(64) NOT NULL,
		"index" INTEGER NOT NULL,
		uuid VARCHAR(128),
		name VARCHAR(128) NOT NULL,
		memory_total_mb INTEGER NOT NULL,
		brand VARCHAR(64),
		status VARCHAR(20) DEFAULT 'available',
		health_status VARCHAR(20) DEFAULT 'healthy',
		allocated_to VARCHAR(64),
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE ssh_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER NOT NULL,
		name VARCHAR(64) NOT NULL,
		fingerprint VARCHAR(128),
		public_key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
		machine_id VARCHAR(64),
		status VARCHAR(32) DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
		username TEXT,
		ip_address TEXT,
		method VARCHAR(10),
		path VARCHAR(512),
		action TEXT NOT NULL,
		resource_type TEXT,
		resource_id TEXT,
		detail TEXT,
		status_code INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	return db
}

// newTestAllocationService 创建测试用分配服务（不注入 AgentClient 和 Redis）
func newTestAllocationService(_ *testing.T, db *gorm.DB) *AllocationService {
	auditSvc := audit.NewAuditService(db)
	return NewAllocationService(db, auditSvc, nil)
}

// ==================== checkGPUQuota 测试 ====================

func TestCheckGPUQuota_Unlimited(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	// quota_gpu=0 表示不限制
	db.Exec(`INSERT INTO customers (username, email, quota_gpu) VALUES ('u1', 'u1@test.com', 0)`)

	err := svc.checkGPUQuota(context.Background(), 1, "h1")
	assert.NoError(t, err)
}

func TestCheckGPUQuota_WithinLimit(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	// 配额 4 个 GPU
	db.Exec(`INSERT INTO customers (username, email, quota_gpu) VALUES ('u1', 'u1@test.com', 4)`)
	// 目标机器有 2 个 GPU
	db.Exec(`INSERT INTO hosts (id, name, ip_address, allocation_status) VALUES ('h1', 'host1', '10.0.0.1', 'idle')`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 0, 'A100', 40960)`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 1, 'A100', 40960)`)

	err := svc.checkGPUQuota(context.Background(), 1, "h1")
	assert.NoError(t, err)
}

func TestCheckGPUQuota_Exceeded(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	db.Exec(`INSERT INTO customers (username, email, quota_gpu) VALUES ('u1', 'u1@test.com', 2)`)
	db.Exec(`INSERT INTO hosts (id, name, ip_address, allocation_status) VALUES ('h1', 'host1', '10.0.0.1', 'allocated')`)
	db.Exec(`INSERT INTO allocations (id, customer_id, host_id, start_time, end_time, status) VALUES ('a1', 1, 'h1', datetime('now'), datetime('now', '+1 month'), 'active')`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 0, 'A100', 40960)`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 1, 'A100', 40960)`)
	db.Exec(`INSERT INTO hosts (id, name, ip_address, allocation_status) VALUES ('h2', 'host2', '10.0.0.2', 'idle')`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h2', 0, 'A100', 40960)`)

	// 已有 2 个 GPU，配额 2，再分配 1 个应超限
	err := svc.checkGPUQuota(context.Background(), 1, "h2")
	assert.Error(t, err)
}

// ==================== AllocateMachine 测试 ====================

func TestAllocateMachine_Success(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	db.Exec(`INSERT INTO customers (username, email, quota_gpu) VALUES ('u1', 'u1@test.com', 0)`)
	db.Exec(`INSERT INTO hosts (id, name, ip_address, allocation_status) VALUES ('h1', 'host1', '10.0.0.1', 'idle')`)

	alloc, err := svc.AllocateMachine(context.Background(), 1, "h1", 3, "测试分配")
	require.NoError(t, err)
	assert.NotEmpty(t, alloc.ID)
	assert.Equal(t, uint(1), alloc.CustomerID)
	assert.Equal(t, "h1", alloc.HostID)
	assert.Equal(t, "active", alloc.Status)

	// 验证机器状态已更新
	var host entity.Host
	db.First(&host, "id = ?", "h1")
	assert.Equal(t, "allocated", host.AllocationStatus)
}

func TestAllocateMachine_InvalidDuration(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	_, err := svc.AllocateMachine(context.Background(), 1, "h1", 0, "")
	assert.Error(t, err)
}

func TestAllocateMachine_MachineNotFound(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	db.Exec(`INSERT INTO customers (username, email) VALUES ('u1', 'u1@test.com')`)

	_, err := svc.AllocateMachine(context.Background(), 1, "nonexistent", 1, "")
	assert.Error(t, err)
}

func TestAllocateMachine_MachineNotIdle(t *testing.T) {
	db := setupAllocationTestDB(t)
	svc := newTestAllocationService(t, db)

	db.Exec(`INSERT INTO customers (username, email) VALUES ('u1', 'u1@test.com')`)
	db.Exec(`INSERT INTO hosts (id, name, ip_address, allocation_status) VALUES ('h1', 'host1', '10.0.0.1', 'allocated')`)

	_, err := svc.AllocateMachine(context.Background(), 1, "h1", 1, "")
	assert.Error(t, err)
}
