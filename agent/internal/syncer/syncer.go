package syncer

import (
	"log/slog"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/client"
	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
	"github.com/YoungBoyGod/remotegpu-agent/internal/store"
)

// Syncer 离线结果同步器
type Syncer struct {
	store    *store.SQLiteStore
	client   *client.ServerClient
	interval time.Duration
	stopCh   chan struct{}
}

// NewSyncer 创建同步器
func NewSyncer(st *store.SQLiteStore, c *client.ServerClient, interval time.Duration) *Syncer {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Syncer{
		store:    st,
		client:   c,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start 启动同步循环
func (s *Syncer) Start() {
	go s.syncLoop()
}

// Stop 停止同步器
func (s *Syncer) Stop() {
	close(s.stopCh)
}

func (s *Syncer) syncLoop() {
	// 启动时立即同步一次
	s.syncUnreported()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.syncUnreported()
		}
	}
}

func (s *Syncer) syncUnreported() {
	tasks, err := s.store.ListUnsynced()
	if err != nil {
		slog.Error("list unsynced tasks failed", "error", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	slog.Info("found unsynced tasks", "count", len(tasks))

	for _, task := range tasks {
		if err := s.syncTask(task); err != nil {
			slog.Error("sync task failed", "task_id", task.ID, "error", err)
			continue
		}
		if err := s.store.MarkSynced(task.ID); err != nil {
			slog.Error("mark synced failed", "task_id", task.ID, "error", err)
		}
	}
}

func (s *Syncer) syncTask(task *models.Task) error {
	// 只有有 AttemptID 的任务才需要上报 Server
	if task.AttemptID == "" {
		return nil
	}
	return s.client.ReportComplete(task)
}
