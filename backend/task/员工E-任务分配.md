# 员工E - 任务分配

## 基本信息
- **负责人**: 员工E
- **主要模块**: Epic 3 - Feature 3.1 环境数据模型
- **预估工时**: 12小时
- **优先级**: P0（最高）
- **开始时间**: 第1周
- **依赖**: Workspace模型（员工A）、Host模型（员工C）

---

## 任务概述

你负责实现**环境数据模型**，这是整个系统的核心业务模型。包括：
1. Environment（开发环境）数据模型和DAO
2. PortMapping（端口映射）数据模型和DAO

---

## 详细任务列表

### Story 3.1.1: 创建Environment实体模型（4小时）

#### Task 3.1.1.1: 创建基础模型（2小时）

**文件路径**: `internal/model/entity/environment.go`

**Subtask清单**:
- [ ] 定义Environment结构体（1小时）
  - 参考SQL文件`sql/05_environments.sql`
  - 基本字段：ID (string), CustomerID, WorkspaceID, HostID
  - 配置字段：Name, Description, Image, CPU, Memory, GPU, Storage
  - 状态字段：Status, ErrorMessage
  - 容器字段：ContainerID, PodName, Namespace
  - 端口字段：SSHPort, RDPPort, JupyterPort
  - 时间字段：StartedAt, StoppedAt, CreatedAt, UpdatedAt

- [ ] 添加GORM标签（30分钟）
  - ID使用varchar(64)
  - 外键字段配置
  - 可空字段使用指针类型

- [ ] 添加关联关系（30分钟）
  - 关联Customer
  - 关联Workspace
  - 关联Host

**验收标准**:
- [ ] 字段完整
- [ ] 关联关系正确
- [ ] GORM标签配置正确

---

#### Task 3.1.1.2: 创建EnvironmentDao（2小时）

**文件路径**: `internal/dao/environment.go`

**Subtask清单**:
- [ ] 实现基础CRUD（1小时）
  - Create, GetByID, Update, Delete

- [ ] 编写单元测试（1小时）

**验收标准**:
- [ ] CRUD方法正确
- [ ] 单元测试通过

---

### Story 3.1.2: 创建PortMapping模型（4小时）

#### Task 3.1.2.1: 创建PortMapping实体（2小时）

**文件路径**: `internal/model/entity/port_mapping.go`

**Subtask清单**:
- [ ] 定义结构体（1小时）
  - 字段：ID, EnvID, ServiceType, ExternalPort, InternalPort, Protocol

- [ ] 添加唯一约束（1小时）
  - ExternalPort必须唯一

**验收标准**:
- [ ] 唯一约束配置正确

---

#### Task 3.1.2.2: 创建PortMappingDao（2小时）

**文件路径**: `internal/dao/port_mapping.go`

**Subtask清单**:
- [ ] 实现CRUD方法（1小时）
- [ ] 实现端口分配逻辑（1小时）
  - AllocatePort方法
  - 查找可用端口

**验收标准**:
- [ ] 端口分配逻辑正确
- [ ] 避免端口冲突

---

### Story 3.1.3: 实现查询功能（4小时）

#### Task 3.1.3.1: 实现用户查询（2小时）

**Subtask清单**:
- [ ] GetByCustomerID方法（1小时）
- [ ] GetByWorkspaceID方法（1小时）

---

#### Task 3.1.3.2: 实现状态查询（2小时）

**Subtask清单**:
- [ ] GetByStatus方法（1小时）
- [ ] GetByHostID方法（1小时）

---

## 协作和依赖

### 依赖其他员工的工作
- **员工A**: 必须等待Workspace模型完成（WorkspaceID外键）
- **员工C**: 必须等待Host模型完成（HostID外键）
- **员工D**: 环境创建时需要调用GPU分配方法

### 其他员工依赖你的工作
- **员工B**: ResourceQuota的GetUsedResources方法需要查询Environment表

### 协作点
- 与**员工A**协作：等待Workspace模型完成
- 与**员工C**协作：等待Host模型完成
- 与**员工D**协作：定义GPU分配接口

---

## 验收检查清单

### 代码质量
- [ ] 代码格式化和检查通过
- [ ] 所有方法有注释

### 测试
- [ ] 单元测试通过
- [ ] 测试覆盖率>80%

### 功能
- [ ] 环境CRUD正确
- [ ] 端口映射正确
- [ ] 查询功能完善

### Git提交
- [ ] 提交信息：`feat: 实现Environment和PortMapping模型`

---

## 注意事项

1. **ID类型**: Environment的ID是string类型
2. **端口分配**: 需要避免端口冲突
3. **状态管理**: Status字段值：creating/running/stopped/error/deleting
4. **外键关系**: 依赖多个其他表

---

## 进度报告

请在完成每个Task后更新`团队协作文档.md`。

---

## 联系方式

遇到问题请联系：
- 员工A（Workspace依赖）
- 员工C（Host依赖）
- 员工D（GPU协作）
- 项目经理
