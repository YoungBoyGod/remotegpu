# 推理服务 API

> 推理服务部署、管理相关接口

---

## 1. 创建推理服务

### 接口信息

- **接口路径**: `/api/v1/inference/services`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "ResNet Inference",
  "description": "图像分类推理服务",
  "workspace_id": 456,
  "model_id": 789,
  "model_version": "v1.0",
  "framework": "pytorch",
  "cpu": 4,
  "memory": 17179869184,
  "gpu": 1,
  "replicas": 2
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "推理服务创建成功",
  "data": {
    "service_id": "svc-abc123",
    "name": "ResNet Inference",
    "status": "creating",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取推理服务列表

### 接口信息

- **接口路径**: `/api/v1/inference/services`
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
        "service_id": "svc-abc123",
        "name": "ResNet Inference",
        "status": "running",
        "endpoint_url": "https://api.remotegpu.com/inference/svc-abc123",
        "replicas": 2,
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

## 3. 获取推理服务详情

### 接口信息

- **接口路径**: `/api/v1/inference/services/{service_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "service_id": "svc-abc123",
    "name": "ResNet Inference",
    "description": "图像分类推理服务",
    "status": "running",
    "model_id": 789,
    "model_version": "v1.0",
    "framework": "pytorch",
    "endpoint_url": "https://api.remotegpu.com/inference/svc-abc123",
    "cpu": 4,
    "memory": 17179869184,
    "gpu": 1,
    "replicas": 2,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 调用推理服务

### 接口信息

- **接口路径**: `/api/v1/inference/services/{service_id}/predict`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "inputs": {
    "image": "base64_encoded_image_data"
  }
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "predictions": [
      {
        "class": "cat",
        "confidence": 0.95
      },
      {
        "class": "dog",
        "confidence": 0.03
      }
    ]
  }
}
```

---

## 5. 停止推理服务

### 接口信息

- **接口路径**: `/api/v1/inference/services/{service_id}/stop`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "推理服务停止中"
}
```

---

## 6. 删除推理服务

### 接口信息

- **接口路径**: `/api/v1/inference/services/{service_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "推理服务删除成功"
}
```
