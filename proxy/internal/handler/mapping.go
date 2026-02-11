package handler

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/remotegpu-proxy/internal/forwarder"
	"github.com/YoungBoyGod/remotegpu-proxy/internal/models"
	"github.com/gin-gonic/gin"
)

// MappingHandler 映射管理 API Handler
type MappingHandler struct {
	manager *forwarder.Manager
}

// NewMappingHandler 创建 Handler
func NewMappingHandler(manager *forwarder.Manager) *MappingHandler {
	return &MappingHandler{manager: manager}
}

// respond 统一响应
func respond(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, models.MappingResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// AddMapping 添加映射 POST /api/v1/mappings
func (h *MappingHandler) AddMapping(c *gin.Context) {
	var req models.MappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, 1, "参数错误: "+err.Error(), nil)
		return
	}

	info, err := h.manager.AddMapping(&req)
	if err != nil {
		respond(c, 2, "添加映射失败: "+err.Error(), nil)
		return
	}

	respond(c, 0, "success", info)
}

// RemoveMapping 移除映射 DELETE /api/v1/mappings/:port
func (h *MappingHandler) RemoveMapping(c *gin.Context) {
	port, err := strconv.Atoi(c.Param("port"))
	if err != nil {
		respond(c, 1, "端口参数无效", nil)
		return
	}

	if err := h.manager.RemoveMapping(port); err != nil {
		respond(c, 2, "移除映射失败: "+err.Error(), nil)
		return
	}

	respond(c, 0, "success", nil)
}

// RemoveByEnvID 移除环境所有映射 DELETE /api/v1/mappings/env/:id
func (h *MappingHandler) RemoveByEnvID(c *gin.Context) {
	envID := c.Param("id")
	if envID == "" {
		respond(c, 1, "环境 ID 不能为空", nil)
		return
	}

	if err := h.manager.RemoveByEnvID(envID); err != nil {
		respond(c, 2, "移除环境映射失败: "+err.Error(), nil)
		return
	}

	respond(c, 0, "success", nil)
}

// ListMappings 列出所有映射 GET /api/v1/mappings
func (h *MappingHandler) ListMappings(c *gin.Context) {
	list := h.manager.ListMappings()
	respond(c, 0, "success", list)
}

// GetStats 获取统计信息 GET /api/v1/stats
func (h *MappingHandler) GetStats(c *gin.Context) {
	stats := h.manager.Stats()
	respond(c, 0, "success", stats)
}

// Ping 健康检查 GET /api/v1/ping
func (h *MappingHandler) Ping(c *gin.Context) {
	respond(c, 0, "pong", nil)
}
