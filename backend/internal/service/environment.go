package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/deployment"
	"github.com/YoungBoyGod/remotegpu/pkg/host"
	"github.com/YoungBoyGod/remotegpu/pkg/k8s"
	"github.com/YoungBoyGod/remotegpu/pkg/network"
	"github.com/YoungBoyGod/remotegpu/pkg/remote"
	"github.com/YoungBoyGod/remotegpu/pkg/security"
	"github.com/YoungBoyGod/remotegpu/pkg/volume"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Redis Key 前缀和过期时间
const (
	// AccessInfoKeyPrefix 访问信息 Redis Key 前缀
	AccessInfoKeyPrefix = "env:access_info:"
	// AccessInfoExpiration 访问信息过期时间(24小时)
	AccessInfoExpiration = 24 * time.Hour
)

// EnvironmentService 环境服务
type EnvironmentService struct {
	envDao            EnvironmentDaoInterface
	portMappingDao    PortMappingDaoInterface
	hostDao           HostDaoInterface
	gpuDao            *dao.GPUDao
	quotaService      ResourceQuotaServiceInterface
	accessInfoService AccessInfoServiceInterface

	// 模块化管理器
	volumeManager     *volume.ConfigManager
	remoteManager     *remote.AccessManager
	networkManager    *network.NetworkManager
	securityManager   *security.SecurityManager
	deploymentManager *deployment.DeploymentManager
	cacheManager      *cache.CacheManager
	hostManager       *host.HostManager

	k8sClient         K8sClientInterface
	db                DBInterface
}

// NewEnvironmentService 创建环境服务实例
func NewEnvironmentService() *EnvironmentService {
	// 初始化 K8s 客户端(可选,某些部署模式不需要)
	k8sClient, err := k8s.GetClient()
	if err != nil {
		fmt.Printf("警告: K8s 客户端初始化失败: %v (某些功能可能不可用)\n", err)
		k8sClient = nil
	}

	// 初始化缓存管理器
	cacheManager := cache.NewCacheManager()

	// 从配置文件读取 Redis 配置
	redisConfig := &cache.CacheConfig{
		Type:     cache.CacheTypeRedis,
		Addr:     fmt.Sprintf("%s:%d", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port),
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
		PoolSize: config.GlobalConfig.Redis.PoolSize,
		Timeout:  time.Duration(config.GlobalConfig.Redis.Timeout) * time.Second,
	}

	// 尝试初始化 Redis,失败时使用 MemoryCache 作为 fallback
	if redisCache, err := cache.NewCache(redisConfig); err == nil {
		cacheManager.Register("default", redisCache)
	} else {
		fmt.Printf("警告: Redis 初始化失败: %v,使用内存缓存作为备用\n", err)
		// 使用 MemoryCache 作为 fallback
		memConfig := &cache.CacheConfig{
			Type: cache.CacheTypeMemory,
		}
		if memCache, err := cache.NewCache(memConfig); err == nil {
			cacheManager.Register("default", memCache)
		} else {
			fmt.Printf("错误: 内存缓存初始化也失败: %v\n", err)
		}
	}

	// 初始化主机管理器
	selector := host.NewLeastUsedSelector()
	checker := host.NewSimpleHealthChecker(5 * time.Second)
	healthMonitor := host.NewHealthMonitor(checker, 30*time.Second)
	hostManager := host.NewHostManager(selector, healthMonitor)

	// 创建服务实例
	service := &EnvironmentService{
		envDao:            dao.NewEnvironmentDao(),
		portMappingDao:    dao.NewPortMappingDao(),
		hostDao:           dao.NewHostDao(),
		gpuDao:            dao.NewGPUDao(),
		quotaService:      NewResourceQuotaService(),
		accessInfoService: NewAccessInfoService(),

		// 模块化管理器
		volumeManager:     volume.NewConfigManager(),
		remoteManager:     remote.NewAccessManager(),
		networkManager:    network.NewNetworkManager(
			network.FirewallTypeIPTables,
			network.DNSProviderCloudflare,
			"remotegpu.com",
		),
		securityManager:   security.NewSecurityManager("jwt-secret-key"),
		deploymentManager: deployment.NewDeploymentManager(),
		cacheManager:      cacheManager,
		hostManager:       hostManager,

		k8sClient: k8sClient,
		db:        NewGormDBWrapper(database.GetDB()),
	}

	// 同步数据库中的主机到 hostManager
	if err := service.syncHostsToManager(); err != nil {
		// 记录错误但不影响服务创建
		fmt.Printf("同步主机到 hostManager 失败: %v\n", err)
	}

	return service
}

// CreateEnvironmentRequest 创建环境请求
type CreateEnvironmentRequest struct {
	UserID           uint              `json:"user_id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Image            string            `json:"image"`
	CPU              int               `json:"cpu"`
	Memory           int64             `json:"memory"`
	GPU              int               `json:"gpu"`
	Storage          *int64            `json:"storage"`
	Command          []string          `json:"command"`
	Args             []string          `json:"args"`
	Env              map[string]string `json:"env"`
	DeploymentMode   string            `json:"deployment_mode"`   // k8s_pod/docker_local/vm
	EnvironmentType  string            `json:"environment_type"`  // ide/terminal/desktop/training
	GPUMode          string            `json:"gpu_mode"`          // exclusive/vgpu/mig/shared
	StorageType      string            `json:"storage_type"`      // local/nfs/ceph/s3/pvc
	UseJumpserver    bool              `json:"use_jumpserver"`    // 是否使用 Jumpserver
	UseGuacamole     bool              `json:"use_guacamole"`     // 是否使用 Guacamole
}

// selectHost 选择最优主机
// 使用 hostManager 的选择策略
func (s *EnvironmentService) selectHost(cpu int, memory int64, gpu int) (*entity.Host, error) {
	// 构建资源需求
	req := &host.ResourceRequirement{
		CPU:    cpu,
		Memory: memory,
		GPU:    gpu,
	}

	// 使用 hostManager 选择主机
	hostInfo, err := s.hostManager.SelectHost(req)
	if err != nil {
		return nil, fmt.Errorf("选择主机失败: %w", err)
	}

	// 从数据库获取完整的主机信息
	hostEntity, err := s.hostDao.GetByID(hostInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("获取主机详情失败: %w", err)
	}

	return hostEntity, nil
}

// allocateGPUs 分配 GPU
// 使用事务和行锁防止并发冲突
func (s *EnvironmentService) allocateGPUs(tx *gorm.DB, hostID string, envID string, count int) ([]*entity.GPU, error) {
	if count == 0 {
		return nil, nil
	}

	// 查询主机上可用的 GPU（使用行锁）
	var gpus []*entity.GPU
	err := tx.Where("host_id = ? AND status = ?", hostID, "available").
		Limit(count).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&gpus).Error

	if err != nil {
		return nil, fmt.Errorf("查询可用 GPU 失败: %w", err)
	}

	if len(gpus) < count {
		return nil, fmt.Errorf("主机 %s 上可用 GPU 不足 (需要:%d, 可用:%d)", hostID, count, len(gpus))
	}

	// 分配 GPU
	now := time.Now()
	for _, gpu := range gpus {
		gpu.Status = "allocated"
		gpu.AllocatedTo = envID
		gpu.AllocatedAt = &now

		if err := tx.Save(gpu).Error; err != nil {
			return nil, fmt.Errorf("分配 GPU %d 失败: %w", gpu.ID, err)
		}
	}

	return gpus, nil
}

// releaseGPUs 释放 GPU
func (s *EnvironmentService) releaseGPUs(tx *gorm.DB, envID string) error {
	// 查询环境分配的所有 GPU
	var gpus []*entity.GPU
	err := tx.Where("allocated_to = ?", envID).Find(&gpus).Error
	if err != nil {
		return fmt.Errorf("查询已分配 GPU 失败: %w", err)
	}

	// 释放 GPU
	for _, gpu := range gpus {
		gpu.Status = "available"
		gpu.AllocatedTo = ""
		gpu.AllocatedAt = nil

		if err := tx.Save(gpu).Error; err != nil {
			return fmt.Errorf("释放 GPU %d 失败: %w", gpu.ID, err)
		}
	}

	return nil
}

// CreateEnvironment 创建环境
func (s *EnvironmentService) CreateEnvironment(req *CreateEnvironmentRequest) (*entity.Environment, error) {
	// 1. 验证输入参数
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 2. 检查配额
	quotaReq := &ResourceRequest{
		CPU:     req.CPU,
		Memory:  req.Memory,
		GPU:     req.GPU,
		Storage: 0,
	}
	if req.Storage != nil {
		quotaReq.Storage = *req.Storage
	}

	ok, err := s.quotaService.CheckQuota(req.UserID, quotaReq)
	if err != nil {
		return nil, fmt.Errorf("配额检查失败: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("资源配额不足")
	}

	// 3. 选择主机
	host, err := s.selectHost(req.CPU, req.Memory, req.GPU)
	if err != nil {
		return nil, fmt.Errorf("选择主机失败: %w", err)
	}

	// 4. 在事务中创建环境
	var env *entity.Environment
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 生成环境 ID
		envID := uuid.New().String()

		// 分配 GPU
		gpus, err := s.allocateGPUs(tx, host.ID, envID, req.GPU)
		if err != nil {
			return err
		}

		// 创建 Environment 记录
		env = &entity.Environment{
			ID:              envID,
			UserID:          req.UserID,
			HostID:          host.ID,
			Name:            req.Name,
			Description:     req.Description,
			Image:           req.Image,
			Status:          "creating",
			CPU:             req.CPU,
			Memory:          req.Memory,
			GPU:             req.GPU,
			Storage:         req.Storage,
			DeploymentMode:  req.DeploymentMode,
			EnvironmentType: req.EnvironmentType,
			GPUMode:         req.GPUMode,
			StorageType:     req.StorageType,
			UseJumpserver:   req.UseJumpserver,
			UseGuacamole:    req.UseGuacamole,
		}

		// 设置默认值
		if env.DeploymentMode == "" {
			env.DeploymentMode = "k8s_pod"
		}
		if env.EnvironmentType == "" {
			env.EnvironmentType = "ide"
		}
		if env.GPUMode == "" {
			env.GPUMode = "exclusive"
		}
		if env.StorageType == "" {
			env.StorageType = "local"
		}
		if env.LifecyclePolicy == "" {
			env.LifecyclePolicy = "persistent"
		}

		if err := tx.Create(env).Error; err != nil {
			return fmt.Errorf("创建环境记录失败: %w", err)
		}

		// 分配端口 - 使用新的 networkManager
		portRequests := s.buildPortRequests(req)
		if len(portRequests) > 0 {
			for _, portReq := range portRequests {
				// 转换 ServiceType
				var serviceType network.ServiceType
				switch portReq.ServiceType {
				case "ssh":
					serviceType = network.ServiceTypeSSH
				case "rdp":
					serviceType = network.ServiceTypeRDP
				case "vnc":
					serviceType = network.ServiceTypeVNC
				case "jupyter":
					serviceType = network.ServiceTypeJupyter
				default:
					serviceType = network.ServiceTypeCustom
				}

				// 使用 networkManager 分配端口
				mapping, err := s.networkManager.AllocateServicePort(
					envID,
					serviceType,
					portReq.Protocol,
					portReq.Description,
				)
				if err != nil {
					return fmt.Errorf("分配端口失败: %w", err)
				}

				// 更新环境的端口字段
				switch portReq.ServiceType {
				case "ssh":
					env.SSHPort = &mapping.ExternalPort
				case "rdp":
					env.RDPPort = &mapping.ExternalPort
				case "jupyter":
					env.JupyterPort = &mapping.ExternalPort
				case "vnc":
					env.VNCPort = &mapping.ExternalPort
				}
			}

			// 保存端口更新
			if err := tx.Save(env).Error; err != nil {
				return fmt.Errorf("更新环境端口失败: %w", err)
			}
		}

		// 根据部署模式创建环境
		switch env.DeploymentMode {
		case "k8s_pod", "k8s_stateful":
			// 注册 K8s 部署配置到 deploymentManager
			k8sConfig := &deployment.K8sDeploymentConfig{
				Namespace: "default",
				Name:      fmt.Sprintf("env-%s", env.ID[:8]),
				Image:     req.Image,
				Command:   req.Command,
				Args:      req.Args,
				Resources: &deployment.ResourceRequirements{
					CPUCores:  req.CPU,
					MemoryGB:  int(req.Memory / 1024 / 1024 / 1024),
					GPUCount:  req.GPU,
					StorageGB: int(*req.Storage / 1024 / 1024 / 1024),
				},
				Env:           req.Env,
				RestartPolicy: "Always",
			}

			if err := s.deploymentManager.Register(envID, k8sConfig); err != nil {
				return fmt.Errorf("注册部署配置失败: %w", err)
			}

			// K8S 部署
			podName, err := s.createK8sPod(env, host, gpus, req)
			if err != nil {
				return fmt.Errorf("创建 K8s Pod 失败: %w", err)
			}
			env.PodName = podName

		case "docker_local":
			// 注册 Docker 部署配置到 deploymentManager
			dockerConfig := &deployment.DockerDeploymentConfig{
				ContainerName: fmt.Sprintf("env-%s", env.ID[:8]),
				Image:         req.Image,
				Command:       req.Command,
				Args:          req.Args,
				Resources: &deployment.ResourceRequirements{
					CPUCores:  req.CPU,
					MemoryGB:  int(req.Memory / 1024 / 1024 / 1024),
					GPUCount:  req.GPU,
					StorageGB: int(*req.Storage / 1024 / 1024 / 1024),
				},
				Env:           req.Env,
				RestartPolicy: "unless-stopped",
				Runtime:       "nvidia",
			}

			if err := s.deploymentManager.Register(envID, dockerConfig); err != nil {
				return fmt.Errorf("注册部署配置失败: %w", err)
			}

			// Docker 部署
			// TODO: 实现 Docker 容器创建
			return fmt.Errorf("TODO: Docker 部署尚未实现")

		default:
			return fmt.Errorf("不支持的部署模式: %s", env.DeploymentMode)
		}

		// 更新环境记录
		if err := tx.Save(env).Error; err != nil {
			return fmt.Errorf("更新环境记录失败: %w", err)
		}

		// 更新主机资源使用量
		if err := s.updateHostResources(tx, host.ID, req.CPU, req.Memory, req.GPU, true); err != nil {
			return fmt.Errorf("更新主机资源失败: %w", err)
		}

		return nil
	})

	if err != nil {
		// 如果事务失败，尝试清理 K8s 资源
		if env != nil && env.PodName != "" {
			_ = s.k8sClient.DeletePod("default", env.PodName)
		}
		return nil, err
	}

	// 更新环境状态为 running
	env.Status = "running"
	now := time.Now()
	env.StartedAt = &now
	if err := s.envDao.Update(env); err != nil {
		return nil, fmt.Errorf("更新环境状态失败: %w", err)
	}

	// 生成并保存连接信息
	if err := s.GenerateAndSaveAccessInfo(env.ID); err != nil {
		// 记录错误但不影响创建流程
		fmt.Printf("生成连接信息失败: %v\n", err)
	}

	return env, nil
}

// validateCreateRequest 验证创建请求
func (s *EnvironmentService) validateCreateRequest(req *CreateEnvironmentRequest) error {
	if req.Name == "" {
		return fmt.Errorf("环境名称不能为空")
	}
	if req.Image == "" {
		return fmt.Errorf("镜像不能为空")
	}
	if req.CPU <= 0 {
		return fmt.Errorf("CPU 必须大于 0")
	}
	if req.Memory <= 0 {
		return fmt.Errorf("内存必须大于 0")
	}
	if req.GPU < 0 {
		return fmt.Errorf("GPU 不能为负数")
	}
	return nil
}

// updateHostResources 更新主机资源使用量
// 同时更新数据库和 hostManager
func (s *EnvironmentService) updateHostResources(tx *gorm.DB, hostID string, cpu int, memory int64, gpu int, add bool) error {
	// 1. 更新数据库
	var host entity.Host
	if err := tx.Where("id = ?", hostID).First(&host).Error; err != nil {
		return err
	}

	if add {
		host.UsedCPU += cpu
		host.UsedMemory += memory
		host.UsedGPU += gpu
	} else {
		host.UsedCPU -= cpu
		host.UsedMemory -= memory
		host.UsedGPU -= gpu
	}

	if err := tx.Save(&host).Error; err != nil {
		return err
	}

	// 2. 同步更新 hostManager
	if err := s.hostManager.UpdateHostResources(hostID, cpu, int64(memory), gpu, add); err != nil {
		// 记录错误但不影响主流程(因为数据库已更新)
		fmt.Printf("更新 hostManager 资源失败: %v\n", err)
	}

	return nil
}

// createK8sPod 创建 K8s Pod
func (s *EnvironmentService) createK8sPod(env *entity.Environment, host *entity.Host, gpus []*entity.GPU, req *CreateEnvironmentRequest) (string, error) {
	// 构建 Pod 名称
	podName := fmt.Sprintf("env-%s", env.ID[:8])

	// 构建 PodConfig
	podConfig := &k8s.PodConfig{
		Name:      podName,
		Namespace: "default",
		Image:     req.Image,
		Command:   req.Command,
		Args:      req.Args,
		CPU:       int64(req.CPU),
		Memory:    req.Memory,
	}

	// 创建 Pod
	pod, err := s.k8sClient.CreatePod(podConfig)
	if err != nil {
		return "", fmt.Errorf("创建 Pod 失败: %w", err)
	}

	return pod.Name, nil
}

// DeleteEnvironment 删除环境
func (s *EnvironmentService) DeleteEnvironment(id string) error {
	// 1. 获取环境信息
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取环境失败: %w", err)
	}

	// 2. 在事务中删除环境
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 根据部署模式删除资源
		switch env.DeploymentMode {
		case "k8s_pod", "k8s_stateful":
			// 删除 K8s 资源
			if env.PodName != "" {
				if err := s.k8sClient.DeletePod("default", env.PodName); err != nil {
					return fmt.Errorf("删除 K8s Pod 失败: %w", err)
				}
			}

		case "docker_local":
			// 删除 Docker 容器
			if env.ContainerID != "" {
				// TODO: 实现 Docker 容器删除
				// if err := s.dockerService.DeleteEnvironmentContainer(context.Background(), env.ContainerID); err != nil {
				// 	return fmt.Errorf("删除 Docker 容器失败: %w", err)
				// }
			}

		default:
			// 未知部署模式,尝试清理可能存在的资源
			if env.PodName != "" {
				_ = s.k8sClient.DeletePod("default", env.PodName)
			}
		}

		// 释放 GPU
		if err := s.releaseGPUs(tx, env.ID); err != nil {
			return fmt.Errorf("释放 GPU 失败: %w", err)
		}

		// 释放网络资源 - 使用 networkManager 清理端口、防火墙规则、DNS 记录
		if err := s.networkManager.CleanupEnvironment(env.ID); err != nil {
			return fmt.Errorf("释放网络资源失败: %w", err)
		}

		// 清理部署信息
		s.deploymentManager.Delete(env.ID)

		// 更新主机资源使用量
		if err := s.updateHostResources(tx, env.HostID, env.CPU, env.Memory, env.GPU, false); err != nil {
			return fmt.Errorf("更新主机资源失败: %w", err)
		}

		// 删除环境记录
		if err := tx.Delete(&entity.Environment{}, "id = ?", env.ID).Error; err != nil {
			return fmt.Errorf("删除环境记录失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 3. 删除 Redis 缓存
	if err := s.DeleteAccessInfoCache(id); err != nil {
		// Redis 删除失败不影响主流程,只记录日志
		fmt.Printf("删除环境 %s 的 Redis 缓存失败: %v\n", id, err)
	}

	return nil
}

// StartEnvironment 启动环境
func (s *EnvironmentService) StartEnvironment(id string) error {
	// 1. 获取环境信息
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取环境失败: %w", err)
	}

	// 2. 检查环境状态
	if env.Status != "stopped" {
		return fmt.Errorf("只能启动已停止的环境 (当前状态: %s)", env.Status)
	}

	// 3. 启动 K8s Pod（这里简化处理，实际可能需要重新创建 Pod）
	// TODO: 实现 Pod 启动逻辑

	// 4. 更新环境状态
	env.Status = "running"
	now := time.Now()
	env.StartedAt = &now

	if err := s.envDao.Update(env); err != nil {
		return fmt.Errorf("更新环境状态失败: %w", err)
	}

	// 5. 生成并保存连接信息
	if err := s.GenerateAndSaveAccessInfo(id); err != nil {
		// 记录错误但不影响启动流程
		fmt.Printf("生成连接信息失败: %v\n", err)
	}

	return nil
}

// StopEnvironment 停止环境
func (s *EnvironmentService) StopEnvironment(id string) error {
	// 1. 获取环境信息
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取环境失败: %w", err)
	}

	// 2. 检查环境状态
	if env.Status != "running" {
		return fmt.Errorf("只能停止运行中的环境 (当前状态: %s)", env.Status)
	}

	// 3. 停止 K8s Pod（这里简化处理）
	// TODO: 实现 Pod 停止逻辑

	// 4. 更新环境状态
	env.Status = "stopped"
	now := time.Now()
	env.StoppedAt = &now

	if err := s.envDao.Update(env); err != nil {
		return fmt.Errorf("更新环境状态失败: %w", err)
	}

	return nil
}

// RestartEnvironment 重启环境
func (s *EnvironmentService) RestartEnvironment(id string) error {
	// 先停止，再启动
	if err := s.StopEnvironment(id); err != nil {
		return fmt.Errorf("停止环境失败: %w", err)
	}

	if err := s.StartEnvironment(id); err != nil {
		return fmt.Errorf("启动环境失败: %w", err)
	}

	return nil
}

// GetEnvironment 获取环境信息
func (s *EnvironmentService) GetEnvironment(id string) (*entity.Environment, error) {
	return s.envDao.GetByID(id)
}


// GenerateAndSaveAccessInfo 生成并保存连接信息
func (s *EnvironmentService) GenerateAndSaveAccessInfo(envID string) error {
	// 获取环境信息
	env, err := s.envDao.GetByID(envID)
	if err != nil {
		return fmt.Errorf("获取环境失败: %w", err)
	}

	// 获取主机信息
	host, err := s.hostDao.GetByID(env.HostID)
	if err != nil {
		return fmt.Errorf("获取主机失败: %w", err)
	}

	// 生成连接信息
	accessInfo, err := s.accessInfoService.GenerateAccessInfo(env, host)
	if err != nil {
		return fmt.Errorf("生成连接信息失败: %w", err)
	}

	// 将连接信息序列化为JSON
	accessInfoJSON, err := json.Marshal(accessInfo)
	if err != nil {
		return fmt.Errorf("序列化连接信息失败: %w", err)
	}

	// 保存到数据库
	db := database.GetDB()
	if err := db.Model(&entity.Environment{}).Where("id = ?", envID).Update("access_info", accessInfoJSON).Error; err != nil {
		return fmt.Errorf("保存连接信息失败: %w", err)
	}

	// 保存到缓存
	cache := s.cacheManager.GetOrDefault("default")
	ctx := context.Background()
	redisKey := AccessInfoKeyPrefix + envID
	if err := cache.Set(ctx, redisKey, accessInfoJSON, AccessInfoExpiration); err != nil {
		// 缓存保存失败不影响主流程,只记录日志
		fmt.Printf("保存连接信息到缓存失败: %v\n", err)
	}

	return nil
}

// ListEnvironments 列出环境
func (s *EnvironmentService) ListEnvironments(userID uint) ([]*entity.Environment, error) {
	return s.envDao.GetByUserID(userID)
}

// GetStatus 获取环境状态
func (s *EnvironmentService) GetStatus(id string) (string, error) {
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("获取环境失败: %w", err)
	}

	// 如果有 K8s Pod，获取 Pod 状态
	if env.PodName != "" && s.k8sClient != nil {
		podStatus, err := s.k8sClient.GetPodStatus("default", env.PodName)
		if err == nil {
			return podStatus, nil
		}
	}

	return env.Status, nil
}

// GetLogs 获取环境日志
func (s *EnvironmentService) GetLogs(id string, tailLines int64) (string, error) {
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("获取环境失败: %w", err)
	}

	if env.PodName == "" {
		return "", fmt.Errorf("环境没有关联的 Pod")
	}

	if s.k8sClient == nil {
		return "", fmt.Errorf("K8s 客户端未初始化")
	}

	opts := &k8s.LogOptions{
		TailLines:  tailLines,
		Timestamps: true,
	}

	return s.k8sClient.GetPodLogs("default", env.PodName, opts)
}

// GetAccessInfo 获取环境访问信息
func (s *EnvironmentService) GetAccessInfo(id string) (map[string]interface{}, error) {
	var accessInfoJSON []byte
	var accessInfo map[string]interface{}

	// 1. 先从缓存获取
	cache := s.cacheManager.GetOrDefault("default")
	ctx := context.Background()
	redisKey := AccessInfoKeyPrefix + id
	val, err := cache.Get(ctx, redisKey)
	if err == nil && val != "" {
		// 缓存命中,直接返回
		accessInfoJSON = []byte(val)
		if err := json.Unmarshal(accessInfoJSON, &accessInfo); err != nil {
			return nil, fmt.Errorf("解析缓存连接信息失败: %w", err)
		}

		// 获取环境基本信息
		env, err := s.envDao.GetByID(id)
		if err != nil {
			return nil, fmt.Errorf("获取环境失败: %w", err)
		}

		return map[string]interface{}{
			"environment_id": env.ID,
			"status":         env.Status,
			"pod_name":       env.PodName,
			"access_info":    accessInfo,
		}, nil
	}

	// 2. 缓存未命中,从数据库获取
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取环境失败: %w", err)
	}

	// 3. 如果数据库中 access_info 为空,生成连接信息
	if env.AccessInfo == nil || len(env.AccessInfo) == 0 {
		if err := s.GenerateAndSaveAccessInfo(id); err != nil {
			return nil, fmt.Errorf("生成连接信息失败: %w", err)
		}
		// 重新获取环境信息
		env, err = s.envDao.GetByID(id)
		if err != nil {
			return nil, fmt.Errorf("获取环境失败: %w", err)
		}
	}

	// 4. 解析 access_info JSON
	if env.AccessInfo != nil {
		accessInfoJSON = env.AccessInfo
		if err := json.Unmarshal(accessInfoJSON, &accessInfo); err != nil {
			return nil, fmt.Errorf("解析连接信息失败: %w", err)
		}

		// 5. 缓存到缓存管理器
		cache := s.cacheManager.GetOrDefault("default")
		ctx := context.Background()
		redisKey := AccessInfoKeyPrefix + id
		if err := cache.Set(ctx, redisKey, accessInfoJSON, AccessInfoExpiration); err != nil {
			// 缓存保存失败不影响主流程,只记录日志
			fmt.Printf("缓存连接信息失败: %v\n", err)
		}
	}

	// 6. 返回结果
	result := map[string]interface{}{
		"environment_id": env.ID,
		"status":         env.Status,
		"pod_name":       env.PodName,
		"access_info":    accessInfo,
	}

	return result, nil
}

// UpdateAccessInfo 更新环境访问信息(当连接信息发生变动时调用)
func (s *EnvironmentService) UpdateAccessInfo(envID string) error {
	// 先删除缓存
	cache := s.cacheManager.GetOrDefault("default")
	ctx := context.Background()
	redisKey := AccessInfoKeyPrefix + envID
	if err := cache.Delete(ctx, redisKey); err != nil {
		fmt.Printf("删除缓存失败: %v\n", err)
	}

	// 重新生成并保存连接信息
	return s.GenerateAndSaveAccessInfo(envID)
}

// DeleteAccessInfoCache 删除环境访问信息缓存
func (s *EnvironmentService) DeleteAccessInfoCache(envID string) error {
	cache := s.cacheManager.GetOrDefault("default")
	ctx := context.Background()
	redisKey := AccessInfoKeyPrefix + envID
	if err := cache.Delete(ctx, redisKey); err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}
	return nil
}

// buildPortRequests 根据环境类型构建端口分配请求
func (s *EnvironmentService) buildPortRequests(req *CreateEnvironmentRequest) []PortRequest {
	var requests []PortRequest

	// 默认环境类型为 IDE
	envType := req.EnvironmentType
	if envType == "" {
		envType = "ide"
	}

	// 根据环境类型分配端口
	switch envType {
	case "ide":
		// IDE 环境: SSH + Jupyter
		requests = append(requests, PortRequest{
			ServiceType:  "ssh",
			InternalPort: 22,
			Protocol:     "tcp",
			Description:  "SSH 访问",
		})
		requests = append(requests, PortRequest{
			ServiceType:  "jupyter",
			InternalPort: 8888,
			Protocol:     "tcp",
			Description:  "Jupyter Notebook",
		})

	case "terminal":
		// Terminal 环境: 仅 SSH
		requests = append(requests, PortRequest{
			ServiceType:  "ssh",
			InternalPort: 22,
			Protocol:     "tcp",
			Description:  "SSH 访问",
		})

	case "desktop":
		// Desktop 环境: SSH + RDP + VNC
		requests = append(requests, PortRequest{
			ServiceType:  "ssh",
			InternalPort: 22,
			Protocol:     "tcp",
			Description:  "SSH 访问",
		})
		requests = append(requests, PortRequest{
			ServiceType:  "rdp",
			InternalPort: 3389,
			Protocol:     "tcp",
			Description:  "RDP 远程桌面",
		})
		requests = append(requests, PortRequest{
			ServiceType:  "vnc",
			InternalPort: 5900,
			Protocol:     "tcp",
			Description:  "VNC 远程桌面",
		})

	case "data_process":
		// Data Process 环境: SSH + Jupyter
		requests = append(requests, PortRequest{
			ServiceType:  "ssh",
			InternalPort: 22,
			Protocol:     "tcp",
			Description:  "SSH 访问",
		})
		requests = append(requests, PortRequest{
			ServiceType:  "jupyter",
			InternalPort: 8888,
			Protocol:     "tcp",
			Description:  "Jupyter Notebook",
		})

	default:
		// 默认: 仅 SSH
		requests = append(requests, PortRequest{
			ServiceType:  "ssh",
			InternalPort: 22,
			Protocol:     "tcp",
			Description:  "SSH 访问",
		})
	}

	return requests
}

// toHostInfo 将 entity.Host 转换为 host.HostInfo
func toHostInfo(h *entity.Host) *host.HostInfo {
	// 解析 Labels JSON
	labels := make(map[string]string)
	if h.Labels != nil && len(h.Labels) > 0 {
		// 尝试解析 JSON 到 map
		var labelsMap map[string]interface{}
		if err := json.Unmarshal(h.Labels, &labelsMap); err == nil {
			for k, v := range labelsMap {
				if str, ok := v.(string); ok {
					labels[k] = str
				}
			}
		}
	}

	return &host.HostInfo{
		ID:          h.ID,
		Name:        h.Name,
		IPAddress:   h.IPAddress,
		Status:      h.Status,
		TotalCPU:    h.TotalCPU,
		UsedCPU:     h.UsedCPU,
		TotalMemory: h.TotalMemory,
		UsedMemory:  h.UsedMemory,
		TotalGPU:    h.TotalGPU,
		UsedGPU:     h.UsedGPU,
		Labels:      labels,
	}
}

// syncHostsToManager 同步数据库中的主机到 hostManager
func (s *EnvironmentService) syncHostsToManager() error {
	// 查询所有主机
	hosts, err := s.hostDao.ListByStatus("active")
	if err != nil {
		return fmt.Errorf("查询主机失败: %w", err)
	}

	// 注册到 hostManager
	for _, h := range hosts {
		hostInfo := toHostInfo(h)
		s.hostManager.RegisterHost(hostInfo)
	}

	return nil
}
