package ops

import (
	"context"
)

type AgentService struct {
	// grpcClient
}

func NewAgentService() *AgentService {
	return &AgentService{}
}

func (s *AgentService) ResetSSH(ctx context.Context, hostID string) error {
	// Call agent gRPC to reset SSH
	return nil
}

func (s *AgentService) MountDataset(ctx context.Context, hostID string, datasetID uint, path string) error {
	// Call agent gRPC to mount
	return nil
}
