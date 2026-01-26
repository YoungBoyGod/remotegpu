# GPU 设备管理

> 所属模块：模块 1 - CMDB 设备管理模块
>
> 功能编号：1.2
>
> 优先级：P0（必须）

---

## 1. 功能概述

### 1.1 功能描述

GPU 设备管理功能负责自动发现、注册和管理服务器上的 GPU 设备，采集 GPU 的详细信息（型号、显存、UUID、计算能力等），并提供 GPU 的分配、释放、健康检查等核心能力，为资源调度模块提供准确的 GPU 资源信息。

### 1.2 业务价值

- ✅ 自动发现 GPU 设备，无需手动配置
- ✅ 精确管理 GPU 资源分配状态
- ✅ 实时监控 GPU 健康状态
- ✅ 支持多种 GPU 型号（NVIDIA、AMD）
- ✅ 支持 GPU 拓扑关系管理（NVLink、PCIe）

### 1.3 适用场景

- 服务器上线后自动发现 GPU
- GPU 资源分配给开发环境
- GPU 故障检测和告警
- GPU 使用情况统计

---

## 2. 功能需求

### 2.1 GPU 自动发现

**需求描述：**
在主机注册完成后，自动检测并注册主机上的所有 GPU 设备。

**实现方式：**

**Linux 平台：**
```bash
# 使用 nvidia-smi 查询 GPU 信息
nvidia-smi --query-gpu=index,uuid,name,memory.total,compute_cap,pci.bus_id \
    --format=csv,noheader,nounits

# 输出示例：
# 0, GPU-12345678-1234-1234-1234-123456789012, Tesla V100-SXM2-32GB, 32510, 7.0, 00000000:3B:00.0
# 1, GPU-87654321-4321-4321-4321-210987654321, Tesla V100-SXM2-32GB, 32510, 7.0, 00000000:3C:00.0
```

**Windows 平台：**
```powershell
# 使用 nvidia-smi.exe
nvidia-smi.exe --query-gpu=index,uuid,name,memory.total,compute_cap,pci.bus_id `
    --format=csv,noheader,nounits
```

**接口定义：**
```go
// 内部函数：发现 GPU
func DiscoverGPUs(hostID string, sshClient *ssh.Client) ([]GPUInfo, error) {
    cmd := "nvidia-smi --query-gpu=index,uuid,name,memory.total,compute_cap,pci.bus_id --format=csv,noheader,nounits"

    output, err := sshClient.Run(cmd)
    if err != nil {
        return nil, fmt.Errorf("nvidia-smi 执行失败: %v", err)
    }

    var gpus []GPUInfo
    lines := strings.Split(strings.TrimSpace(output), "\n")

    for _, line := range lines {
        fields := strings.Split(line, ",")
        if len(fields) < 6 {
            continue
        }

        gpu := GPUInfo{
            Index:             parseInt(strings.TrimSpace(fields[0])),
            UUID:              strings.TrimSpace(fields[1]),
            Name:              strings.TrimSpace(fields[2]),
            MemoryTotal:       parseInt64(strings.TrimSpace(fields[3])) * 1024 * 1024, // MiB to Bytes
            ComputeCapability: strings.TrimSpace(fields[4]),
            PCIBusID:          strings.TrimSpace(fields[5]),
        }
        gpus = append(gpus, gpu)
    }

    return gpus, nil
}
```

**验收标准：**
- [ ] 自动发现所有 GPU 设备
- [ ] 发现成功率 > 99%
- [ ] 发现耗时 < 5 秒
- [ ] 支持 NVIDIA GPU（优先）
- [ ] 支持 AMD GPU（可选）

### 2.2 GPU 信息采集

**需求描述：**
采集 GPU 的详细信息，包括型号、显存、架构、计算能力等。

**采集项清单：**

| 采集项 | nvidia-smi 参数 | 说明 |
|--------|----------------|------|
| GPU 索引 | index | 主机上的 GPU 编号（0, 1, 2...） |
| GPU UUID | uuid | 全局唯一标识符 |
| GPU 名称 | name | 型号名称（Tesla V100, RTX 4090） |
| 显存总量 | memory.total | 显存容量（MiB） |
| 计算能力 | compute_cap | CUDA 计算能力（7.0, 8.0, 8.6） |
| PCI 总线 ID | pci.bus_id | PCIe 总线地址 |
| 驱动版本 | driver_version | NVIDIA 驱动版本 |
| CUDA 版本 | cuda_version | CUDA 运行时版本 |
| 功耗限制 | power.limit | 最大功耗（瓦） |
| 温度限制 | temperature.gpu.tlimit | 温度阈值（摄氏度） |

**数据结构：**
```go
type GPUInfo struct {
    Index             int
    UUID              string
    Name              string
    Brand             string  // NVIDIA, AMD
    Architecture      string  // Ampere, Turing, Volta
    MemoryTotal       int64   // 字节
    CUDACores         int
    ComputeCapability string
    PCIBusID          string
    DriverVersion     string
    CUDAVersion       string
    PowerLimit        int     // 瓦
    TemperatureLimit  int     // 摄氏度
}
```

**验收标准：**
- [ ] 所有必需字段采集成功
- [ ] 显存容量精确到字节
- [ ] 支持识别主流 GPU 型号

### 2.3 GPU 分配与释放

**需求描述：**
为开发环境分配 GPU 资源，并在环境销毁时释放 GPU。

**分配流程：**
1. 调度模块请求分配 GPU
2. 查询可用 GPU（status = 'available'）
3. 根据策略选择 GPU（最少负载、轮询等）
4. 更新 GPU 状态为 'allocated'
5. 记录分配信息（allocated_to, allocated_at）
6. 返回 GPU 信息

**释放流程：**
1. 环境销毁时触发释放
2. 更新 GPU 状态为 'available'
3. 清除分配信息
4. 通知调度模块

**接口定义：**
```go
// 分配 GPU
POST /api/cmdb/gpus/allocate
Body: {
  "host_id": "host-abc123",
  "env_id": "env-xyz789",
  "gpu_count": 2,
  "gpu_model": "Tesla V100",  // 可选：指定型号
  "strategy": "least_loaded"   // 分配策略
}
Response: {
  "gpus": [
    {
      "gpu_id": 1,
      "gpu_index": 0,
      "uuid": "GPU-12345678...",
      "name": "Tesla V100-SXM2-32GB",
      "memory_total": 34359738368
    },
    {
      "gpu_id": 2,
      "gpu_index": 1,
      "uuid": "GPU-87654321...",
      "name": "Tesla V100-SXM2-32GB",
      "memory_total": 34359738368
    }
  ]
}

// 释放 GPU
POST /api/cmdb/gpus/release
Body: {
  "env_id": "env-xyz789"
}
Response: {
  "released_count": 2,
  "status": "success"
}
```

**验收标准：**
- [ ] 分配响应时间 < 1 秒
- [ ] 支持并发分配（使用分布式锁）
- [ ] 分配失败时自动回滚
- [ ] 释放后状态正确更新

### 2.4 GPU 健康检查

**需求描述：**
定期检查 GPU 的健康状态，包括温度、功耗、ECC 错误等。

**检查项清单：**

| 检查项 | 正常范围 | 告警阈值 | 检查频率 |
|--------|---------|---------|---------|
| GPU 温度 | < 80°C | > 85°C | 每 30 秒 |
| GPU 功耗 | < 额定功率 | > 额定功率 * 1.1 | 每 30 秒 |
| ECC 错误 | 0 | > 0 | 每 5 分钟 |
| GPU 可访问性 | 正常 | nvidia-smi 无响应 | 每 30 秒 |
| 显存 ECC | 无错误 | 有错误 | 每 5 分钟 |

**健康检查逻辑：**
```go
func CheckGPUHealth(gpuID int64) (*HealthCheckResult, error) {
    var gpu GPU
    db.Where("id = ?", gpuID).First(&gpu)

    result := &HealthCheckResult{
        GPUID:     gpuID,
        Status:    "healthy",
        Issues:    []string{},
        CheckedAt: time.Now(),
    }

    // 查询最新监控数据
    var metric GPUMetric
    db.Where("gpu_id = ?", gpuID).
        Order("collected_at DESC").
        First(&metric)

    // 检查温度
    if metric.Temperature > 85 {
        result.Issues = append(result.Issues, "温度过高")
        result.Status = "unhealthy"
    } else if metric.Temperature > 80 {
        result.Issues = append(result.Issues, "温度偏高")
        result.Status = "degraded"
    }

    // 检查功耗
    if metric.PowerDraw > float64(gpu.PowerLimit) * 1.1 {
        result.Issues = append(result.Issues, "功耗异常")
        result.Status = "degraded"
    }

    // 检查可访问性
    if time.Since(metric.CollectedAt) > 5*time.Minute {
        result.Issues = append(result.Issues, "GPU 无响应")
        result.Status = "unhealthy"
    }

    // 更新 GPU 健康状态
    db.Model(&gpu).Update("health_status", result.Status)

    return result, nil
}
```

**验收标准：**
- [ ] 健康检查覆盖所有关键指标
- [ ] 异常检测准确率 > 95%
- [ ] 检查耗时 < 1 秒
- [ ] 异常时自动告警

### 2.5 GPU 拓扑关系管理

**需求描述：**
管理 GPU 之间的拓扑关系，包括 NVLink 连接和 PCIe 拓扑。

**拓扑信息采集：**
```bash
# 查询 GPU 拓扑
nvidia-smi topo -m

# 输出示例：
#         GPU0    GPU1    GPU2    GPU3    CPU Affinity    NUMA Affinity
# GPU0     X      NV12    NV12    NV12    0-19,40-59      0
# GPU1    NV12     X      NV12    NV12    0-19,40-59      0
# GPU2    NV12    NV12     X      NV12    20-39,60-79     1
# GPU3    NV12    NV12    NV12     X      20-39,60-79     1
```

**数据结构：**
```go
type GPUTopology struct {
    GPUID        int64
    ConnectedGPUs []GPUConnection
    CPUAffinity   string
    NUMAAffinity  int
}

type GPUConnection struct {
    TargetGPUID int64
    LinkType    string  // NVLink, PCIe
    LinkSpeed   string  // NV12, NV6, PHB
}
```

**验收标准：**
- [ ] 正确识别 NVLink 连接
- [ ] 正确识别 PCIe 拓扑
- [ ] 支持 NUMA 亲和性信息

---

## 3. 数据模型

### 3.1 GPU 设备表 (cmdb_gpus)

```sql
CREATE TABLE cmdb_gpus (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(64) NOT NULL,
    gpu_index INT NOT NULL,

    -- GPU 信息
    uuid VARCHAR(128) UNIQUE,
    name VARCHAR(128),
    brand VARCHAR(64),                      -- NVIDIA, AMD
    architecture VARCHAR(64),               -- Ampere, Turing, Volta

    -- 规格
    memory_total BIGINT,                    -- 显存总量（字节）
    cuda_cores INT,
    compute_capability VARCHAR(32),         -- 7.5, 8.0, 8.6
    pci_bus_id VARCHAR(32),

    -- 状态
    status VARCHAR(20) DEFAULT 'available', -- available, allocated, maintenance, error
    health_status VARCHAR(20) DEFAULT 'healthy',

    -- 分配信息
    allocated_to VARCHAR(64),               -- 环境 ID
    allocated_at TIMESTAMP,

    -- 性能信息
    power_limit INT,                        -- 功耗限制（瓦）
    temperature_limit INT,                  -- 温度限制（摄氏度）

    -- 驱动信息
    driver_version VARCHAR(32),
    cuda_version VARCHAR(32),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (host_id) REFERENCES cmdb_assets(id) ON DELETE CASCADE,
    UNIQUE(host_id, gpu_index),
    INDEX idx_status (status),
    INDEX idx_allocated_to (allocated_to),
    INDEX idx_health_status (health_status)
);
```

### 3.2 GPU 拓扑表 (cmdb_gpu_topology)

```sql
CREATE TABLE cmdb_gpu_topology (
    id BIGSERIAL PRIMARY KEY,
    gpu_id BIGINT NOT NULL,
    target_gpu_id BIGINT NOT NULL,
    link_type VARCHAR(20),                  -- nvlink, pcie
    link_speed VARCHAR(20),                 -- NV12, NV6, PHB
    cpu_affinity VARCHAR(64),
    numa_affinity INT,
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (gpu_id) REFERENCES cmdb_gpus(id) ON DELETE CASCADE,
    FOREIGN KEY (target_gpu_id) REFERENCES cmdb_gpus(id) ON DELETE CASCADE,
    UNIQUE(gpu_id, target_gpu_id)
);
```

---

## 4. API 接口

### 4.1 查询 GPU 列表

```go
// GET /api/cmdb/gpus
func ListGPUs(c *gin.Context) {
    hostID := c.Query("host_id")
    status := c.Query("status")
    gpuModel := c.Query("gpu_model")

    var gpus []GPU
    query := db.Model(&GPU{})

    if hostID != "" {
        query = query.Where("host_id = ?", hostID)
    }
    if status != "" {
        query = query.Where("status = ?", status)
    }
    if gpuModel != "" {
        query = query.Where("name LIKE ?", "%"+gpuModel+"%")
    }

    query.Find(&gpus)

    c.JSON(200, gin.H{"gpus": gpus})
}
```

### 4.2 分配 GPU

```go
// POST /api/cmdb/gpus/allocate
func AllocateGPU(c *gin.Context) {
    var req struct {
        HostID    string `json:"host_id"`
        EnvID     string `json:"env_id" binding:"required"`
        GPUCount  int    `json:"gpu_count" binding:"required"`
        GPUModel  string `json:"gpu_model"`
        Strategy  string `json:"strategy"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 查询可用 GPU
    var availableGPUs []GPU
    query := db.Where("status = ?", "available")

    if req.HostID != "" {
        query = query.Where("host_id = ?", req.HostID)
    }
    if req.GPUModel != "" {
        query = query.Where("name LIKE ?", "%"+req.GPUModel+"%")
    }

    query.Limit(req.GPUCount).Find(&availableGPUs)

    if len(availableGPUs) < req.GPUCount {
        c.JSON(503, gin.H{"error": "可用 GPU 不足"})
        return
    }

    // 分配 GPU
    var allocatedGPUs []GPU
    for i := 0; i < req.GPUCount; i++ {
        gpu := &availableGPUs[i]

        // 使用事务更新状态
        err := db.Transaction(func(tx *gorm.DB) error {
            return tx.Model(gpu).Updates(map[string]interface{}{
                "status":       "allocated",
                "allocated_to": req.EnvID,
                "allocated_at": time.Now(),
            }).Error
        })

        if err != nil {
            // 回滚已分配的 GPU
            rollbackGPUs(allocatedGPUs)
            c.JSON(500, gin.H{"error": "分配失败"})
            return
        }

        allocatedGPUs = append(allocatedGPUs, *gpu)
    }

    c.JSON(200, gin.H{"gpus": allocatedGPUs})
}
```

---

## 5. 前端界面

### 5.1 GPU 列表页面

**页面路径：** `/cmdb/gpus`

**功能要求：**
- 显示所有 GPU 列表
- 支持按主机、状态、型号筛选
- 显示 GPU 使用率、温度、功耗
- 支持批量操作（维护、退役）

### 5.2 GPU 详情页面

**页面路径：** `/cmdb/gpus/:id`

**功能要求：**
- 显示 GPU 基本信息
- 显示实时监控数据（温度、功耗、使用率）
- 显示分配历史
- 显示拓扑关系图

---

## 6. 测试用例

| 用例编号 | 测试场景 | 预期结果 |
|---------|---------|---------|
| TC-01 | 主机注册后自动发现 GPU | 所有 GPU 正确注册 |
| TC-02 | 分配可用 GPU | 状态更新为 allocated |
| TC-03 | 分配不存在的 GPU | 返回错误 |
| TC-04 | 释放已分配的 GPU | 状态恢复为 available |
| TC-05 | GPU 温度超过阈值 | 健康状态变为 degraded |
| TC-06 | GPU 无响应 | 健康状态变为 unhealthy |

---

## 7. 实施计划

**总工作量：** 8 天

| 任务 | 工作量 | 依赖 |
|------|--------|------|
| 数据库表设计 | 0.5 天 | - |
| GPU 发现功能 | 2 天 | 主机注册 |
| GPU 分配/释放 | 2 天 | GPU 发现 |
| 健康检查功能 | 2 天 | GPU 发现 |
| 前端页面 | 2 天 | 后端 API |
| 测试 | 1.5 天 | 所有功能 |

---

**文档版本：** v1.0
**创建日期：** 2026-01-26
