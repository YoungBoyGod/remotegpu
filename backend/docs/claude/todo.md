# 待完成任务清单

本文档记录项目中所有待完成的功能模块，包括未完成原因、实现建议和依赖项。

**作者**: Claude
**创建日期**: 2026-02-04

---

## 优先级说明

- **P0**: 核心功能，影响系统正常运行
- **P1**: 重要功能，影响用户体验
- **P2**: 增强功能，可延后实现
- **P3**: 优化项，非必需

---

## 1. 监控系统 (P1)

### 1.1 Redis 缓存接入 ✅ 已完成

**位置**: `internal/service/ops/monitor_service.go:27`

**完成时间**: 2026-02-05

**实现内容**:
- 在 `MonitorService` 中注入 Redis 客户端
- 缓存键: `monitor:snapshot`
- 缓存过期时间: 30 秒
- 缓存未命中时查询数据库并写入缓存

**修改文件**:
- `internal/service/ops/monitor_service.go`
- `internal/router/router.go`

---

### 1.2 GPU 利用率监控 ✅ 已完成

**位置**: `internal/service/ops/monitor_service.go`

**完成时间**: 2026-02-05

**实现内容**:
- 创建 Prometheus 客户端 (`pkg/prometheus/client.go`)
- 创建 GPU 指标查询辅助方法 (`pkg/prometheus/gpu.go`)
- MonitorService 集成 Prometheus 获取 GPU 利用率
- DashboardService 集成 Prometheus 获取 GPU 趋势
- 支持多种 GPU exporter 指标格式 (DCGM, nvidia_gpu_exporter, nvidia_smi)

**Prometheus 地址**: `192.168.10.210:19090`

**修改文件**:
- `pkg/prometheus/client.go` (新建)
- `pkg/prometheus/gpu.go` (新建)
- `internal/service/ops/monitor_service.go`
- `internal/router/router.go`

---

## 2. Agent 通信 (P0) ✅ 已完成

### 2.1 任务进程管理 ✅ 已完成

**位置**: `internal/service/task/task_service.go:33-45`

**完成时间**: 2026-02-05

**实现内容**:
- `TaskService.StopTask()` 调用 `AgentService.StopProcess()` 停止进程
- `TaskService.StopTaskWithAuth()` 带权限校验的停止任务方法
- 支持 HTTP 和 gRPC 双协议通信

**修改文件**:
- `internal/service/task/task_service.go`
- `internal/service/ops/agent_service.go`
- `internal/agent/http_client.go`
- `internal/agent/grpc_client.go`

---

### 2.2 SSH 重置 ✅ 已完成

**位置**: `internal/service/allocation/allocation_service.go:91-95`

**完成时间**: 2026-02-05

**实现内容**:
- `AgentService.ResetSSH()` 方法实现
- `AllocationService.AllocateMachine()` 中异步调用 ResetSSH
- HTTP 和 gRPC 客户端都实现了 ResetSSH 方法

**修改文件**:
- `internal/service/ops/agent_service.go`
- `internal/service/allocation/allocation_service.go`
- `internal/agent/http_client.go`
- `internal/agent/grpc_client.go`

---

### 2.3 机器回收清理流程 ✅ 已完成

**位置**: `internal/service/allocation/allocation_service.go:152-157`

**完成时间**: 2026-02-05

**实现内容**:
- `AgentService.CleanupMachine()` 方法实现
- `AllocationService.ReclaimMachine()` 中异步调用 CleanupMachine
- 清理类型包括: process, data, ssh, docker
- HTTP 和 gRPC 客户端都实现了 CleanupMachine 方法

**修改文件**:
- `internal/service/ops/agent_service.go`
- `internal/service/allocation/allocation_service.go`
- `internal/agent/http_client.go`
- `internal/agent/grpc_client.go`

### Agent 客户端架构

**实现内容**:
- `internal/agent/client.go` - 定义 Client 接口
- `internal/agent/http_client.go` - HTTP 客户端完整实现
- `internal/agent/grpc_client.go` - gRPC 客户端完整实现
- `api/proto/agent/agent.proto` - Protocol Buffers 定义
- 支持双协议 (HTTP/gRPC)，通过配置切换
- 支持 TLS 加密通信

---

## 3. 镜像同步 (P1)

### 3.1 Harbor 镜像同步

**位置**: `internal/service/image/image_service.go:42`

**当前状态**: 空实现，返回 TODO 错误

**未完成原因**:
- 需要 Harbor API 凭证配置
- 需要设计同步策略（全量/增量）
- Harbor 服务尚未部署

**实现步骤**:
1. 添加 Harbor 配置到 `config.yaml`
2. 实现 Harbor API 客户端
3. 获取 Harbor 仓库镜像列表
4. 对比本地数据库，同步差异
5. 添加同步状态跟踪和失败重试

**依赖**:
- Harbor API v2.0
- `github.com/goharbor/go-client` - Harbor Go 客户端

**配置示例**:
```yaml
harbor:
  url: "https://harbor.example.com"
  username: "admin"
  password: "${HARBOR_PASSWORD}"
  project: "remotegpu"
```

---

## 4. 远程访问服务 (P2)

### 4.1 VNC 服务

**位置**: `internal/service/vnc.go:12-144`

**当前状态**: 框架代码，所有方法返回 TODO 错误

**未完成原因**:
- VNC 服务架构未确定
- 需要与容器环境集成
- 桌面环境配置复杂

**实现步骤**:
1. 选择 VNC 服务器 (TigerVNC/x11vnc)
2. 实现 VNC 配置生成
3. 集成到容器启动流程
4. 实现 noVNC Web 代理

**依赖**:
- TigerVNC 或 x11vnc
- noVNC (Web VNC 客户端)
- WebSocket 代理

---

### 4.2 RDP 服务

**位置**: `internal/service/rdp.go:12-185`

**当前状态**: 框架代码，所有方法返回 TODO 错误

**未完成原因**:
- 需要 xrdp 配置
- Linux 桌面环境支持复杂
- 需要与容器/虚拟机集成

**实现步骤**:
1. 实现 xrdp 配置生成
2. 集成到容器/虚拟机启动流程
3. 实现 RDP 网关代理

**依赖**:
- xrdp
- Apache Guacamole (可选，统一网关)

---

### 4.3 Guacamole 集成

**位置**: `internal/service/guacamole.go:12-175`

**当前状态**: 框架代码，所有方法返回 TODO 错误

**未完成原因**:
- 需要部署 Guacamole 服务
- 需要实现 Guacamole API 客户端
- 统一网关方案尚未确定

**实现步骤**:
1. 部署 Guacamole Server + guacd
2. 实现 Guacamole REST API 客户端
3. 实现连接创建/删除/更新
4. 生成访问 Token 和 URL

**依赖**:
- Apache Guacamole Server
- guacd (Guacamole 代理守护进程)
- PostgreSQL/MySQL (Guacamole 数据库)

---

## 5. 网络服务 (P2)

### 5.1 DNS 管理

**位置**: `internal/service/dns.go:12-225`

**当前状态**: 框架代码，支持多云厂商但未实现

**未完成原因**:
- 需要各云厂商 API 凭证
- DNS 记录管理策略未确定
- 多云厂商适配工作量大

**实现步骤**:
1. 添加 DNS 配置到 `config.yaml`
2. 实现各云厂商 SDK 集成：
   - Cloudflare: `github.com/cloudflare/cloudflare-go`
   - 阿里云: `github.com/aliyun/alibaba-cloud-sdk-go`
   - 腾讯云: `github.com/tencentcloud/tencentcloud-sdk-go`
   - AWS: `github.com/aws/aws-sdk-go-v2`
3. 实现统一的 DNS 操作接口

**依赖**:
- 云厂商 SDK
- DNS 域名和 API 凭证

---

### 5.2 防火墙管理

**位置**: `internal/service/firewall.go:12-178`

**当前状态**: 框架代码，支持多种防火墙但未实现

**未完成原因**:
- 需要确定防火墙类型 (iptables/firewalld/云厂商)
- 需要 root 权限或 Agent 支持
- 安全性考虑需要谨慎设计

**实现步骤**:
1. 确定防火墙管理方式
2. 实现 iptables 规则管理（通过 Agent）
3. 或实现云厂商安全组 API 调用
4. 添加规则同步和清理机制

**依赖**:
- Agent 防火墙管理接口
- 或云厂商安全组 API

---

## 6. 容器管理 (P2)

### 6.1 Docker 操作

**位置**: `internal/service/docker.go:238-402`

**当前状态**: 框架代码，大部分方法返回 TODO 错误

**未完成原因**:
- 需要 Docker API 客户端
- 需要与 Agent 集成
- GPU 支持需要 nvidia-container-toolkit

**实现步骤**:
1. 使用 Docker SDK: `github.com/docker/docker/client`
2. 实现容器 CRUD 操作
3. 实现镜像管理
4. 添加 GPU 支持 (nvidia-docker)

**依赖**:
- Docker SDK for Go
- nvidia-container-toolkit (GPU 支持)

**代码示例**:
```go
import (
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
)

func (s *DockerService) CreateContainer(ctx context.Context, config *ContainerConfig) (string, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return "", err
    }

    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: config.Image,
        Env:   config.Env,
    }, &container.HostConfig{
        Resources: container.Resources{
            DeviceRequests: []container.DeviceRequest{
                {
                    Driver: "nvidia",
                    Count:  -1, // all GPUs
                },
            },
        },
    }, nil, nil, config.Name)

    return resp.ID, err
}
```

---

## 7. 机器管理 (P3)

### 7.1 IP 唯一性校验

**位置**: `internal/service/machine/machine_service.go:26`

**当前状态**: TODO 标记，未实现

**未完成原因**:
- 简单功能，优先级较低
- 需要确定校验范围（全局/租户级别）

**实现步骤**:
1. 在 `MachineDao` 添加 `FindByIP` 方法
2. 在 `CreateMachine` 中调用检查
3. IP 重复时返回错误

**依赖**: 无

**代码示例**:
```go
func (s *MachineService) CreateMachine(ctx context.Context, machine *entity.Host) error {
    // 检查 IP 唯一性
    existing, err := s.machineDao.FindByIP(ctx, machine.IP)
    if err == nil && existing != nil {
        return errors.New(errors.ErrorInvalidParams, "IP address already exists")
    }

    return s.machineDao.Create(ctx, machine)
}
```

---

## 总结

| 模块 | 优先级 | 主要依赖 | 状态 |
|------|--------|----------|------|
| Redis 缓存 | P1 | pkg/cache | 待开发 |
| GPU 监控 | P1 | Prometheus | 待开发 |
| Agent 通信 | P0 | Agent 服务 | 待开发 |
| Harbor 同步 | P1 | Harbor API | 待开发 |
| VNC/RDP | P2 | Guacamole | 待开发 |
| DNS 管理 | P2 | 云厂商 SDK | 待开发 |
| 防火墙 | P2 | Agent | 待开发 |
| Docker | P2 | Docker SDK | 待开发 |
| IP 校验 | P3 | 无 | 待开发 |

**建议实现顺序**: Agent 通信 → Redis 缓存 → GPU 监控 → Harbor 同步 → Docker → 其他
