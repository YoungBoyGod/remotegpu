# 镜像管理 API

> 镜像列表、详情查询相关接口

---

## 1. 获取镜像列表

### 接口信息

- **接口路径**: `/api/v1/images`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| category | string | 否 | 分类过滤 |
| framework | string | 否 | 框架过滤 |
| is_official | bool | 否 | 是否官方镜像 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "image_id": 1,
        "name": "ubuntu20-pytorch:2.0",
        "display_name": "Ubuntu 20.04 + PyTorch 2.0",
        "description": "深度学习训练环境",
        "category": "pytorch",
        "framework": "pytorch",
        "framework_version": "2.0.0",
        "cuda_version": "11.8",
        "python_version": "3.10",
        "is_official": true,
        "size": 5368709120,
        "status": "active"
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

## 2. 获取镜像详情

### 接口信息

- **接口路径**: `/api/v1/images/{image_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "image_id": 1,
    "name": "ubuntu20-pytorch:2.0",
    "display_name": "Ubuntu 20.04 + PyTorch 2.0",
    "description": "深度学习训练环境，包含 PyTorch 2.0、CUDA 11.8、Python 3.10",
    "category": "pytorch",
    "framework": "pytorch",
    "framework_version": "2.0.0",
    "cuda_version": "11.8",
    "python_version": "3.10",
    "is_official": true,
    "size": 5368709120,
    "registry_url": "harbor.remotegpu.com/library/ubuntu20-pytorch:2.0",
    "status": "active",
    "created_at": "2026-01-20T10:00:00Z"
  }
}
```

---

## 3. 创建自定义镜像（管理员）

### 接口信息

- **接口路径**: `/api/v1/images`
- **请求方法**: `POST`
- **是否需要认证**: 是（需要管理员权限）

### 请求参数

```json
{
  "name": "custom-pytorch:1.0",
  "display_name": "自定义 PyTorch 环境",
  "description": "基于 PyTorch 2.0 的自定义环境",
  "category": "custom",
  "framework": "pytorch",
  "registry_url": "harbor.remotegpu.com/custom/pytorch:1.0"
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "镜像创建成功",
  "data": {
    "image_id": 100,
    "name": "custom-pytorch:1.0"
  }
}
```
