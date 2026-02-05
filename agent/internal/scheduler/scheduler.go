package scheduler

import (
	"log"
	"sync"
	"time"

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

// Start 启动调度器
func (s *Scheduler) Start() error {
	// 恢复未完成的任务
	if err := s.recover(); err != nil {
		log.Printf("recover tasks error: %v", err)
	}

	// 启动调度循环
	s.wg.Add(1)
	go s.scheduleLoop()

	return nil
}

// recover 恢复未完成的任务
func (s *Scheduler) recover() error {
	// 恢复 pending 状态的任务到队列
	tasks, err := s.store.ListByStatus(models.TaskStatusPending)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		s.queue.Push(task)
	}
	log.Printf("recovered %d pending tasks", len(tasks))
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

		// 异步执行
		go func(t *models.Task) {
			s.executor.Execute(t)
			// 保存结果
			if err := s.store.Save(t); err != nil {
				log.Printf("save task result error: %v", err)
			}
		}(task)
	}
}

// Submit 提交任务
func (s *Scheduler) Submit(task *models.Task) error {
	task.Status = models.TaskStatusPending
	task.CreatedAt = time.Now()

	if err := s.store.Save(task); err != nil {
		return err
	}

	s.queue.Push(task)
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.store.Close()
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
