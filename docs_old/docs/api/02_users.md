# 用户管理 API

> 用户信息、配额管理相关接口

---

## 1. 获取当前用户信息

### 接口信息

- **接口路径**: `/api/v1/users/me`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "customer_id": 123,
    "uuid": "cust-abc123",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "company": "Example Corp",
    "user_type": "external",
    "account_type": "individual",
    "status": "active",
    "email_verified": true,
    "phone_verified": false,
    "avatar_url": "https://cdn.remotegpu.com/avatars/123.jpg",
    "created_at": "2026-01-20T10:00:00Z",
    "last_login_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 更新用户信息

### 接口信息

- **接口路径**: `/api/v1/users/me`
- **请求方法**: `PUT`
- **是否需要认证**: 是

### 请求参数

```json
{
  "full_name": "John Smith",
  "company": "New Company",
  "phone": "+86 138****1234",
  "avatar_url": "https://cdn.remotegpu.com/avatars/new.jpg"
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "customer_id": 123,
    "full_name": "John Smith",
    "company": "New Company"
  }
}
```

---

## 3. 获取用户配额

### 接口信息

- **接口路径**: `/api/v1/users/me/quota`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "quota_level": "basic",
    "cpu_quota": 8,
    "cpu_used": 4,
    "memory_quota": 16384,
    "memory_used": 8192,
    "gpu_quota": 1,
    "gpu_used": 1,
    "storage_quota": 53687091200,
    "storage_used": 10737418240,
    "environment_quota": 3,
    "environment_used": 2
  }
}
```

---

## 4. 获取用户列表（管理员）

### 接口信息

- **接口路径**: `/api/v1/users`
- **请求方法**: `GET`
- **是否需要认证**: 是（需要管理员权限）

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| user_type | string | 否 | 用户类型过滤 |
| status | string | 否 | 状态过滤 |
| keyword | string | 否 | 关键词搜索 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "customer_id": 123,
        "username": "john_doe",
        "email": "john@example.com",
        "user_type": "external",
        "status": "active",
        "created_at": "2026-01-20T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 45,
      "total_pages": 3
    }
  }
}
```

---

## 5. 获取指定用户信息（管理员）

### 接口信息

- **接口路径**: `/api/v1/users/{customer_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是（需要管理员权限）

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| customer_id | int | 用户 ID |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "customer_id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "user_type": "external",
    "account_type": "individual",
    "status": "active",
    "created_at": "2026-01-20T10:00:00Z"
  }
}
```

---

## 6. 更新用户配额（管理员）

### 接口信息

- **接口路径**: `/api/v1/users/{customer_id}/quota`
- **请求方法**: `PUT`
- **是否需要认证**: 是（需要管理员权限）

### 请求参数

```json
{
  "quota_level": "pro",
  "cpu_quota": 16,
  "memory_quota": 32768,
  "gpu_quota": 2,
  "storage_quota": 214748364800,
  "environment_quota": 10
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "配额更新成功"
}
```

---

## 7. 更新用户状态（管理员）

### 接口信息

- **接口路径**: `/api/v1/users/{customer_id}/status`
- **请求方法**: `PUT`
- **是否需要认证**: 是（需要管理员权限）

### 请求参数

```json
{
  "status": "suspended"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 是 | 状态：active, suspended, deleted |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "状态更新成功"
}
```
