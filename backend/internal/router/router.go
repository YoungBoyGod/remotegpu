package router

import (
	"github.com/gin-gonic/gin"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/storage"

	// 服务层
	serviceAllocation "github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	serviceAudit "github.com/YoungBoyGod/remotegpu/internal/service/audit"
	serviceAuth "github.com/YoungBoyGod/remotegpu/internal/service/auth"
	serviceCustomer "github.com/YoungBoyGod/remotegpu/internal/service/customer"
	serviceDataset "github.com/YoungBoyGod/remotegpu/internal/service/dataset"
	serviceImage "github.com/YoungBoyGod/remotegpu/internal/service/image"
	serviceMachine "github.com/YoungBoyGod/remotegpu/internal/service/machine"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	serviceSSHKey "github.com/YoungBoyGod/remotegpu/internal/service/sshkey"
	serviceStorage "github.com/YoungBoyGod/remotegpu/internal/service/storage"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"

	// 控制器层
	ctrlAuth "github.com/YoungBoyGod/remotegpu/internal/controller/v1/auth"
	ctrlCustomer "github.com/YoungBoyGod/remotegpu/internal/controller/v1/customer"
	ctrlDataset "github.com/YoungBoyGod/remotegpu/internal/controller/v1/dataset"
	ctrlMachine "github.com/YoungBoyGod/remotegpu/internal/controller/v1/machine"
	ctrlOps "github.com/YoungBoyGod/remotegpu/internal/controller/v1/ops"
	ctrlTask "github.com/YoungBoyGod/remotegpu/internal/controller/v1/task"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	db := database.GetDB()

	// 设置Swagger文档
	if err := middleware.SetupSwagger(r); err != nil {
		panic(err)
	}

	// 存储设置
	storageMgr, _ := storage.NewManager(config.GlobalConfig.Storage)

	// 本地存储的静态文件服务（仅开发环境）
	if config.GlobalConfig.Server.Mode == "debug" {
		// 假设 "local-main" 是默认或其中一个后端
		// 更健壮的方式是遍历后端并找到本地的
		for _, backend := range config.GlobalConfig.Storage.Backends {
			if backend.Type == "local" && backend.Enabled {
				// 在 /uploads 路径下提供服务，例如 http://localhost:8080/uploads/filename.jpg
				r.Static("/uploads", backend.Path)
				break 
			}
		}
	}

	// --- 服务层初始化 ---
	authSvc := serviceAuth.NewAuthService(db, cache.GetCache())
	machineSvc := serviceMachine.NewMachineService(db)
	allocSvc := serviceAllocation.NewAllocationService(db)
	custSvc := serviceCustomer.NewCustomerService(db)
	taskSvc := serviceTask.NewTaskService(db)
	datasetSvc := serviceDataset.NewDatasetService(db)
	opsSvc := serviceOps.NewOpsService(db)
	agentSvc := serviceOps.NewAgentService()
	monitorSvc := serviceOps.NewMonitorService()
	storageSvc := serviceStorage.NewStorageService(storageMgr)
	sshKeySvc := serviceSSHKey.NewSSHKeyService(db)
	auditSvc := serviceAudit.NewAuditService(db)
	imageSvc := serviceImage.NewImageService(db)
	
	dashboardSvc := serviceOps.NewDashboardService(machineSvc, custSvc, allocSvc)

	// --- 控制器层初始化 ---
	authController := ctrlAuth.NewAuthController(authSvc)
	dashboardController := ctrlOps.NewDashboardController(dashboardSvc)
	machineController := ctrlMachine.NewMachineController(machineSvc, allocSvc)
	customerController := ctrlCustomer.NewCustomerController(custSvc)
	monitorController := ctrlOps.NewMonitorController(monitorSvc)
	alertController := ctrlOps.NewAlertController(opsSvc)
	
	myMachineController := ctrlCustomer.NewMyMachineController(machineSvc, agentSvc, allocSvc)
	taskController := ctrlTask.NewTaskController(taskSvc)
	datasetController := ctrlDataset.NewDatasetController(datasetSvc, storageSvc, agentSvc)
	sshKeyController := ctrlCustomer.NewSSHKeyController(sshKeySvc)
	auditController := ctrlOps.NewAuditController(auditSvc)
	imageController := ctrlOps.NewImageController(imageSvc)

	// API v1 路由
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 1. Auth Module
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/refresh", authController.Refresh)
			authGroup.POST("/logout", authController.Logout)
			
			// 受保护的个人资料
			authGroup.GET("/profile", middleware.Auth(db), authController.GetProfile)
		}

		// 2. Admin Module (Protected + Role Check)
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(middleware.Auth(db), middleware.RequireAdmin())
		{
			// 仪表板
			adminGroup.GET("/dashboard/stats", dashboardController.GetStats)
			adminGroup.GET("/dashboard/gpu-trend", dashboardController.GetGPUTrend)
			adminGroup.GET("/allocations/recent", dashboardController.GetRecentAllocations)

			// 机器管理
			adminGroup.GET("/machines", machineController.List)
			adminGroup.POST("/machines", machineController.Create)
			adminGroup.POST("/machines/import", machineController.Import)
			adminGroup.POST("/machines/:id/allocate", machineController.Allocate)
			adminGroup.POST("/machines/:id/reclaim", machineController.Reclaim)

			// 客户管理
			adminGroup.GET("/customers", customerController.List)
			adminGroup.POST("/customers", customerController.Create)
			adminGroup.POST("/customers/:id/disable", customerController.Disable)

			// 监控与运维
			adminGroup.GET("/monitoring/realtime", monitorController.GetRealtime)
			adminGroup.GET("/alerts", alertController.List)

			// 审计日志
			adminGroup.GET("/audit/logs", auditController.List)

			// 镜像管理
			adminGroup.GET("/images", imageController.List)
		}

		// 3. Customer Module (Protected)
		custGroup := apiV1.Group("/customer")
		custGroup.Use(middleware.Auth(db))
		{
			// 机器管理
			custGroup.GET("/machines", myMachineController.List)
			custGroup.GET("/machines/:id/connection", myMachineController.GetConnection)
			custGroup.POST("/machines/:id/ssh-reset", myMachineController.ResetSSH)

			// 任务管理
			custGroup.GET("/tasks", taskController.List)
			custGroup.POST("/tasks/training", taskController.CreateTraining)
			custGroup.POST("/tasks/:id/stop", taskController.Stop)

			// 数据集管理
			custGroup.GET("/datasets", datasetController.List)
			custGroup.POST("/datasets/init-multipart", datasetController.InitUpload)
			custGroup.POST("/datasets/:id/mount", datasetController.Mount)

			// SSH 密钥管理
			custGroup.GET("/keys", sshKeyController.List)
			custGroup.POST("/keys", sshKeyController.Create)
			custGroup.DELETE("/keys/:id", sshKeyController.Delete)
		}
	}
}