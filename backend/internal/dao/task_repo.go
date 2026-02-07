package dao

import (
	"context"
	"time"

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

func (d *TaskDao) UpdateErrorMsg(ctx context.Context, id string, msg string) error {
	return d.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).Update("error_msg", msg).Error
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

// ClaimTasks 原子性认领任务
func (d *TaskDao) ClaimTasks(ctx context.Context, machineID, agentID string, limit int) ([]entity.Task, error) {
	var tasks []entity.Task

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询待认领的任务 ID
		var pendingTasks []entity.Task
		if err := tx.Where("machine_id = ? AND status = ?", machineID, "pending").
			Order("priority, created_at").
			Limit(limit).
			Find(&pendingTasks).Error; err != nil {
			return err
		}

		if len(pendingTasks) == 0 {
			return nil
		}

		ids := make([]string, len(pendingTasks))
		for i, t := range pendingTasks {
			ids[i] = t.ID
		}

		// 2. 批量更新状态
		now := time.Now()
		leaseExpires := now.Add(5 * time.Minute)
		if err := tx.Model(&entity.Task{}).
			Where("id IN ? AND status = ?", ids, "pending").
			Updates(map[string]interface{}{
				"status":            "assigned",
				"assigned_agent_id": agentID,
				"assigned_at":       now,
				"lease_expires_at":  leaseExpires,
				"attempt_id":        gorm.Expr("gen_random_uuid()"),
			}).Error; err != nil {
			return err
		}

		// 3. 在事务内重新查询，确保拿到 attempt_id 等更新后的字段
		return tx.Where("id IN ?", ids).Find(&tasks).Error
	})

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// StartTask 标记任务开始
func (d *TaskDao) StartTask(ctx context.Context, id, agentID, attemptID string) error {
	now := time.Now()
	result := d.db.WithContext(ctx).Model(&entity.Task{}).
		Where("id = ? AND assigned_agent_id = ? AND attempt_id = ?", id, agentID, attemptID).
		Updates(map[string]interface{}{
			"status":     "running",
			"started_at": now,
		})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// RenewLease 续约租约
func (d *TaskDao) RenewLease(ctx context.Context, id, agentID, attemptID string, extendSec int) error {
	if extendSec <= 0 {
		extendSec = 300
	}
	leaseExpires := time.Now().Add(time.Duration(extendSec) * time.Second)

	result := d.db.WithContext(ctx).Model(&entity.Task{}).
		Where("id = ? AND assigned_agent_id = ? AND attempt_id = ? AND status = ?", id, agentID, attemptID, "running").
		Update("lease_expires_at", leaseExpires)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// CompleteTask 完成任务
func (d *TaskDao) CompleteTask(ctx context.Context, id, agentID, attemptID string, exitCode int, errMsg string) error {
	now := time.Now()
	status := "completed"
	if exitCode != 0 {
		status = "failed"
	}

	result := d.db.WithContext(ctx).Model(&entity.Task{}).
		Where("id = ? AND assigned_agent_id = ? AND attempt_id = ?", id, agentID, attemptID).
		Updates(map[string]interface{}{
			"status":    status,
			"exit_code": exitCode,
			"error_msg": errMsg,
			"ended_at":  now,
		})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// ListAll 管理员查询所有任务（支持分页和状态过滤）
func (d *TaskDao) ListAll(ctx context.Context, status string, page, pageSize int) ([]entity.Task, int64, error) {
	var tasks []entity.Task
	var total int64

	db := d.db.WithContext(ctx).Model(&entity.Task{})
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("Image").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// CancelTask 取消任务（仅限 queued/pending/assigned 状态）
func (d *TaskDao) CancelTask(ctx context.Context, id string) error {
	now := time.Now()
	result := d.db.WithContext(ctx).Model(&entity.Task{}).
		Where("id = ? AND status IN ?", id, []string{"queued", "pending", "assigned"}).
		Updates(map[string]interface{}{
			"status":   "cancelled",
			"ended_at": now,
		})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// RetryTask 重试任务（重置状态为 queued，增加重试计数）
func (d *TaskDao) RetryTask(ctx context.Context, id string) error {
	result := d.db.WithContext(ctx).Model(&entity.Task{}).
		Where("id = ? AND status IN ?", id, []string{"failed", "cancelled", "stopped"}).
		Updates(map[string]interface{}{
			"status":      "queued",
			"retry_count": gorm.Expr("retry_count + 1"),
			"exit_code":   0,
			"error_msg":   "",
			"ended_at":    nil,
			"started_at":  nil,
			"assigned_at": nil,
			"attempt_id":  "",
		})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// CountRunningTasksByMachineID 统计指定机器上运行中的任务数量
// @author Claude
// @description 用于机器回收前检查是否有运行中的任务
// @param machineID 机器ID
// @return 运行中任务数量、错误
// @modified 2026-02-06
func (d *TaskDao) CountRunningTasksByMachineID(ctx context.Context, machineID string) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Task{}).
		Where("machine_id = ? AND status IN ?", machineID, []string{"running", "assigned"}).
		Count(&count).Error
	return count, err
}
