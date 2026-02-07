# 第二轮代码质量与安全审查报告

> 审查人：architect
> 日期：2026-02-07
> 范围：本轮所有新增/修改代码

---

## 一、编译状态

| 模块 | 状态 | 说明 |
|------|------|------|
| Backend | ✅ 通过 | `go build ./cmd/...` 无错误 |
| Agent | ✅ 通过 | `go build -o remotegpu-agent ./cmd/` 无错误 |
| Frontend | ✅ 通过 | `vue-tsc --noEmit` 类型检查通过 |

---

## 二、架构合规性检查

### 2.1 后端分层架构 ✅ 全部通过

| 检查项 | 状态 | 说明 |
|--------|------|------|
| Entity 层 | ✅ | Document 实体正确实现 TableName() |
| DAO 层 | ✅ | DocumentDao 命名规范，使用 WithContext |
| Service 层 | ✅ | 通过 DAO 访问数据，不直接操作 db |
| Controller 层 | ✅ | 嵌入 BaseController，使用 Success/Error |
| API 层 | ✅ | 纯数据结构，无逻辑 |
| 路由注册 | ✅ | Import 别名规范（serviceXxx/ctrlXxx） |
| 数据库迁移 | ✅ | 编号递增，表结构完整 |

### 2.2 前端架构 ✅ 全部通过

| 检查项 | 状态 | 说明 |
|--------|------|------|
| 目录结构 | ✅ | 视图文件位于 views/admin/ |
| 组合式 API | ✅ | 使用 script setup lang="ts" |
| 路由注册 | ✅ | 路径、名称、元数据完整 |
| 路由守卫 | ✅ | 新增 403 页面，保存重定向路径 |

---

## 三、安全审查

### 3.1 严重问题（P0）

#### 问题 1：敏感凭据在 API 响应中明文返回

- **文件**：`backend/internal/service/machine/machine_service.go`
- **行号**：258-274（GetConnectionInfo）、314-340（GetMachineDetail）
- **描述**：SSH 密码、VNC 密码、Jupyter Token 在数据库中加密存储，但 API 响应中解密后明文返回
- **影响**：凭据泄露风险
- **建议**：对 GetMachineDetail 中的 VNC 密码也进行解密处理（当前直接返回密文），并确保这些接口仅限授权用户访问

#### 问题 2：Entity JSON Tag 泄露敏感字段

- **文件**：`backend/internal/model/entity/resource.go`
- **行号**：26、28
- **描述**：`JupyterToken` 和 `VNCPassword` 使用了 `json:"jupyter_token"` 和 `json:"vnc_password"` 标签，当 Host 实体直接序列化时会暴露
- **影响**：SSHPassword 和 SSHKey 已正确使用 `json:"-"`，但这两个字段遗漏了
- **建议**：将这两个字段的 JSON 标签改为 `json:"-"`

---

### 3.2 高优先级问题（P1）

#### 问题 3：strconv.ParseUint 错误被忽略

- **文件**：`backend/internal/controller/v1/customer/customer_controller.go:149`
- **描述**：`id, _ := strconv.ParseUint(idStr, 10, 64)` 忽略了解析错误
- **影响**：解析失败时 id 为 0，可能影响错误的记录

#### 问题 4：Dataset Controller 多处相同问题

- **文件**：`backend/internal/controller/v1/dataset/dataset_controller.go`
- **行号**：90、129、182、212
- **描述**：多处 `strconv.ParseUint()` 错误被忽略

#### 问题 5：分页参数解析错误被忽略

- **文件**：多个 Controller
  - `customer/my_machine_controller.go:41-42`
  - `machine/machine_controller.go:47-48`
  - `dataset/dataset_controller.go:44-45`
  - `ops/ops_controller.go:24-25`
- **描述**：`strconv.Atoi` 错误被忽略，但由于 DefaultQuery 提供了默认值，实际风险较低

---

### 3.3 中优先级问题（P2）

#### 问题 6：错误处理不一致

- **文件**：`customer_controller.go`
- **描述**：Disable() 方法忽略 ParseUint 错误，Enable() 方法正确检查，风格不一致

#### 问题 7：错误类型使用字符串比较

- **文件**：`dataset_controller.go:94,133`、`task_controller.go:118`
- **描述**：使用 `err.Error() == "无权限访问该资源"` 进行错误判断，应使用 errors.Is 或自定义错误类型

#### 问题 8：LIKE 查询模式构建

- **文件**：`audit_repo.go:46`
- **描述**：`"%"+params.Username+"%"` 虽然通过 GORM 参数化查询安全，但用户输入中的 `%` 和 `_` 通配符未转义
- **风险**：低，仅影响搜索结果准确性

---

## 四、问题汇总

| 编号 | 严重度 | 类型 | 文件 | 状态 |
|------|--------|------|------|------|
| 1 | P0 | 敏感数据泄露 | machine_service.go | 待修复 |
| 2 | P0 | 敏感数据泄露 | resource.go | 待修复 |
| 3 | P1 | 输入验证 | customer_controller.go | 待修复 |
| 4 | P1 | 输入验证 | dataset_controller.go | 待修复 |
| 5 | P1 | 输入验证 | 多个 Controller | 待修复 |
| 6 | P2 | 错误处理 | customer_controller.go | 待修复 |
| 7 | P2 | 错误处理 | dataset/task controller | 待修复 |
| 8 | P2 | 查询安全 | audit_repo.go | 待修复 |

---

## 五、修复建议优先级

1. **立即修复（P0）**：
   - 将 `JupyterToken` 和 `VNCPassword` 的 JSON 标签改为 `json:"-"`
   - 确保 GetMachineDetail 中 VNC 密码正确解密（与 GetConnectionInfo 一致）

2. **尽快修复（P1）**：
   - 所有 `strconv.ParseUint` / `strconv.Atoi` 调用添加错误检查
   - 解析失败时返回 400 错误

3. **计划修复（P2）**：
   - 定义 `ErrForbidden` 等自定义错误类型，替代字符串比较
   - 统一所有 Controller 的错误处理风格

---

## 六、总体评价

本轮开发的代码在**架构合规性**方面表现优秀，所有新增模块严格遵循分层架构规则。

安全方面发现 2 个 P0 问题（敏感数据泄露）和 3 个 P1 问题（输入验证），建议在合并前修复 P0 问题。P1/P2 问题可在后续迭代中逐步修复。
