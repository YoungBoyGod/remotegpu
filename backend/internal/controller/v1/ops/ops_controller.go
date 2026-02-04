package ops

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	common.BaseController
	dashboardService *serviceOps.DashboardService
}

func NewDashboardController(ds *serviceOps.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: ds,
	}
}

// GetStats 获取仪表盘统计数据
// @author Claude
// @modified 2026-02-04
func (c *DashboardController) GetStats(ctx *gin.Context) {
	stats, err := c.dashboardService.GetAggregatedStats(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取统计数据失败")
		return
	}
	c.Success(ctx, stats)
}

// GetGPUTrend 获取 GPU 使用趋势
// @author Claude
// @modified 2026-02-04
// @reason 统一错误信息为中文，移除多余的 data 包装保持返回格式一致
func (c *DashboardController) GetGPUTrend(ctx *gin.Context) {
	trend, err := c.dashboardService.GetGPUTrend(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取趋势数据失败")
		return
	}
	c.Success(ctx, trend)
}

// GetRecentAllocations 获取最近分配记录
// @author Claude
// @modified 2026-02-04
// @reason 统一错误信息为中文，移除多余的 data 包装
func (c *DashboardController) GetRecentAllocations(ctx *gin.Context) {
	allocations, err := c.dashboardService.GetRecentAllocations(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取分配记录失败")
		return
	}
	c.Success(ctx, allocations)
}

type MonitorController struct {
	common.BaseController
	monitorService *serviceOps.MonitorService
}

func NewMonitorController(ms *serviceOps.MonitorService) *MonitorController {
	return &MonitorController{
		monitorService: ms,
	}
}

func (c *MonitorController) GetRealtime(ctx *gin.Context) {
	snapshot, err := c.monitorService.GetGlobalSnapshot(ctx)
	if err != nil {
		c.Error(ctx, 500, "Failed to get snapshot")
		return
	}
	c.Success(ctx, snapshot)
}

type AlertController struct {
	common.BaseController
	opsService *serviceOps.OpsService
}

func NewAlertController(os *serviceOps.OpsService) *AlertController {
	return &AlertController{
		opsService: os,
	}
}

// List 获取告警列表
// @author Claude
// @description 支持分页和筛选的告警列表查询
// @reason 原实现无分页和筛选，现添加完整功能
// @modified 2026-02-04
func (c *AlertController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	severity := ctx.Query("severity")

	var acknowledged *bool
	if ackStr := ctx.Query("acknowledged"); ackStr != "" {
		ack := ackStr == "true"
		acknowledged = &ack
	}

	alerts, total, err := c.opsService.ListAlerts(ctx, page, pageSize, severity, acknowledged)
	if err != nil {
		c.Error(ctx, 500, "获取告警列表失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      alerts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Acknowledge 确认告警
// @author Claude
// @description 将告警标记为已确认
// @modified 2026-02-04
func (c *AlertController) Acknowledge(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的告警ID")
		return
	}

	if err := c.opsService.AcknowledgeAlert(ctx, uint(id)); err != nil {
		c.Error(ctx, 500, "确认告警失败")
		return
	}

	c.Success(ctx, gin.H{"message": "告警已确认"})
}