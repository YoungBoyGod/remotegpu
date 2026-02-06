package system_config

import (
	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceConfig "github.com/YoungBoyGod/remotegpu/internal/service/system_config"
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

// GetConfigs 获取所有系统配置
func (c *SystemConfigController) GetConfigs(ctx *gin.Context) {
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

	if err := c.configService.UpdateConfigs(ctx, req.Configs); err != nil {
		c.Error(ctx, 500, "更新配置失败")
		return
	}
	c.Success(ctx, nil)
}
