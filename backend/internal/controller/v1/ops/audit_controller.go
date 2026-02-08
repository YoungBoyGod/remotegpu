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
// @Summary 获取审计日志列表
// @Description 分页获取审计日志，支持按操作类型、资源类型、用户名和时间范围筛选
// @Tags Admin - Audit
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param action query string false "操作类型筛选"
// @Param resource_type query string false "资源类型筛选"
// @Param username query string false "用户名筛选"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/audit/logs [get]
func (c *AuditController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	// CodeX 2026-02-04: normalize paging input to avoid invalid offsets.
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}

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
