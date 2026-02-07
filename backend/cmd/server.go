package main

import (
	"fmt"
	"log"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/router"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	pkgCache "github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/graceful"
	"github.com/YoungBoyGod/remotegpu/pkg/hotreload"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 API 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func startServer() {
	// 加载配置
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(mode); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.GetLogger().Sync()

	// 创建重启信号通道
	restartChan := make(chan struct{}, 1)

	// 初始化热更新管理器
	hotReloadMgr, err := initHotReload(restartChan)
	if err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("初始化热更新失败: %v", err))
	}
	if hotReloadMgr != nil {
		if err := hotReloadMgr.Start(); err != nil {
			logger.GetLogger().Fatal(fmt.Sprintf("启动热更新失败: %v", err))
		}
		defer hotReloadMgr.Stop()
		logger.GetLogger().Info("热更新已启动")
	}

	// 初始化并运行服务器
	if err := runServer(restartChan); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("服务器运行失败: %v", err))
	}
}

// runServer 运行服务器
func runServer(restartChan chan struct{}) error {
	// 初始化基础设施
	if err := initInfrastructure(); err != nil {
		return err
	}

	// 创建 Gin 引擎
	r := createGinEngine()

	// 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	httpServer := graceful.NewHTTPServer(r, addr)

	// 创建优雅启动管理器
	gracefulCfg := graceful.Config{
		ShutdownTimeout: time.Duration(config.GlobalConfig.Graceful.ShutdownTimeout) * time.Second,
		RetryInterval:   time.Duration(config.GlobalConfig.Graceful.RetryInterval) * time.Second,
		MaxRetries:      config.GlobalConfig.Graceful.MaxRetries,
	}
	gracefulMgr := graceful.NewManager(gracefulCfg, httpServer)

	// 运行服务器
	return gracefulMgr.Run(restartChan)
}

// initInfrastructure 初始化基础设施（数据库、Redis等）
func initInfrastructure() error {
	// 初始化 JWT
	if err := auth.InitJWT(config.GlobalConfig.JWT.Secret, config.GlobalConfig.JWT.ExpireTime); err != nil {
		return fmt.Errorf("初始化 JWT 失败: %w", err)
	}
	logger.GetLogger().Info("JWT 初始化完成")

	// 初始化数据库
	dbConfig := database.Config{
		Host:     config.GlobalConfig.Database.Host,
		Port:     config.GlobalConfig.Database.Port,
		User:     config.GlobalConfig.Database.User,
		Password: config.GlobalConfig.Database.Password,
		DBName:   config.GlobalConfig.Database.DBName,
	}
	if err := database.InitDB(dbConfig); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 自动迁移数据库表 (V2.0 Schema)
	err := database.GetDB().AutoMigrate(
		&entity.Customer{},
		&entity.SSHKey{},
		&entity.Workspace{},
		&entity.Host{},
		&entity.GPU{},
		&entity.Allocation{},
		&entity.Image{},
		&entity.Dataset{},
		&entity.DatasetMount{},
		&entity.Task{},
		&entity.AuditLog{},
		&entity.AlertRule{},
		&entity.ActiveAlert{},
		&entity.MachineEnrollment{},
		&entity.HostMetric{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	logger.GetLogger().Info("数据库表迁移完成 (V2.0)")

	// 初始化 Redis (Using pkg/cache)
	redisConfig := pkgCache.Config{
		Host:     config.GlobalConfig.Redis.Host,
		Port:     config.GlobalConfig.Redis.Port,
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
	}
	if err := pkgCache.InitRedis(redisConfig); err != nil {
		return fmt.Errorf("初始化 Redis 失败: %w", err)
	}

	return nil
}

// createGinEngine 创建 Gin 引擎
func createGinEngine() *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(getGinMode(mode))

	// 创建 Gin 引擎（不使用默认中间件）
	r := gin.New()

	// 添加自定义中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(logger.GetLogger()))
	r.Use(middleware.Recovery(logger.GetLogger()))

	// 初始化路由
	router.InitRouter(r)

	return r
}

// initHotReload 初始化热更新管理器
func initHotReload(restartChan chan struct{}) (*hotreload.Manager, error) {
	if !config.GlobalConfig.HotReload.Enabled {
		return nil, nil
	}

	cfg := hotreload.Config{
		Enabled:       config.GlobalConfig.HotReload.Enabled,
		WatchDirs:     config.GlobalConfig.HotReload.WatchDirs,
		WatchExts:     config.GlobalConfig.HotReload.WatchExts,
		ExcludeDirs:   config.GlobalConfig.HotReload.ExcludeDirs,
		BuildCmd:      config.GlobalConfig.HotReload.BuildCmd,
		Debounce:      time.Duration(config.GlobalConfig.HotReload.Debounce) * time.Second,
		RestartSignal: restartChan,
	}

	return hotreload.NewManager(cfg)
}

func getGinMode(mode string) string {
	switch mode {
	case "release":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
