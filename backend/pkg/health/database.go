package health

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgreSQLChecker PostgreSQL健康检查器
type PostgreSQLChecker struct {
	config config.DatabaseConfig
}

// NewPostgreSQLChecker 创建PostgreSQL健康检查器
func NewPostgreSQLChecker(cfg config.DatabaseConfig) *PostgreSQLChecker {
	return &PostgreSQLChecker{config: cfg}
}

// Name 返回服务名称
func (c *PostgreSQLChecker) Name() string {
	return "postgresql"
}

// Check 执行健康检查
func (c *PostgreSQLChecker) Check(ctx context.Context) *CheckResult {
	start := time.Now()
	result := &CheckResult{
		Service:   c.Name(),
		Timestamp: start,
	}

	// 构建DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.config.Host, c.config.Port, c.config.User, c.config.Password, c.config.DBName)

	// 尝试连接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("连接失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}

	// 获取底层连接
	sqlDB, err := db.DB()
	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("获取连接失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer sqlDB.Close()

	// 执行ping测试
	if err := sqlDB.PingContext(ctx); err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("Ping失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}

	result.Status = StatusHealthy
	result.Message = "连接正常"
	result.Latency = time.Since(start)
	result.Details = map[string]interface{}{
		"host":   c.config.Host,
		"port":   c.config.Port,
		"dbname": c.config.DBName,
	}

	return result
}

// RedisChecker Redis健康检查器
type RedisChecker struct {
	config config.RedisConfig
}

// NewRedisChecker 创建Redis健康检查器
func NewRedisChecker(cfg config.RedisConfig) *RedisChecker {
	return &RedisChecker{config: cfg}
}

// Name 返回服务名称
func (c *RedisChecker) Name() string {
	return "redis"
}

// Check 执行健康检查
func (c *RedisChecker) Check(ctx context.Context) *CheckResult {
	start := time.Now()
	result := &CheckResult{
		Service:   c.Name(),
		Timestamp: start,
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
		Password: c.config.Password,
		DB:       c.config.DB,
	})
	defer client.Close()

	// 执行ping测试
	if err := client.Ping(ctx).Err(); err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("Ping失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}

	result.Status = StatusHealthy
	result.Message = "连接正常"
	result.Latency = time.Since(start)
	result.Details = map[string]interface{}{
		"host": c.config.Host,
		"port": c.config.Port,
		"db":   c.config.DB,
	}

	return result
}
