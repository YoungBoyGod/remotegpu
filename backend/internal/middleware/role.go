package middleware

import (
	"net/http"

	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequireRole 角色权限中间件
// 使用方式: RequireRole(auth.RoleAdmin) 或 RequireRole(auth.RoleAdmin, auth.RoleEnterprise)
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "未找到用户角色信息")
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "角色信息格式错误")
			c.Abort()
			return
		}

		// 检查角色是否在允许列表中
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		// 角色不匹配
		response.Error(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}

// RequireAdmin 要求管理员权限的中间件（快捷方式）
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(auth.RoleAdmin)
}
