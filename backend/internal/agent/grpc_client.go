package agent

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/YoungBoyGod/remotegpu/api/proto/agent"
	"github.com/YoungBoyGod/remotegpu/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient gRPC 客户端实现
type GRPCClient struct {
	config  *config.AgentConfig
	connMap map[string]*grpc.ClientConn
	hostMap map[string]string // hostID -> address
	mu      sync.RWMutex
}

// NewGRPCClient 创建 gRPC 客户端
func NewGRPCClient(cfg *config.AgentConfig) *GRPCClient {
	return &GRPCClient{
		config:  cfg,
		connMap: make(map[string]*grpc.ClientConn),
		hostMap: make(map[string]string),
	}
}

// RegisterHost 注册主机地址
func (c *GRPCClient) RegisterHost(hostID, address string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hostMap[hostID] = address
}

// getClient 获取或创建 gRPC 客户端
func (c *GRPCClient) getClient(hostID string) (pb.AgentServiceClient, error) {
	c.mu.RLock()
	address, ok := c.hostMap[hostID]
	conn, hasConn := c.connMap[hostID]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("host %s not registered", hostID)
	}

	if hasConn {
		return pb.NewAgentServiceClient(conn), nil
	}

	// 创建新连接
	c.mu.Lock()
	defer c.mu.Unlock()

	// 双重检查
	if conn, ok := c.connMap[hostID]; ok {
		return pb.NewAgentServiceClient(conn), nil
	}

	target := fmt.Sprintf("%s:%d", address, c.config.GRPCPort)

	var opts []grpc.DialOption
	if c.config.TLSEnabled {
		creds, err := credentials.NewClientTLSFromFile(c.config.TLSCertFile, "")
		if err != nil {
			return nil, fmt.Errorf("load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("create client %s: %w", target, err)
	}

	c.connMap[hostID] = conn
	return pb.NewAgentServiceClient(conn), nil
}

// Close 关闭所有连接
func (c *GRPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, conn := range c.connMap {
		conn.Close()
	}
	c.connMap = make(map[string]*grpc.ClientConn)
	return nil
}

// StopProcess 停止进程
func (c *GRPCClient) StopProcess(ctx context.Context, req *StopProcessRequest) (*Response, error) {
	client, err := c.getClient(req.HostID)
	if err != nil {
		return nil, err
	}

	pbReq := &pb.StopProcessRequest{
		HostId:    req.HostID,
		ProcessId: int32(req.ProcessID),
		Signal:    req.Signal,
	}

	resp, err := client.StopProcess(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	if !resp.Success || resp.Code != 0 {
		return nil, fmt.Errorf("agent error: code=%d message=%s", resp.Code, resp.Message)
	}

	return &Response{
		Success: resp.Success,
		Code:    int(resp.Code),
		Message: resp.Message,
	}, nil
}

// ResetSSH 重置SSH
func (c *GRPCClient) ResetSSH(ctx context.Context, req *ResetSSHRequest) (*Response, error) {
	client, err := c.getClient(req.HostID)
	if err != nil {
		return nil, err
	}

	pbReq := &pb.ResetSSHRequest{
		HostId:    req.HostID,
		PublicKey: req.PublicKey,
		Username:  req.Username,
	}

	resp, err := client.ResetSSH(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	if !resp.Success || resp.Code != 0 {
		return nil, fmt.Errorf("agent error: code=%d message=%s", resp.Code, resp.Message)
	}

	return &Response{
		Success: resp.Success,
		Code:    int(resp.Code),
		Message: resp.Message,
	}, nil
}

// CleanupMachine 清理机器
func (c *GRPCClient) CleanupMachine(ctx context.Context, req *CleanupRequest) (*Response, error) {
	client, err := c.getClient(req.HostID)
	if err != nil {
		return nil, err
	}

	pbReq := &pb.CleanupRequest{
		HostId:       req.HostID,
		CleanupTypes: req.CleanupTypes,
	}

	resp, err := client.CleanupMachine(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	if !resp.Success || resp.Code != 0 {
		return nil, fmt.Errorf("agent error: code=%d message=%s", resp.Code, resp.Message)
	}

	return &Response{
		Success: resp.Success,
		Code:    int(resp.Code),
		Message: resp.Message,
	}, nil
}

// MountDataset 挂载数据集
func (c *GRPCClient) MountDataset(ctx context.Context, req *MountDatasetRequest) (*Response, error) {
	client, err := c.getClient(req.HostID)
	if err != nil {
		return nil, err
	}

	pbReq := &pb.MountDatasetRequest{
		HostId:     req.HostID,
		DatasetId:  uint32(req.DatasetID),
		SourcePath: req.SourcePath,
		MountPoint: req.MountPoint,
		ReadOnly:   req.ReadOnly,
	}

	resp, err := client.MountDataset(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	if !resp.Success || resp.Code != 0 {
		return nil, fmt.Errorf("agent error: code=%d message=%s", resp.Code, resp.Message)
	}

	return &Response{
		Success: resp.Success,
		Code:    int(resp.Code),
		Message: resp.Message,
	}, nil
}

// GetSystemInfo 获取系统信息
func (c *GRPCClient) GetSystemInfo(ctx context.Context, hostID string) (*SystemInfo, error) {
	client, err := c.getClient(hostID)
	if err != nil {
		return nil, err
	}

	pbReq := &pb.SystemInfoRequest{HostId: hostID}
	resp, err := client.GetSystemInfo(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	info := &SystemInfo{
		Hostname:    resp.Hostname,
		OS:          resp.Os,
		Kernel:      resp.Kernel,
		CPUCores:    int(resp.CpuCores),
		MemoryTotal: resp.MemoryTotal,
		MemoryFree:  resp.MemoryFree,
		DiskTotal:   resp.DiskTotal,
		DiskFree:    resp.DiskFree,
		GPUCount:    int(resp.GpuCount),
		Uptime:      resp.Uptime,
	}

	for _, g := range resp.GpuInfo {
		info.GPUInfo = append(info.GPUInfo, GPUInfo{
			Index:       int(g.Index),
			Name:        g.Name,
			MemoryTotal: g.MemoryTotal,
			MemoryUsed:  g.MemoryUsed,
			Utilization: int(g.Utilization),
			Temperature: int(g.Temperature),
		})
	}

	return info, nil
}

// ExecuteCommand 执行命令
func (c *GRPCClient) ExecuteCommand(ctx context.Context, req *ExecuteCommandRequest) (*ExecuteCommandResponse, error) {
	client, err := c.getClient(req.HostID)
	if err != nil {
		return nil, err
	}

	pbReq := &pb.ExecuteCommandRequest{
		HostId:  req.HostID,
		Command: req.Command,
		Timeout: int32(req.Timeout),
	}

	resp, err := client.ExecuteCommand(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	return &ExecuteCommandResponse{
		ExitCode: int(resp.ExitCode),
		Stdout:   resp.Stdout,
		Stderr:   resp.Stderr,
	}, nil
}

// Ping 健康检查
func (c *GRPCClient) Ping(ctx context.Context, hostID string) error {
	client, err := c.getClient(hostID)
	if err != nil {
		return err
	}

	pbReq := &pb.PingRequest{HostId: hostID}
	resp, err := client.Ping(ctx, pbReq)
	if err != nil {
		return err
	}

	if !resp.Ok {
		return fmt.Errorf("ping failed")
	}
	return nil
}
