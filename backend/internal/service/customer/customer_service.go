package customer

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
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
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	customer.PasswordHash = hash
	return s.customerDao.Create(ctx, customer)
}

func (s *CustomerService) ListCustomers(ctx context.Context, page, pageSize int) ([]entity.Customer, int64, error) {
	return s.customerDao.List(ctx, page, pageSize)
}

func (s *CustomerService) UpdateStatus(ctx context.Context, id uint, status string) error {
	return s.customerDao.UpdateStatus(ctx, id, status)
}

// CountActive 统计活跃客户数量
// @modified 2026-02-04
func (s *CustomerService) CountActive(ctx context.Context) (int64, error) {
	return s.customerDao.CountActive(ctx)
}