package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitJWT_Success 测试JWT初始化成功
func TestInitJWT_Success(t *testing.T) {
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	expireHours := 24

	err := InitJWT(secret, expireHours)
	assert.NoError(t, err)
	assert.NotNil(t, jwtSecret)
	assert.Equal(t, time.Duration(expireHours)*time.Hour, jwtExpireTime)
}

// TestInitJWT_EmptySecret 测试空密钥
func TestInitJWT_EmptySecret(t *testing.T) {
	err := InitJWT("", 24)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestInitJWT_ShortSecret 测试密钥长度不足
func TestInitJWT_ShortSecret(t *testing.T) {
	err := InitJWT("short", 24)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 32 characters")
}

// TestInitJWT_InvalidExpireTime 测试无效的过期时间
func TestInitJWT_InvalidExpireTime(t *testing.T) {
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"

	// 测试零值
	err := InitJWT(secret, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")

	// 测试负值
	err = InitJWT(secret, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
}

// TestGenerateToken_Success 测试Token生成成功
func TestGenerateToken_Success(t *testing.T) {
	// 初始化JWT
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	err := InitJWT(secret, 24)
	require.NoError(t, err)

	// 生成Token
	token, err := GenerateToken(1, "testuser", "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// TestGenerateToken_NotInitialized 测试未初始化时生成Token
func TestGenerateToken_NotInitialized(t *testing.T) {
	// 重置JWT配置
	jwtSecret = nil
	jwtExpireTime = 0

	token, err := GenerateToken(1, "testuser", "admin")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "not initialized")
}

// TestParseToken_Success 测试Token解析成功
func TestParseToken_Success(t *testing.T) {
	// 初始化JWT
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	err := InitJWT(secret, 24)
	require.NoError(t, err)

	// 生成Token
	userID := uint(123)
	username := "testuser"
	role := "admin"
	token, err := GenerateToken(userID, username, role)
	require.NoError(t, err)

	// 解析Token
	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
}

// TestParseToken_Invalid 测试无效Token
func TestParseToken_Invalid(t *testing.T) {
	// 初始化JWT
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	err := InitJWT(secret, 24)
	require.NoError(t, err)

	// 测试完全无效的Token
	claims, err := ParseToken("invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, claims)

	// 测试空Token
	claims, err = ParseToken("")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestParseToken_Expired 测试过期Token
func TestParseToken_Expired(t *testing.T) {
	// 初始化JWT，设置很短的过期时间
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	jwtSecret = []byte(secret)
	jwtExpireTime = 1 * time.Millisecond

	// 生成Token
	claims := Claims{
		UserID:   1,
		Username: "testuser",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 已过期
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	// 解析过期Token
	parsedClaims, err := ParseToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.Contains(t, err.Error(), "expired")
}

// TestParseToken_WrongSecret 测试使用错误密钥签名的Token
func TestParseToken_WrongSecret(t *testing.T) {
	// 使用一个密钥生成Token
	secret1 := "this-is-the-first-secret-key-with-32-characters-or-more"
	jwtSecret = []byte(secret1)
	jwtExpireTime = 24 * time.Hour

	token, err := GenerateToken(1, "testuser", "admin")
	require.NoError(t, err)

	// 使用另一个密钥解析Token
	secret2 := "this-is-the-second-secret-key-with-32-characters-or-more"
	err = InitJWT(secret2, 24)
	require.NoError(t, err)

	claims, err := ParseToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestParseToken_NotInitialized 测试未初始化时解析Token
func TestParseToken_NotInitialized(t *testing.T) {
	// 重置JWT配置
	jwtSecret = nil
	jwtExpireTime = 0

	claims, err := ParseToken("some.token.string")
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "not initialized")
}

// TestParseToken_WrongSigningMethod 测试错误的签名算法
func TestParseToken_WrongSigningMethod(t *testing.T) {
	// 初始化JWT
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	err := InitJWT(secret, 24)
	require.NoError(t, err)

	// 使用不同的签名算法生成Token (使用RS256而不是HS256)
	// 注意：这里我们模拟一个使用不同算法的token
	claims := Claims{
		UserID:   1,
		Username: "testuser",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建一个使用None算法的token（这会被拒绝）
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	// 尝试解析
	parsedClaims, err := ParseToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}

// TestTokenExpiration 测试Token过期时间设置
func TestTokenExpiration(t *testing.T) {
	// 初始化JWT，设置1小时过期
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	err := InitJWT(secret, 1)
	require.NoError(t, err)

	// 生成Token
	token, err := GenerateToken(1, "testuser", "admin")
	require.NoError(t, err)

	// 解析Token并检查过期时间
	claims, err := ParseToken(token)
	require.NoError(t, err)

	// 验证过期时间大约是1小时后
	expectedExpiry := time.Now().Add(1 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time

	// 允许5秒的误差
	diff := actualExpiry.Sub(expectedExpiry)
	assert.True(t, diff < 5*time.Second && diff > -5*time.Second,
		"Expected expiry around %v, got %v", expectedExpiry, actualExpiry)
}

// TestClaims_AllFields 测试Claims的所有字段
func TestClaims_AllFields(t *testing.T) {
	// 初始化JWT
	secret := "this-is-a-valid-secret-key-with-32-characters-or-more"
	err := InitJWT(secret, 24)
	require.NoError(t, err)

	// 生成Token
	userID := uint(999)
	username := "admin_user"
	role := "super_admin"
	token, err := GenerateToken(userID, username, role)
	require.NoError(t, err)

	// 解析Token
	claims, err := ParseToken(token)
	require.NoError(t, err)

	// 验证所有字段
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.True(t, claims.IssuedAt.Before(claims.ExpiresAt.Time))
}
