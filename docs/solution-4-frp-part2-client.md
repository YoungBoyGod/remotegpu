# 方案4：frp 内网穿透 - 第二部分：客户端配置

## 第二部分：frp客户端配置

### 步骤4：安装frp客户端

**执行位置**：内网服务器（192.168.10.210）

```bash
# 下载frp
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz

# 解压
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64

# 移动到系统目录
sudo mkdir -p /opt/frp
sudo cp frpc /opt/frp/
sudo cp frpc.toml /opt/frp/
```

### 步骤5：配置frp客户端

**执行位置**：内网服务器（192.168.10.210）

```bash
# 编辑配置文件
sudo nano /opt/frp/frpc.toml
```

配置内容：
```toml
serverAddr = "公网服务器IP"  # 替换为frp服务端的公网IP
serverPort = 7000

# 认证（与服务端一致）
auth.method = "token"
auth.token = "your_secure_token_here"  # 与服务端相同

# 日志
log.to = "/var/log/frp/frpc.log"
log.level = "info"
log.maxDays = 7

# 代理配置 - Backend API
[[proxies]]
name = "remotegpu-backend"
type = "http"
localIP = "127.0.0.1"
localPort = 8080
customDomains = ["remotegpu.yourdomain.com"]  # 替换为您的域名
locations = ["/api"]

# 代理配置 - Frontend
[[proxies]]
name = "remotegpu-frontend"
type = "http"
localIP = "127.0.0.1"
localPort = 9980
customDomains = ["remotegpu.yourdomain.com"]  # 替换为您的域名
```

### 步骤6：创建systemd服务

**执行位置**：内网服务器（192.168.10.210）

```bash
# 创建日志目录
sudo mkdir -p /var/log/frp

# 创建服务文件
sudo nano /etc/systemd/system/frpc.service
```

服务文件内容：
```ini
[Unit]
Description=frp client
After=network.target

[Service]
Type=simple
ExecStart=/opt/frp/frpc -c /opt/frp/frpc.toml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl start frpc
sudo systemctl enable frpc
sudo systemctl status frpc
```

### 步骤7：配置DNS

**执行位置**：域名服务商控制台

添加A记录：
- 类型：A
- 主机记录：remotegpu
- 记录值：公网frp服务器IP
- TTL：600

### 步骤8：验证访问

**执行位置**：任意设备浏览器

访问 `http://remotegpu.yourdomain.com`

---

## 故障排查

```bash
# 查看服务端日志（公网服务器）
sudo tail -f /var/log/frp/frps.log

# 查看客户端日志（内网服务器）
sudo tail -f /var/log/frp/frpc.log

# 检查端口监听（公网服务器）
netstat -tlnp | grep frps

# 检查连接状态（内网服务器）
sudo systemctl status frpc
```

---

## 成本估算

- 公网服务器：约￥50-100/月（1核2G配置）
- 域名：约￥50-100/年
- 总计：约￥600-1200/年
