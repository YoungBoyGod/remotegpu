# RemoteGPU API 文档

> RemoteGPU 系统 RESTful API 接口文档
>
> 版本：v1.0
> 基础 URL：`https://api.remotegpu.com/v1`

---

## 📋 目录结构

```
api/
├── README.md                    # 本文件
├── common.md                    # 公共规范（认证、分页、错误码等）
├── 01_auth.md                   # 认证授权接口
├── 02_users.md                  # 用户管理接口
├── 03_workspaces.md             # 工作空间接口
├── 04_hosts.md                  # 主机管理接口
├── 05_environments.md           # 环境管理接口
├── 06_datasets.md               # 数据集管理接口
├── 07_models.md                 # 模型管理接口
├── 08_images.md                 # 镜像管理接口
├── 09_training.md               # 训练任务接口
├── 10_inference.md              # 推理服务接口
├── 11_monitoring.md             # 监控接口
├── 12_billing.md                # 计费接口
├── 13_notifications.md          # 通知接口
├── 14_alerts.md                 # 告警接口
├── 15_webhooks.md               # Webhook 接口
└── 16_issues.md                 # 问题单/需求单接口
```

---

## 🚀 快速开始

### 1. 获取 API Token

```bash
curl -X POST https://api.remotegpu.com/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

### 2. 使用 Token 调用 API

```bash
curl -X GET https://api.remotegpu.com/v1/environments \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📖 文档说明

### 接口格式

每个接口文档包含以下内容：

- **接口描述**：接口的功能说明
- **请求方法**：GET、POST、PUT、DELETE 等
- **请求路径**：API 端点路径
- **请求参数**：路径参数、查询参数、请求体
- **响应示例**：成功和失败的响应示例
- **错误码**：可能返回的错误码

### 通用规范

所有接口遵循以下规范：

- **认证方式**：JWT Token（Bearer Token）
- **请求格式**：JSON
- **响应格式**：JSON
- **字符编码**：UTF-8
- **时间格式**：ISO 8601（`2026-01-26T10:00:00Z`）

详见 [公共规范文档](./common.md)

---

## 🔐 认证说明

RemoteGPU API 使用 JWT Token 进行认证。

### 获取 Token

通过登录接口获取 Access Token 和 Refresh Token：

```
POST /auth/login
```

### 使用 Token

在请求头中携带 Token：

```
Authorization: Bearer {access_token}
```

### Token 有效期

- **Access Token**：24 小时
- **Refresh Token**：30 天

详见 [认证接口文档](./01_auth.md)

---

## 📊 模块概览

| 模块 | 文档 | 说明 |
|------|------|------|
| 认证授权 | [01_auth.md](./01_auth.md) | 登录、注册、Token 刷新 |
| 用户管理 | [02_users.md](./02_users.md) | 用户信息、配额管理 |
| 工作空间 | [03_workspaces.md](./03_workspaces.md) | 工作空间、成员管理 |
| 主机管理 | [04_hosts.md](./04_hosts.md) | 主机注册、监控 |
| 环境管理 | [05_environments.md](./05_environments.md) | 创建、启动、停止环境 |
| 数据集 | [06_datasets.md](./06_datasets.md) | 数据集上传、版本管理 |
| 模型 | [07_models.md](./07_models.md) | 模型上传、版本管理 |
| 镜像 | [08_images.md](./08_images.md) | 镜像列表、详情 |
| 训练任务 | [09_training.md](./09_training.md) | 创建、监控训练任务 |
| 推理服务 | [10_inference.md](./10_inference.md) | 部署、管理推理服务 |
| 监控 | [11_monitoring.md](./11_monitoring.md) | 获取监控数据 |
| 计费 | [12_billing.md](./12_billing.md) | 计费记录、账单 |
| 通知 | [13_notifications.md](./13_notifications.md) | 通知列表、标记已读 |
| 告警 | [14_alerts.md](./14_alerts.md) | 告警规则、记录 |
| Webhook | [15_webhooks.md](./15_webhooks.md) | Webhook 配置 |
| 工单 | [16_issues.md](./16_issues.md) | 问题单、需求单 |

---

## 🔗 相关资源

- [系统架构文档](../design/system_architecture.md)
- [数据库设计文档](../design/database_design.md)
- [需求文档](../requirements/)

---

**创建日期**：2026-01-26
**维护者**：RemoteGPU 团队
