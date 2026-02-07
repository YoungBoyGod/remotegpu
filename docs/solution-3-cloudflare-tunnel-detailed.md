# 方案3：Cloudflare Tunnel 详细实施流程

## 方案概述

**架构**：用户 → Cloudflare CDN → Cloudflare Tunnel → 内网服务（无需开放端口）

**优点**：
- 完全免费（包含SSL证书）
- 无需开放80/443端口
- 自动DDoS防护
- 全球CDN加速
- 零信任安全模型

**前置条件**：
- 需要域名（可使用Cloudflare免费域名，或自己的域名）
- 需要Cloudflare账号（免费）

---

## 实施步骤

### 步骤1：注册Cloudflare账号并添加域名

**执行位置**：浏览器

1. 访问 https://dash.cloudflare.com/sign-up
2. 注册免费账号
3. 添加您的域名（或使用Cloudflare提供的免费子域名）
4. 如果使用自己的域名，需要修改域名DNS服务器为Cloudflare提供的NS记录

---

### 步骤2：安装cloudflared客户端

**执行位置**：内网服务器（192.168.10.210）

```bash
# 下载cloudflared
cd /tmp
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb

# 安装
sudo dpkg -i cloudflared-linux-amd64.deb

# 验证安装
cloudflared --version
```

---

### 步骤3：登录Cloudflare账号

**执行位置**：内网服务器（192.168.10.210）

```bash
# 登录（会打开浏览器进行授权）
cloudflared tunnel login

# 授权后，会在 ~/.cloudflared/ 目录下生成证书文件
# 文件路径：~/.cloudflared/cert.pem
```

---

### 步骤4：创建Tunnel

**执行位置**：内网服务器（192.168.10.210）

```bash
# 创建tunnel（替换 remotegpu-tunnel 为您想要的名称）
cloudflared tunnel create remotegpu-tunnel

# 记录输出的Tunnel ID（类似：a1b2c3d4-e5f6-7890-abcd-ef1234567890）
# 会在 ~/.cloudflared/ 目录下生成 <tunnel-id>.json 文件
```

---

### 步骤5：配置Tunnel

**执行位置**：内网服务器（192.168.10.210）

创建配置文件：

```bash
# 创建配置目录（如果不存在）
mkdir -p ~/.cloudflared

# 创建配置文件
nano ~/.cloudflared/config.yml
```

配置文件内容：

```yaml
tunnel: <your-tunnel-id>  # 替换为步骤4中的Tunnel ID
credentials-file: /home/luo/.cloudflared/<your-tunnel-id>.json  # 替换为实际路径

ingress:
  # 前端服务（根路径）
  - hostname: remotegpu.yourdomain.com  # 替换为您的域名
    service: http://localhost:9980
    originRequest:
      noTLSVerify: true

  # 后端API服务
  - hostname: remotegpu.yourdomain.com
    path: /api/*
    service: http://localhost:8080
    originRequest:
      noTLSVerify: true

  # 默认规则（必须）
  - service: http_status:404
```

**配置说明**：
- `tunnel`: 您的Tunnel ID
- `credentials-file`: 凭证文件路径（绝对路径）
- `hostname`: 您的域名
- `service`: 内网服务地址（localhost:9980 和 localhost:8080）

---

### 步骤6：配置DNS记录

**执行位置**：内网服务器（192.168.10.210）

```bash
# 将域名指向tunnel
cloudflared tunnel route dns remotegpu-tunnel remotegpu.yourdomain.com

# 这会自动在Cloudflare DNS中创建CNAME记录
```

或者手动在Cloudflare控制台添加：
- 类型：CNAME
- 名称：remotegpu（或@）
- 目标：<tunnel-id>.cfargotunnel.com
- 代理状态：已代理（橙色云朵）

---

### 步骤7：测试Tunnel

**执行位置**：内网服务器（192.168.10.210）

```bash
# 前台运行测试
cloudflared tunnel run remotegpu-tunnel

# 如果看到 "Connection registered" 表示成功
# 在浏览器访问 https://remotegpu.yourdomain.com 测试
```

---

### 步骤8：配置为系统服务（开机自启）

**执行位置**：内网服务器（192.168.10.210）

```bash
# 安装为系统服务
sudo cloudflared service install

# 启动服务
sudo systemctl start cloudflared

# 设置开机自启
sudo systemctl enable cloudflared

# 查看服务状态
sudo systemctl status cloudflared

# 查看日志
sudo journalctl -u cloudflared -f
```

---

### 步骤9：验证访问

**执行位置**：任意设备浏览器

访问 `https://remotegpu.yourdomain.com`

应该看到：
- ✅ 自动HTTPS（绿色锁）
- ✅ 前端页面正常加载
- ✅ API正常工作
- ✅ 无需开放服务器端口

---

## 故障排查

### 问题1：Tunnel连接失败

```bash
# 检查cloudflared服务状态
sudo systemctl status cloudflared

# 查看详细日志
sudo journalctl -u cloudflared -n 100

# 检查配置文件语法
cloudflared tunnel ingress validate
```

### 问题2：502 Bad Gateway

**原因**：内网服务未运行或端口错误

```bash
# 检查backend和frontend是否运行
ps aux | grep "go run"
ps aux | grep "vite"

# 检查端口监听
netstat -tlnp | grep -E "8080|9980"
```

### 问题3：DNS解析错误

```bash
# 检查DNS记录
nslookup remotegpu.yourdomain.com

# 应该返回 CNAME 记录指向 *.cfargotunnel.com
```

---

## 维护命令

```bash
# 重启服务
sudo systemctl restart cloudflared

# 停止服务
sudo systemctl stop cloudflared

# 查看tunnel列表
cloudflared tunnel list

# 删除tunnel
cloudflared tunnel delete remotegpu-tunnel
```

---

## 优势总结

1. **零配置端口**：无需开放80/443端口，防火墙可以完全关闭
2. **自动HTTPS**：Cloudflare自动提供SSL证书
3. **全球加速**：通过Cloudflare CDN加速访问
4. **DDoS防护**：自动防御DDoS攻击
5. **隐藏源站IP**：攻击者无法获取真实服务器IP

---

## 成本

完全免费（Cloudflare Free计划）
