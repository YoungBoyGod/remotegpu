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

### P1 - 健壮性增强（提升运维信心）

| # | 功能 | 说明 | 涉及文件 | 依赖 |
|---|------|------|---------|------|
| 5 | 结构化日志 | 替换 log.Printf 为 slog，支持 level/JSON 输出 | 全部文件 | P0-4 |
| 6 | Token 认证 | Agent→Server 请求携带 Token，Server 校验 | client.go, cmd/main.go | P0-4 |
| 7 | 进度上报 | 长任务定期上报进度到 /progress | client.go, scheduler.go | 无 |
| 8 | 核心单元测试 | scheduler/executor/queue/store 的关键路径测试 | 新建 *_test.go | 无 |

### P2 - 调度增强（提升调度能力）

| # | 功能 | 说明 | 涉及文件 | 依赖 |
|---|------|------|---------|------|
| 9 | 任务依赖检查 | DependsOn 中的任务全部完成后才可执行 | scheduler.go, store.go | 无 |
| 10 | 抢占式调度 | 高优先级任务抢占低优先级（优先级差>=3） | scheduler.go, executor.go | 无 |
| 11 | 命令白名单 | 可配置的命令白名单/黑名单 | 新建 security/validator.go, executor.go | P0-4 |

### P3 - 高级功能（按需实现）

| # | 功能 | 说明 | 涉及文件 | 依赖 |
|---|------|------|---------|------|
| 12 | 任务组编排 | 串行/并行/DAG 模式 | 新建 orchestrator/ | P2-9 |
| 13 | 任务模板 | 模板变量替换，从模板创建任务 | 新建 template/ | 无 |
| 14 | 资源限制 | cgroup 或 ulimit 限制 CPU/内存 | executor.go | 无 |
| 15 | 审计日志 | 关键操作记录审计事件 | 新建 audit/ | P1-5 |
| 16 | 沙箱隔离 | Docker/nsjail 隔离执行 | executor.go | P2-11 |

---

## 四、实现顺序与详细方案

### 第一批：P0 核心可靠性

#### P0-1: 失败重试机制

**目标**: 任务执行失败后，如果 retry_count < max_retries，延迟 retry_delay 秒后重新入队。

**修改文件**: `internal/scheduler/scheduler.go`

**方案**:
- 在 `runTask()` 末尾，任务状态为 failed 时检查重试条件
- 满足条件则 reset 状态为 pending，retry_count++，延迟后重新 Push 到队列
- 不满足条件则保持 failed 并上报

```go
// runTask() 末尾添加重试逻辑
if task.Status == models.TaskStatusFailed && task.RetryCount < task.MaxRetries {
    task.RetryCount++
    task.Status = models.TaskStatusPending
    task.Error = ""
    task.ExitCode = 0
    s.store.Save(task)

    delay := time.Duration(task.RetryDelay) * time.Second
    if delay == 0 {
        delay = 60 * time.Second
    }
    time.AfterFunc(delay, func() {
        s.queue.Push(task)
        log.Printf("task %s retry %d/%d after %v", task.ID, task.RetryCount, task.MaxRetries, delay)
    })
    return // 不上报 complete，等重试
}
```

---

#### P0-2: 离线结果同步

**目标**: Agent 与 Server 断连期间，已完成任务的结果暂存本地（synced=0），重连后批量上报。

**新建文件**: `internal/syncer/syncer.go`

**方案**:
- 新建 Syncer 组件，定期扫描 `store.ListUnsynced()` 获取未同步任务
- 对每个未同步的 completed/failed 任务调用 `client.ReportComplete()`
- 成功后调用 `store.MarkSynced(id)`
- 在 `scheduler.runTask()` 中，ReportComplete 失败时不 panic，任务保持 synced=0

```go
type Syncer struct {
    store    *store.SQLiteStore
    client   *client.ServerClient
    interval time.Duration
    stopCh   chan struct{}
}

func (s *Syncer) Start() {
    go s.syncLoop()
}

func (s *Syncer) syncLoop() {
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

#### P0-3: 输出大小限制

**目标**: 限制 stdout/stderr 捕获大小，防止内存溢出。

**修改文件**: `internal/executor/executor.go`

**方案**:
- 使用 `io.LimitReader` 包装 stdout/stderr pipe
- 默认限制 1MB，超出截断并追加 `\n...[truncated]`

---

#### P0-4: 配置文件支持

**目标**: 支持 YAML 配置文件，环境变量作为 fallback。

**新建文件**: `internal/config/config.go`

**方案**:
```go
type Config struct {
    Port       int    `yaml:"port"`
    DBPath     string `yaml:"db_path"`
    MaxWorkers int    `yaml:"max_workers"`

    Server struct {
        URL       string        `yaml:"url"`
        AgentID   string        `yaml:"agent_id"`
        MachineID string        `yaml:"machine_id"`
        Token     string        `yaml:"token"`
        Timeout   time.Duration `yaml:"timeout"`
    } `yaml:"server"`

    Poll struct {
        Interval  time.Duration `yaml:"interval"`
        BatchSize int           `yaml:"batch_size"`
    } `yaml:"poll"`

    Log struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"` // text / json
    } `yaml:"log"`

    Limits struct {
        MaxOutputSize int `yaml:"max_output_size"` // bytes
    } `yaml:"limits"`
}
```

配置文件路径优先级: `./agent.yaml` > `/etc/remotegpu-agent/agent.yaml` > 环境变量

---

### 第二批：P1 健壮性增强

#### P1-5: 结构化日志

**修改**: 全部文件中的 `log.Printf` → `slog.Info/Warn/Error`

**方案**: 使用 Go 1.21+ 标准库 `log/slog`，支持 JSON 和 text 两种输出格式，通过配置切换。

---

#### P1-6: Token 认证

**修改文件**: `internal/client/client.go`

**方案**:
- Config 增加 Token 字段
- 所有 HTTP 请求添加 `Authorization: Bearer <token>` header
- 封装 `doRequest()` 方法统一处理

---

#### P1-7: 进度上报

**新增**: `client.ReportProgress(taskID, attemptID, percent int, message string)`

**修改**: `scheduler.runTask()` 中，对长任务（timeout > 300s）启动进度上报 goroutine，定期读取 stdout 最后几行作为进度信息。

---

#### P1-8: 核心单元测试

**新建文件**:
- `internal/queue/manager_test.go` — Push/Pop 顺序、Remove、并发安全
- `internal/store/sqlite_test.go` — Save/Get/ListByStatus/ListUnsynced
- `internal/scheduler/scheduler_test.go` — Submit 本地 vs Server 任务、recover 逻辑
- `internal/executor/executor_test.go` — 正常执行、超时、取消

---

### 第三批：P2 调度增强

#### P2-9: 任务依赖检查

**修改文件**: `internal/scheduler/scheduler.go`

**方案**:
- `tryExecute()` 中 Pop 任务后，检查 `task.DependsOn` 列表
- 查询 Store 中依赖任务的状态，全部 completed 才执行
- 否则放回队列末尾（降低优先级或延迟重试）

---

#### P2-10: 抢占式调度

**修改文件**: `internal/scheduler/scheduler.go`, `internal/executor/executor.go`

**方案**:
- Submit 高优先级任务时，检查当前 running 任务的优先级
- 优先级差 >= 3 时触发抢占：Cancel 当前任务 → 标记 preempted → 重新入队
- 被抢占任务在高优先级任务完成后自动恢复

---

#### P2-11: 命令白名单

**新建文件**: `internal/security/validator.go`

**方案**:
- 配置文件定义 allowed_commands 和 blocked_patterns
- Executor.Execute() 前调用 Validate(command, args)
- 不在白名单或匹配黑名单则拒绝执行

---

### 第四批：P3 高级功能

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

### 新建文件

| 文件 | 批次 | 说明 |
|------|------|------|
| `internal/config/config.go` | P0 | YAML 配置加载 |
| `internal/syncer/syncer.go` | P0 | 离线结果同步 |
| `internal/security/validator.go` | P2 | 命令白名单校验 |
| `internal/queue/manager_test.go` | P1 | 队列单元测试 |
| `internal/store/sqlite_test.go` | P1 | 存储单元测试 |
| `internal/scheduler/scheduler_test.go` | P1 | 调度器单元测试 |
| `internal/executor/executor_test.go` | P1 | 执行器单元测试 |
| `agent.yaml.example` | P0 | 配置文件示例 |

### 修改文件

| 文件 | 批次 | 修改内容 |
|------|------|---------|
| `internal/scheduler/scheduler.go` | P0/P2 | 重试逻辑、依赖检查、抢占 |
| `internal/executor/executor.go` | P0/P2 | 输出限制、资源限制 |
| `internal/client/client.go` | P1 | Token 认证、进度上报 |
| `cmd/main.go` | P0/P1 | 配置加载、Syncer 初始化、slog |

---

## 六、验证测试清单

| 场景 | 预期结果 | 对应功能 |
|------|---------|---------|
| 任务执行失败，retry_count < max_retries | 延迟后自动重试 | P0-1 |
| 任务执行失败，retry_count >= max_retries | 标记 failed，上报 Server | P0-1 |
| Agent 断网期间任务完成 | 结果暂存 synced=0 | P0-2 |
| Agent 重连后 | 自动批量上报未同步结果 | P0-2 |
| 任务 stdout 超过 1MB | 截断并标记 truncated | P0-3 |
| 提供 agent.yaml 启动 | 正确读取配置 | P0-4 |
| 无配置文件，有环境变量 | fallback 到环境变量 | P0-4 |
| 请求携带 Token | Server 校验通过 | P1-6 |
| 依赖任务未完成 | 当前任务等待，不执行 | P2-9 |
| 依赖任务全部完成 | 当前任务正常执行 | P2-9 |
| 高优先级任务到达（差>=3） | 抢占当前任务 | P2-10 |
| 命令不在白名单 | 拒绝执行，返回错误 | P2-11 |
