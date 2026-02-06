---
paths:
  - "frontend/src/**/*.{ts,vue}"
---

# 前端架构规则

## 目录结构

```
frontend/src/
├── api/            # API 请求函数（按模块分文件）
├── components/     # 公共组件
│   ├── common/     # 通用组件（DataTable 等）
│   └── layout/     # 布局组件（AdminLayout、AdminSidebar 等）
├── config/         # 配置（tableColumns 等）
├── router/         # 路由定义（index.ts）
├── stores/         # Pinia 状态管理
├── types/          # TypeScript 类型定义（按模块分文件）
└── views/          # 页面视图
    ├── admin/      # 管理员页面
    └── customer/   # 客户页面
```

## 新增页面的标准流程

1. `types/` — 定义接口类型
2. `api/` — 添加 API 请求函数
3. `views/` — 创建页面组件
4. `router/index.ts` — 注册路由（放在通配路由 `pathMatch` 之前）
5. 如需侧边栏入口 — 修改对应 Sidebar 组件
