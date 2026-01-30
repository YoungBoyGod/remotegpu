package v1

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// ResourceQuotaController 资源配额控制器
type ResourceQuotaController struct {
	quotaService *service.ResourceQuotaService
}

// NewResourceQuotaController 创建资源配额控制器
func NewResourceQuotaController() *ResourceQuotaController {
	return &ResourceQuotaController{
		quotaService: service.NewResourceQuotaService(),
	}
}

// SetQuotaRequest 设置配额请求
type SetQuotaRequest struct {
	CustomerID  uint  `json:"customer_id" binding:"required"`
	WorkspaceID *uint `json:"workspace_id"`
	CPU         int   `json:"cpu" binding:"required,min=0"`
	Memory      int64 `json:"memory" binding:"required,min=0"`
	GPU         int   `json:"gpu" binding:"required,min=0"`
	Storage     int64 `json:"storage" binding:"required,min=0"`
}

// UpdateQuotaRequest 更新配额请求
type UpdateQuotaRequest struct {
	CPU     int   `json:"cpu" binding:"required,min=0"`
	Memory  int64 `json:"memory" binding:"required,min=0"`
	GPU     int   `json:"gpu" binding:"required,min=0"`
	Storage int64 `json:"storage" binding:"required,min=0"`
}

// SetQuota 设置配额
func (ctrl *ResourceQuotaController) SetQuota(c *gin.Context) {
	var req SetQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	quota := &entity.ResourceQuota{
		CustomerID:  req.CustomerID,
		WorkspaceID: req.WorkspaceID,
		CPU:         req.CPU,
		Memory:      req.Memory,
		GPU:         req.GPU,
		Storage:     req.Storage,
	}

	if err := ctrl.quotaService.SetQuota(quota); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, quota)
}

// GetQuota 获取配额
func (ctrl *ResourceQuotaController) GetQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配额ID")
		return
	}

	quota, err := ctrl.quotaService.GetQuotaByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "配额不存在")
		return
	}

	response.Success(c, quota)
}

// UpdateQuota 更新配额
func (ctrl *ResourceQuotaController) UpdateQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配额ID")
		return
	}

	var req UpdateQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	quota := &entity.ResourceQuota{
		ID:      uint(id),
		CPU:     req.CPU,
		Memory:  req.Memory,
		GPU:     req.GPU,
		Storage: req.Storage,
	}

	if err := ctrl.quotaService.UpdateQuota(quota); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, quota)
}

// DeleteQuota 删除配额
func (ctrl *ResourceQuotaController) DeleteQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配额ID")
		return
	}

	if err := ctrl.quotaService.DeleteQuota(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUsage 获取资源使用情况
func (ctrl *ResourceQuotaController) GetUsage(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Query("customer_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的客户ID")
		return
	}

	var workspaceID *uint
	if wsID := c.Query("workspace_id"); wsID != "" {
		id, err := strconv.ParseUint(wsID, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "无效的工作空间ID")
			return
		}
		wsIDUint := uint(id)
		workspaceID = &wsIDUint
	}

	// 获取配额信息
	quota, err := ctrl.quotaService.GetQuota(uint(customerID), workspaceID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "配额不存在")
		return
	}

	// 获取已使用资源
	used, err := ctrl.quotaService.GetUsedResources(uint(customerID), workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资源使用情况失败: "+err.Error())
		return
	}

	// 获取可用配额
	available, err := ctrl.quotaService.GetAvailableQuota(uint(customerID), workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取可用配额失败: "+err.Error())
		return
	}

	// 构造响应
	usageResponse := map[string]interface{}{
		"quota":     quota,
		"used":      used,
		"available": available,
	}

	response.Success(c, usageResponse)
}
