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
// @Summary 获取存储后端列表
// @Description 获取系统配置的所有存储后端信息
// @Tags Admin - Storage
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /admin/storage/backends [get]
func (c *StorageController) ListBackends(ctx *gin.Context) {
	backends := c.storageSvc.ListBackends()
	c.Success(ctx, gin.H{"backends": backends})
}

// GetStats 获取存储使用统计
// @Summary 获取存储使用统计
// @Description 获取指定存储后端的使用统计信息
// @Tags Admin - Storage
// @Produce json
// @Param backend query string false "存储后端名称"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/storage/stats [get]
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
// @Summary 列出存储文件
// @Description 列出指定存储后端和前缀下的文件
// @Tags Admin - Storage
// @Produce json
// @Param backend query string false "存储后端名称"
// @Param prefix query string false "文件路径前缀"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/storage/files [get]
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
// @Summary 删除存储文件
// @Description 删除指定存储后端中的文件
// @Tags Admin - Storage
// @Accept json
// @Produce json
// @Param request body object true "删除文件请求（backend, path）"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/storage/files/delete [post]
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
// @Summary 获取文件下载链接
// @Description 获取指定存储后端中文件的下载 URL
// @Tags Admin - Storage
// @Produce json
// @Param backend query string false "存储后端名称"
// @Param path query string true "文件路径"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/storage/files/download-url [get]
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
