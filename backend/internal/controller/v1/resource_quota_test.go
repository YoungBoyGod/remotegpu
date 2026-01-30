package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestResourceQuotaController_SetQuota 测试设置配额
func TestResourceQuotaController_SetQuota(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	quotaCtrl := NewResourceQuotaController()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-quota-ctrl-" + uuid.New().String()[:8],
		Email:        "quota-ctrl-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Quota Controller User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	r.POST("/quotas", quotaCtrl.SetQuota)

	// 测试设置用户级别配额
	req := SetQuotaRequest{
		CustomerID: customer.ID,
		CPU:        16,
		Memory:     32768,
		GPU:        4,
		Storage:    1000,
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/quotas", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code, "设置配额应该成功")
	t.Log("✅ 设置配额成功")

	// 解析响应
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 验证响应数据
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(16), data["cpu"])
	assert.Equal(t, float64(32768), data["memory"])
	assert.Equal(t, float64(4), data["gpu"])
	assert.Equal(t, float64(1000), data["storage"])
}

// TestResourceQuotaController_GetQuota 测试获取配额
func TestResourceQuotaController_GetQuota(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	quotaCtrl := NewResourceQuotaController()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-get-quota-" + uuid.New().String()[:8],
		Email:        "get-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Get Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 先设置配额
	quotaService := quotaCtrl.quotaService
	quota := &entity.ResourceQuota{
		CustomerID: customer.ID,
		CPU:        8,
		Memory:     16384,
		GPU:        2,
		Storage:    500,
	}
	err = quotaService.SetQuota(quota)
	assert.NoError(t, err)

	r.GET("/quotas/:id", quotaCtrl.GetQuota)

	// 测试获取配额
	httpReq := httptest.NewRequest("GET", "/quotas/"+fmt.Sprint(quota.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code, "获取配额应该成功")
	t.Log("✅ 获取配额成功")

	// 解析响应
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 验证响应数据
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(8), data["cpu"])
	assert.Equal(t, float64(16384), data["memory"])
}

// TestResourceQuotaController_UpdateQuota 测试更新配额
func TestResourceQuotaController_UpdateQuota(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	quotaCtrl := NewResourceQuotaController()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-update-quota-" + uuid.New().String()[:8],
		Email:        "update-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Update Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 先设置配额
	quotaService := quotaCtrl.quotaService
	quota := &entity.ResourceQuota{
		CustomerID: customer.ID,
		CPU:        8,
		Memory:     16384,
		GPU:        2,
		Storage:    500,
	}
	err = quotaService.SetQuota(quota)
	assert.NoError(t, err)

	r.PUT("/quotas/:id", quotaCtrl.UpdateQuota)

	// 测试更新配额
	updateReq := UpdateQuotaRequest{
		CPU:     16,
		Memory:  32768,
		GPU:     4,
		Storage: 1000,
	}

	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest("PUT", "/quotas/"+fmt.Sprint(quota.ID), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code, "更新配额应该成功")
	t.Log("✅ 更新配额成功")
}

// TestResourceQuotaController_DeleteQuota 测试删除配额
func TestResourceQuotaController_DeleteQuota(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	quotaCtrl := NewResourceQuotaController()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-delete-quota-" + uuid.New().String()[:8],
		Email:        "delete-quota-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Delete Quota User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 先设置配额
	quotaService := quotaCtrl.quotaService
	quota := &entity.ResourceQuota{
		CustomerID: customer.ID,
		CPU:        8,
		Memory:     16384,
		GPU:        2,
		Storage:    500,
	}
	err = quotaService.SetQuota(quota)
	assert.NoError(t, err)

	r.DELETE("/quotas/:id", quotaCtrl.DeleteQuota)

	// 测试删除配额
	httpReq := httptest.NewRequest("DELETE", "/quotas/"+fmt.Sprint(quota.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code, "删除配额应该成功")
	t.Log("✅ 删除配额成功")

	// 验证配额已被删除
	_, err = quotaService.GetQuotaByID(quota.ID)
	assert.Error(t, err, "配额应该已被删除")
	t.Log("✅ 验证配额已被删除")
}

// TestResourceQuotaController_GetUsage 测试获取资源使用情况
func TestResourceQuotaController_GetUsage(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	quotaCtrl := NewResourceQuotaController()
	customerDao := dao.NewCustomerDao()
	workspaceDao := dao.NewWorkspaceDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-usage-" + uuid.New().String()[:8],
		Email:        "usage-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Usage User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 创建测试工作空间
	workspace := &entity.Workspace{
		Name:        "test-workspace-" + uuid.New().String()[:8],
		Description: "Test Workspace",
		OwnerID:     customer.ID,
		Status:      "active",
	}
	err = workspaceDao.Create(workspace)
	assert.NoError(t, err)
	defer workspaceDao.Delete(workspace.ID)

	// 设置配额
	quotaService := quotaCtrl.quotaService
	quota := &entity.ResourceQuota{
		CustomerID: customer.ID,
		CPU:        16,
		Memory:     32768,
		GPU:        4,
		Storage:    1000,
	}
	err = quotaService.SetQuota(quota)
	assert.NoError(t, err)
	defer quotaService.DeleteQuota(quota.ID)

	r.GET("/quotas/usage", quotaCtrl.GetUsage)

	// 测试获取资源使用情况（无环境时）
	httpReq := httptest.NewRequest("GET", "/quotas/usage?customer_id="+fmt.Sprint(customer.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code, "获取资源使用情况应该成功")
	t.Log("✅ 获取资源使用情况成功")

	// 解析响应
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 验证响应数据结构
	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["quota"], "应该包含quota字段")
	assert.NotNil(t, data["used"], "应该包含used字段")
	assert.NotNil(t, data["available"], "应该包含available字段")

	// 验证配额值
	quotaData := data["quota"].(map[string]interface{})
	assert.Equal(t, float64(16), quotaData["cpu"])
	assert.Equal(t, float64(32768), quotaData["memory"])

	// 验证已使用资源（无环境时应该为0）
	usedData := data["used"].(map[string]interface{})
	assert.Equal(t, float64(0), usedData["cpu"])
	assert.Equal(t, float64(0), usedData["memory"])

	// 验证可用配额（应该等于总配额）
	availableData := data["available"].(map[string]interface{})
	assert.Equal(t, float64(16), availableData["cpu"])
	assert.Equal(t, float64(32768), availableData["memory"])

	t.Log("✅ 验证资源使用情况数据正确")
}

// TestResourceQuotaController_EnvironmentIntegration 测试ResourceQuota与Environment模块的集成
func TestResourceQuotaController_EnvironmentIntegration(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	quotaCtrl := NewResourceQuotaController()
	customerDao := dao.NewCustomerDao()
	workspaceDao := dao.NewWorkspaceDao()
	envDao := dao.NewEnvironmentDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-env-integration-" + uuid.New().String()[:8],
		Email:        "env-integration-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Environment Integration User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 创建测试工作空间
	workspace := &entity.Workspace{
		Name:        "test-workspace-" + uuid.New().String()[:8],
		Description: "Test Workspace",
		OwnerID:     customer.ID,
		Status:      "active",
	}
	err = workspaceDao.Create(workspace)
	assert.NoError(t, err)
	defer workspaceDao.Delete(workspace.ID)

	t.Log("✅ 测试环境准备完成")

	// 场景1：设置配额并验证
	t.Run("Scenario1: Set quota via API", func(t *testing.T) {
		r.POST("/quotas", quotaCtrl.SetQuota)

		req := SetQuotaRequest{
			CustomerID: customer.ID,
			CPU:        16,
			Memory:     32768,
			GPU:        4,
			Storage:    1000,
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/quotas", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code, "设置配额应该成功")
		t.Log("✅ 场景1：通过API设置配额成功")
	})

	// 场景2：创建环境并验证配额使用情况跟踪
	t.Run("Scenario2: Create environments and verify quota usage", func(t *testing.T) {
		// 创建测试主机
		db := database.GetDB()
		host := &entity.Host{
			ID:             "test-host-integration-" + uuid.New().String()[:8],
			Name:           "Test Integration Host",
			IPAddress:      "192.168.1.100",
			OSType:         "linux",
			DeploymentMode: "k8s",
			Status:         "active",
			TotalCPU:       32,
			TotalMemory:    64000,
			TotalGPU:       8,
		}
		err := db.Create(host).Error
		assert.NoError(t, err)
		defer db.Where("id = ?", host.ID).Delete(&entity.Host{})

		// 创建第一个环境（使用部分配额）
		storage1 := int64(500)
		env1 := &entity.Environment{
			ID:          "env-integration-1-" + uuid.New().String()[:8],
			Name:        "env-integration-1",
			CustomerID:  customer.ID,
			WorkspaceID: &workspace.ID,
			HostID:      host.ID,
			Status:      "running",
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     &storage1,
			Image:       "test-image:latest",
		}
		err = envDao.Create(env1)
		assert.NoError(t, err)
		defer envDao.Delete(env1.ID)

		t.Log("✅ 创建第一个环境成功")

		// 通过API获取资源使用情况
		r.GET("/quotas/usage", quotaCtrl.GetUsage)
		httpReq := httptest.NewRequest("GET", "/quotas/usage?customer_id="+fmt.Sprint(customer.ID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code, "获取资源使用情况应该成功")

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})

		// 验证已使用资源
		usedData := data["used"].(map[string]interface{})
		assert.Equal(t, float64(8), usedData["cpu"], "应该显示CPU已使用8")
		assert.Equal(t, float64(16384), usedData["memory"], "应该显示Memory已使用16384")
		assert.Equal(t, float64(2), usedData["gpu"], "应该显示GPU已使用2")

		// 验证可用配额
		availableData := data["available"].(map[string]interface{})
		assert.Equal(t, float64(8), availableData["cpu"], "剩余CPU应该是8")
		assert.Equal(t, float64(16384), availableData["memory"], "剩余Memory应该是16384")
		assert.Equal(t, float64(2), availableData["gpu"], "剩余GPU应该是2")

		t.Log("✅ 场景2：配额使用情况跟踪正确")
	})
}
