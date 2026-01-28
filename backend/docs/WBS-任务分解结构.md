# RemoteGPU 项目任务分解结构（WBS）

## 文档说明
本文档采用5级任务分解结构，从Epic到Subtask，确保每个任务都足够细化，可以直接分配和执行。

### 任务层级定义
- **Level 1 - Epic（史诗）**：大的功能模块，通常需要2-4周完成
- **Level 2 - Feature（特性）**：具体功能，通常需要3-5天完成
- **Level 3 - Story（用户故事）**：用户场景，通常需要1-2天完成
- **Level 4 - Task（任务）**：具体开发任务，通常需要2-4小时完成
- **Level 5 - Subtask（子任务）**：最小可执行单元，通常需要30分钟-1小时完成

### 任务状态
- `TODO` - 待开始
- `IN_PROGRESS` - 进行中
- `REVIEW` - 代码审查
- `TESTING` - 测试中
- `DONE` - 已完成
- `BLOCKED` - 被阻塞

---

## Epic 1: 用户和权限管理系统

**优先级**: P0（最高）
**预估工时**: 40小时
**依赖**: 数据库已配置
**负责人**: 待分配

---

### Feature 1.1: 工作空间管理

**优先级**: P0
**预估工时**: 12小时
**依赖**: Customer模型已完成

#### Story 1.1.1: 创建Workspace数据模型

**用户故事**: 作为开发者，我需要创建Workspace数据模型，以便存储工作空间信息

**验收标准**:
- Workspace实体模型创建完成
- WorkspaceDao实现完成
- 单元测试通过
- 代码审查通过

**预估工时**: 4小时

##### Task 1.1.1.1: 创建Workspace实体模型

**描述**: 在entity包中创建Workspace结构体

**输入**: SQL文件 `sql/03_users_and_permissions.sql`
**输出**: `internal/model/entity/workspace.go`

**Subtask 1.1.1.1.1**: 创建文件和包声明
- 操作: 在`internal/model/entity/`目录创建`workspace.go`
- 添加package声明和必要的import
- 预估: 5分钟

**Subtask 1.1.1.1.2**: 定义Workspace结构体
- 操作: 定义结构体，包含所有字段（ID, UUID, OwnerID等）
- 参考SQL表结构确保字段完整
- 预估: 15分钟

**Subtask 1.1.1.1.3**: 添加GORM标签
- 操作: 为每个字段添加正确的GORM标签
- 主键、唯一索引、外键、默认值等
- 预估: 20分钟

**Subtask 1.1.1.1.4**: 实现TableName方法
- 操作: 实现TableName()方法返回"workspaces"
- 预估: 5分钟

**Subtask 1.1.1.1.5**: 添加字段注释
- 操作: 为结构体和字段添加Go文档注释
- 预估: 10分钟

**Subtask 1.1.1.1.6**: 代码格式化和检查
- 操作: 运行`go fmt`和`go vet`
- 预估: 5分钟

**任务总计**: 1小时

---

##### Task 1.1.1.2: 创建WorkspaceDao基础结构

**描述**: 创建WorkspaceDao结构体和构造函数

**输入**: Workspace实体模型
**输出**: `internal/dao/workspace.go`（基础结构）

**Subtask 1.1.1.2.1**: 创建文件和包声明
- 操作: 创建`workspace.go`文件
- 添加package和import
- 预估: 5分钟

**Subtask 1.1.1.2.2**: 定义WorkspaceDao结构体
- 操作: 定义结构体，包含db字段
- 预估: 5分钟

**Subtask 1.1.1.2.3**: 实现构造函数
- 操作: 实现NewWorkspaceDao()
- 从database包获取db连接
- 预估: 10分钟

**Subtask 1.1.1.2.4**: 添加结构体注释
- 操作: 添加Go文档注释
- 预估: 5分钟

**任务总计**: 25分钟

---

##### Task 1.1.1.3: 实现WorkspaceDao的Create方法

**描述**: 实现创建工作空间的方法

**输入**: WorkspaceDao基础结构
**输出**: Create方法实现

**Subtask 1.1.1.3.1**: 定义方法签名
- 操作: 定义Create(workspace *entity.Workspace) error
- 预估: 5分钟

**Subtask 1.1.1.3.2**: 实现创建逻辑
- 操作: 使用db.Create()创建记录
- 处理错误返回
- 预估: 15分钟

**Subtask 1.1.1.3.3**: 添加方法注释
- 操作: 添加方法文档注释
- 预估: 5分钟

**Subtask 1.1.1.3.4**: 单元测试
- 操作: 编写TestCreate测试用例
- 预估: 20分钟

**任务总计**: 45分钟

---

##### Task 1.1.1.4: 实现WorkspaceDao的查询方法

**描述**: 实现GetByID、GetByUUID等查询方法

**输入**: WorkspaceDao基础结构
**输出**: 查询方法实现

**Subtask 1.1.1.4.1**: 实现GetByID方法
- 操作: 使用db.First()查询
- 处理记录不存在的情况
- 预估: 20分钟

**Subtask 1.1.1.4.2**: 实现GetByUUID方法
- 操作: 使用db.Where().First()查询
- 预估: 15分钟

**Subtask 1.1.1.4.3**: 实现GetByOwnerID方法
- 操作: 使用db.Where().Find()查询列表
- 预估: 15分钟

**Subtask 1.1.1.4.4**: 单元测试
- 操作: 编写查询方法的测试用例
- 预估: 30分钟

**任务总计**: 1小时20分钟

---

##### Task 1.1.1.5: 实现WorkspaceDao的更新和删除方法

**描述**: 实现Update和Delete方法

**输入**: WorkspaceDao基础结构
**输出**: 更新删除方法实现

**Subtask 1.1.1.5.1**: 实现Update方法
- 操作: 使用db.Save()更新记录
- 预估: 15分钟

**Subtask 1.1.1.5.2**: 实现Delete方法
- 操作: 使用db.Delete()删除记录
- 预估: 15分钟

**Subtask 1.1.1.5.3**: 单元测试
- 操作: 编写更新删除的测试用例
- 预估: 20分钟

**任务总计**: 50分钟

---

#### Story 1.1.2: 创建WorkspaceMember数据模型

**用户故事**: 作为开发者，我需要创建WorkspaceMember数据模型，以便管理工作空间成员关系

**验收标准**:
- WorkspaceMember实体模型创建完成
- WorkspaceMemberDao实现完成
- 唯一约束正确配置
- 单元测试通过

**预估工时**: 4小时

##### Task 1.1.2.1: 创建WorkspaceMember实体模型

**Subtask 1.1.2.1.1**: 创建文件和结构体
- 预估: 20分钟

**Subtask 1.1.2.1.2**: 添加GORM标签和唯一索引
- 预估: 20分钟

**Subtask 1.1.2.1.3**: 实现TableName方法
- 预估: 5分钟

**任务总计**: 45分钟

---

##### Task 1.1.2.2: 创建WorkspaceMemberDao

**Subtask 1.1.2.2.1**: 创建基础结构和构造函数
- 预估: 20分钟

**Subtask 1.1.2.2.2**: 实现Create方法
- 预估: 30分钟

**Subtask 1.1.2.2.3**: 实现查询方法
- GetByWorkspaceID
- GetByCustomerID
- GetByWorkspaceAndCustomer
- 预估: 1小时

**Subtask 1.1.2.2.4**: 实现Update和Delete方法
- 预估: 30分钟

**Subtask 1.1.2.2.5**: 编写单元测试
- 预估: 1小时

**任务总计**: 3小时15分钟

---

### Feature 1.2: 资源配额管理

**优先级**: P0
**预估工时**: 8小时
**依赖**: Workspace模型完成

#### Story 1.2.1: 创建ResourceQuota数据模型

**预估工时**: 4小时

##### Task 1.2.1.1: 创建ResourceQuota实体模型（1小时）
- Subtask: 创建文件和结构体（20分钟）
- Subtask: 添加GORM标签（20分钟）
- Subtask: 实现TableName方法（5分钟）
- Subtask: 添加注释和格式化（15分钟）

##### Task 1.2.1.2: 创建ResourceQuotaDao（3小时）
- Subtask: 创建基础结构（20分钟）
- Subtask: 实现Create方法（30分钟）
- Subtask: 实现查询方法（GetByCustomerID, GetByWorkspaceID）（1小时）
- Subtask: 实现Update方法（20分钟）
- Subtask: 编写单元测试（50分钟）

#### Story 1.2.2: 实现配额检查逻辑

**预估工时**: 4小时

##### Task 1.2.2.1: 实现CheckQuota方法（2小时）
- Subtask: 定义方法签名和参数（15分钟）
- Subtask: 实现配额查询逻辑（30分钟）
- Subtask: 实现配额比较逻辑（30分钟）
- Subtask: 编写单元测试（45分钟）

##### Task 1.2.2.2: 实现配额统计方法（2小时）
- Subtask: 实现GetUsedResources方法（1小时）
- Subtask: 实现GetAvailableQuota方法（30分钟）
- Subtask: 编写单元测试（30分钟）

---

## Epic 2: 资源管理系统

**优先级**: P0
**预估工时**: 60小时
**依赖**: Epic 1完成

### Feature 2.1: 主机管理

**优先级**: P0
**预估工时**: 16小时

#### Story 2.1.1: 创建Host数据模型

**预估工时**: 5小时

##### Task 2.1.1.1: 创建Host实体模型（1.5小时）
- Subtask: 创建文件和基础结构（20分钟）
- Subtask: 定义所有字段（ID, Name, IPAddress等）（30分钟）
- Subtask: 添加GORM标签（JSONB, 数组类型）（30分钟）
- Subtask: 添加注释和验证（10分钟）

##### Task 2.1.1.2: 创建HostDao基础方法（3.5小时）
- Subtask: 创建基础结构和构造函数（20分钟）
- Subtask: 实现Create方法（30分钟）
- Subtask: 实现GetByID方法（20分钟）
- Subtask: 实现Update方法（20分钟）
- Subtask: 实现Delete方法（20分钟）
- Subtask: 编写基础CRUD测试（1小时）

#### Story 2.1.2: 实现主机查询功能

**预估工时**: 4小时

##### Task 2.1.2.1: 实现状态查询方法（2小时）
- Subtask: 实现GetByStatus方法（30分钟）
- Subtask: 实现GetByHealthStatus方法（30分钟）
- Subtask: 实现GetByDeploymentMode方法（30分钟）
- Subtask: 编写测试用例（30分钟）

##### Task 2.1.2.2: 实现资源查询方法（2小时）
- Subtask: 实现GetAvailableHosts方法（1小时）
- Subtask: 实现资源排序逻辑（30分钟）
- Subtask: 编写测试用例（30分钟）

#### Story 2.1.3: 实现主机资源管理

**预估工时**: 4小时

##### Task 2.1.3.1: 实现资源更新方法（2小时）
- Subtask: 实现UpdateUsedResources方法（1小时）
- Subtask: 实现资源增减逻辑（30分钟）
- Subtask: 编写测试用例（30分钟）

##### Task 2.1.3.2: 实现心跳管理（2小时）
- Subtask: 实现UpdateHeartbeat方法（30分钟）
- Subtask: 实现UpdateStatus方法（30分钟）
- Subtask: 编写测试用例（1小时）

#### Story 2.1.4: 集成测试和文档

**预估工时**: 3小时

##### Task 2.1.4.1: 编写集成测试（2小时）
- Subtask: 测试主机注册流程（30分钟）
- Subtask: 测试资源分配流程（30分钟）
- Subtask: 测试心跳更新流程（30分钟）
- Subtask: 测试异常场景（30分钟）

##### Task 2.1.4.2: 编写API文档（1小时）
- Subtask: 编写方法说明（30分钟）
- Subtask: 编写使用示例（30分钟）

---

### Feature 2.2: GPU管理

**优先级**: P0
**预估工时**: 20小时

#### Story 2.2.1: 创建GPU数据模型

**预估工时**: 4小时

##### Task 2.2.1.1: 创建GPU实体模型（1小时）
- Subtask: 创建文件和结构体（20分钟）
- Subtask: 定义所有字段（20分钟）
- Subtask: 添加GORM标签和索引（15分钟）
- Subtask: 添加注释（5分钟）

##### Task 2.2.1.2: 创建GPUDao基础方法（3小时）
- Subtask: 创建基础结构（20分钟）
- Subtask: 实现CRUD方法（1小时）
- Subtask: 编写单元测试（1小时40分钟）

#### Story 2.2.2: 实现GPU查询功能

**预估工时**: 4小时

##### Task 2.2.2.1: 实现基础查询方法（2小时）
- Subtask: 实现GetByHostID方法（30分钟）
- Subtask: 实现GetByStatus方法（30分钟）
- Subtask: 实现GetByAllocatedTo方法（30分钟）
- Subtask: 编写测试用例（30分钟）

##### Task 2.2.2.2: 实现可用GPU查询（2小时）
- Subtask: 实现GetAvailableGPUs方法（1小时）
- Subtask: 添加行锁支持（30分钟）
- Subtask: 编写测试用例（30分钟）

#### Story 2.2.3: 实现GPU分配功能（核心）

**预估工时**: 8小时

##### Task 2.2.3.1: 实现单GPU分配（3小时）
- Subtask: 实现Allocate方法（1小时）
- Subtask: 实现事务支持（1小时）
- Subtask: 编写测试用例（1小时）

##### Task 2.2.3.2: 实现批量分配（3小时）
- Subtask: 实现BatchAllocate方法（1.5小时）
- Subtask: 实现原子性保证（1小时）
- Subtask: 编写测试用例（30分钟）

##### Task 2.2.3.3: 实现GPU释放（2小时）
- Subtask: 实现Release方法（1小时）
- Subtask: 实现资源回收逻辑（30分钟）
- Subtask: 编写测试用例（30分钟）

#### Story 2.2.4: 并发安全测试

**预估工时**: 4小时

##### Task 2.2.4.1: 编写并发测试（3小时）
- Subtask: 测试并发分配场景（1小时）
- Subtask: 测试资源竞争场景（1小时）
- Subtask: 测试死锁场景（1小时）

##### Task 2.2.4.2: 性能测试（1小时）
- Subtask: 测试分配性能（30分钟）
- Subtask: 测试查询性能（30分钟）

---

## Epic 3: 环境管理系统（核心）

**优先级**: P0
**预估工时**: 80小时
**依赖**: Epic 1, Epic 2完成

### Feature 3.1: Environment数据模型

**优先级**: P0
**预估工时**: 12小时

#### Story 3.1.1: 创建Environment实体模型（4小时）

##### Task 3.1.1.1: 创建基础模型（2小时）
- Subtask: 定义Environment结构体（1小时）
- Subtask: 添加GORM标签（30分钟）
- Subtask: 添加关联关系（30分钟）

##### Task 3.1.1.2: 创建EnvironmentDao（2小时）
- Subtask: 实现基础CRUD（1小时）
- Subtask: 编写单元测试（1小时）

#### Story 3.1.2: 创建PortMapping模型（4小时）

##### Task 3.1.2.1: 创建PortMapping实体（2小时）
- Subtask: 定义结构体（1小时）
- Subtask: 添加唯一约束（1小时）

##### Task 3.1.2.2: 创建PortMappingDao（2小时）
- Subtask: 实现CRUD方法（1小时）
- Subtask: 实现端口分配逻辑（1小时）

#### Story 3.1.3: 实现查询功能（4小时）

##### Task 3.1.3.1: 实现用户查询（2小时）
- Subtask: GetByCustomerID方法（1小时）
- Subtask: GetByWorkspaceID方法（1小时）

##### Task 3.1.3.2: 实现状态查询（2小时）
- Subtask: GetByStatus方法（1小时）
- Subtask: GetByHostID方法（1小时）

### Feature 3.2: 环境生命周期管理（核心）

**优先级**: P0
**预估工时**: 40小时

#### Story 3.2.1: 实现环境创建流程（16小时）

##### Task 3.2.1.1: 实现资源分配逻辑（6小时）
- Subtask: 实现主机选择算法（2小时）
- Subtask: 实现GPU分配逻辑（2小时）
- Subtask: 实现端口分配逻辑（1小时）
- Subtask: 编写单元测试（1小时）

##### Task 3.2.1.2: 实现配额检查（4小时）
- Subtask: 集成ResourceQuota检查（2小时）
- Subtask: 实现资源预留机制（1小时）
- Subtask: 编写测试用例（1小时）

##### Task 3.2.1.3: 实现事务处理（6小时）
- Subtask: 设计事务边界（1小时）
- Subtask: 实现创建事务（2小时）
- Subtask: 实现回滚逻辑（2小时）
- Subtask: 编写测试用例（1小时）

#### Story 3.2.2: 实现环境启动停止（12小时）

##### Task 3.2.2.1: 实现启动逻辑（4小时）
- Subtask: 实现Start方法（2小时）
- Subtask: 实现状态转换（1小时）
- Subtask: 编写测试用例（1小时）

##### Task 3.2.2.2: 实现停止逻辑（4小时）
- Subtask: 实现Stop方法（2小时）
- Subtask: 实现资源释放（1小时）
- Subtask: 编写测试用例（1小时）

##### Task 3.2.2.3: 实现重启逻辑（4小时）
- Subtask: 实现Restart方法（2小时）
- Subtask: 实现状态管理（1小时）
- Subtask: 编写测试用例（1小时）

#### Story 3.2.3: 实现环境删除（12小时）

##### Task 3.2.3.1: 实现删除逻辑（6小时）
- Subtask: 实现Delete方法（2小时）
- Subtask: 实现资源回收（2小时）
- Subtask: 实现级联删除（1小时）
- Subtask: 编写测试用例（1小时）

##### Task 3.2.3.2: 实现强制删除（管理员）（6小时）
- Subtask: 实现ForceDelete方法（2小时）
- Subtask: 实现权限检查（1小时）
- Subtask: 实现清理逻辑（2小时）
- Subtask: 编写测试用例（1小时）

---

## 任务分配建议

### 开发团队配置
- **后端开发**: 2-3人
- **测试工程师**: 1人
- **项目经理**: 1人

### 并行开发策略

**第1周**：
- 开发者A: Epic 1 - Feature 1.1 (工作空间管理)
- 开发者B: Epic 1 - Feature 1.2 (资源配额管理)
- 测试: 准备测试环境

**第2周**：
- 开发者A: Epic 2 - Feature 2.1 (主机管理)
- 开发者B: Epic 2 - Feature 2.2 (GPU管理)
- 测试: Epic 1集成测试

**第3-4周**：
- 开发者A: Epic 3 - Feature 3.1 (Environment模型)
- 开发者B: Epic 3 - Feature 3.2 (生命周期管理)
- 测试: Epic 2集成测试

---

## 质量保证

### 代码审查检查点
- [ ] 代码符合Go规范
- [ ] 所有方法有注释
- [ ] 错误处理完善
- [ ] 单元测试覆盖率>80%
- [ ] 无安全漏洞

### 测试检查点
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 并发测试通过
- [ ] 性能测试达标

---

## 风险管理

### 高风险任务
1. **GPU并发分配** (Epic 2, Story 2.2.3)
   - 风险: 并发冲突导致重复分配
   - 缓解: 使用数据库行锁和事务

2. **环境创建事务** (Epic 3, Story 3.2.1)
   - 风险: 事务失败导致资源泄漏
   - 缓解: 完善回滚逻辑和资源清理

3. **资源配额检查** (Epic 1, Story 1.2.2)
   - 风险: 并发创建导致超配额
   - 缓解: 使用悲观锁

---

**文档版本**: v1.0
**创建时间**: 2026-01-28
**总预估工时**: 约200小时
**预计完成时间**: 8-10周（2-3人团队）
