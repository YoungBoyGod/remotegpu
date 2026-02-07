package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
)

// add_test_machines - 批量添加测试机器
// @author Claude
// @description 添加本地 Docker 测试机器并验证加密
// @usage ENCRYPTION_KEY="..." go run cmd/add_test_machines/main.go
// @modified 2026-02-06

func main() {
	// 检查环境变量
	encKey := os.Getenv("ENCRYPTION_KEY")
	if encKey == "" {
		log.Fatal("错误：必须设置 ENCRYPTION_KEY 环境变量（32 字节）")
	}

	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	dbConfig := database.Config{
		Host:     config.GlobalConfig.Database.Host,
		Port:     config.GlobalConfig.Database.Port,
		User:     config.GlobalConfig.Database.User,
		Password: config.GlobalConfig.Database.Password,
		DBName:   config.GlobalConfig.Database.DBName,
	}
	if err := database.InitDB(dbConfig); err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	ctx := context.Background()
	machineSvc := machine.NewMachineService(database.DB)

	// 定义测试机器
	testMachines := []entity.Host{
		{
			Name:        "test-machine-01",
			Hostname:    "test-machine-01",
			Region:      "default",
			IPAddress:   "localhost",
			SSHPort:     2201,
			SSHUsername: "root",
			SSHPassword: "root",
			Status:      "offline",
		},
		{
			Name:        "test-machine-02",
			Hostname:    "test-machine-02",
			Region:      "default",
			IPAddress:   "localhost",
			SSHPort:     2202,
			SSHUsername: "root",
			SSHPassword: "root",
			Status:      "offline",
		},
		{
			Name:        "test-machine-03",
			Hostname:    "test-machine-03",
			Region:      "default",
			IPAddress:   "localhost",
			SSHPort:     2203,
			SSHUsername: "root",
			SSHPassword: "root",
			Status:      "offline",
		},
	}

	fmt.Println("\n=== 批量添加测试机器 ===")
	fmt.Printf("将添加 %d 台测试机器\n\n", len(testMachines))

	successCount := 0
	failCount := 0

	for i, host := range testMachines {
		fmt.Printf("[%d/%d] 添加 %s (localhost:%d) ... ", i+1, len(testMachines), host.Name, host.SSHPort)

		if err := machineSvc.CreateMachine(ctx, &host); err != nil {
			fmt.Printf("失败: %v\n", err)
			failCount++
			continue
		}

		fmt.Println("✅ 成功")
		successCount++
	}

	fmt.Println("\n=== 添加完成 ===")
	fmt.Printf("成功: %d\n", successCount)
	fmt.Printf("失败: %d\n", failCount)
}
