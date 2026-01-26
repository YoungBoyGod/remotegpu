# 主机管理 API

> 主机注册、监控、GPU 管理相关接口

---

## 1. 注册主机

### 接口信息

- **接口路径**: `/api/v1/hosts`
- **请求方法**: `POST`
- **是否需要认证**: 是（需要管理员权限）

### 请求参数

```json
{
  "name": "GPU-Server-01",
  "ip_address": "192.168.1.10",
  "os_type": "linux",
  "os_version": "Ubuntu 20.04",
  "deployment_mode": "traditional",
  "total_cpu": 32,
  "total_memory": 137438953472,
  "total_disk": 2199023255552
}
```

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "主机注册成功",
  "data": {
    "host_id": "host-abc123",
    "name": "GPU-Server-01",
    "status": "offline",
    "agent_token": "agent-token-xyz789"
  }
}
```

---

## 2. 获取主机列表

### 接口信息

- **接口路径**: `/api/v1/hosts`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态过滤 |
| os_type | string | 否 | 操作系统类型 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "host_id": "host-abc123",
        "name": "GPU-Server-01",
        "ip_address": "192.168.1.10",
        "os_type": "linux",
        "status": "online",
        "total_cpu": 32,
        "used_cpu": 8,
        "total_memory": 137438953472,
        "used_memory": 34359738368,
        "total_gpu": 4,
        "used_gpu": 2,
        "last_heartbeat": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 10,
      "total_pages": 1
    }
  }
}
```

---

## 3. 获取主机详情

### 接口信息

- **接口路径**: `/api/v1/hosts/{host_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "host_id": "host-abc123",
    "name": "GPU-Server-01",
    "hostname": "gpu01.example.com",
    "ip_address": "192.168.1.10",
    "os_type": "linux",
    "os_version": "Ubuntu 20.04",
    "deployment_mode": "traditional",
    "status": "online",
    "health_status": "healthy",
    "total_cpu": 32,
    "used_cpu": 8,
    "total_memory": 137438953472,
    "used_memory": 34359738368,
    "total_gpu": 4,
    "used_gpu": 2,
    "labels": {
      "region": "us-west",
      "zone": "a"
    },
    "last_heartbeat": "2026-01-26T10:00:00Z",
    "created_at": "2026-01-20T10:00:00Z"
  }
}
```

---

## 4. 更新主机信息

### 接口信息

- **接口路径**: `/api/v1/hosts/{host_id}`
- **请求方法**: `PUT`
- **是否需要认证**: 是（需要管理员权限）

### 请求参数

```json
{
  "name": "GPU-Server-01-Updated",
  "labels": {
    "region": "us-west",
    "zone": "b"
  }
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

## 5. 删除主机

### 接口信息

- **接口路径**: `/api/v1/hosts/{host_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是（需要管理员权限）

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

---

## 6. 获取主机 GPU 列表

### 接口信息

- **接口路径**: `/api/v1/hosts/{host_id}/gpus`
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
        "gpu_id": 1,
        "gpu_index": 0,
        "name": "Tesla V100-SXM2-32GB",
        "memory_total": 34359738368,
        "status": "allocated",
        "allocated_to": "env-xyz789",
        "temperature": 68.5,
        "utilization": 85.5
      }
    ]
  }
}
```

---

## 7. 主机心跳上报

### 接口信息

- **接口路径**: `/api/v1/hosts/{host_id}/heartbeat`
- **请求方法**: `POST`
- **是否需要认证**: 是（Agent Token）

### 请求参数

```json
{
  "status": "online",
  "health_status": "healthy",
  "used_cpu": 8,
  "used_memory": 34359738368,
  "used_gpu": 2
}
```

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "心跳上报成功"
}
```
