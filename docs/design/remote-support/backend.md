# 远程客户支持平台 — 后端技术方案

> 版本：v1.0 | 日期：2026-02-06 | 作者：后端工程师

## 1. 概述

### 1.1 文档目的

本文档描述远程客户支持平台后端的详细技术方案，包括 API 设计、数据库设计、会话管理、代理转发、认证授权、审计日志等核心模块的实现方案，以及与现有 RemoteGPU 系统的集成策略。

### 1.2 设计原则

- **与现有架构一致**：遵循 Controller → Service → DAO 三层分离架构
- **最小侵入**：新增模块不修改现有代码，通过新增文件和路由扩展
- **复用优先**：复用现有的认证中间件、审计中间件、响应格式等基础设施
- **安全第一**：所有远程访问操作必须经过认证、授权和审计

### 1.3 技术栈

沿用现有技术栈：
- **语言**：Go 1.21+
- **Web 框架**：Gin
- **ORM**：GORM
- **数据库**：PostgreSQL
- **缓存**：Redis
- **远程访问网关**：Apache Guacamole（新增）
- **反向代理**：Nginx / OpenResty（新增）

---

## 2. 数据库设计

### 2.1 新增表概览

| 表名 | 说明 |
|------|------|
| `remote_access_configs` | 机器远程访问配置（域名、端口、协议） |
| `remote_sessions` | 远程访问会话记录 |
| `remote_session_events` | 会话事件日志（连接、断开、错误等） |
| `proxy_routes` | Nginx 反向代理路由配置 |

### 2.2 remote_access_configs — 远程访问配置表

存储每台机器的远程访问配置，与 `hosts` 表一对一关联。

```sql
CREATE TABLE remote_access_configs (
    id          SERIAL PRIMARY KEY,
    host_id     VARCHAR(64) NOT NULL REFERENCES hosts(id),
    enabled     BOOLEAN NOT NULL DEFAULT FALSE,
    protocol    VARCHAR(16) NOT NULL DEFAULT 'ssh',  -- ssh/vnc/rdp/http/https
    public_domain VARCHAR(255),                       -- 对外访问域名
    public_port   INT,                                -- 对外开放端口
    target_port   INT,                                -- 目标服务端口
    extra_ports   VARCHAR(255),                       -- 额外开放端口，逗号分隔
    auth_mode     VARCHAR(32) DEFAULT 'password',     -- password/key/guacamole
    max_sessions  INT DEFAULT 5,                      -- 最大并发会话数
    idle_timeout  INT DEFAULT 1800,                   -- 空闲超时（秒）
    remark        TEXT,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uk_remote_access_host UNIQUE (host_id),
    CONSTRAINT uk_remote_access_domain_port UNIQUE (public_domain, public_port)
);

CREATE INDEX idx_remote_access_enabled ON remote_access_configs(enabled);
```

### 2.3 remote_sessions — 远程会话表

记录每次远程访问会话的生命周期。

```sql
CREATE TABLE remote_sessions (
    id            VARCHAR(64) PRIMARY KEY,            -- UUID
    host_id       VARCHAR(64) NOT NULL,
    customer_id   INT NOT NULL,
    config_id     INT NOT NULL REFERENCES remote_access_configs(id),

    -- 会话信息
    protocol      VARCHAR(16) NOT NULL,               -- ssh/vnc/rdp
    client_ip     VARCHAR(64),                        -- 客户端 IP
    client_info   VARCHAR(512),                       -- 客户端信息（User-Agent 等）

    -- Guacamole 集成
    guac_connection_id VARCHAR(128),                  -- Guacamole 连接 ID
    guac_token         VARCHAR(512),                  -- Guacamole 认证 token

    -- 状态
    status        VARCHAR(20) NOT NULL DEFAULT 'connecting',
                  -- connecting / active / idle / disconnected / terminated / error
    error_msg     TEXT,

    -- 时间
    connected_at    TIMESTAMP,
    last_active_at  TIMESTAMP,
    disconnected_at TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),

    -- 外键
    CONSTRAINT fk_session_host FOREIGN KEY (host_id) REFERENCES hosts(id),
    CONSTRAINT fk_session_customer FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE INDEX idx_session_host ON remote_sessions(host_id);
CREATE INDEX idx_session_customer ON remote_sessions(customer_id);
CREATE INDEX idx_session_status ON remote_sessions(status);
CREATE INDEX idx_session_created ON remote_sessions(created_at);
```

### 2.4 remote_session_events — 会话事件表

记录会话生命周期中的关键事件，用于审计和排障。

```sql
CREATE TABLE remote_session_events (
    id          BIGSERIAL PRIMARY KEY,
    session_id  VARCHAR(64) NOT NULL REFERENCES remote_sessions(id),
    event_type  VARCHAR(32) NOT NULL,
    -- connect / disconnect / idle_timeout / error / terminate / resume
    detail      JSONB,
    client_ip   VARCHAR(64),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_session_event_session ON remote_session_events(session_id);
CREATE INDEX idx_session_event_type ON remote_session_events(event_type);
CREATE INDEX idx_session_event_created ON remote_session_events(created_at);
```

### 2.5 proxy_routes — 代理路由表

存储 Nginx 反向代理的路由配置，支持动态生成 Nginx 配置。

```sql
CREATE TABLE proxy_routes (
    id            SERIAL PRIMARY KEY,
    config_id     INT NOT NULL REFERENCES remote_access_configs(id),
    host_id       VARCHAR(64) NOT NULL,

    -- 路由信息
    listen_port   INT NOT NULL,
    server_name   VARCHAR(255),                       -- Nginx server_name
    upstream_addr VARCHAR(255) NOT NULL,              -- 上游地址 ip:port
    protocol      VARCHAR(16) NOT NULL DEFAULT 'tcp', -- tcp/http/https

    -- SSL
    ssl_enabled   BOOLEAN DEFAULT FALSE,
    ssl_cert_path VARCHAR(512),
    ssl_key_path  VARCHAR(512),

    -- 状态
    status        VARCHAR(20) DEFAULT 'pending',      -- pending/active/error/disabled
    last_synced_at TIMESTAMP,
    sync_error    TEXT,

    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),

    CONSTRAINT uk_proxy_route_listen UNIQUE (listen_port, server_name)
);

CREATE INDEX idx_proxy_route_host ON proxy_routes(host_id);
CREATE INDEX idx_proxy_route_status ON proxy_routes(status);
```

### 2.6 ER 关系图

```
hosts (1) ──── (1) remote_access_configs (1) ──── (N) proxy_routes
  │                       │
  │                       │
  └───── (N) remote_sessions (1) ──── (N) remote_session_events
                │
                │
         customers (1) ──── (N) remote_sessions
```

---

## 3. 实体模型（Entity）

### 3.1 RemoteAccessConfig

文件路径：`backend/internal/model/entity/remote_access.go`

```go
// RemoteAccessConfig 远程访问配置实体
type RemoteAccessConfig struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    HostID      string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"host_id"`
    Enabled     bool      `gorm:"default:false" json:"enabled"`
    Protocol    string    `gorm:"type:varchar(16);not null;default:'ssh'" json:"protocol"`
    PublicDomain string   `gorm:"type:varchar(255)" json:"public_domain"`
    PublicPort  int       `json:"public_port"`
    TargetPort  int       `json:"target_port"`
    ExtraPorts  string    `gorm:"type:varchar(255)" json:"extra_ports"`
    AuthMode    string    `gorm:"type:varchar(32);default:'password'" json:"auth_mode"`
    MaxSessions int       `gorm:"default:5" json:"max_sessions"`
    IdleTimeout int       `gorm:"default:1800" json:"idle_timeout"`
    Remark      string    `gorm:"type:text" json:"remark"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // Relations
    Host        Host         `gorm:"foreignKey:HostID" json:"host,omitempty"`
    ProxyRoutes []ProxyRoute `gorm:"foreignKey:ConfigID" json:"proxy_routes,omitempty"`
}

func (RemoteAccessConfig) TableName() string { return "remote_access_configs" }
```

### 3.2 RemoteSession

```go
// RemoteSession 远程会话实体
type RemoteSession struct {
    ID         string `gorm:"primarykey;type:varchar(64)" json:"id"`
    HostID     string `gorm:"type:varchar(64);not null" json:"host_id"`
    CustomerID uint   `gorm:"not null" json:"customer_id"`
    ConfigID   uint   `gorm:"not null" json:"config_id"`

    Protocol   string `gorm:"type:varchar(16);not null" json:"protocol"`
    ClientIP   string `gorm:"type:varchar(64)" json:"client_ip"`
    ClientInfo string `gorm:"type:varchar(512)" json:"client_info"`

    GuacConnectionID string `gorm:"type:varchar(128)" json:"guac_connection_id"`
    GuacToken        string `gorm:"type:varchar(512)" json:"-"`

    Status   string `gorm:"type:varchar(20);not null;default:'connecting'" json:"status"`
    ErrorMsg string `gorm:"type:text" json:"error_msg"`

    ConnectedAt    *time.Time `json:"connected_at"`
    LastActiveAt   *time.Time `json:"last_active_at"`
    DisconnectedAt *time.Time `json:"disconnected_at"`
    CreatedAt      time.Time  `json:"created_at"`

    // Relations
    Host     Host              `gorm:"foreignKey:HostID" json:"host,omitempty"`
    Customer Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
    Config   RemoteAccessConfig `gorm:"foreignKey:ConfigID" json:"config,omitempty"`
    Events   []RemoteSessionEvent `gorm:"foreignKey:SessionID" json:"events,omitempty"`
}

func (RemoteSession) TableName() string { return "remote_sessions" }
```

### 3.3 RemoteSessionEvent

```go
// RemoteSessionEvent 会话事件实体
type RemoteSessionEvent struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    SessionID string         `gorm:"type:varchar(64);not null" json:"session_id"`
    EventType string         `gorm:"type:varchar(32);not null" json:"event_type"`
    Detail    datatypes.JSON `gorm:"type:jsonb" json:"detail"`
    ClientIP  string         `gorm:"type:varchar(64)" json:"client_ip"`
    CreatedAt time.Time      `json:"created_at"`
}

func (RemoteSessionEvent) TableName() string { return "remote_session_events" }
```

### 3.4 ProxyRoute

```go
// ProxyRoute 代理路由实体
type ProxyRoute struct {
    ID           uint   `gorm:"primarykey" json:"id"`
    ConfigID     uint   `gorm:"not null" json:"config_id"`
    HostID       string `gorm:"type:varchar(64);not null" json:"host_id"`
    ListenPort   int    `gorm:"not null" json:"listen_port"`
    ServerName   string `gorm:"type:varchar(255)" json:"server_name"`
    UpstreamAddr string `gorm:"type:varchar(255);not null" json:"upstream_addr"`
    Protocol     string `gorm:"type:varchar(16);not null;default:'tcp'" json:"protocol"`

    SSLEnabled  bool   `gorm:"default:false" json:"ssl_enabled"`
    SSLCertPath string `gorm:"type:varchar(512)" json:"ssl_cert_path"`
    SSLKeyPath  string `gorm:"type:varchar(512)" json:"ssl_key_path"`

    Status       string     `gorm:"type:varchar(20);default:'pending'" json:"status"`
    LastSyncedAt *time.Time `json:"last_synced_at"`
    SyncError    string     `gorm:"type:text" json:"sync_error"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (ProxyRoute) TableName() string { return "proxy_routes" }
```

---

## 4. API 设计

### 4.1 API 概览

所有新增 API 遵循现有 `/api/v1` 前缀，按角色分组。

| 分组 | 前缀 | 中间件 | 说明 |
|------|------|--------|------|
| 管理员 — 远程访问配置 | `/api/v1/admin/remote-access` | Auth + RequireAdmin + Audit | 配置管理 |
| 管理员 — 会话管理 | `/api/v1/admin/sessions` | Auth + RequireAdmin + Audit | 会话监控与管理 |
| 管理员 — 代理路由 | `/api/v1/admin/proxy-routes` | Auth + RequireAdmin + Audit | Nginx 路由管理 |
| 客户 — 远程连接 | `/api/v1/customer/remote` | Auth | 发起和管理远程连接 |

### 4.2 管理员 — 远程访问配置 API

#### GET /api/v1/admin/machines/:id/remote-access

获取指定机器的远程访问配置。

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "host_id": "node-01",
    "enabled": true,
    "protocol": "ssh",
    "public_domain": "node-01.remote.example.com",
    "public_port": 2222,
    "target_port": 22,
    "extra_ports": "6006,8888",
    "auth_mode": "password",
    "max_sessions": 5,
    "idle_timeout": 1800,
    "remark": "SSH + Jupyter + TensorBoard"
  }
}
```

#### PUT /api/v1/admin/machines/:id/remote-access

创建或更新机器的远程访问配置（Upsert 语义）。

**请求**：
```json
{
  "enabled": true,
  "protocol": "ssh",
  "public_domain": "node-01.remote.example.com",
  "public_port": 2222,
  "target_port": 22,
  "extra_ports": "6006,8888",
  "auth_mode": "password",
  "max_sessions": 5,
  "idle_timeout": 1800,
  "remark": "SSH + Jupyter + TensorBoard"
}
```

**校验规则**：
- `protocol` 必须为 ssh/vnc/rdp/http/https 之一
- `public_domain` + `public_port` 组合唯一
- `public_port` 范围 1-65535，不能与已有路由冲突
- `max_sessions` 范围 1-100
- `idle_timeout` 范围 60-86400（1分钟到24小时）

#### POST /api/v1/admin/machines/:id/remote-access/sync

触发 Nginx 配置同步，将远程访问配置生成 Nginx 配置并 reload。

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "status": "synced",
    "proxy_route_id": 1
  }
}
```

### 4.3 管理员 — 会话管理 API

#### GET /api/v1/admin/sessions

查询所有远程会话列表（分页）。

**查询参数**：
- `page` / `pageSize` — 分页
- `host_id` — 按机器筛选
- `customer_id` — 按客户筛选
- `status` — 按状态筛选（active/disconnected/terminated）
- `protocol` — 按协议筛选

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "sess-uuid-001",
        "host_id": "node-01",
        "host_name": "GPU-Node-01",
        "customer_id": 5,
        "customer_name": "张三",
        "protocol": "ssh",
        "client_ip": "192.168.1.100",
        "status": "active",
        "connected_at": "2026-02-06T10:00:00Z",
        "last_active_at": "2026-02-06T10:30:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "pageSize": 20
  }
}
```

#### POST /api/v1/admin/sessions/:id/terminate

管理员强制终止指定会话。

**请求**：
```json
{
  "reason": "维护需要，强制断开"
}
```

#### GET /api/v1/admin/sessions/:id/events

查询指定会话的事件日志。

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "event_type": "connect",
      "detail": {"client_ip": "192.168.1.100"},
      "created_at": "2026-02-06T10:00:00Z"
    },
    {
      "id": 2,
      "event_type": "idle_timeout",
      "detail": {"idle_seconds": 1800},
      "created_at": "2026-02-06T10:30:00Z"
    }
  ]
}
```

### 4.4 客户 — 远程连接 API

#### POST /api/v1/customer/machines/:id/remote/connect

客户发起远程连接请求。后端校验权限后，创建会话并返回连接信息。

**请求**：
```json
{
  "protocol": "ssh"
}
```

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "session_id": "sess-uuid-001",
    "protocol": "ssh",
    "connection": {
      "host": "node-01.remote.example.com",
      "port": 2222,
      "username": "customer_user"
    },
    "guacamole": {
      "url": "https://guac.example.com/#/client/sess-uuid-001",
      "token": "guac-auth-token"
    }
  }
}
```

**业务逻辑**：
1. 校验客户对该机器的分配权限（通过 AllocationService）
2. 校验机器远程访问配置已启用
3. 校验当前活跃会话数未超过 `max_sessions`
4. 创建 RemoteSession 记录
5. 如果使用 Guacamole，调用 Guacamole API 创建连接
6. 记录 connect 事件
7. 返回连接信息

#### GET /api/v1/customer/remote/sessions

查询当前客户的远程会话列表。

**查询参数**：
- `status` — 按状态筛选（active/disconnected）
- `host_id` — 按机器筛选

#### POST /api/v1/customer/remote/sessions/:id/disconnect

客户主动断开远程会话。

#### GET /api/v1/customer/machines/:id/remote/status

查询指定机器的远程访问状态（是否启用、当前会话数等）。

**响应**：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "enabled": true,
    "protocol": "ssh",
    "active_sessions": 2,
    "max_sessions": 5,
    "connection_info": {
      "host": "node-01.remote.example.com",
      "port": 2222
    }
  }
}
```

---

## 5. 会话管理实现

### 5.1 会话生命周期

```
客户发起连接请求
  → 创建 RemoteSession (status=connecting)
  → 校验权限和配置
  → 建立连接（直连 / Guacamole）
  → 连接成功 → status=active，记录 connect 事件
  → 连接失败 → status=error，记录 error 事件
  → 活跃中 → 定期更新 last_active_at
  → 空闲超时 → status=idle → 自动断开 → status=disconnected
  → 客户主动断开 → status=disconnected
  → 管理员强制终止 → status=terminated
```

### 5.2 SessionService 核心接口

文件路径：`backend/internal/service/remote/session_service.go`

```go
type SessionService struct {
    sessionDao    *dao.RemoteSessionDao
    configDao     *dao.RemoteAccessConfigDao
    eventDao      *dao.RemoteSessionEventDao
    allocSvc      *allocation.AllocationService
    guacClient    *guacamole.Client
    db            *gorm.DB
}

// CreateSession 创建远程会话
func (s *SessionService) CreateSession(ctx context.Context, customerID uint, hostID, protocol string, clientIP string) (*entity.RemoteSession, error)

// DisconnectSession 断开会话
func (s *SessionService) DisconnectSession(ctx context.Context, sessionID string, customerID uint) error

// TerminateSession 管理员强制终止会话
func (s *SessionService) TerminateSession(ctx context.Context, sessionID string, reason string) error

// ListSessions 查询会话列表（分页）
func (s *SessionService) ListSessions(ctx context.Context, filter SessionFilter) ([]entity.RemoteSession, int64, error)

// GetSessionEvents 查询会话事件
func (s *SessionService) GetSessionEvents(ctx context.Context, sessionID string) ([]entity.RemoteSessionEvent, error)

// UpdateHeartbeat 更新会话心跳（前端定期调用）
func (s *SessionService) UpdateHeartbeat(ctx context.Context, sessionID string) error

// CleanupIdleSessions 清理空闲超时的会话（定时任务）
func (s *SessionService) CleanupIdleSessions(ctx context.Context) error
```

### 5.3 空闲会话清理（定时任务）

参考现有 `AllocationService.StartWorker` 模式，在 `SessionService` 中启动后台 goroutine 定期清理空闲会话。

```go
// StartIdleCleanupWorker 启动空闲会话清理定时任务
func (s *SessionService) StartIdleCleanupWorker(ctx context.Context) {
    ticker := time.NewTicker(60 * time.Second)
    go func() {
        for {
            select {
            case <-ctx.Done():
                ticker.Stop()
                return
            case <-ticker.C:
                _ = s.CleanupIdleSessions(ctx)
            }
        }
    }()
}
```

**清理逻辑**：
1. 查询 `status=active` 且 `last_active_at < NOW() - idle_timeout` 的会话
2. 将状态更新为 `disconnected`
3. 如果有 Guacamole 连接，调用 Guacamole API 断开
4. 记录 `idle_timeout` 事件

---

## 6. 代理转发机制

### 6.1 架构概览

```
客户浏览器
    │
    ▼
Nginx 反向代理（公网入口）
    │
    ├── SSH: TCP stream 转发 → 目标机器:22
    ├── VNC: WebSocket 转发 → Guacamole → 目标机器:5900
    ├── RDP: WebSocket 转发 → Guacamole → 目标机器:3389
    └── HTTP/HTTPS: HTTP 反向代理 → 目标机器:8888 (Jupyter 等)
```

### 6.2 ProxyService 核心接口

文件路径：`backend/internal/service/remote/proxy_service.go`

```go
type ProxyService struct {
    routeDao  *dao.ProxyRouteDao
    configDao *dao.RemoteAccessConfigDao
    db        *gorm.DB
    nginxConf NginxConfig
}

type NginxConfig struct {
    ConfDir    string // Nginx 配置片段目录，如 /etc/nginx/conf.d/remote/
    ReloadCmd  string // reload 命令，如 "nginx -s reload"
    StreamDir  string // TCP stream 配置目录
    BaseDomain string // 基础域名
}

// SyncRoute 同步单台机器的代理路由到 Nginx
func (s *ProxyService) SyncRoute(ctx context.Context, hostID string) error

// SyncAllRoutes 全量同步所有启用的代理路由
func (s *ProxyService) SyncAllRoutes(ctx context.Context) error

// RemoveRoute 移除指定机器的代理路由
func (s *ProxyService) RemoveRoute(ctx context.Context, hostID string) error

// GenerateNginxConfig 生成 Nginx 配置片段
func (s *ProxyService) GenerateNginxConfig(route *entity.ProxyRoute) (string, error)

// ReloadNginx 执行 Nginx reload
func (s *ProxyService) ReloadNginx(ctx context.Context) error
```

### 6.3 Nginx 配置模板

#### HTTP/HTTPS 反向代理模板

```nginx
# 自动生成 — 请勿手动修改
# host_id: {{.HostID}} | config_id: {{.ConfigID}}
server {
    listen {{.ListenPort}} {{if .SSLEnabled}}ssl{{end}};
    server_name {{.ServerName}};

    {{if .SSLEnabled}}
    ssl_certificate     {{.SSLCertPath}};
    ssl_certificate_key {{.SSLKeyPath}};
    {{end}}

    location / {
        proxy_pass http://{{.UpstreamAddr}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持（Jupyter/VNC 等需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400s;
    }
}
```

#### TCP Stream 转发模板（SSH）

```nginx
# 自动生成 — host_id: {{.HostID}}
stream {
    server {
        listen {{.ListenPort}};
        proxy_pass {{.UpstreamAddr}};
        proxy_timeout 600s;
        proxy_connect_timeout 10s;
    }
}
```

### 6.4 配置同步流程

```
管理员保存远程访问配置
  → ProxyService.SyncRoute(hostID)
  → 查询 RemoteAccessConfig
  → 生成 ProxyRoute 记录（Upsert）
  → 根据 protocol 选择模板生成 Nginx 配置
  → 写入配置文件到 ConfDir/StreamDir
  → 执行 nginx -t 校验配置
  → 校验通过 → nginx -s reload
  → 更新 ProxyRoute.status = active
  → 校验失败 → 记录错误，status = error
```

---

## 7. Guacamole 集成

### 7.1 集成方式

通过 Guacamole REST API 实现 Web 远程桌面（VNC/RDP/SSH）。后端作为中间层，负责认证和连接管理。

```
客户浏览器 → 后端 API（创建会话）→ Guacamole REST API（创建连接）
客户浏览器 → Guacamole Web（WebSocket 通道）→ guacd → 目标机器
```

### 7.2 GuacamoleClient

文件路径：`backend/pkg/guacamole/client.go`

```go
type Client struct {
    baseURL  string // Guacamole API 地址
    username string // 管理员账号
    password string // 管理员密码
    http     *http.Client
}

// Authenticate 获取 Guacamole auth token
func (c *Client) Authenticate(ctx context.Context) (string, error)

// CreateConnection 创建连接
func (c *Client) CreateConnection(ctx context.Context, req CreateConnectionRequest) (string, error)

// DeleteConnection 删除连接
func (c *Client) DeleteConnection(ctx context.Context, connID string) error

// GetActiveConnections 获取活跃连接列表
func (c *Client) GetActiveConnections(ctx context.Context) ([]ActiveConnection, error)

// KillConnection 强制断开连接
func (c *Client) KillConnection(ctx context.Context, connID string) error
```

### 7.3 连接参数映射

| 协议 | Guacamole 参数 | 来源 |
|------|---------------|------|
| SSH | hostname, port, username, password | Host 表 + RemoteAccessConfig |
| VNC | hostname, port, password | Host 表 + RemoteAccessConfig |
| RDP | hostname, port, username, password, domain | Host 表 + RemoteAccessConfig |

---

## 8. 认证授权实现

### 8.1 复用现有认证体系

远程访问模块完全复用现有的 JWT 认证体系，不引入新的认证机制。

| 组件 | 复用方式 |
|------|----------|
| JWT 中间件 | 复用 `middleware.Auth(db)` |
| 角色中间件 | 复用 `middleware.RequireAdmin()` / `middleware.RequireRole()` |
| Token 刷新 | 复用现有 refresh_token 机制 |
| 黑名单 | 复用 Redis token 黑名单 |

### 8.2 远程访问权限校验

客户发起远程连接时，需要进行多层权限校验：

```
1. JWT 认证 → 确认用户身份
2. 账号状态校验 → status=active
3. 分配权限校验 → 通过 AllocationService 确认客户拥有该机器的有效分配
4. 远程访问配置校验 → 确认机器已启用远程访问
5. 会话数限制校验 → 当前活跃会话数 < max_sessions
```

### 8.3 路由注册

```go
// router.go 中新增路由

// 管理员 — 远程访问管理
adminGroup.GET("/machines/:id/remote-access", remoteAccessCtrl.Get)
adminGroup.PUT("/machines/:id/remote-access", remoteAccessCtrl.Upsert)
adminGroup.POST("/machines/:id/remote-access/sync", remoteAccessCtrl.Sync)
adminGroup.GET("/sessions", sessionCtrl.AdminList)
adminGroup.POST("/sessions/:id/terminate", sessionCtrl.Terminate)
adminGroup.GET("/sessions/:id/events", sessionCtrl.Events)

// 客户 — 远程连接
custGroup.POST("/machines/:id/remote/connect", sessionCtrl.Connect)
custGroup.GET("/remote/sessions", sessionCtrl.MyList)
custGroup.POST("/remote/sessions/:id/disconnect", sessionCtrl.Disconnect)
custGroup.GET("/machines/:id/remote/status", remoteAccessCtrl.Status)
custGroup.POST("/remote/sessions/:id/heartbeat", sessionCtrl.Heartbeat)
```

---

## 9. 审计日志实现

### 9.1 复用现有审计中间件

管理员路由已挂载 `middleware.AuditMiddleware(auditSvc)`，所有写操作自动记录到 `audit_logs` 表。新增的管理员远程访问 API 自动享有审计能力。

需要扩展 `parseResource` 函数以识别新的资源类型：

```go
// middleware/audit.go 中扩展
case contains(path, "remote-access"):
    resourceType = "remote_access"
case contains(path, "sessions"):
    resourceType = "remote_session"
case contains(path, "proxy-routes"):
    resourceType = "proxy_route"
```

### 9.2 会话级审计

除了 HTTP 请求级别的审计，远程会话还需要更细粒度的事件审计，通过 `remote_session_events` 表实现。

| 事件类型 | 触发时机 | 记录内容 |
|----------|----------|----------|
| `connect` | 会话建立成功 | client_ip, protocol, user_agent |
| `disconnect` | 客户主动断开 | 断开原因 |
| `terminate` | 管理员强制终止 | 操作者、终止原因 |
| `idle_timeout` | 空闲超时自动断开 | 空闲时长 |
| `error` | 连接异常 | 错误信息 |
| `resume` | 会话恢复 | 恢复原因 |

### 9.3 审计查询扩展

在现有审计日志查询 API 中，新增 `resource_type` 筛选值：
- `remote_access` — 远程访问配置变更
- `remote_session` — 会话管理操作（终止等）
- `proxy_route` — 代理路由变更

---

## 10. 与现有 RemoteGPU 系统的集成

### 10.1 集成策略

采用**扩展式集成**，不修改现有代码，通过新增文件和路由实现。

| 集成点 | 方式 | 说明 |
|--------|------|------|
| 认证 | 复用 | 直接使用现有 `middleware.Auth` |
| 审计 | 复用 + 扩展 | 复用中间件，扩展 `parseResource` |
| 分配权限 | 依赖 | 调用 `AllocationService` 校验机器归属 |
| 机器信息 | 依赖 | 调用 `MachineService` 获取机器详情 |
| Agent 通信 | 依赖 | 调用 `AgentService` 执行远程操作 |
| 路由注册 | 扩展 | 在 `router.go` 中新增路由组 |

### 10.2 新增文件清单

```
backend/
├── internal/
│   ├── model/entity/
│   │   └── remote_access.go          # 新增实体定义
│   ├── dao/
│   │   ├── remote_access_config_repo.go
│   │   ├── remote_session_repo.go
│   │   ├── remote_session_event_repo.go
│   │   └── proxy_route_repo.go
│   ├── service/remote/
│   │   ├── config_service.go          # 远程访问配置服务
│   │   ├── session_service.go         # 会话管理服务
│   │   └── proxy_service.go           # 代理路由服务
│   └── controller/v1/remote/
│       ├── remote_access_controller.go
│       └── session_controller.go
├── pkg/guacamole/
│   └── client.go                      # Guacamole REST API 客户端
├── api/v1/
│   └── remote_access.go               # 请求/响应结构体
└── sql/
    └── 18_remote_access.sql           # 数据库迁移脚本
```

### 10.3 需修改的现有文件

| 文件 | 修改内容 |
|------|----------|
| `internal/router/router.go` | 新增 Service/Controller 初始化和路由注册 |
| `internal/middleware/audit.go` | `parseResource` 新增 remote_access/session/proxy_route 类型 |

### 10.4 router.go 集成示例

```go
// --- 远程访问模块 Service 初始化 ---
remoteConfigSvc := serviceRemote.NewConfigService(db)
remoteSessionSvc := serviceRemote.NewSessionService(db, allocSvc, guacClient)
remoteSessionSvc.StartIdleCleanupWorker(context.Background())
remoteProxySvc := serviceRemote.NewProxyService(db, nginxConf)

// --- 远程访问模块 Controller 初始化 ---
remoteAccessCtrl := ctrlRemote.NewRemoteAccessController(remoteConfigSvc, remoteProxySvc)
sessionCtrl := ctrlRemote.NewSessionController(remoteSessionSvc)

// 管理员路由（已有 adminGroup 下追加）
adminGroup.GET("/machines/:id/remote-access", remoteAccessCtrl.Get)
adminGroup.PUT("/machines/:id/remote-access", remoteAccessCtrl.Upsert)
adminGroup.POST("/machines/:id/remote-access/sync", remoteAccessCtrl.Sync)
adminGroup.GET("/sessions", sessionCtrl.AdminList)
adminGroup.POST("/sessions/:id/terminate", sessionCtrl.Terminate)
adminGroup.GET("/sessions/:id/events", sessionCtrl.Events)

// 客户路由（已有 custGroup 下追加）
custGroup.POST("/machines/:id/remote/connect", sessionCtrl.Connect)
custGroup.GET("/remote/sessions", sessionCtrl.MyList)
custGroup.POST("/remote/sessions/:id/disconnect", sessionCtrl.Disconnect)
custGroup.GET("/machines/:id/remote/status", remoteAccessCtrl.Status)
custGroup.POST("/remote/sessions/:id/heartbeat", sessionCtrl.Heartbeat)
```

---

## 11. 配置管理

### 11.1 新增配置项

在 `config/config.go` 中新增远程访问相关配置：

```go
type RemoteAccessConfig struct {
    BaseDomain string          `yaml:"base_domain"`
    Nginx      NginxConfig     `yaml:"nginx"`
    Guacamole  GuacamoleConfig `yaml:"guacamole"`
}

type NginxConfig struct {
    ConfDir   string `yaml:"conf_dir"`
    StreamDir string `yaml:"stream_dir"`
    ReloadCmd string `yaml:"reload_cmd"`
}

type GuacamoleConfig struct {
    Enabled  bool   `yaml:"enabled"`
    BaseURL  string `yaml:"base_url"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}
```

### 11.2 配置文件示例

```yaml
# config.yaml 新增部分
remote_access:
  base_domain: "remote.example.com"
  nginx:
    conf_dir: "/etc/nginx/conf.d/remote/"
    stream_dir: "/etc/nginx/stream.d/remote/"
    reload_cmd: "nginx -s reload"
  guacamole:
    enabled: true
    base_url: "http://guacamole:8080/guacamole"
    username: "guacadmin"
    password: "${GUAC_PASSWORD}"
```

---

## 12. 安全考虑

### 12.1 传输安全

- 所有公网 API 通过 HTTPS 传输
- Guacamole WebSocket 通道使用 WSS
- Nginx 反向代理支持 SSL 终止
- SSH 转发使用 TCP stream，依赖 SSH 协议自身加密

### 12.2 访问控制

- 远程访问配置仅管理员可操作
- 客户只能连接自己被分配的机器（通过 AllocationService 校验）
- 会话数限制防止资源滥用
- 管理员可随时强制终止任意会话

### 12.3 审计追踪

- 所有管理员操作通过审计中间件自动记录
- 会话级事件通过 `remote_session_events` 表记录
- 支持按时间、用户、机器、操作类型查询审计日志

### 12.4 输入校验

- 域名格式校验（防止注入）
- 端口范围校验（1-65535）
- 协议白名单校验
- Nginx 配置生成使用模板引擎，避免拼接注入

---

## 13. 实现计划

### 13.1 阶段一：基础设施（P0）

| 任务 | 说明 |
|------|------|
| 数据库迁移脚本 | 创建 `18_remote_access.sql`，包含 4 张新表 |
| 实体定义 | `remote_access.go` 中定义 4 个实体 |
| DAO 层 | 4 个 DAO 文件，继承 BaseDao 模式 |
| 远程访问配置 API | GET/PUT 配置，管理员侧 |

### 13.2 阶段二：会话管理（P1）

| 任务 | 说明 |
|------|------|
| SessionService | 会话创建、断开、终止、心跳、清理 |
| 会话 API | 客户连接/断开，管理员列表/终止/事件查询 |
| 空闲清理定时任务 | 后台 goroutine 定期清理超时会话 |
| 审计中间件扩展 | `parseResource` 新增资源类型 |

### 13.3 阶段三：代理转发（P2）

| 任务 | 说明 |
|------|------|
| ProxyService | Nginx 配置生成、同步、reload |
| 配置模板 | HTTP/HTTPS 反向代理 + TCP stream 模板 |
| 同步 API | POST sync 触发配置同步 |
| 域名自动生成 | 基于 `base_domain` 自动生成机器域名 |

### 13.4 阶段四：Guacamole 集成（P2）

| 任务 | 说明 |
|------|------|
| GuacamoleClient | REST API 客户端封装 |
| 连接创建/销毁 | 会话创建时自动创建 Guacamole 连接 |
| Web 远程桌面 | VNC/RDP 通过 Guacamole 实现 |

---

## 14. 请求/响应结构体定义

文件路径：`backend/api/v1/remote_access.go`

```go
// UpsertRemoteAccessRequest 创建/更新远程访问配置请求
type UpsertRemoteAccessRequest struct {
    Enabled     bool   `json:"enabled"`
    Protocol    string `json:"protocol" binding:"required,oneof=ssh vnc rdp http https"`
    PublicDomain string `json:"public_domain" binding:"required"`
    PublicPort  int    `json:"public_port" binding:"required,min=1,max=65535"`
    TargetPort  int    `json:"target_port" binding:"min=0,max=65535"`
    ExtraPorts  string `json:"extra_ports"`
    AuthMode    string `json:"auth_mode" binding:"oneof=password key guacamole"`
    MaxSessions int    `json:"max_sessions" binding:"min=1,max=100"`
    IdleTimeout int    `json:"idle_timeout" binding:"min=60,max=86400"`
    Remark      string `json:"remark"`
}

// ConnectRequest 发起远程连接请求
type ConnectRequest struct {
    Protocol string `json:"protocol" binding:"required,oneof=ssh vnc rdp"`
}

// ConnectResponse 远程连接响应
type ConnectResponse struct {
    SessionID  string          `json:"session_id"`
    Protocol   string          `json:"protocol"`
    Connection *ConnectionInfo `json:"connection,omitempty"`
    Guacamole  *GuacamoleInfo  `json:"guacamole,omitempty"`
}

type ConnectionInfo struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
}

type GuacamoleInfo struct {
    URL   string `json:"url"`
    Token string `json:"token"`
}

// TerminateRequest 终止会话请求
type TerminateRequest struct {
    Reason string `json:"reason" binding:"required"`
}

// RemoteStatusResponse 远程访问状态响应
type RemoteStatusResponse struct {
    Enabled        bool            `json:"enabled"`
    Protocol       string          `json:"protocol"`
    ActiveSessions int             `json:"active_sessions"`
    MaxSessions    int             `json:"max_sessions"`
    ConnectionInfo *ConnectionInfo `json:"connection_info,omitempty"`
}
```

---

## 15. 数据库迁移脚本

文件路径：`backend/sql/18_remote_access.sql`

迁移脚本整合第 2 章中定义的 4 张表的 DDL，按依赖顺序执行：

```sql
-- 18_remote_access.sql
-- 远程访问模块数据库迁移

-- 1. 远程访问配置表
CREATE TABLE IF NOT EXISTS remote_access_configs ( ... );
-- （完整 DDL 见 2.2 节）

-- 2. 远程会话表
CREATE TABLE IF NOT EXISTS remote_sessions ( ... );
-- （完整 DDL 见 2.3 节）

-- 3. 会话事件表
CREATE TABLE IF NOT EXISTS remote_session_events ( ... );
-- （完整 DDL 见 2.4 节）

-- 4. 代理路由表
CREATE TABLE IF NOT EXISTS proxy_routes ( ... );
-- （完整 DDL 见 2.5 节）
```

---

## 16. 总结

本方案基于现有 RemoteGPU 后端架构，通过扩展式集成实现远程客户支持平台的后端功能。

**核心设计决策**：

| 决策 | 选择 | 理由 |
|------|------|------|
| 架构模式 | 扩展式（新增文件，最小修改） | 降低对现有功能的影响 |
| 认证方式 | 复用现有 JWT 体系 | 统一认证，减少维护成本 |
| 远程桌面 | Guacamole 网关 | 成熟方案，支持 SSH/VNC/RDP |
| 反向代理 | Nginx 配置动态生成 | 灵活、高性能 |
| 会话管理 | 数据库 + 定时清理 | 可靠、可审计 |
| 审计 | 复用中间件 + 会话事件表 | 两层审计，粒度可控 |

**新增代码量估算**：

| 模块 | 文件数 | 说明 |
|------|--------|------|
| Entity | 1 | 4 个实体定义 |
| DAO | 4 | 4 个数据访问对象 |
| Service | 3 | 配置/会话/代理 3 个服务 |
| Controller | 2 | 远程访问/会话 2 个控制器 |
| API 结构体 | 1 | 请求/响应定义 |
| Guacamole 客户端 | 1 | REST API 封装 |
| SQL 迁移 | 1 | 4 张表 DDL |
| **合计** | **13** | |
