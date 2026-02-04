package allocation

import (
	"context"
	stderrors "errors"
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

// ReclaimMachine 回收机器
// @author Claude
// @description 回收已分配的机器，更新分配状态和机器状态，并记录审计日志
// @reason 修复原实现中审计代码不可达的bug
// @modified 2026-02-04
func (s *AllocationService) ReclaimMachine(ctx context.Context, hostID string) error {
	var allocID string

	err := s.db.Transaction(func(tx *gorm.DB) error {
		machineDao := dao.NewMachineDao(tx)
		allocationDao := dao.NewAllocationDao(tx)

		// 1. 查找活跃分配
		alloc, err := allocationDao.FindActiveByHostID(ctx, hostID)
		if err != nil {
			if stderrors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(errors.ErrorAllocationNotFound, "no active allocation found for this host")
			}
			return errors.Wrap(errors.ErrorDatabase, err)
		}
		allocID = alloc.ID

		// 2. 更新分配状态
		now := time.Now()
		alloc.Status = "reclaimed"
		alloc.ActualEndTime = &now
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

	// 4. 记录审计日志
	_ = s.auditService.CreateLog(
		ctx,
		nil, // System action, no customer ID
		"system", "127.0.0.1", "POST", fmt.Sprintf("/admin/machines/%s/reclaim", hostID),
		"reclaim_machine", "machine", hostID,
		map[string]interface{}{"allocation_id": allocID, "reason": "admin_request"},
		200,
	)

	// TODO: 触发异步清理流程（重置SSH、清理用户数据等）

	return nil
}

func (s *AllocationService) GetRecent(ctx context.Context) ([]entity.Allocation, error) {
	return s.allocationDao.FindRecent(ctx, 5)
}

// ValidateHostOwnership 确认机器归属当前用户
// CodeX 2026-02-04: enforce dataset mount authorization by allocation check.
func (s *AllocationService) ValidateHostOwnership(ctx context.Context, hostID string, customerID uint) error {
	_, err := s.allocationDao.FindActiveByHostAndCustomer(ctx, hostID, customerID)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrUnauthorized
		}
		return err
	}
	return nil
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
