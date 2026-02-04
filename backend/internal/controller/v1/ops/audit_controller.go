package ops

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"github.com/gin-gonic/gin"
)

type AuditController struct {
	common.BaseController
	auditService *audit.AuditService
}

func NewAuditController(svc *audit.AuditService) *AuditController {
	return &AuditController{
		auditService: svc,
	}
}

// List 获取审计日志列表
func (c *AuditController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	params := dao.AuditListParams{
		Page:         page,
		PageSize:     pageSize,
		Action:       ctx.Query("action"),
		ResourceType: ctx.Query("resource_type"),
		Username:     ctx.Query("username"),
		StartTime:    ctx.Query("start_time"),
		EndTime:      ctx.Query("end_time"),
	}

	logs, total, err := c.auditService.ListLogs(ctx, params)
	if err != nil {
		c.Error(ctx, 500, "获取审计日志失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
