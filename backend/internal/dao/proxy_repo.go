package dao

import (
	"context"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProxyDao Proxy 节点数据访问层
type ProxyDao struct {
	db *gorm.DB
}

func NewProxyDao(db *gorm.DB) *ProxyDao {
	return &ProxyDao{db: db}
}

// Create 创建 Proxy 节点
func (d *ProxyDao) Create(ctx context.Context, node *entity.ProxyNode) error {
	return d.db.WithContext(ctx).Create(node).Error
}

// Upsert 创建或更新 Proxy 节点（注册时使用）
func (d *ProxyDao) Upsert(ctx context.Context, node *entity.ProxyNode) error {
	return d.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "host", "api_port", "http_port",
			"range_start", "range_end", "version",
			"status", "last_heartbeat", "updated_at",
		}),
	}).Create(node).Error
}

// Update 更新 Proxy 节点
func (d *ProxyDao) Update(ctx context.Context, node *entity.ProxyNode) error {
	return d.db.WithContext(ctx).Save(node).Error
}

// Delete 删除 Proxy 节点
func (d *ProxyDao) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&entity.ProxyNode{}, "id = ?", id).Error
}

// FindByID 根据 ID 查找 Proxy 节点
func (d *ProxyDao) FindByID(ctx context.Context, id string) (*entity.ProxyNode, error) {
	var node entity.ProxyNode
	if err := d.db.WithContext(ctx).First(&node, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

// List 获取所有 Proxy 节点
func (d *ProxyDao) List(ctx context.Context) ([]entity.ProxyNode, error) {
	var nodes []entity.ProxyNode
	err := d.db.WithContext(ctx).Order("created_at desc").Find(&nodes).Error
	return nodes, err
}

// UpdateHeartbeat 更新心跳时间和状态
func (d *ProxyDao) UpdateHeartbeat(ctx context.Context, id string, activeMappings, usedPorts int) error {
	now := time.Now()
	return d.db.WithContext(ctx).Model(&entity.ProxyNode{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_heartbeat":  now,
			"status":          "online",
			"active_mappings": activeMappings,
			"used_ports":      usedPorts,
			"updated_at":      now,
		}).Error
}

// UpdateStatus 更新 Proxy 节点状态
func (d *ProxyDao) UpdateStatus(ctx context.Context, id string, status string) error {
	return d.db.WithContext(ctx).Model(&entity.ProxyNode{}).Where("id = ?", id).
		Update("status", status).Error
}
