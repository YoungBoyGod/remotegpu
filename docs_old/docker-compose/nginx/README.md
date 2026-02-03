# Nginx 反向代理

## 启动服务

```bash
docker-compose up -d
```

## 验证服务

```bash
# 检查配置
docker exec remotegpu-nginx nginx -t

# 重载配置
docker exec remotegpu-nginx nginx -s reload

# 访问测试
curl http://localhost
```

## 配置说明

### 添加反向代理

在 `conf.d/` 目录下创建新的配置文件：

```nginx
# conf.d/api.conf
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://backend:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### SSL/TLS 配置

将证书文件放在 `ssl/` 目录下，然后配置：

```nginx
server {
    listen 443 ssl;
    server_name example.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    # 其他配置...
}
```

## 查看日志

```bash
# 访问日志
docker exec remotegpu-nginx tail -f /var/log/nginx/access.log

# 错误日志
docker exec remotegpu-nginx tail -f /var/log/nginx/error.log
```

## 停止服务

```bash
docker-compose down
```
