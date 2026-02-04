package ops

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
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
// @Summary 获取仪表盘统计数据
// @Description 获取系统整体统计信息，包括机器总数、活跃客户等
// @Tags Admin - Dashboard
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/dashboard/stats [get]
func (c *DashboardController) GetStats(ctx *gin.Context) {
	stats, err := c.dashboardService.GetAggregatedStats(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取统计数据失败")
		return
	}
	c.Success(ctx, stats)
}

// GetGPUTrend 获取 GPU 使用趋势
// @Summary 获取 GPU 使用趋势
// @Description 获取最近时间段的 GPU 利用率趋势数据
// @Tags Admin - Dashboard
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/dashboard/gpu-trend [get]
func (c *DashboardController) GetGPUTrend(ctx *gin.Context) {
	trend, err := c.dashboardService.GetGPUTrend(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取趋势数据失败")
		return
	}
	c.Success(ctx, trend)
}

// GetRecentAllocations 获取最近分配记录
// @Summary 获取最近分配记录
// @Description 获取最近的机器分配记录
// @Tags Admin - Dashboard
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} entity.Allocation
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/allocations/recent [get]
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

// GetRealtime 获取实时监控数据
// @Summary 获取实时监控数据
// @Description 获取系统实时监控快照，包括机器在线状态等
// @Tags Admin - Monitoring
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/monitoring/realtime [get]
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
// @Summary 获取告警列表
// @Description 获取系统告警信息，支持分页和筛选
// @Tags Admin - Monitoring
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param severity query string false "严重程度筛选"
// @Param acknowledged query boolean false "是否已确认"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/alerts [get]
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
// @Summary 确认告警
// @Description 将指定告警标记为已确认
// @Tags Admin - Monitoring
// @Accept json
// @Produce json
// @Param id path int true "告警ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/alerts/{id}/acknowledge [post]
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