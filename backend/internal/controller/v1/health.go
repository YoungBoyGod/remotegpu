package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/pkg/health"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct {
	manager *health.Manager
}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	manager := health.InitManager(config.GlobalConfig)
	return &HealthController{
		manager: manager,
	}
}

// CheckAll 检查所有服务健康状态
func (h *HealthController) CheckAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	results := h.manager.CheckAll(ctx)

	// 统计健康状态
	summary := map[string]int{
		"total":     len(results),
		"healthy":   0,
		"unhealthy": 0,
		"disabled":  0,
		"unknown":   0,
	}

	for _, result := range results {
		switch result.Status {
		case health.StatusHealthy:
			summary["healthy"]++
		case health.StatusUnhealthy:
			summary["unhealthy"]++
		case health.StatusDisabled:
			summary["disabled"]++
		case health.StatusUnknown:
			summary["unknown"]++
		}
	}

	response.Success(c, gin.H{
		"summary": summary,
		"results": results,
	})
}

// CheckService 检查指定服务健康状态
func (h *HealthController) CheckService(c *gin.Context) {
	serviceName := c.Param("service")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	result := h.manager.CheckService(ctx, serviceName)

	if result.Status == health.StatusUnknown {
		response.Error(c, http.StatusNotFound, "服务未找到")
		return
	}

	response.Success(c, result)
}
