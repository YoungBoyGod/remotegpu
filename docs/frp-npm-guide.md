# Nginx Proxy Manager + frp 完整实施指南

## 概述

本文档提供基于 Nginx Proxy Manager (NPM) 的 frp 方案完整实施指南，适用于需要为大量 GPU 机器配置外网访问的场景。

### 方案架构

```
外网用户
    ↓ HTTPS
    ↓ https://gpu1-jupyter.gpu.domain.com
    ↓
云服务器 (公网IP)
    ├─ Nginx Proxy Manager (端口 80/443)
    │   ├─ SSL 证书管理
    │   └─ 反向代理配置
    ├─ frps (端口 7000)
    │   └─ 接收 frpc 连接
    └─ frp 暴露端口 (7001-7400)
        ↓
        ↓ 通过 frp 隧道
        ↓
GPU 机器 (内网)
    ├─ frpc (连接到 frps)
    ├─ Jupyter (端口 8888)
    └─ TensorBoard (端口 6006)
```

---

## 方案对比

### NPM 方案 vs 手动 Nginx 方案

| 特性 | NPM 方案 | 手动 Nginx 方案 |
|------|----------|----------------|
| **配置方式** | 图形界面 | 手动编辑配置文件 |
| **SSL 管理** | 自动申请和续期 | 手动执行 certbot |
| **学习曲线** | 低，适合新手 | 高，需要熟悉 nginx |
| **批量配置** | 需要脚本辅助 | 需要脚本辅助 |
| **灵活性** | 中等 | 高，完全自定义 |
| **资源占用** | 较高（Docker + 数据库） | 较低（仅 nginx） |
| **适用场景** | 开发/测试环境 | 生产环境 |

### 推荐选择

- **开发环境**：推荐 NPM 方案（简单易用）
- **生产环境**：推荐手动 Nginx 方案（性能更好）
- **混合环境**：可以先用 NPM 快速搭建，后期迁移到手动 Nginx

---

## 实施步骤概览

1. **准备工作** - 域名、DNS、服务器
2. **安装 frps** - 在云服务器安装 frp 服务端
3. **安装 NPM** - 在云服务器安装 Nginx Proxy Manager
4. **配置 frpc** - 在 GPU 机器安装 frp 客户端
5. **配置 NPM** - 添加反向代理和 SSL 证书
6. **批量配置** - 使用脚本批量添加代理
7. **测试验证** - 验证访问和 SSL 证书

---

## 第一步：准备工作

### 1.1 域名准备

注册域名（如 `domain.com`），并配置 DNS 解析。

#### 配置泛域名解析

在域名 DNS 管理中添加记录：

| 类型 | 主机记录 | 记录值 | TTL |
|------|----------|--------|-----|
| A | `gpu` | `云服务器公网IP` | 600 |
| A | `*.gpu` | `云服务器公网IP` | 600 |

**验证 DNS 解析**：
```bash
nslookup gpu.domain.com
nslookup gpu1-jupyter.gpu.domain.com
```

### 1.2 服务器准备

#### 云服务器要求

- **配置**：2核4GB 以上
- **带宽**：5Mbps 以上（根据实际需求）
- **系统**：Ubuntu 20.04+ / CentOS 7+
- **公网 IP**：固定 IP

#### 安全组配置

开放以下端口：

| 端口 | 协议 | 说明 |
|------|------|------|
| 80 | TCP | HTTP（NPM 和证书验证） |
| 443 | TCP | HTTPS（NPM） |
| 81 | TCP | NPM 管理界面 |
| 7000 | TCP | frps 控制端口 |
| 7001-7400 | TCP | frp 数据端口（200台机器×2服务） |

### 1.3 GPU 机器准备

确保 GPU 机器可以访问云服务器：
```bash
# 在 GPU 机器上测试连接
ping 云服务器IP
telnet 云服务器IP 7000
```

---

## 第二步：安装 frps

### 2.1 下载 frp

```bash
# 在云服务器上执行
cd ~
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64
```

### 2.2 配置 frps

创建配置文件 `frps.toml`：

```toml
# frps.toml
bindPort = 7000
vhostHTTPPort = 80
vhostHTTPSPort = 443

# 认证
auth.method = "token"
auth.token = "your-secure-token-here"

# 日志
log.to = "/var/log/frps.log"
log.level = "info"
log.maxDays = 7

# 管理界面（可选）
webServer.addr = "0.0.0.0"
webServer.port = 7500
webServer.user = "admin"
webServer.password = "admin123"
```

### 2.3 启动 frps

#### 使用 systemd 管理

创建服务文件 `/etc/systemd/system/frps.service`：

```ini
[Unit]
Description=frp server
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/root/frp_0.52.3_linux_amd64/frps -c /root/frp_0.52.3_linux_amd64/frps.toml

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable frps
sudo systemctl start frps
sudo systemctl status frps
```

---

## 第三步：安装 NPM

详细步骤请参考：**[frp-npm-installation.md](frp-npm-installation.md)**

### 3.1 快速安装

```bash
# 创建目录
mkdir -p ~/npm
cd ~/npm

# 创建 docker-compose.yml
cat > docker-compose.yml <<EOF
version: '3.8'
services:
  npm:
    image: jc21/nginx-proxy-manager:latest
    container_name: npm
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
      - "81:81"
    volumes:
      - npm_data:/data
      - npm_letsencrypt:/etc/letsencrypt
volumes:
  npm_data:
  npm_letsencrypt:
EOF

# 启动 NPM
docker-compose up -d
```

### 3.2 初始登录

访问 `http://云服务器IP:81`，使用默认凭据登录：
- Email: `admin@example.com`
- Password: `changeme`

首次登录后修改管理员信息。

---

## 第四步：配置 frpc

### 4.1 在 GPU 机器上安装 frpc

```bash
# 下载 frp
cd ~
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64
```

### 4.2 配置 frpc

创建配置文件 `frpc.toml`（以 GPU1 为例）：

```toml
# frpc.toml - GPU1
serverAddr = "云服务器IP"
serverPort = 7000

auth.method = "token"
auth.token = "your-secure-token-here"

log.to = "/var/log/frpc.log"
log.level = "info"

# Jupyter
[[proxies]]
name = "gpu1-jupyter"
type = "tcp"
localIP = "127.0.0.1"
localPort = 8888
remotePort = 7001

# TensorBoard
[[proxies]]
name = "gpu1-tensorboard"
type = "tcp"
localIP = "127.0.0.1"
localPort = 6006
remotePort = 7002
```

### 4.3 启动 frpc

创建服务文件 `/etc/systemd/system/frpc.service`：

```ini
[Unit]
Description=frp client
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/root/frp_0.52.3_linux_amd64/frpc -c /root/frp_0.52.3_linux_amd64/frpc.toml

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable frpc
sudo systemctl start frpc
sudo systemctl status frpc
```

---

## 第五步：配置 NPM 代理

详细步骤请参考：**[frp-npm-proxy-config.md](frp-npm-proxy-config.md)**

### 5.1 申请通配符证书（推荐）

1. 登录 NPM 管理界面
2. 点击 **SSL Certificates** → **Add SSL Certificate**
3. 选择 **Let's Encrypt**
4. Domain Names 填写：`*.gpu.domain.com`
5. 勾选 **Use a DNS Challenge**
6. 选择 DNS 提供商并填写 API 凭据
7. 点击 **Save**

### 5.2 添加代理主机（示例）

1. 点击 **Proxy Hosts** → **Add Proxy Host**
2. 填写配置：
   - Domain Names: `gpu1-jupyter.gpu.domain.com`
   - Forward Hostname/IP: `127.0.0.1`
   - Forward Port: `7001`
   - ☑ Block Common Exploits
   - ☑ Websockets Support
3. 切换到 **SSL** 标签页：
   - SSL Certificate: 选择通配符证书
   - ☑ Force SSL
   - ☑ HTTP/2 Support
   - ☑ HSTS Enabled
4. 点击 **Save**

---

## 第六步：批量配置

详细步骤请参考：**[frp-npm-batch-config.md](frp-npm-batch-config.md)**

### 6.1 准备配置数据

创建机器列表 `machines.json`：

```json
[
  {"id": 1, "jupyter_port": 7001, "tensorboard_port": 7002},
  {"id": 2, "jupyter_port": 7003, "tensorboard_port": 7004},
  ...
]
```

### 6.2 使用 API 批量添加

下载批量配置脚本：
```bash
wget https://raw.githubusercontent.com/.../npm_batch_add.py
```

修改配置并运行：
```bash
python3 npm_batch_add.py
```

---

## 第七步：测试验证

### 7.1 验证 frp 连接

```bash
# 在云服务器上检查 frp 端口
netstat -tlnp | grep 7001

# 测试端口访问
curl http://127.0.0.1:7001
```

### 7.2 验证 HTTPS 访问

在浏览器中访问：
```
https://gpu1-jupyter.gpu.domain.com
```

检查：
- ✓ 页面正常加载
- ✓ 浏览器显示绿色锁图标
- ✓ 证书有效且匹配域名

### 7.3 批量测试脚本

```bash
#!/bin/bash
for i in {1..200}; do
  url="https://gpu${i}-jupyter.gpu.domain.com"
  status=$(curl -s -o /dev/null -w "%{http_code}" "$url")
  echo "GPU${i}: $status"
done
```

---

## 维护和监控

### 日常维护

1. **检查 frps 状态**：`systemctl status frps`
2. **检查 NPM 状态**：`docker ps | grep npm`
3. **查看日志**：`docker logs npm`
4. **备份配置**：定期备份 NPM 数据卷

### SSL 证书续期

NPM 会自动续期证书，无需手动操作。

查看证书状态：
1. 登录 NPM 管理界面
2. 点击 **SSL Certificates**
3. 查看到期时间

### 性能监控

```bash
# 检查 NPM 资源占用
docker stats npm

# 检查连接数
netstat -an | grep :443 | wc -l
```

---

## 故障排查

### 无法访问 HTTPS

**检查清单**：
1. DNS 解析是否正确：`nslookup gpu1-jupyter.gpu.domain.com`
2. frpc 是否连接：在云服务器检查 `netstat -tlnp | grep 7001`
3. NPM 代理配置是否正确：检查 Forward Port
4. SSL 证书是否有效：检查 SSL Certificates 页面

### 502 Bad Gateway

**可能原因**：
1. frp 端口配置错误
2. GPU 机器上的服务未启动
3. frpc 未连接到 frps

**解决方案**：
```bash
# 检查 frpc 状态
systemctl status frpc

# 检查本地服务
curl http://localhost:8888

# 查看 frpc 日志
journalctl -u frpc -f
```

---

## 文档索引

### 安装和配置

- **[NPM 安装指南](frp-npm-installation.md)** - NPM 安装和初始配置
- **[NPM 代理配置](frp-npm-proxy-config.md)** - 配置反向代理和 SSL
- **[NPM 批量配置](frp-npm-batch-config.md)** - 批量添加代理配置

### frp 相关

- **[frp 实施指南](frp-implementation-guide.md)** - frp 完整实施步骤
- **[frp 服务端配置](frp-step3-frps.md)** - frps 详细配置
- **[frp 客户端配置](frp-step5-frpc.md)** - frpc 详细配置

### 其他方案

- **[手动 Nginx 配置](frp-step4-nginx.md)** - 手动配置 nginx 反向代理
- **[批量配置脚本](frp-batch-scripts.md)** - 各种批量配置脚本

---

## 优势和注意事项

### NPM 方案的优势

✓ **简单易用**：图形界面，无需学习 nginx 配置语法
✓ **自动化**：SSL 证书自动申请和续期
✓ **集中管理**：统一管理所有代理配置
✓ **快速部署**：适合快速搭建开发环境

### 注意事项

⚠ **资源占用**：NPM 需要 Docker 和数据库，资源占用较高
⚠ **性能**：相比手动 nginx，性能略低
⚠ **批量配置**：仍需要脚本辅助，不能完全图形化
⚠ **生产环境**：建议生产环境使用手动 nginx 方案

---

## 下一步

完成 NPM + frp 方案后，可以考虑：

1. **监控告警**：配置 Prometheus + Grafana 监控
2. **访问控制**：配置 NPM Access List 限制访问
3. **负载均衡**：为高负载服务配置负载均衡
4. **迁移到生产**：迁移到手动 Nginx 方案

---

## 参考资源

- [Nginx Proxy Manager 官方文档](https://nginxproxymanager.com/)
- [frp 官方文档](https://gofrp.org/docs/)
- [Let's Encrypt 文档](https://letsencrypt.org/docs/)
