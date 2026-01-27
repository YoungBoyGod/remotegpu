# RemoteGPU Docker Compose 配置集合

本目录包含 RemoteGPU 系统所需的所有基础设施服务的独立 Docker Compose 配置。

## 📋 服务列表

### 核心服务（必需）

| 服务 | 目录 | 端口 | 说明 |
|------|------|------|------|
| PostgreSQL | `postgresql/` | 5432 | 主数据库 |
| Redis | `redis/` | 6379 | 缓存和会话存储 |
| Etcd | `etcd/` | 2379, 2380 | 配置中心和服务发现 |
| RustFS | `rustfs/` | 9000, 9001 | 对象存储 |
| Nginx | `nginx/` | 80, 443 | 反向代理 |

### 监控服务（推荐）

| 服务 | 目录 | 端口 | 说明 |
|------|------|------|------|
| Prometheus | `prometheus/` | 9090 | 监控指标采集 |
| Grafana | `grafana/` | 3000 | 监控可视化 |
| Uptime Kuma | `uptime-kuma/` | 3001 | 服务监控 |

### 可选服务（参考配置）

| 服务 | 目录 | 端口 | 说明 | 备注 |
|------|------|------|------|------|
| JumpServer | `jumpserver/` | 8080, 2222 | 堡垒机 | ⚠️ 使用外部服务 |
| Guacamole | `guacamole/` | 8081 | 远程桌面网关 | 可选部署 |
| Harbor | `harbor/` | 8082 | 镜像仓库 | 可选部署 |

### 外部服务

以下服务在本项目中使用外部部署：
- **Kubernetes**: 容器编排平台（外部K8s集群）
- **JumpServer**: 堡垒机（使用外部实例）

## 🚀 快速开始

### 1. 启动单个服务

```bash
# 进入服务目录
cd postgresql/

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 2. 启动所有服务

```bash
# 在 docker-compose 目录下执行
for dir in */; do
  echo "Starting $dir..."
  cd "$dir"
  docker-compose up -d
  cd ..
done
```

### 3. 停止所有服务

```bash
# 在 docker-compose 目录下执行
for dir in */; do
  echo "Stopping $dir..."
  cd "$dir"
  docker-compose down
  cd ..
done
```

## 🔧 配置说明

### 环境变量

大多数服务支持通过环境变量配置密码和参数。建议在每个服务目录下创建 `.env` 文件：

```bash
# 示例：postgresql/.env
POSTGRES_PASSWORD=your_secure_password
```

### 默认密码

**⚠️ 安全警告：生产环境必须修改所有默认密码！**

| 服务 | 默认用户名 | 默认密码 |
|------|-----------|---------|
| PostgreSQL | remotegpu_user | changeme_secure_password |
| Redis | - | changeme_redis_password |
| MinIO | admin | changeme_minio_password |
| Grafana | admin | changeme_grafana_password |
| JumpServer | admin | admin |
| Guacamole | guacadmin | guacadmin |
| Harbor | admin | Harbor12345 |

## 📦 服务依赖关系

某些服务依赖其他服务，建议按以下顺序启动：

1. **第一层（基础服务）**
   - PostgreSQL
   - Redis
   - Etcd

2. **第二层（存储和网络）**
   - MinIO
   - Nginx

3. **第三层（监控和管理）**
   - Prometheus
   - Grafana
   - Uptime Kuma

4. **第四层（安全和镜像）**
   - JumpServer（需要 PostgreSQL 和 Redis）
   - Guacamole（需要 PostgreSQL）
   - Harbor（需要 PostgreSQL）

## 🔍 健康检查

每个服务都配置了健康检查，可以使用以下命令查看状态：

```bash
# 查看所有容器状态
docker ps -a

# 查看特定服务健康状态
docker inspect --format='{{.State.Health.Status}}' remotegpu-postgresql
```

## 📊 监控和日志

### 查看日志

```bash
# 查看实时日志
docker-compose logs -f [service_name]

# 查看最近100行日志
docker-compose logs --tail=100 [service_name]
```

### 监控指标

- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- Uptime Kuma: http://localhost:3001

## 🔐 安全建议

1. **修改所有默认密码**
2. **使用 HTTPS**：为 Nginx 配置 SSL 证书
3. **网络隔离**：生产环境使用独立网络
4. **定期备份**：特别是 PostgreSQL 和 MinIO 数据
5. **限制访问**：使用防火墙规则限制端口访问

## 💾 数据持久化

所有服务都使用 Docker volumes 持久化数据：

```bash
# 查看所有 volumes
docker volume ls | grep remotegpu

# 备份 volume
docker run --rm -v remotegpu-postgresql_data:/data -v $(pwd):/backup alpine tar czf /backup/postgresql-backup.tar.gz /data
```

## 🛠️ 故障排查

### 服务无法启动

```bash
# 查看详细日志
docker-compose logs [service_name]

# 检查端口占用
netstat -tulpn | grep [port]

# 重新创建容器
docker-compose down
docker-compose up -d --force-recreate
```

### 数据库连接失败

```bash
# 检查数据库是否就绪
docker exec remotegpu-postgresql pg_isready -U remotegpu_user

# 测试连接
docker exec -it remotegpu-postgresql psql -U remotegpu_user -d remotegpu
```

### 清理和重置

```bash
# 停止并删除容器、网络
docker-compose down

# 同时删除 volumes（⚠️ 会丢失数据）
docker-compose down -v
```

## 📚 更多信息

每个服务目录下都有详细的 README.md 文件，包含：
- 服务配置说明
- 使用示例
- 常见问题解决

## 🤝 贡献

如有问题或建议，请提交 Issue 或 Pull Request。

---

**文档版本：** v1.0
**创建日期：** 2026-01-27
**维护者：** RemoteGPU 团队
