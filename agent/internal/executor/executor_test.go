package executor

import (
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

func TestExecuteSuccess(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "t1",
		Command: "echo hello",
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusCompleted {
		t.Errorf("expected completed, got %s", task.Status)
	}
	if task.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", task.ExitCode)
	}
	if task.Stdout == "" {
		t.Error("expected non-empty stdout")
	}
}

func TestExecuteFailure(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "t2",
		Command: "false",
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusFailed {
		t.Errorf("expected failed, got %s", task.Status)
	}
	if task.ExitCode == 0 {
		t.Error("expected non-zero exit code")
	}
}

func TestExecuteTimeout(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "t3",
		Command: "sleep 60",
		Timeout: 1,
	}
	start := time.Now()
	e.Execute(task)
	elapsed := time.Since(start)

	if task.Status != models.TaskStatusFailed {
		t.Errorf("expected failed, got %s", task.Status)
	}
	if elapsed > 5*time.Second {
		t.Errorf("timeout took too long: %v", elapsed)
	}
}
