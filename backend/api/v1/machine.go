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

// CreateMachineRequest 创建机器请求（基础信息 + 登录信息）
type CreateMachineRequest struct {
	Name        string `json:"name" binding:"required"`
	Hostname    string `json:"hostname"`
	Region      string `json:"region" binding:"required"`
	IPAddress   string `json:"ip_address"`
	PublicIP    string `json:"public_ip"`
	SSHPort     int    `json:"ssh_port" binding:"required"`
	SSHUsername string `json:"ssh_username" binding:"required"`
	SSHPassword string `json:"ssh_password"`
	SSHKey      string `json:"ssh_key"`
}

// UpdateMachineRequest 更新机器请求
type UpdateMachineRequest struct {
	Name        string `json:"name"`
	Region      string `json:"region"`
	PublicIP    string `json:"public_ip"`
	SSHPort     int    `json:"ssh_port"`
	SSHUsername string `json:"ssh_username"`
	SSHPassword string `json:"ssh_password"`
	SSHKey      string `json:"ssh_key"`
}

// ImportMachineItem 批量导入机器条目
type ImportMachineItem struct {
	HostIP      string `json:"host_ip" binding:"required"`
	SSHPort     int    `json:"ssh_port" binding:"required"`
	Region      string `json:"region" binding:"required"`
	GPUModel    string `json:"gpu_model" binding:"required"`
	GPUCount    int    `json:"gpu_count" binding:"required"`
	CPUCores    int    `json:"cpu_cores" binding:"required"`
	RAMSize     int    `json:"ram_size" binding:"required"`     // GB
	DiskSize    int    `json:"disk_size" binding:"required"`    // GB
	PriceHourly int    `json:"price_hourly" binding:"required"` // cents
}

// ImportMachineRequest 批量导入机器请求
type ImportMachineRequest struct {
	Machines []ImportMachineItem `json:"machines" binding:"required,dive"`
}

// CreateMachineEnrollmentRequest 用户添加机器请求
type CreateMachineEnrollmentRequest struct {
	Name        string `json:"name"`
	Hostname    string `json:"hostname"`
	Region      string `json:"region" binding:"required"`
	IPAddress   string `json:"ip_address"`
	SSHPort     int    `json:"ssh_port" binding:"required"`
	SSHUsername string `json:"ssh_username" binding:"required"`
	SSHPassword string `json:"ssh_password"`
	SSHKey      string `json:"ssh_key"`
}
