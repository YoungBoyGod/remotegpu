package v1

// HeartbeatRequest Agent 心跳请求
type HeartbeatRequest struct {
	AgentID   string           `json:"agent_id" binding:"required"`
	MachineID string           `json:"machine_id" binding:"required"`
	Metrics   *HeartbeatMetrics `json:"metrics,omitempty"`
}

// HeartbeatMetrics 心跳携带的监控指标
type HeartbeatMetrics struct {
	CPUUsagePercent    *float64          `json:"cpu_usage_percent,omitempty"`
	MemoryUsagePercent *float64          `json:"memory_usage_percent,omitempty"`
	MemoryUsedGB       *int64            `json:"memory_used_gb,omitempty"`
	DiskUsagePercent   *float64          `json:"disk_usage_percent,omitempty"`
	DiskUsedGB         *int64            `json:"disk_used_gb,omitempty"`
	GPUMetrics         []GPUMetric       `json:"gpu_metrics,omitempty"`
}

// GPUMetric 单个 GPU 的监控指标
type GPUMetric struct {
	Index          int      `json:"index"`
	UUID           string   `json:"uuid,omitempty"`
	Name           string   `json:"name,omitempty"`
	UtilPercent    *float64 `json:"util_percent,omitempty"`
	MemoryUsedMB   *int     `json:"memory_used_mb,omitempty"`
	MemoryTotalMB  *int     `json:"memory_total_mb,omitempty"`
	TemperatureC   *int     `json:"temperature_c,omitempty"`
	PowerUsageW    *float64 `json:"power_usage_w,omitempty"`
}

// RegisterRequest Agent 注册请求
type RegisterRequest struct {
	AgentID    string `json:"agent_id" binding:"required"`
	MachineID  string `json:"machine_id" binding:"required"`
	Version    string `json:"version,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	IPAddress  string `json:"ip_address,omitempty"`
	AgentPort  int    `json:"agent_port,omitempty"`
	MaxWorkers int    `json:"max_workers,omitempty"`
}
