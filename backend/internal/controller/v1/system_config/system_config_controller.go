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
func (c *SystemConfigController) ListGroups(ctx *gin.Context) {
	groups, err := c.configService.ListGroups(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取配置分组失败")
		return
	}
	c.Success(ctx, groups)
}

// GetConfig 获取单条配置
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
