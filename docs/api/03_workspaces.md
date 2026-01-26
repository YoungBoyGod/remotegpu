# 工作空间 API

> 工作空间、成员管理相关接口

---

## 1. 创建工作空间

### 接口信息

- **接口路径**: `/api/v1/workspaces`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "AI Research Team",
  "description": "深度学习研究团队",
  "type": "team"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 工作空间名称 |
| description | string | 否 | 描述 |
| type | string | 是 | 类型：personal, team, enterprise |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "workspace_id": 456,
    "uuid": "ws-xyz789",
    "name": "AI Research Team",
    "type": "team",
    "owner_id": 123,
    "member_count": 1,
    "status": "active",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取工作空间列表

### 接口信息

- **接口路径**: `/api/v1/workspaces`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
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
        "workspace_id": 456,
        "name": "AI Research Team",
        "type": "team",
        "member_count": 5,
        "role": "owner",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

---

## 3. 获取工作空间详情

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "workspace_id": 456,
    "uuid": "ws-xyz789",
    "name": "AI Research Team",
    "description": "深度学习研究团队",
    "type": "team",
    "owner_id": 123,
    "member_count": 5,
    "status": "active",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 更新工作空间

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}`
- **请求方法**: `PUT`
- **是否需要认证**: 是（需要 owner 权限）

### 请求参数

```json
{
  "name": "AI Research Lab",
  "description": "Updated description"
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

## 5. 删除工作空间

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是（需要 owner 权限）

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

---

## 6. 添加成员

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}/members`
- **请求方法**: `POST`
- **是否需要认证**: 是（需要 owner/admin 权限）

### 请求参数

```json
{
  "customer_id": 789,
  "role": "member"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_id | int | 是 | 用户 ID |
| role | string | 是 | 角色：admin, member, viewer |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "成员添加成功"
}
```

---

## 7. 获取成员列表

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}/members`
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
        "customer_id": 123,
        "username": "john_doe",
        "email": "john@example.com",
        "role": "owner",
        "status": "active",
        "joined_at": "2026-01-26T10:00:00Z"
      }
    ]
  }
}
```

---

## 8. 更新成员角色

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}/members/{customer_id}`
- **请求方法**: `PUT`
- **是否需要认证**: 是（需要 owner/admin 权限）

### 请求参数

```json
{
  "role": "admin"
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "角色更新成功"
}
```

---

## 9. 移除成员

### 接口信息

- **接口路径**: `/api/v1/workspaces/{workspace_id}/members/{customer_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是（需要 owner/admin 权限）

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "成员移除成功"
}
```
