# Webhook 管理

> 所属模块：模块 15 - Webhook 管理模块
>
> 功能编号：15.1
>
> 优先级：P1（重要）

---

## 1. 功能概述

### 1.1 功能描述

Webhook 管理功能允许用户配置自定义的 HTTP 回调接口，当系统中发生特定事件时（如任务完成、环境创建、告警触发），自动向用户指定的 URL 发送 HTTP 请求，实现与第三方系统的集成。

### 1.2 业务价值

- ✅ 与第三方系统集成
- ✅ 事件驱动的自动化
- ✅ 灵活的事件订阅
- ✅ 请求重试和失败处理

---

## 2. 支持的事件类型

### 2.1 事件分类

| 事件类型 | 说明 | 触发时机 |
|---------|------|---------|
| environment.created | 环境创建 | 环境创建成功后 |
| environment.started | 环境启动 | 环境启动成功后 |
| environment.stopped | 环境停止 | 环境停止后 |
| job.completed | 任务完成 | 训练任务完成后 |
| job.failed | 任务失败 | 训练任务失败后 |
| alert.triggered | 告警触发 | 监控告警触发时 |
| quota.exceeded | 配额超限 | 资源配额超限时 |

---

## 3. 数据模型

```sql
-- Webhook 配置表
CREATE TABLE webhooks (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    url VARCHAR(512) NOT NULL,
    secret VARCHAR(128),
    events TEXT[] NOT NULL,
    status VARCHAR(32) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- Webhook 日志表
CREATE TABLE webhook_logs (
    id BIGSERIAL PRIMARY KEY,
    webhook_id BIGINT NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    payload JSONB NOT NULL,
    response_status INT,
    response_body TEXT,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (webhook_id) REFERENCES webhooks(id)
);
```

---

## 4. Webhook 实现

### 4.1 Webhook 触发器

```go
// Webhook 触发器
type WebhookTrigger struct {
    db         *gorm.DB
    httpClient *http.Client
}

// 触发 Webhook
func (t *WebhookTrigger) Trigger(event string, payload interface{}) error {
    // 查询订阅该事件的 Webhook
    var webhooks []Webhook
    t.db.Where("? = ANY(events) AND status = 'active'", event).Find(&webhooks)

    for _, webhook := range webhooks {
        go t.sendWebhook(webhook, event, payload)
    }

    return nil
}

// 发送 Webhook 请求
func (t *WebhookTrigger) sendWebhook(webhook Webhook, event string, payload interface{}) {
    data, _ := json.Marshal(map[string]interface{}{
        "event":     event,
        "timestamp": time.Now().Unix(),
        "data":      payload,
    })

    req, _ := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(data))
    req.Header.Set("Content-Type", "application/json")
    
    // 添加签名
    if webhook.Secret != "" {
        signature := generateSignature(data, webhook.Secret)
        req.Header.Set("X-Webhook-Signature", signature)
    }

    resp, err := t.httpClient.Do(req)
    
    // 记录日志
    log := &WebhookLog{
        WebhookID: webhook.ID,
        EventType: event,
        Payload:   data,
    }
    
    if err == nil {
        log.ResponseStatus = resp.StatusCode
    }
    
    t.db.Create(log)
}
```

---

## 5. API 接口

### 5.1 创建 Webhook

```go
POST /api/webhooks
Body: {
  "name": "任务完成通知",
  "url": "https://example.com/webhook",
  "secret": "your_secret_key",
  "events": ["job.completed", "job.failed"]
}

Response: {
  "webhook_id": 1,
  "status": "active"
}
```

### 5.2 Webhook 请求格式

```json
{
  "event": "job.completed",
  "timestamp": 1706256000,
  "data": {
    "job_id": "job-abc123",
    "status": "completed",
    "duration": 3600
  }
}
```

---

## 6. 测试用例

| 用例 | 场景 | 预期结果 |
|------|------|---------|
| TC-01 | 创建 Webhook | 创建成功 |
| TC-02 | 触发事件 | Webhook 请求发送成功 |
| TC-03 | 验证签名 | 签名验证通过 |

---

**文档版本：** v1.0
**创建日期：** 2026-01-26
