# nginx 配置（子域名匹配）

## 配置原理

nginx根据不同的子域名，代理到对应的frp端口。

---

## SSH服务配置（Stream模块）

**执行位置**：云服务器

### 启用Stream模块

编辑 `/etc/nginx/nginx.conf`，在http块外添加：

```nginx
stream {
    # GPU1 SSH
    server {
        listen 22;
        server_name gpu1.gpu.domain.com;
        proxy_pass 127.0.0.1:2201;
    }

    # GPU2 SSH
    server {
        listen 22;
        server_name gpu2.gpu.domain.com;
        proxy_pass 127.0.0.1:2202;
    }

    # ... GPU3-200 同理
}
```

**注意**：SSH需要使用stream模块，不能用http模块。

---

## Web服务配置（HTTP模块）

**执行位置**：云服务器

创建配置文件：
```bash
sudo nano /etc/nginx/sites-available/gpu-services
```

配置内容（示例）：
```nginx
# GPU1 Jupyter
server {
    listen 443 ssl http2;
    server_name gpu1-jupyter.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# GPU1 TensorBoard
server {
    listen 443 ssl http2;
    server_name gpu1-tensorboard.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:9001;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }
}

# GPU2-200 同理...
```

**批量生成配置见后续脚本**
