package auth

// 角色常量定义
const (
	RoleAdmin          = "admin"
	RoleCustomerOwner  = "customer_owner"
	RoleCustomerMember = "customer_member"
)

// IsValidRole 检查角色是否有效
func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleCustomerOwner, RoleCustomerMember:
		return true
	default:
		return false
	}
}

// IsAdmin 检查是否为管理员
func IsAdmin(role string) bool {
	return role == RoleAdmin
}