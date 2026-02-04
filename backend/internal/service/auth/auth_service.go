package auth

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"gorm.io/gorm"
)

const (
	// token 黑名单 key 前缀
	tokenBlacklistPrefix = "auth:token:blacklist:"
	// token 黑名单过期时间
	tokenBlacklistTTL = time.Hour
	// 刷新 token key 前缀
	refreshTokenPrefix = "auth:refresh_token:"
	// 刷新 token 过期时间 (7天)
	refreshTokenTTL = 7 * 24 * time.Hour
)

type AuthService struct {
	customerDao *dao.CustomerDao
	db          *gorm.DB
	cache       cache.Cache
}

func NewAuthService(db *gorm.DB, c cache.Cache) *AuthService {
	return &AuthService{
		customerDao: dao.NewCustomerDao(db),
		db:          db,
		cache:       c,
	}
}

// storeRefreshToken 存储刷新 token
func (s *AuthService) storeRefreshToken(ctx context.Context, refreshToken string, userID uint) error {
	if s.cache == nil {
		return nil
	}
	key := refreshTokenPrefix + refreshToken
	// 存储 UserID
	return s.cache.Set(ctx, key, userID, refreshTokenTTL)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, int64, error) {
	customer, err := s.customerDao.FindByUsername(ctx, username)
	if err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, customer.PasswordHash) {
		return "", "", 0, errors.New("invalid credentials")
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(customer).Update("last_login_at", now)

	// 生成 Token
	accessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, err
	}

	// 生成刷新 Token
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, err
	}

	// 存储刷新 Token
	if err := s.storeRefreshToken(ctx, refreshToken, customer.ID); err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, 3600, nil // 1小时过期
}

func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*entity.Customer, error) {
	return s.customerDao.FindByID(ctx, userID)
}

// AdminLogin Admin 专用登录，验证角色
func (s *AuthService) AdminLogin(ctx context.Context, username, password string) (string, string, int64, error) {
	customer, err := s.customerDao.FindByUsername(ctx, username)
	if err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	// 验证密码
	if !auth.CheckPasswordHash(password, customer.PasswordHash) {
		return "", "", 0, errors.New("invalid credentials")
	}

	// 验证是否是 admin 角色
	if customer.Role != "admin" {
		return "", "", 0, errors.New("permission denied: admin role required")
	}

	// 验证账号状态
	if customer.Status != "active" {
		return "", "", 0, errors.New("account is disabled")
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(customer).Update("last_login_at", now)

	// 生成 Token
	accessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, err
	}

	// 生成刷新 Token
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, err
	}

	// 存储刷新 Token
	if err := s.storeRefreshToken(ctx, refreshToken, customer.ID); err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, 3600, nil
}

// RefreshToken 使用刷新令牌获取新的访问令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, int64, error) {
	if s.cache == nil {
		return "", "", 0, errors.New("cache service not available")
	}

	// 1. Check if token exists
	key := refreshTokenPrefix + refreshToken
	userIDStr, err := s.cache.Get(ctx, key)
	if err != nil || userIDStr == "" {
		return "", "", 0, errors.New("invalid or expired refresh token")
	}

	// 2. Parse UserID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return "", "", 0, errors.New("invalid token data")
	}

	// 3. Get User (check status)
	customer, err := s.customerDao.FindByID(ctx, uint(userID))
	if err != nil {
		return "", "", 0, errors.New("user not found")
	}
	if customer.Status != "active" {
		return "", "", 0, errors.New("account is disabled")
	}

	// 4. Generate new tokens
	newAccessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, err
	}
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, err
	}

	// 5. Rotate tokens (delete old, save new)
	_ = s.cache.Delete(ctx, key)

	if err := s.storeRefreshToken(ctx, newRefreshToken, customer.ID); err != nil {
		return "", "", 0, err
	}

	return newAccessToken, newRefreshToken, 3600, nil
}

// Logout 登出，将 token 加入 Redis 黑名单
func (s *AuthService) Logout(ctx context.Context, token string) error {
	if s.cache == nil {
		return nil
	}
	key := tokenBlacklistPrefix + token
	return s.cache.Set(ctx, key, "1", tokenBlacklistTTL)
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func (s *AuthService) IsTokenBlacklisted(ctx context.Context, token string) bool {
	if s.cache == nil {
		return false
	}
	key := tokenBlacklistPrefix + token
	count, err := s.cache.Exists(ctx, key)
	if err != nil {
		return false
	}
	return count > 0
}