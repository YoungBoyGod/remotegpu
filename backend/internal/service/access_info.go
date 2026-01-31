package service

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

// AccessInfoService 连接信息服务
type AccessInfoService struct{}

// NewAccessInfoService 创建连接信息服务
func NewAccessInfoService() *AccessInfoService {
	return &AccessInfoService{}
}

// AccessInfo 连接信息结构
type AccessInfo struct {
	SSH     *SSHAccessInfo     `json:"ssh,omitempty"`
	RDP     *RDPAccessInfo     `json:"rdp,omitempty"`
	Jupyter *JupyterAccessInfo `json:"jupyter,omitempty"`
	VNC     *VNCAccessInfo     `json:"vnc,omitempty"`
}

// SSHAccessInfo SSH连接信息
type SSHAccessInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Command  string `json:"command"`
}

// RDPAccessInfo RDP连接信息
type RDPAccessInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Command  string `json:"command"`
}

// JupyterAccessInfo Jupyter连接信息
type JupyterAccessInfo struct {
	URL      string `json:"url"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

// VNCAccessInfo VNC连接信息
type VNCAccessInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

// GenerateAccessInfo 生成环境的连接信息
func (s *AccessInfoService) GenerateAccessInfo(env *entity.Environment, host *entity.Host) (*AccessInfo, error) {
	accessInfo := &AccessInfo{}

	// 根据端口配置生成不同类型的连接信息
	if env.SSHPort != nil && *env.SSHPort > 0 {
		accessInfo.SSH = s.generateSSHInfo(env, host)
	}

	if env.RDPPort != nil && *env.RDPPort > 0 {
		accessInfo.RDP = s.generateRDPInfo(env, host)
	}

	if env.JupyterPort != nil && *env.JupyterPort > 0 {
		accessInfo.Jupyter = s.generateJupyterInfo(env, host)
	}

	// 如果没有配置任何端口,提供VNC作为备选
	if accessInfo.SSH == nil && accessInfo.RDP == nil && accessInfo.Jupyter == nil {
		accessInfo.VNC = s.generateVNCInfo(env, host)
	}

	return accessInfo, nil
}

// generateSSHInfo 生成SSH连接信息
func (s *AccessInfoService) generateSSHInfo(env *entity.Environment, host *entity.Host) *SSHAccessInfo {
	// 调用第三方服务获取真实的连接信息
	// 这里模拟生成连接信息
	username := "root"
	password := s.generatePassword(env.ID)

	return &SSHAccessInfo{
		Host:     host.IPAddress,
		Port:     int(*env.SSHPort),
		Username: username,
		Password: password,
		Command:  fmt.Sprintf("ssh %s@%s -p %d", username, host.IPAddress, *env.SSHPort),
	}
}

// generateRDPInfo 生成RDP连接信息
func (s *AccessInfoService) generateRDPInfo(env *entity.Environment, host *entity.Host) *RDPAccessInfo {
	username := "Administrator"
	password := s.generatePassword(env.ID)

	return &RDPAccessInfo{
		Host:     host.IPAddress,
		Port:     int(*env.RDPPort),
		Username: username,
		Password: password,
		Command:  fmt.Sprintf("mstsc /v:%s:%d", host.IPAddress, *env.RDPPort),
	}
}

// generateJupyterInfo 生成Jupyter连接信息
func (s *AccessInfoService) generateJupyterInfo(env *entity.Environment, host *entity.Host) *JupyterAccessInfo {
	token := s.generateToken(env.ID)

	return &JupyterAccessInfo{
		URL:      fmt.Sprintf("http://%s:%d", host.IPAddress, *env.JupyterPort),
		Token:    token,
		Password: "",
	}
}

// generateVNCInfo 生成VNC连接信息
func (s *AccessInfoService) generateVNCInfo(env *entity.Environment, host *entity.Host) *VNCAccessInfo {
	password := s.generatePassword(env.ID)

	return &VNCAccessInfo{
		Host:     host.IPAddress,
		Port:     5900,
		Password: password,
		URL:      fmt.Sprintf("vnc://%s:5900", host.IPAddress),
	}
}

// generatePassword 生成密码(模拟调用第三方服务)
func (s *AccessInfoService) generatePassword(envID string) string {
	// 实际应该调用第三方服务生成安全密码
	// 这里简化处理,使用环境ID生成固定密码
	return fmt.Sprintf("Pass_%s", envID[:8])
}

// generateToken 生成Token
func (s *AccessInfoService) generateToken(envID string) string {
	// 实际应该调用第三方服务生成Token
	return fmt.Sprintf("token_%s", envID[:8])
}

