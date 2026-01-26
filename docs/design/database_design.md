# RemoteGPU 数据库设计文档

> 完整的数据库表结构设计与字段说明
>
> 创建日期：2026-01-26
>
> 前端技术栈：Vue 3 + Element Plus

---

## 目录

1. [数据库概览](#1-数据库概览)
2. [表关系图](#2-表关系图)
3. [基础设施相关表](#3-基础设施相关表)
4. [设备管理相关表](#4-设备管理相关表)
5. [用户与权限相关表](#5-用户与权限相关表)
6. [环境管理相关表](#6-环境管理相关表)
7. [数据与镜像相关表](#7-数据与镜像相关表)
8. [训练与推理相关表](#8-训练与推理相关表)
9. [计费相关表](#9-计费相关表)
10. [监控相关表](#10-监控相关表)

---

## 1. 数据库概览

### 1.1 数据库信息

```yaml
数据库名称: remotegpu
数据库类型: PostgreSQL 14+
字符集: UTF-8
时区: UTC
```

### 1.2 表统计

```yaml
总表数: 35+
核心业务表: 20
监控数据表: 8
关联关系表: 7
```

### 1.3 命名规范

```yaml
表名: 小写 + 下划线分隔 (snake_case)
字段名: 小写 + 下划线分隔 (snake_case)
主键: id (BIGSERIAL)
外键: {表名}_id
时间戳: created_at, updated_at
软删除: deleted_at (可选)
```

---

## 2. 表关系图

### 2.1 核心表关系

```
┌─────────────────────────────────────────────────────────────┐
│                        核心表关系图                            │
└─────────────────────────────────────────────────────────────┘

customers (客户表)
    │
    ├──> workspaces (工作空间)
    │       │
    │       └──> workspace_members (成员)
    │
    ├──> environments (开发环境)
    │       │
    │       ├──> port_mappings (端口映射)
    │       └──> dataset_usage (数据集使用)
    │
    ├──> datasets (数据集)
    │       │
    │       └──> dataset_versions (版本)
    │
    ├──> models (模型)
    │       │
    │       └──> model_versions (版本)
    │
    └──> billing_records (计费记录)
            │
            └──> invoices (账单)

hosts (主机表)
    │
    ├──> gpus (GPU设备)
    │       │
    │       └──> gpu_metrics (GPU监控)
    │
    ├──> host_metrics (主机监控)
    │
    └──> environments (开发环境)
```

---

## 3. 基础设施相关表

### 3.1 系统配置表 (system_configs)

**表说明：** 存储系统全局配置

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 主键 | 1 |
| config_key | VARCHAR(128) | UNIQUE NOT NULL | 配置键 | "ssh_port_range_start" |
| config_value | TEXT | NOT NULL | 配置值 | "30000" |
| config_type | VARCHAR(32) | NOT NULL | 值类型 | "integer", "string", "json" |
| description | TEXT | | 配置说明 | "SSH端口范围起始值" |
| is_public | BOOLEAN | DEFAULT false | 是否公开 | false |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 10:00:00 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_config_key` ON (config_key)

**关联关系：** 无

---

## 4. 设备管理相关表

### 4.1 主机表 (hosts)

**表说明：** 存储物理机/虚拟机信息，支持 Linux 和 Windows

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | VARCHAR(64) | PRIMARY KEY | 主机ID | "host-abc123" |
| name | VARCHAR(128) | NOT NULL | 主机名称 | "GPU-Server-01" |
| hostname | VARCHAR(256) | | 主机名 | "gpu01.example.com" |
| ip_address | VARCHAR(64) | NOT NULL | 内网IP | "192.168.1.10" |
| public_ip | VARCHAR(64) | | 公网IP | "1.2.3.4" |
| os_type | VARCHAR(20) | NOT NULL | 操作系统类型 | "linux", "windows" |
| os_version | VARCHAR(64) | | 操作系统版本 | "Ubuntu 20.04" |
| arch | VARCHAR(20) | DEFAULT 'x86_64' | CPU架构 | "x86_64", "arm64" |
| deployment_mode | VARCHAR(20) | NOT NULL | 部署模式 | "traditional", "kubernetes" |
| k8s_node_name | VARCHAR(128) | | K8s节点名 | "node-01" |
| status | VARCHAR(20) | DEFAULT 'offline' | 状态 | "online", "offline", "maintenance" |
| health_status | VARCHAR(20) | DEFAULT 'unknown' | 健康状态 | "healthy", "degraded", "unhealthy" |
| total_cpu | INT | NOT NULL | CPU总核心数 | 32 |
| total_memory | BIGINT | NOT NULL | 内存总量(字节) | 137438953472 (128GB) |
| total_disk | BIGINT | | 磁盘总量(字节) | 2199023255552 (2TB) |
| total_gpu | INT | DEFAULT 0 | GPU总数量 | 4 |
| used_cpu | INT | DEFAULT 0 | 已用CPU | 8 |
| used_memory | BIGINT | DEFAULT 0 | 已用内存 | 34359738368 (32GB) |
| used_disk | BIGINT | DEFAULT 0 | 已用磁盘 | 549755813888 (512GB) |
| used_gpu | INT | DEFAULT 0 | 已用GPU | 2 |
| ssh_port | INT | DEFAULT 22 | SSH端口 | 22 |
| winrm_port | INT | | WinRM端口 | 5985 |
| agent_port | INT | DEFAULT 8080 | Agent端口 | 8080 |
| labels | JSONB | | 标签 | {"region": "us-west", "zone": "a"} |
| tags | TEXT[] | | 标签数组 | ["gpu-server", "high-memory"] |
| last_heartbeat | TIMESTAMP | | 最后心跳时间 | 2026-01-26 10:00:00 |
| registered_at | TIMESTAMP | DEFAULT NOW() | 注册时间 | 2026-01-26 09:00:00 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_os_type` ON (os_type)
- `idx_status` ON (status)
- `idx_deployment_mode` ON (deployment_mode)
- `idx_labels` ON (labels) USING GIN

**关联关系：**
- 一对多：`gpus` (一个主机有多个GPU)
- 一对多：`host_metrics` (一个主机有多条监控记录)
- 一对多：`environments` (一个主机运行多个环境)

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: '主机名称',
    type: 'input',
    required: true,
    placeholder: '请输入主机名称'
  },
  ip_address: {
    label: 'IP地址',
    type: 'input',
    required: true,
    rules: [{ type: 'ip', message: '请输入有效的IP地址' }]
  },
  os_type: {
    label: '操作系统',
    type: 'select',
    required: true,
    options: [
      { label: 'Linux', value: 'linux' },
      { label: 'Windows', value: 'windows' }
    ]
  },
  deployment_mode: {
    label: '部署模式',
    type: 'select',
    required: true,
    options: [
      { label: '传统架构', value: 'traditional' },
      { label: 'Kubernetes', value: 'kubernetes' }
    ]
  }
}
```

---

### 4.2 GPU设备表 (gpus)

**表说明：** 存储GPU设备信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | GPU ID | 1 |
| host_id | VARCHAR(64) | NOT NULL, FK | 所属主机ID | "host-abc123" |
| gpu_index | INT | NOT NULL | GPU索引 | 0, 1, 2, 3 |
| uuid | VARCHAR(128) | UNIQUE | GPU UUID | "GPU-12345678-1234..." |
| name | VARCHAR(128) | | GPU型号 | "Tesla V100-SXM2-32GB" |
| brand | VARCHAR(64) | | 品牌 | "NVIDIA", "AMD" |
| architecture | VARCHAR(64) | | 架构 | "Ampere", "Turing" |
| memory_total | BIGINT | | 显存总量(字节) | 34359738368 (32GB) |
| cuda_cores | INT | | CUDA核心数 | 5120 |
| compute_capability | VARCHAR(32) | | 计算能力 | "7.0", "8.0" |
| status | VARCHAR(20) | DEFAULT 'available' | 状态 | "available", "allocated", "maintenance" |
| health_status | VARCHAR(20) | DEFAULT 'healthy' | 健康状态 | "healthy", "degraded", "error" |
| allocated_to | VARCHAR(64) | | 分配给的环境ID | "env-xyz789" |
| allocated_at | TIMESTAMP | | 分配时间 | 2026-01-26 10:00:00 |
| power_limit | INT | | 功耗限制(瓦) | 300 |
| temperature_limit | INT | | 温度限制(℃) | 85 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_host_id` ON (host_id)
- `idx_status` ON (status)
- `idx_allocated_to` ON (allocated_to)
- `UNIQUE(host_id, gpu_index)`

**外键约束：**
- `host_id` REFERENCES `hosts(id)` ON DELETE CASCADE

**关联关系：**
- 多对一：`hosts` (多个GPU属于一个主机)
- 一对多：`gpu_metrics` (一个GPU有多条监控记录)

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: 'GPU型号',
    type: 'input',
    disabled: true,
    placeholder: '自动检测'
  },
  memory_total: {
    label: '显存容量',
    type: 'input',
    disabled: true,
    formatter: (value) => `${(value / 1024 / 1024 / 1024).toFixed(0)} GB`
  },
  status: {
    label: '状态',
    type: 'tag',
    colorMap: {
      available: 'success',
      allocated: 'warning',
      maintenance: 'info',
      error: 'danger'
    }
  }
}
```

---

### 4.3 主机监控数据表 (host_metrics)

**表说明：** 存储主机监控指标（时序数据）

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| host_id | VARCHAR(64) | NOT NULL, FK | 主机ID | "host-abc123" |
| cpu_usage_percent | FLOAT | | CPU使用率(%) | 45.5 |
| cpu_load_1m | FLOAT | | 1分钟负载 | 2.5 |
| cpu_load_5m | FLOAT | | 5分钟负载 | 2.3 |
| cpu_load_15m | FLOAT | | 15分钟负载 | 2.1 |
| memory_used | BIGINT | | 已用内存(字节) | 68719476736 (64GB) |
| memory_available | BIGINT | | 可用内存(字节) | 68719476736 (64GB) |
| memory_usage_percent | FLOAT | | 内存使用率(%) | 50.0 |
| disk_used | BIGINT | | 已用磁盘(字节) | 1099511627776 (1TB) |
| disk_available | BIGINT | | 可用磁盘(字节) | 1099511627776 (1TB) |
| disk_usage_percent | FLOAT | | 磁盘使用率(%) | 50.0 |
| disk_io_read_bytes | BIGINT | | 磁盘读取(字节) | 1073741824 |
| disk_io_write_bytes | BIGINT | | 磁盘写入(字节) | 536870912 |
| network_rx_bytes | BIGINT | | 网络接收(字节) | 10737418240 |
| network_tx_bytes | BIGINT | | 网络发送(字节) | 5368709120 |
| network_rx_packets | BIGINT | | 接收包数 | 1000000 |
| network_tx_packets | BIGINT | | 发送包数 | 500000 |
| gpu_avg_utilization | FLOAT | | GPU平均使用率(%) | 75.0 |
| gpu_avg_memory_used | BIGINT | | GPU平均显存使用 | 17179869184 (16GB) |
| gpu_avg_temperature | FLOAT | | GPU平均温度(℃) | 65.0 |
| gpu_avg_power | FLOAT | | GPU平均功耗(瓦) | 250.0 |
| collected_at | TIMESTAMP | DEFAULT NOW() | 采集时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_host_id_time` ON (host_id, collected_at DESC)

**外键约束：**
- `host_id` REFERENCES `hosts(id)` ON DELETE CASCADE

**关联关系：**
- 多对一：`hosts` (多条监控记录属于一个主机)

**数据保留策略：**
- 详细数据：保留 7 天
- 聚合数据（小时级）：保留 90 天
- 聚合数据（天级）：保留 1 年

**Vue 图表配置：**
```javascript
{
  chartType: 'line',
  metrics: [
    { key: 'cpu_usage_percent', label: 'CPU使用率', unit: '%', color: '#409EFF' },
    { key: 'memory_usage_percent', label: '内存使用率', unit: '%', color: '#67C23A' },
    { key: 'disk_usage_percent', label: '磁盘使用率', unit: '%', color: '#E6A23C' },
    { key: 'gpu_avg_utilization', label: 'GPU使用率', unit: '%', color: '#F56C6C' }
  ],
  timeRange: ['1h', '6h', '24h', '7d', '30d']
}
```

---

### 4.4 GPU监控数据表 (gpu_metrics)

**表说明：** 存储单个GPU的监控指标（时序数据）

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| gpu_id | BIGINT | NOT NULL, FK | GPU ID | 1 |
| host_id | VARCHAR(64) | NOT NULL, FK | 主机ID | "host-abc123" |
| utilization_percent | FLOAT | | GPU使用率(%) | 85.5 |
| memory_used | BIGINT | | 显存使用(字节) | 25769803776 (24GB) |
| memory_usage_percent | FLOAT | | 显存使用率(%) | 75.0 |
| temperature | FLOAT | | 温度(℃) | 68.5 |
| power_draw | FLOAT | | 功耗(瓦) | 280.5 |
| fan_speed_percent | FLOAT | | 风扇转速(%) | 60.0 |
| sm_clock | INT | | SM时钟频率(MHz) | 1530 |
| memory_clock | INT | | 显存时钟频率(MHz) | 877 |
| process_count | INT | | 运行进程数 | 2 |
| collected_at | TIMESTAMP | DEFAULT NOW() | 采集时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_gpu_id_time` ON (gpu_id, collected_at DESC)
- `idx_host_id_time` ON (host_id, collected_at DESC)

**外键约束：**
- `gpu_id` REFERENCES `gpus(id)` ON DELETE CASCADE
- `host_id` REFERENCES `hosts(id)` ON DELETE CASCADE

**关联关系：**
- 多对一：`gpus` (多条监控记录属于一个GPU)
- 多对一：`hosts` (多条监控记录属于一个主机)

**Vue 实时监控组件：**
```javascript
{
  refreshInterval: 30000, // 30秒刷新
  gaugeCharts: [
    { key: 'utilization_percent', label: 'GPU使用率', max: 100, unit: '%' },
    { key: 'memory_usage_percent', label: '显存使用率', max: 100, unit: '%' },
    { key: 'temperature', label: '温度', max: 100, unit: '℃', warning: 80, danger: 85 },
    { key: 'power_draw', label: '功耗', max: 350, unit: 'W' }
  ]
}
```

---

## 5. 用户与权限相关表

### 5.1 客户表 (customers)

**表说明：** 存储客户/用户信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 客户ID | 1 |
| uuid | VARCHAR(64) | UNIQUE NOT NULL | 客户UUID | "cust-abc123" |
| username | VARCHAR(128) | UNIQUE NOT NULL | 用户名 | "john_doe" |
| email | VARCHAR(256) | UNIQUE NOT NULL | 邮箱 | "john@example.com" |
| phone | VARCHAR(32) | | 手机号 | "+86 138****1234" |
| password_hash | VARCHAR(256) | NOT NULL | 密码哈希 | "$2a$10$..." |
| full_name | VARCHAR(256) | | 全名 | "John Doe" |
| company | VARCHAR(256) | | 公司名称 | "Example Corp" |
| account_type | VARCHAR(20) | DEFAULT 'individual' | 账户类型 | "individual", "enterprise" |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 | "active", "suspended", "deleted" |
| email_verified | BOOLEAN | DEFAULT false | 邮箱已验证 | true |
| phone_verified | BOOLEAN | DEFAULT false | 手机已验证 | false |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_username` ON (username)
- `idx_email` ON (email)
- `idx_status` ON (status)

**关联关系：**
- 一对多：`workspaces` (一个客户可以创建多个工作空间)
- 一对多：`environments` (一个客户可以创建多个开发环境)

**Vue 表单字段配置：**
```javascript
{
  username: {
    label: '用户名',
    type: 'input',
    required: true,
    placeholder: '请输入用户名'
  },
  email: {
    label: '邮箱',
    type: 'input',
    required: true,
    rules: [{ type: 'email', message: '请输入有效的邮箱地址' }]
  },
  account_type: {
    label: '账户类型',
    type: 'select',
    options: [
      { label: '个人用户', value: 'individual' },
      { label: '企业用户', value: 'enterprise' }
    ]
  }
}
```

---

### 5.2 工作空间表 (workspaces)

**表说明：** 存储工作空间信息，支持团队协作

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 工作空间ID | 1 |
| uuid | VARCHAR(64) | UNIQUE NOT NULL | 工作空间UUID | "ws-xyz789" |
| owner_id | BIGINT | NOT NULL, FK | 所有者ID | 1 |
| name | VARCHAR(128) | NOT NULL | 工作空间名称 | "AI Research Team" |
| description | TEXT | | 描述 | "深度学习研究团队" |
| type | VARCHAR(20) | DEFAULT 'team' | 类型 | "personal", "team", "enterprise" |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 | "active", "archived" |
| member_count | INT | DEFAULT 1 | 成员数量 | 5 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_owner_id` ON (owner_id)

**外键约束：**
- `owner_id` REFERENCES `customers(id)` ON DELETE CASCADE

**关联关系：**
- 多对一：`customers` (多个工作空间属于一个所有者)
- 一对多：`workspace_members` (一个工作空间有多个成员)

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: '工作空间名称',
    type: 'input',
    required: true,
    placeholder: '请输入工作空间名称'
  },
  type: {
    label: '类型',
    type: 'select',
    options: [
      { label: '个人空间', value: 'personal' },
      { label: '团队空间', value: 'team' }
    ]
  }
}
```

---

### 5.3 工作空间成员表 (workspace_members)

**表说明：** 存储工作空间成员关系和权限

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| workspace_id | BIGINT | NOT NULL, FK | 工作空间ID | 1 |
| customer_id | BIGINT | NOT NULL, FK | 客户ID | 2 |
| role | VARCHAR(20) | NOT NULL | 角色 | "owner", "admin", "member", "viewer" |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 | "active", "invited", "suspended" |
| joined_at | TIMESTAMP | | 加入时间 | 2026-01-26 10:00:00 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_workspace_id` ON (workspace_id)
- `idx_customer_id` ON (customer_id)
- `UNIQUE(workspace_id, customer_id)`

**外键约束：**
- `workspace_id` REFERENCES `workspaces(id)` ON DELETE CASCADE
- `customer_id` REFERENCES `customers(id)` ON DELETE CASCADE

**关联关系：**
- 多对一：`workspaces` (多个成员属于一个工作空间)
- 多对一：`customers` (多个成员记录关联到客户)

**Vue 表单字段配置：**
```javascript
{
  role: {
    label: '角色',
    type: 'select',
    required: true,
    options: [
      { label: '管理员', value: 'admin' },
      { label: '成员', value: 'member' },
      { label: '访客', value: 'viewer' }
    ]
  }
}
```

---

## 6. 环境管理相关表

### 6.1 开发环境表 (environments)

**表说明：** 存储开发环境实例信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | VARCHAR(64) | PRIMARY KEY | 环境ID | "env-abc123" |
| customer_id | BIGINT | NOT NULL, FK | 客户ID | 1 |
| workspace_id | BIGINT | FK | 工作空间ID | 1 |
| host_id | VARCHAR(64) | NOT NULL, FK | 所在主机ID | "host-xyz789" |
| name | VARCHAR(128) | NOT NULL | 环境名称 | "PyTorch Training" |
| image | VARCHAR(256) | NOT NULL | 镜像名称 | "ubuntu20-pytorch:2.0" |
| status | VARCHAR(20) | DEFAULT 'creating' | 状态 | "creating", "running", "stopped" |
| cpu | INT | NOT NULL | CPU核心数 | 4 |
| memory | BIGINT | NOT NULL | 内存(字节) | 17179869184 (16GB) |
| gpu | INT | DEFAULT 0 | GPU数量 | 1 |
| ssh_port | INT | | SSH端口 | 30001 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_customer_id` ON (customer_id)
- `idx_host_id` ON (host_id)
- `idx_status` ON (status)

**外键约束：**
- `customer_id` REFERENCES `customers(id)` ON DELETE CASCADE
- `host_id` REFERENCES `hosts(id)` ON DELETE RESTRICT

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: '环境名称',
    type: 'input',
    required: true
  },
  image: {
    label: '镜像',
    type: 'select',
    required: true,
    options: [
      { label: 'Ubuntu 20.04 + PyTorch 2.0', value: 'ubuntu20-pytorch:2.0' },
      { label: 'Ubuntu 20.04 + TensorFlow 2.12', value: 'ubuntu20-tf:2.12' }
    ]
  },
  cpu: {
    label: 'CPU核心数',
    type: 'input-number',
    min: 1,
    max: 64,
    default: 4
  },
  memory: {
    label: '内存(GB)',
    type: 'input-number',
    min: 1,
    max: 512,
    default: 16
  },
  gpu: {
    label: 'GPU数量',
    type: 'input-number',
    min: 0,
    max: 8,
    default: 1
  }
}
```

---

### 6.2 端口映射表 (port_mappings)

**表说明：** 存储SSH/RDP/JupyterLab端口映射关系

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| env_id | VARCHAR(64) | NOT NULL, FK | 环境ID | "env-abc123" |
| service_type | VARCHAR(32) | NOT NULL | 服务类型 | "ssh", "rdp", "jupyter" |
| external_port | INT | NOT NULL UNIQUE | 外部端口 | 30001 |
| internal_port | INT | NOT NULL | 内部端口 | 22 |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 | "active", "released" |
| allocated_at | TIMESTAMP | DEFAULT NOW() | 分配时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_env_id` ON (env_id)
- `idx_external_port` ON (external_port)

**外键约束：**
- `env_id` REFERENCES `environments(id)` ON DELETE CASCADE

**Vue 表单字段配置：**
```javascript
{
  service_type: {
    label: '服务类型',
    type: 'select',
    options: [
      { label: 'SSH', value: 'ssh' },
      { label: 'JupyterLab', value: 'jupyter' }
    ]
  }
}
```

---

## 7. 数据与镜像相关表

### 7.1 数据集表 (datasets)

**表说明：** 存储数据集元信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 数据集ID | 1 |
| uuid | VARCHAR(64) | UNIQUE NOT NULL | 数据集UUID | "dataset-abc123" |
| customer_id | BIGINT | NOT NULL, FK | 客户ID | 1 |
| name | VARCHAR(256) | NOT NULL | 数据集名称 | "ImageNet-1K" |
| description | TEXT | | 描述 | "图像分类数据集" |
| storage_path | VARCHAR(512) | NOT NULL | 存储路径 | "datasets/customer-1/" |
| total_size | BIGINT | DEFAULT 0 | 总大小(字节) | 1099511627776 (1TB) |
| status | VARCHAR(20) | DEFAULT 'uploading' | 状态 | "uploading", "ready" |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_customer_id` ON (customer_id)
- `idx_status` ON (status)

**外键约束：**
- `customer_id` REFERENCES `customers(id)` ON DELETE CASCADE

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: '数据集名称',
    type: 'input',
    required: true,
    placeholder: '请输入数据集名称'
  },
  description: {
    label: '描述',
    type: 'textarea',
    rows: 3
  }
}
```

---

### 7.2 模型表 (models)

**表说明：** 存储模型文件元信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 模型ID | 1 |
| uuid | VARCHAR(64) | UNIQUE NOT NULL | 模型UUID | "model-xyz789" |
| customer_id | BIGINT | NOT NULL, FK | 客户ID | 1 |
| name | VARCHAR(256) | NOT NULL | 模型名称 | "ResNet-50" |
| framework | VARCHAR(64) | | 框架 | "pytorch", "tensorflow" |
| storage_path | VARCHAR(512) | NOT NULL | 存储路径 | "models/customer-1/" |
| total_size | BIGINT | DEFAULT 0 | 总大小(字节) | 102400000 (100MB) |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_customer_id` ON (customer_id)

**外键约束：**
- `customer_id` REFERENCES `customers(id)` ON DELETE CASCADE

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: '模型名称',
    type: 'input',
    required: true
  },
  framework: {
    label: '框架',
    type: 'select',
    options: [
      { label: 'PyTorch', value: 'pytorch' },
      { label: 'TensorFlow', value: 'tensorflow' }
    ]
  }
}
```

---

### 7.3 镜像表 (images)

**表说明：** 存储Docker镜像信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 镜像ID | 1 |
| name | VARCHAR(256) | UNIQUE NOT NULL | 镜像名称 | "ubuntu20-pytorch:2.0" |
| description | TEXT | | 描述 | "Ubuntu 20.04 + PyTorch 2.0" |
| category | VARCHAR(64) | | 分类 | "base", "pytorch", "tensorflow" |
| is_official | BOOLEAN | DEFAULT false | 是否官方镜像 | true |
| size | BIGINT | | 镜像大小(字节) | 5368709120 (5GB) |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_category` ON (category)
- `idx_is_official` ON (is_official)

**Vue 表单字段配置：**
```javascript
{
  name: {
    label: '镜像名称',
    type: 'input',
    required: true
  },
  category: {
    label: '分类',
    type: 'select',
    options: [
      { label: '基础镜像', value: 'base' },
      { label: 'PyTorch', value: 'pytorch' },
      { label: 'TensorFlow', value: 'tensorflow' }
    ]
  }
}
```

---

## 8. 计费相关表

### 8.1 计费记录表 (billing_records)

**表说明：** 存储资源使用计费记录

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| customer_id | BIGINT | NOT NULL, FK | 客户ID | 1 |
| env_id | VARCHAR(64) | FK | 环境ID | "env-abc123" |
| resource_type | VARCHAR(32) | NOT NULL | 资源类型 | "cpu", "memory", "gpu" |
| quantity | FLOAT | NOT NULL | 数量 | 4.0 |
| unit_price | DECIMAL(10,4) | NOT NULL | 单价 | 0.5000 |
| amount | DECIMAL(10,4) | NOT NULL | 金额 | 2.0000 |
| start_time | TIMESTAMP | NOT NULL | 开始时间 | 2026-01-26 09:00:00 |
| end_time | TIMESTAMP | NOT NULL | 结束时间 | 2026-01-26 10:00:00 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_customer_id` ON (customer_id)
- `idx_env_id` ON (env_id)
- `idx_start_time` ON (start_time DESC)

**外键约束：**
- `customer_id` REFERENCES `customers(id)` ON DELETE CASCADE

**Vue 表单字段配置：**
```javascript
{
  resource_type: {
    label: '资源类型',
    type: 'tag',
    colorMap: {
      cpu: 'primary',
      memory: 'success',
      gpu: 'warning'
    }
  },
  amount: {
    label: '金额',
    type: 'text',
    formatter: (value) => `¥${value.toFixed(2)}`
  }
}
```

---

## 9. 监控相关表

### 9.1 环境监控数据表 (environment_metrics)

**表说明：** 存储环境级别的监控指标

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| env_id | VARCHAR(64) | NOT NULL, FK | 环境ID | "env-abc123" |
| cpu_usage_percent | FLOAT | | CPU使用率(%) | 45.5 |
| memory_usage_percent | FLOAT | | 内存使用率(%) | 60.0 |
| gpu_usage_percent | FLOAT | | GPU使用率(%) | 85.0 |
| network_rx_bytes | BIGINT | | 网络接收(字节) | 1073741824 |
| network_tx_bytes | BIGINT | | 网络发送(字节) | 536870912 |
| collected_at | TIMESTAMP | DEFAULT NOW() | 采集时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_env_id_time` ON (env_id, collected_at DESC)

**外键约束：**
- `env_id` REFERENCES `environments(id)` ON DELETE CASCADE

**Vue 图表配置：**
```javascript
{
  chartType: 'line',
  metrics: [
    { key: 'cpu_usage_percent', label: 'CPU使用率', unit: '%', color: '#409EFF' },
    { key: 'memory_usage_percent', label: '内存使用率', unit: '%', color: '#67C23A' },
    { key: 'gpu_usage_percent', label: 'GPU使用率', unit: '%', color: '#F56C6C' }
  ]
}
```

---

**文档结束**

本文档提供了 RemoteGPU 系统的完整数据库表结构设计，包含了所有核心业务表的字段说明、索引、外键约束、关联关系以及 Vue + Element Plus 的表单配置。
