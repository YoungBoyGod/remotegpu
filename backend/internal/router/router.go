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
		}

		// 管理员路由（需要管理员权限）
		admin := apiV1.Group("/admin")
		admin.Use(middleware.Auth(), middleware.RequireAdmin())
		{
			// 健康检查
			admin.GET("/health/all", healthController.CheckAll)
			admin.GET("/health/:service", healthController.CheckService)
		}
	}
}
