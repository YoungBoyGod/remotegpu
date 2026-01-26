# 数据集管理 API

> 数据集上传、版本管理相关接口

---

## 1. 创建数据集

### 接口信息

- **接口路径**: `/api/v1/datasets`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "ImageNet-1K",
  "description": "图像分类数据集",
  "workspace_id": 456,
  "visibility": "private"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 数据集名称 |
| description | string | 否 | 描述 |
| workspace_id | int | 否 | 工作空间 ID |
| visibility | string | 否 | 可见性：private, workspace, public |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "数据集创建成功",
  "data": {
    "dataset_id": 789,
    "uuid": "dataset-abc123",
    "name": "ImageNet-1K",
    "status": "uploading",
    "upload_url": "https://upload.remotegpu.com/datasets/789",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取数据集列表

### 接口信息

- **接口路径**: `/api/v1/datasets`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| workspace_id | int | 否 | 工作空间过滤 |
| visibility | string | 否 | 可见性过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "dataset_id": 789,
        "name": "ImageNet-1K",
        "description": "图像分类数据集",
        "total_size": 1099511627776,
        "file_count": 1281167,
        "status": "ready",
        "visibility": "private",
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

## 3. 获取数据集详情

### 接口信息

- **接口路径**: `/api/v1/datasets/{dataset_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "dataset_id": 789,
    "uuid": "dataset-abc123",
    "name": "ImageNet-1K",
    "description": "图像分类数据集",
    "storage_path": "datasets/customer-123/imagenet",
    "storage_type": "minio",
    "total_size": 1099511627776,
    "file_count": 1281167,
    "status": "ready",
    "visibility": "private",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 上传数据集文件

### 接口信息

- **接口路径**: `/api/v1/datasets/{dataset_id}/upload`
- **请求方法**: `POST`
- **是否需要认证**: 是
- **Content-Type**: `multipart/form-data`

### 请求参数

```
file: (binary)
path: "train/images/001.jpg"
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "文件上传成功",
  "data": {
    "filename": "001.jpg",
    "size": 102400,
    "path": "train/images/001.jpg"
  }
}
```

---

## 5. 删除数据集

### 接口信息

- **接口路径**: `/api/v1/datasets/{dataset_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "数据集删除成功"
}
```

---

## 6. 创建数据集版本

### 接口信息

- **接口路径**: `/api/v1/datasets/{dataset_id}/versions`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "version": "v1.1",
  "description": "添加了新的图片"
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "版本创建成功",
  "data": {
    "version_id": 123,
    "version": "v1.1",
    "dataset_id": 789,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 7. 获取数据集版本列表

### 接口信息

- **接口路径**: `/api/v1/datasets/{dataset_id}/versions`
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
        "version_id": 123,
        "version": "v1.1",
        "description": "添加了新的图片",
        "size": 1099511627776,
        "file_count": 1281167,
        "created_at": "2026-01-26T10:00:00Z"
      }
    ]
  }
}
```
