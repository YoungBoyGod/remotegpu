package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/lib/pq"
)

func main() {
	// 初始化数据库连接（使用配置文件中的数据库配置）
	dbConfig := database.Config{
		Host:     "192.168.10.210",
		Port:     5432,
		User:     "remotegpu_user",
		Password: "remotegpu_password",
		DBName:   "remotegpu",
	}

	if err := database.InitDB(dbConfig); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	db := database.GetDB()

	// 生成100台设备数据
	hosts := make([]*entity.Host, 0, 100)

	for i := 1; i <= 100; i++ {
		// 生成设备ID
		hostID := fmt.Sprintf("host-%03d", i)

		// 生成IP地址
		ipAddress := fmt.Sprintf("192.168.1.%d", i)
		publicIP := fmt.Sprintf("10.0.1.%d", i)

		// 随机选择操作系统类型
		osTypes := []string{"linux", "windows"}
		osType := osTypes[i%2]

		var osVersion string
		if osType == "linux" {
			osVersion = "Ubuntu 22.04"
		} else {
			osVersion = "Windows Server 2022"
		}

		// 随机选择部署模式
		deploymentModes := []string{"docker", "k8s_pod", "vm"}
		deploymentMode := deploymentModes[i%3]

		// 生成K8s节点名称（如果是k8s部署模式）
		var k8sNodeName string
		if deploymentMode == "k8s_pod" {
			k8sNodeName = fmt.Sprintf("k8s-node-%03d", i)
		}

		// 随机分配资源
		totalCPU := 16 + (i%16)*4  // 16-76核
		totalMemory := int64(32768 + (i%32)*2048)  // 32GB-96GB
		totalDisk := int64(500 + (i%50)*100)  // 500GB-5TB
		totalGPU := i % 9  // 0-8个GPU

		// 生成标签
		labels := map[string]interface{}{
			"region":      fmt.Sprintf("region-%d", (i%5)+1),
			"datacenter":  fmt.Sprintf("dc-%d", (i%3)+1),
			"environment": []string{"dev", "test", "prod"}[i%3],
		}
		labelsJSON, _ := json.Marshal(labels)

		// 生成标签数组
		tags := pq.StringArray{
			fmt.Sprintf("gpu-%d", totalGPU),
			fmt.Sprintf("cpu-%d", totalCPU),
			deploymentMode,
		}

		// 创建设备对象
		host := &entity.Host{
			ID:             hostID,
			Name:           fmt.Sprintf("GPU-Server-%03d", i),
			Hostname:       fmt.Sprintf("gpu-server-%03d.example.com", i),
			IPAddress:      ipAddress,
			PublicIP:       publicIP,
			OSType:         osType,
			OSVersion:      osVersion,
			Arch:           "x86_64",
			DeploymentMode: deploymentMode,
			K8sNodeName:    k8sNodeName,
			Status:         "active",
			HealthStatus:   "healthy",
			TotalCPU:       totalCPU,
			TotalMemory:    totalMemory,
			TotalDisk:      totalDisk,
			TotalGPU:       totalGPU,
			UsedCPU:        0,
			UsedMemory:     0,
			UsedDisk:       0,
			UsedGPU:        0,
			SSHPort:        22,
			AgentPort:      8080,
			Labels:         labelsJSON,
			Tags:           tags,
		}

		hosts = append(hosts, host)
	}

	// 批量插入数据
	log.Printf("开始插入 %d 台设备数据...", len(hosts))

	// 使用事务批量插入
	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("开始事务失败: %v", tx.Error)
	}

	// 分批插入，每次插入20条
	batchSize := 20
	for i := 0; i < len(hosts); i += batchSize {
		end := i + batchSize
		if end > len(hosts) {
			end = len(hosts)
		}

		batch := hosts[i:end]
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			log.Fatalf("插入数据失败: %v", err)
		}

		log.Printf("已插入 %d/%d 台设备", end, len(hosts))
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

	log.Printf("成功插入 %d 台设备数据！", len(hosts))

	// 验证插入结果
	var count int64
	db.Model(&entity.Host{}).Count(&count)
	log.Printf("数据库中现有设备总数: %d", count)
}
