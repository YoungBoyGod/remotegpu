package hotreload

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Config 热更新配置
type Config struct {
	Enabled       bool          // 是否启用热更新
	WatchDirs     []string      // 监控的目录
	WatchExts     []string      // 监控的文件扩展名
	ExcludeDirs   []string      // 排除的目录
	BuildCmd      string        // 构建命令
	Debounce      time.Duration // 防抖时间
	RestartSignal chan struct{} // 重启信号通道
}

// Manager 热更新管理器
type Manager struct {
	config  Config
	watcher *fsnotify.Watcher
	mu      sync.Mutex
	timer   *time.Timer
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewManager 创建热更新管理器
func NewManager(config Config) (*Manager, error) {
	if !config.Enabled {
		return nil, nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("创建文件监控器失败: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config:  config,
		watcher: watcher,
		ctx:     ctx,
		cancel:  cancel,
	}

	// 设置默认值
	if m.config.Debounce == 0 {
		m.config.Debounce = 1 * time.Second
	}
	if len(m.config.WatchExts) == 0 {
		m.config.WatchExts = []string{".go", ".yaml", ".yml"}
	}
	if len(m.config.ExcludeDirs) == 0 {
		m.config.ExcludeDirs = []string{"vendor", "node_modules", ".git", "tmp", "logs"}
	}

	return m, nil
}

// Start 启动热更新监控
func (m *Manager) Start() error {
	if m == nil {
		return nil
	}

	// 添加监控目录
	for _, dir := range m.config.WatchDirs {
		if err := m.addWatchDir(dir); err != nil {
			return fmt.Errorf("添加监控目录失败 %s: %w", dir, err)
		}
	}

	go m.watch()
	return nil
}

// Stop 停止热更新监控
func (m *Manager) Stop() error {
	if m == nil {
		return nil
	}

	m.cancel()
	if m.timer != nil {
		m.timer.Stop()
	}
	return m.watcher.Close()
}

// addWatchDir 递归添加监控目录
func (m *Manager) addWatchDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// 检查是否在排除列表中
		for _, exclude := range m.config.ExcludeDirs {
			if strings.Contains(path, exclude) {
				return filepath.SkipDir
			}
		}

		return m.watcher.Add(path)
	})
}

// watch 监控文件变化
func (m *Manager) watch() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}
			m.handleEvent(event)
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("文件监控错误: %v\n", err)
		}
	}
}

// handleEvent 处理文件变化事件
func (m *Manager) handleEvent(event fsnotify.Event) {
	// 只处理写入和创建事件
	if event.Op&fsnotify.Write != fsnotify.Write && event.Op&fsnotify.Create != fsnotify.Create {
		return
	}

	// 检查文件扩展名
	ext := filepath.Ext(event.Name)
	matched := false
	for _, watchExt := range m.config.WatchExts {
		if ext == watchExt {
			matched = true
			break
		}
	}
	if !matched {
		return
	}

	fmt.Printf("检测到文件变化: %s\n", event.Name)

	// 使用防抖机制
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.timer != nil {
		m.timer.Stop()
	}

	m.timer = time.AfterFunc(m.config.Debounce, func() {
		m.rebuild()
	})
}

// rebuild 重新构建并重启
func (m *Manager) rebuild() {
	fmt.Println("开始重新构建...")

	if m.config.BuildCmd != "" {
		// 执行构建命令
		parts := strings.Fields(m.config.BuildCmd)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("构建失败: %v\n", err)
			return
		}
		fmt.Println("构建成功")
	}

	// 发送重启信号
	if m.config.RestartSignal != nil {
		select {
		case m.config.RestartSignal <- struct{}{}:
		default:
		}
	}
}
