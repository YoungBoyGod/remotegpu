# API 接口目录结构

本目录按照模块划分组织所有的 API 接口定义。

## 目录结构

```
api/
├── common/              # 公共模块
│   ├── types.ts        # 通用类型定义（分页、响应等）
│   ├── request.ts      # axios 请求封装
│   └── index.ts        # 公共模块导出
├── auth/               # 用户与权限模块
│   ├── types.ts        # 类型定义
│   ├── index.ts        # API 接口
│   └── README.md       # 模块说明
├── cmdb/               # CMDB 设备管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── environment/        # 环境管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── scheduler/          # 资源调度模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── storage/            # 数据与存储模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── image/              # 镜像管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── training/           # 训练与推理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── billing/            # 计费管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── monitoring/         # 监控告警模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── artifact/           # 制品管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── issue/              # 问题单管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── requirement/        # 需求单管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── notification/       # 通知管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
├── webhook/            # Webhook 管理模块
│   ├── types.ts
│   ├── index.ts
│   └── README.md
└── index.ts            # 统一导出所有模块
```

## 使用方式

```typescript
// 导入特定模块的 API
import { login, register } from '@/api/auth'
import { getEnvironmentList, createEnvironment } from '@/api/environment'

// 或者导入所有 API
import * as authApi from '@/api/auth'
import * as envApi from '@/api/environment'
```

## 模块说明

### 1. common - 公共模块
- 通用类型定义（分页、响应结构等）
- axios 请求封装
- 通用工具函数

### 2. auth - 用户与权限模块
- 用户注册、登录、登出
- 工作空间管理
- 权限控制
- 配额管理

### 3. cmdb - CMDB 设备管理模块
- 资产管理
- 服务器管理
- GPU 设备管理
- 设备状态管理

### 4. environment - 环境管理模块
- 开发环境创建、启动、停止
- SSH/RDP 访问管理
- 数据集挂载
- 环境配置

### 5. scheduler - 资源调度模块
- 资源调度策略
- 端口管理
- 调度历史查询

### 6. storage - 数据与存储模块
- 数据集管理
- 模型管理
- 文件上传下载
- 对象存储

### 7. image - 镜像管理模块
- 官方镜像查询
- 自定义镜像构建
- 镜像版本管理

### 8. training - 训练与推理模块
- 训练任务管理
- 推理服务部署
- 实验管理

### 9. billing - 计费管理模块
- 账户余额查询
- 计费记录
- 账单管理
- 充值支付

### 10. monitoring - 监控告警模块
- 资源监控
- 告警规则
- 告警历史

### 11. artifact - 制品管理模块
- 软件包管理
- 版本控制
- 依赖管理

### 12. issue - 问题单管理模块
- 问题创建与跟踪
- 工单流转
- 评论与附件

### 13. requirement - 需求单管理模块
- 需求管理
- 需求评审
- Sprint 管理

### 14. notification - 通知管理模块
- 消息推送
- 通知历史
- 通知设置

### 15. webhook - Webhook 管理模块
- Webhook 配置
- 事件订阅
- 回调管理
