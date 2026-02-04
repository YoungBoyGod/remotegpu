# 待办事项与技术债务 (Technical Debt & TODOs)

以下是在功能开发过程中识别出的已打通接口但内部逻辑尚需完善的模块。这些任务已被标记为已完成以符合当前阶段的目标（接口与流程闭环），但在生产环境上线前建议优先处理。

## 1. GPU 趋势图接入真实监控数据

- **涉及接口**: `GET /api/v1/admin/dashboard/gpu-trend`
- **当前状态**: 返回静态 Mock 数据。
- **未完成原因**: 需要依赖外部监控系统（如 Prometheus + Grafana 或 InfluxDB）以及相应的 Exporter 数据采集，这超出了当前后端逻辑开发的范畴，属于基础设施集成。
- **实现建议**:
  1.  **部署 Prometheus**: 确保集群中部署了 Prometheus Server。
  2.  **部署 DCGM Exporter**: 在每台 GPU 机器上部署 `NVIDIA/dcgm-exporter` 以采集 GPU 利用率、显存使用量等指标。
  3.  **后端集成**:
      - 引入 Prometheus Client SDK (`github.com/prometheus/client_golang/api`).
      - 在 `backend/internal/service/ops/monitor_service.go` 中实现 `queryPrometheus` 方法。
      - 使用 PromQL 查询（如 `avg(DCGM_FI_DEV_GPU_UTIL) by (host)`）获取历史趋势数据。

## 2. 镜像同步对接 Harbor

- **涉及接口**: `POST /api/v1/admin/images/sync`
- **当前状态**: 接口已定义，Service 层为空实现 (Stub)。
- **未完成原因**: 需要搭建 Harbor 实例并配置相应的 API 访问凭证（Endpoint, Username, Password），目前缺少这些环境配置。
- **实现建议**:
  1.  **配置管理**: 在 `config.yaml` 中增加 Harbor 配置项 (URL, AdminSecret)。
  2.  **客户端实现**:
      - 使用 Harbor Go SDK 或直接封装 HTTP Client 调用 Harbor v2.0 API (`/projects`, `/repositories`, `/artifacts`).
  3.  **同步逻辑**:
      - 遍历 Harbor 中的 Project 和 Repository。
      - 获取 Artifact (Tag) 的详细信息（大小、创建时间）。
      - 与本地数据库 `images` 表比对，执行 `Insert` 或 `Update` 操作。
      - (可选) 增加 Webhook 支持，让 Harbor 主动推送镜像更新事件。

## 3. 机器回收的异步清理流程

- **涉及接口**: `POST /api/v1/admin/machines/:id/reclaim`
- **当前状态**: 仅更新数据库中的机器状态为 `maintenance` 并记录审计日志。
- **未完成原因**: 需要与计算节点上的 Agent 组件通信，涉及 RPC/gRPC 调用或消息队列集成，且需要具体的清理脚本（如 `rm -rf /home/user`, 重置 SSH `authorized_keys`）。
- **实现建议**:
  1.  **Agent 通信**:
      - 定义 gRPC 接口 `ResetMachine(ctx, hostID)`.
      - 或者使用 Redis Pub/Sub / MQTT 发送清理指令。
  2.  **Agent 端实现**:
      - 接收指令后，执行系统清理脚本（清除用户进程、临时文件、Docker 容器）。
      - 重新生成 Host Key（可选）。
  3.  **状态回调**:
      - Agent 清理完成后，回调后端接口或更新状态，将机器状态从 `maintenance` 变更为 `idle`（可分配）。

## 4. 实时监控快照增加缓存

- **涉及接口**: `GET /api/v1/admin/monitoring/realtime`
- **当前状态**: 每次请求都直接查询数据库统计 (`Count` 聚合查询)。
- **未完成原因**: 当前数据量较小，直连数据库性能尚可，且 Redis 缓存逻辑属于优化项。
- **实现建议**:
  1.  **引入 Redis**: 确保 `backend/pkg/cache` 模块可用。
  2.  **缓存策略**:
      - 在 `GetGlobalSnapshot` 中，先读 Redis Key `monitor:snapshot:global`。
      - 如果未命中，查询 DB，计算结果，并写入 Redis，设置过期时间（如 10秒 或 30秒）。
  3.  **失效机制**:
      - (可选) 在机器状态变更（分配/回收）时主动失效缓存，或仅依赖自然过期（最终一致性）。

## 依赖汇总

- **Prometheus Client**: `github.com/prometheus/client_golang/api` (用于 GPU 监控)
- **Harbor SDK**: `github.com/goharbor/go-client` (可选，用于镜像同步)
- **RPC/Messaging**: `google.golang.org/grpc` 或 `redis` (用于 Agent 通信)

---

# 补充建议 (Additional Suggestions)

基于对 `backend/docs/claude/todo.md` 和 `backend/docs/codex/todo.md` 的分析，建议补充以下任务以完善企业级能力：

## 5. 通知系统 (Notification System)

- **涉及功能**: 告警触发、机器分配通知、密码重置邮件。
- **缺失原因**: 当前系统仅有数据层面的状态变更，缺乏主动触达用户的渠道。
- **实现建议**:
  1.  **渠道集成**: 实现 Email (SMTP), Slack/DingTalk Webhook 集成。
  2.  **事件驱动**: 在 `AlertService` 和 `AllocationService` 中埋点，通过 Channel 或 MQ 发送事件。
  3.  **用户配置**: 允许用户在 `/profile` 中设置接收偏好。
- **依赖**: `gopkg.in/gomail.v2` (邮件) 或 HTTP Client。

## 6. 用量报表与导出 (Usage Reporting)

- **涉及功能**: B2B 结算、资源使用审计。
- **缺失原因**: 虽不涉及在线支付，但企业客户需要月度资源使用清单进行线下结算或成本核算。
- **实现建议**:
  1.  **定时任务**: 每日/每月运行 Cron Job，统计 `allocations` 表时长。
  2.  **报表生成**: 生成 CSV/Excel 文件 (`github.com/xuri/excelize`)。
  3.  **API**: 提供 `GET /admin/reports/usage` 下载接口。

## 7. 系统加固 (System Hardening)

- **涉及功能**: 接口限流 (Rate Limiting)、单元测试。
- **缺失原因**: 当前主要关注功能实现，非功能性质量保障（测试覆盖率、防攻击）尚未提上日程。
- **实现建议**:
  1.  **限流**: 引入 `gin-contrib/rate-limiter` 或基于 Redis 的限流中间件，保护 `/auth/*` 和 `/admin/*` 接口。
  2.  **测试**: 为核心 Service (`AuthService`, `AllocationService`) 编写 Unit Test，目标覆盖率 > 60%。