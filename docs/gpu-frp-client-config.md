# frp 客户端配置（GPU机器）

## 配置示例（GPU1）

**执行位置**：GPU1 机器

### 安装 frpc

```bash
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz

sudo mkdir -p /opt/frp
sudo cp frp_0.52.3_linux_amd64/frpc /opt/frp/
```

### 配置文件

```bash
sudo nano /opt/frp/frpc.toml
```

配置内容（GPU1示例）：
```toml
serverAddr = "云服务器IP"
serverPort = 7000

auth.method = "token"
auth.token = "your_secure_token_change_this"

log.to = "/var/log/frp/frpc.log"
log.level = "info"

# SSH
[[proxies]]
name = "gpu1-ssh"
type = "tcp"
localIP = "127.0.0.1"
localPort = 22
remotePort = 2201

# Jupyter
[[proxies]]
name = "gpu1-jupyter"
type = "tcp"
localIP = "127.0.0.1"
localPort = 8888
remotePort = 8001

# TensorBoard
[[proxies]]
name = "gpu1-tensorboard"
type = "tcp"
localIP = "127.0.0.1"
localPort = 6006
remotePort = 9001

# 服务1（根据实际端口修改）
[[proxies]]
name = "gpu1-service1"
type = "tcp"
localIP = "127.0.0.1"
localPort = 本地端口1
remotePort = 10001

# 服务2（根据实际端口修改）
[[proxies]]
name = "gpu1-service2"
type = "tcp"
localIP = "127.0.0.1"
localPort = 本地端口2
remotePort = 11001
```

### 创建系统服务

```bash
sudo mkdir -p /var/log/frp

sudo tee /etc/systemd/system/frpc.service > /dev/null <<EOF
[Unit]
Description=frp client
After=network.target

[Service]
Type=simple
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

---

## GPU2-200配置

**端口规则**：
- GPU编号 N 的SSH端口：2200 + N
- GPU编号 N 的Jupyter端口：8000 + N
- GPU编号 N 的TensorBoard端口：9000 + N
- GPU编号 N 的服务1端口：10000 + N
- GPU编号 N 的服务2端口：11000 + N

**批量配置脚本见后续文档**
