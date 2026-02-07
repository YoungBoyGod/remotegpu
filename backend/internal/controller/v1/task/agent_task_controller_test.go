package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testAgentToken = "test-agent-token-for-testing"

// agentTaskTestEnv Agent 任务接口测试环境
type agentTaskTestEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

// setupAgentTaskTestEnv 初始化 Agent 任务接口测试环境
func setupAgentTaskTestEnv(t *testing.T) *agentTaskTestEnv {
	gin.SetMode(gin.TestMode)

	config.GlobalConfig = &config.Config{
		Agent: config.AgentConfig{
			Enabled: true,
			Token:   testAgentToken,
		},
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 创建 tasks 表（包含 stdout/stderr/progress/progress_message 字段）
	err = db.Exec(`CREATE TABLE tasks (
		id VARCHAR(64) PRIMARY KEY,
		customer_id INTEGER NOT NULL,
		name VARCHAR(256) NOT NULL,
		type VARCHAR(32) NOT NULL DEFAULT 'shell',
		command TEXT NOT NULL DEFAULT '',
		args TEXT,
		work_dir VARCHAR(500),
		env_vars TEXT,
		timeout INTEGER DEFAULT 3600,
		priority INTEGER DEFAULT 5,
		retry_count INTEGER DEFAULT 0,
		retry_delay INTEGER DEFAULT 60,
		max_retries INTEGER DEFAULT 3,
		status VARCHAR(20) DEFAULT 'pending',
		exit_code INTEGER DEFAULT 0,
		error_msg TEXT,
		stdout TEXT,
		stderr TEXT,
		progress INTEGER DEFAULT 0,
		progress_message VARCHAR(500),
		machine_id VARCHAR(64),
		group_id VARCHAR(64),
		parent_id VARCHAR(64),
		assigned_agent_id VARCHAR(64),
		lease_expires_at DATETIME,
		attempt_id VARCHAR(64),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		assigned_at DATETIME,
		started_at DATETIME,
		ended_at DATETIME,
		host_id VARCHAR(64),
		process_id INTEGER DEFAULT 0,
		image_id INTEGER
	)`).Error
	require.NoError(t, err)

	// 创建 images 表（Task Preload Image 需要）
	err = db.Exec(`CREATE TABLE images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		name VARCHAR(128) NOT NULL DEFAULT '',
		tag VARCHAR(64) DEFAULT 'latest',
		registry VARCHAR(256),
		repository VARCHAR(256),
		description TEXT,
		status VARCHAR(20) DEFAULT 'active'
	)`).Error
	require.NoError(t, err)

	taskService := serviceTask.NewTaskService(db, nil)
	controller := NewAgentTaskController(taskService)

	router := gin.New()
	agentGroup := router.Group("/api/v1/agent")
	agentGroup.Use(middleware.AgentAuth())
	{
		agentGroup.POST("/tasks/claim", controller.ClaimTasks)
		agentGroup.POST("/tasks/:id/start", controller.StartTask)
		agentGroup.POST("/tasks/:id/lease/renew", controller.RenewLease)
		agentGroup.POST("/tasks/:id/complete", controller.CompleteTask)
		agentGroup.POST("/tasks/:id/progress", controller.ReportProgress)
	}

	return &agentTaskTestEnv{db: db, router: router}
}

// ==================== 任务认领测试 ====================

func TestClaimTasks_NoPendingTasks(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"machine_id": "host-001",
		"limit":      5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/claim", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	// 没有待认领任务，返回 null 或空列表
	tasks, ok := data["tasks"].([]any)
	if ok {
		assert.Len(t, tasks, 0)
	} else {
		assert.Nil(t, data["tasks"])
	}
}

func TestClaimTasks_MissingFields(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	reqBody := map[string]string{"agent_id": "agent-001"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/claim", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// ==================== 任务开始测试 ====================

func TestStartTask_Success(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	// 预设一个 assigned 状态的任务（模拟认领后的状态）
	env.db.Create(&entity.Task{
		ID: "t-start", CustomerID: 1, Name: "待启动任务",
		Command: "echo start", Status: "assigned",
		MachineID: "host-001", AssignedAgentID: "agent-001",
		AttemptID: "attempt-001",
	})

	reqBody := map[string]string{
		"agent_id":   "agent-001",
		"attempt_id": "attempt-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-start/start", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证状态已更新为 running
	var task entity.Task
	env.db.First(&task, "id = ?", "t-start")
	assert.Equal(t, "running", task.Status)
	assert.NotNil(t, task.StartedAt)
}

func TestStartTask_WrongAgent(t *testing.T) {
	// 使用错误的 agent_id 尝试启动任务
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-start-wrong", CustomerID: 1, Name: "任务",
		Command: "echo", Status: "assigned",
		AssignedAgentID: "agent-001", AttemptID: "attempt-001",
	})

	reqBody := map[string]string{
		"agent_id":   "agent-999",
		"attempt_id": "attempt-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-start-wrong/start", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 409, resp.Code)
}

// ==================== 续约租约测试 ====================

func TestRenewLease_Success(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-renew", CustomerID: 1, Name: "运行中任务",
		Command: "echo renew", Status: "running",
		AssignedAgentID: "agent-001", AttemptID: "attempt-001",
	})

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"attempt_id": "attempt-001",
		"extend_sec": 600,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-renew/lease/renew", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestRenewLease_WrongAttempt(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-renew-wrong", CustomerID: 1, Name: "任务",
		Command: "echo", Status: "running",
		AssignedAgentID: "agent-001", AttemptID: "attempt-001",
	})

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"attempt_id": "wrong-attempt",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-renew-wrong/lease/renew", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 410, resp.Code)
}

// ==================== 完成任务测试 ====================

func TestCompleteTask_Success(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-complete", CustomerID: 1, Name: "运行中任务",
		Command: "echo done", Status: "running",
		AssignedAgentID: "agent-001", AttemptID: "attempt-001",
	})

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"attempt_id": "attempt-001",
		"exit_code":  0,
		"stdout":     "hello world",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-complete/complete", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证状态已更新为 completed
	var task entity.Task
	env.db.First(&task, "id = ?", "t-complete")
	assert.Equal(t, "completed", task.Status)
	assert.Equal(t, 0, task.ExitCode)
	assert.NotNil(t, task.EndedAt)
}

func TestCompleteTask_Failed(t *testing.T) {
	// exit_code != 0 时状态应为 failed
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-fail", CustomerID: 1, Name: "失败任务",
		Command: "exit 1", Status: "running",
		AssignedAgentID: "agent-001", AttemptID: "attempt-002",
	})

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"attempt_id": "attempt-002",
		"exit_code":  1,
		"error":      "segfault",
		"stderr":     "Segmentation fault",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-fail/complete", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var task entity.Task
	env.db.First(&task, "id = ?", "t-fail")
	assert.Equal(t, "failed", task.Status)
	assert.Equal(t, 1, task.ExitCode)
}

// ==================== 进度上报测试 ====================

func TestReportProgress_Success(t *testing.T) {
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-progress", CustomerID: 1, Name: "进度任务",
		Command: "train.py", Status: "running",
		AssignedAgentID: "agent-001", AttemptID: "attempt-001",
	})

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"attempt_id": "attempt-001",
		"percent":    50,
		"message":    "训练进度 50%",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-progress/progress", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证进度已更新
	var task entity.Task
	env.db.First(&task, "id = ?", "t-progress")
	assert.Equal(t, 50, task.Progress)
	assert.Equal(t, "训练进度 50%", task.ProgressMessage)
}

func TestReportProgress_NotRunning(t *testing.T) {
	// 非 running 状态的任务不能上报进度
	env := setupAgentTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-prog-pend", CustomerID: 1, Name: "待处理任务",
		Command: "echo", Status: "pending",
		AssignedAgentID: "agent-001", AttemptID: "attempt-001",
	})

	reqBody := map[string]any{
		"agent_id":   "agent-001",
		"attempt_id": "attempt-001",
		"percent":    10,
		"message":    "进度",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/tasks/t-prog-pend/progress", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", testAgentToken)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 409, resp.Code)
}