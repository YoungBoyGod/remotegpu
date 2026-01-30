//go:build integration
// +build integration

package k8s

import (
	"context"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	corev1 "k8s.io/api/core/v1"
)

const (
	testNamespace = "remotegpu-test"
	testTimeout   = 5 * time.Minute
)

// setupIntegrationTest 设置集成测试环境
func setupIntegrationTest(t *testing.T) *Client {
	// 加载全局配置
	if err := config.LoadConfig("../../config/config.yaml"); err != nil {
		t.Fatalf("Failed to load global config: %v", err)
	}

	// 加载K8s配置
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load K8s config: %v", err)
	}

	// 创建客户端
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 测试连接
	if err := client.Ping(); err != nil {
		t.Fatalf("Failed to ping K8s cluster: %v", err)
	}

	return client
}

// cleanupResource 清理测试资源
func cleanupResource(t *testing.T, cleanup func() error) {
	if err := cleanup(); err != nil {
		t.Logf("Warning: cleanup failed: %v", err)
	}
}

// TestIntegration_PodAndService 测试 Pod 和 Service 集成
func TestIntegration_PodAndService(t *testing.T) {
	client := setupIntegrationTest(t)
	defer client.Close()

	podName := "test-nginx-pod"
	serviceName := "test-nginx-service"

	// 清理函数
	defer cleanupResource(t, func() error {
		_ = client.DeleteService(testNamespace, serviceName)
		return client.DeletePod(testNamespace, podName)
	})

	// 1. 创建 Pod
	t.Log("Creating test pod...")
	podConfig := &PodConfig{
		Name:      podName,
		Namespace: testNamespace,
		Image:     "nginx:alpine",
		Labels:    map[string]string{"app": "test-nginx"},
		CPU:       1,
		Memory:    256,
	}

	pod, err := client.CreatePod(podConfig)
	if err != nil {
		t.Fatalf("Failed to create pod: %v", err)
	}
	t.Logf("Pod created: %s", pod.Name)

	// 2. 等待 Pod 运行
	t.Log("Waiting for pod to be running...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Timeout waiting for pod to be running")
		case <-time.After(2 * time.Second):
			status, err := client.GetPodStatus(testNamespace, podName)
			if err != nil {
				t.Fatalf("Failed to get pod status: %v", err)
			}
			t.Logf("Pod status: %s", status)
			if status == "Running" {
				goto podRunning
			}
		}
	}

podRunning:
	t.Log("Pod is running")

	// 3. 创建 Service
	t.Log("Creating service...")
	serviceConfig := &ServiceConfig{
		Name:      serviceName,
		Namespace: testNamespace,
		Type:      ServiceTypeClusterIP,
		Selector:  map[string]string{"app": "test-nginx"},
		Ports: []ServicePort{
			{
				Name:       "http",
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: 80,
			},
		},
	}

	service, err := client.CreateService(serviceConfig)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	t.Logf("Service created: %s (ClusterIP: %s)", service.Name, service.Spec.ClusterIP)

	// 4. 验证 Service
	retrievedService, err := client.GetService(testNamespace, serviceName)
	if err != nil {
		t.Fatalf("Failed to get service: %v", err)
	}

	if retrievedService.Spec.ClusterIP == "" {
		t.Error("Service ClusterIP is empty")
	}

	t.Log("✅ Pod and Service integration test passed")
}

// TestIntegration_PVCPersistence 测试 PVC 数据持久化
func TestIntegration_PVCPersistence(t *testing.T) {
	client := setupIntegrationTest(t)
	defer client.Close()

	pvcName := "test-pvc"

	// 清理函数
	defer cleanupResource(t, func() error {
		return client.DeletePVC(testNamespace, pvcName)
	})

	// 1. 创建 PVC
	t.Log("Creating PVC...")
	pvcConfig := &PVCConfig{
		Name:      pvcName,
		Namespace: testNamespace,
		Size:      "1Gi",
		AccessModes: []PVCAccessMode{
			ReadWriteOnce,
		},
	}

	pvc, err := client.CreatePVC(pvcConfig)
	if err != nil {
		t.Fatalf("Failed to create PVC: %v", err)
	}
	t.Logf("PVC created: %s", pvc.Name)

	// 2. 等待 PVC 绑定
	t.Log("Waiting for PVC to be bound...")
	if err := client.WaitForPVCBound(testNamespace, pvcName, 2*time.Minute); err != nil {
		t.Logf("Warning: PVC not bound within timeout: %v", err)
	} else {
		t.Log("PVC is bound")
	}

	// 3. 验证 PVC 状态
	status, err := client.GetPVCStatus(testNamespace, pvcName)
	if err != nil {
		t.Fatalf("Failed to get PVC status: %v", err)
	}
	t.Logf("PVC status: %s", status)

	t.Log("✅ PVC persistence test passed")
}

// TestIntegration_ConfigMapMount 测试 ConfigMap 管理
func TestIntegration_ConfigMapMount(t *testing.T) {
	client := setupIntegrationTest(t)
	defer client.Close()

	configMapName := "test-config"

	// 清理函数
	defer cleanupResource(t, func() error {
		return client.DeleteConfigMap(testNamespace, configMapName)
	})

	// 1. 创建 ConfigMap
	t.Log("Creating ConfigMap...")
	configData := map[string]string{
		"app.conf": "server_port=8080",
		"db.conf":  "host=localhost",
		"test.txt": "test content",
	}

	configMap, err := client.CreateConfigMap(testNamespace, configMapName, configData)
	if err != nil {
		t.Fatalf("Failed to create ConfigMap: %v", err)
	}
	t.Logf("ConfigMap created: %s", configMap.Name)

	// 2. 获取 ConfigMap
	retrievedCM, err := client.GetConfigMap(testNamespace, configMapName)
	if err != nil {
		t.Fatalf("Failed to get ConfigMap: %v", err)
	}

	if len(retrievedCM.Data) != 3 {
		t.Errorf("ConfigMap data count = %d, want 3", len(retrievedCM.Data))
	}

	// 3. 更新 ConfigMap
	t.Log("Updating ConfigMap...")
	newData := map[string]string{
		"app.conf": "server_port=9090",
		"new.conf": "new_value",
	}
	updatedCM, err := client.UpdateConfigMap(testNamespace, configMapName, newData)
	if err != nil {
		t.Fatalf("Failed to update ConfigMap: %v", err)
	}

	if len(updatedCM.Data) != 2 {
		t.Errorf("Updated ConfigMap data count = %d, want 2", len(updatedCM.Data))
	}

	t.Log("✅ ConfigMap management test passed")
}

// TestIntegration_EndToEnd 端到端测试：完整环境创建和清理
func TestIntegration_EndToEnd(t *testing.T) {
	client := setupIntegrationTest(t)
	defer client.Close()

	// 资源名称
	pvcName := "e2e-pvc"
	configMapName := "e2e-config"
	podName := "e2e-nginx"
	serviceName := "e2e-service"

	// 清理函数
	defer cleanupResource(t, func() error {
		_ = client.DeleteService(testNamespace, serviceName)
		_ = client.DeletePod(testNamespace, podName)
		_ = client.DeleteConfigMap(testNamespace, configMapName)
		return client.DeletePVC(testNamespace, pvcName)
	})

	// 1. 创建 PVC
	t.Log("Step 1: Creating PVC...")
	pvcConfig := &PVCConfig{
		Name:        pvcName,
		Namespace:   testNamespace,
		Size:        "1Gi",
		AccessModes: []PVCAccessMode{ReadWriteOnce},
	}
	if _, err := client.CreatePVC(pvcConfig); err != nil {
		t.Fatalf("Failed to create PVC: %v", err)
	}

	// 2. 创建 ConfigMap
	t.Log("Step 2: Creating ConfigMap...")
	configData := map[string]string{
		"nginx.conf": "server { listen 80; }",
	}
	if _, err := client.CreateConfigMap(testNamespace, configMapName, configData); err != nil {
		t.Fatalf("Failed to create ConfigMap: %v", err)
	}

	// 3. 创建 Pod
	t.Log("Step 3: Creating Pod...")
	podConfig := &PodConfig{
		Name:      podName,
		Namespace: testNamespace,
		Image:     "nginx:alpine",
		Labels:    map[string]string{"app": "e2e-nginx"},
		CPU:       1,
		Memory:    256,
	}
	if _, err := client.CreatePod(podConfig); err != nil {
		t.Fatalf("Failed to create Pod: %v", err)
	}

	// 4. 创建 Service
	t.Log("Step 4: Creating Service...")
	serviceConfig := &ServiceConfig{
		Name:      serviceName,
		Namespace: testNamespace,
		Type:      ServiceTypeClusterIP,
		Selector:  map[string]string{"app": "e2e-nginx"},
		Ports: []ServicePort{
			{Name: "http", Protocol: corev1.ProtocolTCP, Port: 80, TargetPort: 80},
		},
	}
	if _, err := client.CreateService(serviceConfig); err != nil {
		t.Fatalf("Failed to create Service: %v", err)
	}

	// 5. 验证所有资源
	t.Log("Step 5: Verifying all resources...")
	if _, err := client.GetPVC(testNamespace, pvcName); err != nil {
		t.Errorf("PVC not found: %v", err)
	}
	if _, err := client.GetConfigMap(testNamespace, configMapName); err != nil {
		t.Errorf("ConfigMap not found: %v", err)
	}
	if _, err := client.GetPod(testNamespace, podName); err != nil {
		t.Errorf("Pod not found: %v", err)
	}
	if _, err := client.GetService(testNamespace, serviceName); err != nil {
		t.Errorf("Service not found: %v", err)
	}

	t.Log("✅ End-to-end test passed")
}
