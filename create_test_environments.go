package main

import (
	"fmt"
	"log"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 辅助函数：将 int 转换为指针
func intPtr(i int) *int {
	return &i
}

// 辅助函数：将 int64 转换为指针
func int64Ptr(i int64) *int64 {
	return &i
}

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.GlobalConfig.Database.Host,
		config.GlobalConfig.Database.Port,
		config.GlobalConfig.Database.User,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 获取普通用户ID (假设使用user账号,ID为1602)
	var user entity.User
	if err := db.Where("username = ?", "user").First(&user).Error; err != nil {
		log.Fatalf("查找用户失败: %v", err)
	}

	fmt.Printf("找到用户: %s (ID: %d)\n", user.Username, user.ID)

	// 创建测试环境数据
	now := time.Now()
	environments := []entity.Environment{
		// Linux SSH 环境 - 运行中
		{
			UserID:  user.ID,
			HostID:      "host-001",
			Name:        "Ubuntu 开发环境",
			Description: "Ubuntu 22.04 LTS 开发环境",
			Image:       "ubuntu:22.04",
			Status:      "running",
			CPU:         4,
			Memory:      8192, // 8GB
			GPU:         1,
			Storage:     int64Ptr(100),
			SSHPort:     intPtr(22001),
			CreatedAt:   now.Add(-48 * time.Hour),
			UpdatedAt:   now,
			StartedAt:   &now,
		},
		// Windows RDP 环境 - 运行中
		{
			UserID:  user.ID,
			HostID:      "host-002",
			Name:        "Windows Server 2022",
			Description: "Windows Server 2022 远程桌面环境",
			Image:       "windows-server:2022",
			Status:      "running",
			CPU:         8,
			Memory:      16384, // 16GB
			GPU:         2,
			Storage:     int64Ptr(200),
			RDPPort:     intPtr(3389),
			CreatedAt:   now.Add(-24 * time.Hour),
			UpdatedAt:   now,
			StartedAt:   &now,
		},
		// Jupyter 数据科学环境 - 运行中
		{
			UserID:  user.ID,
			HostID:      "host-003",
			Name:        "PyTorch 训练环境",
			Description: "PyTorch 深度学习训练环境",
			Image:       "pytorch/pytorch:2.0.0-cuda11.7-cudnn8-runtime",
			Status:      "running",
			CPU:         16,
			Memory:      32768, // 32GB
			GPU:         4,
			Storage:     int64Ptr(500),
			SSHPort:     intPtr(22002),
			JupyterPort: intPtr(8888),
			CreatedAt:   now.Add(-72 * time.Hour),
			UpdatedAt:   now,
			StartedAt:   &now,
		},
		// TensorFlow 环境 - 已停止
		{
			UserID:  user.ID,
			HostID:      "host-004",
			Name:        "TensorFlow 开发环境",
			Description: "TensorFlow 2.x 开发环境",
			Image:       "tensorflow/tensorflow:latest-gpu",
			Status:      "stopped",
			CPU:         8,
			Memory:      16384, // 16GB
			GPU:         2,
			Storage:     int64Ptr(300),
			SSHPort:     intPtr(22003),
			JupyterPort: intPtr(8889),
			CreatedAt:   now.Add(-96 * time.Hour),
			UpdatedAt:   now,
		},
		// CentOS 环境 - 已停止
		{
			UserID:  user.ID,
			HostID:      "host-005",
			Name:        "CentOS 测试环境",
			Description: "CentOS 7 测试环境",
			Image:       "centos:7",
			Status:      "stopped",
			CPU:         2,
			Memory:      4096, // 4GB
			GPU:         0,
			Storage:     int64Ptr(50),
			SSHPort:     intPtr(22004),
			CreatedAt:   now.Add(-120 * time.Hour),
			UpdatedAt:   now,
		},
		// Windows 10 桌面环境 - 运行中
		{
			UserID:  user.ID,
			HostID:      "host-006",
			Name:        "Windows 10 工作站",
			Description: "Windows 10 Pro 图形工作站",
			Image:       "windows-10:pro",
			Status:      "running",
			CPU:         6,
			Memory:      12288, // 12GB
			GPU:         1,
			Storage:     int64Ptr(150),
			RDPPort:     intPtr(3390),
			CreatedAt:   now.Add(-36 * time.Hour),
			UpdatedAt:   now,
			StartedAt:   &now,
		},
	}

	// 插入环境数据
	fmt.Println("\n开始创建测试环境...")
	for i, env := range environments {
		if err := db.Create(&env).Error; err != nil {
			log.Printf("创建环境 %s 失败: %v", env.Name, err)
			continue
		}
		fmt.Printf("%d. 创建环境: %s (状态: %s)\n", i+1, env.Name, env.Status)
	}

	fmt.Println("\n✓ 测试环境创建完成!")
	fmt.Printf("共创建 %d 个测试环境\n", len(environments))
}
