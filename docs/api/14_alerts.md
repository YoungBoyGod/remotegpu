# 告警管理 API

> 告警规则、告警记录相关接口

---

## 1. 创建告警规则

### 接口信息

- **接口路径**: `/api/v1/alerts/rules`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "GPU 温度告警",
  "description": "GPU 温度超过 80℃ 时告警",
  "metric_type": "gpu_temperature",
  "threshold": 80.0,
  "comparison": ">",
  "duration": 300,
  "severity": "warning",
  "enabled": true
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "告警规则创建成功",
  "data": {
    "rule_id": 123,
    "name": "GPU 温度告警",
    "enabled": true,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取告警规则列表

### 接口信息

- **接口路径**: `/api/v1/alerts/rules`
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
        "rule_id": 123,
        "name": "GPU 温度告警",
        "metric_type": "gpu_temperature",
        "threshold": 80.0,
        "severity": "warning",
        "enabled": true,
        "created_at": "2026-01-26T10:00:00Z"
      }
    ]
  }
}
```

---

## 3. 更新告警规则

### 接口信息

- **接口路径**: `/api/v1/alerts/rules/{rule_id}`
- **请求方法**: `PUT`
- **是否需要认证**: 是

### 请求参数

```json
{
  "threshold": 85.0,
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

## 4. 删除告警规则

### 接口信息

- **接口路径**: `/api/v1/alerts/rules/{rule_id}`
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

## 5. 获取告警记录

### 接口信息

- **接口路径**: `/api/v1/alerts/records`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态过滤 |
| severity | string | 否 | 严重程度过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "record_id": 456,
        "rule_id": 123,
        "metric_type": "gpu_temperature",
        "current_value": 85.5,
        "threshold": 80.0,
        "severity": "warning",
        "status": "firing",
        "message": "GPU 温度超过阈值",
        "triggered_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 25,
      "total_pages": 2
    }
  }
}
```
