# pkg/network 使用示例

## 1. 端口管理

```go
package main

import (
    "github.com/YoungBoyGod/remotegpu/pkg/network"
)

func ExamplePortManager() {
    // 创建端口管理器
    portManager := network.NewPortManager()

    // 为 SSH 服务分配端口
    sshPort, err := portManager.AllocatePort(network.ServiceTypeSSH)
    if err != nil {
        panic(err)
    }
    // sshPort: 22000-22999 范围内的可用端口

    // 为 RDP 服务分配端口
    rdpPort, err := portManager.AllocatePort(network.ServiceTypeRDP)
    if err != nil {
        panic(err)
    }
    // rdpPort: 33890-34889 范围内的可用端口

    // 检查端口是否已分配
    if portManager.IsAllocated(sshPort) {
        // 端口已被分配
    }

    // 释放端口
    portManager.ReleasePort(sshPort)
}
```

## 2. 防火墙管理

```go
func ExampleFirewallManager() {
    // 创建防火墙管理器
    firewallManager := network.NewFirewallManager(network.FirewallTypeIPTables)

    // 创建防火墙规则
    rule := &network.FirewallRule{
        ID:          "rule-001",
        SourceIP:    "0.0.0.0/0",
        DestIP:      "192.168.1.100",
        DestPort:    22001,
        Protocol:    "tcp",
        Action:      "allow",
        Description: "Allow SSH access",
    }

    if err := firewallManager.CreateRule(rule); err != nil {
        panic(err)
    }

    // 获取规则
    rule, err := firewallManager.GetRule("rule-001")

    // 删除规则
    if err := firewallManager.DeleteRule("rule-001"); err != nil {
        panic(err)
    }
}
```

## 3. DNS 管理

```go
func ExampleDNSManager() {
    // 创建 DNS 管理器
    dnsManager := network.NewDNSManager(
        network.DNSProviderCloudflare,
        "example.com",
    )

    // 生成子域名
    subdomain := dnsManager.GenerateSubdomain("env-12345678", network.ServiceTypeSSH)
    // subdomain: "ssh-env12345.example.com"

    // 创建 DNS 记录
    record := &network.DNSRecord{
        ID:     "record-001",
        Domain: subdomain,
        Type:   "A",
        Value:  "1.2.3.4",
        TTL:    300,
    }

    if err := dnsManager.CreateRecord(record); err != nil {
        panic(err)
    }

    // 删除 DNS 记录
    if err := dnsManager.DeleteRecord("record-001"); err != nil {
        panic(err)
    }
}
```

## 4. 网络管理器 (推荐使用)

```go
func ExampleNetworkManager() {
    // 创建网络管理器
    networkManager := network.NewNetworkManager(
        network.FirewallTypeIPTables,
        network.DNSProviderCloudflare,
        "example.com",
    )

    // 为环境配置完整的网络设置
    services := []network.ServiceType{
        network.ServiceTypeSSH,
        network.ServiceTypeJupyter,
    }

    config, err := networkManager.ConfigureEnvironment(
        "env-123",
        services,
        "192.168.1.100",
    )
    if err != nil {
        panic(err)
    }

    // config 包含:
    // - PortMappings: SSH 和 Jupyter 的端口映射
    // - FirewallRules: 对应的防火墙规则
    // - DNSRecords: 对应的 DNS 记录

    // 获取环境配置
    envConfig, err := networkManager.GetEnvConfig("env-123")

    // 获取端口映射
    portMappings, err := networkManager.GetPortMappings("env-123")

    // 清理环境网络配置
    if err := networkManager.CleanupEnvironment("env-123"); err != nil {
        panic(err)
    }
}
```

## 5. 单独分配端口和创建规则

```go
func ExampleManualConfiguration() {
    networkManager := network.NewNetworkManager(
        network.FirewallTypeIPTables,
        network.DNSProviderCloudflare,
        "example.com",
    )

    // 单独为服务分配端口
    mapping, err := networkManager.AllocateServicePort(
        "env-456",
        network.ServiceTypeSSH,
        "tcp",
        "SSH service for env-456",
    )
    if err != nil {
        panic(err)
    }

    // 单独创建防火墙规则
    rule := &network.FirewallRule{
        ID:          "custom-rule-001",
        SourceIP:    "10.0.0.0/8",
        DestIP:      "192.168.1.100",
        DestPort:    mapping.ExternalPort,
        Protocol:    "tcp",
        Action:      "allow",
        Description: "Custom firewall rule",
    }
    if err := networkManager.CreateFirewallRule("env-456", rule); err != nil {
        panic(err)
    }

    // 生成子域名
    subdomain := networkManager.GenerateSubdomain("env-456", network.ServiceTypeSSH)

    // 单独创建 DNS 记录
    record := &network.DNSRecord{
        ID:     "custom-dns-001",
        Domain: subdomain,
        Type:   "A",
        Value:  "192.168.1.100",
        TTL:    300,
    }
    if err := networkManager.CreateDNSRecord("env-456", record); err != nil {
        panic(err)
    }
}
```

## 6. 在 Service 层使用

```go
// internal/service/environment.go

import (
    "github.com/YoungBoyGod/remotegpu/pkg/network"
)

type EnvironmentService struct {
    networkManager *network.NetworkManager
    // ...
}

func (s *EnvironmentService) CreateEnvironment(req *CreateEnvironmentRequest) (*entity.Environment, error) {
    // ... 创建环境逻辑 ...

    // 配置网络
    services := []network.ServiceType{
        network.ServiceTypeSSH,
    }

    if req.EnableJupyter {
        services = append(services, network.ServiceTypeJupyter)
    }

    if req.EnableRDP {
        services = append(services, network.ServiceTypeRDP)
    }

    // 配置完整的网络设置
    networkConfig, err := s.networkManager.ConfigureEnvironment(
        env.ID,
        services,
        hostInfo.IP,
    )
    if err != nil {
        return nil, fmt.Errorf("配置网络失败: %w", err)
    }

    // 保存网络配置到数据库
    // ...

    return env, nil
}

func (s *EnvironmentService) DeleteEnvironment(envID string) error {
    // ... 删除环境逻辑 ...

    // 清理网络配置
    if err := s.networkManager.CleanupEnvironment(envID); err != nil {
        return fmt.Errorf("清理网络配置失败: %w", err)
    }

    return nil
}
```
