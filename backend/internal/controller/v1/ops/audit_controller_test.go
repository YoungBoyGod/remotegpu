package ops

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
)

func setupAuditTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.Exec(`CREATE TABLE audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
		username VARCHAR(128),
		ip_address VARCHAR(64),
		method VARCHAR(10),
		path VARCHAR(512),
		action VARCHAR(128) NOT NULL,
		resource_type VARCHAR(64),
		resource_id VARCHAR(128),
		detail TEXT,
		status_code INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	return db
}

func setupAuditRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := audit.NewAuditService(db)
	ctrl := NewAuditController(svc)

	api := r.Group("/api/v1/admin")
	api.GET("/audit/logs", ctrl.List)

	return r
}

func TestAudit_List_Empty(t *testing.T) {
	db := setupAuditTestDB(t)
	r := setupAuditRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
}

func TestAudit_List_WithData(t *testing.T) {
	db := setupAuditTestDB(t)
	r := setupAuditRouter(db)

	// 插入测试数据
	db.Create(&entity.AuditLog{
		Username:     "admin",
		Action:       "create",
		ResourceType: "machine",
		IPAddress:    "127.0.0.1",
	})

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, float64(1), data["total"])
}

func TestAudit_List_WithFilter(t *testing.T) {
	db := setupAuditTestDB(t)
	r := setupAuditRouter(db)

	// 插入测试数据
	db.Create(&entity.AuditLog{Username: "admin", Action: "create", ResourceType: "machine"})
	db.Create(&entity.AuditLog{Username: "user1", Action: "delete", ResourceType: "ssh_key"})

	req, _ := http.NewRequest("GET", "/api/v1/admin/audit/logs?action=create", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, float64(1), data["total"])
}
