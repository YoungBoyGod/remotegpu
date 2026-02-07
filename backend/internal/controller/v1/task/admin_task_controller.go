package task

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceTask "github.com/YoungBoyGod/remotegpu/internal/service/task"
	"github.com/gin-gonic/gin"
)

// AdminTaskController 管理员任务控制器
type AdminTaskController struct {
	common.BaseController
	taskService *serviceTask.TaskService
}

func NewAdminTaskController(ts *serviceTask.TaskService) *AdminTaskController {
	return &AdminTaskController{taskService: ts}
}

// List 管理员查询所有任务
func (c *AdminTaskController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	status := ctx.Query("status")

	tasks, total, err := c.taskService.ListAllTasks(ctx, status, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "获取任务列表失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Detail 管理员查看任务详情
func (c *AdminTaskController) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.taskService.GetTask(ctx, id)
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}
	c.Success(ctx, task)
}

// Create 管理员创建任务
func (c *AdminTaskController) Create(ctx *gin.Context) {
	var task entity.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.taskService.SubmitTask(ctx, &task); err != nil {
		c.Error(ctx, 500, "创建任务失败")
		return
	}
	c.Success(ctx, task)
}

// Stop 管理员停止任务
func (c *AdminTaskController) Stop(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.StopTask(ctx, id); err != nil {
		c.Error(ctx, 500, "停止任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已停止"})
}

// Cancel 管理员取消任务
func (c *AdminTaskController) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.CancelTask(ctx, id); err != nil {
		c.Error(ctx, 500, "取消任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已取消"})
}

// Retry 管理员重试任务
func (c *AdminTaskController) Retry(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.RetryTask(ctx, id); err != nil {
		c.Error(ctx, 500, "重试任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已重新排队"})
}

// Logs 获取任务日志
func (c *AdminTaskController) Logs(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.taskService.GetTask(ctx, id)
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}
	c.Success(ctx, gin.H{
		"task_id":  task.ID,
		"status":   task.Status,
		"error_msg": task.ErrorMsg,
	})
}

// Result 获取任务结果
func (c *AdminTaskController) Result(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.taskService.GetTask(ctx, id)
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}
	c.Success(ctx, gin.H{
		"task_id":   task.ID,
		"status":    task.Status,
		"exit_code": task.ExitCode,
		"error_msg": task.ErrorMsg,
		"ended_at":  task.EndedAt,
	})
}
