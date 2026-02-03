package connector

import "fmt"

// SSHConfig SSH 访问配置
type SSHConfig struct {
	// 基础配置
	Port        int    `json:"port"`         // SSH 端口
	Username    string `json:"username"`     // 用户名
	Password    string `json:"password"`     // 密码
	PrivateKey  string `json:"private_key"`  // 私钥

	// 访问方式
	AccessType  AccessType `json:"access_type"` // 访问方式

	// Jumpserver 配置
	JumpserverAssetID string `json:"jumpserver_asset_id"` // Jumpserver 资产 ID

	// Guacamole 配置
	GuacamoleConnID   string `json:"guacamole_conn_id"`   // Guacamole 连接 ID
}

// GetProtocol 获取访问协议
func (c *SSHConfig) GetProtocol() Protocol {
	return ProtocolSSH
}

// GetAccessType 获取访问方式
func (c *SSHConfig) GetAccessType() AccessType {
	if c.AccessType == "" {
		return AccessTypeDirect
	}
	return c.AccessType
}

// Validate 验证配置
func (c *SSHConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("SSH 端口必须大于 0")
	}
	if c.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if c.Password == "" && c.PrivateKey == "" {
		return fmt.Errorf("密码和私钥至少需要提供一个")
	}
	return nil
}

// GenerateAccessInfo 生成访问信息
func (c *SSHConfig) GenerateAccessInfo(envID string, host *HostInfo) *AccessInfo {
	info := &AccessInfo{
		Protocol:   ProtocolSSH,
		AccessType: c.GetAccessType(),
		Credentials: &Credentials{
			Username:   c.Username,
			Password:   c.Password,
			PrivateKey: c.PrivateKey,
		},
	}

	// 生成内网访问地址
	if host.InternalIP != "" {
		info.InternalURL = fmt.Sprintf("ssh://%s@%s:%d", c.Username, host.InternalIP, c.Port)
		info.Command = fmt.Sprintf("ssh %s@%s -p %d", c.Username, host.InternalIP, c.Port)
	}

	// 生成公网访问地址
	if host.PublicDomain != "" {
		info.PublicURL = fmt.Sprintf("ssh://%s@%s:%d", c.Username, host.PublicDomain, c.Port)
	} else if host.PublicIP != "" {
		info.PublicURL = fmt.Sprintf("ssh://%s@%s:%d", c.Username, host.PublicIP, c.Port)
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

	info.Description = "SSH 终端访问"
	return info
}
