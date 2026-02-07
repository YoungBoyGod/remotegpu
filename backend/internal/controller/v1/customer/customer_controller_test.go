package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceCustomer "github.com/YoungBoyGod/remotegpu/internal/service/customer"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// customerTestEnv 客户管理测试环境
type customerTestEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

// setupCustomerTestEnv 初始化客户管理测试环境
func setupCustomerTestEnv(t *testing.T) *customerTestEnv {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 创建 customers 表
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
		last_login_at DATETIME
	)`).Error
	require.NoError(t, err)

	// 创建 allocations 表（GetCustomerDetail 需要）
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

	// 创建 hosts 表（Allocation Preload Host 需要）
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

	customerService := serviceCustomer.NewCustomerService(db)
	controller := NewCustomerController(customerService)

	router := gin.New()
	group := router.Group("/api/v1/admin/customers")
	{
		group.GET("", controller.List)
		group.POST("", controller.Create)
		group.GET("/:id", controller.Detail)
		group.PUT("/:id", controller.Update)
		group.POST("/:id/disable", controller.Disable)
		group.POST("/:id/enable", controller.Enable)
	}

	return &customerTestEnv{db: db, router: router}
}

// ==================== 创建客户测试 ====================

func TestCreateCustomer_Success(t *testing.T) {
	env := setupCustomerTestEnv(t)

	reqBody := apiV1.CreateCustomerRequest{
		Username:    "newcustomer",
		Email:       "new@example.com",
		CompanyCode: "COMP001",
		Password:    "Secure123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证数据库中已创建
	var customer entity.Customer
	err = env.db.First(&customer, "username = ?", "newcustomer").Error
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", customer.Email)
	assert.Equal(t, "customer_owner", customer.Role)
}

func TestCreateCustomer_DefaultPassword(t *testing.T) {
	// 测试不提供密码时使用默认密码，且标记需要修改密码
	env := setupCustomerTestEnv(t)

	reqBody := apiV1.CreateCustomerRequest{
		Username:    "defaultpw",
		Email:       "default@example.com",
		CompanyCode: "COMP001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	var customer entity.Customer
	err = env.db.First(&customer, "username = ?", "defaultpw").Error
	require.NoError(t, err)
	assert.True(t, customer.MustChangePassword)
}

func TestCreateCustomer_MissingRequired(t *testing.T) {
	// 测试缺少必填字段返回 400
	env := setupCustomerTestEnv(t)

	reqBody := map[string]string{"username": "onlyname"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(400), resp["code"])
}

// ==================== 客户列表测试 ====================

func TestListCustomers_Empty(t *testing.T) {
	// 测试空列表返回
	env := setupCustomerTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/customers", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])
}

func TestListCustomers_WithData(t *testing.T) {
	// 测试有数据时的列表返回
	env := setupCustomerTestEnv(t)

	env.db.Create(&entity.Customer{
		Username: "user1", Email: "u1@test.com",
		PasswordHash: "hash1", CompanyCode: "C1",
	})
	env.db.Create(&entity.Customer{
		Username: "user2", Email: "u2@test.com",
		PasswordHash: "hash2", CompanyCode: "C2",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/customers", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
}

func TestListCustomers_Pagination(t *testing.T) {
	// 测试分页
	env := setupCustomerTestEnv(t)

	for i := 0; i < 5; i++ {
		env.db.Create(&entity.Customer{
			Username:     fmt.Sprintf("user%d", i),
			Email:        fmt.Sprintf("u%d@test.com", i),
			PasswordHash: "hash",
			CompanyCode:  "C1",
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/customers?page=1&page_size=2", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(5), data["total"])
	list := data["list"].([]interface{})
	assert.Len(t, list, 2)
}

// ==================== 客户详情测试 ====================

func TestDetailCustomer_Success(t *testing.T) {
	env := setupCustomerTestEnv(t)

	env.db.Create(&entity.Customer{
		Username: "detail_user", Email: "detail@test.com",
		PasswordHash: "hash", CompanyCode: "C1",
		DisplayName: "详情用户",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/customers/1", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestDetailCustomer_NotFound(t *testing.T) {
	env := setupCustomerTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/customers/999", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(404), resp["code"])
}

func TestDetailCustomer_InvalidID(t *testing.T) {
	env := setupCustomerTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/customers/abc", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(400), resp["code"])
}

// ==================== 更新客户测试 ====================

func TestUpdateCustomer_Success(t *testing.T) {
	env := setupCustomerTestEnv(t)

	env.db.Create(&entity.Customer{
		Username: "update_user", Email: "old@test.com",
		PasswordHash: "hash", CompanyCode: "C1",
	})

	reqBody := apiV1.UpdateCustomerRequest{
		Email:       "new@test.com",
		DisplayName: "新名称",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证数据库已更新
	var customer entity.Customer
	env.db.First(&customer, 1)
	assert.Equal(t, "new@test.com", customer.Email)
	assert.Equal(t, "新名称", customer.DisplayName)
}

func TestUpdateCustomer_NoFields(t *testing.T) {
	// 测试没有更新字段时返回 400
	env := setupCustomerTestEnv(t)

	env.db.Create(&entity.Customer{
		Username: "noupdate", Email: "no@test.com",
		PasswordHash: "hash", CompanyCode: "C1",
	})

	reqBody := apiV1.UpdateCustomerRequest{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(400), resp["code"])
}

// ==================== 禁用/启用客户测试 ====================

func TestDisableCustomer_Success(t *testing.T) {
	env := setupCustomerTestEnv(t)

	env.db.Create(&entity.Customer{
		Username: "to_disable", Email: "dis@test.com",
		PasswordHash: "hash", CompanyCode: "C1", Status: "active",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/customers/1/disable", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证状态已更新为 suspended
	var customer entity.Customer
	env.db.First(&customer, 1)
	assert.Equal(t, "suspended", customer.Status)
}

func TestEnableCustomer_Success(t *testing.T) {
	env := setupCustomerTestEnv(t)

	env.db.Create(&entity.Customer{
		Username: "to_enable", Email: "en@test.com",
		PasswordHash: "hash", CompanyCode: "C1", Status: "suspended",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/customers/1/enable", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证状态已更新为 active
	var customer entity.Customer
	env.db.First(&customer, 1)
	assert.Equal(t, "active", customer.Status)
}

func TestEnableCustomer_InvalidID(t *testing.T) {
	env := setupCustomerTestEnv(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/customers/abc/enable", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(400), resp["code"])
}
