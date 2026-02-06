package models

import "time"

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusAssigned  TaskStatus = "assigned"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
	TaskStatusPreempted TaskStatus = "preempted"
	TaskStatusSuspended TaskStatus = "suspended"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeShell  TaskType = "shell"
	TaskTypePython TaskType = "python"
	TaskTypeScript TaskType = "script"
)

// Task 任务模型（按设计文档 2.1 节定义）
type Task struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Type    TaskType          `json:"type"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	WorkDir string            `json:"workdir"`
	Env     map[string]string `json:"env"`
	Timeout int               `json:"timeout"`

	// 优先级和重试
	Priority   int `json:"priority"`
	RetryCount int `json:"retry_count"`
	RetryDelay int `json:"retry_delay"`
	MaxRetries int `json:"max_retries"`

	// 状态相关
	Status   TaskStatus `json:"status"`
	ExitCode int        `json:"exit_code"`
	Stdout   string     `json:"stdout"`
	Stderr   string     `json:"stderr"`
	Error    string     `json:"error"`

	// 时间戳
	CreatedAt  time.Time `json:"created_at"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	EndedAt    time.Time `json:"ended_at,omitempty"`

	// 关联
	MachineID string   `json:"machine_id"`
	GroupID   string   `json:"group_id,omitempty"`
	ParentID  string   `json:"parent_id,omitempty"`
	DependsOn []string `json:"depends_on,omitempty"`

	// 调度与租约
	AssignedAgentID string    `json:"assigned_agent_id"`
	LeaseExpiresAt  time.Time `json:"lease_expires_at,omitempty"`
	AttemptID       string    `json:"attempt_id"`

	// 本地同步标记
	Synced bool `json:"synced"`
}
