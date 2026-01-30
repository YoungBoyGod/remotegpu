package v1

import (
	"net/http"
	"strconv"

	v1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// SetQuota 设置配额
func (ctrl *ResourceQuotaController) SetQuota(c *gin.Context) {
	var req v1.SetQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	quota := &entity.ResourceQuota{
		UserID:       req.UserID,
		WorkspaceID:      req.WorkspaceID,
		CPU:              req.MaxCPU,
		Memory:           req.MaxMemory,
		GPU:              req.MaxGPU,
		Storage:          req.MaxStorage,
		EnvironmentQuota: req.MaxEnvironments,
		QuotaLevel:       req.QuotaLevel,
	}

	if quota.QuotaLevel == "" {
		quota.QuotaLevel = "free"
	}

	if err := ctrl.quotaService.SetQuota(quota); err != nil {
		response.Error(c, http.StatusInternalServerError, "设置配额失败: "+err.Error())
		return
	}

	quotaInfo := ctrl.entityToQuotaInfo(quota)
	response.Success(c, quotaInfo)
}

// List 获取配额列表
func (ctrl *ResourceQuotaController) List(c *gin.Context) {
	response.Success(c, v1.QuotaListResponse{
		Items:    []*v1.QuotaInfo{},
		Total:    0,
		Page:     1,
		PageSize: 10,
	})
}

// GetQuota 获取配额详情
func (ctrl *ResourceQuotaController) GetQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配额ID")
		return
	}

	quota, err := ctrl.quotaService.GetQuotaByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusNotFound, "配额不存在")
		} else {
			response.Error(c, http.StatusInternalServerError, "获取配额失败: "+err.Error())
		}
		return
	}

	quotaInfo := ctrl.entityToQuotaInfo(quota)
	response.Success(c, quotaInfo)
}

// UpdateQuota 更新配额
func (ctrl *ResourceQuotaController) UpdateQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配额ID")
		return
	}

	var req v1.UpdateQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	quota := &entity.ResourceQuota{
		ID:               uint(id),
		CPU:              req.MaxCPU,
		Memory:           req.MaxMemory,
		GPU:              req.MaxGPU,
		Storage:          req.MaxStorage,
		EnvironmentQuota: req.MaxEnvironments,
		QuotaLevel:       req.QuotaLevel,
	}

	if err := ctrl.quotaService.UpdateQuota(quota); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新配额失败: "+err.Error())
		return
	}

	updatedQuota, err := ctrl.quotaService.GetQuotaByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取更新后的配额失败: "+err.Error())
		return
	}

	quotaInfo := ctrl.entityToQuotaInfo(updatedQuota)
	response.Success(c, quotaInfo)
}

// DeleteQuota 删除配额
func (ctrl *ResourceQuotaController) DeleteQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配额ID")
		return
	}

	if err := ctrl.quotaService.DeleteQuota(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除配额失败: "+err.Error())
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

	quota, err := ctrl.quotaService.GetQuota(uint(customerID), workspaceID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "配额不存在")
		return
	}

	used, err := ctrl.quotaService.GetUsedResources(uint(customerID), workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资源使用情况失败: "+err.Error())
		return
	}

	available, err := ctrl.quotaService.GetAvailableQuota(uint(customerID), workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取可用配额失败: "+err.Error())
		return
	}

	usageResponse := v1.QuotaUsageResponse{
		Quota: v1.QuotaDetail{
			MaxGPU:          quota.GPU,
			MaxCPU:          quota.CPU,
			MaxMemory:       quota.Memory,
			MaxStorage:      quota.Storage,
			MaxEnvironments: quota.EnvironmentQuota,
		},
		Used: v1.UsedResources{
			UsedGPU:          used.GPU,
			UsedCPU:          used.CPU,
			UsedMemory:       used.Memory,
			UsedStorage:      used.Storage,
			UsedEnvironments: 0, // TODO: 实现环境数量统计
		},
		Available: v1.AvailableResources{
			AvailableGPU:          available.GPU,
			AvailableCPU:          available.CPU,
			AvailableMemory:       available.Memory,
			AvailableStorage:      available.Storage,
			AvailableEnvironments: available.EnvironmentQuota,
		},
		UsagePercentage: v1.UsagePercentageDetail{
			GPU:          calculatePercentage(used.GPU, quota.GPU),
			CPU:          calculatePercentage(used.CPU, quota.CPU),
			Memory:       calculatePercentage(int(used.Memory), int(quota.Memory)),
			Storage:      calculatePercentage(int(used.Storage), int(quota.Storage)),
			Environments: 0, // TODO: 实现环境数量统计
		},
	}

	response.Success(c, usageResponse)
}

// entityToQuotaInfo 将实体转换为API响应格式
func (ctrl *ResourceQuotaController) entityToQuotaInfo(quota *entity.ResourceQuota) *v1.QuotaInfo {
	return &v1.QuotaInfo{
		ID:              quota.ID,
		UserID:      quota.UserID,
		WorkspaceID:     quota.WorkspaceID,
		QuotaLevel:      quota.QuotaLevel,
		MaxGPU:          quota.GPU,
		MaxCPU:          quota.CPU,
		MaxMemory:       quota.Memory,
		MaxStorage:      quota.Storage,
		MaxEnvironments: quota.EnvironmentQuota,
		CreatedAt:       quota.CreatedAt,
		UpdatedAt:       quota.UpdatedAt,
	}
}

// calculatePercentage 计算百分比
func calculatePercentage(used, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * 100
}
