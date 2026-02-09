package environment

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceEnvironment "github.com/YoungBoyGod/remotegpu/internal/service/environment"
	"github.com/gin-gonic/gin"
)

// EnvironmentController 环境管理控制器
type EnvironmentController struct {
	common.BaseController
	environmentService *serviceEnvironment.EnvironmentService
}

func NewEnvironmentController(svc *serviceEnvironment.EnvironmentService) *EnvironmentController {
	return &EnvironmentController{environmentService: svc}
}

// getCustomerID 从上下文获取当前登录用户 ID
func (c *EnvironmentController) getCustomerID(ctx *gin.Context) (uint, bool) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "未认证")
		return 0, false
	}
	return userID, true
}

// Create 创建环境
// @Summary 创建开发环境
// @Description 用户创建新的容器化开发环境
// @Tags Customer - Environments
// @Accept json
// @Produce json
// @Param request body v1.CreateEnvironmentRequest true "创建环境请求"
// @Security Bearer
// @Success 200 {object} entity.Environment
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/environments [post]
func (c *EnvironmentController) Create(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	var req apiV1.CreateEnvironmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	env := &entity.Environment{
		UserID:      customerID,
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		CPU:         req.CPU,
		Memory:      req.Memory,
		GPU:         req.GPU,
		Storage:     func() int64 { if req.Storage != nil { return *req.Storage }; return 0 }(),
		Status:      "creating",
	}

	if err := c.environmentService.Create(ctx, env); err != nil {
		c.Error(ctx, 500, "创建环境失败: "+err.Error())
		return
	}
	c.Success(ctx, env)
}

// List 环境列表
// @Summary 获取环境列表
// @Description 获取当前用户的开发环境列表，支持按工作空间筛选
// @Tags Customer - Environments
// @Produce json
// @Param workspace_id query int false "工作空间 ID"
// @Security Bearer
// @Success 200 {array} entity.Environment
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/environments [get]
func (c *EnvironmentController) List(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	var workspaceID *uint
	if wsStr := ctx.Query("workspace_id"); wsStr != "" {
		id, err := strconv.ParseUint(wsStr, 10, 64)
		if err != nil {
			c.Error(ctx, 400, "无效的工作空间 ID")
			return
		}
		wsID := uint(id)
		workspaceID = &wsID
	}

	envs, err := c.environmentService.List(ctx, customerID, workspaceID)
	if err != nil {
		c.Error(ctx, 500, "获取环境列表失败")
		return
	}
	c.Success(ctx, envs)
}

// Detail 环境详情
// @Summary 获取环境详情
// @Description 根据环境 ID 获取环境详细信息
// @Tags Customer - Environments
// @Produce json
// @Param id path string true "环境 ID"
// @Security Bearer
// @Success 200 {object} entity.Environment
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /customer/environments/{id} [get]
func (c *EnvironmentController) Detail(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	env, err := c.environmentService.GetByID(ctx, ctx.Param("id"))
	if err != nil {
		c.Error(ctx, 404, "环境不存在")
		return
	}
	if env.UserID != customerID {
		c.Error(ctx, 403, "无权查看该环境")
		return
	}
	c.Success(ctx, env)
}

// Start 启动环境
// @Summary 启动环境
// @Description 启动指定的已停止环境
// @Tags Customer - Environments
// @Produce json
// @Param id path string true "环境 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /customer/environments/{id}/start [post]
func (c *EnvironmentController) Start(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	if err := c.environmentService.Start(ctx, ctx.Param("id"), customerID); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "running"})
}

// Stop 停止环境
// @Summary 停止环境
// @Description 停止指定的运行中环境
// @Tags Customer - Environments
// @Produce json
// @Param id path string true "环境 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /customer/environments/{id}/stop [post]
func (c *EnvironmentController) Stop(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	if err := c.environmentService.Stop(ctx, ctx.Param("id"), customerID); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "stopped"})
}

// Delete 删除环境
// @Summary 删除环境
// @Description 删除指定的已停止或异常环境
// @Tags Customer - Environments
// @Produce json
// @Param id path string true "环境 ID"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /customer/environments/{id} [delete]
func (c *EnvironmentController) Delete(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	if err := c.environmentService.Delete(ctx, ctx.Param("id"), customerID); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	c.Success(ctx, nil)
}

// AccessInfo 获取环境访问信息
// @Summary 获取环境访问信息
// @Description 获取指定运行中环境的 SSH/Jupyter/VNC 访问信息
// @Tags Customer - Environments
// @Produce json
// @Param id path string true "环境 ID"
// @Security Bearer
// @Success 200 {object} environment.AccessInfo
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /customer/environments/{id}/access [get]
func (c *EnvironmentController) AccessInfo(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	info, err := c.environmentService.GetAccessInfo(ctx, ctx.Param("id"), customerID)
	if err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	c.Success(ctx, info)
}
