package task

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type TaskService struct {
	taskDao *dao.TaskDao
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		taskDao: dao.NewTaskDao(db),
	}
}

func (s *TaskService) ListTasks(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Task, int64, error) {
	return s.taskDao.ListByCustomerID(ctx, customerID, page, pageSize)
}

func (s *TaskService) SubmitTask(ctx context.Context, task *entity.Task) error {
	task.Status = "queued"
	return s.taskDao.Create(ctx, task)
}

func (s *TaskService) StopTask(ctx context.Context, id string) error {
	// TODO: Call Agent to kill process
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
	// TODO: Call Agent to kill process
	return s.taskDao.UpdateStatus(ctx, id, "stopped")
}