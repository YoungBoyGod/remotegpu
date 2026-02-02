package deployment

import (
	"fmt"
	"sync"
)

// DeploymentManager 部署管理器
type DeploymentManager struct {
	deployments map[string]*DeploymentInfo
	mu          sync.RWMutex
}

// DeploymentInfo 部署信息
type DeploymentInfo struct {
	EnvID  string           `json:"env_id"`
	Type   DeploymentType   `json:"type"`
	Config DeploymentConfig `json:"config"`
	Status DeploymentStatus `json:"status"`
	Error  string           `json:"error"`
}

// NewDeploymentManager 创建部署管理器
func NewDeploymentManager() *DeploymentManager {
	return &DeploymentManager{
		deployments: make(map[string]*DeploymentInfo),
	}
}

// Register 注册部署配置
func (m *DeploymentManager) Register(envID string, config DeploymentConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.deployments[envID] = &DeploymentInfo{
		EnvID:  envID,
		Type:   config.GetType(),
		Config: config,
		Status: DeploymentStatusPending,
	}

	return nil
}

// Get 获取部署信息
func (m *DeploymentManager) Get(envID string) (*DeploymentInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	info, ok := m.deployments[envID]
	if !ok {
		return nil, fmt.Errorf("部署信息不存在: %s", envID)
	}

	return info, nil
}

// UpdateStatus 更新部署状态
func (m *DeploymentManager) UpdateStatus(envID string, status DeploymentStatus, errorMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, ok := m.deployments[envID]
	if !ok {
		return fmt.Errorf("部署信息不存在: %s", envID)
	}

	info.Status = status
	info.Error = errorMsg

	return nil
}

// Delete 删除部署信息
func (m *DeploymentManager) Delete(envID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.deployments, envID)
}

// List 列出所有部署
func (m *DeploymentManager) List() []*DeploymentInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*DeploymentInfo, 0, len(m.deployments))
	for _, info := range m.deployments {
		list = append(list, info)
	}

	return list
}
