package ops

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/service/image"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	common.BaseController
	imageService *image.ImageService
}

func NewImageController(svc *image.ImageService) *ImageController {
	return &ImageController{imageService: svc}
}

// List 获取镜像列表
// @Summary 获取镜像列表
// @Description 获取系统可用镜像列表，支持分页和筛选
// @Tags Admin - Images
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "分类筛选"
// @Param framework query string false "框架筛选"
// @Param status query string false "状态筛选"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/images [get]
func (c *ImageController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	// CodeX 2026-02-04: normalize paging input to avoid invalid offsets.
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}

	params := dao.ImageListParams{
		Page:      page,
		PageSize:  pageSize,
		Category:  ctx.Query("category"),
		Framework: ctx.Query("framework"),
		Status:    ctx.Query("status"),
	}

	images, total, err := c.imageService.List(ctx, params)
	if err != nil {
		c.Error(ctx, 500, "获取镜像列表失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Sync 同步镜像
// @Summary 同步镜像
// @Description 触发镜像仓库同步任务
// @Tags Admin - Images
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/images/sync [post]
func (c *ImageController) Sync(ctx *gin.Context) {
	if err := c.imageService.Sync(ctx); err != nil {
		c.Error(ctx, 500, "同步镜像失败")
		return
	}
	c.Success(ctx, gin.H{"message": "同步任务已触发"})
}
