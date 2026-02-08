package v1

// CreateWorkspaceRequest 创建工作空间请求
type CreateWorkspaceRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
}

// UpdateWorkspaceRequest 更新工作空间请求
type UpdateWorkspaceRequest struct {
	Name        string `json:"name" binding:"omitempty,max=128"`
	Description string `json:"description"`
}

// AddMemberRequest 添加工作空间成员请求
type AddMemberRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"omitempty,oneof=admin member viewer"`
}
