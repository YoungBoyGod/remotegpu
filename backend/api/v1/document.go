package v1

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	Title    string `json:"title" binding:"required"`
	Category string `json:"category"`
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}
