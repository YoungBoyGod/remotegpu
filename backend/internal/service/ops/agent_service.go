package ops

import (
	"context"
	"fmt"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/agent"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"gorm.io/gorm"
)

// AgentService Agent 服务
type AgentService struct {
	client     agent.Client
	config     *config.AgentConfig
	machineDao *dao.MachineDao
}

// NewAgentService 创建 Agent 服务
func NewAgentService(db *gorm.DB, cfg *config.AgentConfig) *AgentService {
	var client agent.Client
	if cfg.Protocol == "grpc" {
		client = agent.NewGRPCClient(cfg)
	} else {
		client = agent.NewHTTPClient(cfg)
	}

	return &AgentService{
		client:     client,
		config:     cfg,
		machineDao: dao.NewMachineDao(db),
	}
}

type SystemInfoSnapshot struct {
	Hostname      string
	OSType        string
	Kernel        string
	CPUCores      int
	MemoryTotalGB int64
	DiskTotalGB   int64
	Collected     bool
}

func (s *AgentService) GetSystemInfo(ctx context.Context, hostID, address string) (*SystemInfoSnapshot, error) {
	if s.config == nil || !s.config.Enabled {
		return nil, fmt.Errorf("agent disabled")
	}
	if address == "" {
		return nil, fmt.Errorf("host address required")
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, address)
	}
	if grpcClient, ok := s.client.(*agent.GRPCClient); ok {
		grpcClient.RegisterHost(hostID, address)
	}

	info, err := s.client.GetSystemInfo(ctx, hostID)
	if err != nil {
		return nil, err
	}

	return &SystemInfoSnapshot{
		Hostname:      info.Hostname,
		OSType:        info.OS,
		Kernel:        info.Kernel,
		CPUCores:      info.CPUCores,
		MemoryTotalGB: bytesToGB(info.MemoryTotal),
		DiskTotalGB:   bytesToGB(info.DiskTotal),
		Collected:     true,
	}, nil
}

func bytesToGB(value uint64) int64 {
	const gb = 1024 * 1024 * 1024
	if value == 0 {
		return 0
	}
	return int64(value / gb)
}

// getHostAddress 获取主机地址
func (s *AgentService) getHostAddress(ctx context.Context, hostID string) (string, error) {
	host, err := s.machineDao.FindByID(ctx, hostID)
	if err != nil {
		return "", err
	}
	return host.IPAddress, nil
}

// ResetSSH 重置SSH密钥
func (s *AgentService) ResetSSH(ctx context.Context, hostID string) error {
	addr, err := s.getHostAddress(ctx, hostID)
	if err != nil {
		return err
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, addr)
	}

	_, err = s.client.ResetSSH(ctx, &agent.ResetSSHRequest{
		HostID: hostID,
	})
	return err
}

// StopProcess 停止进程
func (s *AgentService) StopProcess(ctx context.Context, hostID string, processID int) error {
	if s.config == nil || !s.config.Enabled {
		return fmt.Errorf("agent disabled")
	}
	if hostID == "" {
		return fmt.Errorf("host id required")
	}
	if processID <= 0 {
		return fmt.Errorf("process id required")
	}
	addr, err := s.getHostAddress(ctx, hostID)
	if err != nil {
		return err
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, addr)
	}

	_, err = s.client.StopProcess(ctx, &agent.StopProcessRequest{
		HostID:    hostID,
		ProcessID: processID,
	})
	return err
}

// MountDataset 挂载数据集
func (s *AgentService) MountDataset(ctx context.Context, hostID string, datasetID uint, path string) error {
	addr, err := s.getHostAddress(ctx, hostID)
	if err != nil {
		return err
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, addr)
	}

	_, err = s.client.MountDataset(ctx, &agent.MountDatasetRequest{
		HostID:     hostID,
		DatasetID:  datasetID,
		MountPoint: path,
	})
	return err
}

// CleanupMachine 清理机器（回收时调用）
func (s *AgentService) CleanupMachine(ctx context.Context, hostID string) error {
	addr, err := s.getHostAddress(ctx, hostID)
	if err != nil {
		return err
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, addr)
	}

	_, err = s.client.CleanupMachine(ctx, &agent.CleanupRequest{
		HostID:       hostID,
		CleanupTypes: []string{"process", "data", "ssh", "docker"},
	})
	return err
}

// CheckAgentHealth 检查 Agent 健康状态
// @author Claude
// @description 验证与指定主机上 Agent 的连接是否正常
// @modified 2026-02-05
func (s *AgentService) CheckAgentHealth(ctx context.Context, hostID string) error {
	if s.config == nil || !s.config.Enabled {
		return fmt.Errorf("agent disabled")
	}

	addr, err := s.getHostAddress(ctx, hostID)
	if err != nil {
		return fmt.Errorf("get host address: %w", err)
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, addr)
	}
	if grpcClient, ok := s.client.(*agent.GRPCClient); ok {
		grpcClient.RegisterHost(hostID, addr)
	}

	return s.client.Ping(ctx, hostID)
}

// IsEnabled 检查 Agent 服务是否启用
func (s *AgentService) IsEnabled() bool {
	return s.config != nil && s.config.Enabled
}
