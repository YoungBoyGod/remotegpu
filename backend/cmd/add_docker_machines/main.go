package main

import (
	"context"
	"fmt"
	"log"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
)

func main() {
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

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

	testMachines := []entity.Host{
		{
			Name:        "test-machine-02",
			Hostname:    "test-machine-02",
			Region:      "default",
			IPAddress:   "127.0.0.1",
			SSHPort:     2202,
			SSHUsername: "root",
			SSHPassword: "root",
			Status:      "offline",
		},
		{
			Name:        "test-machine-03",
			Hostname:    "test-machine-03",
			Region:      "default",
			IPAddress:   "127.0.0.2",
			SSHPort:     2203,
			SSHUsername: "root",
			SSHPassword: "root",
			Status:      "offline",
		},
	}

	fmt.Println("添加测试机器...")
	for _, host := range testMachines {
		if err := machineSvc.CreateMachine(ctx, &host); err != nil {
			fmt.Printf("❌ %s: %v\n", host.Name, err)
			continue
		}
		fmt.Printf("✅ %s 已添加\n", host.Name)
	}
}
