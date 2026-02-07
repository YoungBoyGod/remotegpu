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

// test_add_machine - 测试添加机器并验证 SSH 密码加密
// @author Claude
// @description 测试 P0 安全修复：SSH 密码加密功能
// @usage ENCRYPTION_KEY="..." go run cmd/test_add_machine/main.go
// @modified 2026-02-06

func main() {
	// 检查环境变量
	encKey := os.Getenv("ENCRYPTION_KEY")
	if encKey == "" {
		log.Fatal("错误：必须设置 ENCRYPTION_KEY 环境变量（32 字节）")
	}
	if len(encKey) != 32 {
		log.Fatal("错误：ENCRYPTION_KEY 必须是 32 字节")
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

	db := database.DB
	ctx := context.Background()

	// 创建 MachineService
	machineSvc := machine.NewMachineService(db)

	// 准备机器信息
	host := &entity.Host{
		Name:        "test-202",
		Hostname:    "192.168.10.202",
		Region:      "default",
		IPAddress:   "192.168.10.202",
		SSHPort:     22,
		SSHUsername: "luo",
		SSHPassword: "luo", // 明文密码，将被自动加密
		Status:      "offline",
	}

	fmt.Println("\n=== 测试添加机器（SSH 密码加密）===")
	fmt.Printf("机器信息:\n")
	fmt.Printf("  IP: %s\n", host.IPAddress)
	fmt.Printf("  用户名: %s\n", host.SSHUsername)
	fmt.Printf("  密码: %s (明文，将被加密)\n", host.SSHPassword)
	fmt.Println()

	// 添加机器（密码会自动加密）
	if err := machineSvc.CreateMachine(ctx, host); err != nil {
		log.Fatalf("添加机器失败: %v", err)
	}

	fmt.Println("✅ 机器添加成功！")
	fmt.Printf("机器 ID: %s\n", host.ID)
	fmt.Println()

	// 验证密码是否已加密
	savedHost, err := machineSvc.GetHost(ctx, host.ID)
	if err != nil {
		log.Fatalf("查询机器失败: %v", err)
	}

	fmt.Println("=== 验证密码加密 ===")
	fmt.Printf("数据库中的密码长度: %d 字节\n", len(savedHost.SSHPassword))
	fmt.Printf("密码前缀: %s...\n", savedHost.SSHPassword[:min(30, len(savedHost.SSHPassword))])

	if len(savedHost.SSHPassword) > 40 {
		fmt.Println("✅ 密码已成功加密（长度 > 40 字节）")
	} else {
		fmt.Println("❌ 警告：密码可能未加密")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
