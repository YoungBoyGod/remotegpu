package system_config

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceConfig "github.com/YoungBoyGod/remotegpu/internal/service/system_config"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/gin-gonic/gin"
)

type SystemConfigController struct {
	common.BaseController
	configService *serviceConfig.SystemConfigService
}

func NewSystemConfigController(cs *serviceConfig.SystemConfigService) *SystemConfigController {
	return &SystemConfigController{
		configService: cs,
	}
}

// GetConfigs 获取所有系统配置（支持按分组过滤）
// @Summary 获取系统配置列表
// @Description 获取所有系统配置，支持按分组筛选
// @Tags Admin - System Config
// @Produce json
// @Param group query string false "配置分组筛选"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/settings/configs [get]
func (c *SystemConfigController) GetConfigs(ctx *gin.Context) {
	group := ctx.Query("group")
	if group != "" {
		configs, err := c.configService.GetConfigsByGroup(ctx, group)
		if err != nil {
			c.Error(ctx, 500, "获取配置失败")
			return
		}
		c.Success(ctx, configs)
		return
	}

	configs, err := c.configService.GetAllConfigs(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取配置失败")
		return
	}
	c.Success(ctx, configs)
}

// UpdateConfigs 批量更新系统配置
// @Summary 批量更新系统配置
// @Description 批量更新多个系统配置项的值
// @Tags Admin - System Config
// @Accept json
// @Produce json
// @Param request body v1.UpdateSystemConfigsRequest true "批量更新配置请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/settings/configs [put]
func (c *SystemConfigController) UpdateConfigs(ctx *gin.Context) {
	var req apiV1.UpdateSystemConfigsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	operator := ctx.GetString("username")
	if err := c.configService.UpdateConfigs(ctx, req.Configs, operator); err != nil {
		c.Error(ctx, 500, "更新配置失败")
		return
	}
	c.Success(ctx, nil)
}

// ListGroups 获取所有配置分组
// @Summary 获取配置分组列表
// @Description 获取系统配置的所有分组名称
// @Tags Admin - System Config
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/settings/configs/groups [get]
func (c *SystemConfigController) ListGroups(ctx *gin.Context) {
	groups, err := c.configService.ListGroups(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取配置分组失败")
		return
	}
	c.Success(ctx, groups)
}

// GetConfig 获取单条配置
// @Summary 获取单条系统配置
// @Description 根据配置 ID 获取单条系统配置详情
// @Tags Admin - System Config
// @Produce json
// @Param id path int true "配置 ID"
// @Security Bearer
// @Success 200 {object} entity.SystemConfig
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/settings/configs/{id} [get]
func (c *SystemConfigController) GetConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的配置 ID")
		return
	}

	config, err := c.configService.GetConfig(ctx, uint(id))
	if err != nil {
		c.Error(ctx, 404, err.Error())
		return
	}
	c.Success(ctx, config)
}

// CreateConfig 创建配置项
// @Summary 创建系统配置
// @Description 创建新的系统配置项
// @Tags Admin - System Config
// @Accept json
// @Produce json
// @Param request body v1.CreateSystemConfigRequest true "创建配置请求"
// @Security Bearer
// @Success 200 {object} entity.SystemConfig
// @Failure 400 {object} common.ErrorResponse
// @Failure 409 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/settings/configs [post]
func (c *SystemConfigController) CreateConfig(ctx *gin.Context) {
	var req apiV1.CreateSystemConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	config := &entity.SystemConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		ConfigGroup: req.ConfigGroup,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	}

	operator := ctx.GetString("username")
	if err := c.configService.CreateConfig(ctx, config, operator); err != nil {
		if err == serviceConfig.ErrConfigKeyConflict {
			c.Error(ctx, 409, err.Error())
			return
		}
		c.Error(ctx, 500, "创建配置失败")
		return
	}
	c.Success(ctx, config)
}

// UpdateConfig 更新单条配置
// @Summary 更新单条系统配置
// @Description 根据配置 ID 更新系统配置项的值
// @Tags Admin - System Config
// @Accept json
// @Produce json
// @Param id path int true "配置 ID"
// @Param request body v1.UpdateSystemConfigRequest true "更新配置请求"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/settings/configs/{id} [put]
func (c *SystemConfigController) UpdateConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的配置 ID")
		return
	}

	var req apiV1.UpdateSystemConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	fields := make(map[string]interface{})
	if req.ConfigValue != "" {
		fields["config_value"] = req.ConfigValue
	}
	if req.ConfigType != "" {
		fields["config_type"] = req.ConfigType
	}
	if req.ConfigGroup != "" {
		fields["config_group"] = req.ConfigGroup
	}
	if req.Description != "" {
		fields["description"] = req.Description
	}
	if req.IsPublic != nil {
		fields["is_public"] = *req.IsPublic
	}

	if len(fields) == 0 {
		c.Error(ctx, 400, "没有需要更新的字段")
		return
	}

	operator := ctx.GetString("username")
	if err := c.configService.UpdateConfig(ctx, uint(id), fields, operator); err != nil {
		if err == serviceConfig.ErrConfigNotFound {
			c.Error(ctx, 404, err.Error())
			return
		}
		c.Error(ctx, 500, "更新配置失败")
		return
	}
	c.Success(ctx, nil)
}

// DeleteConfig 删除配置项
// @Summary 删除系统配置
// @Description 根据配置 ID 删除系统配置项
// @Tags Admin - System Config
// @Produce json
// @Param id path int true "配置 ID"
// @Security Bearer
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/settings/configs/{id} [delete]
func (c *SystemConfigController) DeleteConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的配置 ID")
		return
	}

	operator := ctx.GetString("username")
	if err := c.configService.DeleteConfig(ctx, uint(id), operator); err != nil {
		if err == serviceConfig.ErrConfigNotFound {
			c.Error(ctx, 404, err.Error())
			return
		}
		c.Error(ctx, 500, "删除配置失败")
		return
	}
	c.Success(ctx, nil)
}
