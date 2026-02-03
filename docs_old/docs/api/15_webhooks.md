# Webhook 管理 API

> Webhook 配置、日志查询相关接口

---

## 1. 创建 Webhook

### 接口信息

- **接口路径**: `/api/v1/webhooks`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "环境事件通知",
  "url": "https://your-server.com/webhook",
  "secret": "your-secret-key",
  "events": ["environment.created", "environment.started", "environment.stopped"],
  "enabled": true
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | Webhook 名称 |
| url | string | 是 | 回调 URL |
| secret | string | 否 | 签名密钥 |
| events | array | 是 | 订阅的事件列表 |
| enabled | bool | 否 | 是否启用 |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "Webhook 创建成功",
  "data": {
    "webhook_id": 123,
    "name": "环境事件通知",
    "url": "https://your-server.com/webhook",
    "enabled": true,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取 Webhook 列表

### 接口信息

- **接口路径**: `/api/v1/webhooks`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "webhook_id": 123,
        "name": "环境事件通知",
        "url": "https://your-server.com/webhook",
        "events": ["environment.created", "environment.started"],
        "enabled": true,
        "created_at": "2026-01-26T10:00:00Z"
      }
    ]
  }
}
```

---

## 3. 获取 Webhook 详情

### 接口信息

- **接口路径**: `/api/v1/webhooks/{webhook_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "webhook_id": 123,
    "name": "环境事件通知",
    "url": "https://your-server.com/webhook",
    "secret": "your-secret-key",
    "events": ["environment.created", "environment.started", "environment.stopped"],
    "enabled": true,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 更新 Webhook

### 接口信息

- **接口路径**: `/api/v1/webhooks/{webhook_id}`
- **请求方法**: `PUT`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "环境事件通知（更新）",
  "enabled": false
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "更新成功"
}
```

---

## 5. 删除 Webhook

### 接口信息

- **接口路径**: `/api/v1/webhooks/{webhook_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

---

## 6. 测试 Webhook

### 接口信息

- **接口路径**: `/api/v1/webhooks/{webhook_id}/test`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "测试请求已发送",
  "data": {
    "status_code": 200,
    "response_time": 125
  }
}
```

---

## 7. 获取 Webhook 日志

### 接口信息

- **接口路径**: `/api/v1/webhooks/{webhook_id}/logs`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "log_id": 456,
        "event_type": "environment.created",
        "status_code": 200,
        "response_body": "OK",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

---

## 8. 支持的事件类型

| 事件类型 | 说明 |
|---------|------|
| environment.created | 环境创建 |
| environment.started | 环境启动 |
| environment.stopped | 环境停止 |
| environment.deleted | 环境删除 |
| training.started | 训练任务开始 |
| training.completed | 训练任务完成 |
| training.failed | 训练任务失败 |
| alert.triggered | 告警触发 |
| alert.resolved | 告警解决 |
