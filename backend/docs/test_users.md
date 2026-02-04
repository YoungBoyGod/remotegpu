# 测试用户账号

本文档记录 Auth 模块测试中使用的用户账号信息。

## 用户列表

| 用户名 | 密码 | 角色 | 邮箱 | 说明 |
|-------|------|------|------|------|
| admin | luoyang@123 | admin | admin@example.com | 管理员账号 |
| user | user@123 | customer_owner | user@example.com | 普通用户 |
| testuser | Test123456 | customer_owner | test@example.com | 测试用户 |

## 测试文件

- `auth_controller_test.go` - 模拟 HTTP 测试
- `auth_controller_http_test.go` - 真实 HTTP 集成测试

## 运行测试

```bash
# 运行全部测试（忽略缓存）
go test -v -count=1 ./internal/controller/v1/auth/...

# 只运行 HTTP 测试
go test -v -count=1 ./internal/controller/v1/auth/... -run TestHTTP

# 只运行 Admin 测试
go test -v -count=1 ./internal/controller/v1/auth/... -run TestAdmin
```
