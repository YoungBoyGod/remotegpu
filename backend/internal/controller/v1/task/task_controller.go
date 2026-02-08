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

// List 获取当前用户的任务列表
// @Summary 获取任务列表
// @Description 根据当前登录用户获取其任务列表，支持分页
// @Tags Customer - Tasks
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security Bearer
// @Success 200 {object} map[string]interface
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/tasks [get]
//
// @author Claude
// @description 根据JWT中的userID过滤，只返回当前用户的任务
// @reason 原实现使用硬编码userID，存在数据泄露风险
// @modified 2026-02-04
func (c *TaskController) List(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	tasks, total, err := c.taskService.ListTasks(ctx, userID, page, pageSize)
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

// CreateTraining 创建训练任务
// @Summary 创建训练任务
// @Description 创建新的训练任务并绑定当前用户
// @Tags Customer - Tasks
// @Accept json
// @Produce json
// @Param request body entity.Task true "任务信息"
// @Security Bearer
// @Success 200 {object} entity.Task
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/tasks/training [post]
//
// @author Claude
// @description 创建训练任务并绑定当前用户，从JWT获取userID
// @reason 原实现使用硬编码CustomerID，存在安全风险
// @modified 2026-02-04
func (c *TaskController) CreateTraining(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	var task entity.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	task.Type = "training"
	task.CustomerID = userID

	if err := c.taskService.SubmitTask(ctx, &task); err != nil {
		c.Error(ctx, 500, "创建任务失败")
		return
	}
	c.Success(ctx, task)
}

// Detail 获取任务详情
// @Summary 获取任务详情
// @Description 根据任务 ID 获取任务详细信息，校验任务归属
// @Tags Customer - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} entity.Task
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /customer/tasks/{id} [get]
func (c *TaskController) Detail(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	id := ctx.Param("id")
	task, err := c.taskService.GetTask(ctx, id)
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}

	// 校验任务归属
	if task.CustomerID != userID {
		c.Error(ctx, 403, "无权访问该任务")
		return
	}

	c.Success(ctx, task)
}

// Stop 停止任务
// @Summary 停止任务
// @Description 停止指定任务，校验任务归属后执行
// @Tags Customer - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/tasks/{id}/stop [post]
//
// @author Claude
// @description 停止任务前校验任务是否属于当前用户，防止越权操作
// @reason 原实现无权限校验，存在越权风险
// @modified 2026-02-04
func (c *TaskController) Stop(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	id := ctx.Param("id")
	if err := c.taskService.StopTaskWithAuth(ctx, id, userID); err != nil {
		if err.Error() == "无权限访问该资源" {
			c.Error(ctx, 403, "无权限操作该任务")
			return
		}
		c.Error(ctx, 500, "停止任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已停止"})
}

// Cancel 取消任务（带权限校验）
// @Summary 取消任务
// @Description 取消指定任务，校验任务归属后执行
// @Tags Customer - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/tasks/{id}/cancel [post]
func (c *TaskController) Cancel(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	id := ctx.Param("id")
	if err := c.taskService.CancelTaskWithAuth(ctx, id, userID); err != nil {
		if err == entity.ErrUnauthorized {
			c.Error(ctx, 403, "无权限操作该任务")
			return
		}
		c.Error(ctx, 500, "取消任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已取消"})
}

// Retry 重试任务（带权限校验）
// @Summary 重试任务
// @Description 重试指定任务，校验任务归属后重新排队
// @Tags Customer - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /customer/tasks/{id}/retry [post]
func (c *TaskController) Retry(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	id := ctx.Param("id")
	if err := c.taskService.RetryTaskWithAuth(ctx, id, userID); err != nil {
		if err == entity.ErrUnauthorized {
			c.Error(ctx, 403, "无权限操作该任务")
			return
		}
		c.Error(ctx, 500, "重试任务失败")
		return
	}
	c.Success(ctx, gin.H{"message": "任务已重新排队"})
}

// Logs 获取任务日志（带权限校验）
// @Summary 获取任务日志
// @Description 获取指定任务的日志信息，校验任务归属
// @Tags Customer - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /customer/tasks/{id}/logs [get]
func (c *TaskController) Logs(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	id := ctx.Param("id")
	task, err := c.taskService.GetTask(ctx, id)
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}
	if task.CustomerID != userID {
		c.Error(ctx, 403, "无权访问该任务")
		return
	}

	c.Success(ctx, gin.H{
		"task_id":   task.ID,
		"status":    task.Status,
		"error_msg": task.ErrorMsg,
	})
}

// Result 获取任务结果（带权限校验）
// @Summary 获取任务结果
// @Description 获取指定任务的执行结果，包含退出码和错误信息
// @Tags Customer - Tasks
// @Produce json
// @Param id path string true "任务 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /customer/tasks/{id}/result [get]
func (c *TaskController) Result(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	id := ctx.Param("id")
	task, err := c.taskService.GetTask(ctx, id)
	if err != nil {
		c.Error(ctx, 404, "任务不存在")
		return
	}
	if task.CustomerID != userID {
		c.Error(ctx, 403, "无权访问该任务")
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