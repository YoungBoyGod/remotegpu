# Claude 操作审计日志

本文档记录 Claude 在项目中的所有修改操作，用于代码审计和变更追踪。

---

## 2026-02-04 操作记录

### 1. 按用户过滤机器列表 (`/customer/machines`)

**目标**: 实现客户端机器列表按用户过滤功能

**修改文件**:

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `internal/dao/allocation_repo.go` | 新增方法 | 添加 `FindActiveByCustomerID` 方法，支持分页查询用户的活跃分配 |
| `internal/service/allocation/allocation_service.go` | 新增方法 | 添加 `ListByCustomerID` 方法 |
| `internal/controller/v1/customer/my_machine_controller.go` | 重构 | 从JWT获取userID进行过滤，修复原实现返回所有机器的数据泄露问题 |
| `internal/router/router.go` | 修复 | 添加 AllocationService 参数传递 |

**代码审查发现的问题**:
- 原 `List` 方法返回所有机器，存在数据泄露风险
- Host 实体没有 GPUModel/GPUCount 字段，改用 TotalCPU/TotalMemoryGB

**Git 提交**:
- `672f585` feat(dao): 添加根据客户ID查询活跃分配方法
- `ffdea23` feat(service): 添加按客户ID查询分配列表方法
- `0639f87` feat(controller): 实现按用户过滤机器列表
- `5c38f64` fix(router): 修复MyMachineController初始化参数

---

### 2. 任务创建/停止绑定用户 (`/customer/tasks/*`)

**目标**: 实现任务的用户绑定和权限校验

**修改文件**:

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `internal/model/entity/common.go` | 新增 | 添加 `ErrUnauthorized` 和 `ErrNotFound` 错误定义 |
| `internal/dao/task_repo.go` | 新增方法 | 添加 `FindByID` 方法用于权限校验 |
| `internal/service/task/task_service.go` | 新增方法 | 添加 `GetTask` 和 `StopTaskWithAuth` 方法 |
| `internal/controller/v1/task/task_controller.go` | 重构 | List/CreateTraining/Stop 均从JWT获取userID |

**代码审查发现的问题**:
- 原 `List` 方法使用硬编码 `userID = uint(1)`
- 原 `CreateTraining` 方法使用硬编码 `CustomerID = 1`
- 原 `Stop` 方法无权限校验，存在越权风险

**Git 提交**:
- `6471ebf` feat(entity): 添加通用错误定义
- `49a09d8` feat(dao): 添加根据ID查询任务方法
- `6bfa465` feat(service): 添加任务权限校验方法
- `18cdd15` feat(controller): 实现任务创建/停止绑定用户

---

### 3. 数据集隔离与挂载 (`/customer/datasets/*`)

**目标**: 实现数据集的用户隔离和挂载权限校验

**修改文件**:

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `internal/dao/dataset_repo.go` | 新增方法 | 添加 `FindByID` 方法 |
| `internal/service/dataset/dataset_service.go` | 新增方法 | 添加 `GetDataset` 和 `ValidateOwnership` 方法 |
| `internal/controller/v1/dataset/dataset_controller.go` | 重构 | List/Mount 从JWT获取userID，Mount添加所有权校验 |

**代码审查发现的问题**:
- 原 `List` 方法使用硬编码 `userID = uint(1)`
- 原 `Mount` 方法无权限校验，存在越权风险

**Git 提交**:
- `d12e74f` feat(dao): 添加根据ID查询数据集方法
- `c629ea1` feat(service): 添加数据集权限校验方法
- `ef59d95` feat(controller): 实现数据集隔离与挂载权限校验

---

### 4. 监控快照接入真实数据源 (`/admin/monitoring/realtime`)

**目标**: 将监控快照从Mock数据改为真实数据源

**修改文件**:

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `internal/service/ops/monitor_service.go` | 重构 | MonitorService 依赖 MachineService，GetGlobalSnapshot 返回真实数据 |
| `internal/router/router.go` | 修复 | 传入 MachineService 依赖 |

**代码审查发现的问题**:
- 原 `GetGlobalSnapshot` 返回硬编码的Mock数据

**待完成 (TODO)**:
- 接入 Redis 缓存，设置采样频率（如30秒）避免频繁查询数据库
- 从监控系统获取 GPU 利用率

**Git 提交**:
- `18d1867` feat(service): 监控快照接入真实数据源
- `c4776f6` fix(router): 修复MonitorService初始化参数

---

### 5. 告警列表完善 (`/admin/alerts`)

**目标**: 完善告警列表功能，添加分页、筛选和确认机制

**修改文件**:

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `internal/dao/ops_repo.go` | 新增方法 | 添加 `ListAlerts` 和 `AcknowledgeAlert` 方法 |
| `internal/service/ops/ops_service.go` | 新增方法 | 添加 `ListAlerts` 和 `AcknowledgeAlert` 方法 |
| `internal/controller/v1/ops/ops_controller.go` | 重构 | List 添加分页筛选，新增 Acknowledge 方法 |
| `internal/router/router.go` | 新增路由 | 添加 `POST /admin/alerts/:id/acknowledge` |

**代码审查发现的问题**:
- 原 `List` 方法无分页和筛选功能
- 缺少告警确认机制

**Git 提交**:
- `b9573e6` feat(dao): 添加告警分页查询和确认方法
- `27faf7b` feat(service): 添加告警分页查询和确认方法
- `0aa5860` feat(controller): 完善告警列表功能
- `91ea958` feat(router): 添加告警确认路由

---

## 注释规范

所有关键函数添加以下格式的注释：

```go
// FunctionName 函数描述
// @author Claude
// @description 详细描述
// @param paramName 参数说明
// @return 返回值说明
// @reason 修改原因（如有）
// @modified 2026-02-04
// TODO: 待完成事项（如有）
func FunctionName() {}
```

---

## 待完成任务

根据 `docs/任务清单.md`，以下任务待完成：

- [x] 告警列表 `/admin/alerts` - 已完成分页、筛选、确认功能
- [ ] 回收机器 `/admin/machines/:id/reclaim`
- [ ] Swagger/OpenAPI 文档完善
