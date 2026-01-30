package dao

import (
	"testing"
	"time"

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

func TestHostDao_CRUD(t *testing.T) {
	setupTestDB(t)

	dao := NewHostDao()
	testID := "test-host-" + time.Now().Format("20060102150405")

	// 创建测试主机
	host := &entity.Host{
		ID:             testID,
		Name:           "Test Host",
		IPAddress:      "192.168.1.100",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}

	// 测试 Create
	err := dao.Create(host)
	if err != nil {
		t.Fatalf("创建主机失败: %v", err)
	}
	t.Log("创建主机成功")

	// 测试 GetByID
	found, err := dao.GetByID(testID)
	if err != nil {
		t.Fatalf("获取主机失败: %v", err)
	}
	if found.Name != "Test Host" {
		t.Fatalf("主机名称不匹配: got %s, want Test Host", found.Name)
	}
	t.Log("获取主机成功")

	// 测试 Update
	found.Name = "Updated Host"
	err = dao.Update(found)
	if err != nil {
		t.Fatalf("更新主机失败: %v", err)
	}
	t.Log("更新主机成功")

	// 测试 Delete
	err = dao.Delete(testID)
	if err != nil {
		t.Fatalf("删除主机失败: %v", err)
	}
	t.Log("删除主机成功")

	t.Log("Host DAO CRUD 测试通过")
}
