package v1

import (
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	BaseController
	dashboardService *service.DashboardService
}

func NewDashboardController(ds *service.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: ds,
	}
}

func (c *DashboardController) GetStats(ctx *gin.Context) {
	stats, err := c.dashboardService.GetAggregatedStats(ctx)
	if err != nil {
		c.Error(ctx, 500, "Failed to get stats")
		return
	}
	c.Success(ctx, stats)
}

func (c *DashboardController) GetGPUTrend(ctx *gin.Context) {
	// Should call MetricService
	c.Success(ctx, []map[string]interface{}{})
}

func (c *DashboardController) GetRecentAllocations(ctx *gin.Context) {
	// Should call AllocationService
	c.Success(ctx, []map[string]interface{}{})
}

type MonitorController struct {
	BaseController
	monitorService *service.MonitorService
}

func NewMonitorController(ms *service.MonitorService) *MonitorController {
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
	BaseController
	opsService *service.OpsService
}

func NewAlertController(os *service.OpsService) *AlertController {
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
