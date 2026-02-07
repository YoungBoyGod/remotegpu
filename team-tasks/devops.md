# DevOps (devops) 任务清单

## 职责
- 基础设施部署与维护
- CI/CD 流水线
- 容器化与编排
- 监控系统部署

## 任务列表

### T-D01 [P0] 基础设施部署
- 部署/验证 PostgreSQL
- 部署 Redis 缓存服务
- 部署 Prometheus + Grafana
- 验证 docker-compose 配置

### T-D02 [P0] Agent 部署方案
- 制定 Agent 部署脚本
- Agent 自动更新机制
- Agent 配置管理

### T-D03 [P1] GPU Exporter 部署
- 在 GPU 机器上部署 nvidia_gpu_exporter
- 配置 Prometheus 抓取目标
- 验证 GPU 指标采集

### T-D04 [P1] Nginx 反向代理配置
- 前端静态资源部署
- 后端 API 反向代理
- SSL/TLS 证书配置

### T-D05 [P2] CI/CD 流水线
- 后端自动构建与测试
- 前端自动构建与部署
- Agent 自动构建与发布
