package host

import (
	"fmt"
	"sync"
)

// HostManager 主机管理器
type HostManager struct {
	selector      HostSelector
	healthMonitor *HealthMonitor
	hosts         map[string]*HostInfo
	mu            sync.RWMutex
}

// NewHostManager 创建主机管理器
func NewHostManager(selector HostSelector, healthMonitor *HealthMonitor) *HostManager {
	return &HostManager{
		selector:      selector,
		healthMonitor: healthMonitor,
		hosts:         make(map[string]*HostInfo),
	}
}

// RegisterHost 注册主机
func (m *HostManager) RegisterHost(host *HostInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hosts[host.ID] = host
}

// UnregisterHost 注销主机
func (m *HostManager) UnregisterHost(hostID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.hosts, hostID)
}

// GetHost 获取主机信息
func (m *HostManager) GetHost(hostID string) (*HostInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	host, ok := m.hosts[hostID]
	if !ok {
		return nil, fmt.Errorf("主机不存在: %s", hostID)
	}
	return host, nil
}

// ListHosts 列出所有主机
func (m *HostManager) ListHosts() []*HostInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hosts := make([]*HostInfo, 0, len(m.hosts))
	for _, host := range m.hosts {
		hosts = append(hosts, host)
	}
	return hosts
}

// SelectHost 选择主机
func (m *HostManager) SelectHost(req *ResourceRequirement) (*HostInfo, error) {
	hosts := m.ListHosts()
	return m.selector.Select(hosts, req)
}

// UpdateHostResources 更新主机资源使用量
func (m *HostManager) UpdateHostResources(hostID string, cpu int, memory int64, gpu int, add bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	host, ok := m.hosts[hostID]
	if !ok {
		return fmt.Errorf("主机不存在: %s", hostID)
	}

	if add {
		host.UsedCPU += cpu
		host.UsedMemory += memory
		host.UsedGPU += gpu
	} else {
		host.UsedCPU -= cpu
		host.UsedMemory -= memory
		host.UsedGPU -= gpu
	}

	return nil
}

// GetHealthStatus 获取主机健康状态
func (m *HostManager) GetHealthStatus(hostID string) (*HealthCheckResult, error) {
	if m.healthMonitor == nil {
		return nil, fmt.Errorf("健康监控器未启用")
	}

	result, ok := m.healthMonitor.GetResult(hostID)
	if !ok {
		return nil, fmt.Errorf("主机健康状态未知: %s", hostID)
	}

	return result, nil
}
