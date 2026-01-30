package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
)

func TestGPUController_Create(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	hostSvc := service.NewHostService()
	gpuCtrl := NewGPUController()

	// 先创建测试主机
	host := &entity.Host{
		Name:           "Test Host for GPU Controller",
		IPAddress:      "192.168.1.251",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}
	if err := hostSvc.Create(host); err != nil {
		t.Fatalf("创建测试主机失败: %v", err)
	}
	defer hostSvc.Delete(host.ID)

	r.POST("/gpus", gpuCtrl.Create)
	r.DELETE("/gpus/:id", gpuCtrl.Delete)

	gpu := entity.GPU{
		HostID:      host.ID,
		GPUIndex:    0,
		Name:        "RTX 4080",
		Brand:       "NVIDIA",
		MemoryTotal: 17179869184,
	}

	body, _ := json.Marshal(gpu)
	req := httptest.NewRequest("POST", "/gpus", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("创建GPU失败, 状态码: %d", w.Code)
	}
	t.Log("创建GPU成功")

	// 清理
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if data, ok := resp["data"].(map[string]any); ok {
		if id, ok := data["id"].(float64); ok {
			delReq := httptest.NewRequest("DELETE", fmt.Sprintf("/gpus/%d", int(id)), nil)
			delW := httptest.NewRecorder()
			r.ServeHTTP(delW, delReq)
		}
	}
}
