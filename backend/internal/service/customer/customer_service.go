package customer

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"gorm.io/gorm"
)

type CustomerService struct {
	customerDao   *dao.CustomerDao
	allocationDao *dao.AllocationDao
	db            *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{
		customerDao:   dao.NewCustomerDao(db),
		allocationDao: dao.NewAllocationDao(db),
		db:            db,
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

func (s *CustomerService) GetCustomer(ctx context.Context, id uint) (*entity.Customer, error) {
	return s.customerDao.FindByID(ctx, id)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id uint, fields map[string]interface{}) error {
	return s.customerDao.UpdateFields(ctx, id, fields)
}

func (s *CustomerService) UpdateStatus(ctx context.Context, id uint, status string) error {
	return s.customerDao.UpdateStatus(ctx, id, status)
}

// CountActive 统计活跃客户数量
// @modified 2026-02-04
func (s *CustomerService) CountActive(ctx context.Context) (int64, error) {
	return s.customerDao.CountActive(ctx)
}

// GetCustomerDetail 获取客户详情（包含机器分配信息）
func (s *CustomerService) GetCustomerDetail(ctx context.Context, id uint) (map[string]interface{}, error) {
	customer, err := s.customerDao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allocations, err := s.allocationDao.FindAllActiveByCustomerID(ctx, id)
	if err != nil {
		return nil, err
	}

	allocList := make([]map[string]interface{}, 0, len(allocations))
	for _, alloc := range allocations {
		item := map[string]interface{}{
			"allocation_id": alloc.ID,
			"machine_id":    alloc.HostID,
			"allocated_at":  alloc.StartTime,
			"end_time":      alloc.EndTime,
		}
		if alloc.Host.ID != "" {
			item["machine_name"] = alloc.Host.Name
			item["ssh_host"] = alloc.Host.SSHHost
			item["ssh_port"] = alloc.Host.SSHPort
			item["jupyter_url"] = alloc.Host.JupyterURL
			item["vnc_url"] = alloc.Host.VNCURL
			item["status"] = alloc.Host.Status
		}
		allocList = append(allocList, item)
	}

	return map[string]interface{}{
		"customer":    customer,
		"allocations": allocList,
	}, nil
}