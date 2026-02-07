# 云服务器 nginx 配置（企业直连方案）

## 配置原理

nginx根据子域名，直接代理到企业公网IP的对应端口。

---

## SSH访问说明

**重要**: SSH和Web服务的配置完全不同!

- SSH是四层协议,nginx的stream模块**不支持**根据域名路由
- **推荐方案**: 用户直接SSH到企业公网IP的对应端口(不经过云服务器)
- 访问方式: `ssh -p 2201 user@企业公网IP` (GPU1)

**详细说明和配置方案请参考**: `ssh-vs-web-config.md`

---

## Web服务配置（HTTP模块）

**执行位置**：云服务器

```bash
sudo nano /etc/nginx/sites-available/gpu-enterprise
```

配置内容：
```nginx
# GPU1 Jupyter
server {
    listen 443 ssl http2;
    server_name gpu1-jupyter.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://企业公网IP:8001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }
}

# GPU1 TensorBoard
server {
    listen 443 ssl http2;
    server_name gpu1-tensorboard.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://企业公网IP:9001;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }
}
```

**批量生成见后续脚本**

---

## 启用配置

```bash
sudo ln -s /etc/nginx/sites-available/gpu-enterprise /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```
