package machine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/pkg/cache"
)

const (
	// hostStatusKeyPrefix Redis key 前缀
	hostStatusKeyPrefix = "host:status:"
	// hostStatusTTL 设备状态缓存 TTL，超过此时间未续期视为离线
	hostStatusTTL = 90 * time.Second
)

// HostStatusCache 设备实时状态缓存
type HostStatusCache struct {
	cache cache.Cache
}

// CachedHostStatus 缓存中的设备状态
type CachedHostStatus struct {
	HostID         string    `json:"host_id"`
	DeviceStatus   string    `json:"device_status"`
	LastHeartbeat  time.Time `json:"last_heartbeat"`
	CPUUsage       *float64  `json:"cpu_usage,omitempty"`
	MemoryUsage    *float64  `json:"memory_usage,omitempty"`
	DiskUsage      *float64  `json:"disk_usage,omitempty"`
	GPUCount       int       `json:"gpu_count,omitempty"`
}

// NewHostStatusCache 创建设备状态缓存
func NewHostStatusCache(c cache.Cache) *HostStatusCache {
	return &HostStatusCache{cache: c}
}

func hostStatusKey(hostID string) string {
	return hostStatusKeyPrefix + hostID
}

// SetOnline 心跳到达时写入 Redis，刷新 TTL
func (h *HostStatusCache) SetOnline(ctx context.Context, status *CachedHostStatus) error {
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("序列化设备状态失败: %w", err)
	}
	return h.cache.Set(ctx, hostStatusKey(status.HostID), string(data), hostStatusTTL)
}

// Get 从 Redis 读取设备实时状态，key 不存在返回 nil
func (h *HostStatusCache) Get(ctx context.Context, hostID string) (*CachedHostStatus, error) {
	val, err := h.cache.Get(ctx, hostStatusKey(hostID))
	if err != nil {
		return nil, err
	}
	var status CachedHostStatus
	if err := json.Unmarshal([]byte(val), &status); err != nil {
		return nil, fmt.Errorf("反序列化设备状态失败: %w", err)
	}
	return &status, nil
}

// IsOnline 检查设备是否在线（key 存在即在线）
func (h *HostStatusCache) IsOnline(ctx context.Context, hostID string) bool {
	count, err := h.cache.Exists(ctx, hostStatusKey(hostID))
	return err == nil && count > 0
}

// Delete 删除设备状态缓存
func (h *HostStatusCache) Delete(ctx context.Context, hostID string) error {
	return h.cache.Delete(ctx, hostStatusKey(hostID))
}
