# 已实现功能合理性审查报告

> 版本：v1.0 | 日期：2026-02-06 | 作者：产品经理

## 1. 审查范围

对照需求文档 `docs/design/requirements.md`，对已实现功能的代码进行合理性审查，重点关注：
- 并发安全
- 权限隔离
- 状态管理
- 错误处理
- 安全合规

---

## 2. 高严重度问题

### 2.1 机器分配缺少行级锁（并发安全）

- **文件**: `backend/internal/service/allocation/allocation_service.go`
- **问题**: `AllocateMachine` 在事务中检查机器状态时没有使用 `SELECT ... FOR UPDATE`，两个并发请求可能同时读到同一台机器为 "idle" 状态，导致同一台机器被分配给两个客户
- **代码注释已标注**: 注释写了"如果可能应加锁"，但未实际实现
- **建议**: 使用 `tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&host, ...)` 加行级锁

### 2.2 机器回收未检查运行中的任务（业务逻辑）

- **文件**: `backend/internal/service/allocation/allocation_service.go`
- **问题**: `ReclaimMachine` 直接将分配状态改为 "reclaimed"，没有检查该机器上是否有正在运行的任务
- **风险**: 回收后任务仍在执行，但分配关系已断开，可能导致数据不一致
- **建议**: 回收前查询该机器上 status 为 running/assigned 的任务，要求先停止或强制停止

### 2.3 任务停止无超时控制（可用性）

- **文件**: `backend/internal/service/task/task_service.go`
- **问题**: `stopTaskProcess` 调用 `s.agentService.StopProcess()` 时没有设置超时，如果 Agent 无响应，HTTP 请求会一直挂起
- **建议**: 使用 `context.WithTimeout(ctx, 10*time.Second)` 包装 Agent 调用

### 2.4 customer_owner / customer_member 权限未区分（权限）

- **文件**: `backend/internal/router/router.go`, `backend/internal/middleware/role.go`
- **问题**: 实体层定义了 `admin`、`customer_owner`、`customer_member` 三种角色，但路由层只有 `RequireAdmin()` 中间件，Customer 路由没有区分 owner 和 member 的权限
- **影响**: customer_member 可以执行所有 customer_owner 的操作（如管理 SSH 密钥、创建任务等）
- **建议**: 根据业务需要，明确 owner 和 member 的权限边界，添加对应的中间件

### 2.5 机器列表无数据隔离（权限）

- **文件**: `backend/internal/dao/machine_repo.go`
- **问题**: `MachineDao.List()` 返回所有机器，没有按客户权限过滤
- **说明**: 客户端的机器列表通过 `AllocationService` 查询已分配的机器来实现隔离，这个设计是合理的。但 `MachineDao.List()` 本身没有隔离，如果被错误调用可能泄露数据
- **建议**: 确认所有客户端路由都通过 AllocationService 而非直接调用 MachineDao.List()

---

## 3. 中严重度问题

### 3.1 机器创建存在竞态条件

- **文件**: `backend/internal/service/machine/machine_service.go`
- **问题**: `CreateMachine` 先查询 IP/Hostname 是否存在，再执行创建，两步之间存在时间窗口
- **建议**: 依赖数据库 UNIQUE 约束作为最后防线，捕获唯一约束冲突错误并返回友好提示

### 3.2 机器状态转换无校验

- **文件**: `backend/internal/dao/machine_repo.go`
- **问题**: `UpdateStatus` 允许任意状态转换，没有状态机校验
- **风险**: 可能出现非法状态转换（如从 offline 直接到 allocated）
- **建议**: 在 Service 层添加状态转换合法性校验，定义允许的转换路径

### 3.3 改密接口未校验新密码强度

- **文件**: `backend/internal/service/auth/auth_service.go`
- **问题**: `ChangePassword` 没有调用已有的 `ValidateStrength()` 方法，新密码没有强度检查
- **风险**: 用户可以将密码改为 "123456" 等弱密码
- **建议**: 在 ChangePassword 中调用 `auth.ValidateStrength(newPassword, PasswordStrengthWeak)`

### 3.4 JWT Claims 缺少租户标识

- **文件**: `backend/pkg/auth/jwt.go`
- **问题**: Claims 只包含 `UserID`、`Username`、`Role`，没有 `company_code` 或 `tenant_id`
- **影响**: 多租户隔离只能通过 userID 关联查询实现，无法在中间件层直接过滤
- **建议**: 如果未来需要更严格的多租户隔离，考虑在 Claims 中加入 `company_code`

### 3.5 默认密码硬编码

- **文件**: `backend/internal/controller/v1/customer/customer_controller.go`
- **问题**: 默认密码 `ChangeME_123` 硬编码在控制器代码中
- **建议**: 移到系统配置表 `system_configs` 或环境变量中

### 3.6 客户创建缺少输入验证

- **文件**: `backend/internal/controller/v1/customer/customer_controller.go`
- **问题**: `Create` 方法没有验证 Username、Email 等字段的格式有效性
- **说明**: `CreateCustomerRequest` 使用了 `binding:"required,email"` 等 tag，Gin 的 `ShouldBindJSON` 会做基础校验，但错误信息直接暴露了内部验证细节
- **建议**: 统一错误信息格式，避免暴露内部字段名

### 3.7 任务停止缺少状态检查

- **文件**: `backend/internal/service/task/task_service.go`
- **问题**: `StopTask` 和 `StopTaskWithAuth` 没有检查任务当前状态是否允许停止
- **风险**: 已完成（completed）或已失败（failed）的任务也可以被"停止"
- **建议**: 添加状态检查，只允许 pending/assigned/running 状态的任务被停止

---

## 4. 低严重度问题

### 4.1 密码最小长度偏弱

- **文件**: `backend/pkg/auth/password.go`
- **问题**: 最小密码长度为 8 位
- **说明**: 8 位在当前安全标准下偏弱，但考虑到有大小写+数字的复杂度要求，可接受
- **建议**: 后续可考虑提升到 10-12 位

### 4.2 登录错误信息统一（安全设计合理）

- **文件**: `backend/internal/service/auth/auth_service.go`
- **说明**: 用户不存在和密码错误都返回 `ErrorPasswordIncorrect`，这是防止用户枚举的安全设计，**合理**
- **但**: 账号禁用返回了不同的错误码 `ErrorUserDisabled`，这会泄露账号存在信息
- **建议**: 评估是否需要统一所有登录失败的错误码

---

## 5. 已实现功能合理性总结

### 5.1 设计合理的部分

| 模块 | 评价 |
|------|------|
| JWT + Refresh Token 机制 | 合理，refresh_token 存 Redis，支持轮换和黑名单 |
| 首次登录强制改密 | 合理，`must_change_password` 标记 + 前端路由守卫 |
| 审计日志中间件 | 合理，Admin 路由统一挂载审计中间件 |
| Agent 任务租约机制 | 合理，`lease_expires_at` 防止任务卡死 |
| 客户端机器隔离 | 合理，通过 AllocationService 查询已分配机器 |
| 密码哈希存储 | 合理，使用 bcrypt + DefaultCost |
| 登录防枚举 | 合理，统一错误码 |
| 分层架构 | 合理，Controller → Service → DAO 职责清晰 |

### 5.2 需要改进的部分

| 优先级 | 问题 | 影响 |
|--------|------|------|
| P0 | 分配机器缺少行级锁 | 同一机器可能被分配给两个客户 |
| P0 | 回收机器未检查运行中任务 | 数据不一致 |
| P1 | 任务停止无超时控制 | Agent 无响应时请求挂起 |
| P1 | owner/member 权限未区分 | 权限模型不完整 |
| P1 | 改密无密码强度校验 | 弱密码风险 |
| P2 | 状态转换无校验 | 非法状态转换 |
| P2 | 默认密码硬编码 | 可维护性差 |
| P2 | JWT 缺少租户标识 | 多租户扩展受限 |

---

## 6. 建议修复顺序

1. **立即修复（P0）**: 分配机器加行级锁、回收前检查任务
2. **尽快修复（P1）**: 任务停止加超时、改密加强度校验、明确 owner/member 权限边界
3. **计划修复（P2）**: 状态机校验、默认密码配置化、JWT 加租户标识
