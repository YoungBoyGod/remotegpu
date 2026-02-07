package agent

import (
	v1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/gin-gonic/gin"
)

// HeartbeatController Agent 心跳与注册控制器
type HeartbeatController struct {
	common.BaseController
	machineSvc *serviceMachine.MachineService
}

func NewHeartbeatController(machineSvc *serviceMachine.MachineService) *HeartbeatController {
	return &HeartbeatController{machineSvc: machineSvc}
}

// Heartbeat 处理 Agent 心跳上报（支持携带监控指标）
func (c *HeartbeatController) Heartbeat(ctx *gin.Context) {
	var req v1.HeartbeatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	// 转换指标数据
	var metrics *serviceMachine.HeartbeatMetrics
	if req.Metrics != nil {
		metrics = &serviceMachine.HeartbeatMetrics{
			CPUUsagePercent:    req.Metrics.CPUUsagePercent,
			MemoryUsagePercent: req.Metrics.MemoryUsagePercent,
			MemoryUsedGB:       req.Metrics.MemoryUsedGB,
			DiskUsagePercent:   req.Metrics.DiskUsagePercent,
			DiskUsedGB:         req.Metrics.DiskUsedGB,
		}
		for _, g := range req.Metrics.GPUMetrics {
			metrics.GPUMetrics = append(metrics.GPUMetrics, serviceMachine.GPUMetricData{
				Index:         g.Index,
				UUID:          g.UUID,
				Name:          g.Name,
				UtilPercent:   g.UtilPercent,
				MemoryUsedMB:  g.MemoryUsedMB,
				MemoryTotalMB: g.MemoryTotalMB,
				TemperatureC:  g.TemperatureC,
				PowerUsageW:   g.PowerUsageW,
			})
		}
	}

	if err := c.machineSvc.HeartbeatWithMetrics(ctx, req.MachineID, metrics); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "ok"})
}

// Register 处理 Agent 注册
func (c *HeartbeatController) Register(ctx *gin.Context) {
	var req v1.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	info := &serviceMachine.AgentRegistration{
		AgentID:    req.AgentID,
		MachineID:  req.MachineID,
		Version:    req.Version,
		Hostname:   req.Hostname,
		IPAddress:  req.IPAddress,
		AgentPort:  req.AgentPort,
		MaxWorkers: req.MaxWorkers,
	}

	if err := c.machineSvc.RegisterAgent(ctx, info); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "registered", "machine_id": req.MachineID})
}
