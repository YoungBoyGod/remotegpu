package machine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testResponse 测试响应结构
type testResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// machineTestEnv 机器管理测试环境
type machineTestEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

// setupMachineTestEnv 初始化机器管理测试环境
func setupMachineTestEnv(t *testing.T) *machineTestEnv {
	gin.SetMode(gin.TestMode)

	// 使用 SQLite 内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 创建 hosts 表
	err = db.Exec(`CREATE TABLE hosts (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(128) NOT NULL,
		hostname VARCHAR(256),
		region VARCHAR(64) DEFAULT 'default',
		ip_address VARCHAR(64) NOT NULL,
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
		external_ip VARCHAR(64),
		external_ssh_port INTEGER DEFAULT 0,
		external_jupyter_port INTEGER DEFAULT 0,
		external_vnc_port INTEGER DEFAULT 0,
		nginx_domain VARCHAR(255),
		nginx_config_path VARCHAR(500),
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

	// 创建 gpus 表
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

	// 创建 allocations 表（MachineDao.FindByID 会 Preload）
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

	// 创建 customers 表（Allocation Preload Customer 需要）
	err = db.Exec(`CREATE TABLE customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		uuid TEXT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
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
		balance REAL DEFAULT 0,
		currency TEXT DEFAULT 'CNY',
		last_heartbeat DATETIME,
		must_change_password INTEGER DEFAULT 0
	)`).Error
	require.NoError(t, err)

	// 创建 host_metrics 表
	err = db.Exec(`CREATE TABLE host_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_id VARCHAR(64) NOT NULL,
		cpu_usage_percent REAL,
		memory_usage_percent REAL,
		disk_usage_percent REAL,
		collected_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	// 初始化服务和控制器（不注入 allocationService 和 agentService）
	machineService := serviceMachine.NewMachineService(db)
	controller := NewMachineController(machineService, nil, nil)

	// 设置路由
	router := gin.New()
	machineGroup := router.Group("/api/v1/admin/machines")
	{
		machineGroup.GET("", controller.List)
		machineGroup.POST("", controller.Create)
		machineGroup.GET("/:id", controller.Detail)
		machineGroup.PUT("/:id", controller.Update)
		machineGroup.DELETE("/:id", controller.Delete)
		machineGroup.POST("/:id/maintenance", controller.SetMaintenance)
		machineGroup.GET("/:id/usage", controller.Usage)
	}

	return &machineTestEnv{
		db:     db,
		router: router,
	}
}

// insertTestHost 插入测试主机数据
func insertTestHost(t *testing.T, db *gorm.DB, host *entity.Host) {
	err := db.Create(host).Error
	require.NoError(t, err)
}

// ==================== 创建机器测试 ====================

func TestCreateMachine_Success(t *testing.T) {
	env := setupMachineTestEnv(t)

	reqBody := apiV1.CreateMachineRequest{
		Name:        "测试机器",
		Hostname:    "test-node-01",
		Region:      "beijing",
		IPAddress:   "192.168.1.100",
		SSHPort:     22,
		SSHUsername: "root",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/machines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库中已创建
	var host entity.Host
	err = env.db.First(&host, "id = ?", "test-node-01").Error
	require.NoError(t, err)
	assert.Equal(t, "测试机器", host.Name)
	assert.Equal(t, "beijing", host.Region)
	assert.Equal(t, "192.168.1.100", host.IPAddress)
}

func TestCreateMachine_MissingAddress(t *testing.T) {
	// 测试缺少 IP 和 Hostname 时返回 400
	env := setupMachineTestEnv(t)

	reqBody := apiV1.CreateMachineRequest{
		Name:        "测试机器",
		Region:      "beijing",
		SSHPort:     22,
		SSHUsername: "root",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/machines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

func TestCreateMachine_DuplicateIP(t *testing.T) {
	// 测试重复 IP 地址返回 409
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID:        "existing-node",
		Name:      "已有机器",
		IPAddress: "192.168.1.100",
	})

	reqBody := apiV1.CreateMachineRequest{
		Name:        "新机器",
		Hostname:    "new-node",
		Region:      "beijing",
		IPAddress:   "192.168.1.100",
		SSHPort:     22,
		SSHUsername: "root",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/machines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 409, resp.Code)
}

func TestCreateMachine_MissingRequiredFields(t *testing.T) {
	// 测试缺少必填字段返回 400
	env := setupMachineTestEnv(t)

	reqBody := map[string]string{"hostname": "test-node"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/machines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// ==================== 列表查询测试 ====================

func TestListMachines_Empty(t *testing.T) {
	// 测试空列表返回
	env := setupMachineTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/machines", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]interface{}
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, float64(0), data["total"])
}

func TestListMachines_WithData(t *testing.T) {
	// 测试有数据时的列表返回
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "机器1", IPAddress: "10.0.0.1", Region: "beijing",
	})
	insertTestHost(t, env.db, &entity.Host{
		ID: "node-02", Name: "机器2", IPAddress: "10.0.0.2", Region: "shanghai",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/machines", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]interface{}
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, float64(2), data["total"])
}

func TestListMachines_Pagination(t *testing.T) {
	// 测试分页参数
	env := setupMachineTestEnv(t)

	for i := 0; i < 5; i++ {
		insertTestHost(t, env.db, &entity.Host{
			ID: fmt.Sprintf("node-%02d", i), Name: fmt.Sprintf("机器%d", i),
			IPAddress: fmt.Sprintf("10.0.0.%d", i+1),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/machines?page=1&page_size=2", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]interface{}
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, float64(5), data["total"])
	assert.Equal(t, float64(2), data["page_size"])

	list := data["list"].([]interface{})
	assert.Len(t, list, 2)
}

// ==================== 详情查询测试 ====================

func TestDetailMachine_Success(t *testing.T) {
	// 测试获取机器详情
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "测试机器", IPAddress: "10.0.0.1",
		Region: "beijing", TotalCPU: 8, TotalMemoryGB: 32,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/machines/node-01", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]interface{}
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, "node-01", data["id"])
	assert.Equal(t, "测试机器", data["name"])
}

func TestDetailMachine_NotFound(t *testing.T) {
	// 测试查询不存在的机器返回 404
	env := setupMachineTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/machines/nonexistent", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
}

// ==================== 更新机器测试 ====================

func TestUpdateMachine_Success(t *testing.T) {
	// 测试更新机器信息
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "旧名称", IPAddress: "10.0.0.1", Region: "beijing",
	})

	reqBody := apiV1.UpdateMachineRequest{
		Name:   "新名称",
		Region: "shanghai",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/machines/node-01", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库已更新
	var host entity.Host
	err = env.db.First(&host, "id = ?", "node-01").Error
	require.NoError(t, err)
	assert.Equal(t, "新名称", host.Name)
	assert.Equal(t, "shanghai", host.Region)
}

func TestUpdateMachine_NoFields(t *testing.T) {
	// 测试没有更新字段时返回 400
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "机器", IPAddress: "10.0.0.1",
	})

	reqBody := apiV1.UpdateMachineRequest{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/machines/node-01", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// ==================== 删除机器测试 ====================

func TestDeleteMachine_Success(t *testing.T) {
	// 测试删除机器
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "待删除机器", IPAddress: "10.0.0.1",
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/machines/node-01", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库中已删除
	var count int64
	env.db.Model(&entity.Host{}).Where("id = ?", "node-01").Count(&count)
	assert.Equal(t, int64(0), count)
}

// ==================== 维护状态测试 ====================

func TestSetMaintenance_Enable(t *testing.T) {
	// 测试设置机器为维护状态
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "机器", IPAddress: "10.0.0.1",
		AllocationStatus: "idle",
	})

	reqBody := map[string]bool{"maintenance": true}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/machines/node-01/maintenance", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证分配状态已更新为 maintenance
	var host entity.Host
	env.db.First(&host, "id = ?", "node-01")
	assert.Equal(t, "maintenance", host.AllocationStatus)
}

func TestSetMaintenance_Disable(t *testing.T) {
	// 测试取消维护状态
	env := setupMachineTestEnv(t)

	insertTestHost(t, env.db, &entity.Host{
		ID: "node-01", Name: "机器", IPAddress: "10.0.0.1",
		AllocationStatus: "maintenance",
	})

	reqBody := map[string]bool{"maintenance": false}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/machines/node-01/maintenance", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var host entity.Host
	env.db.First(&host, "id = ?", "node-01")
	assert.Equal(t, "idle", host.AllocationStatus)
}
