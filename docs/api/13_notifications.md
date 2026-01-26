# 通知管理 API

> 通知查询、标记已读相关接口

---

## 1. 获取通知列表

### 接口信息

- **接口路径**: `/api/v1/notifications`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| is_read | bool | 否 | 是否已读 |
| type | string | 否 | 类型过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "notification_id": 123,
        "title": "环境创建成功",
        "content": "您的环境 PyTorch Training 已创建成功",
        "type": "environment",
        "level": "info",
        "is_read": false,
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 45,
      "total_pages": 3
    },
    "unread_count": 12
  }
}
```

---

## 2. 获取通知详情

### 接口信息

- **接口路径**: `/api/v1/notifications/{notification_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "notification_id": 123,
    "title": "环境创建成功",
    "content": "您的环境 PyTorch Training 已创建成功，可以通过 SSH 访问",
    "type": "environment",
    "level": "info",
    "is_read": false,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 3. 标记通知为已读

### 接口信息

- **接口路径**: `/api/v1/notifications/{notification_id}/read`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "标记成功"
}
```

---

## 4. 批量标记为已读

### 接口信息

- **接口路径**: `/api/v1/notifications/read-all`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "全部标记为已读"
}
```

---

## 5. 删除通知

### 接口信息

- **接口路径**: `/api/v1/notifications/{notification_id}`
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

## 6. 获取未读通知数量

### 接口信息

- **接口路径**: `/api/v1/notifications/unread-count`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "unread_count": 12
  }
}
```
