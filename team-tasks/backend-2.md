# 后端开发 2 (backend-2) 任务清单

## 职责
- 监控与缓存模块
- 镜像管理
- 存储与数据集服务

## 任务列表

### T-B2-01 [P1] Redis 缓存接入
- 监控快照缓存（30s TTL）
- Token 黑名单优化
- 位置: internal/service/ops/monitor_service.go

### T-B2-02 [P1] GPU 监控数据接入
- 接入 Prometheus 查询 API
- 实现 GPU 利用率实时数据
- 实现 GPU 趋势图历史数据
- 位置: internal/service/ops/monitor_service.go

### T-B2-03 [P1] Harbor 镜像同步
- 接入 Harbor API
- 实现镜像列表同步
- 去重与状态更新
- 位置: internal/service/image/image_service.go

### T-B2-04 [P1] 机器 IP 唯一性校验
- 添加/导入机器时校验 IP 唯一
- 位置: internal/service/machine/machine_service.go

### T-B2-05 [P2] 文档管理模块完善
- 完善 document CRUD 接口
- 文件上传与存储
- 位置: backend/internal/service/document/
