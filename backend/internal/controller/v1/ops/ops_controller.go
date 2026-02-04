package ops

import (
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

func (c *AlertController) List(ctx *gin.Context) {
	alerts, err := c.opsService.ListActiveAlerts(ctx)
	if err != nil {
		c.Error(ctx, 500, "Failed to list alerts")
		return
	}
	c.Success(ctx, alerts)
}