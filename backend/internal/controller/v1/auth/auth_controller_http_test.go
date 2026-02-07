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

// httpTestEnv HTTP 测试环境
type httpTestEnv struct {
	db     *gorm.DB
	server *httptest.Server
}

// setupHTTPTestEnv 初始化 HTTP 测试环境
func setupHTTPTestEnv(t *testing.T) *httpTestEnv {
	gin.SetMode(gin.TestMode)

	// 初始化 JWT
	err := pkgAuth.InitJWT("test-secret-key-must-be-at-least-32-characters", 1)
	require.NoError(t, err)

	// 使用 SQLite 内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 手动创建表
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

	// 启动真实 HTTP 服务器
	server := httptest.NewServer(router)

	return &httpTestEnv{
		db:     db,
		server: server,
	}
}

// cleanup 清理测试环境
func (env *httpTestEnv) cleanup() {
	env.server.Close()
}

// ==================== HTTP 登录测试 ====================

func TestHTTPLogin_Success(t *testing.T) {
	env := setupHTTPTestEnv(t)
	defer env.cleanup()

	reqBody := apiV1.LoginRequest{
		Username: "testuser",
		Password: "Test123456",
	}
	body, _ := json.Marshal(reqBody)

	// 发送真实 HTTP 请求
	resp, err := http.Post(
		env.server.URL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// 验证 HTTP 状态码
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result testResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 0, result.Code)
	assert.Equal(t, "success", result.Msg)
}

func TestHTTPLogin_InvalidPassword(t *testing.T) {
	env := setupHTTPTestEnv(t)
	defer env.cleanup()

	reqBody := apiV1.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		env.server.URL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result testResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, errors.ErrorPasswordIncorrect, result.Code)
}

// ==================== HTTP Admin 登录测试 ====================

func TestHTTPAdminLogin_Success(t *testing.T) {
	env := setupHTTPTestEnv(t)
	defer env.cleanup()

	reqBody := apiV1.LoginRequest{
		Username: "admin",
		Password: "luoyang@123",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		env.server.URL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result testResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 0, result.Code)

	var loginResp apiV1.LoginResponse
	err = json.Unmarshal(result.Data, &loginResp)
	require.NoError(t, err)

	// 验证 token 中的角色
	claims, err := pkgAuth.ParseToken(loginResp.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "admin", claims.Role)
}

// ==================== HTTP GetProfile 测试 ====================

func TestHTTPGetProfile_Success(t *testing.T) {
	env := setupHTTPTestEnv(t)
	defer env.cleanup()

	token, _ := pkgAuth.GenerateToken(1, "testuser", "customer_owner")

	req, _ := http.NewRequest(http.MethodGet, env.server.URL+"/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result testResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 0, result.Code)
}

func TestHTTPGetProfile_NoToken(t *testing.T) {
	env := setupHTTPTestEnv(t)
	defer env.cleanup()

	resp, err := http.Get(env.server.URL + "/api/v1/auth/profile")
	require.NoError(t, err)
	defer resp.Body.Close()

	var result testResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 401, result.Code)
}

// ==================== HTTP 普通用户登录测试 ====================

func TestHTTPNormalUserLogin_Success(t *testing.T) {
	env := setupHTTPTestEnv(t)
	defer env.cleanup()

	reqBody := apiV1.LoginRequest{
		Username: "user",
		Password: "user@123",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		env.server.URL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result testResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 0, result.Code)

	var loginResp apiV1.LoginResponse
	err = json.Unmarshal(result.Data, &loginResp)
	require.NoError(t, err)

	claims, err := pkgAuth.ParseToken(loginResp.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "customer_owner", claims.Role)
}
