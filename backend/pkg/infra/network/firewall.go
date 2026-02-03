package network

import "fmt"

// FirewallType 防火墙类型
type FirewallType string

const (
	// FirewallTypeIPTables iptables
	FirewallTypeIPTables FirewallType = "iptables"

	// FirewallTypeFirewalld firewalld
	FirewallTypeFirewalld FirewallType = "firewalld"

	// FirewallTypeCloud 云厂商防火墙
	FirewallTypeCloud FirewallType = "cloud"
)

// FirewallManager 防火墙管理器
type FirewallManager struct {
	firewallType FirewallType
	rules        map[string]*FirewallRule
}

// NewFirewallManager 创建防火墙管理器
func NewFirewallManager(firewallType FirewallType) *FirewallManager {
	return &FirewallManager{
		firewallType: firewallType,
		rules:        make(map[string]*FirewallRule),
	}
}

// CreateRule 创建防火墙规则
func (m *FirewallManager) CreateRule(rule *FirewallRule) error {
	// TODO: 实现防火墙规则创建
	// 根据 firewallType 调用不同的实现
	m.rules[rule.ID] = rule
	return fmt.Errorf("TODO: 实现防火墙规则创建")
}

// DeleteRule 删除防火墙规则
func (m *FirewallManager) DeleteRule(ruleID string) error {
	// TODO: 实现防火墙规则删除
	delete(m.rules, ruleID)
	return fmt.Errorf("TODO: 实现防火墙规则删除")
}

// GetRule 获取防火墙规则
func (m *FirewallManager) GetRule(ruleID string) (*FirewallRule, error) {
	rule, ok := m.rules[ruleID]
	if !ok {
		return nil, fmt.Errorf("规则不存在: %s", ruleID)
	}
	return rule, nil
}
