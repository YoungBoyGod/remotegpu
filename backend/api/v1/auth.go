package v1

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	ExpiresIn          int64  `json:"expires_in"`
	MustChangePassword bool   `json:"must_change_password"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// RequestPasswordResetRequest 请求密码重置
type RequestPasswordResetRequest struct {
	Username string `json:"username" binding:"required"`
}

// ConfirmPasswordResetRequest 确认密码重置
type ConfirmPasswordResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone"`
	Company     string `json:"company"`
	FullName    string `json:"full_name"`
	AvatarURL   string `json:"avatar_url"`
}
