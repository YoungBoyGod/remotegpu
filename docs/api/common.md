# API 公共规范

> 本文档定义 RemoteGPU API 的公共规范，包括认证、分页、排序、错误码等

---

## 1. 认证方式

### 1.1 JWT Token 认证

所有需要认证的接口都需要在请求头中携带 JWT Token：

```http
Authorization: Bearer {access_token}
```

### 1.2 Token 获取

通过登录接口获取：

```http
POST /api/v1/auth/login
```

### 1.3 Token 刷新

Access Token 过期后，使用 Refresh Token 刷新：

```http
POST /api/v1/auth/refresh
```

---

## 2. 请求格式

### 2.1 Content-Type

所有 POST/PUT 请求必须使用 JSON 格式：

```http
Content-Type: application/json
```

### 2.2 请求示例

```http
POST /api/v1/environments
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "name": "PyTorch Training",
  "image": "ubuntu20-pytorch:2.0",
  "cpu": 4,
  "memory": 16384,
  "gpu": 1
}
```

---

## 3. 响应格式

### 3.1 成功响应

**格式：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

**示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "env-abc123",
    "name": "PyTorch Training",
    "status": "running"
  }
}
```

### 3.2 错误响应

**格式：**

```json
{
  "code": 40001,
  "message": "错误描述",
  "error": "详细错误信息"
}
```

**示例：**

```json
{
  "code": 40001,
  "message": "参数错误",
  "error": "cpu 字段必须是正整数"
}
```

---

## 4. 分页

### 4.1 分页参数

所有列表接口支持分页，使用以下查询参数：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码（从 1 开始） |
| page_size | int | 否 | 20 | 每页数量（最大 100） |

**请求示例：**

```http
GET /api/v1/environments?page=2&page_size=10
```

### 4.2 分页响应

**格式：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [],
    "pagination": {
      "page": 2,
      "page_size": 10,
      "total": 45,
      "total_pages": 5
    }
  }
}
```

**字段说明：**

- `items`: 当前页数据列表
- `page`: 当前页码
- `page_size`: 每页数量
- `total`: 总记录数
- `total_pages`: 总页数

---

## 5. 排序

### 5.1 排序参数

列表接口支持排序，使用以下查询参数：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| sort_by | string | 否 | 排序字段 |
| order | string | 否 | 排序方向：asc（升序）、desc（降序） |

**请求示例：**

```http
GET /api/v1/environments?sort_by=created_at&order=desc
```

### 5.2 支持的排序字段

不同接口支持的排序字段不同，详见各接口文档。

常见排序字段：

- `created_at`: 创建时间
- `updated_at`: 更新时间
- `name`: 名称
- `status`: 状态

---

## 6. 过滤

### 6.1 过滤参数

列表接口支持过滤，使用查询参数：

**请求示例：**

```http
GET /api/v1/environments?status=running&workspace_id=123
```

### 6.2 常见过滤字段

- `status`: 状态过滤
- `workspace_id`: 工作空间 ID
- `customer_id`: 客户 ID
- `created_after`: 创建时间起始
- `created_before`: 创建时间结束

---

## 7. 错误码

### 7.1 HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 204 | 删除成功（无内容） |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

### 7.2 业务错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 系统错误 |
| 10002 | 数据库错误 |
| 10003 | 缓存错误 |
| 20001 | 认证失败 |
| 20002 | Token 无效 |
| 20003 | Token 过期 |
| 20004 | 权限不足 |
| 30001 | 用户不存在 |
| 30002 | 用户已存在 |
| 30003 | 密码错误 |
| 30004 | 邮箱未验证 |
| 40001 | 参数错误 |
| 40002 | 参数缺失 |
| 40003 | 参数格式错误 |
| 50001 | 资源不存在 |
| 50002 | 资源已存在 |
| 50003 | 资源状态错误 |
| 60001 | 配额不足 |
| 60002 | 余额不足 |
| 60003 | 资源不足 |
| 70001 | 操作失败 |
| 70002 | 操作超时 |
| 70003 | 操作冲突 |

---

## 8. 时间格式

### 8.1 ISO 8601 格式

所有时间字段使用 ISO 8601 格式：

```
2026-01-26T10:00:00Z
```

### 8.2 时区

- 服务器时区：UTC
- 客户端需要根据本地时区转换显示

---

## 9. 字段命名规范

### 9.1 命名风格

- 使用 **snake_case**（下划线命名）
- 示例：`created_at`、`workspace_id`、`cpu_usage_percent`

### 9.2 常见字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | string/int | 资源 ID |
| uuid | string | 全局唯一标识符 |
| name | string | 名称 |
| description | string | 描述 |
| status | string | 状态 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |
| deleted_at | string | 删除时间（软删除） |

---

## 10. 批量操作

### 10.1 批量删除

```http
DELETE /api/v1/environments/batch
Content-Type: application/json

{
  "ids": ["env-1", "env-2", "env-3"]
}
```

### 10.2 批量更新

```http
PUT /api/v1/environments/batch
Content-Type: application/json

{
  "ids": ["env-1", "env-2"],
  "status": "stopped"
}
```

---

## 11. 文件上传

### 11.1 单文件上传

```http
POST /api/v1/datasets/upload
Content-Type: multipart/form-data

file: (binary)
name: "dataset-name"
description: "dataset description"
```

### 11.2 分片上传

大文件支持分片上传：

**1. 初始化上传**

```http
POST /api/v1/datasets/upload/init
{
  "filename": "large-dataset.zip",
  "size": 10737418240,
  "chunk_size": 10485760
}
```

**2. 上传分片**

```http
POST /api/v1/datasets/upload/chunk
Content-Type: multipart/form-data

upload_id: "upload-abc123"
chunk_index: 0
chunk: (binary)
```

**3. 完成上传**

```http
POST /api/v1/datasets/upload/complete
{
  "upload_id": "upload-abc123"
}
```

---

## 12. 限流

### 12.1 限流规则

- **普通用户**：100 请求/分钟
- **企业用户**：1000 请求/分钟
- **管理员**：无限制

### 12.2 限流响应

超过限流时返回 429 状态码：

```json
{
  "code": 42901,
  "message": "请求过于频繁，请稍后再试",
  "retry_after": 60
}
```

响应头包含限流信息：

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1706256000
```

---

## 13. 版本控制

### 13.1 API 版本

当前版本：`v1`

基础 URL：`https://api.remotegpu.com/v1`

### 13.2 版本升级

新版本发布时，旧版本会保留至少 6 个月。

---

## 14. CORS

### 14.1 跨域支持

API 支持 CORS，允许的域名：

- `https://app.remotegpu.com`
- `http://localhost:3000`（开发环境）

### 14.2 预检请求

浏览器会发送 OPTIONS 预检请求：

```http
OPTIONS /api/v1/environments
Access-Control-Request-Method: POST
Access-Control-Request-Headers: authorization, content-type
```

---

**文档版本**：v1.0
**最后更新**：2026-01-26
