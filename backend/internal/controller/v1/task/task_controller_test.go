package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"
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

// taskTestEnv 任务管理测试环境
type taskTestEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

// setupTaskTestEnv 初始化任务管理测试环境
func setupTaskTestEnv(t *testing.T) *taskTestEnv {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 创建 tasks 表
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
		progress_message TEXT,
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

	// agentService 传 nil，任务停止相关测试会走 "agent service unavailable" 分支
	taskService := serviceTask.NewTaskService(db, nil)
	controller := NewTaskController(taskService)

	router := gin.New()

	// 模拟认证中间件：从 header 中读取 userID 并设置到 context
	authMiddleware := func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "用户未认证"})
			c.Abort()
			return
		}
		var userID uint
		for _, ch := range userIDStr {
			userID = userID*10 + uint(ch-'0')
		}
		c.Set("userID", userID)
		c.Next()
	}

	group := router.Group("/api/v1/tasks", authMiddleware)
	{
		group.GET("", controller.List)
		group.POST("/training", controller.CreateTraining)
		group.GET("/:id", controller.Detail)
		group.POST("/:id/stop", controller.Stop)
		group.POST("/:id/cancel", controller.Cancel)
		group.POST("/:id/retry", controller.Retry)
		group.GET("/:id/logs", controller.Logs)
		group.GET("/:id/result", controller.Result)
	}

	return &taskTestEnv{db: db, router: router}
}

// ==================== 创建训练任务测试 ====================

func TestCreateTraining_Success(t *testing.T) {
	env := setupTaskTestEnv(t)

	reqBody := map[string]interface{}{
		"id":         "task-001",
		"name":       "训练任务1",
		"command":    "python train.py",
		"machine_id": "host-001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/training", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库中任务已创建，且类型为 training、状态为 queued
	var task entity.Task
	err = env.db.First(&task, "id = ?", "task-001").Error
	require.NoError(t, err)
	assert.Equal(t, "training", task.Type)
	assert.Equal(t, "queued", task.Status)
	assert.Equal(t, uint(1), task.CustomerID)
}

func TestCreateTraining_NoAuth(t *testing.T) {
	env := setupTaskTestEnv(t)

	reqBody := map[string]interface{}{
		"id":      "task-002",
		"name":    "无认证任务",
		"command": "echo hello",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/training", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// 不设置 X-User-ID
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.Code)
}

func TestCreateTraining_InvalidJSON(t *testing.T) {
	env := setupTaskTestEnv(t)

	// 发送无效 JSON，应返回 400
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/training", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
}

// ==================== 任务列表测试 ====================

func TestListTasks_Empty(t *testing.T) {
	env := setupTaskTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, float64(0), data["total"])
}

func TestListTasks_WithData(t *testing.T) {
	env := setupTaskTestEnv(t)

	// 创建属于用户1的任务
	env.db.Create(&entity.Task{ID: "t-1", CustomerID: 1, Name: "任务1", Command: "echo 1", Status: "queued"})
	env.db.Create(&entity.Task{ID: "t-2", CustomerID: 1, Name: "任务2", Command: "echo 2", Status: "running"})
	// 创建属于用户2的任务（用户1不应看到）
	env.db.Create(&entity.Task{ID: "t-3", CustomerID: 2, Name: "任务3", Command: "echo 3", Status: "queued"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	// 用户1只能看到自己的2个任务
	assert.Equal(t, float64(2), data["total"])
}

func TestListTasks_Pagination(t *testing.T) {
	env := setupTaskTestEnv(t)

	for i := 0; i < 5; i++ {
		env.db.Create(&entity.Task{
			ID: fmt.Sprintf("t-%d", i), CustomerID: 1,
			Name: fmt.Sprintf("任务%d", i), Command: "echo", Status: "queued",
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?page=1&page_size=2", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, float64(5), data["total"])
	list := data["list"].([]any)
	assert.Len(t, list, 2)
}

// ==================== 任务详情测试 ====================

func TestDetailTask_Success(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-detail", CustomerID: 1, Name: "详情任务",
		Command: "echo detail", Status: "running",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/t-detail", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestDetailTask_NotFound(t *testing.T) {
	env := setupTaskTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/nonexistent", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
}

func TestDetailTask_Forbidden(t *testing.T) {
	// 用户2尝试查看用户1的任务，应返回 403
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-forbidden", CustomerID: 1, Name: "他人任务",
		Command: "echo forbidden", Status: "queued",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/t-forbidden", nil)
	req.Header.Set("X-User-ID", "2")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.Code)
}

// ==================== 取消任务测试 ====================

func TestCancelTask_Success(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-cancel", CustomerID: 1, Name: "待取消任务",
		Command: "echo cancel", Status: "pending",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/t-cancel/cancel", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证数据库中状态已更新为 cancelled
	var task entity.Task
	env.db.First(&task, "id = ?", "t-cancel")
	assert.Equal(t, "cancelled", task.Status)
}

func TestCancelTask_Forbidden(t *testing.T) {
	// 用户2尝试取消用户1的任务
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-cancel-f", CustomerID: 1, Name: "他人任务",
		Command: "echo", Status: "pending",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/t-cancel-f/cancel", nil)
	req.Header.Set("X-User-ID", "2")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.Code)
}

func TestCancelTask_InvalidStatus(t *testing.T) {
	// 已完成的任务不能取消
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-cancel-done", CustomerID: 1, Name: "已完成任务",
		Command: "echo", Status: "completed",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/t-cancel-done/cancel", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
}

// ==================== 重试任务测试 ====================

func TestRetryTask_Success(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-retry", CustomerID: 1, Name: "失败任务",
		Command: "echo retry", Status: "failed", ExitCode: 1,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/t-retry/retry", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证状态已重置为 queued
	var task entity.Task
	env.db.First(&task, "id = ?", "t-retry")
	assert.Equal(t, "queued", task.Status)
}

func TestRetryTask_Forbidden(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-retry-f", CustomerID: 1, Name: "他人失败任务",
		Command: "echo", Status: "failed",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks/t-retry-f/retry", nil)
	req.Header.Set("X-User-ID", "2")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.Code)
}

// ==================== 任务日志测试 ====================

func TestLogs_Success(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-logs", CustomerID: 1, Name: "日志任务",
		Command: "echo logs", Status: "running", ErrorMsg: "some output",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/t-logs/logs", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, "t-logs", data["task_id"])
	assert.Equal(t, "running", data["status"])
}

func TestLogs_Forbidden(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-logs-f", CustomerID: 1, Name: "他人日志",
		Command: "echo", Status: "running",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/t-logs-f/logs", nil)
	req.Header.Set("X-User-ID", "2")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.Code)
}

// ==================== 任务结果测试 ====================

func TestResult_Success(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-result", CustomerID: 1, Name: "已完成任务",
		Command: "echo done", Status: "completed", ExitCode: 0,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/t-result/result", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	require.NoError(t, err)
	assert.Equal(t, "t-result", data["task_id"])
	assert.Equal(t, float64(0), data["exit_code"])
}

func TestResult_Forbidden(t *testing.T) {
	env := setupTaskTestEnv(t)

	env.db.Create(&entity.Task{
		ID: "t-result-f", CustomerID: 1, Name: "他人结果",
		Command: "echo", Status: "completed",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/t-result-f/result", nil)
	req.Header.Set("X-User-ID", "2")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.Code)
}

func TestResult_NotFound(t *testing.T) {
	env := setupTaskTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/nonexistent/result", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
}
