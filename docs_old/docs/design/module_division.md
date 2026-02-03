# RemoteGPU 系统模块划分

> 完整的系统模块划分与职责定义
>
> 创建日期：2026-01-26
>
> 前端技术栈：Vue 3 + Element Plus

---

## 目录

1. [模块总览](#1-模块总览)
2. [模块 1：CMDB 设备管理模块](#2-模块-1cmdb-设备管理模块)
3. [模块 2：用户与权限模块](#3-模块-2用户与权限模块)
4. [模块 3：环境管理模块](#4-模块-3环境管理模块)
5. [模块 4：资源调度模块](#5-模块-4资源调度模块)
6. [模块 5：数据与存储模块](#6-模块-5数据与存储模块)
7. [模块 6：镜像管理模块](#7-模块-6镜像管理模块)
8. [模块 7：训练与推理模块](#8-模块-7训练与推理模块)
9. [模块 8：计费管理模块](#9-模块-8计费管理模块)
10. [模块 9：监控告警模块](#10-模块-9监控告警模块)
11. [模块 10：网关与认证模块](#11-模块-10网关与认证模块)
12. [模块 11：制品管理模块](#12-模块-11制品管理模块)
13. [模块 12：问题单管理模块](#13-模块-12问题单管理模块)
14. [模块 13：需求单管理模块](#14-模块-13需求单管理模块)
15. [模块 14：通知管理模块](#15-模块-14通知管理模块)
16. [模块 15：Webhook 管理模块](#16-模块-15webhook-管理模块)
17. [模块依赖关系](#17-模块依赖关系)
18. [开发优先级](#18-开发优先级)

---

## 1. 模块总览

### 1.1 模块架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        前端展示层                                 │
│                    Vue 3 + Element Plus                          │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   模块 10: 网关与认证模块                          │
│              API Gateway + Authentication                        │
└────────────────────────┬────────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ 模块 2:      │  │ 模块 3:      │  │ 模块 5:      │
│ 用户与权限    │  │ 环境管理      │  │ 数据与存储    │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └────────┬────────┴────────┬────────┘
                │                 │
                ▼                 ▼
        ┌──────────────┐  ┌──────────────┐
        │ 模块 4:      │  │ 模块 6:      │
        │ 资源调度      │  │ 镜像管理      │
        └──────┬───────┘  └──────┬───────┘
               │                 │
               ▼                 ▼
        ┌──────────────┐  ┌──────────────┐
        │ 模块 1:      │  │ 模块 11:     │
        │ CMDB设备管理  │  │ 制品管理      │
        └──────┬───────┘  └──────────────┘
               │
        ┌──────┴───────┬────────────┬────────────┐
        │              │            │            │
        ▼              ▼            ▼            ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ 模块 7:      │ │ 模块 8:      │ │ 模块 9:      │ │ 模块 12/13:  │
│ 训练与推理    │ │ 计费管理      │ │ 监控告警      │ │ 问题/需求单  │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
```

### 1.2 模块列表

| 模块编号 | 模块名称 | 英文名称 | 核心职责 | 优先级 |
|---------|---------|---------|---------|--------|
| 模块 1 | CMDB 设备管理模块 | CMDB & Device Management | 硬件资产管理、设备生命周期 | P0 |
| 模块 2 | 用户与权限模块 | User & Auth | 用户认证、工作空间、RBAC | P0 |
| 模块 3 | 环境管理模块 | Environment Management | 开发环境创建、SSH/RDP 访问 | P0 |
| 模块 4 | 资源调度模块 | Resource Scheduler | 统一调度、端口管理、资源分配 | P0 |
| 模块 5 | 数据与存储模块 | Data & Storage | 数据集管理、对象存储 | P1 |
| 模块 6 | 镜像管理模块 | Image Management | 官方镜像、自定义镜像构建 | P1 |
| 模块 7 | 训练与推理模块 | Training & Inference | 离线训练、推理服务 | P2 |
| 模块 8 | 计费管理模块 | Billing Management | 资源计费、账单生成 | P1 |
| 模块 9 | 监控告警模块 | Monitoring & Alerting | 资源监控、性能监控、告警 | P1 |
| 模块 10 | 网关与认证模块 | API Gateway & Auth | 统一入口、认证、限流 | P0 |
| 模块 11 | 制品管理模块 | Artifact Management | 软件包管理、版本控制、制品仓库 | P1 |
| 模块 12 | 问题单管理模块 | Issue Management | Bug跟踪、问题处理、工单流转 | P2 |
| 模块 13 | 需求单管理模块 | Requirement Management | 需求收集、需求评审、需求跟踪 | P2 |
| 模块 14 | 通知管理模块 | Notification Management | 多渠道通知、消息推送 | P1 |
| 模块 15 | Webhook 管理模块 | Webhook Management | 事件回调、第三方集成 | P1 |

---

## 2. 模块 1：CMDB 设备管理模块

### 2.1 模块概述

**职责：** 作为独立的配置管理数据库，管理所有硬件资产的完整生命周期，为其他业务模块提供统一的设备信息查询和状态管理接口。

**核心价值：**
- ✅ 单一数据源（Single Source of Truth）
- ✅ 设备状态一致性保证
- ✅ 完整的变更历史追溯
- ✅ 支持多系统集成

### 2.2 核心功能

```yaml
功能列表:
  1. 资产管理:
     1.1 主机注册与发现:
         - 自动发现主机（通过 Agent 心跳）
         - 手动注册主机（填写 IP、SSH 凭证）
         - 主机信息采集（CPU、内存、磁盘、网络）
         - 操作系统识别（Linux/Windows）
         - 部署模式识别（传统架构/Kubernetes）

     1.2 GPU 设备管理:
         - GPU 自动发现（nvidia-smi/DCGM）
         - GPU 信息采集（型号、显存、UUID、计算能力）
         - GPU 分配与释放
         - GPU 健康检查（温度、功耗、ECC 错误）
         - GPU 拓扑关系管理（NVLink、PCIe）

     1.3 资产信息维护:
         - 资产编号生成与管理
         - 资产基本信息编辑（名称、位置、负责人）
         - 采购信息管理（采购日期、价格、保修期）
         - 物理位置管理（机房、机柜、U 位）
         - 资产批量导入导出

     1.4 资产标签与分类:
         - 自定义标签管理（环境标签、业务标签）
         - 资产分组（按区域、用途、项目）
         - 标签搜索与筛选
         - 标签统计与报表

  2. 状态管理:
     2.1 运营状态管理:
         - available（可用）: 设备正常，可分配资源
         - maintenance（维护中）: 计划内维护，不可分配
         - faulty（故障）: 设备故障，需要修复
         - retired（已退役）: 设备下线，不再使用
         - reserved（预留）: 为特定用户/项目预留
         - 状态转换规则校验
         - 状态变更通知（邮件、钉钉、企业微信）

     2.2 使用状态管理:
         - idle（空闲）: 无资源占用
         - partial（部分使用）: 部分资源被占用
         - full（完全使用）: 资源已满
         - overcommit（超分配）: 允许超分配场景
         - 资源使用率计算（CPU、内存、GPU）
         - 资源水位告警（使用率 > 80%）

     2.3 健康状态监测:
         - healthy（健康）: 所有指标正常
         - degraded（降级）: 部分指标异常
         - unhealthy（不健康）: 严重异常
         - 心跳超时检测（30 秒无心跳标记为 offline）
         - 磁盘空间检查（可用空间 < 10% 告警）
         - 温度监控（GPU 温度 > 85°C 告警）
         - 自动隔离不健康设备

     2.4 状态转换控制:
         - 状态机管理（定义允许的状态转换路径）
         - 转换前置条件检查（如：有运行环境时不能进入维护）
         - 转换审批流程（关键状态变更需审批）
         - 转换回滚机制

  3. 生命周期管理:
     3.1 设备上线流程:
         - 设备验收（硬件检测、性能测试）
         - 初始化配置（安装 Agent、配置网络）
         - 状态设置为 available
         - 上线通知与记录

     3.2 维护管理:
         - 计划维护申请（提前通知用户）
         - 维护窗口管理（设置维护时间段）
         - 维护任务跟踪（维护内容、负责人、进度）
         - 维护完成验证（性能测试、功能验证）
         - 自动恢复到 available 状态

     3.3 故障处理:
         - 故障自动检测（心跳超时、健康检查失败）
         - 故障报警（多渠道通知）
         - 故障工单创建（自动关联问题单模块）
         - 故障诊断信息收集（日志、监控数据）
         - 故障修复跟踪
         - 修复后验证测试

     3.4 设备退役:
         - 退役申请与审批
         - 数据清理（用户数据、配置信息）
         - 资源回收（释放所有分配）
         - 退役记录归档
         - 资产转移或报废

  4. 变更管理:
     4.1 变更记录:
         - 自动记录所有变更（状态、配置、资源分配）
         - 记录变更时间、操作人、变更原因
         - 记录变更前后值（old_value, new_value）
         - 变更类型分类（status_change, config_change, resource_change）

     4.2 变更审计:
         - 变更历史查询（按时间、操作人、变更类型）
         - 变更统计分析（变更频率、变更类型分布）
         - 异常变更检测（频繁变更、未授权变更）
         - 变更合规性检查

     4.3 变更追溯:
         - 完整变更链路追踪
         - 变更影响分析（哪些环境受影响）
         - 变更回滚支持（恢复到历史状态）
         - 变更报表生成（月度、季度变更报告）

  5. 资源查询与分配接口:
     5.1 资源查询:
         - 查询可用服务器（按资源需求筛选）
         - 查询服务器详情（完整配置信息）
         - 查询 GPU 列表（按型号、状态筛选）
         - 查询资源使用情况（实时使用率）
         - 批量查询接口（支持多条件组合）

     5.2 资源分配:
         - 分配服务器资源（CPU、内存、GPU）
         - 资源预留（为特定环境预留资源）
         - 资源锁定（防止并发分配冲突）
         - 分配失败回滚

     5.3 资源释放:
         - 释放服务器资源（归还资源池）
         - 资源使用统计更新
         - 释放通知（通知调度模块）
         - 资源清理（清理临时数据）

  6. 集成与同步:
     6.1 监控系统集成:
         - 推送设备状态到监控系统
         - 接收监控告警并更新健康状态
         - 同步设备指标数据

     6.2 调度系统集成:
         - 提供资源查询接口
         - 接收资源分配/释放请求
         - 推送资源变更事件

     6.3 计费系统集成:
         - 提供设备使用记录
         - 推送资源分配/释放事件
         - 提供设备成本信息
```

### 2.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| cmdb_assets | 资产主表 | id, asset_number, operational_status, usage_status |
| cmdb_servers | 服务器详细信息 | asset_id, hostname, ip_address, cpu_cores, memory_total |
| cmdb_gpus | GPU 设备信息 | asset_id, server_id, gpu_index, uuid, status |
| cmdb_change_logs | 变更历史 | asset_id, change_type, old_value, new_value, operator |

### 2.4 对外 API 接口

```go
// CMDB 模块对外提供的核心 API

// 1. 查询可用服务器
GET /api/cmdb/servers/available
Query: cpu, memory, gpu, os_type
Response: []Server

// 2. 查询服务器详情
GET /api/cmdb/servers/:id
Response: ServerDetail

// 3. 分配服务器资源
POST /api/cmdb/servers/:id/allocate
Body: {env_id, cpu, memory, gpu}
Response: {status: "success"}

// 4. 释放服务器资源
POST /api/cmdb/servers/:id/release
Body: {env_id}
Response: {status: "success"}

// 5. 更新设备状态
PUT /api/cmdb/assets/:id/status
Body: {operational_status, reason, operator}
Response: {status: "success"}

// 6. 查询变更历史
GET /api/cmdb/assets/:id/changes
Response: []ChangeLog
```

### 2.5 依赖关系

**依赖的模块：**
- 无（基础设施层，不依赖其他业务模块）

**被依赖的模块：**
- 模块 4：资源调度模块（查询可用设备）
- 模块 3：环境管理模块（获取主机信息）
- 模块 9：监控告警模块（设备健康状态）
- 模块 8：计费管理模块（设备使用记录）

### 2.6 Vue 前端页面

```javascript
// CMDB 模块前端页面列表
pages: [
  {
    path: '/cmdb/assets',
    name: '资产列表',
    component: 'AssetList.vue',
    features: ['列表展示', '筛选', '搜索', '导出']
  },
  {
    path: '/cmdb/assets/:id',
    name: '资产详情',
    component: 'AssetDetail.vue',
    features: ['基本信息', '配置信息', '状态历史', '变更记录']
  },
  {
    path: '/cmdb/servers',
    name: '服务器管理',
    component: 'ServerList.vue',
    features: ['服务器列表', '资源使用率', '健康状态']
  },
  {
    path: '/cmdb/gpus',
    name: 'GPU 管理',
    component: 'GPUList.vue',
    features: ['GPU 列表', '分配状态', '性能监控']
  }
]
```

---

## 3. 模块 2：用户与权限模块

### 3.1 模块概述

**职责：** 管理用户账户、工作空间、权限控制和配额管理，为整个系统提供统一的身份认证和授权服务。

**核心价值：**
- ✅ 多租户隔离
- ✅ 细粒度权限控制
- ✅ 工作空间协作
- ✅ 资源配额管理

### 3.2 核心功能

```yaml
功能列表:
  1. 用户管理:
     - 用户注册与登录
     - 用户信息维护
     - 密码管理
     - 邮箱/手机验证

  2. 工作空间管理:
     - 工作空间创建
     - 成员管理
     - 角色分配
     - 工作空间设置

  3. 权限控制 (RBAC):
     - 角色定义 (owner, admin, member, viewer)
     - 权限分配
     - 资源访问控制
     - 操作审计

  4. 配额管理:
     - CPU/内存/GPU 配额
     - 存储配额
     - 环境数量限制
     - 配额使用监控
```

### 3.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| customers | 用户表 | id, uuid, username, email, account_type, status |
| workspaces | 工作空间表 | id, uuid, owner_id, name, type, member_count |
| workspace_members | 工作空间成员表 | workspace_id, customer_id, role, status |
| resource_quotas | 资源配额表 | customer_id, cpu_quota, memory_quota, gpu_quota |

### 3.4 对外 API 接口

```go
// 用户与权限模块对外提供的核心 API

// 1. 用户注册
POST /api/auth/register
Body: {username, email, password}
Response: {user_id, token}

// 2. 用户登录
POST /api/auth/login
Body: {username, password}
Response: {token, user_info}

// 3. 获取当前用户信息
GET /api/users/me
Response: UserInfo

// 4. 创建工作空间
POST /api/workspaces
Body: {name, description, type}
Response: {workspace_id}

// 5. 添加工作空间成员
POST /api/workspaces/:id/members
Body: {customer_id, role}
Response: {status: "success"}

// 6. 检查权限
POST /api/auth/check-permission
Body: {resource_type, resource_id, action}
Response: {allowed: true/false}

// 7. 查询配额
GET /api/users/me/quota
Response: QuotaInfo

// 8. 查询配额使用情况
GET /api/users/me/quota/usage
Response: QuotaUsage
```

### 3.5 依赖关系

**依赖的模块：**
- 无（基础模块）

**被依赖的模块：**
- 模块 10：网关与认证模块（用户认证）
- 模块 3：环境管理模块（权限检查、配额检查）
- 模块 5：数据与存储模块（权限检查）
- 模块 8：计费管理模块（用户信息）

### 3.6 Vue 前端页面

```javascript
// 用户与权限模块前端页面列表
pages: [
  {
    path: '/login',
    name: '登录',
    component: 'Login.vue',
    features: ['用户名登录', '邮箱登录', '忘记密码']
  },
  {
    path: '/register',
    name: '注册',
    component: 'Register.vue',
    features: ['用户注册', '邮箱验证']
  },
  {
    path: '/profile',
    name: '个人中心',
    component: 'Profile.vue',
    features: ['基本信息', '密码修改', '配额查看']
  },
  {
    path: '/workspaces',
    name: '工作空间列表',
    component: 'WorkspaceList.vue',
    features: ['工作空间列表', '创建工作空间']
  },
  {
    path: '/workspaces/:id',
    name: '工作空间详情',
    component: 'WorkspaceDetail.vue',
    features: ['成员管理', '角色分配', '设置']
  }
]
```

---

## 4. 模块 3：环境管理模块

### 4.1 模块概述

**职责：** 管理用户的开发环境，包括环境创建、启动、停止、删除，以及 SSH/RDP 远程访问配置。

**核心价值：**
- ✅ 一键创建开发环境
- ✅ 支持 Linux 和 Windows
- ✅ 灵活的资源配置
- ✅ 多种访问方式（SSH/RDP/JupyterLab）

### 4.2 核心功能

```yaml
功能列表:
  1. 环境管理:
     - 创建开发环境（Linux/Windows）
     - 启动/停止/重启环境
     - 删除环境
     - 环境列表查询

  2. 访问管理:
     - SSH 端口分配
     - RDP 端口分配（Windows）
     - JupyterLab 配置
     - 访问凭证管理

  3. 资源管理:
     - CPU/内存/GPU 分配
     - 存储挂载
     - 数据集挂载
     - 网络配置

  4. 环境配置:
     - 镜像选择
     - 环境变量设置
     - 启动脚本配置
     - 端口映射
```

### 4.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| environments | 环境表 | id, customer_id, host_id, name, image, status, cpu, memory, gpu |
| port_mappings | 端口映射表 | env_id, service_type, external_port, internal_port, status |
| dataset_usage | 数据集使用记录 | dataset_id, env_id, mount_path, mounted_at |
| environment_configs | 环境配置表 | env_id, config_key, config_value |

### 4.4 对外 API 接口

```go
// 环境管理模块对外提供的核心 API

// 1. 创建环境
POST /api/environments
Body: {name, image, resources: {cpu, memory, gpu}, datasets: []}
Response: {env_id, status}

// 2. 查询环境列表
GET /api/environments
Query: status, workspace_id
Response: []Environment

// 3. 查询环境详情
GET /api/environments/:id
Response: EnvironmentDetail

// 4. 启动环境
POST /api/environments/:id/start
Response: {status: "starting"}

// 5. 停止环境
POST /api/environments/:id/stop
Response: {status: "stopping"}

// 6. 删除环境
DELETE /api/environments/:id
Response: {status: "deleted"}

// 7. 获取访问信息
GET /api/environments/:id/access
Response: {ssh_host, ssh_port, username, password, rdp_port}

// 8. 挂载数据集
POST /api/environments/:id/mount-dataset
Body: {dataset_id, mount_path}
Response: {status: "success"}
```

### 4.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证、权限检查、配额检查）
- 模块 4：资源调度模块（主机选择、资源分配）
- 模块 1：CMDB 设备管理模块（主机信息查询）
- 模块 6：镜像管理模块（镜像信息）
- 模块 5：数据与存储模块（数据集挂载）

**被依赖的模块：**
- 模块 7：训练与推理模块（环境信息）
- 模块 8：计费管理模块（资源使用记录）
- 模块 9：监控告警模块（环境监控）

### 4.6 Vue 前端页面

```javascript
// 环境管理模块前端页面列表
pages: [
  {
    path: '/environments',
    name: '环境列表',
    component: 'EnvironmentList.vue',
    features: ['环境列表', '状态展示', '快速操作']
  },
  {
    path: '/environments/create',
    name: '创建环境',
    component: 'EnvironmentCreate.vue',
    features: ['镜像选择', '资源配置', '数据集选择', '高级配置']
  },
  {
    path: '/environments/:id',
    name: '环境详情',
    component: 'EnvironmentDetail.vue',
    features: ['基本信息', '访问信息', '资源监控', '操作日志']
  },
  {
    path: '/environments/:id/terminal',
    name: 'Web 终端',
    component: 'WebTerminal.vue',
    features: ['SSH 终端', '文件浏览', '命令执行']
  }
]
```

---

## 5. 模块 4：资源调度模块

### 5.1 模块概述

**职责：** 统一管理资源调度，包括主机选择、端口分配、GPU 分配等，为环境创建提供资源保障。

**核心价值：**
- ✅ 智能调度算法
- ✅ 负载均衡
- ✅ 资源利用率优化
- ✅ 多策略支持

### 5.2 核心功能

```yaml
功能列表:
  1. 主机调度:
     - 可用主机查询
     - 负载均衡调度
     - 亲和性调度
     - 优先级调度

  2. 端口管理:
     - 端口池管理
     - 端口分配
     - 端口释放
     - 端口冲突检测

  3. GPU 调度:
     - GPU 可用性查询
     - GPU 分配策略
     - GPU 释放
     - GPU 健康检查

  4. 调度策略:
     - 最少使用策略 (Least Used)
     - 轮询策略 (Round Robin)
     - 随机策略 (Random)
     - 自定义策略
```

### 5.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| scheduler_policies | 调度策略配置 | name, strategy, enabled, priority, config |
| scheduler_history | 调度历史 | env_id, selected_host_id, strategy_used, status |
| port_pools | 端口池 | name, port_range_start, port_range_end, used_ports |
| resource_locks | 资源锁 | resource_type, resource_id, locked_by, expires_at |

### 5.4 对外 API 接口

```go
// 资源调度模块对外提供的核心 API

// 1. 调度资源（内部 API）
POST /internal/scheduler/schedule
Body: {env_id, customer_id, resources: {cpu, memory, gpu}, constraints}
Response: {host_id, gpus: [], ports: {}}

// 2. 释放资源（内部 API）
POST /internal/scheduler/release
Body: {env_id, host_id}
Response: {status: "success"}

// 3. 查询调度历史
GET /api/scheduler/history
Query: env_id, customer_id, start_date, end_date
Response: []ScheduleHistory

// 4. 查询端口使用情况
GET /api/scheduler/ports/usage
Response: {total, used, available}

// 5. 查询调度策略
GET /api/scheduler/policies
Response: []Policy

// 6. 更新调度策略（管理员）
PUT /api/scheduler/policies/:id
Body: {enabled, priority, config}
Response: {status: "success"}
```

### 5.5 依赖关系

**依赖的模块：**
- 模块 1：CMDB 设备管理模块（查询可用主机、GPU 信息）

**被依赖的模块：**
- 模块 3：环境管理模块（资源调度）
- 模块 7：训练与推理模块（资源调度）

### 5.6 Vue 前端页面

```javascript
// 资源调度模块前端页面列表（管理员）
pages: [
  {
    path: '/admin/scheduler/dashboard',
    name: '调度仪表板',
    component: 'SchedulerDashboard.vue',
    features: ['调度统计', '资源使用率', '调度成功率']
  },
  {
    path: '/admin/scheduler/policies',
    name: '调度策略',
    component: 'SchedulerPolicies.vue',
    features: ['策略列表', '策略配置', '策略测试']
  },
  {
    path: '/admin/scheduler/history',
    name: '调度历史',
    component: 'SchedulerHistory.vue',
    features: ['历史记录', '失败分析', '性能分析']
  }
]
```

---

## 6. 模块 5：数据与存储模块

### 6.1 模块概述

**职责：** 管理数据集、模型文件的上传、存储、版本管理和挂载，提供统一的对象存储接口。

**核心价值：**
- ✅ 统一的数据管理
- ✅ 版本控制
- ✅ 高效的文件传输
- ✅ 灵活的挂载方式

### 6.2 核心功能

```yaml
功能列表:
  1. 数据集管理:
     - 数据集创建与上传
     - 数据集版本管理
     - 数据集浏览与下载
     - 数据集共享

  2. 模型管理:
     - 模型上传与存储
     - 模型版本管理
     - 预训练模型库
     - 模型下载

  3. 对象存储:
     - MinIO 集成
     - 预签名 URL 生成
     - 大文件分片上传
     - 断点续传

  4. 文件系统:
     - 数据集挂载到环境
     - 模型挂载到环境
     - 存储配额管理
     - 文件浏览器
```

### 6.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| datasets | 数据集表 | id, uuid, customer_id, name, storage_path, total_size, status |
| dataset_versions | 数据集版本表 | dataset_id, version, storage_path, size, is_default |
| models | 模型表 | id, uuid, customer_id, name, framework, storage_path, status |
| model_versions | 模型版本表 | model_id, version, storage_path, size, metrics |

### 6.4 对外 API 接口

```go
// 数据与存储模块对外提供的核心 API

// 1. 创建数据集
POST /api/datasets
Body: {name, description, visibility, tags}
Response: {dataset_id, storage_path}

// 2. 获取上传凭证
POST /api/datasets/:id/upload-url
Body: {file_name, file_size}
Response: {upload_url, expires_in}

// 3. 完成上传
POST /api/datasets/:id/complete
Body: {files: [{file_name, file_size}]}
Response: {status: "success"}

// 4. 查询数据集列表
GET /api/datasets
Query: visibility, tag
Response: []Dataset

// 5. 浏览数据集文件
GET /api/datasets/:id/files
Query: prefix
Response: {files: []}

// 6. 下载文件
GET /api/datasets/:id/download
Query: file
Response: {download_url, expires_in}

// 7. 创建模型
POST /api/models
Body: {name, framework, description}
Response: {model_id}

// 8. 同步预训练模型
POST /api/models/pretrained/:name/sync
Response: {model_id}
```

### 6.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证、权限检查）

**被依赖的模块：**
- 模块 3：环境管理模块（数据集挂载）
- 模块 7：训练与推理模块（数据集、模型使用）

### 6.6 Vue 前端页面

```javascript
// 数据与存储模块前端页面列表
pages: [
  {
    path: '/datasets',
    name: '数据集列表',
    component: 'DatasetList.vue',
    features: ['数据集列表', '筛选', '搜索', '上传']
  },
  {
    path: '/datasets/upload',
    name: '上传数据集',
    component: 'DatasetUpload.vue',
    features: ['文件选择', '批量上传', '进度显示', '断点续传']
  },
  {
    path: '/datasets/:id',
    name: '数据集详情',
    component: 'DatasetDetail.vue',
    features: ['基本信息', '文件浏览', '版本管理', '使用记录']
  },
  {
    path: '/models',
    name: '模型库',
    component: 'ModelList.vue',
    features: ['模型列表', '预训练模型', '上传模型']
  }
]
```

---

## 7. 模块 6：镜像管理模块

### 7.1 模块概述

**职责：** 管理 Docker 镜像，包括官方镜像库维护和用户自定义镜像构建。

**核心价值：**
- ✅ 丰富的官方镜像
- ✅ 自定义镜像构建
- ✅ 镜像版本管理
- ✅ 快速镜像拉取

### 7.2 核心功能

```yaml
功能列表:
  1. 官方镜像管理:
     - 基础镜像维护
     - 框架镜像（PyTorch/TensorFlow）
     - 镜像更新与发布
     - 镜像文档

  2. 自定义镜像:
     - Dockerfile 编辑器
     - 镜像构建（Kaniko）
     - 构建日志查看
     - 镜像测试

  3. 镜像仓库:
     - Harbor/Registry 集成
     - 镜像推送与拉取
     - 镜像标签管理
     - 镜像清理

  4. 镜像使用:
     - 镜像选择器
     - 镜像搜索
     - 镜像详情查看
     - 使用统计
```

### 7.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| images | 官方镜像表 | id, name, description, category, is_official, size |
| custom_images | 自定义镜像表 | id, uuid, customer_id, name, base_image, dockerfile, status |
| image_builds | 构建历史表 | image_id, build_number, status, build_log |

### 7.4 对外 API 接口

```go
// 镜像管理模块对外提供的核心 API

// 1. 查询官方镜像列表
GET /api/images/official
Query: category
Response: []Image

// 2. 查询镜像详情
GET /api/images/:name
Response: ImageDetail

// 3. 创建自定义镜像
POST /api/images/custom
Body: {name, base_image, dockerfile}
Response: {image_id, status: "building"}

// 4. 查询自定义镜像列表
GET /api/images/custom
Response: []CustomImage

// 5. 查询构建状态
GET /api/images/custom/:id/build-status
Response: {status, build_log}

// 6. 删除自定义镜像
DELETE /api/images/custom/:id
Response: {status: "deleted"}
```

### 7.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证）

**被依赖的模块：**
- 模块 3：环境管理模块（镜像选择）
- 模块 7：训练与推理模块（镜像使用）

### 7.6 Vue 前端页面

```javascript
// 镜像管理模块前端页面列表
pages: [
  {
    path: '/images',
    name: '镜像库',
    component: 'ImageList.vue',
    features: ['官方镜像', '自定义镜像', '镜像搜索']
  },
  {
    path: '/images/custom/create',
    name: '构建镜像',
    component: 'ImageBuilder.vue',
    features: ['基础镜像选择', 'Dockerfile 编辑', '构建配置']
  },
  {
    path: '/images/custom/:id',
    name: '镜像详情',
    component: 'CustomImageDetail.vue',
    features: ['基本信息', '构建历史', '构建日志']
  }
]
```

---

## 8. 模块 7：训练与推理模块

### 8.1 模块概述

**职责：** 管理离线训练任务和推理服务部署，提供完整的模型训练和部署生命周期管理。

**核心价值：**
- ✅ 离线训练任务管理
- ✅ 分布式训练支持
- ✅ 模型部署与服务化
- ✅ 实验管理

### 8.2 核心功能

```yaml
功能列表:
  1. 训练任务管理:
     - 创建训练任务
     - 任务调度与执行
     - 训练日志查看
     - 训练结果保存

  2. 分布式训练:
     - 多机多卡训练
     - 参数服务器模式
     - 数据并行
     - 模型并行

  3. 推理服务:
     - 模型部署
     - 服务管理（启动/停止）
     - 负载均衡
     - API 接口

  4. 实验管理:
     - 实验记录
     - 超参数管理
     - 指标对比
     - 可视化
```

### 8.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| training_jobs | 训练任务表 | id, uuid, customer_id, name, status, config, result |
| inference_services | 推理服务表 | id, uuid, model_id, name, status, endpoint, replicas |
| experiments | 实验表 | id, customer_id, name, description, hyperparameters |
| experiment_runs | 实验运行记录 | experiment_id, run_number, metrics, artifacts |

### 8.4 对外 API 接口

```go
// 训练与推理模块对外提供的核心 API

// 1. 创建训练任务
POST /api/training/jobs
Body: {name, image, script, datasets, resources}
Response: {job_id, status}

// 2. 查询训练任务列表
GET /api/training/jobs
Response: []TrainingJob

// 3. 查询任务详情
GET /api/training/jobs/:id
Response: JobDetail

// 4. 停止训练任务
POST /api/training/jobs/:id/stop
Response: {status: "stopping"}

// 5. 查询训练日志
GET /api/training/jobs/:id/logs
Response: {logs: "..."}

// 6. 部署推理服务
POST /api/inference/services
Body: {name, model_id, replicas, resources}
Response: {service_id, endpoint}

// 7. 查询推理服务列表
GET /api/inference/services
Response: []InferenceService

// 8. 调用推理接口
POST /api/inference/services/:id/predict
Body: {input_data}
Response: {predictions}
```

### 8.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证）
- 模块 4：资源调度模块（资源分配）
- 模块 5：数据与存储模块（数据集、模型）
- 模块 6：镜像管理模块（镜像使用）

**被依赖的模块：**
- 模块 8：计费管理模块（资源使用记录）
- 模块 9：监控告警模块（任务监控）

### 8.6 Vue 前端页面

```javascript
// 训练与推理模块前端页面列表
pages: [
  {
    path: '/training/jobs',
    name: '训练任务',
    component: 'TrainingJobList.vue',
    features: ['任务列表', '创建任务', '状态监控']
  },
  {
    path: '/training/jobs/create',
    name: '创建训练任务',
    component: 'TrainingJobCreate.vue',
    features: ['配置选择', '脚本上传', '参数设置']
  },
  {
    path: '/training/jobs/:id',
    name: '任务详情',
    component: 'TrainingJobDetail.vue',
    features: ['基本信息', '日志查看', '结果下载', '指标可视化']
  },
  {
    path: '/inference/services',
    name: '推理服务',
    component: 'InferenceServiceList.vue',
    features: ['服务列表', '部署服务', '服务监控']
  }
]
```

---

## 9. 模块 8：计费管理模块

### 9.1 模块概述

**职责：** 管理资源使用计费、账户余额、账单生成和支付集成。

**核心价值：**
- ✅ 精确的资源计费
- ✅ 灵活的计费策略
- ✅ 自动账单生成
- ✅ 多种支付方式

### 9.2 核心功能

```yaml
功能列表:
  1. 资源计费:
     - CPU/内存/GPU 按时计费
     - 存储按量计费
     - 网络流量计费
     - 计费规则配置

  2. 账户管理:
     - 账户余额
     - 充值记录
     - 消费记录
     - 余额预警

  3. 账单管理:
     - 账单生成
     - 账单查询
     - 账单导出
     - 发票管理

  4. 支付集成:
     - 支付宝
     - 微信支付
     - 银行卡支付
     - 企业转账
```

### 9.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| billing_records | 计费记录表 | customer_id, env_id, resource_type, quantity, amount |
| accounts | 账户表 | customer_id, balance, credit_limit, status |
| invoices | 账单表 | customer_id, billing_period, total_amount, status |
| payments | 支付记录表 | customer_id, amount, payment_method, status |

### 9.4 对外 API 接口

```go
// 计费管理模块对外提供的核心 API

// 1. 查询账户余额
GET /api/billing/account
Response: {balance, credit_limit}

// 2. 查询计费记录
GET /api/billing/records
Query: start_date, end_date, resource_type
Response: []BillingRecord

// 3. 查询账单列表
GET /api/billing/invoices
Response: []Invoice

// 4. 查询账单详情
GET /api/billing/invoices/:id
Response: InvoiceDetail

// 5. 创建充值订单
POST /api/billing/recharge
Body: {amount, payment_method}
Response: {order_id, payment_url}

// 6. 查询支付状态
GET /api/billing/payments/:order_id
Response: {status, paid_at}
```

### 9.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户信息）
- 模块 3：环境管理模块（资源使用记录）
- 模块 7：训练与推理模块（任务使用记录）

**被依赖的模块：**
- 无

### 9.6 Vue 前端页面

```javascript
// 计费管理模块前端页面列表
pages: [
  {
    path: '/billing/overview',
    name: '计费概览',
    component: 'BillingOverview.vue',
    features: ['账户余额', '本月消费', '消费趋势']
  },
  {
    path: '/billing/records',
    name: '计费记录',
    component: 'BillingRecords.vue',
    features: ['记录列表', '筛选', '导出']
  },
  {
    path: '/billing/invoices',
    name: '账单管理',
    component: 'InvoiceList.vue',
    features: ['账单列表', '账单详情', '下载发票']
  },
  {
    path: '/billing/recharge',
    name: '账户充值',
    component: 'Recharge.vue',
    features: ['充值金额', '支付方式', '充值记录']
  }
]
```

---

## 10. 模块 9：监控告警模块

### 10.1 模块概述

**职责：** 监控系统资源使用、性能指标和业务指标，提供实时告警和可视化展示。

**核心价值：**
- ✅ 全方位监控
- ✅ 实时告警
- ✅ 性能分析
- ✅ 故障定位

### 10.2 核心功能

```yaml
功能列表:
  1. 资源监控:
     - 主机资源监控（CPU/内存/磁盘/网络）
     - GPU 监控（使用率/显存/温度/功耗）
     - 环境资源监控
     - 存储使用监控

  2. 性能监控:
     - API 响应时间
     - 数据库性能
     - 任务执行时间
     - 系统吞吐量

  3. 告警管理:
     - 告警规则配置
     - 告警触发与通知
     - 告警历史查询
     - 告警静默

  4. 日志管理:
     - 应用日志收集
     - 系统日志收集
     - 日志查询与分析
     - 日志归档
```

### 10.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| alert_rules | 告警规则表 | name, metric, threshold, severity, enabled |
| alert_history | 告警历史表 | rule_id, triggered_at, resolved_at, status, message |
| system_metrics | 系统指标表 | metric_name, metric_value, tags, collected_at |

### 10.4 对外 API 接口

```go
// 监控告警模块对外提供的核心 API

// 1. 查询主机监控数据
GET /api/monitoring/hosts/:id/metrics
Query: metric_name, start_time, end_time
Response: {metrics: []}

// 2. 查询 GPU 监控数据
GET /api/monitoring/gpus/:id/metrics
Query: start_time, end_time
Response: {metrics: []}

// 3. 查询环境监控数据
GET /api/monitoring/environments/:id/metrics
Query: start_time, end_time
Response: {metrics: []}

// 4. 查询告警规则
GET /api/monitoring/alert-rules
Response: []AlertRule

// 5. 创建告警规则
POST /api/monitoring/alert-rules
Body: {name, metric, threshold, severity}
Response: {rule_id}

// 6. 查询告警历史
GET /api/monitoring/alerts
Query: status, severity, start_date, end_date
Response: []Alert

// 7. 确认告警
POST /api/monitoring/alerts/:id/acknowledge
Response: {status: "acknowledged"}
```

### 10.5 依赖关系

**依赖的模块：**
- 模块 1：CMDB 设备管理模块（设备监控数据）
- 模块 3：环境管理模块（环境监控数据）
- 模块 7：训练与推理模块（任务监控数据）

**被依赖的模块：**
- 无

### 10.6 Vue 前端页面

```javascript
// 监控告警模块前端页面列表
pages: [
  {
    path: '/monitoring/dashboard',
    name: '监控仪表板',
    component: 'MonitoringDashboard.vue',
    features: ['系统概览', '资源使用率', '实时告警']
  },
  {
    path: '/monitoring/hosts',
    name: '主机监控',
    component: 'HostMonitoring.vue',
    features: ['主机列表', '资源图表', '历史数据']
  },
  {
    path: '/monitoring/gpus',
    name: 'GPU 监控',
    component: 'GPUMonitoring.vue',
    features: ['GPU 列表', '使用率图表', '温度监控']
  },
  {
    path: '/monitoring/alerts',
    name: '告警管理',
    component: 'AlertManagement.vue',
    features: ['告警列表', '告警规则', '告警历史']
  }
]
```

---

## 11. 模块 10：网关与认证模块

### 11.1 模块概述

**职责：** 作为系统统一入口，提供 API 网关、身份认证、请求限流和安全防护功能。

**核心价值：**
- ✅ 统一入口
- ✅ 安全认证
- ✅ 流量控制
- ✅ 请求路由

### 11.2 核心功能

```yaml
功能列表:
  1. API 网关:
     - 请求路由
     - 负载均衡
     - 协议转换
     - API 版本管理

  2. 身份认证:
     - JWT Token 认证
     - OAuth2 集成
     - SSO 单点登录
     - API Key 认证

  3. 访问控制:
     - 权限验证
     - IP 白名单
     - 跨域配置（CORS）
     - 请求签名验证

  4. 流量控制:
     - 请求限流
     - 熔断降级
     - 超时控制
     - 重试机制
```

### 11.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| api_keys | API 密钥表 | customer_id, key, secret, status, expires_at |
| access_logs | 访问日志表 | customer_id, api_path, method, status_code, response_time |
| rate_limits | 限流配置表 | customer_id, api_path, limit_per_minute, limit_per_hour |

### 11.4 对外 API 接口

```go
// 网关与认证模块对外提供的核心 API

// 1. 用户登录（获取 Token）
POST /api/auth/login
Body: {username, password}
Response: {token, expires_in, refresh_token}

// 2. 刷新 Token
POST /api/auth/refresh
Body: {refresh_token}
Response: {token, expires_in}

// 3. 退出登录
POST /api/auth/logout
Response: {status: "success"}

// 4. 创建 API Key
POST /api/auth/api-keys
Body: {name, expires_in}
Response: {api_key, api_secret}

// 5. 查询 API Key 列表
GET /api/auth/api-keys
Response: []APIKey

// 6. 删除 API Key
DELETE /api/auth/api-keys/:id
Response: {status: "deleted"}
```

### 11.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证、权限验证）

**被依赖的模块：**
- 所有业务模块（统一入口）

### 11.6 技术实现

```go
// 网关中间件示例
package gateway

// JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()
            return
        }

        // 验证 Token
        claims, err := ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Token 无效"})
            c.Abort()
            return
        }

        // 设置用户信息到上下文
        c.Set("customer_id", claims.CustomerID)
        c.Next()
    }
}

// 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        customerID := c.GetInt64("customer_id")
        apiPath := c.Request.URL.Path

        // 检查限流
        if !CheckRateLimit(customerID, apiPath) {
            c.JSON(429, gin.H{"error": "请求过于频繁"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

---

## 12. 模块 11：制品管理模块

### 12.1 模块概述

**职责：** 管理软件包、依赖库、编译产物等制品的存储、版本控制和分发，为开发环境和训练任务提供统一的制品仓库。

**核心价值：**
- ✅ 统一的制品管理
- ✅ 版本控制和追溯
- ✅ 加速依赖下载
- ✅ 私有包管理

### 12.2 核心功能

```yaml
功能列表:
  1. 软件包管理:
     - Python 包管理（PyPI 私有仓库）
     - NPM 包管理（Node.js）
     - Maven/Gradle 包管理（Java）
     - Docker 镜像管理（集成 Harbor）

  2. 版本控制:
     - 语义化版本管理
     - 版本发布与回滚
     - 版本依赖关系
     - 版本标签管理

  3. 制品仓库:
     - 制品上传与下载
     - 制品缓存加速
     - 制品权限控制
     - 制品清理策略

  4. 依赖管理:
     - 依赖解析
     - 依赖锁定（lock file）
     - 依赖安全扫描
     - 依赖更新提醒
```

### 12.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|------------|
| artifacts | 制品表 | id, name, type, version, storage_path, size |
| artifact_versions | 制品版本表 | artifact_id, version, release_date, changelog |
| artifact_dependencies | 依赖关系表 | artifact_id, dependency_id, version_constraint |
| package_repositories | 仓库配置表 | name, type, url, credentials, enabled |

### 12.4 对外 API 接口

```go
// 制品管理模块对外提供的核心 API

// 1. 上传制品
POST /api/artifacts
Body: {name, type, version, file}
Response: {artifact_id, download_url}

// 2. 查询制品列表
GET /api/artifacts
Query: type, keyword
Response: []Artifact

// 3. 下载制品
GET /api/artifacts/:id/download
Response: Binary file

// 4. 查询制品版本
GET /api/artifacts/:id/versions
Response: []Version

// 5. 配置私有仓库
POST /api/repositories
Body: {name, type, url, credentials}
Response: {repository_id}
```

### 12.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（权限控制）
- 模块 5：数据与存储模块（对象存储）

**被依赖的模块：**
- 模块 3：环境管理模块（依赖安装）
- 模块 7：训练与推理模块（依赖管理）

### 12.6 Vue 前端页面

```javascript
pages: [
  {
    path: '/artifacts',
    name: '制品列表',
    component: 'ArtifactList.vue',
    features: ['制品列表', '搜索', '上传']
  },
  {
    path: '/artifacts/:id',
    name: '制品详情',
    component: 'ArtifactDetail.vue',
    features: ['版本列表', '依赖关系', '下载']
  },
  {
    path: '/repositories',
    name: '仓库管理',
    component: 'RepositoryList.vue',
    features: ['仓库配置', '镜像源设置']
  }
]
```

---

## 13. 模块 12：问题单管理模块

### 13.1 模块概述

**职责：** 提供完整的问题跟踪和工单管理系统，支持 Bug 报告、问题处理、工单流转和统计分析。

**核心价值：**
- ✅ 问题全生命周期管理
- ✅ 工单流转和协作
- ✅ 问题统计和分析
- ✅ SLA 管理

### 13.2 核心功能

```yaml
功能列表:
  1. 问题管理:
     - 问题创建与提交
     - 问题分类（Bug/Feature/Task）
     - 问题优先级管理
     - 问题状态跟踪

  2. 工单流转:
     - 工单分配
     - 工单转派
     - 工单关闭
     - 工单重新打开

  3. 协作功能:
     - 评论和讨论
     - @提醒
     - 附件上传
     - 关联问题

  4. 统计分析:
     - 问题趋势分析
     - 处理时效统计
     - 团队工作量统计
     - SLA 达成率
```

### 13.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| issues | 问题表 | id, uuid, customer_id, title, type, priority, status, assignee_id |
| issue_comments | 评论表 | id, issue_id, customer_id, content, created_at |
| issue_attachments | 附件表 | id, issue_id, file_name, file_path, file_size |
| issue_status_history | 状态历史表 | id, issue_id, old_status, new_status, operator_id, changed_at |
| issue_labels | 标签表 | id, name, color, description |
| issue_label_relations | 问题标签关联表 | issue_id, label_id |

### 13.4 对外 API 接口

```go
// 问题单管理模块对外提供的核心 API

// 1. 创建问题
POST /api/issues
Body: {title, description, type, priority, assignee_id, labels}
Response: {issue_id}

// 2. 查询问题列表
GET /api/issues
Query: status, type, priority, assignee_id, keyword
Response: []Issue

// 3. 查询问题详情
GET /api/issues/:id
Response: IssueDetail

// 4. 更新问题
PUT /api/issues/:id
Body: {title, description, priority, status}
Response: {status: "success"}

// 5. 分配问题
POST /api/issues/:id/assign
Body: {assignee_id}
Response: {status: "success"}

// 6. 添加评论
POST /api/issues/:id/comments
Body: {content}
Response: {comment_id}

// 7. 上传附件
POST /api/issues/:id/attachments
Body: multipart/form-data
Response: {attachment_id, file_url}

// 8. 关闭问题
POST /api/issues/:id/close
Body: {resolution, comment}
Response: {status: "success"}
```

### 13.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证、权限检查）
- 模块 5：数据与存储模块（附件存储）

**被依赖的模块：**
- 无

### 13.6 Vue 前端页面

```javascript
// 问题单管理模块前端页面列表
pages: [
  {
    path: '/issues',
    name: '问题列表',
    component: 'IssueList.vue',
    features: ['问题列表', '筛选', '搜索', '创建问题']
  },
  {
    path: '/issues/create',
    name: '创建问题',
    component: 'IssueCreate.vue',
    features: ['问题表单', '附件上传', '标签选择']
  },
  {
    path: '/issues/:id',
    name: '问题详情',
    component: 'IssueDetail.vue',
    features: ['基本信息', '评论列表', '状态历史', '附件列表']
  },
  {
    path: '/issues/dashboard',
    name: '问题统计',
    component: 'IssueDashboard.vue',
    features: ['问题趋势', '处理时效', '团队工作量']
  }
]
```

---

## 14. 模块 13：需求单管理模块

### 14.1 模块概述

**职责：** 管理产品需求的收集、评审、优先级排序和开发跟踪，支持敏捷开发流程。

**核心价值：**
- ✅ 需求全流程管理
- ✅ 需求优先级排序
- ✅ 开发进度跟踪
- ✅ 需求变更管理

### 14.2 核心功能

```yaml
功能列表:
  1. 需求管理:
     - 需求创建与提交
     - 需求分类（功能/优化/重构）
     - 需求描述和验收标准
     - 需求附件管理

  2. 需求评审:
     - 需求评审流程
     - 评审意见记录
     - 需求评分
     - 需求批准/拒绝

  3. 优先级管理:
     - 优先级设置（P0/P1/P2）
     - 需求排期
     - 里程碑管理
     - Sprint 规划

  4. 开发跟踪:
     - 需求状态跟踪
     - 开发进度更新
     - 需求关联代码
     - 需求验收
```

### 14.3 核心数据表

| 表名 | 说明 | 关键字段 |
|------|------|---------|
| requirements | 需求表 | id, uuid, customer_id, title, type, priority, status, owner_id |
| requirement_reviews | 评审记录表 | id, requirement_id, reviewer_id, score, comment, status |
| requirement_milestones | 里程碑表 | id, name, start_date, end_date, status |
| requirement_sprints | Sprint表 | id, milestone_id, name, start_date, end_date, capacity |
| requirement_sprint_items | Sprint需求关联表 | sprint_id, requirement_id, story_points |
| requirement_attachments | 附件表 | id, requirement_id, file_name, file_path, file_size |

### 14.4 对外 API 接口

```go
// 需求单管理模块对外提供的核心 API

// 1. 创建需求
POST /api/requirements
Body: {title, description, type, priority, acceptance_criteria}
Response: {requirement_id}

// 2. 查询需求列表
GET /api/requirements
Query: status, type, priority, milestone_id, keyword
Response: []Requirement

// 3. 查询需求详情
GET /api/requirements/:id
Response: RequirementDetail

// 4. 更新需求
PUT /api/requirements/:id
Body: {title, description, priority, status}
Response: {status: "success"}

// 5. 提交评审
POST /api/requirements/:id/submit-review
Response: {status: "success"}

// 6. 评审需求
POST /api/requirements/:id/review
Body: {score, comment, approved}
Response: {status: "success"}

// 7. 创建 Sprint
POST /api/sprints
Body: {name, milestone_id, start_date, end_date, capacity}
Response: {sprint_id}

// 8. 添加需求到 Sprint
POST /api/sprints/:id/items
Body: {requirement_id, story_points}
Response: {status: "success"}
```

### 14.5 依赖关系

**依赖的模块：**
- 模块 2：用户与权限模块（用户认证、权限检查）
- 模块 5：数据与存储模块（附件存储）

**被依赖的模块：**
- 无

### 14.6 Vue 前端页面

```javascript
// 需求单管理模块前端页面列表
pages: [
  {
    path: '/requirements',
    name: '需求列表',
    component: 'RequirementList.vue',
    features: ['需求列表', '筛选', '搜索', '创建需求']
  },
  {
    path: '/requirements/create',
    name: '创建需求',
    component: 'RequirementCreate.vue',
    features: ['需求表单', '验收标准', '附件上传']
  },
  {
    path: '/requirements/:id',
    name: '需求详情',
    component: 'RequirementDetail.vue',
    features: ['基本信息', '评审记录', '开发进度', '关联代码']
  },
  {
    path: '/requirements/roadmap',
    name: '需求路线图',
    component: 'RequirementRoadmap.vue',
    features: ['里程碑视图', '甘特图', '优先级排序']
  },
  {
    path: '/sprints',
    name: 'Sprint 管理',
    component: 'SprintList.vue',
    features: ['Sprint 列表', '燃尽图', '速度图表']
  }
]
```

---

## 15. 模块依赖关系

### 12.1 依赖关系图

```
┌─────────────────────────────────────────────────────────────────┐
│                     模块依赖关系图                                 │
└─────────────────────────────────────────────────────────────────┘

层级 0（基础层）：
┌──────────────┐
│ 模块 1:      │
│ CMDB设备管理  │  ← 无依赖，基础设施层
└──────────────┘

层级 1（核心服务层）：
┌──────────────┐  ┌──────────────┐
│ 模块 2:      │  │ 模块 10:     │
│ 用户与权限    │  │ 网关与认证    │  ← 依赖：模块 2
└──────────────┘  └──────────────┘

层级 2（调度与资源层）：
┌──────────────┐
│ 模块 4:      │
│ 资源调度      │  ← 依赖：模块 1
└──────────────┘

层级 3（业务功能层）：
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ 模块 3:      │  │ 模块 5:      │  │ 模块 6:      │
│ 环境管理      │  │ 数据与存储    │  │ 镜像管理      │
│              │  │              │  │              │
│ 依赖：       │  │ 依赖：       │  │ 依赖：       │
│ - 模块 2     │  │ - 模块 2     │  │ - 模块 2     │
│ - 模块 4     │  └──────────────┘  └──────────────┘
│ - 模块 1     │
│ - 模块 6     │
│ - 模块 5     │
└──────────────┘

层级 4（高级功能层）：
┌──────────────┐
│ 模块 7:      │
│ 训练与推理    │
│              │
│ 依赖：       │
│ - 模块 2     │
│ - 模块 4     │
│ - 模块 5     │
│ - 模块 6     │
└──────────────┘

层级 5（支撑服务层）：
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ 模块 8:      │  │ 模块 9:      │  │ 模块 11:     │
│ 计费管理      │  │ 监控告警      │  │ 制品管理      │
│              │  │              │  │              │
│ 依赖：       │  │ 依赖：       │  │ 依赖：       │
│ - 模块 2     │  │ - 模块 1     │  │ - 模块 2     │
│ - 模块 3     │  │ - 模块 3     │  │ - 模块 5     │
│ - 模块 7     │  │ - 模块 7     │  │              │
└──────────────┘  └──────────────┘  └──────────────┘

层级 6（辅助服务层）：
┌──────────────┐  ┌──────────────┐
│ 模块 12:     │  │ 模块 13:     │
│ 问题单管理    │  │ 需求单管理    │
│              │  │              │
│ 依赖：       │  │ 依赖：       │
│ - 模块 2     │  │ - 模块 2     │
│ - 模块 5     │  │ - 模块 5     │
└──────────────┘  └──────────────┘
```

### 12.2 依赖关系矩阵

| 模块 | 依赖的模块 | 被依赖的模块 |
|------|-----------|-------------|
| 模块 1: CMDB 设备管理 | 无 | 模块 3, 4, 9 |
| 模块 2: 用户与权限 | 无 | 模块 3, 5, 6, 7, 8, 10, 11, 12, 13 |
| 模块 3: 环境管理 | 模块 1, 2, 4, 5, 6 | 模块 7, 8, 9 |
| 模块 4: 资源调度 | 模块 1 | 模块 3, 7 |
| 模块 5: 数据与存储 | 模块 2 | 模块 3, 7, 11, 12, 13 |
| 模块 6: 镜像管理 | 模块 2 | 模块 3, 7 |
| 模块 7: 训练与推理 | 模块 2, 4, 5, 6 | 模块 8, 9 |
| 模块 8: 计费管理 | 模块 2, 3, 7 | 无 |
| 模块 9: 监控告警 | 模块 1, 3, 7 | 无 |
| 模块 10: 网关与认证 | 模块 2 | 所有模块 |
| 模块 11: 制品管理 | 模块 2, 5 | 模块 3, 7 |
| 模块 12: 问题单管理 | 模块 2, 5 | 无 |
| 模块 13: 需求单管理 | 模块 2, 5 | 无 |

### 12.3 模块间通信方式

```yaml
通信方式:
  1. 同步调用（HTTP/gRPC）:
     - 模块 3 → 模块 4（资源调度）
     - 模块 3 → 模块 1（查询主机信息）
     - 模块 7 → 模块 5（获取数据集）

  2. 异步消息（消息队列）:
     - 模块 3 → 模块 8（计费事件）
     - 模块 7 → 模块 8（任务计费）
     - 模块 9 → 告警通知

  3. 数据库共享:
     - 所有模块共享 PostgreSQL
     - 通过外键关联

  4. 缓存共享:
     - 所有模块共享 Redis
     - 用户会话、配置缓存
```

---

## 13. 开发优先级

### 13.1 优先级定义

```yaml
优先级说明:
  P0 (必须): MVP 核心功能，系统运行的基础
  P1 (重要): 扩展功能，提升用户体验
  P2 (可选): 高级功能，增强竞争力
```

### 13.2 模块优先级分配

| 模块编号 | 模块名称 | 优先级 | 开发周期 | 依赖模块 |
|---------|---------|--------|---------|---------|
| 模块 1 | CMDB 设备管理模块 | **P0** | 2-3 周 | 无 |
| 模块 2 | 用户与权限模块 | **P0** | 2-3 周 | 无 |
| 模块 10 | 网关与认证模块 | **P0** | 1-2 周 | 模块 2 |
| 模块 4 | 资源调度模块 | **P0** | 2-3 周 | 模块 1 |
| 模块 3 | 环境管理模块 | **P0** | 3-4 周 | 模块 1, 2, 4 |
| 模块 6 | 镜像管理模块 | **P1** | 2-3 周 | 模块 2 |
| 模块 5 | 数据与存储模块 | **P1** | 2-3 周 | 模块 2 |
| 模块 11 | 制品管理模块 | **P1** | 2-3 周 | 模块 2, 5 |
| 模块 9 | 监控告警模块 | **P1** | 2-3 周 | 模块 1, 3 |
| 模块 8 | 计费管理模块 | **P1** | 2-3 周 | 模块 2, 3 |
| 模块 7 | 训练与推理模块 | **P2** | 3-4 周 | 模块 2, 4, 5, 6 |
| 模块 12 | 问题单管理模块 | **P2** | 2-3 周 | 模块 2, 5 |
| 模块 13 | 需求单管理模块 | **P2** | 2-3 周 | 模块 2, 5 |

### 13.3 开发阶段划分

#### 阶段 1：MVP 核心功能（10-12 周）

**目标：** 实现基本的 GPU 云平台功能，用户可以创建和使用开发环境

**包含模块：**
```yaml
第 1-3 周: 基础设施搭建
  - 数据库设计与初始化
  - 对象存储部署（MinIO）
  - 缓存服务部署（Redis）
  - 基础监控搭建

第 4-6 周: 核心模块开发
  - 模块 1: CMDB 设备管理模块（主机注册、GPU 发现）
  - 模块 2: 用户与权限模块（用户注册、登录、基础权限）
  - 模块 10: 网关与认证模块（JWT 认证、API 网关）

第 7-9 周: 调度与环境模块
  - 模块 4: 资源调度模块（主机选择、端口分配）
  - 模块 3: 环境管理模块（Linux 环境创建、SSH 访问）

第 10-12 周: 前端与测试
  - 前端基础界面（登录、环境列表、环境创建）
  - 集成测试
  - 性能优化
```

**交付物：**
- ✅ 用户可以注册和登录
- ✅ 管理员可以注册 Linux 主机和 GPU
- ✅ 用户可以创建 Linux 开发环境（Docker 容器）
- ✅ 用户可以通过 SSH 访问环境
- ✅ 基础的资源监控

#### 阶段 2：功能扩展（6-8 周）

**目标：** 增强平台功能，提升用户体验

**包含模块：**
```yaml
第 13-15 周: 数据与镜像管理
  - 模块 6: 镜像管理模块（官方镜像库、自定义镜像构建）
  - 模块 5: 数据与存储模块（数据集上传、模型管理）

第 16-18 周: 监控与计费
  - 模块 9: 监控告警模块（完善监控、告警规则）
  - 模块 8: 计费管理模块（资源计费、账单生成）

第 19-20 周: 前端完善与测试
  - 完善前端界面（数据集管理、镜像管理、监控仪表板）
  - 集成测试
  - 用户体验优化
```

**交付物：**
- ✅ 用户可以上传和管理数据集
- ✅ 用户可以使用官方镜像或构建自定义镜像
- ✅ 完善的资源监控和告警系统
- ✅ 资源使用计费和账单查询
- ✅ 完整的前端管理界面

#### 阶段 3：高级功能（3-4 周）

**目标：** 增加训练和推理功能，增强平台竞争力

**包含模块：**
```yaml
第 21-24 周: 训练与推理
  - 模块 7: 训练与推理模块（离线训练任务、推理服务部署）
  - 分布式训练支持（可选）
  - 实验管理功能
```

**交付物：**
- ✅ 用户可以提交离线训练任务
- ✅ 用户可以部署推理服务
- ✅ 实验管理和结果对比
- ✅ 训练日志和指标可视化

### 13.4 开发顺序建议

**按依赖关系的开发顺序：**

```
第一批（并行开发，无依赖）:
├── 模块 1: CMDB 设备管理模块
└── 模块 2: 用户与权限模块

第二批（依赖第一批）:
├── 模块 10: 网关与认证模块（依赖模块 2）
└── 模块 4: 资源调度模块（依赖模块 1）

第三批（依赖第一、二批）:
├── 模块 3: 环境管理模块（依赖模块 1, 2, 4）
├── 模块 5: 数据与存储模块（依赖模块 2）
└── 模块 6: 镜像管理模块（依赖模块 2）

第四批（依赖第三批）:
├── 模块 11: 制品管理模块（依赖模块 2, 5）
├── 模块 9: 监控告警模块（依赖模块 1, 3）
├── 模块 8: 计费管理模块（依赖模块 2, 3）
└── 模块 7: 训练与推理模块（依赖模块 2, 4, 5, 6）

第五批（辅助功能，可选）:
├── 模块 12: 问题单管理模块（依赖模块 2, 5）
└── 模块 13: 需求单管理模块（依赖模块 2, 5）
```

### 13.5 资源分配建议

**团队配置建议：**

```yaml
后端开发团队（4-5 人）:
  - 后端架构师 x1: 负责整体架构设计、技术选型
  - Go 后端工程师 x2-3: 负责各模块后端开发
  - DevOps 工程师 x1: 负责基础设施、CI/CD、监控

前端开发团队（2-3 人）:
  - 前端架构师 x1: 负责前端架构设计
  - Vue 前端工程师 x1-2: 负责前端页面开发

测试团队（1-2 人）:
  - 测试工程师 x1-2: 负责功能测试、性能测试

产品与设计（1-2 人）:
  - 产品经理 x1: 负责需求管理、产品规划
  - UI/UX 设计师 x1: 负责界面设计（可选）
```

**模块分工建议：**

```yaml
后端工程师 A:
  - 模块 1: CMDB 设备管理模块
  - 模块 4: 资源调度模块
  - 模块 9: 监控告警模块

后端工程师 B:
  - 模块 2: 用户与权限模块
  - 模块 10: 网关与认证模块
  - 模块 8: 计费管理模块

后端工程师 C:
  - 模块 3: 环境管理模块
  - 模块 6: 镜像管理模块
  - 模块 7: 训练与推理模块

后端工程师 D（可选）:
  - 模块 5: 数据与存储模块
  - 协助其他模块开发
```

### 13.6 关键里程碑

**里程碑 1：基础平台可用（第 12 周）**

```yaml
目标: MVP 上线，核心功能可用
验收标准:
  - ✅ 用户可以注册、登录系统
  - ✅ 管理员可以注册 Linux 主机和 GPU 设备
  - ✅ 用户可以创建 Linux 开发环境（Docker 容器）
  - ✅ 用户可以通过 SSH 访问开发环境
  - ✅ 基础的资源监控功能
  - ✅ 前端基础界面可用
```

**里程碑 2：功能完善（第 20 周）**

```yaml
目标: 平台功能完善，用户体验提升
验收标准:
  - ✅ 数据集上传和管理功能
  - ✅ 官方镜像库和自定义镜像构建
  - ✅ 完善的监控告警系统
  - ✅ 资源计费和账单管理
  - ✅ 完整的前端管理界面
  - ✅ 工作空间和团队协作功能
```

**里程碑 3：高级功能上线（第 24 周）**

```yaml
目标: 增加训练和推理功能，增强竞争力
验收标准:
  - ✅ 离线训练任务提交和管理
  - ✅ 推理服务部署和管理
  - ✅ 实验管理和结果对比
  - ✅ 训练日志和指标可视化
  - ✅ 分布式训练支持（可选）
```

### 13.7 风险与注意事项

**技术风险：**

```yaml
1. GPU 资源调度复杂性:
   - 风险: GPU 分配策略可能不够优化，导致资源利用率低
   - 缓解: 在模块 4 开发时充分测试各种调度策略，收集实际使用数据优化

2. 多租户隔离安全性:
   - 风险: Docker 容器隔离可能存在安全漏洞
   - 缓解: 使用最新版本 Docker，配置安全策略，定期安全审计

3. 端口管理冲突:
   - 风险: 端口分配可能出现冲突或耗尽
   - 缓解: 实现完善的端口池管理，支持动态扩展端口范围

4. 数据存储性能:
   - 风险: 大规模数据集上传下载可能影响性能
   - 缓解: 使用对象存储（MinIO），支持分片上传和 CDN 加速
```

**项目风险：**

```yaml
1. 开发周期延期:
   - 风险: 模块间依赖可能导致开发阻塞
   - 缓解: 严格按照依赖顺序开发，提前定义好模块接口

2. 需求变更:
   - 风险: 用户需求变化导致返工
   - 缓解: MVP 阶段聚焦核心功能，后续迭代增加新功能

3. 团队协作:
   - 风险: 多人协作可能出现代码冲突
   - 缓解: 使用 Git 分支管理，定期代码审查，统一编码规范
```

### 13.8 总结

本文档完成了 RemoteGPU 系统的完整模块划分，将整个系统划分为 **15 个核心模块**：

**基础层（P0）：**
1. 模块 1: CMDB 设备管理模块 - 硬件资产管理的基石
2. 模块 2: 用户与权限模块 - 多租户和权限控制
3. 模块 10: 网关与认证模块 - 统一入口和安全保障

**核心服务层（P0）：**
4. 模块 4: 资源调度模块 - 智能调度和资源分配
5. 模块 3: 环境管理模块 - 开发环境管理

**业务功能层（P1）：**
6. 模块 6: 镜像管理模块 - 镜像库和构建
7. 模块 5: 数据与存储模块 - 数据集和模型管理
8. 模块 11: 制品管理模块 - 软件包和依赖管理

**支撑服务层（P1）：**
9. 模块 9: 监控告警模块 - 系统监控和告警
10. 模块 8: 计费管理模块 - 资源计费和账单

**高级功能层（P2）：**
11. 模块 7: 训练与推理模块 - AI 训练和推理服务

**辅助服务层（P2）：**
12. 模块 12: 问题单管理模块 - Bug 跟踪和工单管理
13. 模块 13: 需求单管理模块 - 需求收集和敏捷开发

**通知与集成层（P1）：**
14. 模块 14: 通知管理模块 - 多渠道通知和消息推送
15. 模块 15: Webhook 管理模块 - 事件回调和第三方集成

**开发建议：**
- **MVP 阶段（10-12 周）**：聚焦 P0 模块，实现核心功能
- **扩展阶段（6-8 周）**：完善 P1 模块，提升用户体验
- **高级阶段（3-4 周）**：增加 P2 模块，增强竞争力
- **辅助功能（2-3 周）**：可选的问题单和需求单管理

**关键成功因素：**
- ✅ 严格按照模块依赖关系进行开发
- ✅ 提前定义好模块间的 API 接口
- ✅ 持续集成和自动化测试
- ✅ 定期代码审查和技术分享
- ✅ 及时收集用户反馈并迭代优化

---

**文档结束**

本文档为 RemoteGPU 系统提供了完整的模块划分方案，包括每个模块的职责、功能、数据表设计、API 接口、依赖关系和前端页面配置。按照本文档的规划，可以系统化地完成整个平台的开发工作。
