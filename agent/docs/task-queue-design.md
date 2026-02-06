# 分布式任务队列系统设计文档

## 1. 概述

### 1.1 背景
RemoteGPU 平台需要在多台 GPU 机器上执行各种任务（Shell 脚本、Python 脚本等）。需要一个分布式任务队列系统来管理任务的分发、执行、监控和恢复。

### 1.2 目标
- 支持任务优先级调度
- 支持任务持久化和故障恢复
- 支持任务编排（依赖关系、串行/并行）
- 支持任务状态追踪和结果获取
- 支持分布式执行

### 1.3 架构角色
- **Server（后端）**: 任务调度中心，负责任务创建、分发、状态管理
- **Agent（客户端）**: 任务执行器，负责接收任务、执行任务、上报结果

---

## 2. 核心概念

### 2.1 任务（Task）
```
Task {
    id:           string      // 唯一标识 (UUID)
    name:         string      // 任务名称
    type:         string      // 任务类型: shell, python, script
    command:      string      // 执行命令
    args:         []string    // 命令参数
    workdir:      string      // 工作目录
    env:          map         // 环境变量
    timeout:      int         // 超时时间(秒)
    priority:     int         // 优先级 (1-10, 1最高)
    retry_count:  int         // 当前重试次数
    retry_delay:  int         // 重试间隔(秒)
    max_retries:  int         // 最大重试次数

    // 状态相关
    status:       string      // pending/assigned/running/completed/failed/cancelled/preempted/suspended
    exit_code:    int         // 退出码
    stdout:       string      // 标准输出
    stderr:       string      // 标准错误
    error:        string      // 错误信息

    // 时间戳
    created_at:   timestamp
    assigned_at:  timestamp
    started_at:   timestamp
    ended_at:     timestamp

    // 关联
    machine_id:   string      // 目标机器
    group_id:     string      // 任务组ID（用于编排）
    parent_id:    string      // 父任务ID（用于依赖）
    depends_on:   []string    // 依赖的任务ID列表

    // 调度与租约
    assigned_agent_id: string   // 领取任务的 Agent
    lease_expires_at:  timestamp // 租约到期时间
    attempt_id:        string   // 本次领取的 attempt ID
}
```

### 2.2 任务组（TaskGroup）
用于任务编排，支持串行、并行、DAG 执行模式。

```
TaskGroup {
    id:           string
    name:         string
    mode:         string      // serial(串行), parallel(并行), dag(有向无环图)
    tasks:        []Task
    status:       string
    created_at:   timestamp
    ended_at:     timestamp
}
```

### 2.3 任务状态流转
```
                    ┌─────────────┐
                    │   pending   │  (创建后等待调度)
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  assigned   │  (已被 Agent 领取)
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
          ┌────────│   running   │────────┐
          │        └─────────────┘        │
          │                               │
          ▼                               ▼
   ┌─────────────┐                 ┌─────────────┐
   │  completed  │                 │   failed    │
   └─────────────┘                 └──────┬──────┘
                                          │
                                          │ (如果有重试)
                                          ▼
                                   ┌─────────────┐
                                   │   pending   │
                                   └─────────────┘

   running ───────────────► preempted (被抢占) ─► pending
   running ───────────────► suspended (主动暂停) ─► pending
   assigned ──────────────► pending (租约过期/未续约)

   任何状态 ──────────────────────► cancelled (手动取消)
```

---

## 3. 存储设计

### 3.1 存储方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| 内存 | 快速、简单 | 重启丢失 | 开发测试 |
| SQLite | 轻量、持久化 | 单机、并发有限 | 单Agent |
| Redis | 快速、支持队列 | 需要额外服务 | 高性能场景 |
| PostgreSQL | 可靠、功能全 | 较重 | 生产环境 |
| 文件(JSON/YAML) | 简单、可读 | 性能差 | 小规模 |

### 3.2 推荐方案：双层存储

**Server 端（PostgreSQL）**:
- 任务主表：存储所有任务元数据
- 任务结果表：存储执行结果（stdout/stderr 可能很大）
- 任务组表：存储任务编排信息

**Agent 端（SQLite + 内存）**:
- SQLite：持久化本地任务队列，支持重启恢复
- 内存：运行时优先级队列，快速调度

### 3.3 数据库表设计

#### Server 端表结构

```sql
-- 任务表
CREATE TABLE tasks (
    id              UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL DEFAULT 'shell',
    command         TEXT NOT NULL,
    args            JSONB,
    workdir         VARCHAR(500),
    env             JSONB,
    timeout         INT DEFAULT 3600,
    priority        INT DEFAULT 5 CHECK (priority BETWEEN 1 AND 10),
    retry_count     INT DEFAULT 0,
    retry_delay     INT DEFAULT 60,
    max_retries     INT DEFAULT 3,

    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    exit_code       INT,
    error           TEXT,

    machine_id      UUID REFERENCES machines(id),
    group_id        UUID,
    parent_id       UUID REFERENCES tasks(id),

    created_at      TIMESTAMP DEFAULT NOW(),
    assigned_at     TIMESTAMP,
    started_at      TIMESTAMP,
    ended_at        TIMESTAMP,

    assigned_agent_id UUID,
    lease_expires_at  TIMESTAMP,
    attempt_id        UUID,

    created_by      UUID REFERENCES users(id)
);

-- 任务结果表（大文本分离）
CREATE TABLE task_results (
    task_id         UUID PRIMARY KEY REFERENCES tasks(id),
    stdout          TEXT,
    stderr          TEXT,
    artifacts       JSONB,  -- 产出文件列表
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 任务依赖表
CREATE TABLE task_dependencies (
    task_id         UUID REFERENCES tasks(id),
    depends_on_id   UUID REFERENCES tasks(id),
    PRIMARY KEY (task_id, depends_on_id)
);

-- 任务组表
CREATE TABLE task_groups (
    id              UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    mode            VARCHAR(20) NOT NULL DEFAULT 'serial',
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMP DEFAULT NOW(),
    ended_at        TIMESTAMP,
    created_by      UUID REFERENCES users(id)
);

-- 索引
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_machine ON tasks(machine_id);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_group ON tasks(group_id);
```

#### Agent 端表结构（SQLite）

```sql
-- 本地任务队列
CREATE TABLE local_tasks (
    id              TEXT PRIMARY KEY,
    command         TEXT NOT NULL,
    workdir         TEXT,
    env             TEXT,  -- JSON
    timeout         INTEGER DEFAULT 3600,
    priority        INTEGER DEFAULT 5,

    status          TEXT DEFAULT 'pending',
    attempt_id      TEXT,
    assigned_agent_id TEXT,
    lease_expires_at  TEXT,
    exit_code       INTEGER,
    stdout          TEXT,
    stderr          TEXT,
    error           TEXT,

    received_at     TEXT,
    assigned_at     TEXT,
    started_at      TEXT,
    ended_at        TEXT,

    synced          INTEGER DEFAULT 0  -- 是否已同步到Server
);

CREATE INDEX idx_local_tasks_status ON local_tasks(status);
CREATE INDEX idx_local_tasks_priority ON local_tasks(priority);
```

---

## 4. 通信协议设计

### 4.1 通信模式对比

| 模式 | 描述 | 优点 | 缺点 |
|------|------|------|------|
| Push | Server 主动推送任务到 Agent | 实时性好 | 需要 Agent 可达 |
| Pull | Agent 主动从 Server 拉取任务 | Agent 可在 NAT 后 | 有延迟 |
| 混合 | Push + Pull 结合 | 兼顾两者优点 | 复杂度高 |

### 4.2 推荐方案：纯 Pull 模式

```
┌─────────────────────────────────────────────────────────────┐
│                         Server                               │
│  ┌─────────┐  ┌─────────┐  ┌─────────────────┐             │
│  │ Task DB │  │ Queue   │  │ REST API        │             │
│  └────┬────┘  └────┬────┘  └────────┬────────┘             │
│       │            │                │                      │
└───────┼────────────┼────────────────┼──────────────────────┘
        │            │                │
┌───────┼────────────┼────────────────┼──────────────────────┐
│       │            │                │                Agent │
│  ┌────▼────┐  ┌────▼────┐  ┌────────▼────────┐              │
│  │ SQLite  │  │ Memory  │  │ HTTP Client    │              │
│  │ (持久化)│  │ Queue   │  │ (拉取/上报)    │              │
│  └─────────┘  └─────────┘  └────────────────┘              │
└────────────────────────────────────────────────────────────┘
```

### 4.3 通信流程

1. **Agent 启动**: 启动轮询并注册 Agent 信息（可选）
2. **Agent 认领任务**: 周期性调用 REST API 认领任务（有副作用）
3. **Agent 执行任务**: 本地执行，定期上报进度
4. **Agent 完成任务**: 上报结果到 Server

---

## 5. API 设计

### 5.1 Server 端 API

#### 任务管理
```
POST   /api/v1/tasks              # 创建任务
GET    /api/v1/tasks              # 任务列表
GET    /api/v1/tasks/:id          # 任务详情
PUT    /api/v1/tasks/:id          # 更新任务
DELETE /api/v1/tasks/:id          # 删除任务
POST   /api/v1/tasks/:id/cancel   # 取消任务
POST   /api/v1/tasks/:id/retry    # 重试任务
GET    /api/v1/tasks/:id/result   # 获取结果
```

#### 任务组管理
```
POST   /api/v1/task-groups        # 创建任务组
GET    /api/v1/task-groups/:id    # 任务组详情
POST   /api/v1/task-groups/:id/start  # 启动任务组
```

#### Agent 专用 API
```
POST   /api/v1/agent/tasks/claim          # 认领任务（有副作用）
POST   /api/v1/agent/tasks/:id/start      # 标记任务开始
POST   /api/v1/agent/tasks/:id/lease/renew # 租约续约（心跳）
POST   /api/v1/agent/tasks/:id/progress   # 上报进度
POST   /api/v1/agent/tasks/:id/complete   # 上报完成
```

**幂等与错误码约定（Agent 专用 API）**

- 幂等键：`task_id + attempt_id`
- `/claim`:
  - 幂等键：`agent_id + request_id`（建议服务端生成并回传）。
  - 对相同请求重复调用应返回相同任务集合。
- `/start`:
  - 同一 `attempt_id` 重复调用应返回成功（幂等）。
  - 不同 `attempt_id` 视为冲突，返回 409。
- `/lease/renew`:
  - 同一 `attempt_id` 续约应返回新的 `lease_expires_at`。
  - 已过期续约返回 410。
- `/progress`:
  - 旧 `attempt_id` 的进度应被忽略或返回 409。
- `/complete`:
  - 同一 `attempt_id` 重复上报应返回成功（幂等）。
  - 不同 `attempt_id` 视为冲突，返回 409。

错误码约定（业务错误码与 HTTP 状态码分离）：

| 业务码 | 含义 | HTTP 状态 |
|-------|------|-----------|
| 30001 | attempt 不匹配 | 409 |
| 30002 | 任务已完成/不可变更 | 409 |
| 30003 | 租约已过期 | 410 |
| 30004 | 任务不存在 | 404 |
| 30005 | 参数无效 | 400 |
| 30099 | 内部错误 | 500 |

错误码定义位置：`/home/luo/code/remotegpu/agent/internal/errors/errors.go`

错误响应示例：
```json
{
  "code": 30001,
  "message": "attempt mismatch",
  "data": null
}
```

### 5.2 Agent 端 API

```
POST   /api/v1/tasks              # 本地创建任务（调试/本地执行）
GET    /api/v1/tasks              # 本地任务列表
GET    /api/v1/tasks/:id          # 任务详情
POST   /api/v1/tasks/:id/cancel   # 取消任务
GET    /api/v1/queue/status       # 队列状态
```

---

## 6. 故障恢复设计

### 6.1 Agent 重启恢复

```
Agent 启动流程:
1. 加载 SQLite 中的本地任务
2. 检查 running/assigned 状态的任务 → 向 Server 校验 attempt_id/lease
3. 若租约过期 → 停止进程并标记 failed；未过期则继续执行并续约
4. 检查 pending 状态的任务 → 重新加入内存队列
5. 连接 Server，同步状态
6. 开始正常调度
```

### 6.2 网络故障处理

```
场景1: Agent 与 Server 断开连接
- Agent 继续执行本地队列中的任务
- 结果暂存 SQLite，标记 synced=0
- 重连后批量同步未上报的结果

场景2: 任务执行超时
- 强制终止进程
- 标记任务为 failed
- 根据重试策略决定是否重试

场景3: Server 重启
- Agent 定期心跳检测
- 断开后自动重连
- 重连后同步本地状态
```

---

## 6.3 长时间运行任务与 Sidecar 服务

### 6.3.1 Sidecar（心跳/监控类）处理原则

- 这类任务属于“长期常驻服务”，不进入任务队列。
- 由 Agent 作为服务管理器负责启动与保活（systemd 或 docker）。
- 失败自动重启（`restart=always`），状态通过心跳上报给 Server。
- 仅上报运行状态与最后心跳时间，不参与任务重试/调度逻辑。

建议抽象为 `ManagedService`：
```
service_id, name, version, desired_state, current_state, last_heartbeat, restart_policy
```

### 6.3.2 长时间运行任务处理原则

1. **领取时分配租约**：claim 返回 `attempt_id` 与 `lease_expires_at`。
2. **执行期定期续约**：调用 `/lease/renew`，确保任务不会被重复领取。
3. **进度与心跳**：`/progress` 更新 `last_heartbeat_at` 与进度字段。
4. **取消处理**：Server 标记取消 → Agent SIGTERM → grace_period → SIGKILL。
5. **重启恢复**：Agent 重启后先校验租约与 attempt_id，过期则停止并上报失败。

推荐参数：
```
lease_ttl_sec = 300
renew_interval_sec = 60
grace_period_sec = 30
```

### 6.3.3 Agent 本地队列 YAML 示例

```yaml
local_queue:
  agent_id: "agent-001"
  machine_id: "machine-001"
  updated_at: "2026-02-05T12:05:00Z"

  tasks:
    - id: "task-train-001"
      name: "pytorch-train"
      type: "python"
      command: "python"
      args:
        - "train.py"
        - "--epochs=100"
      workdir: "/workspace"
      env:
        CUDA_VISIBLE_DEVICES: "0"
      timeout: 86400
      priority: 3
      max_retries: 3
      retry_delay: 60

      status: "assigned"
      attempt_id: "attempt-uuid-1"
      assigned_agent_id: "agent-001"
      assigned_at: "2026-02-05T12:00:00Z"
      lease_expires_at: "2026-02-05T12:05:00Z"

      started_at: "2026-02-05T12:00:05Z"
      last_heartbeat_at: "2026-02-05T12:04:55Z"
      exit_code: null
      error: null

    - id: "task-cleanup-002"
      name: "cleanup-workdir"
      type: "shell"
      command: "bash"
      args: ["-c", "rm -rf /workspace/tmp/*"]
      workdir: "/workspace"
      env: {}
      timeout: 600
      priority: 5

      status: "pending"
      attempt_id: ""
      assigned_agent_id: ""
      assigned_at: null
      lease_expires_at: null
      started_at: null
      last_heartbeat_at: null
      exit_code: null
      error: null
```

## 7. 任务编排设计

### 7.1 编排模式

#### MVP 用户输入（简化）

- 用户只需提供目标机器（ID 或 IP 列表）与任务基本信息。
- `repeat_count` 默认为 `1`，用于简单循环次数控制。

示例：
```yaml
task:
  name: "stability-test"
  type: "shell"
  command: "bash"
  args: ["-c", "./stress_test.sh"]
  timeout: 3600
  priority: 3
  repeat_count: 1

  targets:
    machine_ids: ["machine-001", "machine-002"]
    ips: []
```

#### 单任务配置示例
```yaml
task:
  id: "task-train-001"
  name: "pytorch-train"
  type: "python"
  command: "python"
  args:
    - "train.py"
    - "--epochs=100"
  workdir: "/workspace"
  env:
    CUDA_VISIBLE_DEVICES: "0"
  timeout: 86400
  priority: 3
  max_retries: 3
  retry_delay: 60

  machine_id: "machine-001"
  group_id: ""
  parent_id: ""
  depends_on: []

  lease_ttl_sec: 300
  renew_interval_sec: 60
  heartbeat_interval_sec: 30
  grace_period_sec: 30
```

#### 循环任务支持策略

**方案 A：任务内部循环（推荐）**
- 由脚本自身完成循环与重试，任务只启动一次。
- 优点：调度压力低、状态一致性好、租约续约更简单。

示例：
```yaml
task:
  id: "task-stability-001"
  name: "stability-test"
  type: "shell"
  command: "bash"
  args: ["-c", "for i in {1..1000}; do ./stress_test.sh || exit 1; done"]
  timeout: 86400
  priority: 3
```

**方案 B：模板 + 重复生成子任务**
- Server 端按模板生成 N 个子任务（可串行或并行）。
- 优点：每次执行可单独追踪与重试；可控制并发。
- 失败策略：超过 `max_failures` 后终止本次 repeat，按 `backoff_sec` 执行退避重试。
- 熔断策略：连续失败达到 `circuit_breaker_threshold` 时停止触发新任务并发出告警。
- 恢复策略：进入熔断后必须人工恢复（运维或管理员确认后解除熔断）。
  - 建议接口：`POST /api/v1/task-groups/:id/repeat/reset` 或 `POST /api/v1/tasks/repeat/reset`
  - 要求：记录操作者、原因与时间（审计日志）

示例：
```yaml
repeat:
  count: 100
  parallel: 5
  interval_sec: 10
  max_failures: 5
  backoff_sec: 30
  circuit_breaker_threshold: 3
  template:
    name: "stability-test"
    type: "shell"
    command: "bash"
    args: ["-c", "./stress_test.sh"]
    timeout: 1800
```

#### 串行执行 (Serial)
```yaml
group:
  mode: serial
  tasks:
    - name: "下载数据"
      command: "wget http://..."
    - name: "处理数据"
      command: "python process.py"
    - name: "上传结果"
      command: "aws s3 cp ..."
```

#### 并行执行 (Parallel)
```yaml
group:
  mode: parallel
  tasks:
    - name: "处理分片1"
      command: "python process.py --shard=1"
    - name: "处理分片2"
      command: "python process.py --shard=2"
    - name: "处理分片3"
      command: "python process.py --shard=3"
```

#### DAG 执行 (有向无环图)
```yaml
group:
  mode: dag
  tasks:
    - id: download
      name: "下载数据"
      command: "wget http://..."
    - id: process1
      name: "处理分片1"
      command: "python process.py --shard=1"
      depends_on: [download]
    - id: process2
      name: "处理分片2"
      command: "python process.py --shard=2"
      depends_on: [download]
    - id: merge
      name: "合并结果"
      command: "python merge.py"
      depends_on: [process1, process2]
```

---

## 8. 安全设计

### 8.1 认证授权
- Agent 使用 Token 认证连接 Server
- Token 在 Agent 注册时由 Server 生成
- 支持 Token 轮换和过期机制
- **Token 泄露防护**：每个 Token 绑定 Agent ID + IP，异常访问自动失效
- **动态 IP 兼容**：允许通过心跳更新 IP 或配置 IP 白名单范围，避免 NAT/弹性 IP 场景误杀

### 8.2 命令白名单
```yaml
# 可选：限制可执行的命令
allowed_commands:
  - "python"
  - "bash"
  - "nvidia-smi"
blocked_patterns:
  - "rm -rf /"
  - "dd if="
```

### 8.3 资源限制
- 单任务内存限制
- 单任务 CPU 限制
- 并发任务数限制
- 输出大小限制（防止日志爆炸）

#### 日志限制与保留

建议配置：
```yaml
log:
  max_size_mb: 100
  retain_count: 7
  local_dir: "/var/log/remotegpu-agent"
```

- `max_size_mb`：单文件最大大小，超过后滚动切分
- `retain_count`：保留最近 N 个日志文件
- `local_dir`：本地日志目录

#### Webhook 通知（用户可配置）

仅用于任务完成/失败提醒，内容最小化为任务 ID 与执行结果。

配置示例：
```yaml
webhook:
  url: "https://example.com/webhook"
```

回调示例：
```json
{
  "task_id": "task-001",
  "status": "completed"
}
```

说明：
- `url` 由用户填写，用于接收通知。

### 8.4 执行隔离

**避免 shell 拼接注入：**
```go
// 错误：直接拼接命令
cmd := exec.Command("bash", "-c", userCommand)  // 危险！

// 正确：参数化执行
cmd := exec.Command(command, args...)  // 安全
```

**可选沙箱隔离：**
```yaml
isolation:
  enabled: false          # 是否启用沙箱
  type: "none"            # none / docker / nsjail
  user: "nobody"          # 最小权限用户
  readonly_paths:         # 只读挂载
    - "/usr"
    - "/lib"
  writable_paths:         # 可写目录
    - "/tmp"
    - "/workspace"
```

---

## 9. 实现计划

### 9.1 阶段划分

#### 阶段一：基础任务队列（MVP）
- [ ] Agent 端优先级队列实现
- [ ] Agent 端 SQLite 持久化
- [ ] Agent 端任务执行器
- [ ] Agent 端基础 API
- [ ] 重启恢复功能

#### 阶段二：Server 集成
- [ ] Server 端任务表
- [ ] Server 端任务 API
- [ ] Agent 与 Server 通信
- [ ] 任务状态同步

#### 阶段三：高级功能
- [ ] 任务编排（串行/并行/DAG）
- [ ] 抢占式调度
- [ ] 安全白名单
- [ ] 资源限制

### 9.2 验证测试清单

#### 核心场景测试
- [ ] 任务重复领取：多 Agent 同时拉取，验证不重复
- [ ] 租约过期：Agent 崩溃后任务自动释放
- [ ] 断线重连：Agent 断网后恢复，任务状态正确
- [ ] 取消/超时：任务取消和超时处理正确
- [ ] 结果一致性：大文件上传 S3 后可正确获取
- [ ] 抢占恢复：被抢占任务正确重新入队

### 9.3 技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| Agent 持久化 | SQLite | 轻量、无需额外服务 |
| 优先级队列 | container/heap | Go 标准库 |
| HTTP 框架 | Gin | 已在使用 |
| HTTP 客户端 | net/http | Go 标准库 |

---

## 10. 详细设计讨论

### 10.1 任务分发策略

#### 设计决策：纯 Pull 模式

**选择原因：**
1. Agent 可能在 NAT/防火墙后，无法被 Server 直接访问
2. 简化实现，无需维护长连接
3. Server 可无状态水平扩展
4. 任务实时性要求不高，几秒延迟可接受
5. 避免 WebSocket 的复杂性

#### Pull 模式架构

```
┌──────────────────────────────────────────────────────┐
│                      Server                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐  │
│  │ Task DB │  │ Queue   │  │ REST API            │  │
│  └─────────┘  └─────────┘  └──────────┬──────────┘  │
└───────────────────────────────────────┼─────────────┘
                                        │
                              HTTP (每5秒轮询)
                                        │
┌───────────────────────────────────────┼─────────────┐
│                                       │      Agent  │
│  ┌─────────┐  ┌─────────┐  ┌─────────▼─────────┐   │
│  │ SQLite  │  │ Memory  │  │ HTTP Client       │   │
│  │ (持久化)│  │ Queue   │  │ (拉取/上报)       │   │
│  └─────────┘  └─────────┘  └───────────────────┘   │
└────────────────────────────────────────────────────┘
```

#### Pull 模式配置

```yaml
agent:
  server_url: "http://server:8080"
  poll_interval: 5s       # 轮询间隔
  batch_size: 10          # 每次拉取数量
  timeout: 30s            # 请求超时
  retry_interval: 10s     # 失败重试间隔
  max_retries: 3          # 最大重试次数
```

#### Pull 模式流程

```
Agent 启动
    │
    ▼
┌─────────────┐
│ 轮询 Server │◄─────────────┐
└──────┬──────┘              │
       │                     │
       ▼                     │
  有新任务？                  │
   │      │                  │
  是      否                 │
   │      │                  │
   ▼      └──等待5秒─────────┘
执行任务
   │
   ▼
上报结果
```

#### 认领 API

```
POST /api/v1/agent/tasks/claim

Request:
{
  "agent_id": "agent-001",
  "machine_id": "machine-xxx",
  "limit": 10,
  "request_id": "req-uuid-1"
}

Response:
{
  "tasks": [
    {
      "id": "task-001",
      "command": "python",
      "args": ["train.py"],
      "priority": 3,
      "status": "assigned",
      "assigned_agent_id": "agent-001",
      "assigned_at": "2026-02-05T12:00:00Z",
      "lease_expires_at": "2026-02-05T12:05:00Z",
      "attempt_id": "attempt-uuid-1"
    }
  ]
}
```

说明：
- `claim` 为有副作用操作，服务端应返回 `Cache-Control: no-store`，并按 `request_id` 提供幂等结果。

---

### 10.1.1 任务分配原子性保证

#### 问题描述

```
Agent A ──拉取任务──► Server ◄──拉取任务── Agent B
                      │
                同一任务被分配给两个 Agent？
```

#### 解决方案对比

| 方案 | 优点 | 缺点 |
|------|------|------|
| 数据库乐观锁 | 简单、无额外依赖 | 高并发时冲突多 |
| 数据库悲观锁 | 强一致性 | 性能差、可能死锁 |
| Redis 分布式锁 | 高性能 | 额外依赖 |
| 任务认领机制 | 灵活 | 两阶段，稍复杂 |

#### 推荐方案：数据库乐观锁

**原理：利用数据库原子更新**

```sql
-- Agent 拉取任务时，原子性地更新状态
UPDATE tasks
SET status = 'assigned',
    assigned_agent_id = 'agent-001',
    assigned_at = NOW(),
    lease_expires_at = NOW() + INTERVAL '5 minutes',
    attempt_id = gen_random_uuid()
WHERE id = (
    SELECT id FROM tasks
    WHERE status = 'pending'
    AND machine_id = 'xxx'
    ORDER BY priority, created_at
    LIMIT 1
)
AND status = 'pending'
RETURNING *;
```

**流程图：**

```
Agent A                Server                Agent B
   │                     │                     │
   ├──POST /tasks/claim──►│                     │
   │                     │◄─POST /tasks/claim───┤
   │                     │                     │
   │              ┌──────┴──────┐              │
   │              │ UPDATE ... │              │
   │              │ WHERE ...  │              │
   │              │ 原子操作    │              │
   │              └──────┬──────┘              │
   │                     │                     │
   │◄──返回 Task 1───────┤                     │
   │                     ├───返回 Task 2──────►│
   │                     │                     │
   │  (各自获得不同任务)  │                     │
```

**关键点：**
1. `UPDATE ... WHERE status = 'pending'` 保证原子性
2. 只有一个 Agent 能成功更新
3. 其他 Agent 获取下一个任务

#### 租约机制（Lease）

**问题：Agent 领取任务后崩溃，任务永远卡住**

**解决方案：引入租约过期机制**

```sql
-- 任务表增加租约字段
ALTER TABLE tasks ADD COLUMN assigned_agent_id UUID;
ALTER TABLE tasks ADD COLUMN assigned_at TIMESTAMP;
ALTER TABLE tasks ADD COLUMN lease_expires_at TIMESTAMP;
ALTER TABLE tasks ADD COLUMN attempt_id UUID;  -- 每次领取生成新ID
```

**租约流程：**

```
1. Agent 领取任务 → 设置 lease_expires_at = NOW() + 5分钟
2. Agent 执行中 → 定期续约（心跳 /lease/renew）
3. Agent 完成 → 清除租约，更新状态
4. Agent 崩溃 → 租约过期，任务可被重新领取
```

**幂等校验 API：**

```go
// POST /api/v1/agent/tasks/:id/start
type StartRequest struct {
    AgentID   string `json:"agent_id"`
    AttemptID string `json:"attempt_id"`  // 必须匹配
}

Response:
{
  "task_id": "task-001",
  "status": "running",
  "attempt_id": "attempt-uuid-1",
  "started_at": "2026-02-05T12:00:05Z"
}

// POST /api/v1/agent/tasks/:id/lease/renew
type LeaseRenewRequest struct {
    AgentID   string `json:"agent_id"`
    AttemptID string `json:"attempt_id"`
    ExtendSec int    `json:"extend_sec,omitempty"` // 默认 300s
}

Response:
{
  "task_id": "task-001",
  "status": "running",
  "attempt_id": "attempt-uuid-1",
  "lease_expires_at": "2026-02-05T12:10:05Z"
}

// POST /api/v1/agent/tasks/:id/complete
type CompleteRequest struct {
    AgentID   string `json:"agent_id"`
    AttemptID string `json:"attempt_id"`  // 必须匹配
    ExitCode  int    `json:"exit_code"`
    // ...
}

Response:
{
  "task_id": "task-001",
  "status": "completed",
  "attempt_id": "attempt-uuid-1",
  "ended_at": "2026-02-05T12:10:00Z"
}
```

---

### 10.2 结果存储策略

#### 文件大小分级

| 级别 | 大小范围 | 存储位置 | 说明 |
|------|----------|----------|------|
| 小 | < 64KB | 数据库 | 直接存 PostgreSQL |
| 大 | >= 64KB | S3/MinIO | 统一上传对象存储 |

#### 结果获取机制

**问题：Agent 离线时如何获取结果？**

**方案：统一上传到 Server/S3**

```
任务完成
    │
    ▼
结果大小判断
    │
    ├─ < 64KB → 直接返回给 Server（存DB）
    │
    └─ >= 64KB → 上传到 S3 → 返回 S3 key 给 Server
```

**结果获取 API：**

```
GET /api/v1/tasks/:id/result

Response:
{
  "storage_type": "s3",
  "presigned_url": "https://s3.../task-xxx?token=...",  // 预签名URL
  "expires_in": 3600
}
```

#### S3 配置

```yaml
storage:
  s3:
    endpoint: "s3.amazonaws.com"  # 或 MinIO 地址
    bucket: "remotegpu-results"
    region: "us-east-1"
    prefix: "task-results/"
    access_key: "${S3_ACCESS_KEY}"
    secret_key: "${S3_SECRET_KEY}"
```

#### 清理策略

```yaml
cleanup:
  local:
    max_age: 7d           # 本地文件保留7天
    max_size: 10GB        # 本地最大占用
  s3:
    max_age: 30d          # S3保留30天
    lifecycle_rule: true  # 使用S3生命周期规则
```

---

### 10.3 抢占式调度

#### 重要限制

**通用 shell/python 任务无法真正 checkpoint！**

| 任务类型 | 抢占方式 | 说明 |
|---------|---------|------|
| 协作式任务 | 暂停+恢复 | 任务主动支持 checkpoint |
| 通用任务 | 强杀+重试 | SIGTERM → 重新执行 |

**默认行为：强杀 + 重试**

```
抢占发生时：
1. 发送 SIGTERM 给当前任务
2. 等待 grace_period（默认30秒）
3. 若未退出，发送 SIGKILL
4. 当前任务标记为 preempted
5. 高优先级任务执行完成后，preempted 任务重新入队
```

#### 抢占规则

| 当前任务优先级 | 新任务优先级 | 是否抢占 |
|---------------|-------------|---------|
| 5 (普通) | 1 (紧急) | 是 |
| 3 (高) | 1 (紧急) | 是 |
| 1 (紧急) | 1 (紧急) | 否 |
| 任意 | 差值 < 3 | 否 |

#### 抢占流程

```
新高优先级任务到达
        │
        ▼
┌───────────────────┐
│ 检查是否满足抢占  │
│ 条件（优先级差≥3）│
└─────────┬─────────┘
          │
    ┌─────┴─────┐
    │           │
   是          否
    │           │
    ▼           ▼
┌─────────────┐  ┌─────────┐
│发送 SIGTERM │  │加入队列 │
│等待退出     │  │等待     │
└──────┬──────┘  └─────────┘
       │
       ▼
┌─────────────┐
│未退出则     │
│发送 SIGKILL │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│标记         │
│preempted    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│执行高优先级 │
│任务         │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│preempted    │
│任务重新入队 │
└─────────────┘
```

#### 任务状态扩展

```go
const (
    TaskStatusPreempted  = "preempted"  // 被抢占
    TaskStatusSuspended  = "suspended"  // 主动暂停
)

type PreemptedTask struct {
    TaskID     string
    Checkpoint string    // 检查点数据
    PreemptedAt time.Time
    PreemptedBy string   // 抢占者任务ID
}
```

---

### 10.4 多机器调度

#### 调度策略

| 策略 | 描述 | 适用场景 |
|------|------|----------|
| 指定机器 | 任务指定到特定机器 | 需要特定硬件 |
| 负载均衡 | 分配到最空闲机器 | 通用任务 |
| 标签匹配 | 按机器标签筛选 | GPU型号要求 |
| 广播 | 所有机器都执行 | 系统维护 |

#### 任务分发配置

```go
type TaskDistribution struct {
    Mode      string   `json:"mode"`      // single, broadcast, scatter
    MachineID string   `json:"machine_id,omitempty"`
    Labels    []string `json:"labels,omitempty"`
    Count     int      `json:"count,omitempty"` // scatter模式分发数量
}
```

#### 示例

```yaml
# 指定机器
distribution:
  mode: single
  machine_id: "uuid-xxx"

# 标签匹配
distribution:
  mode: single
  labels: ["gpu:a100", "region:us-east"]

# 广播到所有机器
distribution:
  mode: broadcast

# 分散到3台机器
distribution:
  mode: scatter
  count: 3
  labels: ["gpu:v100"]
```

---

### 10.5 任务模板

#### 模板结构

```go
type TaskTemplate struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Command     string            `json:"command"`
    Args        []string          `json:"args"`
    WorkDir     string            `json:"workdir"`
    Env         map[string]string `json:"env"`
    Timeout     int               `json:"timeout"`
    Priority    int               `json:"priority"`

    // 模板变量
    Variables   []TemplateVar     `json:"variables"`

    // 元数据
    Category    string            `json:"category"`
    Tags        []string          `json:"tags"`
    CreatedBy   string            `json:"created_by"`
    CreatedAt   time.Time         `json:"created_at"`
}
```

#### 模板变量

```go
type TemplateVar struct {
    Name        string `json:"name"`
    Type        string `json:"type"`    // string, int, bool, select
    Required    bool   `json:"required"`
    Default     string `json:"default"`
    Description string `json:"description"`
    Options     []string `json:"options,omitempty"` // select类型的选项
}
```

#### 模板示例

```yaml
# GPU 训练模板
id: "train-pytorch"
name: "PyTorch 训练任务"
description: "使用 PyTorch 进行模型训练"
command: "python"
args: ["train.py", "--model", "{{model}}", "--epochs", "{{epochs}}"]
workdir: "/workspace"
env:
  CUDA_VISIBLE_DEVICES: "{{gpu_id}}"
timeout: 86400
priority: 5
variables:
  - name: model
    type: select
    required: true
    options: ["resnet50", "vgg16", "bert"]
  - name: epochs
    type: int
    default: "100"
  - name: gpu_id
    type: string
    default: "0"
```

#### 模板 API

```
POST   /api/v1/templates           # 创建模板
GET    /api/v1/templates           # 模板列表
GET    /api/v1/templates/:id       # 模板详情
PUT    /api/v1/templates/:id       # 更新模板
DELETE /api/v1/templates/:id       # 删除模板
POST   /api/v1/templates/:id/run   # 从模板创建任务
```

---

### 10.6 审计日志

#### 日志结构

```go
type AuditLog struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Action    string    `json:"action"`
    Resource  string    `json:"resource"`
    ResourceID string   `json:"resource_id"`
    UserID    string    `json:"user_id"`
    MachineID string    `json:"machine_id,omitempty"`
    Details   string    `json:"details"`
    IP        string    `json:"ip"`
}
```

#### 审计事件类型

| 事件 | 描述 |
|------|------|
| task.created | 任务创建 |
| task.started | 任务开始执行 |
| task.completed | 任务完成 |
| task.failed | 任务失败 |
| task.cancelled | 任务取消 |
| task.preempted | 任务被抢占 |
| template.created | 模板创建 |
| template.updated | 模板更新 |
| template.deleted | 模板删除 |

#### 存储方案

```sql
CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY,
    timestamp   TIMESTAMP NOT NULL DEFAULT NOW(),
    action      VARCHAR(50) NOT NULL,
    resource    VARCHAR(50) NOT NULL,
    resource_id UUID,
    user_id     UUID,
    machine_id  UUID,
    details     JSONB,
    ip          VARCHAR(45)
);

CREATE INDEX idx_audit_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_resource ON audit_logs(resource, resource_id);
CREATE INDEX idx_audit_user ON audit_logs(user_id);
```

#### 查询 API

```
GET /api/v1/audit-logs?resource=task&action=created&from=2024-01-01&to=2024-01-31
```

#### 保留策略

```yaml
audit:
  retention: 90d      # 保留90天
  archive: true       # 过期后归档到S3
  compress: true      # 归档时压缩
```

---

## 11. 参考资料

- [Celery](https://docs.celeryq.dev/) - Python 分布式任务队列
- [Temporal](https://temporal.io/) - 工作流编排引擎
- [Asynq](https://github.com/hibiken/asynq) - Go 异步任务队列
- [Machinery](https://github.com/RichardKnop/machinery) - Go 异步任务队列
