package notification

import (
	"fmt"
	"sync"
)

// SSEEvent SSE 事件
type SSEEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// SSEClient SSE 客户端连接
type SSEClient struct {
	CustomerID uint
	Channel    chan SSEEvent
}

// SSEHub 管理所有 SSE 连接
type SSEHub struct {
	mu      sync.RWMutex
	clients map[uint]map[*SSEClient]struct{} // customerID -> clients
}

// NewSSEHub 创建 SSE Hub
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[uint]map[*SSEClient]struct{}),
	}
}

// Register 注册客户端
func (h *SSEHub) Register(client *SSEClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.CustomerID] == nil {
		h.clients[client.CustomerID] = make(map[*SSEClient]struct{})
	}
	h.clients[client.CustomerID][client] = struct{}{}
}

// Unregister 注销客户端
func (h *SSEHub) Unregister(client *SSEClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[client.CustomerID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.clients, client.CustomerID)
		}
	}
	close(client.Channel)
}

// Send 向指定用户推送事件
func (h *SSEHub) Send(customerID uint, event SSEEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[customerID]
	if !ok {
		return
	}
	for client := range clients {
		select {
		case client.Channel <- event:
		default:
			// 通道满则跳过，避免阻塞
		}
	}
}

// Broadcast 向所有在线用户广播事件
func (h *SSEHub) Broadcast(event SSEEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.clients {
		for client := range clients {
			select {
			case client.Channel <- event:
			default:
			}
		}
	}
}

// OnlineCount 获取指定用户的在线连接数
func (h *SSEHub) OnlineCount(customerID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clients[customerID])
}

// FormatSSE 格式化 SSE 数据
func FormatSSE(event SSEEvent) string {
	return fmt.Sprintf("event: %s\ndata: %s\n\n", event.Event, event.Data)
}
