# frp 服务端配置（云服务器）

## 安装 frps

**执行位置**：云服务器

```bash
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz

sudo mkdir -p /opt/frp
sudo cp frp_0.52.3_linux_amd64/frps /opt/frp/
```

---

## 配置 frps

**执行位置**：云服务器

```bash
sudo nano /opt/frp/frps.toml
```

配置内容：
```toml
bindPort = 7000

# 认证
auth.method = "token"
auth.token = "your_secure_token_change_this"

# 日志
log.to = "/var/log/frp/frps.log"
log.level = "info"

# 允许的端口范围
allowPorts = [
  { start = 2201, end = 2400 },   # SSH
  { start = 8001, end = 8200 },   # Jupyter
  { start = 9001, end = 9200 },   # TensorBoard
  { start = 10001, end = 10200 }, # 服务1
  { start = 11001, end = 11200 }  # 服务2
]
```

---

## 创建系统服务

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

---

## 配置防火墙

```bash
# frp控制端口
sudo ufw allow 7000/tcp

# 服务端口范围
sudo ufw allow 2201:2400/tcp
sudo ufw allow 8001:8200/tcp
sudo ufw allow 9001:9200/tcp
sudo ufw allow 10001:10200/tcp
sudo ufw allow 11001:11200/tcp
```
