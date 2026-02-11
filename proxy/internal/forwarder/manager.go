package forwarder

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/YoungBoyGod/remotegpu-proxy/internal/models"
	"github.com/YoungBoyGod/remotegpu-proxy/internal/portpool"
	"github.com/google/uuid"
)

// MappingInfo 映射信息（用于查询返回）
type MappingInfo struct {
	ID           string    `json:"id"`
	EnvID        string    `json:"env_id"`
	ServiceType  string    `json:"service_type"`
	ExternalPort int       `json:"external_port"`
	TargetHost   string    `json:"target_host"`
	TargetPort   int       `json:"target_port"`
	Protocol     string    `json:"protocol"`
	ConnCount    int64     `json:"conn_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// ManagerStats 管理器统计信息
type ManagerStats struct {
	ActiveMappings int                `json:"active_mappings"`
	PoolStats      portpool.PoolStats `json:"pool_stats"`
}

// mappingEntry 内部映射条目
type mappingEntry struct {
	info      MappingInfo
	forwarder *TCPForwarder
}

// Manager 转发管理器，管理所有 TCP 转发器和 HTTP 代理路由
type Manager struct {
	mu        sync.RWMutex
	pool      *portpool.Pool
	mappings  map[int]*mappingEntry // externalPort -> entry
	httpProxy *HTTPProxy
}

// NewManager 创建转发管理器
func NewManager(pool *portpool.Pool, httpProxy *HTTPProxy) *Manager {
	return &Manager{
		pool:      pool,
		mappings:  make(map[int]*mappingEntry),
		httpProxy: httpProxy,
	}
}

// AddMapping 添加映射：分配端口 + 启动转发
func (m *Manager) AddMapping(req *models.MappingRequest) (*MappingInfo, error) {
	// 分配端口
	port, err := m.pool.Allocate(req.EnvID)
	if err != nil {
		return nil, fmt.Errorf("分配端口失败: %w", err)
	}

	protocol := req.Protocol
	if protocol == "" {
		protocol = "tcp"
	}

	// 启动 TCP 转发器
	fwd := NewTCPForwarder(port, req.TargetHost, req.TargetPort)
	if err := fwd.Start(); err != nil {
		m.pool.Release(port)
		return nil, fmt.Errorf("启动转发器失败: %w", err)
	}

	// 如果启用了 HTTP 代理，同时添加 HTTP 路由
	if m.httpProxy != nil {
		m.httpProxy.AddRoute(port, req.TargetHost, req.TargetPort)
	}

	info := MappingInfo{
		ID:           uuid.New().String(),
		EnvID:        req.EnvID,
		ServiceType:  req.ServiceType,
		ExternalPort: port,
		TargetHost:   req.TargetHost,
		TargetPort:   req.TargetPort,
		Protocol:     protocol,
		CreatedAt:    time.Now(),
	}

	m.mu.Lock()
	m.mappings[port] = &mappingEntry{info: info, forwarder: fwd}
	m.mu.Unlock()

	slog.Info("映射已添加", "id", info.ID, "env", req.EnvID, "port", port, "target", fmt.Sprintf("%s:%d", req.TargetHost, req.TargetPort))
	return &info, nil
}

// RemoveMapping 移除指定端口的映射：停止转发 + 释放端口
func (m *Manager) RemoveMapping(externalPort int) error {
	m.mu.Lock()
	entry, ok := m.mappings[externalPort]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("端口 %d 无映射记录", externalPort)
	}
	delete(m.mappings, externalPort)
	m.mu.Unlock()

	entry.forwarder.Stop()
	m.pool.Release(externalPort)

	if m.httpProxy != nil {
		m.httpProxy.RemoveRoute(externalPort)
	}

	slog.Info("映射已移除", "port", externalPort, "env", entry.info.EnvID)
	return nil
}

// RemoveByEnvID 移除指定环境的所有映射
func (m *Manager) RemoveByEnvID(envID string) error {
	m.mu.RLock()
	var ports []int
	for port, entry := range m.mappings {
		if entry.info.EnvID == envID {
			ports = append(ports, port)
		}
	}
	m.mu.RUnlock()

	for _, port := range ports {
		if err := m.RemoveMapping(port); err != nil {
			slog.Error("移除映射失败", "port", port, "error", err)
		}
	}

	slog.Info("环境映射已清理", "env", envID, "count", len(ports))
	return nil
}

// ListMappings 列出所有活跃映射
func (m *Manager) ListMappings() []MappingInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]MappingInfo, 0, len(m.mappings))
	for _, entry := range m.mappings {
		info := entry.info
		info.ConnCount = entry.forwarder.ConnCount()
		list = append(list, info)
	}
	return list
}

// Stats 返回管理器统计信息
func (m *Manager) Stats() ManagerStats {
	m.mu.RLock()
	count := len(m.mappings)
	m.mu.RUnlock()

	return ManagerStats{
		ActiveMappings: count,
		PoolStats:      m.pool.Stats(),
	}
}

// StopAll 停止所有转发器
func (m *Manager) StopAll() {
	m.mu.Lock()
	entries := make([]*mappingEntry, 0, len(m.mappings))
	for _, entry := range m.mappings {
		entries = append(entries, entry)
	}
	m.mappings = make(map[int]*mappingEntry)
	m.mu.Unlock()

	for _, entry := range entries {
		entry.forwarder.Stop()
		m.pool.Release(entry.info.ExternalPort)
		if m.httpProxy != nil {
			m.httpProxy.RemoveRoute(entry.info.ExternalPort)
		}
	}
	slog.Info("所有转发器已停止", "count", len(entries))
}
