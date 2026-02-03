package v1

import (
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/gin-gonic/gin"
)

type DatasetController struct {
	BaseController
	datasetService *service.DatasetService
	storageService *service.StorageService
	agentService   *service.AgentService
}

func NewDatasetController(ds *service.DatasetService, ss *service.StorageService, as *service.AgentService) *DatasetController {
	return &DatasetController{
		datasetService: ds,
		storageService: ss,
		agentService:   as,
	}
}

func (c *DatasetController) List(ctx *gin.Context) {
	// userID := ctx.GetUint("userID")
	userID := uint(1) // Mock
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	datasets, total, err := c.datasetService.ListDatasets(ctx, userID, page, pageSize)
	if err != nil {
		c.Error(ctx, 500, "Failed to list datasets")
		return
	}

	c.Success(ctx, gin.H{
		"list":      datasets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

type InitMultipartRequest struct {
	Filename string `json:"filename" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
	MD5      string `json:"md5"`
}

func (c *DatasetController) InitUpload(ctx *gin.Context) {
	var req InitMultipartRequest
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
		"upload_id": uploadID,
		"chunk_size": 5 * 1024 * 1024,
	})
}

type MountRequest struct {
	MachineID string `json:"machine_id" binding:"required"`
	MountPoint string `json:"mount_point" binding:"required"`
	ReadOnly  bool   `json:"read_only"`
}

func (c *DatasetController) Mount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	datasetID, _ := strconv.ParseUint(idStr, 10, 64)

	var req MountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	if err := c.agentService.MountDataset(ctx, req.MachineID, uint(datasetID), req.MountPoint); err != nil {
		c.Error(ctx, 500, "Failed to mount dataset")
		return
	}

	c.Success(ctx, gin.H{"message": "Mount command sent"})
}
