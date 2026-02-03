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