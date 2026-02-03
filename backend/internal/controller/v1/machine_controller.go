package v1

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/gin-gonic/gin"
)

type MachineController struct {
	BaseController
	machineService    *service.MachineService
	allocationService *service.AllocationService
}

func NewMachineController(ms *service.MachineService, as *service.AllocationService) *MachineController {
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

type AllocateRequest struct {
	CustomerID     uint   `json:"customer_id" binding:"required"` // Should use CustomerID logic
	HostID         string `json:"host_id" binding:"required"`
	DurationMonths int    `json:"duration_months" binding:"required,min=1"`
	Remark         string `json:"remark"`
}

func (c *MachineController) Allocate(ctx *gin.Context) {
	// machineID from URL param
	hostID := ctx.Param("id")

	var req AllocateRequest
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

type ReclaimRequest struct {
	Reason string `json:"reason"`
	Force  bool   `json:"force"`
}

func (c *MachineController) Reclaim(ctx *gin.Context) {
	hostID := ctx.Param("id")
	
	// Optional: bind body for reason
	var req ReclaimRequest
	ctx.ShouldBindJSON(&req)

	if err := c.allocationService.ReclaimMachine(ctx, hostID); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, gin.H{"message": "Reclaim process started"})
}

func (c *MachineController) Import(ctx *gin.Context) {
	// TODO: Handle Excel file upload and parsing
	c.Success(ctx, gin.H{"message": "Import feature pending implementation"})
}
