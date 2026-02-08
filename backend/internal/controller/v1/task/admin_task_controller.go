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
// @Summary 管理员获取任务列表
// @Description 管理员分页获取所有任务，支持按状态筛选
// @Tags Admin - Tasks
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "状态筛选"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/tasks [get]
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
// @Summary 管理员获取任务详情
// @Description 管理员根据任务 ID 获取任务详细信息
// @Tags Admin - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} entity.Task
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/tasks/{id} [get]
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
// @Summary 管理员创建任务
// @Description 管理员创建新任务
// @Tags Admin - Tasks
// @Accept json
// @Produce json
// @Param request body entity.Task true "任务信息"
// @Security Bearer
// @Success 200 {object} entity.Task
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/tasks [post]
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
// @Summary 管理员停止任务
// @Description 管理员停止指定任务
// @Tags Admin - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/tasks/{id}/stop [post]
func (c *AdminTaskController) Stop(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.StopTask(ctx, id); err != nil {
		c.Error(ctx, 500, "停止任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已停止"})
}

// Cancel 管理员取消任务
// @Summary 管理员取消任务
// @Description 管理员取消指定任务
// @Tags Admin - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/tasks/{id}/cancel [post]
func (c *AdminTaskController) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.CancelTask(ctx, id); err != nil {
		c.Error(ctx, 500, "取消任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已取消"})
}

// Retry 管理员重试任务
// @Summary 管理员重试任务
// @Description 管理员重试指定任务，重新排队执行
// @Tags Admin - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/tasks/{id}/retry [post]
func (c *AdminTaskController) Retry(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.taskService.RetryTask(ctx, id); err != nil {
		c.Error(ctx, 500, "重试任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已重新排队"})
}

// Logs 获取任务日志
// @Summary 管理员获取任务日志
// @Description 管理员获取指定任务的日志信息
// @Tags Admin - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/tasks/{id}/logs [get]
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
// @Summary 管理员获取任务结果
// @Description 管理员获取指定任务的执行结果
// @Tags Admin - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/tasks/{id}/result [get]
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
