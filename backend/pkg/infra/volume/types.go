package volume

// VolumeType 存储卷类型
type VolumeType string

const (
	// VolumeTypeLocal 本地存储
	VolumeTypeLocal VolumeType = "local"

	// VolumeTypeNFS NFS 网络文件系统
	VolumeTypeNFS VolumeType = "nfs"

	// VolumeTypeS3 S3 对象存储
	VolumeTypeS3 VolumeType = "s3"

	// VolumeTypeCeph Ceph 分布式存储
	VolumeTypeCeph VolumeType = "ceph"

	// VolumeTypeJuiceFS JuiceFS 分布式文件系统
	VolumeTypeJuiceFS VolumeType = "juicefs"

	// VolumeTypePVC Kubernetes PVC
	VolumeTypePVC VolumeType = "pvc"
)

// VolumeConfig 存储卷配置接口
type VolumeConfig interface {
	// GetType 获取存储卷类型
	GetType() VolumeType

	// Validate 验证配置
	Validate() error

	// GetMountPoint 获取挂载点
	GetMountPoint() string

	// GetK8sVolume 获取 K8S Volume 配置
	GetK8sVolume(envID string) *K8sVolume

	// GetDockerVolume 获取 Docker Volume 配置
	GetDockerVolume(envID string) *DockerVolume
}

// K8sVolume K8S Volume 配置
type K8sVolume struct {
	Name        string                 `json:"name"`
	VolumeSpec  map[string]interface{} `json:"volume_spec"`
	MountPath   string                 `json:"mount_path"`
	SubPath     string                 `json:"sub_path,omitempty"`
	ReadOnly    bool                   `json:"read_only"`
}

// DockerVolume Docker Volume 配置
type DockerVolume struct {
	Type     string                 `json:"type"`      // bind/volume/tmpfs
	Source   string                 `json:"source"`    // 源路径或卷名
	Target   string                 `json:"target"`    // 容器内路径
	ReadOnly bool                   `json:"read_only"`
	Options  map[string]interface{} `json:"options,omitempty"`
}
