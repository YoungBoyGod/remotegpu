# 远程客户支持平台 — 技术架构设计

> 版本：v1.0 | 日期：2026-02-06 | 作者：架构师

## 1. 文档概述

### 1.1 目的

本文档基于需求分析文档（`docs/design/requirements.md`）和已实现功能审查报告（`docs/design/implementation-review.md`），设计远程客户支持平台的整体技术架构，为后续后端、前端、运维方案提供架构指导。

### 1.2 读者

- 后端工程师
- 前端工程师
- 运维工程师
- 测试工程师
- 项目管理

### 1.3 设计原则

| 原则 | 说明 |
|------|------|
| 渐进增强 | 在现有架构基础上扩展，避免大规模重构 |
| 安全优先 | 远程访问场景安全性为第一优先级 |
| 松耦合 | 各组件通过明确接口交互，支持独立部署和升级 |
| 可观测 | 全链路日志、指标、审计，支持问题排查和合规审计 |
| 最小权限 | 每个组件和用户只拥有完成任务所需的最小权限 |

---

## 2. 系统架构总览

### 2.1 架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                          客户端层                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │ Web 浏览器 │  │ SSH 客户端│  │ VNC 客户端│  │ RDP 客户端       │   │
│  └─────┬────┘  └─────┬────┘  └─────┬────┘  └────────┬─────────┘   │
└────────┼─────────────┼─────────────┼────────────────┼──────────────┘
         │             │             │                │
         ▼             ▼             ▼                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        接入层 (Nginx / OpenResty)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │ HTTPS 终止    │  │ TCP 流代理    │  │ WebSocket 代理           │  │
│  │ (Web + API)  │  │ (SSH/RDP)    │  │ (Web Terminal / VNC)     │  │
│  └──────┬───────┘  └──────┬───────┘  └────────────┬─────────────┘  │
└─────────┼──────────────────┼──────────────────────┼────────────────┘
          │                  │                      │
          ▼                  ▼                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        网关层                                        │
│  ┌──────────────────────┐  ┌──────────────────────────────────────┐ │
│  │   API Gateway        │  │   Apache Guacamole                   │ │
│  │   (Backend API)      │  │   (远程访问网关)                      │ │
│  │   - 认证鉴权          │  │   - SSH/VNC/RDP 协议转换             │ │
│  │   - 路由分发          │  │   - Web 终端渲染                     │ │
│  │   - 限流熔断          │  │   - 会话录制                         │ │
│  └──────────┬───────────┘  └──────────────┬───────────────────────┘ │
└─────────────┼──────────────────────────────┼───────────────────────┘
              │                              │
              ▼                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        应用层                                        │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────────┐   │
│  │ 认证服务    │ │ 机器管理    │ │ 任务管理    │ │ 远程访问管理    │   │
│  │ AuthSvc    │ │ MachineSvc │ │ TaskSvc    │ │ RemoteAccSvc  │   │
│  └────────────┘ └────────────┘ └────────────┘ └────────────────┘   │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────────┐   │
│  │ 客户管理    │ │ 分配管理    │ │ 审计服务    │ │ 监控告警       │   │
│  │ CustSvc    │ │ AllocSvc   │ │ AuditSvc   │ │ MonitorSvc    │   │
│  └────────────┘ └────────────┘ └────────────┘ └────────────────┘   │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        数据层                                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │ PostgreSQL   │  │ Redis        │  │ 文件存储                  │  │
│  │ (主数据库)    │  │ (缓存/会话)   │  │ (数据集/日志)             │  │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        GPU 机器层                                    │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  Agent (Go)                                                  │   │
│  │  - 心跳上报  - 硬件采集  - 任务执行  - SSH 管理  - 监控上报   │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.2 分层职责

| 层级 | 职责 | 核心组件 |
|------|------|----------|
| 客户端层 | 用户交互入口 | Web 浏览器、SSH/VNC/RDP 原生客户端 |
| 接入层 | 流量接入、SSL 终止、协议路由 | Nginx / OpenResty |
| 网关层 | 认证鉴权、协议转换、会话管理 | Backend API、Apache Guacamole |
| 应用层 | 业务逻辑处理 | Go 微服务（单体部署） |
| 数据层 | 数据持久化与缓存 | PostgreSQL、Redis、文件存储 |
| GPU 机器层 | 资源执行与上报 | Agent（Go） |

---

## 3. 核心组件设计

### 3.1 接入层 — Nginx / OpenResty

#### 3.1.1 职责

- HTTPS/SSL 终止：统一管理 TLS 证书
- HTTP 反向代理：将 Web 和 API 请求转发到后端
- TCP 流代理：SSH（端口 2222）、RDP（端口 13389）等 TCP 协议转发
- WebSocket 代理：Web 终端和 VNC 的 WebSocket 长连接
- 静态资源服务：前端 SPA 静态文件

#### 3.1.2 配置架构

```
/etc/nginx/
├── nginx.conf                    # 主配置
├── conf.d/
│   ├── api.conf                  # API 反向代理
│   ├── frontend.conf             # 前端静态资源
│   └── guacamole.conf            # Guacamole WebSocket 代理
├── stream.d/
│   └── tcp-proxy.conf            # TCP 流代理（SSH/RDP）
└── remote-access.d/              # 动态生成的机器远程访问配置
    ├── gpu-001.conf
    ├── gpu-002.conf
    └── ...
```

#### 3.1.3 动态配置管理

后端通过 `RemoteAccessService` 生成 Nginx 配置片段：

1. 管理员配置机器远程访问参数（域名、端口、协议）
2. 后端生成对应的 Nginx server/upstream 配置
3. 写入 `remote-access.d/` 目录
4. 调用 `nginx -s reload` 热加载配置
5. 域名 + 端口唯一约束，防止冲突

### 3.2 远程访问网关 — Apache Guacamole

#### 3.2.1 选型理由

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| Apache Guacamole | 成熟稳定、支持 SSH/VNC/RDP、Web 化、会话录制 | Java 技术栈、需额外部署 | **采用** |
| WebSSH2 | 轻量、Node.js | 仅支持 SSH，无 VNC/RDP | 不采用 |
| noVNC + 自研 | 灵活 | 开发量大、仅 VNC | 不采用 |

#### 3.2.2 架构角色

```
浏览器 ──WebSocket──▶ Nginx ──▶ guacamole-client (Tomcat)
                                       │
                                       ▼
                               guacd (守护进程)
                                       │
                          ┌────────────┼────────────┐
                          ▼            ▼            ▼
                        SSH          VNC          RDP
                      (GPU机器)    (GPU机器)    (GPU机器)
```

#### 3.2.3 集成方式

- **认证对接**：通过 Guacamole REST API 或自定义 Auth Extension，与平台 JWT 打通
- **连接管理**：后端通过 Guacamole API 动态创建/销毁连接，无需手动配置
- **会话录制**：启用 `recording-path` 参数，录制文件存储到共享存储
- **权限控制**：后端校验用户对机器的访问权限后，生成临时连接 Token 传给前端

### 3.3 后端应用层 — Go + Gin

#### 3.3.1 现有架构（保持不变）

```
backend/internal/
├── controller/v1/     # 控制器层：请求解析、参数校验、响应封装
├── service/           # 服务层：业务逻辑
├── dao/               # 数据访问层：数据库操作
├── model/entity/      # 数据模型
├── middleware/         # 中间件：认证、审计、CORS
├── router/            # 路由定义
└── agent/             # Agent 通信客户端
```

#### 3.3.2 新增服务模块

| 模块 | 包路径 | 职责 |
|------|--------|------|
| RemoteAccessService | `service/remote_access/` | 远程访问配置管理、Nginx 配置生成 |
| GuacamoleService | `service/guacamole/` | Guacamole API 对接、连接生命周期管理 |
| SessionService | `service/session/` | 远程会话管理、会话审计 |
| NginxConfigService | `service/nginx/` | Nginx 配置文件生成与热加载 |

### 3.4 Agent 组件

#### 3.4.1 现有能力

- 心跳上报（HTTP 轮询）
- 硬件信息采集（CPU/内存/磁盘/GPU）
- 任务领取与执行（租约机制）
- 进程管理（启动/停止/状态查询）
- 本地 SQLite 存储（断线恢复）

#### 3.4.2 远程访问相关扩展

| 扩展项 | 说明 |
|--------|------|
| SSH 密钥注入 | 接收后端指令，将客户公钥写入 `~/.ssh/authorized_keys` |
| SSH 密码重置 | 通过 `chpasswd` 重置指定用户的 SSH 密码 |
| VNC 服务管理 | 启动/停止 VNC Server（如 TigerVNC） |
| 端口探测 | 检测指定端口的服务可用性，上报给后端 |
| 环境清理 | 机器回收时清理客户数据、重置环境 |

#### 3.4.3 Agent 通信协议

```
Agent ──HTTP POST──▶ Backend API
       /api/v1/agent/heartbeat     # 心跳上报
       /api/v1/agent/tasks/poll    # 任务轮询
       /api/v1/agent/tasks/:id     # 任务状态上报
       /api/v1/agent/ssh/inject    # SSH 密钥注入结果
       /api/v1/agent/ssh/reset     # SSH 密码重置结果

Backend ──HTTP──▶ Agent (反向调用)
       /api/v1/process/stop        # 停止进程
       /api/v1/ssh/inject-key      # 注入 SSH 密钥
       /api/v1/ssh/reset-password  # 重置 SSH 密码
       /api/v1/vnc/start           # 启动 VNC 服务
       /api/v1/vnc/stop            # 停止 VNC 服务
```

---

## 4. 技术选型

### 4.1 核心技术栈（已确定）

| 组件 | 技术 | 版本要求 | 说明 |
|------|------|----------|------|
| 后端框架 | Go + Gin | Go 1.21+ | 现有技术栈，保持不变 |
| ORM | GORM | v2 | 现有技术栈，保持不变 |
| 主数据库 | PostgreSQL | 14+ | 现有技术栈，保持不变 |
| 缓存/会话 | Redis | 7+ | 现有技术栈，保持不变 |
| 前端框架 | Vue 3 + TypeScript | Vue 3.3+ | 现有技术栈，保持不变 |
| UI 组件库 | Element Plus | 2.x | 现有技术栈，保持不变 |
| 构建工具 | Vite | 5.x | 现有技术栈，保持不变 |
| Agent | Go | Go 1.21+ | 现有技术栈，保持不变 |

### 4.2 新增技术组件

| 组件 | 技术 | 用途 | 部署方式 |
|------|------|------|----------|
| 反向代理 | Nginx / OpenResty | SSL 终止、TCP/HTTP 代理、静态资源 | Docker 或系统安装 |
| 远程访问网关 | Apache Guacamole 1.5+ | SSH/VNC/RDP Web 化访问 | Docker Compose |
| 监控采集 | Prometheus | GPU/系统指标采集与存储 | Docker |
| 监控面板 | Grafana | 监控数据可视化（可选） | Docker |
| 镜像仓库 | Harbor | Docker 镜像管理 | Docker Compose |
| VNC Server | TigerVNC | GPU 机器桌面服务 | 系统安装（Agent 管理） |

### 4.3 关键技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 远程桌面方案 | Guacamole（非自研） | 成熟稳定，支持多协议，社区活跃 |
| 配置管理方式 | 文件生成 + reload（非 API） | Nginx 原生支持，简单可靠 |
| Agent 通信 | HTTP 轮询（保持现有） | 已验证可用，后续可升级 gRPC |
| 会话录制存储 | 本地文件 + 定期归档 | 避免引入对象存储依赖 |
| 证书管理 | Let's Encrypt + certbot | 免费、自动续期 |

---

## 5. 安全架构

### 5.1 安全分层模型

```
┌─────────────────────────────────────────────┐
│  网络安全层                                   │
│  - 防火墙规则（仅开放必要端口）                 │
│  - DDoS 防护                                 │
│  - IP 白名单（可选）                           │
├─────────────────────────────────────────────┤
│  传输安全层                                   │
│  - TLS 1.2+ 全链路加密                        │
│  - 证书自动管理（Let's Encrypt）               │
│  - HSTS 强制 HTTPS                           │
├─────────────────────────────────────────────┤
│  认证安全层                                   │
│  - JWT 双令牌机制（access + refresh）          │
│  - 首次登录强制改密                            │
│  - 密码强度校验                               │
│  - Token 黑名单（Redis）                      │
├─────────────────────────────────────────────┤
│  授权安全层                                   │
│  - RBAC 角色权限控制                          │
│  - 资源归属校验（多租户隔离）                   │
│  - owner/member 权限区分                      │
├─────────────────────────────────────────────┤
│  审计安全层                                   │
│  - 操作审计日志                               │
│  - 远程会话录制                               │
│  - 登录/登出记录                              │
└─────────────────────────────────────────────┘
```

### 5.2 认证与鉴权架构

#### 5.2.1 JWT 认证流程（现有，需增强）

```
客户端                    后端                      Redis
  │                       │                         │
  │── POST /auth/login ──▶│                         │
  │                       │── 校验密码 ──▶           │
  │                       │── 生成 access_token ──▶  │
  │                       │── 生成 refresh_token ──▶ │── 存储 refresh_token
  │◀── {access, refresh} ─│                         │
  │                       │                         │
  │── GET /api (Bearer) ──▶│                         │
  │                       │── 校验 JWT 签名          │
  │                       │── 校验黑名单 ──────────▶ │── 查询黑名单
  │                       │── 校验用户状态           │
  │◀── 200 OK ────────────│                         │
```

#### 5.2.2 需增强的安全措施

| 措施 | 当前状态 | 改进方案 |
|------|----------|----------|
| JWT Claims 加租户标识 | 缺少 company_code | 在 Claims 中加入 `company_code` 字段 |
| owner/member 权限区分 | 未实现 | 添加 `RequireOwner()` 中间件 |
| 密码强度校验 | 改密未校验 | ChangePassword 调用 ValidateStrength |
| 登录频率限制 | 未实现 | Redis 滑动窗口限流（5次/分钟） |
| Agent 认证 | 简单 Token | 引入 Agent 专用 API Key + HMAC 签名 |

### 5.3 远程访问安全

#### 5.3.1 SSH 访问安全

| 安全措施 | 说明 |
|----------|------|
| 密钥认证优先 | 推荐使用 SSH 密钥登录，减少密码暴力破解风险 |
| 密码强度要求 | SSH 密码重置时强制校验密码强度 |
| 端口隔离 | 每台机器使用独立的公网端口映射，避免端口冲突 |
| 访问日志 | Nginx stream 模块记录 TCP 连接日志 |
| fail2ban | GPU 机器部署 fail2ban，防止 SSH 暴力破解 |

#### 5.3.2 Guacamole 访问安全

| 安全措施 | 说明 |
|----------|------|
| 临时 Token | 后端生成一次性连接 Token，有效期 60 秒 |
| 权限校验前置 | 后端校验用户对机器的访问权限后才创建 Guacamole 连接 |
| 会话超时 | 空闲会话自动断开（默认 30 分钟） |
| 并发限制 | 单用户单机器最多 1 个活跃会话 |
| 会话录制 | 所有 VNC/RDP 会话自动录制，支持回放审计 |

---

## 6. 网络拓扑

### 6.1 网络分区

```
┌──────────────────────────────────────────────────────────────┐
│                      公网 (Internet)                          │
│  客户浏览器 / SSH 客户端 / VNC 客户端                          │
└──────────────────────┬───────────────────────────────────────┘
                       │ 443 (HTTPS)
                       │ 2222 (SSH)
                       │ 13389 (RDP)
                       │ 15900 (VNC)
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                      DMZ 区                                   │
│  ┌─────────────┐  ┌──────────────┐                           │
│  │ Nginx       │  │ Guacamole    │                           │
│  │ (反向代理)   │  │ (远程网关)    │                           │
│  └──────┬──────┘  └──────┬───────┘                           │
└─────────┼────────────────┼───────────────────────────────────┘
          │                │
          ▼                ▼
┌──────────────────────────────────────────────────────────────┐
│                      应用区                                   │
│  ┌─────────────┐  ┌──────────┐  ┌──────────┐                │
│  │ Backend API │  │ Redis    │  │ PostgreSQL│                │
│  └──────┬──────┘  └──────────┘  └──────────┘                │
└─────────┼────────────────────────────────────────────────────┘
          │
          ▼
┌──────────────────────────────────────────────────────────────┐
│                      GPU 机器区                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│  │ GPU-001  │  │ GPU-002  │  │ GPU-00N  │                   │
│  │ (Agent)  │  │ (Agent)  │  │ (Agent)  │                   │
│  └──────────┘  └──────────┘  └──────────┘                   │
└──────────────────────────────────────────────────────────────┘
```

### 6.2 端口规划

| 端口 | 协议 | 用途 | 暴露范围 |
|------|------|------|----------|
| 443 | HTTPS | Web 前端 + API | 公网 |
| 80 | HTTP | 重定向到 HTTPS | 公网 |
| 2222 | TCP | SSH 代理入口 | 公网 |
| 13389 | TCP | RDP 代理入口 | 公网 |
| 15900 | TCP | VNC 代理入口 | 公网 |
| 8080 | HTTP | Backend API（内部） | 应用区 |
| 8443 | HTTPS | Guacamole Web | 应用区 |
| 4822 | TCP | guacd 守护进程 | 应用区 |
| 5432 | TCP | PostgreSQL | 应用区 |
| 6379 | TCP | Redis | 应用区 |
| 9090 | HTTP | Prometheus | 应用区 |

### 6.3 防火墙规则

| 源 | 目标 | 端口 | 说明 |
|----|------|------|------|
| 公网 | DMZ/Nginx | 80, 443 | Web 访问 |
| 公网 | DMZ/Nginx | 2222, 13389, 15900 | 远程访问代理 |
| DMZ/Nginx | 应用区/Backend | 8080 | API 转发 |
| DMZ/Nginx | 应用区/Guacamole | 8443 | 远程桌面转发 |
| 应用区/Backend | 应用区/PostgreSQL | 5432 | 数据库访问 |
| 应用区/Backend | 应用区/Redis | 6379 | 缓存访问 |
| 应用区/Backend | GPU机器区/Agent | 8090 | Agent 反向调用 |
| 应用区/Guacamole | GPU机器区 | 22, 5900, 3389 | 远程协议连接 |
| GPU机器区/Agent | 应用区/Backend | 8080 | 心跳/任务轮询 |

---

## 7. 数据流设计

### 7.1 远程 SSH 访问数据流

```
客户 SSH 客户端
  │
  │── TCP 连接 ──▶ Nginx (公网:2222)
  │                  │
  │                  │── stream proxy ──▶ GPU 机器 (内网:22)
  │                  │
  │◀── SSH 会话 ────▶│◀── SSH 会话 ────▶│
```

**流程说明**：
1. 客户使用 SSH 客户端连接 `<machine-slug>.<base-domain>:2222`
2. Nginx stream 模块根据 SNI 或端口映射转发到目标 GPU 机器的 22 端口
3. GPU 机器上的 SSHD 进行密钥/密码认证
4. Nginx 记录 TCP 连接日志（源 IP、目标、时长）

### 7.2 Web 远程桌面数据流（VNC/RDP）

```
客户浏览器
  │
  │── HTTPS ──▶ Nginx (公网:443)
  │                │
  │                │── /api/remote/connect ──▶ Backend API
  │                │                            │
  │                │                            │── 校验权限
  │                │                            │── 创建 Guacamole 连接
  │                │                            │── 返回连接 Token
  │                │
  │── WebSocket ──▶│── /guacamole/websocket ──▶ guacamole-client
  │                                               │
  │                                               │── guacd
  │                                               │    │
  │                                               │    ▼
  │                                               │  GPU 机器 (VNC:5900 / RDP:3389)
  │◀── 桌面画面流 ──────────────────────────────────│
```

**流程说明**：
1. 客户在 Web 界面点击"远程桌面"
2. 前端调用后端 API 请求建立远程连接
3. 后端校验用户对该机器的访问权限
4. 后端通过 Guacamole REST API 创建连接，返回临时 Token
5. 前端使用 Token 建立 WebSocket 连接到 Guacamole
6. Guacamole 通过 guacd 连接到 GPU 机器的 VNC/RDP 服务
7. 桌面画面通过 WebSocket 实时传输到浏览器

### 7.3 机器分配与远程访问配置数据流

```
管理员浏览器
  │
  │── POST /admin/allocations ──▶ Backend API
  │                                  │
  │                                  │── 1. 行级锁检查机器状态（SELECT FOR UPDATE）
  │                                  │── 2. 创建分配记录
  │                                  │── 3. 更新机器状态为 allocated
  │                                  │── 4. 异步：通知 Agent 注入 SSH 密钥
  │                                  │── 5. 异步：生成 Nginx 远程访问配置
  │                                  │── 6. 异步：触发 Nginx reload
  │                                  │── 7. 记录审计日志
  │◀── 200 OK ──────────────────────│
```

---

## 8. 数据模型设计

### 8.1 新增数据表

#### 8.1.1 远程访问配置表 `machine_remote_access`

```sql
CREATE TABLE machine_remote_access (
  id            SERIAL PRIMARY KEY,
  host_id       VARCHAR(64) NOT NULL REFERENCES hosts(id),
  enabled       BOOLEAN NOT NULL DEFAULT FALSE,
  protocol      VARCHAR(16) NOT NULL DEFAULT 'tcp',
  public_domain VARCHAR(255) NOT NULL,
  public_port   INT NOT NULL,
  target_port   INT,
  extra_ports   VARCHAR(255),
  remark        TEXT,
  nginx_synced  BOOLEAN NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMP DEFAULT NOW(),
  updated_at    TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_remote_access_domain_port
  ON machine_remote_access(public_domain, public_port);
```

#### 8.1.2 远程会话表 `remote_sessions`

```sql
CREATE TABLE remote_sessions (
  id              SERIAL PRIMARY KEY,
  host_id         VARCHAR(64) NOT NULL REFERENCES hosts(id),
  customer_id     INT NOT NULL REFERENCES customers(id),
  protocol        VARCHAR(16) NOT NULL,       -- ssh/vnc/rdp
  status          VARCHAR(16) NOT NULL DEFAULT 'active',  -- active/closed
  source_ip       VARCHAR(45),
  started_at      TIMESTAMP NOT NULL DEFAULT NOW(),
  ended_at        TIMESTAMP,
  duration_sec    INT,
  recording_path  VARCHAR(512),               -- 会话录制文件路径
  created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_remote_sessions_host ON remote_sessions(host_id);
CREATE INDEX idx_remote_sessions_customer ON remote_sessions(customer_id);
```

### 8.2 现有表改动

| 表 | 改动 | 说明 |
|----|------|------|
| `customers` | 新增 `company_code` 字段 | 多租户标识（已有迁移脚本） |
| `customers` | 新增 `must_change_password` 字段 | 首次登录改密标记（已有） |
| `hosts` | 关联 `machine_remote_access` | 一对一关系 |

---

## 9. 部署架构

### 9.1 部署拓扑

```
管理服务器（1台）
├── Docker Compose
│   ├── nginx          (反向代理 + SSL 终止)
│   ├── backend-api    (Go 后端服务)
│   ├── frontend       (Vue 静态资源，由 Nginx 托管)
│   ├── guacamole      (guacamole-client + guacd)
│   ├── postgresql     (主数据库)
│   ├── redis          (缓存/会话)
│   └── prometheus     (监控采集，可选)
│
GPU 机器（N台）
├── Agent (Go 二进制，systemd 管理)
├── SSHD (系统自带)
├── TigerVNC Server (按需启动)
└── NVIDIA Driver + CUDA
```

### 9.2 容器编排要点

| 服务 | 资源限制 | 持久化卷 | 健康检查 |
|------|----------|----------|----------|
| nginx | 512MB RAM | 配置目录、证书目录、日志 | `nginx -t` |
| backend-api | 1GB RAM | 日志、Nginx 配置输出目录 | HTTP `/health` |
| guacamole-client | 1GB RAM | 无 | HTTP `/guacamole/` |
| guacd | 512MB RAM | 录制文件目录 | TCP 4822 |
| postgresql | 2GB RAM | 数据目录 | `pg_isready` |
| redis | 512MB RAM | 数据目录 | `redis-cli ping` |

---

## 10. 实施路线图

### 10.1 阶段一：安全加固与基础完善（P0/P1）

**目标**：修复已知安全问题，完善权限模型

| 任务 | 优先级 | 依赖 |
|------|--------|------|
| 机器分配加行级锁 | P0 | 无 |
| 回收前检查运行中任务 | P0 | 无 |
| 任务停止加超时控制 | P1 | 无 |
| owner/member 权限区分 | P1 | 无 |
| 改密加密码强度校验 | P1 | 无 |
| JWT Claims 加 company_code | P1 | 无 |
| 登录频率限制 | P1 | Redis |

### 10.2 阶段二：远程访问核心功能（P2）

**目标**：实现 SSH 增强和 Web 远程桌面

| 任务 | 优先级 | 依赖 |
|------|--------|------|
| 远程访问配置后端 API | P2 | 阶段一完成 |
| Nginx 配置自动生成与热加载 | P2 | 远程访问 API |
| SSH 密钥自动注入（Agent） | P2 | Agent API |
| Guacamole 部署与集成 | P2 | Docker Compose |
| Web 远程桌面前端页面 | P2 | Guacamole 集成 |
| 远程会话管理与审计 | P2 | Guacamole 集成 |
| SSL 证书自动管理 | P2 | Nginx |

### 10.3 阶段三：监控增强与高级功能（P3）

**目标**：完善监控体系，实现高级任务管理

| 任务 | 优先级 | 依赖 |
|------|--------|------|
| Prometheus 监控接入 | P3 | Prometheus 部署 |
| GPU 利用率趋势图表 | P3 | Prometheus |
| Harbor 镜像同步 | P3 | Harbor 部署 |
| 任务编排（依赖/串并行） | P3 | 任务引擎设计 |
| 到期自动回收 | P3 | 定时任务框架 |
| DNS 自动管理 | P3 | DNS API |

---

## 11. 风险评估与应对

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| Guacamole 性能瓶颈 | 高并发远程桌面卡顿 | 中 | 水平扩展 guacd 实例，限制单节点并发会话数 |
| Nginx 配置冲突 | 远程访问不可用 | 低 | 配置生成前校验，reload 前 `nginx -t` 检查 |
| Agent 通信中断 | SSH 密钥注入/密码重置失败 | 中 | 异步重试机制，心跳超时告警 |
| GPU 机器网络不稳定 | 远程会话断开 | 中 | Guacamole 自动重连，Agent 断线恢复 |
| 证书过期 | HTTPS 访问失败 | 低 | certbot 自动续期 + 过期前 7 天告警 |
| 数据库并发冲突 | 机器重复分配 | 高 | 行级锁（SELECT FOR UPDATE）+ 唯一约束 |

---

## 12. 架构决策记录（ADR）

### ADR-001：远程桌面方案选择 Guacamole

- **背景**：需要支持 SSH/VNC/RDP 三种协议的 Web 化访问
- **决策**：采用 Apache Guacamole 作为远程访问网关
- **理由**：成熟稳定、多协议支持、内置会话录制、社区活跃
- **代价**：引入 Java 技术栈依赖，需额外运维 Tomcat + guacd

### ADR-002：Nginx 配置管理采用文件生成方式

- **背景**：需要动态管理每台机器的远程访问代理配置
- **决策**：后端生成 Nginx 配置文件 + `nginx -s reload` 热加载
- **备选**：OpenResty Lua 动态路由、Nginx Plus API
- **理由**：实现简单，Nginx 原生支持，无额外依赖
- **代价**：reload 有短暂性能抖动（毫秒级，可接受）

### ADR-003：保持单体部署架构

- **背景**：后端当前为单体 Go 应用，是否需要拆分微服务
- **决策**：保持单体部署，通过包级别模块化实现关注点分离
- **理由**：当前规模不需要微服务，单体部署运维简单，团队规模小
- **代价**：未来如需水平扩展特定模块，需要额外拆分工作

---

## 13. 总结

本架构设计遵循"渐进增强"原则，在现有 RemoteGPU 平台基础上扩展远程访问能力：

1. **接入层**：Nginx/OpenResty 统一处理 HTTPS、TCP 流代理和 WebSocket，支持动态配置热加载
2. **远程访问网关**：Apache Guacamole 提供 SSH/VNC/RDP 的 Web 化访问，内置会话录制
3. **应用层**：保持现有 Go + Gin 分层架构，新增 RemoteAccess、Guacamole、Session、NginxConfig 四个服务模块
4. **Agent 层**：扩展 SSH 密钥注入、密码重置、VNC 管理等远程访问相关能力
5. **安全架构**：五层安全模型（网络、传输、认证、授权、审计），修复已知安全问题
6. **数据层**：新增远程访问配置表和远程会话表，现有表结构小幅调整

实施分三个阶段推进：先修复安全问题和完善权限模型，再实现远程访问核心功能，最后接入监控和高级特性。
