package v1

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// HostController 主机控制器
type HostController struct {
	hostService *service.HostService
}

// NewHostController 创建主机控制器
func NewHostController() *HostController {
	return &HostController{
		hostService: service.NewHostService(),
	}
}

// Create 创建主机
func (ctrl *HostController) Create(c *gin.Context) {
	var host entity.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.hostService.Create(&host); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, host)
}

// GetByID 获取主机详情
func (ctrl *HostController) GetByID(c *gin.Context) {
	id := c.Param("id")
	host, err := ctrl.hostService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "主机不存在")
		return
	}
	response.Success(c, host)
}

// List 获取主机列表
func (ctrl *HostController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	hosts, total, err := ctrl.hostService.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  hosts,
		"total": total,
	})
}

// Update 更新主机
func (ctrl *HostController) Update(c *gin.Context) {
	id := c.Param("id")
	var host entity.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	host.ID = id

	if err := ctrl.hostService.Update(&host); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete 删除主机
func (ctrl *HostController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.hostService.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// Heartbeat 心跳
func (ctrl *HostController) Heartbeat(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.hostService.Heartbeat(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}
