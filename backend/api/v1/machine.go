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
	Name         string `json:"name" binding:"required"`
	Hostname     string `json:"hostname"`
	Region       string `json:"region" binding:"required"`
	IPAddress    string `json:"ip_address"`
	PublicIP     string `json:"public_ip"`
	SSHHost      string `json:"ssh_host"`
	SSHPort      int    `json:"ssh_port" binding:"required"`
	SSHUsername  string `json:"ssh_username" binding:"required"`
	SSHPassword  string `json:"ssh_password"`
	SSHKey       string `json:"ssh_key"`
	JupyterURL   string `json:"jupyter_url"`
	JupyterToken string `json:"jupyter_token"`
	VNCURL       string `json:"vnc_url"`
	VNCPassword  string `json:"vnc_password"`
	// 外映射配置
	ExternalIP          string `json:"external_ip"`
	ExternalSSHPort     int    `json:"external_ssh_port"`
	ExternalJupyterPort int    `json:"external_jupyter_port"`
	ExternalVNCPort     int    `json:"external_vnc_port"`
	NginxDomain         string `json:"nginx_domain"`
	NginxConfigPath     string `json:"nginx_config_path"`
}

// UpdateMachineRequest 更新机器请求
type UpdateMachineRequest struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	PublicIP     string `json:"public_ip"`
	SSHHost      string `json:"ssh_host"`
	SSHPort      int    `json:"ssh_port"`
	SSHUsername  string `json:"ssh_username"`
	SSHPassword  string `json:"ssh_password"`
	SSHKey       string `json:"ssh_key"`
	JupyterURL   string `json:"jupyter_url"`
	JupyterToken string `json:"jupyter_token"`
	VNCURL       string `json:"vnc_url"`
	VNCPassword  string `json:"vnc_password"`
	// 外映射配置
	ExternalIP          string `json:"external_ip"`
	ExternalSSHPort     int    `json:"external_ssh_port"`
	ExternalJupyterPort int    `json:"external_jupyter_port"`
	ExternalVNCPort     int    `json:"external_vnc_port"`
	NginxDomain         string `json:"nginx_domain"`
	NginxConfigPath     string `json:"nginx_config_path"`
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

// BatchSetMaintenanceRequest 批量启用/禁用机器请求
type BatchSetMaintenanceRequest struct {
	HostIDs     []string `json:"host_ids" binding:"required,min=1"`
	Maintenance bool     `json:"maintenance"`
}

// BatchAllocateRequest 批量分配机器请求
type BatchAllocateRequest struct {
	HostIDs        []string `json:"host_ids" binding:"required,min=1"`
	CustomerID     uint     `json:"customer_id" binding:"required"`
	DurationMonths int      `json:"duration_months" binding:"required,min=1"`
	Remark         string   `json:"remark"`
}

// BatchReclaimRequest 批量回收机器请求
type BatchReclaimRequest struct {
	HostIDs []string `json:"host_ids" binding:"required,min=1"`
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
