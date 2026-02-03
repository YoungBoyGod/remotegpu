package task

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	common.BaseController
	taskService *serviceTask.TaskService
}

func NewTaskController(ts *serviceTask.TaskService) *TaskController {
	return &TaskController{
		taskService: ts,
	}
}

func (c *TaskController) List(ctx *gin.Context) {
	// userID := ctx.GetUint("userID")
	userID := uint(1) // Mock
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	tasks, total, err := c.taskService.ListTasks(ctx, userID, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "Failed to list tasks")
		return
	}

	c.Success(ctx, gin.H{
		"list":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (c *TaskController) CreateTraining(ctx *gin.Context) {
	var task entity.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}
	
	task.Type = "training"
	// task.CustomerID = ctx.GetUint("userID")
	task.CustomerID = 1 // Mock

	if err := c.taskService.SubmitTask(ctx, &task); err != nil {
		c.Error(ctx, 500, "Failed to submit task")
		return
	}
	c.Success(ctx, task)
}

func (c *TaskController) Stop(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.StopTask(ctx, id); err != nil {
		c.Error(ctx, 500, "Failed to stop task")
		return
	}
	c.Success(ctx, gin.H{"message": "Task stopped"})
}