package router

import (
	v1 "github.com/YoungBoyGod/remotegpu/internal/controller/v1"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/storage"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	db := database.GetDB()
	// Storage manager should be initialized in main and passed here ideally,
	// or retrieved from a global/package level if using singleton pattern.
	// For now assuming it's available via package function or constructing new (lightweight)
	// In a real app, use dependency injection.
	storageMgr, _ := storage.NewManager(storage.Config{Type: "local", Local: storage.LocalConfig{RootPath: "./uploads"}})

	// --- Services ---
	authService := service.NewAuthService(db)
	machineService := service.NewMachineService(db)
	allocService := service.NewAllocationService(db)
	custService := service.NewCustomerService(db)
	taskService := service.NewTaskService(db)
	datasetService := service.NewDatasetService(db)
	opsService := service.NewOpsService(db)
	agentService := service.NewAgentService()
	monitorService := service.NewMonitorService()
	storageService := service.NewStorageService(storageMgr)
	dashboardService := service.NewDashboardService(machineService, custService, allocService)

	// --- Controllers ---
	authCtrl := v1.NewAuthController(authService)
	dashboardCtrl := v1.NewDashboardController(dashboardService)
	machineCtrl := v1.NewMachineController(machineService, allocService)
	custCtrl := v1.NewCustomerController(custService)
	monitorCtrl := v1.NewMonitorController(monitorService)
	alertCtrl := v1.NewAlertController(opsService)
	
	myMachineCtrl := v1.NewMyMachineController(machineService, agentService)
	taskCtrl := v1.NewTaskController(taskService)
	datasetCtrl := v1.NewDatasetController(datasetService, storageService, agentService)


	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 1. Auth Module
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", authCtrl.Login)
			authGroup.POST("/refresh", authCtrl.Refresh)
			authGroup.POST("/logout", authCtrl.Logout)
			
			// Protected Profile
			authGroup.GET("/profile", middleware.Auth(), authCtrl.GetProfile)
		}

		// 2. Admin Module (Protected + Role Check)
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(middleware.Auth()) // Add middleware.RequireRole("admin")
		{
			// Dashboard
			adminGroup.GET("/dashboard/stats", dashboardCtrl.GetStats)
			adminGroup.GET("/dashboard/gpu-trend", dashboardCtrl.GetGPUTrend)
			adminGroup.GET("/allocations/recent", dashboardCtrl.GetRecentAllocations)

			// Machines
			adminGroup.GET("/machines", machineCtrl.List)
			adminGroup.POST("/machines", machineCtrl.Create)
			adminGroup.POST("/machines/import", machineCtrl.Import)
			adminGroup.POST("/machines/:id/allocate", machineCtrl.Allocate)
			adminGroup.POST("/machines/:id/reclaim", machineCtrl.Reclaim)

			// Customers
			adminGroup.GET("/customers", custCtrl.List)
			adminGroup.POST("/customers", custCtrl.Create)
			adminGroup.POST("/customers/:id/disable", custCtrl.Disable)

			// Monitoring & Ops
			adminGroup.GET("/monitoring/realtime", monitorCtrl.GetRealtime)
			adminGroup.GET("/alerts", alertCtrl.List)
		}

		// 3. Customer Module (Protected)
		custGroup := apiV1.Group("/customer")
		custGroup.Use(middleware.Auth())
		{
			// Machines
			custGroup.GET("/machines", myMachineCtrl.List)
			custGroup.GET("/machines/:id/connection", myMachineCtrl.GetConnection)
			custGroup.POST("/machines/:id/ssh-reset", myMachineCtrl.ResetSSH)

			// Tasks
			custGroup.GET("/tasks", taskCtrl.List)
			custGroup.POST("/tasks/training", taskCtrl.CreateTraining)
			custGroup.POST("/tasks/:id/stop", taskCtrl.Stop)

			// Datasets
			custGroup.GET("/datasets", datasetCtrl.List)
			custGroup.POST("/datasets/init-multipart", datasetCtrl.InitUpload)
			custGroup.POST("/datasets/:id/mount", datasetCtrl.Mount)
		}
	}
}