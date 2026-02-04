# 项目开发计划 (Project Plan)

本文档用于追踪 RemoteGPU 后端项目的整体进度与里程碑。根据实际开发情况动态调整。

## 📅 阶段一：核心功能闭环 (Core MVP) - [已完成]
**目标**: 实现用户登录、机器管理、基础监控、API 文档，打通 B2B 业务主流程。

- [x] **鉴权模块**
  - [x] 登录/登出/刷新令牌
  - [x] 账号状态管控 (Active/Disabled)
  - [x] 角色权限 (Admin/Customer)
- [x] **管理端基础**
  - [x] 仪表盘数据聚合
  - [x] 镜像管理 (接口桩)
  - [x] 审计日志记录
- [x] **机器管理**
  - [x] 机器批量导入
  - [x] 机器分配 (事务保障)
  - [x] 机器回收 (状态流转)
- [x] **监控与告警**
  - [x] 实时状态快照 (DB聚合)
  - [x] 告警列表与确认
- [x] **文档**
  - [x] Swagger/OpenAPI 注解覆盖

## 📅 阶段二：基础设施集成 (Infrastructure Integration) - [待启动]
**目标**: 将 Mock/Stub 数据替换为真实基础设施对接，提升系统可用性。

- [ ] **监控系统对接** (Priority: High)
  - [ ] 部署 Prometheus + DCGM Exporter
  - [ ] 后端集成 Prometheus Client
  - [ ] 实现真实的 GPU 趋势图 (`/dashboard/gpu-trend`)
- [ ] **容器仓库对接** (Priority: Medium)
  - [ ] 部署/配置 Harbor
  - [ ] 后端集成 Harbor API
  - [ ] 实现真实的镜像同步逻辑 (`/images/sync`)
- [ ] **计算节点控制** (Priority: High)
  - [ ] 开发/集成 Node Agent
  - [ ] 实现机器回收时的异步清理 (Reset SSH, Wipe Data)
  - [ ] 实现机器分配时的环境初始化

## 📅 阶段三：性能与稳定性 (Performance & Stability) - [规划中]
**目标**: 应对高并发场景，保障数据一致性与系统响应速度。

- [ ] **缓存层建设**
  - [ ] 引入 Redis 缓存实时监控数据 (`/monitoring/realtime`)
  - [ ] 缓存高频配置数据 (如镜像列表)
- [ ] **消息队列**
  - [ ] 引入 MQ (Kafka/RabbitMQ) 解耦审计日志写入
  - [ ] 异步处理耗时的资源操作
- [ ] **测试与CI/CD**
  - [ ] 单元测试覆盖率 > 60%
  - [ ] 集成测试 (API Level)

## 📅 阶段四：企业级增强 (Enterprise Features) - [规划中]
- [ ] 多租户/组织架构支持
- [ ] 财务报表与账单导出
- [ ] 细粒度的 RBAC 权限系统
