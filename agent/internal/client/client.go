package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
	"github.com/google/uuid"
)

// ServerClient Server 通信客户端
type ServerClient struct {
	baseURL    string
	agentID    string
	machineID  string
	token      string
	httpClient *http.Client
}

// Config 客户端配置
type Config struct {
	ServerURL string
	AgentID   string
	MachineID string
	Token     string
	Timeout   time.Duration
}

// NewServerClient 创建客户端
func NewServerClient(cfg *Config) *ServerClient {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &ServerClient{
		baseURL:   cfg.ServerURL,
		agentID:   cfg.AgentID,
		machineID: cfg.MachineID,
		token:     cfg.Token,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// doPost 发送 POST 请求，自动添加认证 header
func (c *ServerClient) doPost(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.httpClient.Do(req)
}

// ClaimResponse 认领响应
type ClaimResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    *ClaimData `json:"data"`
}

type ClaimData struct {
	Tasks []*models.Task `json:"tasks"`
}

// ClaimTasks 从 Server 认领任务
func (c *ServerClient) ClaimTasks(limit int) ([]*models.Task, error) {
	reqBody := map[string]interface{}{
		"agent_id":   c.agentID,
		"machine_id": c.machineID,
		"limit":      limit,
		"request_id": uuid.New().String(),
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/agent/tasks/claim", c.baseURL)
	resp, err := c.doPost(url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ClaimResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("claim failed: %s", result.Message)
	}

	if result.Data == nil {
		return nil, nil
	}
	return result.Data.Tasks, nil
}

// ReportStart 上报任务开始
func (c *ServerClient) ReportStart(taskID, attemptID string) error {
	reqBody := map[string]string{
		"agent_id":   c.agentID,
		"attempt_id": attemptID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/agent/tasks/%s/start", c.baseURL, taskID)
	resp, err := c.doPost(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("report start failed: %s", result.Message)
	}
	return nil
}

// RenewLease 续约租约
func (c *ServerClient) RenewLease(taskID, attemptID string) error {
	reqBody := map[string]interface{}{
		"agent_id":   c.agentID,
		"attempt_id": attemptID,
		"extend_sec": 300,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/agent/tasks/%s/lease/renew", c.baseURL, taskID)
	resp, err := c.doPost(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("renew lease failed: %s", result.Message)
	}
	return nil
}

// ReportComplete 上报任务完成
func (c *ServerClient) ReportComplete(task *models.Task) error {
	reqBody := map[string]interface{}{
		"agent_id":   c.agentID,
		"attempt_id": task.AttemptID,
		"exit_code":  task.ExitCode,
		"stdout":     task.Stdout,
		"stderr":     task.Stderr,
		"error":      task.Error,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/agent/tasks/%s/complete", c.baseURL, task.ID)
	resp, err := c.doPost(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("report complete failed: %s", result.Message)
	}
	return nil
}
