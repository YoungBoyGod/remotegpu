# Service 层集成方案

## 一、当前架构分析

### EnvironmentService 当前依赖的服务

```go
type EnvironmentService struct {
    // DAO 层
    envDao            EnvironmentDaoInterface
    portMappingDao    PortMappingDaoInterface
    hostDao           HostDaoInterface
    gpuDao            *dao.GPUDao

    // 业务服务
    quotaService      ResourceQuotaServiceInterface
    accessInfoService AccessInfoServiceInterface

    // 独立功能服务(需要替换)
    portPoolService   *PortPoolService      // → pkg/network
    dockerService     *DockerService        // → pkg/deployment
    nfsService        *NFSService           // → pkg/volume
    vncService        *VNCService           // → pkg/remote
    rdpService        *RDPService           // → pkg/remote
    jumpserverService *JumpserverService    // → pkg/remote
    s3Service         *S3Service            // → pkg/volume
    juicefsService    *JuiceFSService       // → pkg/volume

    // 基础设施
    k8sClient         K8sClientInterface
    db                DBInterface
}
```

## 二、模块映射关系

### 1. 存储相关服务 → pkg/volume
- **NFSService** → volume.NFSVolumeConfig
- **S3Service** → volume.S3VolumeConfig
- **JuiceFSService** → volume.JuiceFSVolumeConfig
- **统一管理** → volume.ConfigManager

### 2. 远程访问服务 → pkg/remote
- **VNCService** → remote.VNCConfig
- **RDPService** → remote.RDPConfig
- **JumpserverService** → remote.SSHConfig (with Jumpserver)
- **统一管理** → remote.AccessManager

### 3. 网络服务 → pkg/network
- **PortPoolService** → network.PortManager
- **防火墙规则** → network.FirewallManager
- **DNS 配置** → network.DNSManager
- **统一管理** → network.NetworkManager

### 4. 部署服务 → pkg/deployment
- **DockerService** → deployment.DockerDeploymentConfig
- **K8s 部署** → deployment.K8sDeploymentConfig
- **统一管理** → deployment.DeploymentManager

### 5. 安全服务 → pkg/security
- **密码生成** → security.PasswordGenerator
- **SSH 密钥** → security.SSHKeyManager
- **Token 管理** → security.TokenManager
- **统一管理** → security.SecurityManager

## 三、重构后的 EnvironmentService 结构

```go
type EnvironmentService struct {
    // DAO 层(保持不变)
    envDao         EnvironmentDaoInterface
    portMappingDao PortMappingDaoInterface
    hostDao        HostDaoInterface
    gpuDao         *dao.GPUDao

    // 业务服务(保持不变)
    quotaService      ResourceQuotaServiceInterface
    accessInfoService AccessInfoServiceInterface

    // 新的模块化管理器
    volumeManager     *volume.ConfigManager
    remoteManager     *remote.AccessManager
    networkManager    *network.NetworkManager
    securityManager   *security.SecurityManager
    deploymentManager *deployment.DeploymentManager

    // 基础设施(保持不变)
    k8sClient K8sClientInterface
    db        DBInterface
}
```

## 四、集成步骤

### 步骤1: 创建新的 EnvironmentService 构造函数

```go
func NewEnvironmentService() *EnvironmentService {
    k8sClient, _ := k8s.GetClient()

    return &EnvironmentService{
        // DAO 层
        envDao:         dao.NewEnvironmentDao(),
        portMappingDao: dao.NewPortMappingDao(),
        hostDao:        dao.NewHostDao(),
        gpuDao:         dao.NewGPUDao(),

        // 业务服务
        quotaService:      NewResourceQuotaService(),
        accessInfoService: NewAccessInfoService(),

        // 新的模块化管理器
        volumeManager:     volume.NewConfigManager(),
        remoteManager:     remote.NewAccessManager(),
        networkManager:    network.NewNetworkManager(
            network.FirewallTypeIPTables,
            network.DNSProviderCloudflare,
            "example.com",
        ),
        securityManager:   security.NewSecurityManager("jwt-secret"),
        deploymentManager: deployment.NewDeploymentManager(),

        // 基础设施
        k8sClient: k8sClient,
        db:        database.GetDB(),
    }
}
```

### 步骤2: 重构 CreateEnvironment 方法

将原有的分散逻辑整合到模块化架构中:

1. **资源配额检查** (保持不变)
2. **主机选择** (保持不变)
3. **安全配置生成** (使用 securityManager)
4. **网络配置** (使用 networkManager)
5. **存储配置** (使用 volumeManager)
6. **远程访问配置** (使用 remoteManager)
7. **部署配置** (使用 deploymentManager)
8. **创建环境记录** (保持不变)

## 五、迁移策略

### 策略1: 渐进式迁移(推荐)

**阶段1: 添加新管理器,保留旧服务**
- 在 EnvironmentService 中同时保留新旧两套服务
- 新功能使用新管理器
- 旧功能继续使用旧服务
- 逐步迁移功能到新管理器

**阶段2: 功能迁移**
- 逐个功能模块迁移
- 每次迁移后进行测试
- 确保功能正常后再迁移下一个

**阶段3: 清理旧代码**
- 所有功能迁移完成后
- 删除旧的服务代码
- 清理未使用的依赖

### 策略2: 一次性替换

**优点**: 快速完成迁移
**缺点**: 风险较大,可能引入大量问题
**不推荐**: 除非有充分的测试覆盖

## 六、兼容性考虑

### 1. API 接口兼容性
- CreateEnvironmentRequest 结构保持不变
- 响应格式保持不变
- 确保前端无需修改

### 2. 数据库兼容性
- 数据库表结构保持不变
- 新增字段使用迁移脚本
- 保证数据向后兼容

### 3. 配置兼容性
- 支持旧的配置格式
- 提供配置迁移工具
- 逐步过渡到新配置格式

## 七、实施计划

### 第一步: 添加新管理器到 EnvironmentService
- 修改 EnvironmentService 结构体
- 更新构造函数
- 保留旧服务,确保现有功能不受影响

### 第二步: 重构 CreateEnvironment 方法
- 使用 securityManager 生成密码和密钥
- 使用 networkManager 分配端口和配置网络
- 使用 volumeManager 配置存储卷
- 使用 remoteManager 配置远程访问
- 使用 deploymentManager 管理部署配置

### 第三步: 重构其他方法
- UpdateEnvironment
- DeleteEnvironment
- StartEnvironment
- StopEnvironment
- GetEnvironmentAccessInfo

### 第四步: 测试验证
- 单元测试
- 集成测试
- 功能测试
- 性能测试

### 第五步: 清理旧代码
- 删除旧的服务文件
- 清理未使用的依赖
- 更新文档

## 八、风险评估

### 高风险项
- 大规模代码重构可能引入 bug
- 现有功能可能受到影响
- 测试覆盖可能不足

### 风险缓解措施
- 采用渐进式迁移策略
- 保留旧代码作为备份
- 充分的测试覆盖
- 分阶段发布

## 九、总结

本次集成将把分散的功能服务整合到统一的模块化架构中,预期收益:
- 代码结构更清晰
- 维护成本降低
- 功能扩展更容易
- 配置更加灵活

建议采用渐进式迁移策略,确保系统稳定性。

