package v1

import (
	"fmt"
	"net/http"

	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// EnvironmentController 环境控制器
type EnvironmentController struct {
	envService *service.EnvironmentService
}

// NewEnvironmentController 创建环境控制器
func NewEnvironmentController() *EnvironmentController {
	return &EnvironmentController{
		envService: service.NewEnvironmentService(),
	}
}

// Create 创建环境
func (ctrl *EnvironmentController) Create(c *gin.Context) {
	var req service.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	env, err := ctrl.envService.CreateEnvironment(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, env)
}

// GetByID 获取环境详情
func (ctrl *EnvironmentController) GetByID(c *gin.Context) {
	id := c.Param("id")
	env, err := ctrl.envService.GetEnvironment(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "环境不存在")
		return
	}
	response.Success(c, env)
}

// List 列出环境
func (ctrl *EnvironmentController) List(c *gin.Context) {
	customerID := c.GetUint("customer_id")
	var workspaceID *uint
	if wsID := c.Query("workspace_id"); wsID != "" {
		var id uint
		if _, err := fmt.Sscanf(wsID, "%d", &id); err == nil {
			workspaceID = &id
		}
	}

	envs, err := ctrl.envService.ListEnvironments(customerID, workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, envs)
}

// Delete 删除环境
func (ctrl *EnvironmentController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.envService.DeleteEnvironment(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// Start 启动环境
func (ctrl *EnvironmentController) Start(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.envService.StartEnvironment(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "启动成功"})
}

// Stop 停止环境
func (ctrl *EnvironmentController) Stop(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.envService.StopEnvironment(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "停止成功"})
}
