package v1

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// GPUController GPU控制器
type GPUController struct {
	gpuService *service.GPUService
}

// NewGPUController 创建GPU控制器
func NewGPUController() *GPUController {
	return &GPUController{
		gpuService: service.NewGPUService(),
	}
}

// Create 创建GPU
func (ctrl *GPUController) Create(c *gin.Context) {
	var gpu entity.GPU
	if err := c.ShouldBindJSON(&gpu); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.gpuService.Create(&gpu); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gpu)
}

// GetByID 获取GPU详情
func (ctrl *GPUController) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	gpu, err := ctrl.gpuService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "GPU不存在")
		return
	}
	response.Success(c, gpu)
}

// GetByHostID 获取主机的GPU列表
func (ctrl *GPUController) GetByHostID(c *gin.Context) {
	hostID := c.Param("host_id")
	gpus, err := ctrl.gpuService.GetByHostID(hostID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gpus)
}

// Delete 删除GPU
func (ctrl *GPUController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.gpuService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// List 获取GPU列表
func (ctrl *GPUController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	gpus, total, err := ctrl.gpuService.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  gpus,
		"total": total,
	})
}

// Update 更新GPU
func (ctrl *GPUController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var gpu entity.GPU
	if err := c.ShouldBindJSON(&gpu); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	gpu.ID = uint(id)

	if err := ctrl.gpuService.Update(&gpu); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// Allocate 分配GPU
func (ctrl *GPUController) Allocate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		EnvID string `json:"env_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.gpuService.Allocate(uint(id), req.EnvID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// Release 释放GPU
func (ctrl *GPUController) Release(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.gpuService.Release(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}
