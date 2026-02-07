# RemoteGPU 基础设施配置审查报告

**审查日期**: 2026-02-07
**审查人**: DevOps Engineer (devops)

---

## 一、当前基础设施配置状态

### 1.1 服务总览

项目采用 Docker Compose 分服务部署架构，所有服务配置位于 `docker-compose/` 目录下，每个服务独立一个子目录。

| 服务 | 镜像版本 | 端口映射 | 健康检查 | 状态 |
|------|---------|---------|---------|------|
| PostgreSQL | postgres:17 | 5432:5432 | ✅ pg_isready | 配置完整 |
| Redis | redis:8.4.0-alpine | 6379:6379 | ✅ redis-cli ping | 配置完整 |
| Nginx | nginx:1.24-alpine | 80:80, 443:443 | ✅ nginx -t | 配置完整 |
| Prometheus | prom/prometheus:v2.48.0 | 19090:9090 | ✅ wget spider | 配置完整 |
| Grafana | grafana/grafana:10.2.0 | 13000:3000 | ✅ wget spider | 配置完整 |
| Etcd | coreos/etcd:v3.5.13 | 2379:2379 | ✅ etcdctl | 配置完整 |
| RustFS | rustfs/rustfs:latest | 9000:9000, 9001:9001 | ✅ curl health | 配置完整 |
| Uptime Kuma | louislam/uptime-kuma:1 | 13001:3001 | ✅ node healthcheck | 配置完整 |
| Guacamole | guacamole:1.5.4 | 8081:8080 | ✅ curl | 配置完整 |
| Harbor | goharbor:v2.9.0 | 8082:8080 | ❌ 无 | 需补充 |
| JumpServer | jms_all:v3.10.0 | 8080:80 | ✅ curl | 配置完整 |
| Exporters | 多个 | 9100/9187/9121/9113 | ❌ 无 | 需补充 |

### 1.2 运维脚本

| 脚本 | 路径 | 功能 |
|------|------|------|
| start-all.sh | docker-compose/start-all.sh | 按依赖层级启动所有服务 |
| stop-all.sh | docker-compose/stop-all.sh | 按反向依赖停止所有服务 |
| check-status.sh | docker-compose/check-status.sh | 检查所有容器和端口状态 |

### 1.3 监控告警规则

Prometheus 已配置三组告警规则：
- `service_alerts.yml` — 服务宕机、配置重载失败
- `node_alerts.yml` — CPU/内存/磁盘告警
- `app_alerts.yml` — API 错误率、PG 连接数、Redis 内存

### 1.4 后端配置 (backend/config/config.yaml)

后端服务连接到 `192.168.10.210` 上的 PostgreSQL、Redis、Etcd、Prometheus 等服务，当前为开发/测试环境配置，`mode: debug`。

### 1.5 Agent 配置 (agent/agent.yaml)

Agent 连接到公网地址 `101.43.50.104:60180`，使用 Token 认证，轮询间隔 5s。

---

## 二、发现的问题与风险

### 2.1 🔴 严重问题

#### 2.1.1 配置文件中存在明文凭据

`backend/config/config.yaml` 中包含多个明文密码和密钥：

- 数据库密码: `remotegpu_password`
- Redis 密码: `remotegpu_password`
- JWT Secret: 明文硬编码
- AES 加密密钥: 明文硬编码
- S3 Access/Secret Key: 明文硬编码
- Harbor 密码: `Harbor12345`
- Agent Token: 明文硬编码
- Guacamole 密码: `guacadmin`（默认密码）

**建议**: 生产环境必须使用环境变量或密钥管理服务（如 Vault）替代明文密码。`config.yaml` 应加入 `.gitignore`，仅保留 `config.yaml.example`。

#### 2.1.2 docker-compose-infrastructure.yml 与独立配置冲突

根目录的 `docker-compose-infrastructure.yml` 定义了 Redis、Prometheus、Grafana，但与 `docker-compose/` 下的独立配置存在**版本和端口冲突**：

| 服务 | infrastructure.yml | 独立配置 |
|------|-------------------|---------|
| Redis | redis:7-alpine, 端口 6379, 无密码 | redis:8.4.0-alpine, 端口 6379, 有密码 |
| Prometheus | prom/prometheus:latest, 端口 9090 | prom/prometheus:v2.48.0, 端口 19090 |
| Grafana | grafana/grafana:latest, 端口 3000, 密码 admin123 | grafana:10.2.0, 端口 13000, 密码通过 .env |

**建议**: 废弃 `docker-compose-infrastructure.yml`，统一使用 `docker-compose/` 下的独立配置。

#### 2.1.3 网络隔离不足

所有服务都使用各自独立定义的 `remotegpu-network` bridge 网络，但由于每个 docker-compose 文件独立运行，实际上会创建**多个不同的网络**（如 `postgresql_remotegpu-network`、`redis_remotegpu-network`），导致容器间无法通信。

**建议**: 使用外部网络（`external: true`），先手动创建共享网络：
```bash
docker network create remotegpu-network
```
然后在所有 docker-compose 文件中引用：
```yaml
networks:
  remotegpu-network:
    external: true
```

### 2.2 🟡 中等问题

#### 2.2.1 Exporters 缺少健康检查

`docker-compose/exporters/docker-compose.yml` 中的 4 个 exporter（node-exporter、postgres-exporter、redis-exporter、nginx-exporter）均未配置健康检查。

#### 2.2.2 Exporters 使用硬编码 IP 地址

Exporters 配置中硬编码了 `192.168.10.210` 作为各服务地址，不利于环境迁移。

**建议**: 使用环境变量或 `.env` 文件管理 IP 地址。

#### 2.2.3 PostgreSQL 内存配置偏高

`postgresql.conf` 配置了 `shared_buffers = 4GB`、`effective_cache_size = 12GB`，这要求宿主机至少 16GB 内存。对于开发环境可能过高。

**建议**: 提供开发环境和生产环境两套配置。

#### 2.2.4 Nginx SSL 目录挂载但无证书

Nginx 配置挂载了 `./ssl:/etc/nginx/ssl:ro`，但未见 SSL 证书文件，且 `default.conf` 仅监听 80 端口，未配置 HTTPS。

**建议**: 如需 HTTPS，需补充 SSL 证书配置；如暂不需要，移除 443 端口映射和 ssl 卷挂载。

#### 2.2.5 Nginx 前端静态文件未挂载

`default.conf` 中 `root /usr/share/nginx/html`，但 docker-compose 中未挂载前端构建产物目录。

**建议**: 添加前端构建产物的卷挂载：
```yaml
volumes:
  - /path/to/frontend/dist:/usr/share/nginx/html:ro
```

#### 2.2.6 Harbor 配置不完整

Harbor 的 docker-compose 缺少必要的配置文件挂载（harbor.yml），且 harbor-nginx 缺少配置文件。实际部署 Harbor 建议使用官方安装器。

#### 2.2.7 RustFS 使用弱密码

RustFS 的 Access Key 和 Secret Key 均为 `rustfsadmin`，与 `config.yaml` 中配置的密钥不一致。

### 2.3 🟢 轻微问题

#### 2.3.1 docker-compose 版本声明不一致

部分文件使用 `version: "3.8"`，部分省略（新版 Docker Compose 不再需要）。建议统一。

#### 2.3.2 Prometheus scrape_interval 不一致

`prometheus.yml` 中全局 `scrape_interval: 30s`，对于 GPU 监控场景可能偏慢。

#### 2.3.3 GPU 监控未启用

`prometheus.yml` 中 NVIDIA GPU 监控（dcgm-exporter）的 targets 被注释掉，尚未配置。

#### 2.3.4 test-machines 使用过时镜像

`docker-compose-test-machines.yml` 使用 `rastasheep/ubuntu-sshd:18.04`（Ubuntu 18.04 已 EOL），且密码为弱密码。

---

## 三、部署建议与步骤

### 3.1 部署前准备

#### 步骤 1：创建共享网络
```bash
docker network create remotegpu-network
```

#### 步骤 2：修改所有 docker-compose 文件的网络配置
将每个文件中的网络定义改为外部网络：
```yaml
networks:
  remotegpu-network:
    external: true
```

#### 步骤 3：配置环境变量
为每个服务创建 `.env` 文件（参考 `.env.example`），设置强密码。

#### 步骤 4：准备 SSL 证书（如需 HTTPS）
将证书放入 `docker-compose/nginx/ssl/` 目录，并更新 Nginx 配置。

### 3.2 推荐部署顺序

使用 `start-all.sh` 脚本，按以下层级启动：

1. **第一层 — 基础服务**: PostgreSQL → Redis → Etcd
2. **第二层 — 存储和网络**: RustFS → Nginx
3. **第三层 — 监控**: Prometheus → Grafana → Uptime Kuma → Exporters
4. **第四层 — 可选服务**: Guacamole、Harbor、JumpServer

### 3.3 部署后验证

```bash
# 检查所有容器状态
bash docker-compose/check-status.sh

# 验证 PostgreSQL 连接
docker exec remotegpu-postgresql pg_isready -U remotegpu_user

# 验证 Redis 连接
docker exec remotegpu-redis redis-cli -a remotegpu_password ping

# 验证 Prometheus targets
curl http://localhost:19090/api/v1/targets

# 验证 Grafana
curl http://localhost:13000/api/health
```

### 3.4 后续改进建议

1. **统一编排**: 考虑将所有服务合并到一个 docker-compose 文件中（使用 profiles 区分必选/可选服务），避免网络隔离问题
2. **密钥管理**: 引入 HashiCorp Vault 或 Docker Secrets 管理敏感信息
3. **日志聚合**: 添加 Loki + Promtail 或 ELK 进行集中日志管理
4. **备份策略**: 为 PostgreSQL 和 Redis 配置定期备份（pg_dump + redis-cli bgsave）
5. **GPU 监控**: 在 GPU 节点部署 dcgm-exporter，并取消 Prometheus 中的注释
6. **Alertmanager**: Prometheus 已配置告警规则但未部署 Alertmanager，无法实际发送告警通知
7. **资源限制**: 为所有容器添加 `deploy.resources.limits` 防止资源耗尽

---

## 四、配置文件清单

| 文件路径 | 用途 |
|---------|------|
| docker-compose/postgresql/docker-compose.yml | PostgreSQL 数据库 |
| docker-compose/postgresql/postgresql.conf | PG 性能调优配置 |
| docker-compose/postgresql/init.sql | 数据库初始化脚本 |
| docker-compose/redis/docker-compose.yml | Redis 缓存 |
| docker-compose/redis/redis.conf | Redis 配置（含密码、持久化） |
| docker-compose/nginx/docker-compose.yml | Nginx 反向代理 |
| docker-compose/nginx/nginx.conf | Nginx 主配置 |
| docker-compose/nginx/conf.d/default.conf | Nginx 站点配置 |
| docker-compose/nginx/logrotate.conf | 日志轮转配置 |
| docker-compose/prometheus/docker-compose.yml | Prometheus 监控 |
| docker-compose/prometheus/prometheus.yml | 抓取目标配置 |
| docker-compose/prometheus/rules/*.yml | 告警规则（3个文件） |
| docker-compose/grafana/docker-compose.yml | Grafana 可视化 |
| docker-compose/exporters/docker-compose.yml | 4个 Exporter |
| docker-compose/etcd/docker-compose.yml | Etcd 键值存储 |
| docker-compose/rustfs/docker-compose.yml | RustFS 对象存储 |
| docker-compose/uptime-kuma/docker-compose.yml | Uptime Kuma 可用性监控 |
| docker-compose/guacamole/docker-compose.yml | Guacamole 远程桌面 |
| docker-compose/harbor/docker-compose.yml | Harbor 镜像仓库 |
| docker-compose/jumpserver/docker-compose.yml | JumpServer 堡垒机 |
| docker-compose/jupyter-ssh/docker-compose.yml | Jupyter+SSH 容器 |
| docker-compose-infrastructure.yml | 基础设施合并配置（建议废弃） |
| docker-compose-test-machines.yml | 测试用 SSH 机器 |
| backend/config/config.yaml | 后端应用配置 |
| backend/config/config.yaml.example | 后端配置模板 |
| agent/agent.yaml | Agent 客户端配置 |
