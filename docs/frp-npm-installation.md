# Nginx Proxy Manager 安装指南

## 概述

Nginx Proxy Manager (NPM) 是一个基于 Nginx 的反向代理管理工具，提供图形化界面，简化了 SSL 证书管理和反向代理配置。

### NPM 的优势

- **图形化界面**：无需手动编辑 nginx 配置文件
- **自动 SSL 管理**：自动申请和续期 Let's Encrypt 证书
- **简单易用**：适合不熟悉 nginx 配置的用户
- **集中管理**：统一管理所有代理配置

### 适用场景

本文档适用于在云服务器上安装 NPM，用于管理 frp 反向代理配置。

---

## 前置要求

### 服务器要求

- 云服务器（公网 IP）
- 操作系统：Ubuntu 20.04+ / CentOS 7+
- 内存：至少 1GB RAM
- 已安装 Docker 和 Docker Compose

### 域名要求

- 已注册域名（如 `gpu.domain.com`）
- 域名 DNS 已指向云服务器 IP
- 准备配置泛域名解析 `*.gpu.domain.com`

### 端口要求

NPM 需要使用以下端口：
- **80**：HTTP 访问
- **443**：HTTPS 访问
- **81**：NPM 管理界面

---

## 安装步骤

### 1. 安装 Docker

如果服务器未安装 Docker，执行以下命令：

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo systemctl enable docker
sudo systemctl start docker

# 添加当前用户到 docker 组
sudo usermod -aG docker $USER
```

### 2. 创建 NPM 目录

```bash
mkdir -p ~/npm
cd ~/npm
```

### 3. 创建 Docker Compose 配置

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  npm:
    image: jc21/nginx-proxy-manager:latest
    container_name: npm
    restart: unless-stopped
    ports:
      - "80:80"      # HTTP
      - "443:443"    # HTTPS
      - "81:81"      # 管理界面
    volumes:
      - npm_data:/data
      - npm_letsencrypt:/etc/letsencrypt
    environment:
      # 可选：设置数据库连接（默认使用内置 SQLite）
      DB_SQLITE_FILE: "/data/database.sqlite"

volumes:
  npm_data:
    driver: local
  npm_letsencrypt:
    driver: local
```

### 4. 启动 NPM

```bash
docker-compose up -d
```

### 5. 检查运行状态

```bash
# 查看容器状态
docker ps | grep npm

# 查看日志
docker logs npm
```

---

## 初始配置

### 1. 访问管理界面

在浏览器中访问：

```
http://your-server-ip:81
```

### 2. 默认登录凭据

- **Email**: `admin@example.com`
- **Password**: `changeme`

### 3. 修改管理员信息

首次登录后，系统会要求修改管理员信息：

1. 输入新的邮箱地址
2. 设置新密码（建议使用强密码）
3. 输入姓名（可选）

### 4. 修改管理界面端口（可选）

如果需要修改管理界面端口（默认 81），可以编辑 `docker-compose.yml`：

```yaml
ports:
  - "8081:81"  # 将管理界面改为 8081 端口
```

然后重启容器：

```bash
docker-compose down
docker-compose up -d
```

---

## 防火墙配置

### 云服务器安全组

在云服务商控制台配置安全组规则：

| 端口 | 协议 | 来源 | 说明 |
|------|------|------|------|
| 80 | TCP | 0.0.0.0/0 | HTTP 访问 |
| 443 | TCP | 0.0.0.0/0 | HTTPS 访问 |
| 81 | TCP | 你的IP | NPM 管理界面（建议限制来源 IP） |
| 7000-7200 | TCP | GPU机器IP段 | frp 客户端连接（根据实际情况调整） |

### 系统防火墙

如果服务器启用了 firewalld 或 ufw，需要开放端口：

```bash
# firewalld (CentOS/RHEL)
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=81/tcp
sudo firewall-cmd --permanent --add-port=7000-7200/tcp
sudo firewall-cmd --reload

# ufw (Ubuntu/Debian)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 81/tcp
sudo ufw allow 7000:7200/tcp
```

---

## 验证安装

### 1. 检查 NPM 服务

```bash
# 检查容器运行状态
docker ps | grep npm

# 应该看到类似输出：
# CONTAINER ID   IMAGE                             STATUS         PORTS
# abc123def456   jc21/nginx-proxy-manager:latest   Up 5 minutes   0.0.0.0:80-81->80-81/tcp, 0.0.0.0:443->443/tcp
```

### 2. 访问管理界面

访问 `http://your-server-ip:81`，应该能看到 NPM 登录界面。

### 3. 检查端口监听

```bash
sudo netstat -tlnp | grep -E ':(80|443|81)'

# 应该看到：
# tcp6  0  0 :::80    :::*    LISTEN  12345/docker-proxy
# tcp6  0  0 :::443   :::*    LISTEN  12345/docker-proxy
# tcp6  0  0 :::81    :::*    LISTEN  12345/docker-proxy
```

---

## 常见问题

### 端口被占用

**问题**：启动 NPM 时提示端口 80/443 被占用

**解决方案**：

```bash
# 查看占用端口的进程
sudo netstat -tlnp | grep -E ':(80|443)'

# 如果是 nginx 或 apache，停止服务
sudo systemctl stop nginx
sudo systemctl stop apache2

# 禁止开机自启
sudo systemctl disable nginx
sudo systemctl disable apache2
```

### 无法访问管理界面

**问题**：浏览器无法访问 `http://server-ip:81`

**检查清单**：
1. 检查容器是否运行：`docker ps | grep npm`
2. 检查防火墙规则是否开放 81 端口
3. 检查云服务商安全组是否开放 81 端口
4. 尝试从服务器本地访问：`curl http://localhost:81`

### Docker 容器无法启动

**问题**：`docker-compose up -d` 失败

**解决方案**：

```bash
# 查看详细日志
docker-compose logs npm

# 常见原因：
# 1. 端口冲突 - 修改 docker-compose.yml 中的端口映射
# 2. 权限问题 - 确保当前用户在 docker 组中
# 3. 磁盘空间不足 - 检查磁盘空间：df -h
```

---

## 数据备份

NPM 的配置数据存储在 Docker volume 中，建议定期备份。

### 备份命令

```bash
# 备份数据卷
docker run --rm \
  -v npm_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/npm-backup-$(date +%Y%m%d).tar.gz /data

# 备份 Let's Encrypt 证书
docker run --rm \
  -v npm_letsencrypt:/letsencrypt \
  -v $(pwd):/backup \
  alpine tar czf /backup/npm-letsencrypt-$(date +%Y%m%d).tar.gz /letsencrypt
```

### 恢复命令

```bash
# 恢复数据卷
docker run --rm \
  -v npm_data:/data \
  -v $(pwd):/backup \
  alpine sh -c "cd / && tar xzf /backup/npm-backup-20260207.tar.gz"

# 恢复证书
docker run --rm \
  -v npm_letsencrypt:/letsencrypt \
  -v $(pwd):/backup \
  alpine sh -c "cd / && tar xzf /backup/npm-letsencrypt-20260207.tar.gz"

# 重启容器
docker-compose restart
```

---

## 升级 NPM

### 升级步骤

```bash
cd ~/npm

# 拉取最新镜像
docker-compose pull

# 重启容器
docker-compose down
docker-compose up -d

# 检查版本
docker logs npm | grep "Version"
```

---

## 下一步

安装完成后，继续阅读：

- **[NPM 代理配置指南](frp-npm-proxy-config.md)** - 配置反向代理
- **[NPM 批量配置方案](frp-npm-batch-config.md)** - 批量添加 200 台机器的配置
- **[NPM 完整实施指南](frp-npm-guide.md)** - 完整的 frp + NPM 方案

---

## 参考资源

- [Nginx Proxy Manager 官方文档](https://nginxproxymanager.com/guide/)
- [NPM GitHub 仓库](https://github.com/NginxProxyManager/nginx-proxy-manager)
- [Docker 官方文档](https://docs.docker.com/)
