package v1

// UpdateSystemConfigsRequest 批量更新系统配置请求
type UpdateSystemConfigsRequest struct {
	Configs map[string]string `json:"configs" binding:"required"`
}
