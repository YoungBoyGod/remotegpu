# RemoteGPU 系统架构设计

> 完整的系统模块划分与架构设计
>
> 创建日期：2026-01-26
>
> 支持：传统架构 + Kubernetes 架构

---

## 目录

1. [整体架构](#1-整体架构)
2. [模块划分](#2-模块划分)
3. [基础设施层](#3-基础设施层)
4. [核心服务层](#4-核心服务层)
5. [业务功能层](#5-业务功能层)
6. [前端展示层](#6-前端展示层)
7. [开发路线图](#7-开发路线图)

---

## 1. 整体架构

### 1.1 系统分层架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        前端展示层                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ Web 控制台│  │ 管理后台  │  │ Web 终端 │  │ Web RDP  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTPS / WebSocket
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      业务功能层（API 服务）                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 用户管理  │  │ 环境管理  │  │ 数据管理  │  │ 镜像管理  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 训练管理  │  │ 推理服务  │  │ 计费管理  │  │ 监控告警  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└────────────────────────┬────────────────────────────────────────┘
                         │ gRPC / HTTP
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      核心服务层                                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 调度器    │  │ 资源管理  │  │ 端口管理  │  │ 存储管理  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└────────────────────────┬────────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ 传统架构层    │  │ Kubernetes层  │  │ 混合云层      │
├──────────────┤  ├──────────────┤  ├──────────────┤
│ Linux Agent  │  │ K8s API      │  │ 多云适配器    │
│ Windows Agent│  │ GPU Operator │  │              │
│ Docker       │  │ Helm Charts  │  │              │
└──────────────┘  └──────────────┘  └──────────────┘
        │                │                │
        └────────────────┼────────────────┘
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      基础设施层                                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 数据库    │  │ 对象存储  │  │ 消息队列  │  │ 缓存      │        │
│  │PostgreSQL│  │ MinIO    │  │ RabbitMQ │  │ Redis    │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 监控      │  │ 日志      │  │ 镜像仓库  │  │ 网络      │        │
│  │Prometheus│  │ ELK      │  │ Harbor   │  │ Nginx    │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 技术栈总览

```yaml
前端技术栈:
  框架: React 18 + TypeScript
  状态管理: Redux Toolkit / Zustand
  UI 库: Ant Design
  图表: ECharts
  编辑器: Monaco Editor
  终端: xterm.js
  构建: Vite

后端技术栈:
  主语言: Go 1.21+
  框架: Gin (HTTP) / gRPC
  ORM: GORM
  数据库: PostgreSQL 14+
  缓存: Redis 7+
  消息队列: RabbitMQ / Kafka

基础设施:
  容器编排: Kubernetes 1.28+ / Docker
  GPU 管理: NVIDIA GPU Operator
  对象存储: MinIO / S3
  镜像仓库: Harbor / Registry
  监控: Prometheus + Grafana
  日志: ELK Stack
  负载均衡: Nginx / Traefik
```

---

## 2. 模块划分

### 2.1 核心模块列表

以下模块编号与 `docs/design/module_division.md` 保持一致。

| 模块编号 | 模块名称 | 核心职责 |
|---------|---------|---------|
| 模块 1 | CMDB 设备管理模块 | 主机/GPU 管理、设备生命周期 |
| 模块 2 | 用户与权限模块 | 认证、工作空间、RBAC |
| 模块 3 | 环境管理模块 | 开发环境创建、SSH/RDP 访问 |
| 模块 4 | 资源调度模块 | 统一调度、端口管理、资源分配 |
| 模块 5 | 数据与存储模块 | 数据集管理、对象存储 |
| 模块 6 | 镜像管理模块 | 官方镜像、自定义镜像构建 |
| 模块 7 | 训练与推理模块 | 训练任务、推理服务 |
| 模块 8 | 计费管理模块 | 计费规则、账单管理 |
| 模块 9 | 监控告警模块 | 资源监控、告警 |
| 模块 10 | 网关与认证模块 | 统一入口、认证、限流 |
| 模块 11 | 制品管理模块 | 制品仓库、版本管理 |
| 模块 12 | 问题单管理模块 | 问题跟踪、工单流转 |
| 模块 13 | 需求单管理模块 | 需求收集、评审、跟踪 |
| 模块 14 | 通知管理模块 | 多渠道通知、消息推送 |
| 模块 15 | Webhook 管理模块 | 事件回调、第三方集成 |

### 2.2 模块依赖关系

```
┌─────────────────────────────────────────────────────────────┐
│                     前端模块                                  │
│  Web Console + Admin Panel + Web Terminal                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   API Gateway                               │
│              (统一入口、认证、限流)                           │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ 用户模块      │  │ 环境模块      │  │ 数据模块      │
│              │  │              │  │              │
│ 依赖:        │  │ 依赖:        │  │ 依赖:        │
│ - 基础设施    │  │ - 用户模块    │  │ - 用户模块    │
│              │  │ - 调度模块    │  │ - 存储服务    │
│              │  │ - 设备模块    │  │              │
└──────────────┘  └──────────────┘  └──────────────┘
        │                │                │
        └────────────────┼────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   核心服务层                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │ 调度器    │  │ 设备管理  │  │ 存储管理  │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   基础设施层                                  │
│  Database + Storage + Cache + MQ + Monitoring              │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 基础设施层

### 3.1 模块：数据库服务

**职责：** 存储所有元数据和业务数据

**技术选型：** PostgreSQL 14+

**核心表设计：**

```sql
-- 用户相关
- users (用户表)
- workspaces (工作空间表)
- workspace_members (成员表)

-- 设备相关
- hosts (主机表)
- gpus (GPU 设备表)
- host_metrics (主机监控数据)

-- 环境相关
- environments (环境表)
- port_mappings (端口映射表)
- ssh_credentials (SSH 凭证表)

-- 数据相关
- datasets (数据集表)
- dataset_versions (数据集版本表)
- models (模型表)
- model_versions (模型版本表)

-- 镜像相关
- images (镜像表)
- custom_images (自定义镜像表)
- image_builds (构建历史表)

-- 任务相关
- training_jobs (训练任务表)
- inference_services (推理服务表)

-- 计费相关
- billing_records (计费记录表)
- invoices (账单表)
- payments (支付记录表)
```

**注意事项：**
- ✅ 使用连接池（最大连接数 100）
- ✅ 定期备份（每日全量 + 每小时增量）
- ✅ 读写分离（主从复制）
- ✅ 索引优化（常用查询字段）
- ✅ 分区表（大表按时间分区）

**部署方案：**
```yaml
# docker-compose.yml
postgres:
  image: postgres:14
  environment:
    POSTGRES_DB: remotegpu
    POSTGRES_USER: admin
    POSTGRES_PASSWORD: ${DB_PASSWORD}
  volumes:
    - postgres-data:/var/lib/postgresql/data
  ports:
    - "5432:5432"
```

---

### 3.2 模块：对象存储服务

**职责：** 存储数据集、模型、训练输出等大文件

**技术选型：** MinIO（开源 S3 兼容）

**存储桶规划：**
```
remotegpu/
├── datasets/          # 数据集
├── models/            # 模型文件
├── artifacts/         # 训练输出
├── images/            # 自定义镜像构建产物
└── backups/           # 备份文件
```

**注意事项：**
- ✅ 启用版本控制
- ✅ 配置生命周期策略（自动归档/删除）
- ✅ 启用加密（静态加密）
- ✅ 配置访问策略（IAM）
- ✅ 监控存储使用量

**部署方案：**
```yaml
# 单机部署
minio:
  image: minio/minio
  command: server /data --console-address ":9001"
  environment:
    MINIO_ROOT_USER: admin
    MINIO_ROOT_PASSWORD: ${MINIO_PASSWORD}
  volumes:
    - minio-data:/data
  ports:
    - "9000:9000"
    - "9001:9001"

# 分布式部署（生产）
# 4 节点，每节点 2 块盘
minio1:
  command: server http://minio{1...4}/data{1...2}
```

---

### 3.3 模块：缓存服务

**职责：** 缓存热点数据、会话管理、分布式锁

**技术选型：** Redis 7+

**使用场景：**
```
1. 会话存储 (Session)
   - Key: session:{token}
   - TTL: 24h

2. 用户信息缓存
   - Key: user:{user_id}
   - TTL: 1h

3. 主机状态缓存
   - Key: host:{host_id}:status
   - TTL: 1m

4. 端口分配锁
   - Key: lock:port:{port}
   - TTL: 10s

5. 任务队列
   - List: queue:training
   - List: queue:building
```

**注意事项：**
- ✅ 启用持久化（AOF + RDB）
- ✅ 配置主从复制
- ✅ 设置合理的 TTL
- ✅ 监控内存使用
- ✅ 配置淘汰策略（LRU）

**部署方案：**
```yaml
redis:
  image: redis:7-alpine
  command: redis-server --appendonly yes
  volumes:
    - redis-data:/data
  ports:
    - "6379:6379"
```

---

### 3.4 模块：消息队列服务

**职责：** 异步任务处理、事件通知

**技术选型：** RabbitMQ

**队列设计：**
```
1. 环境创建队列
   - Queue: env.create
   - Consumer: Environment Service

2. 镜像构建队列
   - Queue: image.build
   - Consumer: Image Builder

3. 训练任务队列
   - Queue: training.submit
   - Consumer: Training Service

4. 通知队列
   - Queue: notification
   - Consumer: Notification Service
```

**注意事项：**
- ✅ 配置死信队列（DLQ）
- ✅ 设置消息持久化
- ✅ 配置消息确认机制
- ✅ 限制队列长度
- ✅ 监控队列积压

**部署方案：**
```yaml
rabbitmq:
  image: rabbitmq:3-management
  environment:
    RABBITMQ_DEFAULT_USER: admin
    RABBITMQ_DEFAULT_PASS: ${MQ_PASSWORD}
  ports:
    - "5672:5672"
    - "15672:15672"
```

---

### 3.5 模块：监控日志服务

**职责：** 系统监控、日志收集、告警

**技术选型：** Prometheus + Grafana + ELK

**监控指标：**
```
系统指标:
- CPU 使用率
- 内存使用率
- 磁盘使用率
- 网络流量

GPU 指标:
- GPU 使用率
- GPU 显存使用
- GPU 温度
- GPU 功耗

业务指标:
- API 请求量
- API 响应时间
- 环境创建成功率
- 任务队列长度
```

**日志收集：**
```
应用日志:
- API 服务日志
- Worker 日志
- 任务执行日志

系统日志:
- 主机系统日志
- 容器日志
- 审计日志
```

**注意事项：**
- ✅ 配置日志轮转
- ✅ 设置日志保留期限
- ✅ 配置告警规则
- ✅ 日志脱敏（敏感信息）
- ✅ 分级存储（热/冷数据）

---

## 4. 设备管理模块 (Device Management)

### 4.1 模块：主机管理 (Host Management)

**职责：** 管理 Linux 和 Windows 物理机/虚拟机

**支持的主机类型：**
```
1. Linux 主机
   - Ubuntu 20.04/22.04
   - CentOS 7/8
   - 运行 Docker 或作为 K8s Node

2. Windows 主机
   - Windows Server 2019/2022
   - 支持 Hyper-V
   - 支持 Windows Server 容器运行时（如 Mirantis Container Runtime）

3. Kubernetes 节点
   - K8s Worker Node
   - GPU Operator 管理
```

**数据库设计：**

```sql
-- 主机表
CREATE TABLE hosts (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    hostname VARCHAR(256),
    ip_address VARCHAR(64) NOT NULL,
    public_ip VARCHAR(64),

    -- 主机类型
    os_type VARCHAR(20) NOT NULL,              -- linux, windows
    os_version VARCHAR(64),                    -- Ubuntu 20.04, Windows Server 2022
    arch VARCHAR(20) DEFAULT 'x86_64',         -- x86_64, arm64

    -- 部署模式
    deployment_mode VARCHAR(20) NOT NULL,      -- traditional, kubernetes, hybrid
    k8s_node_name VARCHAR(128),                -- K8s 节点名称

    -- 状态
    status VARCHAR(20) DEFAULT 'offline',      -- online, offline, maintenance, error
    health_status VARCHAR(20) DEFAULT 'unknown', -- healthy, degraded, unhealthy

    -- 资源信息
    total_cpu INT NOT NULL,                    -- CPU 核心数
    total_memory BIGINT NOT NULL,              -- 内存（字节）
    total_disk BIGINT,                         -- 磁盘（字节）
    total_gpu INT DEFAULT 0,                   -- GPU 数量

    used_cpu INT DEFAULT 0,
    used_memory BIGINT DEFAULT 0,
    used_disk BIGINT DEFAULT 0,
    used_gpu INT DEFAULT 0,

    -- 网络信息
    ssh_port INT DEFAULT 22,
    winrm_port INT,                            -- Windows WinRM 端口
    agent_port INT DEFAULT 8080,               -- Agent 监听端口

    -- 标签和元数据
    labels JSONB,                              -- {"region": "us-west", "zone": "a"}
    tags TEXT[],                               -- ["gpu-server", "high-memory"]

    -- 时间戳
    last_heartbeat TIMESTAMP,
    registered_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    INDEX idx_os_type (os_type),
    INDEX idx_status (status),
    INDEX idx_deployment_mode (deployment_mode),
    INDEX idx_labels (labels) USING GIN
);

-- GPU 设备表
CREATE TABLE gpus (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(64) NOT NULL,
    gpu_index INT NOT NULL,                    -- 主机上的 GPU 索引 (0, 1, 2...)

    -- GPU 信息
    uuid VARCHAR(128) UNIQUE,                  -- NVIDIA GPU UUID
    name VARCHAR(128),                         -- Tesla V100, RTX 4090
    brand VARCHAR(64),                         -- NVIDIA, AMD
    architecture VARCHAR(64),                  -- Ampere, Turing

    -- 规格
    memory_total BIGINT,                       -- 显存总量（字节）
    cuda_cores INT,                            -- CUDA 核心数
    compute_capability VARCHAR(32),            -- 7.5, 8.0, 8.6

    -- 状态
    status VARCHAR(20) DEFAULT 'available',    -- available, allocated, maintenance, error
    health_status VARCHAR(20) DEFAULT 'healthy',

    -- 分配信息
    allocated_to VARCHAR(64),                  -- 环境 ID
    allocated_at TIMESTAMP,

    -- 性能信息
    power_limit INT,                           -- 功耗限制（瓦）
    temperature_limit INT,                     -- 温度限制（摄氏度）

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE,
    UNIQUE(host_id, gpu_index),
    INDEX idx_status (status),
    INDEX idx_allocated_to (allocated_to)
);

-- 主机监控数据表
CREATE TABLE host_metrics (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(64) NOT NULL,

    -- CPU 指标
    cpu_usage_percent FLOAT,                   -- CPU 使用率
    cpu_load_1m FLOAT,                         -- 1分钟负载
    cpu_load_5m FLOAT,
    cpu_load_15m FLOAT,

    -- 内存指标
    memory_used BIGINT,
    memory_available BIGINT,
    memory_usage_percent FLOAT,

    -- 磁盘指标
    disk_used BIGINT,
    disk_available BIGINT,
    disk_usage_percent FLOAT,
    disk_io_read_bytes BIGINT,
    disk_io_write_bytes BIGINT,

    -- 网络指标
    network_rx_bytes BIGINT,                   -- 接收字节数
    network_tx_bytes BIGINT,                   -- 发送字节数
    network_rx_packets BIGINT,
    network_tx_packets BIGINT,

    -- GPU 指标（聚合）
    gpu_avg_utilization FLOAT,
    gpu_avg_memory_used BIGINT,
    gpu_avg_temperature FLOAT,
    gpu_avg_power FLOAT,

    collected_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE,
    INDEX idx_host_id_time (host_id, collected_at DESC)
);

-- GPU 监控数据表
CREATE TABLE gpu_metrics (
    id BIGSERIAL PRIMARY KEY,
    gpu_id BIGINT NOT NULL,
    host_id VARCHAR(64) NOT NULL,

    -- 使用率
    utilization_percent FLOAT,                 -- GPU 使用率
    memory_used BIGINT,                        -- 显存使用（字节）
    memory_usage_percent FLOAT,

    -- 温度和功耗
    temperature FLOAT,                         -- 温度（摄氏度）
    power_draw FLOAT,                          -- 功耗（瓦）
    fan_speed_percent FLOAT,                   -- 风扇转速

    -- 性能
    sm_clock INT,                              -- SM 时钟频率（MHz）
    memory_clock INT,                          -- 显存时钟频率（MHz）

    -- 进程信息
    process_count INT,                         -- 运行的进程数

    collected_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (gpu_id) REFERENCES gpus(id) ON DELETE CASCADE,
    INDEX idx_gpu_id_time (gpu_id, collected_at DESC),
    INDEX idx_host_id_time (host_id, collected_at DESC)
);
```

**主机注册流程：**

```go
// device/host_manager.go
package device

type HostManager struct {
    db *gorm.DB
}

// 注册主机
func (m *HostManager) RegisterHost(req RegisterHostRequest) (*Host, error) {
    // 1. 收集主机信息
    hostInfo := m.collectHostInfo(req.IPAddress)

    // 2. 检测 GPU
    gpus := m.detectGPUs(req.IPAddress)

    // 3. 创建主机记录
    host := &Host{
        ID:             GenerateHostID(),
        Name:           req.Name,
        IPAddress:      req.IPAddress,
        OSType:         hostInfo.OSType,
        OSVersion:      hostInfo.OSVersion,
        DeploymentMode: req.DeploymentMode,
        TotalCPU:       hostInfo.CPUCores,
        TotalMemory:    hostInfo.TotalMemory,
        TotalGPU:       len(gpus),
        Status:         "online",
    }

    if err := m.db.Create(host).Error; err != nil {
        return nil, err
    }

    // 4. 创建 GPU 记录
    for i, gpu := range gpus {
        gpuRecord := &GPU{
            HostID:    host.ID,
            GPUIndex:  i,
            UUID:      gpu.UUID,
            Name:      gpu.Name,
            Brand:     gpu.Brand,
            MemoryTotal: gpu.Memory,
            Status:    "available",
        }
        m.db.Create(gpuRecord)
    }

    return host, nil
}

// 收集主机信息（通过 SSH/WinRM）
func (m *HostManager) collectHostInfo(ipAddress string) HostInfo {
    // Linux: 通过 SSH 执行命令
    // Windows: 通过 WinRM 执行 PowerShell

    // 示例：获取 CPU 核心数
    // Linux: nproc
    // Windows: (Get-WmiObject Win32_Processor).NumberOfLogicalProcessors

    return HostInfo{
        OSType:      "linux",
        OSVersion:   "Ubuntu 20.04",
        CPUCores:    32,
        TotalMemory: 128 * 1024 * 1024 * 1024, // 128GB
    }
}

// 检测 GPU（使用 nvidia-smi）
func (m *HostManager) detectGPUs(ipAddress string) []GPUInfo {
    // 执行: nvidia-smi --query-gpu=uuid,name,memory.total --format=csv

    return []GPUInfo{
        {
            UUID:   "GPU-12345678-1234-1234-1234-123456789012",
            Name:   "Tesla V100",
            Brand:  "NVIDIA",
            Memory: 32 * 1024 * 1024 * 1024, // 32GB
        },
    }
}
```

**注意事项：**
- ✅ 主机注册前验证网络连通性
- ✅ 自动检测主机类型（Linux/Windows）
- ✅ 自动发现 GPU 设备
- ✅ 支持标签和分组管理
- ✅ 定期心跳检测（30秒）

---

### 4.2 模块：GPU 资源管理

**职责：** GPU 设备发现、分配、监控

**GPU 发现方式：**

```bash
# 方式 1：nvidia-smi（传统架构）
nvidia-smi --query-gpu=index,uuid,name,memory.total,compute_cap \
    --format=csv,noheader

# 输出示例：
# 0, GPU-12345678-1234-1234-1234-123456789012, Tesla V100-SXM2-32GB, 32510 MiB, 7.0
# 1, GPU-87654321-4321-4321-4321-210987654321, Tesla V100-SXM2-32GB, 32510 MiB, 7.0

# 方式 2：DCGM（NVIDIA Data Center GPU Manager）
dcgmi discovery -l

# 方式 3：Kubernetes GPU Operator
kubectl get nodes -o json | jq '.items[].status.capacity."nvidia.com/gpu"'
```

**GPU 分配策略：**

```go
// device/gpu_allocator.go
package device

type GPUAllocator struct {
    db *gorm.DB
}

// GPU 分配策略
type AllocationStrategy string

const (
    StrategyFirstFit    AllocationStrategy = "first_fit"    // 首次适配
    StrategyBestFit     AllocationStrategy = "best_fit"     // 最佳适配
    StrategyLeastLoaded AllocationStrategy = "least_loaded" // 最少负载
    StrategyRoundRobin  AllocationStrategy = "round_robin"  // 轮询
)

// 分配 GPU
func (a *GPUAllocator) AllocateGPU(req GPUAllocationRequest) (*GPU, error) {
    // 1. 查询可用 GPU
    var availableGPUs []GPU
    query := a.db.Where("status = ?", "available")

    // 筛选条件
    if req.MinMemory > 0 {
        query = query.Where("memory_total >= ?", req.MinMemory)
    }
    if req.GPUModel != "" {
        query = query.Where("name LIKE ?", "%"+req.GPUModel+"%")
    }
    if req.HostID != "" {
        query = query.Where("host_id = ?", req.HostID)
    }

    query.Find(&availableGPUs)

    if len(availableGPUs) == 0 {
        return nil, errors.New("no available GPU")
    }

    // 2. 根据策略选择 GPU
    var selectedGPU *GPU
    switch req.Strategy {
    case StrategyFirstFit:
        selectedGPU = &availableGPUs[0]
    case StrategyLeastLoaded:
        selectedGPU = a.selectLeastLoadedGPU(availableGPUs)
    default:
        selectedGPU = &availableGPUs[0]
    }

    // 3. 标记为已分配
    a.db.Model(selectedGPU).Updates(map[string]interface{}{
        "status":       "allocated",
        "allocated_to": req.EnvID,
        "allocated_at": time.Now(),
    })

    return selectedGPU, nil
}

// 选择负载最低的 GPU
func (a *GPUAllocator) selectLeastLoadedGPU(gpus []GPU) *GPU {
    var minLoad float64 = 100
    var selected *GPU

    for i := range gpus {
        // 查询最近的 GPU 使用率
        var metric GPUMetric
        a.db.Where("gpu_id = ?", gpus[i].ID).
            Order("collected_at DESC").
            First(&metric)

        if metric.UtilizationPercent < minLoad {
            minLoad = metric.UtilizationPercent
            selected = &gpus[i]
        }
    }

    return selected
}

// 释放 GPU
func (a *GPUAllocator) ReleaseGPU(gpuID int64) error {
    return a.db.Model(&GPU{}).Where("id = ?", gpuID).Updates(map[string]interface{}{
        "status":       "available",
        "allocated_to": nil,
        "allocated_at": nil,
    }).Error
}
```

**GPU 监控实现：**

```go
// device/gpu_monitor.go
package device

// GPU 监控器
type GPUMonitor struct {
    db *gorm.DB
}

// 收集 GPU 指标
func (m *GPUMonitor) CollectMetrics(hostID string) error {
    // 1. 获取主机上的所有 GPU
    var gpus []GPU
    m.db.Where("host_id = ?", hostID).Find(&gpus)

    // 2. 执行 nvidia-smi 获取实时数据
    metrics := m.queryNvidiaSMI(hostID)

    // 3. 保存到数据库
    for _, metric := range metrics {
        m.db.Create(&GPUMetric{
            GPUID:              metric.GPUID,
            HostID:             hostID,
            UtilizationPercent: metric.Utilization,
            MemoryUsed:         metric.MemoryUsed,
            MemoryUsagePercent: metric.MemoryUsagePercent,
            Temperature:        metric.Temperature,
            PowerDraw:          metric.PowerDraw,
            CollectedAt:        time.Now(),
        })
    }

    return nil
}

// 查询 nvidia-smi
func (m *GPUMonitor) queryNvidiaSMI(hostID string) []GPUMetricData {
    // 执行命令：
    // nvidia-smi --query-gpu=index,utilization.gpu,memory.used,memory.total,temperature.gpu,power.draw \
    //     --format=csv,noheader,nounits

    // 示例输出：
    // 0, 45, 8192, 32768, 65, 180
    // 1, 30, 4096, 32768, 58, 150

    return []GPUMetricData{
        {
            GPUID:              1,
            Utilization:        45.0,
            MemoryUsed:         8192 * 1024 * 1024,
            MemoryUsagePercent: 25.0,
            Temperature:        65.0,
            PowerDraw:          180.0,
        },
    }
}

// 启动定期监控
func (m *GPUMonitor) StartMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            // 获取所有在线主机
            var hosts []Host
            m.db.Where("status = ?", "online").Find(&hosts)

            // 并发收集指标
            for _, host := range hosts {
                go m.CollectMetrics(host.ID)
            }
        }
    }()
}
```

**注意事项：**
- ✅ 监控频率：每 30 秒收集一次
- ✅ 数据保留：7 天详细数据 + 90 天聚合数据
- ✅ 告警阈值：温度 > 85°C、使用率 > 95%
- ✅ 支持 DCGM 集成（更精确的监控）

---

### 4.3 模块：设备健康检查

**职责：** 主机和 GPU 健康状态监测

**健康检查项：**

```yaml
主机健康检查:
  - 心跳检测（30秒超时）
  - SSH/WinRM 连通性
  - 磁盘空间（> 10% 可用）
  - 内存可用（> 5% 可用）
  - 系统负载（< 80%）

GPU 健康检查:
  - GPU 可访问性（nvidia-smi 响应）
  - 温度正常（< 85°C）
  - 功耗正常（< 额定功率）
  - ECC 错误检查
  - 驱动版本检查
```

**实现代码：**

```go
// device/health_checker.go
package device

type HealthChecker struct {
    db *gorm.DB
}

// 健康检查结果
type HealthCheckResult struct {
    HostID       string
    Status       string // healthy, degraded, unhealthy
    Issues       []string
    CheckedAt    time.Time
}

// 执行健康检查
func (h *HealthChecker) CheckHost(hostID string) (*HealthCheckResult, error) {
    var host Host
    if err := h.db.Where("id = ?", hostID).First(&host).Error; err != nil {
        return nil, err
    }

    result := &HealthCheckResult{
        HostID:    hostID,
        Status:    "healthy",
        Issues:    []string{},
        CheckedAt: time.Now(),
    }

    // 1. 检查心跳
    if time.Since(host.LastHeartbeat) > 60*time.Second {
        result.Issues = append(result.Issues, "心跳超时")
        result.Status = "unhealthy"
    }

    // 2. 检查磁盘空间
    var metrics HostMetric
    h.db.Where("host_id = ?", hostID).Order("collected_at DESC").First(&metrics)

    if metrics.DiskUsagePercent > 90 {
        result.Issues = append(result.Issues, "磁盘空间不足")
        result.Status = "degraded"
    }

    // 3. 检查 GPU 健康
    var gpus []GPU
    h.db.Where("host_id = ?", hostID).Find(&gpus)

    for _, gpu := range gpus {
        if !h.checkGPUHealth(gpu.ID) {
            result.Issues = append(result.Issues, fmt.Sprintf("GPU %d 异常", gpu.GPUIndex))
            result.Status = "degraded"
        }
    }

    // 4. 更新主机健康状态
    h.db.Model(&host).Update("health_status", result.Status)

    return result, nil
}

// 检查 GPU 健康
func (h *HealthChecker) checkGPUHealth(gpuID int64) bool {
    var metric GPUMetric
    h.db.Where("gpu_id = ?", gpuID).Order("collected_at DESC").First(&metric)

    // 检查温度
    if metric.Temperature > 85 {
        return false
    }

    // 检查是否有数据（GPU 可访问）
    if time.Since(metric.CollectedAt) > 5*time.Minute {
        return false
    }

    return true
}

// 定期健康检查
func (h *HealthChecker) StartHealthCheck(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            var hosts []Host
            h.db.Find(&hosts)

            for _, host := range hosts {
                result, _ := h.CheckHost(host.ID)

                // 如果状态变化，发送告警
                if result.Status != "healthy" {
                    h.sendAlert(result)
                }
            }
        }
    }()
}

// 发送告警
func (h *HealthChecker) sendAlert(result *HealthCheckResult) {
    // 发送到告警系统
    alertManager.Send(Alert{
        Level:   "warning",
        Title:   fmt.Sprintf("主机 %s 健康检查异常", result.HostID),
        Message: strings.Join(result.Issues, ", "),
    })
}
```

**注意事项：**
- ✅ 健康检查频率：每 1 分钟
- ✅ 自动隔离不健康主机
- ✅ 告警通知（邮件/钉钉/企业微信）
- ✅ 健康历史记录

---

### 4.4 传统架构 vs Kubernetes 实现对比

**传统架构（Docker）：**

```go
// 主机 Agent 上报心跳和指标
type WorkerAgent struct {
    hostID string
    masterURL string
}

func (a *WorkerAgent) Start() {
    // 1. 注册到 Master
    a.register()

    // 2. 启动心跳
    go a.heartbeat()

    // 3. 启动指标收集
    go a.collectMetrics()

    // 4. 监听 Master 指令
    a.listenCommands()
}

func (a *WorkerAgent) heartbeat() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        http.Post(
            fmt.Sprintf("%s/api/hosts/%s/heartbeat", a.masterURL, a.hostID),
            "application/json",
            nil,
        )
    }
}
```

**Kubernetes 架构：**

```yaml
# 使用 DaemonSet 部署监控 Agent
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: gpu-monitor
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: gpu-monitor
  template:
    metadata:
      labels:
        app: gpu-monitor
    spec:
      hostNetwork: true
      hostPID: true
      containers:
      - name: monitor
        image: remotegpu/gpu-monitor:latest
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        volumeMounts:
        - name: nvidia
          mountPath: /usr/local/nvidia
      volumes:
      - name: nvidia
        hostPath:
          path: /usr/local/nvidia
```

**对比总结：**

| 维度 | 传统架构 | Kubernetes 架构 |
|------|---------|----------------|
| **部署方式** | 手动部署 Agent | DaemonSet 自动部署 |
| **服务发现** | 手动注册 | K8s Service Discovery |
| **健康检查** | 自定义心跳 | Liveness/Readiness Probe |
| **指标收集** | 自定义上报 | Prometheus + DCGM Exporter |
| **GPU 管理** | nvidia-smi | GPU Operator + Device Plugin |

---

## 5. 核心服务层 (Core Services)

### 5.1 模块：统一调度器 (Unified Scheduler)

**职责：** 统一调度 Linux 和 Windows 环境，支持传统架构和 K8s 架构

**调度策略：**

```yaml
调度算法:
  1. 资源匹配
     - CPU、内存、GPU 满足需求
     - 磁盘空间充足

  2. 负载均衡
     - 最少使用策略（Least Used）
     - 轮询策略（Round Robin）
     - 随机策略（Random）

  3. 亲和性调度
     - 节点亲和性（指定主机/区域）
     - GPU 型号亲和性（指定 GPU 型号）
     - 数据亲和性（数据集所在主机优先）

  4. 优先级调度
     - VIP 用户优先
     - 付费用户优先
     - 紧急任务优先
```

**数据库设计：**

```sql
-- 调度策略配置表
CREATE TABLE scheduler_policies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    strategy VARCHAR(64) NOT NULL,           -- least_used, round_robin, priority
    enabled BOOLEAN DEFAULT true,
    priority INT DEFAULT 0,                  -- 策略优先级
    config JSONB,                            -- 策略配置
    created_at TIMESTAMP DEFAULT NOW()
);

-- 调度历史表
CREATE TABLE scheduler_history (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    customer_id BIGINT NOT NULL,

    -- 请求信息
    requested_cpu INT,
    requested_memory BIGINT,
    requested_gpu INT,

    -- 调度结果
    selected_host_id VARCHAR(64),
    selected_gpus JSONB,                     -- [{"gpu_id": 1, "gpu_index": 0}]
    strategy_used VARCHAR(64),

    -- 调度耗时
    scheduling_duration_ms INT,

    -- 状态
    status VARCHAR(20),                      -- success, failed
    failure_reason TEXT,

    scheduled_at TIMESTAMP DEFAULT NOW(),

    INDEX idx_env_id (env_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_scheduled_at (scheduled_at DESC)
);
```

**调度器实现：**

```go
// scheduler/unified_scheduler.go
package scheduler

type UnifiedScheduler struct {
    db           *gorm.DB
    hostManager  *device.HostManager
    gpuAllocator *device.GPUAllocator
}

// 调度请求
type ScheduleRequest struct {
    EnvID      string
    CustomerID int64

    // 资源需求
    CPU    int
    Memory int64
    GPU    int

    // 约束条件
    OSType      string   // linux, windows
    GPUModel    string   // Tesla V100, RTX 4090
    HostID      string   // 指定主机（可选）
    Region      string   // 区域（可选）

    // 调度策略
    Strategy string // least_used, round_robin, priority
}

// 调度结果
type ScheduleResult struct {
    Host     *device.Host
    GPUs     []*device.GPU
    Strategy string
    Duration time.Duration
}

// 执行调度
func (s *UnifiedScheduler) Schedule(req ScheduleRequest) (*ScheduleResult, error) {
    startTime := time.Now()

    // 1. 查询可用主机
    hosts, err := s.findAvailableHosts(req)
    if err != nil {
        return nil, err
    }

    if len(hosts) == 0 {
        return nil, errors.New("no available host")
    }

    // 2. 根据策略选择主机
    var selectedHost *device.Host
    switch req.Strategy {
    case "least_used":
        selectedHost = s.selectLeastUsedHost(hosts)
    case "round_robin":
        selectedHost = s.selectRoundRobinHost(hosts)
    case "priority":
        selectedHost = s.selectPriorityHost(hosts, req.CustomerID)
    default:
        selectedHost = hosts[0]
    }

    // 3. 分配 GPU
    var allocatedGPUs []*device.GPU
    for i := 0; i < req.GPU; i++ {
        gpu, err := s.gpuAllocator.AllocateGPU(device.GPUAllocationRequest{
            EnvID:     req.EnvID,
            HostID:    selectedHost.ID,
            GPUModel:  req.GPUModel,
            Strategy:  device.StrategyLeastLoaded,
        })
        if err != nil {
            // 回滚已分配的 GPU
            s.rollbackGPUs(allocatedGPUs)
            return nil, err
        }
        allocatedGPUs = append(allocatedGPUs, gpu)
    }

    // 4. 更新主机资源使用
    s.updateHostResources(selectedHost.ID, req.CPU, req.Memory, req.GPU)

    // 5. 记录调度历史
    s.recordScheduleHistory(req, selectedHost, allocatedGPUs, time.Since(startTime))

    return &ScheduleResult{
        Host:     selectedHost,
        GPUs:     allocatedGPUs,
        Strategy: req.Strategy,
        Duration: time.Since(startTime),
    }, nil
}

// 查找可用主机
func (s *UnifiedScheduler) findAvailableHosts(req ScheduleRequest) ([]*device.Host, error) {
    var hosts []*device.Host

    query := s.db.Where("status = ? AND health_status = ?", "online", "healthy")

    // 操作系统类型
    if req.OSType != "" {
        query = query.Where("os_type = ?", req.OSType)
    }

    // 指定主机
    if req.HostID != "" {
        query = query.Where("id = ?", req.HostID)
    }

    // 资源过滤
    query = query.Where("total_cpu - used_cpu >= ?", req.CPU)
    query = query.Where("total_memory - used_memory >= ?", req.Memory)
    query = query.Where("total_gpu - used_gpu >= ?", req.GPU)

    query.Find(&hosts)

    return hosts, nil
}

// 选择负载最低的主机
func (s *UnifiedScheduler) selectLeastUsedHost(hosts []*device.Host) *device.Host {
    var minUsage float64 = 100
    var selected *device.Host

    for i := range hosts {
        usage := float64(hosts[i].UsedCPU) / float64(hosts[i].TotalCPU) * 100
        if usage < minUsage {
            minUsage = usage
            selected = &hosts[i]
        }
    }

    return selected
}
```

**注意事项：**
- ✅ 调度延迟 < 100ms
- ✅ 支持调度策略热更新
- ✅ 记录调度历史用于分析
- ✅ 支持调度失败重试

---

### 5.2 模块：端口管理 (Port Manager)

**职责：** 管理 SSH/RDP/JupyterLab 端口分配

**端口范围规划：**

```yaml
端口分配:
  SSH (Linux):     30000-31000  (1000 个端口)
  SSH (Windows):   31000-32000  (1000 个端口)
  RDP (Windows):   33000-34000  (1000 个端口)
  JupyterLab:      38000-39000  (1000 个端口)
  自定义服务:       40000-41000  (1000 个端口)
```

**数据库设计：**

```sql
-- 端口池表
CREATE TABLE port_pools (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    port_range_start INT NOT NULL,
    port_range_end INT NOT NULL,
    total_ports INT NOT NULL,
    used_ports INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 端口分配表（已在前面定义，这里补充）
CREATE TABLE port_mappings (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    service_type VARCHAR(32) NOT NULL,        -- ssh, rdp, jupyter
    external_port INT NOT NULL UNIQUE,
    internal_port INT NOT NULL DEFAULT 22,
    host_id VARCHAR(64),
    status VARCHAR(20) DEFAULT 'active',
    allocated_at TIMESTAMP DEFAULT NOW(),
    released_at TIMESTAMP,
    INDEX idx_env_id (env_id),
    INDEX idx_external_port (external_port),
    INDEX idx_status (status)
);
```

**端口管理器实现：**

```go
// port/port_manager.go
package port

type PortManager struct {
    db    *gorm.DB
    pools map[string]*PortPool
    mu    sync.Mutex
}

type PortPool struct {
    Name      string
    StartPort int
    EndPort   int
    UsedPorts map[int]bool
}

// 初始化端口池
func NewPortManager(db *gorm.DB) *PortManager {
    pm := &PortManager{
        db:    db,
        pools: make(map[string]*PortPool),
    }

    // 初始化各类端口池
    pm.pools["ssh-linux"] = &PortPool{
        Name:      "ssh-linux",
        StartPort: 30000,
        EndPort:   31000,
        UsedPorts: make(map[int]bool),
    }

    pm.pools["ssh-windows"] = &PortPool{
        Name:      "ssh-windows",
        StartPort: 31000,
        EndPort:   32000,
        UsedPorts: make(map[int]bool),
    }

    pm.pools["rdp"] = &PortPool{
        Name:      "rdp",
        StartPort: 33000,
        EndPort:   34000,
        UsedPorts: make(map[int]bool),
    }

    // 从数据库加载已分配的端口
    pm.loadAllocatedPorts()

    return pm
}

// 分配端口
func (pm *PortManager) AllocatePort(poolName string) (int, error) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    pool, exists := pm.pools[poolName]
    if !exists {
        return 0, errors.New("port pool not found")
    }

    // 查找可用端口
    for port := pool.StartPort; port < pool.EndPort; port++ {
        if !pool.UsedPorts[port] {
            pool.UsedPorts[port] = true
            return port, nil
        }
    }

    return 0, errors.New("no available port")
}

// 释放端口
func (pm *PortManager) ReleasePort(port int) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // 查找端口所属的池
    for _, pool := range pm.pools {
        if port >= pool.StartPort && port < pool.EndPort {
            delete(pool.UsedPorts, port)

            // 更新数据库
            pm.db.Model(&PortMapping{}).
                Where("external_port = ?", port).
                Updates(map[string]interface{}{
                    "status":      "released",
                    "released_at": time.Now(),
                })

            return nil
        }
    }

    return errors.New("port not found")
}
```

---

### 5.3 模块：存储管理 (Storage Manager)

**职责：** 管理持久化存储、数据集挂载

**存储类型：**

```yaml
存储类型:
  1. 用户代码存储
     - 路径: /gemini/code
     - 大小: 10GB-100GB
     - 持久化: 是

  2. 输出存储
     - 路径: /gemini/output
     - 大小: 50GB-500GB
     - 持久化: 是

  3. 数据集挂载
     - 路径: /gemini/data-*
     - 只读: 是
     - 来源: 对象存储

  4. 模型挂载
     - 路径: /gemini/pretrain*
     - 只读: 是
     - 来源: 对象存储
```

**传统架构实现：**

```go
// storage/local_storage.go
package storage

// 本地存储管理器（传统架构）
type LocalStorageManager struct {
    baseDir string // /data/environments
}

// 创建环境存储
func (m *LocalStorageManager) CreateEnvStorage(envID string) error {
    envDir := filepath.Join(m.baseDir, envID)

    // 创建目录结构
    dirs := []string{
        filepath.Join(envDir, "code"),
        filepath.Join(envDir, "output"),
    }

    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return err
        }
    }

    return nil
}

// 挂载数据集
func (m *LocalStorageManager) MountDataset(envID, datasetID string, index int) error {
    // 1. 从对象存储下载数据集到本地
    localPath := filepath.Join(m.baseDir, envID, fmt.Sprintf("data-%d", index))

    // 使用 rclone 同步
    cmd := exec.Command(
        "rclone", "sync",
        fmt.Sprintf("minio:remotegpu/datasets/%s", datasetID),
        localPath,
    )

    return cmd.Run()
}
```

**Kubernetes 实现：**

```go
// storage/k8s_storage.go
package storage

// K8s 存储管理器
type K8sStorageManager struct {
    clientset *kubernetes.Clientset
}

// 创建 PVC
func (m *K8sStorageManager) CreatePVC(envID string, size string) error {
    pvc := &corev1.PersistentVolumeClaim{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("code-%s", envID),
            Namespace: "dev-environments",
        },
        Spec: corev1.PersistentVolumeClaimSpec{
            AccessModes: []corev1.PersistentVolumeAccessMode{
                corev1.ReadWriteOnce,
            },
            Resources: corev1.ResourceRequirements{
                Requests: corev1.ResourceList{
                    "storage": resource.MustParse(size),
                },
            },
        },
    }

    _, err := m.clientset.CoreV1().
        PersistentVolumeClaims("dev-environments").
        Create(context.TODO(), pvc, metav1.CreateOptions{})

    return err
}
```

**注意事项：**
- ✅ 定期清理未使用的存储
- ✅ 存储配额限制
- ✅ 数据备份策略
- ✅ 快照功能（可选）

---

## 6. 开发路线图 (Development Roadmap)

### 6.1 阶段划分

```
第一阶段：基础设施（2-3 周）
├── 数据库设计与初始化
├── 对象存储部署（MinIO）
├── 缓存服务部署（Redis）
├── 消息队列部署（RabbitMQ）
└── 监控系统搭建（Prometheus + Grafana）

第二阶段：设备管理（2-3 周）
├── 主机注册与管理
├── GPU 设备发现
├── 设备监控与健康检查
└── 主机 Agent 开发

第三阶段：核心服务（3-4 周）
├── 统一调度器开发
├── 端口管理器实现
├── 存储管理器实现
└── 资源分配逻辑

第四阶段：环境管理（3-4 周）
├── 开发环境创建（Linux）
├── 开发环境创建（Windows）
├── SSH 公网访问
├── RDP 远程桌面
└── JupyterLab 集成

第五阶段：数据与镜像（2-3 周）
├── 数据集管理
├── 模型管理
├── 镜像管理
└── 文件上传下载

第六阶段：用户与权限（2 周）
├── 用户认证系统
├── 工作空间管理
├── RBAC 权限控制
└── 配额管理

第七阶段：训练与推理（3-4 周）
├── 离线训练任务
├── 分布式训练
├── 模型部署
└── 推理服务

第八阶段：计费与监控（2-3 周）
├── 资源计费
├── 账单生成
├── 监控告警
└── 日志收集

第九阶段：前端开发（4-5 周）
├── 控制台界面
├── 管理后台
├── Web 终端
└── Web RDP

第十阶段：测试与优化（2-3 周）
├── 功能测试
├── 性能测试
├── 安全测试
└── 文档完善
```

### 6.2 优先级排序

**P0（必须）- MVP 核心功能：**
1. 基础设施搭建
2. 主机管理（Linux）
3. GPU 资源管理
4. 统一调度器
5. 开发环境创建（Linux + Docker）
6. SSH 公网访问
7. 用户认证系统
8. 基础前端界面

**P1（重要）- 扩展功能：**
1. Windows 主机支持
2. RDP 远程桌面
3. 数据集管理
4. 镜像管理
5. JupyterLab 集成
6. 监控告警
7. 资源配额

**P2（可选）- 高级功能：**
1. 离线训练任务
2. 推理服务
3. 计费系统
4. 分布式训练
5. Web 终端
6. 自定义镜像构建

### 6.3 模块依赖关系

```
基础设施层（无依赖）
    ↓
设备管理模块（依赖：基础设施）
    ↓
核心服务层（依赖：基础设施 + 设备管理）
    ↓
环境管理模块（依赖：核心服务 + 设备管理）
    ↓
数据与镜像模块（依赖：基础设施 + 环境管理）
    ↓
用户与权限模块（依赖：基础设施）
    ↓
训练与推理模块（依赖：环境管理 + 数据与镜像）
    ↓
计费与监控模块（依赖：所有业务模块）
    ↓
前端展示层（依赖：所有后端模块）
```

### 6.4 关键里程碑

**里程碑 1：基础平台可用（4-6 周）**
- ✅ 基础设施部署完成
- ✅ 主机管理功能可用
- ✅ 可以创建 Linux 开发环境
- ✅ SSH 公网访问可用
- ✅ 基础前端界面

**里程碑 2：MVP 上线（8-10 周）**
- ✅ 用户认证系统
- ✅ 数据集管理
- ✅ 镜像管理
- ✅ JupyterLab 集成
- ✅ 基础监控

**里程碑 3：功能完善（12-15 周）**
- ✅ Windows 主机支持
- ✅ 离线训练任务
- ✅ 推理服务
- ✅ 计费系统
- ✅ 完整前端

**里程碑 4：生产就绪（16-20 周）**
- ✅ 性能优化
- ✅ 安全加固
- ✅ 高可用部署
- ✅ 完整文档

---

## 7. 总结

### 7.1 核心设计原则

1. **模块化设计**
   - 每个模块职责清晰
   - 模块间低耦合
   - 易于扩展和维护

2. **架构灵活性**
   - 支持传统架构和 K8s 架构
   - 支持 Linux 和 Windows 主机
   - 支持混合云部署

3. **可扩展性**
   - 水平扩展能力
   - 支持多区域部署
   - 支持大规模并发

4. **安全性**
   - 多租户隔离
   - 数据加密
   - 访问控制
   - 审计日志

### 7.2 技术栈总结

**后端技术栈：**
```
语言：Go 1.21+
框架：Gin (HTTP), gRPC
数据库：PostgreSQL 14+, Redis 7+
消息队列：RabbitMQ
对象存储：MinIO
容器：Docker, Kubernetes 1.28+
GPU：NVIDIA GPU Operator, DCGM
```

**前端技术栈：**
```
框架：React 18 + TypeScript
状态管理：Redux Toolkit
UI 库：Ant Design
图表：ECharts
编辑器：Monaco Editor
终端：xterm.js
构建：Vite
```

**基础设施：**
```
监控：Prometheus + Grafana
日志：ELK Stack
镜像仓库：Harbor
负载均衡：Nginx / Traefik
```

### 7.3 下一步行动

1. **立即开始：**
   - 搭建开发环境
   - 初始化代码仓库
   - 设计数据库表结构
   - 部署基础设施（PostgreSQL, Redis, MinIO）

2. **第一周任务：**
   - 完成数据库设计
   - 实现主机注册 API
   - 开发 GPU 发现功能
   - 搭建前端框架

3. **持续关注：**
   - 性能监控
   - 安全审计
   - 用户反馈
   - 技术债务

---

**文档结束**

本文档提供了 RemoteGPU 系统的完整架构设计，涵盖了从基础设施到业务功能的所有模块。按照本文档的规划，可以系统化地完成整个平台的开发工作。
