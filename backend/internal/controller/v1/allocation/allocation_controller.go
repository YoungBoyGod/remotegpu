package allocation

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceAllocation "github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	"github.com/gin-gonic/gin"
)

type AllocationController struct {
	common.BaseController
	allocationService *serviceAllocation.AllocationService
}

func NewAllocationController(as *serviceAllocation.AllocationService) *AllocationController {
	return &AllocationController{
		allocationService: as,
	}
}

// List 获取分配记录列表
// @Summary 获取分配记录列表
// @Description 分页查询分配记录，支持按客户/机器/状态筛选
// @Tags Admin - Allocations
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param customer_id query int false "客户ID筛选"
// @Param host_id query string false "机器ID筛选"
// @Param status query string false "状态筛选 (active, expired, reclaimed)"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/allocations [get]
func (c *AllocationController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if customerID := ctx.Query("customer_id"); customerID != "" {
		if id, err := strconv.ParseUint(customerID, 10, 64); err == nil {
			filters["customer_id"] = uint(id)
		}
	}
	if hostID := ctx.Query("host_id"); hostID != "" {
		filters["host_id"] = hostID
	}
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}

	allocations, total, err := c.allocationService.ListAllocations(ctx, page, pageSize, filters)
	if err != nil {
		c.Error(ctx, 500, "Failed to list allocations")
		return
	}

	c.Success(ctx, gin.H{
		"list":      allocations,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
