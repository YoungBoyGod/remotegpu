// Package k8s 提供Kubernetes ConfigMap管理功能
package k8s

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateConfigMap 创建ConfigMap
// 根据命名空间、名称和数据创建ConfigMap
func (c *Client) CreateConfigMap(namespace, name string, data map[string]string) (*corev1.ConfigMap, error) {
	if namespace == "" || name == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace and name are required")
	}
	if data == nil {
		data = make(map[string]string)
	}

	// 构建ConfigMap对象
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	// 调用K8s API创建ConfigMap
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	createdConfigMap, err := c.clientset.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to create ConfigMap",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return nil, WrapError(err, fmt.Sprintf("failed to create ConfigMap: %v", err))
	}

	logger.GetLogger().Info("ConfigMap created successfully",
		zap.String("name", createdConfigMap.Name),
		zap.String("namespace", createdConfigMap.Namespace),
		zap.Int("data_keys", len(createdConfigMap.Data)))

	return createdConfigMap, nil
}

// GetConfigMap 获取ConfigMap
// 根据命名空间和名称获取ConfigMap信息
func (c *Client) GetConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
	if namespace == "" || name == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	configMap, err := c.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to get ConfigMap",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return nil, WrapError(ErrPodNotFound, fmt.Sprintf("failed to get ConfigMap: %v", err))
	}

	return configMap, nil
}

// UpdateConfigMap 更新ConfigMap
// 更新现有ConfigMap的数据
func (c *Client) UpdateConfigMap(namespace, name string, data map[string]string) (*corev1.ConfigMap, error) {
	if namespace == "" || name == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	// 获取现有ConfigMap
	configMap, err := c.GetConfigMap(namespace, name)
	if err != nil {
		return nil, err
	}

	// 更新数据
	if data != nil {
		configMap.Data = data
	}

	// 调用K8s API更新ConfigMap
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	updatedConfigMap, err := c.clientset.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to update ConfigMap",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return nil, WrapError(err, fmt.Sprintf("failed to update ConfigMap: %v", err))
	}

	logger.GetLogger().Info("ConfigMap updated successfully",
		zap.String("name", updatedConfigMap.Name),
		zap.String("namespace", updatedConfigMap.Namespace))

	return updatedConfigMap, nil
}

// DeleteConfigMap 删除ConfigMap
// 根据命名空间和名称删除ConfigMap
func (c *Client) DeleteConfigMap(namespace, name string) error {
	if namespace == "" || name == "" {
		return WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	err := c.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to delete ConfigMap",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return WrapError(ErrPodDeletionFailed, fmt.Sprintf("failed to delete ConfigMap: %v", err))
	}

	logger.GetLogger().Info("ConfigMap deleted successfully",
		zap.String("name", name),
		zap.String("namespace", namespace))

	return nil
}
