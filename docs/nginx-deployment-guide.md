# RemoteGPU Nginx 外网访问部署指南

## 前置条件

1. **公网服务器**：
   - 有公网 IP 地址
   - 开放 80 和 443 端口
   - 已安装 nginx（推荐版本 1.18+）

2. **内网服务器**：
   - 运行 RemoteGPU 服务（backend + frontend）
   - 公网服务器可以访问内网服务器

3. **域名**（可选但推荐）：
   - 已解析到公网服务器 IP

## 部署步骤

### 步骤 1：安装 Nginx（公网服务器）

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx

# 启动 nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 步骤 2：配置防火墙

```bash
# Ubuntu/Debian (ufw)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 步骤 3：获取 SSL 证书（推荐使用 Let's Encrypt）

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx  # Ubuntu/Debian
# 或
sudo yum install certbot python3-certbot-nginx  # CentOS/RHEL

# 获取证书（自动配置 nginx）
sudo certbot --nginx -d your-domain.com

# 证书会自动保存到 /etc/letsencrypt/live/your-domain.com/
```

### 步骤 4：修改配置文件

1. 复制配置文件到 nginx 目录：
```bash
sudo cp nginx-external-access.conf /etc/nginx/sites-available/remotegpu
```

2. 修改配置文件中的占位符：
```bash
sudo nano /etc/nginx/sites-available/remotegpu
```

需要修改的内容：
- `内网服务器IP` → 替换为实际的内网服务器 IP（如 192.168.1.100）
- `前端端口` → 替换为前端服务端口（如 5173 或 80）
- `your-domain.com` → 替换为您的域名
- SSL 证书路径（如果使用 certbot，路径会自动配置）

3. 创建符号链接启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/remotegpu /etc/nginx/sites-enabled/
```

4. 删除默认配置（可选）：
```bash
sudo rm /etc/nginx/sites-enabled/default
```

### 步骤 5：测试配置

```bash
# 测试 nginx 配置语法
sudo nginx -t

# 如果测试通过，重新加载 nginx
sudo systemctl reload nginx
```

### 步骤 6：验证访问

1. 在浏览器中访问：`https://your-domain.com`
2. 检查前端页面是否正常加载
3. 测试 API 访问：`https://your-domain.com/api/v1/health`

## 故障排查

### 问题 1：502 Bad Gateway

**原因**：nginx 无法连接到内网服务器

**解决方法**：
```bash
# 检查内网服务器是否可达
ping 内网服务器IP

# 检查端口是否开放
telnet 内网服务器IP 8080

# 检查 nginx 错误日志
sudo tail -f /var/log/nginx/remotegpu_error.log
```

### 问题 2：CORS 错误

**原因**：后端未正确配置 CORS

**解决方法**：在 nginx 配置中添加 CORS 头（已包含在配置文件中）

### 问题 3：SSL 证书错误

**原因**：证书路径不正确或证书过期

**解决方法**：
```bash
# 检查证书文件
sudo ls -la /etc/letsencrypt/live/your-domain.com/

# 续期证书
sudo certbot renew
```

## 安全建议

1. **定期更新证书**：Let's Encrypt 证书有效期 90 天，建议设置自动续期
2. **限制访问**：使用防火墙规则限制只允许必要的 IP 访问
3. **启用日志监控**：定期检查访问日志，发现异常访问
4. **使用强密码**：确保 RemoteGPU 平台使用强密码策略
5. **定期备份**：备份 nginx 配置和 SSL 证书

## 维护命令

```bash
# 重新加载配置（不中断服务）
sudo systemctl reload nginx

# 重启 nginx
sudo systemctl restart nginx

# 查看 nginx 状态
sudo systemctl status nginx

# 查看访问日志
sudo tail -f /var/log/nginx/remotegpu_access.log

# 查看错误日志
sudo tail -f /var/log/nginx/remotegpu_error.log
```
