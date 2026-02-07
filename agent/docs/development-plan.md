# RemoteGPU Agent 开发文档与功能清单

> 初始日期: 2026-02-06
> 最后更新: 2026-02-07
> 基于: task-queue-design.md + 现有代码审计 + 实现状态核实

---

## 一、架构概览

```
┌──────────────────────────────────────────────────────────┐
│                        Server (Go/Gin)                    │
│  PostgreSQL ← TaskDAO ← TaskService ← AgentTaskController │
│                                                           │
│  API: /api/v1/agent/tasks/{claim,start,lease/renew,complete} │
└──────────────────────────┬───────────────────────────────┘
                           │ HTTP (Pull, 每5秒轮询)
┌──────────────────────────┼───────────────────────────────┐
│                          │                        Agent   │
│  ┌────────┐  ┌─────────┐  ┌──────────┐  ┌───────────┐  │
│  │ SQLite │←→│Scheduler│←→│ Executor │  │ Poller    │  │
│  │ (Store)│  │         │  │          │  │ (Client)  │  │
│  └────────┘  └─────────┘  └──────────┘  └───────────┘  │
│       ↕            ↕            ↕                        │
│  ┌────────┐  ┌──────────┐ ┌──────────┐                  │
│  │ Syncer │  │ Priority │ │ Security │                  │
│  │        │  │  Queue   │ │Validator │                  │
│  └────────┘  └──────────┘ └──────────┘                  │
│                                                           │
│  本地 API: /api/v1/{tasks,queue/status,ping,system/info}  │
└──────────────────────────────────────────────────────────┘
```

### 核心数据流

```
Server pending 任务
    → Poller.ClaimTasks()
    → Scheduler.Submit() (保留 assigned 状态)
    → Queue.Push()
    → Scheduler.tryExecute()
    → runTask(): ReportStart → renewLoop → Executor.Execute → ReportComplete
    → Store.Save()
```

---

## 二、模块清单与实现状态

### 2.1 已完整实现的模块

| 模块 | 文件 | 功能 |
|------|------|------|
| 数据模型 | `internal/models/task.go` | Task 结构体、8 种状态常量、3 种任务类型 |
| 优先级队列 | `internal/queue/priority_queue.go` | container/heap 实现，按 priority+created_at 排序 |
| 队列管理器 | `internal/queue/manager.go` | 线程安全封装，Push/Pop/Remove/Get/NotifyChan |
| 任务执行器 | `internal/executor/executor.go` | 超时控制、环境变量、进程组、SIGTERM/SIGKILL 取消 |
| SQLite 持久化 | `internal/store/sqlite.go` | Save/Get/ListByStatus/Delete/ListUnsynced/MarkSynced |
| 任务调度器 | `internal/scheduler/scheduler.go` | 调度循环、Submit、recover、runTask、renewLoop |
| Server 客户端 | `internal/client/client.go` | ClaimTasks/ReportStart/RenewLease/ReportComplete |
| 轮询器 | `internal/poller/poller.go` | 定期拉取任务，回调提交到 Scheduler |
| 错误码 | `internal/errors/errors.go` | 30001-30099 业务错误码 |
| 任务 API | `internal/handler/task.go` | POST 创建 / GET 查询 / POST 取消 / GET 队列状态 |
| 系统管理 API | `cmd/handlers.go` | ping / system/info / process/stop / ssh/reset / cleanup / exec |
| 启动入口 | `cmd/main.go` | 环境变量配置、组件初始化、优雅关闭 |
| 路由注册 | `cmd/routes.go` | Gin 路由绑定 |
| 响应格式 | `cmd/response.go` | 统一 JSON 响应 |
| YAML 配置 | `internal/config/config.go` | YAML 文件 + 环境变量 fallback，支持 Server/Poll/Limits/Security 配置 |
| 离线结果同步 | `internal/syncer/syncer.go` | 定期扫描未同步任务，断网重连后批量上报 Server |
| 命令白名单 | `internal/security/validator.go` | 可配置的命令白名单和黑名单模式匹配 |

### 2.2 已实现但文档未记录的功能（本次核实更新）

| 功能 | 实现文件 | 代码现状 |
|------|---------|---------|
| 失败重试 | `scheduler.go:198-219` | 已实现：failed 后检查 max_retries，延迟 retry_delay 秒重新入队 |
| 离线结果同步 | `syncer/syncer.go` | 已实现：定期扫描 ListUnsynced，上报后 MarkSynced |
| 输出大小限制 | `executor.go:16-44` | 已实现：limitedWriter 限制 1MB，超出截断并标记 truncated |
| 配置文件支持 | `config/config.go` | 已实现：YAML 配置 + 环境变量 fallback |
| 结构化日志 | 全部文件 | 已实现：使用 log/slog，支持 Info/Warn/Error/Debug 级别 |
| Token 认证 | `client/client.go:51-61` | 已实现：doPost 统一添加 Authorization: Bearer header |
| 命令白名单 | `security/validator.go` | 已实现：白名单 + 黑名单模式匹配 |
| 任务依赖检查 | `scheduler.go:158-174` | 已实现：dependenciesMet 检查 DependsOn 列表 |
| 抢占式调度 | `scheduler.go:283-312` | 已实现：优先级差 >= 3 时触发抢占，被抢占任务延迟重新入队 |
| 进度上报 | `client/client.go:176-207` | 已实现：ReportProgress API，renewLoop 中定期上报 |
| 心跳上报 | `cmd/main.go:87-105` | 已实现：30 秒间隔定时心跳 |

### 2.3 尚未实现的功能

| 功能 | 说明 | 优先级 |
|------|------|--------|
| 任务组编排 | GroupID/ParentID 字段已定义，无 TaskGroup 实现 | P3 |
| 单元测试 | 无 *_test.go 文件 | P1 |
| 任务模板 | 无模板模块 | P3 |
| 审计日志 | 无审计模块 | P3 |
| 资源限制 | 无 cgroup/ulimit 集成 | P3 |
| 沙箱隔离 | 无 Docker/nsjail 集成 | P3 |

---

## 三、功能清单与优先级排序

### P0 - 核心可靠性 ✅ 全部已实现

| # | 功能 | 状态 | 实现文件 |
|---|------|------|---------|
| 1 | 失败重试机制 | ✅ 已实现 | `scheduler.go:198-219` |
| 2 | 离线结果同步 | ✅ 已实现 | `syncer/syncer.go` |
| 3 | 输出大小限制 | ✅ 已实现 | `executor.go:16-44` |
| 4 | 配置文件支持 | ✅ 已实现 | `config/config.go` |

### P1 - 健壮性增强（大部分已实现）

| # | 功能 | 状态 | 实现文件 |
|---|------|------|---------|
| 5 | 结构化日志 | ✅ 已实现 | 全部文件使用 `log/slog` |
| 6 | Token 认证 | ✅ 已实现 | `client/client.go:51-61` |
| 7 | 进度上报 | ✅ 已实现 | `client/client.go:176-207`, `scheduler.go:236-256` |
| 8 | 核心单元测试 | ❌ 未实现 | 无 *_test.go 文件 |

### P2 - 调度增强（大部分已实现）

| # | 功能 | 状态 | 实现文件 |
|---|------|------|---------|
| 9 | 任务依赖检查 | ✅ 已实现 | `scheduler.go:158-174` |
| 10 | 抢占式调度 | ✅ 已实现 | `scheduler.go:283-312` |
| 11 | 命令白名单 | ✅ 已实现 | `security/validator.go` |

### P3 - 高级功能（均未实现，按需开发）

| # | 功能 | 说明 | 涉及文件 | 依赖 |
|---|------|------|---------|------|
| 12 | 任务组编排 | 串行/并行/DAG 模式 | 新建 orchestrator/ | P2-9 |
| 13 | 任务模板 | 模板变量替换，从模板创建任务 | 新建 template/ | 无 |
| 14 | 资源限制 | cgroup 或 ulimit 限制 CPU/内存 | executor.go | 无 |
| 15 | 审计日志 | 关键操作记录审计事件 | 新建 audit/ | P1-5 |
| 16 | 沙箱隔离 | Docker/nsjail 隔离执行 | executor.go | P2-11 |

---

## 四、已实现功能的详细说明

### P0 核心可靠性（全部已实现）

#### P0-1: 失败重试机制 ✅

**实现文件**: `internal/scheduler/scheduler.go:198-219`

**实现方式**:
- 在 `runTask()` 末尾，任务状态为 failed 时检查 `MaxRetries > 0 && RetryCount < MaxRetries`
- 满足条件则 reset 状态为 pending，清除 AttemptID/Error/ExitCode/Stdout/Stderr
- 延迟 `RetryDelay` 秒后重新 Push 到队列（默认 60 秒）
- 不满足条件则保持 failed 并上报 Server

```go
// scheduler.go:198-219 实际实现
if task.Status == models.TaskStatusFailed && task.MaxRetries > 0 && task.RetryCount < task.MaxRetries {
    task.RetryCount++
    task.Status = models.TaskStatusPending
    task.AttemptID = ""
    task.Error = ""
    task.ExitCode = 0
    task.Stdout = ""
    task.Stderr = ""
    s.store.Save(task)

    delay := time.Duration(task.RetryDelay) * time.Second
    if delay <= 0 {
        delay = 60 * time.Second
    }
    slog.Info("task failed, scheduling retry", "task_id", task.ID,
        "retry", task.RetryCount, "max_retries", task.MaxRetries, "delay", delay)
    time.AfterFunc(delay, func() {
        s.queue.Push(task)
    })
    return
}
```

---

#### P0-2: 离线结果同步 ✅

**实现文件**: `internal/syncer/syncer.go`

**实现方式**:
- Syncer 组件定期（默认 30 秒）扫描 `store.ListUnsynced()` 获取未同步任务
- 对每个有 AttemptID 的 completed/failed 任务调用 `client.ReportComplete()`
- 成功后调用 `store.MarkSynced(id)`
- 启动时立即同步一次，之后按间隔定期同步
- 在 `cmd/main.go` 中初始化并集成到优雅关闭流程

```go
// syncer.go 实际实现
type Syncer struct {
    store    *store.SQLiteStore
    client   *client.ServerClient
    interval time.Duration
    stopCh   chan struct{}
}

func (s *Syncer) syncLoop() {
    s.syncUnreported() // 启动时立即同步一次
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()
    for {
        select {
        case <-s.stopCh:
            return
        case <-ticker.C:
            s.syncUnreported()
        }
    }
}
```

---

#### P0-3: 输出大小限制 ✅

**实现文件**: `internal/executor/executor.go:16-44`

**实现方式**:
- 自定义 `limitedWriter` 包装 stdout/stderr，限制 1MB（`maxOutputSize = 1 << 20`）
- 超出限制后静默丢弃（不中断进程），并在输出末尾追加 `\n...[truncated, output exceeded 1MB limit]`
- 通过 `cmd.Stdout = &limitedWriter{limit: maxOutputSize}` 直接替代 pipe

---

#### P0-4: 配置文件支持 ✅

**实现文件**: `internal/config/config.go`

**实现方式**:
- 支持 YAML 配置文件，路径优先级: `./agent.yaml` > `/etc/remotegpu-agent/agent.yaml`
- 环境变量作为 fallback 覆盖（AGENT_PORT, AGENT_DB_PATH, SERVER_URL 等）
- 配置结构包含 Server/Poll/Limits/Security 四个子配置
```go
// config.go 实际实现
type Config struct {
    Port       int    `yaml:"port"`
    DBPath     string `yaml:"db_path"`
    MaxWorkers int    `yaml:"max_workers"`

    Server   ServerConfig   `yaml:"server"`
    Poll     PollConfig     `yaml:"poll"`
    Limits   LimitsConfig   `yaml:"limits"`
    Security SecurityConfig `yaml:"security"`
}
```

配置文件路径优先级: `./agent.yaml` > `/etc/remotegpu-agent/agent.yaml` > 环境变量

---

### P1 健壮性增强（大部分已实现）

#### P1-5: 结构化日志 ✅

**实现方式**: 全部文件已使用 Go 标准库 `log/slog`，支持 Info/Warn/Error/Debug 级别和结构化键值对。

---

#### P1-6: Token 认证 ✅

**实现文件**: `internal/client/client.go:51-61`

**实现方式**:
- Config 包含 Token 字段，通过 YAML 或环境变量 `AGENT_TOKEN` 配置
- `doPost()` 方法统一为所有 HTTP 请求添加 `Authorization: Bearer <token>` header

---

#### P1-7: 进度上报 ✅

**实现文件**: `internal/client/client.go:176-207`, `internal/scheduler/scheduler.go:236-256`

**实现方式**:
- `client.ReportProgress(taskID, attemptID, percent, message)` API 已实现
- `renewLoop` 中每 30 秒自动上报进度（与租约续约 60 秒交替进行）

---

#### P1-8: 核心单元测试 ❌ 未实现

**待新建文件**:
- `internal/queue/manager_test.go`
- `internal/store/sqlite_test.go`
- `internal/scheduler/scheduler_test.go`
- `internal/executor/executor_test.go`

---

### P2 调度增强（全部已实现）

#### P2-9: 任务依赖检查 ✅

**实现文件**: `internal/scheduler/scheduler.go:158-174`

**实现方式**:
- `tryExecute()` 中 Pop 任务后调用 `dependenciesMet(task)`
- 遍历 `task.DependsOn` 列表，查询 Store 中依赖任务状态
- 全部 completed 才执行，否则延迟 2 秒后重新入队

---

#### P2-10: 抢占式调度 ✅

**实现文件**: `internal/scheduler/scheduler.go:283-312`

**实现方式**:
- `Submit()` 中调用 `tryPreempt(newTask)` 检查是否需要抢占
- 比较新任务与当前运行中最低优先级任务的优先级差值
- 差值 >= 3 时触发抢占：Cancel 当前任务 → 标记 preempted → 1 秒后重新入队

---

#### P2-11: 命令白名单 ✅

**实现文件**: `internal/security/validator.go`, `internal/executor/executor.go:87-96`

**实现方式**:
- `Validator` 支持 `allowedCommands`（白名单）和 `blockedPatterns`（黑名单）
- 在 `Executor.Execute()` 开头调用 `validator.Validate(command, args)`
- 不在白名单或匹配黑名单则拒绝执行，任务标记为 failed
- 通过 YAML 配置文件的 `security` 段配置

---

### P3 高级功能（均未实现，按需开发）

P3 功能按需实现，此处仅列出方向，不展开详细方案。

| # | 功能 | 方向 |
|---|------|------|
| 12 | 任务组编排 | 新建 orchestrator 包，实现 serial/parallel/dag 三种模式 |
| 13 | 任务模板 | Server 端功能为主，Agent 端仅需支持接收模板展开后的任务 |
| 14 | 资源限制 | Linux cgroup v2 或 `ulimit` 限制单任务资源 |
| 15 | 审计日志 | Server 端功能为主，Agent 端记录本地操作日志 |
| 16 | 沙箱隔离 | 可选 Docker 或 nsjail，通过配置开关 |

---

## 五、文件变更汇总

### 已实现的文件

| 文件 | 批次 | 说明 |
|------|------|------|
| `internal/config/config.go` | P0 ✅ | YAML 配置加载 + 环境变量 fallback |
| `internal/syncer/syncer.go` | P0 ✅ | 离线结果同步 |
| `internal/security/validator.go` | P2 ✅ | 命令白名单/黑名单校验 |
| `agent.yaml.example` | P0 ✅ | 配置文件示例 |

### 已修改的文件

| 文件 | 批次 | 修改内容 |
|------|------|---------|
| `internal/scheduler/scheduler.go` | P0/P2 ✅ | 重试逻辑、依赖检查、抢占调度 |
| `internal/executor/executor.go` | P0/P2 ✅ | 输出限制（limitedWriter）、命令白名单校验 |
| `internal/client/client.go` | P1 ✅ | Token 认证（doPost）、进度上报、心跳 |
| `cmd/main.go` | P0/P1 ✅ | 配置加载、Syncer/Poller 初始化、心跳、slog |

### 待新建的文件（未实现）

| 文件 | 批次 | 说明 |
|------|------|------|
| `internal/queue/manager_test.go` | P1 | 队列单元测试 |
| `internal/store/sqlite_test.go` | P1 | 存储单元测试 |
| `internal/scheduler/scheduler_test.go` | P1 | 调度器单元测试 |
| `internal/executor/executor_test.go` | P1 | 执行器单元测试 |

---

## 六、验证测试清单

| 场景 | 预期结果 | 对应功能 | 代码已实现 | 有单元测试 |
|------|---------|---------|:---------:|:---------:|
| 任务执行失败，retry_count < max_retries | 延迟后自动重试 | P0-1 | ✅ | ❌ |
| 任务执行失败，retry_count >= max_retries | 标记 failed，上报 Server | P0-1 | ✅ | ❌ |
| Agent 断网期间任务完成 | 结果暂存 synced=0 | P0-2 | ✅ | ❌ |
| Agent 重连后 | 自动批量上报未同步结果 | P0-2 | ✅ | ❌ |
| 任务 stdout 超过 1MB | 截断并标记 truncated | P0-3 | ✅ | ❌ |
| 提供 agent.yaml 启动 | 正确读取配置 | P0-4 | ✅ | ❌ |
| 无配置文件，有环境变量 | fallback 到环境变量 | P0-4 | ✅ | ❌ |
| 请求携带 Token | Server 校验通过 | P1-6 | ✅ | ❌ |
| 依赖任务未完成 | 当前任务等待，不执行 | P2-9 | ✅ | ❌ |
| 依赖任务全部完成 | 当前任务正常执行 | P2-9 | ✅ | ❌ |
| 高优先级任务到达（差>=3） | 抢占当前任务 | P2-10 | ✅ | ❌ |
| 命令不在白名单 | 拒绝执行，返回错误 | P2-11 | ✅ | ❌ |
| 心跳定时上报 | 30 秒间隔发送心跳 | 心跳 | ✅ | ❌ |
| Agent 重启后恢复任务 | pending 重新入队，running 标记 failed | recover | ✅ | ❌ |
