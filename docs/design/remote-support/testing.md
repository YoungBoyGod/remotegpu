# 远程客户支持平台 — 测试策略与 Review 方案

> 版本：v1.0 | 日期：2026-02-06 | 作者：测试工程师

## 1. 文档概述

### 1.1 目的

本文档为远程客户支持平台制定全面的测试策略和 Code Review 流程，确保平台在功能正确性、安全合规、性能稳定性等方面达到企业级标准。

### 1.2 范围

覆盖以下模块的测试策略：
- 后端 API 服务（Go + Gin + GORM + PostgreSQL）
- 前端应用（Vue 3 + TypeScript + Element Plus）
- Agent 客户端（Go + SQLite + HTTP/gRPC）
- 模块间集成测试
- 端到端测试

### 1.3 参考文档

| 文档 | 路径 |
|------|------|
| 需求分析 | `docs/design/requirements.md` |
| 实现审查 | `docs/design/implementation-review.md` |
| 远程访问设计 | `docs/design/machine-remote-access.md` |

---

## 2. 测试现状分析

### 2.1 已有测试覆盖

| 模块 | 测试文件 | 覆盖范围 | 评价 |
|------|----------|----------|------|
| 后端 JWT | `pkg/auth/jwt_test.go` | Token 生成、解析、验证 | 良好 |
| 后端认证控制器 | `controller/v1/auth/auth_controller_test.go` | 登录、改密、Token 刷新 | 良好 |
| 后端认证 HTTP | `controller/v1/auth/auth_controller_http_test.go` | 真实 HTTP 请求测试 | 良好 |
| 后端审计 | `controller/v1/ops/audit_controller_test.go` | 审计日志查询 | 基础 |
| 后端镜像 | `controller/v1/ops/image_controller_test.go` | 镜像列表 | 基础 |
| 后端 SSH 密钥 | `controller/v1/customer/sshkey_controller_test.go` | SSH 密钥 CRUD | 良好 |
| Agent 执行器 | `agent/internal/executor/executor_test.go` | 命令执行、超时 | 良好 |
| Agent 队列 | `agent/internal/queue/manager_test.go` | 优先级队列、并发 | 良好 |
| Agent 存储 | `agent/internal/store/sqlite_test.go` | SQLite 持久化 | 良好 |
| 前端组件 | `frontend/src/__tests__/App.spec.ts` | 基础组件挂载 | 不足 |
| 前端 E2E | `frontend/e2e/vue.spec.ts` | 端到端基础 | 不足 |

### 2.2 测试技术栈现状

| 模块 | 框架 | 断言库 | Mock 方式 | 数据库 |
|------|------|--------|-----------|--------|
| 后端 | Go `testing` | `testify/assert` + `require` | 依赖注入 + `httptest` | SQLite 内存 |
| Agent | Go `testing` | 标准库 | 直接创建对象 | 临时文件 SQLite |
| 前端 | Vitest | `expect` | 无 | 无 |

### 2.3 测试缺口

基于代码审查，以下关键模块缺少测试：

| 模块 | 缺失测试 | 风险等级 |
|------|----------|----------|
| 机器分配服务 | `AllocationService` 无单元测试 | **高** |
| 机器回收流程 | `ReclaimMachine` 无测试 | **高** |
| 任务停止流程 | `StopTask` / `StopTaskWithAuth` 无测试 | **高** |
| 客户管理控制器 | `CustomerController` 无测试 | 中 |
| 机器管理控制器 | `MachineController` 无测试 | 中 |
| 权限中间件 | `RequireAdmin` 无独立测试 | 中 |
| 前端业务组件 | 所有业务页面无测试 | 中 |
| 前端 API 层 | 请求拦截器、Token 刷新无测试 | 中 |

---

## 3. 功能测试策略

### 3.1 测试分层模型

采用经典测试金字塔模型，按投入比例分配：

```
        /  E2E 测试  \          ~10%  关键业务流程
       / 集成测试      \        ~30%  模块间交互、API 测试
      / 单元测试          \     ~60%  函数/方法级别
```

### 3.2 后端单元测试

#### 3.2.1 测试规范

- 框架：Go 标准 `testing` + `testify/assert` + `testify/require`
- 数据库：SQLite 内存数据库（`:memory:`），与现有测试保持一致
- Mock：通过接口注入，不引入 gomock 等外部 mock 框架
- 命名：`TestXxx_场景描述`，如 `TestAllocateMachine_Success`
- 组织：每个 Service/Controller 对应一个 `_test.go` 文件

#### 3.2.2 必须覆盖的测试用例

**认证模块（`service/auth`）**

| 用例 | 验证点 |
|------|--------|
| 登录成功 | 返回 access_token + refresh_token |
| 密码错误 | 返回 `ErrorPasswordIncorrect` |
| 账号禁用 | 返回 `ErrorUserDisabled` |
| Token 刷新成功 | 旧 token 失效，新 token 有效 |
| Token 刷新过期 | 返回 `ErrorTokenInvalid` |
| 改密成功 | `must_change_password` 置为 false |
| 改密旧密码错误 | 返回错误 |
| 登出后 Token 黑名单 | Token 加入黑名单，后续请求被拒 |

**机器分配模块（`service/allocation`）— 高优先级补充**

| 用例 | 验证点 |
|------|--------|
| 分配成功 | 机器状态变为 `allocated`，创建分配记录 |
| 分配不可用机器 | 状态非 idle/online 时返回错误 |
| 并发分配同一机器 | 只有一个请求成功（需行级锁） |
| 租期参数校验 | `durationMonths < 1` 返回错误 |
| 回收成功 | 分配状态变为 `reclaimed`，机器进入 `maintenance` |
| 回收无分配记录 | 返回 `ErrorAllocationNotFound` |
| 回收时有运行中任务 | 应先停止任务或返回警告 |

**任务管理模块（`service/task`）— 高优先级补充**

| 用例 | 验证点 |
|------|--------|
| 提交任务 | 状态为 `queued` |
| 停止任务（有权限） | 任务状态变为 `stopped` |
| 停止任务（无权限） | 返回 `ErrUnauthorized` |
| 停止已完成任务 | 应返回状态错误 |
| Agent 认领任务 | 返回按优先级排序的任务 |
| 任务租约续期 | `lease_expires_at` 延长 |
| 任务完成 | 状态更新，记录 exit_code |

**客户管理模块（`service/customer` + `controller`）**

| 用例 | 验证点 |
|------|--------|
| 创建客户 | 默认密码哈希存储，`must_change_password=true` |
| 用户名重复 | 返回唯一约束错误 |
| 禁用客户 | 状态变为 `disabled` |
| 禁用当前登录账号 | 应被阻止 |
| 客户列表分页 | 分页参数正确，总数准确 |

**中间件测试（`middleware/`）**

| 用例 | 验证点 |
|------|--------|
| Auth 中间件 — 无 Token | 返回 401 |
| Auth 中间件 — 无效 Token | 返回 401 |
| Auth 中间件 — 黑名单 Token | 返回 401 |
| Auth 中间件 — 账号已禁用 | 返回 403 |
| Auth 中间件 — 正常 Token | 上下文注入 userID/username/role |
| RequireAdmin — 非 admin 角色 | 返回 403 |
| RequireAdmin — admin 角色 | 放行 |
| AuditMiddleware | 操作记录正确写入审计表 |

**远程访问配置模块（新增功能）**

| 用例 | 验证点 |
|------|--------|
| 创建远程访问配置 | 域名+端口唯一约束 |
| 更新配置 | 字段正确更新 |
| 域名重复 | 返回唯一约束错误 |
| 端口范围校验 | 非法端口返回错误 |
| 启用/禁用切换 | `enabled` 状态正确 |
| Nginx 配置生成 | 生成正确的 server block |

### 3.3 前端单元测试

#### 3.3.1 测试规范

- 框架：Vitest + `@vue/test-utils`
- Mock：使用 Vitest 内置 `vi.mock()` 模拟 API 请求
- 组织：每个组件/composable 对应一个 `.spec.ts` 文件，放在同级 `__tests__/` 目录
- 命名：`describe('组件名')` + `it('应该...')`

#### 3.3.2 必须覆盖的测试用例

**认证相关**

| 用例 | 验证点 |
|------|--------|
| 登录表单校验 | 空用户名/密码提示错误 |
| 登录成功跳转 | admin 跳转管理后台，customer 跳转客户首页 |
| 首次登录改密拦截 | `must_change_password=true` 时跳转改密页 |
| Token 过期自动刷新 | 401 响应触发 refresh，重试原请求 |
| 登出清理 | 清除 localStorage 中的 token，跳转登录页 |

**管理员页面**

| 用例 | 验证点 |
|------|--------|
| 客户列表渲染 | 表格数据正确展示，分页组件工作 |
| 创建客户表单 | 必填字段校验，邮箱格式校验 |
| 机器列表筛选 | 状态/区域筛选参数正确传递 |
| 机器分配对话框 | 客户选择、租期设置、确认提交 |
| 机器回收确认 | 二次确认弹窗，回收原因填写 |

**客户页面**

| 用例 | 验证点 |
|------|--------|
| 我的机器列表 | 只展示已分配的机器 |
| SSH 连接信息展示 | 地址、端口、用户名正确 |
| SSH 密钥管理 | 添加/删除密钥，格式校验 |
| 任务列表 | 状态筛选，停止操作 |
| 改密页面 | 旧密码+新密码+确认密码校验 |

**远程访问组件（新增功能）**

| 用例 | 验证点 |
|------|--------|
| WebTerminal 组件挂载 | xterm.js 实例正确初始化 |
| WebSocket 连接状态 | 状态指示器正确显示 connecting/connected/disconnected |
| 断线重连逻辑 | 指数退避重试，最多 5 次 |
| RemoteDesktop 组件 | Guacamole 客户端正确初始化 |
| 会话列表渲染 | 活跃/历史会话正确展示 |
| 强制断开操作 | 管理员可终止任意会话 |

### 3.4 Agent 单元测试

#### 3.4.1 测试规范

- 框架：Go 标准 `testing`
- 数据库：临时文件 SQLite，`t.Cleanup()` 清理
- 命名：`TestXxx_场景描述`

#### 3.4.2 必须覆盖的测试用例

| 用例 | 验证点 |
|------|--------|
| 命令执行成功 | 返回正确的 exit_code 和 stdout |
| 命令执行超时 | 超时后进程被终止 |
| 优先级队列排序 | 高优先级任务先出队 |
| 并发队列操作 | 多 goroutine 安全 push/pop |
| SQLite 任务持久化 | 写入后可正确读取 |
| 心跳上报 | 正确发送心跳请求 |
| 任务认领 | 轮询领取并更新本地状态 |
| 断线恢复 | 重启后从 SQLite 恢复未完成任务 |

---

## 4. 安全测试方案

### 4.1 认证与授权测试

#### 4.1.1 JWT Token 安全

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| 过期 Token 访问 | 使用过期 JWT 请求 API | 返回 401 |
| 篡改 Token 签名 | 修改 JWT payload 后请求 | 返回 401 |
| 空 Token | 不携带 Authorization 头 | 返回 401 |
| 格式错误 Token | 发送非 Bearer 格式 | 返回 401 |
| 黑名单 Token | 登出后使用原 Token | 返回 401 |
| Refresh Token 轮换 | 使用旧 refresh_token | 返回无效 |
| Token 中角色篡改 | 修改 claims 中 role 字段 | 签名校验失败 |

#### 4.1.2 权限隔离测试

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| Customer 访问 Admin API | 用 customer Token 请求 `/admin/*` | 返回 403 |
| Customer 访问他人机器 | 请求非自己分配的机器详情 | 返回 403 或空 |
| Customer 停止他人任务 | 用自己 Token 停止他人任务 | 返回 `ErrUnauthorized` |
| Customer 查看他人数据集 | 请求非自己的数据集 | 返回 403 |
| 禁用账号访问 | 账号 status=disabled 后请求 | 返回 403 |
| Agent API 无认证 | Agent 端点无 JWT 保护（当前设计） | 需评估是否加认证 |

#### 4.1.3 密码安全测试

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| 弱密码设置 | 改密时使用 "123456" | 应被拒绝（待修复） |
| 密码哈希存储 | 查看数据库中密码字段 | 存储 bcrypt 哈希，非明文 |
| 默认密码强制改密 | 首次登录后访问其他页面 | 被拦截到改密页 |
| 暴力破解防护 | 连续多次错误密码 | 建议增加限流/锁定 |

### 4.2 输入校验与注入防护

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| SQL 注入 — 登录 | 用户名输入 `' OR 1=1 --` | 登录失败，无 SQL 错误泄露 |
| SQL 注入 — 搜索 | 列表搜索参数注入 SQL 片段 | GORM 参数化查询防护 |
| XSS — 客户名称 | 创建客户时 username 含 `<script>` | 前端转义展示，不执行脚本 |
| XSS — 机器备注 | 机器 remark 字段含 HTML 标签 | 前端转义展示 |
| 路径遍历 | 数据集上传路径含 `../../` | 路径规范化，拒绝越界 |
| 超长输入 | 各字段输入超长字符串 | 返回参数校验错误 |
| 特殊字符 | 用户名/密码含特殊字符 | 正确处理，不引发异常 |

### 4.3 远程访问安全测试

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| WebSocket 无 Token 连接 | 不携带 JWT 建立 WS | 连接被拒绝 |
| WebSocket Token 过期 | 连接中 Token 过期 | 会话超时断开 |
| 会话隔离 | 客户 A 尝试恢复客户 B 的会话 | 返回 403 |
| 并发会话限制 | 超过策略允许的最大并发数 | 新连接被拒绝 |
| IP 白名单 | 非白名单 IP 发起连接 | 连接被拒绝 |
| 协议限制 | 策略禁止 VNC 时尝试 VNC 连接 | 返回 403 |
| 会话超时 | 空闲超过配置时间 | 自动断开 |
| 强制断开 | 管理员终止活跃会话 | 客户端收到断开通知 |

### 4.4 审计合规测试

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| 登录审计 | 登录后查询审计日志 | 记录登录事件 |
| 分配审计 | 分配机器后查询 | 记录分配操作和详情 |
| 回收审计 | 回收机器后查询 | 记录回收操作和原因 |
| 远程会话审计 | 创建/断开会话后查询 | 记录会话生命周期事件 |
| 审计日志不可篡改 | 尝试通过 API 修改审计记录 | 无修改/删除接口 |
| 审计日志筛选 | 按时间/类型/用户筛选 | 返回正确结果集 |

### 4.5 数据隔离测试

| 测试项 | 方法 | 预期结果 |
|--------|------|----------|
| 客户 A 查看客户 B 机器 | 用 A 的 Token 请求 B 的机器 | 返回空或 403 |
| 客户 A 查看客户 B 任务 | 用 A 的 Token 请求 B 的任务列表 | 只返回 A 的任务 |
| 客户 A 挂载到 B 的机器 | 数据集挂载指定 B 的机器 | 返回 `ErrUnauthorized` |
| 多租户列表隔离 | 不同客户查询同一列表接口 | 各自只看到自己的数据 |

---

## 5. 性能测试方案

### 5.1 测试工具

| 工具 | 用途 |
|------|------|
| `k6` | HTTP API 压力测试、WebSocket 负载测试 |
| `go test -bench` | Go 函数级基准测试 |
| `pprof` | Go 运行时性能分析（CPU/内存） |
| `pgbench` | PostgreSQL 数据库压力测试 |

### 5.2 API 性能基准

| 接口 | 并发数 | 目标 QPS | 目标 P99 延迟 |
|------|--------|----------|---------------|
| `POST /auth/login` | 50 | ≥200 | ≤500ms |
| `GET /admin/machines` | 100 | ≥500 | ≤200ms |
| `GET /customer/machines` | 100 | ≥500 | ≤200ms |
| `GET /customer/tasks` | 100 | ≥500 | ≤200ms |
| `POST /machines/:id/allocate` | 20 | ≥50 | ≤1000ms |
| `POST /agent/heartbeat` | 200 | ≥1000 | ≤100ms |
| `POST /agent/tasks/claim` | 100 | ≥500 | ≤200ms |

### 5.3 WebSocket 性能测试

| 场景 | 指标 | 目标 |
|------|------|------|
| 并发 WebSocket 连接 | 最大同时连接数 | ≥500 |
| 消息吞吐量 | 每秒消息数（单连接） | ≥1000 msg/s |
| 连接建立延迟 | WebSocket 握手完成时间 | ≤200ms |
| 长连接稳定性 | 持续连接 24h 无异常断开 | 0 次意外断开 |
| 断线重连风暴 | 500 连接同时重连 | 服务不崩溃，逐步恢复 |

### 5.4 数据库性能测试

| 场景 | 方法 | 关注指标 |
|------|------|----------|
| 大表查询 | 机器表 10000+ 条记录分页查询 | 查询延迟 ≤100ms |
| 并发写入 | 100 并发创建任务 | 无死锁，写入成功率 100% |
| 索引有效性 | `EXPLAIN ANALYZE` 关键查询 | 使用索引扫描，非全表扫描 |
| 连接池压力 | 超过连接池上限的并发请求 | 排队等待，不报错 |

### 5.5 Go 基准测试

需要为以下关键函数编写 `Benchmark` 测试：

```go
// 密码哈希性能（bcrypt 较慢，需确认可接受）
func BenchmarkHashPassword(b *testing.B)
func BenchmarkCheckPasswordHash(b *testing.B)

// JWT 生成和解析
func BenchmarkGenerateToken(b *testing.B)
func BenchmarkParseToken(b *testing.B)

// 任务队列操作
func BenchmarkQueuePushPop(b *testing.B)
```

---

## 6. 集成测试方案

### 6.1 后端 API 集成测试

使用 `httptest.NewServer` 启动完整的 Gin 路由，配合 SQLite 内存数据库，测试完整的 HTTP 请求链路。

#### 6.1.1 测试环境搭建模式

沿用现有 `testEnv` 模式，统一封装：

```go
// backend/internal/testutil/setup.go
type TestEnv struct {
    DB     *gorm.DB
    Router *gin.Engine
    Server *httptest.Server
}

func SetupTestEnv(t *testing.T) *TestEnv {
    t.Helper()
    // 1. 初始化 JWT
    // 2. 创建 SQLite 内存数据库
    // 3. AutoMigrate 所有实体
    // 4. 初始化 Service 和 Controller
    // 5. 注册路由
    // 6. 启动 httptest.Server
    t.Cleanup(func() { env.Server.Close() })
    return env
}

func (e *TestEnv) AdminToken(t *testing.T) string { ... }
func (e *TestEnv) CustomerToken(t *testing.T, userID uint) string { ... }
```

#### 6.1.2 关键业务流程集成测试

**流程一：客户入驻全流程**

```
1. POST /admin/customers — 管理员创建客户
2. POST /auth/login — 客户使用默认密码登录
3. 验证返回 must_change_password=true
4. POST /auth/password/change — 客户修改密码
5. POST /auth/login — 使用新密码登录成功
6. GET /auth/profile — 获取个人信息
```

**流程二：机器分配与使用全流程**

```
1. POST /admin/machines — 管理员添加机器
2. POST /admin/machines/:id/allocate — 分配给客户
3. GET /customer/machines — 客户查看已分配机器
4. GET /customer/machines/:id/connection — 获取连接信息
5. POST /customer/tasks/training — 创建训练任务
6. GET /customer/tasks — 查看任务列表
7. POST /customer/tasks/:id/stop — 停止任务
```

**流程三：机器回收全流程**

```
1. POST /admin/machines/:id/reclaim — 管理员回收机器
2. GET /customer/machines — 客户机器列表不再包含该机器
3. GET /admin/audit/logs — 审计日志记录回收操作
```

**流程四：Agent 任务生命周期**

```
1. POST /agent/heartbeat — Agent 上报心跳
2. POST /agent/tasks/claim — Agent 认领任务
3. POST /agent/tasks/:id/start — 标记任务开始
4. POST /agent/tasks/:id/lease/renew — 续约租约
5. POST /agent/tasks/:id/complete — 任务完成
```

**流程五：远程会话全流程（新增功能）**

```
1. POST /customer/remote/sessions — 创建远程会话
2. WebSocket 连接建立 — 验证 Token 认证
3. GET /admin/remote/sessions — 管理员查看活跃会话
4. POST /admin/remote/sessions/:id/terminate — 强制断开
5. GET /admin/audit/logs?resource_type=remote_session — 审计记录
```

### 6.2 前后端联调测试

在开发环境中，前端连接真实后端 API，验证端到端数据流：

| 场景 | 验证点 |
|------|--------|
| 登录流程 | Token 存储、角色路由跳转、改密拦截 |
| Token 自动刷新 | access_token 过期后自动续期，请求无感重试 |
| 列表分页 | 分页参数传递、总数展示、翻页加载 |
| 表单提交 | 校验错误展示、成功后列表刷新 |
| 错误处理 | 后端错误码映射为前端提示信息 |
| 实时数据 | WebSocket 连接状态、消息收发 |

### 6.3 端到端（E2E）测试

#### 6.3.1 工具选型

| 工具 | 用途 |
|------|------|
| Playwright | 浏览器自动化 E2E 测试（推荐，替代现有 vue.spec.ts） |
| Docker Compose | 搭建完整测试环境（后端 + 数据库 + Redis） |

#### 6.3.2 E2E 测试场景

| 场景 | 步骤 | 验证点 |
|------|------|--------|
| 管理员登录 | 打开登录页 → 输入凭据 → 点击登录 | 跳转到管理后台仪表板 |
| 客户首次登录改密 | 登录 → 自动跳转改密页 → 修改密码 → 跳转首页 | 改密后正常使用 |
| 创建客户 | 管理后台 → 客户管理 → 新建 → 填写表单 → 提交 | 列表中出现新客户 |
| 机器分配 | 机器详情 → 分配 → 选择客户 → 确认 | 机器状态变为 allocated |
| 客户查看机器 | 客户登录 → 我的机器 → 查看列表 | 展示已分配的机器 |
| Web 终端连接 | 机器详情 → 点击 SSH 终端 → 等待连接 | 终端界面正常显示 |

---

## 7. Code Review 流程

### 7.1 Review 原则

| 原则 | 说明 |
|------|------|
| 安全优先 | 涉及认证、权限、远程访问的代码必须重点审查 |
| 功能正确 | 业务逻辑是否符合需求文档描述 |
| 代码规范 | 是否遵循项目架构规则和代码风格 |
| 测试覆盖 | 新增功能是否附带对应测试 |
| 性能影响 | 是否引入 N+1 查询、内存泄漏等性能问题 |

### 7.2 Review 检查清单

#### 7.2.1 后端代码检查清单

**架构合规**
- [ ] 是否遵循 Controller → Service → DAO 分层
- [ ] Controller 是否只做参数绑定和响应，不含业务逻辑
- [ ] Service 是否通过 DAO 访问数据，不直接操作 `*gorm.DB`
- [ ] 新增实体是否定义了 `TableName()` 方法
- [ ] 新增路由是否在 `router.go` 中正确注册

**安全检查**
- [ ] Admin API 是否挂载了 `RequireAdmin()` 中间件
- [ ] Customer API 是否挂载了 `Auth()` 中间件
- [ ] 资源访问是否校验了归属关系（customerID 匹配）
- [ ] 敏感操作是否记录审计日志
- [ ] 用户输入是否通过 `binding` tag 校验
- [ ] 密码是否使用 bcrypt 哈希存储
- [ ] 错误信息是否避免泄露内部细节

**数据库与并发**
- [ ] 涉及状态变更的操作是否使用事务
- [ ] 并发敏感操作是否加行级锁（`SELECT ... FOR UPDATE`）
- [ ] 新增表/字段是否附带 SQL 迁移脚本
- [ ] 查询是否使用 `WithContext(ctx)` 支持取消
- [ ] 列表查询是否支持分页，避免全量加载
- [ ] 是否存在 N+1 查询问题

**测试要求**
- [ ] 新增 Service 方法是否有对应单元测试
- [ ] 新增 API 端点是否有 HTTP 集成测试
- [ ] 测试是否覆盖成功路径和主要失败路径
- [ ] 测试是否使用 SQLite 内存数据库，不依赖外部服务

#### 7.2.2 前端代码检查清单

**架构合规**
- [ ] 新增页面是否放在正确的 `views/admin/` 或 `views/customer/` 目录
- [ ] 新增组件是否放在 `components/` 对应子目录
- [ ] API 请求是否通过 `api/` 模块统一管理
- [ ] 类型定义是否放在 `types/` 目录
- [ ] 路由是否在 `router/index.ts` 中正确注册

**安全检查**
- [ ] 路由守卫是否正确配置 `requiresRole`
- [ ] 敏感操作是否有二次确认弹窗
- [ ] 用户输入是否做前端校验（格式、长度、必填）
- [ ] 是否避免在前端存储敏感信息（密码明文等）
- [ ] XSS 防护：动态内容是否使用 `v-text` 而非 `v-html`

**类型安全与代码质量**
- [ ] 新增接口和数据结构是否使用 TypeScript interface 定义
- [ ] 是否通过 `npm run type-check` 无报错
- [ ] 是否复用现有通用组件（DataTable、StatusTag、PageHeader 等）
- [ ] 组件 props 是否有明确的类型定义
- [ ] 是否处理了加载态、空状态、错误态

### 7.3 Review 流程

```
开发者提交 PR
  → CI 自动检查（编译、lint、测试）
  → 至少 1 名 Reviewer 审查
  → 涉及安全/权限模块需 2 名 Reviewer
  → Review 通过 + CI 通过 → 合并到 main
```

**PR 提交规范**：
- PR 标题遵循 `type(scope): 描述` 格式
- PR 描述包含：变更说明、测试方法、影响范围
- 关联 Issue 或需求编号
- 附带测试结果截图（前端 UI 变更时）

---

## 8. 测试环境搭建

### 8.1 本地开发测试环境

#### 8.1.1 后端测试环境

```bash
# 单元测试（使用 SQLite 内存数据库，无外部依赖）
cd backend && go test ./...

# 指定模块测试
cd backend && go test ./internal/service/allocation/...
cd backend && go test ./internal/controller/v1/auth/...

# 带覆盖率
cd backend && go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 基准测试
cd backend && go test -bench=. -benchmem ./pkg/auth/...
```

#### 8.1.2 前端测试环境

```bash
# 单元测试
cd frontend && npm run test:unit

# 类型检查
cd frontend && npm run type-check

# E2E 测试（需先启动后端）
cd frontend && npx playwright test
```

#### 8.1.3 Agent 测试环境

```bash
# 单元测试
cd agent && go test ./...

# 带覆盖率
cd agent && go test -cover ./...
```

### 8.2 集成测试环境（Docker Compose）

用于运行需要真实 PostgreSQL 和 Redis 的集成测试：

```yaml
# docker-compose.test.yml
services:
  postgres-test:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: remotegpu_test
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
    ports:
      - "15432:5432"
    tmpfs:
      - /var/lib/postgresql/data  # 内存存储，测试后自动清理

  redis-test:
    image: redis:7-alpine
    ports:
      - "16379:6379"
```

```bash
# 启动测试环境
docker compose -f docker-compose.test.yml up -d

# 运行集成测试
DATABASE_URL="postgres://test:test@localhost:15432/remotegpu_test" \
REDIS_URL="redis://localhost:16379" \
go test -tags=integration ./...

# 清理
docker compose -f docker-compose.test.yml down
```

### 8.3 CI/CD 集成

#### 8.3.1 CI 流水线阶段

```
PR 提交 → Lint → 编译 → 单元测试 → 集成测试 → 覆盖率报告
```

**阶段一：静态检查**
- 后端：`go vet ./...` + `golangci-lint run`
- 前端：`npm run type-check` + `npm run lint`

**阶段二：编译**
- 后端：`go build ./cmd/...`
- 前端：`npm run build`
- Agent：`go build ./cmd/...`

**阶段三：单元测试**
- 后端：`go test -cover ./...`
- 前端：`npm run test:unit`
- Agent：`go test -cover ./...`

**阶段四：集成测试**
- 启动 PostgreSQL + Redis 容器
- 运行带 `-tags=integration` 的测试

#### 8.3.2 测试覆盖率目标

| 模块 | 当前覆盖率（估计） | 目标覆盖率 |
|------|---------------------|------------|
| 后端 Service 层 | ~20% | ≥60% |
| 后端 Controller 层 | ~30% | ≥50% |
| 后端 pkg 工具包 | ~40% | ≥70% |
| Agent 核心模块 | ~50% | ≥70% |
| 前端组件 | ~5% | ≥30% |

**覆盖率门禁**：CI 中设置覆盖率阈值，新增代码覆盖率低于 50% 时 PR 标记警告。

---

## 9. 已知问题修复验证

基于 `docs/design/implementation-review.md` 中发现的问题，需在修复后编写回归测试：

### 9.1 高严重度问题回归测试

| 问题 | 修复方案 | 回归测试 |
|------|----------|----------|
| 机器分配缺少行级锁 | 添加 `SELECT ... FOR UPDATE` | 并发测试：两个 goroutine 同时分配同一机器，只有一个成功 |
| 回收未检查运行中任务 | 回收前查询 running/assigned 任务 | 有运行中任务时回收返回错误 |
| 任务停止无超时控制 | `context.WithTimeout` 包装 Agent 调用 | Mock 慢响应 Agent，验证超时返回错误 |
| owner/member 权限未区分 | 添加 `RequireOwner()` 中间件 | member 调用 owner 专属 API 返回 403 |

### 9.2 中严重度问题回归测试

| 问题 | 修复方案 | 回归测试 |
|------|----------|----------|
| 机器创建竞态条件 | 依赖 DB UNIQUE 约束 | 并发创建同 IP 机器，只有一个成功 |
| 状态转换无校验 | Service 层添加状态机 | offline → allocated 返回非法转换错误 |
| 改密无强度校验 | 调用 `ValidateStrength()` | 弱密码被拒绝，强密码通过 |
| 默认密码硬编码 | 移到 system_configs 表 | 修改配置后新客户使用新默认密码 |
| 任务停止缺少状态检查 | 添加状态校验 | 停止 completed 任务返回错误 |

---

## 10. 测试实施计划

### 10.1 阶段划分

**阶段一：补齐现有功能测试**

- 补充 `AllocationService` 单元测试（分配、回收、并发）
- 补充 `TaskService` 单元测试（停止、权限校验、状态检查）
- 补充 `CustomerController` 集成测试
- 补充中间件独立测试（Auth、RequireAdmin）
- 搭建 CI 流水线基础框架

**阶段二：新功能测试（随开发同步）**

- 远程访问配置 API 测试
- WebSocket 连接管理测试
- 会话管理 CRUD 测试
- 远程访问策略测试
- 前端远程访问组件测试（WebTerminal、RemoteDesktop）

**阶段三：安全与性能测试**

- 执行安全测试清单（第 4 章）
- API 性能基准测试（k6）
- WebSocket 压力测试
- 数据库查询性能分析

**阶段四：E2E 与回归测试**

- 搭建 Playwright E2E 测试框架
- 编写核心业务流程 E2E 用例
- 已知问题修复后的回归测试
- 建立回归测试套件，纳入 CI

### 10.2 测试优先级矩阵

| 优先级 | 测试类型 | 覆盖模块 |
|--------|----------|----------|
| P0 | 认证鉴权单元测试 | AuthService、JWT、中间件 |
| P0 | 权限隔离测试 | Admin/Customer 路由隔离、数据隔离 |
| P0 | 机器分配并发测试 | AllocationService |
| P1 | 任务管理测试 | TaskService、Agent API |
| P1 | 客户管理测试 | CustomerController |
| P1 | 远程访问安全测试 | WebSocket 认证、会话隔离 |
| P2 | 前端组件测试 | 登录、列表、表单组件 |
| P2 | API 性能基准 | 核心接口 QPS/延迟 |
| P3 | E2E 测试 | 核心业务流程 |

---

## 11. 总结

### 11.1 核心策略

1. **测试金字塔**：60% 单元测试 + 30% 集成测试 + 10% E2E 测试
2. **安全优先**：认证、权限、数据隔离是测试重点
3. **沿用现有模式**：后端使用 testify + SQLite 内存数据库，前端使用 Vitest
4. **CI 门禁**：编译、lint、测试、覆盖率全部通过才允许合并
5. **Code Review 必审**：安全相关代码需 2 名 Reviewer

### 11.2 关键风险与应对

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| 机器分配并发冲突 | 同一机器分配给多个客户 | 行级锁 + 并发测试 |
| 远程会话安全漏洞 | 未授权访问 GPU 机器 | WebSocket 认证 + 会话隔离测试 |
| 数据泄露 | 客户间数据交叉访问 | 多租户隔离测试 + 权限 Review |
| Agent 通信中断 | 任务卡死、心跳丢失 | 租约机制 + 超时测试 + 断线恢复测试 |
| 前端 XSS | 恶意脚本注入 | 输入校验 + 输出转义 + 安全测试 |

### 11.3 交付物清单

| 交付物 | 说明 |
|--------|------|
| 本文档 | 测试策略与 Review 方案 |
| 后端单元测试 | `*_test.go` 文件，覆盖核心 Service 和 Controller |
| 前端单元测试 | `*.spec.ts` 文件，覆盖核心组件和 composable |
| 集成测试套件 | 业务流程集成测试 |
| E2E 测试套件 | Playwright 端到端测试 |
| CI 配置 | 流水线配置文件 |
| 性能测试脚本 | k6 压力测试脚本 |
| 测试报告模板 | 覆盖率报告、安全测试报告模板 |
