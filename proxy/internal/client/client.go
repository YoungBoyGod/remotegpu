package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ServerClient 后端通信客户端
type ServerClient struct {
	baseURL    string
	proxyID    string
	token      string
	httpClient *http.Client
}

// Config 客户端配置
type Config struct {
	ServerURL string
	ProxyID   string
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
		baseURL: cfg.ServerURL,
		proxyID: cfg.ProxyID,
		token:   cfg.Token,
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
		req.Header.Set("X-Proxy-Token", c.token)
	}
	return c.httpClient.Do(req)
}

// Register 向后端注册 Proxy
func (c *ServerClient) Register(innerIP, outerIP string, portStart, portEnd int) error {
	reqBody := map[string]interface{}{
		"proxy_id":    c.proxyID,
		"inner_ip":    innerIP,
		"outer_ip":    outerIP,
		"port_start":  portStart,
		"port_end":    portEnd,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/proxy/register", c.baseURL)
	resp, err := c.doPost(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("register failed: %s", result.Msg)
	}
	return nil
}

// PortUsage 端口使用情况
type PortUsage struct {
	Total     int `json:"total"`
	Used      int `json:"used"`
	Available int `json:"available"`
}

// Heartbeat 上报心跳，携带端口使用情况
func (c *ServerClient) Heartbeat(usage *PortUsage) error {
	reqBody := map[string]interface{}{
		"proxy_id": c.proxyID,
	}
	if usage != nil {
		reqBody["port_usage"] = usage
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/proxy/heartbeat", c.baseURL)
	resp, err := c.doPost(url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("heartbeat failed: %s", result.Msg)
	}
	return nil
}
