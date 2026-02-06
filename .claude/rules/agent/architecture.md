---
paths:
  - "agent/**/*.go"
---

# Agent 模块架构规则

## 目录结构

```
agent/
├── cmd/            # 入口
├── internal/
│   ├── client/     # 与后端 API 通信的 HTTP 客户端
│   ├── executor/   # 任务执行器
│   ├── models/     # 数据模型
│   ├── queue/      # 任务队列管理
│   ├── scheduler/  # 任务调度器
│   └── store/      # 本地 SQLite 存储
└── docs/           # 设计文档
```

## 通信协议

- Agent 通过 HTTP 轮询后端 API（`/api/v1/agent/` 前缀）
- 任务生命周期：claim → start → renew lease → complete
- 配置文件：`agent.yaml`（参考 `agent.yaml.example`）
