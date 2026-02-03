package host

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthStatus 健康状态
type HealthStatus string

const (
	// HealthStatusHealthy 健康
	HealthStatusHealthy HealthStatus = "healthy"
	// HealthStatusUnhealthy 不健康
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	// HealthStatusUnknown 未知
	HealthStatusUnknown HealthStatus = "unknown"
)

// HealthCheckResult 健康检查结果
type HealthCheckResult struct {
	HostID    string       `json:"host_id"`
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message"`
	CheckedAt time.Time    `json:"checked_at"`
}

// HealthChecker 健康检查器接口
type HealthChecker interface {
	// Check 执行健康检查
	Check(ctx context.Context, host *HostInfo) (*HealthCheckResult, error)
}

// SimpleHealthChecker 简单健康检查器
type SimpleHealthChecker struct {
	timeout time.Duration
}

// NewSimpleHealthChecker 创建简单健康检查器
func NewSimpleHealthChecker(timeout time.Duration) *SimpleHealthChecker {
	return &SimpleHealthChecker{
		timeout: timeout,
	}
}

// Check 执行健康检查
func (c *SimpleHealthChecker) Check(ctx context.Context, host *HostInfo) (*HealthCheckResult, error) {
	result := &HealthCheckResult{
		HostID:    host.ID,
		CheckedAt: time.Now(),
	}

	// 检查主机状态
	if host.Status != "active" {
		result.Status = HealthStatusUnhealthy
		result.Message = fmt.Sprintf("主机状态异常: %s", host.Status)
		return result, nil
	}

	// 检查资源使用率
	usageRate := host.UsageRate()
	if usageRate > 0.95 {
		result.Status = HealthStatusUnhealthy
		result.Message = fmt.Sprintf("资源使用率过高: %.2f%%", usageRate*100)
		return result, nil
	}

	result.Status = HealthStatusHealthy
	result.Message = "主机运行正常"
	return result, nil
}

// HealthMonitor 健康监控器
type HealthMonitor struct {
	checker  HealthChecker
	interval time.Duration
	results  map[string]*HealthCheckResult
	mu       sync.RWMutex
	stopCh   chan struct{}
}

// NewHealthMonitor 创建健康监控器
func NewHealthMonitor(checker HealthChecker, interval time.Duration) *HealthMonitor {
	return &HealthMonitor{
		checker:  checker,
		interval: interval,
		results:  make(map[string]*HealthCheckResult),
		stopCh:   make(chan struct{}),
	}
}

// Start 启动健康监控
func (m *HealthMonitor) Start(hosts []*HostInfo) {
	go m.monitor(hosts)
}

// Stop 停止健康监控
func (m *HealthMonitor) Stop() {
	close(m.stopCh)
}

// GetResult 获取健康检查结果
func (m *HealthMonitor) GetResult(hostID string) (*HealthCheckResult, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result, ok := m.results[hostID]
	return result, ok
}

// monitor 监控循环
func (m *HealthMonitor) monitor(hosts []*HostInfo) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkAllHosts(hosts)
		case <-m.stopCh:
			return
		}
	}
}

// checkAllHosts 检查所有主机
func (m *HealthMonitor) checkAllHosts(hosts []*HostInfo) {
	ctx := context.Background()

	for _, host := range hosts {
		result, err := m.checker.Check(ctx, host)
		if err != nil {
			result = &HealthCheckResult{
				HostID:    host.ID,
				Status:    HealthStatusUnknown,
				Message:   fmt.Sprintf("健康检查失败: %v", err),
				CheckedAt: time.Now(),
			}
		}

		m.mu.Lock()
		m.results[host.ID] = result
		m.mu.Unlock()
	}
}
