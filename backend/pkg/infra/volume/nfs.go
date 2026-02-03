package volume

import "fmt"

// NFSVolumeConfig NFS 存储卷配置
type NFSVolumeConfig struct {
	// 基础配置
	Server      string   `json:"server"`       // NFS 服务器地址
	ServerPath  string   `json:"server_path"`  // 服务器路径
	MountPath   string   `json:"mount_path"`   // 容器内挂载路径

	// 挂载选项
	ReadOnly    bool     `json:"read_only"`    // 只读模式
	Options     []string `json:"options"`      // 挂载选项
}

// GetType 获取存储卷类型
func (c *NFSVolumeConfig) GetType() VolumeType {
	return VolumeTypeNFS
}

// GetMountPoint 获取挂载点
func (c *NFSVolumeConfig) GetMountPoint() string {
	return c.MountPath
}

// Validate 验证配置
func (c *NFSVolumeConfig) Validate() error {
	if c.Server == "" {
		return fmt.Errorf("NFS 服务器地址不能为空")
	}
	if c.ServerPath == "" {
		return fmt.Errorf("NFS 服务器路径不能为空")
	}
	if c.MountPath == "" {
		return fmt.Errorf("挂载路径不能为空")
	}
	return nil
}

// GetK8sVolume 获取 K8S Volume 配置
func (c *NFSVolumeConfig) GetK8sVolume(envID string) *K8sVolume {
	return &K8sVolume{
		Name: fmt.Sprintf("nfs-%s", envID),
		VolumeSpec: map[string]interface{}{
			"nfs": map[string]interface{}{
				"server":   c.Server,
				"path":     c.ServerPath,
				"readOnly": c.ReadOnly,
			},
		},
		MountPath: c.MountPath,
		ReadOnly:  c.ReadOnly,
	}
}

// GetDockerVolume 获取 Docker Volume 配置
func (c *NFSVolumeConfig) GetDockerVolume(envID string) *DockerVolume {
	return &DockerVolume{
		Type:     "volume",
		Source:   fmt.Sprintf("nfs-%s", envID),
		Target:   c.MountPath,
		ReadOnly: c.ReadOnly,
		Options: map[string]interface{}{
			"type": "nfs",
			"o":    fmt.Sprintf("addr=%s,rw", c.Server),
			"device": fmt.Sprintf(":%s", c.ServerPath),
		},
	}
}
