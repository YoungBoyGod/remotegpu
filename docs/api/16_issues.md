# 问题单和需求单 API

> 问题单、需求单、评论管理相关接口

---

## 1. 创建问题单

### 接口信息

- **接口路径**: `/api/v1/issues`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "title": "环境无法启动",
  "description": "环境 env-abc123 无法启动，一直处于 creating 状态",
  "type": "bug",
  "priority": "high"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题 |
| description | string | 是 | 描述 |
| type | string | 是 | 类型：bug, question, feature |
| priority | string | 否 | 优先级：low, medium, high, critical |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "问题单创建成功",
  "data": {
    "issue_id": 123,
    "issue_no": "ISSUE-2026-001",
    "title": "环境无法启动",
    "status": "open",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取问题单列表

### 接口信息

- **接口路径**: `/api/v1/issues`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态过滤 |
| type | string | 否 | 类型过滤 |
| priority | string | 否 | 优先级过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "issue_id": 123,
        "issue_no": "ISSUE-2026-001",
        "title": "环境无法启动",
        "type": "bug",
        "priority": "high",
        "status": "open",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 15,
      "total_pages": 1
    }
  }
}
```

---

## 3. 获取问题单详情

### 接口信息

- **接口路径**: `/api/v1/issues/{issue_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "issue_id": 123,
    "issue_no": "ISSUE-2026-001",
    "title": "环境无法启动",
    "description": "环境 env-abc123 无法启动，一直处于 creating 状态",
    "type": "bug",
    "priority": "high",
    "status": "open",
    "assignee_id": 456,
    "created_at": "2026-01-26T10:00:00Z",
    "updated_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 更新问题单

### 接口信息

- **接口路径**: `/api/v1/issues/{issue_id}`
- **请求方法**: `PUT`
- **是否需要认证**: 是

### 请求参数

```json
{
  "status": "resolved",
  "assignee_id": 456
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

## 5. 创建需求单

### 接口信息

- **接口路径**: `/api/v1/requirements`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "title": "支持 Windows 环境",
  "description": "希望能够创建 Windows 开发环境",
  "priority": "medium"
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "需求单创建成功",
  "data": {
    "requirement_id": 789,
    "requirement_no": "REQ-2026-001",
    "title": "支持 Windows 环境",
    "status": "submitted",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 6. 获取需求单列表

### 接口信息

- **接口路径**: `/api/v1/requirements`
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
        "requirement_id": 789,
        "requirement_no": "REQ-2026-001",
        "title": "支持 Windows 环境",
        "priority": "medium",
        "status": "submitted",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 8,
      "total_pages": 1
    }
  }
}
```

---

## 7. 添加评论

### 接口信息

- **接口路径**: `/api/v1/issues/{issue_id}/comments`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "content": "我也遇到了同样的问题"
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "评论添加成功",
  "data": {
    "comment_id": 456,
    "content": "我也遇到了同样的问题",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 8. 获取评论列表

### 接口信息

- **接口路径**: `/api/v1/issues/{issue_id}/comments`
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
        "comment_id": 456,
        "customer_id": 123,
        "username": "john_doe",
        "content": "我也遇到了同样的问题",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ]
  }
}
```
