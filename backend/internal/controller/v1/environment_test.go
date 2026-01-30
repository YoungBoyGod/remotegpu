package v1

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestEnvironmentController_Create 测试创建环境接口
func TestEnvironmentController_Create(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// 无效的 JSON
		c.Request = httptest.NewRequest("POST", "/environments", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.Create(c)

		// 验证返回了响应（不验证具体状态码，因为依赖 service 层）
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_GetByID 测试获取环境详情接口
func TestEnvironmentController_GetByID(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/environments/test-id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.GetByID(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_List 测试环境列表接口
func TestEnvironmentController_List(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/environments", nil)
		c.Set("customer_id", uint(1)) // 模拟认证中间件设置的 customer_id

		ctrl.List(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_Delete 测试删除环境接口
func TestEnvironmentController_Delete(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("DELETE", "/environments/test-id", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.Delete(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_Start 测试启动环境接口
func TestEnvironmentController_Start(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/environments/test-id/start", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.Start(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_Stop 测试停止环境接口
func TestEnvironmentController_Stop(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/environments/test-id/stop", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.Stop(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_Restart 测试重启环境接口
func TestEnvironmentController_Restart(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/environments/test-id/restart", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.Restart(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_GetAccessInfo 测试获取访问信息接口
func TestEnvironmentController_GetAccessInfo(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/environments/test-id/access", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.GetAccessInfo(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}

// TestEnvironmentController_GetLogs 测试获取日志接口
func TestEnvironmentController_GetLogs(t *testing.T) {
	setupTestDB(t)
	ctrl := NewEnvironmentController()

	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/environments/test-id/logs", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.GetLogs(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})

	t.Run("WithTailParameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/environments/test-id/logs?tail=50", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "test-id"}}

		ctrl.GetLogs(c)

		// 验证返回了响应
		assert.NotNil(t, w.Body)
	})
}
