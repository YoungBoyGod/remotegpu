package auth

// 角色常量定义
const (
	RoleAdmin      = "admin"      // 管理员
	RoleUser       = "user"       // 普通用户
	RoleEnterprise = "enterprise" // 企业用户
)

// IsValidRole 检查角色是否有效
func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleEnterprise:
		return true
	default:
		return false
	}
}

// IsAdmin 检查是否为管理员
func IsAdmin(role string) bool {
	return role == RoleAdmin
}

// IsEnterprise 检查是否为企业用户
func IsEnterprise(role string) bool {
	return role == RoleEnterprise
}
