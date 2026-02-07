package scheduler

import (
	"os"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

// tempScheduler 创建使用临时数据库的调度器
func tempScheduler(t *testing.T, maxWorkers int) *Scheduler {
	t.Helper()
	f, err := os.CreateTemp("", "scheduler-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	s, err := NewScheduler(f.Name(), maxWorkers)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Stop() })
	return s
}

func TestNewScheduler(t *testing.T) {
	s := tempScheduler(t, 2)
	if s.queue == nil || s.store == nil || s.executor == nil {
		t.Error("调度器组件未正确初始化")
	}
}

func TestSubmitAndGetTask(t *testing.T) {
	s := tempScheduler(t, 2)

	task := &models.Task{
		ID:      "t1",
		Command: "echo hello",
	}
	if err := s.Submit(task); err != nil {
		t.Fatalf("Submit 失败: %v", err)
	}

	// 任务应被保存到 store
	got, err := s.GetTask("t1")
	if err != nil {
		t.Fatalf("GetTask 失败: %v", err)
	}
	if got.ID != "t1" {
		t.Errorf("期望 ID=t1，实际为 %s", got.ID)
	}
	if got.Status != models.TaskStatusPending {
		t.Errorf("本地提交的任务状态应为 pending，实际为 %s", got.Status)
	}
}

func TestSubmitSetsDefaults(t *testing.T) {
	s := tempScheduler(t, 2)

	task := &models.Task{
		ID:      "t1",
		Command: "echo test",
	}
	s.Submit(task)

	if task.Status != models.TaskStatusPending {
		t.Errorf("默认状态应为 pending，实际为 %s", task.Status)
	}
	if task.CreatedAt.IsZero() {
		t.Error("CreatedAt 应被自动设置")
	}
}

func TestSubmitServerTask(t *testing.T) {
	s := tempScheduler(t, 2)

	// Server 下发的任务有 AttemptID，不应覆盖状态
	task := &models.Task{
		ID:        "t1",
		Command:   "echo test",
		AttemptID: "attempt-1",
		Status:    models.TaskStatusAssigned,
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}
	s.Submit(task)

	if task.Status != models.TaskStatusAssigned {
		t.Errorf("Server 任务状态不应被覆盖，期望 assigned，实际为 %s", task.Status)
	}
}

func TestCancelTaskFromQueue(t *testing.T) {
	s := tempScheduler(t, 0) // maxWorkers=0，任务不会被执行

	task := &models.Task{
		ID:      "t1",
		Command: "sleep 60",
	}
	s.Submit(task)

	// 任务在队列中，应能取消
	if !s.CancelTask("t1") {
		t.Error("CancelTask 应返回 true（任务在队列中）")
	}

	// 取消不存在的任务
	if s.CancelTask("nonexistent") {
		t.Error("CancelTask 不存在的任务应返回 false")
	}
}

func TestGetQueueStatus(t *testing.T) {
	s := tempScheduler(t, 2)

	status := s.GetQueueStatus()
	if status.Pending != 0 {
		t.Errorf("初始 Pending 应为 0，实际为 %d", status.Pending)
	}
	if status.Running != 0 {
		t.Errorf("初始 Running 应为 0，实际为 %d", status.Running)
	}

	// 提交任务后检查
	s.Submit(&models.Task{ID: "t1", Command: "echo 1"})
	s.Submit(&models.Task{ID: "t2", Command: "echo 2"})

	status = s.GetQueueStatus()
	if status.Pending < 0 {
		t.Errorf("Pending 不应为负数: %d", status.Pending)
	}
}

func TestRecoverPendingTasks(t *testing.T) {
	// 创建临时数据库并预存 pending 任务
	f, err := os.CreateTemp("", "recover-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	// 第一阶段：创建调度器并保存 pending 任务
	s1, err := NewScheduler(f.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}
	s1.store.Save(&models.Task{
		ID:      "t1",
		Command: "echo 1",
		Status:  models.TaskStatusPending,
	})
	s1.store.Save(&models.Task{
		ID:      "t2",
		Command: "echo 2",
		Status:  models.TaskStatusPending,
	})
	s1.Stop()

	// 第二阶段：新调度器 recover 应恢复 pending 任务到队列
	s2, err := NewScheduler(f.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Stop()

	if err := s2.recover(); err != nil {
		t.Fatalf("recover 失败: %v", err)
	}
	if s2.queue.Len() != 2 {
		t.Errorf("recover 后队列应有 2 个任务，实际为 %d", s2.queue.Len())
	}
}

func TestRecoverRunningTasksMarkedFailed(t *testing.T) {
	f, err := os.CreateTemp("", "recover-running-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	// 预存 running 任务（模拟 Agent 崩溃后遗留）
	s1, err := NewScheduler(f.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}
	s1.store.Save(&models.Task{
		ID:      "r1",
		Command: "sleep 100",
		Status:  models.TaskStatusRunning,
	})
	s1.Stop()

	// 新调度器 recover 应将 running 标记为 failed
	s2, err := NewScheduler(f.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Stop()

	s2.recover()

	got, err := s2.store.Get("r1")
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}
	if got.Status != models.TaskStatusFailed {
		t.Errorf("running 任务 recover 后应为 failed，实际为 %s", got.Status)
	}
	if got.Error == "" {
		t.Error("recover 后 Error 字段应包含原因")
	}
}

func TestDependenciesMetNoDeps(t *testing.T) {
	s := tempScheduler(t, 2)

	task := &models.Task{ID: "t1", Command: "echo 1"}
	if !s.dependenciesMet(task) {
		t.Error("无依赖的任务 dependenciesMet 应返回 true")
	}
}

func TestDependenciesMetAllCompleted(t *testing.T) {
	s := tempScheduler(t, 2)

	// 保存已完成的依赖任务
	s.store.Save(&models.Task{
		ID:      "dep1",
		Command: "echo dep1",
		Status:  models.TaskStatusCompleted,
	})
	s.store.Save(&models.Task{
		ID:      "dep2",
		Command: "echo dep2",
		Status:  models.TaskStatusCompleted,
	})

	task := &models.Task{
		ID:        "t1",
		Command:   "echo 1",
		DependsOn: []string{"dep1", "dep2"},
	}
	if !s.dependenciesMet(task) {
		t.Error("所有依赖已完成，dependenciesMet 应返回 true")
	}
}

func TestDependenciesMetNotCompleted(t *testing.T) {
	s := tempScheduler(t, 2)

	s.store.Save(&models.Task{
		ID:      "dep1",
		Command: "echo dep1",
		Status:  models.TaskStatusCompleted,
	})
	s.store.Save(&models.Task{
		ID:      "dep2",
		Command: "echo dep2",
		Status:  models.TaskStatusRunning, // 未完成
	})

	task := &models.Task{
		ID:        "t1",
		Command:   "echo 1",
		DependsOn: []string{"dep1", "dep2"},
	}
	if s.dependenciesMet(task) {
		t.Error("有依赖未完成，dependenciesMet 应返回 false")
	}
}

func TestDependenciesMetMissing(t *testing.T) {
	s := tempScheduler(t, 2)

	// 依赖任务不存在
	task := &models.Task{
		ID:        "t1",
		Command:   "echo 1",
		DependsOn: []string{"nonexistent"},
	}
	if s.dependenciesMet(task) {
		t.Error("依赖任务不存在，dependenciesMet 应返回 false")
	}
}

func TestStartExecutesTask(t *testing.T) {
	s := tempScheduler(t, 2)
	s.Start()

	task := &models.Task{
		ID:      "exec1",
		Command: "echo integration",
		Timeout: 5,
	}
	s.Submit(task)

	// 等待任务执行完成
	deadline := time.After(5 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("任务执行超时")
		default:
			got, err := s.store.Get("exec1")
			if err == nil && (got.Status == models.TaskStatusCompleted || got.Status == models.TaskStatusFailed) {
				if got.Status != models.TaskStatusCompleted {
					t.Errorf("任务应为 completed，实际为 %s, error: %s", got.Status, got.Error)
				}
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TestRecoverAssignedLeaseExpired(t *testing.T) {
	f, err := os.CreateTemp("", "recover-assigned-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	s1, err := NewScheduler(f.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}
	// 保存一个租约已过期的 assigned 任务
	s1.store.Save(&models.Task{
		ID:             "a1",
		Command:        "echo assigned",
		Status:         models.TaskStatusAssigned,
		LeaseExpiresAt: time.Now().Add(-1 * time.Hour),
	})
	s1.Stop()

	s2, err := NewScheduler(f.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Stop()

	s2.recover()

	got, err := s2.store.Get("a1")
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}
	if got.Status != models.TaskStatusFailed {
		t.Errorf("租约过期的 assigned 任务应标记为 failed，实际为 %s", got.Status)
	}
}
