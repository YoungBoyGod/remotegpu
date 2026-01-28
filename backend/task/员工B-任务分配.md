# 员工B - 任务分配

## 基本信息
- **负责人**: 员工B
- **主要模块**: Epic 1 - Feature 1.2 资源配额管理
- **预估工时**: 8小时
- **优先级**: P0（最高）
- **开始时间**: 第1周
- **依赖**: Workspace模型完成（员工A负责）

---

## 任务概述

你负责实现**资源配额管理**功能，这是控制用户资源使用的核心模块。包括：
1. ResourceQuota（资源配额）数据模型和DAO
2. 配额检查逻辑实现

---

## 详细任务列表

### Story 1.2.1: 创建ResourceQuota数据模型（4小时）

#### Task 1.2.1.1: 创建ResourceQuota实体模型（1小时）

**文件路径**: `internal/model/entity/resource_quota.go`

**Subtask清单**:
- [ ] 创建文件和结构体（20分钟）
  - 参考SQL文件`sql/03_users_and_permissions.sql`
  - 定义字段：ID, CustomerID, WorkspaceID, QuotaLevel, CPUQuota, MemoryQuota, GPUQuota, StorageQuota, EnvironmentQuota, CreatedAt, UpdatedAt

- [ ] 添加GORM标签（20分钟）
  - WorkspaceID使用指针类型`*uint`（允许NULL）
  - 数值字段添加默认值
  - 添加索引

- [ ] 实现TableName方法（5分钟）
  - 返回"resource_quotas"

- [ ] 添加注释和格式化（15分钟）

**验收标准**:
- [ ] WorkspaceID字段使用指针类型
- [ ] 字段类型与SQL表一致
- [ ] 可空字段正确配置

---

#### Task 1.2.1.2: 创建ResourceQuotaDao（3小时）

**文件路径**: `internal/dao/resource_quota.go`

**Subtask清单**:
- [ ] 创建基础结构（20分钟）
  - 定义ResourceQuotaDao结构体
  - 实现构造函数

- [ ] 实现Create方法（30分钟）
  - 处理WorkspaceID为NULL的情况

- [ ] 实现查询方法（1小时）
  - GetByCustomerID(customerID uint)
  - GetByWorkspaceID(workspaceID uint)
  - GetByQuotaLevel(level string)

- [ ] 实现Update方法（20分钟）

- [ ] 编写单元测试（50分钟）
  - 测试创建配额
  - 测试查询配额
  - 测试更新配额
  - 测试NULL值处理

**验收标准**:
- [ ] 所有CRUD方法实现正确
- [ ] NULL值处理正确
- [ ] 单元测试通过

---

### Story 1.2.2: 实现配额检查逻辑（4小时）

#### Task 1.2.2.1: 实现CheckQuota方法（2小时）

**文件路径**: `internal/dao/resource_quota.go`

**Subtask清单**:
- [ ] 定义方法签名和参数（15分钟）
  - `CheckQuota(customerID uint, cpu, memory, gpu int, storage int64) (bool, error)`

- [ ] 实现配额查询逻辑（30分钟）
  - 查询用户的资源配额
  - 查询用户当前已使用的资源

- [ ] 实现配额比较逻辑（30分钟）
  - 比较CPU配额
  - 比较内存配额
  - 比较GPU配额
  - 比较存储配额

- [ ] 编写单元测试（45分钟）
  - 测试配额充足的情况
  - 测试配额不足的情况
  - 测试边界情况

**验收标准**:
- [ ] 配额检查逻辑正确
- [ ] 所有资源类型都检查
- [ ] 单元测试覆盖所有场景

---

#### Task 1.2.2.2: 实现配额统计方法（2小时）

**Subtask清单**:
- [ ] 实现GetUsedResources方法（1小时）
  - 统计用户所有运行中环境的资源使用
  - 返回已用CPU、内存、GPU、存储

- [ ] 实现GetAvailableQuota方法（30分钟）
  - 计算可用配额 = 总配额 - 已用配额

- [ ] 编写单元测试（30分钟）
  - 测试资源统计准确性

**验收标准**:
- [ ] 统计逻辑正确
- [ ] 考虑所有环境状态
- [ ] 单元测试通过

---

## 协作和依赖

### 依赖其他员工的工作
- **员工A**: 必须等待Workspace模型完成后才能开始（WorkspaceID外键）
- **员工E**: GetUsedResources方法需要查询Environment表（可以先mock）

### 其他员工依赖你的工作
- **员工E**: 创建环境时需要调用CheckQuota方法检查配额

### 协作点
- 与**员工A**协作：等待Workspace模型完成的通知
- 与**员工E**协作：定义CheckQuota方法的接口，确保环境创建时能正确调用

---

## 开发顺序建议

1. **第一阶段**（可立即开始）：
   - 创建ResourceQuota实体模型
   - 实现基础CRUD方法
   - 编写单元测试

2. **第二阶段**（等待员工A完成Workspace）：
   - 测试WorkspaceID外键关系
   - 完善查询方法

3. **第三阶段**（可并行）：
   - 实现CheckQuota方法
   - 实现配额统计方法

---

## 验收检查清单

### 代码质量
- [ ] 所有代码通过`go fmt`格式化
- [ ] 所有代码通过`go vet`检查
- [ ] 所有导出的结构体和方法有注释
- [ ] 错误处理完善

### 测试
- [ ] 所有单元测试通过
- [ ] 测试覆盖率>80%
- [ ] 配额检查逻辑测试充分

### 功能
- [ ] 支持用户级配额
- [ ] 支持工作空间级配额
- [ ] 配额检查准确
- [ ] 资源统计正确

### Git提交
- [ ] 提交信息清晰：`feat: 实现ResourceQuota模型和配额检查逻辑`
- [ ] 代码已推送到远程仓库

---

## 注意事项

1. **指针类型**: WorkspaceID必须使用`*uint`类型，允许NULL
2. **配额级别**: QuotaLevel字段值为：free/basic/pro/enterprise
3. **资源单位**:
   - CPU: 核心数（int）
   - Memory: MB（int）
   - GPU: 数量（int）
   - Storage: 字节（int64）
4. **并发安全**: CheckQuota方法可能被并发调用，需要考虑事务
5. **性能优化**: GetUsedResources方法需要优化查询性能

---

## 技术要点

### 配额检查示例
```go
func (d *ResourceQuotaDao) CheckQuota(customerID uint, cpu, memory, gpu int, storage int64) (bool, error) {
    // 1. 查询配额
    quota, err := d.GetByCustomerID(customerID)
    if err != nil {
        return false, err
    }

    // 2. 查询已用资源
    used, err := d.GetUsedResources(customerID)
    if err != nil {
        return false, err
    }

    // 3. 检查是否超限
    if used.CPU + cpu > quota.CPUQuota {
        return false, nil
    }
    // ... 其他资源检查

    return true, nil
}
```

---

## 进度报告

请在完成每个Task后更新`团队协作文档.md`中的进度表。

---

## 联系方式

遇到问题请及时沟通：
- 依赖问题：联系员工A（Workspace）、员工E（Environment）
- 技术问题：查看`实施步骤指南.md`
- 其他问题：联系项目经理
