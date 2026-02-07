package task

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	"gorm.io/gorm"
)

type TaskService struct {
	taskDao      *dao.TaskDao
	agentService *serviceOps.AgentService
}

func NewTaskService(db *gorm.DB, agentSvc *serviceOps.AgentService) *TaskService {
	return &TaskService{
		taskDao:      dao.NewTaskDao(db),
		agentService: agentSvc,
	}
}

func (s *TaskService) ListTasks(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Task, int64, error) {
	return s.taskDao.ListByCustomerID(ctx, customerID, page, pageSize)
}

func (s *TaskService) SubmitTask(ctx context.Context, task *entity.Task) error {
	task.Status = "queued"
	if err := s.taskDao.Create(ctx, task); err != nil {
		return err
	}
	middleware.TasksCreatedTotal.Inc()
	return nil
}

func (s *TaskService) StopTask(ctx context.Context, id string) error {
	task, err := s.taskDao.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.stopTaskProcess(ctx, task); err != nil {
		return err
	}

	return s.taskDao.UpdateStatus(ctx, id, "stopped")
}

// GetTask 根据ID获取任务
// @author Claude
// @description 获取任务详情，用于权限校验
// @modified 2026-02-04
func (s *TaskService) GetTask(ctx context.Context, id string) (*entity.Task, error) {
	return s.taskDao.FindByID(ctx, id)
}

// StopTaskWithAuth 停止任务（带权限校验）
// @author Claude
// @description 停止任务前校验任务是否属于当前用户，防止越权操作
// @param id 任务ID
// @param customerID 当前用户ID（从JWT获取）
// @return 错误
// @modified 2026-02-04
func (s *TaskService) StopTaskWithAuth(ctx context.Context, id string, customerID uint) error {
	task, err := s.taskDao.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if task.CustomerID != customerID {
		return entity.ErrUnauthorized
	}

	if err := s.stopTaskProcess(ctx, task); err != nil {
		return err
	}

	return s.taskDao.UpdateStatus(ctx, id, "stopped")
}

func (s *TaskService) stopTaskProcess(ctx context.Context, task *entity.Task) error {
	// CodeX 2026-02-05: validate process_id/host_id before stop and record failures.
	if s.agentService == nil {
		err := fmt.Errorf("agent service unavailable")
		_ = s.taskDao.UpdateErrorMsg(ctx, task.ID, err.Error())
		return err
	}
	if task.HostID == "" {
		err := fmt.Errorf("host id missing")
		_ = s.taskDao.UpdateErrorMsg(ctx, task.ID, err.Error())
		return err
	}
	if task.ProcessID <= 0 {
		err := fmt.Errorf("process id missing")
		_ = s.taskDao.UpdateErrorMsg(ctx, task.ID, err.Error())
		return err
	}

	// 修复 P0 问题：添加超时控制，防止 Agent 无响应时请求挂起
	stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.agentService.StopProcess(stopCtx, task.HostID, task.ProcessID); err != nil {
		_ = s.taskDao.UpdateErrorMsg(ctx, task.ID, err.Error())
		return err
	}
	return nil
}

// === Agent 专用 API ===

// ClaimTasks Agent 认领任务
func (s *TaskService) ClaimTasks(ctx context.Context, machineID, agentID string, limit int) ([]entity.Task, error) {
	return s.taskDao.ClaimTasks(ctx, machineID, agentID, limit)
}

// StartTask 标记任务开始
func (s *TaskService) StartTask(ctx context.Context, id, agentID, attemptID string) error {
	return s.taskDao.StartTask(ctx, id, agentID, attemptID)
}

// RenewLease 续约租约
func (s *TaskService) RenewLease(ctx context.Context, id, agentID, attemptID string, extendSec int) error {
	return s.taskDao.RenewLease(ctx, id, agentID, attemptID, extendSec)
}

// CompleteTask 完成任务
func (s *TaskService) CompleteTask(ctx context.Context, id, agentID, attemptID string, exitCode int, errMsg string) error {
	if err := s.taskDao.CompleteTask(ctx, id, agentID, attemptID, exitCode, errMsg); err != nil {
		return err
	}
	middleware.TasksCompletedTotal.Inc()
	return nil
}
