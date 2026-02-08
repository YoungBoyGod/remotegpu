package notification

import (
	"context"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupNotificationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER NOT NULL,
		title VARCHAR(256) NOT NULL,
		content TEXT,
		type VARCHAR(32) NOT NULL,
		level VARCHAR(20) DEFAULT 'info',
		is_read INTEGER DEFAULT 0,
		read_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	require.NoError(t, err)

	return db
}

// newTestNotificationService 创建测试用通知服务
func newTestNotificationService(t *testing.T) (*NotificationService, *gorm.DB) {
	db := setupNotificationTestDB(t)
	hub := NewSSEHub()
	svc := NewNotificationService(db, hub)
	return svc, db
}

// ==================== CreateAndPush 测试 ====================

func TestCreateAndPush_Success(t *testing.T) {
	svc, _ := newTestNotificationService(t)

	n := &entity.Notification{
		CustomerID: 1,
		Title:      "测试通知",
		Content:    "这是一条测试通知",
		Type:       "system",
		Level:      "info",
	}
	err := svc.CreateAndPush(context.Background(), n)
	require.NoError(t, err)
	assert.NotZero(t, n.ID)
}

// ==================== List 测试 ====================

func TestList_Empty(t *testing.T) {
	svc, _ := newTestNotificationService(t)

	list, total, err := svc.List(context.Background(), 1, false, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}

func TestList_WithData(t *testing.T) {
	svc, _ := newTestNotificationService(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_ = svc.CreateAndPush(ctx, &entity.Notification{
			CustomerID: 1, Title: "通知", Type: "system",
		})
	}

	list, total, err := svc.List(ctx, 1, false, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)
}

// ==================== CountUnread 测试 ====================

func TestCountUnread(t *testing.T) {
	svc, db := newTestNotificationService(t)
	ctx := context.Background()

	// 创建 3 条通知（默认未读）
	for i := 0; i < 3; i++ {
		_ = svc.CreateAndPush(ctx, &entity.Notification{
			CustomerID: 1, Title: "通知", Type: "system",
		})
	}
	// 标记第 1 条为已读
	db.Exec(`UPDATE notifications SET is_read = 1 WHERE id = 1`)

	count, err := svc.CountUnread(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// ==================== MarkRead 测试 ====================

func TestMarkRead_Success(t *testing.T) {
	svc, _ := newTestNotificationService(t)
	ctx := context.Background()

	_ = svc.CreateAndPush(ctx, &entity.Notification{
		CustomerID: 1, Title: "待读", Type: "system",
	})

	err := svc.MarkRead(ctx, 1, 1)
	require.NoError(t, err)

	count, err := svc.CountUnread(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestMarkRead_WrongCustomer(t *testing.T) {
	svc, _ := newTestNotificationService(t)
	ctx := context.Background()

	_ = svc.CreateAndPush(ctx, &entity.Notification{
		CustomerID: 1, Title: "通知", Type: "system",
	})

	// 用错误的 customerID 标记，应返回 not found
	err := svc.MarkRead(ctx, 1, 999)
	assert.Error(t, err)
}

// ==================== MarkAllRead 测试 ====================

func TestMarkAllRead_Success(t *testing.T) {
	svc, _ := newTestNotificationService(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_ = svc.CreateAndPush(ctx, &entity.Notification{
			CustomerID: 1, Title: "通知", Type: "system",
		})
	}

	err := svc.MarkAllRead(ctx, 1)
	require.NoError(t, err)

	count, err := svc.CountUnread(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
