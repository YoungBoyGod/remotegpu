package dao

import (
	"testing"

	"github.com/YoungBoyGod/remotegpu/pkg/database"
)

// setupTestDB 设置测试数据库连接（用于需要真实数据库的集成测试）
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
