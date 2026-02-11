package proxy

import (
	"github.com/gin-gonic/gin"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceProxy "github.com/YoungBoyGod/remotegpu/internal/service/proxy"
)

// ProxyController Proxy 管理控制器
type ProxyController struct {
	common.BaseController
	proxySvc *serviceProxy.ProxyService
}

func NewProxyController(proxySvc *serviceProxy.ProxyService) *ProxyController {
	return &ProxyController{proxySvc: proxySvc}
}

// Register 处理 Proxy 注册
func (c *ProxyController) Register(ctx *gin.Context) {
	var req serviceProxy.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	node, err := c.proxySvc.Register(ctx, &req)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, node)
}

// Heartbeat 处理 Proxy 心跳
func (c *ProxyController) Heartbeat(ctx *gin.Context) {
	var req serviceProxy.HeartbeatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.proxySvc.Heartbeat(ctx, &req); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "ok"})
}

// ListNodes 列出所有 Proxy 节点（Admin）
func (c *ProxyController) ListNodes(ctx *gin.Context) {
	nodes, err := c.proxySvc.ListNodes(ctx)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, nodes)
}

// GetNode 获取 Proxy 节点详情（Admin）
func (c *ProxyController) GetNode(ctx *gin.Context) {
	id := ctx.Param("id")
	node, err := c.proxySvc.GetNode(ctx, id)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, node)
}

// DeleteNode 删除 Proxy 节点（Admin）
func (c *ProxyController) DeleteNode(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.proxySvc.DeleteNode(ctx, id); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"status": "deleted"})
}

// ListMappings 列出所有端口映射（Admin）
func (c *ProxyController) ListMappings(ctx *gin.Context) {
	mappings, err := c.proxySvc.ListMappings(ctx)
	if err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, mappings)
}
