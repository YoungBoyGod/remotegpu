package v1

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/YoungBoyGod/remotegpu/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StorageController 存储控制器
type StorageController struct {
	manager *storage.Manager
}

// NewStorageController 创建存储控制器
func NewStorageController(manager *storage.Manager) *StorageController {
	return &StorageController{manager: manager}
}

// ListBackends 列出所有存储后端
func (s *StorageController) ListBackends(c *gin.Context) {
	backends := s.manager.List()
	response.Success(c, gin.H{
		"backends":        backends,
		"max_upload_size": config.GlobalConfig.Storage.MaxUploadSize,
	})
}

// Upload 上传文件
func (s *StorageController) Upload(c *gin.Context) {
	backendName := c.PostForm("backend")
	backend, err := s.manager.Get(backendName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}
	defer file.Close()

	if header.Size > config.GlobalConfig.Storage.MaxUploadSize {
		response.Error(c, http.StatusBadRequest, "文件大小超过限制")
		return
	}

	// 生成存储路径
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	path := c.PostForm("path")
	if path != "" {
		filename = filepath.Join(path, filename)
	}

	opts := &storage.UploadOptions{
		ContentType: header.Header.Get("Content-Type"),
	}

	if err := backend.Upload(c.Request.Context(), filename, file, header.Size, opts); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"filename": filename,
		"size":     header.Size,
		"backend":  backend.Name(),
	})
}

// Download 下载文件
func (s *StorageController) Download(c *gin.Context) {
	backendName := c.Query("backend")
	path := c.Query("path")

	if path == "" {
		response.Error(c, http.StatusBadRequest, "缺少文件路径")
		return
	}

	backend, err := s.manager.Get(backendName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	reader, info, err := backend.Download(c.Request.Context(), path)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
	c.Header("Content-Type", info.ContentType)
	c.Header("Content-Length", strconv.FormatInt(info.Size, 10))

	io.Copy(c.Writer, reader)
}

// Delete 删除文件
func (s *StorageController) Delete(c *gin.Context) {
	backendName := c.Query("backend")
	path := c.Query("path")

	if path == "" {
		response.Error(c, http.StatusBadRequest, "缺少文件路径")
		return
	}

	backend, err := s.manager.Get(backendName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := backend.Delete(c.Request.Context(), path); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// List 列出文件
func (s *StorageController) List(c *gin.Context) {
	backendName := c.Query("backend")
	prefix := c.Query("prefix")

	backend, err := s.manager.Get(backendName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	files, err := backend.List(c.Request.Context(), prefix)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"files": files})
}

// GetURL 获取文件预签名URL
func (s *StorageController) GetURL(c *gin.Context) {
	backendName := c.Query("backend")
	path := c.Query("path")
	expiresStr := c.DefaultQuery("expires", "3600")

	if path == "" {
		response.Error(c, http.StatusBadRequest, "缺少文件路径")
		return
	}

	expires, _ := strconv.Atoi(expiresStr)
	if expires <= 0 {
		expires = 3600
	}

	backend, err := s.manager.Get(backendName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	url, err := backend.GetURL(c.Request.Context(), path, time.Duration(expires)*time.Second)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"url": url, "expires": expires})
}
