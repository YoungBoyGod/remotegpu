package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis 缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建 Redis 缓存
func NewRedisCache(config *CacheConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Get 获取缓存值
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("键不存在: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("获取缓存失败: %w", err)
	}
	return val, nil
}

// Set 设置缓存值
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case []byte:
		val = string(v)
	default:
		// 其他类型序列化为 JSON
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("序列化值失败: %w", err)
		}
		val = string(data)
	}

	if err := c.client.Set(ctx, key, val, expiration).Err(); err != nil {
		return fmt.Errorf("设置缓存失败: %w", err)
	}
	return nil
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}
	return nil
}

// Exists 检查键是否存在
func (c *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := c.client.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("检查键存在失败: %w", err)
	}
	return count, nil
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := c.client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("设置过期时间失败: %w", err)
	}
	return nil
}

// TTL 获取剩余过期时间
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取过期时间失败: %w", err)
	}
	return ttl, nil
}

// Close 关闭缓存连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}
