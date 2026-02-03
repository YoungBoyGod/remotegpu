# 模型管理 API

> 模型上传、版本管理相关接口

---

## 1. 创建模型

### 接口信息

- **接口路径**: `/api/v1/models`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "ResNet-50",
  "description": "图像分类模型",
  "framework": "pytorch",
  "workspace_id": 456
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "模型创建成功",
  "data": {
    "model_id": 456,
    "uuid": "model-xyz789",
    "name": "ResNet-50",
    "framework": "pytorch",
    "status": "uploading",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取模型列表

### 接口信息

- **接口路径**: `/api/v1/models`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| framework | string | 否 | 框架过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "model_id": 456,
        "name": "ResNet-50",
        "framework": "pytorch",
        "total_size": 102400000,
        "status": "ready",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

## 3. 获取模型详情

### 接口信息

- **接口路径**: `/api/v1/models/{model_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model_id": 456,
    "uuid": "model-xyz789",
    "name": "ResNet-50",
    "description": "图像分类模型",
    "framework": "pytorch",
    "storage_path": "models/customer-123/resnet50",
    "total_size": 102400000,
    "status": "ready",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 上传模型文件

### 接口信息

- **接口路径**: `/api/v1/models/{model_id}/upload`
- **请求方法**: `POST`
- **是否需要认证**: 是
- **Content-Type**: `multipart/form-data`

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "文件上传成功"
}
```

---

## 5. 删除模型

### 接口信息

- **接口路径**: `/api/v1/models/{model_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "模型删除成功"
}
```
