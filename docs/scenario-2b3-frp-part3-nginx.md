# 场景2B/3：nginx 配置（frp 场景）

## 第三部分：云服务器 nginx 配置

### 步骤9：安装 nginx 和 certbot

**执行位置**：云服务器

```bash
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx -y
```

### 步骤10：配置 DNS

**执行位置**：域名服务商控制台

添加 A 记录：
- 类型：A
- 主机记录：remotegpu（或您想要的子域名）
- 记录值：云服务器公网 IP
- TTL：600

等待 DNS 生效（5-10分钟）：
```bash
nslookup remotegpu.yourdomain.com
```

### 步骤11：获取 SSL 证书

**执行位置**：云服务器

```bash
sudo certbot certonly --nginx -d remotegpu.yourdomain.com
```

### 步骤12：配置 nginx

**执行位置**：云服务器

```bash
sudo nano /etc/nginx/sites-available/remotegpu
```

配置内容：
```nginx
# HTTP 重定向
server {
    listen 80;
    server_name remotegpu.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS 服务器
server {
    listen 443 ssl http2;
    server_name remotegpu.yourdomain.com;

    # SSL 证书
    ssl_certificate /etc/letsencrypt/live/remotegpu.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/remotegpu.yourdomain.com/privkey.pem;

    # 日志
    access_log /var/log/nginx/remotegpu_access.log;
    error_log /var/log/nginx/remotegpu_error.log;

    # 上传大小
    client_max_body_size 10G;

    # 后端 API 代理（通过 frp）
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
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

### 步骤13：启用配置

**执行位置**：云服务器

```bash
# 创建符号链接
sudo ln -s /etc/nginx/sites-available/remotegpu /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重新加载
sudo systemctl reload nginx
```

### 步骤14：配置证书自动续期

**执行位置**：云服务器

```bash
# 添加定时任务
sudo crontab -e

# 添加以下行
0 2 * * * certbot renew --quiet && systemctl reload nginx
```

---

## 验证

**执行位置**：任意设备浏览器

访问：`https://remotegpu.yourdomain.com/api/v1/health`

应该返回：`{"status":"ok"}`
