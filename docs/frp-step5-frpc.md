# frp方案 - 第五步:配置frpc客户端

## 目标

在GPU机器上安装frpc客户端,连接到云服务器的frps。

---

## 前置准备

- GPU机器可以访问外网
- 知道云服务器IP和frps端口(7000)
- 知道frps的token

---

## 下载和安装frpc

### 1. 下载frp

```bash
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64
```

### 2. 安装frpc

```bash
sudo cp frpc /usr/local/bin/
sudo chmod +x /usr/local/bin/frpc
sudo mkdir -p /etc/frp
```

---

## 配置frpc

### GPU1配置示例

创建配置文件:
```bash
sudo nano /etc/frp/frpc.ini
```

配置内容:
```ini
[common]
server_addr = 云服务器IP
server_port = 7000
authentication_method = token
token = your_secure_token_here

[gpu1-ssh]
type = tcp
local_ip = 127.0.0.1
local_port = 22
remote_port = 10001

[gpu1-jupyter]
type = tcp
local_ip = 127.0.0.1
local_port = 8888
remote_port = 11001

[gpu1-tensorboard]
type = tcp
local_ip = 127.0.0.1
local_port = 6006
remote_port = 12001

[gpu1-service1]
type = tcp
local_ip = 127.0.0.1
local_port = 本地端口1
remote_port = 13001

[gpu1-service2]
type = tcp
local_ip = 127.0.0.1
local_port = 本地端口2
remote_port = 14001
```

**配置说明**:
- `server_addr`: 云服务器公网IP
- `token`: 与frps配置相同
- `local_port`: GPU机器本地服务端口
- `remote_port`: 映射到云服务器的端口

---

## 创建systemd服务

```bash
sudo nano /etc/systemd/system/frpc.service
```

内容:
```ini
[Unit]
Description=frp client
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frpc -c /etc/frp/frpc.ini
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
```

---

## 启动frpc

```bash
sudo systemctl daemon-reload
sudo systemctl start frpc
sudo systemctl enable frpc
sudo systemctl status frpc
```

---

## 验证连接

### 1. 查看frpc日志

```bash
sudo journalctl -u frpc -f
```

正常输出:
```
[I] [service.go:xxx] login to server success
[I] [proxy_manager.go:xxx] proxy added: [gpu1-ssh gpu1-jupyter ...]
```

### 2. 在云服务器查看Dashboard

访问 `http://云服务器IP:7500`,应该看到GPU1的所有代理。

### 3. 测试端口

在云服务器上测试:
```bash
# 测试SSH端口
telnet 127.0.0.1 10001

# 测试Jupyter端口
curl http://127.0.0.1:11001
```

---

## GPU2-200配置

每台GPU机器的配置类似,只需修改:
1. 代理名称(gpu2-ssh, gpu3-ssh...)
2. remote_port(GPU2用10002, 11002, 12002...)

详见批量配置脚本: `frp-batch-scripts.md`

---

## 常见问题

### Q1: frpc连接失败?

排查:
```bash
# 检查网络连通性
ping 云服务器IP
telnet 云服务器IP 7000

# 检查token是否正确
grep token /etc/frp/frpc.ini
```

### Q2: 端口冲突?

确保remote_port在frps的allow_ports范围内(10000-15000)。

---

## 下一步

frpc配置完成后,进入下一步:

👉 **第六步**: `frp-step6-test.md` - 测试验证
