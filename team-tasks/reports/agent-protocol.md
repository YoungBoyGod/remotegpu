# Agent 与 Backend 通信协议规范

> 版本：v1.0
> 日期：2026-02-07
> 状态：基于现有代码整理，含改进建议

---

## 一、概述

Agent 部署在 GPU 机器上，通过 HTTP 轮询方式与 Backend 通信。所有 Agent API 位于 `/api/v1/agent/` 前缀下，使用独立的 Token 认证机制。

### 1.1 通信模式

- **方向**：Agent → Backend（Agent 主动发起请求）
- **协议**：HTTP/HTTPS
- **数据格式**：JSON
- **轮询间隔**：可配置，默认 5 秒

### 1.2 核心流程

```
Agent 启动
    │
    ├─ 心跳上报（每 30 秒）
    │
    ├─ 任务轮询（每 5 秒）
    │   └─ ClaimTasks → StartTask → RenewLease(循环) → CompleteTask
    │
    └─ 离线同步（定期同步未上报结果）
```

---

## 二、认证机制

### 2.1 当前实现

**Backend 端**（`middleware/agent_auth.go`）：
- 从请求头 `X-Agent-Token` 读取 Token
- 使用 `crypto/subtle.ConstantTimeCompare` 与配置中的 `agent.token` 比对
- 所有 `/api/v1/agent/*` 路由均受此中间件保护

**Agent 端**（`client/client.go`）：
- 使用 `Authorization: Bearer <token>` 发送 Token

### 2.2 已知问题：认证 Header 不一致

Backend 中间件检查 `X-Agent-Token`，但 Agent 客户端发送的是 `Authorization: Bearer <token>`。**两端使用了不同的 Header 名称，需要统一。**

**建议**：统一使用 `X-Agent-Token` Header，Agent 客户端 `doPost` 方法需修改：

```go
// 修改前
req.Header.Set("Authorization", "Bearer "+c.token)

// 修改后
req.Header.Set("X-Agent-Token", c.token)
```

### 2.3 规范定义

| 项目 | 值 |
|------|-----|
| Header 名称 | `X-Agent-Token` |
| Token 格式 | 纯字符串（不含 Bearer 前缀） |
| 配置来源 | Backend: `config.agent.token`；Agent: `agent.yaml` 中 `server.token` |
| 校验方式 | 常量时间比较（防时序攻击） |

---

## 三、统一响应格式

所有 API 响应均使用以下 JSON 格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 0 表示成功，非零为错误码 |
| msg | string | 成功时为 "success"，失败时为错误描述 |
| data | object/null | 业务数据，失败时为 null |

### 3.1 错误码定义

| HTTP 状态码 | 业务 code | 场景 |
|-------------|-----------|------|
| 200 | 0 | 成功 |
| 200 | 400 | 请求参数错误 |
| 200 | 409 | 冲突（任务状态不匹配、attempt_id 不匹配） |
| 200 | 410 | 租约已过期 |
| 200 | 500 | 服务器内部错误 |
| 401 | - | 认证失败（HTTP 层直接返回） |

---

## 四、API 端点定义

### 4.1 心跳上报

**用途**：Agent 定期上报存活状态，Backend 更新机器的 `last_heartbeat` 和 `device_status`。

```
POST /api/v1/agent/heartbeat
```

**请求体**：

```json
{
  "agent_id": "agent-uuid-001",
  "machine_id": "host-uuid-001"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| agent_id | string | 是 | Agent 唯一标识 |
| machine_id | string | 是 | 机器唯一标识 |

**成功响应**：

```json
{
  "code": 0,
  "msg": "success",
  "data": { "status": "ok" }
}
```

**Agent 行为**：
- 每 30 秒发送一次
- 失败时仅记录日志，不中断其他流程

**Backend 行为**：
- 更新 `hosts.last_heartbeat = NOW()`
- 更新 `hosts.device_status = "online"`

---

### 4.2 任务认领

**用途**：Agent 批量认领分配给本机器的待执行任务。

```
POST /api/v1/agent/tasks/claim
```

**请求体**：

```json
{
  "agent_id": "agent-uuid-001",
  "machine_id": "host-uuid-001",
  "limit": 10,
  "request_id": "req-uuid-001"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| agent_id | string | 是 | Agent 唯一标识 |
| machine_id | string | 是 | 机器唯一标识 |
| limit | int | 否 | 最大认领数量，默认 10 |
| request_id | string | 否 | 请求幂等 ID（Agent 端生成） |

**成功响应**：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "tasks": [
      {
        "id": "task-uuid-001",
        "name": "训练任务",
        "type": "shell",
        "command": "python",
        "args": ["train.py"],
        "workdir": "/workspace",
        "env": { "CUDA_VISIBLE_DEVICES": "0" },
        "timeout": 3600,
        "priority": 5,
        "max_retries": 3,
        "status": "assigned",
        "machine_id": "host-uuid-001",
        "assigned_agent_id": "agent-uuid-001",
        "attempt_id": "attempt-uuid-001",
        "lease_expires_at": "2026-02-07T10:05:00Z"
      }
    ]
  }
}
```

**Backend 行为**（事务内原子操作）：
1. 查询 `machine_id` 对应的 `status = 'pending'` 任务
2. 按 `priority ASC, created_at ASC` 排序，取 `limit` 条
3. 批量更新：`status → assigned`，设置 `assigned_agent_id`、`assigned_at`、`lease_expires_at`（+5分钟）、`attempt_id`（新生成 UUID）
4. 返回更新后的任务列表

**Agent 行为**：
- 每 5 秒轮询一次
- 收到任务后调用 `scheduler.Submit()` 入队
- 无任务时返回空数组，不视为错误

---

### 4.3 任务开始

**用途**：Agent 通知 Backend 任务已开始执行。

```
POST /api/v1/agent/tasks/:id/start
```

**路径参数**：

| 参数 | 说明 |
|------|------|
| id | 任务 ID |

**请求体**：

```json
{
  "agent_id": "agent-uuid-001",
  "attempt_id": "attempt-uuid-001"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| agent_id | string | 是 | Agent 唯一标识 |
| attempt_id | string | 是 | 认领时获得的 attempt_id，用于幂等校验 |

**成功响应**：

```json
{
  "code": 0,
  "msg": "success",
  "data": { "task_id": "task-uuid-001", "status": "running" }
}
```

**错误响应**（attempt_id 不匹配）：

```json
{
  "code": 409,
  "msg": "record not found",
  "data": null
}
```

**Backend 行为**：
- 校验 `id + agent_id + attempt_id` 三者匹配
- 更新 `status → running`，设置 `started_at = NOW()`
- 若 `RowsAffected == 0`，返回 409

---

### 4.4 租约续约

**用途**：Agent 在任务执行期间定期续约，防止租约过期后任务被重新分配。

```
POST /api/v1/agent/tasks/:id/lease/renew
```

**路径参数**：

| 参数 | 说明 |
|------|------|
| id | 任务 ID |

**请求体**：

```json
{
  "agent_id": "agent-uuid-001",
  "attempt_id": "attempt-uuid-001",
  "extend_sec": 300
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| agent_id | string | 是 | Agent 唯一标识 |
| attempt_id | string | 是 | 认领时获得的 attempt_id |
| extend_sec | int | 否 | 续约秒数，默认 300（5分钟） |

**成功响应**：

```json
{
  "code": 0,
  "msg": "success",
  "data": { "task_id": "task-uuid-001", "renewed": true }
}
```

**Backend 行为**：
- 校验 `id + agent_id + attempt_id + status='running'` 四者匹配
- 更新 `lease_expires_at = NOW() + extend_sec`
- 若 `RowsAffected == 0`，返回 410（租约已过期或任务不存在）

**Agent 行为**：
- 每 60 秒续约一次
- 续约失败时记录日志，不中断任务执行

---

### 4.5 任务完成

**用途**：Agent 上报任务执行结果。

```
POST /api/v1/agent/tasks/:id/complete
```

**路径参数**：

| 参数 | 说明 |
|------|------|
| id | 任务 ID |

**请求体**：

```json
{
  "agent_id": "agent-uuid-001",
  "attempt_id": "attempt-uuid-001",
  "exit_code": 0,
  "stdout": "训练完成，准确率 95.2%",
  "stderr": "",
  "error": ""
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| agent_id | string | 是 | Agent 唯一标识 |
| attempt_id | string | 是 | 认领时获得的 attempt_id |
| exit_code | int | 是 | 进程退出码，0 为成功 |
| stdout | string | 否 | 标准输出（限制 1MB） |
| stderr | string | 否 | 标准错误（限制 1MB） |
| error | string | 否 | Agent 层面的错误信息 |

**成功响应**：

```json
{
  "code": 0,
  "msg": "success",
  "data": { "task_id": "task-uuid-001", "status": "completed" }
}
```

**Backend 行为**：
- 校验 `id + agent_id + attempt_id` 三者匹配
- 根据 `exit_code` 决定最终状态：`exit_code == 0` → `completed`，否则 → `failed`
- 更新 `exit_code`、`error_msg`、`ended_at = NOW()`
- 若 `RowsAffected == 0`，返回 409

**Agent 行为**：
- 任务执行完毕后立即调用
- 若上报失败，保存到本地 SQLite，标记 `synced = false`
- Syncer 定期重试未同步的结果

---

### 4.6 进度上报（待实现）

**用途**：Agent 在任务执行期间上报进度百分比和消息。

> 注意：Agent 客户端已实现 `ReportProgress` 方法，但 Backend 尚未注册对应路由。

```
POST /api/v1/agent/tasks/:id/progress
```

**请求体**：

```json
{
  "agent_id": "agent-uuid-001",
  "attempt_id": "attempt-uuid-001",
  "percent": 45,
  "message": "正在训练第 45/100 个 epoch"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| agent_id | string | 是 | Agent 唯一标识 |
| attempt_id | string | 是 | 认领时获得的 attempt_id |
| percent | int | 否 | 进度百分比（0-100） |
| message | string | 否 | 进度描述 |

**Agent 行为**：
- 每 30 秒上报一次（与租约续约交替）

---

## 五、任务状态机

### 5.1 状态定义

| 状态 | 说明 | 所在端 |
|------|------|--------|
| queued | 已入队，等待分配到机器 | Backend |
| pending | 已分配到机器，等待 Agent 认领 | Backend |
| assigned | Agent 已认领，等待执行 | Backend + Agent |
| running | 正在执行 | Backend + Agent |
| completed | 执行成功（exit_code == 0） | Backend + Agent |
| failed | 执行失败（exit_code != 0） | Backend + Agent |
| cancelled | 被用户取消 | Backend |
| preempted | 被高优先级任务抢占 | Agent |
| suspended | 主动暂停 | Agent |

### 5.2 状态转换图

```
queued ──(分配到机器)──► pending ──(Agent ClaimTasks)──► assigned
                                                            │
                                                    (Agent StartTask)
                                                            │
                                                            ▼
                                                        running
                                                       /   |   \
                                              (成功)  / (失败) \  (抢占)
                                                    /     |     \
                                                   ▼      ▼      ▼
                                             completed  failed  preempted
                                                          │        │
                                                    (重试) │   (重新入队)
                                                          ▼        │
                                                       pending ◄───┘

cancelled ◄──(用户取消)── queued/pending/assigned
```

---

## 六、租约机制

### 6.1 设计目的

防止 Agent 崩溃后任务永久卡在 `running` 状态。

### 6.2 租约生命周期

| 阶段 | 时间 | 说明 |
|------|------|------|
| 认领时设置 | +5 分钟 | ClaimTasks 时设置初始租约 |
| 执行中续约 | 每 60 秒 | Agent 调用 RenewLease 延长 5 分钟 |
| 过期回收 | 租约到期后 | Backend HeartbeatMonitor 检测并重置任务 |

### 6.3 幂等性保证

`attempt_id` 是租约机制的核心：

1. 每次 ClaimTasks 时，Backend 为每个任务生成新的 `attempt_id`
2. 后续所有操作（Start、RenewLease、Complete）都必须携带匹配的 `attempt_id`
3. 若 Agent 崩溃后任务被重新认领，旧 Agent 的 `attempt_id` 失效，无法干扰新执行

---

## 七、Backend → Agent 反向通信

### 7.1 当前实现

Backend 通过 Agent 暴露的本地 HTTP API 与 Agent 通信（用于运维操作）。

Backend 的 `AgentService` 根据机器的 `agent_address`（IP:Port）直接调用 Agent 的 HTTP 端点。

### 7.2 Agent 本地 API 端点

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/ping` | 健康检查 |
| GET | `/api/v1/system/info` | 获取系统信息（CPU、内存、磁盘、GPU） |
| POST | `/api/v1/process/stop` | 停止指定进程 |
| POST | `/api/v1/ssh/reset` | 重置 SSH 授权密钥 |
| POST | `/api/v1/machine/cleanup` | 清理机器（Docker 容器、SSH 密钥） |
| POST | `/api/v1/command/exec` | 执行 shell 命令 |

### 7.3 安全建议

Agent 本地 API 暴露了高危操作端点，需要加强安全防护：

1. **添加认证**：Agent 本地 API 应使用与 Backend 相同的 Token 认证
2. **限制监听地址**：仅监听 `127.0.0.1` 或内网 IP
3. **命令白名单**：`/command/exec` 应受 `security.allowed_commands` 约束
4. **审计日志**：记录所有反向调用操作

---

## 八、配置参考

### 8.1 Agent 端配置（agent.yaml）

```yaml
port: 8081                    # Agent HTTP 服务端口
db_path: "./agent.db"         # SQLite 数据库路径
max_workers: 4                # 最大并发任务数

server:
  url: "http://backend:8080"  # Backend 地址
  agent_id: "agent-uuid-001"  # Agent 唯一标识
  machine_id: "host-uuid-001" # 机器唯一标识
  token: "shared-secret"      # 认证 Token
  timeout: 30s                # HTTP 请求超时

poll:
  interval: 5s                # 轮询间隔
  batch_size: 10              # 每次认领数量

limits:
  max_output_size: 1048576    # 输出限制（1MB）

security:
  allowed_commands:            # 命令白名单（空则不限制）
    - python
    - bash
  blocked_patterns:            # 命令黑名单模式
    - "rm -rf /"
```

### 8.2 Backend 端配置（config.yaml 中 agent 部分）

```yaml
agent:
  token: "shared-secret"      # 与 Agent 端相同的 Token
  port: 8081                  # Agent HTTP 服务端口（用于反向调用）
  timeout: 30                 # 调用 Agent API 的超时秒数
```

---

## 九、已知问题与改进建议

### 9.1 认证 Header 不一致（P0）

- **问题**：Backend 检查 `X-Agent-Token`，Agent 发送 `Authorization: Bearer`
- **影响**：Agent 无法通过认证
- **修复**：统一使用 `X-Agent-Token`

### 9.2 进度上报路由缺失（P1）

- **问题**：Agent 客户端实现了 `ReportProgress`，但 Backend 未注册 `/api/v1/agent/tasks/:id/progress` 路由
- **影响**：进度上报请求返回 404
- **修复**：在 Backend 添加对应的 Controller 和路由

### 9.3 CompleteTask 缺少 stdout/stderr 存储（P1）

- **问题**：Agent 上报 `stdout`/`stderr`，但 Backend `AgentTaskController.CompleteTask` 的请求体未包含这两个字段
- **影响**：任务输出丢失
- **修复**：Backend 请求体增加 `stdout`/`stderr` 字段，DAO 层同步更新

### 9.4 缺少 Agent 注册机制（P2）

- **问题**：Agent 启动后直接开始心跳和轮询，无注册流程
- **影响**：Backend 无法感知新 Agent 上线，无法管理 Agent 生命周期
- **建议**：增加 `POST /api/v1/agent/register` 端点，Agent 启动时注册自身信息

### 9.5 轮询模式的局限性（P3）

- **问题**：Agent 通过 HTTP 轮询获取任务，存在延迟（最大等于轮询间隔）
- **影响**：任务下发不够实时
- **建议**：未来可考虑 WebSocket 或 Server-Sent Events 实现推送

---

## 十、API 端点汇总

### 10.1 Agent → Backend（已实现）

| 方法 | 路径 | 功能 | 认证 |
|------|------|------|------|
| POST | `/api/v1/agent/heartbeat` | 心跳上报 | X-Agent-Token |
| POST | `/api/v1/agent/tasks/claim` | 认领任务 | X-Agent-Token |
| POST | `/api/v1/agent/tasks/:id/start` | 任务开始 | X-Agent-Token |
| POST | `/api/v1/agent/tasks/:id/lease/renew` | 租约续约 | X-Agent-Token |
| POST | `/api/v1/agent/tasks/:id/complete` | 任务完成 | X-Agent-Token |

### 10.2 Agent → Backend（待实现）

| 方法 | 路径 | 功能 | 认证 |
|------|------|------|------|
| POST | `/api/v1/agent/register` | Agent 注册 | X-Agent-Token |
| POST | `/api/v1/agent/tasks/:id/progress` | 进度上报 | X-Agent-Token |

### 10.3 Backend → Agent（反向调用）

| 方法 | 路径 | 功能 | 认证 |
|------|------|------|------|
| GET | `/api/v1/ping` | 健康检查 | 待添加 |
| GET | `/api/v1/system/info` | 系统信息 | 待添加 |
| POST | `/api/v1/process/stop` | 停止进程 | 待添加 |
| POST | `/api/v1/ssh/reset` | 重置 SSH | 待添加 |
| POST | `/api/v1/machine/cleanup` | 机器清理 | 待添加 |
| POST | `/api/v1/command/exec` | 执行命令 | 待添加 |
