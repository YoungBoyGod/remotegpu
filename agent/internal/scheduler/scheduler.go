package scheduler

import (
	"log/slog"
	"sync"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/client"
	"github.com/YoungBoyGod/remotegpu-agent/internal/executor"
	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
	"github.com/YoungBoyGod/remotegpu-agent/internal/queue"
	"github.com/YoungBoyGod/remotegpu-agent/internal/store"
)

// Scheduler 任务调度器
type Scheduler struct {
	queue    *queue.Manager
	store    *store.SQLiteStore
	executor *executor.Executor
	client   *client.ServerClient

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewScheduler 创建调度器
func NewScheduler(dbPath string, maxWorkers int) (*Scheduler, error) {
	st, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}

	s := &Scheduler{
		queue:    queue.NewManager(),
		store:    st,
		executor: executor.NewExecutor(maxWorkers),
		stopCh:   make(chan struct{}),
	}

	return s, nil
}

// SetClient 设置 Server 客户端
func (s *Scheduler) SetClient(c *client.ServerClient) {
	s.client = c
}

// GetExecutor 返回执行器（供外部设置 validator 等）
func (s *Scheduler) GetExecutor() *executor.Executor {
	return s.executor
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	// 恢复未完成的任务
	if err := s.recover(); err != nil {
		slog.Error("recover tasks failed", "error", err)
	}

	// 启动调度循环
	s.wg.Add(1)
	go s.scheduleLoop()

	return nil
}

// recover 恢复未完成的任务
func (s *Scheduler) recover() error {
	// 1. 恢复 pending 任务到队列（无需校验）
	pendingTasks, err := s.store.ListByStatus(models.TaskStatusPending)
	if err != nil {
		return err
	}
	for _, task := range pendingTasks {
		s.queue.Push(task)
	}
	slog.Info("recovered pending tasks", "count", len(pendingTasks))

	// 2. 处理 assigned 任务：检查租约是否过期
	assignedTasks, err := s.store.ListByStatus(models.TaskStatusAssigned)
	if err != nil {
		return err
	}
	for _, task := range assignedTasks {
		if !task.LeaseExpiresAt.IsZero() && time.Now().After(task.LeaseExpiresAt) {
			task.Status = models.TaskStatusFailed
			task.Error = "lease expired during agent restart"
			task.EndedAt = time.Now()
			s.store.Save(task)
			slog.Warn("task lease expired, marked as failed", "task_id", task.ID)
		} else {
			s.queue.Push(task)
			slog.Info("task lease still valid, re-queued", "task_id", task.ID)
		}
	}

	// 3. 处理 running 任务：进程已丢失，标记为 failed
	runningTasks, err := s.store.ListByStatus(models.TaskStatusRunning)
	if err != nil {
		return err
	}
	for _, task := range runningTasks {
		task.Status = models.TaskStatusFailed
		task.Error = "process lost during agent restart"
		task.EndedAt = time.Now()
		s.store.Save(task)
		if s.client != nil && task.AttemptID != "" {
			if err := s.client.ReportComplete(task); err != nil {
				slog.Error("report failed task error", "task_id", task.ID, "error", err)
			}
		}
		slog.Warn("task was running, marked as failed (process lost)", "task_id", task.ID)
	}

	return nil
}

// scheduleLoop 调度循环
func (s *Scheduler) scheduleLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-s.queue.NotifyChan():
			s.tryExecute()
		case <-ticker.C:
			s.tryExecute()
		}
	}
}

// tryExecute 尝试执行队列中的任务
func (s *Scheduler) tryExecute() {
	for s.executor.CanAccept() {
		task := s.queue.Pop()
		if task == nil {
			return
		}

		// 检查任务依赖是否满足
		if !s.dependenciesMet(task) {
			// 依赖未满足，延迟后重新入队
			time.AfterFunc(2*time.Second, func() {
				s.queue.Push(task)
			})
			continue
		}

		go s.runTask(task)
	}
}

// dependenciesMet 检查任务的所有依赖是否已完成
func (s *Scheduler) dependenciesMet(task *models.Task) bool {
	if len(task.DependsOn) == 0 {
		return true
	}
	for _, depID := range task.DependsOn {
		dep, err := s.store.Get(depID)
		if err != nil {
			// 依赖任务不在本地存储中，视为未满足
			return false
		}
		if dep.Status != models.TaskStatusCompleted {
			return false
		}
	}
	return true
}

// runTask 执行单个任务（包含状态同步）
func (s *Scheduler) runTask(task *models.Task) {
	// 上报任务开始
	if s.client != nil && task.AttemptID != "" {
		if err := s.client.ReportStart(task.ID, task.AttemptID); err != nil {
			slog.Error("report start error", "task_id", task.ID, "error", err)
		}
	}

	// 启动租约续约
	stopRenew := make(chan struct{})
	if s.client != nil && task.AttemptID != "" {
		go s.renewLoop(task, stopRenew)
	}

	// 执行任务
	s.executor.Execute(task)

	// 停止续约
	close(stopRenew)

	// 检查是否需要重试
	if task.Status == models.TaskStatusFailed && task.MaxRetries > 0 && task.RetryCount < task.MaxRetries {
		task.RetryCount++
		task.Status = models.TaskStatusPending
		task.AttemptID = "" // 清除旧的 attempt，重试作为本地任务重新调度
		task.Error = ""
		task.ExitCode = 0
		task.Stdout = ""
		task.Stderr = ""

		if err := s.store.Save(task); err != nil {
			slog.Error("save retry task error", "task_id", task.ID, "error", err)
		}

		delay := time.Duration(task.RetryDelay) * time.Second
		if delay <= 0 {
			delay = 60 * time.Second
		}
		slog.Info("task failed, scheduling retry", "task_id", task.ID, "retry", task.RetryCount, "max_retries", task.MaxRetries, "delay", delay)
		time.AfterFunc(delay, func() {
			s.queue.Push(task)
		})
		return
	}

	// 保存结果
	if err := s.store.Save(task); err != nil {
		slog.Error("save task error", "task_id", task.ID, "error", err)
	}

	// 上报完成
	if s.client != nil && task.AttemptID != "" {
		if err := s.client.ReportComplete(task); err != nil {
			slog.Error("report complete error", "task_id", task.ID, "error", err)
		}
	}
}

// renewLoop 租约续约 + 进度上报循环
func (s *Scheduler) renewLoop(task *models.Task, stop <-chan struct{}) {
	renewTicker := time.NewTicker(60 * time.Second)
	progressTicker := time.NewTicker(30 * time.Second)
	defer renewTicker.Stop()
	defer progressTicker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-renewTicker.C:
			if err := s.client.RenewLease(task.ID, task.AttemptID); err != nil {
				slog.Error("renew lease error", "task_id", task.ID, "error", err)
			}
		case <-progressTicker.C:
			if err := s.client.ReportProgress(task.ID, task.AttemptID, 0, "running"); err != nil {
				slog.Debug("report progress error", "task_id", task.ID, "error", err)
			}
		}
	}
}

// Submit 提交任务
func (s *Scheduler) Submit(task *models.Task) error {
	// 仅对本地提交的任务设置默认值；Server 下发的任务（有 AttemptID）保留原始状态
	if task.AttemptID == "" {
		if task.Status == "" {
			task.Status = models.TaskStatusPending
		}
		if task.CreatedAt.IsZero() {
			task.CreatedAt = time.Now()
		}
	}

	if err := s.store.Save(task); err != nil {
		return err
	}

	s.queue.Push(task)

	// 抢占检查：新任务优先级比当前运行中最低优先级任务高 3 级以上时触发抢占
	s.tryPreempt(task)

	return nil
}

// tryPreempt 尝试抢占低优先级任务
func (s *Scheduler) tryPreempt(newTask *models.Task) {
	lowest := s.executor.LowestPriorityRunning()
	if lowest == nil {
		return
	}
	// 优先级数字越小越高，差值 >= 3 才抢占
	if lowest.Priority-newTask.Priority >= 3 {
		slog.Info("preempting task",
			"preempted_id", lowest.ID,
			"preempted_priority", lowest.Priority,
			"new_id", newTask.ID,
			"new_priority", newTask.Priority,
		)
		s.executor.Cancel(lowest.ID)
		// 被抢占的任务标记为 preempted，稍后重新入队
		lowest.Status = models.TaskStatusPreempted
		lowest.Error = "preempted by higher priority task " + newTask.ID
		s.store.Save(lowest)
		time.AfterFunc(1*time.Second, func() {
			lowest.Status = models.TaskStatusPending
			lowest.Error = ""
			lowest.ExitCode = 0
			lowest.Stdout = ""
			lowest.Stderr = ""
			s.store.Save(lowest)
			s.queue.Push(lowest)
			slog.Info("preempted task re-queued", "task_id", lowest.ID)
		})
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.store.Close()
}

// GetStore 返回底层存储（供 Syncer 等外部组件使用）
func (s *Scheduler) GetStore() *store.SQLiteStore {
	return s.store
}

// GetTask 获取任务
func (s *Scheduler) GetTask(id string) (*models.Task, error) {
	return s.store.Get(id)
}

// CancelTask 取消任务
func (s *Scheduler) CancelTask(id string) bool {
	// 先尝试从队列移除
	if s.queue.Remove(id) {
		return true
	}
	// 再尝试取消正在执行的任务
	return s.executor.Cancel(id)
}

// QueueStatus 队列状态
type QueueStatus struct {
	Pending  int `json:"pending"`
	Running  int `json:"running"`
	Capacity int `json:"capacity"`
}

// GetQueueStatus 获取队列状态
func (s *Scheduler) GetQueueStatus() *QueueStatus {
	return &QueueStatus{
		Pending:  s.queue.Len(),
		Running:  s.executor.RunningCount(),
		Capacity: 4,
	}
}
