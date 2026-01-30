# ResourceQuota API 文档

## 概述

ResourceQuota 服务提供资源配额管理功能，支持用户级别和工作空间级别的配额设置、查询、更新和检查。

## 目录

- [数据结构](#数据结构)
- [服务方法](#服务方法)
- [使用示例](#使用示例)
- [错误处理](#错误处理)
- [并发安全](#并发安全)

---

## 数据结构

### ResourceQuota

资源配额实体，存储用户或工作空间的资源限制。

```go
type ResourceQuota struct {
    ID          uint      // 配额ID
    CustomerID  uint      // 客户ID（必填）
    WorkspaceID *uint     // 工作空间ID（可选，为nil表示用户级别配额）
    CPU         int       // CPU核心数配额
    Memory      int64     // 内存配额（MB）
    GPU         int       // GPU数量配额
    Storage     int64     // 存储配额（GB）
    CreatedAt   time.Time // 创建时间
    UpdatedAt   time.Time // 更新时间
}
```

**配额级别说明：**
- **用户级别配额**：`WorkspaceID` 为 `nil`，限制用户在所有工作空间中的总资源使用
- **工作空间级别配额**：`WorkspaceID` 不为 `nil`，限制特定工作空间的资源使用

### ResourceRequest

资源请求结构，用于配额检查。

```go
type ResourceRequest struct {
    CPU     int   // 请求的CPU核心数
    Memory  int64 // 请求的内存（MB）
    GPU     int   // 请求的GPU数量
    Storage int64 // 请求的存储（GB）
}
```

### UsedResources

已使用资源结构，表示当前已分配的资源。

```go
type UsedResources struct {
    CPU     int   // 已使用的CPU核心数
    Memory  int64 // 已使用的内存（MB）
    GPU     int   // 已使用的GPU数量
    Storage int64 // 已使用的存储（GB）
}
```

### QuotaExceededError

配额超限错误，当资源请求超过可用配额时返回。

```go
type QuotaExceededError struct {
    Resource  string // 超限的资源类型（CPU/Memory/GPU/Storage）
    Requested int64  // 请求的资源量
    Available int64  // 可用的资源量
}
```

---

## 服务方法

### 1. SetQuota

设置或更新资源配额。如果配额已存在则更新，否则创建新配额。

**方法签名：**
```go
func (s *ResourceQuotaService) SetQuota(quota *entity.ResourceQuota) error
```

**参数：**
- `quota`: 资源配额对象，必须包含 `CustomerID` 和资源值

**返回值：**
- `error`: 错误信息，成功时为 `nil`

**验证规则：**
- 所有资源值（CPU、Memory、GPU、Storage）不能为负数
- 自动检测是否已存在配额（根据 CustomerID 和 WorkspaceID）

**示例：**
```go
// 设置用户级别配额
quota := &entity.ResourceQuota{
    CustomerID:  100,
    WorkspaceID: nil,
    CPU:         16,
    Memory:      32768,
    GPU:         4,
    Storage:     1000,
}
err := service.SetQuota(quota)
if err != nil {
    log.Printf("设置配额失败: %v", err)
}
```

---

### 2. GetQuota

获取资源配额。

**方法签名：**
```go
func (s *ResourceQuotaService) GetQuota(customerID uint, workspaceID *uint) (*entity.ResourceQuota, error)
```

**参数：**
- `customerID`: 客户ID
- `workspaceID`: 工作空间ID指针，为 `nil` 时获取用户级别配额

**返回值：**
- `*entity.ResourceQuota`: 配额对象
- `error`: 错误信息，配额不存在时返回 `gorm.ErrRecordNotFound`

**示例：**
```go
// 获取用户级别配额
quota, err := service.GetQuota(customerID, nil)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        log.Println("配额不存在")
    }
    return err
}

// 获取工作空间级别配额
workspaceID := uint(10)
quota, err := service.GetQuota(customerID, &workspaceID)
```

---

### 3. GetQuotaInTx

在事务中获取资源配额，使用悲观锁（FOR UPDATE）保证并发安全。

**方法签名：**
```go
func (s *ResourceQuotaService) GetQuotaInTx(tx *gorm.DB, customerID uint, workspaceID *uint) (*entity.ResourceQuota, error)
```

**参数：**
- `tx`: GORM 事务对象（不能为 `nil`）
- `customerID`: 客户ID
- `workspaceID`: 工作空间ID指针

**返回值：**
- `*entity.ResourceQuota`: 配额对象（已加锁）
- `error`: 错误信息

**使用场景：**
- 环境创建时的配额检查
- 需要原子性操作的场景
- 防止并发创建导致的配额超限

**示例：**
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

quota, err := service.GetQuotaInTx(tx, customerID, nil)
if err != nil {
    tx.Rollback()
    return err
}

// 执行其他操作...
tx.Commit()
```

---

### 4. UpdateQuota

更新资源配额。

**方法签名：**
```go
func (s *ResourceQuotaService) UpdateQuota(quota *entity.ResourceQuota) error
```

**参数：**
- `quota`: 要更新的配额对象，必须包含有效的 `ID`

**返回值：**
- `error`: 错误信息

**验证规则：**
- 所有资源值不能为负数
- 新配额不能小于已使用的资源量

**示例：**
```go
quota, _ := service.GetQuota(customerID, nil)
quota.CPU = 32
quota.Memory = 65536

err := service.UpdateQuota(quota)
if err != nil {
    log.Printf("更新配额失败: %v", err)
}
```

---

### 5. DeleteQuota

删除资源配额。

**方法签名：**
```go
func (s *ResourceQuotaService) DeleteQuota(id uint) error
```

**参数：**
- `id`: 配额ID

**返回值：**
- `error`: 错误信息

**示例：**
```go
err := service.DeleteQuota(quotaID)
if err != nil {
    log.Printf("删除配额失败: %v", err)
}
```

---

### 6. CheckQuota

检查资源配额是否足够。

**⚠️ 注意：此方法不保证并发安全。如需在环境创建等场景中使用，请使用 `CheckQuotaInTx`。**

**方法签名：**
```go
func (s *ResourceQuotaService) CheckQuota(customerID uint, workspaceID *uint, request *ResourceRequest) (bool, error)
```

**参数：**
- `customerID`: 客户ID
- `workspaceID`: 工作空间ID指针
- `request`: 资源请求对象

**返回值：**
- `bool`: 配额是否足够（`true` 表示足够，`false` 表示不足）
- `error`: 错误信息，配额不足时返回 `QuotaExceededError`

**示例：**
```go
request := &ResourceRequest{
    CPU:     8,
    Memory:  16384,
    GPU:     2,
    Storage: 500,
}

ok, err := service.CheckQuota(customerID, nil, request)
if !ok {
    if quotaErr, ok := err.(*QuotaExceededError); ok {
        log.Printf("配额不足: %s", quotaErr.Error())
    }
    return err
}
```

---

### 7. CheckQuotaInTx

在事务中检查资源配额是否足够，使用悲观锁保证并发安全。

**方法签名：**
```go
func (s *ResourceQuotaService) CheckQuotaInTx(tx *gorm.DB, customerID uint, workspaceID *uint, request *ResourceRequest) (bool, error)
```

**参数：**
- `tx`: GORM 事务对象（不能为 `nil`）
- `customerID`: 客户ID
- `workspaceID`: 工作空间ID指针
- `request`: 资源请求对象

**返回值：**
- `bool`: 配额是否足够
- `error`: 错误信息

**使用场景：**
- 环境创建时的配额检查
- 需要原子性操作的场景

**示例：**
```go
tx := db.Begin()

request := &ResourceRequest{
    CPU:     8,
    Memory:  16384,
    GPU:     2,
    Storage: 500,
}

ok, err := service.CheckQuotaInTx(tx, customerID, nil, request)
if !ok {
    tx.Rollback()
    return err
}

// 创建环境...
env := &entity.Environment{...}
if err := tx.Create(env).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

---

### 8. GetUsedResources

获取已使用的资源量。

**方法签名：**
```go
func (s *ResourceQuotaService) GetUsedResources(customerID uint, workspaceID *uint) (*UsedResources, error)
```

**参数：**
- `customerID`: 客户ID
- `workspaceID`: 工作空间ID指针

**返回值：**
- `*UsedResources`: 已使用资源对象
- `error`: 错误信息

**计算规则：**
- 统计所有状态为 `running` 和 `creating` 的环境
- 用户级别配额：统计该用户所有工作空间的资源使用
- 工作空间级别配额：仅统计指定工作空间的资源使用

**示例：**
```go
used, err := service.GetUsedResources(customerID, nil)
if err != nil {
    return err
}

log.Printf("已使用: CPU=%d, Memory=%d, GPU=%d, Storage=%d",
    used.CPU, used.Memory, used.GPU, used.Storage)
```

---

### 9. GetAvailableQuota

获取可用配额（总配额 - 已使用资源）。

**方法签名：**
```go
func (s *ResourceQuotaService) GetAvailableQuota(customerID uint, workspaceID *uint) (*entity.ResourceQuota, error)
```

**参数：**
- `customerID`: 客户ID
- `workspaceID`: 工作空间ID指针

**返回值：**
- `*entity.ResourceQuota`: 可用配额对象
- `error`: 错误信息

**计算规则：**
- 可用配额 = 总配额 - 已使用资源
- 如果计算结果为负数，则设为 0（表示已超额使用）

**示例：**
```go
available, err := service.GetAvailableQuota(customerID, nil)
if err != nil {
    return err
}

log.Printf("可用配额: CPU=%d, Memory=%d, GPU=%d, Storage=%d",
    available.CPU, available.Memory, available.GPU, available.Storage)
```

---

## 使用示例

### 完整的环境创建流程（带配额检查）

```go
func CreateEnvironment(customerID uint, workspaceID *uint, envReq *EnvironmentRequest) error {
    service := NewResourceQuotaService()
    db := database.GetDB()

    // 开始事务
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 1. 检查配额（使用事务和悲观锁）
    request := &ResourceRequest{
        CPU:     envReq.CPU,
        Memory:  envReq.Memory,
        GPU:     envReq.GPU,
        Storage: envReq.Storage,
    }

    ok, err := service.CheckQuotaInTx(tx, customerID, workspaceID, request)
    if !ok {
        tx.Rollback()
        return fmt.Errorf("配额检查失败: %v", err)
    }

    // 2. 创建环境
    env := &entity.Environment{
        CustomerID:  customerID,
        WorkspaceID: workspaceID,
        CPU:         envReq.CPU,
        Memory:      envReq.Memory,
        GPU:         envReq.GPU,
        Storage:     &envReq.Storage,
        Status:      "creating",
    }

    if err := tx.Create(env).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 3. 提交事务
    if err := tx.Commit().Error; err != nil {
        return err
    }

    return nil
}
```

### 配额管理示例

```go
func ManageQuota(customerID uint) error {
    service := NewResourceQuotaService()

    // 1. 设置用户级别配额
    quota := &entity.ResourceQuota{
        CustomerID:  customerID,
        WorkspaceID: nil,
        CPU:         32,
        Memory:      65536,
        GPU:         8,
        Storage:     2000,
    }

    if err := service.SetQuota(quota); err != nil {
        return fmt.Errorf("设置配额失败: %v", err)
    }

    // 2. 查询已使用资源
    used, err := service.GetUsedResources(customerID, nil)
    if err != nil {
        return err
    }

    log.Printf("已使用资源: CPU=%d, Memory=%d, GPU=%d, Storage=%d",
        used.CPU, used.Memory, used.GPU, used.Storage)

    // 3. 查询可用配额
    available, err := service.GetAvailableQuota(customerID, nil)
    if err != nil {
        return err
    }

    log.Printf("可用配额: CPU=%d, Memory=%d, GPU=%d, Storage=%d",
        available.CPU, available.Memory, available.GPU, available.Storage)

    return nil
}
```

---

## 错误处理

### 常见错误类型

1. **配额不存在**
   ```go
   if errors.Is(err, gorm.ErrRecordNotFound) {
       // 处理配额不存在的情况
   }
   ```

2. **配额超限**
   ```go
   if quotaErr, ok := err.(*QuotaExceededError); ok {
       log.Printf("资源 %s 配额不足: 需要 %d, 可用 %d",
           quotaErr.Resource, quotaErr.Requested, quotaErr.Available)
   }
   ```

3. **验证错误**
   - "配额值不能为负数"
   - "CPU配额不能小于已使用量"
   - "内存配额不能小于已使用量"
   - "GPU配额不能小于已使用量"
   - "存储配额不能小于已使用量"

4. **事务错误**
   - "事务不能为空"

---

## 并发安全

### 并发场景说明

在多用户同时创建环境的场景下，可能出现以下问题：

1. **竞态条件**：两个请求同时检查配额，都认为配额足够，但实际创建后超过配额
2. **超额分配**：配额为 10，两个请求各需要 6，同时通过检查，实际分配了 12

### 解决方案

使用 `CheckQuotaInTx` 和事务来保证并发安全：

```go
// ❌ 错误示例：不安全
ok, err := service.CheckQuota(customerID, nil, request)
if ok {
    // 创建环境（可能导致超额）
    CreateEnvironment(...)
}

// ✅ 正确示例：安全
tx := db.Begin()
ok, err := service.CheckQuotaInTx(tx, customerID, nil, request)
if ok {
    // 在同一事务中创建环境
    tx.Create(&env)
    tx.Commit()
} else {
    tx.Rollback()
}
```

### 悲观锁机制

`GetQuotaInTx` 使用 `SELECT ... FOR UPDATE` 锁定配额记录：

```sql
SELECT * FROM resource_quotas
WHERE customer_id = ? AND workspace_id IS NULL
FOR UPDATE;
```

这确保在事务提交前，其他事务无法读取或修改该配额记录。

---

## 最佳实践

1. **配额设置**
   - 用户级别配额应大于所有工作空间配额之和
   - 定期审查和调整配额

2. **配额检查**
   - 环境创建时必须使用 `CheckQuotaInTx`
   - 简单查询可使用 `CheckQuota`

3. **错误处理**
   - 始终检查 `QuotaExceededError` 并提供友好的错误信息
   - 记录配额超限事件用于审计

4. **性能优化**
   - 避免频繁调用 `GetUsedResources`
   - 考虑使用缓存存储配额信息

5. **监控告警**
   - 监控配额使用率
   - 当使用率超过 80% 时发送告警

---

## 版本历史

- **v1.0** (2026-01-30): 初始版本，包含基本的配额管理功能
- 支持用户级别和工作空间级别配额
- 实现并发安全的配额检查机制

---

## 相关文档

- [Environment API 文档](./environment.md)
- [Workspace API 文档](./workspace.md)
- [数据库设计文档](../database/schema.md)

---

**文档维护者**: B1 (A2 员工)
**最后更新**: 2026-01-30
**联系方式**: 如有问题请提交 Issue
