package agent

import "time"

// Request 通用请求结构
type Request struct {
	HostID    string            `json:"host_id"`
	Action    string            `json:"action"`
	Params    map[string]string `json:"params,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Response 通用响应结构
type Response struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// StopProcessRequest 停止进程请求
type StopProcessRequest struct {
	HostID    string `json:"host_id"`
	ProcessID int    `json:"process_id"`
	Signal    string `json:"signal,omitempty"` // SIGTERM, SIGKILL
}

// ResetSSHRequest 重置SSH请求
type ResetSSHRequest struct {
	HostID    string `json:"host_id"`
	PublicKey string `json:"public_key"`
	Username  string `json:"username,omitempty"`
}

// CleanupRequest 清理机器请求
type CleanupRequest struct {
	HostID       string   `json:"host_id"`
	CleanupTypes []string `json:"cleanup_types"` // process, data, ssh, docker
}

// MountDatasetRequest 挂载数据集请求
type MountDatasetRequest struct {
	HostID     string `json:"host_id"`
	DatasetID  uint   `json:"dataset_id"`
	SourcePath string `json:"source_path"`
	MountPoint string `json:"mount_point"`
	ReadOnly   bool   `json:"read_only"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	Hostname    string    `json:"hostname"`
	OS          string    `json:"os"`
	Kernel      string    `json:"kernel"`
	CPUCores    int       `json:"cpu_cores"`
	MemoryTotal uint64    `json:"memory_total"`
	MemoryFree  uint64    `json:"memory_free"`
	DiskTotal   uint64    `json:"disk_total"`
	DiskFree    uint64    `json:"disk_free"`
	GPUCount    int       `json:"gpu_count"`
	GPUInfo     []GPUInfo `json:"gpu_info,omitempty"`
	Uptime      int64     `json:"uptime"`
}

// GPUInfo GPU信息
type GPUInfo struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	MemoryTotal uint64 `json:"memory_total"`
	MemoryUsed  uint64 `json:"memory_used"`
	Utilization int    `json:"utilization"`
	Temperature int    `json:"temperature"`
}

// ExecuteCommandRequest 执行命令请求
type ExecuteCommandRequest struct {
	HostID  string `json:"host_id"`
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}

// ExecuteCommandResponse 执行命令响应
type ExecuteCommandResponse struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}
