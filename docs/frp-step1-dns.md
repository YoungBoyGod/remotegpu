# frp方案 - 第一步:DNS配置

## 目标

配置泛域名解析,让所有 `*.gpu.domain.com` 的子域名都指向云服务器IP。

---

## 前置准备

- 已有域名(如 `domain.com`)
- 有域名管理权限
- 知道云服务器的公网IP

---

## 配置步骤

### 1. 登录域名服务商管理后台

常见域名服务商:
- 阿里云(万网)
- 腾讯云DNSPod
- Cloudflare
- GoDaddy
- Namecheap

### 2. 添加DNS记录

**记录类型**: A记录

**配置示例**:

| 记录类型 | 主机记录 | 记录值 | TTL |
|---------|---------|--------|-----|
| A | *.gpu | 云服务器公网IP | 600 |

**说明**:
- **主机记录**: `*.gpu` (泛域名,匹配所有 `xxx.gpu.domain.com`)
- **记录值**: 云服务器的公网IP地址(如 `1.2.3.4`)
- **TTL**: 600秒(10分钟),可以设置更短以便测试

### 3. 不同服务商的具体操作

#### 阿里云(万网)

1. 登录 [阿里云控制台](https://dns.console.aliyun.com/)
2. 进入"云解析DNS" → 选择域名 → "解析设置"
3. 点击"添加记录"
4. 填写:
   - 记录类型: `A`
   - 主机记录: `*.gpu`
   - 解析线路: `默认`
   - 记录值: `云服务器IP`
   - TTL: `10分钟`
5. 点击"确认"

#### 腾讯云DNSPod

1. 登录 [DNSPod控制台](https://console.dnspod.cn/)
2. 选择域名 → "记录管理"
3. 点击"添加记录"
4. 填写:
   - 主机记录: `*.gpu`
   - 记录类型: `A`
   - 线路类型: `默认`
   - 记录值: `云服务器IP`
   - TTL: `600`
5. 点击"保存"

#### Cloudflare

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 选择域名 → "DNS" → "Records"
3. 点击"Add record"
4. 填写:
   - Type: `A`
   - Name: `*.gpu`
   - IPv4 address: `云服务器IP`
   - Proxy status: `DNS only` (关闭代理,显示灰色云朵)
   - TTL: `Auto`
5. 点击"Save"

**重要**: Cloudflare必须关闭代理(DNS only),否则SSL证书验证会失败!

---

## 验证DNS配置

### 方法1: 使用dig命令

```bash
# 测试泛域名解析
dig gpu1.gpu.domain.com
dig gpu2-jupyter.gpu.domain.com
dig test.gpu.domain.com

# 应该都返回云服务器IP
```

### 方法2: 使用nslookup命令

```bash
nslookup gpu1.gpu.domain.com
nslookup gpu2-jupyter.gpu.domain.com

# 应该都返回云服务器IP
```

### 方法3: 使用ping命令

```bash
ping gpu1.gpu.domain.com
ping gpu2-jupyter.gpu.domain.com

# 应该都能ping通云服务器IP
```

### 方法4: 在线DNS查询工具

- https://www.nslookup.io/
- https://dnschecker.org/
- https://mxtoolbox.com/

输入 `gpu1.gpu.domain.com`,查看是否解析到云服务器IP。

---

## 常见问题

### Q1: DNS配置后多久生效?

**答**:
- 新增记录: 通常5-10分钟
- 修改记录: 取决于TTL设置,最长可能需要24-48小时
- 建议: 配置后等待10-15分钟再进行验证

### Q2: 泛域名解析是否影响其他子域名?

**答**:
- 泛域名 `*.gpu` 只匹配 `xxx.gpu.domain.com` 格式
- 不影响 `www.domain.com`、`api.domain.com` 等其他子域名
- 如果已有 `test.gpu.domain.com` 的A记录,会优先使用具体记录

### Q3: 可以使用二级泛域名吗?

**答**:
- 可以,如 `*.*.gpu.domain.com`
- 但部分DNS服务商不支持
- 建议使用一级泛域名 `*.gpu.domain.com`

### Q4: Cloudflare的代理模式可以开启吗?

**答**:
- **不可以**!必须设置为"DNS only"(灰色云朵)
- 原因: 后续SSL证书申请需要验证域名所有权
- 如果开启代理,Let's Encrypt无法验证,会导致证书申请失败

---

## 配置示例

假设:
- 域名: `example.com`
- 云服务器IP: `123.45.67.89`

**DNS配置**:
```
类型: A
主机记录: *.gpu
记录值: 123.45.67.89
TTL: 600
```

**生效后**:
- `gpu1.gpu.example.com` → `123.45.67.89`
- `gpu2-jupyter.gpu.example.com` → `123.45.67.89`
- `gpu100-tensorboard.gpu.example.com` → `123.45.67.89`
- 所有 `*.gpu.example.com` 都指向 `123.45.67.89`

---

## 下一步

DNS配置完成并验证通过后,进入下一步:

👉 **第二步**: `frp-step2-ssl.md` - 获取SSL证书
