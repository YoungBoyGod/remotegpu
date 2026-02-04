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
func (c *ImageController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

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
