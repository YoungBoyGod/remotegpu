package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient gRPC 客户端实现
type GRPCClient struct {
	config  *config.AgentConfig
	connMap map[string]*grpc.ClientConn
}

// NewGRPCClient 创建 gRPC 客户端
func NewGRPCClient(cfg *config.AgentConfig) *GRPCClient {
	return &GRPCClient{
		config:  cfg,
		connMap: make(map[string]*grpc.ClientConn),
	}
}

// getConn 获取或创建连接
func (c *GRPCClient) getConn(hostID, address string) (*grpc.ClientConn, error) {
	if conn, ok := c.connMap[hostID]; ok {
		return conn, nil
	}

	target := fmt.Sprintf("%s:%d", address, c.config.GRPCPort)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.Timeout)*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", target, err)
	}

	c.connMap[hostID] = conn
	return conn, nil
}

// Close 关闭所有连接
func (c *GRPCClient) Close() error {
	for _, conn := range c.connMap {
		conn.Close()
	}
	c.connMap = make(map[string]*grpc.ClientConn)
	return nil
}

// TODO: 以下方法需要在生成 proto 代码后实现
// StopProcess, ResetSSH, CleanupMachine, MountDataset
// GetSystemInfo, ExecuteCommand, Ping

// StopProcess 停止进程 (占位实现)
func (c *GRPCClient) StopProcess(ctx context.Context, req *StopProcessRequest) (*Response, error) {
	return nil, fmt.Errorf("gRPC client not implemented yet")
}

// ResetSSH 重置SSH (占位实现)
func (c *GRPCClient) ResetSSH(ctx context.Context, req *ResetSSHRequest) (*Response, error) {
	return nil, fmt.Errorf("gRPC client not implemented yet")
}

// CleanupMachine 清理机器 (占位实现)
func (c *GRPCClient) CleanupMachine(ctx context.Context, req *CleanupRequest) (*Response, error) {
	return nil, fmt.Errorf("gRPC client not implemented yet")
}

// MountDataset 挂载数据集 (占位实现)
func (c *GRPCClient) MountDataset(ctx context.Context, req *MountDatasetRequest) (*Response, error) {
	return nil, fmt.Errorf("gRPC client not implemented yet")
}

// GetSystemInfo 获取系统信息 (占位实现)
func (c *GRPCClient) GetSystemInfo(ctx context.Context, hostID string) (*SystemInfo, error) {
	return nil, fmt.Errorf("gRPC client not implemented yet")
}

// ExecuteCommand 执行命令 (占位实现)
func (c *GRPCClient) ExecuteCommand(ctx context.Context, req *ExecuteCommandRequest) (*ExecuteCommandResponse, error) {
	return nil, fmt.Errorf("gRPC client not implemented yet")
}

// Ping 健康检查 (占位实现)
func (c *GRPCClient) Ping(ctx context.Context, hostID string) error {
	return fmt.Errorf("gRPC client not implemented yet")
}
