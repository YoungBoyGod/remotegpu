package auth

import (
	"context"
	"errors"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"gorm.io/gorm"
)

type AuthService struct {
	customerDao *dao.CustomerDao
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		customerDao: dao.NewCustomerDao(db),
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, int64, error) {
	customer, err := s.customerDao.FindByUsername(ctx, username)
	if err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, customer.PasswordHash) {
		return "", "", 0, errors.New("invalid credentials")
	}

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
	
	return accessToken, refreshToken, 3600, nil // 1小时过期
}

func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*entity.Customer, error) {
	return s.customerDao.FindByID(ctx, userID)
}