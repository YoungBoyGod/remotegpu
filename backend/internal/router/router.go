package router

import (
	v1 "github.com/YoungBoyGod/remotegpu/internal/controller/v1"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	// 初始化控制器
	userController := v1.NewUserController()
	healthController := v1.NewHealthController()
	hostController := v1.NewHostController()
	gpuController := v1.NewGPUController()
	environmentController := v1.NewEnvironmentController()
	quotaController := v1.NewResourceQuotaController()

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 健康检查
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// 用户路由（公开）
		user := apiV1.Group("/user")
		{
			user.POST("/register", userController.Register)
			user.POST("/login", userController.Login)
			user.GET("/:id", userController.GetUserByID)
		}

		// 需要认证的路由
		auth := apiV1.Group("")
		auth.Use(middleware.Auth())
		{
			// 用户相关（需要认证）
			auth.GET("/user/info", userController.GetUserInfo)
			auth.PUT("/user/info", userController.UpdateUser)

			// 资源配额查询（需要认证）
			auth.GET("/quotas/usage", quotaController.GetUsage)

			// 环境管理（需要认证，普通用户只能操作自己的环境）
			auth.POST("/environments", environmentController.Create)
			auth.GET("/environments", environmentController.List)
			auth.GET("/environments/:id", environmentController.GetByID)
			auth.DELETE("/environments/:id", environmentController.Delete)
			auth.POST("/environments/:id/start", environmentController.Start)
			auth.POST("/environments/:id/stop", environmentController.Stop)
			auth.POST("/environments/:id/restart", environmentController.Restart)
			auth.GET("/environments/:id/access", environmentController.GetAccessInfo)
			auth.GET("/environments/:id/logs", environmentController.GetLogs)
		}

		// 管理员路由（需要管理员权限）
		admin := apiV1.Group("/admin")
		admin.Use(middleware.Auth(), middleware.RequireAdmin())
		{
			// 健康检查
			admin.GET("/health/all", healthController.CheckAll)
			admin.GET("/health/:service", healthController.CheckService)

			// 主机管理
			admin.POST("/hosts", hostController.Create)
			admin.GET("/hosts", hostController.List)
			admin.GET("/hosts/:id", hostController.GetByID)
			admin.PUT("/hosts/:id", hostController.Update)
			admin.DELETE("/hosts/:id", hostController.Delete)
			admin.POST("/hosts/:id/heartbeat", hostController.Heartbeat)

			// GPU管理
			admin.POST("/gpus", gpuController.Create)
			admin.GET("/gpus", gpuController.List)
			admin.GET("/gpus/:id", gpuController.GetByID)
			admin.PUT("/gpus/:id", gpuController.Update)
			admin.DELETE("/gpus/:id", gpuController.Delete)
			admin.POST("/gpus/:id/allocate", gpuController.Allocate)
			admin.POST("/gpus/:id/release", gpuController.Release)
			admin.GET("/hosts/:id/gpus", gpuController.GetByHostID)

			// 环境管理
			admin.POST("/environments", environmentController.Create)
			admin.GET("/environments", environmentController.List)
			admin.GET("/environments/:id", environmentController.GetByID)
			admin.DELETE("/environments/:id", environmentController.Delete)
			admin.POST("/environments/:id/start", environmentController.Start)
			admin.POST("/environments/:id/stop", environmentController.Stop)
			admin.POST("/environments/:id/restart", environmentController.Restart)
			admin.GET("/environments/:id/access", environmentController.GetAccessInfo)
			admin.GET("/environments/:id/logs", environmentController.GetLogs)

			// 资源配额管理
			admin.POST("/quotas", quotaController.SetQuota)
			admin.GET("/quotas", quotaController.List)
			admin.GET("/quotas/:id", quotaController.GetQuota)
			admin.PUT("/quotas/:id", quotaController.UpdateQuota)
			admin.DELETE("/quotas/:id", quotaController.DeleteQuota)
		}
	}
}
