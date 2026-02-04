package allocation

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"github.com/YoungBoyGod/remotegpu/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AllocationService struct {
	db            *gorm.DB
	allocationDao *dao.AllocationDao
	machineDao    *dao.MachineDao
	auditService  *audit.AuditService
}

func NewAllocationService(db *gorm.DB, auditSvc *audit.AuditService) *AllocationService {
	return &AllocationService{
		db:            db,
		allocationDao: dao.NewAllocationDao(db),
		machineDao:    dao.NewMachineDao(db),
		auditService:  auditSvc,
	}
}

func (s *AllocationService) AllocateMachine(ctx context.Context, customerID uint, hostID string, durationMonths int, remark string) (*entity.Allocation, error) {
	if durationMonths < 1 {
		return nil, errors.New(errors.ErrorInvalidParams, "lease duration must be at least 1 month")
	}

	// 使用事务确保原子性
	var allocation *entity.Allocation
	
	err := s.db.Transaction(func(tx *gorm.DB) error {
		machineDao := dao.NewMachineDao(tx)
		allocationDao := dao.NewAllocationDao(tx)

		// 1. 检查机器状态 (如果可能应加锁，此处仅检查状态)
		host, err := machineDao.FindByID(ctx, hostID)
		if err != nil {
			return errors.Wrap(errors.ErrorHostNotFound, err)
		}
		if host.Status != "idle" && host.Status != "online" { // 假设 'online' 表示空闲/可用
			return errors.New(errors.ErrorMachineNotAvailable, "machine is not available for allocation")
		}

		// 2. 更新机器状态
		if err := machineDao.UpdateStatus(ctx, hostID, "allocated"); err != nil {
			return errors.Wrap(errors.ErrorDatabase, err)
		}

		// 3. 创建分配记录
		startTime := time.Now()
		endTime := startTime.AddDate(0, durationMonths, 0)
		allocation = &entity.Allocation{
			ID:         "alloc-" + uuid.New().String(),
			CustomerID: customerID,
			HostID:     hostID,
			StartTime:  startTime,
			EndTime:    endTime,
			Status:     "active",
			Remark:     remark,
		}
		if err := allocationDao.Create(ctx, allocation); err != nil {
			return errors.Wrap(errors.ErrorDatabase, err)
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	// TODO: 异步触发 Agent 重置 SSH
	// go s.agentService.ResetMachine(hostID)

	return allocation, nil
}

func (s *AllocationService) ReclaimMachine(ctx context.Context, hostID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		machineDao := dao.NewMachineDao(tx)
		allocationDao := dao.NewAllocationDao(tx)

		// 1. 查找活跃分配
		alloc, err := allocationDao.FindActiveByHostID(ctx, hostID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(errors.ErrorAllocationNotFound, "no active allocation found for this host")
			}
			return errors.Wrap(errors.ErrorDatabase, err)
		}

		// 2. 更新分配状态
		now := time.Now()
		alloc.Status = "reclaimed"
		alloc.ActualEndTime = &now
		// 在实际 dao 中，应该有特定的更新方法或使用通用更新
		if err := tx.Save(alloc).Error; err != nil {
			return errors.Wrap(errors.ErrorDatabase, err)
		}

		// 3. 更新机器状态 (进入维护/清理状态)
		if err := machineDao.UpdateStatus(ctx, hostID, "maintenance"); err != nil {
			return errors.Wrap(errors.ErrorDatabase, err)
		}

		return nil
	})
	
	if err != nil {
		return err
	}

	// 4. Log Audit (Fire and Forget or handle error)
	_ = s.auditService.CreateLog(
		ctx,
		nil, // System action, no customer ID
		"system", "127.0.0.1", "POST", "/reclaim",
		"reclaim_machine", "machine", hostID,
		map[string]interface{}{"reason": "admin_request"},
		200,
	)

	// TODO: 触发清理流程
	return nil
}

func (s *AllocationService) GetRecent(ctx context.Context) ([]entity.Allocation, error) {
	return s.allocationDao.FindRecent(ctx, 5)
}

// ListByCustomerID 获取指定客户的活跃分配列表
// @author Claude
// @description 根据客户ID查询其所有活跃的机器分配记录，用于客户端"我的机器"列表
// @param customerID 客户ID（从JWT中获取）
// @param page 页码
// @param pageSize 每页数量
// @return 分配列表、总数、错误
// @modified 2026-02-04
func (s *AllocationService) ListByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Allocation, int64, error) {
	return s.allocationDao.FindActiveByCustomerID(ctx, customerID, page, pageSize)
}