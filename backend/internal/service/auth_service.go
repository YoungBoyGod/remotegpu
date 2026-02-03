package service

import (
	"context"
	"errors"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"golang.org/x/crypto/bcrypt"
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

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password)); err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	// Generate Tokens
	accessToken, err := auth.GenerateToken(customer.ID, customer.Username, customer.Role)
	if err != nil {
		return "", "", 0, err
	}
	
	// In a real scenario, Refresh Token should be stored in DB/Redis
	refreshToken := "dummy_refresh_token_" + accessToken[len(accessToken)-10:] 
	
	return accessToken, refreshToken, 3600, nil // 1 hour expiry
}

func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*entity.Customer, error) {
	return s.customerDao.FindByID(ctx, userID)
}
