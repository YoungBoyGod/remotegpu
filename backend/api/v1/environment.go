package v1

// CreateEnvironmentRequest 创建环境请求
type CreateEnvironmentRequest struct {
	WorkspaceID *uint             `json:"workspace_id"`
	Name        string            `json:"name" binding:"required,max=128"`
	Description string            `json:"description"`
	Image       string            `json:"image" binding:"required"`
	CPU         int               `json:"cpu" binding:"required,min=1"`
	Memory      int64             `json:"memory" binding:"required,min=512"`
	GPU         int               `json:"gpu" binding:"min=0"`
	Storage     *int64            `json:"storage"`
	Env         map[string]string `json:"env"`
}
