# 场景2B/3：frp 内网穿透 - 第二部分：本地配置

## 第二部分：本地服务器配置

### 步骤5：安装 frpc

**执行位置**：本地服务器（192.168.10.210）

```bash
# 下载 frp
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz

# 安装
sudo mkdir -p /opt/frp
sudo cp frp_0.52.3_linux_amd64/frpc /opt/frp/
```

### 步骤6：配置 frpc

**执行位置**：本地服务器

```bash
sudo nano /opt/frp/frpc.toml
```

配置内容：
```toml
serverAddr = "云服务器IP"
serverPort = 7000

auth.method = "token"
auth.token = "your_secure_token_123"  # 与服务端相同

log.to = "/var/log/frp/frpc.log"
log.level = "info"

# Backend API 代理
[[proxies]]
name = "remotegpu-backend"
type = "tcp"
localIP = "127.0.0.1"
localPort = 8080
remotePort = 8080
```

**重要**：将 `云服务器IP` 替换为实际的云服务器公网 IP。

### 步骤7：创建系统服务

**执行位置**：本地服务器

```bash
sudo mkdir -p /var/log/frp

sudo tee /etc/systemd/system/frpc.service > /dev/null <<EOF
[Unit]
Description=frp client
After=network.target

[Service]
Type=simple
User=luo
ExecStart=/opt/frp/frpc -c /opt/frp/frpc.toml
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start frpc
sudo systemctl enable frpc
```

### 步骤8：验证连接

**执行位置**：本地服务器

```bash
# 查看 frpc 状态
sudo systemctl status frpc

# 查看日志
sudo tail -f /var/log/frp/frpc.log

# 应该看到 "login to server success"
```

**执行位置**：云服务器

```bash
# 测试通过 frp 访问本地 Backend
curl http://127.0.0.1:8080/api/v1/health

# 应该返回 {"status":"ok"}
```

---

## 故障排查

### 问题1：frpc 无法连接

```bash
# 检查云服务器 IP 是否正确
ping 云服务器IP

# 检查 7000 端口是否开放
telnet 云服务器IP 7000

# 检查 token 是否一致
```

### 问题2：Backend 无法访问

```bash
# 检查 Backend 是否运行
ps aux | grep "go run"

# 检查端口监听
netstat -tlnp | grep 8080
```

---

## 优点

- ✅ 解决 NAT 问题
- ✅ 无需固定 IP
- ✅ 自动重连
- ✅ 加密传输

## 注意事项

- ⚠️ frpc 需要一直运行
- ⚠️ 有轻微性能损耗（通常 <10ms）
