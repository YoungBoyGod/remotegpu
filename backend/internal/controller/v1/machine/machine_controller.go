package machine

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceAllocation "github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/gin-gonic/gin"
)

type MachineController struct {
	common.BaseController
	machineService    *serviceMachine.MachineService
	allocationService *serviceAllocation.AllocationService
}

func NewMachineController(ms *serviceMachine.MachineService, as *serviceAllocation.AllocationService) *MachineController {
	return &MachineController{
		machineService:    ms,
		allocationService: as,
	}
}

// List 获取机器列表
// @Summary 获取机器列表
// @Description 获取所有机器的列表，支持分页和筛选
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "状态筛选 (idle, allocated, maintenance, offline)"
// @Param region query string false "区域筛选"
// @Param gpu_model query string false "GPU型号筛选"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines [get]
func (c *MachineController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}
	if region := ctx.Query("region"); region != "" {
		filters["region"] = region
	}
	if gpuModel := ctx.Query("gpu_model"); gpuModel != "" {
		filters["gpu_model"] = gpuModel
	}

	machines, total, err := c.machineService.ListMachines(ctx, page, pageSize, filters)
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

// Create 创建机器
// @Summary 创建机器
// @Description 添加新的机器到系统
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body entity.Host true "机器信息"
// @Security Bearer
// @Success 200 {object} entity.Host
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines [post]
func (c *MachineController) Create(ctx *gin.Context) {
	var host entity.Host
	if err := ctx.ShouldBindJSON(&host); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.machineService.CreateMachine(ctx, &host); err != nil {
		if err == serviceMachine.ErrHostDuplicateIP || err == serviceMachine.ErrHostDuplicateHostname {
			c.Error(ctx, 409, "Host already exists")
			return
		}
		c.Error(ctx, 500, "Failed to create machine")
		return
	}

	c.Success(ctx, host)
}

// Allocate 分配机器
// @Summary 分配机器
// @Description 将机器分配给客户
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器ID"
// @Param request body v1.AllocateRequest true "分配请求"
// @Security Bearer
// @Success 200 {object} entity.Allocation
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/allocate [post]
func (c *MachineController) Allocate(ctx *gin.Context) {
	// machineID from URL param
	hostID := ctx.Param("id")

	var req apiV1.AllocateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	// Override HostID from URL if needed, or validate consistency
	if req.HostID != "" && req.HostID != hostID {
		c.Error(ctx, 400, "Host ID mismatch")
		return
	}

	alloc, err := c.allocationService.AllocateMachine(ctx, req.CustomerID, hostID, req.DurationMonths, req.Remark)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, alloc)
}

// Reclaim 回收机器
// @Summary 回收机器
// @Description 从客户处回收机器
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param id path string true "机器ID"
// @Param request body v1.ReclaimRequest false "回收请求"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/{id}/reclaim [post]
func (c *MachineController) Reclaim(ctx *gin.Context) {
	hostID := ctx.Param("id")

	// Optional: bind body for reason
	var req apiV1.ReclaimRequest
	ctx.ShouldBindJSON(&req)

	if err := c.allocationService.ReclaimMachine(ctx, hostID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, gin.H{"message": "Reclaim process started"})
}

// Import 批量导入机器
// @Summary 批量导入机器
// @Description 批量导入机器信息
// @Tags Admin - Machines
// @Accept json
// @Produce json
// @Param request body v1.ImportMachineRequest true "导入请求"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/machines/import [post]
func (c *MachineController) Import(ctx *gin.Context) {
	var req apiV1.ImportMachineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid import data")
		return
	}

	var hosts []entity.Host
	for _, m := range req.Machines {
		hosts = append(hosts, entity.Host{
			IPAddress:   m.HostIP,
			SSHPort:     m.SSHPort,
			Region:      m.Region,
			GPUModel:    m.GPUModel,
			GPUCount:    m.GPUCount,
			CPUCores:    m.CPUCores,
			RAMSize:     m.RAMSize,
			DiskSize:    m.DiskSize,
			Status:      "idle", // Default status
			PriceHourly: float64(m.PriceHourly) / 100.0,
		})
	}

	if err := c.machineService.ImportMachines(ctx, hosts); err != nil {
		c.Error(ctx, 500, "Failed to import machines")
		return
	}

	c.Success(ctx, gin.H{
		"message": "Imported successfully",
		"count":   len(hosts),
	})
}
