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
// @Summary 获取文档列表
// @Description 分页获取文档列表，支持按分类和关键词筛选
// @Tags Admin - Documents
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category query string false "分类筛选"
// @Param keyword query string false "关键词搜索"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/documents [get]
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
// @Summary 获取文档详情
// @Description 根据文档 ID 获取文档详细信息
// @Tags Admin - Documents
// @Produce json
// @Param id path int true "文档 ID"
// @Security Bearer
// @Success 200 {object} entity.Document
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/documents/{id} [get]
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
// @Summary 上传文档
// @Description 通过 multipart/form-data 上传文档文件
// @Tags Admin - Documents
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "文档标题"
// @Param category formData string false "文档分类" default(general)
// @Param file formData file true "文档文件"
// @Security Bearer
// @Success 200 {object} entity.Document
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/documents [post]
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
// @Summary 更新文档信息
// @Description 根据文档 ID 更新文档的标题或分类
// @Tags Admin - Documents
// @Accept json
// @Produce json
// @Param id path int true "文档 ID"
// @Param request body v1.UpdateDocumentRequest true "更新文档请求"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/documents/{id} [put]
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
// @Summary 删除文档
// @Description 根据文档 ID 删除文档
// @Tags Admin - Documents
// @Produce json
// @Param id path int true "文档 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/documents/{id} [delete]
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
// @Summary 获取文档分类列表
// @Description 获取所有可用的文档分类
// @Tags Admin - Documents
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/documents/categories [get]
func (c *DocumentController) Categories(ctx *gin.Context) {
	categories, err := c.documentSvc.ListCategories(ctx)
	if err != nil {
		c.Error(ctx, 500, "获取分类列表失败")
		return
	}
	c.Success(ctx, categories)
}

// Download 获取文档下载链接
// @Summary 获取文档下载链接
// @Description 根据文档 ID 获取文档的下载 URL
// @Tags Admin - Documents
// @Produce json
// @Param id path int true "文档 ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/documents/{id}/download [get]
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
