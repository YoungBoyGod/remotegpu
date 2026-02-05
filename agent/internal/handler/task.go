package handler

import (
	"net/http"

	"github.com/YoungBoyGod/remotegpu-agent/internal/errors"
	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
	"github.com/YoungBoyGod/remotegpu-agent/internal/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	scheduler *scheduler.Scheduler
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(s *scheduler.Scheduler) *TaskHandler {
	return &TaskHandler{scheduler: s}
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		Name    string            `json:"name"`
		Type    string            `json:"type"`
		Command string            `json:"command" binding:"required"`
		Args    []string          `json:"args"`
		WorkDir string            `json:"workdir"`
		Env     map[string]string `json:"env"`
		Timeout int               `json:"timeout"`
		Priority int              `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errors.ErrInvalidParams,
			"message": err.Error(),
		})
		return
	}

	task := &models.Task{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Type:     models.TaskType(req.Type),
		Command:  req.Command,
		Args:     req.Args,
		WorkDir:  req.WorkDir,
		Env:      req.Env,
		Timeout:  req.Timeout,
		Priority: req.Priority,
	}

	if task.Type == "" {
		task.Type = models.TaskTypeShell
	}
	if task.Priority == 0 {
		task.Priority = 5
	}

	if err := h.scheduler.Submit(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errors.ErrInternal,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    task,
	})
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.scheduler.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    errors.ErrTaskNotFound,
			"message": "task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    task,
	})
}

// CancelTask 取消任务
func (h *TaskHandler) CancelTask(c *gin.Context) {
	id := c.Param("id")
	if h.scheduler.CancelTask(id) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    errors.ErrTaskNotFound,
			"message": "task not found",
		})
	}
}

// GetQueueStatus 获取队列状态
func (h *TaskHandler) GetQueueStatus(c *gin.Context) {
	status := h.scheduler.GetQueueStatus()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    status,
	})
}
