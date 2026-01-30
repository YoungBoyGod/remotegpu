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
