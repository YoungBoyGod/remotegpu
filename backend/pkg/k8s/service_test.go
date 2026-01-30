package k8s

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestServiceConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ServiceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &ServiceConfig{
				Name:      "test-service",
				Namespace: "default",
				Type:      ServiceTypeClusterIP,
				Selector:  map[string]string{"app": "test"},
				Ports: []ServicePort{
					{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
						TargetPort: 8080,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: &ServiceConfig{
				Namespace: "default",
				Selector:  map[string]string{"app": "test"},
				Ports: []ServicePort{
					{Port: 80, TargetPort: 8080},
				},
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			config: &ServiceConfig{
				Name:     "test-service",
				Selector: map[string]string{"app": "test"},
				Ports: []ServicePort{
					{Port: 80, TargetPort: 8080},
				},
			},
			wantErr: true,
		},
		{
			name: "missing selector",
			config: &ServiceConfig{
				Name:      "test-service",
				Namespace: "default",
				Ports: []ServicePort{
					{Port: 80, TargetPort: 8080},
				},
			},
			wantErr: true,
		},
		{
			name: "missing ports",
			config: &ServiceConfig{
				Name:      "test-service",
				Namespace: "default",
				Selector:  map[string]string{"app": "test"},
				Ports:     []ServicePort{},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: &ServiceConfig{
				Name:      "test-service",
				Namespace: "default",
				Selector:  map[string]string{"app": "test"},
				Ports: []ServicePort{
					{Port: -1, TargetPort: 8080},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid target port",
			config: &ServiceConfig{
				Name:      "test-service",
				Namespace: "default",
				Selector:  map[string]string{"app": "test"},
				Ports: []ServicePort{
					{Port: 80, TargetPort: -1},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateService(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	t.Run("create ClusterIP service", func(t *testing.T) {
		config := &ServiceConfig{
			Name:      "test-service",
			Namespace: "default",
			Type:      ServiceTypeClusterIP,
			Selector:  map[string]string{"app": "test"},
			Ports: []ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: 8080,
				},
			},
			Labels: map[string]string{"env": "test"},
		}

		service, err := client.CreateService(config)
		if err != nil {
			t.Fatalf("CreateService() error = %v", err)
		}

		if service.Name != config.Name {
			t.Errorf("CreateService() name = %v, want %v", service.Name, config.Name)
		}
		if service.Namespace != config.Namespace {
			t.Errorf("CreateService() namespace = %v, want %v", service.Namespace, config.Namespace)
		}
		if len(service.Spec.Ports) != 1 {
			t.Errorf("CreateService() ports count = %v, want 1", len(service.Spec.Ports))
		}
	})

	t.Run("create NodePort service", func(t *testing.T) {
		config := &ServiceConfig{
			Name:      "test-nodeport",
			Namespace: "default",
			Type:      ServiceTypeNodePort,
			Selector:  map[string]string{"app": "test"},
			Ports: []ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: 8080,
				},
			},
		}

		service, err := client.CreateService(config)
		if err != nil {
			t.Fatalf("CreateService() error = %v", err)
		}

		if string(service.Spec.Type) != string(corev1.ServiceTypeNodePort) {
			t.Errorf("CreateService() type = %v, want NodePort", service.Spec.Type)
		}
	})

	t.Run("create with invalid config", func(t *testing.T) {
		config := &ServiceConfig{
			Name:      "",
			Namespace: "default",
			Selector:  map[string]string{"app": "test"},
			Ports:     []ServicePort{{Port: 80}},
		}

		_, err := client.CreateService(config)
		if err == nil {
			t.Error("CreateService() should return error for invalid config")
		}
	})
}

func TestClient_GetService(t *testing.T) {
	// 创建fake clientset并预先创建一个Service
	fakeClientset := fake.NewSimpleClientset(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	service, err := client.GetService("default", "test-service")
	if err != nil {
		t.Fatalf("GetService() error = %v", err)
	}

	if service.Name != "test-service" {
		t.Errorf("GetService() name = %v, want test-service", service.Name)
	}
}

func TestClient_DeleteService(t *testing.T) {
	// 创建fake clientset并预先创建一个Service
	fakeClientset := fake.NewSimpleClientset(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	err := client.DeleteService("default", "test-service")
	if err != nil {
		t.Fatalf("DeleteService() error = %v", err)
	}

	// 验证Service已被删除
	_, err = client.GetService("default", "test-service")
	if err == nil {
		t.Error("DeleteService() service still exists after deletion")
	}
}

func TestClient_DeleteService_ErrorCases(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	tests := []struct {
		name      string
		namespace string
		svcName   string
		wantErr   bool
	}{
		{
			name:      "empty namespace",
			namespace: "",
			svcName:   "test-service",
			wantErr:   true,
		},
		{
			name:      "empty name",
			namespace: "default",
			svcName:   "",
			wantErr:   true,
		},
		{
			name:      "both empty",
			namespace: "",
			svcName:   "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteService(tt.namespace, tt.svcName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListServices(t *testing.T) {
	// 创建fake clientset并预先创建多个Service
	fakeClientset := fake.NewSimpleClientset(
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service1",
				Namespace: "default",
				Labels:    map[string]string{"app": "test"},
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service2",
				Namespace: "default",
				Labels:    map[string]string{"app": "test"},
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service3",
				Namespace: "default",
				Labels:    map[string]string{"app": "other"},
			},
		},
	)
	client := NewClientWithClientset(fakeClientset, "default")

	// 测试列出所有Service
	services, err := client.ListServices("default", "")
	if err != nil {
		t.Fatalf("ListServices() error = %v", err)
	}
	if len(services.Items) != 3 {
		t.Errorf("ListServices() count = %v, want 3", len(services.Items))
	}

	// 测试使用标签选择器
	services, err = client.ListServices("default", "app=test")
	if err != nil {
		t.Fatalf("ListServices() with selector error = %v", err)
	}
	if len(services.Items) != 2 {
		t.Errorf("ListServices() with selector count = %v, want 2", len(services.Items))
	}

	// 测试空命名空间错误
	_, err = client.ListServices("", "")
	if err == nil {
		t.Error("ListServices() should return error for empty namespace")
	}
}

func TestClient_UpdateService(t *testing.T) {
	// 创建fake clientset并预先创建一个Service
	fakeClientset := fake.NewSimpleClientset(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
			Labels:    map[string]string{"version": "v1"},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	// 获取Service并更新
	service, err := client.GetService("default", "test-service")
	if err != nil {
		t.Fatalf("GetService() error = %v", err)
	}

	service.Labels["version"] = "v2"
	updatedService, err := client.UpdateService(service)
	if err != nil {
		t.Fatalf("UpdateService() error = %v", err)
	}

	if updatedService.Labels["version"] != "v2" {
		t.Errorf("UpdateService() label = %v, want v2", updatedService.Labels["version"])
	}
}

func TestClient_UpdateService_ErrorCases(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	tests := []struct {
		name    string
		service *corev1.Service
		wantErr bool
	}{
		{
			name:    "nil service",
			service: nil,
			wantErr: true,
		},
		{
			name: "empty name",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
		{
			name: "empty namespace",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.UpdateService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
