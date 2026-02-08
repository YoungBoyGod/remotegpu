# 后端代码质量审查报告 (R5)

> 审查范围：`backend/internal/controller/`、`backend/internal/router/`、`backend/internal/model/entity/`
> 审查日期：2026-02-07

---

## 1. Controller 错误处理一致性

### 1.1 错误消息语言混用

项目中 Controller 层的错误消息存在 **中英文混用** 问题，缺乏统一规范。

| Controller | 语言 | 示例 |
|---|---|---|
| `auth_controller.go` | 英文 | `"Invalid request parameters"`, `"Authentication failed"` |
| `customer_controller.go` | 混合 | List/Create/Detail 用英文，UpdateQuota/ResourceUsage 用中文 |
| `machine_controller.go` | 混合 | 单机操作用英文，批量操作用中文（`"批量操作失败"`, `"批量分配失败"`） |
| `task_controller.go` | 中文 | `"用户未认证"`, `"获取任务列表失败"` |
| `admin_task_controller.go` | 中文 | `"获取任务列表失败"`, `"任务不存在"` |
| `my_machine_controller.go` | 混合 | 认证/权限用中文，业务错误用英文（`"Failed to get connection info"`） |
| `notification_controller.go` | 中文 | `"未认证"`, `"无效的通知 ID"` |
| `sshkey_controller.go` | 中文 | `"未授权"`, `"获取 SSH 密钥列表失败"` |
| `dataset_controller.go` | 混合 | 认证/权限用中文，InitUpload 用英文（`"Failed to init upload"`） |
| `ops_controller.go` | 混合 | Dashboard 用中文，Monitor 用英文（`"Failed to get snapshot"`） |

**建议**：统一使用中文错误消息（面向国内用户），或统一使用英文（面向国际化）。

### 1.2 错误码使用不一致

部分 Controller 对同类错误使用了不同的 HTTP 状态码：

| 问题 | 位置 | 说明 |
|---|---|---|
| 所有错误返回 `code=0` 的 200 | `base_controller.go:Error()` | `Error()` 方法使用 `ctx.JSON(http.StatusOK, ...)` 返回，HTTP 状态码始终为 200，仅通过 body 中的 code 字段区分。这是设计选择，但需确保前端统一按 code 字段判断 |
| service 错误一律返回 500 | 多处 | 如 `machine_controller.go:Detail()` 对任何错误都返回 404，无法区分"不存在"和"内部错误" |
| `agent_task_controller.go` 使用 409/410 | `StartTask`/`RenewLease` | 与其他 Controller 的错误码风格不同，但语义更准确 |

**建议**：在 Service 层定义明确的错误类型（如 `ErrNotFound`、`ErrConflict`），Controller 层根据错误类型映射到对应的错误码。

### 1.3 参数解析错误未处理

| 位置 | 问题 | 严重程度 |
|---|---|---|
| `customer_controller.go:149` | `Disable()` 中 `id, _ := strconv.ParseUint(...)` 忽略错误，无效 ID 会传 0 给 Service | **高** |
| `dataset_controller.go:90` | `CompleteUpload()` 中 `datasetID, _ := strconv.ParseUint(...)` 忽略错误 | **高** |
| `dataset_controller.go:129` | `Mount()` 中 `datasetID, _ := strconv.ParseUint(...)` 忽略错误 | **高** |
| `dataset_controller.go:182` | `Unmount()` 中 `datasetID, _ := strconv.ParseUint(...)` 忽略错误 | **高** |
| `dataset_controller.go:213` | `ListMounts()` 中 `datasetID, _ := strconv.ParseUint(...)` 忽略错误 | **高** |
| 多处 `strconv.Atoi` 分页参数 | `page, _ := strconv.Atoi(...)` 忽略错误 | **低**（有默认值兜底） |

**建议**：对路径参数的 `ParseUint` 必须检查错误并返回 400。分页参数因有 `DefaultQuery` 兜底，风险较低。

### 1.4 用户认证获取方式不一致

Controller 层获取当前用户 ID 存在两种方式，混合使用：

| 方式 | 使用位置 | 风险 |
|---|---|---|
| `ctx.Get("userID")` + 类型断言 `userID.(uint)` | `task_controller.go`, `my_machine_controller.go`, `dataset_controller.go`, `notification_controller.go` | 若中间件未设置 userID，类型断言会 **panic** |
| `ctx.GetUint("userID")` | `sshkey_controller.go`, `machine_enrollment_controller.go`, `auth_controller.go` | 安全，不存在时返回 0 |

**建议**：统一使用 `ctx.GetUint("userID")`，避免类型断言 panic 风险。`sshkey_controller.go` 中的 `getUserID()` 辅助方法是良好实践，可推广到其他 Controller。

### 1.5 权限校验方式不一致

| 方式 | 使用位置 | 说明 |
|---|---|---|
| 字符串比较 `err.Error() == "无权限访问该资源"` | `task_controller.go:118`, `dataset_controller.go:94,133,149` | 脆弱，依赖错误消息文本 |
| 哨兵错误比较 `err == entity.ErrUnauthorized` | `task_controller.go:138,158` | 正确做法 |
| `errors.Is(err, sshkey.ErrXxx)` | `sshkey_controller.go:65,69,96,100` | 最佳实践 |

**建议**：统一使用 `errors.Is()` 进行错误类型判断，避免字符串比较。

---

## 2. 路由冲突与问题

### 2.1 重复路由注册

`router.go` 中存在重复路由：

```
custGroup.POST("/machines", enrollmentController.Create)       // 第 312 行
custGroup.POST("/machines/enroll", enrollmentController.Create) // 第 340 行
```

两个路由指向同一个 Handler，`POST /customer/machines` 同时承担"创建 enrollment"的职责，与 RESTful 语义冲突（通常 `POST /machines` 表示创建机器本身）。

**建议**：移除 `POST /customer/machines` 路由，仅保留 `POST /customer/machines/enroll`。

### 2.2 SSE 通知路由重复

```
apiV1.GET("/notifications/stream", middleware.Auth(db), notificationController.SSE)  // 第 194 行
custGroup.GET("/notifications/sse", notificationController.SSE)                       // 第 345 行
```

两个 SSE 端点功能相同，路径不同。前者在 `/api/v1/notifications/stream`，后者在 `/api/v1/customer/notifications/sse`。

**建议**：保留一个即可，推荐保留 `custGroup` 下的路径以保持路由分组一致性。

---

## 3. 安全问题

### 3.1 硬编码默认密码

`customer_controller.go:50` 中硬编码了默认密码：

```go
if req.Password == "" {
    req.Password = "ChangeME_123"
    mustChangePassword = true
}
```

虽然设置了 `mustChangePassword = true` 强制用户修改，但默认密码是公开的常量字符串。

**建议**：使用随机生成的临时密码，或通过邮件/短信发送一次性链接。

### 3.2 分页参数未做上限校验

大部分 Controller 的分页参数没有上限校验，用户可传入 `page_size=999999` 导致大量数据查询。

**已做校验的**：`audit_controller.go`、`image_controller.go`（`pageSize > 200` 时重置为 20）

**未做校验的**：`machine_controller.go`、`customer_controller.go`、`task_controller.go`、`admin_task_controller.go`、`notification_controller.go`、`document_controller.go` 等

**建议**：在 `BaseController` 中添加通用的分页参数解析方法，统一校验 page >= 1 和 pageSize 上限。

---

## 4. Entity / GORM 标签审查

### 4.1 billing.go 残留

`backend/internal/model/entity/billing.go` 仍然存在，包含 7 个计费相关实体（BillingRecord、Invoice、BillingPlan、BillingRule、CustomerSubscription、BalanceTransaction、RechargeOrder）。此文件在之前的清理任务中已被删除，但被其他队友重新创建。

**建议**：再次删除 `billing.go`，并确认无其他代码引用这些实体。

### 4.2 customer.go 计费字段残留

`customer.go:39-43` 中仍保留了计费相关字段：

```go
// Billing
Balance       float64 `gorm:"type:decimal(10,4);default:0.00" json:"balance"`
Currency      string  `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
CreditLimit   float64 `gorm:"type:decimal(10,4);default:0" json:"credit_limit"`
BillingPlanID *uint   `json:"billing_plan_id,omitempty"`
```

这些字段在之前的清理任务中已被删除，但被重新引入。

**建议**：再次移除这 4 个字段及 `// Billing` 注释。

### 4.3 GORM 标签一致性

| 问题 | 位置 | 说明 |
|---|---|---|
| 主键类型不统一 | `Host.ID` 为 `varchar(64)` 字符串，`Customer.ID` 为 `uint` | 设计选择，但需注意外键关联时类型匹配 |
| `Task.ID` 为字符串主键 | `task.go:11` | 与 `Allocation.ID` 一致（varchar(64)），但与 Customer/GPU 等 uint 主键不同 |
| 缺少 `TableName()` | `GPU`、`SSHKey`、`Workspace` | 依赖 GORM 自动推断表名，建议显式声明 |
| `Host.SSHPassword` 和 `Host.SSHKey` 用 `json:"-"` 隐藏 | `resource.go:21-22` | 正确做法，敏感字段不暴露 |
| `Host.JupyterToken` 和 `Host.VNCPassword` 用 `json:"-"` | `resource.go:26-28` | 正确 |

### 4.4 machine_controller.go Import 字段映射错误

`machine_controller.go:352` 中批量导入时将 GPU 型号存入了 CPU 信息字段：

```go
CPUInfo: m.GPUModel, // 暂存 GPU 型号信息
```

**建议**：Host 实体应增加 `GPUModel` 字段，或使用独立的 GPU 实体关联。

---

## 5. 其他代码质量问题

### 5.1 notification_controller.go 缺少认证检查

`notification_controller.go` 的 `List`、`UnreadCount`、`MarkRead`、`MarkAllRead` 方法中：

```go
userID, _ := ctx.Get("userID")
customerID := userID.(uint)
```

使用 `_` 忽略了 `exists` 返回值，若 `userID` 不存在会直接 panic。与同文件中 `SSE` 方法的处理方式不一致（SSE 方法正确检查了 `exists`）。

**建议**：所有方法统一检查 `exists`，或使用 `ctx.GetUint("userID")`。

### 5.2 router.go 中 storageMgr 错误被忽略

`router.go:86`：

```go
storageMgr, _ := storage.NewManager(config.GlobalConfig.Storage)
```

存储管理器初始化失败时错误被忽略，后续使用 `storageMgr` 可能为 nil 导致 panic。

**建议**：检查错误并在存储不可用时记录日志或降级处理。

### 5.3 Swagger 注释覆盖不完整

| Controller | Swagger 注释覆盖情况 |
|---|---|
| `auth_controller.go` | 完整 |
| `machine_controller.go` | 部分（Update/Usage/Batch* 缺少） |
| `customer_controller.go` | 缺少所有方法的 Swagger 注释 |
| `task_controller.go` | 缺少 |
| `admin_task_controller.go` | 缺少 |
| `dataset_controller.go` | 缺少 |
| `notification_controller.go` | 缺少 |
| `sshkey_controller.go` | 缺少 |
| `document_controller.go` | 缺少 |
| `system_config_controller.go` | 缺少 |
| `storage_controller.go` | 缺少 |

**建议**：补充所有 Controller 方法的 Swagger 注释（已有 #89 任务在进行中）。

---

## 6. 优秀实践（值得推广）

| 实践 | 位置 | 说明 |
|---|---|---|
| `getUserID()` 辅助方法 | `sshkey_controller.go:24-32` | 封装认证检查，避免重复代码 |
| `errors.Is()` 错误判断 | `sshkey_controller.go` | 使用哨兵错误 + `errors.Is()` 判断，类型安全 |
| 分页参数校验 | `audit_controller.go:28-33` | 对 page/pageSize 做范围校验 |
| 敏感字段 `json:"-"` | `resource.go` Host 实体 | SSHPassword、SSHKey、JupyterToken、VNCPassword 不暴露 |
| 权限校验 | `task_controller.go` | 所有操作前校验任务归属，防止越权 |
| `AppError` 模式 | `auth_controller.go` | Service 层返回 AppError，Controller 层提取 code/msg |

---

## 7. 问题汇总与优先级

### P0 — 必须修复（可能导致 panic 或安全问题）

| # | 问题 | 位置 |
|---|---|---|
| 1 | `ctx.Get("userID")` 类型断言可能 panic | `task_controller.go`, `my_machine_controller.go`, `dataset_controller.go`, `notification_controller.go` |
| 2 | `ParseUint` 错误被忽略（传 0 给 Service） | `customer_controller.go:149`, `dataset_controller.go:90,129,182,213` |
| 3 | `notification_controller.go` List/UnreadCount/MarkRead/MarkAllRead 未检查 userID 存在性 | `notification_controller.go:66,83,96,114` |

### P1 — 应当修复（代码质量 / 一致性问题）

| # | 问题 | 位置 |
|---|---|---|
| 4 | 错误消息中英文混用 | 全部 Controller |
| 5 | 权限校验使用字符串比较而非 `errors.Is()` | `task_controller.go:118`, `dataset_controller.go:94,133,149` |
| 6 | 重复路由 `POST /customer/machines` | `router.go:312,340` |
| 7 | SSE 通知路由重复 | `router.go:194,345` |
| 8 | `billing.go` 和 `customer.go` 计费字段残留 | `entity/billing.go`, `entity/customer.go:39-43` |
| 9 | `storageMgr` 初始化错误被忽略 | `router.go:86` |

### P2 — 建议改进（提升代码质量）

| # | 问题 | 位置 |
|---|---|---|
| 10 | 分页参数缺少上限校验 | 大部分 Controller |
| 11 | 硬编码默认密码 `"ChangeME_123"` | `customer_controller.go:50` |
| 12 | GPU 型号存入 CPUInfo 字段 | `machine_controller.go:352` |
| 13 | 缺少 `TableName()` 方法 | `GPU`、`SSHKey`、`Workspace` 实体 |
| 14 | Swagger 注释覆盖不完整 | 大部分 Controller |
