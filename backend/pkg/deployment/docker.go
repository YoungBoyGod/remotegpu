package deployment

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/pkg/network"
	"github.com/YoungBoyGod/remotegpu/pkg/volume"
)

// DockerDeploymentConfig Docker 部署配置
type DockerDeploymentConfig struct {
	ContainerName string                    `json:"container_name"`
	Image         string                    `json:"image"`
	Command       []string                  `json:"command"`
	Args          []string                  `json:"args"`
	Resources     *ResourceRequirements     `json:"resources"`
	Volumes       []volume.VolumeConfig     `json:"volumes"`
	PortMappings  []*network.PortMapping    `json:"port_mappings"`
	Env           map[string]string         `json:"env"`
	Labels        map[string]string         `json:"labels"`
	RestartPolicy string                    `json:"restart_policy"`
	NetworkMode   string                    `json:"network_mode"`
	Privileged    bool                      `json:"privileged"`
	Runtime       string                    `json:"runtime"` // nvidia, runc
}

// GetType 获取部署类型
func (c *DockerDeploymentConfig) GetType() DeploymentType {
	return DeploymentTypeDocker
}

// Validate 验证配置
func (c *DockerDeploymentConfig) Validate() error {
	if c.ContainerName == "" {
		return fmt.Errorf("container_name 不能为空")
	}
	if c.Image == "" {
		return fmt.Errorf("image 不能为空")
	}
	if c.Resources == nil {
		return fmt.Errorf("resources 不能为空")
	}
	return nil
}

// GetResourceRequirements 获取资源需求
func (c *DockerDeploymentConfig) GetResourceRequirements() *ResourceRequirements {
	return c.Resources
}
