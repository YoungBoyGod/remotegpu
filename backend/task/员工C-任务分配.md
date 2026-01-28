# 员工C - 任务分配

## 基本信息
- **负责人**: 员工C
- **主要模块**: Epic 2 - Feature 2.1 主机管理
- **预估工时**: 16小时
- **优先级**: P0（最高）
- **开始时间**: 第1周
- **依赖**: 无（可立即开始）

---

## 任务概述

你负责实现**主机管理**功能，这是资源调度的基础模块。包括：
1. Host（主机）数据模型和DAO
2. 主机资源管理方法
3. 主机状态和心跳管理

---

## 详细任务列表

### Story 2.1.1: 创建Host数据模型（5小时）

#### Task 2.1.1.1: 创建Host实体模型（1.5小时）

**文件路径**: `internal/model/entity/host.go`

**Subtask清单**:
- [ ] 创建文件和基础结构（20分钟）
  - 参考SQL文件`sql/04_hosts_and_devices.sql`

- [ ] 定义所有字段（30分钟）
  - ID (string) - 主键，非自增
  - Name, Hostname, IPAddress, PublicIP
  - OSType, OSVersion, Arch
  - DeploymentMode, K8sNodeName
  - Status, HealthStatus
  - 资源字段：TotalCPU, TotalMemory, TotalDisk, TotalGPU
  - 已用资源：UsedCPU, UsedMemory, UsedDisk, UsedGPU
  - 端口：SSHPort, WinRMPort, AgentPort
  - Labels (JSONB), Tags (数组)
  - LastHeartbeat, RegisteredAt, CreatedAt, UpdatedAt

- [ ] 添加GORM标签（30分钟）
  - ID使用`gorm:"primaryKey;type:varchar(64)"`
  - Labels使用`gorm:"type:jsonb"`
  - Tags使用`gorm:"type:text[]"`
  - 可空字段使用指针类型

- [ ] 添加注释和验证（10分钟）

**验收标准**:
- [ ] ID使用string类型
- [ ] JSONB和数组类型正确配置
- [ ] 所有字段与SQL表一致

**注意事项**:
- 需要导入`gorm.io/datatypes`包（JSONB）
- 需要导入`github.com/lib/pq`包（StringArray）

---

#### Task 2.1.1.2: 创建HostDao基础方法（3.5小时）

**文件路径**: `internal/dao/host.go`

**Subtask清单**:
- [ ] 创建基础结构和构造函数（20分钟）

- [ ] 实现Create方法（30分钟）
  - 处理ID字段（非自增）

- [ ] 实现GetByID方法（20分钟）

- [ ] 实现Update方法（20分钟）

- [ ] 实现Delete方法（20分钟）

- [ ] 编写基础CRUD测试（1小时）

**验收标准**:
- [ ] 所有CRUD方法实现正确
- [ ] ID字段处理正确
- [ ] 单元测试通过

---

### Story 2.1.2: 实现主机查询功能（4小时）

#### Task 2.1.2.1: 实现状态查询方法（2小时）

**Subtask清单**:
- [ ] 实现GetByStatus方法（30分钟）
  - 按状态筛选：online/offline/maintenance

- [ ] 实现GetByHealthStatus方法（30分钟）
  - 按健康状态筛选：healthy/degraded/unhealthy/unknown

- [ ] 实现GetByDeploymentMode方法（30分钟）
  - 按部署模式筛选：traditional/kubernetes

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 所有查询方法实现正确
- [ ] 支持多条件筛选
- [ ] 单元测试通过

---

#### Task 2.1.2.2: 实现资源查询方法（2小时）

**Subtask清单**:
- [ ] 实现GetAvailableHosts方法（1小时）
  - 查询条件：status='online' AND health_status='healthy'
  - 资源条件：(total_cpu - used_cpu) >= cpu
  - 资源条件：(total_memory - used_memory) >= memory
  - 资源条件：(total_gpu - used_gpu) >= gpu
  - 按负载排序：ORDER BY (used_cpu::float / total_cpu) ASC

- [ ] 实现资源排序逻辑（30分钟）
  - 负载均衡算法

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 资源查询逻辑正确
- [ ] 负载均衡算法有效
- [ ] 单元测试通过

---

### Story 2.1.3: 实现主机资源管理（4小时）

#### Task 2.1.3.1: 实现资源更新方法（2小时）

**Subtask清单**:
- [ ] 实现UpdateUsedResources方法（1小时）
  - 更新已用CPU、内存、GPU
  - 支持增量更新

- [ ] 实现资源增减逻辑（30分钟）
  - 分配资源时增加used_*
  - 释放资源时减少used_*

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 资源更新逻辑正确
- [ ] 支持并发更新
- [ ] 单元测试通过

---

#### Task 2.1.3.2: 实现心跳管理（2小时）

**Subtask清单**:
- [ ] 实现UpdateHeartbeat方法（30分钟）
  - 更新last_heartbeat字段

- [ ] 实现UpdateStatus方法（30分钟）
  - 更新status和health_status

- [ ] 编写测试用例（1小时）

**验收标准**:
- [ ] 心跳更新正确
- [ ] 状态转换逻辑正确
- [ ] 单元测试通过

---

### Story 2.1.4: 集成测试和文档（3小时）

#### Task 2.1.4.1: 编写集成测试（2小时）

**Subtask清单**:
- [ ] 测试主机注册流程（30分钟）
- [ ] 测试资源分配流程（30分钟）
- [ ] 测试心跳更新流程（30分钟）
- [ ] 测试异常场景（30分钟）

---

#### Task 2.1.4.2: 编写API文档（1小时）

**Subtask清单**:
- [ ] 编写方法说明（30分钟）
- [ ] 编写使用示例（30分钟）

---

## 协作和依赖

### 依赖其他员工的工作
- **无直接依赖** - 可以独立开发

### 其他员工依赖你的工作
- **员工D**: GPU模型需要关联Host（host_id外键）
- **员工E**: Environment模型需要关联Host（host_id外键）

### 协作点
- 与**员工D**协作：完成Host模型后立即通知
- 与**员工E**协作：确保host_id字段类型一致（string）

---

## 验收检查清单

### 代码质量
- [ ] 所有代码通过`go fmt`和`go vet`
- [ ] 所有导出方法有注释
- [ ] 错误处理完善

### 测试
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 测试覆盖率>80%

### 功能
- [ ] 支持主机注册
- [ ] 支持资源查询
- [ ] 支持心跳管理
- [ ] 支持状态管理

### Git提交
- [ ] 提交信息：`feat: 实现Host模型和主机管理功能`

---

## 注意事项

1. **ID类型**: Host的ID是string类型，不是自增主键
2. **JSONB字段**: Labels字段需要GIN索引支持高效查询
3. **数组字段**: Tags字段使用PostgreSQL数组类型
4. **资源计算**: 注意浮点数精度问题
5. **并发安全**: 资源更新需要考虑并发场景

---

## 进度报告

请在完成每个Task后更新`团队协作文档.md`。

---

## 联系方式

遇到问题请联系：
- 员工D（GPU依赖）
- 员工E（Environment依赖）
- 项目经理
