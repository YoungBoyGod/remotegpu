# 监控 API

> 监控数据查询相关接口

---

## 1. 获取主机监控数据

### 接口信息

- **接口路径**: `/api/v1/monitoring/hosts/{host_id}/metrics`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | string | 是 | 开始时间（ISO 8601） |
| end_time | string | 是 | 结束时间（ISO 8601） |
| metrics | string | 否 | 指标列表（逗号分隔） |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "host_id": "host-abc123",
    "metrics": [
      {
        "timestamp": "2026-01-26T10:00:00Z",
        "cpu_usage_percent": 45.5,
        "memory_usage_percent": 60.0,
        "gpu_avg_utilization": 75.0
      }
    ]
  }
}
```

---

## 2. 获取 GPU 监控数据

### 接口信息

- **接口路径**: `/api/v1/monitoring/gpus/{gpu_id}/metrics`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | string | 是 | 开始时间 |
| end_time | string | 是 | 结束时间 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "gpu_id": 1,
    "metrics": [
      {
        "timestamp": "2026-01-26T10:00:00Z",
        "utilization_percent": 85.5,
        "memory_usage_percent": 75.0,
        "temperature": 68.5,
        "power_draw": 280.5
      }
    ]
  }
}
```

---

## 3. 获取环境监控数据

### 接口信息

- **接口路径**: `/api/v1/monitoring/environments/{env_id}/metrics`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | string | 是 | 开始时间 |
| end_time | string | 是 | 结束时间 |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env_id": "env-abc123",
    "metrics": [
      {
        "timestamp": "2026-01-26T10:00:00Z",
        "cpu_usage_percent": 45.5,
        "memory_usage_percent": 60.0,
        "gpu_usage_percent": 85.0
      }
    ]
  }
}
```

---

## 4. 获取系统概览

### 接口信息

- **接口路径**: `/api/v1/monitoring/overview`
- **请求方法**: `GET`
- **是否需要认证**: 是（需要管理员权限）

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_hosts": 10,
    "online_hosts": 8,
    "total_gpus": 40,
    "available_gpus": 15,
    "total_environments": 25,
    "running_environments": 18,
    "total_users": 150,
    "active_users": 45
  }
}
```
