package store

import (
	"os"
	"testing"

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
