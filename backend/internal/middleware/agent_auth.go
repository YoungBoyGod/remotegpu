package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// AgentAuth Agent Token 认证中间件
// 同时支持 X-Agent-Token 和 Authorization: Bearer 两种认证方式
func AgentAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Agent-Token")
		// 兼容 Agent 客户端使用 Authorization: Bearer 发送 Token
		if token == "" {
			if auth := c.GetHeader("Authorization"); len(auth) > 7 && auth[:7] == "Bearer " {
				token = auth[7:]
			}
		}
		if token == "" {
			response.Error(c, http.StatusUnauthorized, "缺少 Agent 认证令牌")
			c.Abort()
			return
		}

		expectedToken := ""
		if config.GlobalConfig != nil {
			expectedToken = config.GlobalConfig.Agent.Token
		}

		if expectedToken == "" || subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) != 1 {
			response.Error(c, http.StatusUnauthorized, "无效的 Agent 认证令牌")
			c.Abort()
			return
		}

		c.Next()
	}
}
