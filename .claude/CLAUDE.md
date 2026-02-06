# RemoteGPU 项目指南

## 项目概述

RemoteGPU 是一个 GPU 远程管理平台，包含三个子模块：
- `backend/` — Go 后端 API 服务（Gin + GORM + PostgreSQL）
- `frontend/` — Vue 3 前端（TypeScript + Element Plus + Vite）
- `agent/` — Go Agent 客户端（部署在 GPU 机器上）

## 常用命令

```bash
# 后端
cd backend && go build ./cmd/...        # 编译
cd backend && go test ./...              # 测试
cd backend && go run ./cmd/main.go server  # 启动

# 前端
cd frontend && npm install               # 安装依赖
cd frontend && npm run dev               # 开发服务器
cd frontend && npm run build             # 构建
cd frontend && npm run type-check        # 类型检查

# Agent
cd agent && go build ./cmd/...           # 编译
cd agent && go test ./...                # 测试
```

## 语言

- 代码注释、commit message、文档统一使用中文
- 变量名、函数名、类型名使用英文

## 架构规则

详见各子目录的规则文件：
- @.claude/rules/backend/architecture.md
- @.claude/rules/frontend/architecture.md
- @.claude/rules/agent/architecture.md

## 数据库

- 使用 PostgreSQL，ORM 为 GORM
- 迁移脚本在 `backend/sql/` 下，按编号递增命名：`NN_描述.sql`
- 新增表或字段必须同时添加迁移脚本

## Git 约定

- commit message 格式：`type(scope): 描述`
- type: feat / fix / docs / refactor / test / chore
- scope: backend / frontend / agent / 具体模块名
