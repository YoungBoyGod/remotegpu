# CodeX 操作流水账

> 说明：记录 CodeX 在仓库中的关键变更，便于审计追踪。

## 2026-02-05

- 任务停止 Agent 集成说明
  - 新增说明与 Review 文档：`docs/CodeX_任务停止Agent集成.md`

- 任务停止 Agent 集成实现
  - Task 增加 `process_id` 字段：`internal/model/entity/task.go`
  - 任务停止校验与错误回写：`internal/service/task/task_service.go`, `internal/dao/task_repo.go`
  - Agent StopProcess 传入 `process_id` 并校验：`internal/service/ops/agent_service.go`
  - SQL 迁移补充：`sql/15_task_process_id.sql`

- 机器详情完善
  - 详情 API 预加载分配与客户信息：`internal/dao/machine_repo.go`
  - 详情页字段补全与格式化：`../frontend/src/views/admin/MachineDetailView.vue`
  - 机器类型字段补齐：`../frontend/src/types/machine.ts`
  - Review 文档：`docs/CodeX_机器详情Review.md`

- 测试执行
  - 执行 `go test ./...` 失败，报错：`pkg/infra/volume/example_test.go` 引用缺失模块 `github.com/YoungBoyGod/remotegpu/pkg/volume`

- 测试修复
  - 修正 volume 示例测试引用路径：`pkg/infra/volume/example_test.go`
  - 更新登录失败测试断言为业务错误码：`internal/controller/v1/auth/auth_controller_test.go`, `internal/controller/v1/auth/auth_controller_http_test.go`
  - 执行 `go test ./...` 通过

- OpenAPI 文档更新
  - 错误响应补充业务错误码与示例：`docs/openapi.yaml`

- 文档 Review
  - 任务队列设计 Review 文档：`/home/luo/code/remotegpu/agent/docs/CodeX_task-queue-design_review.md`

- 任务队列设计修订
  - 统一为纯 Pull 模式，移除 WebSocket 描述与推送接口说明：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - 统一任务状态机与租约字段，补齐 assigned/lease/attempt 描述：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - Agent 本地 schema 增加 attempt_id/assigned_at，API 示例补充租约字段：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - Agent API 幂等与错误码约定：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - Agent API 独立业务错误码表与示例响应：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - Agent 专用错误码文件：`/home/luo/code/remotegpu/agent/internal/errors/errors.go`
  - Agent API 响应封装与错误码接入：`/home/luo/code/remotegpu/agent/cmd/handlers.go`, `/home/luo/code/remotegpu/agent/cmd/response.go`
  - 错误码表补充内部错误并修正响应字段：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - 后端 Agent 客户端校验业务错误码：`internal/agent/http_client.go`, `internal/agent/grpc_client.go`
  - 认领 API 改为 POST + 租约续约接口 + 结果存储/抢占/状态机补齐：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - 补充 Sidecar 与长时间任务处理原则：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - Agent 本地队列 YAML 示例：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - 任务配置 YAML 示例（单任务/DAG）：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`
  - 循环任务支持策略与示例：`/home/luo/code/remotegpu/agent/docs/task-queue-design.md`

- 规划与待办更新
  - 计划状态调整：`docs/codex/Plan.md`
  - 进行中状态调整：`docs/codex/doing.md`
  - 待办清单更新：`docs/codex/todo.md`

## 2026-02-04

- 登录/鉴权 Review 与改进
  - 新增审计文档：`docs/登录鉴权Review.md`
  - 统一登录错误类型与刷新令牌逻辑，增加 refresh token 轮换与兜底：`internal/service/auth/auth_service.go`
  - 控制器统一错误映射并支持 refresh token：`internal/controller/v1/auth/auth_controller.go`
  - 鉴权中间件加入账号状态校验：`internal/middleware/auth.go`

- SSHKey / 审计 / 镜像 Review 与改进
  - 新增审计文档：`docs/SSHKey审计镜像Review.md`
  - SSHKey 获取 userID 方式修正，避免断言 panic：`internal/controller/v1/customer/sshkey_controller.go`
  - SSH 公钥解析改为 `ssh.ParseAuthorizedKey`：`internal/service/sshkey/sshkey_service.go`
  - 审计详情 JSON 序列化失败直接返回错误：`internal/service/audit/audit_service.go`
  - 审计与镜像列表分页参数归一化：`internal/controller/v1/ops/audit_controller.go`, `internal/controller/v1/ops/image_controller.go`

- 审计日志自动记录
  - Admin 路由挂载审计中间件：`internal/router/router.go`

- 数据集挂载权限校验
  - 挂载前校验机器归属：`internal/controller/v1/dataset/dataset_controller.go`
  - 分配记录新增按用户与主机查询：`internal/dao/allocation_repo.go`
  - 分配服务增加归属校验方法：`internal/service/allocation/allocation_service.go`

- 未完成任务清单
  - 新增待办与实现建议：`docs/codex/todo.md`

- 机器创建校验 / 导入幂等
  - DAO 增加 IP/hostname 查询与批量存在性检查：`internal/dao/machine_repo.go`
  - 创建时校验重复并在导入中跳过重复：`internal/service/machine/machine_service.go`
  - 创建接口重复返回 409：`internal/controller/v1/machine/machine_controller.go`

- 计划与进行中记录
  - 新增计划表：`docs/codex/Plan.md`
  - 新增进行中记录：`docs/codex/doing.md`

- 操作流水账迁移
  - 迁移记录至：`docs/codex/operations.md`

- TODO 汇总补充
  - 同步 gemini/claude 目录中的待办到 `docs/codex/todo.md`

- 前端需求草案
  - 新增前端最小需求文档：`docs/codex/frontend_requirements.md`

- 前端工程检查
  - 核对 `../frontend/src` 结构与已有视图

- 前端认证接口对齐
  - Auth API 路径改为 `/api/v1/auth/*` 并处理 snake_case 响应：`frontend/src/api/auth.ts`
  - 登录后拉取 profile、刷新时更新 refresh token：`frontend/src/stores/auth.ts`

- 前端机器添加表单调整
  - 仅保留基本信息与登录信息字段，移除硬件配置输入：`frontend/src/views/admin/MachineAddView.vue`
  - 新增创建机器 payload 类型：`frontend/src/api/admin.ts`

- 前端 SSH 认证校验
  - 机器添加页要求 SSH 密钥或密码至少填写一个：`frontend/src/views/admin/MachineAddView.vue`
  - 创建 payload 增加 `ssh_key`：`frontend/src/api/admin.ts`

- 前端连接地址校验
  - 机器添加页支持 IP/域名/hostname，校验至少填写一个地址：`frontend/src/views/admin/MachineAddView.vue`

- 前端 SSH 认证校验调整
  - SSH 密钥/密码必填其一改为提交时统一校验：`frontend/src/views/admin/MachineAddView.vue`

- 机器新增后端对齐
  - 新建创建请求结构并支持 hostname/domain：`api/v1/machine.go`
  - Host 记录保存 SSH 登录信息：`internal/model/entity/resource.go`
  - 创建时尝试采集硬件信息（Agent）：`internal/controller/v1/machine/machine_controller.go`
  - AgentService 补充系统信息采集方法：`internal/service/ops/agent_service.go`
  - OpenAPI 更新创建机器字段：`docs/openapi.yaml`

- 用户添加机器流程
  - 新增 enrollment 实体与数据访问：`internal/model/entity/machine_enrollment.go`, `internal/dao/machine_enrollment_repo.go`
  - 后端异步采集与入库逻辑：`internal/service/machine/enrollment_service.go`
  - 客户端新增接口：`internal/controller/v1/customer/machine_enrollment_controller.go`
  - 路由注册与自动迁移：`internal/router/router.go`, `cmd/server.go`, `cmd/tools.go`

- 用户添加机器采集校验
  - 补充采集规格校验与 hostname 回退：`internal/service/machine/enrollment_service.go`
  - 仅在 Agent 数据完整时采用采集结果：`internal/service/machine/enrollment_service.go`
  - Host 记录保存 SSH 登录信息：`internal/service/machine/enrollment_service.go`
  - 新增 Review 文档：`docs/CodeX_用户机器添加Review.md`

- 前端需求补充
  - 机器添加要求 SSH 私钥或密码至少填写一个：`docs/codex/frontend_requirements.md`
  - 机器添加页提示 SSH 私钥：`frontend/src/views/admin/MachineAddView.vue`

- TODO 补充
  - 用户添加机器队列化处理：`docs/codex/todo.md`

- SQL 与迁移补充
  - 新增机器添加任务表与主机 SSH 字段：`sql/14_machine_enrollments.sql`
  - 补充触发器函数保障独立执行：`sql/14_machine_enrollments.sql`

- 客户添加机器前端
  - 新增添加表单与进度页面：`frontend/src/views/customer/MachineEnrollView.vue`, `frontend/src/views/customer/MachineEnrollmentListView.vue`
  - 客户 API 增加 enrollment 接口：`frontend/src/api/customer.ts`
  - 客户路由与侧边栏补充入口：`frontend/src/router/index.ts`, `frontend/src/components/layout/CustomerSidebar.vue`
  - 机器列表页加入快捷入口：`frontend/src/views/customer/MachineListView.vue`

- 任务队列说明
  - 补充用户添加机器队列 worker 作用与场景：`docs/codex/todo.md`

- 任务队列实现
  - 用户添加机器改为 Redis 队列消费：`internal/service/machine/enrollment_service.go`
  - enrollment 待处理重入队：`internal/service/machine/enrollment_service.go`, `internal/dao/machine_enrollment_repo.go`
  - 启动时拉起队列 worker：`internal/router/router.go`
  - 失败重试与延迟重入队：`internal/service/machine/enrollment_service.go`

- 运维工具补充
  - 新增 Redis 连接检查命令：`cmd/tools.go`

- 队列配置
  - enrollment 重试次数与延迟写入配置：`config/config.yaml`, `config/config.go`, `internal/service/machine/enrollment_service.go`
  - 支持 max_retries=0 禁用重试：`internal/service/machine/enrollment_service.go`, `config/config.yaml`

- Redis 检查
  - 使用 `tools redis-check` 尝试连通性，沙箱网络限制导致无法直连外部 Redis（需在 docker 环境执行）。

- 机器分配异步动作
  - 分配/回收动作改为 Redis 队列消费：`internal/service/allocation/allocation_service.go`
  - 动作重试与延迟重入队：`internal/service/allocation/allocation_service.go`
  - 启动时拉起动作 worker：`internal/router/router.go`
  - 队列配置新增 machine_action：`config/config.go`, `config/config.yaml`

- 用户添加机器流程补齐
  - 客户提交接口新增 `POST /customer/machines`：`internal/router/router.go`, `frontend/src/api/customer.ts`
  - 前端需求文档补充提交接口：`docs/codex/frontend_requirements.md`

- 用户添加机器临时跳过采集
  - 新增 `machine_enrollment.skip_collect` 配置：`config/config.yaml`, `config/config.go`
  - 跳过采集直接入库：`internal/service/machine/enrollment_service.go`

- 补采硬件信息
  - Host 增加 needs_collect 标记：`internal/model/entity/resource.go`, `sql/14_machine_enrollments.sql`
  - 支持管理员触发补采：`internal/controller/v1/machine/machine_controller.go`, `internal/router/router.go`
  - 采集更新写入 DB：`internal/service/machine/machine_service.go`, `internal/dao/machine_repo.go`
  - OpenAPI 补充补采与添加任务接口：`docs/openapi.yaml`
  - 前端补采按钮与标记展示：`frontend/src/views/admin/MachineListView.vue`, `frontend/src/api/admin.ts`, `frontend/src/types/machine.ts`
