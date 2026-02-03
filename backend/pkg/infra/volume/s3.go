package volume

import "fmt"

// S3VolumeConfig S3 对象存储卷配置
type S3VolumeConfig struct {
	// S3 连接配置
	Endpoint        string `json:"endpoint"`          // S3 端点
	AccessKeyID     string `json:"access_key_id"`     // Access Key ID
	SecretAccessKey string `json:"secret_access_key"` // Secret Access Key
	Region          string `json:"region"`            // 区域
	BucketName      string `json:"bucket_name"`       // 存储桶名称

	// 挂载配置
	MountPath       string `json:"mount_path"`        // 容器内挂载路径
	SubPath         string `json:"sub_path"`          // 子路径
	ReadOnly        bool   `json:"read_only"`         // 只读模式

	// S3FS 配置
	UseS3FS         bool   `json:"use_s3fs"`          // 是否使用 s3fs
	CacheDir        string `json:"cache_dir"`         // 缓存目录
}

// GetType 获取存储卷类型
func (c *S3VolumeConfig) GetType() VolumeType {
	return VolumeTypeS3
}

// GetMountPoint 获取挂载点
func (c *S3VolumeConfig) GetMountPoint() string {
	return c.MountPath
}

// Validate 验证配置
func (c *S3VolumeConfig) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("S3 端点不能为空")
	}
	if c.BucketName == "" {
		return fmt.Errorf("存储桶名称不能为空")
	}
	if c.MountPath == "" {
		return fmt.Errorf("挂载路径不能为空")
	}
	return nil
}

// GetK8sVolume 获取 K8S Volume 配置
// 注意: K8S 不直接支持 S3,需要使用 CSI 驱动或通过 hostPath + s3fs
func (c *S3VolumeConfig) GetK8sVolume(envID string) *K8sVolume {
	// 使用 hostPath 方式(需要在主机上先挂载 s3fs)
	return &K8sVolume{
		Name: fmt.Sprintf("s3-%s", envID),
		VolumeSpec: map[string]interface{}{
			"hostPath": map[string]interface{}{
				"path": fmt.Sprintf("/mnt/s3/%s", envID),
				"type": "DirectoryOrCreate",
			},
		},
		MountPath: c.MountPath,
		SubPath:   c.SubPath,
		ReadOnly:  c.ReadOnly,
	}
}

// GetDockerVolume 获取 Docker Volume 配置
func (c *S3VolumeConfig) GetDockerVolume(envID string) *DockerVolume {
	// Docker 使用 bind mount 方式(需要在主机上先挂载 s3fs)
	return &DockerVolume{
		Type:     "bind",
		Source:   fmt.Sprintf("/mnt/s3/%s", envID),
		Target:   c.MountPath,
		ReadOnly: c.ReadOnly,
	}
}
