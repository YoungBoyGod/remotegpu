# 阿里云DNS解析配置方案

## 一、方案选择

### 方案A：通配符解析（推荐）

**适用场景：**
- 所有环境共享负载均衡器IP
- 通过反向代理（Nginx/Traefik）路由到不同容器
- 简单、稳定、成本低

**优点：**
- ✅ 一次配置，永久生效
- ✅ 无需API调用
- ✅ 无额外费用
- ✅ 配置简单

**缺点：**
- ❌ 所有子域名指向同一IP
- ❌ 需要反向代理做二次路由

### 方案B：API动态管理

**适用场景：**
- 每个环境需要独立IP
- 需要精细控制DNS记录
- 需要自动化管理

**优点：**
- ✅ 灵活控制每个域名
- ✅ 支持独立IP
- ✅ 完全自动化

**缺点：**
- ❌ 需要开发和维护
- ❌ 依赖API可用性
- ❌ 有API调用费用

---

## 二、方案A：通配符解析配置（推荐）

### 2.1 阿里云控制台配置

**步骤1：登录阿里云控制台**
```
访问：https://dns.console.aliyun.com/
```

**步骤2：选择域名**
```
1. 在域名列表中找到 remotegpu.com
2. 点击"解析设置"
```

**步骤3：添加通配符记录**
```
记录类型：A
主机记录：*
解析线路：默认
记录值：47.100.1.100（你的负载均衡器公网IP）
TTL：10分钟（600秒）
```

**配置示例：**
```
主机记录    记录类型    解析线路    记录值           TTL
*          A          默认       47.100.1.100     600
@          A          默认       47.100.1.100     600
www        A          默认       47.100.1.100     600
```

**步骤4：验证配置**
```bash
# 等待5-10分钟后验证
nslookup env-test.remotegpu.com
nslookup env-abc123.remotegpu.com

# 应该都返回：47.100.1.100
```

### 2.2 Nginx反向代理配置

**nginx.conf配置：**
```nginx
# HTTP重定向到HTTPS
server {
    listen 80;
    server_name *.remotegpu.com;
    return 301 https://$host$request_uri;
}

# HTTPS主配置
server {
    listen 443 ssl http2;
    server_name ~^env-(?<env_id>[a-z0-9-]+)\.remotegpu\.com$;

    # SSL证书（通配符证书）
    ssl_certificate /etc/nginx/ssl/remotegpu.com.crt;
    ssl_certificate_key /etc/nginx/ssl/remotegpu.com.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 日志
    access_log /var/log/nginx/env-$env_id-access.log;
    error_log /var/log/nginx/env-$env_id-error.log;

    # Jupyter服务
    location /jupyter {
        # 根据env_id查找后端容器
        set $backend "http://env-$env_id:8888";
        proxy_pass $backend;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400;
    }

    # Code Server
    location /code {
        set $backend "http://env-$env_id:8080";
        proxy_pass $backend;

        proxy_set_header Host $host;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 默认路由
    location / {
        return 200 "Environment: $env_id\n";
        add_header Content-Type text/plain;
    }
}
```

**SSH端口转发（stream模块）：**
```nginx
# 在nginx.conf的顶层添加stream块
stream {
    # 定义上游服务器映射
    map $ssl_preread_server_name $backend_ssh {
        ~^env-(?<env_id>[a-z0-9-]+)\.remotegpu\.com$ env-$env_id:22;
        default 127.0.0.1:22;
    }

    # SSH端口转发
    server {
        listen 22;
        proxy_pass $backend_ssh;
        proxy_protocol on;
    }
}
```

### 2.3 Docker网络配置

**docker-compose.yml示例：**
```yaml
version: '3.8'

services:
  nginx:
    image: nginx:latest
    ports:
      - "80:80"
      - "443:443"
      - "22:22"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    networks:
      - remotegpu-network

  # 环境容器示例
  env-abc123:
    image: remotegpu/workspace:latest
    hostname: env-abc123
    environment:
      - USERNAME=remotegpu
      - PASSWORD=${PASSWORD}
      - JUPYTER_TOKEN=${JUPYTER_TOKEN}
    networks:
      - remotegpu-network

networks:
  remotegpu-network:
    driver: bridge
```

### 2.4 SSL证书申请

**使用Let's Encrypt申请通配符证书：**
```bash
# 安装certbot
apt-get install certbot

# 申请通配符证书（需要DNS验证）
certbot certonly --manual \
  --preferred-challenges=dns \
  --email admin@remotegpu.com \
  --server https://acme-v02.api.letsencrypt.org/directory \
  --agree-tos \
  -d remotegpu.com \
  -d *.remotegpu.com

# 按照提示在阿里云DNS添加TXT记录
# 记录类型：TXT
# 主机记录：_acme-challenge
# 记录值：（certbot提供的值）

# 证书路径
# /etc/letsencrypt/live/remotegpu.com/fullchain.pem
# /etc/letsencrypt/live/remotegpu.com/privkey.pem

# 自动续期
certbot renew --dry-run
```

**阿里云DNS添加TXT记录：**
```
记录类型：TXT
主机记录：_acme-challenge
记录值：xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TTL：600
```

---

## 三、方案B：API动态管理实现

### 3.1 阿里云SDK集成

**安装依赖：**
```bash
go get github.com/aliyun/alibaba-cloud-sdk-go/services/alidns
```

**配置AccessKey：**
```yaml
# config/config.yaml
aliyun:
  access_key_id: "LTAI5t..."
  access_key_secret: "xxxxxxxxxxxxx"
  region_id: "cn-hangzhou"
  dns:
    domain: "remotegpu.com"
```

### 3.2 DNS管理器实现

**pkg/dns/aliyun_dns.go：**
```go
package dns

import (
    "fmt"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// AliyunDNSManager 阿里云DNS管理器
type AliyunDNSManager struct {
    client *alidns.Client
    domain string
}

// NewAliyunDNSManager 创建DNS管理器
func NewAliyunDNSManager(accessKeyID, accessKeySecret, regionID, domain string) (*AliyunDNSManager, error) {
    client, err := alidns.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
    if err != nil {
        return nil, fmt.Errorf("创建阿里云DNS客户端失败: %w", err)
    }

    return &AliyunDNSManager{
        client: client,
        domain: domain,
    }, nil
}

// AddRecord 添加DNS A记录
func (m *AliyunDNSManager) AddRecord(subdomain, ip string) (string, error) {
    request := alidns.CreateAddDomainRecordRequest()
    request.Scheme = "https"
    request.DomainName = m.domain
    request.RR = subdomain // 如：env-abc123
    request.Type = "A"
    request.Value = ip
    request.TTL = "600"

    response, err := m.client.AddDomainRecord(request)
    if err != nil {
        return "", fmt.Errorf("添加DNS记录失败: %w", err)
    }

    return response.RecordId, nil
}

// DeleteRecord 删除DNS记录
func (m *AliyunDNSManager) DeleteRecord(recordID string) error {
    request := alidns.CreateDeleteDomainRecordRequest()
    request.Scheme = "https"
    request.RecordId = recordID

    _, err := m.client.DeleteDomainRecord(request)
    if err != nil {
        return fmt.Errorf("删除DNS记录失败: %w", err)
    }

    return nil
}

// UpdateRecord 更新DNS记录
func (m *AliyunDNSManager) UpdateRecord(recordID, subdomain, ip string) error {
    request := alidns.CreateUpdateDomainRecordRequest()
    request.Scheme = "https"
    request.RecordId = recordID
    request.RR = subdomain
    request.Type = "A"
    request.Value = ip
    request.TTL = "600"

    _, err := m.client.UpdateDomainRecord(request)
    if err != nil {
        return fmt.Errorf("更新DNS记录失败: %w", err)
    }

    return nil
}

// GetRecord 查询DNS记录
func (m *AliyunDNSManager) GetRecord(subdomain string) (*alidns.Record, error) {
    request := alidns.CreateDescribeDomainRecordsRequest()
    request.Scheme = "https"
    request.DomainName = m.domain
    request.RRKeyWord = subdomain
    request.TypeKeyWord = "A"

    response, err := m.client.DescribeDomainRecords(request)
    if err != nil {
        return nil, fmt.Errorf("查询DNS记录失败: %w", err)
    }

    if len(response.DomainRecords.Record) == 0 {
        return nil, fmt.Errorf("DNS记录不存在")
    }

    return &response.DomainRecords.Record[0], nil
}

// ListRecords 列出所有DNS记录
func (m *AliyunDNSManager) ListRecords() ([]alidns.Record, error) {
    request := alidns.CreateDescribeDomainRecordsRequest()
    request.Scheme = "https"
    request.DomainName = m.domain
    request.PageSize = "100"

    response, err := m.client.DescribeDomainRecords(request)
    if err != nil {
        return nil, fmt.Errorf("列出DNS记录失败: %w", err)
    }

    return response.DomainRecords.Record, nil
}
```

### 3.3 集成到环境服务

**internal/service/environment.go：**
```go
// EnvironmentService 添加DNS管理器
type EnvironmentService struct {
    // ... 其他字段
    dnsManager *dns.AliyunDNSManager
}

// CreateEnvironment 创建环境时添加DNS记录
func (s *EnvironmentService) CreateEnvironment(req *CreateEnvironmentRequest) (*Environment, error) {
    // 1. 生成环境ID
    envID := fmt.Sprintf("env-%s", uuid.New().String()[:8])

    // 2. 创建容器并获取IP
    container, err := s.deploymentManager.CreateContainer(...)
    if err != nil {
        return nil, err
    }

    containerIP := container.NetworkSettings.IPAddress

    // 3. 添加DNS记录
    recordID, err := s.dnsManager.AddRecord(envID, containerIP)
    if err != nil {
        // 回滚容器创建
        s.deploymentManager.DeleteContainer(container.ID)
        return nil, fmt.Errorf("添加DNS记录失败: %w", err)
    }

    // 4. 保存环境信息
    env := &entity.Environment{
        ID:          envID,
        Domain:      fmt.Sprintf("%s.remotegpu.com", envID),
        DNSRecordID: recordID,
        ContainerID: container.ID,
        // ... 其他字段
    }

    err = s.envDao.Create(env)
    if err != nil {
        // 回滚DNS记录
        s.dnsManager.DeleteRecord(recordID)
        s.deploymentManager.DeleteContainer(container.ID)
        return nil, err
    }

    return env, nil
}

// DeleteEnvironment 删除环境时删除DNS记录
func (s *EnvironmentService) DeleteEnvironment(envID string) error {
    env, err := s.envDao.GetByID(envID)
    if err != nil {
        return err
    }

    // 1. 删除DNS记录
    if env.DNSRecordID != "" {
        err = s.dnsManager.DeleteRecord(env.DNSRecordID)
        if err != nil {
            // 记录错误但继续删除
            log.Printf("删除DNS记录失败: %v", err)
        }
    }

    // 2. 删除容器
    err = s.deploymentManager.DeleteContainer(env.ContainerID)
    if err != nil {
        return err
    }

    // 3. 删除数据库记录
    return s.envDao.Delete(envID)
}
```

### 3.4 数据库表扩展

```sql
-- 添加DNS记录ID字段
ALTER TABLE environments ADD COLUMN dns_record_id VARCHAR(64);
CREATE INDEX idx_environments_dns_record ON environments(dns_record_id);
```

### 3.5 错误处理和重试

```go
// AddRecordWithRetry 添加DNS记录（带重试）
func (m *AliyunDNSManager) AddRecordWithRetry(subdomain, ip string, maxRetries int) (string, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        recordID, err := m.AddRecord(subdomain, ip)
        if err == nil {
            return recordID, nil
        }

        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }

    return "", fmt.Errorf("添加DNS记录失败（重试%d次）: %w", maxRetries, lastErr)
}
```

---

## 四、最佳实践

### 4.1 推荐配置

**对于RemoteGPU项目，推荐使用方案A（通配符解析）：**

1. **阿里云DNS配置：**
   - 添加通配符A记录：`* → 负载均衡器IP`
   - TTL设置为600秒（10分钟）

2. **负载均衡器：**
   - 使用阿里云SLB或自建Nginx
   - 配置健康检查
   - 启用HTTPS

3. **反向代理：**
   - Nginx/Traefik根据域名路由
   - 支持WebSocket
   - 配置SSL证书

### 4.2 性能优化

**DNS缓存：**
```go
// 在应用层缓存DNS查询结果
type DNSCache struct {
    cache map[string]string
    mu    sync.RWMutex
    ttl   time.Duration
}

func (c *DNSCache) Get(domain string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    ip, ok := c.cache[domain]
    return ip, ok
}
```

**连接池：**
```go
// 复用阿里云SDK客户端
var (
    dnsClient *alidns.Client
    once      sync.Once
)

func GetDNSClient() *alidns.Client {
    once.Do(func() {
        dnsClient, _ = alidns.NewClientWithAccessKey(...)
    })
    return dnsClient
}
```

### 4.3 监控和告警

**监控指标：**
- DNS解析成功率
- DNS解析延迟
- API调用失败次数

**告警规则：**
```yaml
# Prometheus告警规则
groups:
  - name: dns
    rules:
      - alert: DNSResolveFailed
        expr: dns_resolve_errors_total > 10
        for: 5m
        annotations:
          summary: "DNS解析失败次数过多"
```

---

## 五、故障排查

### 5.1 常见问题

**问题1：域名无法解析**
```bash
# 检查DNS配置
nslookup env-test.remotegpu.com

# 检查阿里云DNS记录
dig @223.5.5.5 env-test.remotegpu.com

# 清除本地DNS缓存
# Linux
systemd-resolve --flush-caches
# macOS
dscacheutil -flushcache
```

**问题2：SSL证书错误**
```bash
# 检查证书有效期
openssl x509 -in /etc/nginx/ssl/remotegpu.com.crt -noout -dates

# 测试SSL连接
openssl s_client -connect env-test.remotegpu.com:443
```

**问题3：API调用失败**
```go
// 启用SDK调试日志
client.SetDebug(true)

// 检查AccessKey权限
// 需要AliyunDNSFullAccess权限
```

### 5.2 日志分析

**Nginx访问日志：**
```bash
# 查看特定环境的访问日志
tail -f /var/log/nginx/env-abc123-access.log

# 统计访问量
awk '{print $1}' access.log | sort | uniq -c | sort -rn
```

**阿里云API日志：**
```go
// 记录API调用日志
log.Printf("DNS API调用: 操作=%s, 域名=%s, IP=%s", operation, subdomain, ip)
```

---

## 六、成本估算

### 6.1 方案A成本

**阿里云DNS：**
- 域名解析：免费（基础版）
- 通配符记录：无额外费用

**总成本：** 0元/月

### 6.2 方案B成本

**阿里云DNS API：**
- API调用：0.01元/万次
- 假设每天创建100个环境：100次 × 30天 = 3000次/月
- 成本：约0.3元/月

**总成本：** 可忽略不计

---

## 七、总结

**推荐方案：方案A（通配符解析）**

**理由：**
1. ✅ 配置简单，一次配置永久生效
2. ✅ 无需开发和维护API代码
3. ✅ 成本为零
4. ✅ 稳定可靠，不依赖API
5. ✅ 符合RemoteGPU的架构设计

**实施步骤：**
1. 在阿里云DNS添加通配符A记录
2. 配置Nginx反向代理
3. 申请SSL通配符证书
4. 测试验证

**预计时间：** 1-2小时即可完成配置
