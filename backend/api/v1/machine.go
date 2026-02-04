package v1

// AllocateRequest 分配机器请求
type AllocateRequest struct {
	CustomerID     uint   `json:"customer_id" binding:"required"`
	HostID         string `json:"host_id" binding:"required"`
	DurationMonths int    `json:"duration_months" binding:"required,min=1"`
	Remark         string `json:"remark"`
}

// ReclaimRequest 回收机器请求
type ReclaimRequest struct {
	Reason string `json:"reason"`
	Force  bool   `json:"force"`
}

// CreateMachineRequest 创建机器请求
type CreateMachineRequest struct {
	HostIP      string `json:"host_ip" binding:"required"`
	SSHPort     int    `json:"ssh_port" binding:"required"`
	Region      string `json:"region" binding:"required"`
	GPUModel    string `json:"gpu_model" binding:"required"`
	GPUCount    int    `json:"gpu_count" binding:"required"`
	CPUCores    int    `json:"cpu_cores" binding:"required"`
	RAMSize     int    `json:"ram_size" binding:"required"` // GB
	DiskSize    int    `json:"disk_size" binding:"required"` // GB
	PriceHourly int    `json:"price_hourly" binding:"required"` // cents
}

// ImportMachineRequest 批量导入机器请求
type ImportMachineRequest struct {
	Machines []CreateMachineRequest `json:"machines" binding:"required,dive"`
}
