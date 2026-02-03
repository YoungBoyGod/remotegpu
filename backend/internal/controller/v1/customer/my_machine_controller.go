package customer

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	"github.com/gin-gonic/gin"
)

type MyMachineController struct {
	common.BaseController
	machineService *serviceMachine.MachineService
	agentService   *serviceOps.AgentService
}

func NewMyMachineController(ms *serviceMachine.MachineService, as *serviceOps.AgentService) *MyMachineController {
	return &MyMachineController{
		machineService: ms,
		agentService:   as,
	}
}

func (c *MyMachineController) List(ctx *gin.Context) {
	// 实际逻辑中，应该根据 token 中的 tenantID 进行过滤
	// tenantID := ctx.GetUint("tenantID")
	// MVP 阶段，暂时返回所有数据或使用模拟逻辑
	
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	// 假设 ListMachines 支持租户过滤，但服务层尚未实现
	// 所以目前调用通用的 List 方法
	machines, total, err := c.machineService.ListMachines(ctx, page, pageSize, nil)
	if err != nil {
		c.Error(ctx, 500, "Failed to list machines")
		return
	}

	c.Success(ctx, gin.H{
		"list":      machines,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (c *MyMachineController) GetConnection(ctx *gin.Context) {
	hostID := ctx.Param("id")
	info, err := c.machineService.GetConnectionInfo(ctx, hostID)
	if err != nil {
		c.Error(ctx, 500, "Failed to get connection info")
		return
	}
	c.Success(ctx, info)
}

func (c *MyMachineController) ResetSSH(ctx *gin.Context) {
	hostID := ctx.Param("id")
	if err := c.agentService.ResetSSH(ctx, hostID); err != nil {
		c.Error(ctx, 500, "Failed to reset SSH")
		return
	}
	c.Success(ctx, gin.H{"message": "SSH reset triggered"})
}