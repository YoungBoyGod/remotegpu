package agent

import (
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/gin-gonic/gin"
)

// HeartbeatController Agent 心跳控制器
type HeartbeatController struct {
	common.BaseController
	machineSvc *serviceMachine.MachineService
}

func NewHeartbeatController(machineSvc *serviceMachine.MachineService) *HeartbeatController {
	return &HeartbeatController{machineSvc: machineSvc}
}

// Heartbeat 处理 Agent 心跳上报
func (c *HeartbeatController) Heartbeat(ctx *gin.Context) {
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
		MachineID string `json:"machine_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.machineSvc.Heartbeat(ctx, req.MachineID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "ok"})
}
