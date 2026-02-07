# Nginx Proxy Manager 代理配置指南

## 概述

本文档介绍如何在 Nginx Proxy Manager (NPM) 中配置反向代理，将域名请求转发到 frp 暴露的端口。

### 配置目标

为每台 GPU 机器配置代理，例如：
- `https://gpu1-jupyter.gpu.domain.com` → frp 端口 7001
- `https://gpu1-tensorboard.gpu.domain.com` → frp 端口 7002
- `https://gpu2-jupyter.gpu.domain.com` → frp 端口 7003
- ...

---

## 前置要求

1. 已完成 NPM 安装（参考 [frp-npm-installation.md](frp-npm-installation.md)）
2. frps 服务已运行，GPU 机器的 frpc 已连接
3. DNS 泛域名解析已配置（`*.gpu.domain.com` → 云服务器 IP）

---

## 配置步骤

### 1. 登录 NPM 管理界面

访问 `http://your-server-ip:81`，使用管理员账号登录。

### 2. 添加代理主机（Proxy Host）

#### 2.1 进入代理主机页面

点击顶部菜单 **"Hosts"** → **"Proxy Hosts"**

#### 2.2 点击 "Add Proxy Host"

点击右上角的 **"Add Proxy Host"** 按钮。

#### 2.3 填写基本信息

在 **"Details"** 标签页中填写：

| 字段 | 值 | 说明 |
|------|-----|------|
| Domain Names | `gpu1-jupyter.gpu.domain.com` | 访问域名 |
| Scheme | `http` | frp 后端协议（通常是 http） |
| Forward Hostname / IP | `127.0.0.1` | frp 在本机 |
| Forward Port | `7001` | frp 暴露的端口 |
| Cache Assets | ☐ 不勾选 | 动态内容不缓存 |
| Block Common Exploits | ☑ 勾选 | 启用安全防护 |
| Websockets Support | ☑ 勾选 | Jupyter 需要 WebSocket |

#### 2.4 配置 SSL 证书

切换到 **"SSL"** 标签页：

| 字段 | 值 | 说明 |
|------|-----|------|
| SSL Certificate | `Request a new SSL Certificate` | 自动申请证书 |
| Force SSL | ☑ 勾选 | 强制 HTTPS |
| HTTP/2 Support | ☑ 勾选 | 启用 HTTP/2 |
| HSTS Enabled | ☑ 勾选 | 启用 HSTS |
| HSTS Subdomains | ☐ 不勾选 | 仅当前域名 |
| Use a DNS Challenge | ☐ 不勾选 | 使用 HTTP 验证 |
| Email Address | `your-email@example.com` | Let's Encrypt 通知邮箱 |
| I Agree to the Let's Encrypt Terms of Service | ☑ 勾选 | 同意服务条款 |

#### 2.5 保存配置

点击 **"Save"** 按钮，NPM 会自动：
1. 申请 Let's Encrypt SSL 证书
2. 配置 nginx 反向代理
3. 重载 nginx 配置

---

## 配置示例

### 示例 1：GPU1 的 Jupyter

**域名**: `gpu1-jupyter.gpu.domain.com`
**转发目标**: `http://127.0.0.1:7001`

**Details 标签页**:
```
Domain Names: gpu1-jupyter.gpu.domain.com
Scheme: http
Forward Hostname / IP: 127.0.0.1
Forward Port: 7001
☑ Block Common Exploits
☑ Websockets Support
```

**SSL 标签页**:
```
SSL Certificate: Request a new SSL Certificate
☑ Force SSL
☑ HTTP/2 Support
☑ HSTS Enabled
Email: admin@example.com
☑ I Agree to the Let's Encrypt Terms of Service
```

### 示例 2：GPU1 的 TensorBoard

**域名**: `gpu1-tensorboard.gpu.domain.com`
**转发目标**: `http://127.0.0.1:7002`

配置与示例 1 类似，只需修改：
- Domain Names: `gpu1-tensorboard.gpu.domain.com`
- Forward Port: `7002`

### 示例 3：GPU2 的 Jupyter

**域名**: `gpu2-jupyter.gpu.domain.com`
**转发目标**: `http://127.0.0.1:7003`

配置与示例 1 类似，只需修改：
- Domain Names: `gpu2-jupyter.gpu.domain.com`
- Forward Port: `7003`

---

## 高级配置

### 自定义 Nginx 配置

如果需要添加自定义 nginx 配置，可以在 **"Advanced"** 标签页中添加。

#### 示例：增加超时时间

```nginx
# 适用于长时间运行的请求（如 Jupyter Notebook）
proxy_read_timeout 3600s;
proxy_connect_timeout 3600s;
proxy_send_timeout 3600s;
```

#### 示例：自定义请求头

```nginx
# 传递真实客户端 IP
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
```

#### 示例：限制访问 IP

```nginx
# 仅允许特定 IP 访问
allow 1.2.3.4;
allow 5.6.7.0/24;
deny all;
```

### 配置访问控制

NPM 支持为代理主机添加访问控制（Access List）。

#### 创建访问控制列表

1. 点击顶部菜单 **"Access Lists"**
2. 点击 **"Add Access List"**
3. 填写名称（如 `Internal Only`）
4. 在 **"Authorization"** 中添加用户名和密码
5. 或在 **"Access"** 中配置 IP 白名单

#### 应用访问控制

编辑代理主机，在 **"Details"** 标签页中：
- **Access List**: 选择刚创建的访问控制列表

---

## SSL 证书管理

### 查看证书状态

1. 点击顶部菜单 **"SSL Certificates"**
2. 查看所有证书的状态和到期时间

### 证书自动续期

NPM 会自动续期即将到期的证书（到期前 30 天）。

### 手动续期证书

如果需要手动续期：
1. 进入 **"SSL Certificates"** 页面
2. 点击证书右侧的 **"..."** 菜单
3. 选择 **"Renew Certificate"**

### 使用通配符证书

如果有大量子域名，可以申请通配符证书（需要 DNS 验证）。

#### 申请通配符证书

1. 点击 **"SSL Certificates"** → **"Add SSL Certificate"**
2. 选择 **"Let's Encrypt"**
3. Domain Names 填写：`*.gpu.domain.com`
4. 勾选 **"Use a DNS Challenge"**
5. 选择 DNS 提供商（如 Cloudflare、阿里云等）
6. 填写 DNS API 凭据
7. 点击 **"Save"**

#### 使用通配符证书

创建代理主机时，在 **"SSL"** 标签页中：
- SSL Certificate: 选择已申请的通配符证书
- 不需要再次申请证书

---

## 验证配置

### 1. 检查代理主机状态

在 **"Proxy Hosts"** 页面，确认：
- Status 列显示绿色的 **"Online"**
- SSL 列显示绿色的锁图标

### 2. 测试 HTTPS 访问

在浏览器中访问：
```
https://gpu1-jupyter.gpu.domain.com
```

应该能够：
- 正常访问 Jupyter 页面
- 浏览器地址栏显示绿色锁图标（证书有效）
- 没有证书警告

### 3. 测试 HTTP 重定向

访问 HTTP 地址：
```
http://gpu1-jupyter.gpu.domain.com
```

应该自动重定向到 HTTPS。

### 4. 检查 SSL 证书

使用在线工具检查证书：
```
https://www.ssllabs.com/ssltest/analyze.html?d=gpu1-jupyter.gpu.domain.com
```

应该获得 A 或 A+ 评级。

---

## 常见问题

### SSL 证书申请失败

**问题**：点击 Save 后提示证书申请失败

**可能原因**：
1. DNS 解析未生效 - 等待 DNS 传播（最多 48 小时）
2. 80 端口未开放 - Let's Encrypt 需要通过 80 端口验证
3. 域名已有证书 - Let's Encrypt 限制每周申请次数

**解决方案**：
```bash
# 检查 DNS 解析
nslookup gpu1-jupyter.gpu.domain.com

# 检查 80 端口
sudo netstat -tlnp | grep :80

# 查看 NPM 日志
docker logs npm
```

### 无法访问后端服务

**问题**：HTTPS 可以访问，但显示 502 Bad Gateway

**可能原因**：
1. frp 端口配置错误
2. frpc 未连接到 frps
3. GPU 机器上的服务未启动

**解决方案**：
```bash
# 在云服务器上检查 frp 端口
curl http://127.0.0.1:7001

# 检查 frps 日志
docker logs frps

# 在 GPU 机器上检查 frpc 状态
systemctl status frpc
```

### WebSocket 连接失败

**问题**：Jupyter 页面可以打开，但无法执行代码

**解决方案**：
1. 确保代理主机配置中勾选了 **"Websockets Support"**
2. 在 **"Advanced"** 标签页添加：
```nginx
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
```

### 访问速度慢

**问题**：页面加载缓慢

**解决方案**：
1. 检查网络延迟：`ping gpu.domain.com`
2. 检查 frp 带宽限制
3. 在 **"Advanced"** 标签页添加缓存配置：
```nginx
proxy_buffering on;
proxy_buffer_size 4k;
proxy_buffers 8 4k;
```

---

## 下一步

单个代理配置完成后，继续阅读：

- **[NPM 批量配置方案](frp-npm-batch-config.md)** - 批量添加 200 台机器的配置
- **[NPM 完整实施指南](frp-npm-guide.md)** - 完整的 frp + NPM 方案

---

## 参考资源

- [NPM 代理主机文档](https://nginxproxymanager.com/guide/#proxy-hosts)
- [Let's Encrypt 速率限制](https://letsencrypt.org/docs/rate-limits/)
- [Nginx 反向代理配置](http://nginx.org/en/docs/http/ngx_http_proxy_module.html)
