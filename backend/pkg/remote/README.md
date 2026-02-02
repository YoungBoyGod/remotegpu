# pkg/remote 使用示例

## 1. SSH 访问配置

```go
package main

import (
    "github.com/YoungBoyGod/remotegpu/pkg/remote"
)

func ExampleSSH() {
    // 创建 SSH 配置
    sshConfig := &remote.SSHConfig{
        Port:       22001,
        Username:   "root",
        Password:   "Pass_12345678",
        AccessType: remote.AccessTypeDirect,
    }

    // 验证配置
    if err := sshConfig.Validate(); err != nil {
        panic(err)
    }

    // 生成访问信息
    host := &remote.HostInfo{
        InternalIP:   "192.168.1.100",
        PublicIP:     "1.2.3.4",
        PublicDomain: "ssh-env123.example.com",
    }

    accessInfo := sshConfig.GenerateAccessInfo("env-123", host)

    // 输出访问信息
    // accessInfo.InternalURL: ssh://root@192.168.1.100:22001
    // accessInfo.PublicURL: ssh://root@ssh-env123.example.com:22001
    // accessInfo.Command: ssh root@192.168.1.100 -p 22001
}
```

## 2. RDP 访问配置

```go
func ExampleRDP() {
    rdpConfig := &remote.RDPConfig{
        Port:       3389,
        Username:   "Administrator",
        Password:   "Pass_12345678",
        AccessType: remote.AccessTypeDirect,
    }

    host := &remote.HostInfo{
        InternalIP: "192.168.1.100",
        PublicIP:   "1.2.3.4",
    }

    accessInfo := rdpConfig.GenerateAccessInfo("env-456", host)
    // accessInfo.InternalURL: rdp://192.168.1.100:3389
}
```

## 3. 使用访问管理器

```go
func ExampleManager() {
    manager := remote.NewAccessManager()

    // 注册 SSH 配置
    sshConfig := &remote.SSHConfig{
        Port:     22001,
        Username: "root",
        Password: "Pass_12345678",
    }
    manager.Register("env-123", sshConfig)

    // 注册 RDP 配置
    rdpConfig := &remote.RDPConfig{
        Port:     3389,
        Username: "Administrator",
        Password: "Pass_12345678",
    }
    manager.Register("env-123", rdpConfig)

    // 获取所有访问配置
    configs := manager.Get("env-123")
    // len(configs) == 2

    // 获取指定协议的配置
    sshCfg := manager.GetByProtocol("env-123", remote.ProtocolSSH)

    // 生成所有访问信息
    host := &remote.HostInfo{
        InternalIP: "192.168.1.100",
        PublicIP:   "1.2.3.4",
    }
    allAccessInfo := manager.GenerateAllAccessInfo("env-123", host)
    // len(allAccessInfo) == 2
}
```

## 4. Jumpserver 集成

```go
func ExampleJumpserver() {
    sshConfig := &remote.SSHConfig{
        Port:              22001,
        Username:          "root",
        Password:          "Pass_12345678",
        AccessType:        remote.AccessTypeJumpserver,
        JumpserverAssetID: "asset-123",
    }

    host := &remote.HostInfo{
        InternalIP: "192.168.1.100",
    }

    accessInfo := sshConfig.GenerateAccessInfo("env-123", host)
    // accessInfo.WebURL: /luna/?asset=asset-123
}
```

## 5. Guacamole 集成

```go
func ExampleGuacamole() {
    rdpConfig := &remote.RDPConfig{
        Port:            3389,
        Username:        "Administrator",
        Password:        "Pass_12345678",
        AccessType:      remote.AccessTypeGuacamole,
        GuacamoleConnID: "conn-456",
    }

    host := &remote.HostInfo{
        InternalIP: "192.168.1.100",
    }

    accessInfo := rdpConfig.GenerateAccessInfo("env-456", host)
    // accessInfo.WebURL: /#/client/conn-456
}
```

## 6. 在 Service 层使用

```go
// internal/service/environment.go

import (
    "github.com/YoungBoyGod/remotegpu/pkg/remote"
)

type EnvironmentService struct {
    remoteManager *remote.AccessManager
    // ...
}

func (s *EnvironmentService) CreateEnvironment(req *CreateEnvironmentRequest) (*entity.Environment, error) {
    // ... 创建环境逻辑 ...

    // 配置 SSH 访问
    sshConfig := &remote.SSHConfig{
        Port:     *env.SSHPort,
        Username: "root",
        Password: generatePassword(),
    }
    s.remoteManager.Register(env.ID, sshConfig)

    // 配置 RDP 访问 (如果需要)
    if env.RDPPort != nil {
        rdpConfig := &remote.RDPConfig{
            Port:     *env.RDPPort,
            Username: "Administrator",
            Password: generatePassword(),
        }
        s.remoteManager.Register(env.ID, rdpConfig)
    }

    // 生成访问信息
    host := &remote.HostInfo{
        InternalIP:   hostInfo.IP,
        PublicIP:     hostInfo.PublicIP,
        PublicDomain: hostInfo.PublicDomain,
    }
    accessInfos := s.remoteManager.GenerateAllAccessInfo(env.ID, host)

    // 保存访问信息到数据库
    // ...

    return env, nil
}
```
