package main

import (
	"github.com/YoungBoyGod/remotegpu-agent/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, taskHandler *handler.TaskHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/ping", handlePing)
		api.GET("/system/info", handleSystemInfo)
		api.POST("/process/stop", handleStopProcess)
		api.POST("/ssh/reset", handleResetSSH)
		api.POST("/machine/cleanup", handleCleanup)
		api.POST("/command/exec", handleExecCommand)

		// 任务队列 API
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.POST("/tasks/:id/cancel", taskHandler.CancelTask)
		api.GET("/queue/status", taskHandler.GetQueueStatus)
	}
}
