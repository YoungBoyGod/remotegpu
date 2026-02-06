# CodeX 未完成任务清单

> 说明：记录当前尚未完成的技术债务、原因、建议方案与依赖。

## 1. GPU 趋势真实数据

- 位置：`internal/service/ops/monitor_service.go` (`GetGPUTrend`)
- 原因：目前返回固定模拟数据，缺少监控系统数据源。
- 实现建议：
  1. 在 `config` 中补充监控系统配置（Prometheus/InfluxDB）。
  2. 部署 GPU exporter（如 `dcgm-exporter`）并确保指标可采集。
  3. 封装监控查询客户端（例如 Prometheus HTTP API）。
  4. 添加时间范围与步进参数（如最近 24h，每 1h 一点）。
  5. 将采样结果映射为 `time/usage` 数组。
- 依赖：监控系统部署 + GPU exporter + 可访问 API；必要时加网络权限。

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

## 4. 任务启动写入 process_id

- 位置：任务运行器/调度器（待确定）与 `tasks.process_id`
- 原因：任务停止已依赖 `process_id`，但当前缺少回写来源。
- 实现建议：
  1. 任务启动时由 Agent/调度器返回 `process_id`。
  2. 回写 `tasks.process_id`，并同步更新 `status` 为 `running`。
  3. 若启动失败，写入 `error_msg` 并保持状态一致。
- 依赖：Agent 任务启动接口或任务运行器落地。

## 5. 远程访问服务（VNC/RDP/Guacamole）

- 位置：`internal/service/vnc.go`, `internal/service/rdp.go`, `internal/service/guacamole.go`
- 原因：当前仅有框架代码，缺少实际实现与运行环境集成。
- 实现建议：
  1. 选型并落地 VNC/RDP 服务（TigerVNC/xrdp）。
  2. 在 Agent 或容器侧生成配置文件与启动脚本。
  3. 若使用 Guacamole，封装 Guacamole API 客户端并创建连接。
  4. 返回统一访问地址并记录到数据库。
- 依赖：VNC/RDP 服务、Guacamole 服务（可选）、容器运行环境。

## 6. DNS 管理

- 位置：`internal/service/dns.go`
- 原因：多云 DNS 适配尚未实现。
- 实现建议：
  1. 定义 DNS Provider 接口并在配置中启用。
  2. 分别实现 Cloudflare/阿里云/腾讯云/AWS 等适配器。
  3. 对接环境创建/释放时自动生成或清理记录。
- 依赖：各云厂商 SDK 与 API 凭证。

## 7. 防火墙管理

- 位置：`internal/service/firewall.go`
- 原因：iptables/firewalld/云厂商安全组逻辑未实现。
- 实现建议：
  1. 明确执行环境（Agent/云厂商 API）。
  2. 实现端口映射 CRUD（创建/更新/删除）。
  3. 对接环境生命周期（分配/回收时自动配置）。
- 依赖：Agent 权限或云厂商安全组 API。

## 8. Docker 容器管理

- 位置：`internal/service/docker.go`
- 原因：核心容器 CRUD/日志/资源统计仍为 TODO。
- 实现建议：
  1. 引入 Docker SDK 并封装基础操作。
  2. 增加 GPU 与存储挂载配置。
  3. 对接任务生命周期（创建/停止/清理）。
- 依赖：Docker API 权限、nvidia-container-toolkit（GPU）。

## 9. Agent 任务恢复租约校验

- 位置：`/home/luo/code/remotegpu/agent/internal/scheduler/scheduler.go` (`recover`)
- 原因：当前仅恢复 pending 任务，未校验 assigned/running 任务的租约有效性。
- 实现建议：
  1. 启动时加载 assigned/running 任务。
  2. 通过 Server 校验 attempt_id/lease_expires_at。
  3. 过期则停止进程并标记 failed；有效则继续执行并续约。
- 依赖：Server 提供任务租约校验或状态查询接口。
