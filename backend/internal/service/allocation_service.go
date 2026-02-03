package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AllocationService struct {
	db            *gorm.DB
	allocationDao *dao.AllocationDao
	machineDao    *dao.MachineDao
}

func NewAllocationService(db *gorm.DB) *AllocationService {
	return &AllocationService{
		db:            db,
		allocationDao: dao.NewAllocationDao(db),
		machineDao:    dao.NewMachineDao(db),
	}
}

func (s *AllocationService) AllocateMachine(ctx context.Context, customerID uint, hostID string, durationMonths int, remark string) (*entity.Allocation, error) {
	// Transaction to ensure atomicity
	var allocation *entity.Allocation
	
	err := s.db.Transaction(func(tx *gorm.DB) error {
		machineDao := dao.NewMachineDao(tx)
		allocationDao := dao.NewAllocationDao(tx)

		// 1. Check machine status (Lock row if possible, but status check is minimal)
		host, err := machineDao.FindByID(ctx, hostID)
		if err != nil {
			return err
		}
		if host.Status != "idle" && host.Status != "online" { // Assuming 'online' means available/idle for simplicity or explicit 'idle' state
			return errors.New("machine is not available")
		}

		// 2. Update machine status
		if err := machineDao.UpdateStatus(ctx, hostID, "allocated"); err != nil {
			return err
		}

		// 3. Create allocation record
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
			return err
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	// TODO: Trigger Agent to reset SSH (Async)
	// go s.agentService.ResetMachine(hostID)

	return allocation, nil
}

func (s *AllocationService) ReclaimMachine(ctx context.Context, hostID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		machineDao := dao.NewMachineDao(tx)
		allocationDao := dao.NewAllocationDao(tx)

		// 1. Find active allocation
		alloc, err := allocationDao.FindActiveByHostID(ctx, hostID)
		if err != nil {
			return fmt.Errorf("active allocation not found: %w", err)
		}

		// 2. Update allocation status
		now := time.Now()
		alloc.Status = "reclaimed"
		alloc.ActualEndTime = &now
		// In a real dao, we would have a specific update method for this, or use generic Update
		if err := tx.Save(alloc).Error; err != nil {
			return err
		}

		// 3. Update machine status (to maintenance/wiping)
		if err := machineDao.UpdateStatus(ctx, hostID, "maintenance"); err != nil {
			return err
		}

		return nil
	})
	
	// TODO: Trigger Wipe process
}

func (s *AllocationService) GetRecent(ctx context.Context) ([]entity.Allocation, error) {
	return s.allocationDao.FindRecent(ctx, 5)
}
