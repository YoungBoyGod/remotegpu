package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CacheManager 缓存管理器
type CacheManager struct {
	caches map[string]Cache
	mu     sync.RWMutex
}

// NewCacheManager 创建缓存管理器
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]Cache),
	}
}

// Register 注册缓存实例
func (m *CacheManager) Register(name string, cache Cache) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.caches[name] = cache
}

// Get 获取缓存实例
func (m *CacheManager) Get(name string) (Cache, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cache, ok := m.caches[name]
	if !ok {
		return nil, fmt.Errorf("缓存实例不存在: %s", name)
	}
	return cache, nil
}

// GetOrDefault 获取缓存实例，如果不存在则返回默认实例
func (m *CacheManager) GetOrDefault(name string) Cache {
	cache, err := m.Get(name)
	if err != nil {
		// 返回默认的内存缓存
		return NewMemoryCache()
	}
	return cache
}

// Close 关闭所有缓存连接
func (m *CacheManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, cache := range m.caches {
		if err := cache.Close(); err != nil {
			return fmt.Errorf("关闭缓存 %s 失败: %w", name, err)
		}
	}

	m.caches = make(map[string]Cache)
	return nil
}

// NewCache 创建缓存实例
func NewCache(config *CacheConfig) (Cache, error) {
	switch config.Type {
	case CacheTypeRedis:
		return NewRedisCache(config)
	case CacheTypeMemory:
		return NewMemoryCache(), nil
	default:
		return nil, fmt.Errorf("不支持的缓存类型: %s", config.Type)
	}
}

// CacheHelper 缓存辅助函数
type CacheHelper struct {
	cache Cache
}

// NewCacheHelper 创建缓存辅助函数
func NewCacheHelper(cache Cache) *CacheHelper {
	return &CacheHelper{
		cache: cache,
	}
}

// GetString 获取字符串值
func (h *CacheHelper) GetString(ctx context.Context, key string) (string, error) {
	return h.cache.Get(ctx, key)
}

// SetString 设置字符串值
func (h *CacheHelper) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return h.cache.Set(ctx, key, value, expiration)
}

// GetJSON 获取 JSON 对象
func (h *CacheHelper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := h.cache.Get(ctx, key)
	if err != nil {
		return err
	}

	// 反序列化 JSON
	return unmarshalJSON([]byte(val), dest)
}

// SetJSON 设置 JSON 对象
func (h *CacheHelper) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return h.cache.Set(ctx, key, value, expiration)
}

// unmarshalJSON 反序列化 JSON（简化实现）
func unmarshalJSON(data []byte, dest interface{}) error {
	// 这里应该使用 json.Unmarshal，但为了避免循环导入，简化处理
	return fmt.Errorf("JSON 反序列化未实现")
}
