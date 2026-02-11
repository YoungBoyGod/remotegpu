package proxy

import (
	"context"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// ProxyService Proxy 节点业务逻辑层
type ProxyService struct {
	proxyDao *dao.ProxyDao
	db       *gorm.DB
}

func NewProxyService(db *gorm.DB) *ProxyService {
	return &ProxyService{
		proxyDao: dao.NewProxyDao(db),
		db:       db,
	}
}

// RegisterRequest Proxy 注册请求
type RegisterRequest struct {
	ID         string `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Host       string `json:"host" binding:"required"`
	APIPort    int    `json:"api_port"`
	HTTPPort   int    `json:"http_port"`
	RangeStart int    `json:"range_start"`
	RangeEnd   int    `json:"range_end"`
	Version    string `json:"version"`
}

// HeartbeatRequest Proxy 心跳请求
type HeartbeatRequest struct {
	ID             string `json:"id" binding:"required"`
	ActiveMappings int    `json:"active_mappings"`
	UsedPorts      int    `json:"used_ports"`
}

// Register 注册 Proxy 节点（upsert 语义）
func (s *ProxyService) Register(ctx context.Context, req *RegisterRequest) (*entity.ProxyNode, error) {
	now := time.Now()
	node := &entity.ProxyNode{
		ID:            req.ID,
		Name:          req.Name,
		Host:          req.Host,
		APIPort:       req.APIPort,
		HTTPPort:      req.HTTPPort,
		RangeStart:    req.RangeStart,
		RangeEnd:      req.RangeEnd,
		Version:       req.Version,
		Status:        "online",
		LastHeartbeat: &now,
	}

	if node.APIPort == 0 {
		node.APIPort = 9090
	}
	if node.HTTPPort == 0 {
		node.HTTPPort = 9091
	}
	if node.RangeStart == 0 {
		node.RangeStart = 20000
	}
	if node.RangeEnd == 0 {
		node.RangeEnd = 60000
	}

	if err := s.proxyDao.Upsert(ctx, node); err != nil {
		return nil, err
	}
	return node, nil
}

// Heartbeat 处理 Proxy 心跳
func (s *ProxyService) Heartbeat(ctx context.Context, req *HeartbeatRequest) error {
	return s.proxyDao.UpdateHeartbeat(ctx, req.ID, req.ActiveMappings, req.UsedPorts)
}

// ListNodes 列出所有 Proxy 节点
func (s *ProxyService) ListNodes(ctx context.Context) ([]entity.ProxyNode, error) {
	return s.proxyDao.List(ctx)
}

// GetNode 获取单个 Proxy 节点详情
func (s *ProxyService) GetNode(ctx context.Context, id string) (*entity.ProxyNode, error) {
	return s.proxyDao.FindByID(ctx, id)
}

// DeleteNode 删除 Proxy 节点
func (s *ProxyService) DeleteNode(ctx context.Context, id string) error {
	return s.proxyDao.Delete(ctx, id)
}

// ListMappings 列出所有端口映射（带 proxy 信息）
func (s *ProxyService) ListMappings(ctx context.Context) ([]entity.PortMapping, error) {
	var mappings []entity.PortMapping
	err := s.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("allocated_at desc").
		Find(&mappings).Error
	return mappings, err
}
