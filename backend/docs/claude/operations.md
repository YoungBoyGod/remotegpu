# Claude 操作流水账

本文档记录 Claude 在项目中的所有操作，用于审计和变更追踪。

**维护者**: Claude
**创建日期**: 2026-02-05

---

## 2026-02-05

### [DOC] 创建项目管理文件结构

**时间**: 2026-02-05

**操作**:
1. 创建 `docs/claude/Plan.md` - 项目规划
2. 创建 `docs/claude/doing.md` - 当前任务
3. 创建 `docs/claude/operations.md` - 操作流水账

**Git 提交**:
- `7288b9b` docs(claude): 添加待完成任务清单
- `7d8ff52` docs(claude): 添加项目管理文件

---

## 2026-02-04

### [FEAT] 1. 按用户过滤机器列表

**目标**: `/customer/machines` 接口按用户过滤

**修改文件**:
| 文件 | 类型 | 说明 |
|------|------|------|
| `internal/dao/allocation_repo.go` | 新增 | `FindActiveByCustomerID` |
| `internal/service/allocation/allocation_service.go` | 新增 | `ListByCustomerID` |
| `internal/controller/v1/customer/my_machine_controller.go` | 重构 | JWT用户过滤 |
| `internal/router/router.go` | 修复 | 参数传递 |

**发现问题**: 原实现返回所有机器，存在数据泄露

**Git**: `672f585`, `ffdea23`, `0639f87`, `5c38f64`

---

### [FEAT] 2. 任务创建/停止绑定用户

**目标**: `/customer/tasks/*` 接口用户绑定

**修改文件**:
| 文件 | 类型 | 说明 |
|------|------|------|
| `internal/model/entity/common.go` | 新增 | 错误定义 |
| `internal/dao/task_repo.go` | 新增 | `FindByID` |
| `internal/service/task/task_service.go` | 新增 | 权限校验方法 |
| `internal/controller/v1/task/task_controller.go` | 重构 | JWT用户绑定 |

**发现问题**: 硬编码 `userID = 1`，无权限校验

**Git**: `6471ebf`, `49a09d8`, `6bfa465`, `18cdd15`

---

### [FEAT] 3. 数据集隔离与挂载

**目标**: `/customer/datasets/*` 接口权限校验

**修改文件**:
| 文件 | 类型 | 说明 |
|------|------|------|
| `internal/dao/dataset_repo.go` | 新增 | `FindByID` |
| `internal/service/dataset/dataset_service.go` | 新增 | 所有权校验 |
| `internal/controller/v1/dataset/dataset_controller.go` | 重构 | 权限校验 |

**发现问题**: 硬编码用户ID，无挂载权限校验

**Git**: `d12e74f`, `c629ea1`, `ef59d95`

---

### [FEAT] 4. 监控快照接入真实数据

**目标**: `/admin/monitoring/realtime` 返回真实数据

**修改文件**:
| 文件 | 类型 | 说明 |
|------|------|------|
| `internal/service/ops/monitor_service.go` | 重构 | 依赖 MachineService |
| `internal/router/router.go` | 修复 | 参数传递 |

**发现问题**: 返回硬编码 Mock 数据

**Git**: `18d1867`, `c4776f6`

---

### [FEAT] 5. 告警列表完善

**目标**: `/admin/alerts` 分页筛选和确认

**修改文件**:
| 文件 | 类型 | 说明 |
|------|------|------|
| `internal/dao/ops_repo.go` | 新增 | 分页查询、确认方法 |
| `internal/service/ops/ops_service.go` | 新增 | 对应服务方法 |
| `internal/controller/v1/ops/ops_controller.go` | 重构 | 分页筛选 |
| `internal/router/router.go` | 新增 | 确认路由 |

**Git**: `b9573e6`, `27faf7b`, `0aa5860`, `91ea958`

---

### [FIX] 6. ReclaimMachine Bug 修复

**目标**: 修复审计代码不可达 bug

**修改文件**:
| 文件 | 类型 | 说明 |
|------|------|------|
| `internal/service/allocation/allocation_service.go` | 修复 | 事务返回值处理 |

**问题**: `return s.db.Transaction(...)` 导致审计代码永不执行

**Git**: `4b39131`

---

## 注释规范

```go
// FunctionName 描述
// @author Claude
// @modified 2026-02-04
func FunctionName() {}
```
