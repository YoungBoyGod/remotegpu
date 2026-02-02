package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager Token 管理器
type TokenManager struct {
	jwtSecret []byte
}

// NewTokenManager 创建 Token 管理器
func NewTokenManager(jwtSecret string) *TokenManager {
	return &TokenManager{
		jwtSecret: []byte(jwtSecret),
	}
}

// GenerateJWT 生成 JWT Token
func (m *TokenManager) GenerateJWT(claims jwt.MapClaims, expiresIn time.Duration) (*Token, error) {
	// 设置过期时间
	expiresAt := time.Now().Add(expiresIn)
	claims["exp"] = expiresAt.Unix()
	claims["iat"] = time.Now().Unix()

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名
	tokenString, err := token.SignedString(m.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("生成 JWT Token 失败: %w", err)
	}

	return &Token{
		Type:      TokenTypeJWT,
		Value:     tokenString,
		ExpiresAt: expiresAt,
		Metadata:  convertClaimsToMetadata(claims),
		CreatedAt: time.Now(),
	}, nil
}

// ValidateJWT 验证 JWT Token
func (m *TokenManager) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("验证 JWT Token 失败: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的 Token")
}

// GenerateAPIKey 生成 API Key
func (m *TokenManager) GenerateAPIKey(prefix string, length int) (*Token, error) {
	if length < 16 {
		return nil, fmt.Errorf("API Key 长度至少为 16")
	}

	// 生成随机字节
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("生成随机字节失败: %w", err)
	}

	// Base64 编码
	apiKey := base64.URLEncoding.EncodeToString(bytes)
	if prefix != "" {
		apiKey = fmt.Sprintf("%s_%s", prefix, apiKey)
	}

	return &Token{
		Type:      TokenTypeAPIKey,
		Value:     apiKey,
		ExpiresAt: time.Time{}, // API Key 不过期
		Metadata:  map[string]string{"prefix": prefix},
		CreatedAt: time.Now(),
	}, nil
}

// GenerateSessionToken 生成 Session Token
func (m *TokenManager) GenerateSessionToken(length int, expiresIn time.Duration) (*Token, error) {
	if length < 32 {
		return nil, fmt.Errorf("Session Token 长度至少为 32")
	}

	// 生成随机字节
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("生成随机字节失败: %w", err)
	}

	// Base64 编码
	sessionToken := base64.URLEncoding.EncodeToString(bytes)

	return &Token{
		Type:      TokenTypeSession,
		Value:     sessionToken,
		ExpiresAt: time.Now().Add(expiresIn),
		Metadata:  map[string]string{},
		CreatedAt: time.Now(),
	}, nil
}

// IsExpired 检查 Token 是否过期
func (m *TokenManager) IsExpired(token *Token) bool {
	if token.ExpiresAt.IsZero() {
		return false // 永不过期
	}
	return time.Now().After(token.ExpiresAt)
}

// convertClaimsToMetadata 转换 Claims 为 Metadata
func convertClaimsToMetadata(claims jwt.MapClaims) map[string]string {
	metadata := make(map[string]string)
	for key, value := range claims {
		if key == "exp" || key == "iat" {
			continue
		}
		metadata[key] = fmt.Sprintf("%v", value)
	}
	return metadata
}
