package k8s

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPVCConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *PVCConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &PVCConfig{
				Name:      "test-pvc",
				Namespace: "default",
				Size:      "10Gi",
				AccessModes: []PVCAccessMode{
					ReadWriteOnce,
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: &PVCConfig{
				Namespace: "default",
				Size:      "10Gi",
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			config: &PVCConfig{
				Name: "test-pvc",
				Size: "10Gi",
			},
			wantErr: true,
		},
		{
			name: "missing size",
			config: &PVCConfig{
				Name:      "test-pvc",
				Namespace: "default",
			},
			wantErr: true,
		},
		{
			name: "invalid size format",
			config: &PVCConfig{
				Name:      "test-pvc",
				Namespace: "default",
				Size:      "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PVCConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreatePVC(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	client := NewClientWithClientset(fakeClientset, "default")

	config := &PVCConfig{
		Name:             "test-pvc",
		Namespace:        "default",
		Size:             "10Gi",
		StorageClassName: "standard",
		AccessModes: []PVCAccessMode{
			ReadWriteOnce,
		},
		Labels: map[string]string{"env": "test"},
	}

	pvc, err := client.CreatePVC(config)
	if err != nil {
		t.Fatalf("CreatePVC() error = %v", err)
	}

	if pvc.Name != config.Name {
		t.Errorf("CreatePVC() name = %v, want %v", pvc.Name, config.Name)
	}
	if pvc.Namespace != config.Namespace {
		t.Errorf("CreatePVC() namespace = %v, want %v", pvc.Namespace, config.Namespace)
	}
}

func TestClient_GetPVC(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pvc",
			Namespace: "default",
		},
		Status: corev1.PersistentVolumeClaimStatus{
			Phase: corev1.ClaimBound,
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	pvc, err := client.GetPVC("default", "test-pvc")
	if err != nil {
		t.Fatalf("GetPVC() error = %v", err)
	}

	if pvc.Name != "test-pvc" {
		t.Errorf("GetPVC() name = %v, want test-pvc", pvc.Name)
	}
}

func TestClient_GetPVCStatus(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pvc",
			Namespace: "default",
		},
		Status: corev1.PersistentVolumeClaimStatus{
			Phase: corev1.ClaimBound,
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	status, err := client.GetPVCStatus("default", "test-pvc")
	if err != nil {
		t.Fatalf("GetPVCStatus() error = %v", err)
	}

	if status != "Bound" {
		t.Errorf("GetPVCStatus() status = %v, want Bound", status)
	}
}

func TestClient_DeletePVC(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset(&corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pvc",
			Namespace: "default",
		},
	})
	client := NewClientWithClientset(fakeClientset, "default")

	err := client.DeletePVC("default", "test-pvc")
	if err != nil {
		t.Fatalf("DeletePVC() error = %v", err)
	}

	_, err = client.GetPVC("default", "test-pvc")
	if err == nil {
		t.Error("DeletePVC() PVC still exists after deletion")
	}
}
