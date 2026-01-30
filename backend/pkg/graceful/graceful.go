package graceful

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config 优雅启动配置
type Config struct {
	ShutdownTimeout time.Duration // 关闭超时时间
	RetryInterval   time.Duration // 重试间隔
	MaxRetries      int           // 最大重试次数，0表示无限重试
}

// Server 服务器接口
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// Manager 优雅启动管理器
type Manager struct {
	config Config
	server Server
}

// NewManager 创建优雅启动管理器
func NewManager(config Config, server Server) *Manager {
	// 设置默认值
	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 10 * time.Second
	}
	if config.RetryInterval == 0 {
		config.RetryInterval = 5 * time.Second
	}

	return &Manager{
		config: config,
		server: server,
	}
}

// Run 运行服务器，支持优雅启动和关闭
func (m *Manager) Run(restartChan <-chan struct{}) error {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器（带重试）
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- m.startWithRetry()
	}()

	// 等待信号
	for {
		select {
		case err := <-serverErr:
			if err != nil {
				return err
			}
		case <-quit:
			fmt.Println("\n收到关闭信号，开始优雅关闭...")
			return m.shutdown()
		case <-restartChan:
			fmt.Println("\n收到重启信号，开始重启...")
			if err := m.shutdown(); err != nil {
				fmt.Printf("关闭服务失败: %v\n", err)
			}
			// 重新启动
			go func() {
				serverErr <- m.startWithRetry()
			}()
		}
	}
}

// startWithRetry 启动服务器，失败时重试
func (m *Manager) startWithRetry() error {
	retries := 0
	for {
		err := m.server.Start()
		if err == nil {
			return nil
		}

		// 检查是否是正常关闭
		if err == http.ErrServerClosed {
			return nil
		}

		retries++
		if m.config.MaxRetries > 0 && retries >= m.config.MaxRetries {
			return fmt.Errorf("启动失败，已达到最大重试次数 %d: %w", m.config.MaxRetries, err)
		}

		fmt.Printf("启动失败 (第 %d 次): %v\n", retries, err)
		fmt.Printf("等待 %v 后重试...\n", m.config.RetryInterval)
		time.Sleep(m.config.RetryInterval)
	}
}

// shutdown 优雅关闭服务器
func (m *Manager) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.ShutdownTimeout)
	defer cancel()

	if err := m.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	fmt.Println("服务器已优雅关闭")
	return nil
}
