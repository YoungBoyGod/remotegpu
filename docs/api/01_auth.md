# 认证授权 API

> 用户认证、注册、Token 管理相关接口

---

## 1. 用户注册

### 接口信息

- **接口路径**: `/api/v1/auth/register`
- **请求方法**: `POST`
- **是否需要认证**: 否

### 请求参数

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe",
  "company": "Example Corp"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（3-64字符） |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码（8位以上，包含大小写字母和数字） |
| full_name | string | 否 | 全名 |
| company | string | 否 | 公司名称 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "customer_id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "user_type": "external",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**错误响应 (400)**

```json
{
  "code": 30002,
  "message": "用户已存在",
  "error": "用户名或邮箱已被注册"
}
```

---

## 2. 用户登录

### 接口信息

- **接口路径**: `/api/v1/auth/login`
- **请求方法**: `POST`
- **是否需要认证**: 否

### 请求参数

```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名或邮箱 |
| password | string | 是 | 密码 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "customer_id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "user_type": "external",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**错误响应 (401)**

```json
{
  "code": 30003,
  "message": "用户名或密码错误"
}
```

---

## 3. 刷新 Token

### 接口信息

- **接口路径**: `/api/v1/auth/refresh`
- **请求方法**: `POST`
- **是否需要认证**: 否

### 请求参数

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "Token 刷新成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

---

## 4. 用户登出

### 接口信息

- **接口路径**: `/api/v1/auth/logout`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "登出成功"
}
```

---

## 5. 修改密码

### 接口信息

- **接口路径**: `/api/v1/auth/change-password`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "old_password": "OldPass123!",
  "new_password": "NewPass456!"
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "密码修改成功"
}
```

---

## 6. 忘记密码

### 接口信息

- **接口路径**: `/api/v1/auth/forgot-password`
- **请求方法**: `POST`
- **是否需要认证**: 否

### 请求参数

```json
{
  "email": "john@example.com"
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "重置密码邮件已发送"
}
```

---

## 7. 重置密码

### 接口信息

- **接口路径**: `/api/v1/auth/reset-password`
- **请求方法**: `POST`
- **是否需要认证**: 否

### 请求参数

```json
{
  "token": "reset-token-abc123",
  "new_password": "NewPass456!"
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "密码重置成功"
}
```

---

## 8. 验证邮箱

### 接口信息

- **接口路径**: `/api/v1/auth/verify-email`
- **请求方法**: `POST`
- **是否需要认证**: 否

### 请求参数

```json
{
  "token": "verify-token-xyz789"
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "邮箱验证成功"
}
```
