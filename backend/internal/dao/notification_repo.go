package dao

import (
	"context"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type NotificationDao struct {
	db *gorm.DB
}

func NewNotificationDao(db *gorm.DB) *NotificationDao {
	return &NotificationDao{db: db}
}

func (d *NotificationDao) Create(ctx context.Context, n *entity.Notification) error {
	return d.db.WithContext(ctx).Create(n).Error
}

// ListByCustomerID 查询用户通知（分页）
func (d *NotificationDao) ListByCustomerID(ctx context.Context, customerID uint, onlyUnread bool, page, pageSize int) ([]entity.Notification, int64, error) {
	var list []entity.Notification
	var total int64

	db := d.db.WithContext(ctx).Model(&entity.Notification{}).Where("customer_id = ?", customerID)
	if onlyUnread {
		db = db.Where("is_read = ?", false)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// CountUnread 统计未读数量
func (d *NotificationDao) CountUnread(ctx context.Context, customerID uint) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("customer_id = ? AND is_read = ?", customerID, false).
		Count(&count).Error
	return count, err
}

// MarkRead 标记单条已读
func (d *NotificationDao) MarkRead(ctx context.Context, id, customerID uint) error {
	now := time.Now()
	result := d.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ? AND customer_id = ?", id, customerID).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// MarkAllRead 标记全部已读
func (d *NotificationDao) MarkAllRead(ctx context.Context, customerID uint) error {
	now := time.Now()
	return d.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("customer_id = ? AND is_read = ?", customerID, false).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		}).Error
}
