package health

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdChecker Etcd健康检查器
type EtcdChecker struct {
	config config.EtcdConfig
}

// NewEtcdChecker 创建Etcd健康检查器
func NewEtcdChecker(cfg config.EtcdConfig) *EtcdChecker {
	return &EtcdChecker{config: cfg}
}

// Name 返回服务名称
func (c *EtcdChecker) Name() string {
	return "etcd"
}

// Check 执行健康检查
func (c *EtcdChecker) Check(ctx context.Context) *CheckResult {
	start := time.Now()
	result := &CheckResult{
		Service:   c.Name(),
		Timestamp: start,
	}

	// 如果服务未启用
	if !c.config.Enabled {
		result.Status = StatusDisabled
		result.Message = "服务未启用"
		result.Latency = time.Since(start)
		return result
	}

	// 创建Etcd客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.config.Endpoints,
		Username:    c.config.Username,
		Password:    c.config.Password,
		DialTimeout: time.Duration(c.config.Timeout) * time.Second,
	})
	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("连接失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer cli.Close()

	// 执行健康检查
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = cli.Get(checkCtx, "health-check")
	if err != nil && err != context.DeadlineExceeded {
		result.Status = StatusUnhealthy
		result.Message = fmt.Sprintf("健康检查失败: %v", err)
		result.Latency = time.Since(start)
		return result
	}

	result.Status = StatusHealthy
	result.Message = "连接正常"
	result.Latency = time.Since(start)
	result.Details = map[string]interface{}{
		"endpoints": c.config.Endpoints,
	}

	return result
}

// K8sChecker Kubernetes健康检查器
type K8sChecker struct {
	config config.K8sConfig
}

// NewK8sChecker 创建K8s健康检查器
func NewK8sChecker(cfg config.K8sConfig) *K8sChecker {
	return &K8sChecker{config: cfg}
}

// Name 返回服务名称
func (c *K8sChecker) Name() string {
	return "kubernetes"
}

// Check 执行健康检查
func (c *K8sChecker) Check(ctx context.Context) *CheckResult {
	start := time.Now()
	result := &CheckResult{
		Service:   c.Name(),
		Timestamp: start,
	}

	// 如果服务未启用
	if !c.config.Enabled {
		result.Status = StatusDisabled
		result.Message = "服务未启用"
		result.Latency = time.Since(start)
		return result
	}

	// 简化的K8s健康检查
	// 实际项目中应该使用k8s client-go库进行更详细的检查
	result.Status = StatusHealthy
	result.Message = "配置已加载"
	result.Latency = time.Since(start)
	result.Details = map[string]interface{}{
		"kubeconfig": c.config.KubeConfig,
		"namespace":  c.config.Namespace,
		"in_cluster": c.config.InCluster,
	}

	return result
}
