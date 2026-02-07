package executor

import (
	"strings"
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
	"github.com/YoungBoyGod/remotegpu-agent/internal/security"
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

func TestLimitedWriterWithinLimit(t *testing.T) {
	w := &limitedWriter{limit: 100}
	data := []byte("hello world")
	n, err := w.Write(data)
	if err != nil {
		t.Fatalf("Write 失败: %v", err)
	}
	if n != len(data) {
		t.Errorf("写入字节数应为 %d，实际为 %d", len(data), n)
	}
	if w.dropped {
		t.Error("未超限时 dropped 应为 false")
	}
	if !strings.Contains(w.String(), "hello world") {
		t.Errorf("输出应包含 hello world，实际为 %q", w.String())
	}
}

func TestLimitedWriterExceedsLimit(t *testing.T) {
	w := &limitedWriter{limit: 10}
	data := []byte("this is a long string that exceeds the limit")
	n, err := w.Write(data)
	if err != nil {
		t.Fatalf("Write 不应返回错误: %v", err)
	}
	// 首次截断时返回实际写入的字节数（remaining），后续满后才返回 len(p)
	if n != 10 {
		t.Errorf("首次截断应返回 remaining 字节数 10，实际为 %d", n)
	}
	if !w.dropped {
		t.Error("超限后 dropped 应为 true")
	}
	output := w.String()
	if !strings.Contains(output, "truncated") {
		t.Errorf("超限输出应包含 truncated 标记，实际为 %q", output)
	}
}

func TestLimitedWriterMultipleWrites(t *testing.T) {
	w := &limitedWriter{limit: 10}
	w.Write([]byte("12345"))
	w.Write([]byte("67890"))
	// 已满，再写应丢弃
	w.Write([]byte("extra"))

	if !w.dropped {
		t.Error("超限后 dropped 应为 true")
	}
	if w.buf.Len() != 10 {
		t.Errorf("缓冲区应为 10 字节，实际为 %d", w.buf.Len())
	}
}

func TestValidatorRejectsCommand(t *testing.T) {
	e := NewExecutor(2)
	v := security.NewValidator([]string{"echo", "ls"}, nil)
	e.SetValidator(v)

	task := &models.Task{
		ID:      "v1",
		Command: "rm",
		Args:    []string{"-rf", "/"},
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusFailed {
		t.Errorf("被拒绝的命令应为 failed，实际为 %s", task.Status)
	}
	if !strings.Contains(task.Error, "rejected") {
		t.Errorf("错误信息应包含 rejected，实际为 %q", task.Error)
	}
}

func TestValidatorAllowsCommand(t *testing.T) {
	e := NewExecutor(2)
	v := security.NewValidator([]string{"echo"}, nil)
	e.SetValidator(v)

	task := &models.Task{
		ID:      "v2",
		Command: "echo",
		Args:    []string{"allowed"},
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusCompleted {
		t.Errorf("白名单命令应为 completed，实际为 %s, error: %s", task.Status, task.Error)
	}
}

func TestExecuteWithEnv(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "env1",
		Command: "echo $MY_VAR",
		Env:     map[string]string{"MY_VAR": "test_value"},
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusCompleted {
		t.Fatalf("任务应为 completed，实际为 %s", task.Status)
	}
	if !strings.Contains(task.Stdout, "test_value") {
		t.Errorf("stdout 应包含环境变量值 test_value，实际为 %q", task.Stdout)
	}
}

func TestExecuteWithWorkDir(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "wd1",
		Command: "pwd",
		WorkDir: "/tmp",
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusCompleted {
		t.Fatalf("任务应为 completed，实际为 %s", task.Status)
	}
	if !strings.Contains(task.Stdout, "/tmp") {
		t.Errorf("stdout 应包含 /tmp，实际为 %q", task.Stdout)
	}
}

func TestCanAcceptAndRunningCount(t *testing.T) {
	e := NewExecutor(1)

	if !e.CanAccept() {
		t.Error("空执行器应能接受任务")
	}
	if e.RunningCount() != 0 {
		t.Errorf("初始 RunningCount 应为 0，实际为 %d", e.RunningCount())
	}
}

func TestExecuteWithArgs(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "args1",
		Command: "echo",
		Args:    []string{"hello", "world"},
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusCompleted {
		t.Fatalf("任务应为 completed，实际为 %s", task.Status)
	}
	if !strings.Contains(task.Stdout, "hello world") {
		t.Errorf("stdout 应包含 hello world，实际为 %q", task.Stdout)
	}
}

func TestExecuteDefaultTimeout(t *testing.T) {
	e := NewExecutor(2)

	// Timeout=0 应使用默认值 3600 秒，不会立即超时
	task := &models.Task{
		ID:      "dt1",
		Command: "echo default_timeout",
		Timeout: 0,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusCompleted {
		t.Errorf("默认超时任务应为 completed，实际为 %s", task.Status)
	}
}

func TestExecuteSetsTimestamps(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "ts1",
		Command: "echo timestamps",
		Timeout: 10,
	}
	e.Execute(task)

	if task.StartedAt.IsZero() {
		t.Error("StartedAt 应被设置")
	}
	if task.EndedAt.IsZero() {
		t.Error("EndedAt 应被设置")
	}
	if task.EndedAt.Before(task.StartedAt) {
		t.Error("EndedAt 应在 StartedAt 之后")
	}
}

func TestLowestPriorityRunningEmpty(t *testing.T) {
	e := NewExecutor(2)

	if e.LowestPriorityRunning() != nil {
		t.Error("无运行任务时 LowestPriorityRunning 应返回 nil")
	}
}

func TestCancelNonexistent(t *testing.T) {
	e := NewExecutor(2)

	if e.Cancel("nonexistent") {
		t.Error("取消不存在的任务应返回 false")
	}
}

func TestBlockedPatternValidator(t *testing.T) {
	e := NewExecutor(2)
	v := security.NewValidator(nil, []string{"rm -rf", "mkfs"})
	e.SetValidator(v)

	task := &models.Task{
		ID:      "bp1",
		Command: "rm",
		Args:    []string{"-rf", "/tmp/test"},
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusFailed {
		t.Errorf("匹配黑名单的命令应为 failed，实际为 %s", task.Status)
	}
}

func TestExecuteStderr(t *testing.T) {
	e := NewExecutor(2)

	task := &models.Task{
		ID:      "se1",
		Command: "echo error_msg >&2 && false",
		Timeout: 10,
	}
	e.Execute(task)

	if task.Status != models.TaskStatusFailed {
		t.Errorf("任务应为 failed，实际为 %s", task.Status)
	}
	if !strings.Contains(task.Stderr, "error_msg") {
		t.Errorf("stderr 应包含 error_msg，实际为 %q", task.Stderr)
	}
}
