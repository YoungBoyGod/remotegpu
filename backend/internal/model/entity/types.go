package entity

// DeploymentMode 部署模式
type DeploymentMode string

const (
	// DeploymentK8sPod K8S Pod 部署
	DeploymentK8sPod DeploymentMode = "k8s_pod"
	// DeploymentK8sStateful K8S StatefulSet 部署
	DeploymentK8sStateful DeploymentMode = "k8s_stateful"
	// DeploymentDockerLocal 裸金属 Docker 部署
	DeploymentDockerLocal DeploymentMode = "docker_local"
	// DeploymentVM 虚拟机部署
	DeploymentVM DeploymentMode = "vm"
	// DeploymentBareMetal 裸金属直接安装
	DeploymentBareMetal DeploymentMode = "bare_metal"
)

// EnvironmentType 环境类型
type EnvironmentType string

const (
	// EnvTypeIDE IDE 环境
	EnvTypeIDE EnvironmentType = "ide"
	// EnvTypeTerminal 命令行环境
	EnvTypeTerminal EnvironmentType = "terminal"
	// EnvTypeDesktop 桌面环境
	EnvTypeDesktop EnvironmentType = "desktop"
	// EnvTypeDataProcess 数据处理
	EnvTypeDataProcess EnvironmentType = "data_process"
)

// GPUMode GPU 分配模式
type GPUMode string

const (
	// GPUModeExclusive 独占模式
	GPUModeExclusive GPUMode = "exclusive"
	// GPUModeVGPU vGPU 虚拟化
	GPUModeVGPU GPUMode = "vgpu"
	// GPUModeMIG MIG 分片
	GPUModeMIG GPUMode = "mig"
	// GPUModeShared 时间片共享
	GPUModeShared GPUMode = "shared"
	// GPUModeNone 无 GPU
	GPUModeNone GPUMode = "none"
)

// StorageType 存储类型
type StorageType string

const (
	// StorageLocal 本地存储
	StorageLocal StorageType = "local"
	// StorageNFS NFS 网络存储
	StorageNFS StorageType = "nfs"
	// StorageCeph Ceph 分布式存储
	StorageCeph StorageType = "ceph"
	// StorageS3 S3 对象存储
	StorageS3 StorageType = "s3"
	// StoragePVC K8S PVC
	StoragePVC StorageType = "pvc"
	// StorageJuiceFS JuiceFS 分布式文件系统
	StorageJuiceFS StorageType = "juicefs"
)

// LifecyclePolicy 生命周期策略
type LifecyclePolicy string

const (
	// LifecyclePersistent 持久化运行
	LifecyclePersistent LifecyclePolicy = "persistent"
	// LifecycleEphemeral 临时环境
	LifecycleEphemeral LifecyclePolicy = "ephemeral"
	// LifecycleScheduled 定时任务
	LifecycleScheduled LifecyclePolicy = "scheduled"
	// LifecycleOnDemand 按需启动
	LifecycleOnDemand LifecyclePolicy = "on_demand"
)

// StorageConfig 存储配置结构
type StorageConfig struct {
	// NFS 配置
	NFS *NFSConfig `json:"nfs,omitempty"`
	// S3 配置
	S3 *S3Config `json:"s3,omitempty"`
	// JuiceFS 配置
	JuiceFS *JuiceFSConfig `json:"juicefs,omitempty"`
}

// NFSConfig NFS 存储配置
type NFSConfig struct {
	Server    string   `json:"server"`              // NFS 服务器地址
	Path      string   `json:"path"`                // 挂载路径
	MountPath string   `json:"mount_path"`          // 容器内挂载点
	Options   []string `json:"options,omitempty"`   // 挂载选项
	ReadOnly  bool     `json:"read_only,omitempty"` // 只读模式
}

// S3Config S3 对象存储配置
type S3Config struct {
	Endpoint        string `json:"endpoint"`                   // S3 端点
	AccessKey       string `json:"access_key"`                 // Access Key
	SecretKey       string `json:"secret_key"`                 // Secret Key
	Bucket          string `json:"bucket"`                     // 存储桶名称
	Region          string `json:"region,omitempty"`           // 区域
	UseSSL          bool   `json:"use_ssl,omitempty"`          // 使用 SSL
	PathStyleAccess bool   `json:"path_style_access,omitempty"` // 路径风格访问
}

// JuiceFSConfig JuiceFS 配置
type JuiceFSConfig struct {
	Name         string    `json:"name"`                    // 文件系统名称
	MetaURL      string    `json:"meta_url"`                // 元数据引擎 URL
	Storage      string    `json:"storage"`                 // 对象存储类型(s3, minio, oss)
	Bucket       string    `json:"bucket"`                  // 存储桶
	AccessKey    string    `json:"access_key,omitempty"`    // Access Key
	SecretKey    string    `json:"secret_key,omitempty"`    // Secret Key
	MountPath    string    `json:"mount_path"`              // 容器内挂载点
	CacheDir     string    `json:"cache_dir,omitempty"`     // 缓存目录
	CacheSize    int64     `json:"cache_size,omitempty"`    // 缓存大小(MB)
	S3Config     *S3Config `json:"s3_config,omitempty"`     // S3 配置
}

// NetworkConfig 网络配置结构
type NetworkConfig struct {
	PublicAccess bool              `json:"public_access,omitempty"` // 公网访问
	Domain       string            `json:"domain,omitempty"`        // 域名
	SSLEnabled   bool              `json:"ssl_enabled,omitempty"`   // SSL 启用
	Ports        map[string]int    `json:"ports,omitempty"`         // 端口映射
	ExtraPorts   []PortConfig      `json:"extra_ports,omitempty"`   // 额外端口
}

// PortConfig 端口配置(用于网络配置)
type PortConfig struct {
	Name         string `json:"name"`                    // 服务名称
	InternalPort int    `json:"internal_port"`           // 容器内端口
	ExternalPort int    `json:"external_port"`           // 对外端口
	Protocol     string `json:"protocol,omitempty"`      // 协议(tcp/udp)
	Description  string `json:"description,omitempty"`   // 描述
}
