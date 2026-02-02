package service

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/network"
	"github.com/YoungBoyGod/remotegpu/pkg/remote"
	"github.com/YoungBoyGod/remotegpu/pkg/security"
)

// AccessInfoService 连接信息服务
type AccessInfoService struct {
	// 模块化管理器
	networkManager  *network.NetworkManager
	remoteManager   *remote.AccessManager
	securityManager *security.SecurityManager
}

// NewAccessInfoService 创建连接信息服务
func NewAccessInfoService() *AccessInfoService {
	return &AccessInfoService{
		// 模块化管理器
		networkManager: network.NewNetworkManager(
			network.FirewallTypeIPTables,
			network.DNSProviderCloudflare,
			"remotegpu.com",
		),
		remoteManager:   remote.NewAccessManager(),
		securityManager: security.NewSecurityManager("jwt-secret-key"),
	}
}

// AccessInfo 连接信息结构
type AccessInfo struct {
	// 访问类型
	AccessType string `json:"access_type"` // internal/public/jumpserver/guacamole

	// 基础服务访问信息
	SSH     *SSHAccessInfo     `json:"ssh,omitempty"`
	RDP     *RDPAccessInfo     `json:"rdp,omitempty"`
	Jupyter *JupyterAccessInfo `json:"jupyter,omitempty"`
	VNC     *VNCAccessInfo     `json:"vnc,omitempty"`

	// Jumpserver 访问信息
	Jumpserver *JumpserverAccessInfo `json:"jumpserver,omitempty"`

	// Guacamole 访问信息
	Guacamole *GuacamoleAccessInfo `json:"guacamole,omitempty"`
}

// JumpserverAccessInfo Jumpserver 访问信息
type JumpserverAccessInfo struct {
	URL       string   `json:"url"`        // Jumpserver 访问地址
	AssetIDs  []string `json:"asset_ids"`  // 资产 ID 列表
	Available bool     `json:"available"`  // 是否可用
}

// GuacamoleAccessInfo Guacamole 访问信息
type GuacamoleAccessInfo struct {
	URL            string   `json:"url"`             // Guacamole 访问地址
	ConnectionIDs  []string `json:"connection_ids"`  // 连接 ID 列表
	Available      bool     `json:"available"`       // 是否可用
}

// SSHAccessInfo SSH连接信息
type SSHAccessInfo struct {
	// 内网访问
	InternalHost string `json:"internal_host"`
	InternalPort int    `json:"internal_port"`
	InternalURL  string `json:"internal_url"`

	// 公网访问
	PublicHost   string `json:"public_host,omitempty"`
	PublicPort   int    `json:"public_port,omitempty"`
	PublicDomain string `json:"public_domain,omitempty"`
	PublicURL    string `json:"public_url,omitempty"`

	// 认证信息
	Username string `json:"username"`
	Password string `json:"password"`
	Command  string `json:"command"`
}

// RDPAccessInfo RDP连接信息
type RDPAccessInfo struct {
	// 内网访问
	InternalHost string `json:"internal_host"`
	InternalPort int    `json:"internal_port"`
	InternalURL  string `json:"internal_url"`

	// 公网访问
	PublicHost   string `json:"public_host,omitempty"`
	PublicPort   int    `json:"public_port,omitempty"`
	PublicDomain string `json:"public_domain,omitempty"`
	PublicURL    string `json:"public_url,omitempty"`

	// 认证信息
	Username string `json:"username"`
	Password string `json:"password"`
	Command  string `json:"command"`
}

// JupyterAccessInfo Jupyter连接信息
type JupyterAccessInfo struct {
	// 内网访问
	InternalURL string `json:"internal_url"`

	// 公网访问
	PublicURL    string `json:"public_url,omitempty"`
	PublicDomain string `json:"public_domain,omitempty"`

	// 认证信息
	Token    string `json:"token"`
	Password string `json:"password"`
}

// VNCAccessInfo VNC连接信息
type VNCAccessInfo struct {
	// 内网访问
	InternalHost string `json:"internal_host"`
	InternalPort int    `json:"internal_port"`
	InternalURL  string `json:"internal_url"`

	// 公网访问
	PublicHost   string `json:"public_host,omitempty"`
	PublicPort   int    `json:"public_port,omitempty"`
	PublicDomain string `json:"public_domain,omitempty"`
	PublicURL    string `json:"public_url,omitempty"`

	// 认证信息
	Password string `json:"password"`
}

// GenerateAccessInfo 生成环境的连接信息
func (s *AccessInfoService) GenerateAccessInfo(env *entity.Environment, host *entity.Host) (*AccessInfo, error) {
	accessInfo := &AccessInfo{}

	// 确定访问类型
	if env.UseGuacamole {
		accessInfo.AccessType = "guacamole"
	} else if env.UseJumpserver {
		accessInfo.AccessType = "jumpserver"
	} else {
		accessInfo.AccessType = "internal" // 默认内网访问
	}

	// 使用 networkManager 获取环境的所有端口映射
	portMappings, err := s.networkManager.GetPortMappings(env.ID)
	if err != nil {
		return nil, fmt.Errorf("获取端口映射失败: %w", err)
	}

	// 根据端口映射生成各服务的访问信息
	for _, pm := range portMappings {
		switch pm.ServiceType {
		case network.ServiceTypeSSH:
			accessInfo.SSH = s.generateSSHInfo(env, host, pm)
		case network.ServiceTypeRDP:
			accessInfo.RDP = s.generateRDPInfo(env, host, pm)
		case network.ServiceTypeJupyter:
			accessInfo.Jupyter = s.generateJupyterInfo(env, host, pm)
		case network.ServiceTypeVNC:
			accessInfo.VNC = s.generateVNCInfo(env, host, pm)
		}
	}

	// 如果启用了 Jumpserver,生成 Jumpserver 访问信息
	if env.UseJumpserver {
		// TODO: 实现 Jumpserver 访问信息生成
		accessInfo.Jumpserver = &JumpserverAccessInfo{
			URL:       "TODO: 从配置获取 Jumpserver URL",
			AssetIDs:  []string{}, // TODO: 从环境获取资产 ID
			Available: false,      // TODO: 检查 Jumpserver 是否可用
		}
	}

	// 如果启用了 Guacamole,生成 Guacamole 访问信息
	if env.UseGuacamole {
		// TODO: 实现 Guacamole 访问信息生成
		accessInfo.Guacamole = &GuacamoleAccessInfo{
			URL:           "TODO: 从配置获取 Guacamole URL",
			ConnectionIDs: []string{env.GuacamoleConnID}, // 从环境获取连接 ID
			Available:     env.GuacamoleConnID != "",
		}
	}

	return accessInfo, nil
}

// generateSSHInfo 生成 SSH 连接信息
func (s *AccessInfoService) generateSSHInfo(env *entity.Environment, host *entity.Host, pm *network.PortMapping) *SSHAccessInfo {
	username := "root"
	// 使用 securityManager 生成密码
	password, _ := s.securityManager.GeneratePassword(security.PasswordStrengthStrong)

	info := &SSHAccessInfo{
		// 内网访问
		InternalHost: host.IPAddress,
		InternalPort: pm.ExternalPort,
		InternalURL:  fmt.Sprintf("ssh://%s:%d", host.IPAddress, pm.ExternalPort),

		// 认证信息
		Username: username,
		Password: password,
		Command:  fmt.Sprintf("ssh %s@%s -p %d", username, host.IPAddress, pm.ExternalPort),
	}

	// 如果有公网访问配置
	if pm.PublicPort != 0 {
		info.PublicPort = pm.PublicPort
		info.PublicURL = fmt.Sprintf("ssh://%s:%d", host.IPAddress, pm.PublicPort)
	}

	return info
}

// generateRDPInfo 生成 RDP 连接信息
func (s *AccessInfoService) generateRDPInfo(env *entity.Environment, host *entity.Host, pm *network.PortMapping) *RDPAccessInfo {
	username := "Administrator"
	// 使用 securityManager 生成密码
	password, _ := s.securityManager.GeneratePassword(security.PasswordStrengthStrong)

	info := &RDPAccessInfo{
		// 内网访问
		InternalHost: host.IPAddress,
		InternalPort: pm.ExternalPort,
		InternalURL:  fmt.Sprintf("rdp://%s:%d", host.IPAddress, pm.ExternalPort),

		// 认证信息
		Username: username,
		Password: password,
		Command:  fmt.Sprintf("mstsc /v:%s:%d", host.IPAddress, pm.ExternalPort),
	}

	// 如果有公网访问配置
	if pm.PublicPort != 0 {
		info.PublicPort = pm.PublicPort
		info.PublicURL = fmt.Sprintf("rdp://%s:%d", host.IPAddress, pm.PublicPort)
	}

	return info
}

// generateJupyterInfo 生成 Jupyter 连接信息
func (s *AccessInfoService) generateJupyterInfo(env *entity.Environment, host *entity.Host, pm *network.PortMapping) *JupyterAccessInfo {
	// 使用 securityManager 生成 API Key 作为 token
	token, _ := s.securityManager.GenerateAPIKey("jupyter", 32)

	info := &JupyterAccessInfo{
		// 内网访问
		InternalURL: fmt.Sprintf("http://%s:%d", host.IPAddress, pm.ExternalPort),

		// 认证信息
		Token:    token.Value,
		Password: "",
	}

	// 如果有公网访问配置
	if pm.PublicPort != 0 {
		info.PublicURL = fmt.Sprintf("https://%s:%d", host.IPAddress, pm.PublicPort)
	}

	return info
}

// generateVNCInfo 生成 VNC 连接信息
func (s *AccessInfoService) generateVNCInfo(env *entity.Environment, host *entity.Host, pm *network.PortMapping) *VNCAccessInfo {
	// 使用 securityManager 生成密码
	password, _ := s.securityManager.GeneratePassword(security.PasswordStrengthStrong)

	info := &VNCAccessInfo{
		// 内网访问
		InternalHost: host.IPAddress,
		InternalPort: pm.ExternalPort,
		InternalURL:  fmt.Sprintf("vnc://%s:%d", host.IPAddress, pm.ExternalPort),

		// 认证信息
		Password: password,
	}

	// 如果有公网访问配置
	if pm.PublicPort != 0 {
		info.PublicPort = pm.PublicPort
		info.PublicURL = fmt.Sprintf("vnc://%s:%d", host.IPAddress, pm.PublicPort)
	}

	return info
}
