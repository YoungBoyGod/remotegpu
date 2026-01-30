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

	t.Run("create with data", func(t *testing.T) {
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
	})

	t.Run("create with nil data", func(t *testing.T) {
		configMap, err := client.CreateConfigMap("default", "test-configmap-nil", nil)
		if err != nil {
			t.Fatalf("CreateConfigMap() error = %v", err)
		}
		if configMap.Data == nil {
			t.Errorf("CreateConfigMap() data should not be nil")
		}
	})

	t.Run("missing namespace", func(t *testing.T) {
		_, err := client.CreateConfigMap("", "test-configmap", nil)
		if err == nil {
			t.Errorf("CreateConfigMap() should return error for empty namespace")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		_, err := client.CreateConfigMap("default", "", nil)
		if err == nil {
			t.Errorf("CreateConfigMap() should return error for empty name")
		}
	})
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

	t.Run("update with new data", func(t *testing.T) {
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
	})

	t.Run("update with empty namespace", func(t *testing.T) {
		_, err := client.UpdateConfigMap("", "test-configmap", map[string]string{"key": "value"})
		if err == nil {
			t.Error("UpdateConfigMap() should return error for empty namespace")
		}
	})

	t.Run("update with empty name", func(t *testing.T) {
		_, err := client.UpdateConfigMap("default", "", map[string]string{"key": "value"})
		if err == nil {
			t.Error("UpdateConfigMap() should return error for empty name")
		}
	})
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

func TestClient_DeleteConfigMap_ErrorCases(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	tests := []struct {
		name      string
		namespace string
		cmName    string
		wantErr   bool
	}{
		{
			name:      "empty namespace",
			namespace: "",
			cmName:    "test-configmap",
			wantErr:   true,
		},
		{
			name:      "empty name",
			namespace: "default",
			cmName:    "",
			wantErr:   true,
		},
		{
			name:      "both empty",
			namespace: "",
			cmName:    "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteConfigMap(tt.namespace, tt.cmName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteConfigMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
