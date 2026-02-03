# RemoteGPU 项目完成计划

> 详细的项目完成顺序和清单
>
> 创建日期：2026-01-28
> 状态：规划中

---

## 📊 项目现状分析

### 已完成部分
- ✅ 前端框架搭建（Vue3 + Element Plus）
- ✅ 前端视图文件（29个页面组件）
- ✅ API 文档（16个模块完整文档）
- ✅ 数据库设计文档
- ✅ 系统架构设计文档
- ✅ 基础设施配置文件（Docker Compose）
- ✅ 后端基础框架（Gin + GORM）
- ✅ 基础中间件（CORS、Logger、Recovery、Auth、Role）
- ✅ 基础 pkg 包（auth、database、errors、health、logger、redis、response）

### 待完成部分
- ❌ 后端业务模块实现（仅完成用户/客户管理）
- ❌ 数据库表结构完整实现
- ❌ 基础设施服务配置完善
- ❌ 前后端接口对接
- ❌ 单元测试和集成测试
- ❌ 部署配置和 CI/CD

### 关键问题
1. **backend/main.go 文件为空**（导致语法错误）
2. **后端实现严重滞后**（仅11个 Go 文件，只有用户管理）
3. **数据库迁移不完整**（只迁移了 Customer 表）
4. **基础设施服务未完全配置**（缺少 .env 文件）

---

## 🎯 完成顺序总览

```
阶段 0: 紧急修复（1天）
  └─ 修复 main.go 语法错误

阶段 1: 基础设施层（3-5天）
  ├─ 数据库表结构实现
  ├─ 基础设施服务配置
  └─ 配置文件完善

阶段 2: 核心业务模块（15-20天）
  ├─ 认证授权模块
  ├─ 用户和工作空间模块
  ├─ 主机管理模块
  ├─ 环境管理模块（核心）
  ├─ 数据集和模型模块
  ├─ 镜像管理模块
  └─ 训练和推理模块

阶段 3: 支撑服务模块（10-15天）
  ├─ 监控模块
  ├─ 计费模块
  ├─ 通知和告警模块
  └─ Webhook 和工单模块

阶段 4: 集成和测试（5-7天）
  ├─ 前后端接口对接
  ├─ 单元测试
  ├─ 集成测试
  └─ 端到端测试

阶段 5: 部署和文档（3-5天）
  ├─ Docker 镜像构建
  ├─ Kubernetes 配置
  ├─ CI/CD 配置
  └─ 部署文档完善
```

---

## 📋 详细完成清单

### 阶段 0: 紧急修复 ⚠️

**优先级：最高 | 预计时间：1天**

#### 0.1 修复 backend/main.go 语法错误

**文件位置：** `backend/main.go`

**问题：** 文件为空，导致语法错误

**解决方案：**
- [ ] 删除空的 `backend/main.go` 文件
- [ ] 确认 `backend/cmd/main.go` 是正确的入口文件
- [ ] 更新 `go.mod` 中的 module 路径（如需要）
- [ ] 测试编译：`cd backend && go build ./cmd/main.go`

**依赖：** 无

**输出：** 后端可以正常编译运行

---

### 阶段 1: 基础设施层 🏗️

**优先级：高 | 预计时间：3-5天**

#### 1.1 数据库表结构实现

**文件位置：** `backend/internal/model/entity/`

**任务清单：**
- [ ] 创建 `user.go` - 用户表（已存在，需完善）
- [ ] 创建 `workspace.go` - 工作空间表
- [ ] 创建 `host.go` - 主机表
- [ ] 创建 `environment.go` - 环境表
- [ ] 创建 `dataset.go` - 数据集表
- [ ] 创建 `model.go` - 模型表
- [ ] 创建 `image.go` - 镜像表
- [ ] 创建 `training_job.go` - 训练任务表
- [ ] 创建 `inference_service.go` - 推理服务表
- [ ] 创建 `billing_record.go` - 计费记录表
- [ ] 创建 `notification.go` - 通知表
- [ ] 创建 `alert.go` - 告警表
- [ ] 创建 `webhook.go` - Webhook 表
- [ ] 创建 `issue.go` - 工单表

**参考文档：** `docs/design/database_design.md`, `docs/sql/`

**依赖：** 阶段 0 完成

#### 1.2 更新数据库迁移

**文件位置：** `backend/cmd/main.go`

**任务清单：**
- [ ] 在 `AutoMigrate` 中添加所有实体表
- [ ] 创建初始化数据脚本 `backend/sql/init.sql`
- [ ] 创建种子数据脚本 `backend/sql/seed.sql`

#### 1.3 基础设施服务配置

**文件位置：** `docker-compose/`

**任务清单：**
- [ ] PostgreSQL: 创建 `.env.example` 和配置文档
- [ ] Redis: 验证配置
- [ ] Etcd: 验证配置
- [ ] MinIO: 配置和测试
- [ ] Harbor: 完善 `.env.example` 和配置
- [ ] Prometheus: 完善配置和告警规则
- [ ] Grafana: 完善仪表板配置
- [ ] Nginx: 完善反向代理配置
- [ ] JumpServer: 配置和集成测试
- [ ] Guacamole: 配置和集成测试
- [ ] Uptime-kuma: 配置健康检查

**参考文档：** `docs/infrastructure/infrastructure.md`

#### 1.4 配置文件完善

**文件位置：** `backend/config/`

**任务清单：**
- [ ] 创建 `config.yaml.example` 配置模板
- [ ] 完善 `config.go` 配置加载逻辑
- [ ] 添加环境变量支持
- [ ] 添加配置验证逻辑

**依赖：** 1.1, 1.2, 1.3 完成

**输出：** 完整的数据库表结构、基础设施服务可用、配置文件完善

---

### 阶段 2: 核心业务模块 🚀

**优先级：高 | 预计时间：15-20天**

#### 2.1 认证授权模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `controller/v1/auth.go` - 认证控制器（登录、注册、登出、刷新Token）
- [ ] `service/auth.go` - 认证服务（JWT生成、验证、刷新）
- [ ] `dao/auth.go` - 认证数据访问（Token存储、黑名单）
- [ ] 完善 `pkg/auth/jwt.go` - JWT工具包
- [ ] 实现 OAuth2 集成（可选）
- [ ] 实现 LDAP 集成（可选）

**参考文档：** `docs/api/01_auth.md`

**依赖：** 阶段 1 完成

#### 2.2 用户和工作空间模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `controller/v1/user.go` - 用户控制器（已存在，需完善）
- [ ] `service/user.go` - 用户服务（已存在，需完善）
- [ ] `dao/user.go` - 用户数据访问（已存在，需完善）
- [ ] `controller/v1/workspace.go` - 工作空间控制器
- [ ] `service/workspace.go` - 工作空间服务
- [ ] `dao/workspace.go` - 工作空间数据访问
- [ ] 实现用户配额管理
- [ ] 实现工作空间成员管理

**参考文档：** `docs/api/02_users.md`, `docs/api/03_workspaces.md`

**依赖：** 2.1 完成

#### 2.3 主机管理模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/host.go` - 主机实体
- [ ] `controller/v1/host.go` - 主机控制器
- [ ] `service/host.go` - 主机服务
- [ ] `dao/host.go` - 主机数据访问
- [ ] 实现主机注册和心跳
- [ ] 实现主机资源监控
- [ ] 实现主机健康检查
- [ ] 集成 Prometheus 监控

**参考文档：** `docs/api/04_hosts.md`

**依赖：** 2.2 完成

#### 2.4 环境管理模块（核心）⭐

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/environment.go` - 环境实体
- [ ] `controller/v1/environment.go` - 环境控制器
- [ ] `service/environment.go` - 环境服务
- [ ] `dao/environment.go` - 环境数据访问
- [ ] 实现环境创建（Docker/K8s）
- [ ] 实现环境启动/停止/重启
- [ ] 实现环境删除和清理
- [ ] 实现 SSH 访问配置
- [ ] 实现 RDP 访问配置（Windows）
- [ ] 实现 JupyterLab 集成
- [ ] 实现端口动态分配
- [ ] 实现 GPU 资源调度
- [ ] 集成 JumpServer 堡垒机
- [ ] 集成 Guacamole 远程桌面

**参考文档：** `docs/api/05_environments.md`, `docs/design/ssh_access_design.md`, `docs/design/windows_remote_access.md`

**依赖：** 2.3 完成

**重要性：** 这是系统的核心模块，需要重点投入

#### 2.5 数据集和模型模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/dataset.go` - 数据集实体
- [ ] `model/entity/model.go` - 模型实体
- [ ] `controller/v1/dataset.go` - 数据集控制器
- [ ] `controller/v1/model.go` - 模型控制器
- [ ] `service/dataset.go` - 数据集服务
- [ ] `service/model.go` - 模型服务
- [ ] `dao/dataset.go` - 数据集数据访问
- [ ] `dao/model.go` - 模型数据访问
- [ ] 实现数据集上传（集成 MinIO）
- [ ] 实现数据集版本管理
- [ ] 实现数据集挂载到环境
- [ ] 实现模型上传和版本管理
- [ ] 实现模型下载和部署

**参考文档：** `docs/api/06_datasets.md`, `docs/api/07_models.md`

**依赖：** 2.4 完成

#### 2.6 镜像管理模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/image.go` - 镜像实体
- [ ] `controller/v1/image.go` - 镜像控制器
- [ ] `service/image.go` - 镜像服务
- [ ] `dao/image.go` - 镜像数据访问
- [ ] 实现镜像列表和详情
- [ ] 实现自定义镜像构建
- [ ] 集成 Harbor 镜像仓库
- [ ] 实现镜像推送和拉取

**参考文档：** `docs/api/08_images.md`, `docs/design/storage_and_image_management.md`

**依赖：** 2.5 完成

#### 2.7 训练和推理模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/training_job.go` - 训练任务实体
- [ ] `model/entity/inference_service.go` - 推理服务实体
- [ ] `controller/v1/training.go` - 训练控制器
- [ ] `controller/v1/inference.go` - 推理控制器
- [ ] `service/training.go` - 训练服务
- [ ] `service/inference.go` - 推理服务
- [ ] `dao/training.go` - 训练数据访问
- [ ] `dao/inference.go` - 推理数据访问
- [ ] 实现训练任务创建和调度
- [ ] 实现训练任务监控
- [ ] 实现推理服务部署
- [ ] 实现推理服务管理和扩缩容

**参考文档：** `docs/api/09_training.md`, `docs/api/10_inference.md`

**依赖：** 2.6 完成

**输出：** 核心业务模块全部实现，系统基本功能可用

---

### 阶段 3: 支撑服务模块 🛠️

**优先级：中 | 预计时间：10-15天**

#### 3.1 监控模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/metric.go` - 监控指标实体
- [ ] `controller/v1/monitoring.go` - 监控控制器
- [ ] `service/monitoring.go` - 监控服务
- [ ] `dao/monitoring.go` - 监控数据访问
- [ ] 实现 Prometheus 指标采集
- [ ] 实现 Grafana 仪表板集成
- [ ] 实现实时监控数据查询
- [ ] 实现历史监控数据查询

**参考文档：** `docs/api/11_monitoring.md`

**依赖：** 阶段 2 完成

#### 3.2 计费模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/billing_record.go` - 计费记录实体
- [ ] `controller/v1/billing.go` - 计费控制器
- [ ] `service/billing.go` - 计费服务
- [ ] `dao/billing.go` - 计费数据访问
- [ ] 实现按量计费逻辑
- [ ] 实现包年包月计费逻辑
- [ ] 实现账单生成和查询
- [ ] 实现支付集成（可选）

**参考文档：** `docs/api/12_billing.md`

**依赖：** 3.1 完成

#### 3.3 通知和告警模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/notification.go` - 通知实体
- [ ] `model/entity/alert.go` - 告警实体
- [ ] `controller/v1/notification.go` - 通知控制器
- [ ] `controller/v1/alert.go` - 告警控制器
- [ ] `service/notification.go` - 通知服务
- [ ] `service/alert.go` - 告警服务
- [ ] `dao/notification.go` - 通知数据访问
- [ ] `dao/alert.go` - 告警数据访问
- [ ] 实现通知发送（邮件、短信、站内信）
- [ ] 实现告警规则配置
- [ ] 实现告警触发和通知

**参考文档：** `docs/api/13_notifications.md`, `docs/api/14_alerts.md`

**依赖：** 3.2 完成

#### 3.4 Webhook 和工单模块

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] `model/entity/webhook.go` - Webhook 实体
- [ ] `model/entity/issue.go` - 工单实体
- [ ] `controller/v1/webhook.go` - Webhook 控制器
- [ ] `controller/v1/issue.go` - 工单控制器
- [ ] `service/webhook.go` - Webhook 服务
- [ ] `service/issue.go` - 工单服务
- [ ] `dao/webhook.go` - Webhook 数据访问
- [ ] `dao/issue.go` - 工单数据访问
- [ ] 实现 Webhook 配置和触发
- [ ] 实现工单创建和管理

**参考文档：** `docs/api/15_webhooks.md`, `docs/api/16_issues.md`

**依赖：** 3.3 完成

**输出：** 支撑服务模块全部实现，系统功能完善

---

### 阶段 4: 集成和测试 🧪

**优先级：高 | 预计时间：5-7天**

#### 4.1 前后端接口对接

**文件位置：** `frontend/src/api/`

**任务清单：**
- [ ] 完善前端 API 调用模块
- [ ] 对接所有后端接口
- [ ] 实现请求拦截器（Token、错误处理）
- [ ] 实现响应拦截器（统一错误处理）
- [ ] 测试所有页面功能

**依赖：** 阶段 2、3 完成

#### 4.2 单元测试

**文件位置：** `backend/internal/`

**任务清单：**
- [ ] 为所有 service 层编写单元测试
- [ ] 为所有 dao 层编写单元测试
- [ ] 为关键 pkg 包编写单元测试
- [ ] 使用 mock 进行依赖隔离
- [ ] 确保测试覆盖率 > 70%

**工具：** `go test`, `testify`, `mockery`

**依赖：** 4.1 完成

#### 4.3 集成测试

**文件位置：** `backend/tests/integration/`

**任务清单：**
- [ ] 创建集成测试框架
- [ ] 编写 API 集成测试
- [ ] 编写数据库集成测试
- [ ] 编写 Redis 集成测试
- [ ] 编写第三方服务集成测试

**依赖：** 4.2 完成

#### 4.4 端到端测试

**文件位置：** `frontend/e2e/`

**任务清单：**
- [ ] 配置 E2E 测试框架（Playwright/Cypress）
- [ ] 编写关键用户流程测试
- [ ] 编写环境创建和管理测试
- [ ] 编写数据集和模型管理测试

**依赖：** 4.3 完成

**输出：** 系统经过完整测试，质量有保障

---

### 阶段 5: 部署和文档 🚀

**优先级：中 | 预计时间：3-5天**

#### 5.1 Docker 镜像构建

**文件位置：** `backend/`, `frontend/`

**任务清单：**
- [ ] 创建 `backend/Dockerfile`
- [ ] 创建 `frontend/Dockerfile`
- [ ] 优化镜像大小（多阶段构建）
- [ ] 创建 `.dockerignore` 文件
- [ ] 测试镜像构建和运行

**依赖：** 阶段 4 完成

#### 5.2 Kubernetes 配置

**文件位置：** `k8s/`

**任务清单：**
- [ ] 创建 `k8s/namespace.yaml`
- [ ] 创建 `k8s/configmap.yaml`
- [ ] 创建 `k8s/secrets.yaml`
- [ ] 创建 `k8s/deployments/backend.yaml`
- [ ] 创建 `k8s/deployments/frontend.yaml`
- [ ] 创建 `k8s/services/backend.yaml`
- [ ] 创建 `k8s/services/frontend.yaml`
- [ ] 创建 `k8s/ingress.yaml`
- [ ] 测试 K8s 部署

**依赖：** 5.1 完成

#### 5.3 CI/CD 配置

**文件位置：** `.github/workflows/` 或 `.gitlab-ci.yml`

**任务清单：**
- [ ] 创建 CI 流水线（构建、测试）
- [ ] 创建 CD 流水线（部署）
- [ ] 配置自动化测试
- [ ] 配置代码质量检查
- [ ] 配置安全扫描

**依赖：** 5.2 完成

#### 5.4 部署文档完善

**文件位置：** `docs/`

**任务清单：**
- [ ] 完善 `docs/infrastructure/infrastructure.md`
- [ ] 创建 `docs/deployment/docker_deployment.md`
- [ ] 创建 `docs/deployment/k8s_deployment.md`
- [ ] 创建 `docs/deployment/troubleshooting.md`
- [ ] 更新 `README.md` 部署说明

**依赖：** 5.3 完成

**输出：** 系统可以一键部署，文档完善

---

## 🎯 关键路径分析

### 最短路径（MVP）

如果需要快速交付最小可行产品（MVP），建议按以下顺序实现：

1. **阶段 0** - 修复 main.go（必须）
2. **阶段 1.1-1.2** - 数据库表结构和迁移（必须）
3. **阶段 2.1** - 认证授权模块（必须）
4. **阶段 2.2** - 用户和工作空间模块（必须）
5. **阶段 2.4** - 环境管理模块（核心功能）
6. **阶段 4.1** - 前后端接口对接（必须）

**MVP 预计时间：** 10-15天

**MVP 功能：** 用户可以登录、创建环境、启动/停止环境

### 完整路径

按照阶段 0 → 1 → 2 → 3 → 4 → 5 的顺序完成所有模块。

**完整项目预计时间：** 37-52天（约 2 个月）

---

## 👥 资源分配建议

### 团队配置建议

**最小团队（3人）：**
- 1 名后端工程师（负责阶段 1、2、3）
- 1 名前端工程师（负责阶段 4.1、前端优化）
- 1 名全栈工程师（负责基础设施、测试、部署）

**理想团队（5-6人）：**
- 2 名后端工程师（并行开发核心模块）
- 1 名前端工程师（前后端对接、优化）
- 1 名 DevOps 工程师（基础设施、CI/CD）
- 1 名测试工程师（测试、质量保障）
- 1 名技术负责人（架构、Code Review）

### 工作分配

**后端工程师 A：**
- 阶段 2.1-2.3（认证、用户、主机）
- 阶段 3.1-3.2（监控、计费）

**后端工程师 B：**
- 阶段 2.4-2.7（环境、数据集、镜像、训练）
- 阶段 3.3-3.4（通知、告警、Webhook、工单）

**前端工程师：**
- 阶段 4.1（前后端对接）
- 前端优化和用户体验改进

**DevOps 工程师：**
- 阶段 1.3（基础设施配置）
- 阶段 5（部署和 CI/CD）

**测试工程师：**
- 阶段 4.2-4.4（单元测试、集成测试、E2E 测试）

---

## ⚠️ 风险和注意事项

### 技术风险

1. **环境管理模块复杂度高**
   - 涉及 Docker/K8s、SSH/RDP、GPU 调度等多个技术栈
   - 建议：预留充足时间，分阶段实现，先支持 Docker，再支持 K8s

2. **第三方服务集成**
   - JumpServer、Guacamole、Harbor 等集成可能遇到兼容性问题
   - 建议：提前进行技术验证，准备备选方案

3. **GPU 资源调度**
   - GPU 调度算法需要考虑多种场景和边界情况
   - 建议：参考成熟方案（如 NVIDIA GPU Operator），逐步优化

### 进度风险

1. **后端实现严重滞后**
   - 当前只有 11 个 Go 文件，需要实现 100+ 个文件
   - 建议：优先实现 MVP，后续迭代完善

2. **依赖关系复杂**
   - 模块间存在较多依赖，可能影响并行开发
   - 建议：明确接口定义，使用 mock 进行解耦

### 质量风险

1. **测试覆盖不足**
   - 项目规模大，测试工作量大
   - 建议：从一开始就编写测试，不要等到最后

2. **文档不完善**
   - 代码文档、API 文档、部署文档需要持续维护
   - 建议：代码和文档同步更新

### 关键注意事项

1. **立即修复 backend/main.go 语法错误**（阻塞问题）
2. **优先实现核心功能**（环境管理模块）
3. **及时进行技术验证**（第三方服务集成）
4. **保持代码质量**（Code Review、测试）
5. **定期同步进度**（避免返工）

---

## 🚀 快速开始指南

### 第一周工作计划

**Day 1-2: 紧急修复和环境准备**
- [ ] 删除空的 `backend/main.go` 文件
- [ ] 验证 `backend/cmd/main.go` 可以正常编译运行
- [ ] 配置开发环境（Go、Node.js、Docker、PostgreSQL、Redis）
- [ ] 启动基础设施服务（PostgreSQL、Redis）

**Day 3-5: 数据库表结构实现**
- [ ] 创建所有实体模型文件（14个表）
- [ ] 更新数据库迁移代码
- [ ] 测试数据库迁移
- [ ] 创建初始化和种子数据脚本

**Day 6-7: 认证授权模块**
- [ ] 实现认证控制器、服务、DAO
- [ ] 实现 JWT 生成和验证
- [ ] 测试登录、注册、Token 刷新功能

### 推荐开发顺序

**第 1 周：** 阶段 0 + 阶段 1（基础设施）
**第 2-3 周：** 阶段 2.1-2.3（认证、用户、主机）
**第 4-5 周：** 阶段 2.4（环境管理 - 核心）
**第 6 周：** 阶段 2.5-2.7（数据集、镜像、训练）
**第 7-8 周：** 阶段 3（支撑服务）
**第 9 周：** 阶段 4（集成和测试）
**第 10 周：** 阶段 5（部署和文档）

---

## 📝 总结

### 项目规模

- **前端：** 29 个页面组件（已完成）
- **后端：** 需要实现约 100+ 个 Go 文件（当前仅 11 个）
- **API 接口：** 16 个模块，约 80+ 个接口
- **数据库表：** 14+ 个表
- **基础设施：** 11 个服务（PostgreSQL、Redis、Etcd、MinIO、Harbor、Prometheus、Grafana、Nginx、JumpServer、Guacamole、Uptime-kuma）

### 关键里程碑

1. **Week 1:** 基础设施就绪，数据库表结构完成
2. **Week 3:** 用户认证和基础功能可用
3. **Week 5:** 环境管理核心功能完成（MVP）
4. **Week 8:** 所有业务模块完成
5. **Week 10:** 系统测试完成，可以部署

### 下一步行动

**立即执行（今天）：**
1. 删除 `backend/main.go` 空文件
2. 测试 `backend/cmd/main.go` 编译运行
3. 启动 PostgreSQL 和 Redis 服务

**本周完成：**
1. 实现所有数据库表结构
2. 完成认证授权模块
3. 配置基础设施服务

**本月目标：**
1. 完成核心业务模块（阶段 2）
2. 实现 MVP 功能
3. 前后端基本对接

---

## 📚 参考文档

- [API 文档](./docs/api/README.md)
- [系统架构设计](./docs/design/system_architecture.md)
- [数据库设计](./docs/design/database_design.md)
- [基础设施配置](./docs/infrastructure/infrastructure.md)
- [需求文档](./docs/requirements/README.md)

---

**文档创建：** 2026-01-28
**最后更新：** 2026-01-28
**维护者：** RemoteGPU 开发团队

