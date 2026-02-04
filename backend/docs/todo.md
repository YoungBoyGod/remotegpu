# 待完成任务清单

本文档记录项目中所有待完成的功能模块，包括未完成原因、实现建议和依赖项。

---

## 优先级说明

- **P0**: 核心功能，影响系统正常运行
- **P1**: 重要功能，影响用户体验
- **P2**: 增强功能，可延后实现
- **P3**: 优化项，非必需

---

## 1. 监控系统集成 (P1)

### 1.1 Redis 缓存接入

**位置**: `internal/service/ops/monitor_service.go:27`

**当前状态**: 每次请求都查询数据库

**未完成原因**:
- 需要设计缓存键结构
- 需要确定缓存过期策略

**实现步骤**:
1. 在 `MonitorService` 中注入 Redis 客户端
2. 设计缓存键: `monitor:snapshot:{timestamp}`
3. 实现缓存读写逻辑，设置 30 秒过期
4. 添加缓存失效机制

**依赖**:
- `pkg/cache` - Redis 客户端（已存在）

**代码示例**:
```go
func (s *MonitorService) GetGlobalSnapshot(ctx context.Context) (map[string]interface, error) {
    cacheKey := "monitor:snapshot"

    // 尝试从缓存读取
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }

    // 查询数据库
    stats, err := s.machineService.GetStatusStats(ctx)
    // ...

    // 写入缓存，30秒过期
    s.cache.Set(ctx, cacheKey, result, 30*time.Second)
    return result, nil
}
```

---

### 1.2 GPU 利用率监控

**位置**: `internal/service/ops/monitor_service.go:48,111-113`

**当前状态**: 返回硬编码值 `0.0`

**未完成原因**:
- 需要部署监控系统 (Prometheus/InfluxDB)
- 需要在 Agent 端采集 GPU 指标

**实现步骤**:
1. 部署 Prometheus + node_exporter + nvidia_gpu_exporter
2. 配置 Prometheus 抓取 GPU 指标
3. 在后端添加 Prometheus 客户端查询接口
4. 实现 `GetGPUTrend` 从 Prometheus 查询历史数据

**依赖**:
- Prometheus Server
- nvidia_gpu_exporter (GPU 指标采集)
- `github.com/prometheus/client_golang` - Prometheus Go 客户端

**配置示例** (prometheus.yml):
```yaml
scrape_configs:
  - job_name: 'gpu_nodes'
    static_configs:
      - targets: ['node1:9835', 'node2:9835']  # nvidia_gpu_exporter 端口
```

---

## 2. 镜像同步 (P1)

### 2.1 Harbor 镜像同步

**位置**: `internal/service/image/image_service.go:42`

**当前状态**: 空实现，返回 TODO 错误

**未完成原因**:
- 需要 Harbor API 凭证配置
- 需要设计同步策略（全量/增量）

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

## 3. Agent 通信 (P0)

### 3.1 任务进程管理

**位置**: `internal/service/task/task_service.go:31,58`

**当前状态**: 仅更新数据库状态，未实际停止进程

**未完成原因**:
- Agent 服务尚未完全实现
- 需要设计 Agent 通信协议

**实现步骤**:
1. 定义 Agent gRPC/HTTP API 接口
2. 实现 `AgentService.StopTask(hostID, taskID)` 方法
3. Agent 端实现进程查找和终止逻辑
4. 添加超时和重试机制

**依赖**:
- Agent 服务 (`internal/service/ops/agent_service.go`)
- gRPC 或 HTTP 客户端

---

### 3.2 SSH 重置

**位置**: `internal/service/allocation/allocation_service.go:82`

**当前状态**: 注释掉的 TODO

**未完成原因**:
- 需要 Agent 支持 SSH 密钥重置
- 需要异步任务队列

**实现步骤**:
1. Agent 实现 SSH 密钥重置接口
2. 后端调用 Agent API 触发重置
3. 使用消息队列实现异步处理
4. 添加重置状态回调

**依赖**:
- Agent SSH 管理模块
- 消息队列 (Redis/RabbitMQ)

---

### 3.3 机器回收清理流程

**位置**: `internal/service/allocation/allocation_service.go:140`

**当前状态**: TODO 标记，未实现

**未完成原因**:
- 清理流程涉及多个步骤
- 需要异步执行避免阻塞

**实现步骤**:
1. 定义清理任务结构
2. 实现清理步骤：
   - 终止所有用户进程
   - 删除用户数据
   - 重置 SSH 密钥
   - 清理临时文件
3. 使用工作队列异步执行
4. 清理完成后更新机器状态为 `idle`

**依赖**:
- Agent 清理接口
- 异步任务队列

---

## 4. 数据集挂载验证 (P1) - 已完成

### 4.1 机器所有权验证

**位置**: `internal/controller/v1/dataset/dataset_controller.go:109`

**状态**: ✅ 已由 CodeX 于 2026-02-04 完成

**实现**:
- `AllocationDao.FindActiveByHostAndCustomer` - 查询用户机器分配
- `AllocationService.ValidateHostOwnership` - 验证机器所有权

---

## 5. 远程访问服务 (P2)

### 5.1 VNC 服务

**位置**: `internal/service/vnc.go`

**当前状态**: 框架代码，所有方法返回 TODO 错误

**未完成原因**:
- VNC 服务架构未确定
- 需要与容器环境集成

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

### 5.2 RDP 服务

**位置**: `internal/service/rdp.go`

**当前状态**: 框架代码，所有方法返回 TODO 错误

**未完成原因**:
- 需要 xrdp 配置
- Windows 环境支持复杂

**实现步骤**:
1. 实现 xrdp 配置生成
2. 集成到容器/虚拟机启动流程
3. 实现 RDP 网关代理

**依赖**:
- xrdp
- Apache Guacamole (可选，统一网关)

---

### 5.3 Guacamole 集成

**位置**: `internal/service/guacamole.go`

**当前状态**: 框架代码，所有方法返回 TODO 错误

**未完成原因**:
- 需要部署 Guacamole 服务
- 需要实现 Guacamole API 客户端

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

## 6. 网络服务 (P2)

### 6.1 DNS 管理

**位置**: `internal/service/dns.go`

**当前状态**: 框架代码，支持多云厂商但未实现

**未完成原因**:
- 需要各云厂商 API 凭证
- DNS 记录管理策略未确定

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

### 6.2 防火墙管理

**位置**: `internal/service/firewall.go`

**当前状态**: 框架代码，支持多种防火墙但未实现

**未完成原因**:
- 需要确定防火墙类型 (iptables/firewalld/云厂商)
- 需要 root 权限或 Agent 支持

**实现步骤**:
1. 确定防火墙管理方式
2. 实现 iptables 规则管理（通过 Agent）
3. 或实现云厂商安全组 API 调用
4. 添加规则同步和清理机制

**依赖**:
- Agent 防火墙管理接口
- 或云厂商安全组 API

---

## 7. 容器管理 (P2)

### 7.1 Docker 操作

**位置**: `internal/service/docker.go`

**当前状态**: 框架代码，大部分方法返回 TODO 错误

**未完成原因**:
- 需要 Docker API 客户端
- 需要与 Agent 集成

**实现步骤**:
1. 使用 Docker SDK: `github.com/docker/docker/client`
2. 实现容器 CRUD 操作
3. 实现镜像管理
4. 添加 GPU 支持 (nvidia-docker)

**依赖**:
- Docker SDK for Go
- nvidia-container-toolkit (GPU 支持)

---

## 8. 机器管理 (P2)

### 8.1 IP 唯一性校验

**位置**: `internal/service/machine/machine_service.go:26`

**当前状态**: TODO 标记，未实现

**未完成原因**:
- 简单功能，优先级较低

**实现步骤**:
1. 在 `MachineDao` 添加 `FindByIP` 方法
2. 在 `CreateMachine` 中调用检查
3. IP 重复时返回错误

**依赖**: 无

---

## 总结

| 模块 | 优先级 | 预估工作量 | 主要依赖 | 状态 |
|------|--------|-----------|----------|------|
| Redis 缓存 | P1 | 1天 | pkg/cache | 待开发 |
| GPU 监控 | P1 | 3天 | Prometheus | 待开发 |
| Harbor 同步 | P1 | 2天 | Harbor API | 待开发 |
| Agent 通信 | P0 | 5天 | Agent 服务 | 待开发 |
| 数据集验证 | P1 | 0.5天 | AllocationDao | ✅ 已完成 |
| VNC/RDP | P2 | 5天 | Guacamole | 待开发 |
| DNS 管理 | P2 | 3天 | 云厂商 SDK | 待开发 |
| 防火墙 | P2 | 3天 | Agent | 待开发 |
| Docker | P2 | 3天 | Docker SDK | 待开发 |

**建议实现顺序**: Agent 通信 → Redis 缓存 → GPU 监控 → Harbor 同步 → 其他
