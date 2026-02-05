# CodeX Review: task-queue-design

- 作者: CodeX
- 日期: 2026-02-05
- 对象: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md`

## 结论

文档结构清晰、覆盖面完整，但存在多处架构决策与状态模型不一致的问题，若不统一将导致实现阶段分歧与状态错乱。建议先统一“通信模式”和“任务状态/租约字段”，再细化 API 与数据模型。

## 主要问题（按严重程度排序）

### High

1) 架构决策自相矛盾（Pull+WebSocket vs 纯 Pull）
- 位置: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md:241`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:266`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:508`
- 问题: 4.2/4.3 明确采用 WebSocket 通知 + Pull，但 10.1 又明确“纯 Pull 模式、避免 WebSocket 复杂性”，API 也保留了“Server 推送”接口。
- 影响: 实现阶段会出现两套互斥机制，导致服务端/Agent 设计反复。
- 建议: 选定一种模式后统一调整 4.x、5.x、10.x 与 API 列表，删除另一模式的接口与流程描述。

2) 任务状态/分配字段不一致
- 位置: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md:40`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:75`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:150`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:602`
- 问题: 状态列表包含 `assigned`，状态流转图却只有 `queued`；SQL 示例用 `status='assigned'`，但主表未定义 `assigned_to/assigned_at/lease_expires_at` 字段，后续租约机制又新增字段但未回写到主表设计。
- 影响: 任务认领、续约、重试等流程难以落地，容易出现重复领取或不可恢复的卡死状态。
- 建议: 定义单一的状态集与流转图（例如 pending → assigned → running → completed/failed/cancelled），并在主表中补齐 `assigned_agent_id/assigned_at/lease_expires_at/attempt_id` 等字段。

### Medium

3) 结果存储位置与 Task 模型重复
- 位置: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md:40`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:166`
- 问题: Task 结构包含 `stdout/stderr`，同时又设计 `task_results` 表存储结果；未明确哪个是权威来源。
- 影响: 可能造成存储冗余与读取口径不一致。
- 建议: 明确结果只在 `task_results` 存储，Task 仅保存摘要或指针；或在文档中标注 Task 的 `stdout/stderr` 仅用于 Agent 本地缓存。

4) Server 端数据模型与现有系统命名不一致
- 位置: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md:154`、`/home/luo/code/remotegpu/agent/docs/task-queue-design.md:163`
- 问题: `machine_id REFERENCES machines(id)`、`created_by REFERENCES users(id)` 与现有后端实体（hosts/customers）不一致。
- 影响: 迁移到现有系统时需要额外映射，易出错。
- 建议: 将表引用与当前系统实体名对齐，或注明“这是逻辑名，后续需映射到 hosts/customers”。

### Low

5) Token 绑定 IP 的可用性风险
- 位置: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md:409`
- 问题: Token 绑定 IP 在 NAT/云环境可能发生变化（弹性 IP、代理、容器网络）。
- 影响: Agent 可能被误踢或无法重连。
- 建议: 改为绑定 AgentID + 短期签名（设备指纹/证书），允许有限 IP 变更白名单。

6) Agent 重启恢复流程缺少“运行中进程”判定策略
- 位置: `/home/luo/code/remotegpu/agent/docs/task-queue-design.md:323`
- 问题: 重启时直接将 running 标记 failed，未考虑进程仍在运行或可恢复的情况。
- 影响: 可能导致重复执行或任务误失败。
- 建议: 增加进程探测/attempt_id 校验，确认任务已停止后再标记失败。

## 建议优先修复顺序

1. 统一通信模式（Pull vs Pull+WS）并更新 API 列表。
2. 统一任务状态机与租约字段，更新主表设计与 API 流程。
3. 明确结果存储的单一来源与读写口径。
