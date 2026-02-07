package collector

// Metrics 采集到的系统和 GPU 指标
type Metrics struct {
	CPUUsagePercent    *float64    `json:"cpu_usage_percent,omitempty"`
	MemoryUsagePercent *float64    `json:"memory_usage_percent,omitempty"`
	MemoryUsedGB       *int64      `json:"memory_used_gb,omitempty"`
	DiskUsagePercent   *float64    `json:"disk_usage_percent,omitempty"`
	DiskUsedGB         *int64      `json:"disk_used_gb,omitempty"`
	GPUMetrics         []GPUMetric `json:"gpu_metrics,omitempty"`
}

// GPUMetric 单个 GPU 的监控指标
type GPUMetric struct {
	Index         int      `json:"index"`
	UUID          string   `json:"uuid,omitempty"`
	Name          string   `json:"name,omitempty"`
	UtilPercent   *float64 `json:"util_percent,omitempty"`
	MemoryUsedMB  *int     `json:"memory_used_mb,omitempty"`
	MemoryTotalMB *int     `json:"memory_total_mb,omitempty"`
	TemperatureC  *int     `json:"temperature_c,omitempty"`
	PowerUsageW   *float64 `json:"power_usage_w,omitempty"`
}

// Collector 指标采集器
type Collector struct{}

// NewCollector 创建采集器
func NewCollector() *Collector {
	return &Collector{}
}

// Collect 采集所有指标（系统 + GPU）
func (c *Collector) Collect() *Metrics {
	m := &Metrics{}
	collectSystem(m)
	collectGPU(m)
	return m
}
