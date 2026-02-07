package dataset

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceAllocation "github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	serviceDataset "github.com/YoungBoyGod/remotegpu/internal/service/dataset"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	serviceStorage "github.com/YoungBoyGod/remotegpu/internal/service/storage"
	"github.com/gin-gonic/gin"
)

type DatasetController struct {
	common.BaseController
	datasetService    *serviceDataset.DatasetService
	storageService    *serviceStorage.StorageService
	agentService      *serviceOps.AgentService
	allocationService *serviceAllocation.AllocationService
}

func NewDatasetController(ds *serviceDataset.DatasetService, ss *serviceStorage.StorageService, as *serviceOps.AgentService, alloc *serviceAllocation.AllocationService) *DatasetController {
	return &DatasetController{
		datasetService:    ds,
		storageService:    ss,
		agentService:      as,
		allocationService: alloc,
	}
}

// List 获取当前用户的数据集列表
// @author Claude
// @description 根据JWT中的userID过滤，只返回当前用户的数据集
// @reason 原实现使用硬编码userID，存在数据泄露风险
// @modified 2026-02-04
func (c *DatasetController) List(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	datasets, total, err := c.datasetService.ListDatasets(ctx, userID.(uint), page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "获取数据集列表失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      datasets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (c *DatasetController) InitUpload(ctx *gin.Context) {
	var req apiV1.InitMultipartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	// Bucket name strategy could be tenant specific
	uploadID, err := c.storageService.InitMultipart(ctx, "datasets", req.Filename)
	if err != nil {
		c.Error(ctx, 500, "Failed to init upload")
		return
	}

	c.Success(ctx, gin.H{
		"upload_id":  uploadID,
		"chunk_size": 5 * 1024 * 1024,
	})
}

// CompleteUpload 完成分片上传
func (c *DatasetController) CompleteUpload(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	idStr := ctx.Param("id")
	datasetID, _ := strconv.ParseUint(idStr, 10, 64)

	// 验证数据集所有权
	if err := c.datasetService.ValidateOwnership(ctx, uint(datasetID), userID.(uint)); err != nil {
		if err.Error() == "无权限访问该资源" {
			c.Error(ctx, 403, "无权限操作该数据集")
			return
		}
		c.Error(ctx, 404, "数据集不存在")
		return
	}

	var req apiV1.CompleteMultipartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.datasetService.CompleteUpload(ctx, uint(datasetID), req.Name, req.Size); err != nil {
		c.Error(ctx, 500, "完成上传失败")
		return
	}

	c.Success(ctx, gin.H{"message": "上传完成"})
}

// Mount 挂载数据集到机器
// @author Claude
// @description 挂载数据集前校验数据集是否属于当前用户
// @reason 原实现无权限校验，存在越权风险
// @modified 2026-02-04
func (c *DatasetController) Mount(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	idStr := ctx.Param("id")
	datasetID, _ := strconv.ParseUint(idStr, 10, 64)

	// 验证数据集所有权
	if err := c.datasetService.ValidateOwnership(ctx, uint(datasetID), userID.(uint)); err != nil {
		if err.Error() == "无权限访问该资源" {
			c.Error(ctx, 403, "无权限操作该数据集")
			return
		}
		c.Error(ctx, 404, "数据集不存在")
		return
	}

	var req apiV1.MountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	// CodeX 2026-02-04: verify machine ownership before mounting.
	if err := c.allocationService.ValidateHostOwnership(ctx, req.MachineID, userID.(uint)); err != nil {
		if err.Error() == "无权限访问该资源" {
			c.Error(ctx, 403, "无权限操作该机器")
			return
		}
		c.Error(ctx, 404, "机器不存在或未分配")
		return
	}

	// 创建挂载记录并校验路径合法性
	mount, err := c.datasetService.MountDataset(ctx, uint(datasetID), req.MachineID, req.MountPoint, req.ReadOnly)
	if err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	// 异步发送挂载命令到 Agent
	if err := c.agentService.MountDataset(ctx, req.MachineID, uint(datasetID), req.MountPoint); err != nil {
		c.Error(ctx, 500, "挂载数据集失败")
		return
	}

	c.Success(ctx, gin.H{"message": "挂载命令已发送", "mount_id": mount.ID})
}

// Unmount 卸载数据集
func (c *DatasetController) Unmount(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	idStr := ctx.Param("id")
	datasetID, _ := strconv.ParseUint(idStr, 10, 64)

	// 验证数据集所有权
	if err := c.datasetService.ValidateOwnership(ctx, uint(datasetID), userID.(uint)); err != nil {
		c.Error(ctx, 403, "无权限操作该数据集")
		return
	}

	mountIDStr := ctx.Param("mount_id")
	mountID, err := strconv.ParseUint(mountIDStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的挂载 ID")
		return
	}

	if err := c.datasetService.UnmountDataset(ctx, uint(mountID)); err != nil {
		c.Error(ctx, 500, err.Error())
		return
	}
	c.Success(ctx, gin.H{"message": "卸载成功"})
}

// ListMounts 获取数据集的挂载列表
func (c *DatasetController) ListMounts(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		c.Error(ctx, 401, "用户未认证")
		return
	}

	idStr := ctx.Param("id")
	datasetID, _ := strconv.ParseUint(idStr, 10, 64)

	if err := c.datasetService.ValidateOwnership(ctx, uint(datasetID), userID.(uint)); err != nil {
		c.Error(ctx, 403, "无权限操作该数据集")
		return
	}

	mounts, err := c.datasetService.ListMountsByDataset(ctx, uint(datasetID))
	if err != nil {
		c.Error(ctx, 500, "获取挂载列表失败")
		return
	}
	c.Success(ctx, mounts)
}
