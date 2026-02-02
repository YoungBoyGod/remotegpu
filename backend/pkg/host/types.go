package host

// SelectStrategy 主机选择策略
type SelectStrategy string

const (
	// SelectStrategyRandom 随机选择
	SelectStrategyRandom SelectStrategy = "random"
	// SelectStrategyLeastUsed 选择资源使用率最低的主机
	SelectStrategyLeastUsed SelectStrategy = "least_used"
	// SelectStrategyMostFree 选择剩余资源最多的主机
	SelectStrategyMostFree SelectStrategy = "most_free"
	// SelectStrategyRoundRobin 轮询选择
	SelectStrategyRoundRobin SelectStrategy = "round_robin"
)

// ResourceRequirement 资源需求
type ResourceRequirement struct {
	CPU    int   `json:"cpu"`    // CPU 核心数
	Memory int64 `json:"memory"` // 内存（字节）
	GPU    int   `json:"gpu"`    // GPU 数量
}

// HostInfo 主机信息
type HostInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	IPAddress   string            `json:"ip_address"`
	Status      string            `json:"status"`
	TotalCPU    int               `json:"total_cpu"`
	UsedCPU     int               `json:"used_cpu"`
	TotalMemory int64             `json:"total_memory"`
	UsedMemory  int64             `json:"used_memory"`
	TotalGPU    int               `json:"total_gpu"`
	UsedGPU     int               `json:"used_gpu"`
	Labels      map[string]string `json:"labels"`
}

// AvailableResources 获取可用资源
func (h *HostInfo) AvailableResources() *ResourceRequirement {
	return &ResourceRequirement{
		CPU:    h.TotalCPU - h.UsedCPU,
		Memory: h.TotalMemory - h.UsedMemory,
		GPU:    h.TotalGPU - h.UsedGPU,
	}
}

// UsageRate 计算资源使用率
func (h *HostInfo) UsageRate() float64 {
	cpuRate := float64(h.UsedCPU) / float64(h.TotalCPU)
	memoryRate := float64(h.UsedMemory) / float64(h.TotalMemory)

	var gpuRate float64
	if h.TotalGPU > 0 {
		gpuRate = float64(h.UsedGPU) / float64(h.TotalGPU)
	}

	// 综合使用率（CPU 40%, Memory 40%, GPU 20%）
	return cpuRate*0.4 + memoryRate*0.4 + gpuRate*0.2
}

// CanAllocate 检查是否可以分配指定资源
func (h *HostInfo) CanAllocate(req *ResourceRequirement) bool {
	available := h.AvailableResources()
	return available.CPU >= req.CPU &&
		available.Memory >= req.Memory &&
		available.GPU >= req.GPU
}

// HostFilter 主机过滤器
type HostFilter func(*HostInfo) bool

// HostSelector 主机选择器接口
type HostSelector interface {
	// Select 选择主机
	Select(hosts []*HostInfo, req *ResourceRequirement) (*HostInfo, error)
}
