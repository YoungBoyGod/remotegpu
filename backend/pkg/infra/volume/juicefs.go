package volume

import "fmt"

// JuiceFSVolumeConfig JuiceFS 分布式文件系统卷配置
type JuiceFSVolumeConfig struct {
	// 元数据引擎配置
	Name        string `json:"name"`         // 文件系统名称
	MetaURL     string `json:"meta_url"`     // 元数据引擎 URL

	// 挂载配置
	MountPath   string `json:"mount_path"`   // 容器内挂载路径
	SubPath     string `json:"sub_path"`     // 子路径
	ReadOnly    bool   `json:"read_only"`    // 只读模式

	// 缓存配置
	CacheDir    string `json:"cache_dir"`    // 缓存目录
	CacheSize   int64  `json:"cache_size"`   // 缓存大小 (MB)
}

// GetType 获取存储卷类型
func (c *JuiceFSVolumeConfig) GetType() VolumeType {
	return VolumeTypeJuiceFS
}

// GetMountPoint 获取挂载点
func (c *JuiceFSVolumeConfig) GetMountPoint() string {
	return c.MountPath
}

// Validate 验证配置
func (c *JuiceFSVolumeConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("文件系统名称不能为空")
	}
	if c.MetaURL == "" {
		return fmt.Errorf("元数据引擎 URL 不能为空")
	}
	if c.MountPath == "" {
		return fmt.Errorf("挂载路径不能为空")
	}
	return nil
}

// GetK8sVolume 获取 K8S Volume 配置
func (c *JuiceFSVolumeConfig) GetK8sVolume(envID string) *K8sVolume {
	return &K8sVolume{
		Name: fmt.Sprintf("juicefs-%s", envID),
		VolumeSpec: map[string]interface{}{
			"hostPath": map[string]interface{}{
				"path": fmt.Sprintf("/mnt/juicefs/%s", envID),
				"type": "DirectoryOrCreate",
			},
		},
		MountPath: c.MountPath,
		SubPath:   c.SubPath,
		ReadOnly:  c.ReadOnly,
	}
}

// GetDockerVolume 获取 Docker Volume 配置
func (c *JuiceFSVolumeConfig) GetDockerVolume(envID string) *DockerVolume {
	return &DockerVolume{
		Type:     "bind",
		Source:   fmt.Sprintf("/mnt/juicefs/%s", envID),
		Target:   c.MountPath,
		ReadOnly: c.ReadOnly,
	}
}
