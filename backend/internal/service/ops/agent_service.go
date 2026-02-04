package ops

import (
	"context"

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
func (s *AgentService) StopProcess(ctx context.Context, hostID string, taskID string) error {
	addr, err := s.getHostAddress(ctx, hostID)
	if err != nil {
		return err
	}

	if httpClient, ok := s.client.(*agent.HTTPClient); ok {
		httpClient.RegisterHost(hostID, addr)
	}

	_, err = s.client.StopProcess(ctx, &agent.StopProcessRequest{
		HostID: hostID,
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
