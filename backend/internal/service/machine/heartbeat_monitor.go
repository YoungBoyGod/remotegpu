package machine

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/middleware"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"gorm.io/gorm"
)

// HeartbeatMonitor 心跳监控服务
// @author Claude
// @description 监控机器心跳，自动检测离线机器
// @modified 2026-02-06
type HeartbeatMonitor struct {
	db          *gorm.DB
	machineDao  *dao.MachineDao
	timeout     time.Duration // 心跳超时时间
	checkInterval time.Duration // 检查间隔
}

// NewHeartbeatMonitor 创建心跳监控服务
func NewHeartbeatMonitor(db *gorm.DB, timeout, checkInterval time.Duration) *HeartbeatMonitor {
	return &HeartbeatMonitor{
		db:            db,
		machineDao:    dao.NewMachineDao(db),
		timeout:       timeout,
		checkInterval: checkInterval,
	}
}

// Start 启动心跳监控
func (m *HeartbeatMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()

	logger.GetLogger().Info("心跳监控服务已启动")

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger().Info("心跳监控服务已停止")
			return
		case <-ticker.C:
			m.checkOfflineHosts(ctx)
		}
	}
}

// checkOfflineHosts 检查离线机器
func (m *HeartbeatMonitor) checkOfflineHosts(ctx context.Context) {
	// 更新在线机器数指标
	onlineHosts, err := m.machineDao.ListOnline(ctx)
	if err == nil {
		middleware.MachinesOnline.Set(float64(len(onlineHosts)))
	}

	// 计算超时时间点
	timeoutAt := time.Now().Add(-m.timeout)

	// 查询所有在线但心跳超时的机器
	var hosts []struct {
		ID            string
		Name          string
		Status        string
		LastHeartbeat *time.Time
	}

	err = m.db.WithContext(ctx).
		Model(&entity.Host{}).
		Select("id, name, status, last_heartbeat").
		Where("status IN ?", []string{"online", "idle"}).
		Where("last_heartbeat IS NOT NULL AND last_heartbeat < ?", timeoutAt).
		Find(&hosts).Error

	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("查询超时机器失败: %v", err))
		return
	}

	if len(hosts) == 0 {
		return
	}

	logger.GetLogger().Info(fmt.Sprintf("发现 %d 台心跳超时的机器", len(hosts)))

	// 批量更新为 offline 状态
	for _, host := range hosts {
		err := m.machineDao.UpdateStatus(ctx, host.ID, "offline")
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("更新机器 %s 状态失败: %v", host.ID, err))
			continue
		}

		logger.GetLogger().Info(fmt.Sprintf("机器 %s (%s) 心跳超时，已标记为 offline", host.ID, host.Name))
	}
}
