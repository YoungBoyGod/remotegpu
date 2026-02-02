package cache

import (
	"context"
	"time"
)

// CacheType 缓存类型
type CacheType string

const (
	// CacheTypeRedis Redis 缓存
	CacheTypeRedis CacheType = "redis"
	// CacheTypeMemory 内存缓存
	CacheTypeMemory CacheType = "memory"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, keys ...string) error

	// Exists 检查键是否存在
	Exists(ctx context.Context, keys ...string) (int64, error)

	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// TTL 获取剩余过期时间
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Close 关闭缓存连接
	Close() error
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type     CacheType     `json:"type"`      // 缓存类型
	Addr     string        `json:"addr"`      // Redis 地址
	Password string        `json:"password"`  // Redis 密码
	DB       int           `json:"db"`        // Redis 数据库
	PoolSize int           `json:"pool_size"` // 连接池大小
	Timeout  time.Duration `json:"timeout"`   // 超时时间
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Type:     CacheTypeRedis,
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
		Timeout:  5 * time.Second,
	}
}
