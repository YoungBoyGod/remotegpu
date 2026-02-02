package remote

import "fmt"

// RDPConfig RDP 访问配置
type RDPConfig struct {
	// 基础配置
	Port        int    `json:"port"`         // RDP 端口
	Username    string `json:"username"`     // 用户名
	Password    string `json:"password"`     // 密码

	// 访问方式
	AccessType  AccessType `json:"access_type"` // 访问方式

	// Jumpserver 配置
	JumpserverAssetID string `json:"jumpserver_asset_id"`

	// Guacamole 配置
	GuacamoleConnID   string `json:"guacamole_conn_id"`
}

// GetProtocol 获取访问协议
func (c *RDPConfig) GetProtocol() Protocol {
	return ProtocolRDP
}

// GetAccessType 获取访问方式
func (c *RDPConfig) GetAccessType() AccessType {
	if c.AccessType == "" {
		return AccessTypeDirect
	}
	return c.AccessType
}

// Validate 验证配置
func (c *RDPConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("RDP 端口必须大于 0")
	}
	if c.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if c.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return nil
}

// GenerateAccessInfo 生成访问信息
func (c *RDPConfig) GenerateAccessInfo(envID string, host *HostInfo) *AccessInfo {
	info := &AccessInfo{
		Protocol:   ProtocolRDP,
		AccessType: c.GetAccessType(),
		Credentials: &Credentials{
			Username: c.Username,
			Password: c.Password,
		},
	}

	// 生成内网访问地址
	if host.InternalIP != "" {
		info.InternalURL = fmt.Sprintf("rdp://%s:%d", host.InternalIP, c.Port)
		info.Command = fmt.Sprintf("mstsc /v:%s:%d", host.InternalIP, c.Port)
	}

	// 生成公网访问地址
	if host.PublicDomain != "" {
		info.PublicURL = fmt.Sprintf("rdp://%s:%d", host.PublicDomain, c.Port)
	} else if host.PublicIP != "" {
		info.PublicURL = fmt.Sprintf("rdp://%s:%d", host.PublicIP, c.Port)
	}

	// 生成 Web 访问地址
	switch c.GetAccessType() {
	case AccessTypeJumpserver:
		if c.JumpserverAssetID != "" {
			info.WebURL = fmt.Sprintf("/luna/?asset=%s", c.JumpserverAssetID)
		}
	case AccessTypeGuacamole:
		if c.GuacamoleConnID != "" {
			info.WebURL = fmt.Sprintf("/#/client/%s", c.GuacamoleConnID)
		}
	}

	info.Description = "RDP 远程桌面访问"
	return info
}
