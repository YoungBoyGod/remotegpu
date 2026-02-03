package network

import (
	"fmt"
	"sync"
)

// EnvNetworkConfig 环境网络配置
type EnvNetworkConfig struct {
	EnvID         string          `json:"env_id"`
	PortMappings  []*PortMapping  `json:"port_mappings"`
	FirewallRules []*FirewallRule `json:"firewall_rules"`
	DNSRecords    []*DNSRecord    `json:"dns_records"`
}

// NetworkManager 网络管理器
type NetworkManager struct {
	portManager     *PortManager
	firewallManager *FirewallManager
	dnsManager      *DNSManager
	envConfigs      map[string]*EnvNetworkConfig
	mu              sync.RWMutex
}

// NewNetworkManager 创建网络管理器
func NewNetworkManager(firewallType FirewallType, dnsProvider DNSProvider, baseDomain string) *NetworkManager {
	return &NetworkManager{
		portManager:     NewPortManager(),
		firewallManager: NewFirewallManager(firewallType),
		dnsManager:      NewDNSManager(dnsProvider, baseDomain),
		envConfigs:      make(map[string]*EnvNetworkConfig),
	}
}

// AllocateServicePort 为服务分配端口
func (m *NetworkManager) AllocateServicePort(envID string, serviceType ServiceType, protocol string, description string) (*PortMapping, error) {
	// 分配端口
	port, err := m.portManager.AllocatePort(serviceType)
	if err != nil {
		return nil, fmt.Errorf("分配端口失败: %w", err)
	}

	// 创建端口映射
	mapping := &PortMapping{
		EnvID:        envID,
		ServiceType:  serviceType,
		InternalPort: port,
		ExternalPort: port,
		PublicPort:   port,
		Protocol:     protocol,
		Description:  description,
	}

	// 保存到环境配置
	m.mu.Lock()
	defer m.mu.Unlock()

	config, ok := m.envConfigs[envID]
	if !ok {
		config = &EnvNetworkConfig{
			EnvID:         envID,
			PortMappings:  make([]*PortMapping, 0),
			FirewallRules: make([]*FirewallRule, 0),
			DNSRecords:    make([]*DNSRecord, 0),
		}
		m.envConfigs[envID] = config
	}

	config.PortMappings = append(config.PortMappings, mapping)

	return mapping, nil
}

// CreateFirewallRule 创建防火墙规则
func (m *NetworkManager) CreateFirewallRule(envID string, rule *FirewallRule) error {
	// 创建防火墙规则
	if err := m.firewallManager.CreateRule(rule); err != nil {
		return fmt.Errorf("创建防火墙规则失败: %w", err)
	}

	// 保存到环境配置
	m.mu.Lock()
	defer m.mu.Unlock()

	config, ok := m.envConfigs[envID]
	if !ok {
		config = &EnvNetworkConfig{
			EnvID:         envID,
			PortMappings:  make([]*PortMapping, 0),
			FirewallRules: make([]*FirewallRule, 0),
			DNSRecords:    make([]*DNSRecord, 0),
		}
		m.envConfigs[envID] = config
	}

	config.FirewallRules = append(config.FirewallRules, rule)

	return nil
}

// CreateDNSRecord 创建 DNS 记录
func (m *NetworkManager) CreateDNSRecord(envID string, record *DNSRecord) error {
	// 创建 DNS 记录
	if err := m.dnsManager.CreateRecord(record); err != nil {
		return fmt.Errorf("创建 DNS 记录失败: %w", err)
	}

	// 保存到环境配置
	m.mu.Lock()
	defer m.mu.Unlock()

	config, ok := m.envConfigs[envID]
	if !ok {
		config = &EnvNetworkConfig{
			EnvID:         envID,
			PortMappings:  make([]*PortMapping, 0),
			FirewallRules: make([]*FirewallRule, 0),
			DNSRecords:    make([]*DNSRecord, 0),
		}
		m.envConfigs[envID] = config
	}

	config.DNSRecords = append(config.DNSRecords, record)

	return nil
}

// GenerateSubdomain 生成子域名
func (m *NetworkManager) GenerateSubdomain(envID string, serviceType ServiceType) string {
	return m.dnsManager.GenerateSubdomain(envID, serviceType)
}

// GetEnvConfig 获取环境网络配置
func (m *NetworkManager) GetEnvConfig(envID string) (*EnvNetworkConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, ok := m.envConfigs[envID]
	if !ok {
		return nil, fmt.Errorf("环境网络配置不存在: %s", envID)
	}

	return config, nil
}

// GetPortMappings 获取环境的端口映射
func (m *NetworkManager) GetPortMappings(envID string) ([]*PortMapping, error) {
	config, err := m.GetEnvConfig(envID)
	if err != nil {
		return nil, err
	}
	return config.PortMappings, nil
}

// CleanupEnvironment 清理环境网络配置
func (m *NetworkManager) CleanupEnvironment(envID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	config, ok := m.envConfigs[envID]
	if !ok {
		return fmt.Errorf("环境网络配置不存在: %s", envID)
	}

	// 释放端口
	for _, mapping := range config.PortMappings {
		m.portManager.ReleasePort(mapping.ExternalPort)
	}

	// 删除防火墙规则
	for _, rule := range config.FirewallRules {
		if err := m.firewallManager.DeleteRule(rule.ID); err != nil {
			return fmt.Errorf("删除防火墙规则失败: %w", err)
		}
	}

	// 删除 DNS 记录
	for _, record := range config.DNSRecords {
		if err := m.dnsManager.DeleteRecord(record.ID); err != nil {
			return fmt.Errorf("删除 DNS 记录失败: %w", err)
		}
	}

	// 删除环境配置
	delete(m.envConfigs, envID)

	return nil
}

// ConfigureEnvironment 为环境配置完整的网络设置
func (m *NetworkManager) ConfigureEnvironment(envID string, services []ServiceType, hostIP string) (*EnvNetworkConfig, error) {
	// 为每个服务分配端口
	for _, serviceType := range services {
		mapping, err := m.AllocateServicePort(envID, serviceType, "tcp", fmt.Sprintf("%s service", serviceType))
		if err != nil {
			// 清理已分配的资源
			m.CleanupEnvironment(envID)
			return nil, fmt.Errorf("分配 %s 端口失败: %w", serviceType, err)
		}

		// 创建防火墙规则允许访问
		rule := &FirewallRule{
			ID:          fmt.Sprintf("%s-%s-allow", envID, serviceType),
			DestIP:      hostIP,
			DestPort:    mapping.ExternalPort,
			Protocol:    "tcp",
			Action:      "allow",
			Description: fmt.Sprintf("Allow %s access for %s", serviceType, envID),
		}

		if err := m.CreateFirewallRule(envID, rule); err != nil {
			// 清理已分配的资源
			m.CleanupEnvironment(envID)
			return nil, fmt.Errorf("创建防火墙规则失败: %w", err)
		}

		// 创建 DNS 记录
		subdomain := m.GenerateSubdomain(envID, serviceType)
		record := &DNSRecord{
			ID:     fmt.Sprintf("%s-%s-dns", envID, serviceType),
			Domain: subdomain,
			Type:   "A",
			Value:  hostIP,
			TTL:    300,
		}

		if err := m.CreateDNSRecord(envID, record); err != nil {
			// 清理已分配的资源
			m.CleanupEnvironment(envID)
			return nil, fmt.Errorf("创建 DNS 记录失败: %w", err)
		}
	}

	return m.GetEnvConfig(envID)
}
