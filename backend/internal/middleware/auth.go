package middleware

import (
	"net/http"
	"strings"

	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

const tokenBlacklistPrefix = "auth:token:blacklist:"

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "请提供认证令牌")
			c.Abort()
			return
		}

		// 验证 token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "无效的认证令牌")
			c.Abort()
			return
		}

		// 检查 token 是否在黑名单中
		if isTokenBlacklisted(c, parts[1]) {
			response.Error(c, http.StatusUnauthorized, "令牌已失效，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息存入上下文 (Fixed Key Name)
		c.Set("userID", claims.UserID) // Changed from user_id
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// isTokenBlacklisted 检查 token 是否在黑名单中
func isTokenBlacklisted(c *gin.Context, token string) bool {
	cacheClient := cache.GetCache()
	if cacheClient == nil {
		return false
	}
	key := tokenBlacklistPrefix + token
	count, err := cacheClient.Exists(c, key)
	if err != nil {
		return false
	}
	return count > 0
}