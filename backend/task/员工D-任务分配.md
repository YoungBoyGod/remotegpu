# 员工D - 任务分配

## 基本信息
- **负责人**: 员工D
- **主要模块**: Epic 2 - Feature 2.2 GPU管理
- **预估工时**: 20小时
- **优先级**: P0（最高）
- **开始时间**: 第1周
- **依赖**: Host模型完成（员工C负责）

---

## 任务概述

你负责实现**GPU管理**功能，这是资源分配的核心模块。包括：
1. GPU数据模型和DAO
2. GPU查询和分配功能
3. 并发安全的GPU分配逻辑

---

## 详细任务列表

### Story 2.2.1: 创建GPU数据模型（4小时）

#### Task 2.2.1.1: 创建GPU实体模型（1小时）

**文件路径**: `internal/model/entity/gpu.go`

**Subtask清单**:
- [ ] 创建文件和结构体（20分钟）
  - 参考SQL文件`sql/04_hosts_and_devices.sql`
  - 定义字段：ID, HostID, GPUIndex, UUID, Name, Brand, Architecture
  - MemoryTotal, CUDACores, ComputeCapability
  - Status, HealthStatus, AllocatedTo, AllocatedAt
  - PowerLimit, TemperatureLimit, CreatedAt, UpdatedAt

- [ ] 添加GORM标签和索引（15分钟）
  - UUID字段添加唯一索引
  - HostID添加外键索引
  - AllocatedTo使用指针类型

- [ ] 添加注释（5分钟）

**验收标准**:
- [ ] 字段完整
- [ ] 索引配置正确
- [ ] UUID唯一约束生效

---

#### Task 2.2.1.2: 创建GPUDao基础方法（3小时）

**文件路径**: `internal/dao/gpu.go`

**Subtask清单**:
- [ ] 创建基础结构（20分钟）
- [ ] 实现CRUD方法（1小时）
  - Create, GetByID, Update, Delete
- [ ] 编写单元测试（1小时40分钟）

**验收标准**:
- [ ] CRUD方法正确
- [ ] 单元测试通过

---

### Story 2.2.2: 实现GPU查询功能（4小时）

#### Task 2.2.2.1: 实现基础查询方法（2小时）

**Subtask清单**:
- [ ] 实现GetByHostID方法（30分钟）
  - 查询指定主机的所有GPU

- [ ] 实现GetByStatus方法（30分钟）
  - 按状态查询：available/allocated/maintenance

- [ ] 实现GetByAllocatedTo方法（30分钟）
  - 查询分配给特定环境的GPU

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 查询方法正确
- [ ] 单元测试通过

---

#### Task 2.2.2.2: 实现可用GPU查询（2小时）

**Subtask清单**:
- [ ] 实现GetAvailableGPUs方法（1小时）
  - 查询条件：host_id = ? AND status = 'available' AND health_status = 'healthy'
  - 限制数量：LIMIT ?
  - 使用FOR UPDATE行锁（在事务中）

- [ ] 添加行锁支持（30分钟）
  - 防止并发分配冲突

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 查询逻辑正确
- [ ] 行锁机制有效
- [ ] 单元测试通过

---

### Story 2.2.3: 实现GPU分配功能（8小时）

#### Task 2.2.3.1: 实现单GPU分配（3小时）

**Subtask清单**:
- [ ] 实现Allocate方法（1小时）
  - 参数：gpuID, envID
  - 更新status为'allocated'
  - 设置allocated_to和allocated_at

- [ ] 实现事务支持（1小时）
  - 使用数据库事务确保原子性
  - 失败时回滚

- [ ] 编写测试用例（1小时）
  - 测试正常分配
  - 测试重复分配
  - 测试事务回滚

**验收标准**:
- [ ] 分配逻辑正确
- [ ] 事务处理完善
- [ ] 单元测试通过

---

#### Task 2.2.3.2: 实现批量分配（3小时）

**Subtask清单**:
- [ ] 实现BatchAllocate方法（1.5小时）
  - 参数：gpuIDs []uint, envID
  - 批量更新GPU状态
  - 使用事务确保全部成功或全部失败

- [ ] 实现原子性保证（1小时）
  - 任何一个GPU分配失败则全部回滚

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 批量分配正确
- [ ] 原子性保证有效
- [ ] 单元测试通过

---

#### Task 2.2.3.3: 实现GPU释放（2小时）

**Subtask清单**:
- [ ] 实现Release方法（1小时）
  - 更新status为'available'
  - 清空allocated_to和allocated_at

- [ ] 实现资源回收逻辑（30分钟）
  - 批量释放支持

- [ ] 编写测试用例（30分钟）

**验收标准**:
- [ ] 释放逻辑正确
- [ ] 单元测试通过

---

### Story 2.2.4: 并发安全测试（4小时）

#### Task 2.2.4.1: 编写并发测试（3小时）

**Subtask清单**:
- [ ] 测试并发分配场景（1小时）
  - 多个goroutine同时分配同一GPU
  - 验证只有一个成功

- [ ] 测试资源竞争场景（1小时）
  - 模拟高并发场景

- [ ] 测试死锁场景（1小时）
  - 验证不会出现死锁

**验收标准**:
- [ ] 并发测试通过
- [ ] 无资源竞争
- [ ] 无死锁

---

#### Task 2.2.4.2: 性能测试（1小时）

**Subtask清单**:
- [ ] 测试分配性能（30分钟）
- [ ] 测试查询性能（30分钟）

**验收标准**:
- [ ] 性能达标

---

## 协作和依赖

### 依赖其他员工的工作
- **员工C**: 必须等待Host模型完成（HostID外键）

### 其他员工依赖你的工作
- **员工E**: Environment创建时需要调用GPU分配方法

### 协作点
- 与**员工C**协作：等待Host模型完成通知
- 与**员工E**协作：定义GPU分配接口，确保环境创建时能正确调用

---

## 验收检查清单

### 代码质量
- [ ] 代码格式化和检查通过
- [ ] 所有方法有注释
- [ ] 错误处理完善

### 测试
- [ ] 单元测试通过
- [ ] 并发测试通过
- [ ] 性能测试达标
- [ ] 测试覆盖率>80%

### 功能
- [ ] GPU分配正确
- [ ] 并发安全
- [ ] 事务处理完善

### Git提交
- [ ] 提交信息：`feat: 实现GPU模型和分配功能`

---

## 注意事项

1. **并发安全**: 使用FOR UPDATE行锁防止并发冲突
2. **事务处理**: 所有分配操作必须在事务中
3. **状态管理**: 严格控制GPU状态转换
4. **性能优化**: 批量操作优于单个操作
5. **错误处理**: 区分不同的失败原因

---

## 技术要点

### GPU分配示例
```go
func (d *GPUDao) Allocate(gpuID uint, envID string) error {
    tx := d.db.Begin()
    defer tx.Rollback()

    // 使用行锁
    var gpu entity.GPU
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("id = ? AND status = ?", gpuID, "available").
        First(&gpu).Error

    if err != nil {
        return err
    }

    // 更新状态
    gpu.Status = "allocated"
    gpu.AllocatedTo = &envID
    now := time.Now()
    gpu.AllocatedAt = &now

    if err := tx.Save(&gpu).Error; err != nil {
        return err
    }

    return tx.Commit().Error
}
```

---

## 进度报告

请在完成每个Task后更新`团队协作文档.md`。

---

## 联系方式

遇到问题请联系：
- 员工C（Host依赖）
- 员工E（Environment协作）
- 项目经理
