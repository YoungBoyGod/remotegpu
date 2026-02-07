# 方案1：公网IP直接访问（最快方案）

## 架构图
```
用户浏览器
    ↓
http://45.78.48.169
    ↓
nginx (Docker容器)
    ↓
backend(8080) + frontend(9980) (宿主机)
```

## 优点
- ✅ 最快实施（10分钟内完成）
- ✅ 零成本（无需购买域名）
- ✅ 利用现有nginx Docker容器
- ✅ 适合快速测试和演示

## 缺点
- ❌ 无HTTPS加密（不安全）
- ❌ 使用IP访问（不专业）
- ❌ 浏览器会显示"不安全"警告
- ❌ 不适合生产环境

## 实施步骤

### 1. 修改nginx配置文件

编辑 `/home/luo/code/remotegpu/docker-compose/nginx/conf.d/default.conf`：

```nginx
server {
    listen 80;
    server_name 45.78.48.169;

    # 日志
    access_log /var/log/nginx/remotegpu_access.log;
    error_log /var/log/nginx/remotegpu_error.log;

    # 上传大小限制
    client_max_body_size 10G;

    # 前端代理（根路径）
    location / {
        # 使用宿主机IP访问宿主机服务
        proxy_pass http://host.docker.internal:9980;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://host.docker.internal:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # API超时配置
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
    }
}
```

### 2. 检查Docker网络配置

确保nginx容器可以访问宿主机：

```bash
# 进入nginx容器测试
docker exec -it remotegpu-nginx sh

# 测试访问宿主机服务
wget -O- http://host.docker.internal:8080/api/v1/health
wget -O- http://host.docker.internal:9980
```

如果 `host.docker.internal` 不可用，使用宿主机IP `192.168.10.210`：
```nginx
proxy_pass http://192.168.10.210:9980;
proxy_pass http://192.168.10.210:8080;
```

### 3. 重启nginx容器

```bash
cd /home/luo/code/remotegpu/docker-compose
docker-compose restart nginx
# 或
docker restart remotegpu-nginx
```

### 4. 验证访问

在浏览器中访问：
- 前端：`http://45.78.48.169`
- API健康检查：`http://45.78.48.169/api/v1/health`

## 故障排查

### 问题1：502 Bad Gateway
```bash
# 检查nginx日志
docker logs remotegpu-nginx

# 检查backend和frontend是否运行
ps aux | grep "go run"
ps aux | grep "vite"
```

### 问题2：CORS错误
在backend的配置中确保允许来自公网IP的请求。

## 后续升级路径

完成测试后，建议升级到方案2（添加域名和HTTPS）。
