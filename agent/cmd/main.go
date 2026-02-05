package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/YoungBoyGod/remotegpu-agent/internal/handler"
	"github.com/YoungBoyGod/remotegpu-agent/internal/scheduler"
	"github.com/gin-gonic/gin"
)

var version = "0.1.0"

func main() {
	port := os.Getenv("AGENT_PORT")
	if port == "" {
		port = "8090"
	}

	dbPath := os.Getenv("AGENT_DB_PATH")
	if dbPath == "" {
		dbPath = "/var/lib/remotegpu-agent/tasks.db"
	}

	// 确保数据目录存在
	os.MkdirAll("/var/lib/remotegpu-agent", 0755)

	// 创建调度器
	sched, err := scheduler.NewScheduler(dbPath, 4)
	if err != nil {
		log.Fatalf("create scheduler error: %v", err)
	}

	// 启动调度器
	if err := sched.Start(); err != nil {
		log.Fatalf("start scheduler error: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册路由
	taskHandler := handler.NewTaskHandler(sched)
	registerRoutes(r, taskHandler)

	fmt.Printf("RemoteGPU Agent v%s starting on :%s\n", version, port)

	// 优雅关闭
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")
		sched.Stop()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
