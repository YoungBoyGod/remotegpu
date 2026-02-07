package v1

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Type       string `json:"type" binding:"required"`
	Level      string `json:"level"`
}
