package network

// PortRange 端口范围
type PortRange struct {
	Start int `json:"start"` // 起始端口
	End   int `json:"end"`   // 结束端口
}

// ServiceType 服务类型
type ServiceType string

const (
	// ServiceTypeSSH SSH 服务
	ServiceTypeSSH ServiceType = "ssh"

	// ServiceTypeRDP RDP 服务
	ServiceTypeRDP ServiceType = "rdp"

	// ServiceTypeVNC VNC 服务
	ServiceTypeVNC ServiceType = "vnc"

	// ServiceTypeJupyter Jupyter 服务
	ServiceTypeJupyter ServiceType = "jupyter"

	// ServiceTypeTensorBoard TensorBoard 服务
	ServiceTypeTensorBoard ServiceType = "tensorboard"

	// ServiceTypeCustom 自定义服务
	ServiceTypeCustom ServiceType = "custom"
)

// PortMapping 端口映射
type PortMapping struct {
	ID           int64       `json:"id"`
	EnvID        string      `json:"env_id"`
	ServiceType  ServiceType `json:"service_type"`
	InternalPort int         `json:"internal_port"`
	ExternalPort int         `json:"external_port"`
	PublicPort   int         `json:"public_port"`
	Protocol     string      `json:"protocol"`
	Description  string      `json:"description"`
}

// FirewallRule 防火墙规则
type FirewallRule struct {
	ID          string `json:"id"`
	SourceIP    string `json:"source_ip"`
	DestIP      string `json:"dest_ip"`
	SourcePort  int    `json:"source_port"`
	DestPort    int    `json:"dest_port"`
	Protocol    string `json:"protocol"`
	Action      string `json:"action"` // allow/deny
	Description string `json:"description"`
}

// DNSRecord DNS 记录
type DNSRecord struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
	Type   string `json:"type"`   // A/AAAA/CNAME
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}
