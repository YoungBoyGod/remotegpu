# RemoteGPU 第六轮问题清单

> 记录时间：2026-02-07 23:10
> 来源：用户反馈 + 代码审查

---

## 一、机器列表页面问题

### 1.1 设备状态大部分显示离线（与实际不符）

**现象：** 机器列表中大部分机器显示"离线"状态，但实际机器是在线的。

**根因分析：**
- Agent 每 30 秒发送心跳，写入 Redis（TTL=90秒）
- HostStatusSyncer 每 5 分钟才同步一次到数据库（已改为 1 分钟）
- 后端重启后 Redis 缓存丢失，所有机器被标记为 offline
- HeartbeatMonitor 和 HostStatusSyncer 的超时判断不一致

**涉及文件：**
- `backend/internal/service/machine/host_status_syncer.go`
- `backend/internal/service/machine/host_status_cache.go`
- `backend/internal/router/router.go`（同步间隔配置）
- `agent/cmd/main.go`（心跳间隔）

**已做修复：**
- 同步间隔从 5 分钟改为 1 分钟
- 启动时立即执行一次同步

**待确认：**
- Agent 进程是否在 GPU 机器上正常运行
- HeartbeatMonitor 配置的 timeout 值是否与 Redis TTL 一致

---

### 1.2 SSH 连接信息需要展开拓展

**现象：** 机器列表中 SSH 连接信息展示不够详细。

**期望：** 参考截图中的展开式连接信息，包含 SSH 地址、端口、用户名、密码（可复制）。

**涉及文件：**
- `frontend/src/views/admin/MachineListView.vue`
- `frontend/src/views/admin/MachineDetailView.vue`

---

### 1.3 分配操作提示逻辑问题

**现象：** 点击分配 → 选择用户 → 确认分配后，提示错误。

**根因分析：**
- 前端分配按钮条件为 `allocation_status !== 'allocated'`，应改为 `=== 'idle'`
- 维护中的机器前端显示分配按钮，但后端要求状态必须是 `idle`
- 前端 `allocateMachine` 发送了 `host_id` 字段，但后端从 URL 参数获取

**涉及文件：**
- `frontend/src/views/admin/MachineListView.vue`（L504 分配按钮条件）
- `frontend/src/api/admin.ts`（L97-99 allocateMachine）
- `backend/internal/controller/v1/machine/machine_controller.go`（L275-298）
- `backend/internal/service/allocation/allocation_service.go`（L274-276）

---

### 1.4 维护模式后状态未改变

**现象：** 设置维护后，设备状态和分配状态都没有改变。

**根因分析：**
- 后端 `SetMaintenance` 只更新 `allocation_status`，不更新 `device_status`
- 前端刷新列表后可能因为缓存问题没有获取到最新状态
- 维护模式应同时影响 `device_status`（显示为"维护"）和 `allocation_status`（显示为"维护"）

**涉及文件：**
- `backend/internal/controller/v1/machine/machine_controller.go`（L499-521）
- `frontend/src/views/admin/MachineListView.vue`（L185-194）

---

### 1.5 删除机器提示资源错误

**现象：** 点击删除后提示资源错误，前后端没有打通。

**根因分析：**
- 后端删除时没有检查机器是否已分配
- 可能存在外键约束导致删除失败（allocations 表引用了 host_id）
- 前端错误处理不够友好

**涉及文件：**
- `backend/internal/controller/v1/machine/machine_controller.go`（L476-485）
- `backend/internal/service/machine/machine_service.go`
- `frontend/src/views/admin/MachineListView.vue`（L164-183）

---

## 二、添加机器页面问题

### 2.1 Nginx 配置字段格式

**现象：** Nginx 配置应该是长文本形式（textarea），支持独立配置。

**涉及文件：**
- `frontend/src/views/admin/MachineAddView.vue`

---

### 2.2 映射端口字段类型

**现象：** SSH 映射端口、PPyTi 映射端口等应该是数字字段，且应有默认值。

**涉及文件：**
- `frontend/src/views/admin/MachineAddView.vue`

---

### 2.3 批量导入提示错误

**现象：** 批量导入功能报错。

**根因分析：**
- 前端发送 `price_hourly`，后端没有处理
- 后端没有处理 `gpu_count`
- GPU 型号被错误存储在 `CPUInfo` 字段中
- 导入时没有设置 `NeedsCollect = true`

**涉及文件：**
- `frontend/src/views/admin/MachineListView.vue`（L42-78 CSV 导入）
- `backend/internal/controller/v1/machine/machine_controller.go`（L339-371）

---

## 三、分配记录页面问题

### 3.1 分配记录提示错误

**现象：** 分配记录页面加载时报错。

**根因分析：**
- 前端期望 allocation 对象包含 `customer` 和 `host` 的完整对象
- 后端可能只返回 ID，没有做 JOIN 查询
- 前端 AllocationListView.vue 中访问 `row.customer?.company` 等嵌套字段

**涉及文件：**
- `frontend/src/views/admin/AllocationListView.vue`
- `backend/internal/dao/allocation_repo.go`
- `backend/internal/controller/v1/machine/machine_controller.go`

---

## 四、前后端数据交互通用问题

### 4.1 Machine 类型定义混乱

**问题：**
- `frontend/src/types/machine.ts` 中字段重复（gpuModel 等与 gpus 数组重复）
- snake_case 和 camelCase 混用
- 存在不存在的字段（allocatedTo 等）

### 4.2 API 过滤参数不一致

**问题：**
- 前端支持 `keyword` 过滤，后端没有实现
- 后端支持 `device_status`/`allocation_status` 过滤，前端没有使用

### 4.3 批量操作 API 缺失

**问题：**
- 后端实现了 BatchSetMaintenance、BatchAllocate、BatchReclaim
- 前端 admin.ts 中没有对应的 API 函数

### 4.4 三个状态字段混淆

**问题：**
- `Status`（兼容旧字段）、`DeviceStatus`、`AllocationStatus` 三个字段含义不清
- 前端在不同地方使用不同字段判断状态

---

## 五、优先级排序

| 优先级 | 问题 | 影响 |
|--------|------|------|
| P0 | 设备状态显示离线 | 核心功能不可用 |
| P0 | 分配操作报错 | 核心功能不可用 |
| P0 | 删除操作报错 | 核心功能不可用 |
| P0 | 维护模式不生效 | 核心功能不可用 |
| P1 | 批量导入报错 | 重要功能不可用 |
| P1 | 分配记录报错 | 重要功能不可用 |
| P1 | SSH 连接信息展示 | 用户体验 |
| P2 | 添加机器表单优化 | 用户体验 |
| P2 | Machine 类型定义清理 | 代码质量 |
| P2 | 批量操作 API 补充 | 功能完善 |

---

## 六、API 路径不匹配问题（第二轮检查发现）

### 6.1 密码重置 API 路径错误
- 前端: `/auth/password-reset/request` 和 `/auth/password-reset/confirm`
- 后端: `/auth/password/request` 和 `/auth/password/confirm`
- **路径不一致**

### 6.2 后端缺失的 API 端点
- `PUT /auth/profile` — 更新个人资料（前端已调用，后端未注册）
- `GET /auth/validate` — Token 验证（前端已调用，后端未注册）
- `GET /admin/allocations` — 分配记录列表（前端已调用，后端只有 /admin/allocations/recent）
- `POST /admin/alerts/batch-acknowledge` — 批量确认告警
- `POST /admin/alert-rules/:id/toggle` — 切换告警规则

### 6.3 Workspace/Environment API 路径缺少前缀
- 前端 `frontend/src/api/workspace/index.ts` 使用 `/workspaces`
- 后端注册的是 `/customer/workspaces`
- Environment 同理

### 6.4 Host API 模块路径错误
- `frontend/src/api/host/index.ts` 使用 `/admin/hosts`
- 后端使用 `/admin/machines`

### 6.5 客户列表字段名不一致
- 前端使用 `customer.contactPerson`（camelCase）
- 后端返回 `contact_person`（snake_case）

### 6.6 客户端仪表板 API 聚合缺少错误处理
- 使用 `Promise.all` 而非 `Promise.allSettled`，任一 API 失败导致整体失败

---

## 七、补充优先级

| 优先级 | 问题 | 影响 |
|--------|------|------|
| P0 | 密码重置 API 路径错误 | 密码重置功能不可用 |
| P0 | Workspace/Environment API 路径缺前缀 | 工作空间功能不可用 |
| P0 | 分配记录 API 端点缺失 | 分配记录页面不可用 |
| P1 | 更新个人资料 API 缺失 | 个人信息页面不可用 |
| P1 | Host API 路径错误 | 部分页面不可用 |
| P2 | 客户字段名不一致 | 数据显示错误 |
| P2 | 仪表板错误处理 | 页面可能白屏 |

---

## 八、管理员后台页面逐页审查（第三轮）

> 审查时间：2026-02-07
> 审查范围：frontend/src/views/admin/ 下所有 10 个页面与后端 API 的数据流对齐

### 8.1 MonitoringView.vue — 监控中心：硬编码假数据

**严重程度：P1**

**问题：** `MonitoringView.vue:50-59` 中设备列表的 CPU/内存/GPU/磁盘使用率和温度全部是 `Math.random()` 生成的假数据。

```javascript
// MonitoringView.vue:52-58
cpuUsage: Math.floor(Math.random() * 60) + 10,
memoryUsage: Math.floor(Math.random() * 70) + 20,
gpuUsage: Math.floor(Math.random() * 80) + 5,
diskUsage: 45,
temperature: 60 + Math.floor(Math.random() * 20),
uptime: 'Running'
```

**影响：** 监控中心显示的设备级指标完全不真实，管理员无法据此判断机器健康状况。

**建议：**
- 后端应提供 `/admin/monitoring/devices` 接口，返回每台机器的实时指标
- 或在 `/admin/machines` 返回数据中嵌入最近一次采集的指标快照

### 8.2 MonitoringView.vue — 状态字段使用旧的 `status` 而非拆分后的双状态

**严重程度：P1**

**问题：** `MonitoringView.vue:169` 使用 `row.status` 显示状态，但后端已拆分为 `device_status` + `allocation_status`。`getStatusType` 函数映射的是 `idle/allocated/maintenance/offline`，这些是 `allocation_status` 的值，不是 `status` 的值。

**涉及文件：**
- `frontend/src/views/admin/MonitoringView.vue:77-85, 169`

**建议：** 改为使用 `row.device_status` 和 `row.allocation_status` 分别显示。

### 8.3 ImageListView.vue — 删除镜像 API 后端未注册路由

**严重程度：P0**

**问题：** `ImageListView.vue:127` 调用 `request.delete('/admin/images/${row.id}')`，但后端路由 `router.go:277-278` 只注册了：
- `GET /admin/images` — 列表
- `POST /admin/images/sync` — 同步

**缺少 `DELETE /admin/images/:id` 路由**，删除操作会返回 404。

**涉及文件：**
- `frontend/src/views/admin/ImageListView.vue:120-133`
- `backend/internal/router/router.go:277-278`
- `backend/internal/controller/v1/ops/image_controller.go`（缺少 Delete 方法）

### 8.4 AlertListView.vue — 告警规则 toggle API 后端未注册路由

**严重程度：P0**

**问题：** `AlertListView.vue:239` 调用 `toggleAlertRule(rule.id, !rule.enabled)`，对应前端 API `admin.ts:353-355` 发送 `POST /admin/alert-rules/:id/toggle`。

但后端路由 `router.go:262-266` 中告警规则只注册了：
- `GET /admin/alert-rules` — 列表
- `GET /admin/alert-rules/:id` — 详情
- `POST /admin/alert-rules` — 创建
- `PUT /admin/alert-rules/:id` — 更新
- `DELETE /admin/alert-rules/:id` — 删除

**缺少 `POST /admin/alert-rules/:id/toggle` 路由**。

**临时方案：** 前端可改用 `updateAlertRule(id, { enabled: !rule.enabled })` 调用 PUT 接口实现同等效果。

**涉及文件：**
- `frontend/src/api/admin.ts:353-355`
- `backend/internal/router/router.go:262-266`

### 8.5 AlertListView.vue — 批量确认告警 API 后端未注册路由

**严重程度：P1**

**问题：** 前端 `admin.ts:291-293` 定义了 `batchAcknowledgeAlerts` 函数，发送 `POST /admin/alerts/batch-acknowledge`，但后端路由中没有注册此端点。

虽然当前 AlertListView.vue 页面未直接调用此函数，但 API 层已定义，后续使用时会 404。

**涉及文件：**
- `frontend/src/api/admin.ts:291-293`
- `backend/internal/router/router.go:258-259`

### 8.6 AlertListView.vue — 告警规则 `duration_seconds` 与后端 `duration` 字段名不一致

**严重程度：P2**

**问题：** 前端 `AlertRuleForm` 使用 `duration_seconds` 字段，但后端 `CreateAlertRuleRequest` 使用 `duration` 字段（单位秒）。

- 前端 `admin.ts:316`: `duration_seconds: number`
- 后端 `ops_controller.go:258`: `rule.Duration = req.Duration`

前端发送 `{ duration_seconds: 300 }` 时，后端无法识别该字段，`duration` 会使用默认值 60。

**涉及文件：**
- `frontend/src/api/admin.ts:304,316`
- `backend/internal/controller/v1/ops/ops_controller.go:246-261`

### 8.7 CustomerListView.vue — 客户列表分页逻辑错误

**严重程度：P2**

**问题：** `CustomerListView.vue` 同时使用了服务端分页和客户端过滤，导致分页数据不准确。

- `loadCustomers()` 调用后端 API 获取分页数据（`response.data.list` + `response.data.total`）
- 但 `filteredCustomers` 在客户端再次过滤，DataTable 的 `:total` 绑定的是 `filteredCustomers.length` 而非服务端 `total`
- 当服务端返回第 1 页 10 条数据，客户端过滤后可能只剩 3 条，分页器显示总数为 3，无法翻页查看更多匹配项

**涉及文件：**
- `frontend/src/views/admin/CustomerListView.vue:26-41, 202-210`

### 8.8 CustomerListView.vue — 启用按钮条件与后端状态不匹配

**严重程度：P2**

**问题：** `CustomerListView.vue:245` 启用按钮条件为 `row.status === 'suspended'`，但禁用操作调用的是 `disableCustomer`（后端将状态设为 `disabled` 而非 `suspended`）。

- 禁用后状态变为 `disabled`，但启用按钮只在 `suspended` 时显示
- 导致禁用后的客户无法通过 UI 重新启用

**涉及文件：**
- `frontend/src/views/admin/CustomerListView.vue:242-246`

### 8.9 CustomerListView.vue — "添加客户"按钮路由未实现

**严重程度：P1**

**问题：** `CustomerListView.vue:178` 有"添加客户"按钮，点击跳转 `/admin/customers/add`，但 `frontend/src/views/admin/` 下没有 `CustomerAddView.vue` 文件，该路由可能未注册或指向空页面。

**涉及文件：**
- `frontend/src/views/admin/CustomerListView.vue:178`
- `frontend/src/router/index.ts`

### 8.10 ImageListView.vue — 直接使用 `request` 而非封装的 API 函数

**严重程度：P2**

**问题：** `ImageListView.vue:6` 导入了 `request from '@/utils/request'`，在 `handleDelete` 中直接调用 `request.delete('/admin/images/${row.id}')`，绕过了 `admin.ts` 的 API 封装层。

这违反了项目的分层架构规范，且 `admin.ts` 中没有定义 `deleteImage` 函数。

**涉及文件：**
- `frontend/src/views/admin/ImageListView.vue:6, 127`
- `frontend/src/api/admin.ts`（缺少 deleteImage 函数）

### 8.11 MachineDetailView.vue — 使用 `(machine as any)` 强制类型断言

**严重程度：P2**

**问题：** `MachineDetailView.vue` 多处使用 `(machine as any).device_status` 和 `(machine as any).allocation_status`（行 150, 153, 175, 182），说明 `Machine` 类型定义中缺少 `device_status` 和 `allocation_status` 字段。

**涉及文件：**
- `frontend/src/views/admin/MachineDetailView.vue:150,153,175,182`
- `frontend/src/types/machine.ts`（需补充字段定义）

### 8.12 TaskListView.vue — 任务统计基于当前页数据而非全局

**严重程度：P2**

**问题：** `TaskListView.vue:393-402` 的 `taskStats` 计算基于 `tasks.value`（当前页数据），而非全局统计。例如 `running` 计数只统计当前页的运行中任务，不是全部。`total` 使用了服务端的 `total.value`，但其他状态计数只基于当前页。

**建议：** 后端应提供 `/admin/tasks/stats` 接口返回全局任务统计。

**涉及文件：**
- `frontend/src/views/admin/TaskListView.vue:393-402`

### 8.13 MonitoringView.vue — 设备表格缺少 `gpu_model` 字段

**严重程度：P2**

**问题：** `MonitoringView.vue:204` 使用 `prop="gpu_model"` 显示 GPU 型号，但后端 `Host` 实体没有 `gpu_model` 顶层字段，GPU 信息在 `gpus` 数组中。该列始终显示为空。

**涉及文件：**
- `frontend/src/views/admin/MonitoringView.vue:204`

### 8.14 AllocationListView.vue — 分配记录列表 API 端点后端缺失（已在 6.2 记录）

**严重程度：P0**

**问题：** 前端调用 `GET /admin/allocations`（带分页参数），但后端路由只注册了 `GET /admin/allocations/recent`（仪表板用，无分页）。分配记录列表页面无法正常加载数据。

**补充说明：** 此问题在 6.2 节已记录，但需要强调这是分配记录页面完全不可用的根本原因。后端需要新增一个支持分页、筛选的 `/admin/allocations` 端点。

---

### 逐页审查结论

| 页面 | 状态 | 发现问题 |
|------|------|----------|
| DashboardView.vue | ✅ 基本正常 | API 对齐，无假数据，错误处理完善 |
| MachineListView.vue | ⚠️ 已有修复中 | 状态显示、分配按钮等问题已在前几轮记录 |
| MachineDetailView.vue | ⚠️ 类型问题 | `(machine as any)` 强转，缺少 device_status/allocation_status 类型定义 |
| MachineAddView.vue | ✅ 基本正常 | 表单字段与后端 CreateMachineRequest 对齐 |
| CustomerListView.vue | ⚠️ 多个问题 | 分页逻辑错误、启用按钮条件不匹配、添加客户页面缺失 |
| AllocationListView.vue | ❌ 不可用 | 后端缺少 GET /admin/allocations 分页端点 |
| TaskListView.vue | ⚠️ 统计不准 | 任务统计基于当前页而非全局 |
| ImageListView.vue | ❌ 删除不可用 | 后端缺少 DELETE /admin/images/:id 路由；绕过 API 层 |
| MonitoringView.vue | ❌ 假数据 | CPU/内存/GPU/磁盘指标全部 Math.random()；状态字段使用旧字段 |
| AlertListView.vue | ⚠️ toggle 缺失 | 后端缺少 toggle 路由；duration_seconds 字段名不一致 |

---

### 第三轮补充优先级

| 优先级 | 问题编号 | 问题 | 影响 |
|--------|----------|------|------|
| P0 | 8.3 | 镜像删除 API 后端未注册 | 删除功能 404 |
| P0 | 8.4 | 告警规则 toggle API 后端未注册 | 启用/禁用规则 404 |
| P0 | 8.14 | 分配记录列表 API 缺失（重复） | 页面完全不可用 |
| P1 | 8.1 | 监控中心设备指标全部假数据 | 监控数据不可信 |
| P1 | 8.2 | 监控中心使用旧 status 字段 | 状态显示错误 |
| P1 | 8.5 | 批量确认告警 API 未注册 | 功能预留但不可用 |
| P1 | 8.8 | 客户启用按钮条件不匹配 | 禁用后无法重新启用 |
| P1 | 8.9 | 添加客户页面缺失 | 按钮点击后空白 |
| P2 | 8.6 | 告警规则 duration 字段名不一致 | 持续时间始终为默认值 |
| P2 | 8.7 | 客户列表分页逻辑错误 | 分页数据不准确 |
| P2 | 8.10 | 镜像删除绕过 API 封装层 | 代码规范问题 |
| P2 | 8.11 | MachineDetail 使用 as any 强转 | 类型安全问题 |
| P2 | 8.12 | 任务统计基于当前页 | 统计数据误导 |
| P2 | 8.13 | 监控表格 gpu_model 字段不存在 | 列始终为空 |
