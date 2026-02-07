package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/pkg/crypto"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"golang.org/x/crypto/ssh"
)

// test_ssh_connection - 测试 SSH 连接和密码解密
// @author Claude
// @description 验证 P0 安全修复：SSH 密码解密功能
// @usage ENCRYPTION_KEY="..." go run cmd/test_ssh_connection/main.go <host_id>
// @modified 2026-02-06

func main() {
	if len(os.Args) < 2 {
		log.Fatal("用法: go run main.go <host_id>")
	}
	hostID := os.Args[1]

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
	machineDao := dao.NewMachineDao(database.DB)

	// 查询机器信息
	host, err := machineDao.FindByID(ctx, hostID)
	if err != nil {
		log.Fatalf("查询机器失败: %v", err)
	}

	fmt.Println("\n=== SSH 连接测试 ===")
	fmt.Printf("机器 ID: %s\n", host.ID)
	fmt.Printf("IP 地址: %s\n", host.IPAddress)
	fmt.Printf("SSH 用户: %s\n", host.SSHUsername)
	fmt.Printf("加密密码长度: %d 字节\n", len(host.SSHPassword))
	fmt.Println()

	// 解密密码
	fmt.Println("步骤 1: 解密 SSH 密码...")
	decryptedPassword, err := crypto.DecryptAES256GCM(host.SSHPassword)
	if err != nil {
		log.Fatalf("❌ 解密失败: %v", err)
	}
	fmt.Printf("✅ 解密成功，密码长度: %d\n", len(decryptedPassword))
	fmt.Println()

	// 测试 SSH 连接
	fmt.Println("步骤 2: 测试 SSH 连接...")
	if err := testSSHConnection(host.IPAddress, host.SSHPort, host.SSHUsername, decryptedPassword); err != nil {
		log.Fatalf("❌ SSH 连接失败: %v", err)
	}
	fmt.Println("✅ SSH 连接成功！")
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
	fmt.Println("✅ 密码加密/解密功能正常")
	fmt.Println("✅ SSH 连接功能正常")
}

func testSSHConnection(host string, port int, username, password string) error {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer client.Close()

	// 执行简单命令验证连接
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput("hostname")
	if err != nil {
		return fmt.Errorf("执行命令失败: %w", err)
	}

	fmt.Printf("  远程主机名: %s", string(output))
	return nil
}
