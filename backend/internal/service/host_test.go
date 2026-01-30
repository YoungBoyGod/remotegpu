package service

import (
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
)

func setupTestDB(t *testing.T) {
	cfg := database.Config{
		Host:     "192.168.10.210",
		Port:     5432,
		User:     "remotegpu_user",
		Password: "remotegpu_password",
		DBName:   "remotegpu",
	}
	if err := database.InitDB(cfg); err != nil {
		t.Skipf("跳过测试，无法连接数据库: %v", err)
	}
}

func TestHostService_CRUD(t *testing.T) {
	setupTestDB(t)

	svc := NewHostService()

	// 创建主机
	host := &entity.Host{
		Name:           "Test Service Host",
		IPAddress:      "192.168.1.200",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       16,
		TotalMemory:    34359738368,
	}

	if err := svc.Create(host); err != nil {
		t.Fatalf("创建主机失败: %v", err)
	}
	t.Logf("创建主机成功, ID: %s", host.ID)

	// 获取主机
	found, err := svc.GetByID(host.ID)
	if err != nil {
		t.Fatalf("获取主机失败: %v", err)
	}
	if found.Name != "Test Service Host" {
		t.Fatalf("主机名称不匹配")
	}
	t.Log("获取主机成功")

	// 删除主机
	if err := svc.Delete(host.ID); err != nil {
		t.Fatalf("删除主机失败: %v", err)
	}
	t.Log("删除主机成功")

	t.Log("Host Service CRUD 测试通过")
}
