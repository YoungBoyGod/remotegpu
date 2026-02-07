package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
)

// HTTPClient HTTP 客户端实现
type HTTPClient struct {
	client  *http.Client
	config  *config.AgentConfig
	hostMap map[string]string // hostID -> host:port
}

// NewHTTPClient 创建 HTTP 客户端
func NewHTTPClient(cfg *config.AgentConfig) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		config:  cfg,
		hostMap: make(map[string]string),
	}
}

// RegisterHost 注册主机地址
func (c *HTTPClient) RegisterHost(hostID, address string) {
	c.hostMap[hostID] = address
}

// getHostURL 获取主机 URL
func (c *HTTPClient) getHostURL(hostID, path string) (string, error) {
	addr, ok := c.hostMap[hostID]
	if !ok {
		return "", fmt.Errorf("host %s not registered", hostID)
	}
	return fmt.Sprintf("http://%s:%d%s", addr, c.config.HTTPPort, path), nil
}

// doRequest 执行 HTTP 请求
func (c *HTTPClient) doRequest(ctx context.Context, method, url string, body any) (*Response, error) {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if !result.Success || result.Code != 0 {
		return nil, fmt.Errorf("agent error: code=%d message=%s", result.Code, result.Message)
	}
	return &result, nil
}

// StopProcess 停止进程
func (c *HTTPClient) StopProcess(ctx context.Context, req *StopProcessRequest) (*Response, error) {
	url, err := c.getHostURL(req.HostID, "/api/v1/process/stop")
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, http.MethodPost, url, req)
}

// ResetSSH 重置SSH密钥
func (c *HTTPClient) ResetSSH(ctx context.Context, req *ResetSSHRequest) (*Response, error) {
	url, err := c.getHostURL(req.HostID, "/api/v1/ssh/reset")
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, http.MethodPost, url, req)
}

// SyncSSHKeys 同步SSH密钥（全量覆盖 authorized_keys）
func (c *HTTPClient) SyncSSHKeys(ctx context.Context, req *SyncSSHKeysRequest) (*Response, error) {
	url, err := c.getHostURL(req.HostID, "/api/v1/ssh/sync-keys")
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, http.MethodPost, url, req)
}

// CleanupMachine 清理机器
func (c *HTTPClient) CleanupMachine(ctx context.Context, req *CleanupRequest) (*Response, error) {
	url, err := c.getHostURL(req.HostID, "/api/v1/machine/cleanup")
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, http.MethodPost, url, req)
}

// MountDataset 挂载数据集
func (c *HTTPClient) MountDataset(ctx context.Context, req *MountDatasetRequest) (*Response, error) {
	url, err := c.getHostURL(req.HostID, "/api/v1/dataset/mount")
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, http.MethodPost, url, req)
}

// GetSystemInfo 获取系统信息
func (c *HTTPClient) GetSystemInfo(ctx context.Context, hostID string) (*SystemInfo, error) {
	url, err := c.getHostURL(hostID, "/api/v1/system/info")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info SystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ExecuteCommand 执行命令
func (c *HTTPClient) ExecuteCommand(ctx context.Context, req *ExecuteCommandRequest) (*ExecuteCommandResponse, error) {
	url, err := c.getHostURL(req.HostID, "/api/v1/command/exec")
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	result := &ExecuteCommandResponse{}
	if v, ok := data["exit_code"].(float64); ok {
		result.ExitCode = int(v)
	}
	if v, ok := data["stdout"].(string); ok {
		result.Stdout = v
	}
	if v, ok := data["stderr"].(string); ok {
		result.Stderr = v
	}
	return result, nil
}

// Ping 检查连接
func (c *HTTPClient) Ping(ctx context.Context, hostID string) error {
	url, err := c.getHostURL(hostID, "/api/v1/ping")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping failed: %d", resp.StatusCode)
	}
	return nil
}

// Close 关闭连接
func (c *HTTPClient) Close() error {
	c.client.CloseIdleConnections()
	return nil
}
