# Gemini 操作日志

## 2026-02-04 Session

### 1. 功能梳理与任务清单确认
- **时间**: Session Start
- **操作**: 
  - 读取并分析 `backend/docs/功能梳理.md`。
  - 识别出缺失的 B2B 关键模块：镜像管理、SSH 密钥管理、系统审计。
  - 根据用户反馈剔除计费相关内容。
  - 更新 `backend/docs/功能梳理.md` 和 `backend/docs/任务清单.md`。
- **修改文件**:
  - `backend/docs/功能梳理.md`
  - `backend/docs/任务清单.md`

### 2. 登录与鉴权闭环 (P0)
- **时间**: P0 Phase
- **操作**:
  - 实现刷新令牌逻辑：`AuthService.RefreshToken`，支持 Redis 存储与轮换。
  - 更新 `AuthController.Refresh` 接口。
  - 在 `AuthService.Login` 和 `middleware.Auth` 中增加账号状态 (`status=active`) 校验。
  - 优化错误处理，引入 `pkg/errors` 统一管理错误码。
- **修改文件**:
  - `backend/api/v1/auth.go` (新增 `RefreshTokenRequest`)
  - `backend/internal/service/auth/auth_service.go`
  - `backend/internal/controller/v1/auth/auth_controller.go`
  - `backend/internal/middleware/auth.go`
  - `backend/internal/router/router.go` (注入 DB 依赖)

### 3. 权限与管理端基础 (P1)
- **时间**: P1 Phase - Dashboard & Ops
- **操作**:
  - 确认 `RequireAdmin` 中间件已存在。
  - 实现仪表盘数据源：`GetGPUTrend` (Mock), `GetRecentAllocations` (Real)。
  - 实现镜像同步接口：`/admin/images/sync` (Stub)。
  - 完善 `MonitorService`，接入 `MachineService` 获取真实机器状态快照。
  - 为 `OpsController` (Dashboard, Monitor, Alert) 添加 Swagger 注解。
- **修改文件**:
  - `backend/internal/service/ops/monitor_service.go`
  - `backend/internal/controller/v1/ops/ops_controller.go`
  - `backend/internal/controller/v1/ops/image_controller.go`
  - `backend/internal/service/image/image_service.go`
  - `backend/internal/router/router.go`

### 4. 机器管理闭环 (P1)
- **时间**: P1 Phase - Machine
- **操作**:
  - 实现机器批量导入：`/admin/machines/import`，支持 JSON 列表导入。
  - 增强机器分配逻辑：`AllocateMachine` 增加事务处理、状态检查 (idle/online) 和租期校验。
  - 增强机器回收逻辑：`ReclaimMachine` 增加审计日志记录 (`AuditService`)。
  - 定义新的错误码：`ErrorMachineNotAvailable`, `ErrorAllocationNotFound`。
- **修改文件**:
  - `backend/api/v1/machine.go`
  - `backend/internal/controller/v1/machine/machine_controller.go`
  - `backend/internal/service/machine/machine_service.go`
  - `backend/internal/service/allocation/allocation_service.go`
  - `backend/pkg/errors/errors.go`

### 5. 文档完善 (P3)
- **时间**: Final Phase
- **操作**:
  - 为主要控制器添加 Swagger 注解：Auth, Machine, Dashboard, Monitor, Alert, Image。
  - 修复 `SSHKeyController` 中潜在的类型断言 panic 问题。
  - 提交所有更改并更新任务清单状态。
- **修改文件**:
  - `backend/internal/controller/v1/auth/auth_controller.go`
  - `backend/internal/controller/v1/machine/machine_controller.go`
  - `backend/internal/controller/v1/ops/ops_controller.go`
  - `backend/internal/controller/v1/ops/image_controller.go`
  - `backend/internal/controller/v1/customer/sshkey_controller.go`
  - `backend/docs/任务清单.md`
