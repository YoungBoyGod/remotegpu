# RemoteGPU 后端开发顺序规划

## 项目概述
这是一个完整的GPU云平台系统，包含用户管理、资源管理、环境管理、计费、监控等核心功能。

## 开发原则
1. **依赖优先**：先实现基础功能，再实现依赖它的高级功能
2. **核心优先**：先实现核心业务（环境管理），再实现辅助功能
3. **分层开发**：数据层 → 业务层 → API层 → 集成层
4. **MVP思维**：每个阶段完成后都能提供可用的功能

---

## 阶段1：基础数据层（Foundation Layer）
**目标**：建立完整的数据模型和DAO层，为后续开发打好基础
**预期成果**：所有实体模型和DAO层完成，可以进行数据库CRUD操作

### 1.1 用户和权限模块
- [x] **Task 1.1.1**: Customer模型和DAO（已完成）
- [ ] **Task 1.1.2**: Workspace模型和DAO
  - 创建 `internal/model/entity/workspace.go`
  - 创建 `internal/dao/workspace.go`
  - 实现基础CRUD方法
- [ ] **Task 1.1.3**: WorkspaceMember模型和DAO
  - 创建 `internal/model/entity/workspace_member.go`
  - 创建 `internal/dao/workspace_member.go`
  - 实现成员关系管理方法
- [ ] **Task 1.1.4**: ResourceQuota模型和DAO
  - 创建 `internal/model/entity/resource_quota.go`
  - 创建 `internal/dao/resource_quota.go`
  - 实现配额查询和更新方法

### 1.2 资源管理模块
- [ ] **Task 1.2.1**: Host模型和DAO
  - 创建 `internal/model/entity/host.go`
  - 创建 `internal/dao/host.go`
  - 实现主机注册、查询、状态更新方法
- [ ] **Task 1.2.2**: GPU模型和DAO
  - 创建 `internal/model/entity/gpu.go`
  - 创建 `internal/dao/gpu.go`
  - 实现GPU查询、分配、释放方法

### 1.3 环境管理模块（核心）
- [ ] **Task 1.3.1**: Environment模型和DAO
  - 创建 `internal/model/entity/environment.go`
  - 创建 `internal/dao/environment.go`
  - 实现环境CRUD和状态管理方法
- [ ] **Task 1.3.2**: PortMapping模型和DAO
  - 创建 `internal/model/entity/port_mapping.go`
  - 创建 `internal/dao/port_mapping.go`
  - 实现端口映射管理方法

### 1.4 镜像和数据模块
- [ ] **Task 1.4.1**: Image模型和DAO
  - 创建 `internal/model/entity/image.go`
  - 创建 `internal/dao/image.go`
  - 实现镜像CRUD方法
- [ ] **Task 1.4.2**: Dataset模型和DAO
  - 创建 `internal/model/entity/dataset.go`
  - 创建 `internal/dao/dataset.go`
  - 实现数据集管理方法

### 1.5 计费模块
- [ ] **Task 1.5.1**: BillingRecord模型和DAO
  - 创建 `internal/model/entity/billing_record.go`
  - 创建 `internal/dao/billing_record.go`
  - 实现计费记录查询和统计方法
- [ ] **Task 1.5.2**: Invoice模型和DAO
  - 创建 `internal/model/entity/invoice.go`
  - 创建 `internal/dao/invoice.go`
  - 实现账单生成和查询方法

### 1.6 监控和日志模块
- [ ] **Task 1.6.1**: Notification模型和DAO
  - 创建 `internal/model/entity/notification.go`
  - 创建 `internal/dao/notification.go`
  - 实现通知CRUD方法
- [ ] **Task 1.6.2**: SystemLog模型和DAO
  - 创建 `internal/model/entity/system_log.go`
  - 创建 `internal/dao/system_log.go`
  - 实现日志记录和查询方法

**阶段1完成标志**：所有实体模型和DAO层完成，单元测试通过

---

## 阶段2：核心业务层（Core Business Layer）
**目标**：实现核心业务逻辑，提供完整的业务功能
**预期成果**：所有Service层完成，业务逻辑可以正常运行

### 2.1 用户管理Service
- [x] **Task 2.1.1**: 基础用户Service（已完成：注册、登录、获取信息、更新）
- [ ] **Task 2.1.2**: 完善用户管理功能
  - 实现用户列表查询（分页、筛选）
  - 实现用户状态管理（激活、暂停、删除）
  - 实现用户权限验证

### 2.2 工作空间管理Service
- [ ] **Task 2.2.1**: 工作空间Service
  - 创建 `internal/service/workspace.go`
  - 实现工作空间CRUD
  - 实现成员管理（添加、移除、角色变更）
  - 实现资源配额管理

### 2.3 主机和GPU管理Service
- [ ] **Task 2.3.1**: 主机管理Service
  - 创建 `internal/service/host.go`
  - 实现主机注册和注销
  - 实现主机状态监控和更新
  - 实现主机资源统计
- [ ] **Task 2.3.2**: GPU管理Service
  - 创建 `internal/service/gpu.go`
  - 实现GPU资源查询和筛选
  - 实现GPU分配和释放逻辑
  - 实现GPU健康检查

### 2.4 环境管理Service（核心）
- [ ] **Task 2.4.1**: 环境管理Service
  - 创建 `internal/service/environment.go`
  - 实现环境创建逻辑（资源分配、端口分配）
  - 实现环境生命周期管理（启动、停止、重启、删除）
  - 实现环境状态查询和监控
  - 实现资源回收逻辑

### 2.5 镜像管理Service
- [ ] **Task 2.5.1**: 镜像管理Service
  - 创建 `internal/service/image.go`
  - 实现镜像CRUD
  - 实现镜像版本管理
  - 预留Harbor集成接口

### 2.6 计费Service
- [ ] **Task 2.6.1**: 计费Service
  - 创建 `internal/service/billing.go`
  - 实现计费记录生成逻辑
  - 实现账单生成和统计
  - 实现费用查询和导出

### 2.7 监控和通知Service
- [ ] **Task 2.7.1**: 通知Service
  - 创建 `internal/service/notification.go`
  - 实现通知发送逻辑
  - 实现通知查询和标记已读

**阶段2完成标志**：所有Service层完成，业务逻辑测试通过

---

## 阶段3：API控制层（API Controller Layer）
**目标**：提供RESTful API接口，供前端调用
**预期成果**：所有API端点完成，前端可以调用

### 3.1 用户Controller
- [x] **Task 3.1.1**: 基础用户API（已完成：注册、登录、获取信息、更新）
- [ ] **Task 3.1.2**: 管理员用户API
  - 实现用户列表API `GET /api/v1/admin/users`
  - 实现用户状态管理API `PUT /api/v1/admin/users/:id/status`
  - 实现用户删除API `DELETE /api/v1/admin/users/:id`

### 3.2 工作空间Controller
- [ ] **Task 3.2.1**: 工作空间API
  - 创建 `internal/controller/v1/workspace.go`
  - 实现工作空间CRUD API
  - 实现成员管理API
  - 实现资源配额API

### 3.3 主机和GPU Controller
- [ ] **Task 3.3.1**: 主机管理API
  - 创建 `internal/controller/v1/host.go`
  - 实现主机列表和详情API
  - 实现主机注册和注销API（管理员）
- [ ] **Task 3.3.2**: GPU管理API
  - 创建 `internal/controller/v1/gpu.go`
  - 实现GPU列表和筛选API
  - 实现GPU分配状态查询API

### 3.4 环境Controller（核心）
- [ ] **Task 3.4.1**: 环境管理API
  - 创建 `internal/controller/v1/environment.go`
  - 实现环境创建API `POST /api/v1/environments`
  - 实现环境列表API `GET /api/v1/environments`
  - 实现环境详情API `GET /api/v1/environments/:id`
  - 实现环境操作API（启动、停止、重启、删除）
  - 实现环境访问信息API（SSH、RDP、Jupyter端口）

### 3.5 镜像Controller
- [ ] **Task 3.5.1**: 镜像管理API
  - 创建 `internal/controller/v1/image.go`
  - 实现镜像CRUD API
  - 实现镜像搜索和筛选API

### 3.6 计费Controller
- [ ] **Task 3.6.1**: 计费API
  - 创建 `internal/controller/v1/billing.go`
  - 实现计费记录查询API
  - 实现账单查询API
  - 实现费用统计API

### 3.7 通知Controller
- [ ] **Task 3.7.1**: 通知API
  - 创建 `internal/controller/v1/notification.go`
  - 实现通知列表API
  - 实现通知标记已读API

### 3.8 路由配置
- [ ] **Task 3.8.1**: 更新路由配置
  - 更新 `internal/router/router.go`
  - 添加所有新的API路由
  - 配置权限中间件

**阶段3完成标志**：所有API端点完成，API文档完成，前端可以集成

---

## 阶段4：基础设施集成（Infrastructure Integration）
**目标**：集成外部基础设施服务，实现完整的云平台功能
**预期成果**：K8s、Harbor、Prometheus等服务集成完成

### 4.1 K8s集成（核心）
- [ ] **Task 4.1.1**: K8s客户端封装
  - 创建 `pkg/k8s/client.go`
  - 实现K8s客户端初始化和连接
- [ ] **Task 4.1.2**: Pod管理
  - 创建 `pkg/k8s/pod.go`
  - 实现Pod创建、删除、查询
  - 实现Pod日志查询
- [ ] **Task 4.1.3**: Service和Ingress管理
  - 创建 `pkg/k8s/service.go`
  - 实现Service创建和管理
  - 实现端口映射配置
- [ ] **Task 4.1.4**: 环境Service集成K8s
  - 更新 `internal/service/environment.go`
  - 实现基于K8s的环境创建
  - 实现环境到Pod的映射

### 4.2 Harbor集成
- [ ] **Task 4.2.1**: Harbor客户端封装
  - 创建 `pkg/harbor/client.go`
  - 实现Harbor API调用
- [ ] **Task 4.2.2**: 镜像Service集成Harbor
  - 更新 `internal/service/image.go`
  - 实现镜像推送和拉取
  - 实现镜像仓库管理

### 4.3 Prometheus集成
- [ ] **Task 4.3.1**: Prometheus客户端封装
  - 创建 `pkg/prometheus/client.go`
  - 实现指标查询API
- [ ] **Task 4.3.2**: 监控数据采集
  - 实现环境资源使用监控
  - 实现GPU使用率监控

### 4.4 远程访问集成
- [ ] **Task 4.4.1**: Guacamole集成
  - 创建 `pkg/guacamole/client.go`
  - 实现Web远程访问
- [ ] **Task 4.4.2**: Jumpserver集成（可选）
  - 创建 `pkg/jumpserver/client.go`
  - 实现SSH/RDP访问管理

### 4.5 存储集成
- [ ] **Task 4.5.1**: RustFS集成
  - 创建 `pkg/storage/client.go`
  - 实现数据卷管理
  - 实现数据集挂载

**阶段4完成标志**：所有基础设施服务集成完成，环境可以正常创建和访问

---

## 阶段5：高级功能和优化（Advanced Features）
**目标**：实现高级功能，优化系统性能和用户体验
**预期成果**：系统功能完善，性能优化

### 5.1 自动化功能
- [ ] **Task 5.1.1**: 自动计费
  - 实现定时任务，自动生成计费记录
  - 实现账单自动生成
- [ ] **Task 5.1.2**: 资源自动回收
  - 实现环境超时自动停止
  - 实现资源泄漏检测和回收

### 5.2 监控和告警
- [ ] **Task 5.2.1**: 告警系统
  - 实现告警规则配置
  - 实现告警通知（邮件、Webhook）
- [ ] **Task 5.2.2**: 系统监控Dashboard
  - 实现系统资源监控
  - 实现用户使用统计

### 5.3 性能优化
- [ ] **Task 5.3.1**: 数据库优化
  - 添加必要的索引
  - 优化慢查询
- [ ] **Task 5.3.2**: 缓存优化
  - 实现Redis缓存
  - 优化热点数据查询

### 5.4 文档和测试
- [ ] **Task 5.4.1**: API文档
  - 使用Swagger生成API文档
  - 编写API使用示例
- [ ] **Task 5.4.2**: 单元测试和集成测试
  - 补充Service层单元测试
  - 编写API集成测试

**阶段5完成标志**：系统功能完善，性能达标，文档完整

---

## 开发顺序总结

### 优先级排序
1. **P0（必须）**：阶段1 + 阶段2 + 阶段3（核心功能）
2. **P1（重要）**：阶段4.1（K8s集成）
3. **P2（需要）**：阶段4.2-4.5（其他集成）
4. **P3（优化）**：阶段5（高级功能）

### 建议开发节奏
- **第1-2周**：完成阶段1（基础数据层）
- **第3-4周**：完成阶段2（核心业务层）
- **第5-6周**：完成阶段3（API控制层）
- **第7-8周**：完成阶段4（基础设施集成）
- **第9-10周**：完成阶段5（高级功能和优化）

### 里程碑
- **Milestone 1**：阶段1完成 - 数据层就绪
- **Milestone 2**：阶段2完成 - 业务逻辑就绪
- **Milestone 3**：阶段3完成 - API就绪，可以前后端联调
- **Milestone 4**：阶段4.1完成 - 环境可以在K8s上创建和运行
- **Milestone 5**：阶段4完成 - 所有基础设施集成完成
- **Milestone 6**：阶段5完成 - 系统功能完善，可以上线

---

## 当前进度
- [x] 数据库迁移到PostgreSQL
- [x] Customer模型和DAO
- [x] 基础用户Service（注册、登录、获取信息、更新）
- [x] 基础用户API（注册、登录、获取信息、更新）
- [x] 健康检查系统
- [x] 基础设施配置

**下一步**：开始阶段1.1.2 - 创建Workspace模型和DAO
