package service

import (
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/k8s"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EnvironmentService 环境服务
type EnvironmentService struct {
	envDao          EnvironmentDaoInterface
	portMappingDao  PortMappingDaoInterface
	hostDao         HostDaoInterface
	gpuDao          *dao.GPUDao
	quotaService    ResourceQuotaServiceInterface
	k8sClient       K8sClientInterface
	db              DBInterface
}

// NewEnvironmentService 创建环境服务实例
func NewEnvironmentService() *EnvironmentService {
	k8sClient, _ := k8s.GetClient()
	return &EnvironmentService{
		envDao:         dao.NewEnvironmentDao(),
		portMappingDao: dao.NewPortMappingDao(),
		hostDao:        dao.NewHostDao(),
		gpuDao:         dao.NewGPUDao(),
		quotaService:   NewResourceQuotaService(),
		k8sClient:      k8sClient,
	}
}

// CreateEnvironmentRequest 创建环境请求
type CreateEnvironmentRequest struct {
	CustomerID  uint    `json:"customer_id"`
	WorkspaceID *uint   `json:"workspace_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	CPU         int     `json:"cpu"`
	Memory      int64   `json:"memory"`
	GPU         int     `json:"gpu"`
	Storage     *int64  `json:"storage"`
	Command     []string `json:"command"`
	Args        []string `json:"args"`
	Env         map[string]string `json:"env"`
}

// selectHost 选择最优主机
// 算法：查询可用主机，过滤资源不足的主机，按资源使用率排序，返回最优主机
func (s *EnvironmentService) selectHost(cpu int, memory int64, gpu int) (*entity.Host, error) {
	// 1. 查询所有活跃的主机
	hosts, err := s.hostDao.ListByStatus("active")
	if err != nil {
		return nil, fmt.Errorf("查询主机失败: %w", err)
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("没有可用的主机")
	}

	// 2. 过滤资源不足的主机
	var availableHosts []*entity.Host
	for _, host := range hosts {
		availableCPU := host.TotalCPU - host.UsedCPU
		availableMemory := host.TotalMemory - host.UsedMemory
		availableGPU := host.TotalGPU - host.UsedGPU

		if availableCPU >= cpu && availableMemory >= memory && availableGPU >= gpu {
			availableHosts = append(availableHosts, host)
		}
	}

	if len(availableHosts) == 0 {
		return nil, fmt.Errorf("没有满足资源要求的主机 (需要 CPU:%d, Memory:%d, GPU:%d)", cpu, memory, gpu)
	}

	// 3. 按资源使用率排序，选择使用率最低的主机（负载均衡）
	var bestHost *entity.Host
	var lowestUsageRate float64 = 1.0

	for _, host := range availableHosts {
		cpuUsageRate := float64(host.UsedCPU) / float64(host.TotalCPU)
		memoryUsageRate := float64(host.UsedMemory) / float64(host.TotalMemory)

		var gpuUsageRate float64
		if host.TotalGPU > 0 {
			gpuUsageRate = float64(host.UsedGPU) / float64(host.TotalGPU)
		}

		// 综合使用率（CPU 40%, Memory 40%, GPU 20%）
		avgUsageRate := cpuUsageRate*0.4 + memoryUsageRate*0.4 + gpuUsageRate*0.2

		if bestHost == nil || avgUsageRate < lowestUsageRate {
			bestHost = host
			lowestUsageRate = avgUsageRate
		}
	}

	return bestHost, nil
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

	ok, err := s.quotaService.CheckQuota(req.CustomerID, req.WorkspaceID, quotaReq)
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
			ID:          envID,
			CustomerID:  req.CustomerID,
			WorkspaceID: req.WorkspaceID,
			HostID:      host.ID,
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
			Status:      "creating",
			CPU:         req.CPU,
			Memory:      req.Memory,
			GPU:         req.GPU,
			Storage:     req.Storage,
		}

		if err := tx.Create(env).Error; err != nil {
			return fmt.Errorf("创建环境记录失败: %w", err)
		}

		// 创建 K8s Pod
		podName, err := s.createK8sPod(env, host, gpus, req)
		if err != nil {
			return fmt.Errorf("创建 K8s Pod 失败: %w", err)
		}

		// 更新 Pod 名称
		env.PodName = podName
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
func (s *EnvironmentService) updateHostResources(tx *gorm.DB, hostID string, cpu int, memory int64, gpu int, add bool) error {
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

	return tx.Save(&host).Error
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
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除 K8s 资源
		if env.PodName != "" {
			if err := s.k8sClient.DeletePod("default", env.PodName); err != nil {
				return fmt.Errorf("删除 K8s Pod 失败: %w", err)
			}
		}

		// 释放 GPU
		if err := s.releaseGPUs(tx, env.ID); err != nil {
			return fmt.Errorf("释放 GPU 失败: %w", err)
		}

		// 释放端口
		if err := s.portMappingDao.DeleteByEnvironmentID(env.ID); err != nil {
			return fmt.Errorf("释放端口失败: %w", err)
		}

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

// ListEnvironments 列出环境
func (s *EnvironmentService) ListEnvironments(customerID uint, workspaceID *uint) ([]*entity.Environment, error) {
	if workspaceID != nil {
		return s.envDao.GetByWorkspaceID(*workspaceID)
	}
	return s.envDao.GetByCustomerID(customerID)
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
	env, err := s.envDao.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取环境失败: %w", err)
	}

	// 获取端口映射
	portMappings, err := s.portMappingDao.GetByEnvironmentID(id)
	if err != nil {
		return nil, fmt.Errorf("获取端口映射失败: %w", err)
	}

	// 构建访问信息
	accessInfo := map[string]interface{}{
		"environment_id": env.ID,
		"status":         env.Status,
		"pod_name":       env.PodName,
		"ports":          portMappings,
	}

	return accessInfo, nil
}
