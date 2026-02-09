package workspace

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceWorkspace "github.com/YoungBoyGod/remotegpu/internal/service/workspace"
	"github.com/gin-gonic/gin"
)

// WorkspaceController 工作空间控制器
type WorkspaceController struct {
	common.BaseController
	workspaceService *serviceWorkspace.WorkspaceService
}

func NewWorkspaceController(svc *serviceWorkspace.WorkspaceService) *WorkspaceController {
	return &WorkspaceController{workspaceService: svc}
}

// getCustomerID 从上下文获取当前登录用户 ID
func (c *WorkspaceController) getCustomerID(ctx *gin.Context) (uint, bool) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "未认证")
		return 0, false
	}
	return userID, true
}

// Create 创建工作空间
func (c *WorkspaceController) Create(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	var req apiV1.CreateWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	ws, err := c.workspaceService.Create(ctx, customerID, req.Name, req.Description)
	if err != nil {
		c.Error(ctx, 500, "创建工作空间失败: "+err.Error())
		return
	}
	c.Success(ctx, ws)
}

// List 工作空间列表
func (c *WorkspaceController) List(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	items, total, err := c.workspaceService.List(ctx, customerID, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "获取工作空间列表失败")
		return
	}

	c.Success(ctx, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Detail 工作空间详情
func (c *WorkspaceController) Detail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的工作空间 ID")
		return
	}

	ws, err := c.workspaceService.GetByID(ctx, uint(id))
	if err != nil {
		c.Error(ctx, 404, "工作空间不存在")
		return
	}
	c.Success(ctx, ws)
}

// Update 更新工作空间
func (c *WorkspaceController) Update(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的工作空间 ID")
		return
	}

	var req apiV1.UpdateWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	fields := make(map[string]interface{})
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.Description != "" {
		fields["description"] = req.Description
	}
	if len(fields) == 0 {
		c.Error(ctx, 400, "没有需要更新的字段")
		return
	}

	if err := c.workspaceService.Update(ctx, uint(id), customerID, fields); err != nil {
		c.Error(ctx, 403, err.Error())
		return
	}
	c.Success(ctx, nil)
}

// Delete 删除工作空间
func (c *WorkspaceController) Delete(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的工作空间 ID")
		return
	}

	if err := c.workspaceService.Delete(ctx, uint(id), customerID); err != nil {
		c.Error(ctx, 403, err.Error())
		return
	}
	c.Success(ctx, nil)
}

// AddMember 添加工作空间成员
func (c *WorkspaceController) AddMember(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的工作空间 ID")
		return
	}

	var req apiV1.AddMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.workspaceService.AddMember(ctx, uint(id), customerID, req.UserID, req.Role); err != nil {
		c.Error(ctx, 403, err.Error())
		return
	}
	c.Success(ctx, nil)
}

// RemoveMember 移除工作空间成员
func (c *WorkspaceController) RemoveMember(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的工作空间 ID")
		return
	}

	targetID, err := strconv.ParseUint(ctx.Param("userId"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的用户 ID")
		return
	}

	if err := c.workspaceService.RemoveMember(ctx, uint(id), customerID, uint(targetID)); err != nil {
		c.Error(ctx, 403, err.Error())
		return
	}
	c.Success(ctx, nil)
}

// ListMembers 获取工作空间成员列表
func (c *WorkspaceController) ListMembers(ctx *gin.Context) {
	customerID, ok := c.getCustomerID(ctx)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的工作空间 ID")
		return
	}

	members, err := c.workspaceService.ListMembers(ctx, uint(id), customerID)
	if err != nil {
		c.Error(ctx, 403, err.Error())
		return
	}
	c.Success(ctx, members)
}
