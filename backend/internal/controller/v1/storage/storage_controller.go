package storage

import (
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceStorage "github.com/YoungBoyGod/remotegpu/internal/service/storage"
	"github.com/gin-gonic/gin"
)

// StorageController 存储管理控制器
type StorageController struct {
	common.BaseController
	storageSvc *serviceStorage.StorageService
}

func NewStorageController(svc *serviceStorage.StorageService) *StorageController {
	return &StorageController{storageSvc: svc}
}

// ListBackends 获取存储池列表
func (c *StorageController) ListBackends(ctx *gin.Context) {
	backends := c.storageSvc.ListBackends()
	c.Success(ctx, gin.H{"backends": backends})
}

// GetStats 获取存储使用统计
func (c *StorageController) GetStats(ctx *gin.Context) {
	backendName := ctx.Query("backend")
	stats, err := c.storageSvc.GetBackendStats(ctx, backendName)
	if err != nil {
		c.Error(ctx, 500, "获取存储统计失败: "+err.Error())
		return
	}
	c.Success(ctx, stats)
}

// ListFiles 列出存储文件
func (c *StorageController) ListFiles(ctx *gin.Context) {
	backendName := ctx.Query("backend")
	prefix := ctx.Query("prefix")

	files, err := c.storageSvc.ListFiles(ctx, backendName, prefix)
	if err != nil {
		c.Error(ctx, 500, "获取文件列表失败: "+err.Error())
		return
	}
	c.Success(ctx, gin.H{
		"files": files,
		"total": len(files),
	})
}

// DeleteFile 删除文件
func (c *StorageController) DeleteFile(ctx *gin.Context) {
	var req struct {
		Backend string `json:"backend"`
		Path    string `json:"path" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.storageSvc.DeleteFile(ctx, req.Backend, req.Path); err != nil {
		c.Error(ctx, 500, "删除文件失败: "+err.Error())
		return
	}
	c.Success(ctx, gin.H{"message": "删除成功"})
}

// GetDownloadURL 获取文件下载链接
func (c *StorageController) GetDownloadURL(ctx *gin.Context) {
	backendName := ctx.Query("backend")
	path := ctx.Query("path")
	if path == "" {
		c.Error(ctx, 400, "path 参数不能为空")
		return
	}

	url, err := c.storageSvc.GetDownloadURL(ctx, backendName, path)
	if err != nil {
		c.Error(ctx, 500, "获取下载链接失败: "+err.Error())
		return
	}
	c.Success(ctx, gin.H{"url": url})
}
