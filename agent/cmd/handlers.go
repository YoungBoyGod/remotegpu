package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	agentErrors "github.com/YoungBoyGod/remotegpu-agent/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func handlePing(c *gin.Context) {
	respondSuccess(c, gin.H{"ok": true})
}

func handleSystemInfo(c *gin.Context) {
	hostInfo, _ := host.Info()
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")
	cpuCount, _ := cpu.Counts(true)

	c.JSON(http.StatusOK, gin.H{
		"hostname":     hostInfo.Hostname,
		"os":           runtime.GOOS,
		"kernel":       hostInfo.KernelVersion,
		"cpu_cores":    cpuCount,
		"memory_total": memInfo.Total,
		"memory_free":  memInfo.Free,
		"disk_total":   diskInfo.Total,
		"disk_free":    diskInfo.Free,
		"uptime":       hostInfo.Uptime,
	})
}

func handleStopProcess(c *gin.Context) {
	var req struct {
		ProcessID int    `json:"process_id"`
		Signal    string `json:"signal"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondErrorCode(c, http.StatusBadRequest, agentErrors.ErrInvalidParams)
		return
	}

	sig := syscall.SIGTERM
	if req.Signal == "SIGKILL" {
		sig = syscall.SIGKILL
	}

	process, err := os.FindProcess(req.ProcessID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, agentErrors.ErrInternal, err.Error())
		return
	}

	if err := process.Signal(sig); err != nil {
		respondError(c, http.StatusInternalServerError, agentErrors.ErrInternal, err.Error())
		return
	}

	respondSuccess(c, nil)
}

func handleResetSSH(c *gin.Context) {
	cmd := exec.Command("bash", "-c", "rm -f ~/.ssh/authorized_keys")
	if err := cmd.Run(); err != nil {
		respondError(c, http.StatusInternalServerError, agentErrors.ErrInternal, err.Error())
		return
	}
	respondSuccess(c, nil)
}

func handleCleanup(c *gin.Context) {
	var req struct {
		CleanupTypes []string `json:"cleanup_types"`
	}
	c.ShouldBindJSON(&req)

	for _, t := range req.CleanupTypes {
		switch strings.ToLower(t) {
		case "docker":
			exec.Command("docker", "system", "prune", "-af").Run()
		case "ssh":
			exec.Command("bash", "-c", "rm -f ~/.ssh/authorized_keys").Run()
		}
	}
	respondSuccess(c, nil)
}

// handleExecCommand 执行 shell 命令
func handleExecCommand(c *gin.Context) {
	var req struct {
		Command string `json:"command" binding:"required"`
		Timeout int    `json:"timeout"`
		WorkDir string `json:"workdir"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondErrorCode(c, http.StatusBadRequest, agentErrors.ErrInvalidParams)
		return
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 60
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", req.Command)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	data := gin.H{
		"exit_code": exitCode,
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
	}
	if err != nil {
		respondWithData(c, false, agentErrors.ErrInternal, err.Error(), data)
		return
	}
	respondWithData(c, true, 0, "ok", data)
}
