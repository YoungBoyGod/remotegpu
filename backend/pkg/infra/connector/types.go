package connector

// Protocol 访问协议
type Protocol string

const (
	// ProtocolSSH SSH 协议
	ProtocolSSH Protocol = "ssh"

	// ProtocolRDP RDP 协议
	ProtocolRDP Protocol = "rdp"

	// ProtocolVNC VNC 协议
	ProtocolVNC Protocol = "vnc"

	// ProtocolHTTP HTTP 协议
	ProtocolHTTP Protocol = "http"

	// ProtocolHTTPS HTTPS 协议
	ProtocolHTTPS Protocol = "https"

	// ProtocolWebSocket WebSocket 协议
	ProtocolWebSocket Protocol = "websocket"
)

// AccessType 访问方式
type AccessType string

const (
	// AccessTypeDirect 直接访问
	AccessTypeDirect AccessType = "direct"

	// AccessTypeJumpserver 堡垒机访问
	AccessTypeJumpserver AccessType = "jumpserver"

	// AccessTypeGuacamole Guacamole Web 访问
	AccessTypeGuacamole AccessType = "guacamole"

	// AccessTypeProxy 代理访问
	AccessTypeProxy AccessType = "proxy"

	// AccessTypeVPN VPN 访问
	AccessTypeVPN AccessType = "vpn"
)

// AccessConfig 访问配置接口
type AccessConfig interface {
	// GetProtocol 获取访问协议
	GetProtocol() Protocol

	// GetAccessType 获取访问方式
	GetAccessType() AccessType

	// Validate 验证配置
	Validate() error

	// GenerateAccessInfo 生成访问信息
	GenerateAccessInfo(envID string, host *HostInfo) *AccessInfo
}

// AccessInfo 访问信息
type AccessInfo struct {
	Protocol      Protocol    `json:"protocol"`       // 访问协议
	AccessType    AccessType  `json:"access_type"`    // 访问方式
	InternalURL   string      `json:"internal_url"`   // 内网访问地址
	PublicURL     string      `json:"public_url"`     // 公网访问地址
	WebURL        string      `json:"web_url"`        // Web 访问地址 (Jumpserver/Guacamole)
	Credentials   *Credentials `json:"credentials"`   // 访问凭证
	Command       string      `json:"command"`        // 访问命令 (如 SSH 命令)
	Description   string      `json:"description"`    // 描述信息
}

// HostInfo 主机信息
type HostInfo struct {
	ID            string `json:"id"`              // 主机 ID
	InternalIP    string `json:"internal_ip"`     // 内网 IP
	PublicIP      string `json:"public_ip"`       // 公网 IP
	PublicDomain  string `json:"public_domain"`   // 公网域名
	Port          int    `json:"port"`            // 端口
}

// Credentials 访问凭证
type Credentials struct {
	Username      string `json:"username"`        // 用户名
	Password      string `json:"password"`        // 密码
	PrivateKey    string `json:"private_key"`     // 私钥
	Token         string `json:"token"`           // 访问令牌
}

// PortMapping 端口映射信息
type PortMapping struct {
	InternalPort  int    `json:"internal_port"`   // 内部端口
	ExternalPort  int    `json:"external_port"`   // 外部端口
	PublicPort    int    `json:"public_port"`     // 公网端口
	Protocol      string `json:"protocol"`        // 协议 (tcp/udp)
}
