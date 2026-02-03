# RemoteGPU 前后端对齐改造清单（后端）

**版本**: v1.0
**日期**: 2026-02-02
**作者**: RemoteGPU开发团队

---

## 一、现状概览

- 后端主要路由集中在 `/api/v1`，前端要求 `/api/admin` 与 `/api/customer`。
- 已实现模块：用户、主机、GPU、环境、配额基础接口。
- 未覆盖模块：客户管理、分配管理、监控告警、镜像、数据集 API 层、任务、统计、系统设置。

涉及主要代码位置：
- `backend/internal/router/router.go`
- `backend/internal/controller/v1/*`
- `backend/internal/service/*`
- `backend/pkg/response/response.go`

---

## 二、接口对齐表（前端接口 → 后端实现/缺口）

### 2.1 认证与用户

| 前端接口 | 后端实现/缺口 | 备注 |
| --- | --- | --- |
| `POST /api/auth/login` | `POST /api/v1/user/login` 已有 | 需接入 AD/LDAP/SSO 并新增路由别名 |
| `POST /api/auth/refresh` | 缺口 | 新增 Refresh Token 与撤销机制 |
| `POST /api/auth/logout` | 缺口 | 需支持 Token 失效策略 |
| `GET /api/auth/profile` | `GET /api/v1/user/info` 已有 | 需路径与返回结构对齐 |
| `POST /api/auth/password-reset/request` | 缺口 | 找回密码流程 |
| `POST /api/auth/password-reset/confirm` | 缺口 | 找回密码流程 |

### 2.2 管理员端

| 前端接口 | 后端实现/缺口 | 备注 |
| --- | --- | --- |
| `/api/admin/dashboard/*` | 缺口 | 统计聚合与趋势接口 |
| `/api/admin/machines*` | 缺口 | 需聚合 `hosts` + `gpus` 形成“机器”视图 |
| `/api/admin/machines/:id/allocate` | 缺口 | 目前仅 `gpus/:id/allocate` |
| `/api/admin/machines/:id/reclaim` | 缺口 | 需回收与释放流程 |
| `/api/admin/machines/import*` | 缺口 | Excel 导入/模板下载 |
| `/api/admin/customers*` | 缺口 | 客户管理与禁用/启用 |
| `/api/admin/allocations*` | 缺口 | 分配记录、续期、回收 |
| `/api/admin/monitoring/realtime` | 缺口 | 实时监控数据聚合 |
| `/ws/admin/monitoring` | 缺口 | WebSocket 推送 |
| `/api/admin/alerts*` | 缺口 | 告警列表/处理/规则 |
| `/api/admin/images*` | 缺口 | 镜像库、审核、上传 |
| `/api/admin/datasets*` | 缺口 | 仅有服务层，无 API 层 |
| `/api/admin/storage/*` | 部分缺口 | 已有配额 CRUD，但缺存储概览/分布 |
| `/api/admin/tasks*` | 缺口 | 任务管理与日志 |
| `/api/admin/statistics/*` | 缺口 | 资源/客户统计 |
| `/api/admin/settings` | 缺口 | 平台配置 |
| `/api/admin/audits` | 缺口 | 审计日志 |

### 2.3 客户端

| 前端接口 | 后端实现/缺口 | 备注 |
| --- | --- | --- |
| `/api/customer/dashboard/*` | 缺口 | 工作台统计 |
| `/api/customer/machines*` | 部分缺口 | 可复用 `environments` 但需语义对齐 |
| `/api/customer/machines/:id/monitoring` | 缺口 | 监控数据接口 |
| `/api/customer/machines/:id/processes` | 缺口 | 进程列表 |
| `/api/customer/machines/:id/tasks` | 缺口 | 机器任务历史 |
| `/api/customer/tasks*` | 缺口 | 训练/推理任务全套接口 |
| `/ws/customer/tasks/:id/logs` | 缺口 | 任务日志推送 |
| `/api/customer/images*` | 缺口 | 镜像市场与我的镜像 |
| `/api/customer/datasets*` | 缺口 | 仅有服务层，无 API 层 |
| `/api/customer/storage/quota` | `GET /api/v1/quotas/usage` | 需路径与响应对齐 |
| `/api/customer/settings/*` | 缺口 | 个人设置/SSH Key/通知 |

---

## 三、模块改造任务清单（按模块拆分）

### 3.1 基础设施与规范（P0）
- 新增 `/api/admin` 与 `/api/customer` 路由分组，并保持 `/api/v1` 兼容层。
- 统一响应结构为 `{code, message, data, traceId}`。
- 统一分页/排序/过滤参数解析。

### 3.2 认证与权限（P0）
- 接入 AD/LDAP/SSO，停用注册流程。
- 增加 Refresh Token 与注销/撤销机制。
- 引入租户/角色模型（admin/owner/member），实现数据隔离中间件。

### 3.3 管理员端模块（P0/P1）
- **Dashboard（P0）**: 资源统计、GPU趋势、活跃客户与告警汇总。
- **机器管理（P0）**: 聚合主机+GPU，支持分配/回收/导入。
- **客户管理（P0）**: 客户 CRUD、禁用/启用、详情统计。
- **分配管理（P0）**: allocation 表与分配流程、续期/回收。
- **监控与告警（P1）**: 实时监控 API + WebSocket、告警规则与处理。
- **镜像管理（P1）**: 镜像上传、审核、下架。
- **数据集与存储（P1）**: 数据集 CRUD、公开/私有、配额与清理建议。
- **任务管理（P1）**: 任务列表、详情、日志、停止。
- **统计与审计（P1）**: 资源统计、客户统计、审计日志。

### 3.4 客户端模块（P0/P1）
- **工作台（P0）**: 统计卡片、机器概览、近期任务。
- **我的机器（P0）**: 机器详情、监控、进程、任务历史。
- **任务管理（P1）**: 训练/推理任务创建与控制、日志/指标。
- **镜像市场（P1）**: 市场列表、部署、我的镜像管理。
- **数据集管理（P1）**: 上传/分享/挂载/详情文件管理。
- **个人设置（P1）**: 资料、SSH Key、通知设置。

---

## 四、数据模型对齐（P0）

- 统一 `customers`/`users` 命名与实体映射，明确 `customer_id` 与 `user_id` 字段规范。
- 补齐 allocation、任务、告警、审计等表与 DAO/Service。
- 对齐数据集、环境等字段与状态枚举（`storage_path/visibility` 等）。

---

## 五、模块接口对齐 → 实现建议

### 5.1 认证与账号

| 接口范围 | 建议实现 | 依赖层 | 备注 |
| --- | --- | --- | --- |
| `/api/auth/login` | `AuthController.Login` | `AuthService` + SSO/LDAP Adapter | 替换本地注册登录 |
| `/api/auth/refresh` | `AuthController.Refresh` | `TokenService` | 引入 Refresh Token 存储/撤销 |
| `/api/auth/logout` | `AuthController.Logout` | `TokenService` | Token 黑名单或版本号策略 |
| `/api/auth/password-reset/*` | `PasswordController` | `UserService` | 若统一认证承接则改为跳转 |

### 5.2 管理员端

| 模块 | 接口范围 | 建议实现 | 依赖层 | 备注 |
| --- | --- | --- | --- | --- |
| Dashboard | `/api/admin/dashboard/*` | `AdminDashboardController` | `DashboardService` | 统计聚合与趋势数据 |
| 机器管理 | `/api/admin/machines*` | `AdminMachineController` | `MachineService` | 聚合 `hosts`+`gpus` |
| 客户管理 | `/api/admin/customers*` | `AdminCustomerController` | `CustomerService` | 客户 CRUD + 团队成员 |
| 分配管理 | `/api/admin/allocations*` | `AdminAllocationController` | `AllocationService` | 整机分配/续期/回收 |
| 监控告警 | `/api/admin/monitoring/*` `/api/admin/alerts*` | `AdminMonitoringController` | `MonitoringService` | WebSocket + 告警规则 |
| 镜像管理 | `/api/admin/images*` | `AdminImageController` | `ImageService` | 审核流转与上传 |
| 数据集存储 | `/api/admin/datasets*` `/api/admin/storage/*` | `AdminDatasetController` | `DatasetService` | 公开/私有、配额 |
| 任务管理 | `/api/admin/tasks*` | `AdminTaskController` | `TaskService` | 训练/推理任务查询 |
| 统计审计 | `/api/admin/statistics/*` `/api/admin/audits` | `AdminStatisticsController` | `StatisticsService` | 资源/客户统计 |
| 系统设置 | `/api/admin/settings` | `AdminSettingsController` | `SettingsService` | 全局配置 |

### 5.3 客户端

| 模块 | 接口范围 | 建议实现 | 依赖层 | 备注 |
| --- | --- | --- | --- | --- |
| 工作台 | `/api/customer/dashboard/*` | `CustomerDashboardController` | `DashboardService` | 个人统计聚合 |
| 我的机器 | `/api/customer/machines*` | `CustomerMachineController` | `EnvironmentService` | 环境语义对齐机器 |
| 任务管理 | `/api/customer/tasks*` | `CustomerTaskController` | `TaskService` | 训练/推理创建 + 日志 |
| 镜像市场 | `/api/customer/images*` | `CustomerImageController` | `ImageService` | 市场/我的镜像 |
| 数据集管理 | `/api/customer/datasets*` | `CustomerDatasetController` | `DatasetService` | 上传/分享/挂载 |
| 个人设置 | `/api/customer/settings/*` | `CustomerSettingsController` | `UserService` | 资料/SSH Key/通知 |

### 5.4 共性能力

| 能力 | 建议实现 | 依赖层 | 备注 |
| --- | --- | --- | --- |
| WebSocket 推送 | `WS Hub` | `MonitoringService`/`TaskService` | 监控与日志流 |
| 审计日志 | `AuditService` | `AuditLogDao` | 管理员/关键操作 |
| 文件上传 | `UploadService` | `StorageAdapter` | 分片/断点续传 |

---

## 六、里程碑与阶段计划（建议）

| 阶段 | 时间范围 | 目标 | 交付物 |
| --- | --- | --- | --- |
| Phase 0 | 1-2 周 | 基础对齐 | 路由前缀统一、响应结构、SSO 登录、租户隔离 |
| Phase 1 | 2-3 周 | MVP 核心能力 | 管理员 Dashboard/机器/客户/分配，客户工作台/我的机器 |
| Phase 2 | 3-4 周 | 增值功能 | 监控告警、镜像、数据集、任务管理 |
| Phase 3 | 2-3 周 | 完善与运营 | 统计、审计、设置、通知、日志推送 |

---

## 七、API 示例（关键路径）

> 示例用于前后端对齐，实际字段以最终数据模型为准。

### 7.1 统一响应结构
```
{
  "code": 0,
  "message": "success",
  "data": {},
  "traceId": "d3b1..."
}
```

### 7.2 认证

**登录** `POST /api/auth/login`
```
{
  "username": "alice",
  "password": "******"
}
```
```
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "...",
    "refreshToken": "...",
    "user": { "id": 1, "name": "Alice", "role": "owner", "tenantId": 1001 }
  },
  "traceId": "..."
}
```

**刷新 Token** `POST /api/auth/refresh`
```
{ "refreshToken": "..." }
```

### 7.3 管理员 - 机器管理

**机器列表** `GET /api/admin/machines?page=1&pageSize=20`
```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "host-001",
        "name": "GPU-Server-01",
        "status": "online",
        "gpuModel": "RTX 4090",
        "gpuCount": 4,
        "cpu": 32,
        "memory": 128,
        "allocationStatus": "allocated"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  },
  "traceId": "..."
}
```

**分配机器** `POST /api/admin/machines/:id/allocate`
```
{ "customerId": 2001, "durationMonths": 3, "remark": "试用" }
```

### 7.4 管理员 - 客户与分配

**客户列表** `GET /api/admin/customers?page=1&pageSize=20`
```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      { "id": 2001, "name": "Acme", "status": "active", "machines": 2 }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  },
  "traceId": "..."
}
```

**分配记录** `GET /api/admin/allocations?page=1&pageSize=20`
```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      { "id": 9001, "machineId": "host-001", "customerId": 2001, "status": "active" }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  },
  "traceId": "..."
}
```

### 7.5 管理员 - 监控与告警

**实时监控** `GET /api/admin/monitoring/realtime`
```
{
  "code": 0,
  "message": "success",
  "data": {
    "machines": [
      { "id": "host-001", "gpuUsage": 70, "cpuUsage": 30, "temperature": 65 }
    ]
  },
  "traceId": "..."
}
```

### 7.6 管理员 - 数据集

**数据集列表** `GET /api/admin/datasets?page=1&pageSize=20`
```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      { "id": 3001, "name": "ImageNet", "visibility": "public", "status": "available" }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  },
  "traceId": "..."
}
```

### 7.7 客户端 - 工作台与机器

**工作台统计** `GET /api/customer/dashboard/stats`
```
{
  "code": 0,
  "message": "success",
  "data": { "machineCount": 2, "runningTasks": 1, "storageUsed": 500 },
  "traceId": "..."
}
```

**我的机器列表** `GET /api/customer/machines?page=1&pageSize=20`
```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      { "id": "host-001", "gpuUsage": 55, "sshHost": "ssh.xxx.com" }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  },
  "traceId": "..."
}
```

### 7.8 客户端 - 任务

**创建训练任务** `POST /api/customer/tasks/training`
```
{
  "name": "train-resnet",
  "machineId": "host-001",
  "gpu": 2,
  "image": "pytorch:2.0-cuda11",
  "command": "python train.py"
}
```

### 7.9 客户端 - 数据集

**上传数据集** `POST /api/customer/datasets/upload`
```
{
  "name": "my-dataset",
  "type": "image",
  "visibility": "private"
}
```

---

**文档结束**
