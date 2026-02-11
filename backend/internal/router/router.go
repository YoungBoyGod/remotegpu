package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	serviceDocument "github.com/YoungBoyGod/remotegpu/internal/service/document"
	serviceNotification "github.com/YoungBoyGod/remotegpu/internal/service/notification"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"
	serviceWorkspace "github.com/YoungBoyGod/remotegpu/internal/service/workspace"
	serviceEnvironment "github.com/YoungBoyGod/remotegpu/internal/service/environment"
	serviceProxy "github.com/YoungBoyGod/remotegpu/internal/service/proxy"

	// 控制器层
	ctrlAuth "github.com/YoungBoyGod/remotegpu/internal/controller/v1/auth"
	ctrlCustomer "github.com/YoungBoyGod/remotegpu/internal/controller/v1/customer"
	ctrlDataset "github.com/YoungBoyGod/remotegpu/internal/controller/v1/dataset"
	ctrlMachine "github.com/YoungBoyGod/remotegpu/internal/controller/v1/machine"
	ctrlOps "github.com/YoungBoyGod/remotegpu/internal/controller/v1/ops"
	ctrlSystemConfig "github.com/YoungBoyGod/remotegpu/internal/controller/v1/system_config"
	ctrlTask "github.com/YoungBoyGod/remotegpu/internal/controller/v1/task"

	ctrlAgent "github.com/YoungBoyGod/remotegpu/internal/controller/v1/agent"
	ctrlDocument "github.com/YoungBoyGod/remotegpu/internal/controller/v1/document"
	ctrlNotification "github.com/YoungBoyGod/remotegpu/internal/controller/v1/notification"
	ctrlStorage "github.com/YoungBoyGod/remotegpu/internal/controller/v1/storage"
	ctrlWorkspace "github.com/YoungBoyGod/remotegpu/internal/controller/v1/workspace"
	ctrlEnvironment "github.com/YoungBoyGod/remotegpu/internal/controller/v1/environment"
	ctrlAllocation "github.com/YoungBoyGod/remotegpu/internal/controller/v1/allocation"
	ctrlProxy "github.com/YoungBoyGod/remotegpu/internal/controller/v1/proxy"
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

	// Prometheus metrics 中间件和端点
	r.Use(middleware.PrometheusMetrics())
	r.GET("/metrics", middleware.Auth(db), middleware.RequireAdmin(), gin.WrapH(promhttp.Handler()))

	// 注册数据库连接池指标
	if sqlDB, err := db.DB(); err == nil {
		middleware.RegisterDBMetrics(sqlDB)
	}

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
	// 注入 Redis 设备状态缓存
	// 使用心跳监控配置的超时值作为 Redis TTL，确保超时判断一致
	heartbeatTimeout := time.Duration(config.GlobalConfig.HeartbeatMonitor.Timeout) * time.Second
	if heartbeatTimeout == 0 {
		heartbeatTimeout = 180 * time.Second // 默认 180 秒
	}
	if c := cache.GetCache(); c != nil {
		hostStatusCache := serviceMachine.NewHostStatusCache(c, heartbeatTimeout)
		machineSvc.SetStatusCache(hostStatusCache)
		// 启动定时同步（每 1 分钟将 Redis 状态同步到 PostgreSQL）
		syncer := serviceMachine.NewHostStatusSyncer(db, hostStatusCache, 1*time.Minute)
		go syncer.Start(context.Background())
	}
	auditSvc := serviceAudit.NewAuditService(db)
	agentSvc := serviceOps.NewAgentService(db, &config.GlobalConfig.Agent)
	allocSvc := serviceAllocation.NewAllocationService(db, auditSvc, agentSvc)
	allocSvc.StartWorker(context.Background())
	custSvc := serviceCustomer.NewCustomerService(db)
	opsSvc := serviceOps.NewOpsService(db)
	taskSvc := serviceTask.NewTaskService(db, agentSvc)
	datasetSvc := serviceDataset.NewDatasetService(db)
	documentSvc := serviceDocument.NewDocumentService(db, storageMgr)
	sseHub := serviceNotification.NewSSEHub()
	notificationSvc := serviceNotification.NewNotificationService(db, sseHub)

	// Prometheus 客户端
	promClient := prometheus.NewClient(&prometheus.Config{
		Enabled:  config.GlobalConfig.Prometheus.Enabled,
		Endpoint: config.GlobalConfig.Prometheus.Endpoint,
	})

	monitorSvc := serviceOps.NewMonitorService(machineSvc, cache.GetCache(), promClient)
	storageSvc := serviceStorage.NewStorageService(storageMgr)
	sshKeySvc := serviceSSHKey.NewSSHKeyService(db)
	sshKeySvc.SetKeySyncer(allocSvc) // 注入密钥同步器，密钥变更时自动同步到已分配机器
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
	systemConfigSvc.SetAuditService(auditSvc) // 注入审计服务，配置变更时记录审计日志
	dashboardSvc := serviceOps.NewDashboardService(machineSvc, custSvc, allocSvc, promClient)
	workspaceSvc := serviceWorkspace.NewWorkspaceService(db)
	environmentSvc := serviceEnvironment.NewEnvironmentService(db)
	proxySvc := serviceProxy.NewProxyService(db)

	// --- 控制器层初始化 ---
	authController := ctrlAuth.NewAuthController(authSvc)
	dashboardController := ctrlOps.NewDashboardController(dashboardSvc)
	machineController := ctrlMachine.NewMachineController(machineSvc, allocSvc, agentSvc)
	customerController := ctrlCustomer.NewCustomerController(custSvc)
	monitorController := ctrlOps.NewMonitorController(monitorSvc)
	alertController := ctrlOps.NewAlertController(opsSvc)
	agentController := ctrlOps.NewAgentController(machineSvc)

	myMachineController := ctrlCustomer.NewMyMachineController(machineSvc, agentSvc, allocSvc)
	taskController := ctrlTask.NewTaskController(taskSvc)
	adminTaskController := ctrlTask.NewAdminTaskController(taskSvc)
	agentTaskController := ctrlTask.NewAgentTaskController(taskSvc)
	agentHeartbeatController := ctrlAgent.NewHeartbeatController(machineSvc)
	datasetController := ctrlDataset.NewDatasetController(datasetSvc, storageSvc, agentSvc, allocSvc)
	sshKeyController := ctrlCustomer.NewSSHKeyController(sshKeySvc)
	enrollmentController := ctrlCustomer.NewMachineEnrollmentController(enrollmentSvc)
	auditController := ctrlOps.NewAuditController(auditSvc)
	imageController := ctrlOps.NewImageController(imageSvc)
	systemConfigController := ctrlSystemConfig.NewSystemConfigController(systemConfigSvc)
	documentController := ctrlDocument.NewDocumentController(documentSvc, storageSvc)
	storageController := ctrlStorage.NewStorageController(storageSvc)
	notificationController := ctrlNotification.NewNotificationController(notificationSvc)
	workspaceController := ctrlWorkspace.NewWorkspaceController(workspaceSvc)
	environmentController := ctrlEnvironment.NewEnvironmentController(environmentSvc)
	allocationController := ctrlAllocation.NewAllocationController(allocSvc)
	proxyController := ctrlProxy.NewProxyController(proxySvc)

	// API v1 路由
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// SSE 实时通知推送（匹配 Nginx 代理路径）
		apiV1.GET("/notifications/stream", middleware.Auth(db), notificationController.SSE)

		// 1. Auth Module
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/admin/login", authController.AdminLogin)
			authGroup.POST("/refresh", authController.Refresh)
			authGroup.POST("/logout", authController.Logout)

			// 受保护的个人资料
			authGroup.GET("/profile", middleware.Auth(db), authController.GetProfile)
			authGroup.PUT("/profile", middleware.Auth(db), authController.UpdateProfile)
			authGroup.POST("/password/change", middleware.Auth(db), authController.ChangePassword)
			authGroup.POST("/password/request", authController.RequestPasswordReset)
			authGroup.POST("/password/confirm", authController.ConfirmPasswordReset)
		}

		// 2. Admin Module (Protected + Role Check)
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(middleware.Auth(db), middleware.RequireAdmin(), middleware.AuditMiddleware(auditSvc))
		{
			// 仪表板
			adminGroup.GET("/dashboard/stats", dashboardController.GetStats)
			adminGroup.GET("/dashboard/gpu-trend", dashboardController.GetGPUTrend)
			adminGroup.GET("/allocations/recent", dashboardController.GetRecentAllocations)
			adminGroup.GET("/allocations", allocationController.List)

			// 机器管理
			adminGroup.GET("/machines", machineController.List)
			adminGroup.GET("/machines/:id", machineController.Detail)
			adminGroup.POST("/machines", machineController.Create)
			adminGroup.PUT("/machines/:id", machineController.Update)
			adminGroup.POST("/machines/import", machineController.Import)
			adminGroup.DELETE("/machines/:id", machineController.Delete)
			adminGroup.POST("/machines/:id/collect", machineController.CollectSpec)
			adminGroup.POST("/machines/:id/allocate", machineController.Allocate)
			adminGroup.POST("/machines/:id/reclaim", machineController.Reclaim)
			adminGroup.POST("/machines/:id/maintenance", machineController.SetMaintenance)
			adminGroup.GET("/machines/:id/usage", machineController.Usage)

			// 机器批量操作
			adminGroup.POST("/machines/batch/maintenance", machineController.BatchSetMaintenance)
			adminGroup.POST("/machines/batch/allocate", machineController.BatchAllocate)
			adminGroup.POST("/machines/batch/reclaim", machineController.BatchReclaim)

			// 客户管理
			adminGroup.GET("/customers", customerController.List)
			adminGroup.GET("/customers/:id", customerController.Detail)
			adminGroup.POST("/customers", customerController.Create)
			adminGroup.PUT("/customers/:id", customerController.Update)
			adminGroup.POST("/customers/:id/disable", customerController.Disable)
			adminGroup.POST("/customers/:id/enable", customerController.Enable)
			adminGroup.PUT("/customers/:id/quota", customerController.UpdateQuota)
			adminGroup.GET("/customers/:id/usage", customerController.ResourceUsage)

			// 监控与运维
			adminGroup.GET("/monitoring/realtime", monitorController.GetRealtime)
			adminGroup.GET("/alerts", alertController.List)
			adminGroup.POST("/alerts/:id/acknowledge", alertController.Acknowledge)

			// 告警规则管理
			adminGroup.GET("/alert-rules", alertController.ListRules)
			adminGroup.GET("/alert-rules/:id", alertController.GetRule)
			adminGroup.POST("/alert-rules", alertController.CreateRule)
			adminGroup.PUT("/alert-rules/:id", alertController.UpdateRule)
			adminGroup.DELETE("/alert-rules/:id", alertController.DeleteRule)
			adminGroup.POST("/alert-rules/:id/toggle", alertController.ToggleRule)

			// 审计日志
			adminGroup.GET("/audit/logs", auditController.List)

			// 镜像管理
			adminGroup.GET("/images", imageController.List)
			adminGroup.POST("/images/sync", imageController.Sync)
			adminGroup.DELETE("/images/:id", imageController.Delete)

			// 系统配置
			adminGroup.GET("/settings/configs", systemConfigController.GetConfigs)
			adminGroup.PUT("/settings/configs", systemConfigController.UpdateConfigs)
			adminGroup.GET("/settings/configs/groups", systemConfigController.ListGroups)
			adminGroup.GET("/settings/configs/:id", systemConfigController.GetConfig)
			adminGroup.POST("/settings/configs", systemConfigController.CreateConfig)
			adminGroup.PUT("/settings/configs/:id", systemConfigController.UpdateConfig)
			adminGroup.DELETE("/settings/configs/:id", systemConfigController.DeleteConfig)

			// 任务管理
			adminGroup.GET("/tasks", adminTaskController.List)
			adminGroup.GET("/tasks/:id", adminTaskController.Detail)
			adminGroup.POST("/tasks", adminTaskController.Create)
			adminGroup.POST("/tasks/:id/stop", adminTaskController.Stop)
			adminGroup.POST("/tasks/:id/cancel", adminTaskController.Cancel)
			adminGroup.POST("/tasks/:id/retry", adminTaskController.Retry)
			adminGroup.GET("/tasks/:id/logs", adminTaskController.Logs)
			adminGroup.GET("/tasks/:id/result", adminTaskController.Result)

			// Agent 管理
			adminGroup.GET("/agents", agentController.List)

			// 文档中心
			adminGroup.GET("/documents", documentController.List)
			adminGroup.GET("/documents/categories", documentController.Categories)
			adminGroup.GET("/documents/:id", documentController.Detail)
			adminGroup.POST("/documents", documentController.Upload)
			adminGroup.PUT("/documents/:id", documentController.Update)
			adminGroup.DELETE("/documents/:id", documentController.Delete)
			adminGroup.GET("/documents/:id/download", documentController.Download)

			// 存储管理
			adminGroup.GET("/storage/backends", storageController.ListBackends)
			adminGroup.GET("/storage/stats", storageController.GetStats)
			adminGroup.GET("/storage/files", storageController.ListFiles)
			adminGroup.POST("/storage/files/delete", storageController.DeleteFile)
			adminGroup.GET("/storage/files/download-url", storageController.GetDownloadURL)

			// Proxy 管理
			adminGroup.GET("/proxy/nodes", proxyController.ListNodes)
			adminGroup.GET("/proxy/nodes/:id", proxyController.GetNode)
			adminGroup.DELETE("/proxy/nodes/:id", proxyController.DeleteNode)
			adminGroup.GET("/proxy/mappings", proxyController.ListMappings)
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
			custGroup.GET("/tasks/:id", taskController.Detail)
			custGroup.POST("/tasks/training", taskController.CreateTraining)
			custGroup.POST("/tasks/:id/stop", taskController.Stop)
			custGroup.POST("/tasks/:id/cancel", taskController.Cancel)
			custGroup.POST("/tasks/:id/retry", taskController.Retry)
			custGroup.GET("/tasks/:id/logs", taskController.Logs)
			custGroup.GET("/tasks/:id/result", taskController.Result)

			// 数据集管理
			custGroup.GET("/datasets", datasetController.List)
			custGroup.POST("/datasets/init-multipart", datasetController.InitUpload)
			custGroup.POST("/datasets/:id/complete", datasetController.CompleteUpload)
			custGroup.POST("/datasets/:id/mount", datasetController.Mount)
			custGroup.GET("/datasets/:id/mounts", datasetController.ListMounts)
			custGroup.POST("/datasets/:id/mounts/:mount_id/unmount", datasetController.Unmount)

			// SSH 密钥管理
			custGroup.GET("/keys", sshKeyController.List)
			custGroup.POST("/keys", sshKeyController.Create)
			custGroup.DELETE("/keys/:id", sshKeyController.Delete)

			// 用户添加机器
			custGroup.POST("/machines/enroll", enrollmentController.Create)
			custGroup.GET("/machines/enrollments", enrollmentController.List)
			custGroup.GET("/machines/enrollments/:id", enrollmentController.Detail)

			// 通知管理
			custGroup.GET("/notifications/sse", notificationController.SSE)
			custGroup.GET("/notifications", notificationController.List)
			custGroup.GET("/notifications/unread-count", notificationController.UnreadCount)
			custGroup.POST("/notifications/:id/read", notificationController.MarkRead)
			custGroup.POST("/notifications/read-all", notificationController.MarkAllRead)

			// 工作空间管理
			custGroup.POST("/workspaces", workspaceController.Create)
			custGroup.GET("/workspaces", workspaceController.List)
			custGroup.GET("/workspaces/:id", workspaceController.Detail)
			custGroup.PUT("/workspaces/:id", workspaceController.Update)
			custGroup.DELETE("/workspaces/:id", workspaceController.Delete)
			custGroup.POST("/workspaces/:id/members", workspaceController.AddMember)
			custGroup.DELETE("/workspaces/:id/members/:userId", workspaceController.RemoveMember)
			custGroup.GET("/workspaces/:id/members", workspaceController.ListMembers)

			// 环境管理
			custGroup.POST("/environments", environmentController.Create)
			custGroup.GET("/environments", environmentController.List)
			custGroup.GET("/environments/:id", environmentController.Detail)
			custGroup.POST("/environments/:id/start", environmentController.Start)
			custGroup.POST("/environments/:id/stop", environmentController.Stop)
			custGroup.DELETE("/environments/:id", environmentController.Delete)
			custGroup.GET("/environments/:id/access", environmentController.AccessInfo)
		}

		// 4. Agent Module (Agent 专用 API，需要 Agent Token 认证)
		agentGroup := apiV1.Group("/agent")
		agentGroup.Use(middleware.AgentAuth())
		{
			agentGroup.POST("/register", agentHeartbeatController.Register)
			agentGroup.POST("/heartbeat", agentHeartbeatController.Heartbeat)
			agentGroup.POST("/tasks/claim", agentTaskController.ClaimTasks)
			agentGroup.POST("/tasks/:id/start", agentTaskController.StartTask)
			agentGroup.POST("/tasks/:id/lease/renew", agentTaskController.RenewLease)
			agentGroup.POST("/tasks/:id/complete", agentTaskController.CompleteTask)
			agentGroup.POST("/tasks/:id/progress", agentTaskController.ReportProgress)
		}

		// 5. Proxy Module (Proxy 专用 API，需要 Agent Token 认证)
		proxyGroup := apiV1.Group("/proxy")
		proxyGroup.Use(middleware.AgentAuth())
		{
			proxyGroup.POST("/register", proxyController.Register)
			proxyGroup.POST("/heartbeat", proxyController.Heartbeat)
		}
	}
}
