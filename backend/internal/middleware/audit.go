package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/audit"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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

		// 构建请求参数摘要
		var detail datatypes.JSON
		if len(bodyBytes) > 0 {
			detail = buildRequestSummary(bodyBytes)
		}

		log := &entity.AuditLog{
			CustomerID:   customerID,
			Username:     usernameStr,
			IPAddress:    c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Detail:       detail,
			StatusCode:   c.Writer.Status(),
			CreatedAt:    time.Now(),
		}

		_ = auditSvc.LogAction(c, log)
	}
}

func parseAction(method, path string) string {
	// 细粒度操作解析：根据路径后缀区分具体操作
	if method == "POST" {
		switch {
		case strings.HasSuffix(path, "/allocate"):
			return "allocate"
		case strings.HasSuffix(path, "/reclaim"):
			return "reclaim"
		case strings.HasSuffix(path, "/disable"):
			return "disable"
		case strings.HasSuffix(path, "/enable"):
			return "enable"
		case strings.HasSuffix(path, "/stop"):
			return "stop"
		case strings.HasSuffix(path, "/cancel"):
			return "cancel"
		case strings.HasSuffix(path, "/retry"):
			return "retry"
		case strings.HasSuffix(path, "/collect"):
			return "collect"
		case strings.HasSuffix(path, "/maintenance"):
			return "maintenance"
		case strings.HasSuffix(path, "/acknowledge"):
			return "acknowledge"
		case strings.HasSuffix(path, "/sync"):
			return "sync"
		case strings.HasSuffix(path, "/import"):
			return "import"
		default:
			return "create"
		}
	}

	switch method {
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

	resourceType := "unknown"
	switch {
	case strings.Contains(path, "machines"):
		resourceType = "machine"
	case strings.Contains(path, "customers"):
		resourceType = "customer"
	case strings.Contains(path, "keys"):
		resourceType = "ssh_key"
	case strings.Contains(path, "datasets"):
		resourceType = "dataset"
	case strings.Contains(path, "tasks"):
		resourceType = "task"
	case strings.Contains(path, "allocations"):
		resourceType = "allocation"
	case strings.Contains(path, "documents"):
		resourceType = "document"
	case strings.Contains(path, "alerts"):
		resourceType = "alert"
	case strings.Contains(path, "images"):
		resourceType = "image"
	case strings.Contains(path, "settings"):
		resourceType = "system_config"
	case strings.Contains(path, "storage"):
		resourceType = "storage"
	}

	return resourceType, resourceID
}

// buildRequestSummary 构建请求参数摘要，过滤敏感字段
func buildRequestSummary(body []byte) datatypes.JSON {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	// 过滤敏感字段
	sensitiveKeys := []string{"password", "ssh_password", "ssh_key", "token", "secret"}
	for _, key := range sensitiveKeys {
		if _, ok := raw[key]; ok {
			raw[key] = "***"
		}
	}

	// 截断过长的值
	for k, v := range raw {
		if s, ok := v.(string); ok && len(s) > 200 {
			raw[k] = s[:200] + "..."
		}
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	return datatypes.JSON(data)
}
