# frp方案 - 第四步:配置nginx

## 目标

配置nginx作为反向代理,将HTTPS请求转发到frp的本地端口。

---

## 前置准备

- 已完成SSL证书申请
- 已安装并启动frps
- 已安装nginx

---

## 安装nginx

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx -y

# CentOS/RHEL
sudo yum install nginx -y

# 启动nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

---

## Web服务配置

### 1. 创建配置文件

```bash
sudo nano /etc/nginx/sites-available/gpu-frp
```

### 2. 配置内容(示例:GPU1-3)

```nginx
# GPU1 Jupyter
server {
    listen 443 ssl http2;
    server_name gpu1-jupyter.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:11001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# GPU1 TensorBoard
server {
    listen 443 ssl http2;
    server_name gpu1-tensorboard.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:12001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# GPU2-200 同理...
```

**端口对应关系**:
- Jupyter: 11001-11200 (GPU1-200)
- TensorBoard: 12001-12200
- 服务1: 13001-13200
- 服务2: 14001-14200

---

## 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/gpu-frp /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载nginx
sudo systemctl reload nginx
```

---

## SSH访问说明

<thinking>
用户强调要分段输出。我现在创建的nginx配置文档已经包含了基本的Web服务配置。

关于SSH访问,我需要说明:
1. SSH不能像Web服务那样通过nginx的http模块代理
2. SSH需要用户直接连接到云服务器的frp端口
3. 或者使用nginx的stream模块(但配置复杂)

我应该简要说明SSH的访问方式,并引用之前创建的ssh-vs-web-config.md文档。

让我继续完成这个文档,但要保持简洁。
</thinking>

**重要**: SSH和Web服务配置不同!

SSH访问方式:
```bash
ssh -p 10001 user@云服务器IP    # GPU1
ssh -p 10002 user@云服务器IP    # GPU2
```

详细说明请参考: `ssh-vs-web-config.md`

---

## 防火墙配置

```bash
# 开放HTTPS端口
sudo ufw allow 443/tcp

# 开放HTTP端口(用于重定向到HTTPS)
sudo ufw allow 80/tcp
```

---

## 验证配置

### 1. 检查nginx状态

```bash
sudo systemctl status nginx
```

### 2. 查看nginx日志

```bash
sudo tail -f /var/log/nginx/error.log
```

### 3. 测试域名解析

```bash
curl -I https://gpu1-jupyter.gpu.domain.com
```

---

## 批量生成配置

对于200台GPU机器,手动配置太繁琐,使用批量脚本生成。

详见: `frp-batch-scripts.md`

---

## 下一步

nginx配置完成后,进入下一步:

👉 **第五步**: `frp-step5-frpc.md` - 在GPU机器上配置frpc客户端
