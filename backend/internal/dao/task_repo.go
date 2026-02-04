package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type TaskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) *TaskDao {
	return &TaskDao{db: db}
}

func (d *TaskDao) Create(ctx context.Context, task *entity.Task) error {
	return d.db.WithContext(ctx).Create(task).Error
}

func (d *TaskDao) ListByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Task, int64, error) {
	var tasks []entity.Task
	var total int64

	db := d.db.WithContext(ctx).Model(&entity.Task{}).Where("customer_id = ?", customerID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("Image").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (d *TaskDao) UpdateStatus(ctx context.Context, id string, status string) error {
	return d.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).Update("status", status).Error
}

// FindByID 根据ID查询任务
// @author Claude
// @description 根据任务ID查询任务详情，用于权限校验和任务详情展示
// @param id 任务ID
// @return 任务实体、错误
// @modified 2026-02-04
func (d *TaskDao) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	var task entity.Task
	if err := d.db.WithContext(ctx).First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}