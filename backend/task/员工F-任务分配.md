# 员工F - 任务分配（整合与基础环境）

## 基本信息
- **负责人**: 员工F
- **主要职责**: 系统整合、基础环境、Service层、Controller层、杂活
- **预估工时**: 30小时
- **优先级**: P0（最高）
- **开始时间**: 第1周
- **依赖**: 部分工作依赖其他所有员工

---

## 角色定位

你是团队的**整合者和支持者**,负责:
1. 搭建和维护基础开发环境
2. 将各个DAO层整合成Service层
3. 开发Controller层和API接口
4. 处理各种杂活和支持工作
5. 协助其他成员解决问题

**特点**:
- 工作内容多样化
- 需要对整个系统有全局理解
- 是连接各个模块的关键角色
- 部分工作可以立即开始,部分需要等待

---

## 任务概述

### 阶段1: 基础环境准备（可立即开始）
- 数据库迁移脚本
- 配置管理
- 开发工具配置
- 文档框架搭建

### 阶段2: Service层开发（等待DAO层）
- WorkspaceService
- ResourceQuotaService
- HostService
- GPUService
- EnvironmentService

### 阶段3: Controller层开发（等待Service层）
- API接口实现
- 路由配置
- 中间件开发
- 请求验证

### 阶段4: 集成和优化
- 集成测试
- 性能优化
- 文档完善
- Bug修复

---

## 详细任务列表

### 阶段1: 基础环境准备（8小时，可立即开始）

#### Task F.1.1: 数据库迁移脚本（2小时）

**目标**: 创建自动化数据库迁移工具

**Subtask清单**:
- [ ] 创建迁移脚本目录结构（30分钟）
  - `scripts/migrate.sh` - 主迁移脚本
  - `scripts/rollback.sh` - 回滚脚本
  - `scripts/seed.sh` - 测试数据脚本

- [ ] 编写迁移脚本（1小时）
  ```bash
  #!/bin/bash
  # 按顺序执行所有SQL文件
  psql -U remotegpu_user -d remotegpu -f sql/01_init.sql
  psql -U remotegpu_user -d remotegpu -f sql/02_customers.sql
  psql -U remotegpu_user -d remotegpu -f sql/03_users_and_permissions.sql
  # ... 其他SQL文件
  ```

- [ ] 编写测试数据脚本（30分钟）
  - 创建测试用户
  - 创建测试工作空间
  - 创建测试主机和GPU

**验收标准**:
- [ ] 一键执行所有数据库迁移
- [ ] 支持回滚功能
- [ ] 可以快速创建测试数据

---

#### Task F.1.2: 配置管理优化（2小时）

**目标**: 优化配置文件管理,支持多环境

**Subtask清单**:
- [ ] 创建多环境配置（1小时）
  - `config/config.dev.yaml` - 开发环境
  - `config/config.test.yaml` - 测试环境
  - `config/config.prod.yaml` - 生产环境

- [ ] 实现配置加载逻辑（30分钟）
  ```go
  // 根据环境变量加载不同配置
  env := os.Getenv("APP_ENV")
  configFile := fmt.Sprintf("config/config.%s.yaml", env)
  ```

- [ ] 添加配置验证（30分钟）
  - 验证必填字段
  - 验证配置格式

**验收标准**:
- [ ] 支持dev/test/prod三种环境
- [ ] 配置加载正确
- [ ] 配置验证完善

---

#### Task F.1.3: 开发工具配置（2小时）

**目标**: 配置开发工具,提升开发效率

**Subtask清单**:
- [ ] 配置Makefile（1小时）
  ```makefile
  .PHONY: build test run migrate

  build:
      go build -o bin/remotegpu cmd/main.go

  test:
      go test ./... -v -cover

  run:
      go run cmd/main.go

  migrate:
      ./scripts/migrate.sh
  ```

- [ ] 配置.editorconfig（15分钟）
  - 统一代码格式

- [ ] 配置pre-commit hooks（45分钟）
  - 自动运行go fmt
  - 自动运行go vet
  - 自动运行测试

**验收标准**:
- [ ] Makefile命令可用
- [ ] 代码格式统一
- [ ] Git提交前自动检查

---

#### Task F.1.4: API文档框架（2小时）

**目标**: 搭建Swagger API文档框架

**Subtask清单**:
- [ ] 安装Swagger工具（30分钟）
  ```bash
  go get -u github.com/swaggo/swag/cmd/swag
  go get -u github.com/swaggo/gin-swagger
  ```

- [ ] 配置Swagger（1小时）
  - 在main.go中添加Swagger注释
  - 配置Swagger路由

- [ ] 生成初始文档（30分钟）
  ```bash
  swag init -g cmd/main.go
  ```

**验收标准**:
- [ ] Swagger文档可访问
- [ ] 文档自动生成
- [ ] 文档格式正确

---

### 阶段2: Service层开发（12小时，等待DAO层完成）

#### Task F.2.1: WorkspaceService（2小时）

**依赖**: 员工A完成WorkspaceDao

**文件路径**: `internal/service/workspace.go`

**Subtask清单**:
- [ ] 创建WorkspaceService结构体（30分钟）
  ```go
  type WorkspaceService struct {
      workspaceDao *dao.WorkspaceDao
      memberDao    *dao.WorkspaceMemberDao
  }
  ```

- [ ] 实现业务方法（1小时）
  - CreateWorkspace - 创建工作空间
  - AddMember - 添加成员
  - RemoveMember - 移除成员
  - UpdateMemberRole - 更新成员角色
  - CheckPermission - 检查权限

- [ ] 编写单元测试（30分钟）

**验收标准**:
- [ ] 业务逻辑正确
- [ ] 权限检查完善
- [ ] 单元测试通过

---

#### Task F.2.2: ResourceQuotaService（2小时）

**依赖**: 员工B完成ResourceQuotaDao

**文件路径**: `internal/service/resource_quota.go`

**Subtask清单**:
- [ ] 创建ResourceQuotaService（30分钟）

- [ ] 实现业务方法（1小时）
  - GetQuota - 获取配额
  - CheckQuota - 检查配额
  - UpdateQuota - 更新配额
  - GetUsage - 获取使用情况

- [ ] 编写单元测试（30分钟）

**验收标准**:
- [ ] 配额检查逻辑正确
- [ ] 单元测试通过

---

#### Task F.2.3: HostService（2小时）

**依赖**: 员工C完成HostDao

**文件路径**: `internal/service/host.go`

**Subtask清单**:
- [ ] 创建HostService（30分钟）

- [ ] 实现业务方法（1小时）
  - RegisterHost - 注册主机
  - UpdateHostStatus - 更新状态
  - GetAvailableHosts - 获取可用主机
  - SelectHost - 选择主机（负载均衡）

- [ ] 编写单元测试（30分钟）

**验收标准**:
- [ ] 主机选择算法正确
- [ ] 单元测试通过

---

#### Task F.2.4: GPUService（3小时）

**依赖**: 员工D完成GPUDao

**文件路径**: `internal/service/gpu.go`

**Subtask清单**:
- [ ] 创建GPUService（30分钟）

- [ ] 实现业务方法（1.5小时）
  - AllocateGPU - 分配GPU
  - BatchAllocateGPU - 批量分配
  - ReleaseGPU - 释放GPU
  - GetAvailableGPUs - 获取可用GPU

- [ ] 实现事务处理（30分钟）
  - 确保分配的原子性

- [ ] 编写单元测试（30分钟）

**验收标准**:
- [ ] GPU分配逻辑正确
- [ ] 事务处理完善
- [ ] 并发安全
- [ ] 单元测试通过

---

#### Task F.2.5: EnvironmentService（3小时）

**依赖**: 员工E完成EnvironmentDao

**文件路径**: `internal/service/environment.go`

**Subtask清单**:
- [ ] 创建EnvironmentService（30分钟）

- [ ] 实现环境创建逻辑（1.5小时）
  - CreateEnvironment - 创建环境
    - 检查配额
    - 选择主机
    - 分配GPU
    - 分配端口
    - 创建环境记录
    - 调用K8s/Docker创建容器

- [ ] 实现环境管理方法（30分钟）
  - StartEnvironment - 启动
  - StopEnvironment - 停止
  - RestartEnvironment - 重启
  - DeleteEnvironment - 删除

- [ ] 编写单元测试（30分钟）

**验收标准**:
- [ ] 环境创建流程完整
- [ ] 事务处理完善
- [ ] 资源回滚正确
- [ ] 单元测试通过

---

### 阶段3: Controller层开发（6小时，等待Service层完成）

#### Task F.3.1: WorkspaceController（1小时）

**文件路径**: `internal/controller/v1/workspace.go`

**Subtask清单**:
- [ ] 实现API接口（40分钟）
  - POST /api/v1/workspaces - 创建工作空间
  - GET /api/v1/workspaces - 获取列表
  - GET /api/v1/workspaces/:id - 获取详情
  - PUT /api/v1/workspaces/:id - 更新
  - DELETE /api/v1/workspaces/:id - 删除
  - POST /api/v1/workspaces/:id/members - 添加成员

- [ ] 添加参数验证（20分钟）

**验收标准**:
- [ ] API接口实现正确
- [ ] 参数验证完善

---

#### Task F.3.2: EnvironmentController（1.5小时）

**文件路径**: `internal/controller/v1/environment.go`

**Subtask清单**:
- [ ] 实现API接口（1小时）
  - POST /api/v1/environments - 创建环境
  - GET /api/v1/environments - 获取列表
  - GET /api/v1/environments/:id - 获取详情
  - PUT /api/v1/environments/:id/start - 启动
  - PUT /api/v1/environments/:id/stop - 停止
  - DELETE /api/v1/environments/:id - 删除

- [ ] 添加参数验证（30分钟）

**验收标准**:
- [ ] API接口实现正确
- [ ] 参数验证完善

---

#### Task F.3.3: 其他Controller（1.5小时）

**Subtask清单**:
- [ ] HostController（30分钟）
  - 管理员主机管理接口

- [ ] GPUController（30分钟）
  - 管理员GPU管理接口

- [ ] ResourceQuotaController（30分钟）
  - 配额查询接口

---

#### Task F.3.4: 路由配置（1小时）

**文件路径**: `internal/router/router.go`

**Subtask清单**:
- [ ] 配置用户路由（30分钟）
  - 工作空间路由
  - 环境路由
  - 配额路由

- [ ] 配置管理员路由（30分钟）
  - 主机管理路由
  - GPU管理路由
  - 用户管理路由

**验收标准**:
- [ ] 路由配置正确
- [ ] 中间件应用正确
- [ ] 权限控制正确

---

#### Task F.3.5: 中间件开发（1小时）

**Subtask清单**:
- [ ] 实现权限检查中间件（30分钟）
  - 检查工作空间权限
  - 检查资源所有权

- [ ] 实现参数验证中间件（30分钟）
  - 统一参数验证
  - 统一错误响应

**验收标准**:
- [ ] 中间件功能正确
- [ ] 错误处理完善

---

### 阶段4: 集成和优化（4小时）

#### Task F.4.1: 集成测试（2小时）

**Subtask清单**:
- [ ] 编写API集成测试（1小时）
  - 测试完整的业务流程
  - 测试错误场景

- [ ] 编写端到端测试（1小时）
  - 测试用户注册到创建环境的完整流程

**验收标准**:
- [ ] 集成测试通过
- [ ] 覆盖主要业务流程

---

#### Task F.4.2: 性能优化（1小时）

**Subtask清单**:
- [ ] 数据库查询优化（30分钟）
  - 添加必要的索引
  - 优化N+1查询

- [ ] API响应优化（30分钟）
  - 添加缓存
  - 优化响应格式

---

#### Task F.4.3: 文档完善（1小时）

**Subtask清单**:
- [ ] 完善API文档（30分钟）
  - 添加Swagger注释
  - 生成完整文档

- [ ] 编写部署文档（30分钟）
  - 部署步骤
  - 配置说明

---

## 杂活任务清单

### 日常支持工作

**每日任务**:
- [ ] 代码审查（每天30分钟）
  - 审查其他成员的PR
  - 提供改进建议

- [ ] 问题支持（按需）
  - 帮助其他成员解决技术问题
  - 解决代码冲突

- [ ] 文档维护（每天15分钟）
  - 更新团队协作文档
  - 记录重要决策

**每周任务**:
- [ ] 代码集成（每周五）
  - 合并各个分支
  - 解决冲突
  - 运行完整测试

- [ ] 环境维护（每周）
  - 更新依赖包
  - 清理测试数据
  - 备份数据库

---

## 协作和依赖

### 依赖其他员工的工作
- **员工A**: WorkspaceDao完成后才能开发WorkspaceService
- **员工B**: ResourceQuotaDao完成后才能开发ResourceQuotaService
- **员工C**: HostDao完成后才能开发HostService
- **员工D**: GPUDao完成后才能开发GPUService
- **员工E**: EnvironmentDao完成后才能开发EnvironmentService

### 其他员工依赖你的工作
- **所有员工**: 依赖你的基础环境配置
- **所有员工**: 依赖你的开发工具配置
- **测试**: 依赖你的集成测试框架

### 协作点
- 与**所有员工**协作：提供技术支持和问题解决
- 与**项目经理**协作：汇报整体进度和风险

---

## 工作时间安排

### 第1-2天（阶段1）
- 搭建基础环境
- 配置开发工具
- 可以与其他员工完全并行

### 第3-4天（阶段2）
- 等待DAO层完成
- 开发Service层
- 部分并行，部分串行

### 第5天（阶段3）
- 开发Controller层
- 配置路由
- 开发中间件

### 第6天（阶段4）
- 集成测试
- 性能优化
- 文档完善

---

## 验收检查清单

### 基础环境
- [ ] 数据库迁移脚本可用
- [ ] 多环境配置正确
- [ ] 开发工具配置完成
- [ ] API文档框架搭建

### Service层
- [ ] 所有Service实现完成
- [ ] 业务逻辑正确
- [ ] 事务处理完善
- [ ] 单元测试通过

### Controller层
- [ ] 所有API接口实现
- [ ] 路由配置正确
- [ ] 参数验证完善
- [ ] 权限控制正确

### 集成和优化
- [ ] 集成测试通过
- [ ] 性能优化完成
- [ ] 文档完善

### Git提交
- [ ] 提交信息清晰
- [ ] 代码已推送

---

## 注意事项

1. **全局视角**: 需要理解整个系统架构
2. **灵活调整**: 根据其他成员进度灵活调整工作顺序
3. **主动支持**: 主动帮助其他成员解决问题
4. **质量把关**: 作为最后的整合者,要确保代码质量
5. **文档维护**: 及时更新文档,记录重要决策

---

## 技能要求

- 熟悉Go语言和Gin框架
- 理解RESTful API设计
- 熟悉数据库和SQL
- 了解Docker和K8s
- 良好的沟通能力
- 问题解决能力

---

## 进度报告

请在完成每个Task后更新`团队协作文档.md`。

---

## 联系方式

遇到问题请联系：
- 所有员工（提供技术支持）
- 项目经理（汇报进度和风险）
