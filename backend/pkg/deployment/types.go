package deployment

// DeploymentType 部署类型
type DeploymentType string

const (
	// DeploymentTypeKubernetes Kubernetes 部署
	DeploymentTypeKubernetes DeploymentType = "kubernetes"

	// DeploymentTypeDocker Docker 部署
	DeploymentTypeDocker DeploymentType = "docker"

	// DeploymentTypeVM 虚拟机部署
	DeploymentTypeVM DeploymentType = "vm"
)

// DeploymentStatus 部署状态
type DeploymentStatus string

const (
	// DeploymentStatusPending 待部署
	DeploymentStatusPending DeploymentStatus = "pending"

	// DeploymentStatusDeploying 部署中
	DeploymentStatusDeploying DeploymentStatus = "deploying"

	// DeploymentStatusRunning 运行中
	DeploymentStatusRunning DeploymentStatus = "running"

	// DeploymentStatusFailed 部署失败
	DeploymentStatusFailed DeploymentStatus = "failed"

	// DeploymentStatusStopped 已停止
	DeploymentStatusStopped DeploymentStatus = "stopped"
)

// DeploymentConfig 部署配置接口
type DeploymentConfig interface {
	GetType() DeploymentType
	Validate() error
	GetResourceRequirements() *ResourceRequirements
}

// ResourceRequirements 资源需求
type ResourceRequirements struct {
	CPUCores   int    `json:"cpu_cores"`
	MemoryGB   int    `json:"memory_gb"`
	GPUCount   int    `json:"gpu_count"`
	GPUModel   string `json:"gpu_model"`
	StorageGB  int    `json:"storage_gb"`
	NetworkMbps int   `json:"network_mbps"`
}
