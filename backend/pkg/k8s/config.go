package k8s

import (
	"fmt"
	"os"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
)

// Config K8s客户端配置
type Config struct {
	// KubeConfig kubeconfig文件路径
	KubeConfig string

	// Namespace 默认命名空间
	Namespace string

	// InCluster 是否在集群内运行
	InCluster bool

	// Timeout 操作超时时间
	Timeout time.Duration
}

// DefaultConfig 默认配置
var DefaultConfig = &Config{
	Namespace: "default",
	InCluster: false,
	Timeout:   30 * time.Second,
}

// LoadConfig 从全局配置加载K8s配置
func LoadConfig() (*Config, error) {
	if config.GlobalConfig == nil {
		return nil, WrapError(ErrInvalidConfig, "global config not loaded")
	}

	k8sCfg := config.GlobalConfig.K8s

	// 如果K8s未启用，返回错误
	if !k8sCfg.Enabled {
		return nil, WrapError(ErrInvalidConfig, "kubernetes is not enabled in config")
	}

	cfg := &Config{
		KubeConfig: k8sCfg.KubeConfig,
		Namespace:  k8sCfg.Namespace,
		InCluster:  k8sCfg.InCluster,
		Timeout:    30 * time.Second,
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 如果不是in-cluster模式，必须提供kubeconfig路径
	if !c.InCluster && c.KubeConfig == "" {
		return WrapError(ErrInvalidConfig, "kubeconfig path is required when not running in-cluster")
	}

	// 如果提供了kubeconfig路径，检查文件是否存在
	if c.KubeConfig != "" {
		if _, err := os.Stat(c.KubeConfig); os.IsNotExist(err) {
			return WrapErrorf(ErrInvalidConfig, "kubeconfig file not found: %s", c.KubeConfig)
		}
	}

	// 验证命名空间
	if c.Namespace == "" {
		c.Namespace = "default"
	}

	// 验证超时时间
	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Second
	}

	return nil
}

// String 返回配置的字符串表示
func (c *Config) String() string {
	return fmt.Sprintf("K8sConfig{Namespace: %s, InCluster: %v, Timeout: %v}",
		c.Namespace, c.InCluster, c.Timeout)
}
