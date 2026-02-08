package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceAuth "github.com/YoungBoyGod/remotegpu/internal/service/auth"
	pkgAuth "github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testResponse 测试响应结构
type testResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// testEnv 测试环境
type testEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

// setupTestEnv 初始化测试环境
func setupTestEnv(t *testing.T) *testEnv {
	gin.SetMode(gin.TestMode)

	// 初始化 JWT
	err := pkgAuth.InitJWT("test-secret-key-must-be-at-least-32-characters", 1)
	require.NoError(t, err)

	// 使用 SQLite 内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 手动创建表（避免 PostgreSQL 特有函数）
	err = db.Exec(`CREATE TABLE customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		uuid TEXT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		display_name TEXT,
		full_name TEXT,
		company_code TEXT,
		company TEXT,
		phone TEXT,
		avatar_url TEXT,
		role TEXT DEFAULT 'customer_owner',
		user_type TEXT DEFAULT 'external',
		account_type TEXT DEFAULT 'individual',
		status TEXT DEFAULT 'active',
		email_verified INTEGER DEFAULT 0,
		phone_verified INTEGER DEFAULT 0,
		must_change_password INTEGER DEFAULT 0,
		quota_gpu INTEGER DEFAULT 0,
		quota_storage INTEGER DEFAULT 0,
		balance REAL DEFAULT 0,
		currency TEXT DEFAULT 'CNY',
		credit_limit REAL DEFAULT 0,
		billing_plan_id INTEGER,
		last_login_at DATETIME
	)`).Error
	require.NoError(t, err)

	// 创建测试用户
	password, _ := pkgAuth.HashPassword("Test123456")
	testUser := &entity.Customer{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: password,
		Role:         "customer_owner",
		Status:       "active",
	}
	db.Create(testUser)

	// 创建 admin 用户
	adminPassword, _ := pkgAuth.HashPassword("luoyang@123")
	adminUser := &entity.Customer{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: adminPassword,
		Role:         "admin",
		Status:       "active",
	}
	db.Create(adminUser)

	// 创建普通用户
	normalPassword, _ := pkgAuth.HashPassword("user@123")
	normalUser := &entity.Customer{
		Username:     "user",
		Email:        "user@example.com",
		PasswordHash: normalPassword,
		Role:         "customer_owner",
		Status:       "active",
	}
	db.Create(normalUser)

	// 初始化服务和控制器
	authService := serviceAuth.NewAuthService(db, nil)
	controller := NewAuthController(authService)

	// 设置路由
	router := gin.New()
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", controller.Login)
		authGroup.POST("/refresh", controller.Refresh)
		authGroup.POST("/logout", controller.Logout)
		authGroup.GET("/profile", testAuthMiddleware(), controller.GetProfile)
		authGroup.POST("/password/change", testAuthMiddleware(), controller.ChangePassword)
	}

	return &testEnv{
		db:     db,
		router: router,
	}
}

// testAuthMiddleware 测试用认证中间件
func testAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请提供认证令牌"})
			c.Abort()
			return
		}

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr := authHeader[7:]
			claims, err := pkgAuth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "无效的认证令牌"})
				c.Abort()
				return
			}
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "认证令牌格式错误"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ==================== 登录测试 ====================

func TestLogin_Success(t *testing.T) {
	env := setupTestEnv(t)

	reqBody := apiV1.LoginRequest{
		Username: "testuser",
		Password: "Test123456",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Msg)

	var loginResp apiV1.LoginResponse
	err = json.Unmarshal(resp.Data, &loginResp)
	require.NoError(t, err)

	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
	assert.Equal(t, int64(3600), loginResp.ExpiresIn)
}

func TestLogin_InvalidUsername(t *testing.T) {
	env := setupTestEnv(t)

	reqBody := apiV1.LoginRequest{
		Username: "nonexistent",
		Password: "Test123456",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, errors.ErrorPasswordIncorrect, resp.Code)
}

func TestLogin_InvalidPassword(t *testing.T) {
	env := setupTestEnv(t)

	reqBody := apiV1.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, errors.ErrorPasswordIncorrect, resp.Code)
}

func TestLogin_MissingParams(t *testing.T) {
	env := setupTestEnv(t)

	reqBody := map[string]string{"username": "testuser"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.Code)
}

// ==================== GetProfile 测试 ====================

func TestGetProfile_Success(t *testing.T) {
	env := setupTestEnv(t)

	token, _ := pkgAuth.GenerateToken(1, "testuser", "customer_owner")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
}

func TestGetProfile_NoToken(t *testing.T) {
	env := setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.Code)
}

func TestGetProfile_InvalidToken(t *testing.T) {
	env := setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.Code)
}

// ==================== ChangePassword 测试 ====================

func TestChangePassword_Success(t *testing.T) {
	env := setupTestEnv(t)

	_ = env.db.Model(&entity.Customer{}).Where("id = ?", 1).Update("must_change_password", true).Error

	token, _ := pkgAuth.GenerateToken(1, "testuser", "customer_owner")
	reqBody := apiV1.ChangePasswordRequest{
		OldPassword: "Test123456",
		NewPassword: "NewPass123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password/change", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	var customer entity.Customer
	err = env.db.First(&customer, 1).Error
	require.NoError(t, err)
	assert.False(t, customer.MustChangePassword)
	assert.True(t, pkgAuth.CheckPasswordHash("NewPass123", customer.PasswordHash))
}

func TestChangePassword_InvalidOldPassword(t *testing.T) {
	env := setupTestEnv(t)

	token, _ := pkgAuth.GenerateToken(1, "testuser", "customer_owner")
	reqBody := apiV1.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "NewPass123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password/change", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, errors.ErrorPasswordIncorrect, resp.Code)
}

// ==================== Logout 测试 ====================

func TestLogout_Success(t *testing.T) {
	env := setupTestEnv(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
}

// ==================== Admin 登录测试 ====================

func TestAdminLogin_Success(t *testing.T) {
	env := setupTestEnv(t)

	reqBody := apiV1.LoginRequest{
		Username: "admin",
		Password: "luoyang@123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Msg)

	var loginResp apiV1.LoginResponse
	err = json.Unmarshal(resp.Data, &loginResp)
	require.NoError(t, err)

	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)

	// 验证 token 中的角色是 admin
	claims, err := pkgAuth.ParseToken(loginResp.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "admin", claims.Role)
}

// ==================== 普通用户登录测试 ====================

func TestNormalUserLogin_Success(t *testing.T) {
	env := setupTestEnv(t)

	reqBody := apiV1.LoginRequest{
		Username: "user",
		Password: "user@123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)

	var resp testResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Msg)

	var loginResp apiV1.LoginResponse
	err = json.Unmarshal(resp.Data, &loginResp)
	require.NoError(t, err)

	assert.NotEmpty(t, loginResp.AccessToken)

	// 验证 token 中的角色
	claims, err := pkgAuth.ParseToken(loginResp.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "customer_owner", claims.Role)
}
