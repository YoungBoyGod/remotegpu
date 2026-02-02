package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// cacheItem 缓存项
type cacheItem struct {
	value      string
	expiration time.Time
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	items map[string]*cacheItem
	mu    sync.RWMutex
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*cacheItem),
	}

	// 启动清理过期项的 goroutine
	go cache.cleanupExpired()

	return cache
}

// cleanupExpired 清理过期项
func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if !item.expiration.IsZero() && now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// Get 获取缓存值
func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return "", fmt.Errorf("键不存在: %s", key)
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return "", fmt.Errorf("键已过期: %s", key)
	}

	return item.value, nil
}

// Set 设置缓存值
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
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

	c.mu.Lock()
	defer c.mu.Unlock()

	item := &cacheItem{
		value: val,
	}

	if expiration > 0 {
		item.expiration = time.Now().Add(expiration)
	}

	c.items[key] = item
	return nil
}

// Delete 删除缓存
func (c *MemoryCache) Delete(ctx context.Context, keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.items, key)
	}
	return nil
}

// Exists 检查键是否存在
func (c *MemoryCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var count int64
	now := time.Now()

	for _, key := range keys {
		if item, ok := c.items[key]; ok {
			// 检查是否过期
			if item.expiration.IsZero() || now.Before(item.expiration) {
				count++
			}
		}
	}

	return count, nil
}

// Expire 设置过期时间
func (c *MemoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return fmt.Errorf("键不存在: %s", key)
	}

	if expiration > 0 {
		item.expiration = time.Now().Add(expiration)
	} else {
		item.expiration = time.Time{}
	}

	return nil
}

// TTL 获取剩余过期时间
func (c *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return 0, fmt.Errorf("键不存在: %s", key)
	}

	if item.expiration.IsZero() {
		return -1, nil // 永不过期
	}

	ttl := time.Until(item.expiration)
	if ttl < 0 {
		return 0, nil // 已过期
	}

	return ttl, nil
}

// Close 关闭缓存连接
func (c *MemoryCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}
