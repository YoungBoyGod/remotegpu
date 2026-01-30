package k8s

import (
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestIngressConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *IngressConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &IngressConfig{
				Name:      "test-ingress",
				Namespace: "default",
				Rules: []IngressRule{
					{
						Host:        "example.com",
						Path:        "/",
						PathType:    PathTypePrefix,
						ServiceName: "test-service",
						ServicePort: 80,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: &IngressConfig{
				Namespace: "default",
				Rules: []IngressRule{
					{
						Host:        "example.com",
						ServiceName: "test-service",
						ServicePort: 80,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			config: &IngressConfig{
				Name: "test-ingress",
				Rules: []IngressRule{
					{
						Host:        "example.com",
						ServiceName: "test-service",
						ServicePort: 80,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing rules",
			config: &IngressConfig{
				Name:      "test-ingress",
				Namespace: "default",
				Rules:     []IngressRule{},
			},
			wantErr: true,
		},
		{
			name: "missing host",
			config: &IngressConfig{
				Name:      "test-ingress",
				Namespace: "default",
				Rules: []IngressRule{
					{
						ServiceName: "test-service",
						ServicePort: 80,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing service name",
			config: &IngressConfig{
				Name:      "test-ingress",
				Namespace: "default",
				Rules: []IngressRule{
					{
						Host:        "example.com",
						ServicePort: 80,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IngressConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateIngress(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	t.Run("create with single rule", func(t *testing.T) {
		config := &IngressConfig{
			Name:      "test-ingress",
			Namespace: "default",
			Rules: []IngressRule{
				{
					Host:        "example.com",
					Path:        "/api",
					PathType:    PathTypePrefix,
					ServiceName: "test-service",
					ServicePort: 80,
				},
			},
			Labels: map[string]string{"env": "test"},
		}

		ingress, err := client.CreateIngress(config)
		if err != nil {
			t.Fatalf("CreateIngress() error = %v", err)
		}

		if ingress.Name != config.Name {
			t.Errorf("CreateIngress() name = %v, want %v", ingress.Name, config.Name)
		}
		if len(ingress.Spec.Rules) != 1 {
			t.Errorf("CreateIngress() rules count = %v, want 1", len(ingress.Spec.Rules))
		}
	})

	t.Run("create with multiple rules", func(t *testing.T) {
		config := &IngressConfig{
			Name:      "test-multi-ingress",
			Namespace: "default",
			Rules: []IngressRule{
				{
					Host:        "api.example.com",
					Path:        "/",
					PathType:    PathTypePrefix,
					ServiceName: "api-service",
					ServicePort: 80,
				},
				{
					Host:        "web.example.com",
					Path:        "/",
					PathType:    PathTypePrefix,
					ServiceName: "web-service",
					ServicePort: 8080,
				},
			},
		}

		ingress, err := client.CreateIngress(config)
		if err != nil {
			t.Fatalf("CreateIngress() error = %v", err)
		}

		if len(ingress.Spec.Rules) != 2 {
			t.Errorf("CreateIngress() rules count = %v, want 2", len(ingress.Spec.Rules))
		}
	})

	t.Run("create with invalid config", func(t *testing.T) {
		config := &IngressConfig{
			Name:      "",
			Namespace: "default",
			Rules:     []IngressRule{{Host: "example.com"}},
		}

		_, err := client.CreateIngress(config)
		if err == nil {
			t.Error("CreateIngress() should return error for invalid config")
		}
	})
}

func TestClient_GetIngress(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ingress",
			Namespace: "default",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	ingress, err := client.GetIngress("default", "test-ingress")
	if err != nil {
		t.Fatalf("GetIngress() error = %v", err)
	}

	if ingress.Name != "test-ingress" {
		t.Errorf("GetIngress() name = %v, want test-ingress", ingress.Name)
	}
}

func TestClient_DeleteIngress(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ingress",
			Namespace: "default",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	err := client.DeleteIngress("default", "test-ingress")
	if err != nil {
		t.Fatalf("DeleteIngress() error = %v", err)
	}

	_, err = client.GetIngress("default", "test-ingress")
	if err == nil {
		t.Error("DeleteIngress() ingress still exists after deletion")
	}
}

func TestClient_DeleteIngress_ErrorCases(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	tests := []struct {
		name        string
		namespace   string
		ingressName string
		wantErr     bool
	}{
		{
			name:        "empty namespace",
			namespace:   "",
			ingressName: "test-ingress",
			wantErr:     true,
		},
		{
			name:        "empty name",
			namespace:   "default",
			ingressName: "",
			wantErr:     true,
		},
		{
			name:        "both empty",
			namespace:   "",
			ingressName: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteIngress(tt.namespace, tt.ingressName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteIngress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
