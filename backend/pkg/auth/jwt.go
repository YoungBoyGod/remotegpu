package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	jwtExpireTime time.Duration
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// InitJWT 初始化 JWT 配置
// secret: JWT 密钥，必须至少 32 字符
// expireHours: Token 过期时间（小时）
func InitJWT(secret string, expireHours int) error {
	if secret == "" {
		return errors.New("JWT secret cannot be empty")
	}
	if len(secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters, got %d", len(secret))
	}
	if expireHours <= 0 {
		return fmt.Errorf("JWT expire time must be positive, got %d", expireHours)
	}

	jwtSecret = []byte(secret)
	jwtExpireTime = time.Duration(expireHours) * time.Hour
	return nil
}

// GenerateToken 生成 JWT token
func GenerateToken(userID uint, username string, role string) (string, error) {
	if jwtSecret == nil {
		return "", errors.New("JWT not initialized, call InitJWT first")
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	if jwtSecret == nil {
		return nil, errors.New("JWT not initialized, call InitJWT first")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
