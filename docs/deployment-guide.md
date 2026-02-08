# RemoteGPU 部署指南

## 目录

- [环境要求](#环境要求)
- [项目结构](#项目结构)
- [快速启动](#快速启动)
- [各服务配置说明](#各服务配置说明)
- [后端服务部署](#后端服务部署)
- [前端服务部署](#前端服务部署)
- [Agent 部署到 GPU 机器](#agent-部署到-gpu-机器)
- [监控体系](#监控体系)
- [测试环境](#测试环境)
- [常见问题排查](#常见问题排查)

---

## 环境要求

### 管理节点（运行平台服务）

| 组件 | 最低版本 | 说明 |
|------|---------|------|
| Docker | 24.0+ | 容器运行时 |
| Docker Compose | v2.20+ | 服务编排（推荐使用 `docker compose` V2 命令） |
| 操作系统 | Ubuntu 22.04 / CentOS 8+ | 推荐 Linux x86_64 |
| 内存 | 8 GB+ | 运行全部基础服务的最低要求 |
| 磁盘 | 50 GB+ | 数据库、对象存储、日志等 |

### GPU 节点（运行 Agent）

| 组件 | 最低版本 | 说明 |
|------|---------|------|
| 操作系统 | Ubuntu 20.04+ | 需支持 NVIDIA 驱动 |
| NVIDIA 驱动 | 535+ | 根据 GPU 型号选择 |
| Go | 1.23+ | 仅源码编译时需要 |

### 开发环境（可选）

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.25+ | 后端开发 |
| Node.js | 20+ | 前端开发 |
| npm | 10+ | 前端包管理 |

---

## 项目结构

```
remotegpu/
├── backend/                # Go 后端 API 服务
├── frontend/               # Vue 3 前端
├── agent/                  # Go Agent（部署在 GPU 机器上）
├── docker-compose/         # 基础设施服务编排
│   ├── start-all.sh        # 一键启动所有服务
│   ├── stop-all.sh         # 一键停止所有服务
│   ├── check-status.sh     # 服务状态检查
│   ├── postgresql/         # PostgreSQL 数据库
│   ├── redis/              # Redis 缓存
│   ├── etcd/               # Etcd 配置中心
│   ├── rustfs/             # RustFS 对象存储
│   ├── nginx/              # Nginx 反向代理
│   ├── prometheus/         # Prometheus 监控
│   ├── grafana/            # Grafana 仪表板
│   ├── exporters/          # 各类 Exporter
│   ├── uptime-kuma/        # 可用性监控
│   ├── guacamole/          # 远程桌面网关
│   ├── harbor/             # 容器镜像仓库
│   └── test-env/           # 测试环境（模拟 GPU 节点）
└── docs/                   # 文档
```

---

## 快速启动

### 1. 克隆项目

```bash
git clone <repository-url>
cd remotegpu
```

### 2. 创建 Docker 网络

所有服务共享同一个 Docker 网络，需要先手动创建：

```bash
docker network create remotegpu-network
```

### 3. 一键启动基础设施

```bash
cd docker-compose
chmod +x start-all.sh stop-all.sh check-status.sh
./start-all.sh
```

启动脚本按依赖顺序分五层启动服务：

| 层级 | 服务 | 说明 |
|------|------|------|
| 第一层 | PostgreSQL、Redis、Etcd | 基础存储服务 |
| 第二层 | RustFS、Nginx | 对象存储和反向代理 |
| 第三层 | Prometheus、Grafana、Uptime Kuma | 监控服务 |
| 第四层 | test-env | 测试环境（交互式确认） |
| 第五层 | Guacamole、Harbor | 可选服务（交互式确认） |

### 4. 检查服务状态

```bash
./check-status.sh
```

### 5. 停止所有服务

```bash
./stop-all.sh
```

---

## 各服务配置说明

### 端口映射总览

| 服务 | 容器名 | 端口 | 说明 |
|------|--------|------|------|
| PostgreSQL | remotegpu-postgresql | 5432 | 数据库 |
| Redis | remotegpu-redis | 6379 | 缓存 |
| Etcd | remotegpu-etcd | 2379, 2380 | 配置中心 |
| RustFS | remotegpu-rustfs | 9000 (API), 9001 (Console) | 对象存储 |
| Nginx | remotegpu-nginx | 80, 443 | 反向代理 |
| Prometheus | remotegpu-prometheus | 19090 | 监控采集 |
| Grafana | remotegpu-grafana | 13000 | 监控仪表板 |
| Uptime Kuma | remotegpu-uptime-kuma | 13001 | 可用性监控 |
| Guacamole | remotegpu-guacamole | 8081 | 远程桌面网关 |
| Harbor | remotegpu-harbor-nginx | 8082 | 镜像仓库 |
| Node Exporter | remotegpu-node-exporter | 9100 | 主机指标 |
| Postgres Exporter | remotegpu-postgres-exporter | 9187 | 数据库指标 |
| Redis Exporter | remotegpu-redis-exporter | 9121 | 缓存指标 |
| Nginx Exporter | remotegpu-nginx-exporter | 9113 | 代理指标 |

### PostgreSQL

- 镜像：`postgres:17`
- 默认用户：`remotegpu_user`，密码：`remotegpu_password`，数据库：`remotegpu`
- 数据持久化：Docker volume `postgresql_data`
- 自定义配置：`docker-compose/postgresql/postgresql.conf`
- 初始化脚本：`docker-compose/postgresql/init.sql`

生产环境务必修改默认密码。

### Redis

- 镜像：`redis:8.4.0-alpine`
- 默认密码：在 `docker-compose/redis/redis.conf` 中配置
- 数据持久化：Docker volume `redis_data`

### Etcd

- 镜像：`quay.io/coreos/etcd:v3.5.13`
- 单节点模式，集群 token：`remotegpu-cluster`
- 数据持久化：Docker volume `etcd_data`

### RustFS（对象存储）

- 镜像：`rustfs/rustfs:latest`
- API 端口 9000，控制台端口 9001
- 默认访问密钥：`rustfsadmin` / `rustfsadmin`（生产环境务必修改）
- 数据持久化：Docker volume `rustfs_data`

### Nginx

- 镜像：`nginx:1.24-alpine`
- 配置文件：`docker-compose/nginx/nginx.conf` 和 `docker-compose/nginx/conf.d/`
- SSL 证书目录：`docker-compose/nginx/ssl/`
- 前端静态文件挂载自 `frontend/dist/`（需先构建前端）

### Guacamole（可选）

- 包含 `guacd`（协议代理）和 `guacamole`（Web 应用）两个容器
- 需要独立的 PostgreSQL 数据库（默认库名 `guacamole`）
- 访问地址：`http://<host>:8081/guacamole/`

### Harbor（可选）

- 包含 core、registry、registryctl、jobservice、portal、nginx 六个容器
- 访问地址：`http://<host>:8082`
- 需要独立的 PostgreSQL 数据库（默认库名 `harbor`）

---

## 后端服务部署

### 方式一：Docker 部署（推荐）

```bash
cd backend
# 复制并修改配置文件
cp config/config.yaml.example config/config.yaml
# 编辑 config.yaml，填入实际的数据库、Redis 等连接信息

# 构建镜像
docker build -t remotegpu-backend .

# 运行
docker run -d \
  --name remotegpu-backend \
  --network remotegpu-network \
  -p 8080:8080 \
  -v $(pwd)/config/config.yaml:/app/config/config.yaml \
  remotegpu-backend
```

### 方式二：源码编译

```bash
cd backend
cp config/config.yaml.example config/config.yaml
# 编辑 config.yaml

go build -o remotegpu ./cmd/
./remotegpu server
```

### 关键配置项

编辑 `backend/config/config.yaml`，以下为必须修改的配置：

```yaml
database:
  host: remotegpu-postgresql   # Docker 网络内使用容器名
  port: 5432
  user: remotegpu_user
  password: "你的数据库密码"
  dbname: remotegpu

redis:
  host: remotegpu-redis
  port: 6379
  password: "你的 Redis 密码"

jwt:
  secret: "至少32字符的随机字符串"

encryption:
  key: "恰好32字节的 AES-256 密钥"
```

其他服务（Etcd、Prometheus、Harbor 等）按需在配置中启用并填写连接信息。

---

## 前端服务部署

### 方式一：Docker 部署（推荐）

```bash
cd frontend

# 构建镜像
docker build -t remotegpu-frontend .

# 运行
docker run -d \
  --name remotegpu-frontend \
  --network remotegpu-network \
  -p 3000:80 \
  remotegpu-frontend
```

### 方式二：集成到 Nginx

```bash
cd frontend
npm ci
npm run build
```

构建产物在 `frontend/dist/` 目录，Nginx 服务已通过 volume 挂载该目录：

```yaml
# docker-compose/nginx/docker-compose.yml 中的挂载
- ../../frontend/dist:/usr/share/nginx/html:ro
```

构建完成后重启 Nginx 即可生效：

```bash
cd docker-compose/nginx
docker compose restart
```

---

## Agent 部署到 GPU 机器

Agent 是部署在每台 GPU 机器上的客户端程序，负责与后端通信、上报机器状态和执行远程操作。

### 方式一：Docker 部署

```bash
cd agent
docker build -t remotegpu-agent .

docker run -d \
  --name remotegpu-agent \
  --restart unless-stopped \
  remotegpu-agent
```

### 方式二：源码编译部署

```bash
# 在 GPU 机器上编译（需要 Go 1.23+）
cd agent
go build -o remotegpu-agent ./cmd/

# 安装为 systemd 服务
sudo cp remotegpu-agent /usr/local/bin/
```

创建 systemd 服务文件 `/etc/systemd/system/remotegpu-agent.service`：

```ini
[Unit]
Description=RemoteGPU Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/remotegpu-agent
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable remotegpu-agent
sudo systemctl start remotegpu-agent
```

### 方式三：交叉编译

在开发机上为 GPU 机器交叉编译，然后通过 scp 分发：

```bash
cd agent
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o remotegpu-agent ./cmd/

# 分发到 GPU 机器
scp remotegpu-agent user@gpu-node:/usr/local/bin/
```

### Agent 默认端口

Agent 默认监听 HTTP 端口 `8090`，后端通过该端口与 Agent 通信。确保 GPU 机器的防火墙允许后端访问此端口。

---

## 监控体系

监控体系由 Prometheus + Grafana + Exporters 组成。

### 架构

```
GPU 节点 / 服务 --> Exporters --> Prometheus --> Grafana
```

### Exporters

Exporters 配置在 `docker-compose/exporters/docker-compose.yml`，包含：

| Exporter | 端口 | 监控目标 |
|----------|------|---------|
| Node Exporter | 9100 | 主机 CPU、内存、磁盘、网络 |
| Postgres Exporter | 9187 | PostgreSQL 连接数、查询性能 |
| Redis Exporter | 9121 | Redis 内存、命中率、连接数 |
| Nginx Exporter | 9113 | Nginx 请求量、连接状态 |

启动 Exporters：

```bash
cd docker-compose/exporters
docker compose up -d
```

### Prometheus

- 访问地址：`http://<host>:19090`
- 配置文件：`docker-compose/prometheus/prometheus.yml`
- 告警规则：`docker-compose/prometheus/rules/`
- 数据保留 30 天

### Grafana

- 访问地址：`http://<host>:13000`
- 默认账号：`admin` / `changeme_grafana_password`
- 数据源和仪表板通过 `docker-compose/grafana/provisioning/` 自动配置

### Uptime Kuma

- 访问地址：`http://<host>:13001`
- 用于监控各服务的可用性

---

## 测试环境

测试环境位于 `docker-compose/test-env/`，模拟多台 GPU 机器，用于开发和测试。

### 测试节点

| 容器名 | SSH 端口 | Jupyter 端口 | Agent 端口 | Jupyter |
|--------|---------|-------------|-----------|---------|
| test-gpu-01 | 2201 | 8891 | 8901 | 启用 |
| test-gpu-02 | 2202 | 8892 | 8902 | 启用 |
| test-gpu-03 | 2203 | 8893 | 8903 | 禁用 |

### 启动测试环境

```bash
cd docker-compose/test-env
docker compose up -d --build
```

首次构建会编译 Agent 二进制（多阶段构建），后续构建会利用 Docker 缓存。

### 默认凭据

- SSH 用户：`root`，密码：`remotegpu123`（可通过 `SSH_PASSWORD` 环境变量修改）
- Jupyter Token：`remotegpu`（可通过 `JUPYTER_TOKEN` 环境变量修改）
- 每个节点内置模拟的 `nvidia-smi` 命令，用于测试 GPU 监控功能

### 连接测试

```bash
# SSH 连接测试节点 01
ssh -p 2201 root@localhost

# 访问 Jupyter
# http://localhost:8891/?token=remotegpu
```

---

## 常见问题排查

### 1. Docker 网络不通

**现象**：容器之间无法通过容器名互相访问。

**排查**：

```bash
# 检查网络是否存在
docker network inspect remotegpu-network

# 检查容器是否加入了网络
docker inspect <container-name> | grep -A 10 Networks
```

**解决**：确保所有服务都加入了 `remotegpu-network` 外部网络，且网络已提前创建。

### 2. 端口冲突

**现象**：容器启动失败，提示端口已被占用。

**排查**：

```bash
# 查看端口占用
sudo netstat -tlnp | grep <port>
# 或
sudo ss -tlnp | grep <port>
```

**解决**：停止占用端口的进程，或修改 `docker-compose.yml` 中的端口映射。

### 3. PostgreSQL 连接失败

**现象**：后端启动报错 `connection refused` 或 `authentication failed`。

**排查**：

```bash
# 检查 PostgreSQL 容器状态
docker ps | grep postgresql

# 查看日志
cd docker-compose/postgresql
docker compose logs

# 测试连接
docker exec -it remotegpu-postgresql psql -U remotegpu_user -d remotegpu
```

**解决**：
- 确认 PostgreSQL 容器已启动且健康检查通过
- 确认 `config.yaml` 中的数据库连接信息与 `docker-compose.yml` 中的环境变量一致
- Docker 网络内使用容器名 `remotegpu-postgresql` 而非 `127.0.0.1`

### 4. Agent 无法连接后端

**现象**：Agent 启动后无法注册到后端平台。

**排查**：

```bash
# 检查 Agent 进程
ps aux | grep remotegpu-agent

# 检查 Agent 端口
ss -tlnp | grep 8090

# 从后端机器测试连通性
curl http://<gpu-node-ip>:8090/health
```

**解决**：
- 确认 GPU 机器防火墙已放行 8090 端口
- 确认后端 `config.yaml` 中 `agent.token` 与 Agent 配置一致
- 检查网络路由，确保后端能访问 GPU 机器

### 5. Nginx 启动失败

**现象**：Nginx 容器启动后立即退出。

**排查**：

```bash
cd docker-compose/nginx
docker compose logs
```

**解决**：
- 检查 `nginx.conf` 语法是否正确
- 确认 SSL 证书文件存在（如果配置了 HTTPS）
- 确认 `frontend/dist/` 目录存在（Nginx 挂载了该目录）

### 6. 数据卷清理

如需完全重置某个服务的数据：

```bash
cd docker-compose/<service>
docker compose down -v   # -v 会删除关联的数据卷
docker compose up -d
```

### 7. 查看容器日志

```bash
# 查看指定服务的实时日志
cd docker-compose/<service>
docker compose logs -f

# 查看最近 100 行日志
docker compose logs --tail 100
```
