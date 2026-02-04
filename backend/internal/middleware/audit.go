package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"github.com/gin-gonic/gin"
)

// AuditMiddleware 审计中间件
func AuditMiddleware(auditSvc *audit.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录写操作
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		// 记录请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		// 获取用户信息
		var customerID *uint
		if id, exists := c.Get("userID"); exists {
			uid := id.(uint)
			customerID = &uid
		}

		username, _ := c.Get("username")
		usernameStr, _ := username.(string)

		// 解析操作类型
		action := parseAction(c.Request.Method, c.FullPath())
		resourceType, resourceID := parseResource(c.FullPath(), c.Params)

		log := &entity.AuditLog{
			CustomerID:   customerID,
			Username:     usernameStr,
			IPAddress:    c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			StatusCode:   c.Writer.Status(),
			CreatedAt:    time.Now(),
		}

		_ = auditSvc.LogAction(c, log)
	}
}

func parseAction(method, path string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return method
	}
}

func parseResource(path string, params gin.Params) (string, string) {
	resourceID := ""
	if id := params.ByName("id"); id != "" {
		resourceID = id
	}

	// 简单解析资源类型
	resourceType := "unknown"
	switch {
	case contains(path, "machines"):
		resourceType = "machine"
	case contains(path, "customers"):
		resourceType = "customer"
	case contains(path, "keys"):
		resourceType = "ssh_key"
	case contains(path, "datasets"):
		resourceType = "dataset"
	case contains(path, "tasks"):
		resourceType = "task"
	}

	return resourceType, resourceID
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
