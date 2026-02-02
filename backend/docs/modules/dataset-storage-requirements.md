# 数据集存储管理需求文档

## 1. 需求概述

### 1.1 背景
RemoteGPU平台需要支持混合云存储架构,允许用户在公有云(阿里云OSS)和私有云(本地MinIO/RustFS)之间灵活存储和访问数据集,同时根据用户类型和节点位置智能选择最优存储路径。

### 1.2 目标
- 支持多S3存储源管理(阿里云OSS、MinIO、RustFS等)
- 实现存储源之间的自动同步和备份
- 根据用户类型和节点位置智能路由存储请求
- 提供统一的数据集访问接口
- 优化存储成本和访问性能

### 1.3 用户角色
- **外部用户**: 通过公网访问,数据优先存储到公有云OSS
- **内部用户**: 通过内网访问,数据直接存储到本地MinIO
- **系统管理员**: 配置存储源、同步策略、路由规则

## 2. 使用场景

### 2.1 场景一: 外部用户上传数据集
**流程**:
1. 外部用户通过Web界面上传数据集
2. 系统检测用户类型为"外部用户"
3. 数据首先上传到阿里云OSS(公有云)
4. 系统触发异步同步任务,将数据下载到本地MinIO
5. 用户创建训练环境时,优先从本地MinIO加载数据(提升性能)
6. 如果本地MinIO数据不可用,回退到阿里云OSS

**关键点**:
- 公有云作为主存储,保证数据可靠性
- 本地MinIO作为缓存,提升访问速度
- 异步同步,不阻塞用户操作

### 2.2 场景二: 内部用户直接上传
**流程**:
1. 内部用户通过内网访问平台
2. 系统检测用户类型为"内部用户"
3. 数据直接上传到本地MinIO(内网速度快)
4. 可选:根据备份策略同步到阿里云OSS(灾备)

**关键点**:
- 内网直连,上传速度快
- 减少公有云流量成本
- 可选的云端备份

### 2.3 场景三: 根据节点位置智能选择存储
**流程**:
1. 用户创建训练环境,选择GPU节点
2. 系统检测节点所在区域(如:北京、上海、深圳)
3. 系统查询该区域最近的存储源
4. 数据自动从最近的存储源加载到节点
5. 如果最近的存储源不可用,自动切换到备用存储源

**关键点**:
- 就近访问,降低网络延迟
- 多存储源容灾
- 自动故障转移

## 3. 功能需求

### 3.1 存储源管理
- **FR-1.1**: 支持配置多个S3存储源(阿里云OSS、MinIO、AWS S3等)
- **FR-1.2**: 每个存储源包含:名称、类型、Endpoint、AccessKey、SecretKey、Bucket、区域
- **FR-1.3**: 支持存储源的启用/禁用
- **FR-1.4**: 支持存储源的健康检查和状态监控

### 3.2 数据集上传
- **FR-2.1**: 根据用户类型自动选择主存储源
- **FR-2.2**: 支持大文件分片上传(>5GB)
- **FR-2.3**: 显示上传进度和速度
- **FR-2.4**: 上传完成后触发同步任务

### 3.3 数据同步
- **FR-3.1**: 支持配置同步策略(实时同步、定时同步、手动同步)
- **FR-3.2**: 支持配置同步方向(单向、双向)
- **FR-3.3**: 同步任务支持断点续传
- **FR-3.4**: 同步失败自动重试(最多3次)
- **FR-3.5**: 同步任务状态监控(进行中、成功、失败)
- **FR-3.6**: 支持增量同步(只同步变更的文件)

### 3.4 智能路由
- **FR-4.1**: 根据用户类型选择存储源(外部用户→OSS,内部用户→MinIO)
- **FR-4.2**: 根据节点区域选择最近的存储源
- **FR-4.3**: 支持配置路由规则(优先级、权重)
- **FR-4.4**: 存储源不可用时自动故障转移
- **FR-4.5**: 记录路由决策日志

### 3.5 数据访问
- **FR-5.1**: 提供统一的数据集访问接口(屏蔽底层存储差异)
- **FR-5.2**: 支持生成预签名URL(临时访问链接)
- **FR-5.3**: 支持数据集版本管理
- **FR-5.4**: 支持数据集共享(用户间、工作空间间)

## 4. 非功能需求

### 4.1 性能要求
- **NFR-1.1**: 上传速度不低于10MB/s(内网环境)
- **NFR-1.2**: 同步任务不影响用户正常操作(异步执行)
- **NFR-1.3**: 路由决策延迟<100ms
- **NFR-1.4**: 支持并发上传(至少10个用户同时上传)

### 4.2 可靠性要求
- **NFR-2.1**: 数据持久性99.999999999%(11个9)
- **NFR-2.2**: 同步失败自动重试,成功率>99%
- **NFR-2.3**: 存储源故障时自动切换,RTO<5分钟

### 4.3 安全要求
- **NFR-3.1**: 数据传输使用HTTPS/TLS加密
- **NFR-3.2**: 存储凭证加密存储
- **NFR-3.3**: 支持数据集访问权限控制
- **NFR-3.4**: 审计日志记录所有数据操作

### 4.4 可扩展性
- **NFR-4.1**: 支持动态添加新的存储源
- **NFR-4.2**: 支持水平扩展同步服务
- **NFR-4.3**: 存储容量支持PB级扩展

## 5. 数据模型设计

### 5.1 存储源(StorageSource)
```go
type StorageSource struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`           // 存储源名称
    Type        string    `json:"type"`           // 类型: oss/minio/s3/rustfs
    Endpoint    string    `json:"endpoint"`       // 访问端点
    AccessKey   string    `json:"-"`              // 访问密钥(加密存储)
    SecretKey   string    `json:"-"`              // 密钥(加密存储)
    Bucket      string    `json:"bucket"`         // 存储桶名称
    Region      string    `json:"region"`         // 区域
    IsPublic    bool      `json:"is_public"`      // 是否公有云
    Priority    int       `json:"priority"`       // 优先级(数字越小优先级越高)
    Status      string    `json:"status"`         // 状态: active/inactive/error
    HealthCheck time.Time `json:"health_check"`   // 最后健康检查时间
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 5.2 数据集(Dataset)
```go
type Dataset struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`           // 数据集名称
    UserID      uint      `json:"user_id"`        // 所属用户
    WorkspaceID *uint     `json:"workspace_id"`   // 所属工作空间(可选)
    Size        int64     `json:"size"`           // 大小(字节)
    FileCount   int       `json:"file_count"`     // 文件数量
    Description string    `json:"description"`    // 描述
    Tags        []string  `json:"tags"`           // 标签
    Version     string    `json:"version"`        // 版本号
    Status      string    `json:"status"`         // 状态: uploading/ready/syncing/error
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 5.3 数据集副本(DatasetReplica)
```go
type DatasetReplica struct {
    ID              uint      `json:"id"`
    DatasetID       uint      `json:"dataset_id"`        // 数据集ID
    StorageSourceID uint      `json:"storage_source_id"` // 存储源ID
    Path            string    `json:"path"`              // 存储路径
    Size            int64     `json:"size"`              // 副本大小
    Status          string    `json:"status"`            // 状态: syncing/ready/error
    IsPrimary       bool      `json:"is_primary"`        // 是否主副本
    LastSyncAt      *time.Time `json:"last_sync_at"`     // 最后同步时间
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`

    // 关联关系
    Dataset       *Dataset       `gorm:"foreignKey:DatasetID" json:"dataset,omitempty"`
    StorageSource *StorageSource `gorm:"foreignKey:StorageSourceID" json:"storage_source,omitempty"`
}
```

### 5.4 同步任务(SyncTask)
```go
type SyncTask struct {
    ID              uint      `json:"id"`
    DatasetID       uint      `json:"dataset_id"`        // 数据集ID
    SourceID        uint      `json:"source_id"`         // 源存储ID
    TargetID        uint      `json:"target_id"`         // 目标存储ID
    Type            string    `json:"type"`              // 类型: full/incremental
    Status          string    `json:"status"`            // 状态: pending/running/success/failed
    Progress        int       `json:"progress"`          // 进度(0-100)
    TransferredSize int64     `json:"transferred_size"`  // 已传输大小
    TotalSize       int64     `json:"total_size"`        // 总大小
    Speed           int64     `json:"speed"`             // 传输速度(字节/秒)
    ErrorMessage    string    `json:"error_message"`     // 错误信息
    RetryCount      int       `json:"retry_count"`       // 重试次数
    StartedAt       *time.Time `json:"started_at"`       // 开始时间
    CompletedAt     *time.Time `json:"completed_at"`     // 完成时间
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`

    // 关联关系
    Dataset *Dataset       `gorm:"foreignKey:DatasetID" json:"dataset,omitempty"`
    Source  *StorageSource `gorm:"foreignKey:SourceID" json:"source,omitempty"`
    Target  *StorageSource `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}
```

## 6. 技术实现方案

### 6.1 架构设计

#### 6.1.1 整体架构
```
┌─────────────────────────────────────────────────────────────┐
│                         用户层                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  外部用户    │  │  内部用户    │  │  管理员      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      API网关层                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  路由规则引擎 (用户类型识别 + 节点位置感知)           │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      服务层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 数据集服务   │  │ 存储源服务   │  │ 同步服务     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      存储适配层                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ OSS适配器    │  │ MinIO适配器  │  │ S3适配器     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      存储层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 阿里云OSS    │  │ 本地MinIO    │  │ RustFS       │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

#### 6.1.2 核心组件

**1. 路由规则引擎**
- 用户类型识别: 根据用户属性判断是外部用户还是内部用户
- 节点位置感知: 根据GPU节点所在区域选择最近的存储源
- 优先级策略: 支持配置多级路由规则
- 故障转移: 主存储源不可用时自动切换到备用存储源

**2. 存储适配器**
- 统一接口: 提供统一的S3操作接口(上传、下载、删除、列表)
- 多协议支持: 支持AWS S3、阿里云OSS、MinIO等不同实现
- 连接池管理: 复用连接,提升性能
- 健康检查: 定期检查存储源可用性

**3. 同步服务**
- 异步任务队列: 使用消息队列(如Redis)管理同步任务
- 断点续传: 支持大文件分片传输和断点续传
- 增量同步: 只同步变更的文件,减少传输量
- 并发控制: 限制同时执行的同步任务数量

### 6.2 关键流程

#### 6.2.1 外部用户上传流程
```
1. 用户上传数据集
   ↓
2. API网关识别用户类型(外部用户)
   ↓
3. 路由引擎选择主存储源(阿里云OSS)
   ↓
4. 通过OSS适配器上传到阿里云OSS
   ↓
5. 创建Dataset记录(status=uploading)
   ↓
6. 上传完成,更新Dataset(status=ready)
   ↓
7. 创建DatasetReplica记录(primary=true, storage=OSS)
   ↓
8. 触发异步同步任务
   ↓
9. 同步服务创建SyncTask(source=OSS, target=MinIO)
   ↓
10. 后台Worker执行同步
   ↓
11. 同步完成,创建DatasetReplica记录(primary=false, storage=MinIO)
```

#### 6.2.2 内部用户上传流程
```
1. 用户上传数据集
   ↓
2. API网关识别用户类型(内部用户)
   ↓
3. 路由引擎选择主存储源(本地MinIO)
   ↓
4. 通过MinIO适配器直接上传到本地MinIO
   ↓
5. 创建Dataset和DatasetReplica记录
   ↓
6. (可选)根据备份策略触发到OSS的同步任务
```

#### 6.2.3 环境加载数据集流程
```
1. 用户创建训练环境,选择数据集
   ↓
2. 系统获取环境所在节点的区域信息
   ↓
3. 查询该数据集的所有副本(DatasetReplica)
   ↓
4. 路由引擎根据节点区域选择最近的存储源
   ↓
5. 生成预签名URL或挂载路径
   ↓
6. 环境启动时自动加载数据集
   ↓
7. 如果主存储源不可用,自动切换到备用存储源
```

### 6.3 技术选型

#### 6.3.1 存储SDK
- **阿里云OSS**: 使用官方Go SDK `github.com/aliyun/aliyun-oss-go-sdk`
- **MinIO**: 使用官方Go SDK `github.com/minio/minio-go/v7`
- **AWS S3**: 使用官方Go SDK `github.com/aws/aws-sdk-go-v2`

#### 6.3.2 任务队列
- **Redis**: 使用`github.com/go-redis/redis/v8`作为消息队列
- **Asynq**: 使用`github.com/hibiken/asynq`作为异步任务处理框架
  - 支持任务重试
  - 支持任务优先级
  - 提供Web UI监控

#### 6.3.3 配置管理
- 存储源配置存储在数据库中,支持动态修改
- 敏感信息(AccessKey/SecretKey)使用AES加密存储
- 路由规则支持热更新,无需重启服务

#### 6.3.4 监控告警
- 使用Prometheus采集指标:
  - 上传/下载速度
  - 同步任务成功率
  - 存储源健康状态
- 使用Grafana展示监控面板
- 关键事件(同步失败、存储源不可用)触发告警

## 7. 模块归属建议

### 7.1 推荐方案: 创建新模块G - 数据集与存储管理

**理由**:
1. **功能独立性**: 数据集存储管理是一个完整的业务域,包含存储源管理、数据集管理、同步服务等多个子模块
2. **复杂度较高**: 涉及多存储源适配、异步同步、智能路由等复杂逻辑,独立模块便于维护
3. **跨模块依赖**: 数据集会被环境管理(模块D)、工作空间(模块B)等多个模块使用,独立出来避免循环依赖
4. **扩展性**: 未来可能支持更多存储类型(HDFS、Ceph等),独立模块便于扩展

**目录结构**:
```
internal/
├── model/
│   └── entity/
│       ├── storage_source.go      # 存储源实体
│       ├── dataset.go              # 数据集实体
│       ├── dataset_replica.go      # 数据集副本实体
│       └── sync_task.go            # 同步任务实体
├── dao/
│   ├── storage_source_dao.go
│   ├── dataset_dao.go
│   ├── dataset_replica_dao.go
│   └── sync_task_dao.go
├── service/
│   ├── storage_source_service.go   # 存储源管理服务
│   ├── dataset_service.go          # 数据集管理服务
│   ├── storage_router_service.go   # 路由规则引擎
│   └── sync_service.go             # 同步服务
├── adapter/
│   ├── storage_adapter.go          # 存储适配器接口
│   ├── oss_adapter.go              # 阿里云OSS适配器
│   ├── minio_adapter.go            # MinIO适配器
│   └── s3_adapter.go               # AWS S3适配器
├── worker/
│   └── sync_worker.go              # 同步任务Worker
└── api/
    └── dataset_handler.go          # 数据集API接口
```

### 7.2 备选方案: 集成到模块F - 公共基础设施

**理由**:
- 存储管理可以视为基础设施的一部分
- 减少模块数量,降低系统复杂度

**缺点**:
- 模块F职责过重,不利于维护
- 数据集业务逻辑与基础设施混合,边界不清晰

### 7.3 实施计划

#### 阶段一: 基础功能(2周)
1. 创建数据模型和数据库表
2. 实现存储适配器(OSS、MinIO)
3. 实现基本的上传/下载功能
4. 实现简单的路由规则(基于用户类型)

#### 阶段二: 同步功能(2周)
1. 集成Asynq任务队列
2. 实现异步同步服务
3. 实现断点续传和重试机制
4. 添加同步任务监控

#### 阶段三: 智能路由(1周)
1. 实现基于节点位置的路由
2. 实现故障转移机制
3. 添加路由决策日志

#### 阶段四: 优化与监控(1周)
1. 性能优化(连接池、并发控制)
2. 添加Prometheus监控指标
3. 完善单元测试和集成测试
4. 编写API文档

## 8. 总结

本需求文档详细描述了RemoteGPU平台的数据集存储管理功能,包括:
- 支持多S3存储源的混合云架构
- 根据用户类型和节点位置的智能路由
- 存储源之间的自动同步和备份
- 完整的数据模型和技术实现方案

建议创建独立的**模块G - 数据集与存储管理**,预计开发周期6周,可以分阶段实施,优先实现核心功能。

