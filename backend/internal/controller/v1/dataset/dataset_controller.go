package dataset

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceDataset "github.com/YoungBoyGod/remotegpu/internal/service/dataset"
	serviceOps "github.com/YoungBoyGod/remotegpu/internal/service/ops"
	serviceStorage "github.com/YoungBoyGod/remotegpu/internal/service/storage"
	"github.com/gin-gonic/gin"
)

type DatasetController struct {
	common.BaseController
	datasetService *serviceDataset.DatasetService
	storageService *serviceStorage.StorageService
	agentService   *serviceOps.AgentService
}

func NewDatasetController(ds *serviceDataset.DatasetService, ss *serviceStorage.StorageService, as *serviceOps.AgentService) *DatasetController {
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
		"upload_id": uploadID,
		"chunk_size": 5 * 1024 * 1024,
	})
}

func (c *DatasetController) Mount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	datasetID, _ := strconv.ParseUint(idStr, 10, 64)

	var req apiV1.MountRequest
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
