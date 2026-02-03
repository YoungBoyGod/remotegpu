# 环境管理 API

> 开发环境创建、管理、访问相关接口

---

## 1. 创建环境

### 接口信息

- **接口路径**: `/api/v1/environments`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "PyTorch Training",
  "description": "深度学习训练环境",
  "workspace_id": 456,
  "image": "ubuntu20-pytorch:2.0",
  "cpu": 4,
  "memory": 17179869184,
  "gpu": 1,
  "storage": 107374182400
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 环境名称 |
| description | string | 否 | 描述 |
| workspace_id | int | 否 | 工作空间 ID |
| image | string | 是 | 镜像名称 |
| cpu | int | 是 | CPU 核心数 |
| memory | int | 是 | 内存（字节） |
| gpu | int | 否 | GPU 数量 |
| storage | int | 否 | 存储空间（字节） |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "环境创建成功",
  "data": {
    "env_id": "env-abc123",
    "name": "PyTorch Training",
    "status": "creating",
    "image": "ubuntu20-pytorch:2.0",
    "cpu": 4,
    "memory": 17179869184,
    "gpu": 1,
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取环境列表

### 接口信息

- **接口路径**: `/api/v1/environments`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态过滤 |
| workspace_id | int | 否 | 工作空间过滤 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "env_id": "env-abc123",
        "name": "PyTorch Training",
        "status": "running",
        "image": "ubuntu20-pytorch:2.0",
        "cpu": 4,
        "memory": 17179869184,
        "gpu": 1,
        "ssh_port": 30001,
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

## 3. 获取环境详情

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env_id": "env-abc123",
    "name": "PyTorch Training",
    "description": "深度学习训练环境",
    "status": "running",
    "image": "ubuntu20-pytorch:2.0",
    "cpu": 4,
    "memory": 17179869184,
    "gpu": 1,
    "storage": 107374182400,
    "host_id": "host-xyz789",
    "ssh_port": 30001,
    "jupyter_port": 50001,
    "access_info": {
      "ssh": "ssh root@1.2.3.4 -p 30001",
      "jupyter": "https://jupyter.remotegpu.com:50001"
    },
    "created_at": "2026-01-26T10:00:00Z",
    "started_at": "2026-01-26T10:02:00Z"
  }
}
```

---

## 4. 启动环境

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}/start`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "环境启动中"
}
```

---

## 5. 停止环境

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}/stop`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "环境停止中"
}
```

---

## 6. 重启环境

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}/restart`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "环境重启中"
}
```

---

## 7. 删除环境

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "环境删除成功"
}
```

---

## 8. 挂载数据集

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}/datasets`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "dataset_id": 789,
  "mount_path": "/data/imagenet",
  "readonly": true
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "数据集挂载成功"
}
```

---

## 9. 获取环境日志

### 接口信息

- **接口路径**: `/api/v1/environments/{env_id}/logs`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| lines | int | 否 | 日志行数（默认 100） |
| follow | bool | 否 | 是否实时跟踪 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "logs": [
      "[2026-01-26 10:00:00] Container started",
      "[2026-01-26 10:00:01] SSH service ready"
    ]
  }
}
```
