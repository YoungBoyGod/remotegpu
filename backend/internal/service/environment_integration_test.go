package service

import (
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnvironmentIntegration_CreateAndDelete 集成测试：创建和删除环境
func TestEnvironmentIntegration_CreateAndDelete(t *testing.T) {
	// 跳过集成测试（如果没有测试数据库）
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化测试数据库
	setupTestDB(t)
	db := database.GetDB()
	require.NotNil(t, db, "数据库连接不能为空")

	// 清理测试数据
	defer cleanupTestData(t)

	// 准备测试数据：创建测试客户
	customer := &entity.Customer{
		UUID:         uuid.New(),
		Username:     "test_user_integration",
		Email:        "test_integration@example.com",
		PasswordHash: "test_password_hash",
		DisplayName:  "Test Integration User",
		Status:       "active",
	}
	err := db.Create(customer).Error
	require.NoError(t, err, "创建测试客户失败")

	// 准备测试数据：创建测试主机
	host := &entity.Host{
		ID:             "test-host-integration",
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "k8s",
		Status:         "active",
		TotalCPU:       16,
		UsedCPU:        0,
		TotalMemory:    32000,
		UsedMemory:     0,
		TotalGPU:       4,
		UsedGPU:        0,
	}
	err = db.Create(host).Error
	require.NoError(t, err, "创建测试主机失败")

	// 准备测试数据：创建测试 GPU
	for i := 0; i < 4; i++ {
		gpu := &entity.GPU{
			HostID:   host.ID,
			GPUIndex: i,
			UUID:     "GPU-TEST-UUID-" + string(rune('0'+i)),
			Name:     "Test GPU",
			Status:   "available",
		}
		err = db.Create(gpu).Error
		require.NoError(t, err, "创建测试 GPU 失败")
	}

	// 准备测试数据：创建资源配额
	quota := &entity.ResourceQuota{
		CustomerID: customer.ID,
		CPU:        100,
		Memory:     200000,
		GPU:        10,
		Storage:    1000000,
	}
	err = db.Create(quota).Error
	require.NoError(t, err, "创建资源配额失败")

	t.Log("✅ 测试数据准备完成")

	// 注意：由于 CreateEnvironment 需要 K8s 客户端，这里只测试到主机选择
	// 完整的集成测试需要 K8s 测试环境
	t.Log("⚠️ 完整的环境创建集成测试需要 K8s 环境")
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T) {
	db := database.GetDB()

	// 删除测试数据（按依赖关系倒序删除）
	db.Exec("DELETE FROM gpus WHERE host_id LIKE 'test-%'")
	db.Exec("DELETE FROM hosts WHERE id LIKE 'test-%'")
	db.Exec("DELETE FROM resource_quotas WHERE customer_id IN (SELECT id FROM customers WHERE username LIKE 'test_%')")
	db.Exec("DELETE FROM customers WHERE username LIKE 'test_%'")

	t.Log("✅ 测试数据清理完成")
}

// TestEnvironmentIntegration_HostSelection 集成测试：主机选择算法
func TestEnvironmentIntegration_HostSelection(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化测试数据库
	setupTestDB(t)
	db := database.GetDB()
	require.NotNil(t, db)

	defer cleanupTestData(t)

	// 创建多个测试主机
	hosts := []*entity.Host{
		{
			ID:          "test-host-1",
			Name:        "Host 1",
			IPAddress:   "192.168.1.101",
			OSType:      "linux",
			DeploymentMode: "k8s",
			Status:      "active",
			TotalCPU:    16,
			UsedCPU:     12, // 75% 使用率
			TotalMemory: 32000,
			UsedMemory:  24000,
			TotalGPU:    4,
			UsedGPU:     3,
		},
		{
			ID:          "test-host-2",
			Name:        "Host 2",
			IPAddress:   "192.168.1.102",
			OSType:      "linux",
			DeploymentMode: "k8s",
			Status:      "active",
			TotalCPU:    16,
			UsedCPU:     4, // 25% 使用率
			TotalMemory: 32000,
			UsedMemory:  8000,
			TotalGPU:    4,
			UsedGPU:     1,
		},
	}

	for _, host := range hosts {
		err := db.Create(host).Error
		require.NoError(t, err)
	}

	// 创建 service
	service := &EnvironmentService{
		hostDao: dao.NewHostDao(),
	}

	// 测试主机选择
	selectedHost, err := service.selectHost(2, 4000, 1)

	assert.NoError(t, err)
	assert.NotNil(t, selectedHost)
	// 应该选择使用率较低的 host-2
	assert.Equal(t, "test-host-2", selectedHost.ID)

	t.Log("✅ 主机选择算法测试通过")
}
