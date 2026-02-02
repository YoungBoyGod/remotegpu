package network

import (
	"fmt"
	"sync"
)

// PortManager 端口管理器
type PortManager struct {
	ranges    map[ServiceType]PortRange
	allocated map[int]bool
	mu        sync.RWMutex
}

// NewPortManager 创建端口管理器
func NewPortManager() *PortManager {
	return &PortManager{
		ranges: map[ServiceType]PortRange{
			ServiceTypeSSH:         {Start: 22000, End: 22999},
			ServiceTypeRDP:         {Start: 33890, End: 34889},
			ServiceTypeVNC:         {Start: 5900, End: 6899},
			ServiceTypeJupyter:     {Start: 8888, End: 9887},
			ServiceTypeTensorBoard: {Start: 6006, End: 7005},
			ServiceTypeCustom:      {Start: 10000, End: 19999},
		},
		allocated: make(map[int]bool),
	}
}

// AllocatePort 分配端口
func (m *PortManager) AllocatePort(serviceType ServiceType) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	portRange, ok := m.ranges[serviceType]
	if !ok {
		return 0, fmt.Errorf("不支持的服务类型: %s", serviceType)
	}

	for port := portRange.Start; port <= portRange.End; port++ {
		if !m.allocated[port] {
			m.allocated[port] = true
			return port, nil
		}
	}

	return 0, fmt.Errorf("端口池已满,无可用端口")
}

// ReleasePort 释放端口
func (m *PortManager) ReleasePort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.allocated, port)
}

// IsAllocated 检查端口是否已分配
func (m *PortManager) IsAllocated(port int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.allocated[port]
}
