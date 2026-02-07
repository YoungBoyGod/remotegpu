package task

import (
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"
	"github.com/gin-gonic/gin"
)

// AgentTaskController Agent 专用任务控制器
type AgentTaskController struct {
	common.BaseController
	taskService *serviceTask.TaskService
}

func NewAgentTaskController(ts *serviceTask.TaskService) *AgentTaskController {
	return &AgentTaskController{taskService: ts}
}

// ClaimTasks Agent 认领任务
func (c *AgentTaskController) ClaimTasks(ctx *gin.Context) {
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
		MachineID string `json:"machine_id" binding:"required"`
		Limit     int    `json:"limit"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	tasks, err := c.taskService.ClaimTasks(ctx, req.MachineID, req.AgentID, req.Limit)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"tasks": tasks})
}

// StartTask 标记任务开始
func (c *AgentTaskController) StartTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
		AttemptID string `json:"attempt_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.taskService.StartTask(ctx, id, req.AgentID, req.AttemptID); err != nil {
		c.Error(ctx, 409, err.Error())
		return
	}
	c.Success(ctx, gin.H{"task_id": id, "status": "running"})
}

// RenewLease 续约租约
func (c *AgentTaskController) RenewLease(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
		AttemptID string `json:"attempt_id" binding:"required"`
		ExtendSec int    `json:"extend_sec"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.taskService.RenewLease(ctx, id, req.AgentID, req.AttemptID, req.ExtendSec); err != nil {
		c.Error(ctx, 410, err.Error())
		return
	}
	c.Success(ctx, gin.H{"task_id": id, "renewed": true})
}

// CompleteTask 完成任务
func (c *AgentTaskController) CompleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
		AttemptID string `json:"attempt_id" binding:"required"`
		ExitCode  int    `json:"exit_code"`
		Error     string `json:"error"`
		Stdout    string `json:"stdout"`
		Stderr    string `json:"stderr"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.taskService.CompleteTask(ctx, id, req.AgentID, req.AttemptID, req.ExitCode, req.Error, req.Stdout, req.Stderr); err != nil {
		c.Error(ctx, 409, err.Error())
		return
	}
	c.Success(ctx, gin.H{"task_id": id, "status": "completed"})
}

// ReportProgress 上报任务进度
func (c *AgentTaskController) ReportProgress(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
		AttemptID string `json:"attempt_id" binding:"required"`
		Percent   int    `json:"percent"`
		Message   string `json:"message"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.taskService.ReportProgress(ctx, id, req.AgentID, req.AttemptID, req.Percent, req.Message); err != nil {
		c.Error(ctx, 409, err.Error())
		return
	}
	c.Success(ctx, gin.H{"task_id": id, "progress": req.Percent})
}
