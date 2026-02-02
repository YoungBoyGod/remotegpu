package remote

import "fmt"

// VNCConfig VNC 访问配置
type VNCConfig struct {
	// 基础配置
	Port        int    `json:"port"`         // VNC 端口
	Password    string `json:"password"`     // VNC 密码

	// 访问方式
	AccessType  AccessType `json:"access_type"` // 访问方式

	// Jumpserver 配置
	JumpserverAssetID string `json:"jumpserver_asset_id"`

	// Guacamole 配置
	GuacamoleConnID   string `json:"guacamole_conn_id"`
}

// GetProtocol 获取访问协议
func (c *VNCConfig) GetProtocol() Protocol {
	return ProtocolVNC
}

// GetAccessType 获取访问方式
func (c *VNCConfig) GetAccessType() AccessType {
	if c.AccessType == "" {
		return AccessTypeDirect
	}
	return c.AccessType
}

// Validate 验证配置
func (c *VNCConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("VNC 端口必须大于 0")
	}
	return nil
}

// GenerateAccessInfo 生成访问信息
func (c *VNCConfig) GenerateAccessInfo(envID string, host *HostInfo) *AccessInfo {
	info := &AccessInfo{
		Protocol:   ProtocolVNC,
		AccessType: c.GetAccessType(),
		Credentials: &Credentials{
			Password: c.Password,
		},
	}

	// 生成内网访问地址
	if host.InternalIP != "" {
		info.InternalURL = fmt.Sprintf("vnc://%s:%d", host.InternalIP, c.Port)
	}

	// 生成公网访问地址
	if host.PublicDomain != "" {
		info.PublicURL = fmt.Sprintf("vnc://%s:%d", host.PublicDomain, c.Port)
	} else if host.PublicIP != "" {
		info.PublicURL = fmt.Sprintf("vnc://%s:%d", host.PublicIP, c.Port)
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

	info.Description = "VNC 远程桌面访问"
	return info
}
