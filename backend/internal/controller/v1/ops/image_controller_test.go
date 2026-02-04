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
	"github.com/YoungBoyGod/remotegpu/internal/service/image"
)

func setupImageTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.Exec(`CREATE TABLE images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(256) NOT NULL UNIQUE,
		display_name VARCHAR(256),
		description TEXT,
		category VARCHAR(64),
		framework VARCHAR(64),
		cuda_version VARCHAR(32),
		registry_url VARCHAR(512),
		is_official INTEGER DEFAULT 0,
		customer_id INTEGER,
		status VARCHAR(20) DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	return db
}

func setupImageRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := image.NewImageService(db)
	ctrl := NewImageController(svc)

	api := r.Group("/api/v1/admin")
	api.GET("/images", ctrl.List)

	return r
}

func TestImage_List_Empty(t *testing.T) {
	db := setupImageTestDB(t)
	r := setupImageRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/admin/images", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
}

func TestImage_List_WithData(t *testing.T) {
	db := setupImageTestDB(t)
	r := setupImageRouter(db)

	db.Create(&entity.Image{
		Name:        "pytorch/pytorch:2.0-cuda11.8",
		DisplayName: "PyTorch 2.0",
		Category:    "deep-learning",
		Framework:   "pytorch",
	})

	req, _ := http.NewRequest("GET", "/api/v1/admin/images", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, float64(1), data["total"])
}

func TestImage_List_WithFilter(t *testing.T) {
	db := setupImageTestDB(t)
	r := setupImageRouter(db)

	db.Create(&entity.Image{Name: "pytorch:2.0", Framework: "pytorch"})
	db.Create(&entity.Image{Name: "tensorflow:2.0", Framework: "tensorflow"})

	req, _ := http.NewRequest("GET", "/api/v1/admin/images?framework=pytorch", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, float64(1), data["total"])
}
