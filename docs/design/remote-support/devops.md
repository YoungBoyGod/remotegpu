# 远程客户支持平台 — 运维部署方案

> 版本：v1.0 | 日期：2026-02-06 | 作者：运维工程师

## 目录

1. [部署架构](#1-部署架构)
2. [容器化方案](#2-容器化方案)
3. [监控告警](#3-监控告警)
4. [日志收集](#4-日志收集)
5. [备份恢复](#5-备份恢复)
6. [扩展性方案](#6-扩展性方案)
7. [运维流程](#7-运维流程)

---

## 1. 部署架构

### 1.1 整体架构概览

```
                        ┌─────────────┐
                        │   客户端     │
                        │ (浏览器/SSH) │
                        └──────┬──────┘
                               │
                        ┌──────▼──────┐
                        │   Nginx     │
                        │ (反向代理)   │
                        │ SSL 终止    │
                        └──────┬──────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
       ┌──────▼──────┐ ┌──────▼──────┐ ┌───────▼──────┐
       │  Frontend   │ │  Backend    │ │  Guacamole   │
       │  (Vue SPA)  │ │  (Go API)  │ │  (远程访问)   │
       └─────────────┘ └──────┬──────┘ └──────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
       ┌──────▼──────┐ ┌─────▼──────┐ ┌──────▼──────┐
       │ PostgreSQL  │ │   Redis    │ │   RustFS    │
       │  (主数据库)  │ │  (缓存)    │ │ (对象存储)   │
       └─────────────┘ └────────────┘ └─────────────┘

       ┌─────────────────────────────────────────────┐
       │              GPU 机器集群                     │
       │  ┌─────────┐ ┌─────────┐ ┌─────────┐       │
       │  │ Agent   │ │ Agent   │ │ Agent   │       │
       │  │ + GPU   │ │ + GPU   │ │ + GPU   │       │
       │  └─────────┘ └─────────┘ └─────────┘       │
       └─────────────────────────────────────────────┘
```

### 1.2 节点规划

| 节点类型 | 数量 | 最低配置 | 推荐配置 | 用途 |
|---------|------|---------|---------|------|
| 管理节点 | 1-2 | 4C/8G/100G | 8C/16G/200G SSD | 运行平台服务（Backend、Frontend、中间件） |
| 数据库节点 | 1-2 | 4C/16G/500G SSD | 8C/32G/1T SSD | PostgreSQL + Redis |
| 监控节点 | 1 | 4C/8G/200G | 8C/16G/500G SSD | Prometheus + Grafana + 日志收集 |
| GPU 计算节点 | N | 按需 | 按需 | Agent + GPU 工作负载 |

### 1.3 网络架构

```
┌─────────────────────────────────────────────────────┐
│                    公网区域                           │
│  客户端 ──── DNS ──── 公网 IP / 负载均衡              │
└────────────────────────┬────────────────────────────┘
                         │ 80/443
┌────────────────────────▼────────────────────────────┐
│                    DMZ 区域                           │
│  Nginx (反向代理 + SSL 终止 + WAF)                    │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│                  应用服务区域                          │
│  Backend API | Frontend | Guacamole | Harbor         │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│                  数据服务区域                          │
│  PostgreSQL | Redis | RustFS | etcd                  │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│                  GPU 计算区域                          │
│  Agent + GPU 工作负载（内网访问）                      │
└─────────────────────────────────────────────────────┘
```

### 1.4 端口规划

| 服务 | 内部端口 | 外部映射端口 | 协议 | 说明 |
|------|---------|-------------|------|------|
| Nginx | 80/443 | 80/443 | HTTP/HTTPS | 统一入口 |
| Backend API | 8080 | - | HTTP | 通过 Nginx 代理 |
| Frontend | 3000 | - | HTTP | 通过 Nginx 代理 |
| PostgreSQL | 5432 | 5432（仅内网） | TCP | 主数据库 |
| Redis | 6379 | 6379（仅内网） | TCP | 缓存/Token 存储 |
| Prometheus | 9090 | 19090（仅内网） | HTTP | 监控指标 |
| Grafana | 3000 | 13000 | HTTP | 监控面板 |
| Guacamole | 8080 | 8081 | HTTP | 远程访问网关 |
| Harbor | 8080 | 8082 | HTTP | 镜像仓库 |
| RustFS | 9000/9001 | 9000/9001 | HTTP | 对象存储 |
| Uptime Kuma | 3001 | 13001 | HTTP | 可用性监控 |
| etcd | 2379/2380 | 2379/2380（仅内网） | HTTP | 服务发现 |

---

## 2. 容器化方案

### 2.1 容器化策略

所有平台服务均通过 Docker 容器化部署，使用 Docker Compose 编排。GPU 计算节点上的 Agent 以二进制方式直接部署（需要访问宿主机 GPU 驱动和硬件信息）。

| 组件 | 容器化方式 | 说明 |
|------|-----------|------|
| Backend API | Docker | 多阶段构建，Go 静态编译 |
| Frontend | Docker (Nginx) | 构建产物通过 Nginx 提供静态文件服务 |
| PostgreSQL | Docker | 官方镜像 postgres:17 |
| Redis | Docker | 官方镜像 redis:8.4.0-alpine |
| Nginx | Docker | 官方镜像 nginx:1.24-alpine |
| Prometheus | Docker | 官方镜像 prom/prometheus:v2.48.0 |
| Grafana | Docker | 官方镜像 grafana/grafana:10.2.0 |
| Guacamole | Docker | 官方镜像 guacamole/guacamole:1.5.4 |
| Harbor | Docker | 官方镜像 goharbor/*:v2.9.0 |
| RustFS | Docker | 官方镜像 rustfs/rustfs:latest |
| Uptime Kuma | Docker | 官方镜像 louislam/uptime-kuma:1 |
| etcd | Docker | 官方镜像 quay.io/coreos/etcd:v3.5.13 |
| Agent | 二进制 | 直接部署在 GPU 宿主机上 |

### 2.2 Backend Dockerfile（多阶段构建）

```dockerfile
# 构建阶段
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/remotegpu ./cmd

# 运行阶段
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/bin/remotegpu /usr/local/bin/remotegpu
EXPOSE 8080
ENTRYPOINT ["remotegpu", "server"]
```

### 2.3 Frontend Dockerfile（多阶段构建）

```dockerfile
# 构建阶段
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# 运行阶段
FROM nginx:1.24-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

### 2.4 统一 Docker Compose 编排

生产环境使用统一的 `docker-compose.prod.yml` 编排所有管理节点服务：

```yaml
version: "3.8"

services:
  nginx:
    image: nginx:1.24-alpine
    container_name: remotegpu-nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
      - nginx_logs:/var/log/nginx
    depends_on:
      backend:
        condition: service_healthy
      frontend:
        condition: service_started
    networks:
      - remotegpu-network

  backend:
    image: ${REGISTRY}/remotegpu-backend:${VERSION:-latest}
    container_name: remotegpu-backend
    restart: always
    env_file: .env.backend
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 15s
      timeout: 5s
      retries: 3
    depends_on:
      postgresql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - remotegpu-network

  frontend:
    image: ${REGISTRY}/remotegpu-frontend:${VERSION:-latest}
    container_name: remotegpu-frontend
    restart: always
    networks:
      - remotegpu-network

  postgresql:
    image: postgres:17
    container_name: remotegpu-postgresql
    restart: always
    env_file: .env.db
    volumes:
      - postgresql_data:/var/lib/postgresql/data
      - ./postgresql/postgresql.conf:/etc/postgresql/postgresql.conf
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $POSTGRES_USER -d $POSTGRES_DB"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - remotegpu-network

  redis:
    image: redis:8.4.0-alpine
    container_name: remotegpu-redis
    restart: always
    command: redis-server /usr/local/etc/redis/redis.conf
    volumes:
      - redis_data:/data
      - ./redis/redis.conf:/usr/local/etc/redis/redis.conf
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
    networks:
      - remotegpu-network

volumes:
  postgresql_data:
  redis_data:
  nginx_logs:

networks:
  remotegpu-network:
    driver: bridge
```

### 2.5 Agent 部署方案

Agent 以 systemd 服务方式部署在每台 GPU 宿主机上：

```ini
# /etc/systemd/system/remotegpu-agent.service
[Unit]
Description=RemoteGPU Agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/remotegpu-agent
Restart=always
RestartSec=10
EnvironmentFile=/etc/remotegpu/agent.env
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

Agent 环境变量配置（`/etc/remotegpu/agent.env`）：

```bash
REMOTEGPU_SERVER_URL=https://api.remotegpu.example.com
REMOTEGPU_AGENT_TOKEN=<agent-auth-token>
REMOTEGPU_LOG_LEVEL=info
REMOTEGPU_DATA_DIR=/var/lib/remotegpu
```

### 2.6 环境变量管理

生产环境使用 `.env.*` 文件管理敏感配置，**禁止提交到 Git 仓库**。

| 文件 | 用途 |
|------|------|
| `.env.backend` | Backend API 配置（数据库连接、JWT 密钥、Redis 地址等） |
| `.env.db` | PostgreSQL 配置（用户名、密码、数据库名） |
| `.env.redis` | Redis 配置（密码） |
| `.env.harbor` | Harbor 配置（密钥、数据库密码） |

`.env.backend` 示例：

```bash
# 数据库
DB_HOST=postgresql
DB_PORT=5432
DB_USER=remotegpu_user
DB_PASSWORD=<strong-password>
DB_NAME=remotegpu
DB_SSLMODE=disable

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=<redis-password>

# JWT
JWT_SECRET=<jwt-secret-key>
JWT_ACCESS_EXPIRE=15m
JWT_REFRESH_EXPIRE=7d

# 服务
SERVER_PORT=8080
GIN_MODE=release
LOG_LEVEL=info
```

### 2.7 镜像构建与发布流程

```
代码提交 → CI 触发构建 → 多阶段 Docker 构建 → 推送到 Harbor 私有仓库
                                                    ↓
                                            管理节点拉取镜像
                                                    ↓
                                          docker compose up -d
```

镜像标签规范：
- `latest` — 最新开发版本
- `v1.0.0` — 正式发布版本（语义化版本号）
- `main-<commit-sha>` — 主分支构建版本

---

## 3. 监控告警

### 3.1 监控体系架构

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ GPU 节点     │  │ 管理节点      │  │ 数据库节点    │
│ node_exporter│  │ node_exporter│  │ node_exporter│
│ nvidia_exporter│ │ cadvisor    │  │ pg_exporter  │
│ agent metrics│  │              │  │ redis_exporter│
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┼─────────────────┘
                         │
                  ┌──────▼───────┐
                  │  Prometheus  │
                  │  (指标采集)   │
                  └──────┬───────┘
                         │
              ┌──────────┼──────────┐
              │          │          │
       ┌──────▼──────┐ ┌▼────────┐ ┌▼───────────┐
       │  Grafana    │ │Alertmgr │ │Uptime Kuma │
       │  (可视化)    │ │(告警)    │ │(可用性监控)  │
       └─────────────┘ └─────────┘ └────────────┘
```

### 3.2 监控指标采集

#### 3.2.1 基础设施指标

| Exporter | 部署位置 | 采集指标 |
|----------|---------|---------|
| node_exporter | 所有节点 | CPU、内存、磁盘、网络、系统负载 |
| nvidia_gpu_exporter | GPU 节点 | GPU 利用率、显存、温度、功耗、风扇转速 |
| cadvisor | 管理节点 | 容器 CPU/内存/网络/磁盘 I/O |
| postgres_exporter | 数据库节点 | 连接数、查询延迟、锁等待、缓存命中率 |
| redis_exporter | 数据库节点 | 内存使用、连接数、命中率、键数量 |

#### 3.2.2 应用指标

| 指标来源 | 采集方式 | 关键指标 |
|---------|---------|---------|
| Backend API | `/metrics` 端点（Prometheus 格式） | 请求 QPS、延迟 P99、错误率、活跃连接数 |
| Agent 心跳 | Backend 内部统计 | 在线 Agent 数、心跳超时数、任务队列深度 |
| Nginx | nginx_exporter / access_log | 请求量、响应时间、4xx/5xx 错误率 |

### 3.3 告警规则

#### 3.3.1 基础设施告警

| 告警名称 | 条件 | 级别 | 说明 |
|---------|------|------|------|
| NodeDown | up == 0 持续 2 分钟 | Critical | 节点不可达 |
| HighCPU | CPU 使用率 > 90% 持续 5 分钟 | Warning | CPU 过载 |
| HighMemory | 内存使用率 > 85% 持续 5 分钟 | Warning | 内存不足 |
| DiskSpaceLow | 磁盘使用率 > 85% | Warning | 磁盘空间不足 |
| DiskSpaceCritical | 磁盘使用率 > 95% | Critical | 磁盘即将满 |

#### 3.3.2 GPU 告警

| 告警名称 | 条件 | 级别 | 说明 |
|---------|------|------|------|
| GPUTemperatureHigh | GPU 温度 > 85°C 持续 5 分钟 | Warning | GPU 过热 |
| GPUTemperatureCritical | GPU 温度 > 95°C | Critical | GPU 严重过热 |
| GPUMemoryFull | 显存使用率 > 95% | Warning | 显存即将耗尽 |
| GPUDown | GPU 指标消失持续 3 分钟 | Critical | GPU 不可用 |

#### 3.3.3 应用服务告警

| 告警名称 | 条件 | 级别 | 说明 |
|---------|------|------|------|
| BackendDown | Backend 健康检查失败持续 2 分钟 | Critical | API 服务不可用 |
| HighErrorRate | 5xx 错误率 > 5% 持续 3 分钟 | Warning | 服务异常 |
| HighLatency | P99 延迟 > 3s 持续 5 分钟 | Warning | 响应缓慢 |
| AgentOffline | Agent 心跳超时 > 5 分钟 | Warning | GPU 节点失联 |
| TaskStuck | 任务状态为 running 且租约过期 | Warning | 任务卡死 |
| PostgreSQLDown | PostgreSQL 健康检查失败 | Critical | 数据库不可用 |
| RedisDown | Redis 健康检查失败 | Critical | 缓存不可用 |
| PostgreSQLConnHigh | 连接数 > 80% 最大值 | Warning | 数据库连接池即将耗尽 |

### 3.4 告警通知渠道

| 渠道 | 适用级别 | 说明 |
|------|---------|------|
| 企业微信/钉钉 Webhook | Warning + Critical | 即时消息通知运维群 |
| 邮件 | Critical | 发送给运维负责人 |
| Uptime Kuma 状态页 | 所有级别 | 对外展示服务可用性状态 |

Alertmanager 告警路由配置示例：

```yaml
route:
  receiver: default
  group_by: ['alertname', 'instance']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  routes:
    - match:
        severity: critical
      receiver: critical-notify
      repeat_interval: 1h

receivers:
  - name: default
    webhook_configs:
      - url: 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=<key>'
  - name: critical-notify
    webhook_configs:
      - url: 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=<key>'
    email_configs:
      - to: 'ops@example.com'
```

### 3.5 Grafana 监控面板

| 面板名称 | 数据源 | 核心指标 |
|---------|--------|---------|
| 集群概览 | Prometheus | 节点数、GPU 总数、在线率、任务数 |
| GPU 监控 | Prometheus | GPU 利用率、显存、温度、功耗趋势 |
| 应用服务 | Prometheus | API QPS、延迟分布、错误率 |
| 数据库 | Prometheus | 连接数、查询延迟、缓存命中率、慢查询 |
| 容器资源 | Prometheus (cAdvisor) | 容器 CPU/内存/网络使用 |
| Agent 状态 | Prometheus | 在线 Agent 数、心跳延迟、任务队列 |

---

## 4. 日志收集

### 4.1 日志架构

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Backend 日志  │  │ Nginx 日志   │  │ Agent 日志   │
│ (stdout/file)│  │ (access/err) │  │ (file)       │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┼─────────────────┘
                         │
                  ┌──────▼───────┐
                  │   Loki       │
                  │  (日志聚合)   │
                  └──────┬───────┘
                         │
                  ┌──────▼───────┐
                  │   Grafana    │
                  │  (日志查询)   │
                  └──────────────┘
```

### 4.2 日志分类与格式

| 日志来源 | 格式 | 存储位置 | 保留策略 |
|---------|------|---------|---------|
| Backend API | JSON 结构化日志 | stdout → Loki | 30 天 |
| Nginx access | JSON 格式 | /var/log/nginx/access.log → Loki | 30 天 |
| Nginx error | 标准格式 | /var/log/nginx/error.log → Loki | 30 天 |
| PostgreSQL | 标准格式 | 容器 stdout → Loki | 14 天 |
| Redis | 标准格式 | 容器 stdout → Loki | 14 天 |
| Agent | JSON 结构化日志 | /var/log/remotegpu/agent.log → Loki | 30 天 |
| 审计日志 | 数据库记录 | PostgreSQL audit_logs 表 | 永久保留 |

### 4.3 日志采集方案

推荐使用 **Grafana Loki + Promtail** 轻量级日志方案，与现有 Grafana 监控面板无缝集成。

#### Promtail 部署

每个节点部署一个 Promtail 实例，负责采集本地日志并推送到 Loki：

```yaml
# docker-compose.loki.yml
services:
  loki:
    image: grafana/loki:2.9.0
    container_name: remotegpu-loki
    restart: always
    ports:
      - "3100:3100"
    volumes:
      - loki_data:/loki
      - ./loki-config.yml:/etc/loki/local-config.yaml
    command: -config.file=/etc/loki/local-config.yaml
    networks:
      - remotegpu-network

  promtail:
    image: grafana/promtail:2.9.0
    container_name: remotegpu-promtail
    restart: always
    volumes:
      - /var/log:/var/log:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - ./promtail-config.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
    networks:
      - remotegpu-network
```

### 4.4 日志格式规范

Backend API 统一使用 JSON 结构化日志：

```json
{
  "time": "2026-02-06T10:30:00.000Z",
  "level": "info",
  "msg": "请求处理完成",
  "method": "POST",
  "path": "/api/v1/tasks",
  "status": 200,
  "latency_ms": 45,
  "user_id": 123,
  "request_id": "req-abc-123"
}
```

关键字段说明：
- `request_id`：请求唯一标识，用于全链路追踪
- `user_id`：操作用户 ID，用于审计关联
- `latency_ms`：请求处理耗时，用于性能分析

---

## 5. 备份恢复

### 5.1 备份策略总览

| 数据类型 | 备份方式 | 频率 | 保留周期 | 存储位置 |
|---------|---------|------|---------|---------|
| PostgreSQL 数据库 | pg_dump 全量备份 | 每日 02:00 | 30 天 | 本地 + 异地对象存储 |
| PostgreSQL WAL 日志 | 持续归档 | 实时 | 7 天 | 本地 |
| Redis 数据 | RDB 快照 | 每小时 | 7 天 | 本地 |
| Nginx 配置 | 文件备份 | 每次变更 | 永久（Git 管理） | Git 仓库 |
| Docker Compose 配置 | 文件备份 | 每次变更 | 永久（Git 管理） | Git 仓库 |
| RustFS 对象数据 | 增量同步 | 每日 03:00 | 30 天 | 异地存储 |
| 审计日志 | 数据库内 | 随数据库备份 | 永久 | 随数据库 |

### 5.2 PostgreSQL 备份方案

#### 自动备份脚本

```bash
#!/bin/bash
# /opt/remotegpu/scripts/backup-db.sh
set -euo pipefail

BACKUP_DIR="/opt/remotegpu/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# 执行全量备份
docker exec remotegpu-postgresql pg_dump \
  -U remotegpu_user \
  -d remotegpu \
  -Fc \
  -f /tmp/backup_${DATE}.dump

# 拷贝到宿主机
docker cp remotegpu-postgresql:/tmp/backup_${DATE}.dump \
  ${BACKUP_DIR}/backup_${DATE}.dump

# 清理容器内临时文件
docker exec remotegpu-postgresql rm /tmp/backup_${DATE}.dump

# 清理过期备份
find ${BACKUP_DIR} -name "backup_*.dump" -mtime +${RETENTION_DAYS} -delete

echo "[$(date)] 备份完成: backup_${DATE}.dump"
```

#### Crontab 配置

```cron
# 每日 02:00 执行数据库备份
0 2 * * * /opt/remotegpu/scripts/backup-db.sh >> /var/log/remotegpu/backup.log 2>&1
```

### 5.3 数据库恢复流程

```bash
#!/bin/bash
# /opt/remotegpu/scripts/restore-db.sh
set -euo pipefail

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
  echo "用法: $0 <backup_file.dump>"
  exit 1
fi

# 1. 停止 Backend 服务，防止写入
docker stop remotegpu-backend

# 2. 恢复数据库
docker cp ${BACKUP_FILE} remotegpu-postgresql:/tmp/restore.dump
docker exec remotegpu-postgresql pg_restore \
  -U remotegpu_user \
  -d remotegpu \
  --clean --if-exists \
  /tmp/restore.dump

# 3. 清理临时文件
docker exec remotegpu-postgresql rm /tmp/restore.dump

# 4. 重启 Backend 服务
docker start remotegpu-backend

echo "[$(date)] 恢复完成: ${BACKUP_FILE}"
```

### 5.4 备份验证

定期验证备份有效性，确保灾难恢复可行：

| 验证项 | 频率 | 方法 |
|--------|------|------|
| 备份文件完整性 | 每日 | 检查 pg_dump 退出码和文件大小 |
| 恢复测试 | 每月 | 在测试环境执行完整恢复流程 |
| 数据一致性 | 每月 | 恢复后对比关键表的记录数和校验和 |

### 5.5 灾难恢复 RTO/RPO 目标

| 场景 | RPO（数据丢失） | RTO（恢复时间） |
|------|----------------|----------------|
| 数据库故障 | < 24 小时（全量备份间隔） | < 1 小时 |
| 管理节点故障 | 0（无状态服务） | < 30 分钟 |
| GPU 节点故障 | 0（Agent 断线恢复） | < 10 分钟 |
| 全站灾难 | < 24 小时 | < 4 小时 |

---

## 6. 扩展性方案

### 6.1 水平扩展能力分析

| 组件 | 扩展方式 | 说明 |
|------|---------|------|
| Backend API | 水平扩展（多实例） | 无状态服务，通过 Nginx 负载均衡分发 |
| Frontend | 水平扩展（多实例） | 静态资源，CDN 加速 |
| PostgreSQL | 垂直扩展 → 读写分离 | 初期垂直扩展，后期主从复制 |
| Redis | 垂直扩展 → Sentinel | 初期单节点，后期 Sentinel 高可用 |
| Agent | 随 GPU 节点线性扩展 | 每台 GPU 机器一个 Agent |
| Nginx | 垂直扩展 | 单节点足够，极端场景可用 keepalived 双活 |

### 6.2 分阶段扩展路线

#### 阶段一：单机部署（GPU 节点 < 50 台）

所有管理服务部署在单台管理节点上：

```
管理节点（1 台）
├── Nginx
├── Backend API（单实例）
├── Frontend
├── PostgreSQL（单实例）
├── Redis（单实例）
├── Prometheus + Grafana
└── Loki + Promtail
```

适用场景：初期部署、小规模客户

#### 阶段二：分离部署（GPU 节点 50-200 台）

数据库和监控服务独立部署：

```
管理节点（1 台）
├── Nginx
├── Backend API（2 实例，Nginx 负载均衡）
├── Frontend
├── Guacamole
└── Harbor

数据库节点（1 台）
├── PostgreSQL（主）
└── Redis（Sentinel 模式）

监控节点（1 台）
├── Prometheus
├── Grafana
├── Loki
└── Alertmanager
```

适用场景：中等规模客户，需要更高可用性

#### 阶段三：高可用部署（GPU 节点 > 200 台）

全面高可用架构：

```
管理节点（2 台，主备）
├── Nginx（keepalived VIP 漂移）
├── Backend API（多实例）
├── Frontend
├── Guacamole（多实例）
└── Harbor

数据库节点（2 台，主从）
├── PostgreSQL（主从流复制 + 自动故障转移）
└── Redis（Sentinel 3 节点）

监控节点（1 台）
├── Prometheus（联邦模式）
├── Grafana
├── Loki
└── Alertmanager
```

适用场景：大规模客户，要求高可用和故障自动恢复

### 6.3 Backend 多实例负载均衡

Nginx 负载均衡配置示例：

```nginx
upstream backend_api {
    least_conn;
    server backend-1:8080 max_fails=3 fail_timeout=30s;
    server backend-2:8080 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

server {
    listen 443 ssl;
    server_name api.remotegpu.example.com;

    ssl_certificate     /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    location /api/ {
        proxy_pass http://backend_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 10s;
        proxy_read_timeout 60s;
    }
}
```

---

## 7. 运维流程

### 7.1 发布流程

#### 7.1.1 发布步骤

```
1. 代码合并到 main 分支
2. CI 自动构建 Docker 镜像并推送到 Harbor
3. 运维人员在管理节点拉取新镜像
4. 执行滚动更新（先更新一个实例，验证后更新其余）
5. 验证服务健康状态
6. 更新发布记录
```

#### 7.1.2 发布脚本

```bash
#!/bin/bash
# /opt/remotegpu/scripts/deploy.sh
set -euo pipefail

VERSION=${1:-latest}
REGISTRY="harbor.remotegpu.example.com/remotegpu"

echo "[$(date)] 开始部署版本: ${VERSION}"

# 1. 拉取新镜像
docker pull ${REGISTRY}/remotegpu-backend:${VERSION}
docker pull ${REGISTRY}/remotegpu-frontend:${VERSION}

# 2. 备份当前版本号
docker inspect remotegpu-backend --format='{{.Config.Image}}' \
  > /opt/remotegpu/backups/last-version.txt 2>/dev/null || true

# 3. 滚动更新
export VERSION=${VERSION}
docker compose -f docker-compose.prod.yml up -d --no-deps backend frontend

# 4. 等待健康检查通过
echo "等待服务健康检查..."
for i in $(seq 1 30); do
  if docker inspect remotegpu-backend \
    --format='{{.State.Health.Status}}' 2>/dev/null | grep -q healthy; then
    echo "[$(date)] 部署成功: ${VERSION}"
    exit 0
  fi
  sleep 2
done

echo "[$(date)] 健康检查超时，执行回滚"
bash /opt/remotegpu/scripts/rollback.sh
exit 1
```

#### 7.1.3 回滚脚本

```bash
#!/bin/bash
# /opt/remotegpu/scripts/rollback.sh
set -euo pipefail

LAST_VERSION_FILE="/opt/remotegpu/backups/last-version.txt"

if [ ! -f "$LAST_VERSION_FILE" ]; then
  echo "未找到上一版本记录，无法回滚"
  exit 1
fi

LAST_IMAGE=$(cat ${LAST_VERSION_FILE})
echo "[$(date)] 回滚到: ${LAST_IMAGE}"

# 使用上一版本镜像重新启动
docker compose -f docker-compose.prod.yml up -d --no-deps backend frontend

echo "[$(date)] 回滚完成"
```

### 7.2 Agent 升级流程

Agent 部署在 GPU 宿主机上，需要逐台滚动升级：

```
1. 构建新版本 Agent 二进制
2. 上传到分发服务器（RustFS 或 HTTP 文件服务）
3. 逐台 GPU 节点执行升级：
   a. 检查该节点是否有运行中的任务
   b. 如有任务，等待任务完成或标记为维护模式
   c. 停止 Agent 服务：systemctl stop remotegpu-agent
   d. 替换二进制文件
   e. 启动 Agent 服务：systemctl start remotegpu-agent
   f. 验证心跳恢复正常
4. 全部节点升级完成后，确认监控面板无异常
```

### 7.3 数据库迁移流程

```
1. 在测试环境验证迁移脚本
2. 备份生产数据库（执行 backup-db.sh）
3. 停止 Backend 服务，防止写入冲突
4. 执行迁移脚本：
   docker exec -i remotegpu-postgresql \
     psql -U remotegpu_user -d remotegpu < backend/sql/NN_描述.sql
5. 验证迁移结果（检查表结构、数据完整性）
6. 启动 Backend 服务
7. 验证服务正常运行
```

迁移脚本命名规范：`NN_描述.sql`，编号递增，与 `backend/sql/` 目录保持一致。

### 7.4 GPU 节点上架流程

```
1. 硬件就绪：确认 GPU 驱动、CUDA、Docker 已安装
2. 网络配置：确认节点可访问管理节点 API
3. 部署 Agent：
   a. 创建配置目录：mkdir -p /etc/remotegpu
   b. 写入环境变量配置：/etc/remotegpu/agent.env
   c. 拷贝 Agent 二进制到 /usr/local/bin/
   d. 安装 systemd 服务文件
   e. 启动服务：systemctl enable --now remotegpu-agent
4. 在管理平台添加机器（填写 SSH 连接信息）
5. 触发硬件信息采集
6. 验证 Agent 心跳正常、硬件信息完整
7. 机器状态变为 online，可分配给客户
```

### 7.5 故障处理流程

#### 7.5.1 故障分级

| 级别 | 定义 | 响应要求 |
|------|------|---------|
| P0 | 平台完全不可用（API 宕机、数据库故障） | 立即响应，所有运维人员参与 |
| P1 | 核心功能受损（部分 GPU 节点离线、任务无法执行） | 30 分钟内响应 |
| P2 | 非核心功能异常（监控面板异常、日志采集中断） | 2 小时内响应 |
| P3 | 性能下降或告警预警 | 下一工作日处理 |

#### 7.5.2 常见故障处理手册

**Backend API 不可用：**

```
1. 检查容器状态：docker ps -a | grep backend
2. 查看容器日志：docker logs --tail 100 remotegpu-backend
3. 检查数据库连接：docker exec remotegpu-postgresql pg_isready
4. 检查 Redis 连接：docker exec remotegpu-redis redis-cli ping
5. 尝试重启服务：docker compose -f docker-compose.prod.yml restart backend
6. 如重启无效，检查磁盘空间和内存使用
7. 如仍无法恢复，执行版本回滚
```

**GPU 节点 Agent 离线：**

```
1. 检查节点网络连通性：ping <node-ip>
2. SSH 登录节点检查 Agent 状态：systemctl status remotegpu-agent
3. 查看 Agent 日志：journalctl -u remotegpu-agent --tail 100
4. 检查 GPU 驱动状态：nvidia-smi
5. 尝试重启 Agent：systemctl restart remotegpu-agent
6. 如 GPU 驱动异常，重启节点
7. 记录故障原因，更新故障知识库
```

**PostgreSQL 数据库故障：**

```
1. 检查容器状态：docker ps -a | grep postgresql
2. 查看数据库日志：docker logs --tail 200 remotegpu-postgresql
3. 检查磁盘空间：df -h（数据库数据卷）
4. 尝试重启：docker compose -f docker-compose.prod.yml restart postgresql
5. 如数据损坏，从最近备份恢复（执行 restore-db.sh）
6. 恢复后验证数据完整性
```

### 7.6 日常巡检

#### 每日巡检清单

| 巡检项 | 方法 | 预期结果 |
|--------|------|---------|
| 所有容器运行状态 | `docker ps` | 所有容器 Up 且 healthy |
| GPU 节点在线率 | Grafana 面板 | 在线率 > 99% |
| 数据库备份状态 | 检查备份日志 | 最近备份成功且文件大小正常 |
| 磁盘使用率 | `df -h` | 所有分区 < 80% |
| 告警列表 | Alertmanager / 企业微信 | 无未处理的 Critical 告警 |
| API 响应时间 | Grafana 面板 | P99 < 1s |

#### 每周巡检清单

| 巡检项 | 方法 | 预期结果 |
|--------|------|---------|
| Docker 镜像清理 | `docker image prune` | 清理无用镜像，释放磁盘 |
| 数据库连接池 | Grafana 面板 | 连接数 < 最大值 60% |
| Redis 内存使用 | `redis-cli info memory` | 内存使用 < 最大值 70% |
| SSL 证书有效期 | 检查证书到期时间 | 距到期 > 30 天 |
| 安全更新 | `apt list --upgradable` | 无高危安全补丁待安装 |

### 7.7 安全运维

#### 7.7.1 访问控制

| 措施 | 说明 |
|------|------|
| SSH 密钥登录 | 管理节点禁用密码登录，仅允许密钥认证 |
| 端口最小化 | 仅开放必要端口，数据库/Redis 仅内网可访问 |
| 防火墙规则 | 使用 iptables/nftables 限制入站流量 |
| 操作审计 | 所有管理操作通过平台审计日志记录 |

#### 7.7.2 密钥与证书管理

| 项目 | 管理方式 | 轮换周期 |
|------|---------|---------|
| SSL/TLS 证书 | Let's Encrypt 自动续期 或 手动管理 | 90 天（自动）/ 1 年（手动） |
| JWT 签名密钥 | 环境变量，不提交 Git | 每季度轮换 |
| 数据库密码 | 环境变量，不提交 Git | 每季度轮换 |
| Redis 密码 | 环境变量，不提交 Git | 每季度轮换 |
| Agent Token | 后端生成，按节点分配 | 按需轮换 |

#### 7.7.3 Docker 安全加固

| 措施 | 说明 |
|------|------|
| 非 root 运行 | 容器内进程以非 root 用户运行 |
| 只读文件系统 | 生产容器启用 `read_only: true`（需要写入的目录挂载 tmpfs） |
| 资源限制 | 所有容器设置 CPU/内存限制，防止资源耗尽 |
| 镜像扫描 | Harbor 启用镜像漏洞扫描，阻止高危镜像部署 |
| 网络隔离 | 使用 Docker network 隔离不同服务组 |

### 7.8 运维工具清单

| 工具 | 用途 | 部署位置 |
|------|------|---------|
| Docker / Docker Compose | 容器编排 | 管理节点 |
| Prometheus + Alertmanager | 监控告警 | 监控节点 |
| Grafana | 可视化面板 | 监控节点 |
| Loki + Promtail | 日志收集 | 所有节点 |
| Uptime Kuma | 可用性监控 | 监控节点 |
| Harbor | 镜像仓库 | 管理节点 |
| pg_dump / pg_restore | 数据库备份恢复 | 数据库节点 |

### 7.9 生产环境目录结构

```
/opt/remotegpu/
├── docker-compose.prod.yml       # 主编排文件
├── .env.backend                  # Backend 环境变量
├── .env.db                       # 数据库环境变量
├── nginx/
│   ├── nginx.conf                # Nginx 主配置
│   ├── conf.d/                   # 站点配置
│   └── ssl/                      # SSL 证书
├── postgresql/
│   └── postgresql.conf           # PostgreSQL 配置
├── redis/
│   └── redis.conf                # Redis 配置
├── prometheus/
│   └── prometheus.yml            # Prometheus 配置
├── scripts/
│   ├── deploy.sh                 # 发布脚本
│   ├── rollback.sh               # 回滚脚本
│   ├── backup-db.sh              # 数据库备份脚本
│   └── restore-db.sh             # 数据库恢复脚本
├── backups/
│   ├── postgresql/               # 数据库备份文件
│   └── last-version.txt          # 上一版本记录
└── logs/                         # 本地日志目录
```

---

## 8. 总结

### 8.1 关键决策

| 决策项 | 选择 | 理由 |
|--------|------|------|
| 容器编排 | Docker Compose | 项目规模适中，无需 K8s 的复杂性 |
| 监控方案 | Prometheus + Grafana | 社区成熟，GPU 监控生态完善 |
| 日志方案 | Loki + Promtail | 轻量级，与 Grafana 无缝集成 |
| 数据库备份 | pg_dump 全量 + WAL 归档 | 简单可靠，满足 RPO 要求 |
| Agent 部署 | systemd 二进制 | 需要直接访问 GPU 硬件，不适合容器化 |
| 镜像仓库 | Harbor | 支持漏洞扫描，企业级功能完善 |

### 8.2 实施优先级

| 阶段 | 内容 | 前置条件 |
|------|------|---------|
| P0 | 基础部署（Docker Compose 编排、Nginx、Backend、Frontend、PostgreSQL、Redis） | 无 |
| P0 | 数据库备份脚本 | 基础部署完成 |
| P1 | 监控告警（Prometheus + Grafana + Alertmanager） | 基础部署完成 |
| P1 | 日志收集（Loki + Promtail） | 基础部署完成 |
| P1 | Agent 部署自动化脚本 | 基础部署完成 |
| P2 | Harbor 镜像仓库 | 基础部署完成 |
| P2 | Guacamole 远程访问网关 | 基础部署完成 |
| P2 | Uptime Kuma 可用性监控 | 监控告警完成 |
| P3 | 高可用架构（主从数据库、多实例 Backend） | 业务规模增长后 |
