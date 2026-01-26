# 训练任务 API

> 训练任务创建、管理、监控相关接口

---

## 1. 创建训练任务

### 接口信息

- **接口路径**: `/api/v1/training/jobs`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 请求参数

```json
{
  "name": "ResNet Training",
  "description": "图像分类模型训练",
  "workspace_id": 456,
  "framework": "pytorch",
  "script_path": "/workspace/train.py",
  "dataset_id": 789,
  "cpu": 8,
  "memory": 34359738368,
  "gpu": 2,
  "priority": 5
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 任务名称 |
| description | string | 否 | 描述 |
| workspace_id | int | 否 | 工作空间 ID |
| framework | string | 是 | 框架：pytorch, tensorflow |
| script_path | string | 是 | 训练脚本路径 |
| dataset_id | int | 否 | 数据集 ID |
| cpu | int | 是 | CPU 核心数 |
| memory | int | 是 | 内存（字节） |
| gpu | int | 是 | GPU 数量 |
| priority | int | 否 | 优先级（0-10） |

### 响应示例

**成功响应 (201)**

```json
{
  "code": 0,
  "message": "训练任务创建成功",
  "data": {
    "job_id": "job-abc123",
    "name": "ResNet Training",
    "status": "pending",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 2. 获取训练任务列表

### 接口信息

- **接口路径**: `/api/v1/training/jobs`
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
        "job_id": "job-abc123",
        "name": "ResNet Training",
        "status": "running",
        "framework": "pytorch",
        "cpu": 8,
        "gpu": 2,
        "started_at": "2026-01-26T10:05:00Z",
        "created_at": "2026-01-26T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 12,
      "total_pages": 1
    }
  }
}
```

---

## 3. 获取训练任务详情

### 接口信息

- **接口路径**: `/api/v1/training/jobs/{job_id}`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "job_id": "job-abc123",
    "name": "ResNet Training",
    "description": "图像分类模型训练",
    "status": "running",
    "framework": "pytorch",
    "script_path": "/workspace/train.py",
    "dataset_id": 789,
    "model_id": 456,
    "cpu": 8,
    "memory": 34359738368,
    "gpu": 2,
    "priority": 5,
    "env_id": "env-xyz789",
    "started_at": "2026-01-26T10:05:00Z",
    "created_at": "2026-01-26T10:00:00Z"
  }
}
```

---

## 4. 停止训练任务

### 接口信息

- **接口路径**: `/api/v1/training/jobs/{job_id}/stop`
- **请求方法**: `POST`
- **是否需要认证**: 是

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "训练任务停止中"
}
```

---

## 5. 删除训练任务

### 接口信息

- **接口路径**: `/api/v1/training/jobs/{job_id}`
- **请求方法**: `DELETE`
- **是否需要认证**: 是

### 响应示例

**成功响应 (204)**

```json
{
  "code": 0,
  "message": "训练任务删除成功"
}
```

---

## 6. 获取训练任务日志

### 接口信息

- **接口路径**: `/api/v1/training/jobs/{job_id}/logs`
- **请求方法**: `GET`
- **是否需要认证**: 是

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| lines | int | 否 | 日志行数（默认 100） |

### 响应示例

**成功响应 (200)**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "logs": [
      "Epoch 1/100, Loss: 0.5234",
      "Epoch 2/100, Loss: 0.4123"
    ]
  }
}
```
