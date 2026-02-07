# 场景2B/3：frp 内网穿透方案

## 适用环境

**场景2B**：有出口 IP + 无防火墙权限
**场景3**：NAT 后 + 无固定公网 IP

**共同特点**：
- ❌ 云服务器无法主动连接本地
- ✅ 本地可以访问公网（出站）

**典型环境**：家庭宽带、企业内网（无管理权限）

---

## 架构图

```
用户（全国）
  ↓
域名 → 云服务器（nginx + frps）
  ↓
frp 隧道（本地主动建立）
  ↓
本地服务器（frpc + backend）
  ↓
200台 GPU
```

---

## 第一部分：云服务器配置

### 步骤1：安装 frps

**执行位置**：云服务器

```bash
# 下载 frp
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz

# 安装
sudo mkdir -p /opt/frp
sudo cp frp_0.52.3_linux_amd64/frps /opt/frp/
```

### 步骤2：配置 frps

**执行位置**：云服务器

```bash
sudo nano /opt/frp/frps.toml
```

配置内容：
```toml
bindPort = 7000
vhostHTTPPort = 8080

auth.method = "token"
auth.token = "your_secure_token_123"  # 修改为强密码

log.to = "/var/log/frp/frps.log"
log.level = "info"
```

### 步骤3：创建系统服务

**执行位置**：云服务器

```bash
sudo mkdir -p /var/log/frp

sudo tee /etc/systemd/system/frps.service > /dev/null <<EOF
[Unit]
Description=frp server
After=network.target

[Service]
Type=simple
ExecStart=/opt/frp/frps -c /opt/frp/frps.toml
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start frps
sudo systemctl enable frps
```

### 步骤4：配置防火墙

**执行位置**：云服务器

```bash
# 开放 frp 端口
sudo ufw allow 7000/tcp
sudo ufw allow 8080/tcp
```
