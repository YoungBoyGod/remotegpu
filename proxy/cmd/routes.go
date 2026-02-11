package main

import (
	"github.com/YoungBoyGod/remotegpu-proxy/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, h *handler.MappingHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/ping", h.Ping)
		api.GET("/stats", h.GetStats)

		// 映射管理
		api.POST("/mappings", h.AddMapping)
		api.GET("/mappings", h.ListMappings)
		api.DELETE("/mappings/:port", h.RemoveMapping)
		api.DELETE("/mappings/env/:id", h.RemoveByEnvID)
	}
}
