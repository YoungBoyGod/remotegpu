package customer

import (
	"context"
	"errors"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"gorm.io/gorm"
)

var (
	ErrQuotaExceeded = errors.New("已超出配额限制")
)

type CustomerService struct {
	customerDao   *dao.CustomerDao
	allocationDao *dao.AllocationDao
	datasetDao    *dao.DatasetDao
	auditDao      *dao.AuditDao
	db            *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{
		customerDao:   dao.NewCustomerDao(db),
		allocationDao: dao.NewAllocationDao(db),
		datasetDao:    dao.NewDatasetDao(db),
		auditDao:      dao.NewAuditDao(db),
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

	// 获取资源使用统计
	usage, _ := s.GetResourceUsage(ctx, id)

	return map[string]interface{}{
		"customer":       customer,
		"allocations":    allocList,
		"resource_usage": usage,
	}, nil
}

// ResourceUsage 客户资源使用统计
type ResourceUsage struct {
	AllocatedMachines int64 `json:"allocated_machines"` // 已分配机器数
	AllocatedGPUs     int64 `json:"allocated_gpus"`     // 已分配 GPU 数
	StorageUsedMB     int64 `json:"storage_used_mb"`    // 存储用量（MB）
	DatasetCount      int64 `json:"dataset_count"`      // 数据集数量
}

// GetResourceUsage 获取客户资源使用统计
func (s *CustomerService) GetResourceUsage(ctx context.Context, customerID uint) (*ResourceUsage, error) {
	machineCount, err := s.allocationDao.CountActiveByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	gpuCount, err := s.allocationDao.CountGPUsByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	storageBytes, err := s.datasetDao.SumStorageByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	_, datasetCount, err := s.datasetDao.ListByCustomerID(ctx, customerID, 1, 1)
	if err != nil {
		return nil, err
	}

	return &ResourceUsage{
		AllocatedMachines: machineCount,
		AllocatedGPUs:     gpuCount,
		StorageUsedMB:     storageBytes / (1024 * 1024),
		DatasetCount:      datasetCount,
	}, nil
}

// UpdateQuota 更新客户配额
func (s *CustomerService) UpdateQuota(ctx context.Context, customerID uint, quotaGPU int, quotaStorage int64) error {
	_, err := s.customerDao.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	return s.customerDao.UpdateQuota(ctx, customerID, quotaGPU, quotaStorage)
}

// CheckGPUQuota 检查客户 GPU 配额是否允许新增分配
func (s *CustomerService) CheckGPUQuota(ctx context.Context, customerID uint, additionalGPUs int64) error {
	customer, err := s.customerDao.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	// quota_gpu 为 0 表示不限制
	if customer.QuotaGPU == 0 {
		return nil
	}

	currentGPUs, err := s.allocationDao.CountGPUsByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	if currentGPUs+additionalGPUs > int64(customer.QuotaGPU) {
		return ErrQuotaExceeded
	}
	return nil
}

// CheckStorageQuota 检查客户存储配额是否允许新增存储
func (s *CustomerService) CheckStorageQuota(ctx context.Context, customerID uint, additionalMB int64) error {
	customer, err := s.customerDao.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	// quota_storage 为 0 表示不限制
	if customer.QuotaStorage == 0 {
		return nil
	}

	currentBytes, err := s.datasetDao.SumStorageByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	currentMB := currentBytes / (1024 * 1024)
	if currentMB+additionalMB > customer.QuotaStorage {
		return ErrQuotaExceeded
	}
	return nil
}