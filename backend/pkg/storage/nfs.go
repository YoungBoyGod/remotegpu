package storage

import "fmt"

// NFSConfig NFS 存储配置
type NFSConfig struct {
	// 基础配置
	Type        StorageType `json:"type"`         // 存储类型
	Server      string      `json:"server"`       // NFS 服务器地址
	ServerPath  string      `json:"server_path"`  // 服务器路径
	MountPoint  string      `json:"mount_point"`  // 本地挂载点

	// 挂载选项
	Options     []string    `json:"options"`      // 挂载选项
	ReadOnly    bool        `json:"read_only"`    // 只读模式
	Version     string      `json:"version"`      // NFS 版本 (v3/v4)

	// 性能配置
	RSize       int         `json:"rsize"`        // 读取块大小
	WSize       int         `json:"wsize"`        // 写入块大小
	Timeout     int         `json:"timeout"`      // 超时时间
	Retrans     int         `json:"retrans"`      // 重传次数
}

// GetType 获取存储类型
func (c *NFSConfig) GetType() StorageType {
	return StorageTypeNFS
}

// GetMountPoint 获取挂载点
func (c *NFSConfig) GetMountPoint() string {
	return c.MountPoint
}

// Validate 验证配置
func (c *NFSConfig) Validate() error {
	if c.Server == "" {
		return fmt.Errorf("NFS 服务器地址不能为空")
	}
	if c.ServerPath == "" {
		return fmt.Errorf("NFS 服务器路径不能为空")
	}
	if c.MountPoint == "" {
		return fmt.Errorf("挂载点不能为空")
	}
	return nil
}

// GetK8sVolume 获取 K8S Volume 配置
func (c *NFSConfig) GetK8sVolume(envID string) map[string]interface{} {
	return map[string]interface{}{
		"volume": map[string]interface{}{
			"name": fmt.Sprintf("nfs-%s", envID),
			"nfs": map[string]interface{}{
				"server":   c.Server,
				"path":     c.ServerPath,
				"readOnly": c.ReadOnly,
			},
		},
		"volumeMount": map[string]interface{}{
			"name":      fmt.Sprintf("nfs-%s", envID),
			"mountPath": c.MountPoint,
			"readOnly":  c.ReadOnly,
		},
	}
}

// GetDockerVolume 获取 Docker Volume 配置
func (c *NFSConfig) GetDockerVolume(envID string) map[string]interface{} {
	return map[string]interface{}{
		"type":   "volume",
		"source": fmt.Sprintf("nfs-%s", envID),
		"target": c.MountPoint,
		"volume_options": map[string]interface{}{
			"driver": "local",
			"driver_opts": map[string]string{
				"type":   "nfs",
				"o":      fmt.Sprintf("addr=%s,rw", c.Server),
				"device": fmt.Sprintf(":%s", c.ServerPath),
			},
		},
	}
}
