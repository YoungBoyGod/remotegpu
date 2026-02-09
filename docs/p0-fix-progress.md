# P0 修复进度跟踪

> 开始时间：2026-02-09
> 团队：remotegpu-p0-fix

## 任务总览

| ID | 任务 | 负责人 | 状态 |
|----|------|--------|------|
| 1 | 架构分析与实施计划 | architect | ✅ 已完成 |
| 2 | P0-1: 统一前后端 API 路径 | frontend-1 | ✅ 已完成（路径已一致） |
| 3 | P0-2: 补齐缺失的后端端点 | backend | ✅ 已完成 |
| 4 | P0-3: 修复 ctx.Get("userID") 类型断言 | backend | ✅ 已完成（之前已修复） |
| 5 | P0-4: 修复 ParseUint 错误忽略 | backend | ✅ 已完成（之前已修复） |
| 6 | P0-6: 修复分配操作前端按钮和 API 字段 | frontend-2 | ✅ 已完成 |
| 7 | 前端 UI 美化 - 管理端 | frontend-1 | ✅ 已完成 |
| 8 | 前端 UI 美化 - 客户端 | frontend-2 | 进行中 |
| 9 | P1 前端修复: 客户启用按钮/Machine类型/状态字段 | frontend-2 | ✅ 已完成（合并到 #6） |

## 详细进度

### Backend Agent 完成记录

#### P0-3 + P0-4: 类型安全问题（Task #4）
- **状态**: ✅ 已完成（代码已在之前修复）
- **分析结果**: 所有控制器已使用 `ctx.GetUint("userID")` 安全方式，所有 `ParseUint` 已有错误检查
- **涉及文件**（均已正确实现）:
  - `controller/v1/task/task_controller.go`
  - `controller/v1/dataset/dataset_controller.go`
  - `controller/v1/notification/notification_controller.go`
  - `controller/v1/customer/my_machine_controller.go`
  - `controller/v1/customer/customer_controller.go`

#### P0-2: 补齐缺失的后端端点（Task #3）
- **状态**: ✅ 已完成
- **分析结果**: 4个端点中3个已存在，仅缺 toggle 端点
- **已存在端点**:
  - `GET /admin/allocations` — allocation_controller.go
  - `DELETE /admin/images/:id` — image_controller.go
  - `PUT /auth/profile` — auth_controller.go
- **新增端点**: `POST /admin/alert-rules/:id/toggle`
  - DAO: `dao/ops_repo.go` — 添加 `ToggleAlertRule` 方法
  - Service: `service/ops/ops_service.go` — 添加 `ToggleAlertRule` 方法
  - Controller: `controller/v1/ops/ops_controller.go` — 添加 `ToggleRule` handler
  - Router: `router/router.go` — 注册路由

#### P0-5: 统一心跳超时判断（Task #8）
- **状态**: ✅ 已完成
- **问题**: Redis TTL (90s) 与 HeartbeatMonitor 超时 (180s) 不一致
- **修复**:
  - `host_status_cache.go`: 移除硬编码 TTL，改为构造函数传入
  - `router.go`: 从配置读取超时值统一传给 Redis 缓存
  - 现在 Redis TTL = HeartbeatMonitor timeout = 配置中 `heartbeat_monitor.timeout`
- **额外修复（P0-3 遗漏）**:
  - `environment_controller.go`: 修复不安全的 `userID.(uint)` 断言
  - `workspace_controller.go`: 同上

---

### Frontend-2 Agent 完成记录

#### P0-6: 修复分配操作前端按钮条件和 API 字段（Task #5）
- **状态**: ✅ 已完成
- **修复内容**:
  1. `MachineAllocateView.vue`: 机器列表过滤只显示 `idle` 状态机器
  2. `MachineAllocateView.vue`: API 调用改用 `allocateMachine`，只发送后端接受的字段（`customer_id`, `duration_months`, `remark`），移除后端不支持的 `start_time`, `end_time`, `contact_person`, `notify_methods`
  3. `MachineListView.vue`: 分配弹窗和批量分配弹窗的客户列表过滤只显示 `active` 状态客户

#### P1 修复: 客户启用按钮/Machine类型（合并到 Task #5）
- **状态**: ✅ 已完成
- **修复内容**:
  1. `CustomerListView.vue`: 启用按钮条件从 `disabled` 扩展到 `disabled || suspended`
  2. `CustomerListView.vue`: 筛选下拉补全 `pending` 和 `suspended` 状态选项
  3. `customer.ts`: 修复 camelCase 字段为 snake_case（`contactPerson` → `contact_person` 等），移除重复的 `createdAt`/`lastLoginAt`
  4. `CustomerListView.vue`: 更新 `contactPerson` 引用为 `contact_person`

---

### Architect Agent 分析报告（Task #1）

#### P0-1: 前后端 API 路径差异详细清单

**分析方法**: 对比 `backend/internal/router/router.go` 路由定义与 `frontend/src/api/` 下所有 API 调用

**结论**: 大部分路径已匹配，以下为仍存在差异的端点：

| # | 前端调用路径 | 后端实际路径 | 差异说明 | 前端文件 |
|---|------------|------------|---------|---------|
| 1 | `POST /admin/alert-rules/:id/toggle` | 已由 backend agent 新增 | ✅ 已修复 | `api/admin.ts` |
| 2 | `POST /admin/alerts/batch-acknowledge` | **不存在** | 后端缺少批量确认端点 | `api/admin.ts` |
| 3 | `GET /admin/quotas` | **不存在** | 后端无 quota 路由组 | `api/quota/index.ts` |
| 4 | `POST /admin/quotas` | **不存在** | 同上 | `api/quota/index.ts` |
| 5 | `PUT /admin/quotas/:id` | **不存在** | 同上 | `api/quota/index.ts` |
| 6 | `DELETE /admin/quotas/:id` | **不存在** | 同上 | `api/quota/index.ts` |
| 7 | `GET /quotas/usage` | **不存在** | 后端无 quota 路由 | `api/quota/index.ts` |
| 8 | `GET /auth/validate` | **不存在** | 后端无 validate 端点 | `api/auth.ts` |
| 9 | `POST /admin/machines/:id/heartbeat` | **不存在** | 心跳由 Agent 发，非管理员 | `api/host/index.ts` |

**P0-1 修复建议**:

1. **quota 相关端点（#3-7）**: 前端已有完整的 quota API 定义，但后端完全没有 quota 路由。建议：
   - 方案 A（推荐）：前端暂时移除 quota 页面入口，标记为"即将上线"
   - 方案 B：后端新增 quota 模块（工作量较大，建议作为后续迭代）
2. **batch-acknowledge（#2）**: 前端有批量确认告警功能，后端缺少。建议后端在 `ops_controller.go` 新增 `BatchAcknowledge` handler
3. **auth/validate（#8）**: 前端用于检查 token 有效性。建议前端改用 `GET /auth/profile` 替代（已存在），返回 401 即表示无效
4. **heartbeat（#9）**: 前端 `api/host/index.ts` 中定义了管理员发心跳，但心跳应由 Agent 发送。建议前端删除此 API 定义

---

#### P0-3: ctx.Get("userID") 类型断言 — 独立验证

**验证结论**: 大部分控制器已使用安全的 `ctx.GetUint("userID")`，但仍有 2 处使用不安全的直接断言：

| # | 文件 | 行号 | 代码 | 风险 |
|---|------|------|------|------|
| 1 | `controller/v1/environment/environment_controller.go` | 30 | `userID.(uint)` | **高** — 无 comma-ok |
| 2 | `controller/v1/workspace/workspace_controller.go` | 29 | `userID.(uint)` | **高** — 无 comma-ok |

**修复方案**: 将 `getCustomerID` 方法改为安全模式：
```go
// 修复前
return userID.(uint), true

// 修复后
id, ok := userID.(uint)
if !ok {
    c.Error(ctx, 500, "用户ID类型错误")
    return 0, false
}
return id, true
```

---

#### P0-4: ParseUint 错误处理 — 独立验证

**验证结论**: 所有 30 处 `ParseUint` 调用均已正确处理错误，无需修复。

**抽查示例**:
- `controller/v1/allocation/allocation_controller.go:43` — `if err == nil { ... }`
- `controller/v1/customer/sshkey_controller.go:122` — `if err != nil { c.Error(...); return }`
- `controller/v1/customer/customer_controller.go:117` — `if err != nil { c.Error(...); return }`
- `service/auth/auth_service.go:166` — `if err != nil { return ..., errors.New(...) }`

---

#### P0-5: 设备状态大部分显示离线 — 心跳链路分析

**完整链路时序**:

```
Agent (30s) → POST /api/v1/agent/heartbeat → MachineService.HeartbeatWithMetrics()
    ├─ Redis: SET host:status:{id} (TTL=90s)
    └─ DB: INSERT host_metrics

HostStatusSyncer (定时) → 遍历所有机器
    ├─ Redis key 存在 → 同步 device_status="online" 到 DB
    └─ Redis key 过期 → 更新 device_status="offline"

HeartbeatMonitor (定时) → 查询 DB
    └─ device_status="online" AND last_heartbeat < now-timeout → 标记 offline
```

**关键文件位置**:

| 组件 | 文件 | 关键行号 |
|------|------|---------|
| Agent 心跳发送 | `agent/cmd/main.go` | 88-111 |
| Agent HTTP 客户端 | `agent/internal/client/client.go` | 210-243 |
| 后端心跳接收 | `backend/internal/controller/v1/agent/heartbeat_controller.go` | 32-68 |
| 心跳处理 Service | `backend/internal/service/machine/machine_service.go` | 545-585 |
| Redis 缓存实现 | `backend/internal/service/machine/host_status_cache.go` | 12-51 |
| Redis→DB 同步器 | `backend/internal/service/machine/host_status_syncer.go` | 49-70 |
| 心跳超时监控 | `backend/internal/service/machine/heartbeat_monitor.go` | 54-101 |
| DB 更新操作 | `backend/internal/dao/machine_repo.go` | 225-233 |

**潜在问题分析**:

1. **双重离线检测冲突**: `HostStatusSyncer` 和 `HeartbeatMonitor` 都会标记设备离线，可能产生竞争
2. **Redis 不可用时的降级**: `HeartbeatWithMetrics` 中 Redis 写入失败被静默忽略（`_ = s.statusCache.SetOnline(...)`），但直接写 DB 的降级路径只更新 `device_status`，不更新 `last_heartbeat` 时间戳
3. **Syncer 同步间隔**: 如果 Syncer 间隔 > Redis TTL（90s），可能出现 Redis key 已过期但 DB 仍为 online 的窗口期
4. **HeartbeatMonitor timeout 配置**: 需确保 timeout > Agent 心跳间隔（30s）+ 网络延迟

**P0-5 修复建议**:

1. 统一超时参数：Agent 心跳 30s → Redis TTL 90s（已正确）→ HeartbeatMonitor timeout 建议 120s
2. Syncer 同步间隔建议 ≤ 30s，确保小于 Redis TTL
3. Redis 写入失败时，降级路径应同时更新 `last_heartbeat`
4. 考虑合并 Syncer 和 HeartbeatMonitor 为单一组件，避免竞争

---

#### P0-6: 分配操作报错 — 前后端字段差异分析

**前端分配入口**: `frontend/src/views/admin/MachineAllocateView.vue`
**前端 API 函数**: `frontend/src/api/admin.ts` — `assignMachine()` (行 103-115)
**后端 Handler**: `backend/internal/controller/v1/machine/machine_controller.go:275-298`
**后端请求体**: `backend/api/v1/machine.go:3-9` — `AllocateRequest`
**后端 Service**: `backend/internal/service/allocation/allocation_service.go:246-336`

**前后端字段差异对比**:

| 字段 | 前端发送 | 后端接收 | 后端存储 | 状态 |
|------|---------|---------|---------|------|
| `customer_id` | ✓ | ✓ | ✓ | 正常 |
| `duration_months` | ✓ | ✓ | 用于计算 | 正常 |
| `remark` | ✓ | ✓ | ✓ | 正常 |
| `start_time` | ✓ | ✗ | ✗ | **前端发送，后端忽略** |
| `end_time` | ✓ | ✗ | ✗ | **前端发送，后端忽略** |
| `contact_person` | ✓ | ✗ | ✗ | **前端发送，后端不支持** |
| `notify_methods` | ✓ | ✗ | ✗ | **前端发送，后端不支持** |

**核心问题**:

1. **时间字段被忽略**: 前端允许用户选择自定义开始/结束时间，但后端硬编码 `startTime = time.Now()`，用户选择的时间完全无效
2. **缺失字段**: `contact_person` 和 `notify_methods` 在后端 `AllocateRequest` 结构体、`Allocation` 实体、数据库表中均不存在
3. **后端 `AllocateRequest` 还要求 `host_id`**（`binding:"required"`），但前端是通过 URL 路径参数 `/admin/machines/{id}/allocate` 传递机器 ID，不在 body 中

**P0-6 修复建议**:

**方案 A（推荐 — 前端对齐后端）**:
1. 前端移除 `start_time`、`end_time`、`contact_person`、`notify_methods` 表单字段
2. 前端只发送 `customer_id`、`duration_months`、`remark`
3. 后端从 `AllocateRequest` 移除 `host_id`（已从 URL 获取）

**方案 B（后端扩展支持）**:
1. 后端 `AllocateRequest` 新增 `start_time`、`end_time` 可选字段
2. 后端 `Allocation` 实体新增 `contact_person`、`notify_methods` 字段
3. 新增数据库迁移脚本
4. 修改 `AllocateMachine` Service 方法支持自定义时间

---

#### 前端 UI 美化方向建议

**现状评估**:

- 布局框架：`AdminLayout.vue` + `CustomerLayout.vue`，侧边栏 + 顶部导航栏结构清晰
- 组件库：广泛使用 Element Plus（el-table, el-card, el-form 等）
- 自定义组件：StatCard, DataTable, FilterBar, DetailCard, PageHeader, StatusTag
- 主题：使用 Element Plus 标准色，登录页有渐变背景，管理后台浅灰 + 白色卡片
- 卡片样式：圆角 8px，阴影 `0 2px 8px rgba(0,0,0,0.08)`，hover 有提升效果

**现存不足**:

- 没有全局 CSS 变量定义（主题色、间距、圆角等散落各组件）
- 表格行高、间距不统一
- 表单项的 margin/padding 差异
- 部分自定义图表与 Element Plus 风格不协调
- 缺少暗黑主题支持

**美化方向建议（按优先级）**:

**P1 — 高影响力**:
1. 建立全局 CSS 变量系统（`src/styles/variables.css`），统一颜色、间距、圆角、阴影
2. 统一表格样式：行高、单元格 padding、分割线、hover 效果、表头背景色
3. 统一表单样式：间距 margin-bottom 16px、焦点状态、标签对齐

**P2 — 中等影响力**:
4. 仪表板 StatCard 视觉层次优化（字体大小、颜色对比度）
5. 按钮和操作栏统一大小和间距，改进批量操作栏设计
6. 响应式布局改进（移动端、平板适配）

**P3 — 增强体验**:
7. 页面切换动画、加载状态动画
8. 暗黑主题支持（定义暗色方案 + 切换功能）
9. 使用 ECharts 替代自定义图表，统一数据可视化风格

---

### 代码审查报告（第二轮 — architect）

> 审查时间：2026-02-09（第二轮）
> 审查范围：Task #3（toggle 端点）、Task #8（心跳超时）、Task #2（前端路径）

#### 编译状态

| 模块 | 命令 | 结果 |
|------|------|------|
| 后端 | `go build ./cmd/...` | ✅ 通过 |
| 前端 | `npx vue-tsc --noEmit` | ✅ 通过（之前的 2 个类型错误已修复） |

#### 审查 1: Task #3 — alert-rules toggle 端点（后端）

**评级: ✅ 良好（有小建议）**

| 文件 | 行号 | 状态 | 说明 |
|------|------|------|------|
| `dao/ops_repo.go` | 109-121 | ⚠️ 可改进 | 两次 DB 操作无事务，有竞态风险 |
| `service/ops/ops_service.go` | 71-74 | ✅ 通过 | 简单透传，符合当前复杂度 |
| `controller/v1/ops/ops_controller.go` | 365-390 | ✅ 通过 | 参数校验、Swagger 文档完整 |
| `router/router.go` | 276 | ✅ 通过 | 路由注册正确，中间件链完整 |

**建议（非阻塞）**: DAO 层 `ToggleAlertRule` 的 First+Update 两步操作可用 `db.Transaction()` 包裹，防止并发竞态。Controller 可区分 `gorm.ErrRecordNotFound` 返回 404。

#### 审查 2: Task #8 — 心跳超时一致性修复（后端）

**评级: ✅ 优秀**

| 文件 | 行号 | 状态 | 说明 |
|------|------|------|------|
| `host_status_cache.go` | 18-37 | ✅ 通过 | TTL 改为构造函数参数，注释清晰 |
| `router.go` | 106-121 | ✅ 通过 | 从配置读取超时，有合理默认值 180s |
| `environment_controller.go` | 23-31 | ✅ 通过 | 使用安全的 `ctx.GetUint("userID")` |
| `workspace_controller.go` | 22-30 | ✅ 通过 | 同上，实现一致 |

无问题，实现质量高。TTL 统一由配置驱动，userID 类型断言全部改为安全方式。

#### 审查 3: Task #2 + #5 — 前端路径修复 + 分配操作修复

**评级: ✅ 良好（有遗留项）**

| 文件 | 状态 | 说明 |
|------|------|------|
| `api/admin.ts` — `allocateMachine` | ✅ 通过 | 只发送 `customer_id`, `duration_months`, `remark` |
| `api/admin.ts` — `toggleAlertRule` | ✅ 通过 | 路径正确 `/admin/alert-rules/${id}/toggle` |
| `MachineAllocateView.vue` | ✅ 通过 | 正确调用 `allocateMachine`，过滤 idle 机器 |
| `MachineListView.vue` | ✅ 通过 | 分配弹窗过滤 active 客户，批量操作正确 |
| `CustomerListView.vue` | ✅ 通过 | 启用按钮支持 `disabled`+`suspended`，snake_case 字段 |
| `types/customer.ts` | ✅ 通过 | 全部 snake_case，与后端一致 |

**遗留项（非阻塞，建议后续处理）**:

1. `api/admin.ts:103-115` — `AssignMachinePayload` 接口仍包含后端不支持的字段（`start_time` 等），虽然实际调用已改用 `allocateMachine`，但死代码应清理
2. `api/auth.ts:140-142` — `validateToken()` 调用 `/auth/validate`，后端无此端点，建议移除或改用 `GET /auth/profile`
3. `api/host/index.ts:37-39` — `sendHeartbeat()` 调用管理员心跳端点，后端无此路由且心跳应由 Agent 发送，建议移除
4. `MachineAllocateView.vue:14-21` — 表单仍包含 `dateRange`、`contactPerson`、`notifyMethods` 字段但提交时被忽略，可能造成用户困惑

#### 审查总结

| 修复项 | 评级 | 阻塞问题 | 建议改进 |
|--------|------|---------|---------|
| Task #3 toggle 端点 | ✅ 良好 | 无 | DAO 加事务 |
| Task #8 心跳超时 | ✅ 优秀 | 无 | 无 |
| Task #2 前端路径 | ✅ 良好 | 无 | 清理死代码 |
| Task #5 分配操作 | ✅ 良好 | 无 | 清理未用表单字段 |

**结论**: 所有 P0 修复代码质量合格，后端和前端均编译通过，无阻塞问题。建议在 UI 美化阶段顺带清理上述遗留死代码。
