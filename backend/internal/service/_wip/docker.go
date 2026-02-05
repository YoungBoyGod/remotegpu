package service

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	apperrors "github.com/YoungBoyGod/remotegpu/pkg/errors"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
)

// DockerService Docker 容器管理服务
// 用于在裸金属主机上直接管理 Docker 容器
type DockerService struct {
	client *client.Client
}

// NewDockerService 创建 Docker 服务
func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ErrorDocker, err)
	}

	return &DockerService{
		client: cli,
	}, nil
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	// 基础配置
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Command     []string          `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	WorkingDir  string            `json:"working_dir"`

	// 资源配置
	CPU    int   `json:"cpu"`     // CPU 核心数
	Memory int64 `json:"memory"`  // 内存(字节)
	GPU    int   `json:"gpu"`     // GPU 数量

	// 端口映射
	PortMappings []*PortMapping `json:"port_mappings"`

	// 存储挂载
	Volumes []*VolumeMount `json:"volumes"`

	// 网络配置
	NetworkMode string `json:"network_mode"` // bridge/host/none
	Hostname    string `json:"hostname"`

	// 其他配置
	Privileged  bool     `json:"privileged"`
	CapAdd      []string `json:"cap_add"`
	Devices     []string `json:"devices"`
	RestartPolicy string `json:"restart_policy"` // no/always/on-failure/unless-stopped
}

// PortMapping 端口映射配置
type PortMapping struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"` // tcp/udp
}

// VolumeMount 存储挂载配置
type VolumeMount struct {
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	ReadOnly      bool   `json:"read_only"`
}

// ContainerInfo 容器信息
type ContainerInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Status    string            `json:"status"`
	State     string            `json:"state"`
	CreatedAt string            `json:"created_at"`
	Ports     []*PortMapping    `json:"ports"`
	Volumes   []*VolumeMount    `json:"volumes"`
	Labels    map[string]string `json:"labels"`
}

// CreateContainer 创建容器
func (s *DockerService) CreateContainer(ctx context.Context, config *ContainerConfig) (*ContainerInfo, error) {
	logger.Info("开始创建容器", "name", config.Name, "image", config.Image)

	// 1. 拉取镜像(如果不存在)
	if err := s.PullImage(ctx, config.Image); err != nil {
		return nil, err
	}

	// 2. 配置端口映射
	portBindings, exposedPorts, err := s.buildPortConfig(config.PortMappings)
	if err != nil {
		return nil, apperrors.WrapWithMessage(apperrors.ErrorDocker, "配置端口映射失败", err)
	}

	// 3. 配置存储挂载
	mounts := s.buildMountConfig(config.Volumes)

	// 4. 配置环境变量
	envVars := make([]string, 0, len(config.Env))
	for k, v := range config.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	// 5. 配置资源限制
	resources := container.Resources{
		Memory:   config.Memory,
		NanoCPUs: int64(config.CPU * 1e9), // 转换为纳秒
	}

	// 6. 配置 GPU (如果需要)
	if config.GPU > 0 {
		deviceRequests := s.buildGPUConfig(config.GPU)
		resources.DeviceRequests = deviceRequests
	}

	// 7. 创建容器配置
	containerConfig := &container.Config{
		Image:        config.Image,
		Cmd:          config.Command,
		Env:          envVars,
		WorkingDir:   config.WorkingDir,
		ExposedPorts: exposedPorts,
		Hostname:     config.Hostname,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Mounts:       mounts,
		Resources:    resources,
		NetworkMode:  container.NetworkMode(config.NetworkMode),
		Privileged:   config.Privileged,
		CapAdd:       config.CapAdd,
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyMode(config.RestartPolicy),
		},
	}

	// 8. 创建容器
	resp, err := s.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, config.Name)
	if err != nil {
		logger.Error("创建容器失败", "name", config.Name, "error", err)
		return nil, apperrors.WrapWithMessage(apperrors.ErrorDocker, "创建容器失败", err)
	}

	// 9. 启动容器
	if err := s.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		logger.Error("启动容器失败", "id", resp.ID, "error", err)
		// 清理已创建的容器
		s.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return nil, apperrors.WrapWithMessage(apperrors.ErrorDocker, "启动容器失败", err)
	}

	logger.Info("容器创建并启动成功", "id", resp.ID, "name", config.Name)

	// 10. 返回容器信息
	return &ContainerInfo{
		ID:     resp.ID,
		Name:   config.Name,
		Image:  config.Image,
		Status: "running",
		State:  "running",
		Ports:  config.PortMappings,
		Volumes: config.Volumes,
	}, nil
}

// StartContainer 启动容器
func (s *DockerService) StartContainer(ctx context.Context, containerID string) error {
	logger.Info("启动容器", "id", containerID)

	if err := s.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		logger.Error("启动容器失败", "id", containerID, "error", err)
		return apperrors.WrapWithMessage(apperrors.ErrorDocker, "启动容器失败", err)
	}

	logger.Info("容器启动成功", "id", containerID)
	return nil
}

// StopContainer 停止容器
func (s *DockerService) StopContainer(ctx context.Context, containerID string, timeout int) error {
	logger.Info("停止容器", "id", containerID, "timeout", timeout)

	stopTimeout := timeout
	if err := s.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &stopTimeout}); err != nil {
		logger.Error("停止容器失败", "id", containerID, "error", err)
		return apperrors.WrapWithMessage(apperrors.ErrorDocker, "停止容器失败", err)
	}

	logger.Info("容器停止成功", "id", containerID)
	return nil
}

// RestartContainer 重启容器
func (s *DockerService) RestartContainer(ctx context.Context, containerID string, timeout int) error {
	logger.Info("重启容器", "id", containerID, "timeout", timeout)

	stopTimeout := timeout
	if err := s.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &stopTimeout}); err != nil {
		logger.Error("重启容器失败", "id", containerID, "error", err)
		return apperrors.WrapWithMessage(apperrors.ErrorDocker, "重启容器失败", err)
	}

	logger.Info("容器重启成功", "id", containerID)
	return nil
}

// DeleteContainer 删除容器
func (s *DockerService) DeleteContainer(ctx context.Context, containerID string, force bool) error {
	logger.Info("删除容器", "id", containerID, "force", force)

	if err := s.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force}); err != nil {
		logger.Error("删除容器失败", "id", containerID, "error", err)
		return apperrors.WrapWithMessage(apperrors.ErrorDocker, "删除容器失败", err)
	}

	logger.Info("容器删除成功", "id", containerID)
	return nil
}

// GetContainer 获取容器信息
func (s *DockerService) GetContainer(ctx context.Context, containerID string) (*ContainerInfo, error) {
	// TODO: 实现获取容器信息

	return nil, fmt.Errorf("TODO: 实现获取容器信息")
}

// ListContainers 列出容器
func (s *DockerService) ListContainers(ctx context.Context, all bool) ([]*ContainerInfo, error) {
	// TODO: 实现列出容器
	// all: 是否包含已停止的容器

	return nil, fmt.Errorf("TODO: 实现列出容器")
}

// GetContainerLogs 获取容器日志
func (s *DockerService) GetContainerLogs(ctx context.Context, containerID string, tail int) (string, error) {
	// TODO: 实现获取容器日志
	// tail: 返回最后 N 行日志

	return "", fmt.Errorf("TODO: 实现获取容器日志")
}

// ExecCommand 在容器中执行命令
func (s *DockerService) ExecCommand(ctx context.Context, containerID string, cmd []string) (string, error) {
	// TODO: 实现在容器中执行命令

	return "", fmt.Errorf("TODO: 实现在容器中执行命令")
}

// GetContainerStats 获取容器资源使用统计
func (s *DockerService) GetContainerStats(ctx context.Context, containerID string) (*ContainerStats, error) {
	// TODO: 实现获取容器资源统计

	return nil, fmt.Errorf("TODO: 实现获取容器资源统计")
}

// ContainerStats 容器资源统计
type ContainerStats struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   int64   `json:"memory_usage"`
	MemoryLimit   int64   `json:"memory_limit"`
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     int64   `json:"network_rx"`
	NetworkTx     int64   `json:"network_tx"`
	BlockRead     int64   `json:"block_read"`
	BlockWrite    int64   `json:"block_write"`
}

// PullImage 拉取镜像
func (s *DockerService) PullImage(ctx context.Context, imageName string) error {
	logger.Info("开始拉取镜像", "image", imageName)

	reader, err := s.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		logger.Error("拉取镜像失败", "image", imageName, "error", err)
		return apperrors.WrapWithMessage(apperrors.ErrorDocker, "拉取镜像失败", err)
	}
	defer reader.Close()

	// 读取拉取进度（避免阻塞）
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		logger.Error("读取镜像拉取进度失败", "image", imageName, "error", err)
		return apperrors.WrapWithMessage(apperrors.ErrorDocker, "读取镜像拉取进度失败", err)
	}

	logger.Info("镜像拉取成功", "image", imageName)
	return nil
}

// ListImages 列出镜像
func (s *DockerService) ListImages(ctx context.Context) ([]string, error) {
	// TODO: 实现列出镜像

	return nil, fmt.Errorf("TODO: 实现列出镜像")
}

// DeleteImage 删除镜像
func (s *DockerService) DeleteImage(ctx context.Context, image string, force bool) error {
	// TODO: 实现镜像删除

	return fmt.Errorf("TODO: 实现镜像删除")
}

// ============ 环境部署相关方法 ============

// CreateEnvironmentContainer 为环境创建容器
// 根据环境配置创建 Docker 容器
func (s *DockerService) CreateEnvironmentContainer(ctx context.Context, env *entity.Environment, host *entity.Host, portMappings []*entity.PortMapping) (string, error) {
	// TODO: 实现环境容器创建
	// 1. 构建容器配置
	// 2. 配置端口映射(从 portMappings 读取)
	// 3. 配置存储挂载(根据 env.StorageType)
	// 4. 配置 GPU(根据 env.GPU 和 env.GPUMode)
	// 5. 创建并启动容器
	// 6. 返回容器 ID

	return "", fmt.Errorf("TODO: 实现环境容器创建")
}

// DeleteEnvironmentContainer 删除环境容器
func (s *DockerService) DeleteEnvironmentContainer(ctx context.Context, containerID string) error {
	// TODO: 实现环境容器删除
	// 1. 停止容器
	// 2. 删除容器
	// 3. 清理相关资源

	return fmt.Errorf("TODO: 实现环境容器删除")
}

// buildContainerConfig 构建容器配置
func (s *DockerService) buildContainerConfig(env *entity.Environment, portMappings []*entity.PortMapping) *ContainerConfig {
	config := &ContainerConfig{
		Name:       fmt.Sprintf("env-%s", env.ID),
		Image:      env.Image,
		CPU:        env.CPU,
		Memory:     env.Memory,
		GPU:        env.GPU,
		NetworkMode: "bridge",
		RestartPolicy: "unless-stopped",
	}

	// 配置端口映射
	for _, pm := range portMappings {
		config.PortMappings = append(config.PortMappings, &PortMapping{
			HostPort:      pm.ExternalPort,
			ContainerPort: pm.InternalPort,
			Protocol:      pm.Protocol,
		})
	}

	// TODO: 配置存储挂载
	// TODO: 配置环境变量
	// TODO: 配置 GPU

	return config
}

// configureGPU 配置 GPU 支持
func (s *DockerService) configureGPU(config *ContainerConfig, gpuCount int, gpuMode string) error {
	// TODO: 实现 GPU 配置
	// 1. 检查主机是否支持 GPU
	// 2. 根据 gpuMode 配置 GPU:
	//    - exclusive: 独占模式,分配指定数量的 GPU
	//    - vgpu: vGPU 虚拟化
	//    - mig: MIG 分片
	//    - shared: 时间片共享
	// 3. 配置 NVIDIA Docker runtime
	// 4. 设置 GPU 设备映射

	return fmt.Errorf("TODO: 实现 GPU 配置")
}

// configureStorage 配置存储挂载
func (s *DockerService) configureStorage(config *ContainerConfig, env *entity.Environment) error {
	// TODO: 实现存储配置
	// 1. 根据 env.StorageType 配置不同的存储:
	//    - local: 本地存储
	//    - nfs: NFS 挂载
	//    - ceph: Ceph 挂载
	//    - s3: S3 挂载(使用 s3fs)
	//    - juicefs: JuiceFS 挂载
	// 2. 创建挂载点
	// 3. 配置权限

	return fmt.Errorf("TODO: 实现存储配置")
}

// buildPortConfig 构建端口映射配置
func (s *DockerService) buildPortConfig(portMappings []*PortMapping) (nat.PortMap, nat.PortSet, error) {
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for _, pm := range portMappings {
		if pm == nil {
			continue
		}

		protocol := pm.Protocol
		if protocol == "" {
			protocol = "tcp"
		}

		containerPort := nat.Port(fmt.Sprintf("%d/%s", pm.ContainerPort, protocol))
		exposedPorts[containerPort] = struct{}{}

		portBindings[containerPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", pm.HostPort),
			},
		}
	}

	return portBindings, exposedPorts, nil
}

// buildMountConfig 构建存储挂载配置
func (s *DockerService) buildMountConfig(volumes []*VolumeMount) []mount.Mount {
	mounts := make([]mount.Mount, 0, len(volumes))

	for _, v := range volumes {
		if v == nil {
			continue
		}

		mounts = append(mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   v.HostPath,
			Target:   v.ContainerPath,
			ReadOnly: v.ReadOnly,
		})
	}

	return mounts
}

// buildGPUConfig 构建 GPU 配置
func (s *DockerService) buildGPUConfig(gpuCount int) []container.DeviceRequest {
	return []container.DeviceRequest{
		{
			Driver:       "nvidia",
			Count:        gpuCount,
			Capabilities: [][]string{{"gpu"}},
		},
	}
}
