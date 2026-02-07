# frp方案 - 第三步:安装和配置frps服务端

## 目标

在云服务器上安装frp服务端(frps),接收来自GPU机器的连接。

---

## 前置准备

- 已完成DNS配置和SSL证书申请
- 云服务器有root权限
- 云服务器防火墙开放必要端口

---

## 下载frp

### 1. 访问frp官方仓库

https://github.com/fatedier/frp/releases

### 2. 下载最新版本

```bash
# 查看系统架构
uname -m
# x86_64 → amd64
# aarch64 → arm64

# 下载frp(以0.52.3版本为例,请使用最新版本)
cd /tmp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz

# 解压
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64
```

### 3. 安装frps

```bash
# 复制frps到系统目录
sudo cp frps /usr/local/bin/
sudo chmod +x /usr/local/bin/frps

# 创建配置目录
sudo mkdir -p /etc/frp

# 验证安装
frps --version
```

---

## 配置frps

### 1. 创建配置文件

```bash
sudo nano /etc/frp/frps.ini
```

### 2. 基础配置

```ini
[common]
# frps监听端口(客户端连接这个端口)
bind_port = 7000

# Dashboard配置(可选,用于查看连接状态)
dashboard_port = 7500
dashboard_user = admin
dashboard_pwd = your_password_here

# 认证token(客户端需要相同的token)
authentication_method = token
token = your_secure_token_here

# 日志配置
log_file = /var/log/frp/frps.log
log_level = info
log_max_days = 3

# 允许的端口范围
allow_ports = 10000-15000
```

<thinking>
用户强调要分段输出,避免一次性输出过大的内容。我现在创建的frp-step3-frps.md文档内容已经比较多了,我应该控制长度。

让我继续完成这个文档,但要注意不要写得太长。我需要包含:
1. 下载和安装frps (已完成)
2. 配置frps (正在进行)
3. 创建systemd服务
4. 启动和验证
5. 防火墙配置
6. 常见问题

我会尽量简洁,把重点内容写清楚即可。
</thinking>

**配置说明**:
- `bind_port`: frpc客户端连接端口,默认7000
- `dashboard_port`: Web管理界面端口
- `token`: 认证密钥,客户端必须使用相同token
- `allow_ports`: 允许客户端使用的端口范围

### 3. 生成安全token

```bash
# 生成随机token
openssl rand -base64 32
```

将生成的token填入配置文件的 `token` 字段。

---

## 创建systemd服务

### 1. 创建服务文件

```bash
sudo nano /etc/systemd/system/frps.service
```

### 2. 服务配置

```ini
[Unit]
Description=frp server
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frps -c /etc/frp/frps.ini
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
```

### 3. 创建日志目录

```bash
sudo mkdir -p /var/log/frp
```

---

## 启动frps

```bash
# 重载systemd配置
sudo systemctl daemon-reload

# 启动frps
sudo systemctl start frps

# 设置开机自启
sudo systemctl enable frps

# 查看状态
sudo systemctl status frps
```

**正常输出**:
```
● frps.service - frp server
   Loaded: loaded (/etc/systemd/system/frps.service; enabled)
   Active: active (running) since ...
```

---

## 验证frps

### 1. 查看日志

```bash
sudo tail -f /var/log/frp/frps.log
```

正常输出应包含:
```
[I] [service.go:xxx] frps started successfully
[I] [service.go:xxx] dashboard listen on 0.0.0.0:7500
```

### 2. 访问Dashboard

浏览器访问: `http://云服务器IP:7500`

- 用户名: `admin`
- 密码: 配置文件中设置的密码

### 3. 检查端口监听

```bash
sudo netstat -tlnp | grep frps
```

应该看到:
```
tcp6  0  0 :::7000  :::*  LISTEN  xxx/frps
tcp6  0  0 :::7500  :::*  LISTEN  xxx/frps
```

---

## 防火墙配置

### Ubuntu/Debian (ufw)

```bash
# 开放frps端口
sudo ufw allow 7000/tcp comment 'frps'
sudo ufw allow 7500/tcp comment 'frps dashboard'

# 开放frp代理端口范围
sudo ufw allow 10000:15000/tcp comment 'frp proxy ports'

# 重载防火墙
sudo ufw reload
```

### CentOS/RHEL (firewalld)

```bash
# 开放端口
sudo firewall-cmd --permanent --add-port=7000/tcp
sudo firewall-cmd --permanent --add-port=7500/tcp
sudo firewall-cmd --permanent --add-port=10000-15000/tcp

# 重载防火墙
sudo firewall-cmd --reload
```

### 云服务器安全组

如果使用阿里云/腾讯云等,还需要在控制台配置安全组规则:

- 入方向规则:
  - 端口: 7000 (frps)
  - 端口: 7500 (dashboard,可选)
  - 端口: 10000-15000 (代理端口)
  - 协议: TCP
  - 来源: 0.0.0.0/0

---

## 常见问题

### Q1: frps启动失败?

**排查步骤**:
```bash
# 查看详细日志
sudo journalctl -u frps -n 50

# 检查配置文件语法
frps verify -c /etc/frp/frps.ini

# 检查端口占用
sudo netstat -tlnp | grep 7000
```

### Q2: Dashboard无法访问?

**答**:
1. 检查防火墙是否开放7500端口
2. 检查云服务器安全组规则
3. 确认frps已启动: `sudo systemctl status frps`

### Q3: 如何修改配置?

**答**:
```bash
# 编辑配置
sudo nano /etc/frp/frps.ini

# 重启服务
sudo systemctl restart frps

# 查看日志确认
sudo tail -f /var/log/frp/frps.log
```

---

## 配置示例(完整)

```ini
[common]
bind_port = 7000
dashboard_port = 7500
dashboard_user = admin
dashboard_pwd = StrongPassword123!

authentication_method = token
token = AbCdEf1234567890XyZ

log_file = /var/log/frp/frps.log
log_level = info
log_max_days = 3

allow_ports = 10000-15000

# 可选:限制最大连接数
max_pool_count = 50

# 可选:心跳配置
heartbeat_timeout = 90
```

---

## 下一步

frps配置完成并验证通过后,进入下一步:

👉 **第四步**: `frp-step4-nginx.md` - 配置nginx反向代理
