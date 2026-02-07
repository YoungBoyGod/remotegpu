package notification

import (
	"context"
	"encoding/json"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// NotificationService 通知服务
type NotificationService struct {
	notificationDao *dao.NotificationDao
	hub             *SSEHub
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB, hub *SSEHub) *NotificationService {
	return &NotificationService{
		notificationDao: dao.NewNotificationDao(db),
		hub:             hub,
	}
}

// GetHub 获取 SSE Hub（供 Controller 使用）
func (s *NotificationService) GetHub() *SSEHub {
	return s.hub
}

// CreateAndPush 创建通知并实时推送
func (s *NotificationService) CreateAndPush(ctx context.Context, n *entity.Notification) error {
	if err := s.notificationDao.Create(ctx, n); err != nil {
		return err
	}
	data, _ := json.Marshal(n)
	s.hub.Send(n.CustomerID, SSEEvent{
		Event: "notification",
		Data:  string(data),
	})
	return nil
}

// List 查询通知列表
func (s *NotificationService) List(ctx context.Context, customerID uint, onlyUnread bool, page, pageSize int) ([]entity.Notification, int64, error) {
	return s.notificationDao.ListByCustomerID(ctx, customerID, onlyUnread, page, pageSize)
}

// CountUnread 统计未读数量
func (s *NotificationService) CountUnread(ctx context.Context, customerID uint) (int64, error) {
	return s.notificationDao.CountUnread(ctx, customerID)
}

// MarkRead 标记单条已读
func (s *NotificationService) MarkRead(ctx context.Context, id, customerID uint) error {
	return s.notificationDao.MarkRead(ctx, id, customerID)
}

// MarkAllRead 标记全部已读
func (s *NotificationService) MarkAllRead(ctx context.Context, customerID uint) error {
	return s.notificationDao.MarkAllRead(ctx, customerID)
}

// PushTaskStatusChange 推送任务状态变更通知
func (s *NotificationService) PushTaskStatusChange(ctx context.Context, customerID uint, taskID, status string) error {
	n := &entity.Notification{
		CustomerID: customerID,
		Title:      "任务状态变更",
		Content:    "任务 " + taskID + " 状态变更为 " + status,
		Type:       "task",
		Level:      "info",
	}
	if status == "failed" {
		n.Level = "error"
	}
	return s.CreateAndPush(ctx, n)
}

// PushAlert 推送告警通知
func (s *NotificationService) PushAlert(ctx context.Context, customerID uint, title, content, level string) error {
	n := &entity.Notification{
		CustomerID: customerID,
		Title:      title,
		Content:    content,
		Type:       "alert",
		Level:      level,
	}
	return s.CreateAndPush(ctx, n)
}

// PushMachineStatusChange 推送机器状态变更通知
func (s *NotificationService) PushMachineStatusChange(ctx context.Context, customerID uint, machineID, status string) error {
	n := &entity.Notification{
		CustomerID: customerID,
		Title:      "机器状态变更",
		Content:    "机器 " + machineID + " 状态变更为 " + status,
		Type:       "machine",
		Level:      "info",
	}
	if status == "offline" || status == "error" {
		n.Level = "warning"
	}
	return s.CreateAndPush(ctx, n)
}
