# 方案2：域名+SSL证书（生产环境推荐）

## 架构图
```
用户浏览器
    ↓
https://remotegpu.yourdomain.com
    ↓
nginx (Docker容器) + Let's Encrypt SSL
    ↓
backend(8080) + frontend(9980) (宿主机)
```

## 优点
- ✅ HTTPS加密，安全可靠
- ✅ 使用域名访问，专业美观
- ✅ 免费SSL证书（Let's Encrypt）
- ✅ 适合生产环境
- ✅ 浏览器显示安全锁

## 缺点
- ❌ 需要购买域名（约￥50-100/年）
- ❌ 需要配置DNS解析
- ❌ 实施时间稍长（30分钟）

## 前置条件
1. 购买域名（推荐：阿里云、腾讯云、Cloudflare）
2. 将域名A记录解析到 `45.78.48.169`

## 实施步骤

### 步骤1：配置DNS解析

在域名服务商添加A记录：
```
类型: A
主机记录: remotegpu (或 @)
记录值: 45.78.48.169
TTL: 600
```

等待DNS生效（通常5-10分钟）：
```bash
# 验证DNS解析
nslookup remotegpu.yourdomain.com
```

### 步骤2：安装certbot到Docker容器

```bash
# 进入nginx容器
docker exec -it remotegpu-nginx sh

# 安装certbot（Alpine Linux）
apk add --no-cache certbot certbot-nginx
```

### 步骤3：修改nginx配置

编辑 `/home/luo/code/remotegpu/docker-compose/nginx/conf.d/default.conf`：

```nginx
# HTTP服务器（用于证书验证和重定向）
server {
    listen 80;
    server_name remotegpu.yourdomain.com;

    # Let's Encrypt验证路径
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    # 其他请求重定向到HTTPS
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS服务器
server {
    listen 443 ssl http2;
    server_name remotegpu.yourdomain.com;

    # SSL证书（稍后由certbot自动配置）
    ssl_certificate /etc/letsencrypt/live/remotegpu.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/remotegpu.yourdomain.com/privkey.pem;

    # SSL安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000" always;

    # 日志
    access_log /var/log/nginx/remotegpu_access.log;
    error_log /var/log/nginx/remotegpu_error.log;

    # 上传大小限制
    client_max_body_size 10G;

    # 前端代理
    location / {
        proxy_pass http://host.docker.internal:9980;
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
        proxy_pass http://host.docker.internal:8080;
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

### 步骤4：获取SSL证书

```bash
# 在宿主机上执行
docker exec -it remotegpu-nginx certbot --nginx -d remotegpu.yourdomain.com

# 按提示输入邮箱和同意条款
```

### 步骤5：配置证书自动续期

```bash
# 添加定时任务（宿主机）
crontab -e

# 添加以下行（每天凌晨2点检查续期）
0 2 * * * docker exec remotegpu-nginx certbot renew --quiet
```

### 步骤6：验证访问

访问 `https://remotegpu.yourdomain.com`，应该看到：
- ✅ 浏览器显示安全锁
- ✅ 前端页面正常加载
- ✅ API正常工作

## 维护

证书有效期90天，自动续期任务会在到期前30天自动续期。
