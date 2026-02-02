package deployment

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/pkg/network"
	"github.com/YoungBoyGod/remotegpu/pkg/volume"
)

// K8sDeploymentConfig Kubernetes 部署配置
type K8sDeploymentConfig struct {
	Namespace      string                    `json:"namespace"`
	Name           string                    `json:"name"`
	Image          string                    `json:"image"`
	Command        []string                  `json:"command"`
	Args           []string                  `json:"args"`
	Resources      *ResourceRequirements     `json:"resources"`
	Volumes        []volume.VolumeConfig     `json:"volumes"`
	PortMappings   []*network.PortMapping    `json:"port_mappings"`
	Env            map[string]string         `json:"env"`
	Labels         map[string]string         `json:"labels"`
	Annotations    map[string]string         `json:"annotations"`
	RestartPolicy  string                    `json:"restart_policy"`
	ServiceAccount string                    `json:"service_account"`
}

// GetType 获取部署类型
func (c *K8sDeploymentConfig) GetType() DeploymentType {
	return DeploymentTypeKubernetes
}

// Validate 验证配置
func (c *K8sDeploymentConfig) Validate() error {
	if c.Namespace == "" {
		return fmt.Errorf("namespace 不能为空")
	}
	if c.Name == "" {
		return fmt.Errorf("name 不能为空")
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
func (c *K8sDeploymentConfig) GetResourceRequirements() *ResourceRequirements {
	return c.Resources
}
