package machine

import (
	"context"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"gorm.io/gorm"
)

// HostStatusSyncer 定时将 Redis 中的设备状态同步到 PostgreSQL
type HostStatusSyncer struct {
	machineDao      *dao.MachineDao
	statusCache     *HostStatusCache
	syncInterval    time.Duration
}

// NewHostStatusSyncer 创建同步器
func NewHostStatusSyncer(db *gorm.DB, statusCache *HostStatusCache, syncInterval time.Duration) *HostStatusSyncer {
	return &HostStatusSyncer{
		machineDao:   dao.NewMachineDao(db),
		statusCache:  statusCache,
		syncInterval: syncInterval,
	}
}

// Start 启动定时同步
func (s *HostStatusSyncer) Start(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	// 启动时立即执行一次同步，避免重启后状态不一致
	s.syncAll(ctx)

	logger.GetLogger().Info("设备状态同步服务已启动")

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger().Info("设备状态同步服务已停止")
			return
		case <-ticker.C:
			s.syncAll(ctx)
		}
	}
}

// syncAll 遍历所有机器，将 Redis 状态同步到数据库
func (s *HostStatusSyncer) syncAll(ctx context.Context) {
	hosts, err := s.machineDao.ListAll(ctx)
	if err != nil {
		logger.GetLogger().Error("同步设备状态: 查询机器列表失败: " + err.Error())
		return
	}

	for _, host := range hosts {
		cached, err := s.statusCache.Get(ctx, host.ID)
		if err != nil {
			// key 不存在说明 TTL 已过期，标记离线
			if host.DeviceStatus == "online" {
				_ = s.machineDao.UpdateDeviceStatus(ctx, host.ID, "offline")
			}
			continue
		}

		// 将 Redis 中的心跳时间和状态同步到数据库
		_ = s.machineDao.UpdateHeartbeat(ctx, host.ID, cached.DeviceStatus)
	}
}
