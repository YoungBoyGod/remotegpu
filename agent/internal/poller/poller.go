package poller

import (
	"log/slog"
	"sync"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/client"
	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

// Poller 任务轮询器
type Poller struct {
	client       *client.ServerClient
	interval     time.Duration
	batchSize    int
	taskCallback func(*models.Task)

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// Config 轮询器配置
type Config struct {
	Client    *client.ServerClient
	Interval  time.Duration
	BatchSize int
	OnTask    func(*models.Task)
}

// NewPoller 创建轮询器
func NewPoller(cfg *Config) *Poller {
	interval := cfg.Interval
	if interval == 0 {
		interval = 5 * time.Second
	}
	batchSize := cfg.BatchSize
	if batchSize == 0 {
		batchSize = 10
	}

	return &Poller{
		client:       cfg.Client,
		interval:     interval,
		batchSize:    batchSize,
		taskCallback: cfg.OnTask,
		stopCh:       make(chan struct{}),
	}
}

// Start 启动轮询
func (p *Poller) Start() {
	p.wg.Add(1)
	go p.pollLoop()
}

// Stop 停止轮询
func (p *Poller) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *Poller) pollLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.poll()
		}
	}
}

func (p *Poller) poll() {
	tasks, err := p.client.ClaimTasks(p.batchSize)
	if err != nil {
		slog.Error("claim tasks failed", "error", err)
		return
	}

	for _, task := range tasks {
		if p.taskCallback != nil {
			p.taskCallback(task)
		}
	}
}
