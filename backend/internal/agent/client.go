package agent

import "context"

// Client Agent 客户端接口
type Client interface {
	// StopProcess 停止进程
	StopProcess(ctx context.Context, req *StopProcessRequest) (*Response, error)

	// ResetSSH 重置SSH密钥
	ResetSSH(ctx context.Context, req *ResetSSHRequest) (*Response, error)

	// CleanupMachine 清理机器
	CleanupMachine(ctx context.Context, req *CleanupRequest) (*Response, error)

	// MountDataset 挂载数据集
	MountDataset(ctx context.Context, req *MountDatasetRequest) (*Response, error)

	// GetSystemInfo 获取系统信息
	GetSystemInfo(ctx context.Context, hostID string) (*SystemInfo, error)

	// ExecuteCommand 执行命令
	ExecuteCommand(ctx context.Context, req *ExecuteCommandRequest) (*ExecuteCommandResponse, error)

	// Ping 检查连接
	Ping(ctx context.Context, hostID string) error

	// Close 关闭连接
	Close() error
}
