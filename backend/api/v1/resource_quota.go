package v1

import "time"

// SetQuotaRequest 设置资源配额请求
type SetQuotaRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	MaxGPU          int    `json:"max_gpu" binding:"min=0"`
	MaxCPU          int    `json:"max_cpu" binding:"min=0"`
	MaxMemory       int64  `json:"max_memory" binding:"min=0"`
	MaxStorage      int64  `json:"max_storage" binding:"min=0"`
	MaxEnvironments int    `json:"max_environments" binding:"min=0"`
	QuotaLevel      string `json:"quota_level"`
}

// UpdateQuotaRequest 更新资源配额请求
type UpdateQuotaRequest struct {
	MaxGPU          int    `json:"max_gpu" binding:"min=0"`
	MaxCPU          int    `json:"max_cpu" binding:"min=0"`
	MaxMemory       int64  `json:"max_memory" binding:"min=0"`
	MaxStorage      int64  `json:"max_storage" binding:"min=0"`
	MaxEnvironments int    `json:"max_environments" binding:"min=0"`
	QuotaLevel      string `json:"quota_level"`
}

// QuotaInfo 资源配额信息
type QuotaInfo struct {
	ID              uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	QuotaLevel      string    `json:"quota_level"`
	MaxGPU          int       `json:"max_gpu"`
	MaxCPU          int       `json:"max_cpu"`
	MaxMemory       int64     `json:"max_memory"`
	MaxStorage      int64     `json:"max_storage"`
	MaxEnvironments int       `json:"max_environments"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// QuotaUsageResponse 配额使用情况响应
type QuotaUsageResponse struct {
	Quota            QuotaDetail            `json:"quota"`
	Used             UsedResources          `json:"used"`
	Available        AvailableResources     `json:"available"`
	UsagePercentage  UsagePercentageDetail  `json:"usage_percentage"`
}

// QuotaDetail 配额详情
type QuotaDetail struct {
	MaxGPU          int   `json:"max_gpu"`
	MaxCPU          int   `json:"max_cpu"`
	MaxMemory       int64 `json:"max_memory"`
	MaxStorage      int64 `json:"max_storage"`
	MaxEnvironments int   `json:"max_environments"`
}

// UsedResources 已使用资源
type UsedResources struct {
	UsedGPU          int   `json:"used_gpu"`
	UsedCPU          int   `json:"used_cpu"`
	UsedMemory       int64 `json:"used_memory"`
	UsedStorage      int64 `json:"used_storage"`
	UsedEnvironments int   `json:"used_environments"`
}

// AvailableResources 可用资源
type AvailableResources struct {
	AvailableGPU          int   `json:"available_gpu"`
	AvailableCPU          int   `json:"available_cpu"`
	AvailableMemory       int64 `json:"available_memory"`
	AvailableStorage      int64 `json:"available_storage"`
	AvailableEnvironments int   `json:"available_environments"`
}

// UsagePercentageDetail 使用百分比详情
type UsagePercentageDetail struct {
	GPU          float64 `json:"gpu"`
	CPU          float64 `json:"cpu"`
	Memory       float64 `json:"memory"`
	Storage      float64 `json:"storage"`
	Environments float64 `json:"environments"`
}

// QuotaListResponse 配额列表响应
type QuotaListResponse struct {
	Items    []*QuotaInfo `json:"items"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}
