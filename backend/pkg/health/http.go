package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
)

// HTTPChecker HTTP服务健康检查器
type HTTPChecker struct {
	name     string
	endpoint string
	enabled  bool
	timeout  time.Duration
}

// NewHTTPChecker 创建HTTP健康检查器
func NewHTTPChecker(name, endpoint string, enabled bool) *HTTPChecker {
	return &HTTPChecker{
		name:     name,
		endpoint: endpoint,
		enabled:  enabled,
		timeout:  5 * time.Second,
	}
}

// Name 返回服务名称
func (c *HTTPChecker) Name() string {
	return c.name
}

// Check 执行健康检查
func (c *HTTPChecker) Check(ctx context.Context) *CheckResult {
	start := time.Now()
	result := &CheckResult{
		Service:   c.Name(),
		Timestamp: start,
	}

	// 如果服务未启用
	if !c.enabled {
		result.Status = StatusDisabled
		result.Message = "服务未启用"
		result.Latency = time.Since(start)
		return result
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: c.timeout,
	}

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", c.endpoint, nil)
	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("创建请求失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("请求失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Status = StatusHealthy
		result.Message = "服务正常"
	} else {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("HTTP状态码: %d", resp.StatusCode)
	}

	result.Latency = time.Since(start)
	result.Details = map[string]interface{}{
		"endpoint":    c.endpoint,
		"status_code": resp.StatusCode,
	}

	return result
}

// NewHarborChecker 创建Harbor健康检查器
func NewHarborChecker(cfg config.HarborConfig) *HTTPChecker {
	endpoint := fmt.Sprintf("%s/api/v2.0/health", cfg.Endpoint)
	return NewHTTPChecker("harbor", endpoint, cfg.Enabled)
}

// NewPrometheusChecker 创建Prometheus健康检查器
func NewPrometheusChecker(cfg config.PrometheusConfig) *HTTPChecker {
	endpoint := fmt.Sprintf("%s/-/healthy", cfg.Endpoint)
	return NewHTTPChecker("prometheus", endpoint, cfg.Enabled)
}

// NewJumpserverChecker 创建Jumpserver健康检查器
func NewJumpserverChecker(cfg config.JumpserverConfig) *HTTPChecker {
	endpoint := fmt.Sprintf("%s/api/health/", cfg.Endpoint)
	return NewHTTPChecker("jumpserver", endpoint, cfg.Enabled)
}

// NewNginxChecker 创建Nginx健康检查器
func NewNginxChecker(cfg config.NginxConfig) *HTTPChecker {
	return NewHTTPChecker("nginx", cfg.Endpoint, cfg.Enabled)
}

// NewUptimeKumaChecker 创建Uptime Kuma健康检查器
func NewUptimeKumaChecker(cfg config.UptimeKumaConfig) *HTTPChecker {
	return NewHTTPChecker("uptime-kuma", cfg.Endpoint, cfg.Enabled)
}

// NewGuacamoleChecker 创建Guacamole健康检查器
func NewGuacamoleChecker(cfg config.GuacamoleConfig) *HTTPChecker {
	return NewHTTPChecker("guacamole", cfg.Endpoint, cfg.Enabled)
}

// NewRustFSChecker 创建RustFS健康检查器
func NewRustFSChecker(cfg config.StorageConfig) *HTTPChecker {
	endpoint := fmt.Sprintf("%s/health", cfg.RustFS.Endpoint)
	enabled := cfg.Type == "rustfs"
	return NewHTTPChecker("rustfs", endpoint, enabled)
}
