# 远程客户支持平台 — 前端技术方案

> 版本：v1.0 | 日期：2026-02-06 | 作者：前端工程师

## 1. 概述

### 1.1 文档目的

本文档描述远程客户支持平台前端模块的技术方案，包括用户界面设计、远程访问客户端集成、会话管理、权限管理、审计日志查看，以及与现有 RemoteGPU 前端的集成策略。

### 1.2 技术栈（沿用现有）

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue 3 | 3.5.x | 前端框架（组合式 API） |
| TypeScript | 5.9.x | 类型安全 |
| Element Plus | 2.13.x | UI 组件库 |
| Vite | 7.3.x | 构建工具 |
| Pinia | 3.0.x | 状态管理 |
| Axios | 1.13.x | HTTP 客户端 |
| ECharts | 6.0.x | 数据可视化 |
| xterm.js | 5.x | Web 终端（新增） |
| noVNC | 1.x | Web VNC 客户端（新增） |

### 1.3 设计原则

- **渐进集成**：在现有 RemoteGPU 前端基础上扩展，不破坏已有功能
- **组件复用**：复用现有 DataTable、StatusTag、PageHeader 等通用组件
- **路由一致**：遵循现有 `/admin/*` 和 `/customer/*` 路由前缀规范
- **类型安全**：所有新增接口和数据结构使用 TypeScript interface 定义
- **权限前置**：路由守卫 + 组件级权限控制双重保障

---

## 2. 新增页面规划

### 2.1 管理员侧页面

| 页面 | 路由 | 说明 |
|------|------|------|
| 远程会话管理 | `/admin/remote/sessions` | 查看所有活跃远程会话，支持强制断开 |
| 远程访问配置 | `/admin/machines/:id/remote-access` | 机器远程访问配置（已有部分实现） |
| 会话审计日志 | `/admin/remote/audit` | 远程会话操作审计记录 |
| 远程访问策略 | `/admin/remote/policies` | 访问时段、IP 白名单、协议限制等策略配置 |

### 2.2 客户侧页面

| 页面 | 路由 | 说明 |
|------|------|------|
| Web 终端 | `/customer/machines/:id/terminal` | 基于 xterm.js 的 Web SSH 终端 |
| 远程桌面 | `/customer/machines/:id/desktop` | 基于 noVNC 的 Web VNC/RDP 桌面 |
| 我的会话 | `/customer/remote/sessions` | 查看个人活跃会话和历史记录 |
| 快速连接 | `/customer/machines/connect` | 已有菜单入口，实现快速选择机器并连接 |

---

## 3. 远程访问客户端集成

### 3.1 Web SSH 终端（xterm.js）

#### 3.1.1 技术选型

```
浏览器 (xterm.js) ←WebSocket→ 后端 WebSocket Proxy ←SSH→ GPU 机器
```

- **xterm.js**：成熟的 Web 终端模拟器，支持完整的 VT100/xterm 终端仿真
- **xterm-addon-fit**：自适应容器尺寸
- **xterm-addon-web-links**：URL 自动识别和点击
- **xterm-addon-search**：终端内搜索

#### 3.1.2 组件设计

```
src/components/remote/
├── WebTerminal.vue          # Web SSH 终端组件
├── TerminalToolbar.vue      # 终端工具栏（复制、粘贴、全屏、字体大小）
├── TerminalTabs.vue         # 多标签终端管理
└── ConnectionStatus.vue     # 连接状态指示器
```

#### 3.1.3 WebTerminal.vue 核心逻辑

```typescript
// 组件 props
interface WebTerminalProps {
  machine_id: string
  session_id?: string       // 恢复已有会话
  font_size?: number        // 默认 14
  theme?: 'dark' | 'light'  // 默认 dark
}

// WebSocket 连接流程
// 1. 调用 POST /api/v1/customer/remote/sessions 创建会话，获取 session_id + ws_url
// 2. 建立 WebSocket 连接到 ws_url
// 3. xterm.js onData → WebSocket send（用户输入）
// 4. WebSocket onMessage → xterm.js write（服务端输出）
// 5. 窗口 resize → WebSocket send resize 事件
// 6. 断线重连：指数退避，最多重试 5 次
```

#### 3.1.4 终端功能清单

| 功能 | 说明 |
|------|------|
| 基础终端 | 完整的 xterm 终端仿真，支持颜色、光标移动 |
| 自适应尺寸 | 窗口 resize 时自动调整行列数 |
| 复制粘贴 | Ctrl+Shift+C/V 或右键菜单 |
| 全屏模式 | F11 或工具栏按钮切换全屏 |
| 字体调整 | Ctrl+加号/减号 调整字体大小 |
| 搜索 | Ctrl+Shift+F 终端内搜索 |
| 多标签 | 同时打开多个终端标签页 |
| 断线重连 | 网络断开后自动重连，保持会话 |
| 连接状态 | 实时显示连接状态（连接中/已连接/断开/重连中） |

### 3.2 Web VNC/RDP 远程桌面（noVNC + Guacamole）

#### 3.2.1 架构方案

```
方案 A（推荐）：Guacamole 统一网关
浏览器 ←WebSocket→ guacd ←VNC/RDP→ GPU 机器

方案 B：noVNC 直连
浏览器 (noVNC) ←WebSocket→ websockify ←VNC→ GPU 机器
```

推荐方案 A，原因：
- Guacamole 同时支持 VNC、RDP、SSH，统一网关
- 内置会话录制、剪贴板共享、文件传输
- 后端只需管理 Guacamole 连接参数，前端通过 guacamole-common-js 接入

#### 3.2.2 组件设计

```
src/components/remote/
├── RemoteDesktop.vue        # 远程桌面主组件
├── DesktopToolbar.vue       # 桌面工具栏（全屏、剪贴板、Ctrl+Alt+Del）
├── ClipboardSync.vue        # 剪贴板同步面板
└── FileTransfer.vue         # 文件传输面板（Guacamole 支持）
```

#### 3.2.3 RemoteDesktop.vue 核心逻辑

```typescript
// 组件 props
interface RemoteDesktopProps {
  machine_id: string
  protocol: 'vnc' | 'rdp'
  session_id?: string
}

// Guacamole 连接流程
// 1. 调用 POST /api/v1/customer/remote/sessions 创建会话
//    请求体包含 protocol: 'vnc' | 'rdp'
//    返回 guacamole_token + tunnel_url
// 2. 使用 guacamole-common-js 建立 Guacamole.WebSocketTunnel
// 3. 创建 Guacamole.Client，绑定到 display 容器
// 4. 处理键盘/鼠标事件转发
// 5. 支持剪贴板双向同步
```

#### 3.2.4 远程桌面功能清单

| 功能 | 说明 |
|------|------|
| VNC 桌面 | 通过 Guacamole 访问 VNC 桌面 |
| RDP 桌面 | 通过 Guacamole 访问 Windows RDP |
| 全屏模式 | 全屏显示远程桌面 |
| 自适应缩放 | 桌面分辨率自适应浏览器窗口 |
| 剪贴板同步 | 本地与远程剪贴板双向同步 |
| Ctrl+Alt+Del | 发送特殊按键组合 |
| 文件传输 | 通过 Guacamole SFTP 上传/下载文件 |
| 连接质量指示 | 显示延迟和帧率 |

### 3.3 连接入口设计

客户在机器详情页和机器列表页均可发起远程连接：

```
机器详情页 → "远程连接" 按钮组
  ├── SSH 终端（Web Terminal）
  ├── VNC 桌面（需机器支持）
  ├── RDP 桌面（需机器支持）
  └── 复制 SSH 命令（已有）

机器列表页 → 操作列 "连接" 下拉菜单
  ├── Web 终端
  ├── 远程桌面
  └── 复制连接信息
```

---

## 4. 会话管理界面

### 4.1 数据模型（前端类型定义）

```typescript
// src/types/remoteSession.ts

interface RemoteSession {
  id: string
  machine_id: string
  machine_name: string
  user_id: number
  username: string
  protocol: 'ssh' | 'vnc' | 'rdp'
  status: 'connecting' | 'active' | 'disconnected' | 'terminated'
  client_ip: string
  started_at: string
  ended_at?: string
  duration_seconds?: number
  bytes_sent?: number
  bytes_received?: number
}

interface SessionListRequest {
  page: number
  page_size: number
  status?: string
  protocol?: string
  machine_id?: string
  user_id?: number
  start_time?: string
  end_time?: string
}

interface SessionListResponse {
  total: number
  items: RemoteSession[]
}
```

### 4.2 管理员会话管理页面

**路由**：`/admin/remote/sessions`

**功能**：
- 活跃会话列表：实时展示所有在线远程会话
- 筛选条件：协议类型、机器、客户、状态
- 强制断开：管理员可强制终止任意会话
- 会话详情：点击查看会话的详细信息（连接时间、流量、客户端 IP）
- 历史记录：切换标签页查看已结束的会话记录

**表格列定义**：

| 列 | 字段 | 说明 |
|----|------|------|
| 会话 ID | id | 短 ID 展示 |
| 机器 | machine_name | 关联机器名称 |
| 用户 | username | 连接用户 |
| 协议 | protocol | SSH/VNC/RDP 标签 |
| 状态 | status | StatusTag 组件展示 |
| 客户端 IP | client_ip | 连接来源 |
| 开始时间 | started_at | 本地时间格式 |
| 持续时间 | duration | 实时计算 |
| 操作 | — | 断开、查看详情 |

### 4.3 客户会话页面

**路由**：`/customer/remote/sessions`

**功能**：
- 我的活跃会话：展示当前用户的在线会话
- 快速恢复：点击活跃会话可恢复连接
- 历史记录：查看个人历史会话
- 会话统计：本月连接次数、总时长

---

## 5. 权限管理界面

### 5.1 远程访问策略配置

**路由**：`/admin/remote/policies`（仅管理员）

#### 5.1.1 数据模型

```typescript
// src/types/remotePolicy.ts

interface RemoteAccessPolicy {
  id: number
  name: string
  description: string
  enabled: boolean
  // 适用范围
  scope_type: 'global' | 'customer' | 'machine'
  scope_id?: number
  // 访问控制
  allowed_protocols: ('ssh' | 'vnc' | 'rdp')[]
  allowed_time_ranges?: TimeRange[]    // 允许访问的时间段
  ip_whitelist?: string[]              // IP 白名单（CIDR 格式）
  max_concurrent_sessions?: number     // 最大并发会话数
  session_timeout_minutes?: number     // 会话超时时间
  idle_timeout_minutes?: number        // 空闲超时时间
  created_at: string
  updated_at: string
}

interface TimeRange {
  start_time: string   // "09:00"
  end_time: string     // "18:00"
  weekdays: number[]   // [1,2,3,4,5] 周一到周五
}
```

#### 5.1.2 页面布局

```
┌─────────────────────────────────────────────┐
│ 远程访问策略                    [+ 新建策略]  │
├─────────────────────────────────────────────┤
│ 筛选：[范围类型 ▼] [状态 ▼]    [搜索...]    │
├─────────────────────────────────────────────┤
│ 策略名称 | 范围 | 协议 | 并发限制 | 状态 | 操作│
│ ─────────────────────────────────────────── │
│ 默认策略  | 全局 | SSH  | 5       | 启用 | 编辑│
│ VIP客户   | 客户 | 全部 | 10      | 启用 | 编辑│
│ 维护窗口  | 全局 | SSH  | 2       | 禁用 | 编辑│
└─────────────────────────────────────────────┘
```

### 5.2 前端权限控制

#### 5.2.1 路由级权限

在现有路由守卫基础上扩展，新增远程访问相关路由的权限检查：

```typescript
// 管理员远程管理路由
{
  path: 'remote/sessions',
  name: 'admin-remote-sessions',
  component: () => import('@/views/admin/RemoteSessionListView.vue'),
  meta: { title: '远程会话管理', requiresRole: 'admin' }
}

// 客户远程连接路由
{
  path: 'machines/:id/terminal',
  name: 'customer-machine-terminal',
  component: () => import('@/views/customer/WebTerminalView.vue'),
  meta: { title: 'Web 终端', requiresRole: ['customer_owner', 'customer_member'] }
}
```

#### 5.2.2 组件级权限

```typescript
// 连接按钮根据机器分配状态和协议支持情况动态显示
// 仅当机器状态为 allocated 且分配给当前客户时显示连接按钮
// VNC/RDP 按钮仅在机器配置了对应协议时显示
```

---

## 6. 审计日志查看

### 6.1 远程会话审计

在现有审计日志页面（`/admin/audit`）基础上扩展，新增远程会话相关的审计事件类型：

#### 6.1.1 新增审计事件类型

| 事件类型 | 说明 |
|----------|------|
| `remote_session.created` | 创建远程会话 |
| `remote_session.connected` | 会话连接成功 |
| `remote_session.disconnected` | 会话断开 |
| `remote_session.terminated` | 管理员强制终止会话 |
| `remote_session.timeout` | 会话超时断开 |
| `remote_access.config_updated` | 远程访问配置变更 |
| `remote_policy.created` | 创建访问策略 |
| `remote_policy.updated` | 更新访问策略 |

#### 6.1.2 审计日志筛选扩展

在现有审计日志页面的筛选栏中，资源类型下拉框新增：
- `remote_session` — 远程会话
- `remote_policy` — 访问策略
- `remote_access` — 远程访问配置

#### 6.1.3 会话回放（P3 远期）

Guacamole 支持会话录制，后续可在审计日志详情中嵌入会话回放播放器：
- 管理员点击审计记录 → 查看会话详情 → 播放会话录像
- 使用 guacamole-common-js 的 SessionRecording 组件

---

## 7. 与现有前端的集成策略

### 7.1 目录结构扩展

```
frontend/src/
├── api/
│   ├── admin.ts              # 扩展：远程会话管理、策略管理 API
│   ├── customer.ts           # 扩展：创建会话、获取连接信息 API
│   └── remote.ts             # 新增：远程访问专用 API 模块
├── components/
│   └── remote/               # 新增：远程访问组件目录
│       ├── WebTerminal.vue
│       ├── TerminalToolbar.vue
│       ├── TerminalTabs.vue
│       ├── RemoteDesktop.vue
│       ├── DesktopToolbar.vue
│       ├── ClipboardSync.vue
│       ├── FileTransfer.vue
│       └── ConnectionStatus.vue
├── types/
│   ├── remoteSession.ts      # 新增：会话类型定义
│   └── remotePolicy.ts       # 新增：策略类型定义
├── composables/
│   ├── useWebSocket.ts       # 新增：WebSocket 连接管理
│   ├── useTerminal.ts        # 新增：终端实例管理
│   └── useRemoteDesktop.ts   # 新增：远程桌面连接管理
├── stores/
│   └── remote.ts             # 新增：远程会话状态管理
└── views/
    ├── admin/
    │   ├── RemoteSessionListView.vue   # 新增
    │   └── RemotePolicyView.vue        # 新增
    └── customer/
        ├── WebTerminalView.vue         # 新增
        ├── RemoteDesktopView.vue       # 新增
        └── RemoteSessionListView.vue   # 新增
```

### 7.2 侧边栏菜单扩展

#### 管理员侧边栏（AdminSidebar.vue）

在现有菜单中新增"远程管理"分组：

```typescript
{
  id: 'remote',
  title: '远程管理',
  icon: Connection,  // 复用已导入的图标
  children: [
    { id: 'remote-sessions', title: '会话管理', path: '/admin/remote/sessions' },
    { id: 'remote-policies', title: '访问策略', path: '/admin/remote/policies' },
    { id: 'remote-audit', title: '会话审计', path: '/admin/remote/audit' }
  ]
}
```

#### 客户侧边栏（CustomerSidebar.vue）

在"我的机器"子菜单中新增连接入口：

```typescript
// 在 machines children 中追加
{ id: 'my-sessions', title: '我的会话', path: '/customer/remote/sessions', icon: '🖥️' }
```

### 7.3 路由注册

在 `router/index.ts` 中，在通配路由 `:pathMatch(.*)*` 之前注册新路由：

```typescript
// 管理员远程管理路由
{
  path: 'remote/sessions',
  name: 'admin-remote-sessions',
  component: () => import('@/views/admin/RemoteSessionListView.vue'),
  meta: { title: '远程会话管理' }
},
{
  path: 'remote/policies',
  name: 'admin-remote-policies',
  component: () => import('@/views/admin/RemotePolicyView.vue'),
  meta: { title: '远程访问策略' }
},

// 客户远程连接路由
{
  path: 'machines/:id/terminal',
  name: 'customer-machine-terminal',
  component: () => import('@/views/customer/WebTerminalView.vue'),
  meta: { title: 'Web 终端' }
},
{
  path: 'machines/:id/desktop',
  name: 'customer-machine-desktop',
  component: () => import('@/views/customer/RemoteDesktopView.vue'),
  meta: { title: '远程桌面' }
},
{
  path: 'remote/sessions',
  name: 'customer-remote-sessions',
  component: () => import('@/views/customer/RemoteSessionListView.vue'),
  meta: { title: '我的会话' }
}
```

### 7.4 新增依赖

```json
{
  "dependencies": {
    "@xterm/xterm": "^5.5.0",
    "@xterm/addon-fit": "^0.10.0",
    "@xterm/addon-web-links": "^0.11.0",
    "@xterm/addon-search": "^0.15.0",
    "@nicedoc/guacamole-common-js": "^1.5.0"
  }
}
```

### 7.5 API 模块扩展

```typescript
// src/api/remote.ts — 新增远程访问 API 模块

import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types/common'
import type { RemoteSession, SessionListRequest } from '@/types/remoteSession'
import type { RemoteAccessPolicy } from '@/types/remotePolicy'

// 会话管理
export function createSession(data: {
  machine_id: string
  protocol: 'ssh' | 'vnc' | 'rdp'
}): Promise<ApiResponse<{ session_id: string; ws_url: string }>> {
  return request.post('/customer/remote/sessions', data)
}

export function getSessionList(
  params: SessionListRequest
): Promise<ApiResponse<PageResponse<RemoteSession>>> {
  return request.get('/admin/remote/sessions', { params })
}

export function terminateSession(
  sessionId: string
): Promise<ApiResponse<null>> {
  return request.post(`/admin/remote/sessions/${sessionId}/terminate`)
}

// 策略管理
export function getPolicyList(): Promise<ApiResponse<RemoteAccessPolicy[]>> {
  return request.get('/admin/remote/policies')
}

export function createPolicy(
  data: Partial<RemoteAccessPolicy>
): Promise<ApiResponse<RemoteAccessPolicy>> {
  return request.post('/admin/remote/policies', data)
}

export function updatePolicy(
  id: number,
  data: Partial<RemoteAccessPolicy>
): Promise<ApiResponse<RemoteAccessPolicy>> {
  return request.put(`/admin/remote/policies/${id}`, data)
}
```

---

## 8. 状态管理设计

### 8.1 远程会话 Store

```typescript
// src/stores/remote.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { RemoteSession } from '@/types/remoteSession'

export const useRemoteStore = defineStore('remote', () => {
  // 活跃会话列表
  const activeSessions = ref<RemoteSession[]>([])

  // 当前活跃会话数
  const activeCount = computed(() =>
    activeSessions.value.filter(s => s.status === 'active').length
  )

  // 添加会话
  function addSession(session: RemoteSession) {
    activeSessions.value.push(session)
  }

  // 移除会话
  function removeSession(sessionId: string) {
    const index = activeSessions.value.findIndex(s => s.id === sessionId)
    if (index > -1) activeSessions.value.splice(index, 1)
  }

  // 更新会话状态
  function updateSessionStatus(sessionId: string, status: RemoteSession['status']) {
    const session = activeSessions.value.find(s => s.id === sessionId)
    if (session) session.status = status
  }

  return { activeSessions, activeCount, addSession, removeSession, updateSessionStatus }
})
```

---

## 9. WebSocket 连接管理

### 9.1 composable 设计

```typescript
// src/composables/useWebSocket.ts

interface UseWebSocketOptions {
  url: string
  protocols?: string[]
  reconnect?: boolean
  maxRetries?: number
  retryDelay?: number  // 初始重试延迟（ms）
  onOpen?: () => void
  onMessage?: (data: ArrayBuffer | string) => void
  onClose?: (event: CloseEvent) => void
  onError?: (event: Event) => void
}

// 返回值
interface UseWebSocketReturn {
  ws: Ref<WebSocket | null>
  status: Ref<'connecting' | 'connected' | 'disconnected' | 'reconnecting'>
  send: (data: string | ArrayBuffer) => void
  close: () => void
  reconnect: () => void
}
```

### 9.2 重连策略

- 指数退避：初始 1s，最大 30s
- 最多重试 5 次
- 重连时显示状态提示
- 用户可手动触发重连
- 页面不可见时暂停重连（visibilitychange 事件）

---

## 10. 实施计划

### 10.1 阶段划分

**阶段一（P1）：Web SSH 终端**
- 实现 WebTerminal 组件和 WebSocket 连接管理
- 客户机器详情页集成终端入口
- 基础会话管理（创建、列表、断开）
- 管理员会话列表页面

**阶段二（P2）：远程桌面 + 策略管理**
- 集成 Guacamole，实现 VNC/RDP 远程桌面
- 远程访问策略配置页面
- 会话审计日志扩展
- 剪贴板同步和文件传输

**阶段三（P3）：增强功能**
- 会话录制与回放
- 多标签终端管理
- 连接质量监控
- 移动端适配

### 10.2 依赖关系

```
后端 WebSocket Proxy 就绪 → 阶段一可开始
Guacamole 部署就绪 → 阶段二可开始
会话录制存储方案就绪 → 阶段三可开始
```

---

## 11. 安全考虑

| 安全项 | 措施 |
|--------|------|
| WebSocket 认证 | 连接时携带 JWT Token，服务端验证后建立连接 |
| 会话隔离 | 客户只能访问自己分配的机器会话 |
| 会话超时 | 空闲超时自动断开，防止资源占用 |
| 操作审计 | 所有远程连接操作记录审计日志 |
| IP 限制 | 支持 IP 白名单策略 |
| 协议限制 | 按策略限制可用的远程协议 |
| 并发控制 | 限制单用户/单机器的最大并发会话数 |
| XSS 防护 | 终端输出内容不渲染为 HTML，xterm.js 天然防护 |
| CSRF 防护 | WebSocket 握手阶段验证 Origin |
