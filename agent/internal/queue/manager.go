package queue

import (
	"container/heap"
	"sync"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

// Manager 优先级队列管理器
type Manager struct {
	mu       sync.RWMutex
	tasks    taskHeap
	taskMap  map[string]*models.Task // 用于快速查找
	notifyCh chan struct{}           // 通知有新任务
}

// NewManager 创建队列管理器
func NewManager() *Manager {
	m := &Manager{
		tasks:    make(taskHeap, 0),
		taskMap:  make(map[string]*models.Task),
		notifyCh: make(chan struct{}, 1),
	}
	heap.Init(&m.tasks)
	return m
}

// Push 添加任务到队列
func (m *Manager) Push(task *models.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.taskMap[task.ID]; exists {
		return // 任务已存在
	}

	heap.Push(&m.tasks, task)
	m.taskMap[task.ID] = task

	// 非阻塞通知
	select {
	case m.notifyCh <- struct{}{}:
	default:
	}
}

// Pop 取出最高优先级任务
func (m *Manager) Pop() *models.Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.tasks.Len() == 0 {
		return nil
	}

	task := heap.Pop(&m.tasks).(*models.Task)
	delete(m.taskMap, task.ID)
	return task
}

// Peek 查看最高优先级任务但不移除
func (m *Manager) Peek() *models.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.tasks.Len() == 0 {
		return nil
	}
	return m.tasks[0]
}

// Remove 从队列中移除指定任务
func (m *Manager) Remove(taskID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.taskMap[taskID]; !exists {
		return false
	}

	// 重建堆（简单实现）
	newTasks := make(taskHeap, 0, len(m.tasks)-1)
	for _, t := range m.tasks {
		if t.ID != taskID {
			newTasks = append(newTasks, t)
		}
	}
	m.tasks = newTasks
	heap.Init(&m.tasks)
	delete(m.taskMap, taskID)
	return true
}

// Len 返回队列长度
func (m *Manager) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tasks.Len()
}

// Get 获取指定任务
func (m *Manager) Get(taskID string) *models.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.taskMap[taskID]
}

// NotifyChan 返回通知通道
func (m *Manager) NotifyChan() <-chan struct{} {
	return m.notifyCh
}

// List 返回所有任务（按优先级排序）
func (m *Manager) List() []*models.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*models.Task, len(m.tasks))
	copy(result, m.tasks)
	return result
}
