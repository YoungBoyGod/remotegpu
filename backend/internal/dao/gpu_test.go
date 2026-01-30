package dao

import (
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

func TestGPUDao_CRUD(t *testing.T) {
	setupTestDB(t)

	hostDao := NewHostDao()
	gpuDao := NewGPUDao()

	// 先创建测试主机
	testHostID := "test-host-gpu-" + time.Now().Format("20060102150405")
	host := &entity.Host{
		ID:             testHostID,
		Name:           "Test Host for GPU",
		IPAddress:      "192.168.1.101",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}
	if err := hostDao.Create(host); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostDao.Delete(testHostID)

	// 测试 Create GPU
	gpu := &entity.GPU{
		HostID:      testHostID,
		GPUIndex:    0,
		Name:        "Tesla V100",
		Brand:       "NVIDIA",
		MemoryTotal: 34359738368,
	}
	if err := gpuDao.Create(gpu); err != nil {
		t.Fatalf("创建GPU失败: %v", err)
	}
	t.Logf("创建GPU成功, ID: %d", gpu.ID)

	// 测试 GetByID
	found, err := gpuDao.GetByID(gpu.ID)
	if err != nil {
		t.Fatalf("获取GPU失败: %v", err)
	}
	if found.Name != "Tesla V100" {
		t.Fatalf("GPU名称不匹配")
	}
	t.Log("获取GPU成功")

	// 测试 GetByHostID
	gpus, err := gpuDao.GetByHostID(testHostID)
	if err != nil {
		t.Fatalf("获取主机GPU列表失败: %v", err)
	}
	if len(gpus) != 1 {
		t.Fatalf("GPU数量不匹配: got %d, want 1", len(gpus))
	}
	t.Log("获取主机GPU列表成功")

	// 测试 Delete
	if err := gpuDao.Delete(gpu.ID); err != nil {
		t.Fatalf("删除GPU失败: %v", err)
	}
	t.Log("删除GPU成功")

	t.Log("GPU DAO CRUD 测试通过")
}

func TestGPUDao_ListAndAllocate(t *testing.T) {
	setupTestDB(t)

	hostDao := NewHostDao()
	gpuDao := NewGPUDao()

	// 创建测试主机
	testHostID := "test-host-gpu-list-" + time.Now().Format("20060102150405")
	host := &entity.Host{
		ID:             testHostID,
		Name:           "Test Host for GPU List",
		IPAddress:      "192.168.1.102",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}
	if err := hostDao.Create(host); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostDao.Delete(testHostID)

	// 创建测试GPU
	gpu := &entity.GPU{
		HostID:      testHostID,
		GPUIndex:    0,
		Name:        "RTX 4090",
		Brand:       "NVIDIA",
		MemoryTotal: 25769803776,
		Status:      "available",
	}
	if err := gpuDao.Create(gpu); err != nil {
		t.Fatalf("创建GPU失败: %v", err)
	}
	defer gpuDao.Delete(gpu.ID)

	// 测试 List
	gpus, total, err := gpuDao.List(1, 10)
	if err != nil {
		t.Fatalf("获取GPU列表失败: %v", err)
	}
	if total < 1 {
		t.Fatalf("GPU总数不正确: %d", total)
	}
	t.Logf("获取GPU列表成功, 总数: %d, 当前页: %d", total, len(gpus))

	// 测试 Allocate
	if err := gpuDao.Allocate(gpu.ID, "env-test-001"); err != nil {
		t.Fatalf("分配GPU失败: %v", err)
	}
	allocated, _ := gpuDao.GetByID(gpu.ID)
	if allocated.Status != "allocated" || allocated.AllocatedTo != "env-test-001" {
		t.Fatalf("GPU分配状态不正确")
	}
	t.Log("分配GPU成功")

	// 测试 Release
	if err := gpuDao.Release(gpu.ID); err != nil {
		t.Fatalf("释放GPU失败: %v", err)
	}
	released, _ := gpuDao.GetByID(gpu.ID)
	if released.Status != "available" || released.AllocatedTo != "" {
		t.Fatalf("GPU释放状态不正确")
	}
	t.Log("释放GPU成功")

	t.Log("GPU DAO List/Allocate/Release 测试通过")
}
