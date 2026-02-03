package v1

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/gin-gonic/gin"
)

type MyMachineController struct {
	BaseController
	machineService *service.MachineService
	agentService   *service.AgentService
}

func NewMyMachineController(ms *service.MachineService, as *service.AgentService) *MyMachineController {
	return &MyMachineController{
		machineService: ms,
		agentService:   as,
	}
}

func (c *MyMachineController) List(ctx *gin.Context) {
	// In real logic, filter by tenantID from token
	// tenantID := ctx.GetUint("tenantID")
	// For MVP, just return all or mock logic
	
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	// Assuming ListMachines supports tenant filter which isn't implemented in service yet
	// So calling generic List for now
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
