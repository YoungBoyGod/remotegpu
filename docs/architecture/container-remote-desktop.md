# 容器管理与远程桌面架构设计

> 文档编号：ARCH-068
> 创建日期：2026-02-07
> 状态：草案

## 1. 概述

本文档基于 `backend/internal/service/_wip/` 目录下的代码审查结果，制定容器管理和远程桌面接入的整体架构方案。涵盖容器生命周期管理、VNC/RDP 远程桌面接入、端口池管理策略，以及与现有 Agent 的集成方式。

### 1.1 涉及模块

| 文件 | 模块 | 职责 |
|------|------|------|
| `docker.go` | DockerService | Docker 容器生命周期管理 |
| `vnc.go` | VNCService | VNC 桌面环境配置与管理 |
| `rdp.go` | RDPService | RDP 远程桌面配置与管理 |
| `guacamole.go` | GuacamoleService | Apache Guacamole Web 远程桌面网关集成 |
| `port_pool.go` | PortPoolService | 端口池分配与回收 |
| `dns.go` | DNSService | DNS 子域名自动配置 |
| `firewall.go` | FirewallService | 防火墙端口映射管理 |

---

## 2. 各模块完成度评估

### 2.1 DockerService（docker.go）— 完成度 50%

**已实现：**
- `NewDockerService()` — Docker 客户端初始化（使用环境变量自动协商 API 版本）
- `CreateContainer()` — 完整的容器创建流程（拉取镜像 → 端口映射 → 存储挂载 → 资源限制 → GPU 配置 → 创建 → 启动）
- `StartContainer()` / `StopContainer()` / `RestartContainer()` / `DeleteContainer()` — 基础生命周期操作
- `PullImage()` — 镜像拉取
- `buildPortConfig()` / `buildMountConfig()` / `buildGPUConfig()` — 配置构建辅助方法
- `buildContainerConfig()` — 从 Environment 实体构建容器配置

**未实现（TODO）：**
- `GetContainer()` — 获取容器详情
- `ListContainers()` — 列出容器
- `GetContainerLogs()` — 获取容器日志
- `ExecCommand()` — 容器内执行命令
- `GetContainerStats()` — 容器资源统计
- `ListImages()` / `DeleteImage()` — 镜像管理
- `CreateEnvironmentContainer()` / `DeleteEnvironmentContainer()` — 环境级容器操作
- `configureGPU()` — GPU 高级配置（vGPU/MIG/共享模式）
- `configureStorage()` — 存储配置（NFS/Ceph/S3/JuiceFS）

**评估：** 核心框架完整，基础 CRUD 可用。缺少查询类方法和环境级集成方法。GPU 高级模式和分布式存储配置是重要缺失。

### 2.2 VNCService（vnc.go）— 完成度 20%

**已实现：**
- 数据结构定义完整（`VNCConfig`、`VNCEnvironmentConfig`）
- `GetVNCEnvironmentVariables()` — 环境变量生成
- `GetVNCPorts()` — 端口需求声明（含 noVNC）
- `GetDesktopPackages()` / `GetNoVNCPackages()` — 包依赖列表

**未实现（TODO）：**
- `GenerateVNCConfig()` — 配置生成
- `GenerateVNCStartupScript()` — 启动脚本生成
- `GenerateVNCDockerfile()` — Dockerfile 生成

**评估：** 仅有数据结构和辅助方法，核心配置生成和脚本生成均未实现。设计思路清晰，支持 TigerVNC + noVNC 的方案合理。

### 2.3 RDPService（rdp.go）— 完成度 20%

**已实现：**
- 数据结构定义完整（`RDPConfig`、`RDPEnvironmentConfig`）
- `GetRDPEnvironmentVariables()` — 环境变量生成
- `GetRDPPorts()` — 端口需求声明
- `GetXRDPPackages()` — 包依赖列表（支持 xfce/lxde/gnome/mate）

**未实现（TODO）：**
- `GenerateRDPConfig()` — 配置生成
- `GenerateRDPStartupScript()` — 启动脚本生成
- `GenerateRDPDockerfile()` — Dockerfile 生成
- `ConfigureXRDP()` — xrdp 服务配置

**评估：** 与 VNCService 类似，仅有骨架。Linux 环境下基于 xrdp 的方案合理，但实现优先级低于 VNC（GPU 场景以 Linux 为主）。

### 2.4 GuacamoleService（guacamole.go）— 完成度 15%

**已实现：**
- 数据结构定义（`GuacamoleConfig`、`ConnectionRequest`、`ConnectionResponse`）
- 协议参数构建方法（`buildSSHParameters()`、`buildRDPParameters()`、`buildVNCParameters()`）

**未实现（TODO）：**
- Guacamole 客户端初始化
- 所有 CRUD 操作（`CreateConnection`、`DeleteConnection`、`UpdateConnection`、`GetConnection`）
- `ConfigureEnvironmentGuacamole()` — 环境级配置
- `CleanupEnvironmentGuacamole()` — 环境清理
- `GetGuacamoleAccessURL()` — 访问 URL 构建

**评估：** 基本只有接口定义和参数模板。Guacamole 作为 Web 远程桌面网关的选型合理，但需要完整实现 REST API 客户端。

### 2.5 PortPoolService（port_pool.go）— 完成度 85%

**已实现：**
- `AllocatePort()` — 单端口分配（带互斥锁保护）
- `AllocatePorts()` — 批量端口分配
- `ReleasePort()` / `ReleasePortsByEnvID()` — 端口释放
- `GetAllocatedPorts()` / `GetPortByServiceType()` — 端口查询
- 完整的端口范围定义（SSH/RDP/Jupyter/VNC/noVNC/自定义）

**未实现：**
- 缺少 `entity.PortMapping` 实体定义（当前 entity 目录中不存在）
- 缺少对应的数据库迁移脚本
- 未考虑多主机场景下的端口隔离（当前查询所有 active 端口，未按 host 过滤）

**评估：** 完成度最高的模块，核心逻辑已可用。需要补充 PortMapping 实体和修复多主机端口隔离问题。

### 2.6 DNSService（dns.go）— 完成度 10%

**已实现：**
- 数据结构定义（`DNSConfig`、`DNSRecordRequest`、`DNSRecordResponse`）
- `GenerateSubdomain()` — 子域名生成规则
- `BuildAccessURL()` — 访问 URL 构建（支持 ssh/rdp/jupyter/vnc/novnc 协议）

**未实现（TODO）：**
- DNS 客户端初始化
- 所有 CRUD 操作
- 多云厂商适配（Cloudflare/阿里云/腾讯云/AWS）
- `ConfigureEnvironmentDNS()` / `CleanupEnvironmentDNS()`

**评估：** 仅有骨架。DNS 自动配置是锦上添花的功能，优先级较低。

### 2.7 FirewallService（firewall.go）— 完成度 10%

**已实现：**
- 数据结构定义（`FirewallConfig`、`PortMappingRequest`、`PortMappingResponse`）

**未实现（TODO）：**
- 防火墙客户端初始化
- 所有 CRUD 操作
- 多类型适配（iptables/firewalld/云厂商防火墙）
- `ConfigureEnvironmentFirewall()` / `CleanupEnvironmentFirewall()`

**评估：** 仅有骨架。在当前 Agent 直连模式下，防火墙配置可通过 Agent 执行命令实现，独立服务优先级低。

---

## 3. 容器生命周期管理方案

### 3.1 架构决策：Agent 侧执行 vs Backend 侧执行

**现状分析：**
- 当前 `DockerService` 直接调用 Docker SDK，意味着它运行在能访问 Docker daemon 的机器上
- 现有 Agent 部署在 GPU 主机上，通过 HTTP 轮询与 Backend 通信
- Backend 通过 `agent.Client` 接口（HTTP/gRPC）向 Agent 下发命令

**推荐方案：Agent 侧执行容器操作**

```
┌─────────┐     API 请求      ┌─────────┐    任务下发     ┌─────────┐
│ Frontend │ ──────────────→  │ Backend │ ──────────────→ │  Agent  │
│  (Vue3)  │                  │  (Gin)  │                 │  (Go)   │
└─────────┘                   └─────────┘                 └────┬────┘
                                   │                          │
                              状态持久化                  Docker SDK
                                   │                          │
                              ┌────▼────┐                ┌────▼────┐
                              │PostgreSQL│                │ Docker  │
                              └─────────┘                │ Daemon  │
                                                         └─────────┘
```

**理由：**
1. Agent 已部署在 GPU 主机上，天然可访问本地 Docker daemon
2. 避免 Backend 直连远程 Docker daemon 的网络安全风险
3. 复用现有 Agent 任务执行框架（轮询 → 执行 → 上报）
4. 容器操作本质上是"在远程主机上执行的任务"，与现有 Task 模型一致

### 3.2 容器生命周期状态机

```
                    ┌──────────┐
                    │ creating │
                    └────┬─────┘
                         │ 创建成功
                    ┌────▼─────┐
              ┌────→│ running  │←────┐
              │     └────┬─────┘     │
              │          │           │
           启动     停止/暂停     重启
              │          │           │
              │     ┌────▼─────┐     │
              └─────│ stopped  │─────┘
                    └────┬─────┘
                         │ 删除
                    ┌────▼─────┐
                    │ deleted  │
                    └──────────┘

异常路径:
  creating → error（创建失败，自动清理）
  running  → error（运行异常，可重启或删除）
```

### 3.3 容器操作流程

#### 创建环境容器

```
1. 用户请求创建环境
2. Backend 校验配额和资源
3. Backend 调用 PortPoolService 分配端口
4. Backend 构建 ContainerConfig（镜像、资源、端口、GPU）
5. Backend 创建 Task（type=container_create），下发给 Agent
6. Agent 认领任务，调用 DockerService.CreateContainer()
7. Agent 上报容器 ID 和状态
8. Backend 更新 Environment 记录（container_id, status=running）
```

#### 停止/启动/重启容器

```
1. 用户请求操作
2. Backend 创建 Task（type=container_stop/start/restart）
3. Agent 执行对应的 DockerService 方法
4. Agent 上报结果
5. Backend 更新 Environment 状态
```

#### 销毁环境容器

```
1. 用户请求销毁环境
2. Backend 创建 Task（type=container_delete）
3. Agent 执行 DockerService.DeleteContainer()
4. Agent 上报结果
5. Backend 调用 PortPoolService.ReleasePortsByEnvID() 释放端口
6. Backend 更新 Environment 状态为 deleted
```

### 3.4 新增 Task 类型

在现有 Task 模型基础上扩展以下类型：

| Task Type | 说明 | 参数（Args JSON） |
|-----------|------|-------------------|
| `container_create` | 创建容器 | ContainerConfig JSON |
| `container_start` | 启动容器 | `{container_id}` |
| `container_stop` | 停止容器 | `{container_id, timeout}` |
| `container_restart` | 重启容器 | `{container_id, timeout}` |
| `container_delete` | 删除容器 | `{container_id, force}` |
| `container_inspect` | 查询容器状态 | `{container_id}` |
| `container_logs` | 获取容器日志 | `{container_id, tail}` |
| `container_exec` | 容器内执行命令 | `{container_id, cmd}` |
| `container_stats` | 获取资源统计 | `{container_id}` |

---

## 4. VNC/RDP 远程桌面接入方案

### 4.1 方案选型

| 方案 | 协议 | 客户端要求 | 适用场景 | 推荐度 |
|------|------|-----------|---------|--------|
| VNC + noVNC | VNC over WebSocket | 浏览器 | Linux GPU 开发环境 | ★★★★★ |
| xRDP | RDP | RDP 客户端/浏览器(Guacamole) | Windows 或需要 RDP 的场景 | ★★★☆☆ |
| Guacamole | SSH/VNC/RDP over HTTP | 浏览器 | 统一 Web 网关 | ★★★★☆ |

**推荐方案：VNC + noVNC 为主，Guacamole 为可选增强**

**理由：**
1. GPU 开发环境以 Linux 为主，VNC 是 Linux 远程桌面的标准方案
2. noVNC 提供纯浏览器访问，无需安装客户端，用户体验好
3. 现有 Host 实体已有 `VNCURL` 和 `VNCPassword` 字段，说明 VNC 是既定方向
4. Guacamole 可作为统一网关层叠加，但引入额外部署复杂度，建议作为第二阶段

### 4.2 VNC + noVNC 架构

```
用户浏览器
    │
    │ WebSocket (wss://)
    │
┌───▼──────────┐
│  Nginx/网关   │  反向代理 + SSL 终止
└───┬──────────┘
    │
    │ WebSocket (ws://)
    │
┌───▼──────────┐
│   noVNC      │  WebSocket → VNC 协议转换
│ (websockify) │  运行在容器内部
└───┬──────────┘
    │
    │ VNC 协议 (localhost)
    │
┌───▼──────────┐
│  TigerVNC    │  VNC 服务器
│   Server     │  运行在容器内部
└───┬──────────┘
    │
┌───▼──────────┐
│  Xfce4 桌面   │  轻量级桌面环境
└──────────────┘
```

**容器内组件栈：**
1. **TigerVNC Server** — VNC 服务端，监听 `:5901`
2. **websockify (noVNC)** — WebSocket 代理，监听 `:6080`，转发到 VNC `:5901`
3. **Xfce4** — 默认桌面环境（轻量，适合远程）
4. **dbus-x11** — D-Bus 会话总线（桌面环境依赖）

### 4.3 VNC 容器镜像策略

**方案：预构建基础镜像 + 运行时配置注入**

提供两类预构建镜像：

| 镜像 | 基础 | 包含组件 | 用途 |
|------|------|---------|------|
| `remotegpu/desktop-cuda:12.x` | nvidia/cuda:12.x-devel-ubuntu22.04 | TigerVNC + noVNC + Xfce4 + CUDA | GPU 桌面开发 |
| `remotegpu/desktop-base:latest` | ubuntu:22.04 | TigerVNC + noVNC + Xfce4 | 非 GPU 桌面 |

运行时通过环境变量注入配置：
- `VNC_PASSWORD` — VNC 访问密码
- `VNC_RESOLUTION` — 分辨率（默认 1920x1080）
- `VNC_COLOR_DEPTH` — 色深（默认 24）
- `DISPLAY` — X11 显示号（默认 :1）

### 4.4 Guacamole 集成方案（第二阶段）

Guacamole 作为可选的统一 Web 网关，提供以下增强能力：
- 统一入口：SSH / VNC / RDP 通过同一个 Web 界面访问
- 会话录制：支持操作审计
- 多用户权限：细粒度的连接访问控制
- 剪贴板/文件传输：跨协议的文件交互

**部署方式：** 独立 Docker Compose 部署（guacd + guacamole-web + PostgreSQL）

**集成点：**
- 环境创建时，Backend 调用 Guacamole REST API 自动创建连接
- 环境销毁时，自动清理连接配置
- 前端嵌入 Guacamole iframe 或使用 guacamole-common-js 库

---

## 5. 端口池管理策略

### 5.1 端口分配规划

现有 `port_pool.go` 定义的端口范围：

| 服务类型 | 端口范围 | 容量 | 容器内端口 |
|---------|---------|------|-----------|
| SSH | 22000-22999 | 1000 | 22 |
| RDP | 33890-34889 | 1000 | 3389 |
| Jupyter | 8888-9887 | 1000 | 8888 |
| VNC | 5900-6899 | 1000 | 5901 |
| noVNC | 6080-7079 | 1000 | 6080 |
| 自定义 | 10000-19999 | 10000 | 用户指定 |

### 5.2 关键问题：多主机端口隔离

**现有问题：** 当前 `AllocatePort()` 查询所有 `status=active` 的端口，未按主机过滤。在多主机场景下，不同主机可以使用相同的宿主机端口，但当前实现会错误地认为端口已被占用。

**修复方案：** PortMapping 实体增加 `host_id` 字段，分配和查询时按主机隔离。

```go
// 修改后的查询逻辑
err := db.Model(&entity.PortMapping{}).
    Where("host_id = ? AND status = ?", hostID, "active").
    Pluck("external_port", &usedPorts).Error
```

### 5.3 PortMapping 实体定义（待新增）

```go
// PortMapping 端口映射实体
type PortMapping struct {
    ID           uint      `gorm:"primarykey" json:"id"`
    EnvID        string    `gorm:"type:varchar(64);not null;index" json:"env_id"`
    HostID       string    `gorm:"type:varchar(64);not null;index" json:"host_id"`
    ServiceType  string    `gorm:"type:varchar(32);not null" json:"service_type"`
    InternalPort int       `gorm:"not null" json:"internal_port"`
    ExternalPort int       `gorm:"not null" json:"external_port"`
    Protocol     string    `gorm:"type:varchar(10);default:'tcp'" json:"protocol"`
    Description  string    `gorm:"type:varchar(256)" json:"description"`
    Status       string    `gorm:"type:varchar(20);default:'active'" json:"status"`
    ReleasedAt   *time.Time `json:"released_at"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

**需要的数据库迁移脚本：** `backend/sql/NN_add_port_mappings.sql`

**索引建议：**
- `idx_port_mapping_host_status` — `(host_id, status)` 用于端口分配查询
- `idx_port_mapping_env` — `(env_id, status)` 用于环境端口查询
- `uk_port_mapping_host_port` — `(host_id, external_port, status)` 唯一约束，防止同主机端口冲突

---

## 6. 与现有 Agent 的集成方式

### 6.1 现有 Agent 架构回顾

当前 Agent 采用 **轮询-执行-上报** 模式：

```
Agent 启动
  │
  ├─→ Poller: 定期调用 ClaimTasks() 认领任务
  ├─→ Scheduler: 从队列取任务，交给 Executor
  ├─→ Executor: 执行 shell/python/script 命令
  ├─→ Syncer: 上报进度和结果
  └─→ Collector: 采集系统/GPU 指标，随心跳上报
```

**Agent 客户端接口（Backend → Agent）：**
- `agent.Client` 接口支持 HTTP 和 gRPC 两种协议
- 已有方法：`StopProcess`、`ResetSSH`、`SyncSSHKeys`、`MountDataset`、`ExecuteCommand`、`Ping`

### 6.2 集成方案：扩展 Agent 支持容器操作

**方案一（推荐）：通过 Task 系统下发容器操作**

利用现有 Task 模型，新增 `container_*` 类型任务。Agent 侧新增 `ContainerHandler` 处理容器相关任务。

```
Backend                          Agent
  │                                │
  │  创建 Task(type=container_*)   │
  │──────────────────────────────→ │
  │                                │ Poller 认领任务
  │                                │ Scheduler 调度
  │                                │ ContainerHandler 执行
  │                                │   └─→ DockerService.*()
  │  ←──────────────────────────── │
  │  上报结果(container_id, status) │
  │                                │
```

**Agent 侧改动：**

```
agent/internal/
├── handler/
│   ├── task.go              # 现有任务处理
│   └── container.go         # 新增：容器任务处理
├── docker/
│   └── docker_service.go    # 从 _wip/docker.go 迁移
└── models/
    └── task.go              # 扩展 Task 类型常量
```

**方案二（备选）：扩展 Agent Client 接口**

在 `agent.Client` 接口中新增容器操作方法，Backend 直接调用 Agent API。

```go
// 扩展 agent.Client 接口
type Client interface {
    // ... 现有方法 ...
    CreateContainer(ctx context.Context, req *CreateContainerRequest) (*CreateContainerResponse, error)
    StopContainer(ctx context.Context, req *StopContainerRequest) (*Response, error)
    DeleteContainer(ctx context.Context, req *DeleteContainerRequest) (*Response, error)
    InspectContainer(ctx context.Context, req *InspectContainerRequest) (*ContainerInfo, error)
}
```

**方案对比：**

| 维度 | 方案一（Task 系统） | 方案二（Agent Client） |
|------|-------------------|---------------------|
| 改动范围 | Agent 新增 handler | Agent 新增 API + Backend 新增 Client 方法 |
| 异步支持 | 天然异步（任务队列） | 需额外处理长时间操作 |
| 可靠性 | 任务有重试、租约机制 | 需自行实现重试 |
| 实时性 | 受轮询间隔影响 | 即时响应 |
| 推荐场景 | 容器创建/删除（耗时操作） | 容器状态查询（快速操作） |

**最终建议：混合使用**
- 耗时操作（创建、删除）通过 Task 系统下发
- 快速查询（状态、日志、统计）通过 Agent Client 直接调用

### 6.3 Agent 侧 DockerService 迁移

将 `_wip/docker.go` 迁移到 Agent 侧，调整如下：

| 原位置 | 目标位置 | 说明 |
|--------|---------|------|
| `backend/internal/service/_wip/docker.go` | `agent/internal/docker/docker_service.go` | 核心容器操作 |
| `backend/internal/service/_wip/vnc.go` | `agent/internal/docker/vnc_config.go` | VNC 配置生成 |
| `backend/internal/service/_wip/rdp.go` | `agent/internal/docker/rdp_config.go` | RDP 配置生成（可延后） |

Backend 侧保留的模块：

| 模块 | 保留位置 | 说明 |
|------|---------|------|
| `port_pool.go` | `backend/internal/service/environment/` | 端口分配由 Backend 集中管理 |
| `guacamole.go` | `backend/internal/service/remote/` | Guacamole 网关由 Backend 管理 |
| `dns.go` | `backend/internal/service/network/` | DNS 配置由 Backend 管理 |
| `firewall.go` | `backend/internal/service/network/` | 防火墙由 Backend 管理 |

---

## 7. 各模块改进建议与优先级排序

### 7.1 优先级总览

| 优先级 | 模块 | 改进内容 | 依赖 |
|--------|------|---------|------|
| P0 | PortPoolService | 补充 PortMapping 实体 + 迁移脚本 + 多主机隔离 | 无 |
| P0 | DockerService | 迁移到 Agent 侧 + 补全查询方法 | PortMapping 实体 |
| P1 | VNCService | 实现配置生成 + 启动脚本 + 预构建镜像 | DockerService |
| P1 | Agent 集成 | 新增 ContainerHandler + Agent Client 扩展 | DockerService |
| P2 | GuacamoleService | 实现 REST API 客户端 + 环境集成 | VNC 可用后 |
| P3 | RDPService | 实现 xrdp 配置生成 | VNC 完成后 |
| P3 | DNSService | 实现 Cloudflare 适配 | 有域名需求时 |
| P3 | FirewallService | 实现 iptables 适配 | 有公网映射需求时 |

### 7.2 P0：PortPoolService 改进

**改进项：**

1. **新增 PortMapping 实体** — 在 `backend/internal/model/entity/` 下创建，包含 `host_id` 字段
2. **数据库迁移脚本** — 创建 `port_mappings` 表，含复合唯一索引
3. **多主机端口隔离** — `AllocatePort()` 和 `AllocatePorts()` 增加 `hostID` 参数
4. **端口冲突检测** — 分配前检查宿主机端口是否被非本系统占用（通过 Agent 探测）
5. **端口范围可配置** — 从硬编码改为从 `system_config` 表读取，支持运行时调整

### 7.3 P0：DockerService 改进

**改进项：**

1. **迁移到 Agent 侧** — 将 `docker.go` 移至 `agent/internal/docker/`，调整 import 路径
2. **补全查询方法** — 实现 `GetContainer()`、`ListContainers()`、`GetContainerLogs()`、`GetContainerStats()`
3. **实现 ExecCommand()** — 使用 Docker Exec API，支持交互式和非交互式模式
4. **实现环境级方法** — `CreateEnvironmentContainer()` 整合端口映射、GPU 配置、存储挂载
5. **GPU 高级模式** — 实现 `configureGPU()` 支持独占/vGPU/MIG/共享四种模式
6. **健康检查** — 容器创建时配置 Docker healthcheck，Agent 定期上报容器健康状态
7. **错误恢复** — 创建失败时自动清理残留资源（已部分实现，需完善）

### 7.4 P1：VNCService 改进

**改进项：**

1. **实现 GenerateVNCConfig()** — 从 Environment 实体读取配置，生成密码，设置默认值
2. **实现 GenerateVNCStartupScript()** — 生成容器启动脚本（设置密码 → 启动 VNC → 启动 noVNC → 启动桌面）
3. **预构建 Docker 镜像** — 创建 `remotegpu/desktop-cuda` 和 `remotegpu/desktop-base` 镜像
4. **GPU 渲染支持** — 配置 VirtualGL 或 EGL，使 GPU 加速的图形应用可在 VNC 中显示
5. **分辨率动态调整** — 支持通过 API 动态修改 VNC 分辨率

### 7.5 P1：Agent 集成改进

**改进项：**

1. **新增 ContainerHandler** — `agent/internal/handler/container.go`，根据 Task type 分发到 DockerService 对应方法
2. **扩展 Agent Client 接口** — 新增 `InspectContainer()`、`ContainerLogs()`、`ContainerStats()` 快速查询方法
3. **容器状态同步** — Agent 心跳中增加容器状态列表，Backend 据此更新 Environment 状态
4. **容器事件监听** — Agent 监听 Docker 事件流（container die/oom/restart），主动上报异常

### 7.6 P2：GuacamoleService 改进

**改进项：**

1. **实现 Guacamole REST API 客户端** — 认证（获取 token）、连接 CRUD、用户/权限管理
2. **环境自动配置** — 环境创建时自动注册 VNC/SSH 连接到 Guacamole
3. **前端集成** — 提供 Guacamole 连接 URL，前端通过 iframe 或 guacamole-common-js 嵌入
4. **会话管理** — 支持查看活跃会话、强制断开、会话录制回放

### 7.7 P3：RDPService / DNSService / FirewallService 改进

**RDPService：**
- 实现 xrdp 配置生成和启动脚本
- 构建包含 xrdp 的 Docker 镜像
- 优先级低，仅在有 Windows 桌面需求时推进

**DNSService：**
- 优先实现 Cloudflare 适配（国内外通用）
- 支持泛域名解析（`*.env.remotegpu.example.com`）
- 可通过 Nginx 反向代理替代，降低实现优先级

**FirewallService：**
- 优先实现 iptables 适配（覆盖大部分 Linux 主机）
- 在 Agent 直连模式下，可通过 Agent `ExecuteCommand` 执行 iptables 命令替代

---

## 8. 实施路线图

### 第一阶段：基础设施（P0）

**目标：** 打通容器创建和管理的完整链路

1. 创建 `PortMapping` 实体和数据库迁移脚本
2. 修复 `PortPoolService` 多主机端口隔离问题
3. 将 `DockerService` 迁移到 Agent 侧
4. 补全 `DockerService` 查询方法
5. Agent 新增 `ContainerHandler` 处理容器任务
6. Backend 新增环境容器管理 API

### 第二阶段：远程桌面（P1）

**目标：** 用户可通过浏览器访问 GPU 桌面环境

1. 构建 `remotegpu/desktop-cuda` 预构建镜像（TigerVNC + noVNC + Xfce4 + CUDA）
2. 实现 `VNCService` 配置生成和启动脚本
3. 环境创建流程集成 VNC 端口分配
4. 前端新增 noVNC 连接入口（iframe 嵌入）
5. Agent Client 扩展容器快速查询接口

### 第三阶段：增强功能（P2）

**目标：** 统一 Web 网关 + 高级 GPU 支持

1. 部署 Apache Guacamole（Docker Compose）
2. 实现 `GuacamoleService` REST API 客户端
3. 环境创建自动注册 Guacamole 连接
4. 实现 GPU 高级模式（vGPU/MIG）
5. 实现分布式存储挂载（NFS/JuiceFS）

### 第四阶段：网络增强（P3）

**目标：** 自动化网络配置

1. 实现 `DNSService` Cloudflare 适配
2. 实现 `FirewallService` iptables 适配
3. 实现 `RDPService` xrdp 配置（按需）
4. Nginx 反向代理自动配置

---

## 9. 总结

### 关键决策

1. **容器操作在 Agent 侧执行**，复用现有任务框架，避免 Backend 直连 Docker daemon
2. **VNC + noVNC 为主要远程桌面方案**，Guacamole 作为可选增强
3. **端口池由 Backend 集中管理**，需修复多主机隔离问题
4. **混合通信模式**：耗时操作走 Task 系统，快速查询走 Agent Client 直调

### _wip 代码处置建议

| 文件 | 处置 | 目标位置 |
|------|------|---------|
| `docker.go` | 迁移 | `agent/internal/docker/` |
| `vnc.go` | 迁移 | `agent/internal/docker/` |
| `rdp.go` | 暂留 | 待 P3 阶段迁移 |
| `guacamole.go` | 迁移 | `backend/internal/service/remote/` |
| `port_pool.go` | 迁移 | `backend/internal/service/environment/` |
| `dns.go` | 暂留 | 待 P3 阶段迁移 |
| `firewall.go` | 暂留 | 待 P3 阶段迁移 |

迁移完成后，`_wip/` 目录可删除。
