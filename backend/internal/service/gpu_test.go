package service

import (
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

func TestGPUService_CRUD(t *testing.T) {
	setupTestDB(t)

	hostSvc := NewHostService()
	gpuSvc := NewGPUService()

	// 先创建测试主机
	host := &entity.Host{
		Name:           "Test Host for GPU Service",
		IPAddress:      "192.168.1.201",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}
	if err := hostSvc.Create(host); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostSvc.Delete(host.ID)

	// 创建GPU
	gpu := &entity.GPU{
		HostID:      host.ID,
		GPUIndex:    0,
		Name:        "RTX 4090",
		Brand:       "NVIDIA",
		MemoryTotal: 25769803776,
	}
	if err := gpuSvc.Create(gpu); err != nil {
		t.Fatalf("创建GPU失败: %v", err)
	}
	t.Logf("创建GPU成功, ID: %d", gpu.ID)

	// 获取GPU
	found, err := gpuSvc.GetByID(gpu.ID)
	if err != nil {
		t.Fatalf("获取GPU失败: %v", err)
	}
	if found.Name != "RTX 4090" {
		t.Fatalf("GPU名称不匹配")
	}
	t.Log("获取GPU成功")

	// 删除GPU
	if err := gpuSvc.Delete(gpu.ID); err != nil {
		t.Fatalf("删除GPU失败: %v", err)
	}
	t.Log("删除GPU成功")

	t.Log("GPU Service CRUD 测试通过")
}
