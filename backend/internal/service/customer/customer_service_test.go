package customer

import (
	"context"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// insertTestCustomer 插入测试客户
func insertTestCustomer(t *testing.T, db *gorm.DB, c *entity.Customer) {
	if c.PasswordHash == "" {
		c.PasswordHash = "testhash"
	}
	err := db.Create(c).Error
	require.NoError(t, err)
}

// setupCustomerServiceTestDB 初始化测试数据库
func setupCustomerServiceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		uuid TEXT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
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

	err = db.Exec(`CREATE TABLE datasets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		uuid TEXT,
		customer_id INTEGER NOT NULL,
		workspace_id INTEGER,
		name TEXT NOT NULL,
		description TEXT,
		storage_path TEXT NOT NULL DEFAULT '',
		storage_type TEXT DEFAULT 'minio',
		total_size INTEGER DEFAULT 0,
		file_count INTEGER DEFAULT 0,
		status TEXT DEFAULT 'ready'
	)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
		username TEXT,
		action TEXT NOT NULL,
		resource_type TEXT,
		resource_id TEXT,
		ip_address TEXT,
		user_agent TEXT,
		request_method TEXT,
		request_path TEXT,
		status_code INTEGER,
		error_message TEXT,
		detail TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	return db
}

// ==================== CreateCustomer 测试 ====================

func TestCreateCustomer_Success(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	customer := &entity.Customer{
		Username: "testuser",
		Email:    "test@example.com",
	}
	err := svc.CreateCustomer(context.Background(), customer, "password123")
	require.NoError(t, err)

	// 验证密码已被哈希
	assert.NotEmpty(t, customer.PasswordHash)
	assert.NotEqual(t, "password123", customer.PasswordHash)

	// 验证数据库中已创建
	got, err := svc.GetCustomer(context.Background(), customer.ID)
	require.NoError(t, err)
	assert.Equal(t, "testuser", got.Username)
}

// ==================== ListCustomers 测试 ====================

func TestListCustomers_Empty(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	list, total, err := svc.ListCustomers(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}

func TestListCustomers_Pagination(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	for i := 0; i < 5; i++ {
		insertTestCustomer(t, db, &entity.Customer{
			Username: "user" + string(rune('a'+i)),
			Email:    "u" + string(rune('a'+i)) + "@test.com",
		})
	}

	list, total, err := svc.ListCustomers(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, list, 2)
}

// ==================== UpdateCustomer 测试 ====================

func TestUpdateCustomer_Success(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "upd", Email: "upd@test.com",
	})

	err := svc.UpdateCustomer(context.Background(), 1, map[string]interface{}{
		"display_name": "新名称",
	})
	require.NoError(t, err)

	got, err := svc.GetCustomer(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "新名称", got.DisplayName)
}

// ==================== UpdateStatus 测试 ====================

func TestUpdateStatus_Success(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "st", Email: "st@test.com", Status: "active",
	})

	err := svc.UpdateStatus(context.Background(), 1, "suspended")
	require.NoError(t, err)

	got, err := svc.GetCustomer(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "suspended", got.Status)
}

// ==================== CountActive 测试 ====================

func TestCountActive(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "a1", Email: "a1@test.com", Status: "active",
	})
	insertTestCustomer(t, db, &entity.Customer{
		Username: "a2", Email: "a2@test.com", Status: "active",
	})
	insertTestCustomer(t, db, &entity.Customer{
		Username: "s1", Email: "s1@test.com", Status: "suspended",
	})

	count, err := svc.CountActive(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// ==================== UpdateQuota 测试 ====================

func TestUpdateQuota_Success(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "q1", Email: "q1@test.com",
	})

	err := svc.UpdateQuota(context.Background(), 1, 8, 10240)
	require.NoError(t, err)

	got, err := svc.GetCustomer(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, 8, got.QuotaGPU)
	assert.Equal(t, int64(10240), got.QuotaStorage)
}

func TestUpdateQuota_CustomerNotFound(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	err := svc.UpdateQuota(context.Background(), 999, 4, 1024)
	assert.Error(t, err)
}

// ==================== CheckGPUQuota 测试 ====================

func TestCheckGPUQuota_Unlimited(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	// quota_gpu=0 表示不限制
	insertTestCustomer(t, db, &entity.Customer{
		Username: "unlim", Email: "unlim@test.com", QuotaGPU: 0,
	})

	err := svc.CheckGPUQuota(context.Background(), 1, 100)
	assert.NoError(t, err)
}

func TestCheckGPUQuota_WithinLimit(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "lim", Email: "lim@test.com", QuotaGPU: 4,
	})

	err := svc.CheckGPUQuota(context.Background(), 1, 2)
	assert.NoError(t, err)
}

func TestCheckGPUQuota_Exceeded(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "exc", Email: "exc@test.com", QuotaGPU: 2,
	})

	// 先创建一个活跃分配，关联 2 个 GPU
	db.Exec(`INSERT INTO hosts (id, name, ip_address) VALUES ('h1', 'host1', '10.0.0.1')`)
	db.Exec(`INSERT INTO allocations (id, customer_id, host_id, start_time, end_time, status) VALUES ('a1', 1, 'h1', datetime('now'), datetime('now', '+1 month'), 'active')`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 0, 'RTX4090', 24576)`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 1, 'RTX4090', 24576)`)

	// 已有 2 个 GPU，配额为 2，再申请 1 个应超限
	err := svc.CheckGPUQuota(context.Background(), 1, 1)
	assert.ErrorIs(t, err, ErrQuotaExceeded)
}

// ==================== CheckStorageQuota 测试 ====================

func TestCheckStorageQuota_Unlimited(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "su", Email: "su@test.com", QuotaStorage: 0,
	})

	err := svc.CheckStorageQuota(context.Background(), 1, 999999)
	assert.NoError(t, err)
}

func TestCheckStorageQuota_Exceeded(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	// 配额 100 MB
	insertTestCustomer(t, db, &entity.Customer{
		Username: "se", Email: "se@test.com", QuotaStorage: 100,
	})

	// 已有 80 MB 数据集（total_size 单位为字节）
	db.Exec(`INSERT INTO datasets (customer_id, name, storage_path, total_size, status)
		VALUES (1, 'ds1', '/data/ds1', ?, 'ready')`, 80*1024*1024)

	// 再申请 30 MB 应超限
	err := svc.CheckStorageQuota(context.Background(), 1, 30)
	assert.ErrorIs(t, err, ErrQuotaExceeded)
}

// ==================== GetResourceUsage 测试 ====================

func TestGetResourceUsage(t *testing.T) {
	db := setupCustomerServiceTestDB(t)
	svc := NewCustomerService(db)

	insertTestCustomer(t, db, &entity.Customer{
		Username: "ru", Email: "ru@test.com",
	})

	// 创建分配和 GPU
	db.Exec(`INSERT INTO hosts (id, name, ip_address) VALUES ('h1', 'host1', '10.0.0.1')`)
	db.Exec(`INSERT INTO allocations (id, customer_id, host_id, start_time, end_time, status)
		VALUES ('a1', 1, 'h1', datetime('now'), datetime('now', '+1 month'), 'active')`)
	db.Exec(`INSERT INTO gpus (host_id, "index", name, memory_total_mb) VALUES ('h1', 0, 'A100', 40960)`)

	usage, err := svc.GetResourceUsage(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), usage.AllocatedMachines)
	assert.Equal(t, int64(1), usage.AllocatedGPUs)
}
