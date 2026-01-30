# Agent 6 - 公共基础设施模块

**任务ID**: #7
**优先级**: 最高（其他所有模块依赖此模块）
**状态**: 进行中

## 🎯 任务目标

开发和完善公共基础设施模块（Module F），为其他所有模块提供基础支持。

## 📋 核心功能

- 数据库连接管理（PostgreSQL + GORM）
- Redis缓存管理
- 日志系统（Zap）
- K8s客户端封装
- 存储后端（Local/S3/RustFS）
- 统一响应格式
- 错误处理
- 健康检查
- 优雅启动/关闭
- 热更新管理

## 📁 涉及文件

```
pkg/
├── auth/         - JWT认证
├── database/     - 数据库连接
├── errors/       - 错误处理
├── graceful/     - 优雅关闭
├── health/       - 健康检查
├── hotreload/    - 热更新
├── k8s/          - K8s客户端
├── logger/       - 日志系统
├── redis/        - Redis缓存
├── response/     - 响应格式
└── storage/      - 存储后端
```

## ✅ 工作清单

- [ ] 审查现有代码实现
- [ ] 补充单元测试（目标覆盖率>80%）
- [ ] 编写集成测试
- [ ] 创建模块文档 `docs/modules/infrastructure-base.md`
- [ ] 创建环境配置指南 `docs/setup/`
- [ ] 提交代码（使用 `[ModuleF-Base]` 前缀）

## 📝 文档要求

- 记录每个子模块的使用方法
- 记录数据库连接和迁移方案
- 记录K8s客户端使用方法
- 记录存储后端配置
- 记录日志规范和错误处理规范

## 🔧 Git提交规范

```bash
[ModuleF-Base] 功能描述

详细说明

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

## 📊 进度记录

在 `progress.md` 中记录每日进度。
