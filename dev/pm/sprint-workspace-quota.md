# Sprint: Workspace & ResourceQuota 模块开发

**项目**: RemoteGPU 企业级GPU云平台
**Sprint周期**: 2026-01-30 开始
**PM**: 项目经理
**状态**: 🟡 进行中

---

## 📋 Sprint 目标

完成 Workspace 和 ResourceQuota 两个核心模块的 Controller 层和 API 路由开发，实现前后端完整对接。

### 业务价值
- **Workspace**: 支持多租户团队协作，提升企业客户价值
- **ResourceQuota**: 实现资源配额管理，防止资源滥用，保障平台稳定性

---

## 🎯 模块一：Workspace 工作空间管理

### 后端开发任务

**负责人**: 后端开发
**状态**: ✅ 已完成
**Service层**: ✅ 已完成（`internal/service/workspace.go`）

#### API 接口定义

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | `/api/v1/workspaces` | 创建工作空间 | 认证用户 |
| GET | `/api/v1/workspaces` | 列出工作空间（分页） | 认证用户 |
| GET | `/api/v1/workspaces/:id` | 获取工作空间详情 | 工作空间成员 |
| PUT | `/api/v1/workspaces/:id` | 更新工作空间 | 工作空间所有者 |
| DELETE | `/api/v1/workspaces/:id` | 删除工作空间 | 工作空间所有者 |
| POST | `/api/v1/workspaces/:id/members` | 添加成员 | 工作空间所有者 |
| DELETE | `/api/v1/workspaces/:id/members/:user_id` | 移除成员 | 工作空间所有者 |
| GET | `/api/v1/workspaces/:id/members` | 列出成员 | 工作空间成员 |

#### 请求/响应格式

**1. 创建工作空间**
```json
// POST /api/v1/workspaces
// Request
{
  "name": "AI研发团队",
  "description": "AI模型训练工作空间",
  "resource_quota": {
    "max_gpu": 4,
    "max_cpu": 16,
    "max_memory": 65536
  }
}

// Response 200
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "AI研发团队",
    "description": "AI模型训练工作空间",
    "owner_id": 1,
    "created_at": "2026-01-30T10:00:00Z"
  }
}
```

**2. 列出工作空间**
```json
// GET /api/v1/workspaces?page=1&page_size=10
// Response 200
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "AI研发团队",
        "description": "AI模型训练工作空间",
        "owner_id": 1,
        "member_count": 5,
        "created_at": "2026-01-30T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

**3. 添加成员**
```json
// POST /api/v1/workspaces/:id/members
// Request
{
  "user_id": 2,
  "role": "member"  // owner, admin, member
}

// Response 200
{
  "code": 200,
  "message": "success"
}
```

#### 后端开发清单

- [x] 创建 `internal/controller/v1/workspace.go`
- [x] 实现 `Create` 方法
- [x] 实现 `List` 方法（支持分页）
- [x] 实现 `GetByID` 方法
- [x] 实现 `Update` 方法
- [x] 实现 `Delete` 方法
- [x] 实现 `AddMember` 方法
- [x] 实现 `RemoveMember` 方法
- [x] 实现 `ListMembers` 方法
- [x] 在 `internal/router/router.go` 添加路由
- [ ] 编写单元测试
- [x] 更新 API 文档（已在 dev/backend/workspace-api-completed.md）

#### 验收标准
- ✅ 所有 API 接口可通过 Postman 测试
- ✅ 权限控制正确（只有所有者可以删除工作空间）
- ✅ 单元测试覆盖率 > 80%
- ✅ 错误处理完善（参数验证、权限检查）

---

### 前端开发任务

**负责人**: 前端开发
**状态**: ⏳ 待开始（依赖后端 API 完成）

#### 页面开发清单

**1. 工作空间列表页面** (`frontend/src/views/workspace/WorkspaceList.vue`)
- [ ] 工作空间列表展示（表格）
- [ ] 分页功能
- [ ] 创建工作空间按钮
- [ ] 编辑/删除操作
- [ ] 搜索过滤功能

**2. 工作空间创建/编辑页面** (`frontend/src/views/workspace/WorkspaceForm.vue`)
- [ ] 工作空间名称输入
- [ ] 描述输入
- [ ] 资源配额设置（可选）
- [ ] 表单验证
- [ ] 提交/取消按钮

**3. 工作空间详情页面** (`frontend/src/views/workspace/WorkspaceDetail.vue`)
- [ ] 工作空间基本信息展示
- [ ] 成员列表展示
- [ ] 添加成员功能
- [ ] 移除成员功能
- [ ] 成员角色管理

**4. API 调用模块** (`frontend/src/api/workspace.ts`)
- [ ] `createWorkspace(data)` - 创建工作空间
- [ ] `getWorkspaces(page, pageSize)` - 获取工作空间列表
- [ ] `getWorkspaceById(id)` - 获取工作空间详情
- [ ] `updateWorkspace(id, data)` - 更新工作空间
- [ ] `deleteWorkspace(id)` - 删除工作空间
- [ ] `addMember(workspaceId, userId, role)` - 添加成员
- [ ] `removeMember(workspaceId, userId)` - 移除成员
- [ ] `getMembers(workspaceId)` - 获取成员列表

#### 验收标准
- ✅ 所有页面功能正常
- ✅ 与后端 API 对接成功
- ✅ UI 符合 Element Plus 设计规范
- ✅ 响应式设计，支持不同屏幕尺寸
- ✅ 错误提示友好

---

## 🎯 模块二：ResourceQuota 资源配额管理

### 后端开发任务

**负责人**: 后端开发
**状态**: ⏳ 待开始
**Service层**: ✅ 已完成（`internal/service/resource_quota.go`）

#### API 接口定义

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | `/api/v1/admin/quotas` | 设置资源配额 | 管理员 |
| GET | `/api/v1/admin/quotas` | 获取配额列表 | 管理员 |
| GET | `/api/v1/admin/quotas/:id` | 获取配额详情 | 管理员 |
| PUT | `/api/v1/admin/quotas/:id` | 更新资源配额 | 管理员 |
| DELETE | `/api/v1/admin/quotas/:id` | 删除资源配额 | 管理员 |
| GET | `/api/v1/quotas/usage` | 获取当前用户配额使用情况 | 认证用户 |

#### 请求/响应格式

**1. 设置资源配额**
```json
// POST /api/v1/admin/quotas
// Request
{
  "customer_id": 1,
  "workspace_id": null,  // null表示用户级配额
  "max_gpu": 8,
  "max_cpu": 32,
  "max_memory": 131072,
  "max_storage": 1048576,
  "max_environments": 10
}

// Response 200
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "customer_id": 1,
    "workspace_id": null,
    "max_gpu": 8,
    "max_cpu": 32,
    "max_memory": 131072,
    "max_storage": 1048576,
    "max_environments": 10,
    "created_at": "2026-01-30T10:00:00Z"
  }
}
```

**2. 获取配额使用情况**
```json
// GET /api/v1/quotas/usage?workspace_id=1
// Response 200
{
  "code": 200,
  "message": "success",
  "data": {
    "quota": {
      "max_gpu": 8,
      "max_cpu": 32,
      "max_memory": 131072,
      "max_storage": 1048576,
      "max_environments": 10
    },
    "used": {
      "used_gpu": 4,
      "used_cpu": 16,
      "used_memory": 65536,
      "used_storage": 524288,
      "used_environments": 5
    },
    "available": {
      "available_gpu": 4,
      "available_cpu": 16,
      "available_memory": 65536,
      "available_storage": 524288,
      "available_environments": 5
    },
    "usage_percentage": {
      "gpu": 50.0,
      "cpu": 50.0,
      "memory": 50.0,
      "storage": 50.0,
      "environments": 50.0
    }
  }
}
```

#### 后端开发清单

- [ ] 创建 `internal/controller/v1/resource_quota.go`
- [ ] 实现 `SetQuota` 方法
- [ ] 实现 `GetQuota` 方法
- [ ] 实现 `UpdateQuota` 方法
- [ ] 实现 `DeleteQuota` 方法
- [ ] 实现 `GetUsage` 方法（获取配额使用情况）
- [ ] 实现 `ListQuotas` 方法（管理员查看所有配额）
- [ ] 在 `internal/router/router.go` 添加路由
- [ ] 编写单元测试
- [ ] 更新 API 文档

#### 验收标准
- ✅ 所有 API 接口可通过 Postman 测试
- ✅ 配额计算准确
- ✅ 权限控制正确（管理员才能设置配额）
- ✅ 单元测试覆盖率 > 80%

---

### 前端开发任务

**负责人**: 前端开发
**状态**: ⏳ 待开始（依赖后端 API 完成）

#### 页面开发清单

**1. 资源配额管理页面** (`frontend/src/views/quota/QuotaManagement.vue`)
- [ ] 配额列表展示（表格）
- [ ] 设置配额按钮（管理员）
- [ ] 编辑/删除配额操作
- [ ] 用户/工作空间筛选

**2. 配额设置表单** (`frontend/src/views/quota/QuotaForm.vue`)
- [ ] 用户选择
- [ ] 工作空间选择（可选）
- [ ] GPU 配额输入
- [ ] CPU 配额输入
- [ ] 内存配额输入
- [ ] 存储配额输入
- [ ] 环境数量配额输入
- [ ] 表单验证

**3. 配额使用统计页面** (`frontend/src/views/quota/QuotaUsage.vue`)
- [ ] 配额使用情况展示（进度条）
- [ ] 各资源使用百分比可视化（ECharts）
- [ ] 实时更新
- [ ] 告警提示（使用率超过 80%）

**4. API 调用模块** (`frontend/src/api/quota.ts`)
- [ ] `setQuota(data)` - 设置资源配额
- [ ] `getQuotas()` - 获取配额列表
- [ ] `getQuotaById(id)` - 获取配额详情
- [ ] `updateQuota(id, data)` - 更新资源配额
- [ ] `deleteQuota(id)` - 删除资源配额
- [ ] `getUsage(workspaceId)` - 获取配额使用情况

#### 验收标准
- ✅ 所有页面功能正常
- ✅ 与后端 API 对接成功
- ✅ 数据可视化清晰
- ✅ 实时更新配额使用情况
- ✅ 告警提示及时

---

## 📅 开发流程

### 阶段 1: Workspace 模块（优先）

1. **后端开发** (预计 1 天)
   - 实现 Workspace Controller
   - 添加 API 路由
   - 编写单元测试
   - 输出 API 文档

2. **前端开发** (预计 1 天)
   - 开发工作空间管理页面
   - 实现 API 对接
   - 前后端联调测试

3. **QA 测试** (预计 0.5 天)
   - 功能测试
   - 集成测试
   - 性能测试

### 阶段 2: ResourceQuota 模块

1. **后端开发** (预计 1 天)
   - 实现 ResourceQuota Controller
   - 添加 API 路由
   - 编写单元测试
   - 输出 API 文档

2. **前端开发** (预计 1 天)
   - 开发资源配额管理页面
   - 实现 API 对接
   - 前后端联调测试

3. **QA 测试** (预计 0.5 天)
   - 功能测试
   - 集成测试
   - 性能测试

---

## 🔧 IT 运维支持

**负责人**: IT 运维
**状态**: ✅ 已完成

### 基础设施检查清单
- [x] PostgreSQL 服务运行正常
- [x] Redis 服务运行正常
- [x] Etcd 服务运行正常
- [x] 后端配置文件检查
- [ ] 测试数据准备
- [ ] 开发环境文档

---

## 📊 进度跟踪

| 模块 | 后端开发 | 前端开发 | QA测试 | 状态 |
|------|---------|---------|--------|------|
| Workspace | ⏳ 0% | ⏳ 0% | ⏳ 0% | 待开始 |
| ResourceQuota | ⏳ 0% | ⏳ 0% | ⏳ 0% | 待开始 |

---

## 📝 沟通协作规则

1. **后端开发完成后**：
   - 在本文档中更新后端开发清单（勾选完成项）
   - 在 `dev/backend/` 目录输出 API 测试文档
   - 通知前端开发可以开始对接

2. **前端开发完成后**：
   - 在本文档中更新前端开发清单（勾选完成项）
   - 在 `dev/frontend/` 目录输出页面截图和说明
   - 通知 QA 可以开始测试

3. **QA 测试完成后**：
   - 在 `dev/qa/` 目录输出测试报告
   - 如有 Bug，在本文档中记录并分配给对应开发人员

4. **遇到问题时**：
   - 在对应目录（backend/frontend/qa/it）创建问题文档
   - 在本文档中记录问题和解决方案

---

## 🐛 问题跟踪

### 待解决问题
暂无

### 已解决问题
暂无

---

**最后更新**: 2026-01-30
**更新人**: 项目经理
