# RemoteGPU CMDB 设计文档

> 配置管理数据库（CMDB）- 硬件资产管理
>
> 创建日期：2026-01-26
>
> 前端技术栈：Vue 3 + Element Plus

---

## 目录

1. [CMDB 概述](#1-cmdb-概述)
2. [设备状态模型](#2-设备状态模型)
3. [核心表设计](#3-核心表设计)
4. [设备生命周期](#4-设备生命周期)
5. [与业务系统集成](#5-与业务系统集成)

---

## 1. CMDB 概述

### 1.1 什么是 CMDB

CMDB（Configuration Management Database）是一个独立的配置管理数据库，用于：

- **资产管理**：记录所有硬件设备的详细信息
- **状态跟踪**：跟踪设备的可用性和使用状态
- **生命周期管理**：从采购到报废的完整流程
- **变更管理**：记录设备的所有变更历史
- **关系管理**：设备之间的依赖和关联关系

### 1.2 为什么需要独立的 CMDB

```
传统方式（不推荐）：
业务系统直接管理设备 → 状态混乱、难以维护

CMDB 方式（推荐）：
┌─────────────────────────────────────────────────────────┐
│                      CMDB 系统                           │
│  - 设备注册                                              │
│  - 状态管理                                              │
│  - 生命周期管理                                          │
└────────────────────┬────────────────────────────────────┘
                     │ API 调用
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
   ┌────────┐  ┌────────┐  ┌────────┐
   │业务系统│  │监控系统│  │运维系统│
   │(调度器)│  │        │  │        │
   └────────┘  └────────┘  └────────┘
```

**优势：**
- ✅ 单一数据源（Single Source of Truth）
- ✅ 状态一致性保证
- ✅ 便于审计和追溯
- ✅ 支持多系统集成
- ✅ 独立的权限控制

---

## 2. 设备状态模型

### 2.1 设备状态定义

我们需要区分两个维度的状态：

#### A. 运营状态（Operational Status）- "是否可用"

| 状态 | 说明 | 可否分配 | 示例场景 |
|------|------|---------|---------|
| **available** | 可用 | ✅ 是 | 设备正常，可以分配给用户 |
| **maintenance** | 维护中 | ❌ 否 | 正在进行硬件维护 |
| **faulty** | 故障 | ❌ 否 | 设备出现硬件故障 |
| **retired** | 已退役 | ❌ 否 | 设备已报废或下线 |
| **reserved** | 预留 | ❌ 否 | 为特定用户/项目预留 |

#### B. 使用状态（Usage Status）- "是否在用"

| 状态 | 说明 | 资源占用 | 示例场景 |
|------|------|---------|---------|
| **idle** | 空闲 | 0% | 设备未分配任何环境 |
| **partial** | 部分使用 | 1-99% | 设备有部分资源被占用 |
| **full** | 完全使用 | 100% | 设备资源已全部分配 |
| **overcommit** | 超分配 | >100% | 允许超分配的情况 |

### 2.2 状态组合矩阵

| 运营状态 | 使用状态 | 是否可分配 | 说明 |
|---------|---------|-----------|------|
| available | idle | ✅ 是 | 最佳状态，可以立即分配 |
| available | partial | ✅ 是 | 有剩余资源，可以继续分配 |
| available | full | ❌ 否 | 资源已满，无法分配 |
| maintenance | * | ❌ 否 | 维护中，不可分配 |
| faulty | * | ❌ 否 | 故障，不可分配 |
| retired | * | ❌ 否 | 已退役，不可分配 |

### 2.3 状态转换图

```
┌──────────────────────────────────────────────────────────┐
│                    设备状态转换                            │
└──────────────────────────────────────────────────────────┘

    [新设备]
       │
       ▼
   ┌─────────┐
   │ pending │ (待上线)
   └────┬────┘
        │ 验收通过
        ▼
   ┌──────────┐
   │available │ (可用)
   └────┬─────┘
        │
        ├──────────────┐
        │              │
        ▼              ▼
   ┌──────────┐  ┌──────────┐
   │maintenance│  │ reserved │
   │  (维护)   │  │  (预留)  │
   └────┬─────┘  └────┬─────┘
        │              │
        └──────┬───────┘
               │ 恢复
               ▼
          ┌─────────┐
          │available│
          └────┬────┘
               │ 故障
               ▼
          ┌─────────┐
          │ faulty  │ (故障)
          └────┬────┘
               │
               ├─────────┐
               │ 修复    │ 无法修复
               ▼         ▼
          ┌─────────┐  ┌─────────┐
          │available│  │ retired │ (退役)
          └─────────┘  └─────────┘
```

---

## 3. 核心表设计

### 3.1 设备资产表 (cmdb_assets)

**表说明：** CMDB 核心表，存储所有硬件设备信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | VARCHAR(64) | PRIMARY KEY | 资产ID | "asset-abc123" |
| asset_number | VARCHAR(64) | UNIQUE NOT NULL | 资产编号 | "GPU-2024-001" |
| asset_type | VARCHAR(32) | NOT NULL | 资产类型 | "server", "gpu", "storage" |
| name | VARCHAR(128) | NOT NULL | 设备名称 | "GPU-Server-01" |
| brand | VARCHAR(64) | | 品牌 | "Dell", "HP", "Supermicro" |
| model | VARCHAR(128) | | 型号 | "PowerEdge R750xa" |
| serial_number | VARCHAR(128) | UNIQUE | 序列号 | "SN123456789" |
| purchase_date | DATE | | 采购日期 | 2024-01-15 |
| warranty_expire | DATE | | 保修到期日 | 2027-01-15 |
| purchase_price | DECIMAL(12,2) | | 采购价格 | 50000.00 |
| location | VARCHAR(256) | | 物理位置 | "机房A-机柜01-U10" |
| idc_name | VARCHAR(128) | | 机房名称 | "北京IDC" |
| rack_number | VARCHAR(64) | | 机柜编号 | "A-01" |
| operational_status | VARCHAR(20) | NOT NULL | 运营状态 | "available", "maintenance", "faulty" |
| usage_status | VARCHAR(20) | NOT NULL | 使用状态 | "idle", "partial", "full" |
| health_score | INT | DEFAULT 100 | 健康分数 | 95 (0-100) |
| owner | VARCHAR(128) | | 负责人 | "张三" |
| department | VARCHAR(128) | | 所属部门 | "运维部" |
| tags | TEXT[] | | 标签 | ["gpu-server", "production"] |
| metadata | JSONB | | 扩展信息 | {"vendor_contact": "xxx"} |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_asset_number` ON (asset_number)
- `idx_asset_type` ON (asset_type)
- `idx_operational_status` ON (operational_status)
- `idx_usage_status` ON (usage_status)
- `idx_location` ON (location)

**Vue 表单字段配置：**
```javascript
{
  asset_number: {
    label: '资产编号',
    type: 'input',
    required: true,
    placeholder: '请输入资产编号',
    rules: [{ pattern: /^[A-Z]+-\d{4}-\d{3}$/, message: '格式：GPU-2024-001' }]
  },
  asset_type: {
    label: '资产类型',
    type: 'select',
    required: true,
    options: [
      { label: '服务器', value: 'server' },
      { label: 'GPU设备', value: 'gpu' },
      { label: '存储设备', value: 'storage' },
      { label: '网络设备', value: 'network' }
    ]
  },
  operational_status: {
    label: '运营状态',
    type: 'select',
    required: true,
    options: [
      { label: '可用', value: 'available' },
      { label: '维护中', value: 'maintenance' },
      { label: '故障', value: 'faulty' },
      { label: '已退役', value: 'retired' },
      { label: '预留', value: 'reserved' }
    ],
    colorMap: {
      available: 'success',
      maintenance: 'warning',
      faulty: 'danger',
      retired: 'info',
      reserved: 'primary'
    }
  },
  usage_status: {
    label: '使用状态',
    type: 'tag',
    colorMap: {
      idle: 'info',
      partial: 'warning',
      full: 'danger',
      overcommit: 'danger'
    }
  },
  health_score: {
    label: '健康分数',
    type: 'progress',
    format: (value) => `${value}%`,
    colorMap: {
      90: 'success',
      70: 'warning',
      0: 'danger'
    }
  }
}
```

---

### 3.2 服务器详细信息表 (cmdb_servers)

**表说明：** 存储服务器的详细配置信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| asset_id | VARCHAR(64) | UNIQUE NOT NULL, FK | 资产ID | "asset-abc123" |
| hostname | VARCHAR(256) | UNIQUE | 主机名 | "gpu01.example.com" |
| ip_address | VARCHAR(64) | NOT NULL | 内网IP | "192.168.1.10" |
| os_type | VARCHAR(20) | NOT NULL | 操作系统类型 | "linux", "windows" |
| cpu_cores | INT | NOT NULL | CPU核心数 | 32 |
| memory_total | BIGINT | NOT NULL | 内存总量(字节) | 137438953472 (128GB) |
| gpu_count | INT | DEFAULT 0 | GPU数量 | 4 |
| deployment_mode | VARCHAR(20) | | 部署模式 | "traditional", "kubernetes" |
| last_heartbeat | TIMESTAMP | | 最后心跳 | 2026-01-26 10:00:00 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_asset_id` ON (asset_id)
- `idx_hostname` ON (hostname)
- `idx_ip_address` ON (ip_address)

**外键约束：**
- `asset_id` REFERENCES `cmdb_assets(id)` ON DELETE CASCADE

**Vue 表单字段配置：**
```javascript
{
  hostname: {
    label: '主机名',
    type: 'input',
    required: true
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
    options: [
      { label: 'Linux', value: 'linux' },
      { label: 'Windows', value: 'windows' }
    ]
  }
}
```

---

### 3.3 GPU设备详细信息表 (cmdb_gpus)

**表说明：** 存储GPU设备的详细信息

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | GPU ID | 1 |
| asset_id | VARCHAR(64) | FK | 所属资产ID | "asset-abc123" |
| server_id | BIGINT | FK | 所属服务器ID | 1 |
| gpu_index | INT | NOT NULL | GPU索引 | 0, 1, 2, 3 |
| uuid | VARCHAR(128) | UNIQUE | GPU UUID | "GPU-12345678-1234..." |
| name | VARCHAR(128) | | GPU型号 | "Tesla V100-SXM2-32GB" |
| memory_total | BIGINT | | 显存总量(字节) | 34359738368 (32GB) |
| operational_status | VARCHAR(20) | DEFAULT 'available' | 运营状态 | "available", "maintenance" |
| usage_status | VARCHAR(20) | DEFAULT 'idle' | 使用状态 | "idle", "allocated" |
| allocated_to | VARCHAR(64) | | 分配给的环境ID | "env-xyz789" |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 | 2026-01-26 09:00:00 |

**索引：**
- `idx_server_id` ON (server_id)
- `idx_operational_status` ON (operational_status)
- `idx_usage_status` ON (usage_status)

**外键约束：**
- `server_id` REFERENCES `cmdb_servers(id)` ON DELETE CASCADE

**Vue 表单字段配置：**
```javascript
{
  operational_status: {
    label: '运营状态',
    type: 'tag',
    colorMap: {
      available: 'success',
      maintenance: 'warning',
      faulty: 'danger'
    }
  },
  usage_status: {
    label: '使用状态',
    type: 'tag',
    colorMap: {
      idle: 'info',
      allocated: 'success'
    }
  }
}
```

---

### 3.4 设备变更历史表 (cmdb_change_logs)

**表说明：** 记录设备的所有变更历史

| 字段名 | 类型 | 约束 | 说明 | 示例值 |
|--------|------|------|------|--------|
| id | BIGSERIAL | PRIMARY KEY | 记录ID | 1 |
| asset_id | VARCHAR(64) | NOT NULL, FK | 资产ID | "asset-abc123" |
| change_type | VARCHAR(32) | NOT NULL | 变更类型 | "status_change", "config_change" |
| field_name | VARCHAR(128) | | 变更字段 | "operational_status" |
| old_value | TEXT | | 旧值 | "available" |
| new_value | TEXT | | 新值 | "maintenance" |
| reason | TEXT | | 变更原因 | "定期维护" |
| operator | VARCHAR(128) | | 操作人 | "admin" |
| created_at | TIMESTAMP | DEFAULT NOW() | 变更时间 | 2026-01-26 10:00:00 |

**索引：**
- `idx_asset_id` ON (asset_id)
- `idx_change_type` ON (change_type)
- `idx_created_at` ON (created_at DESC)

**外键约束：**
- `asset_id` REFERENCES `cmdb_assets(id)` ON DELETE CASCADE

---

## 4. 设备生命周期

### 4.1 生命周期阶段

```
1. 采购阶段 (Procurement)
   - 提交采购申请
   - 审批通过
   - 下单采购

2. 入库阶段 (Receiving)
   - 设备到货
   - 验收检查
   - 录入CMDB
   - 状态：pending

3. 部署阶段 (Deployment)
   - 安装操作系统
   - 配置网络
   - 安装Agent
   - 状态：available

4. 运营阶段 (Operation)
   - 正常使用
   - 定期维护
   - 状态：available/maintenance

5. 维修阶段 (Repair)
   - 故障报修
   - 维修处理
   - 状态：faulty

6. 退役阶段 (Retirement)
   - 性能不足/老化
   - 下线处理
   - 状态：retired
```

### 4.2 状态变更流程

```go
// cmdb/lifecycle.go
package cmdb

type LifecycleManager struct {
    db *gorm.DB
}

// 设备上线
func (m *LifecycleManager) OnlineAsset(assetID string, operator string) error {
    // 1. 检查设备状态
    var asset Asset
    if err := m.db.Where("id = ?", assetID).First(&asset).Error; err != nil {
        return err
    }

    if asset.OperationalStatus != "pending" {
        return errors.New("设备状态不是待上线")
    }

    // 2. 更新状态
    oldStatus := asset.OperationalStatus
    asset.OperationalStatus = "available"
    asset.UsageStatus = "idle"

    // 3. 保存变更
    if err := m.db.Save(&asset).Error; err != nil {
        return err
    }

    // 4. 记录变更历史
    m.logChange(assetID, "status_change", "operational_status",
        oldStatus, "available", "设备上线", operator)

    return nil
}

// 设备维护
func (m *LifecycleManager) MaintenanceAsset(assetID string, reason string, operator string) error {
    var asset Asset
    if err := m.db.Where("id = ?", assetID).First(&asset).Error; err != nil {
        return err
    }

    // 检查是否有正在运行的环境
    var envCount int64
    m.db.Model(&Environment{}).Where("host_id = ? AND status = 'running'", assetID).Count(&envCount)
    if envCount > 0 {
        return errors.New("设备上还有运行中的环境，无法进入维护")
    }

    oldStatus := asset.OperationalStatus
    asset.OperationalStatus = "maintenance"

    m.db.Save(&asset)
    m.logChange(assetID, "status_change", "operational_status",
        oldStatus, "maintenance", reason, operator)

    return nil
}

// 记录变更
func (m *LifecycleManager) logChange(assetID, changeType, fieldName, oldValue, newValue, reason, operator string) {
    log := &ChangeLog{
        AssetID:    assetID,
        ChangeType: changeType,
        FieldName:  fieldName,
        OldValue:   oldValue,
        NewValue:   newValue,
        Reason:     reason,
        Operator:   operator,
    }
    m.db.Create(log)
}
```

---

## 5. 与业务系统集成

### 5.1 集成架构

```
┌─────────────────────────────────────────────────────────┐
│                    CMDB 系统                             │
│  - 设备注册                                              │
│  - 状态管理                                              │
│  - 变更记录                                              │
└────────────────────┬────────────────────────────────────┘
                     │ REST API / gRPC
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
   ┌────────┐  ┌────────┐  ┌────────┐
   │调度系统│  │监控系统│  │计费系统│
   └────────┘  └────────┘  └────────┘
```

### 5.2 调度系统集成

**调度器查询可用设备：**

```go
// scheduler/cmdb_client.go
package scheduler

type CMDBClient struct {
    baseURL string
}

// 查询可用服务器
func (c *CMDBClient) QueryAvailableServers(req ResourceRequirement) ([]*Server, error) {
    // 调用 CMDB API
    resp, err := http.Get(fmt.Sprintf("%s/api/cmdb/servers/available?cpu=%d&memory=%d&gpu=%d",
        c.baseURL, req.CPU, req.Memory, req.GPU))

    if err != nil {
        return nil, err
    }

    var servers []*Server
    json.NewDecoder(resp.Body).Decode(&servers)

    return servers, nil
}

// 分配设备
func (c *CMDBClient) AllocateServer(serverID, envID string) error {
    // 调用 CMDB API 更新使用状态
    data := map[string]interface{}{
        "usage_status": "partial",
        "allocated_to": envID,
    }

    jsonData, _ := json.Marshal(data)
    resp, err := http.Post(
        fmt.Sprintf("%s/api/cmdb/servers/%s/allocate", c.baseURL, serverID),
        "application/json",
        bytes.NewBuffer(jsonData),
    )

    return err
}

// 释放设备
func (c *CMDBClient) ReleaseServer(serverID, envID string) error {
    resp, err := http.Post(
        fmt.Sprintf("%s/api/cmdb/servers/%s/release?env_id=%s", c.baseURL, serverID, envID),
        "application/json",
        nil,
    )

    return err
}
```

### 5.3 CMDB API 设计

```go
// cmdb/api.go
package cmdb

// 查询可用服务器
// GET /api/cmdb/servers/available
func QueryAvailableServers(c *gin.Context) {
    cpu := c.Query("cpu")
    memory := c.Query("memory")
    gpu := c.Query("gpu")

    var servers []Server
    query := db.Joins("JOIN cmdb_assets ON cmdb_servers.asset_id = cmdb_assets.id").
        Where("cmdb_assets.operational_status = ?", "available").
        Where("cmdb_assets.usage_status IN ?", []string{"idle", "partial"})

    if cpu != "" {
        query = query.Where("cmdb_servers.cpu_cores >= ?", cpu)
    }

    query.Find(&servers)

    c.JSON(200, gin.H{"servers": servers})
}

// 分配服务器
// POST /api/cmdb/servers/:id/allocate
func AllocateServer(c *gin.Context) {
    serverID := c.Param("id")

    var req struct {
        EnvID string `json:"env_id"`
        CPU   int    `json:"cpu"`
        Memory int64 `json:"memory"`
        GPU   int    `json:"gpu"`
    }

    c.ShouldBindJSON(&req)

    // 更新使用状态
    var server Server
    db.Where("id = ?", serverID).First(&server)

    // 计算使用率
    usagePercent := calculateUsage(server, req.CPU, req.Memory, req.GPU)

    var newStatus string
    if usagePercent >= 100 {
        newStatus = "full"
    } else if usagePercent > 0 {
        newStatus = "partial"
    } else {
        newStatus = "idle"
    }

    // 更新资产状态
    db.Model(&Asset{}).Where("id = ?", server.AssetID).Update("usage_status", newStatus)

    c.JSON(200, gin.H{"status": "success"})
}
```

---

## 6. 总结

### 6.1 CMDB 的核心价值

1. **单一数据源**：所有设备信息统一管理
2. **状态一致性**：避免业务系统状态不一致
3. **审计追溯**：完整的变更历史记录
4. **生命周期管理**：从采购到退役的全流程
5. **多系统集成**：为调度、监控、计费等系统提供统一接口

### 6.2 实施建议

**阶段 1：基础 CMDB（1-2周）**
- 实现资产表和服务器表
- 实现基本的状态管理
- 提供查询 API

**阶段 2：完善功能（2-3周）**
- 添加 GPU 详细信息表
- 实现变更历史记录
- 实现生命周期管理

**阶段 3：系统集成（1-2周）**
- 调度系统集成
- 监控系统集成
- 前端界面开发

### 6.3 与原 hosts 表的关系

```
原设计：
hosts 表 → 直接被业务系统使用

新设计：
cmdb_assets (资产表) → cmdb_servers (服务器表) → 业务系统通过 API 查询

优势：
✅ 职责分离：CMDB 管理设备，业务系统管理业务
✅ 状态清晰：运营状态 + 使用状态
✅ 易于扩展：可以管理更多类型的设备
✅ 审计完善：所有变更都有记录
```

---

**文档结束**

本文档提供了完整的 CMDB 设计方案，包括设备状态模型、核心表结构、生命周期管理和系统集成方案。
