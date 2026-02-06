package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/client"
	agentcfg "github.com/YoungBoyGod/remotegpu-agent/internal/config"
	"github.com/YoungBoyGod/remotegpu-agent/internal/handler"
	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
	"github.com/YoungBoyGod/remotegpu-agent/internal/poller"
	"github.com/YoungBoyGod/remotegpu-agent/internal/scheduler"
	"github.com/YoungBoyGod/remotegpu-agent/internal/syncer"
	"github.com/gin-gonic/gin"
)

var version = "0.1.0"

func main() {
	cfg := agentcfg.Load()

	port := strconv.Itoa(cfg.Port)

	// 确保数据目录存在
	os.MkdirAll(filepath.Dir(cfg.DBPath), 0755)

	// 创建调度器
	sched, err := scheduler.NewScheduler(cfg.DBPath, cfg.MaxWorkers)
	if err != nil {
		log.Fatalf("create scheduler error: %v", err)
	}

	// 启动调度器
	if err := sched.Start(); err != nil {
		log.Fatalf("start scheduler error: %v", err)
	}

	// 启动 Poller 和 Syncer（如果配置了 Server）
	var p *poller.Poller
	var sy *syncer.Syncer
	if cfg.ServerConfigured() {
		serverClient := client.NewServerClient(&client.Config{
			ServerURL: cfg.Server.URL,
			AgentID:   cfg.Server.AgentID,
			MachineID: cfg.Server.MachineID,
			Token:     cfg.Server.Token,
			Timeout:   cfg.Server.Timeout,
		})
		sched.SetClient(serverClient)

		p = poller.NewPoller(&poller.Config{
			Client:    serverClient,
			Interval:  cfg.Poll.Interval,
			BatchSize: cfg.Poll.BatchSize,
			OnTask: func(task *models.Task) {
				if err := sched.Submit(task); err != nil {
					slog.Error("submit task error", "error", err)
				}
			},
		})
		p.Start()
		slog.Info("poller started", "server", cfg.Server.URL)

		// 启动离线结果同步器
		sy = syncer.NewSyncer(sched.GetStore(), serverClient, 30*time.Second)
		sy.Start()
		slog.Info("syncer started")
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
		slog.Info("shutting down...")
		if sy != nil {
			sy.Stop()
		}
		if p != nil {
			p.Stop()
		}
		sched.Stop()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
