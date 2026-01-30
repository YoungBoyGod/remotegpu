package graceful

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPServer HTTP服务器包装器
type HTTPServer struct {
	engine *gin.Engine
	server *http.Server
	addr   string
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(engine *gin.Engine, addr string) *HTTPServer {
	return &HTTPServer{
		engine: engine,
		addr:   addr,
	}
}

// Start 启动服务器
func (s *HTTPServer) Start() error {
	s.server = &http.Server{
		Addr:           s.addr,
		Handler:        s.engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("服务启动在 %s\n", s.addr)
	return s.server.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}
