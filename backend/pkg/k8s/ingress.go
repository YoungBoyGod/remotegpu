// Package k8s 提供Kubernetes Ingress管理功能
package k8s

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"go.uber.org/zap"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressPathType Ingress路径类型
type IngressPathType string

const (
	// PathTypePrefix 前缀匹配
	PathTypePrefix IngressPathType = "Prefix"
	// PathTypeExact 精确匹配
	PathTypeExact IngressPathType = "Exact"
	// PathTypeImplementationSpecific 实现特定匹配
	PathTypeImplementationSpecific IngressPathType = "ImplementationSpecific"
)

// IngressTLS TLS配置
type IngressTLS struct {
	// Hosts TLS证书覆盖的主机列表
	Hosts []string
	// SecretName 包含TLS证书的Secret名称
	SecretName string
}

// IngressRule Ingress规则
type IngressRule struct {
	// Host 域名
	Host string
	// Path 路径
	Path string
	// PathType 路径类型
	PathType IngressPathType
	// ServiceName 后端Service名称
	ServiceName string
	// ServicePort 后端Service端口
	ServicePort int32
}

// IngressConfig Ingress配置
type IngressConfig struct {
	// Name Ingress名称
	Name string
	// Namespace 命名空间
	Namespace string
	// Rules Ingress规则列表
	Rules []IngressRule
	// TLS TLS配置
	TLS []IngressTLS
	// Labels 标签
	Labels map[string]string
	// Annotations 注解（用于配置Ingress Controller）
	Annotations map[string]string
	// IngressClassName Ingress类名称
	IngressClassName string
}

// Validate 验证Ingress配置
func (c *IngressConfig) Validate() error {
	if c.Name == "" {
		return WrapError(ErrInvalidConfig, "ingress name is required")
	}
	if c.Namespace == "" {
		return WrapError(ErrInvalidConfig, "namespace is required")
	}
	if len(c.Rules) == 0 {
		return WrapError(ErrInvalidConfig, "at least one rule is required")
	}

	// 验证规则配置
	for i, rule := range c.Rules {
		if rule.Host == "" {
			return WrapError(ErrInvalidConfig, fmt.Sprintf("invalid rule at index %d: host is required", i))
		}
		if rule.Path == "" {
			c.Rules[i].Path = "/" // 默认路径
		}
		if rule.PathType == "" {
			c.Rules[i].PathType = PathTypePrefix // 默认前缀匹配
		}
		if rule.ServiceName == "" {
			return WrapError(ErrInvalidConfig, fmt.Sprintf("invalid rule at index %d: service name is required", i))
		}
		if rule.ServicePort <= 0 {
			return WrapError(ErrInvalidConfig, fmt.Sprintf("invalid rule at index %d: service port must be positive", i))
		}
	}

	return nil
}

// CreateIngress 创建Ingress
// 根据配置创建Kubernetes Ingress，用于HTTP/HTTPS路由
func (c *Client) CreateIngress(config *IngressConfig) (*networkingv1.Ingress, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 构建Ingress规则
	ingressRules := make([]networkingv1.IngressRule, 0, len(config.Rules))

	// 按Host分组规则
	hostRules := make(map[string][]IngressRule)
	for _, rule := range config.Rules {
		hostRules[rule.Host] = append(hostRules[rule.Host], rule)
	}

	// 为每个Host构建IngressRule
	for host, rules := range hostRules {
		paths := make([]networkingv1.HTTPIngressPath, 0, len(rules))

		for _, rule := range rules {
			pathType := networkingv1.PathType(rule.PathType)
			path := networkingv1.HTTPIngressPath{
				Path:     rule.Path,
				PathType: &pathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: rule.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Number: rule.ServicePort,
						},
					},
				},
			}
			paths = append(paths, path)
		}

		ingressRule := networkingv1.IngressRule{
			Host: host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		}
		ingressRules = append(ingressRules, ingressRule)
	}

	// 构建TLS配置
	var tlsConfigs []networkingv1.IngressTLS
	if len(config.TLS) > 0 {
		tlsConfigs = make([]networkingv1.IngressTLS, 0, len(config.TLS))
		for _, tls := range config.TLS {
			tlsConfig := networkingv1.IngressTLS{
				Hosts:      tls.Hosts,
				SecretName: tls.SecretName,
			}
			tlsConfigs = append(tlsConfigs, tlsConfig)
		}
	}

	// 构建Ingress对象
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Spec: networkingv1.IngressSpec{
			Rules: ingressRules,
			TLS:   tlsConfigs,
		},
	}

	// 设置IngressClassName（如果指定）
	if config.IngressClassName != "" {
		ingress.Spec.IngressClassName = &config.IngressClassName
	}

	// 调用K8s API创建Ingress
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	createdIngress, err := c.clientset.NetworkingV1().Ingresses(config.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to create ingress",
			zap.String("name", config.Name),
			zap.String("namespace", config.Namespace),
			zap.Error(err))
		return nil, WrapError(err, fmt.Sprintf("failed to create ingress: %v", err))
	}

	logger.GetLogger().Info("Ingress created successfully",
		zap.String("name", createdIngress.Name),
		zap.String("namespace", createdIngress.Namespace),
		zap.Int("rules", len(createdIngress.Spec.Rules)))

	return createdIngress, nil
}

// GetIngress 获取Ingress
// 根据命名空间和名称获取Ingress信息
func (c *Client) GetIngress(namespace, name string) (*networkingv1.Ingress, error) {
	if namespace == "" || name == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	ingress, err := c.clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to get ingress",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return nil, WrapError(ErrPodNotFound, fmt.Sprintf("failed to get ingress: %v", err))
	}

	return ingress, nil
}

// DeleteIngress 删除Ingress
// 根据命名空间和名称删除Ingress
func (c *Client) DeleteIngress(namespace, name string) error {
	if namespace == "" || name == "" {
		return WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	err := c.clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to delete ingress",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return WrapError(ErrPodDeletionFailed, fmt.Sprintf("failed to delete ingress: %v", err))
	}

	logger.GetLogger().Info("Ingress deleted successfully",
		zap.String("name", name),
		zap.String("namespace", namespace))

	return nil
}
