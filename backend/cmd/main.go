package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/router"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	pkgRedis "github.com/YoungBoyGod/remotegpu/pkg/redis"
	"github.com/gin-gonic/gin"
)

var (
	configPath = flag.String("config", "./config/config.yaml", "配置文件路径")
	mode       = flag.String("mode", "debug", "运行模式: debug, release, test")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 加载配置
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(*mode); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.GetLogger().Sync()

	// 初始化数据库
	dbConfig := database.Config{
		Host:     config.GlobalConfig.Database.Host,
		Port:     config.GlobalConfig.Database.Port,
		User:     config.GlobalConfig.Database.User,
		Password: config.GlobalConfig.Database.Password,
		DBName:   config.GlobalConfig.Database.DBName,
	}
	if err := database.InitDB(dbConfig); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 自动迁移数据库表
	if err := database.GetDB().AutoMigrate(&entity.Customer{}); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("数据库迁移失败: %v", err))
	}
	logger.GetLogger().Info("数据库表迁移完成")

	// 初始化 Redis
	redisConfig := pkgRedis.Config{
		Host:     config.GlobalConfig.Redis.Host,
		Port:     config.GlobalConfig.Redis.Port,
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
	}
	if err := pkgRedis.InitRedis(redisConfig); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("初始化 Redis 失败: %v", err))
	}

	// 设置 Gin 模式
	gin.SetMode(getGinMode(*mode))

	// 创建 Gin 引擎（不使用默认中间件）
	r := gin.New()

	// 添加自定义中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(logger.GetLogger()))
	r.Use(middleware.Recovery(logger.GetLogger()))

	// 初始化路由
	router.InitRouter(r)

	// 启动服务
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	logger.GetLogger().Info(fmt.Sprintf("服务启动在 %s", addr))
	if err := r.Run(addr); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("启动服务失败: %v", err))
	}
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
