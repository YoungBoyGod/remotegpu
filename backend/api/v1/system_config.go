package v1

// UpdateSystemConfigsRequest 批量更新系统配置请求
type UpdateSystemConfigsRequest struct {
	Configs map[string]string `json:"configs" binding:"required"`
}

// CreateSystemConfigRequest 创建系统配置请求
type CreateSystemConfigRequest struct {
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
	ConfigType  string `json:"config_type"`
	ConfigGroup string `json:"config_group"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateSystemConfigRequest 更新单条系统配置请求
type UpdateSystemConfigRequest struct {
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	ConfigGroup string `json:"config_group"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"is_public"`
}
