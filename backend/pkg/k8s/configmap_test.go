package k8s

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClient_CreateConfigMap(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	configMap, err := client.CreateConfigMap("default", "test-configmap", data)
	if err != nil {
		t.Fatalf("CreateConfigMap() error = %v", err)
	}

	if configMap.Name != "test-configmap" {
		t.Errorf("CreateConfigMap() name = %v, want test-configmap", configMap.Name)
	}
	if len(configMap.Data) != 2 {
		t.Errorf("CreateConfigMap() data count = %v, want 2", len(configMap.Data))
	}
}

func TestClient_GetConfigMap(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
		Data: map[string]string{
			"key1": "value1",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	configMap, err := client.GetConfigMap("default", "test-configmap")
	if err != nil {
		t.Fatalf("GetConfigMap() error = %v", err)
	}

	if configMap.Name != "test-configmap" {
		t.Errorf("GetConfigMap() name = %v, want test-configmap", configMap.Name)
	}
	if configMap.Data["key1"] != "value1" {
		t.Errorf("GetConfigMap() data[key1] = %v, want value1", configMap.Data["key1"])
	}
}

func TestClient_UpdateConfigMap(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
		Data: map[string]string{
			"key1": "value1",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	newData := map[string]string{
		"key1": "updated-value1",
		"key2": "value2",
	}

	configMap, err := client.UpdateConfigMap("default", "test-configmap", newData)
	if err != nil {
		t.Fatalf("UpdateConfigMap() error = %v", err)
	}

	if len(configMap.Data) != 2 {
		t.Errorf("UpdateConfigMap() data count = %v, want 2", len(configMap.Data))
	}
	if configMap.Data["key1"] != "updated-value1" {
		t.Errorf("UpdateConfigMap() data[key1] = %v, want updated-value1", configMap.Data["key1"])
	}
}

func TestClient_DeleteConfigMap(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	err := client.DeleteConfigMap("default", "test-configmap")
	if err != nil {
		t.Fatalf("DeleteConfigMap() error = %v", err)
	}

	_, err = client.GetConfigMap("default", "test-configmap")
	if err == nil {
		t.Error("DeleteConfigMap() ConfigMap still exists after deletion")
	}
}
