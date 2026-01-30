// Package k8s 提供Kubernetes存储管理功能
package k8s

import (
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVCAccessMode PVC访问模式
type PVCAccessMode string

const (
	// ReadWriteOnce 单节点读写
	ReadWriteOnce PVCAccessMode = "ReadWriteOnce"
	// ReadOnlyMany 多节点只读
	ReadOnlyMany PVCAccessMode = "ReadOnlyMany"
	// ReadWriteMany 多节点读写
	ReadWriteMany PVCAccessMode = "ReadWriteMany"
)

// PVCConfig PVC配置
type PVCConfig struct {
	// Name PVC名称
	Name string
	// Namespace 命名空间
	Namespace string
	// StorageClassName 存储类名称
	StorageClassName string
	// Size 存储大小（如 "10Gi"）
	Size string
	// AccessModes 访问模式列表
	AccessModes []PVCAccessMode
	// Labels 标签
	Labels map[string]string
	// Annotations 注解
	Annotations map[string]string
}

// Validate 验证PVC配置
func (c *PVCConfig) Validate() error {
	if c.Name == "" {
		return WrapError(ErrInvalidConfig, "pvc name is required")
	}
	if c.Namespace == "" {
		return WrapError(ErrInvalidConfig, "namespace is required")
	}
	if c.Size == "" {
		return WrapError(ErrInvalidConfig, "size is required")
	}

	// 验证存储大小格式
	if _, err := resource.ParseQuantity(c.Size); err != nil {
		return WrapError(ErrInvalidConfig, fmt.Sprintf("invalid size format: %v", err))
	}

	// 默认访问模式
	if len(c.AccessModes) == 0 {
		c.AccessModes = []PVCAccessMode{ReadWriteOnce}
	}

	return nil
}

// CreatePVC 创建PVC
// 根据配置创建Kubernetes PersistentVolumeClaim
func (c *Client) CreatePVC(config *PVCConfig) (*corev1.PersistentVolumeClaim, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 解析存储大小
	quantity, err := resource.ParseQuantity(config.Size)
	if err != nil {
		return nil, WrapError(ErrInvalidConfig, fmt.Sprintf("invalid size: %v", err))
	}

	// 构建访问模式列表
	accessModes := make([]corev1.PersistentVolumeAccessMode, 0, len(config.AccessModes))
	for _, mode := range config.AccessModes {
		accessModes = append(accessModes, corev1.PersistentVolumeAccessMode(mode))
	}

	// 构建PVC对象
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: accessModes,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: quantity,
				},
			},
		},
	}

	// 设置StorageClassName（如果指定）
	if config.StorageClassName != "" {
		pvc.Spec.StorageClassName = &config.StorageClassName
	}

	// 调用K8s API创建PVC
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	createdPVC, err := c.clientset.CoreV1().PersistentVolumeClaims(config.Namespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to create PVC",
			zap.String("name", config.Name),
			zap.String("namespace", config.Namespace),
			zap.Error(err))
		return nil, WrapError(err, fmt.Sprintf("failed to create PVC: %v", err))
	}

	logger.GetLogger().Info("PVC created successfully",
		zap.String("name", createdPVC.Name),
		zap.String("namespace", createdPVC.Namespace),
		zap.String("size", config.Size))

	return createdPVC, nil
}

// WaitForPVCBound 等待PVC绑定
// 等待PVC状态变为Bound，可选超时时间
func (c *Client) WaitForPVCBound(namespace, name string, timeout time.Duration) error {
	ctx, cancel := c.GetContextWithTimeout(timeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return WrapError(ctx.Err(), "timeout waiting for PVC to be bound")
		case <-ticker.C:
			pvc, err := c.clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return WrapError(err, "failed to get PVC status")
			}

			if pvc.Status.Phase == corev1.ClaimBound {
				logger.GetLogger().Info("PVC bound successfully",
					zap.String("name", name),
					zap.String("namespace", namespace))
				return nil
			}

			logger.GetLogger().Debug("Waiting for PVC to be bound",
				zap.String("name", name),
				zap.String("phase", string(pvc.Status.Phase)))
		}
	}
}

// GetPVC 获取PVC
// 根据命名空间和名称获取PVC信息
func (c *Client) GetPVC(namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	if namespace == "" || name == "" {
		return nil, WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	pvc, err := c.clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to get PVC",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return nil, WrapError(ErrPodNotFound, fmt.Sprintf("failed to get PVC: %v", err))
	}

	return pvc, nil
}

// GetPVCStatus 获取PVC状态
// 返回PVC的当前状态（Pending、Bound、Lost）
func (c *Client) GetPVCStatus(namespace, name string) (string, error) {
	pvc, err := c.GetPVC(namespace, name)
	if err != nil {
		return "", err
	}

	return string(pvc.Status.Phase), nil
}

// DeletePVC 删除PVC
// 根据命名空间和名称删除PVC
func (c *Client) DeletePVC(namespace, name string) error {
	if namespace == "" || name == "" {
		return WrapError(ErrInvalidConfig, "namespace and name are required")
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	err := c.clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		logger.GetLogger().Error("Failed to delete PVC",
			zap.String("name", name),
			zap.String("namespace", namespace),
			zap.Error(err))
		return WrapError(ErrPodDeletionFailed, fmt.Sprintf("failed to delete PVC: %v", err))
	}

	logger.GetLogger().Info("PVC deleted successfully",
		zap.String("name", name),
		zap.String("namespace", namespace))

	return nil
}
