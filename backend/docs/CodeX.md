# CodeX 操作审计

> 说明：记录 CodeX 在仓库中的关键变更，便于审计追踪。

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

## 说明

- 所有新增注释均带 `CodeX` 前缀与日期。
