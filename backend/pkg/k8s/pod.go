package k8s

import (
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/YoungBoyGod/remotegpu/pkg/logger"
)

// PodConfig Pod配置
type PodConfig struct {
	// Name Pod名称
	Name string

	// Namespace 命名空间
	Namespace string

	// Image 容器镜像
	Image string

	// Command 容器命令
	Command []string

	// Args 容器参数
	Args []string

	// CPU CPU核心数
	CPU int64

	// Memory 内存大小（MB）
	Memory int64

	// GPU GPU数量
	GPU int64

	// Env 环境变量
	Env map[string]string

	// Volumes 卷挂载
	Volumes []VolumeMount

	// Labels 标签
	Labels map[string]string

	// Annotations 注解
	Annotations map[string]string

	// RestartPolicy 重启策略
	RestartPolicy corev1.RestartPolicy
}

// VolumeMount 卷挂载配置
type VolumeMount struct {
	// Name 卷名称
	Name string

	// MountPath 挂载路径
	MountPath string

	// HostPath 主机路径（用于HostPath类型）
	HostPath string

	// ReadOnly 是否只读
	ReadOnly bool
}

// LogOptions 日志选项
type LogOptions struct {
	// Container 容器名称（如果Pod有多个容器）
	Container string

	// TailLines 返回最后N行日志
	TailLines int64

	// Follow 是否持续跟踪日志
	Follow bool

	// Timestamps 是否显示时间戳
	Timestamps bool

	// SinceSeconds 返回最近N秒的日志
	SinceSeconds int64
}

// CreatePod 创建Pod
func (c *Client) CreatePod(config *PodConfig) (*corev1.Pod, error) {
	if config == nil {
		return nil, WrapError(ErrInvalidConfig, "pod config is nil")
	}

	// 验证必填字段
	if config.Name == "" {
		return nil, WrapError(ErrInvalidConfig, "pod name is required")
	}
	if config.Image == "" {
		return nil, WrapError(ErrInvalidConfig, "pod image is required")
	}

	// 使用默认命名空间
	namespace := config.Namespace
	if namespace == "" {
		namespace = c.namespace
	}

	// 构建Pod对象
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: config.RestartPolicy,
			Containers: []corev1.Container{
				{
					Name:    config.Name,
					Image:   config.Image,
					Command: config.Command,
					Args:    config.Args,
				},
			},
		},
	}

	// 设置默认重启策略
	if pod.Spec.RestartPolicy == "" {
		pod.Spec.RestartPolicy = corev1.RestartPolicyNever
	}

	// 配置资源限制
	c.configureResources(pod, config)

	// 配置环境变量
	c.configureEnv(pod, config)

	// 配置卷挂载
	c.configureVolumes(pod, config)

	// 创建Pod
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	createdPod, err := c.clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, WrapErrorf(ErrPodCreationFailed, "failed to create pod %s in namespace %s: %v", config.Name, namespace, err)
	}

	logger.GetLogger().Info("Pod created successfully",
		zap.String("name", config.Name),
		zap.String("namespace", namespace))

	return createdPod, nil
}

// configureResources 配置资源限制
func (c *Client) configureResources(pod *corev1.Pod, config *PodConfig) {
	if len(pod.Spec.Containers) == 0 {
		return
	}

	container := &pod.Spec.Containers[0]
	resources := corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}

	// 配置CPU
	if config.CPU > 0 {
		cpuQuantity := resource.MustParse(fmt.Sprintf("%d", config.CPU))
		resources.Limits[corev1.ResourceCPU] = cpuQuantity
		resources.Requests[corev1.ResourceCPU] = cpuQuantity
	}

	// 配置Memory
	if config.Memory > 0 {
		memQuantity := resource.MustParse(fmt.Sprintf("%dMi", config.Memory))
		resources.Limits[corev1.ResourceMemory] = memQuantity
		resources.Requests[corev1.ResourceMemory] = memQuantity
	}

	// 配置GPU
	if config.GPU > 0 {
		gpuQuantity := resource.MustParse(fmt.Sprintf("%d", config.GPU))
		resources.Limits["nvidia.com/gpu"] = gpuQuantity
	}

	container.Resources = resources
}

// configureEnv 配置环境变量
func (c *Client) configureEnv(pod *corev1.Pod, config *PodConfig) {
	if len(pod.Spec.Containers) == 0 || len(config.Env) == 0 {
		return
	}

	container := &pod.Spec.Containers[0]
	for key, value := range config.Env {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}
}

// configureVolumes 配置卷挂载
func (c *Client) configureVolumes(pod *corev1.Pod, config *PodConfig) {
	if len(pod.Spec.Containers) == 0 || len(config.Volumes) == 0 {
		return
	}

	container := &pod.Spec.Containers[0]
	for _, vol := range config.Volumes {
		// 添加卷定义
		volume := corev1.Volume{
			Name: vol.Name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vol.HostPath,
				},
			},
		}
		pod.Spec.Volumes = append(pod.Spec.Volumes, volume)

		// 添加卷挂载
		volumeMount := corev1.VolumeMount{
			Name:      vol.Name,
			MountPath: vol.MountPath,
			ReadOnly:  vol.ReadOnly,
		}
		container.VolumeMounts = append(container.VolumeMounts, volumeMount)
	}
}

// GetPod 获取Pod
func (c *Client) GetPod(namespace, name string) (*corev1.Pod, error) {
	if name == "" {
		return nil, WrapError(ErrInvalidConfig, "pod name is required")
	}

	// 使用默认命名空间
	if namespace == "" {
		namespace = c.namespace
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, WrapErrorf(ErrPodNotFound, "failed to get pod %s in namespace %s: %v", name, namespace, err)
	}

	return pod, nil
}

// GetPodStatus 获取Pod状态
func (c *Client) GetPodStatus(namespace, name string) (string, error) {
	pod, err := c.GetPod(namespace, name)
	if err != nil {
		return "", err
	}

	return string(pod.Status.Phase), nil
}

// DeletePod 删除Pod
func (c *Client) DeletePod(namespace, name string) error {
	return c.DeletePodGracefully(namespace, name, 0)
}

// DeletePodGracefully 优雅删除Pod
func (c *Client) DeletePodGracefully(namespace, name string, gracePeriod int64) error {
	if name == "" {
		return WrapError(ErrInvalidConfig, "pod name is required")
	}

	// 使用默认命名空间
	if namespace == "" {
		namespace = c.namespace
	}

	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	deleteOptions := metav1.DeleteOptions{}
	if gracePeriod > 0 {
		deleteOptions.GracePeriodSeconds = &gracePeriod
	}

	err := c.clientset.CoreV1().Pods(namespace).Delete(ctx, name, deleteOptions)
	if err != nil {
		return WrapErrorf(ErrPodDeletionFailed, "failed to delete pod %s in namespace %s: %v", name, namespace, err)
	}

	logger.GetLogger().Info("Pod deleted successfully",
		zap.String("name", name),
		zap.String("namespace", namespace))

	return nil
}

// GetPodLogs 获取Pod日志
func (c *Client) GetPodLogs(namespace, name string, opts *LogOptions) (string, error) {
	if name == "" {
		return "", WrapError(ErrInvalidConfig, "pod name is required")
	}

	// 使用默认命名空间
	if namespace == "" {
		namespace = c.namespace
	}

	// 构建日志选项
	podLogOpts := &corev1.PodLogOptions{}
	if opts != nil {
		if opts.Container != "" {
			podLogOpts.Container = opts.Container
		}
		if opts.TailLines > 0 {
			podLogOpts.TailLines = &opts.TailLines
		}
		podLogOpts.Follow = opts.Follow
		podLogOpts.Timestamps = opts.Timestamps
		if opts.SinceSeconds > 0 {
			podLogOpts.SinceSeconds = &opts.SinceSeconds
		}
	}

	// 获取日志流
	ctx, cancel := c.GetContextWithTimeout(c.config.Timeout)
	defer cancel()

	req := c.clientset.CoreV1().Pods(namespace).GetLogs(name, podLogOpts)
	logStream, err := req.Stream(ctx)
	if err != nil {
		return "", WrapErrorf(ErrLogsFetchFailed, "failed to get logs for pod %s in namespace %s: %v", name, namespace, err)
	}
	defer logStream.Close()

	// 读取日志内容
	logs, err := io.ReadAll(logStream)
	if err != nil {
		return "", WrapErrorf(ErrLogsFetchFailed, "failed to read logs for pod %s: %v", name, err)
	}

	return string(logs), nil
}

// WatchPodStatus 监控Pod状态变化
func (c *Client) WatchPodStatus(namespace, name string, callback func(status string)) error {
	if name == "" {
		return WrapError(ErrInvalidConfig, "pod name is required")
	}
	if callback == nil {
		return WrapError(ErrInvalidConfig, "callback function is required")
	}

	// 使用默认命名空间
	if namespace == "" {
		namespace = c.namespace
	}

	// 创建Watch
	ctx, cancel := c.GetContextWithTimeout(5 * time.Minute) // 使用较长的超时时间
	defer cancel()

	watchOpts := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	}

	watcher, err := c.clientset.CoreV1().Pods(namespace).Watch(ctx, watchOpts)
	if err != nil {
		return WrapErrorf(ErrConnectionFailed, "failed to watch pod %s in namespace %s: %v", name, namespace, err)
	}
	defer watcher.Stop()

	logger.GetLogger().Info("Started watching pod status",
		zap.String("name", name),
		zap.String("namespace", namespace))

	// 监听事件
	for event := range watcher.ResultChan() {
		if event.Type == watch.Error {
			return WrapError(ErrPodStatusTimeout, "watch error occurred")
		}

		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}

		status := string(pod.Status.Phase)
		callback(status)

		// 如果Pod已经完成或失败，停止监控
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			logger.GetLogger().Info("Pod reached terminal state",
				zap.String("name", name),
				zap.String("status", status))
			break
		}
	}

	return nil
}
