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
// @Summary 建立 SSE 实时通知连接
// @Description 建立 Server-Sent Events 长连接，实时推送通知事件
// @Tags Customer - Notifications
// @Produce text/event-stream
// @Security Bearer
// @Success 200 {string} string "SSE 事件流"
// @Failure 401 {object} common.ErrorResponse
// @Router /customer/notifications/sse [get]
func (c *NotificationController) SSE(ctx *gin.Context) {
	customerID := ctx.GetUint("userID")
	if customerID == 0 {
		c.Error(ctx, 401, "未认证")
		return
	}

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
// @Summary 获取通知列表
// @Description 分页获取当前用户的通知列表，支持筛选未读通知
// @Tags Customer - Notifications
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param unread query string false "仅未读（true/false）"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/notifications [get]
func (c *NotificationController) List(ctx *gin.Context) {
	customerID := ctx.GetUint("userID")
	if customerID == 0 {
		c.Error(ctx, 401, "未认证")
		return
	}

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
// @Summary 获取未读通知数量
// @Description 获取当前用户的未读通知数量
// @Tags Customer - Notifications
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/notifications/unread-count [get]
func (c *NotificationController) UnreadCount(ctx *gin.Context) {
	customerID := ctx.GetUint("userID")
	if customerID == 0 {
		c.Error(ctx, 401, "未认证")
		return
	}

	count, err := c.notificationService.CountUnread(ctx, customerID)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"count": count})
}

// MarkRead 标记单条已读
// @Summary 标记通知已读
// @Description 将指定通知标记为已读状态
// @Tags Customer - Notifications
// @Produce json
// @Param id path int true "通知 ID"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /customer/notifications/{id}/read [post]
func (c *NotificationController) MarkRead(ctx *gin.Context) {
	customerID := ctx.GetUint("userID")
	if customerID == 0 {
		c.Error(ctx, 401, "未认证")
		return
	}

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
// @Summary 标记全部通知已读
// @Description 将当前用户的所有未读通知标记为已读
// @Tags Customer - Notifications
// @Produce json
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/notifications/read-all [post]
func (c *NotificationController) MarkAllRead(ctx *gin.Context) {
	customerID := ctx.GetUint("userID")
	if customerID == 0 {
		c.Error(ctx, 401, "未认证")
		return
	}

	if err := c.notificationService.MarkAllRead(ctx, customerID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, nil)
}
