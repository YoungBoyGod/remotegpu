# RemoteGPU 架构审查报告

> 审查日期：2026-02-07
> 审查范围：backend / frontend / agent 三个子模块

---

## 一、总体评价

项目整体架构清晰，分层合理，遵循了 Go 后端的标准分层模式（Controller → Service → DAO → Entity）和 Vue 3 前端的标准组织方式。Agent 模块设计完整，具备任务调度、优先级队列、离线同步等核心能力。

以下按模块逐一分析现存问题和改进建议。

---

## 二、Backend 架构审查

### 2.1 分层结构现状

```
backend/
├── api/v1/              # 请求/响应结构体
├── internal/
│   ├── controller/v1/   # HTTP 处理层
│   │   ├── common/      # BaseController
│   │   ├── agent/       # Agent 心跳
│   │   ├── auth/        # 认证
│   │   ├── customer/    # 客户相关（含 SSH Key、机器注册）
│   │   ├── dataset/     # 数据集
│   │   ├── document/    # 文档
│   │   ├── machine/     # 机器管理
│   │   ├── ops/         # 运维（仪表盘、监控、告警、审计、镜像、Agent）
│   │   ├── system_config/ # 系统配置
│   │   └── task/        # 任务（客户/管理员/Agent 三个 Controller）
│   ├── service/         # 业务逻辑层
│   ├── dao/             # 数据访问层
│   ├── model/entity/    # 实体定义
│   ├── middleware/       # 中间件
│   └── router/          # 路由注册
├── pkg/                 # 公共工具包
└── sql/                 # 迁移脚本
```

### 2.2 发现的问题

#### 问题 1：router.go 承担过多职责（高优先级）

`router/router.go` 的 `InitRouter` 函数同时负责：
- 所有 Service 的实例化和依赖注入
- 所有 Controller 的实例化
- 后台 Worker 的启动（`allocSvc.StartWorker`、`enrollmentSvc.StartWorker`、心跳监控、指标采集）
- 路由注册

**建议**：将依赖注入和 Worker 启动拆分到独立的 `wire.go` 或 `bootstrap.go` 中，`router.go` 只负责路由注册。

#### 问题 2：DAO 层 BaseDao 使用率低（中优先级）

已定义泛型 `BaseDao[T]`，但实际 DAO（`MachineDao`、`AllocationDao`、`TaskDao` 等）均未继承它，而是各自持有 `*gorm.DB` 并手写 CRUD。这导致大量重复代码。

**建议**：
- 对简单实体（如 `Document`、`Image`、`SSHKey`、`SystemConfig`）使用 `BaseDao[T]` 组合
- 对复杂实体保留自定义 DAO，但通过嵌入 `BaseDao[T]` 复用基础方法

#### 问题 3：ops 包职责过于宽泛（中优先级）

`service/ops/` 包含了 `OpsService`、`MonitorService`、`AgentService`、`DashboardService` 四个不同职责的服务。对应的 `controller/v1/ops/` 也包含了仪表盘、监控、告警、审计、镜像、Agent 六个 Controller。

**建议**：
- `AgentService` 应独立为 `service/agent/`
- `DashboardService` 应独立为 `service/dashboard/`
- `MonitorService` 可保留在 `service/ops/` 或独立为 `service/monitor/`

#### 问题 4：ID 类型不一致（中优先级）

- `Host` 实体使用 `string` 类型 ID（UUID）
- `Customer`、`Document` 等使用 `uint` 类型 ID（自增）
- `Allocation` 使用 `string` ID 但关联 `CustomerID uint`

**建议**：统一使用 UUID 作为主键类型，或至少在同一聚合根内保持一致。

#### 问题 5：缺少接口抽象（中优先级）

Service 层直接依赖具体 DAO 结构体，而非接口。这导致：
- 单元测试需要真实数据库连接
- 无法方便地 mock DAO 层

**建议**：为每个 DAO 定义接口，Service 依赖接口而非具体实现。

#### 问题 6：_wip 目录残留（低优先级）

`service/_wip/` 下有 6 个未完成的服务文件（dns、firewall、port_pool、vnc、docker、guacamole、rdp），这些文件不应出现在主分支。

**建议**：移至独立分支或删除，避免混淆。

### 2.3 安全相关

- Agent API 使用独立的 `middleware.AgentAuth()` 认证，与用户认证分离，设计合理
- 管理员路由统一使用 `middleware.RequireAdmin()` + `middleware.AuditMiddleware()`，审计覆盖完整
- `BaseController.Error` 直接返回 HTTP 200 + 业务错误码，需注意前端正确处理

---

## 三、Frontend 架构审查

### 3.1 结构现状

```
frontend/src/
├── api/           # API 请求（admin.ts + index.ts 为主）
├── components/    # 公共组件（layout/AdminLayout、AdminSidebar 等）
├── config/        # 表格列配置、表单配置、GPU 选项
├── router/        # 路由定义
├── stores/        # Pinia 状态管理
├── types/         # TypeScript 类型
└── views/         # 页面视图
    ├── admin/     # 管理员页面（~15 个视图）
    └── customer/  # 客户页面（~8 个视图）
```

### 3.2 发现的问题

#### 问题 1：API 模块文件被大量删除（高优先级）

从 git status 可见，`api/` 下按模块拆分的文件（artifact、billing、cmdb、image、issue、monitoring、notification、requirement、scheduler、storage、webhook）全部被删除，合并到了 `api/admin.ts` 和 `api/index.ts`。

**风险**：单文件过大，难以维护。

**建议**：保持按业务模块拆分 API 文件的方式，如 `api/machine.ts`、`api/task.ts`、`api/customer.ts` 等。

#### 问题 2：类型定义分散（中优先级）

`types/` 目录下有 `allocation.ts`、`common.ts` 等文件，但部分类型可能直接定义在组件或 API 文件中。

**建议**：确保所有共享类型集中在 `types/` 目录，API 响应类型与后端 `api/v1/` 结构体保持对应。

#### 问题 3：config 目录中硬编码业务数据（低优先级）

`config/hostSelection.ts` 中硬编码了 GPU 型号列表和地区列表。这些数据应从后端 API 动态获取。

**建议**：将 GPU 型号和地区改为从后端 `/api/v1/admin/settings/configs` 或专用接口获取。

#### 问题 4：客户端多个任务路由指向同一组件（低优先级）

`customer/tasks`、`customer/tasks/training`、`customer/tasks/inference`、`customer/tasks/queue`、`customer/tasks/history` 五个路由全部指向 `TaskListView.vue`，仅通过 meta.title 区分。

**建议**：通过路由参数或 query 传递任务类型过滤条件，在组件内根据路由自动切换过滤。

---

## 四、Agent 架构审查

### 4.1 结构现状

```
agent/
├── cmd/              # 入口 + HTTP 处理器
├── internal/
│   ├── client/       # Server 通信客户端
│   ├── config/       # 配置管理
│   ├── errors/       # 错误码
│   ├── executor/     # 任务执行器
│   ├── handler/      # 任务 HTTP 处理器
│   ├── models/       # 数据模型
│   ├── poller/       # 任务轮询器
│   ├── queue/        # 优先级队列
│   ├── scheduler/    # 任务调度器
│   ├── security/     # 命令校验
│   ├── store/        # SQLite 存储
│   └── syncer/       # 离线结果同步
└── docs/             # 设计文档
```

### 4.2 已实现的核心能力

| 能力 | 状态 | 说明 |
|------|------|------|
| 任务轮询与认领 | ✅ | Poller 定时 ClaimTasks |
| 优先级队列 | ✅ | 堆实现，支持优先级 + 时间排序 |
| 任务执行 | ✅ | 支持 shell/python/script，输出限制 1MB |
| 租约续约 | ✅ | 60 秒续约，30 秒进度上报 |
| 失败重试 | ✅ | 支持 MaxRetries + RetryDelay |
| 离线同步 | ✅ | Syncer 定期同步未上报结果 |
| 抢占式调度 | ✅ | 优先级差 ≥ 3 时触发 |
| 命令白名单 | ✅ | Security Validator |
| 心跳上报 | ✅ | 30 秒间隔 |

### 4.3 发现的问题

#### 问题 1：cmd/ 目录职责混乱（高优先级）

`cmd/` 下同时包含 `main.go`（入口）、`routes.go`（路由）、`handlers.go`（系统管理 HTTP 处理器）、`response.go`（响应格式）。系统管理 API（进程停止、SSH 重置、机器清理、命令执行）直接放在 cmd 包中，违反了分层原则。

**建议**：
- `handlers.go` 中的系统管理功能移至 `internal/handler/system.go`
- `routes.go` 移至 `internal/router/`
- `response.go` 移至 `internal/` 或 `pkg/`

#### 问题 2：Agent 本地 API 缺少认证（高优先级）

Agent 暴露了 `/api/v1/command/exec`（执行任意命令）、`/api/v1/process/stop`（停止进程）等高危端点，但从代码结构看缺少认证中间件。

**建议**：
- 添加 Token 认证中间件
- 限制监听地址为 `127.0.0.1` 或内网 IP
- 对 `/command/exec` 端点增加命令白名单校验

#### 问题 3：任务模型过于庞大（中优先级）

`models/task.go` 的 Task 结构体包含 30+ 字段，混合了任务定义、执行状态、调度信息、本地同步标记等多个关注点。

**建议**：考虑拆分为：
- `TaskSpec`：任务定义（Command、Args、Env、Timeout 等）
- `TaskStatus`：执行状态（Status、ExitCode、Stdout、Stderr 等）
- `TaskSchedule`：调度信息（Priority、AttemptID、LeaseExpiresAt 等）

#### 问题 4：缺少优雅关闭的完整实现（中优先级）

`main.go` 中启动了 Scheduler、Poller、Syncer 和 HTTP 服务器，但需确认所有组件在收到 SIGTERM 时能正确停止正在执行的任务并上报最终状态。

**建议**：添加 shutdown hook，确保：
1. 停止接受新任务
2. 等待正在执行的任务完成（或超时后强制终止）
3. 同步所有未上报的结果
4. 关闭 HTTP 服务器

#### 问题 5：测试覆盖不足（中优先级）

仅有 `executor_test.go`、`sqlite_test.go`、`manager_test.go`、`scheduler_test.go` 四个测试文件，`client/`、`poller/`、`syncer/` 等核心模块缺少测试。

**建议**：优先补充 `client/` 和 `scheduler/` 的单元测试。

---

## 五、跨模块问题

### 5.1 Agent ↔ Backend 通信协议未文档化（高优先级）

Agent 与 Backend 之间的通信协议散落在代码中，缺少统一的协议文档。虽然 `agent/docs/task-queue-design.md` 有详细设计，但与实际实现存在差异。

**建议**：创建独立的 `agent-protocol.md`，明确定义所有 API 端点、请求/响应格式、错误码、认证方式。（将在 Task #2 中完成）

### 5.2 错误码未统一（中优先级）

Backend 使用 `code` 字段返回业务错误码，但缺少统一的错误码定义文件。Agent 有独立的 `errors/errors.go`，两者未对齐。

**建议**：在 `backend/pkg/errors/` 中定义统一错误码常量，前端和 Agent 共享错误码定义。

### 5.3 缺少 API 版本管理策略（低优先级）

当前所有 API 在 `/api/v1/` 下，但缺少版本升级策略。

**建议**：在架构文档中明确 API 版本管理规则（何时升级 v2、如何兼容旧版本）。

---

## 六、改进优先级总结

| 优先级 | 问题 | 模块 |
|--------|------|------|
| P0 | Agent 本地 API 缺少认证 | Agent |
| P0 | router.go 承担过多职责 | Backend |
| P0 | Agent ↔ Backend 通信协议未文档化 | 跨模块 |
| P1 | API 模块文件合并导致单文件过大 | Frontend |
| P1 | cmd/ 目录职责混乱 | Agent |
| P1 | DAO 层 BaseDao 使用率低 | Backend |
| P1 | 缺少接口抽象 | Backend |
| P2 | ops 包职责过于宽泛 | Backend |
| P2 | ID 类型不一致 | Backend |
| P2 | 任务模型过于庞大 | Agent |
| P2 | 测试覆盖不足 | Agent |
| P2 | 错误码未统一 | 跨模块 |
| P3 | _wip 目录残留 | Backend |
| P3 | config 硬编码业务数据 | Frontend |
| P3 | 客户端多路由指向同一组件 | Frontend |
| P3 | 缺少 API 版本管理策略 | 跨模块 |
