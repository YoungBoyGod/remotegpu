package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/YoungBoyGod/remotegpu-proxy/internal/client"
	proxycfg "github.com/YoungBoyGod/remotegpu-proxy/internal/config"
	"github.com/YoungBoyGod/remotegpu-proxy/internal/forwarder"
	"github.com/YoungBoyGod/remotegpu-proxy/internal/handler"
	"github.com/YoungBoyGod/remotegpu-proxy/internal/portpool"
	"github.com/gin-gonic/gin"
)

var version = "0.1.0"

func main() {
	cfg := proxycfg.Load()

	// 初始化端口池
	pool := portpool.NewPool(cfg.PortPool.RangeStart, cfg.PortPool.RangeEnd)

	// 初始化 HTTP 反向代理（可选）
	var httpProxy *forwarder.HTTPProxy
	if cfg.HTTPProxy.Enabled {
		httpProxy = forwarder.NewHTTPProxy(cfg.HTTPProxy.Port)
	}

	// 初始化转发管理器
	mgr := forwarder.NewManager(pool, httpProxy)

	// 启动 HTTP 反向代理
	if httpProxy != nil {
		if err := httpProxy.Start(); err != nil {
			log.Fatalf("启动 HTTP 反向代理失败: %v", err)
		}
		slog.Info("HTTP 反向代理已启动", "port", cfg.HTTPProxy.Port)
	}

	// 启动后端通信客户端和心跳
	var serverClient *client.ServerClient
	var heartbeatTicker *time.Ticker
	if cfg.ServerConfigured() {
		serverClient = client.NewServerClient(&client.Config{
			ServerURL: cfg.Server.URL,
			ProxyID:   cfg.Server.ProxyID,
			Token:     cfg.Server.Token,
			Timeout:   cfg.Server.Timeout,
		})

		// 注册 Proxy
		if err := serverClient.Register(
			cfg.Network.InnerIP,
			cfg.Network.OuterIP,
			cfg.PortPool.RangeStart,
			cfg.PortPool.RangeEnd,
		); err != nil {
			slog.Error("注册 Proxy 失败", "error", err)
		} else {
			slog.Info("Proxy 已注册到后端")
		}

		// 启动心跳
		heartbeatTicker = time.NewTicker(cfg.Heartbeat.Interval)
		go func() {
			for range heartbeatTicker.C {
				stats := pool.Stats()
				usage := &client.PortUsage{
					Total:     stats.Total,
					Used:      stats.Used,
					Available: stats.Available,
				}
				if err := serverClient.Heartbeat(usage); err != nil {
					slog.Error("心跳发送失败", "error", err)
				} else {
					slog.Debug("心跳已发送")
				}
			}
		}()
		slog.Info("心跳已启动", "interval", cfg.Heartbeat.Interval)
	}

	// 启动管理 API
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	mappingHandler := handler.NewMappingHandler(mgr)
	registerRoutes(r, mappingHandler)

	port := strconv.Itoa(cfg.Port)
	fmt.Printf("RemoteGPU Proxy v%s starting on :%s\n", version, port)

	// 优雅关闭
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		slog.Info("正在关闭...")
		if heartbeatTicker != nil {
			heartbeatTicker.Stop()
		}
		if httpProxy != nil {
			httpProxy.Stop()
		}
		mgr.StopAll()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
