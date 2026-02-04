package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwagger 设置Swagger文档路由
func SetupSwagger(r *gin.Engine) error {
	if !config.GlobalConfig.Swagger.Enabled {
		return nil
	}

	// 读取OpenAPI文档
	openapiPath := filepath.Join("docs", "openapi.yaml")
	openapiContent, err := os.ReadFile(openapiPath)
	if err != nil {
		return fmt.Errorf("读取OpenAPI文档失败: %w", err)
	}

	openapiRoute := config.GlobalConfig.Swagger.OpenAPIPath
	if openapiRoute == "" {
		openapiRoute = "/openapi.yaml"
	}
	if !strings.HasPrefix(openapiRoute, "/") {
		openapiRoute = "/" + openapiRoute
	}

	// 提供OpenAPI YAML文件（动态替换服务器地址）
	r.GET(openapiRoute, func(c *gin.Context) {
		// 获取请求的host
		requestHost := c.Request.Host

		// 动态替换服务器地址
		content := string(openapiContent)
		// 简单的字符串替换，将文档中的服务器地址替换为实际的地址
		content = strings.ReplaceAll(content, "http://0.0.0.:8080", fmt.Sprintf("http://%s", requestHost))
		content = strings.ReplaceAll(content, "https://0.0.0.0:8080", fmt.Sprintf("https://%s", requestHost))

		c.Header("Content-Type", "application/x-yaml")
		c.String(http.StatusOK, content)
	})

	// 配置Swagger UI
	url := ginSwagger.URL(openapiRoute)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 重定向根路径到swagger
	if config.GlobalConfig.Server.Mode == "debug" {
		r.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		})
	}

	logHost := config.GlobalConfig.Server.Host
	if logHost == "" {
		logHost = "0.0.0.0"
	}
	if _, _, err := net.SplitHostPort(logHost); err != nil {
		logHost = net.JoinHostPort(logHost, strconv.Itoa(config.GlobalConfig.Server.Port))
	}

	fmt.Printf("Swagger UI 已启动: http://%s/swagger/index.html\n", logHost)
	fmt.Printf("OpenAPI 已启动: http://%s%s\n", logHost, openapiRoute)
	return nil
}
