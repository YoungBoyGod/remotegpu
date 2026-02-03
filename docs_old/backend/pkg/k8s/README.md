# K8s 客户端和 Pod 管理使用文档

## 概述

本包提供了 Kubernetes 客户端封装和 Pod 管理功能，支持：
- K8s 客户端初始化（kubeconfig 和 in-cluster 模式）
- Pod 创建、查询、删除
- Pod 日志获取
- Pod 状态监控

## 快速开始

### 1. 配置

在 `config.yaml` 中配置 K8s 连接信息：

```yaml
k8s:
  enabled: true
  kubeconfig: "/path/to/kubeconfig"  # kubeconfig 文件路径
  namespace: "default"                # 默认命名空间
  in_cluster: false                   # 是否在集群内运行
```

### 2. 初始化客户端

```go
import "github.com/YoungBoyGod/remotegpu/pkg/k8s"

// 方式1：使用全局配置自动初始化
client, err := k8s.GetClient()
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 方式2：使用自定义配置
config := &k8s.Config{
    KubeConfig: "/path/to/kubeconfig",
    Namespace:  "default",
    InCluster:  false,
    Timeout:    30 * time.Second,
}
client, err := k8s.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## Pod 管理

### 创建 Pod

#### 基本示例

```go
config := &k8s.PodConfig{
    Name:      "my-pod",
    Namespace: "default",
    Image:     "nginx:latest",
}

pod, err := client.CreatePod(config)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Pod created: %s\n", pod.Name)
```

#### 带资源限制的 Pod

```go
config := &k8s.PodConfig{
    Name:      "gpu-pod",
    Namespace: "default",
    Image:     "tensorflow/tensorflow:latest-gpu",
    CPU:       4,      // 4 核 CPU
    Memory:    8192,   // 8GB 内存
    GPU:       1,      // 1 个 GPU
}

pod, err := client.CreatePod(config)
```

#### 带环境变量的 Pod

```go
config := &k8s.PodConfig{
    Name:      "app-pod",
    Image:     "myapp:latest",
    Env: map[string]string{
        "DATABASE_URL": "postgres://...",
        "API_KEY":      "secret-key",
    },
}

pod, err := client.CreatePod(config)
```

#### 带卷挂载的 Pod

```go
config := &k8s.PodConfig{
    Name:  "data-pod",
    Image: "ubuntu:latest",
    Volumes: []k8s.VolumeMount{
        {
            Name:      "data",
            MountPath: "/data",
            HostPath:  "/host/data",
            ReadOnly:  false,
        },
    },
}

pod, err := client.CreatePod(config)
```

#### 完整配置示例

```go
config := &k8s.PodConfig{
    Name:      "full-pod",
    Namespace: "production",
    Image:     "myapp:v1.0",
    Command:   []string{"/bin/sh"},
    Args:      []string{"-c", "echo hello"},
    CPU:       2,
    Memory:    4096,
    GPU:       1,
    Env: map[string]string{
        "ENV": "production",
    },
    Volumes: []k8s.VolumeMount{
        {
            Name:      "config",
            MountPath: "/etc/config",
            HostPath:  "/host/config",
            ReadOnly:  true,
        },
    },
    Labels: map[string]string{
        "app":     "myapp",
        "version": "v1.0",
    },
    Annotations: map[string]string{
        "description": "My application pod",
    },
    RestartPolicy: corev1.RestartPolicyAlways,
}

pod, err := client.CreatePod(config)
```

### 查询 Pod

#### 获取 Pod 详细信息

```go
pod, err := client.GetPod("default", "my-pod")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Pod: %s, Status: %s\n", pod.Name, pod.Status.Phase)
```

#### 获取 Pod 状态

```go
status, err := client.GetPodStatus("default", "my-pod")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Pod status: %s\n", status)
```

### 删除 Pod

#### 立即删除

```go
err := client.DeletePod("default", "my-pod")
if err != nil {
    log.Fatal(err)
}
```

#### 优雅删除（带 grace period）

```go
// 等待 30 秒后删除
err := client.DeletePodGracefully("default", "my-pod", 30)
if err != nil {
    log.Fatal(err)
}
```

### 获取 Pod 日志

#### 获取所有日志

```go
logs, err := client.GetPodLogs("default", "my-pod", nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(logs)
```

#### 获取最后 100 行日志

```go
opts := &k8s.LogOptions{
    TailLines: 100,
}
logs, err := client.GetPodLogs("default", "my-pod", opts)
```

#### 获取带时间戳的日志

```go
opts := &k8s.LogOptions{
    Timestamps: true,
    TailLines:  50,
}
logs, err := client.GetPodLogs("default", "my-pod", opts)
```

#### 获取指定容器的日志

```go
opts := &k8s.LogOptions{
    Container:  "sidecar",
    TailLines:  100,
}
logs, err := client.GetPodLogs("default", "my-pod", opts)
```

### 监控 Pod 状态

```go
callback := func(status string) {
    fmt.Printf("Pod status changed: %s\n", status)
}

err := client.WatchPodStatus("default", "my-pod", callback)
if err != nil {
    log.Fatal(err)
}
```

## 错误处理

### 错误类型

```go
import "github.com/YoungBoyGod/remotegpu/pkg/k8s"

// 判断是否为 Pod 不存在错误
if k8s.IsPodNotFound(err) {
    fmt.Println("Pod not found")
}

// 判断是否为连接失败错误
if k8s.IsConnectionFailed(err) {
    fmt.Println("Failed to connect to K8s")
}
```

### 常见错误

| 错误 | 说明 | 解决方法 |
|------|------|----------|
| `ErrPodNotFound` | Pod 不存在 | 检查 Pod 名称和命名空间 |
| `ErrPodCreationFailed` | Pod 创建失败 | 检查配置和资源限制 |
| `ErrPodDeletionFailed` | Pod 删除失败 | 检查 Pod 状态和权限 |
| `ErrConnectionFailed` | K8s 连接失败 | 检查 kubeconfig 和网络 |
| `ErrInvalidConfig` | 配置无效 | 检查配置参数 |
| `ErrLogsFetchFailed` | 日志获取失败 | 检查 Pod 状态和容器名称 |

## 配置说明

### Config 结构体

```go
type Config struct {
    KubeConfig string        // kubeconfig 文件路径
    Namespace  string        // 默认命名空间
    InCluster  bool          // 是否在集群内运行
    Timeout    time.Duration // 操作超时时间
}
```

### PodConfig 结构体

```go
type PodConfig struct {
    Name          string                  // Pod 名称（必填）
    Namespace     string                  // 命名空间（可选，默认使用客户端命名空间）
    Image         string                  // 容器镜像（必填）
    Command       []string                // 容器命令
    Args          []string                // 容器参数
    CPU           int64                   // CPU 核心数
    Memory        int64                   // 内存大小（MB）
    GPU           int64                   // GPU 数量
    Env           map[string]string       // 环境变量
    Volumes       []VolumeMount           // 卷挂载
    Labels        map[string]string       // 标签
    Annotations   map[string]string       // 注解
    RestartPolicy corev1.RestartPolicy    // 重启策略
}
```

## GPU 配置

### NVIDIA GPU 支持

本包使用 NVIDIA Device Plugin 来支持 GPU。确保你的 K8s 集群已安装 NVIDIA Device Plugin。

```go
config := &k8s.PodConfig{
    Name:   "gpu-pod",
    Image:  "nvidia/cuda:11.0-base",
    GPU:    2,  // 请求 2 个 GPU
}

pod, err := client.CreatePod(config)
```

GPU 资源会被自动配置为 `nvidia.com/gpu` 限制。

## 最佳实践

### 1. 使用单例模式

```go
// 推荐：使用全局客户端
client, err := k8s.GetClient()
defer client.Close()
```

### 2. 设置合理的超时时间

```go
config := &k8s.Config{
    Timeout: 60 * time.Second,  // 对于长时间运行的操作
}
```

### 3. 使用标签管理 Pod

```go
config := &k8s.PodConfig{
    Name:  "my-pod",
    Image: "myapp:latest",
    Labels: map[string]string{
        "app":     "myapp",
        "env":     "production",
        "version": "v1.0",
    },
}
```

### 4. 优雅删除 Pod

```go
// 给 Pod 30 秒时间清理资源
err := client.DeletePodGracefully("default", "my-pod", 30)
```

### 5. 错误处理

```go
pod, err := client.GetPod("default", "my-pod")
if err != nil {
    if k8s.IsPodNotFound(err) {
        // 处理 Pod 不存在的情况
        return nil
    }
    // 处理其他错误
    return err
}
```

## 测试

### 单元测试

```bash
# 运行所有测试
go test -v ./pkg/k8s

# 运行特定测试
go test -v ./pkg/k8s -run TestCreatePod

# 查看测试覆盖率
go test -cover ./pkg/k8s
```

### 集成测试

集成测试需要真实的 K8s 集群。使用 minikube 或 kind 搭建测试环境：

```bash
# 启动 minikube
minikube start

# 运行集成测试
go test -v ./pkg/k8s -tags=integration
```

## 常见问题

### Q: 如何在集群内运行？

A: 设置 `InCluster: true`：

```go
config := &k8s.Config{
    InCluster: true,
    Namespace: "default",
}
```

### Q: 如何配置 GPU？

A: 在 PodConfig 中设置 GPU 字段：

```go
config := &k8s.PodConfig{
    GPU: 1,  // 请求 1 个 GPU
}
```

### Q: 如何处理 Pod 创建失败？

A: 检查错误信息和 Pod 事件：

```go
pod, err := client.CreatePod(config)
if err != nil {
    log.Printf("Failed to create pod: %v", err)
    // 检查配置和资源限制
}
```

### Q: 如何监控多个 Pod？

A: 为每个 Pod 启动一个 goroutine：

```go
for _, podName := range podNames {
    go func(name string) {
        err := client.WatchPodStatus("default", name, callback)
        if err != nil {
            log.Printf("Watch failed: %v", err)
        }
    }(podName)
}
```

## 依赖

- k8s.io/client-go v0.35.0
- k8s.io/api v0.35.0
- k8s.io/apimachinery v0.35.0

## 版本历史

- v1.0.0 (2026-01-30)
  - 初始版本
  - 实现 K8s 客户端封装
  - 实现 Pod 管理功能
  - 实现日志获取和状态监控

## 贡献者

- A3 (Claude Sonnet 4.5)

## 许可证

[项目许可证]
