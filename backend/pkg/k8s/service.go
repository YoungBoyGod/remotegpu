// Package k8s 提供Kubernetes Service和网络管理功能
package k8s

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ServiceType Service类型
type ServiceType string

const (
	// ServiceTypeClusterIP 集群内部访问
	ServiceTypeClusterIP ServiceType = "ClusterIP"
	// ServiceTypeNodePort 通过节点端口访问
	ServiceTypeNodePort ServiceType = "NodePort"
	// ServiceTypeLoadBalancer 通过负载均衡器访问
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
)

// ServicePort Service端口配置
type ServicePort struct {
	// Name 端口名称
	Name string
	// Protocol 协议（TCP/UDP）
	Protocol corev1.Protocol
	// Port Service端口
	Port int32
	// TargetPort 目标端口（Pod端口）
	TargetPort int32
	// NodePort 节点端口（仅NodePort类型）
	NodePort int32
}

// ServiceConfig Service配置
type ServiceConfig struct {
	// Name Service名称
	Name string
	// Namespace 命名空间
	Namespace string
	// Type Service类型
	Type ServiceType
	// Selector 标签选择器
	Selector map[string]string
	// Ports 端口映射列表
	Ports []ServicePort
	// Labels 标签
	Labels map[string]string
	// Annotations 注解
	Annotations map[string]string
}

// Validate 验证Service配置
func (c *ServiceConfig) Validate() error {
	if c.Name == "" {
		return WrapError(ErrInvalidConfig, "service name is required")
	}
	if c.Namespace == "" {
		return WrapError(ErrInvalidConfig, "namespace is required")
	}
	if len(c.Selector) == 0 {
		return WrapError(ErrInvalidConfig, "selector is required")
	}
	if len(c.Ports) == 0 {
		return WrapError(ErrInvalidConfig, "at least one port is required")
	}

	// 验证端口配置
	for i, port := range c.Ports {
		if port.Port <= 0 {
			return WrapError(ErrInvalidConfig, fmt.Sprintf("invalid port at index %d: port must be positive", i))
		}
		if port.TargetPort <= 0 {
			return WrapError(ErrInvalidConfig, fmt.Sprintf("invalid port at index %d: target port must be positive", i))
		}
		if port.Protocol == "" {
			c.Ports[i].Protocol = corev1.ProtocolTCP // 默认TCP
		}
	}

	// 验证Service类型
	if c.Type == "" {
		c.Type = ServiceTypeClusterIP // 默认ClusterIP
	}

	return nil
}

// CreateService 创建Service
// 根据配置创建Kubernetes Service，支持ClusterIP、NodePort和LoadBalancer类型
func (c *Client) CreateService(config *ServiceConfig) (*corev1.Service, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 构建Service端口列表
	servicePorts := make([]corev1.ServicePort, 0, len(config.Ports))
	for _, port := range config.Ports {
		servicePort := corev1.ServicePort{
			Name:       port.Name,
			Protocol:   port.Protocol,
			Port:       port.Port,
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		}

		// NodePort类型需要指定NodePort
		if config.Type == ServiceTypeNodePort && port.NodePort > 0 {
			servicePort.NodePort = port.NodePort
		}

		servicePorts = append(servicePorts, servicePort)
	}

	// 构建Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceType(config.Type),
			Selector: config.Selector,
			Ports:    servicePorts,
		},
	}

	// 调用K8s API创建Service
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	createdService, err := c.clientset.CoreV1().Services(config.Namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to create service",
			zap.String("name", config.Name),
			zap.String("namespace", config.Namespace),
			zap.Error(err))
		return nil, WrapError(ErrPodCreationFailed, fmt.Sprintf("failed to create service: %v", err))
	}

	logger.GetLogger().Info("Service created successfully",
		zap.String("name", createdService.Name),
		zap.String("namespace", createdService.Namespace),
		zap.String("type", string(config.Type)),
		zap.String("cluster_ip", createdService.Spec.ClusterIP))

	return createdService, nil
}

// GetService 获取Service
// 根据命名空间和名称获取Service信息
func (c *Client) GetService(namespace, name string) (*corev1.Service, error) {
	if namespace == "" || name == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	service, err := c.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to get service",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return nil, WrapError(ErrPodNotFound, fmt.Sprintf("failed to get service: %v", err))
	}

	return service, nil
}

// ListServices 列出命名空间下的所有Service
// 可选择性地通过标签选择器过滤
func (c *Client) ListServices(namespace string, labelSelector string) (*corev1.ServiceList, error) {
	if namespace == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace is required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	services, err := c.clientset.CoreV1().Services(namespace).List(ctx, listOptions)
	if err != nil {
		logger.GetLogger().Error("Failed to list services",
			zap.String("namespace", namespace),
			zap.String("label_selector", labelSelector),
			zap.Error(err))
		return nil, WrapError(err, "failed to list services")
	}

	return services, nil
}

// UpdateService 更新Service
// 更新现有Service的配置
func (c *Client) UpdateService(service *corev1.Service) (*corev1.Service, error) {
	if service == nil {
		return nil, WrapError(ErrInvalidConfig, "service is nil")
	}
	if service.Name == "" || service.Namespace == "" {
		return nil, WrapError(ErrInvalidConfig, "service name and namespace are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	updatedService, err := c.clientset.CoreV1().Services(service.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to update service",
			zap.String("name", service.Name),
			zap.String("namespace", service.Namespace),
			zap.Error(err))
		return nil, WrapError(err, fmt.Sprintf("failed to update service: %v", err))
	}

	logger.GetLogger().Info("Service updated successfully",
		zap.String("name", updatedService.Name),
		zap.String("namespace", updatedService.Namespace))

	return updatedService, nil
}

// DeleteService 删除Service
// 根据命名空间和名称删除Service
func (c *Client) DeleteService(namespace, name string) error {
	if namespace == "" || name == "" {
		return WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	err := c.clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to delete service",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return WrapError(ErrPodDeletionFailed, fmt.Sprintf("failed to delete service: %v", err))
	}

	logger.GetLogger().Info("Service deleted successfully",
		zap.String("name", name),
		zap.String("namespace", namespace))

	return nil
}
