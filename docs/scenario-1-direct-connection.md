# 场景1：有固定公网 IP + 无防火墙限制

## 适用环境

- ✅ 有固定的公网 IP
- ✅ 可以开放端口
- ✅ 云服务器可以直接访问本地服务器

**典型环境**：企业专线、IDC 机房

---

## 架构图

```
用户（全国）
  ↓
域名 → 云服务器（nginx）
  ↓
直接连接（公网）
  ↓
本地服务器（固定IP）
  ↓
Backend + 200台 GPU
```

---

## 实施步骤

### 步骤1：确认本地公网 IP

**执行位置**：本地服务器

```bash
# 查看本地公网 IP
curl ifconfig.me

# 记录这个 IP，后续配置需要用到
```

---

### 步骤2：配置本地防火墙

**执行位置**：本地服务器

```bash
# 开放 Backend 端口（8080）
sudo ufw allow 8080/tcp

# 或限制只允许云服务器 IP 访问（更安全）
sudo ufw allow from 云服务器IP to any port 8080
```

---

### 步骤3：云服务器 nginx 配置

**执行位置**：云服务器

创建配置文件：
```bash
sudo nano /etc/nginx/sites-available/remotegpu
```

配置内容：
```nginx
# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS 服务器
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # 日志
    access_log /var/log/nginx/remotegpu_access.log;
    error_log /var/log/nginx/remotegpu_error.log;

    # 上传大小
    client_max_body_size 10G;

    # 后端 API 代理
    location /api/ {
        proxy_pass http://本地公网IP:8080;
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

**重要**：将 `本地公网IP` 替换为步骤1中获取的 IP。

---

### 步骤4：测试连接

**执行位置**：云服务器

```bash
# 测试能否访问本地 Backend
curl http://本地公网IP:8080/api/v1/health

# 应该返回 {"status":"ok"}
```

---

## 优点

- ✅ 配置最简单
- ✅ 性能最好（无隧道开销）
- ✅ 延迟最低

## 注意事项

- ⚠️ 需要固定 IP（如果 IP 变化需要更新配置）
- ⚠️ 需要开放端口（注意安全）
