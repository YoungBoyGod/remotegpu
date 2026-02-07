package notification

import (
	"fmt"
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceNotification "github.com/YoungBoyGod/remotegpu/internal/service/notification"
	"github.com/gin-gonic/gin"
)

// NotificationController 通知控制器
type NotificationController struct {
	common.BaseController
	notificationService *serviceNotification.NotificationService
}

// NewNotificationController 创建通知控制器
func NewNotificationController(svc *serviceNotification.NotificationService) *NotificationController {
	return &NotificationController{notificationService: svc}
}

// SSE 建立 Server-Sent Events 连接
func (c *NotificationController) SSE(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		c.Error(ctx, 401, "未认证")
		return
	}
	customerID := userID.(uint)

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")

	hub := c.notificationService.GetHub()
	client := &serviceNotification.SSEClient{
		CustomerID: customerID,
		Channel:    make(chan serviceNotification.SSEEvent, 64),
	}
	hub.Register(client)
	defer hub.Unregister(client)

	// 发送初始连接成功事件
	fmt.Fprintf(ctx.Writer, "event: connected\ndata: {\"status\":\"ok\"}\n\n")
	ctx.Writer.Flush()

	clientGone := ctx.Request.Context().Done()
	for {
		select {
		case <-clientGone:
			return
		case event, ok := <-client.Channel:
			if !ok {
				return
			}
			fmt.Fprintf(ctx.Writer, "event: %s\ndata: %s\n\n", event.Event, event.Data)
			ctx.Writer.Flush()
		}
	}
}

// List 查询通知列表
func (c *NotificationController) List(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	customerID := userID.(uint)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	onlyUnread := ctx.Query("unread") == "true"

	list, total, err := c.notificationService.List(ctx, customerID, onlyUnread, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"list": list, "total": total})
}

// UnreadCount 查询未读数量
func (c *NotificationController) UnreadCount(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	customerID := userID.(uint)

	count, err := c.notificationService.CountUnread(ctx, customerID)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"count": count})
}

// MarkRead 标记单条已读
func (c *NotificationController) MarkRead(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	customerID := userID.(uint)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的通知 ID")
		return
	}

	if err := c.notificationService.MarkRead(ctx, uint(id), customerID); err != nil {
		c.Error(ctx, 404, "通知不存在")
		return
	}
	c.Success(ctx, nil)
}

// MarkAllRead 标记全部已读
func (c *NotificationController) MarkAllRead(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	customerID := userID.(uint)

	if err := c.notificationService.MarkAllRead(ctx, customerID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, nil)
}
