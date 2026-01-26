# RemoteGPU 基础设施配置文档

> 本文档描述 RemoteGPU 系统所需的基础设施组件、配置要求和部署规范
>
> 创建日期：2026-01-26

---

## 1. 基础设施架构概览

```yaml
基础设施分层:
  数据层:
    - PostgreSQL: 主数据库
    - Redis: 缓存和会话存储
    - Etcd: 配置中心和服务发现
    - MinIO/S3: 对象存储

  计算层:
    - Kubernetes: 容器编排
    - Docker: 容器运行时
    - GPU 节点: 计算资源

  网络层:
    - Nginx: 反向代理和负载均衡
    - Traefik: API 网关
    - JumpServer: 堡垒机
    - Guacamole: 远程桌面网关

  监控层:
    - Prometheus: 指标采集
    - Grafana: 可视化
    - Loki: 日志聚合
    - AlertManager: 告警管理
```

---

## 2. 数据库服务

### 2.1 PostgreSQL

**版本要求：** PostgreSQL 14.x 或更高

**硬件配置：**
```yaml
最小配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 100GB SSD

推荐配置:
  CPU: 8 核
  内存: 16GB
  磁盘: 500GB SSD (RAID 10)
```

**配置参数：**
```ini
# postgresql.conf

# 连接配置
max_connections = 200
superuser_reserved_connections = 3

# 内存配置
shared_buffers = 4GB
effective_cache_size = 12GB
maintenance_work_mem = 1GB
work_mem = 64MB

# WAL 配置
wal_level = replica
max_wal_size = 2GB
min_wal_size = 1GB
checkpoint_completion_target = 0.9

# 查询优化
random_page_cost = 1.1
effective_io_concurrency = 200

# 日志配置
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_min_duration_statement = 1000
```

**数据库初始化：**
```sql
-- 创建数据库
CREATE DATABASE remotegpu
    WITH
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

-- 创建用户
CREATE USER remotegpu_user WITH PASSWORD 'your_secure_password';

-- 授权
GRANT ALL PRIVILEGES ON DATABASE remotegpu TO remotegpu_user;

-- 连接到数据库
\c remotegpu

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";
```

---

## 3. Redis 缓存服务

**版本要求：** Redis 7.x 或更高

**硬件配置：**
```yaml
最小配置:
  CPU: 2 核
  内存: 4GB
  磁盘: 20GB SSD

推荐配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 50GB SSD
```

**配置参数：**
```ini
# redis.conf

# 网络配置
bind 0.0.0.0
port 6379
protected-mode yes
requirepass your_redis_password

# 内存配置
maxmemory 4gb
maxmemory-policy allkeys-lru

# 持久化配置
save 900 1
save 300 10
save 60 10000
appendonly yes
appendfilename "appendonly.aof"

# 日志配置
loglevel notice
logfile "/var/log/redis/redis-server.log"
```

**使用场景：**
```yaml
缓存用途:
  - 用户会话存储 (Session)
  - JWT Token 黑名单
  - 端口映射缓存
  - 监控数据临时存储
  - 分布式锁
  - 消息队列
```

---

## 4. Etcd 配置中心

**版本要求：** Etcd 3.5.x 或更高

**硬件配置：**
```yaml
最小配置 (单节点):
  CPU: 2 核
  内存: 2GB
  磁盘: 20GB SSD

推荐配置 (3 节点集群):
  CPU: 4 核
  内存: 4GB
  磁盘: 50GB SSD
```

**配置参数：**
```yaml
# etcd.conf.yml
name: 'etcd-node-1'
data-dir: /var/lib/etcd
listen-client-urls: http://0.0.0.0:2379
advertise-client-urls: http://192.168.1.10:2379
listen-peer-urls: http://0.0.0.0:2380
initial-advertise-peer-urls: http://192.168.1.10:2380
initial-cluster: 'etcd-node-1=http://192.168.1.10:2380'
initial-cluster-token: 'remotegpu-cluster'
initial-cluster-state: 'new'
```

**使用场景：**
```yaml
配置管理:
  - 系统配置动态更新
  - 服务发现
  - 分布式锁
  - 主机状态同步
```

---

## 5. Kubernetes 集群

**版本要求：** Kubernetes 1.28.x 或更高

**集群架构：**
```yaml
Master 节点 (3 个):
  CPU: 4 核
  内存: 8GB
  磁盘: 100GB SSD

Worker 节点 (根据需求):
  CPU 节点:
    CPU: 16 核
    内存: 32GB
    磁盘: 200GB SSD

  GPU 节点:
    CPU: 32 核
    内存: 128GB
    GPU: 4-8 张
    磁盘: 500GB SSD
```

**必需组件：**
```yaml
核心组件:
  - kubeadm: 集群初始化工具
  - kubelet: 节点代理
  - kubectl: 命令行工具
  - containerd: 容器运行时

网络插件:
  - Calico 或 Flannel

存储插件:
  - Rook-Ceph 或 Longhorn

GPU 支持:
  - NVIDIA Device Plugin
  - NVIDIA GPU Operator
```

---

## 6. MinIO 对象存储

**版本要求：** MinIO RELEASE.2024-01-01 或更高

**硬件配置：**
```yaml
最小配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 500GB SSD

推荐配置 (分布式):
  节点数: 4 个
  CPU: 8 核/节点
  内存: 16GB/节点
  磁盘: 2TB SSD/节点
```

**Docker Compose 部署：**
```yaml
version: '3.8'
services:
  minio:
    image: minio/minio:latest
    container_name: minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: admin
      MINIO_ROOT_PASSWORD: your_secure_password
    volumes:
      - /data/minio:/data
    command: server /data --console-address ":9001"
```

**使用场景：**
```yaml
存储用途:
  - 数据集存储
  - 模型文件存储
  - 训练日志存储
  - 镜像备份
```

---

## 6.2 JuiceFS 分布式文件系统

**版本要求：** JuiceFS 1.1.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 根据数据量配置
```

**架构说明：**
```yaml
JuiceFS 架构:
  元数据引擎:
    - Redis (推荐，高性能)
    - PostgreSQL
    - MySQL
    - TiKV

  对象存储:
    - MinIO
    - AWS S3
    - 阿里云 OSS
    - 腾讯云 COS
```

**安装部署：**
```bash
# 下载 JuiceFS
wget https://github.com/juicedata/juicefs/releases/download/v1.1.0/juicefs-1.1.0-linux-amd64.tar.gz
tar -zxf juicefs-1.1.0-linux-amd64.tar.gz
sudo install juicefs /usr/local/bin

# 格式化文件系统
juicefs format \
  --storage minio \
  --bucket http://minio:9000/juicefs \
  --access-key admin \
  --secret-key your_password \
  redis://redis:6379/1 \
  remotegpu-fs

# 挂载文件系统
juicefs mount redis://redis:6379/1 /mnt/juicefs
```

**使用场景：**
```yaml
应用场景:
  - 共享数据集存储
  - 多节点数据共享
  - 训练数据缓存
  - 模型文件共享
  - 支持 POSIX 接口
```

---

## 7. JumpServer 堡垒机

**版本要求：** JumpServer v3.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 8 核
  内存: 16GB
  磁盘: 100GB SSD
```

**Docker Compose 部署：**
```yaml
version: '3.8'
services:
  jumpserver:
    image: jumpserver/jms_all:v3.10.0
    container_name: jumpserver
    ports:
      - "80:80"
      - "2222:2222"
    environment:
      SECRET_KEY: your_secret_key
      BOOTSTRAP_TOKEN: your_bootstrap_token
      DB_ENGINE: postgresql
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: jumpserver
      DB_PASSWORD: your_db_password
      DB_NAME: jumpserver
      REDIS_HOST: redis
      REDIS_PORT: 6379
      REDIS_PASSWORD: your_redis_password
    volumes:
      - /opt/jumpserver/data:/opt/jumpserver/data
```

**使用场景：**
```yaml
功能:
  - SSH 访问审计
  - 主机资产管理
  - 用户权限管理
  - 操作录像回放
```

---

## 8. Apache Guacamole 远程桌面网关

**版本要求：** Guacamole 1.5.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 50GB SSD
```

**Docker Compose 部署：**
```yaml
version: '3.8'
services:
  guacd:
    image: guacamole/guacd:1.5.4
    container_name: guacd
    restart: always

  guacamole:
    image: guacamole/guacamole:1.5.4
    container_name: guacamole
    ports:
      - "8080:8080"
    environment:
      GUACD_HOSTNAME: guacd
      POSTGRES_HOSTNAME: postgres
      POSTGRES_DATABASE: guacamole
      POSTGRES_USER: guacamole
      POSTGRES_PASSWORD: your_password
    depends_on:
      - guacd
```

**使用场景：**
```yaml
功能:
  - RDP 远程桌面访问
  - VNC 访问
  - SSH 访问
  - 浏览器内访问，无需客户端
```

---

## 9. Nginx 反向代理

**版本要求：** Nginx 1.24.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 4 核
  内存: 4GB
  磁盘: 20GB SSD
```

**配置示例：**
```nginx
# /etc/nginx/nginx.conf
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 4096;
    use epoll;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript;

    # 上传大小限制
    client_max_body_size 10G;

    include /etc/nginx/conf.d/*.conf;
}
```

---

## 10. Prometheus 监控

**版本要求：** Prometheus 2.45.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 200GB SSD
```

**配置文件：**
```yaml
# prometheus.yml
global:
  scrape_interval: 30s
  evaluation_interval: 30s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node1:9100', 'node2:9100']

  - job_name: 'nvidia-gpu'
    static_configs:
      - targets: ['gpu-node1:9445', 'gpu-node2:9445']
```

---

## 11. Grafana 可视化

**版本要求：** Grafana 10.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 2 核
  内存: 4GB
  磁盘: 50GB SSD
```

**Docker 部署：**
```yaml
version: '3.8'
services:
  grafana:
    image: grafana/grafana:10.2.0
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
      GF_INSTALL_PLUGINS: grafana-piechart-panel
    volumes:
      - grafana-data:/var/lib/grafana

volumes:
  grafana-data:
```

---

## 12. Uptime Kuma 服务监控

**版本要求：** Uptime Kuma 1.23.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 2 核
  内存: 2GB
  磁盘: 20GB SSD
```

**Docker 部署：**
```yaml
version: '3.8'
services:
  uptime-kuma:
    image: louislam/uptime-kuma:1
    container_name: uptime-kuma
    ports:
      - "3001:3001"
    volumes:
      - uptime-kuma-data:/app/data
    restart: always

volumes:
  uptime-kuma-data:
```

**监控对象：**
```yaml
监控服务:
  HTTP/HTTPS 监控:
    - API 服务健康检查
    - Web 前端可用性
    - JupyterLab 服务状态

  TCP 端口监控:
    - PostgreSQL (5432)
    - Redis (6379)
    - Kubernetes API (6443)

  Ping 监控:
    - GPU 节点存活检查
    - 网络连通性检查

  关键词监控:
    - API 响应内容检查
    - 错误页面检测
```

**告警配置：**
```yaml
通知渠道:
  - 邮件通知
  - 钉钉机器人
  - 企业微信
  - Slack
  - Webhook
```

---

## 13. Harbor 镜像仓库

**版本要求：** Harbor 2.9.x 或更高

**硬件配置：**
```yaml
推荐配置:
  CPU: 4 核
  内存: 8GB
  磁盘: 500GB SSD
```

**Docker Compose 部署：**
```bash
# 下载 Harbor 离线安装包
wget https://github.com/goharbor/harbor/releases/download/v2.9.0/harbor-offline-installer-v2.9.0.tgz
tar xzvf harbor-offline-installer-v2.9.0.tgz
cd harbor

# 配置 harbor.yml
cp harbor.yml.tmpl harbor.yml
# 编辑 harbor.yml，设置 hostname、https、密码等

# 安装
./install.sh
```

**使用场景：**
```yaml
功能:
  - Docker 镜像存储
  - 镜像扫描
  - 镜像签名
  - 镜像复制
```

---

## 13. 网络规划

**网络架构：**
```yaml
网络分层:
  管理网络 (10.0.1.0/24):
    - 用途: 管理节点、堡垒机
    - VLAN: 10

  业务网络 (10.0.2.0/24):
    - 用途: API 服务、Web 服务
    - VLAN: 20

  存储网络 (10.0.3.0/24):
    - 用途: 数据库、对象存储
    - VLAN: 30

  计算网络 (10.0.4.0/24):
    - 用途: GPU 节点、容器网络
    - VLAN: 40
```

**端口规划：**
```yaml
对外服务端口:
  - 80: HTTP (Nginx)
  - 443: HTTPS (Nginx)
  - 2222: SSH (JumpServer)

内部服务端口:
  - 5432: PostgreSQL
  - 6379: Redis
  - 2379-2380: Etcd
  - 6443: Kubernetes API
  - 9000: MinIO
  - 9090: Prometheus
  - 3000: Grafana
  - 3001: Uptime Kuma
```

---

## 14. 部署清单

**基础设施组件清单：**

| 组件 | 版本 | 节点数 | 用途 | 优先级 |
|------|------|--------|------|--------|
| PostgreSQL | 14.x | 1-3 | 主数据库 | P0 |
| Redis | 7.x | 1-3 | 缓存 | P0 |
| Etcd | 3.5.x | 3 | 配置中心 | P0 |
| Kubernetes | 1.28.x | 3+ | 容器编排 | P0 |
| MinIO | Latest | 4+ | 对象存储 | P0 |
| Nginx | 1.24.x | 2+ | 反向代理 | P0 |
| Prometheus | 2.45.x | 1 | 监控 | P1 |
| Grafana | 10.x | 1 | 可视化 | P1 |
| Uptime Kuma | 1.23.x | 1 | 服务监控 | P1 |
| Harbor | 2.9.x | 1 | 镜像仓库 | P1 |
| JumpServer | 3.x | 1 | 堡垒机 | P2 |
| Guacamole | 1.5.x | 1 | 远程桌面 | P2 |

---

## 15. 安全配置

**防火墙规则：**
```bash
# 仅允许必要端口
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw enable
```

**SSL/TLS 证书：**
```yaml
证书管理:
  - 使用 Let's Encrypt 自动续期
  - 或使用企业 CA 签发证书
  
配置位置:
  - Nginx: /etc/nginx/ssl/
  - 有效期: 90 天（自动续期）
```

**密码策略：**
```yaml
要求:
  - 最小长度: 12 位
  - 包含大小写字母、数字、特殊字符
  - 定期更换（90 天）
  - 禁止使用弱密码
```

---

## 16. 部署步骤

**第一阶段：基础环境准备**
```bash
# 1. 操作系统准备（所有节点）
# Ubuntu 22.04 LTS 或 CentOS 8

# 2. 时间同步
apt install -y chrony
systemctl enable chrony

# 3. 关闭 swap
swapoff -a
sed -i '/swap/d' /etc/fstab

# 4. 内核参数优化
cat >> /etc/sysctl.conf << 'EOL'
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOL
sysctl -p
```

**第二阶段：数据库部署**
```bash
# 1. 部署 PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=your_password \
  -v /data/postgres:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:14

# 2. 部署 Redis
docker run -d \
  --name redis \
  -v /data/redis:/data \
  -p 6379:6379 \
  redis:7 redis-server --requirepass your_password

# 3. 部署 Etcd
docker run -d \
  --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  quay.io/coreos/etcd:v3.5.0
```

---

**第三阶段：Kubernetes 集群部署**
```bash
# 1. 安装 Docker/Containerd
apt install -y docker.io
systemctl enable docker

# 2. 安装 kubeadm、kubelet、kubectl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
apt update && apt install -y kubeadm kubelet kubectl

# 3. 初始化 Master 节点
kubeadm init --pod-network-cidr=10.244.0.0/16

# 4. 安装网络插件（Calico）
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml

# 5. Worker 节点加入集群
kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>
```

---

## 17. 监控指标

**关键监控指标：**
```yaml
系统指标:
  - CPU 使用率 < 80%
  - 内存使用率 < 85%
  - 磁盘使用率 < 90%
  - 网络带宽使用率 < 70%

服务指标:
  - API 响应时间 < 500ms
  - 数据库连接数 < 150
  - Redis 内存使用 < 80%
  - Kubernetes Pod 健康状态

GPU 指标:
  - GPU 使用率
  - GPU 温度 < 85°C
  - 显存使用率
```

---

## 18. 备份策略

**数据备份计划：**
```yaml
PostgreSQL 备份:
  - 全量备份: 每天 2:00 AM
  - 增量备份: 每 6 小时
  - 保留周期: 30 天
  - 备份工具: pg_dump + cron

对象存储备份:
  - 数据集: 每周全量备份
  - 模型文件: 每天增量备份
  - 保留周期: 90 天

配置备份:
  - Etcd 快照: 每天
  - Kubernetes 配置: 每天
  - 保留周期: 7 天
```

---

## 19. 故障恢复

**恢复时间目标（RTO）：**
```yaml
服务级别:
  - P0 服务: 30 分钟
  - P1 服务: 2 小时
  - P2 服务: 4 小时

数据恢复:
  - 数据库: 1 小时
  - 对象存储: 4 小时
```

---

**文档版本：** v1.0
**创建日期：** 2026-01-26
**维护者：** RemoteGPU 运维团队
