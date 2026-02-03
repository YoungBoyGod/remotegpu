package deployment

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/pkg/infra/network"
	"github.com/YoungBoyGod/remotegpu/pkg/infra/volume"
)

// VMDeploymentConfig 虚拟机部署配置
type VMDeploymentConfig struct {
	VMName        string                    `json:"vm_name"`
	OSType        string                    `json:"os_type"`        // linux, windows
	OSImage       string                    `json:"os_image"`
	Resources     *ResourceRequirements     `json:"resources"`
	Volumes       []volume.VolumeConfig     `json:"volumes"`
	PortMappings  []*network.PortMapping    `json:"port_mappings"`
	CloudInit     string                    `json:"cloud_init"`
	SSHKeys       []string                  `json:"ssh_keys"`
	Hypervisor    string                    `json:"hypervisor"`     // kvm, vmware, virtualbox
	NetworkConfig string                    `json:"network_config"`
}

// GetType 获取部署类型
func (c *VMDeploymentConfig) GetType() DeploymentType {
	return DeploymentTypeVM
}

// Validate 验证配置
func (c *VMDeploymentConfig) Validate() error {
	if c.VMName == "" {
		return fmt.Errorf("vm_name 不能为空")
	}
	if c.OSType == "" {
		return fmt.Errorf("os_type 不能为空")
	}
	if c.OSImage == "" {
		return fmt.Errorf("os_image 不能为空")
	}
	if c.Resources == nil {
		return fmt.Errorf("resources 不能为空")
	}
	return nil
}

// GetResourceRequirements 获取资源需求
func (c *VMDeploymentConfig) GetResourceRequirements() *ResourceRequirements {
	return c.Resources
}
