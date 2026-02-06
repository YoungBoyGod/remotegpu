package executor

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

const maxOutputSize = 1 << 20 // 1MB

// limitedWriter 限制写入大小的 Writer
type limitedWriter struct {
	buf     bytes.Buffer
	limit   int
	dropped bool
}

func (w *limitedWriter) Write(p []byte) (int, error) {
	if w.buf.Len() >= w.limit {
		w.dropped = true
		return len(p), nil // 丢弃但不报错，避免中断进程
	}
	remaining := w.limit - w.buf.Len()
	if len(p) > remaining {
		w.dropped = true
		p = p[:remaining]
	}
	return w.buf.Write(p)
}

func (w *limitedWriter) String() string {
	s := w.buf.String()
	if w.dropped {
		s += "\n...[truncated, output exceeded 1MB limit]"
	}
	return s
}

// Executor 任务执行器
type Executor struct {
	mu         sync.Mutex
	running    map[string]*runningTask
	maxWorkers int
}

type runningTask struct {
	task   *models.Task
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

// NewExecutor 创建执行器
func NewExecutor(maxWorkers int) *Executor {
	return &Executor{
		running:    make(map[string]*runningTask),
		maxWorkers: maxWorkers,
	}
}

// RunningCount 返回正在运行的任务数
func (e *Executor) RunningCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.running)
}

// CanAccept 检查是否可以接受新任务
func (e *Executor) CanAccept() bool {
	return e.RunningCount() < e.maxWorkers
}

// Execute 执行任务
func (e *Executor) Execute(task *models.Task) {
	timeout := task.Timeout
	if timeout <= 0 {
		timeout = 3600
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	var cmd *exec.Cmd
	if len(task.Args) > 0 {
		cmd = exec.CommandContext(ctx, task.Command, task.Args...)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", task.Command)
	}

	if task.WorkDir != "" {
		cmd.Dir = task.WorkDir
	}

	// 设置环境变量
	if len(task.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range task.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	// 设置进程组，便于杀死子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout := &limitedWriter{limit: maxOutputSize}
	stderr := &limitedWriter{limit: maxOutputSize}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// 记录运行中的任务
	e.mu.Lock()
	e.running[task.ID] = &runningTask{task: task, cmd: cmd, cancel: cancel}
	e.mu.Unlock()

	// 更新任务状态
	task.Status = models.TaskStatusRunning
	task.StartedAt = time.Now()

	// 执行命令
	err := cmd.Run()

	// 清理
	e.mu.Lock()
	delete(e.running, task.ID)
	e.mu.Unlock()
	cancel()

	// 更新结果
	task.EndedAt = time.Now()
	task.Stdout = stdout.String()
	task.Stderr = stderr.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			task.ExitCode = exitErr.ExitCode()
		} else {
			task.ExitCode = -1
		}
		task.Error = err.Error()
		task.Status = models.TaskStatusFailed
	} else {
		task.ExitCode = 0
		task.Status = models.TaskStatusCompleted
	}
}

// Cancel 取消任务
func (e *Executor) Cancel(taskID string) bool {
	e.mu.Lock()
	rt, exists := e.running[taskID]
	e.mu.Unlock()

	if !exists {
		return false
	}

	// 先发送 SIGTERM
	if rt.cmd.Process != nil {
		syscall.Kill(-rt.cmd.Process.Pid, syscall.SIGTERM)
	}

	// 等待一段时间后强制杀死
	go func() {
		time.Sleep(5 * time.Second)
		e.mu.Lock()
		if _, still := e.running[taskID]; still && rt.cmd.Process != nil {
			syscall.Kill(-rt.cmd.Process.Pid, syscall.SIGKILL)
		}
		e.mu.Unlock()
	}()

	return true
}
