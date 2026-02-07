# 方案4：frp 内网穿透详细实施流程

## 方案概述

**架构**：用户 → 公网frp服务器 → frp客户端 → 内网服务

**适用场景**：内网环境，无公网IP时使用

**注意**：您已有公网IP（45.78.48.169），此方案不是最优选择，仅作为备选方案。

---

## 前置条件

1. 一台有公网IP的服务器（作为frp服务端）
2. 内网服务器（192.168.10.210）
3. 域名（可选，用于HTTPS）

---

## 第一部分：frp服务端配置

### 步骤1：安装frp服务端

**执行位置**：公网服务器

```bash
# 下载frp
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz

# 解压
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64

# 移动到系统目录
sudo mkdir -p /opt/frp
sudo cp frps /opt/frp/
sudo cp frps.toml /opt/frp/
```

### 步骤2：配置frp服务端

**执行位置**：公网服务器

```bash
# 编辑配置文件
sudo nano /opt/frp/frps.toml
```

配置内容：
```toml
bindPort = 7000  # frp服务端口
vhostHTTPPort = 80  # HTTP代理端口
vhostHTTPSPort = 443  # HTTPS代理端口

# 认证
auth.method = "token"
auth.token = "your_secure_token_here"  # 修改为强密码

# 日志
log.to = "/var/log/frp/frps.log"
log.level = "info"
log.maxDays = 7
```

### 步骤3：创建systemd服务

**执行位置**：公网服务器

```bash
# 创建日志目录
sudo mkdir -p /var/log/frp

# 创建服务文件
sudo nano /etc/systemd/system/frps.service
```

服务文件内容：
```ini
[Unit]
Description=frp server
After=network.target

[Service]
Type=simple
ExecStart=/opt/frp/frps -c /opt/frp/frps.toml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl start frps
sudo systemctl enable frps
sudo systemctl status frps
```
