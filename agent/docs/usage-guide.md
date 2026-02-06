# RemoteGPU Agent 使用指南

本文档说明如何部署和使用 RemoteGPU Agent 系统，包括服务端（Backend）和客户端（Agent）的配置与使用。

## 目录

1. [系统架构概述](#系统架构概述)
2. [服务端配置](#服务端配置)
3. [Agent 客户端部署](#agent-客户端部署)
4. [配置说明](#配置说明)
5. [使用示例](#使用示例)
6. [故障排查](#故障排查)

---

## 系统架构概述

RemoteGPU Agent 系统采用 Server-Agent 架构：

```
┌─────────────┐         HTTP API          ┌──────────────┐
│   Backend   │ ◄────────────────────────► │    Agent     │
│   Server    │                            │  (GPU 机器)   │
└─────────────┘                            └──────────────┘
      │                                           │
      │                                           │
   PostgreSQL                                  SQLite
   (任务存储)                                  (本地缓存)
```

**工作流程：**
1. 用户通过 Backend API 创建任务
2. Agent 定期轮询 Backend 获取待执行任务（ClaimTasks）
3. Agent 执行任务并上报状态（Start/Progress/Complete）
4. Agent 通过租约机制保持任务所有权
5. 离线时任务结果缓存到本地，恢复后自动同步

---

## 服务端配置

### 1. 数据库准备

Backend 使用 PostgreSQL 存储任务数据。确保已执行所有迁移脚本：

```bash
cd backend
# 按顺序执行 sql/ 目录下的迁移脚本
psql -U postgres -d remotegpu -f sql/01_init.sql
psql -U postgres -d remotegpu -f sql/02_add_tasks.sql
# ... 执行所有迁移脚本
```

### 2. 启动 Backend 服务

```bash
cd backend
go run ./cmd/main.go server
```

Backend 默认监听 `8080` 端口。

### 3. Agent API 端点

Backend 提供以下 Agent API 端点（需要 Token 认证）：

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/agent/tasks/claim` | POST | Agent 认领任务 |
| `/api/v1/agent/tasks/:id/start` | POST | 上报任务开始 |
| `/api/v1/agent/tasks/:id/renew` | POST | 续约任务租约 |
| `/api/v1/agent/tasks/:id/complete` | POST | 上报任务完成 |
| `/api/v1/agent/tasks/:id/progress` | POST | 上报任务进度 |

所有请求需要在 Header 中携带：`Authorization: Bearer <token>`

---

## Agent 客户端部署

### 1. 编译 Agent

```bash
cd agent
go build -o remotegpu-agent ./cmd/main.go
```

### 2. 创建配置文件

将 `agent.yaml.example` 复制为 `agent.yaml`：

```bash
cp agent.yaml.example agent.yaml
```

编辑配置文件，填入必要信息：

```yaml
port: 8090
db_path: /var/lib/remotegpu-agent/tasks.db
max_workers: 4

server:
  url: "http://your-backend-server:8080"
  agent_id: "agent-001"
  machine_id: "machine-001"
  token: "your-agent-token-here"
  timeout: 30s

poll:
  interval: 5s
  batch_size: 10
```

### 3. 启动 Agent

```bash
# 前台运行
./remotegpu-agent

# 后台运行（推荐使用 systemd）
nohup ./remotegpu-agent > agent.log 2>&1 &
```

Agent 启动后会：
- 监听本地 HTTP 端口（默认 8090）
- 定期轮询 Backend 获取任务
- 自动执行任务并上报结果
- 离线时缓存结果，恢复后自动同步

---

## 配置说明

### 基础配置

| 配置项 | 说明 | 默认值 | 环境变量 |
|--------|------|--------|----------|
| `port` | Agent HTTP 监听端口 | 8090 | `AGENT_PORT` |
| `db_path` | SQLite 数据库路径 | `/var/lib/remotegpu-agent/tasks.db` | `AGENT_DB_PATH` |
| `max_workers` | 最大并发任务数 | 4 | `AGENT_MAX_WORKERS` |

### Server 配置

| 配置项 | 说明 | 环境变量 |
|--------|------|----------|
| `server.url` | Backend 服务地址 | `SERVER_URL` |
| `server.agent_id` | Agent 唯一标识 | `AGENT_ID` |
| `server.machine_id` | 机器 ID | `MACHINE_ID` |
| `server.token` | 认证 Token | `AGENT_TOKEN` |
| `server.timeout` | 请求超时时间 | - |

### 安全配置

```yaml
security:
  # 命令白名单：仅允许执行列表中的命令（为空则不限制）
  allowed_commands:
    - "python"
    - "python3"
    - "nvidia-smi"

  # 命令黑名单：禁止执行包含这些模式的命令
  blocked_patterns:
    - "rm -rf /"
    - "mkfs"
```

**说明：**
- `allowed_commands` 为空时不限制命令
- `blocked_patterns` 始终生效，优先级高于白名单
- 命令校验失败的任务会立即标记为 failed

---

## 使用示例

### 示例 1：通过 Backend API 创建任务

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <user_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_id": "machine-001",
    "command": "python3",
    "args": ["train.py", "--epochs", "10"],
    "work_dir": "/workspace",
    "priority": 5,
    "timeout": 3600
  }'
```

Agent 会自动：
1. 通过轮询获取该任务
2. 执行命令并捕获输出
3. 上报执行结果到 Backend

---

### 示例 2：本地直接提交任务到 Agent

```bash
curl -X POST http://localhost:8090/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "id": "local-task-001",
    "command": "nvidia-smi",
    "priority": 1
  }'
```

本地提交的任务不会同步到 Backend，仅在 Agent 本地执行。

### 示例 3：查询 Agent 状态

```bash
# 查看队列状态
curl http://localhost:8090/api/queue/status

# 查看特定任务
curl http://localhost:8090/api/tasks/task-001
```

---

## 故障排查

### 1. Agent 无法连接到 Backend

**症状：** Agent 日志显示连接错误

**排查步骤：**
- 检查 `server.url` 配置是否正确
- 确认 Backend 服务正在运行
- 检查网络连接和防火墙设置
- 验证 Token 是否正确配置

### 2. 任务一直处于 pending 状态

**可能原因：**
- Agent 未启动或未连接到 Backend
- `machine_id` 配置不匹配
- Agent 已达到 `max_workers` 限制
- 任务依赖未满足（检查 `depends_on` 字段）

**解决方法：**
```bash
# 查看 Agent 队列状态
curl http://localhost:8090/api/queue/status
```

### 3. 任务执行失败

**排查步骤：**
- 查看任务的 `stderr` 输出
- 检查命令是否被安全策略拦截
- 确认工作目录和环境变量配置正确
- 检查任务超时设置是否合理

### 4. 租约过期问题

**症状：** 任务被标记为 failed，错误信息为 "lease expired"

**原因：** Agent 与 Backend 网络中断超过 5 分钟

**解决方法：**
- 检查网络稳定性
- 任务会自动重试（如果配置了 `max_retries`）

### 5. 查看日志

Agent 使用结构化日志（slog），日志级别包括：
- `INFO`: 正常运行信息
- `WARN`: 警告信息（如租约过期）
- `ERROR`: 错误信息（如网络故障）
- `DEBUG`: 调试信息

**查看实时日志：**
```bash
# 如果使用 nohup 启动
tail -f agent.log

# 如果使用 systemd
journalctl -u remotegpu-agent -f
```

---

## 高级特性

### 1. 任务优先级和抢占

Agent 支持基于优先级的任务调度：
- 优先级数字越小，优先级越高（1 > 5 > 10）
- 当新任务优先级比运行中最低优先级任务高 3 级以上时，会触发抢占
- 被抢占的任务会重新入队等待执行

**示例：**
```json
{
  "command": "python train.py",
  "priority": 1  // 高优先级任务
}
```

### 2. 任务依赖

支持通过 `depends_on` 字段指定任务依赖关系：

```json
{
  "id": "task-002",
  "command": "python evaluate.py",
  "depends_on": ["task-001"]  // 等待 task-001 完成后执行
}
```

Agent 会自动检查依赖，只有所有依赖任务完成后才会执行当前任务。

### 3. 自动重试机制

任务失败时可以自动重试：

```json
{
  "command": "python train.py",
  "max_retries": 3,      // 最多重试 3 次
  "retry_delay": 60      // 重试间隔 60 秒
}
```

重试逻辑：
- 任务失败后，如果 `retry_count < max_retries`，会自动重新入队
- 重试时会清除之前的错误信息和输出
- 重试作为本地任务执行，不会生成新的 `attempt_id`

---

## 最佳实践

### 1. 生产环境部署建议

**使用 systemd 管理 Agent：**

创建 `/etc/systemd/system/remotegpu-agent.service`：

```ini
[Unit]
Description=RemoteGPU Agent
After=network.target

[Service]
Type=simple
User=remotegpu
WorkingDirectory=/opt/remotegpu-agent
ExecStart=/opt/remotegpu-agent/remotegpu-agent
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
systemctl enable remotegpu-agent
systemctl start remotegpu-agent
```

### 2. 安全建议

- **启用命令白名单：** 在生产环境中务必配置 `security.allowed_commands`
- **配置黑名单模式：** 添加危险命令模式到 `security.blocked_patterns`
- **Token 管理：** 定期轮换 Agent Token，不要在代码中硬编码
- **网络隔离：** Agent 只需访问 Backend API，建议配置防火墙规则

### 3. 性能优化

- **合理设置 max_workers：** 根据机器 CPU/GPU 资源调整并发数
- **调整轮询间隔：** 高负载场景可减小 `poll.interval`，低负载场景可增大以减少网络开销
- **输出大小限制：** 默认 1MB 输出限制可防止内存溢出，如需更大输出可调整 `limits.max_output_size`

---

## 总结

RemoteGPU Agent 提供了完整的远程任务执行能力，主要特性包括：

✅ **可靠性：** 租约机制、自动重试、离线同步
✅ **安全性：** Token 认证、命令白名单/黑名单
✅ **灵活性：** 优先级调度、任务依赖、抢占式执行
✅ **易用性：** YAML 配置、环境变量覆盖、结构化日志

## 相关文档

- [Agent 架构设计](./task-queue-design.md)
- [Backend API 文档](../../backend/docs/openapi.yaml)
- [项目整体指南](../../.claude/CLAUDE.md)

---

**文档版本：** v1.0
**最后更新：** 2026-02-06
