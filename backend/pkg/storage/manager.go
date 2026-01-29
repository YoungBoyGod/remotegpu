package storage

import (
	"fmt"
	"sync"

	"github.com/YoungBoyGod/remotegpu/config"
)

// Manager 存储管理器
type Manager struct {
	backends map[string]Storage
	default_ string
	mu       sync.RWMutex
}

var globalManager *Manager

// NewManager 创建存储管理器
func NewManager(cfg config.StorageConfig) (*Manager, error) {
	m := &Manager{
		backends: make(map[string]Storage),
		default_: cfg.Default,
	}

	for _, backend := range cfg.Backends {
		if !backend.Enabled {
			continue
		}

		var s Storage
		var err error

		switch backend.Type {
		case "local":
			s, err = NewLocalStorage(backend.Name, backend.Path)
		case "rustfs", "s3":
			s, err = NewS3Storage(
				backend.Name,
				backend.Type,
				backend.Endpoint,
				backend.AccessKey,
				backend.SecretKey,
				backend.Bucket,
				backend.Region,
			)
		default:
			return nil, fmt.Errorf("未知的存储类型: %s", backend.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("初始化存储后端 %s 失败: %w", backend.Name, err)
		}

		m.backends[backend.Name] = s
	}

	if len(m.backends) == 0 {
		return nil, fmt.Errorf("没有启用的存储后端")
	}

	if _, ok := m.backends[m.default_]; !ok && m.default_ != "" {
		return nil, fmt.Errorf("默认存储后端 %s 不存在", m.default_)
	}

	globalManager = m
	return m, nil
}

// GetManager 获取全局存储管理器
func GetManager() *Manager {
	return globalManager
}

// Get 获取指定名称的存储后端
func (m *Manager) Get(name string) (Storage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if name == "" {
		name = m.default_
	}

	s, ok := m.backends[name]
	if !ok {
		return nil, fmt.Errorf("存储后端 %s 不存在", name)
	}
	return s, nil
}

// Default 获取默认存储后端
func (m *Manager) Default() (Storage, error) {
	return m.Get(m.default_)
}

// List 列出所有启用的存储后端
func (m *Manager) List() []BackendInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []BackendInfo
	for name, s := range m.backends {
		list = append(list, BackendInfo{
			Name:      name,
			Type:      s.Type(),
			IsDefault: name == m.default_,
		})
	}
	return list
}

// BackendInfo 存储后端信息
type BackendInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
}
