# pkg/deployment 使用示例

## 1. Kubernetes 部署配置

```go
package main

import (
    "github.com/YoungBoyGod/remotegpu/pkg/deployment"
    "github.com/YoungBoyGod/remotegpu/pkg/network"
    "github.com/YoungBoyGod/remotegpu/pkg/volume"
)

func ExampleK8sDeployment() {
    // 创建 Kubernetes 部署配置
    k8sConfig := &deployment.K8sDeploymentConfig{
        Namespace: "remotegpu",
        Name:      "env-123",
        Image:     "nvidia/cuda:11.8.0-runtime-ubuntu22.04",
        Command:   []string{"/bin/bash"},
        Args:      []string{"-c", "sleep infinity"},
        Resources: &deployment.ResourceRequirements{
            CPUCores:  4,
            MemoryGB:  16,
            GPUCount:  1,
            GPUModel:  "RTX-3090",
            StorageGB: 100,
        },
        Env: map[string]string{
            "NVIDIA_VISIBLE_DEVICES": "all",
            "CUDA_VISIBLE_DEVICES":   "0",
        },
        Labels: map[string]string{
            "app":   "remotegpu",
            "env":   "env-123",
        },
        RestartPolicy: "Always",
    }

    // 验证配置
    if err := k8sConfig.Validate(); err != nil {
        panic(err)
    }

    // 获取资源需求
    resources := k8sConfig.GetResourceRequirements()
    // resources.GPUCount: 1
    // resources.GPUModel: "RTX-3090"
}
```

## 2. Docker 部署配置

```go
func ExampleDockerDeployment() {
    // 创建 Docker 部署配置
    dockerConfig := &deployment.DockerDeploymentConfig{
        ContainerName: "env-456",
        Image:         "nvidia/cuda:11.8.0-runtime-ubuntu22.04",
        Command:       []string{"/bin/bash"},
        Args:          []string{"-c", "sleep infinity"},
        Resources: &deployment.ResourceRequirements{
            CPUCores:  8,
            MemoryGB:  32,
            GPUCount:  2,
            GPUModel:  "A100",
            StorageGB: 500,
        },
        Env: map[string]string{
            "NVIDIA_VISIBLE_DEVICES": "all",
        },
        Labels: map[string]string{
            "app": "remotegpu",
            "env": "env-456",
        },
        RestartPolicy: "unless-stopped",
        Runtime:       "nvidia",
        Privileged:    false,
    }

    // 验证配置
    if err := dockerConfig.Validate(); err != nil {
        panic(err)
    }
}
```

## 3. VM 部署配置

```go
func ExampleVMDeployment() {
    // 创建 VM 部署配置
    vmConfig := &deployment.VMDeploymentConfig{
        VMName:  "env-789",
        OSType:  "linux",
        OSImage: "ubuntu-22.04-server-cloudimg-amd64.img",
        Resources: &deployment.ResourceRequirements{
            CPUCores:  16,
            MemoryGB:  64,
            GPUCount:  4,
            GPUModel:  "H100",
            StorageGB: 1000,
        },
        SSHKeys: []string{
            "ssh-rsa AAAAB3NzaC1yc2E...",
        },
        Hypervisor: "kvm",
        CloudInit:  "#cloud-config\n...",
    }

    // 验证配置
    if err := vmConfig.Validate(); err != nil {
        panic(err)
    }
}
```

## 4. 部署管理器 (推荐使用)

```go
func ExampleDeploymentManager() {
    // 创建部署管理器
    deploymentMgr := deployment.NewDeploymentManager()

    // 注册 Kubernetes 部署配置
    k8sConfig := &deployment.K8sDeploymentConfig{
        Namespace: "remotegpu",
        Name:      "env-123",
        Image:     "nvidia/cuda:11.8.0-runtime-ubuntu22.04",
        Resources: &deployment.ResourceRequirements{
            CPUCores: 4,
            MemoryGB: 16,
            GPUCount: 1,
            GPUModel: "RTX-3090",
        },
    }
    deploymentMgr.Register("env-123", k8sConfig)

    // 获取部署信息
    info, _ := deploymentMgr.Get("env-123")
    // info.Type: DeploymentTypeKubernetes
    // info.Status: DeploymentStatusPending

    // 更新部署状态
    deploymentMgr.UpdateStatus("env-123", deployment.DeploymentStatusRunning, "")

    // 列出所有部署
    allDeployments := deploymentMgr.List()

    // 删除部署信息
    deploymentMgr.Delete("env-123")
}
```

## 5. 在 Service 层使用

```go
// internal/service/environment.go

import (
    "github.com/YoungBoyGod/remotegpu/pkg/deployment"
)

type EnvironmentService struct {
    deploymentMgr *deployment.DeploymentManager
    // ...
}

func (s *EnvironmentService) CreateEnvironment(req *CreateEnvironmentRequest) (*entity.Environment, error) {
    // ... 创建环境逻辑 ...

    // 创建部署配置
    var deployConfig deployment.DeploymentConfig

    switch req.DeploymentType {
    case "kubernetes":
        deployConfig = &deployment.K8sDeploymentConfig{
            Namespace: "remotegpu",
            Name:      env.ID,
            Image:     req.Image,
            Resources: &deployment.ResourceRequirements{
                CPUCores:  req.CPUCores,
                MemoryGB:  req.MemoryGB,
                GPUCount:  req.GPUCount,
                GPUModel:  req.GPUModel,
                StorageGB: req.StorageGB,
            },
        }
    case "docker":
        deployConfig = &deployment.DockerDeploymentConfig{
            ContainerName: env.ID,
            Image:         req.Image,
            Resources: &deployment.ResourceRequirements{
                CPUCores:  req.CPUCores,
                MemoryGB:  req.MemoryGB,
                GPUCount:  req.GPUCount,
                GPUModel:  req.GPUModel,
                StorageGB: req.StorageGB,
            },
            Runtime: "nvidia",
        }
    }

    // 注册部署配置
    if err := s.deploymentMgr.Register(env.ID, deployConfig); err != nil {
        return nil, fmt.Errorf("注册部署配置失败: %w", err)
    }

    return env, nil
}
```

