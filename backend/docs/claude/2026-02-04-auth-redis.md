# Auth 模块 Redis 集成工作总结

**日期**: 2026-02-04

## 重点改动

### 1. Service 层 (`internal/service/auth/auth_service.go`)

- 新增 `AdminLogin` 方法：验证 admin 角色
- 新增 `Logout` 方法：使用 Redis 存储 token 黑名单
- 新增 `IsTokenBlacklisted` 方法：检查 token 是否已登出
- 登录时更新 `last_login_at` 字段

### 2. Controller 层 (`internal/controller/v1/auth/auth_controller.go`)

- 新增 `AdminLogin` 接口：非 admin 返回 403
- 完善 `Logout` 接口：提取 token 并加入黑名单

### 3. Cache 层 (`pkg/cache/redis.go`)

- 新增 `GlobalCache` 全局变量
- 新增 `GetCache()` 函数返回 `cache.Cache` 接口
- `InitRedis` 中初始化 `GlobalCache`

## Redis 黑名单设计

```
Key: auth:token:blacklist:{token}
Value: "1"
TTL: 1 小时
```

## 测试覆盖

| 测试类型 | 数量 | 状态 |
|---------|------|------|
| HTTP 集成测试 | 6 | ✅ |
| 模拟测试 | 10 | ✅ |
| **总计** | **16** | ✅ |
