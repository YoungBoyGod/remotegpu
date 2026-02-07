package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/prometheus"
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
	serviceSystemConfig "github.com/YoungBoyGod/remotegpu/internal/service/system_config"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"

	// 控制器层
	ctrlAuth "github.com/YoungBoyGod/remotegpu/internal/controller/v1/auth"
	ctrlCustomer "github.com/YoungBoyGod/remotegpu/internal/controller/v1/customer"
	ctrlDataset "github.com/YoungBoyGod/remotegpu/internal/controller/v1/dataset"
	ctrlMachine "github.com/YoungBoyGod/remotegpu/internal/controller/v1/machine"
	ctrlOps "github.com/YoungBoyGod/remotegpu/internal/controller/v1/ops"
	ctrlSystemConfig "github.com/YoungBoyGod/remotegpu/internal/controller/v1/system_config"
	ctrlTask "github.com/YoungBoyGod/remotegpu/internal/controller/v1/task"

	ctrlAgent "github.com/YoungBoyGod/remotegpu/internal/controller/v1/agent"
)

// agentAdapter 适配器，将 AgentService 转换为 AgentSystemInfoProvider 接口
type agentAdapter struct {
	svc *serviceOps.AgentService
}

func (a *agentAdapter) GetSystemInfo(ctx context.Context, hostID, address string) (*serviceMachine.SystemInfoSnapshot, error) {
	info, err := a.svc.GetSystemInfo(ctx, hostID, address)
	if err != nil {
		return nil, err
	}
	return &serviceMachine.SystemInfoSnapshot{
		Hostname:      info.Hostname,
		CPUCores:      info.CPUCores,
		MemoryTotalGB: info.MemoryTotalGB,
		DiskTotalGB:   info.DiskTotalGB,
		Collected:     info.Collected,
	}, nil
}

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
	auditSvc := serviceAudit.NewAuditService(db)
	agentSvc := serviceOps.NewAgentService(db, &config.GlobalConfig.Agent)
	allocSvc := serviceAllocation.NewAllocationService(db, auditSvc, agentSvc)
	allocSvc.StartWorker(context.Background())
	custSvc := serviceCustomer.NewCustomerService(db)
	opsSvc := serviceOps.NewOpsService(db)
	taskSvc := serviceTask.NewTaskService(db, agentSvc)
	datasetSvc := serviceDataset.NewDatasetService(db)

	// Prometheus 客户端
	promClient := prometheus.NewClient(&prometheus.Config{
		Enabled:  config.GlobalConfig.Prometheus.Enabled,
		Endpoint: config.GlobalConfig.Prometheus.Endpoint,
	})

	monitorSvc := serviceOps.NewMonitorService(machineSvc, cache.GetCache(), promClient)
	storageSvc := serviceStorage.NewStorageService(storageMgr)
	sshKeySvc := serviceSSHKey.NewSSHKeyService(db)
	imageSvc := serviceImage.NewImageService(db)
	enrollmentSvc := serviceMachine.NewMachineEnrollmentService(db, machineSvc, &agentAdapter{svc: agentSvc})
	enrollmentSvc.StartWorker(context.Background())

	// 启动心跳监控服务
	if config.GlobalConfig.HeartbeatMonitor.Enabled {
		heartbeatMonitor := serviceMachine.NewHeartbeatMonitor(
			db,
			time.Duration(config.GlobalConfig.HeartbeatMonitor.Timeout)*time.Second,
			time.Duration(config.GlobalConfig.HeartbeatMonitor.CheckInterval)*time.Second,
		)
		go heartbeatMonitor.Start(context.Background())
	}

	// 启动监控数据采集服务
	if config.GlobalConfig.MetricsCollector.Enabled {
		metricsCollector := serviceMachine.NewMetricsCollector(
			db,
			time.Duration(config.GlobalConfig.MetricsCollector.Interval)*time.Second,
			config.GlobalConfig.MetricsCollector.RetentionDays,
		)
		go metricsCollector.Start(context.Background())
	}

	systemConfigSvc := serviceSystemConfig.NewSystemConfigService(db)
	dashboardSvc := serviceOps.NewDashboardService(machineSvc, custSvc, allocSvc, promClient)

	// --- 控制器层初始化 ---
	authController := ctrlAuth.NewAuthController(authSvc)
	dashboardController := ctrlOps.NewDashboardController(dashboardSvc)
	machineController := ctrlMachine.NewMachineController(machineSvc, allocSvc, agentSvc)
	customerController := ctrlCustomer.NewCustomerController(custSvc)
	monitorController := ctrlOps.NewMonitorController(monitorSvc)
	alertController := ctrlOps.NewAlertController(opsSvc)

	myMachineController := ctrlCustomer.NewMyMachineController(machineSvc, agentSvc, allocSvc)
	taskController := ctrlTask.NewTaskController(taskSvc)
	agentTaskController := ctrlTask.NewAgentTaskController(taskSvc)
	agentHeartbeatController := ctrlAgent.NewHeartbeatController(machineSvc)
	datasetController := ctrlDataset.NewDatasetController(datasetSvc, storageSvc, agentSvc, allocSvc)
	sshKeyController := ctrlCustomer.NewSSHKeyController(sshKeySvc)
	enrollmentController := ctrlCustomer.NewMachineEnrollmentController(enrollmentSvc)
	auditController := ctrlOps.NewAuditController(auditSvc)
	imageController := ctrlOps.NewImageController(imageSvc)
	systemConfigController := ctrlSystemConfig.NewSystemConfigController(systemConfigSvc)

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
			authGroup.POST("/password/change", middleware.Auth(db), authController.ChangePassword)
		}

		// 2. Admin Module (Protected + Role Check)
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(middleware.Auth(db), middleware.RequireAdmin(), middleware.AuditMiddleware(auditSvc))
		{
			// 仪表板
			adminGroup.GET("/dashboard/stats", dashboardController.GetStats)
			adminGroup.GET("/dashboard/gpu-trend", dashboardController.GetGPUTrend)
			adminGroup.GET("/allocations/recent", dashboardController.GetRecentAllocations)

			// 机器管理
			adminGroup.GET("/machines", machineController.List)
			adminGroup.GET("/machines/:id", machineController.Detail)
			adminGroup.POST("/machines", machineController.Create)
			adminGroup.POST("/machines/import", machineController.Import)
			adminGroup.DELETE("/machines/:id", machineController.Delete)
			adminGroup.POST("/machines/:id/collect", machineController.CollectSpec)
			adminGroup.POST("/machines/:id/allocate", machineController.Allocate)
			adminGroup.POST("/machines/:id/reclaim", machineController.Reclaim)
			adminGroup.POST("/machines/:id/maintenance", machineController.SetMaintenance)

			// 客户管理
			adminGroup.GET("/customers", customerController.List)
			adminGroup.POST("/customers", customerController.Create)
			adminGroup.POST("/customers/:id/disable", customerController.Disable)

			// 监控与运维
			adminGroup.GET("/monitoring/realtime", monitorController.GetRealtime)
			adminGroup.GET("/alerts", alertController.List)
			adminGroup.POST("/alerts/:id/acknowledge", alertController.Acknowledge)

			// 审计日志
			adminGroup.GET("/audit/logs", auditController.List)

			// 镜像管理
			adminGroup.GET("/images", imageController.List)
			adminGroup.POST("/images/sync", imageController.Sync)

			// 系统配置
			adminGroup.GET("/settings/configs", systemConfigController.GetConfigs)
			adminGroup.PUT("/settings/configs", systemConfigController.UpdateConfigs)
		}

		// 3. Customer Module (Protected)
		custGroup := apiV1.Group("/customer")
		custGroup.Use(middleware.Auth(db))
		{
			// 机器管理
			custGroup.GET("/machines", myMachineController.List)
			custGroup.POST("/machines", enrollmentController.Create)
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

			// 用户添加机器
			custGroup.POST("/machines/enroll", enrollmentController.Create)
			custGroup.GET("/machines/enrollments", enrollmentController.List)
			custGroup.GET("/machines/enrollments/:id", enrollmentController.Detail)
		}

		// 4. Agent Module (Agent 专用 API)
		agentGroup := apiV1.Group("/agent")
		{
			agentGroup.POST("/heartbeat", agentHeartbeatController.Heartbeat)
			agentGroup.POST("/tasks/claim", agentTaskController.ClaimTasks)
			agentGroup.POST("/tasks/:id/start", agentTaskController.StartTask)
			agentGroup.POST("/tasks/:id/lease/renew", agentTaskController.RenewLease)
			agentGroup.POST("/tasks/:id/complete", agentTaskController.CompleteTask)
		}
	}
}
