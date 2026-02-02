package storage

// StorageType 存储类型
type StorageType string

const (
	// StorageTypeLocal 本地存储
	StorageTypeLocal StorageType = "local"

	// StorageTypeNFS NFS 网络文件系统
	StorageTypeNFS StorageType = "nfs"

	// StorageTypeS3 S3 对象存储
	StorageTypeS3 StorageType = "s3"

	// StorageTypeCeph Ceph 分布式存储
	StorageTypeCeph StorageType = "ceph"

	// StorageTypeJuiceFS JuiceFS 分布式文件系统
	StorageTypeJuiceFS StorageType = "juicefs"

	// StorageTypePVC Kubernetes PVC
	StorageTypePVC StorageType = "pvc"
)

// StorageConfig 存储配置接口
type StorageConfig interface {
	// GetType 获取存储类型
	GetType() StorageType

	// Validate 验证配置
	Validate() error

	// GetMountPoint 获取挂载点
	GetMountPoint() string
}

// VolumeConfig Volume 配置接口
type VolumeConfig interface {
	// GetK8sVolume 获取 K8S Volume 配置
	GetK8sVolume(envID string) map[string]interface{}

	// GetDockerVolume 获取 Docker Volume 配置
	GetDockerVolume(envID string) map[string]interface{}
}
