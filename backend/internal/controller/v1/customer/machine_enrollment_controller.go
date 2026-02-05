package customer

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/gin-gonic/gin"
)

type MachineEnrollmentController struct {
	common.BaseController
	service *serviceMachine.MachineEnrollmentService
}

func NewMachineEnrollmentController(svc *serviceMachine.MachineEnrollmentService) *MachineEnrollmentController {
	return &MachineEnrollmentController{service: svc}
}

// Create 创建用户机器添加任务
func (c *MachineEnrollmentController) Create(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	var req apiV1.CreateMachineEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	if req.IPAddress == "" && req.Hostname == "" {
		c.Error(ctx, 400, "Host address is required")
		return
	}
	if req.SSHPassword == "" && req.SSHKey == "" {
		c.Error(ctx, 400, "SSH key or password required")
		return
	}

	address := req.IPAddress
	if address == "" {
		address = req.Hostname
	}

	enrollment := &entity.MachineEnrollment{
		Name:        req.Name,
		Hostname:    req.Hostname,
		Region:      req.Region,
		Address:     address,
		SSHPort:     req.SSHPort,
		SSHUsername: req.SSHUsername,
		SSHPassword: req.SSHPassword,
		SSHKey:      req.SSHKey,
	}

	result, err := c.service.CreateEnrollment(ctx, userID, enrollment)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}

	c.Success(ctx, result)
}

// List 列出用户机器添加任务
func (c *MachineEnrollmentController) List(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	list, total, err := c.service.ListByCustomer(ctx, userID, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "获取添加任务失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Detail 获取单个任务
func (c *MachineEnrollmentController) Detail(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的任务ID")
		return
	}

	enrollment, err := c.service.GetByID(ctx, uint(ID))
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}
	if enrollment.CustomerID != userID {
		c.Error(ctx, 403, "无权访问该任务")
		return
	}

	c.Success(ctx, enrollment)
}
