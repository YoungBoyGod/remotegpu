package queue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

func newTask(id string, priority int) *models.Task {
	return &models.Task{
		ID:        id,
		Priority:  priority,
		CreatedAt: time.Now(),
	}
}

func TestPushPop(t *testing.T) {
	m := NewManager()

	m.Push(newTask("t1", 5))
	m.Push(newTask("t2", 1))
	m.Push(newTask("t3", 3))

	if m.Len() != 3 {
		t.Fatalf("expected len 3, got %d", m.Len())
	}

	// 应按优先级顺序弹出: t2(1) -> t3(3) -> t1(5)
	got := m.Pop()
	if got.ID != "t2" {
		t.Errorf("expected t2, got %s", got.ID)
	}
	got = m.Pop()
	if got.ID != "t3" {
		t.Errorf("expected t3, got %s", got.ID)
	}
	got = m.Pop()
	if got.ID != "t1" {
		t.Errorf("expected t1, got %s", got.ID)
	}

	if m.Pop() != nil {
		t.Error("expected nil from empty queue")
	}
}

func TestSamePriorityOrderByTime(t *testing.T) {
	m := NewManager()

	t1 := &models.Task{ID: "t1", Priority: 5, CreatedAt: time.Now()}
	time.Sleep(time.Millisecond)
	t2 := &models.Task{ID: "t2", Priority: 5, CreatedAt: time.Now()}

	m.Push(t2)
	m.Push(t1)

	got := m.Pop()
	if got.ID != "t1" {
		t.Errorf("expected t1 (earlier), got %s", got.ID)
	}
}

func TestDuplicatePush(t *testing.T) {
	m := NewManager()
	task := newTask("t1", 5)

	m.Push(task)
	m.Push(task) // 重复 push

	if m.Len() != 1 {
		t.Fatalf("expected len 1 after duplicate push, got %d", m.Len())
	}
}

func TestRemove(t *testing.T) {
	m := NewManager()
	m.Push(newTask("t1", 5))
	m.Push(newTask("t2", 3))

	if !m.Remove("t1") {
		t.Error("expected Remove to return true")
	}
	if m.Len() != 1 {
		t.Fatalf("expected len 1, got %d", m.Len())
	}
	if m.Remove("t1") {
		t.Error("expected Remove to return false for missing task")
	}

	got := m.Pop()
	if got.ID != "t2" {
		t.Errorf("expected t2, got %s", got.ID)
	}
}

func TestPeekAndGet(t *testing.T) {
	m := NewManager()

	if m.Peek() != nil {
		t.Error("expected nil Peek on empty queue")
	}

	m.Push(newTask("t1", 5))
	m.Push(newTask("t2", 1))

	peeked := m.Peek()
	if peeked.ID != "t2" {
		t.Errorf("expected Peek t2, got %s", peeked.ID)
	}
	// Peek 不应移除
	if m.Len() != 2 {
		t.Error("Peek should not remove task")
	}

	got := m.Get("t1")
	if got == nil || got.ID != "t1" {
		t.Error("Get t1 failed")
	}
	if m.Get("nonexistent") != nil {
		t.Error("Get nonexistent should return nil")
	}
}

func TestNotifyChan(t *testing.T) {
	m := NewManager()
	ch := m.NotifyChan()

	m.Push(newTask("t1", 5))

	select {
	case <-ch:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Error("expected notification after Push")
	}
}

func TestConcurrentAccess(t *testing.T) {
	m := NewManager()
	var wg sync.WaitGroup

	// 并发 push
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.Push(&models.Task{
				ID:        fmt.Sprintf("t%d", n),
				Priority:  n % 10,
				CreatedAt: time.Now(),
			})
		}(i)
	}
	wg.Wait()

	if m.Len() != 100 {
		t.Fatalf("expected 100 tasks, got %d", m.Len())
	}

	// 并发 pop
	var popped int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if m.Pop() != nil {
				atomic.AddInt64(&popped, 1)
			}
		}()
	}
	wg.Wait()

	if popped != 100 {
		t.Fatalf("expected 100 pops, got %d", popped)
	}
}
