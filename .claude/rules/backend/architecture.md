---
paths:
  - "backend/**/*.go"
---

# 后端架构规则

## 分层结构

```
backend/
├── api/v1/           # 请求/响应结构体定义（纯数据结构，无逻辑）
├── internal/
│   ├── controller/v1/ # HTTP 处理层（参数绑定、调用 Service、返回响应）
│   ├── service/       # 业务逻辑层（每个模块一个子目录）
│   ├── dao/           # 数据访问层（直接操作 GORM）
│   ├── model/entity/  # 数据库实体定义
│   ├── middleware/     # Gin 中间件
│   └── router/        # 路由注册（router.go）
├── pkg/               # 公共工具包（response、errors、auth、cache 等）
└── sql/               # 数据库迁移脚本
```

## 新增功能的标准流程

1. `model/entity/` — 定义实体，指定 `TableName()`
2. `dao/` — 创建 DAO，简单实体继承 `BaseDao[T]`
3. `service/模块名/` — 创建 Service，注入 DAO
4. `api/v1/` — 定义请求/响应结构体
5. `controller/v1/模块名/` — 创建 Controller，嵌入 `common.BaseController`
6. `router/router.go` — 注册路由，初始化 Service 和 Controller
