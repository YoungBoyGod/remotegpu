package main

import (
	"context"
	"fmt"
	"log"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "运维与辅助工具",
}

var genPassCmd = &cobra.Command{
	Use:   "gen-pass [password]",
	Short: "生成 Bcrypt 密码哈希",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]
		hash, err := auth.HashPassword(password)
		if err != nil {
			fmt.Printf("生成哈希失败: %v\n", err)
			return
		}
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Hash:     %s\n", hash)
	},
}

var resetDbCmd = &cobra.Command{
	Use:   "reset-db",
	Short: "重置数据库 (清空所有数据!)",
	Run: func(cmd *cobra.Command, args []string) {
		initDBOrDie()
		fmt.Println("正在重置数据库...")
		db := database.GetDB()
		
		if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
			log.Fatalf("重置失败: %v", err)
		}
		
		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
			log.Fatalf("恢复扩展失败: %v", err)
		}

		fmt.Println("数据库已重置成功！")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库迁移",
	Run: func(cmd *cobra.Command, args []string) {
		initDBOrDie()
		fmt.Println("正在执行数据库迁移...")
		err := database.GetDB().AutoMigrate(
			&entity.Customer{},
			&entity.SSHKey{},
			&entity.Workspace{},
			&entity.Host{},
			&entity.GPU{},
			&entity.Allocation{},
			&entity.Image{},
			&entity.Dataset{},
			&entity.DatasetMount{},
			&entity.Task{},
			&entity.AuditLog{},
			&entity.AlertRule{},
			&entity.ActiveAlert{},
		)
		if err != nil {
			log.Fatalf("迁移失败: %v", err)
		}
		fmt.Println("数据库迁移完成！")
	},
}

// Admin Commands
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "管理员账号管理",
}

var (
	adminUser string
	adminPass string
	adminEmail string
)

var createAdminCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新的管理员用户",
	Run: func(cmd *cobra.Command, args []string) {
		initDBOrDie()
		db := database.GetDB()
		customerDao := dao.NewCustomerDao(db)

		// Check if exists
		if _, err := customerDao.FindByUsername(context.Background(), adminUser); err == nil {
			log.Fatalf("用户 %s 已存在", adminUser)
		}

		// Hash Password
		hash, err := auth.HashPassword(adminPass)
		if err != nil {
			log.Fatalf("密码加密失败: %v", err)
		}

		admin := &entity.Customer{
			Username:     adminUser,
			Email:        adminEmail,
			PasswordHash: hash,
			Role:         "admin",
			UserType:     "admin",
			DisplayName:  "Administrator",
			Status:       "active",
		}

		if err := customerDao.Create(context.Background(), admin); err != nil {
			log.Fatalf("创建管理员失败: %v", err)
		}

		fmt.Printf("管理员创建成功!\n用户名: %s\n邮箱: %s\n", adminUser, adminEmail)
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
	toolsCmd.AddCommand(genPassCmd)
	toolsCmd.AddCommand(resetDbCmd)
	toolsCmd.AddCommand(migrateCmd) // Corrected line

	rootCmd.AddCommand(adminCmd)
	adminCmd.AddCommand(createAdminCmd)

	createAdminCmd.Flags().StringVarP(&adminUser, "user", "u", "admin", "用户名")
	createAdminCmd.Flags().StringVarP(&adminPass, "password", "p", "admin123", "密码")
	createAdminCmd.Flags().StringVarP(&adminEmail, "email", "e", "admin@localhost", "邮箱")
}

func initDBOrDie() {
	if err := config.LoadConfig(configPath); err != nil {
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
		log.Fatalf("初始化数据库失败: %v", err)
	}
}