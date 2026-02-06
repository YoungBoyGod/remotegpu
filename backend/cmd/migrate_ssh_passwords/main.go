package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/crypto"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
)

// migrate_ssh_passwords - 迁移 SSH 密码加密
// @author Claude
// @description 将数据库中的明文 SSH 密码加密存储
// @usage go run cmd/migrate_ssh_passwords/main.go --dry-run
// @modified 2026-02-06

var (
	dryRun = flag.Bool("dry-run", false, "只检查不实际修改数据库")
	force  = flag.Bool("force", false, "强制执行，跳过确认")
)

func main() {
	flag.Parse()

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

	// 统计需要迁移的记录
	var hosts []entity.Host
	if err := db.WithContext(ctx).Find(&hosts).Error; err != nil {
		log.Fatalf("查询 hosts 失败: %v", err)
	}

	needMigrate := 0
	for _, host := range hosts {
		if host.SSHPassword != "" {
			needMigrate++
		}
	}

	fmt.Printf("\n=== SSH 密码加密迁移工具 ===\n")
	fmt.Printf("总记录数: %d\n", len(hosts))
	fmt.Printf("需要迁移: %d\n", needMigrate)
	fmt.Printf("模式: ")
	if *dryRun {
		fmt.Printf("试运行（不会修改数据）\n\n")
	} else {
		fmt.Printf("实际执行\n\n")
	}

	if needMigrate == 0 {
		fmt.Println("没有需要迁移的记录")
		return
	}

	// 确认执行
	if !*dryRun && !*force {
		fmt.Print("确认执行迁移？这将修改数据库中的 SSH 密码字段。(yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("已取消")
			return
		}
	}

	// 执行迁移
	migrated := 0
	failed := 0
	skipped := 0

	for i, host := range hosts {
		if host.SSHPassword == "" {
			skipped++
			continue
		}

		// 检查是否已经加密（简单判断：加密后的字符串通常较长且包含特殊字符）
		if isLikelyEncrypted(host.SSHPassword) {
			fmt.Printf("[%d/%d] 跳过 %s (已加密)\n", i+1, len(hosts), host.ID)
			skipped++
			continue
		}

		fmt.Printf("[%d/%d] 迁移 %s ... ", i+1, len(hosts), host.ID)

		if *dryRun {
			fmt.Println("OK (试运行)")
			migrated++
			continue
		}

		// 加密密码
		encrypted, err := crypto.EncryptAES256GCM(host.SSHPassword)
		if err != nil {
			fmt.Printf("失败: %v\n", err)
			failed++
			continue
		}

		// 更新数据库
		if err := db.WithContext(ctx).Model(&entity.Host{}).
			Where("id = ?", host.ID).
			Update("ssh_password", encrypted).Error; err != nil {
			fmt.Printf("失败: %v\n", err)
			failed++
			continue
		}

		fmt.Println("OK")
		migrated++
	}

	// 输出结果
	fmt.Printf("\n=== 迁移完成 ===\n")
	fmt.Printf("成功: %d\n", migrated)
	fmt.Printf("失败: %d\n", failed)
	fmt.Printf("跳过: %d\n", skipped)

	if failed > 0 {
		os.Exit(1)
	}
}

// isLikelyEncrypted 简单判断字符串是否可能已加密
// 加密后的 base64 字符串通常较长（>40字符）且只包含 base64 字符
func isLikelyEncrypted(s string) bool {
	if len(s) < 40 {
		return false
	}
	// base64 字符集：A-Z, a-z, 0-9, +, /, =
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '+' || c == '/' || c == '=') {
			return false
		}
	}
	return true
}
