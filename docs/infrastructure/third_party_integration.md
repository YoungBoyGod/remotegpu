# 第三方系统集成文档

> 本文档描述 RemoteGPU 系统与第三方系统的集成方案和配置
>
> 创建日期：2026-01-26

---

## 1. 集成架构概览

```yaml
集成分类:
  认证集成:
    - OAuth2 (GitHub, Google, GitLab)
    - LDAP/AD (企业目录服务)
    - SAML 2.0 (企业单点登录)

  通知集成:
    - 钉钉机器人
    - 企业微信机器人
    - Slack Webhook
    - 飞书机器人

  存储集成:
    - AWS S3
    - 阿里云 OSS
    - 腾讯云 COS
    - MinIO

  监控集成:
    - Prometheus
    - Grafana
    - Uptime Kuma
    - Sentry (错误追踪)

  CI/CD 集成:
    - GitHub Actions
    - GitLab CI
    - Jenkins
```

---

## 2. OAuth2 认证集成

### 2.1 GitHub OAuth

**配置步骤：**
1. 在 GitHub 创建 OAuth App
2. 获取 Client ID 和 Client Secret
3. 配置回调 URL

**配置文件：**
```yaml
# config/oauth.yml
github:
  client_id: your_github_client_id
  client_secret: your_github_client_secret
  redirect_uri: https://remotegpu.example.com/auth/github/callback
  scopes:
    - user:email
    - read:user
```

---

**实现代码：**
```go
// GitHub OAuth 处理器
func HandleGitHubCallback(c *gin.Context) {
    code := c.Query("code")
    
    // 交换 access token
    token, err := exchangeGitHubToken(code)
    if err != nil {
        c.JSON(500, gin.H{"error": "获取 token 失败"})
        return
    }
    
    // 获取用户信息
    user, err := getGitHubUser(token)
    if err != nil {
        c.JSON(500, gin.H{"error": "获取用户信息失败"})
        return
    }
    
    // 创建或更新用户
    customer := createOrUpdateCustomer(user)
    
    // 生成 JWT token
    jwtToken := generateJWT(customer)
    
    c.JSON(200, gin.H{"token": jwtToken})
}
```

---

### 2.2 Google OAuth

**配置文件：**
```yaml
# config/oauth.yml
google:
  client_id: your_google_client_id
  client_secret: your_google_client_secret
  redirect_uri: https://remotegpu.example.com/auth/google/callback
  scopes:
    - https://www.googleapis.com/auth/userinfo.email
    - https://www.googleapis.com/auth/userinfo.profile
```

---

## 3. LDAP/AD 集成

**配置文件：**
```yaml
# config/ldap.yml
ldap:
  host: ldap.example.com
  port: 389
  use_ssl: false
  base_dn: dc=example,dc=com
  bind_dn: cn=admin,dc=example,dc=com
  bind_password: your_password
  user_filter: (uid=%s)
  attributes:
    username: uid
    email: mail
    display_name: cn
```

---

## 4. 钉钉机器人集成

**配置步骤：**
1. 在钉钉群中添加自定义机器人
2. 获取 Webhook URL 和加签密钥
3. 配置到系统中

**配置文件：**
```yaml
# config/dingtalk.yml
dingtalk:
  webhook_url: https://oapi.dingtalk.com/robot/send?access_token=xxx
  secret: your_secret_key
```

**发送消息示例：**
```go
func SendDingTalkMessage(content string) error {
    timestamp := time.Now().UnixMilli()
    sign := generateDingTalkSign(timestamp, secret)
    
    payload := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": content,
        },
    }
    
    url := fmt.Sprintf("%s&timestamp=%d&sign=%s", 
        webhookURL, timestamp, sign)
    
    return sendHTTPRequest(url, payload)
}
```

---

## 5. 企业微信机器人集成

**配置文件：**
```yaml
# config/wechat.yml
wechat:
  webhook_url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
```

**发送消息示例：**
```go
func SendWeChatMessage(content string) error {
    payload := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": content,
        },
    }
    
    return sendHTTPRequest(webhookURL, payload)
}
```

---

## 6. Slack 集成

**配置文件：**
```yaml
# config/slack.yml
slack:
  webhook_url: https://hooks.slack.com/services/xxx/yyy/zzz
```

**发送消息示例：**
```go
func SendSlackMessage(content string) error {
    payload := map[string]interface{}{
        "text": content,
    }
    
    return sendHTTPRequest(webhookURL, payload)
}
```

---

## 7. 对象存储集成

### 7.1 AWS S3

**配置文件：**
```yaml
# config/storage.yml
s3:
  endpoint: s3.amazonaws.com
  region: us-east-1
  access_key: your_access_key
  secret_key: your_secret_key
  bucket: remotegpu-data
```

### 7.2 阿里云 OSS

**配置文件：**
```yaml
oss:
  endpoint: oss-cn-hangzhou.aliyuncs.com
  access_key: your_access_key
  secret_key: your_secret_key
  bucket: remotegpu-data
```

---

## 8. Sentry 错误追踪集成

**配置文件：**
```yaml
# config/sentry.yml
sentry:
  dsn: https://xxx@sentry.io/yyy
  environment: production
  traces_sample_rate: 0.1
```

**初始化代码：**
```go
import "github.com/getsentry/sentry-go"

func InitSentry() {
    sentry.Init(sentry.ClientOptions{
        Dsn: config.Sentry.DSN,
        Environment: config.Sentry.Environment,
    })
}
```

---

## 9. GitHub Actions 集成

**工作流示例：**
```yaml
# .github/workflows/deploy.yml
name: Deploy to RemoteGPU

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to RemoteGPU
        env:
          REMOTEGPU_API_KEY: ${{ secrets.REMOTEGPU_API_KEY }}
        run: |
          curl -X POST https://api.remotegpu.com/deployments \
            -H "Authorization: Bearer $REMOTEGPU_API_KEY" \
            -d '{"image": "myapp:latest"}'
```

---

## 10. GitLab CI 集成

**配置示例：**
```yaml
# .gitlab-ci.yml
deploy:
  stage: deploy
  script:
    - |
      curl -X POST https://api.remotegpu.com/deployments \
        -H "Authorization: Bearer $REMOTEGPU_API_KEY" \
        -d '{"image": "myapp:latest"}'
  only:
    - main
```

---

## 11. 集成清单

| 系统 | 类型 | 用途 | 优先级 |
|------|------|------|--------|
| GitHub OAuth | 认证 | 用户登录 | P0 |
| Google OAuth | 认证 | 用户登录 | P1 |
| LDAP/AD | 认证 | 企业用户 | P1 |
| 钉钉机器人 | 通知 | 消息推送 | P1 |
| 企业微信 | 通知 | 消息推送 | P1 |
| Slack | 通知 | 消息推送 | P2 |
| AWS S3 | 存储 | 对象存储 | P1 |
| 阿里云 OSS | 存储 | 对象存储 | P1 |
| Sentry | 监控 | 错误追踪 | P2 |
| GitHub Actions | CI/CD | 自动部署 | P2 |
| GitLab CI | CI/CD | 自动部署 | P2 |

---

**文档版本：** v1.0
**创建日期：** 2026-01-26
