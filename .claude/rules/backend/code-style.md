---
paths:
  - "backend/**/*.go"
---

# 后端代码风格

## 命名约定

- DAO 文件：`xxx_repo.go`，结构体 `XxxDao`
- Service 文件：`xxx_service.go`，结构体 `XxxService`
- Controller 文件：`xxx_controller.go`，结构体 `XxxController`
- Entity 文件：按领域分组（如 `resource.go` 包含 Host 和 GPU）

## Controller 规范

- 嵌入 `common.BaseController`
- 成功响应：`c.Success(ctx, data)`
- 错误响应：`c.Error(ctx, code, msg)`
- 参数绑定用 `ctx.ShouldBindJSON(&req)`

## DAO 规范

- 简单实体继承 `BaseDao[T]`（提供 Create/Update/Delete/FindByID/FindByUUID）
- 所有方法第一个参数为 `context.Context`
- 使用 `d.db.WithContext(ctx)` 确保可取消
- 复杂操作使用 `d.db.WithContext(ctx).Transaction()`

## Service 规范

- 构造函数 `NewXxxService(db *gorm.DB, ...)` 内部创建 DAO
- 不直接操作 `*gorm.DB`，通过 DAO 访问数据

## 路由注册

- 所有路由在 `router/router.go` 的 `InitRouter` 中注册
- Service 和 Controller 在此函数内初始化
- import 别名：Service 用 `serviceXxx`，Controller 用 `ctrlXxx`

## 统一响应格式

```json
{ "code": 0, "msg": "success", "data": ... }
```

code=0 表示成功，非零为错误码。
