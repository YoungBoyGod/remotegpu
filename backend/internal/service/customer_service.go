package service

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerService struct {
	customerDao *dao.CustomerDao
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{
		customerDao: dao.NewCustomerDao(db),
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *entity.Customer, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	customer.PasswordHash = string(hash)
	return s.customerDao.Create(ctx, customer)
}

func (s *CustomerService) ListCustomers(ctx context.Context, page, pageSize int) ([]entity.Customer, int64, error) {
	return s.customerDao.List(ctx, page, pageSize)
}

func (s *CustomerService) UpdateStatus(ctx context.Context, id uint, status string) error {
	return s.customerDao.UpdateStatus(ctx, id, status)
}
