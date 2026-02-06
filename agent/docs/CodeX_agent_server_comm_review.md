# CodeX Review: Agent 与 Server 通信实现

- 作者: CodeX
- 日期: 2026-02-05
- 范围: `/home/luo/code/remotegpu/agent/internal/*`, `/home/luo/code/remotegpu/agent/cmd/*`

## 结论

当前 Agent 已具备基础任务执行与轮询框架，但与 Server 通信链路未完整接入：任务认领后的状态/租约/完成结果没有被上报，导致 Server 无法正确感知任务生命周期。建议补齐“启动上报 + 租约续约 + 完成上报”链路，并修正认领任务的本地状态覆盖问题。

## 主要问题（按严重程度）

### High

1) 任务认领后未上报启动/完成
- 位置: `/home/luo/code/remotegpu/agent/internal/scheduler/scheduler.go`, `/home/luo/code/remotegpu/agent/internal/executor/executor.go`
- 问题: Scheduler/Executor 完全不调用 ServerClient 的 ReportStart/ReportComplete。
- 影响: Server 无法更新任务状态，租约无法续约，任务可能被重复认领。
- 改进: 在执行前调用 ReportStart；执行中续约租约；执行后调用 ReportComplete。

2) 认领任务状态被本地覆盖
- 位置: `/home/luo/code/remotegpu/agent/internal/scheduler/scheduler.go` (`Submit`)
- 问题: Submit 强制设置 `Status=pending` 与 `CreatedAt=now`，会覆盖 Server 已分配任务（assigned/attempt_id）。
- 影响: 丢失 attempt_id/assigned 状态，破坏租约语义。
- 改进: 仅在字段为空时初始化，保留 Server 下发的 assigned/attempt。

### Medium

3) 续约机制未实现
- 位置: `/home/luo/code/remotegpu/agent/internal/scheduler/scheduler.go`
- 问题: Poller 认领后没有对长任务进行续约。
- 影响: 长任务被 Server 视为租约过期，重复分配。
- 改进: 执行期间启用续约 goroutine，任务完成时停止续约。

4) Claim 请求缺少 request_id 幂等键
- 位置: `/home/luo/code/remotegpu/agent/internal/client/client.go`
- 问题: ClaimTasks 未携带 request_id，重复请求难以幂等。
- 影响: 可能重复分配或丢任务。
- 改进: 生成 request_id 并随请求发送。

## 为什么这样不好 / 我会如何改好

- 不上报状态会让任务生命周期断链，无法保证“认领 → 执行 → 完成”一致性；我的改法将执行前后与续约统一接入 ServerClient，保证可观测与幂等。
- 本地覆盖任务状态会抹掉 attempt_id，导致租约校验失败；我的改法仅在字段为空时设置默认值，避免覆盖 Server 权威数据。
- 无幂等键会让网络重试产生不确定行为；我的改法补充 request_id 作为幂等键。
