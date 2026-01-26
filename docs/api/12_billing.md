# 计费管理 API

> 计费记录、账单查询相关接口

---

## 1. 获取计费记录

### 接口信息

- **接口路径**: `/api/v1/billing/records`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| start_time | string | 否 | 开始时间 |
| end_time | string | 否 | 结束时间 |
| resource_type | string | 否 | 资源类型过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "record_id": 123,
        "env_id": "env-abc123",
        "resource_type": "gpu",
        "quantity": 1.0,
        "unit_price": 2.5000,
        "amount": 2.5000,
        "currency": "CNY",
        "start_time": "2026-01-26T10:00:00Z",
        "end_time": "2026-01-26T11:00:00Z",
        "created_at": "2026-01-26T11:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

---

## 2. 获取账单列表

### 接口信息

- **接口路径**: `/api/v1/billing/invoices`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "invoice_id": 456,
        "invoice_no": "INV-2026-01-001",
        "billing_period_start": "2026-01-01T00:00:00Z",
        "billing_period_end": "2026-01-31T23:59:59Z",
        "total_amount": 1250.5000,
        "currency": "CNY",
        "status": "paid",
        "paid_at": "2026-02-01T10:00:00Z",
        "created_at": "2026-02-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 12,
      "total_pages": 1
    }
  }
}
```

---

## 3. 获取账单详情

### 接口信息

- **接口路径**: `/api/v1/billing/invoices/{invoice_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "invoice_id": 456,
    "invoice_no": "INV-2026-01-001",
    "billing_period_start": "2026-01-01T00:00:00Z",
    "billing_period_end": "2026-01-31T23:59:59Z",
    "total_amount": 1250.5000,
    "currency": "CNY",
    "status": "paid",
    "paid_at": "2026-02-01T10:00:00Z",
    "items": [
      {
        "resource_type": "gpu",
        "quantity": 500.0,
        "unit_price": 2.5000,
        "amount": 1250.0000
      }
    ],
    "created_at": "2026-02-01T00:00:00Z"
  }
}
```

---

## 4. 获取费用统计

### 接口信息

- **接口路径**: `/api/v1/billing/statistics`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | string | 是 | 开始时间 |
| end_time | string | 是 | 结束时间 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_amount": 1250.5000,
    "by_resource_type": {
      "cpu": 150.0000,
      "memory": 200.0000,
      "gpu": 900.5000
    },
    "currency": "CNY"
  }
}
```
