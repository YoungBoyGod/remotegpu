package volume

import (
	"encoding/json"
	"fmt"
)

// ConfigManager 卷配置管理器
type ConfigManager struct {
	configs map[string]VolumeConfig
}

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs: make(map[string]VolumeConfig),
	}
}

// Register 注册卷配置
func (m *ConfigManager) Register(envID string, config VolumeConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}
	m.configs[envID] = config
	return nil
}

// Get 获取卷配置
func (m *ConfigManager) Get(envID string) (VolumeConfig, error) {
	config, ok := m.configs[envID]
	if !ok {
		return nil, fmt.Errorf("环境 %s 的卷配置不存在", envID)
	}
	return config, nil
}

// Remove 移除卷配置
func (m *ConfigManager) Remove(envID string) {
	delete(m.configs, envID)
}

// ParseConfig 从 JSON 解析卷配置
func ParseConfig(volumeType VolumeType, data []byte) (VolumeConfig, error) {
	var config VolumeConfig

	switch volumeType {
	case VolumeTypeNFS:
		var nfsConfig NFSVolumeConfig
		if err := json.Unmarshal(data, &nfsConfig); err != nil {
			return nil, fmt.Errorf("解析 NFS 配置失败: %w", err)
		}
		config = &nfsConfig

	case VolumeTypeS3:
		var s3Config S3VolumeConfig
		if err := json.Unmarshal(data, &s3Config); err != nil {
			return nil, fmt.Errorf("解析 S3 配置失败: %w", err)
		}
		config = &s3Config

	case VolumeTypeJuiceFS:
		var juicefsConfig JuiceFSVolumeConfig
		if err := json.Unmarshal(data, &juicefsConfig); err != nil {
			return nil, fmt.Errorf("解析 JuiceFS 配置失败: %w", err)
		}
		config = &juicefsConfig

	default:
		return nil, fmt.Errorf("不支持的卷类型: %s", volumeType)
	}

	return config, nil
}
