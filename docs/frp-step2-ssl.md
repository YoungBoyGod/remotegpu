# frp方案 - 第二步:获取SSL证书

## 目标

使用Let's Encrypt获取泛域名SSL证书,为所有 `*.gpu.domain.com` 提供HTTPS支持。

---

## 前置准备

- 已完成DNS配置(第一步)
- DNS已生效(可以解析到云服务器IP)
- 云服务器已安装certbot

---

## 安装certbot

### Ubuntu/Debian

```bash
sudo apt update
sudo apt install certbot -y
```

### CentOS/RHEL

```bash
sudo yum install epel-release -y
sudo yum install certbot -y
```

### 验证安装

```bash
certbot --version
# 应该显示版本号,如: certbot 1.x.x
```

---

## 申请泛域名证书

### 方法1: DNS手动验证(推荐)

**适用场景**: 所有DNS服务商

```bash
sudo certbot certonly \
  --manual \
  --preferred-challenges dns \
  -d "*.gpu.domain.com" \
  -d "gpu.domain.com"
```

**说明**:
- `--manual`: 手动模式
- `--preferred-challenges dns`: 使用DNS验证
- `-d "*.gpu.domain.com"`: 泛域名证书
- `-d "gpu.domain.com"`: 同时包含主域名(可选)

### 执行过程

1. **输入邮箱**:
```
Enter email address (used for urgent renewal and security notices):
```
输入您的邮箱地址。

2. **同意服务条款**:
```
Please read the Terms of Service at https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf
(A)gree/(C)ancel:
```
输入 `A` 同意。

3. **DNS验证提示**:
```
Please deploy a DNS TXT record under the name:
_acme-challenge.gpu.domain.com

with the following value:
aBcDeFgHiJkLmNoPqRsTuVwXyZ1234567890

Before continuing, verify the TXT record has been deployed.
Press Enter to Continue
```

**重要**: 不要立即按Enter!

4. **添加DNS TXT记录**:

登录DNS服务商管理后台,添加TXT记录:

| 记录类型 | 主机记录 | 记录值 |
|---------|---------|--------|
| TXT | _acme-challenge.gpu | aBcDeFgHiJkLmNoPqRsTuVwXyZ1234567890 |

**注意**: 记录值使用certbot显示的实际值!

5. **验证TXT记录**:

在另一个终端窗口验证:
```bash
dig TXT _acme-challenge.gpu.domain.com

# 或使用nslookup
nslookup -type=TXT _acme-challenge.gpu.domain.com
```

确认返回正确的TXT记录值后,回到certbot窗口按Enter继续。

6. **等待验证完成**:
```
Successfully received certificate.
Certificate is saved at: /etc/letsencrypt/live/gpu.domain.com/fullchain.pem
Key is saved at:         /etc/letsencrypt/live/gpu.domain.com/privkey.pem
```

---

## 方法2: DNS自动验证(高级)

**适用场景**: DNS服务商支持API(如阿里云、腾讯云、Cloudflare)

### Cloudflare示例

1. **安装Cloudflare插件**:
```bash
sudo apt install python3-certbot-dns-cloudflare -y
```

2. **创建API Token**:
- 登录Cloudflare → My Profile → API Tokens
- Create Token → Edit zone DNS
- 保存Token

3. **创建配置文件**:
```bash
sudo mkdir -p /root/.secrets
sudo nano /root/.secrets/cloudflare.ini
```

内容:
```ini
dns_cloudflare_api_token = your_api_token_here
```

设置权限:
```bash
sudo chmod 600 /root/.secrets/cloudflare.ini
```

4. **申请证书**:
```bash
sudo certbot certonly \
  --dns-cloudflare \
  --dns-cloudflare-credentials /root/.secrets/cloudflare.ini \
  -d "*.gpu.domain.com" \
  -d "gpu.domain.com"
```

---

## 证书文件位置

证书申请成功后,文件保存在:

```
/etc/letsencrypt/live/gpu.domain.com/
├── fullchain.pem    # 完整证书链(nginx使用这个)
├── privkey.pem      # 私钥(nginx使用这个)
├── cert.pem         # 证书
└── chain.pem        # 证书链
```

**nginx配置使用**:
```nginx
ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;
```

---

## 验证证书

### 查看证书信息

```bash
sudo certbot certificates
```

输出示例:
```
Certificate Name: gpu.domain.com
  Domains: *.gpu.domain.com gpu.domain.com
  Expiry Date: 2026-05-07 12:34:56+00:00 (VALID: 89 days)
  Certificate Path: /etc/letsencrypt/live/gpu.domain.com/fullchain.pem
  Private Key Path: /etc/letsencrypt/live/gpu.domain.com/privkey.pem
```

### 测试证书文件

```bash
# 查看证书内容
sudo openssl x509 -in /etc/letsencrypt/live/gpu.domain.com/fullchain.pem -text -noout

# 验证私钥
sudo openssl rsa -in /etc/letsencrypt/live/gpu.domain.com/privkey.pem -check
```

---

## 证书自动续期

Let's Encrypt证书有效期90天,需要定期续期。

### 测试续期

```bash
sudo certbot renew --dry-run
```

如果输出 `Congratulations, all simulated renewals succeeded`,说明自动续期配置正确。

### 自动续期配置

certbot安装时会自动创建定时任务:

**查看定时任务**:
```bash
# systemd timer
sudo systemctl list-timers | grep certbot

# 或cron
sudo cat /etc/cron.d/certbot
```

**手动续期**:
```bash
sudo certbot renew
```

**续期后重启nginx**:

编辑 `/etc/letsencrypt/renewal-hooks/deploy/reload-nginx.sh`:
```bash
#!/bin/bash
systemctl reload nginx
```

设置权限:
```bash
sudo chmod +x /etc/letsencrypt/renewal-hooks/deploy/reload-nginx.sh
```

---

## 常见问题

### Q1: DNS验证一直失败?

**答**:
1. 确认DNS TXT记录已添加
2. 等待DNS传播(5-10分钟)
3. 使用 `dig TXT _acme-challenge.gpu.domain.com` 验证
4. 确认没有多余的TXT记录(删除旧的)

### Q2: Cloudflare代理模式导致验证失败?

**答**:
- 必须关闭Cloudflare代理(DNS only模式)
- 或使用DNS自动验证方法(方法2)

### Q3: 证书包含哪些域名?

**答**:
- 如果只申请 `-d "*.gpu.domain.com"`,只包含泛域名
- 建议同时申请 `-d "*.gpu.domain.com" -d "gpu.domain.com"`
- 这样 `gpu.domain.com` 和 `*.gpu.domain.com` 都可以使用

### Q4: 可以申请多个泛域名吗?

**答**:
- 可以,如: `-d "*.gpu.domain.com" -d "*.api.domain.com"`
- 但需要分别验证每个域名的DNS TXT记录

### Q5: 证书续期失败怎么办?

**答**:
1. 检查DNS TXT记录是否还存在
2. 手动执行 `sudo certbot renew --force-renewal`
3. 查看日志: `sudo tail -f /var/log/letsencrypt/letsencrypt.log`

---

## 配置示例

假设域名为 `example.com`,申请 `*.gpu.example.com` 证书:

```bash
# 1. 申请证书
sudo certbot certonly \
  --manual \
  --preferred-challenges dns \
  -d "*.gpu.example.com" \
  -d "gpu.example.com"

# 2. 按提示添加DNS TXT记录
# 记录类型: TXT
# 主机记录: _acme-challenge.gpu
# 记录值: (certbot显示的值)

# 3. 验证DNS
dig TXT _acme-challenge.gpu.example.com

# 4. 按Enter继续验证

# 5. 证书保存在
# /etc/letsencrypt/live/gpu.example.com/fullchain.pem
# /etc/letsencrypt/live/gpu.example.com/privkey.pem
```

---

## 下一步

SSL证书获取成功后,进入下一步:

👉 **第三步**: `frp-step3-frps.md` - 安装和配置frps服务端
