package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// TestCreatePod 测试创建Pod
func TestCreatePod(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	t.Run("create pod with basic config", func(t *testing.T) {
		config := &PodConfig{
			Name:      "test-pod",
			Namespace: "default",
			Image:     "nginx:latest",
		}

		pod, err := client.CreatePod(config)
		assert.NoError(t, err)
		assert.NotNil(t, pod)
		assert.Equal(t, "test-pod", pod.Name)
		assert.Equal(t, "default", pod.Namespace)
		assert.Equal(t, "nginx:latest", pod.Spec.Containers[0].Image)
	})

	t.Run("create pod with resources", func(t *testing.T) {
		config := &PodConfig{
			Name:      "test-pod-resources",
			Namespace: "default",
			Image:     "nginx:latest",
			CPU:       2,
			Memory:    1024,
			GPU:       1,
		}

		pod, err := client.CreatePod(config)
		assert.NoError(t, err)
		assert.NotNil(t, pod)

		resources := pod.Spec.Containers[0].Resources
		assert.NotNil(t, resources.Limits)
		assert.NotNil(t, resources.Limits[corev1.ResourceCPU])
		assert.NotNil(t, resources.Limits[corev1.ResourceMemory])
		assert.NotNil(t, resources.Limits["nvidia.com/gpu"])
	})

	t.Run("create pod with env vars", func(t *testing.T) {
		config := &PodConfig{
			Name:      "test-pod-env",
			Namespace: "default",
			Image:     "nginx:latest",
			Env: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		}

		pod, err := client.CreatePod(config)
		assert.NoError(t, err)
		assert.NotNil(t, pod)
		assert.Len(t, pod.Spec.Containers[0].Env, 2)
	})

	t.Run("nil config", func(t *testing.T) {
		pod, err := client.CreatePod(nil)
		assert.Error(t, err)
		assert.Nil(t, pod)
		assert.Contains(t, err.Error(), "pod config is nil")
	})

	t.Run("missing name", func(t *testing.T) {
		config := &PodConfig{
			Image: "nginx:latest",
		}
		pod, err := client.CreatePod(config)
		assert.Error(t, err)
		assert.Nil(t, pod)
		assert.Contains(t, err.Error(), "pod name is required")
	})

	t.Run("missing image", func(t *testing.T) {
		config := &PodConfig{
			Name: "test-pod",
		}
		pod, err := client.CreatePod(config)
		assert.Error(t, err)
		assert.Nil(t, pod)
		assert.Contains(t, err.Error(), "pod image is required")
	})
}

// TestGetPod 测试获取Pod
func TestGetPod(t *testing.T) {
	// 创建一个已存在的Pod
	existingPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}

	fakeClientset := fake.NewSimpleClientset(existingPod)
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	t.Run("get existing pod", func(t *testing.T) {
		pod, err := client.GetPod("default", "existing-pod")
		assert.NoError(t, err)
		assert.NotNil(t, pod)
		assert.Equal(t, "existing-pod", pod.Name)
	})

	t.Run("get non-existent pod", func(t *testing.T) {
		pod, err := client.GetPod("default", "non-existent")
		assert.Error(t, err)
		assert.Nil(t, pod)
	})

	t.Run("missing pod name", func(t *testing.T) {
		pod, err := client.GetPod("default", "")
		assert.Error(t, err)
		assert.Nil(t, pod)
		assert.Contains(t, err.Error(), "pod name is required")
	})
}

// TestGetPodStatus 测试获取Pod状态
func TestGetPodStatus(t *testing.T) {
	existingPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}

	fakeClientset := fake.NewSimpleClientset(existingPod)
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	status, err := client.GetPodStatus("default", "test-pod")
	assert.NoError(t, err)
	assert.Equal(t, "Running", status)
}

// TestDeletePod 测试删除Pod
func TestDeletePod(t *testing.T) {
	existingPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	fakeClientset := fake.NewSimpleClientset(existingPod)
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	t.Run("delete existing pod", func(t *testing.T) {
		err := client.DeletePod("default", "test-pod")
		assert.NoError(t, err)
	})

	t.Run("missing pod name", func(t *testing.T) {
		err := client.DeletePod("default", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pod name is required")
	})
}

// TestDeletePodGracefully 测试优雅删除Pod
func TestDeletePodGracefully(t *testing.T) {
	existingPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	fakeClientset := fake.NewSimpleClientset(existingPod)
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	err := client.DeletePodGracefully("default", "test-pod", 30)
	assert.NoError(t, err)
}

// TestGetPodLogs 测试获取Pod日志
func TestGetPodLogs(t *testing.T) {
	existingPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}

	fakeClientset := fake.NewSimpleClientset(existingPod)
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	t.Run("get logs without options", func(t *testing.T) {
		logs, err := client.GetPodLogs("default", "test-pod", nil)
		assert.NoError(t, err)
		assert.NotNil(t, logs)
	})

	t.Run("get logs with options", func(t *testing.T) {
		opts := &LogOptions{
			TailLines:  100,
			Timestamps: true,
		}
		logs, err := client.GetPodLogs("default", "test-pod", opts)
		assert.NoError(t, err)
		assert.NotNil(t, logs)
	})

	t.Run("missing pod name", func(t *testing.T) {
		logs, err := client.GetPodLogs("default", "", nil)
		assert.Error(t, err)
		assert.Empty(t, logs)
		assert.Contains(t, err.Error(), "pod name is required")
	})
}

// TestCreatePodWithVolumes 测试创建带卷挂载的Pod
func TestCreatePodWithVolumes(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	config := &PodConfig{
		Name:      "test-pod-volumes",
		Namespace: "default",
		Image:     "nginx:latest",
		Volumes: []VolumeMount{
			{
				Name:      "data",
				MountPath: "/data",
				HostPath:  "/host/data",
				ReadOnly:  false,
			},
		},
	}

	pod, err := client.CreatePod(config)
	assert.NoError(t, err)
	assert.NotNil(t, pod)
	assert.Len(t, pod.Spec.Volumes, 1)
	assert.Len(t, pod.Spec.Containers[0].VolumeMounts, 1)
}
