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

func (c *MachineController) Create(ctx *gin.Context) {
	var host entity.Host
	if err := ctx.ShouldBindJSON(&host); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.machineService.CreateMachine(ctx, &host); err != nil {
		c.Error(ctx, 500, "Failed to create machine")
		return
	}

	c.Success(ctx, host)
}

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
