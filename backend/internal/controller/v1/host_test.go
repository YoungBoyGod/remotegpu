package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/gin-gonic/gin"
)

func setupTestDB(t *testing.T) {
	cfg := database.Config{
		Host:     "192.168.10.210",
		Port:     5432,
		User:     "remotegpu_user",
		Password: "remotegpu_password",
		DBName:   "remotegpu",
	}
	if err := database.InitDB(cfg); err != nil {
		t.Skipf("跳过测试，无法连接数据库: %v", err)
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestHostController_Create(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()
	ctrl := NewHostController()

	r.POST("/hosts", ctrl.Create)

	host := entity.Host{
		Name:           "Test Controller Host",
		IPAddress:      "192.168.1.250",
		OSType:         "linux",
		DeploymentMode: "traditional",
		TotalCPU:       8,
		TotalMemory:    17179869184,
	}

	body, _ := json.Marshal(host)
	req := httptest.NewRequest("POST", "/hosts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("创建主机失败, 状态码: %d, 响应: %s", w.Code, w.Body.String())
	}
	t.Log("创建主机成功")

	// 清理
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(string); ok {
			r.DELETE("/hosts/:id", ctrl.Delete)
			delReq := httptest.NewRequest("DELETE", "/hosts/"+id, nil)
			delW := httptest.NewRecorder()
			r.ServeHTTP(delW, delReq)
		}
	}
}
