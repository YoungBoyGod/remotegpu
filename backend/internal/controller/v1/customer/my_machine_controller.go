package customer

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceAllocation "github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	"github.com/gin-gonic/gin"
)

type MyMachineController struct {
	common.BaseController
	machineService    *serviceMachine.MachineService
	agentService      *serviceOps.AgentService
	allocationService *serviceAllocation.AllocationService
}

func NewMyMachineController(ms *serviceMachine.MachineService, as *serviceOps.AgentService, alloc *serviceAllocation.AllocationService) *MyMachineController {
	return &MyMachineController{
		machineService:    ms,
		agentService:      as,
		allocationService: alloc,
	}
}

// List 获取当前用户的机器列表
// @Summary 获取我的机器列表
// @Description 根据当前登录用户获取其已分配的机器列表，支持分页
// @Tags Customer - Machines
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/machines [get]
//
// @author Claude
// @description 根据JWT中的userID过滤，只返回分配给当前用户的机器
// @reason 原实现返回所有机器，存在数据泄露风险，现改为按用户过滤
// @modified 2026-02-04
func (c *MyMachineController) List(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	allocations, total, err := c.allocationService.ListByCustomerID(ctx, userID, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "获取机器列表失败")
		return
	}

	// 提取机器信息
	machines := make([]any, 0, len(allocations))
	for _, alloc := range allocations {
		machines = append(machines, map[string]any{
			"id":                alloc.Host.ID,
			"hostname":          alloc.Host.Hostname,
			"ip_address":        alloc.Host.IPAddress,
			"public_ip":         alloc.Host.PublicIP,
			"status":            alloc.Host.Status,
			"device_status":     alloc.Host.DeviceStatus,
			"allocation_status": alloc.Host.AllocationStatus,
			"total_cpu":         alloc.Host.TotalCPU,
			"total_memory_gb":   alloc.Host.TotalMemoryGB,
			"start_time":        alloc.StartTime,
			"end_time":          alloc.EndTime,
		})
	}

	c.Success(ctx, gin.H{
		"list":      machines,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetConnection 获取机器连接信息
// @Summary 获取机器连接信息
// @Description 获取指定机器的 SSH/Jupyter/VNC 等连接信息
// @Tags Customer - Machines
// @Produce json
// @Param id path string true "机器 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/machines/{id}/connection [get]
func (c *MyMachineController) GetConnection(ctx *gin.Context) {
	hostID := ctx.Param("id")

	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}
	if err := c.allocationService.ValidateHostOwnership(ctx, hostID, userID); err != nil {
		c.Error(ctx, 403, "无权访问该机器")
		return
	}

	info, err := c.machineService.GetConnectionInfo(ctx, hostID)
	if err != nil {
		c.Error(ctx, 500, "Failed to get connection info")
		return
	}
	c.Success(ctx, info)
}

// ResetSSH 重置机器 SSH 配置
// @Summary 重置机器 SSH
// @Description 触发指定机器的 SSH 配置重置
// @Tags Customer - Machines
// @Produce json
// @Param id path string true "机器 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/machines/{id}/ssh-reset [post]
func (c *MyMachineController) ResetSSH(ctx *gin.Context) {
	hostID := ctx.Param("id")

	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}
	if err := c.allocationService.ValidateHostOwnership(ctx, hostID, userID); err != nil {
		c.Error(ctx, 403, "无权访问该机器")
		return
	}

	if err := c.agentService.ResetSSH(ctx, hostID); err != nil {
		c.Error(ctx, 500, "Failed to reset SSH")
		return
	}
	c.Success(ctx, gin.H{"message": "SSH reset triggered"})
}