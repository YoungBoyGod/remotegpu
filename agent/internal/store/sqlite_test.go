package store

import (
	"os"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

func tempStore(t *testing.T) *SQLiteStore {
	t.Helper()
	f, err := os.CreateTemp("", "agent-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	st, err := NewSQLiteStore(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

func TestSaveAndGet(t *testing.T) {
	st := tempStore(t)

	task := &models.Task{
		ID:      "t1",
		Command: "echo hello",
		Status:  models.TaskStatusPending,
	}
	if err := st.Save(task); err != nil {
		t.Fatal(err)
	}

	got, err := st.Get("t1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "t1" || got.Command != "echo hello" {
		t.Errorf("unexpected task: %+v", got)
	}
}

func TestGetNotFound(t *testing.T) {
	st := tempStore(t)
	_, err := st.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestListByStatus(t *testing.T) {
	st := tempStore(t)

	st.Save(&models.Task{ID: "t1", Command: "a", Status: models.TaskStatusPending})
	st.Save(&models.Task{ID: "t2", Command: "b", Status: models.TaskStatusRunning})
	st.Save(&models.Task{ID: "t3", Command: "c", Status: models.TaskStatusPending})

	tasks, err := st.ListByStatus(models.TaskStatusPending)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 pending tasks, got %d", len(tasks))
	}
}

func TestListUnsyncedOnlyCompleted(t *testing.T) {
	st := tempStore(t)

	st.Save(&models.Task{ID: "t1", Command: "a", Status: models.TaskStatusCompleted, Synced: false})
	st.Save(&models.Task{ID: "t2", Command: "b", Status: models.TaskStatusFailed, Synced: false})
	st.Save(&models.Task{ID: "t3", Command: "c", Status: models.TaskStatusPending, Synced: false})
	st.Save(&models.Task{ID: "t4", Command: "d", Status: models.TaskStatusCompleted, Synced: true})

	tasks, err := st.ListUnsynced()
	if err != nil {
		t.Fatal(err)
	}
	// 只返回 completed/failed 且 synced=0 的任务
	if len(tasks) != 2 {
		t.Fatalf("expected 2 unsynced tasks, got %d", len(tasks))
	}
}

func TestMarkSynced(t *testing.T) {
	st := tempStore(t)

	st.Save(&models.Task{ID: "t1", Command: "a", Status: models.TaskStatusCompleted, Synced: false})
	st.MarkSynced("t1")

	tasks, err := st.ListUnsynced()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 unsynced after MarkSynced, got %d", len(tasks))
	}
}

func TestDelete(t *testing.T) {
	st := tempStore(t)

	st.Save(&models.Task{ID: "t1", Command: "a", Status: models.TaskStatusPending})
	st.Delete("t1")

	_, err := st.Get("t1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestSaveAndGetWithTimestamps(t *testing.T) {
	st := tempStore(t)

	now := time.Now().Truncate(time.Second)
	task := &models.Task{
		ID:        "ts1",
		Command:   "echo ts",
		Status:    models.TaskStatusCompleted,
		CreatedAt: now,
		StartedAt: now.Add(1 * time.Second),
		EndedAt:   now.Add(5 * time.Second),
	}
	st.Save(task)

	got, err := st.Get("ts1")
	if err != nil {
		t.Fatal(err)
	}
	if got.CreatedAt.Unix() != now.Unix() {
		t.Errorf("CreatedAt 不匹配: 期望 %v，实际 %v", now, got.CreatedAt)
	}
	if got.StartedAt.IsZero() {
		t.Error("StartedAt 不应为零值")
	}
	if got.EndedAt.IsZero() {
		t.Error("EndedAt 不应为零值")
	}
}

func TestSaveAndGetWithArgsAndEnv(t *testing.T) {
	st := tempStore(t)

	task := &models.Task{
		ID:      "ae1",
		Command: "echo",
		Args:    []string{"arg1", "arg2", "arg3"},
		Env:     map[string]string{"KEY1": "val1", "KEY2": "val2"},
		Status:  models.TaskStatusPending,
	}
	st.Save(task)

	got, err := st.Get("ae1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Args) != 3 || got.Args[0] != "arg1" {
		t.Errorf("Args 反序列化不正确: %v", got.Args)
	}
	if len(got.Env) != 2 || got.Env["KEY1"] != "val1" {
		t.Errorf("Env 反序列化不正确: %v", got.Env)
	}
}

func TestSaveUpdatesExisting(t *testing.T) {
	st := tempStore(t)

	// 第一次保存
	task := &models.Task{
		ID:      "u1",
		Command: "echo v1",
		Status:  models.TaskStatusPending,
	}
	st.Save(task)

	// 第二次保存（更新）
	task.Status = models.TaskStatusCompleted
	task.Stdout = "output"
	task.ExitCode = 0
	st.Save(task)

	got, err := st.Get("u1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != models.TaskStatusCompleted {
		t.Errorf("更新后状态应为 completed，实际为 %s", got.Status)
	}
	if got.Stdout != "output" {
		t.Errorf("更新后 Stdout 应为 output，实际为 %q", got.Stdout)
	}
}
