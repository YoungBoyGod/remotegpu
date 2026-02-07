package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/pkg/response"
	"github.com/gin-gonic/gin"
)

// AgentAuth Agent Token 认证中间件
// 从请求头 X-Agent-Token 读取 token，与配置中的 agent.token 比对
func AgentAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Agent-Token")
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
