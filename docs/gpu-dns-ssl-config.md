# DNS和SSL配置指南

## 第一步：DNS配置

### 配置泛域名

**执行位置**：域名服务商控制台

添加A记录：
```
类型：A
主机记录：*.gpu
记录值：云服务器公网IP
TTL：600
```

**效果**：
- gpu1.gpu.domain.com → 云服务器
- gpu1-jupyter.gpu.domain.com → 云服务器
- 所有 *.gpu.domain.com 都指向云服务器

### 验证DNS

**执行位置**：任意终端

```bash
# 测试解析
nslookup gpu1.gpu.domain.com
nslookup gpu1-jupyter.gpu.domain.com

# 应该都返回云服务器IP
```

---

## 第二步：SSL证书

### 获取泛域名证书

**执行位置**：云服务器

```bash
# 安装certbot
sudo apt install certbot -y

# 获取泛域名证书（DNS验证）
sudo certbot certonly --manual \
  --preferred-challenges dns \
  -d "*.gpu.domain.com"
```

### DNS验证步骤

1. certbot会要求添加TXT记录
2. 在域名服务商添加：
   ```
   类型：TXT
   主机记录：_acme-challenge.gpu
   记录值：certbot提供的值
   ```
3. 等待DNS生效（1-5分钟）
4. 按回车继续

### 证书位置

```
证书：/etc/letsencrypt/live/gpu.domain.com/fullchain.pem
密钥：/etc/letsencrypt/live/gpu.domain.com/privkey.pem
```

### 自动续期

```bash
# 添加定时任务
sudo crontab -e

# 每月1号凌晨2点续期
0 2 1 * * certbot renew --quiet
```

---

## 注意事项

- 泛域名证书覆盖所有子域名
- 有效期90天，需要定期续期
- DNS验证需要域名管理权限
