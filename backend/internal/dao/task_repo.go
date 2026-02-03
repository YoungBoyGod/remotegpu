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