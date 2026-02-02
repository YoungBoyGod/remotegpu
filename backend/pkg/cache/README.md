# pkg/cache - 缓存管理模块

## 功能

统一的缓存管理接口，支持多种缓存后端：
- Redis 缓存
- 内存缓存

## 使用示例

### 1. 创建 Redis 缓存

```go
package main

import (
    "context"
    "time"
    "github.com/YoungBoyGod/remotegpu/pkg/cache"
)

func main() {
    // 创建 Redis 缓存配置
    config := &cache.CacheConfig{
        Type:     cache.CacheTypeRedis,
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
        PoolSize: 10,
        Timeout:  5 * time.Second,
    }

    // 创建 Redis 缓存实例
    redisCache, err := cache.NewRedisCache(config)
    if err != nil {
        panic(err)
    }
    defer redisCache.Close()

    ctx := context.Background()

    // 设置缓存
    err = redisCache.Set(ctx, "key1", "value1", 10*time.Minute)
    if err != nil {
        panic(err)
    }

    // 获取缓存
    val, err := redisCache.Get(ctx, "key1")
    if err != nil {
        panic(err)
    }
    println(val) // 输出: value1
}
```

### 2. 创建内存缓存

```go
func main() {
    // 创建内存缓存实例
    memCache := cache.NewMemoryCache()
    defer memCache.Close()

    ctx := context.Background()

    // 设置缓存
    err := memCache.Set(ctx, "key1", "value1", 10*time.Minute)
    if err != nil {
        panic(err)
    }

    // 获取缓存
    val, err := memCache.Get(ctx, "key1")
    if err != nil {
        panic(err)
    }
    println(val) // 输出: value1
}
```

### 3. 使用缓存管理器

```go
func main() {
    // 创建缓存管理器
    manager := cache.NewCacheManager()
    defer manager.Close()

    // 注册 Redis 缓存
    redisConfig := cache.DefaultCacheConfig()
    redisCache, _ := cache.NewCache(redisConfig)
    manager.Register("redis", redisCache)

    // 注册内存缓存
    memCache := cache.NewMemoryCache()
    manager.Register("memory", memCache)

    // 获取缓存实例
    cache, _ := manager.Get("redis")
    ctx := context.Background()
    cache.Set(ctx, "key1", "value1", 10*time.Minute)
}
```

### 4. 在 Service 层使用

```go
// internal/service/environment.go

type EnvironmentService struct {
    cacheManager *cache.CacheManager
    // ...
}

func NewEnvironmentService() *EnvironmentService {
    manager := cache.NewCacheManager()

    // 注册 Redis 缓存
    redisConfig := &cache.CacheConfig{
        Type: cache.CacheTypeRedis,
        Addr: "localhost:6379",
    }
    redisCache, _ := cache.NewCache(redisConfig)
    manager.Register("default", redisCache)

    return &EnvironmentService{
        cacheManager: manager,
        // ...
    }
}

func (s *EnvironmentService) SaveAccessInfo(envID string, info interface{}) error {
    cache, _ := s.cacheManager.Get("default")
    ctx := context.Background()

    key := "env:access_info:" + envID
    return cache.Set(ctx, key, info, 24*time.Hour)
}
```

## API 接口

### Cache 接口

```go
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, keys ...string) (int64, error)
    Expire(ctx context.Context, key string, expiration time.Duration) error
    TTL(ctx context.Context, key string) (time.Duration, error)
    Close() error
}
```

## 特性

- **统一接口**: 所有缓存后端实现相同的接口
- **类型安全**: 支持字符串、字节数组和 JSON 对象
- **过期管理**: 支持设置过期时间和 TTL 查询
- **连接池**: Redis 缓存支持连接池配置
- **自动清理**: 内存缓存自动清理过期项
