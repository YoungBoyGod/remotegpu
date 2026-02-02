package remote

import (
	"encoding/json"
	"fmt"
)

// AccessManager 访问配置管理器
type AccessManager struct {
	configs map[string][]AccessConfig
}

// NewAccessManager 创建访问配置管理器
func NewAccessManager() *AccessManager {
	return &AccessManager{
		configs: make(map[string][]AccessConfig),
	}
}

// Register 注册访问配置
func (m *AccessManager) Register(envID string, config AccessConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}
	m.configs[envID] = append(m.configs[envID], config)
	return nil
}

// Get 获取环境的所有访问配置
func (m *AccessManager) Get(envID string) []AccessConfig {
	return m.configs[envID]
}

// GetByProtocol 获取指定协议的访问配置
func (m *AccessManager) GetByProtocol(envID string, protocol Protocol) AccessConfig {
	configs := m.configs[envID]
	for _, config := range configs {
		if config.GetProtocol() == protocol {
			return config
		}
	}
	return nil
}

// Remove 移除环境的所有访问配置
func (m *AccessManager) Remove(envID string) {
	delete(m.configs, envID)
}

// GenerateAllAccessInfo 生成环境的所有访问信息
func (m *AccessManager) GenerateAllAccessInfo(envID string, host *HostInfo) []*AccessInfo {
	configs := m.configs[envID]
	var infos []*AccessInfo
	for _, config := range configs {
		info := config.GenerateAccessInfo(envID, host)
		infos = append(infos, info)
	}
	return infos
}

// ParseConfig 从 JSON 解析访问配置
func ParseConfig(protocol Protocol, data []byte) (AccessConfig, error) {
	var config AccessConfig

	switch protocol {
	case ProtocolSSH:
		var sshConfig SSHConfig
		if err := json.Unmarshal(data, &sshConfig); err != nil {
			return nil, fmt.Errorf("解析 SSH 配置失败: %w", err)
		}
		config = &sshConfig

	case ProtocolRDP:
		var rdpConfig RDPConfig
		if err := json.Unmarshal(data, &rdpConfig); err != nil {
			return nil, fmt.Errorf("解析 RDP 配置失败: %w", err)
		}
		config = &rdpConfig

	case ProtocolVNC:
		var vncConfig VNCConfig
		if err := json.Unmarshal(data, &vncConfig); err != nil {
			return nil, fmt.Errorf("解析 VNC 配置失败: %w", err)
		}
		config = &vncConfig

	default:
		return nil, fmt.Errorf("不支持的协议类型: %s", protocol)
	}

	return config, nil
}
