# Agent 通信模块文档

**作者**: Claude
**创建日期**: 2026-02-05
**状态**: ✅ 已完成

---

## 概述

Agent 通信模块负责后端服务与远程 GPU 机器上的 Agent 服务之间的通信。支持 HTTP 和 gRPC 双协议，可通过配置切换。

## 架构设计

```
┌─────────────────┐         ┌─────────────────┐
│   Backend API   │         │   GPU Machine   │
│                 │         │                 │
│  AgentService   │◄───────►│  Agent Service  │
│                 │  HTTP/  │                 │
│  ┌───────────┐  │  gRPC   │  ┌───────────┐  │
│  │HTTPClient │  │         │  │ HTTP API  │  │
│  │GRPCClient │  │         │  │ gRPC API  │  │
│  └───────────┘  │         │  └───────────┘  │
└─────────────────┘         └─────────────────┘
```

## 文件结构

```
internal/
├── agent/
│   ├── client.go       # Client 接口定义
│   ├── types.go        # 请求/响应类型定义
│   ├── http_client.go  # HTTP 客户端实现
│   └── grpc_client.go  # gRPC 客户端实现
├── service/
│   └── ops/
│       └── agent_service.go  # Agent 服务层
api/
└── proto/
    └── agent/
        ├── agent.proto       # Protocol Buffers 定义
        ├── agent.pb.go       # 生成的消息类型
        └── agent_grpc.pb.go  # 生成的 gRPC 代码
```

---

## 配置说明

在 `config.yaml` 中配置 Agent 通信参数：

```yaml
agent:
  enabled: true
  protocol: "http"      # http 或 grpc
  http_port: 8081       # HTTP 端口
  grpc_port: 50051      # gRPC 端口
  timeout: 30           # 超时时间（秒）
  tls_enabled: false    # 是否启用 TLS
  tls_cert_file: ""     # TLS 证书文件路径
```

---

## Client 接口

```go
type Client interface {
    // 停止进程
    StopProcess(ctx context.Context, req *StopProcessRequest) (*Response, error)

    // 重置 SSH 密钥
    ResetSSH(ctx context.Context, req *ResetSSHRequest) (*Response, error)

    // 清理机器
    CleanupMachine(ctx context.Context, req *CleanupRequest) (*Response, error)

    // 挂载数据集
    MountDataset(ctx context.Context, req *MountDatasetRequest) (*Response, error)

    // 获取系统信息
    GetSystemInfo(ctx context.Context, hostID string) (*SystemInfo, error)

    // 执行命令
    ExecuteCommand(ctx context.Context, req *ExecuteCommandRequest) (*ExecuteCommandResponse, error)

    // 健康检查
    Ping(ctx context.Context, hostID string) error

    // 关闭连接
    Close() error
}
```

---

## 功能实现

### 1. 任务进程管理

**位置**: `internal/service/task/task_service.go`

```go
func (s *TaskService) StopTask(ctx context.Context, id string) error {
    task, err := s.taskDao.FindByID(ctx, id)
    if err != nil {
        return err
    }

    // 调用 Agent 停止进程
    if s.agentService != nil && task.HostID != "" {
        _ = s.agentService.StopProcess(ctx, task.HostID, task.ID)
    }

    return s.taskDao.UpdateStatus(ctx, id, "stopped")
}
```

### 2. SSH 重置

**位置**: `internal/service/allocation/allocation_service.go`

分配机器时异步重置 SSH：

```go
// 异步触发 Agent 重置 SSH
if s.agentClient != nil {
    go func() {
        _ = s.agentClient.ResetSSH(context.Background(), hostID)
    }()
}
```

### 3. 机器回收清理

**位置**: `internal/service/allocation/allocation_service.go`

回收机器时异步清理：

```go
// 异步触发清理流程
if s.agentClient != nil {
    go func() {
        _ = s.agentClient.CleanupMachine(context.Background(), hostID)
    }()
}
```

清理类型包括：
- `process` - 终止用户进程
- `data` - 删除用户数据
- `ssh` - 重置 SSH 密钥
- `docker` - 清理 Docker 容器

---

## 验证机制

### 1. 编译时验证

```bash
go build ./...
```

### 2. 健康检查

```go
// 检查 Agent 连接状态
err := agentService.CheckAgentHealth(ctx, hostID)
if err != nil {
    log.Printf("Agent 不可达: %v", err)
}
```

### 3. 响应状态检查

```go
type Response struct {
    Success   bool        `json:"success"`
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
}
```

---

## Agent 端 API 规范

Agent 服务需要实现以下 HTTP API：

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/ping` | 健康检查 |
| GET | `/api/v1/system/info` | 获取系统信息 |
| POST | `/api/v1/process/stop` | 停止进程 |
| POST | `/api/v1/ssh/reset` | 重置 SSH |
| POST | `/api/v1/machine/cleanup` | 清理机器 |
| POST | `/api/v1/dataset/mount` | 挂载数据集 |
| POST | `/api/v1/command/exec` | 执行命令 |

---

## 测试方法

### Mock Agent 测试

可以使用简单的 HTTP 服务模拟 Agent：

```go
// mock_agent_test.go
func TestAgentCommunication(t *testing.T) {
    // 启动 Mock Server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(agent.Response{
            Success: true,
            Code:    200,
            Message: "ok",
        })
    }))
    defer server.Close()

    // 测试连接
    // ...
}
```

### 集成测试

1. 部署 Agent 服务到测试机器
2. 配置后端连接参数
3. 调用 `CheckAgentHealth` 验证连接
4. 执行各功能测试

---

## 待改进项

1. **错误处理增强**: 异步调用目前忽略错误，建议添加日志记录
2. **重试机制**: 添加失败重试逻辑
3. **状态回调**: 异步操作完成后回调更新状态
4. **监控指标**: 添加 Prometheus 指标监控通信状态
