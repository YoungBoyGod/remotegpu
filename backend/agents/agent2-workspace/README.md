# Agent 2 - 工作空间与组织管理模块

**任务ID**: #3
**依赖**: 任务#7（模块F）、任务#2（模块A）
**状态**: 等待依赖任务完成

## 🎯 任务目标

开发工作空间与组织管理模块（Module B）。

## 📋 核心功能

- 工作空间创建、更新、删除
- 工作空间成员管理
- 角色权限（owner/admin/member/viewer）
- 工作空间类型（personal/team/enterprise）
- 软删除支持

## 📁 涉及文件

```
internal/
├── model/entity/workspace.go
├── controller/v1/workspace.go
├── service/workspace.go
└── dao/workspace.go
```

## 🔌 API端点

- POST /api/v1/workspaces
- GET /api/v1/workspaces
- GET /api/v1/workspaces/:id
- PUT /api/v1/workspaces/:id
- DELETE /api/v1/workspaces/:id
- POST /api/v1/workspaces/:id/members
- DELETE /api/v1/workspaces/:id/members/:user_id

## ✅ 工作清单

- [ ] 等待任务#7和#2完成
- [ ] 实现工作空间CRUD功能
- [ ] 实现成员管理功能
- [ ] 实现权限控制
- [ ] 编写单元测试（目标覆盖率>80%）
- [ ] 创建模块文档 `docs/modules/workspace.md`
- [ ] 更新API文档 `docs/api/workspace.md`
- [ ] 提交代码（使用 `[ModuleB-Workspace]` 前缀）

## 🔧 Git提交规范

```bash
[ModuleB-Workspace] 功能描述

详细说明

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

## 📊 进度记录

在 `progress.md` 中记录每日进度。
