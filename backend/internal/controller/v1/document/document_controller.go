package document

import (
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	serviceDocument "github.com/YoungBoyGod/remotegpu/internal/service/document"
	serviceStorage "github.com/YoungBoyGod/remotegpu/internal/service/storage"
	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	common.BaseController
	documentSvc *serviceDocument.DocumentService
	storageSvc  *serviceStorage.StorageService
}

func NewDocumentController(ds *serviceDocument.DocumentService, ss *serviceStorage.StorageService) *DocumentController {
	return &DocumentController{
		documentSvc: ds,
		storageSvc:  ss,
	}
}

// List 获取文档列表
func (c *DocumentController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	category := ctx.Query("category")
	keyword := ctx.Query("keyword")

	docs, total, err := c.documentSvc.ListDocuments(ctx, page, pageSize, category, keyword)
	if err != nil {
		c.Error(ctx, 500, "获取文档列表失败")
		return
	}

	c.Success(ctx, gin.H{
		"list":      docs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Detail 获取文档详情
func (c *DocumentController) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的文档 ID")
		return
	}

	doc, err := c.documentSvc.GetDocument(ctx, uint(id))
	if err != nil {
		c.Error(ctx, 404, "文档不存在")
		return
	}
	c.Success(ctx, doc)
}

// Upload 上传文档（multipart/form-data）
func (c *DocumentController) Upload(ctx *gin.Context) {
	title := ctx.PostForm("title")
	category := ctx.PostForm("category")
	if category == "" {
		category = "general"
	}
	if title == "" {
		c.Error(ctx, 400, "标题不能为空")
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		c.Error(ctx, 400, "请选择要上传的文件")
		return
	}
	defer file.Close()

	storagePath := serviceDocument.BuildStoragePath(category, header.Filename)

	doc := &entity.Document{
		Title:          title,
		Category:       category,
		FileName:       header.Filename,
		FilePath:       storagePath,
		FileSize:       header.Size,
		ContentType:    header.Header.Get("Content-Type"),
		StorageBackend: "",
		UploadedBy:     ctx.GetUint("userID"),
	}

	if err := c.documentSvc.CreateDocument(ctx, doc); err != nil {
		c.Error(ctx, 500, "创建文档记录失败")
		return
	}

	c.Success(ctx, doc)
}

// Update 更新文档信息
func (c *DocumentController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的文档 ID")
		return
	}

	var req apiV1.UpdateDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, err.Error())
		return
	}

	fields := make(map[string]any)
	if req.Title != "" {
		fields["title"] = req.Title
	}
	if req.Category != "" {
		fields["category"] = req.Category
	}
	if len(fields) == 0 {
		c.Error(ctx, 400, "没有需要更新的字段")
		return
	}

	if err := c.documentSvc.UpdateDocument(ctx, uint(id), fields); err != nil {
		c.Error(ctx, 500, "更新文档失败")
		return
	}
	c.Success(ctx, gin.H{"message": "文档已更新"})
}

// Delete 删除文档
func (c *DocumentController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的文档 ID")
		return
	}

	if err := c.documentSvc.DeleteDocument(ctx, uint(id)); err != nil {
		c.Error(ctx, 500, "删除文档失败")
		return
	}
	c.Success(ctx, gin.H{"message": "文档已删除"})
}

// Categories 获取分类列表
func (c *DocumentController) Categories(ctx *gin.Context) {
	categories, err := c.documentSvc.ListCategories(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取分类列表失败")
		return
	}
	c.Success(ctx, categories)
}

// Download 获取文档下载链接
func (c *DocumentController) Download(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的文档 ID")
		return
	}

	url, err := c.documentSvc.GetDownloadURL(ctx, uint(id))
	if err != nil {
		c.Error(ctx, 500, "获取下载链接失败")
		return
	}
	c.Success(ctx, gin.H{"url": url})
}
