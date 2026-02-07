package machine

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"gorm.io/gorm"
)

// MetricsCollector 监控数据采集器
// @author Claude
// @description 定期采集机器的监控数据
// @modified 2026-02-06
type MetricsCollector struct {
	db             *gorm.DB
	metricDao      *dao.HostMetricDao
	machineDao     *dao.MachineDao
	interval       time.Duration
	retentionDays  int
	stopCh         chan struct{}
}

// NewMetricsCollector 创建监控数据采集器
func NewMetricsCollector(db *gorm.DB, interval time.Duration, retentionDays int) *MetricsCollector {
	return &MetricsCollector{
		db:            db,
		metricDao:     dao.NewHostMetricDao(db),
		machineDao:    dao.NewMachineDao(db),
		interval:      interval,
		retentionDays: retentionDays,
		stopCh:        make(chan struct{}),
	}
}

// Start 启动采集器
func (c *MetricsCollector) Start(ctx context.Context) {
	logger.GetLogger().Info("监控数据采集器已启动")
	go c.collectLoop(ctx)
	go c.cleanupLoop(ctx)
}

// Stop 停止采集器
func (c *MetricsCollector) Stop() {
	close(c.stopCh)
}

func (c *MetricsCollector) collectLoop(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.collectMetrics(ctx)
		}
	}
}

func (c *MetricsCollector) collectMetrics(ctx context.Context) {
	// 获取所有在线的机器
	hosts, err := c.machineDao.ListOnline(ctx)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("查询在线机器失败: %v", err))
		return
	}

	if len(hosts) == 0 {
		return
	}

	logger.GetLogger().Info(fmt.Sprintf("开始采集 %d 台机器的监控数据", len(hosts)))

	// TODO: 这里暂时生成模拟数据，后续需要从Agent获取真实数据
	metrics := make([]*entity.HostMetric, 0, len(hosts))
	now := time.Now()

	for _, host := range hosts {
		metric := c.generateMockMetric(host.ID, now)
		metrics = append(metrics, metric)
	}

	// 批量保存监控数据
	if err := c.metricDao.BatchCreate(ctx, metrics); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("保存监控数据失败: %v", err))
		return
	}

	logger.GetLogger().Info(fmt.Sprintf("成功采集 %d 条监控数据", len(metrics)))
}

// generateMockMetric 生成模拟监控数据（临时使用）
func (c *MetricsCollector) generateMockMetric(hostID string, collectedAt time.Time) *entity.HostMetric {
	cpuUsage := 45.5
	memUsage := 60.2
	diskUsage := 35.8

	return &entity.HostMetric{
		HostID:             hostID,
		CPUUsagePercent:    &cpuUsage,
		MemoryUsagePercent: &memUsage,
		DiskUsagePercent:   &diskUsage,
		CollectedAt:        collectedAt,
	}
}

func (c *MetricsCollector) cleanupLoop(ctx context.Context) {
	// 每天清理一次旧数据
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.cleanupOldMetrics(ctx)
		}
	}
}

func (c *MetricsCollector) cleanupOldMetrics(ctx context.Context) {
	before := time.Now().AddDate(0, 0, -c.retentionDays)
	err := c.metricDao.DeleteOldRecords(ctx, before)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("清理旧监控数据失败: %v", err))
		return
	}
	logger.GetLogger().Info(fmt.Sprintf("已清理 %d 天前的监控数据", c.retentionDays))
}
