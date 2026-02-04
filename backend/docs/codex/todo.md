# CodeX 未完成任务清单

> 说明：记录当前尚未完成的技术债务、原因、建议方案与依赖。

## 1. GPU 趋势真实数据

- 位置：`internal/service/ops/monitor_service.go` (`GetGPUTrend`)
- 原因：目前返回固定模拟数据，缺少监控系统数据源。
- 实现建议：
  1. 在 `config` 中补充监控系统配置（Prometheus/InfluxDB）。
  2. 封装监控查询客户端（例如 Prometheus HTTP API）。
  3. 添加时间范围与步进参数（如最近 24h，每 1h 一点）。
  4. 将采样结果映射为 `time/usage` 数组。
- 依赖：监控系统部署 + 可访问 API；必要时加网络权限。

## 2. 监控快照缓存

- 位置：`internal/service/ops/monitor_service.go` (`GetGlobalSnapshot`) 的 TODO
- 原因：频繁查询数据库；当前无缓存与采样频率控制。
- 实现建议：
  1. 使用 Redis 缓存快照（key：`monitor:snapshot`）。
  2. 缓存 TTL 30-60s。
  3. 请求时优先读缓存，未命中再查询数据库并写缓存。
- 依赖：Redis 可用（`pkg/cache` 已接入）。

## 3. 镜像同步（Harbor/本地仓库）

- 位置：`internal/service/image/image_service.go` (`Sync`)
- 原因：未接入 Harbor API，当前仅占位。
- 实现建议：
  1. 读取 `config.harbor` 配置。
  2. 调用 Harbor API 获取项目与仓库镜像列表。
  3. 解析镜像名称、tag、大小、创建时间，写入 `images` 表。
  4. 做去重与状态更新（已存在则更新元数据）。
- 依赖：Harbor API 访问凭据、网络可达。

## 4. 机器分配后的异步动作

- 位置：`internal/service/allocation/allocation_service.go` (`AllocateMachine`, `ReclaimMachine`)
- 原因：异步清理/重置 SSH 逻辑未实现。
- 实现建议：
  1. 在分配后调用 Agent 重置 SSH / 初始化环境。
  2. 回收后触发清理任务（数据清理、镜像清理）。
  3. 可用 goroutine 或任务队列（如 Redis + worker）。
- 依赖：Agent 接口可用、异步任务框架（可选）。

## 5. 任务停止的 Agent 集成

- 位置：`internal/service/task/task_service.go` (`StopTask` / `StopTaskWithAuth`)
- 原因：当前仅更新状态，未通知实际运行环境。
- 实现建议：
  1. 调用 Agent API 终止容器/进程。
  2. Agent 返回成功后再更新任务状态。
  3. 失败时保持原状态并记录错误。
- 依赖：Agent 任务管理接口。
