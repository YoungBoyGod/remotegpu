package customer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/sshkey"
)

func setupSSHKeyTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 创建表
	err = db.Exec(`CREATE TABLE ssh_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER NOT NULL,
		name VARCHAR(64) NOT NULL,
		fingerprint VARCHAR(128),
		public_key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	return db
}

func setupSSHKeyRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := sshkey.NewSSHKeyService(db)
	ctrl := NewSSHKeyController(svc)

	// 模拟认证中间件
	authMiddleware := func(userID uint) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		}
	}

	api := r.Group("/api/v1/customer")
	api.Use(authMiddleware(1))
	{
		api.GET("/keys", ctrl.List)
		api.POST("/keys", ctrl.Create)
		api.DELETE("/keys/:id", ctrl.Delete)
	}

	return r
}

// 有效的 SSH 公钥（用于测试）
const testSSHPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl test@example.com"

func TestSSHKey_List_Empty(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/customer/keys", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	assert.Empty(t, list)
}

func TestSSHKey_Create_Success(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	body := map[string]string{
		"name":       "my-key",
		"public_key": testSSHPublicKey,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/customer/keys", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "my-key", data["name"])
	assert.NotEmpty(t, data["fingerprint"])
}

func TestSSHKey_Create_InvalidKey(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	body := map[string]string{
		"name":       "bad-key",
		"public_key": "invalid-key-format",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/customer/keys", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(400), resp["code"])
}

func TestSSHKey_Create_Duplicate(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	body := map[string]string{
		"name":       "my-key",
		"public_key": testSSHPublicKey,
	}
	jsonBody, _ := json.Marshal(body)

	// 第一次创建
	req1, _ := http.NewRequest("POST", "/api/v1/customer/keys", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 第二次创建（重复）
	jsonBody2, _ := json.Marshal(body)
	req2, _ := http.NewRequest("POST", "/api/v1/customer/keys", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(409), resp["code"])
}

func TestSSHKey_Delete_Success(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	// 先创建一个密钥
	key := &entity.SSHKey{
		CustomerID:  1,
		Name:        "to-delete",
		PublicKey:   testSSHPublicKey,
		Fingerprint: "SHA256:test",
	}
	db.Create(key)

	req, _ := http.NewRequest("DELETE", "/api/v1/customer/keys/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
}

func TestSSHKey_Delete_NotFound(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	req, _ := http.NewRequest("DELETE", "/api/v1/customer/keys/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(404), resp["code"])
}

func TestSSHKey_Delete_NotOwned(t *testing.T) {
	db := setupSSHKeyTestDB(t)
	r := setupSSHKeyRouter(db)

	// 创建属于其他用户的密钥
	key := &entity.SSHKey{
		CustomerID:  999, // 不同的用户
		Name:        "other-user-key",
		PublicKey:   testSSHPublicKey,
		Fingerprint: "SHA256:test",
	}
	db.Create(key)

	req, _ := http.NewRequest("DELETE", "/api/v1/customer/keys/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(403), resp["code"])
}
