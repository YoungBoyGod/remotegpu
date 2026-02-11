package forwarder

import (
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// HTTPProxy HTTP 反向代理，支持 WebSocket
type HTTPProxy struct {
	port     int
	routes   map[int]*httputil.ReverseProxy // externalPort -> proxy
	mu       sync.RWMutex
	server   *http.Server
	stopCh   chan struct{}
}

// NewHTTPProxy 创建 HTTP 反向代理
func NewHTTPProxy(port int) *HTTPProxy {
	return &HTTPProxy{
		port:   port,
		routes: make(map[int]*httputil.ReverseProxy),
		stopCh: make(chan struct{}),
	}
}

// AddRoute 添加路由，将 externalPort 映射到目标地址
func (h *HTTPProxy) AddRoute(externalPort int, targetHost string, targetPort int) {
	target := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(targetHost, strconv.Itoa(targetPort)),
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("HTTP 代理错误", "target", target.Host, "error", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	h.mu.Lock()
	h.routes[externalPort] = proxy
	h.mu.Unlock()
}

// RemoveRoute 移除路由
func (h *HTTPProxy) RemoveRoute(externalPort int) {
	h.mu.Lock()
	delete(h.routes, externalPort)
	h.mu.Unlock()
}

// Start 启动 HTTP 反向代理服务
func (h *HTTPProxy) Start() error {
	h.server = &http.Server{
		Addr:         ":" + strconv.Itoa(h.port),
		Handler:      h,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		slog.Info("HTTP 反向代理已启动", "port", h.port)
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP 反向代理异常退出", "error", err)
		}
	}()
	return nil
}

// Stop 停止 HTTP 反向代理服务
func (h *HTTPProxy) Stop() {
	if h.server != nil {
		h.server.Close()
	}
	slog.Info("HTTP 反向代理已停止")
}

// ServeHTTP 实现 http.Handler 接口，路径格式: /proxy/{port}/*
func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 解析路径中的端口号: /proxy/{port}/...
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/proxy/"), "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "invalid proxy path", http.StatusBadRequest)
		return
	}

	port, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid port number", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	proxy, ok := h.routes[port]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "no route for port", http.StatusNotFound)
		return
	}

	// 重写路径，去掉 /proxy/{port} 前缀
	if len(parts) > 1 {
		r.URL.Path = "/" + parts[1]
	} else {
		r.URL.Path = "/"
	}

	// WebSocket 升级检测
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		h.handleWebSocket(w, r, port)
		return
	}

	proxy.ServeHTTP(w, r)
}

// handleWebSocket 处理 WebSocket 升级请求，使用 hijack 实现双向转发
func (h *HTTPProxy) handleWebSocket(w http.ResponseWriter, r *http.Request, port int) {
	h.mu.RLock()
	proxy, ok := h.routes[port]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "no route for port", http.StatusNotFound)
		return
	}

	// 对于 WebSocket，直接使用 ReverseProxy 转发
	// httputil.ReverseProxy 默认支持 WebSocket 升级（Go 1.12+）
	proxy.ServeHTTP(w, r)
}
