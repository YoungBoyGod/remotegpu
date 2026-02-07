package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
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

const testAgentToken = "test-agent-token-for-testing"

// agentTestEnv Agent 通信测试环境
type agentTestEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

// setupAgentTestEnv 初始化 Agent 通信测试环境
func setupAgentTestEnv(t *testing.T) *agentTestEnv {
	gin.SetMode(gin.TestMode)

	// 设置全局配置，用于 AgentAuth 中间件验证 Token
	config.GlobalConfig = &config.Config{
		Agent: config.AgentConfig{
			Enabled: true,
			Token:   testAgentToken,
		},
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 创建 hosts 表
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

	// 创建 host_metrics 表（心跳携带指标时写入）
	err = db.Exec(`CREATE TABLE host_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_id VARCHAR(255) NOT NULL,
		cpu_usage_percent REAL,
		cpu_cores_used REAL,
		memory_total_gb INTEGER,
		memory_used_gb INTEGER,
		memory_usage_percent REAL,
		disk_total_gb INTEGER,
		disk_used_gb INTEGER,
		disk_usage_percent REAL,
		network_rx_bytes INTEGER,
		network_tx_bytes INTEGER,
		collected_at DATETIME NOT NULL
	)`).Error
	require.NoError(t, err)

	// 创建测试机器
	err = db.Exec(`INSERT INTO hosts (id, name, ip_address, total_cpu, total_memory_gb, device_status)
		VALUES ('host-001', '测试机器1', '192.168.1.100', 8, 32, 'offline')`).Error
	require.NoError(t, err)

	machineSvc := serviceMachine.NewMachineService(db)
	controller := NewHeartbeatController(machineSvc)

	router := gin.New()
	agentGroup := router.Group("/api/v1/agent")
	agentGroup.Use(middleware.AgentAuth())
	{
		agentGroup.POST("/register", controller.Register)
		agentGroup.POST("/heartbeat", controller.Heartbeat)
	}

	return &agentTestEnv{db: db, router: router}
}

// ==================== Agent 认证测试 ====================

func TestAgentAuth_NoToken(t *testing.T) {
	env := setupAgentTestEnv(t)

	reqBody := apiV1.HeartbeatRequest{
		AgentID: "agent-001", MachineID: "host-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	// AgentAuth 中间件返回 HTTP 401
	assert.Equal(t, http.StatusOK, w.Code)
	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAgentAuth_InvalidToken(t *testing.T) {
	env := setupAgentTestEnv(t)

	reqBody := apiV1.HeartbeatRequest{
		AgentID: "agent-001", MachineID: "host-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", "wrong-token")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAgentAuth_ValidToken_XAgentToken(t *testing.T) {
	env := setupAgentTestEnv(t)

	reqBody := apiV1.HeartbeatRequest{
		AgentID: "agent-001", MachineID: "host-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestAgentAuth_ValidToken_BearerAuth(t *testing.T) {
	// 验证 Authorization: Bearer 方式也能通过认证
	env := setupAgentTestEnv(t)

	reqBody := apiV1.HeartbeatRequest{
		AgentID: "agent-001", MachineID: "host-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

// ==================== 心跳测试 ====================

func TestHeartbeat_Success(t *testing.T) {
	env := setupAgentTestEnv(t)

	reqBody := apiV1.HeartbeatRequest{
		AgentID: "agent-001", MachineID: "host-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库中设备状态已更新为 online
	var host entity.Host
	env.db.First(&host, "id = ?", "host-001")
	assert.Equal(t, "online", host.DeviceStatus)
	assert.NotNil(t, host.LastHeartbeat)
}

func TestHeartbeat_WithMetrics(t *testing.T) {
	env := setupAgentTestEnv(t)

	cpuUsage := 45.5
	memUsage := 60.2
	reqBody := apiV1.HeartbeatRequest{
		AgentID:   "agent-001",
		MachineID: "host-001",
		Metrics: &apiV1.HeartbeatMetrics{
			CPUUsagePercent:    &cpuUsage,
			MemoryUsagePercent: &memUsage,
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证 host_metrics 表中写入了指标数据
	var count int64
	env.db.Table("host_metrics").Where("host_id = ?", "host-001").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestHeartbeat_MissingFields(t *testing.T) {
	env := setupAgentTestEnv(t)

	// 缺少 required 字段
	reqBody := map[string]string{"agent_id": "agent-001"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

func TestHeartbeat_UnknownMachine(t *testing.T) {
	// 心跳上报不存在的机器，UpdateHeartbeat 不会报错（RowsAffected=0 但无 error）
	env := setupAgentTestEnv(t)

	reqBody := apiV1.HeartbeatRequest{
		AgentID: "agent-001", MachineID: "nonexistent",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/heartbeat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	// UpdateHeartbeat 使用 Updates 不检查 RowsAffected，返回成功
	assert.Equal(t, 0, resp.Code)
}

// ==================== Agent 注册测试 ====================

func TestRegister_Success(t *testing.T) {
	env := setupAgentTestEnv(t)

	reqBody := apiV1.RegisterRequest{
		AgentID:   "agent-001",
		MachineID: "host-001",
		Version:   "1.0.0",
		Hostname:  "gpu-node-1",
		AgentPort: 9090,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库中机器信息已更新
	var host entity.Host
	env.db.First(&host, "id = ?", "host-001")
	assert.Equal(t, "online", host.DeviceStatus)
	assert.Equal(t, "gpu-node-1", host.Hostname)
	assert.Equal(t, 9090, host.AgentPort)
}

func TestRegister_MissingFields(t *testing.T) {
	env := setupAgentTestEnv(t)

	// 缺少 required 字段 machine_id
	reqBody := map[string]string{"agent_id": "agent-001"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}
