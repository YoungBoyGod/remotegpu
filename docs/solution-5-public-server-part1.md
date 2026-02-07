# 方案5：独立公网服务器 + nginx 反向代理

## 方案概述

**架构**：用户 → 公网nginx服务器 → 内网RemoteGPU服务器（192.168.10.210）

**适用场景**：
- 有独立的公网服务器
- 需要统一管理多个内网服务
- 需要更强的安全控制（WAF、限流等）

---

## 前置条件

1. **公网服务器**：有公网IP，已安装nginx
2. **网络连通性**：公网服务器能访问内网服务器（VPN/专线/公网IP）
3. **域名**：用于HTTPS访问

---

## 第一部分：公网服务器配置

### 步骤1：安装nginx

**执行位置**：公网服务器

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install nginx certbot python3-certbot-nginx

# 启动nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 步骤2：配置防火墙

**执行位置**：公网服务器

```bash
# Ubuntu/Debian
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# CentOS/RHEL
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 步骤3：配置DNS

**执行位置**：域名服务商控制台

添加A记录：
- 类型：A
- 主机记录：remotegpu
- 记录值：公网服务器IP
- TTL：600

### 步骤4：获取SSL证书

**执行位置**：公网服务器

```bash
# 使用certbot自动获取证书
sudo certbot --nginx -d remotegpu.yourdomain.com

# 按提示输入邮箱并同意条款
```

### 步骤5：配置nginx反向代理

**执行位置**：公网服务器

```bash
# 创建配置文件
sudo nano /etc/nginx/sites-available/remotegpu
```

配置内容：
```nginx
# HTTP重定向到HTTPS
server {
    listen 80;
    server_name remotegpu.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS服务器
server {
    listen 443 ssl http2;
    server_name remotegpu.yourdomain.com;

    # SSL证书（certbot自动配置）
    ssl_certificate /etc/letsencrypt/live/remotegpu.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/remotegpu.yourdomain.com/privkey.pem;

    # SSL安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 日志
    access_log /var/log/nginx/remotegpu_access.log;
    error_log /var/log/nginx/remotegpu_error.log;

    # 上传大小限制
    client_max_body_size 10G;

    # 前端代理
    location / {
        proxy_pass http://192.168.10.210:9980;  # 内网服务器IP
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://192.168.10.210:8080;  # 内网服务器IP
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
    }
}
```

**重要**：将 `192.168.10.210` 替换为实际的内网服务器地址。

### 步骤6：启用配置

**执行位置**：公网服务器

```bash
# 创建符号链接
sudo ln -s /etc/nginx/sites-available/remotegpu /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重新加载nginx
sudo systemctl reload nginx
```
