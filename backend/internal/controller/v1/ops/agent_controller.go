package ops

import (
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/gin-gonic/gin"
)

// AgentController Agent 管理控制器
type AgentController struct {
	common.BaseController
	machineSvc *serviceMachine.MachineService
}

func NewAgentController(ms *serviceMachine.MachineService) *AgentController {
	return &AgentController{machineSvc: ms}
}

// List 获取 Agent 列表
// @Summary 获取 Agent 列表
// @Description 获取所有部署的 Agent 信息和状态
// @Tags Admin - Agents
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/agents [get]
func (c *AgentController) List(ctx *gin.Context) {
	result, err := c.machineSvc.ListAgents(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取 Agent 列表失败")
		return
	}
	c.Success(ctx, result)
}
