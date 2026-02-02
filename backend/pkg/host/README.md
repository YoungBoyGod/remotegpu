# pkg/host - 主机管理模块

## 功能

统一的主机管理接口，提供：
- 主机选择策略（最低使用率、最多剩余资源等）
- 主机健康检查
- 主机资源管理

## 使用示例

### 1. 创建主机管理器

```go
package main

import (
    "github.com/YoungBoyGod/remotegpu/pkg/host"
    "time"
)

func main() {
    // 创建主机选择器
    selector := host.NewLeastUsedSelector()

    // 创建健康检查器
    checker := host.NewSimpleHealthChecker(5 * time.Second)
    healthMonitor := host.NewHealthMonitor(checker, 30 * time.Second)

    // 创建主机管理器
    manager := host.NewHostManager(selector, healthMonitor)

    // 注册主机
    host1 := &host.HostInfo{
        ID:          "host-1",
        Name:        "GPU Server 1",
        IPAddress:   "192.168.1.100",
        Status:      "active",
        TotalCPU:    32,
        UsedCPU:     8,
        TotalMemory: 128 * 1024 * 1024 * 1024,
        UsedMemory:  32 * 1024 * 1024 * 1024,
        TotalGPU:    4,
        UsedGPU:     1,
    }
    manager.RegisterHost(host1)
}
```

### 2. 选择主机

```go
func selectHost(manager *host.HostManager) {
    // 定义资源需求
    req := &host.ResourceRequirement{
        CPU:    4,
        Memory: 16 * 1024 * 1024 * 1024,
        GPU:    1,
    }

    // 选择主机
    selectedHost, err := manager.SelectHost(req)
    if err != nil {
        panic(err)
    }

    println("选择的主机:", selectedHost.Name)
}
```

### 3. 更新主机资源

```go
func allocateResources(manager *host.HostManager) {
    // 分配资源
    err := manager.UpdateHostResources("host-1", 4, 16*1024*1024*1024, 1, true)
    if err != nil {
        panic(err)
    }

    // 释放资源
    err = manager.UpdateHostResources("host-1", 4, 16*1024*1024*1024, 1, false)
    if err != nil {
        panic(err)
    }
}
```

### 4. 健康检查

```go
func checkHealth(manager *host.HostManager) {
    // 获取主机健康状态
    result, err := manager.GetHealthStatus("host-1")
    if err != nil {
        panic(err)
    }

    println("健康状态:", result.Status)
    println("检查信息:", result.Message)
}
```

## 主机选择策略

### LeastUsedSelector
选择资源使用率最低的主机，适合负载均衡场景。

## API 接口

### HostManager

```go
type HostManager struct {
    // ...
}

func (m *HostManager) RegisterHost(host *HostInfo)
func (m *HostManager) UnregisterHost(hostID string)
func (m *HostManager) GetHost(hostID string) (*HostInfo, error)
func (m *HostManager) ListHosts() []*HostInfo
func (m *HostManager) SelectHost(req *ResourceRequirement) (*HostInfo, error)
func (m *HostManager) UpdateHostResources(hostID string, cpu int, memory int64, gpu int, add bool) error
func (m *HostManager) GetHealthStatus(hostID string) (*HealthCheckResult, error)
```

## 特性

- **灵活的选择策略**: 支持多种主机选择算法
- **健康监控**: 自动监控主机健康状态
- **资源管理**: 统一管理主机资源分配
- **线程安全**: 所有操作都是线程安全的
