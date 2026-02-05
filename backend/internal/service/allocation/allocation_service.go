package allocation

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/errors"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AgentClient Agent 客户端接口（避免循环依赖）
type AgentClient interface {
	ResetSSH(ctx context.Context, hostID string) error
	CleanupMachine(ctx context.Context, hostID string) error
}

const (
	machineActionQueueKey   = "machine:action:queue"
	machineActionRetryKey   = "machine:action:retry"
	machineActionPayloadKey = "machine:action:payload"
)

const (
	machineActionResetSSH = "reset_ssh"
	machineActionCleanup  = "cleanup"
)

type machineActionPayload struct {
	Action string `json:"action"`
	HostID string `json:"host_id"`
}

type AllocationService struct {
	db            *gorm.DB
	allocationDao *dao.AllocationDao
	machineDao    *dao.MachineDao
	auditService  *audit.AuditService
	agentClient   AgentClient
	redisClient   *redis.Client
	actionRetries int
	actionDelay   time.Duration
}

func NewAllocationService(db *gorm.DB, auditSvc *audit.AuditService, agentClient AgentClient) *AllocationService {
	actionRetries := 3
	actionDelay := 10 * time.Second
	if config.GlobalConfig != nil {
		if config.GlobalConfig.MachineAction.MaxRetries >= 0 {
			actionRetries = config.GlobalConfig.MachineAction.MaxRetries
		}
		if config.GlobalConfig.MachineAction.RetryDelay > 0 {
			actionDelay = time.Duration(config.GlobalConfig.MachineAction.RetryDelay) * time.Second
		}
	}
	return &AllocationService{
		db:            db,
		allocationDao: dao.NewAllocationDao(db),
		machineDao:    dao.NewMachineDao(db),
		auditService:  auditSvc,
		agentClient:   agentClient,
		redisClient:   cache.GetRedis(),
		actionRetries: actionRetries,
		actionDelay:   actionDelay,
	}
}

func (s *AllocationService) StartWorker(ctx context.Context) {
	if s.redisClient == nil {
		logger.GetLogger().Warn("Machine action queue disabled: redis client not initialized")
		return
	}
	go s.runWorker(ctx)
}

func (s *AllocationService) enqueueAction(ctx context.Context, payload machineActionPayload) error {
	if s.redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return s.redisClient.RPush(ctx, machineActionQueueKey, string(data)).Err()
}

func (s *AllocationService) runWorker(ctx context.Context) {
	s.requeueRetries(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result, err := s.redisClient.BLPop(ctx, 5*time.Second, machineActionQueueKey).Result()
		if err != nil {
			if stderrors.Is(err, redis.Nil) {
				continue
			}
			logger.GetLogger().Warn(fmt.Sprintf("Machine action queue error: %v", err))
			continue
		}
		if len(result) < 2 {
			continue
		}

		var payload machineActionPayload
		if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Invalid machine action payload: %v", err))
			continue
		}
		if payload.Action == "" || payload.HostID == "" {
			logger.GetLogger().Warn("Invalid machine action payload: missing action or host")
			continue
		}

		if s.agentClient == nil {
			logger.GetLogger().Warn("Machine action skipped: agent client not initialized")
			continue
		}

		if err := s.executeAction(ctx, payload); err != nil {
			s.handleActionFailure(ctx, payload, err)
			continue
		}
		s.clearActionRetry(ctx, payload)
	}
}

func (s *AllocationService) requeueRetries(ctx context.Context) {
	if s.redisClient == nil {
		return
	}
	payloads, err := s.redisClient.HGetAll(ctx, machineActionPayloadKey).Result()
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("Failed to requeue machine actions: %v", err))
		return
	}
	for _, payload := range payloads {
		if err := s.redisClient.RPush(ctx, machineActionQueueKey, payload).Err(); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to requeue machine action: %v", err))
		}
	}
}

func (s *AllocationService) executeAction(ctx context.Context, payload machineActionPayload) error {
	switch payload.Action {
	case machineActionResetSSH:
		return s.agentClient.ResetSSH(ctx, payload.HostID)
	case machineActionCleanup:
		return s.agentClient.CleanupMachine(ctx, payload.HostID)
	default:
		return fmt.Errorf("unknown machine action: %s", payload.Action)
	}
}

func (s *AllocationService) handleActionFailure(ctx context.Context, payload machineActionPayload, err error) {
	if s.redisClient == nil {
		logger.GetLogger().Warn(fmt.Sprintf("Machine action failed: %v", err))
		return
	}
	if s.actionRetries <= 0 {
		_ = s.redisClient.HDel(ctx, machineActionRetryKey, s.retryKey(payload)).Err()
		_ = s.redisClient.HDel(ctx, machineActionPayloadKey, s.retryKey(payload)).Err()
		logger.GetLogger().Warn(fmt.Sprintf("Machine action failed: %v", err))
		return
	}

	key := s.retryKey(payload)
	retryCount, retryErr := s.redisClient.HIncrBy(ctx, machineActionRetryKey, key, 1).Result()
	if retryErr != nil {
		logger.GetLogger().Warn(fmt.Sprintf("Machine action retry error: %v", retryErr))
		return
	}

	if retryCount > int64(s.actionRetries) {
		_ = s.redisClient.HDel(ctx, machineActionRetryKey, key).Err()
		_ = s.redisClient.HDel(ctx, machineActionPayloadKey, key).Err()
		logger.GetLogger().Warn(fmt.Sprintf("Machine action exceeded retries: %v", err))
		return
	}

	if data, marshalErr := json.Marshal(payload); marshalErr == nil {
		_ = s.redisClient.HSet(ctx, machineActionPayloadKey, key, string(data)).Err()
	}
	s.scheduleRetry(payload)
}

func (s *AllocationService) scheduleRetry(payload machineActionPayload) {
	go func() {
		<-time.After(s.actionDelay)
		retryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.enqueueAction(retryCtx, payload); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to requeue machine action: %v", err))
		}
	}()
}

func (s *AllocationService) clearActionRetry(ctx context.Context, payload machineActionPayload) {
	if s.redisClient == nil {
		return
	}
	key := s.retryKey(payload)
	_ = s.redisClient.HDel(ctx, machineActionRetryKey, key).Err()
	_ = s.redisClient.HDel(ctx, machineActionPayloadKey, key).Err()
}

func (s *AllocationService) retryKey(payload machineActionPayload) string {
	return fmt.Sprintf("%s:%s", payload.Action, payload.HostID)
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

	// 异步触发 Agent 重置 SSH
	if s.agentClient != nil {
		enqueueCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := s.enqueueAction(enqueueCtx, machineActionPayload{Action: machineActionResetSSH, HostID: hostID}); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to enqueue reset ssh: %v", err))
		}
		cancel()
	}

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

	// 5. 异步触发清理流程（重置SSH、清理用户数据等）
	if s.agentClient != nil {
		enqueueCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := s.enqueueAction(enqueueCtx, machineActionPayload{Action: machineActionCleanup, HostID: hostID}); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to enqueue cleanup: %v", err))
		}
		cancel()
	}

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
