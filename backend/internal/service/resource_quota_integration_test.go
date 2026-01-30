package service

import (
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestResourceQuotaIntegration_QuotaCheckFlow 测试配额检查流程
func TestResourceQuotaIntegration_QuotaCheckFlow(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-integration-" + uuid.New().String()[:8],
		Email:        "integration-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Integration User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 场景1：设置配额 → 检查配额（足够）
	t.Run("Scenario1: Set quota and check (sufficient)", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			CustomerID:  customer.ID,
			WorkspaceID: nil,
			CPU:         16,
			Memory:      32768,
			GPU:         4,
			Storage:     1000,
		}
		err := service.SetQuota(quota)
		assert.NoError(t, err)

		// 检查配额（请求资源少于配额）
		request := &ResourceRequest{
			CPU:     8,
			Memory:  16384,
			GPU:     2,
			Storage: 500,
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.NoError(t, err)
		assert.True(t, ok)
		t.Log("场景1通过：配额足够")

		// 清理
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	})

	// 场景2：设置配额 → 检查配额（不足）→ 拒绝
	t.Run("Scenario2: Set quota and check (insufficient)", func(t *testing.T) {
		quota := &entity.ResourceQuota{
			CustomerID:  customer.ID,
			WorkspaceID: nil,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}
		err := service.SetQuota(quota)
		assert.NoError(t, err)

		// 检查配额（请求资源超过配额）
		request := &ResourceRequest{
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)
		t.Log("场景2通过：配额不足，正确拒绝")

		// 清理
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	})

	// 场景3：测试配额更新后的检查
	t.Run("Scenario3: Update quota and recheck", func(t *testing.T) {
		// 设置初始配额
		quota := &entity.ResourceQuota{
			CustomerID:  customer.ID,
			WorkspaceID: nil,
			CPU:         8,
			Memory:      16384,
			GPU:         2,
			Storage:     500,
		}
		err := service.SetQuota(quota)
		assert.NoError(t, err)

		// 请求资源（超过当前配额）
		request := &ResourceRequest{
			CPU:     16,
			Memory:  32768,
			GPU:     4,
			Storage: 1000,
		}
		ok, err := service.CheckQuota(customer.ID, nil, request)
		assert.Error(t, err)
		assert.False(t, ok)

		// 更新配额（增加配额）
		retrieved, _ := service.GetQuota(customer.ID, nil)
		retrieved.CPU = 32
		retrieved.Memory = 65536
		retrieved.GPU = 8
		retrieved.Storage = 2000
		err = service.UpdateQuota(retrieved)
		assert.NoError(t, err)

		// 再次检查配额（现在应该足够了）
		ok, err = service.CheckQuota(customer.ID, nil, request)
		assert.NoError(t, err)
		assert.True(t, ok)
		t.Log("场景3通过：更新配额后检查通过")

		// 清理
		service.DeleteQuota(retrieved.ID)
	})
}

// TestResourceQuotaIntegration_AvailableQuota 测试可用配额计算
func TestResourceQuotaIntegration_AvailableQuota(t *testing.T) {
	setupResourceQuotaServiceTest(t)

	service := NewResourceQuotaService()
	customerDao := dao.NewCustomerDao()

	// 创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test-available-" + uuid.New().String()[:8],
		Email:        "available-" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "test-hash",
		DisplayName:  "Test Available User",
		Status:       "active",
	}
	err := customerDao.Create(customer)
	assert.NoError(t, err)
	defer customerDao.Delete(customer.ID)

	// 设置配额
	quota := &entity.ResourceQuota{
		CustomerID:  customer.ID,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
	}
	err = service.SetQuota(quota)
	assert.NoError(t, err)
	defer func() {
		retrieved, _ := service.GetQuota(customer.ID, nil)
		if retrieved != nil {
			service.DeleteQuota(retrieved.ID)
		}
	}()

	// TODO: 需要等待 Environment 实体实现后才能完整测试资源统计功能
	// 当前测试场景：无环境 → 已使用资源为0 → 可用配额等于总配额
	available, err := service.GetAvailableQuota(customer.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 16, available.CPU)
	assert.Equal(t, int64(32768), available.Memory)
	assert.Equal(t, 4, available.GPU)
	assert.Equal(t, int64(1000), available.Storage)
	t.Log("可用配额测试通过（当前无环境，可用配额等于总配额）")

	// TODO: 待 Environment 实体实现后补充以下测试场景：
	// - 场景1：创建环境 → 已使用资源增加 → 可用配额减少
	// - 场景2：删除环境 → 已使用资源减少 → 可用配额增加
	// - 场景3：多个工作空间 → 分别统计各工作空间的资源使用
}
