# 工作空间与环境管理 API 规范

## 1. 概述

本文档定义 RemoteGPU 平台工作空间（Workspace）和环境（Environment）管理的后端 API 规范，与前端已有的 `api/workspace/` 和 `api/environment/` 接口对齐。

### 1.1 现状分析

**前端已定义的接口：**

工作空间（`frontend/src/api/workspace/`）：
- CRUD：创建、列表（分页）、详情、更新、删除
- 成员管理：添加成员、移除成员、成员列表

环境（`frontend/src/api/environment/`）：
- CRUD：创建、列表（按 workspace_id 筛选）、详情、删除
- 生命周期：启动、停止
- 访问信息：获取环境访问信息（SSH/Jupyter/VNC）

**后端现状：**

- `Workspace` 实体已定义（`customer.go`），包含 ID、UUID、OwnerID、Name、Description、Type、Status
- `workspace_members` 表已存在（`03_users_and_permissions.sql`）
- `environments` 表已存在（`05_environments.sql`），但缺少 `workspace_id`、`deployment_mode`、`customer_id` 等字段
- 路由中尚未注册工作空间和环境相关的 API 端点
- 缺少对应的 Service、Controller、DAO 层实现

---

## 2. 数据模型

### 2.1 Workspace 实体（已有，无需修改）

现有 `Workspace` 实体已满足需求：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| uuid | UUID | 外部标识 |
| owner_id | BIGINT | 所有者客户 ID |
| name | VARCHAR(128) | 工作空间名称 |
| description | TEXT | 描述 |
| type | VARCHAR(32) | 类型：personal / team |
| status | VARCHAR(32) | 状态：active / archived |

### 2.2 WorkspaceMember 实体（需新增 Go 结构体）

数据库表 `workspace_members` 已存在，需新增 Go 实体：

```go
// WorkspaceMember 工作空间成员
type WorkspaceMember struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    WorkspaceID uint      `gorm:"not null;uniqueIndex:idx_ws_member" json:"workspace_id"`
    CustomerID  uint      `gorm:"not null;uniqueIndex:idx_ws_member" json:"customer_id"`
    Role        string    `gorm:"type:varchar(32);default:'member'" json:"role"` // owner, admin, member
    Status      string    `gorm:"type:varchar(32);default:'active'" json:"status"`
    JoinedAt    time.Time `gorm:"default:NOW()" json:"joined_at"`
    CreatedAt   time.Time `json:"created_at"`

    // Relations
    Customer  Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
    Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}
```

### 2.3 Environment 实体（需新增 Go 结构体 + 扩展数据库表）

现有 `environments` 表缺少部分字段，需要扩展。新增 Go 实体：

```go
// Environment 开发环境
type Environment struct {
    ID          string  `gorm:"primarykey;type:varchar(64)" json:"id"`
    CustomerID  uint    `gorm:"not null;index" json:"customer_id"`
    WorkspaceID *uint   `gorm:"index" json:"workspace_id,omitempty"`
    HostID      string  `gorm:"type:varchar(64);not null;index" json:"host_id"`
    Name        string  `gorm:"type:varchar(128);not null" json:"name"`
    Description string  `gorm:"type:text" json:"description"`
    Image       string  `gorm:"type:varchar(256);not null" json:"image"`

    // 资源配置
    DeploymentMode string `gorm:"type:varchar(32);default:'docker_local'" json:"deployment_mode"`
    CPU            int    `gorm:"not null" json:"cpu"`
    Memory         int64  `gorm:"not null" json:"memory"`
    GPU            int    `gorm:"default:0" json:"gpu"`
    Storage        *int64 `json:"storage,omitempty"`

    // 端口映射
    SSHPort     *int `json:"ssh_port,omitempty"`
    RDPPort     *int `json:"rdp_port,omitempty"`
    JupyterPort *int `json:"jupyter_port,omitempty"`

    // 容器信息
    ContainerID string `gorm:"type:varchar(128)" json:"container_id"`
    PodName     string `gorm:"type:varchar(128)" json:"pod_name"`

    // 状态与时间
    Status    string     `gorm:"type:varchar(20);default:'creating'" json:"status"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    StartedAt *time.Time `json:"started_at,omitempty"`
    StoppedAt *time.Time `json:"stopped_at,omitempty"`
}
```

**environments 表需要新增的字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `customer_id` | BIGINT | 客户 ID（现有 `user_id` 重命名） |
| `workspace_id` | BIGINT | 所属工作空间 ID |
| `deployment_mode` | VARCHAR(32) | 部署模式：docker_local / k8s_pod / vm |
| `description` | TEXT | 环境描述 |

---

## 3. 工作空间 API 规范

### 3.1 路由总览

所有工作空间 API 挂载在 `/api/v1/customer/workspaces` 下，需要登录认证。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/customer/workspaces` | 创建工作空间 |
| GET | `/customer/workspaces` | 工作空间列表（分页） |
| GET | `/customer/workspaces/:id` | 工作空间详情 |
| PUT | `/customer/workspaces/:id` | 更新工作空间 |
| DELETE | `/customer/workspaces/:id` | 删除工作空间 |
| POST | `/customer/workspaces/:id/members` | 添加成员 |
| DELETE | `/customer/workspaces/:id/members/:userId` | 移除成员 |
| GET | `/customer/workspaces/:id/members` | 成员列表 |

### 3.2 创建工作空间

**POST** `/customer/workspaces`

请求体：
```json
{
  "name": "我的团队",
  "description": "团队工作空间"
}
```

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "uuid": "...",
    "name": "我的团队",
    "description": "团队工作空间",
    "owner_id": 100,
    "member_count": 1,
    "type": "personal",
    "status": "active",
    "created_at": "2026-02-07T10:00:00Z",
    "updated_at": "2026-02-07T10:00:00Z"
  }
}
```

业务规则：
- `owner_id` 自动设为当前登录用户
- 自动在 `workspace_members` 中插入一条 role=owner 的记录
- `member_count` 初始为 1

### 3.3 工作空间列表

**GET** `/customer/workspaces?page=1&page_size=10`

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "items": [{ "id": 1, "name": "...", "..." }],
    "total": 5,
    "page": 1,
    "page_size": 10
  }
}
```

业务规则：
- 返回当前用户拥有的 + 作为成员加入的所有工作空间
- 按 `created_at DESC` 排序

### 3.4 更新工作空间

**PUT** `/customer/workspaces/:id`

请求体：
```json
{
  "name": "新名称",
  "description": "新描述"
}
```

业务规则：
- 仅 owner 或 admin 角色可更新
- 字段为可选，未传的字段不更新

### 3.5 删除工作空间

**DELETE** `/customer/workspaces/:id`

业务规则：
- 仅 owner 可删除
- 删除前检查是否有运行中的环境，有则拒绝删除
- 级联删除 workspace_members 记录

### 3.6 添加成员

**POST** `/customer/workspaces/:id/members`

请求体：
```json
{
  "user_id": 200,
  "role": "member"
}
```

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

业务规则：
- 仅 owner 或 admin 可添加成员
- admin 不能添加 owner 角色
- 不能重复添加已有成员
- 添加成功后更新 workspace.member_count

### 3.7 移除成员

**DELETE** `/customer/workspaces/:id/members/:userId`

业务规则：
- 仅 owner 或 admin 可移除成员
- 不能移除 owner
- admin 不能移除其他 admin
- 移除成功后更新 workspace.member_count

### 3.8 成员列表

**GET** `/customer/workspaces/:id/members`

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "workspace_id": 1,
      "customer_id": 100,
      "username": "user1",
      "email": "user1@example.com",
      "role": "owner",
      "joined_at": "2026-02-07T10:00:00Z"
    }
  ]
}
```

业务规则：
- 工作空间的所有成员均可查看成员列表
- 返回成员的基本信息（username、email）通过 JOIN customers 表获取

---

## 4. 环境管理 API 规范

### 4.1 路由总览

环境 API 挂载在 `/api/v1/customer/environments` 下，需要登录认证。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/customer/environments` | 创建环境 |
| GET | `/customer/environments` | 环境列表 |
| GET | `/customer/environments/:id` | 环境详情 |
| POST | `/customer/environments/:id/start` | 启动环境 |
| POST | `/customer/environments/:id/stop` | 停止环境 |
| DELETE | `/customer/environments/:id` | 删除环境 |
| GET | `/customer/environments/:id/access` | 获取访问信息 |

### 4.2 创建环境

**POST** `/customer/environments`

请求体：
```json
{
  "workspace_id": 1,
  "name": "PyTorch 开发环境",
  "description": "用于模型训练",
  "image": "pytorch/pytorch:2.1-cuda12.1",
  "deployment_mode": "docker_local",
  "cpu": 4,
  "memory": 8192,
  "gpu": 1,
  "storage": 50000,
  "env": {"CUDA_VISIBLE_DEVICES": "0"}
}
```

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "env-abc123",
    "customer_id": 100,
    "workspace_id": 1,
    "host_id": "node-01",
    "name": "PyTorch 开发环境",
    "status": "creating",
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

业务规则：
- `customer_id` 自动设为当前登录用户
- `workspace_id` 可选，不传则不关联工作空间
- 校验客户配额（GPU、存储）是否充足
- 校验客户余额是否充足（余额 > 0 或在信用额度内）
- 自动调度选择可用 Host（基于资源空闲情况）
- 生成唯一环境 ID（格式：`env-` + nanoid）
- 异步创建容器，立即返回 status=creating

### 4.3 环境列表

**GET** `/customer/environments?workspace_id=1`

查询参数：
- `workspace_id`（可选）：按工作空间筛选

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": "env-abc123",
      "name": "PyTorch 开发环境",
      "status": "running",
      "image": "pytorch/pytorch:2.1-cuda12.1",
      "cpu": 4,
      "memory": 8192,
      "gpu": 1,
      "created_at": "2026-02-07T10:00:00Z"
    }
  ]
}
```

业务规则：
- 仅返回当前用户拥有的环境
- 如指定 workspace_id，还需校验用户是否为该工作空间成员

### 4.4 启动环境

**POST** `/customer/environments/:id/start`

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": { "status": "running" }
}
```

业务规则：
- 仅 status=stopped 的环境可启动
- 启动前校验余额是否充足
- 通过 Agent 异步启动容器

### 4.5 停止环境

**POST** `/customer/environments/:id/stop`

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": { "status": "stopped" }
}
```

业务规则：
- 仅 status=running 的环境可停止
- 停止后计费记录截止到当前时间
- 通过 Agent 异步停止容器

### 4.6 删除环境

**DELETE** `/customer/environments/:id`

业务规则：
- 仅 status=stopped 或 status=error 的环境可删除
- 运行中的环境需先停止再删除
- 通过 Agent 清理容器和相关资源
- 释放端口映射

### 4.7 获取访问信息

**GET** `/customer/environments/:id/access`

响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "ssh": {
      "host": "gpu01.example.com",
      "port": 30022,
      "username": "root"
    },
    "jupyter": {
      "url": "https://gpu01.example.com:30088",
      "token": "abc123"
    },
    "vnc": {
      "url": "https://gpu01.example.com:30090"
    }
  }
}
```

业务规则：
- 仅 status=running 的环境可获取访问信息
- 根据环境的端口映射和 Host 的外部映射配置组装访问地址
- SSH/Jupyter/VNC 信息根据实际分配的端口返回，未分配的字段不返回

---

## 5. 后端实现层设计

### 5.1 新增文件清单

#### Entity 层

| 文件 | 说明 |
|------|------|
| `model/entity/workspace.go` | 将 Workspace、WorkspaceMember 从 customer.go 拆出 |
| `model/entity/environment.go` | 新增 Environment 实体 |

#### DAO 层

| 文件 | 说明 |
|------|------|
| `dao/workspace_repo.go` | WorkspaceDao：CRUD + 按成员查询 |
| `dao/workspace_member_repo.go` | WorkspaceMemberDao：成员增删查 |
| `dao/environment_repo.go` | EnvironmentDao：CRUD + 按客户/工作空间查询 |

#### Service 层

| 文件 | 说明 |
|------|------|
| `service/workspace/workspace_service.go` | 工作空间 CRUD + 权限校验 |
| `service/workspace/member_service.go` | 成员管理逻辑 |
| `service/environment/environment_service.go` | 环境生命周期管理 |

#### Controller 层

| 文件 | 说明 |
|------|------|
| `controller/v1/workspace/workspace_controller.go` | 工作空间 HTTP 处理 |
| `controller/v1/environment/environment_controller.go` | 环境 HTTP 处理 |

#### API 请求/响应结构体

| 文件 | 说明 |
|------|------|
| `api/v1/workspace.go` | 工作空间相关请求/响应结构体 |
| `api/v1/environment.go` | 环境相关请求/响应结构体 |

### 5.2 路由注册

在 `router/router.go` 的 `custGroup` 下新增：

```go
// 工作空间管理
custGroup.POST("/workspaces", workspaceController.Create)
custGroup.GET("/workspaces", workspaceController.List)
custGroup.GET("/workspaces/:id", workspaceController.Detail)
custGroup.PUT("/workspaces/:id", workspaceController.Update)
custGroup.DELETE("/workspaces/:id", workspaceController.Delete)
custGroup.POST("/workspaces/:id/members", workspaceController.AddMember)
custGroup.DELETE("/workspaces/:id/members/:userId", workspaceController.RemoveMember)
custGroup.GET("/workspaces/:id/members", workspaceController.ListMembers)
```

```go
// 环境管理
custGroup.POST("/environments", environmentController.Create)
custGroup.GET("/environments", environmentController.List)
custGroup.GET("/environments/:id", environmentController.Detail)
custGroup.POST("/environments/:id/start", environmentController.Start)
custGroup.POST("/environments/:id/stop", environmentController.Stop)
custGroup.DELETE("/environments/:id", environmentController.Delete)
custGroup.GET("/environments/:id/access", environmentController.AccessInfo)
```

---

## 6. 数据库迁移脚本清单

| 编号 | 文件名 | 说明 |
|------|--------|------|
| 38 | `38_extend_environments.sql` | 扩展 environments 表：重命名 user_id→customer_id，新增 workspace_id、deployment_mode、description |
| 39 | `39_add_environment_indexes.sql` | 为 environments 表新增 workspace_id、customer_id 索引 |

> **注意**：workspaces 和 workspace_members 表已在 `03_users_and_permissions.sql` 中创建，无需新增迁移脚本。

---

## 7. 前后端接口对齐说明

前端已定义的 API 路径使用的是不带 `/customer` 前缀的路径（如 `/workspaces`、`/environments`），这是因为前端的 `request` 实例已配置了 `baseURL`（通常为 `/api/v1/customer`）。

### 7.1 工作空间接口对齐

| 前端调用 | 后端完整路径 | 状态 |
|----------|-------------|------|
| `POST /workspaces` | `/api/v1/customer/workspaces` | 待实现 |
| `GET /workspaces?page=&page_size=` | `/api/v1/customer/workspaces` | 待实现 |
| `GET /workspaces/:id` | `/api/v1/customer/workspaces/:id` | 待实现 |
| `PUT /workspaces/:id` | `/api/v1/customer/workspaces/:id` | 待实现 |
| `DELETE /workspaces/:id` | `/api/v1/customer/workspaces/:id` | 待实现 |
| `POST /workspaces/:id/members` | `/api/v1/customer/workspaces/:id/members` | 待实现 |
| `DELETE /workspaces/:id/members/:userId` | `/api/v1/customer/workspaces/:id/members/:userId` | 待实现 |
| `GET /workspaces/:id/members` | `/api/v1/customer/workspaces/:id/members` | 待实现 |

### 7.2 环境接口对齐

| 前端调用 | 后端完整路径 | 状态 |
|----------|-------------|------|
| `GET /environments?workspace_id=` | `/api/v1/customer/environments` | 待实现 |
| `POST /environments` | `/api/v1/customer/environments` | 待实现 |
| `GET /environments/:id` | `/api/v1/customer/environments/:id` | 待实现 |
| `POST /environments/:id/start` | `/api/v1/customer/environments/:id/start` | 待实现 |
| `POST /environments/:id/stop` | `/api/v1/customer/environments/:id/stop` | 待实现 |
| `DELETE /environments/:id` | `/api/v1/customer/environments/:id` | 待实现 |
| `GET /environments/:id/access` | `/api/v1/customer/environments/:id/access` | 待实现 |
