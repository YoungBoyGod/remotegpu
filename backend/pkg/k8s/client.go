// Package k8s 提供Kubernetes客户端封装和Pod管理功能
package k8s

import (
	"context"
	"sync"
	"time"

	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client K8s客户端封装
type Client struct {
	// clientset K8s客户端集
	clientset kubernetes.Interface

	// config 客户端配置
	config *Config

	// restConfig REST配置
	restConfig *rest.Config

	// namespace 默认命名空间
	namespace string

	// ctx 上下文
	ctx context.Context

	// cancel 取消函数
	cancel context.CancelFunc
}

var (
	// globalClient 全局客户端实例（单例）
	globalClient *Client

	// clientMutex 客户端互斥锁
	clientMutex sync.Mutex
)

// NewClient 创建新的K8s客户端
// 支持两种模式：
// 1. kubeconfig模式：通过kubeconfig文件连接集群
// 2. in-cluster模式：在Pod内部运行时自动获取配置
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, WrapError(ErrInvalidConfig, "config is nil")
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 构建REST配置
	restConfig, err := buildRestConfig(cfg)
	if err != nil {
		return nil, WrapError(ErrConnectionFailed, err.Error())
	}

	// 设置超时
	restConfig.Timeout = cfg.Timeout

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, WrapError(ErrConnectionFailed, err.Error())
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		clientset:  clientset,
		config:     cfg,
		restConfig: restConfig,
		namespace:  cfg.Namespace,
		ctx:        ctx,
		cancel:     cancel,
	}

	logger.GetLogger().Info("K8s client initialized successfully",
		zap.String("namespace", cfg.Namespace),
		zap.Bool("in_cluster", cfg.InCluster))

	return client, nil
}

// NewClientWithClientset 使用提供的clientset创建客户端（用于测试）
func NewClientWithClientset(clientset kubernetes.Interface, namespace string) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		clientset:  clientset,
		config:     &Config{Namespace: namespace, Timeout: 30 * time.Second},
		restConfig: nil,
		namespace:  namespace,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// buildRestConfig 构建REST配置
func buildRestConfig(cfg *Config) (*rest.Config, error) {
	var restConfig *rest.Config
	var err error

	if cfg.InCluster {
		// In-cluster模式：从Pod环境变量和ServiceAccount获取配置
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, WrapError(err, "failed to get in-cluster config")
		}
		logger.GetLogger().Info("Using in-cluster kubernetes config")
	} else {
		// Kubeconfig模式：从文件加载配置
		restConfig, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
		if err != nil {
			return nil, WrapErrorf(err, "failed to build config from kubeconfig: %s", cfg.KubeConfig)
		}
		logger.GetLogger().Info("Using kubeconfig file", zap.String("path", cfg.KubeConfig))
	}

	return restConfig, nil
}

// GetClient 获取全局K8s客户端实例（单例模式）
// 如果客户端未初始化，会自动从全局配置加载并初始化
func GetClient() (*Client, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if globalClient != nil {
		return globalClient, nil
	}

	// 从全局配置加载
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// 创建客户端
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	globalClient = client
	return globalClient, nil
}

// InitClient 初始化全局客户端（使用自定义配置）
func InitClient(cfg *Config) error {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if globalClient != nil {
		globalClient.Close()
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	globalClient = client
	return nil
}

// Close 关闭客户端，释放资源
func (c *Client) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	logger.GetLogger().Info("K8s client closed")
}

// GetClientset 获取kubernetes clientset
func (c *Client) GetClientset() kubernetes.Interface {
	return c.clientset
}

// GetNamespace 获取默认命名空间
func (c *Client) GetNamespace() string {
	return c.namespace
}

// GetContext 获取上下文
func (c *Client) GetContext() context.Context {
	return c.ctx
}

// GetContextWithTimeout 获取带超时的上下文
func (c *Client) GetContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.ctx, timeout)
}

// Ping 测试K8s连接
func (c *Client) Ping() error {
	_, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return WrapError(ErrConnectionFailed, err.Error())
	}

	return nil
}
