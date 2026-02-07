package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/errors"
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
	// 密码重置 token key 前缀
	passwordResetPrefix = "auth:password_reset:"
	// 密码重置 token 过期时间 (30分钟)
	passwordResetTTL = 30 * time.Minute
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

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, int64, bool, error) {
	customer, err := s.customerDao.FindByUsername(ctx, username)
	if err != nil {
		return "", "", 0, false, errors.New(errors.ErrorPasswordIncorrect, "invalid credentials")
	}

	if !auth.CheckPasswordHash(password, customer.PasswordHash) {
		return "", "", 0, false, errors.New(errors.ErrorPasswordIncorrect, "invalid credentials")
	}

	// 验证账号状态
	if customer.Status != "active" {
		return "", "", 0, false, errors.New(errors.ErrorUserDisabled, "")
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(customer).Update("last_login_at", now)

	// 生成 Token
	accessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, false, err
	}

	// 生成刷新 Token
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, false, err
	}

	// 存储刷新 Token
	if err := s.storeRefreshToken(ctx, refreshToken, customer.ID); err != nil {
		return "", "", 0, false, err
	}

	return accessToken, refreshToken, 3600, customer.MustChangePassword, nil // 1小时过期
}

func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*entity.Customer, error) {
	return s.customerDao.FindByID(ctx, userID)
}

// AdminLogin Admin 专用登录，验证角色
func (s *AuthService) AdminLogin(ctx context.Context, username, password string) (string, string, int64, bool, error) {
	customer, err := s.customerDao.FindByUsername(ctx, username)
	if err != nil {
		return "", "", 0, false, errors.New(errors.ErrorPasswordIncorrect, "invalid credentials")
	}

	// 验证密码
	if !auth.CheckPasswordHash(password, customer.PasswordHash) {
		return "", "", 0, false, errors.New(errors.ErrorPasswordIncorrect, "invalid credentials")
	}

	// 验证是否是 admin 角色
	if customer.Role != "admin" {
		return "", "", 0, false, errors.New(errors.ErrorForbidden, "permission denied: admin role required")
	}

	// 验证账号状态
	if customer.Status != "active" {
		return "", "", 0, false, errors.New(errors.ErrorUserDisabled, "account is disabled")
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(customer).Update("last_login_at", now)

	// 生成 Token
	accessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, false, err
	}

	// 生成刷新 Token
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, false, err
	}

	// 存储刷新 Token
	if err := s.storeRefreshToken(ctx, refreshToken, customer.ID); err != nil {
		return "", "", 0, false, err
	}

	return accessToken, refreshToken, 3600, customer.MustChangePassword, nil
}

// RefreshToken 使用刷新令牌获取新的访问令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, int64, bool, error) {
	if s.cache == nil {
		return "", "", 0, false, errors.New(errors.ErrorServerError, "cache service not available")
	}

	// 1. Check if token exists
	key := refreshTokenPrefix + refreshToken
	userIDStr, err := s.cache.Get(ctx, key)
	if err != nil || userIDStr == "" {
		return "", "", 0, false, errors.New(errors.ErrorTokenInvalid, "invalid or expired refresh token")
	}

	// 2. Parse UserID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return "", "", 0, false, errors.New(errors.ErrorTokenInvalid, "invalid token data")
	}

	// 3. Get User (check status)
	customer, err := s.customerDao.FindByID(ctx, uint(userID))
	if err != nil {
		return "", "", 0, false, errors.New(errors.ErrorUserNotFound, "user not found")
	}
	if customer.Status != "active" {
		return "", "", 0, false, errors.New(errors.ErrorUserDisabled, "account is disabled")
	}

	// 4. Generate new tokens
	newAccessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, false, err
	}
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, false, err
	}

	// 5. Rotate tokens (delete old, save new)
	_ = s.cache.Delete(ctx, key)

	if err := s.storeRefreshToken(ctx, newRefreshToken, customer.ID); err != nil {
		return "", "", 0, false, err
	}

	return newAccessToken, newRefreshToken, 3600, customer.MustChangePassword, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	customer, err := s.customerDao.FindByID(ctx, userID)
	if err != nil {
		return errors.New(errors.ErrorUserNotFound, "user not found")
	}

	if !auth.CheckPasswordHash(oldPassword, customer.PasswordHash) {
		return errors.New(errors.ErrorPasswordIncorrect, "invalid credentials")
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return errors.WrapWithMessage(errors.ErrorServerError, "failed to hash password", err)
	}

	return s.db.WithContext(ctx).Model(&entity.Customer{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash":        newHash,
		"must_change_password": false,
	}).Error
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

// RequestPasswordReset 请求密码重置，生成重置 token
func (s *AuthService) RequestPasswordReset(ctx context.Context, username string) (string, error) {
	if s.cache == nil {
		return "", errors.New(errors.ErrorServerError, "cache service not available")
	}

	customer, err := s.customerDao.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New(errors.ErrorUserNotFound, "user not found")
	}
	if customer.Status != "active" {
		return "", errors.New(errors.ErrorUserDisabled, "account is disabled")
	}

	// 生成重置 token
	token, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", errors.WrapWithMessage(errors.ErrorServerError, "failed to generate reset token", err)
	}

	// 存储 token -> userID 映射
	key := passwordResetPrefix + token
	if err := s.cache.Set(ctx, key, customer.ID, passwordResetTTL); err != nil {
		return "", errors.WrapWithMessage(errors.ErrorServerError, "failed to store reset token", err)
	}

	return token, nil
}

// ConfirmPasswordReset 确认密码重置
func (s *AuthService) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if s.cache == nil {
		return errors.New(errors.ErrorServerError, "cache service not available")
	}

	key := passwordResetPrefix + token
	userIDStr, err := s.cache.Get(ctx, key)
	if err != nil || userIDStr == "" {
		return errors.New(errors.ErrorTokenInvalid, "invalid or expired reset token")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return errors.New(errors.ErrorTokenInvalid, "invalid token data")
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return errors.WrapWithMessage(errors.ErrorServerError, "failed to hash password", err)
	}

	err = s.db.WithContext(ctx).Model(&entity.Customer{}).
		Where("id = ?", uint(userID)).
		Updates(map[string]interface{}{
			"password_hash":        newHash,
			"must_change_password": false,
		}).Error
	if err != nil {
		return errors.WrapWithMessage(errors.ErrorServerError, "failed to update password", err)
	}

	// 删除已使用的 token
	_ = s.cache.Delete(ctx, key)

	return nil
}
