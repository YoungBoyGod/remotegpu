package ops

import (
	"context"
	"encoding/json"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	"github.com/YoungBoyGod/remotegpu/internal/service/customer"
	"github.com/YoungBoyGod/remotegpu/internal/service/machine"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/prometheus"
)

type MonitorService struct {
	machineService *machine.MachineService
	cache          cache.Cache
	promClient     *prometheus.Client
}

func NewMonitorService(ms *machine.MachineService, c cache.Cache, prom *prometheus.Client) *MonitorService {
	return &MonitorService{
		machineService: ms,
		cache:          c,
		promClient:     prom,
	}
}

// GetGlobalSnapshot 获取全局监控快照
// @author Claude
// @description 获取系统实时监控数据，包括机器状态统计，支持 Redis 缓存
// @modified 2026-02-05
func (s *MonitorService) GetGlobalSnapshot(ctx context.Context) (map[string]interface{}, error) {
	const cacheKey = "monitor:snapshot"
	const cacheTTL = 30 * time.Second

	// 尝试从缓存读取
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, nil
			}
		}
	}

	// 缓存未命中，查询数据库
	stats, err := s.machineService.GetStatusStats(ctx)
	if err != nil {
		return nil, err
	}

	// 计算总数和在线数
	var totalMachines int64
	for _, count := range stats {
		totalMachines += count
	}
	onlineMachines := stats["idle"] + stats["allocated"]

	result := map[string]interface{}{
		"total_machines":     totalMachines,
		"online_machines":    onlineMachines,
		"idle_machines":      stats["idle"],
		"allocated_machines": stats["allocated"],
		"offline_machines":   stats["offline"],
		"avg_gpu_util":       s.getGPUUtilization(ctx),
	}

	// 写入缓存
	if s.cache != nil {
		if data, err := json.Marshal(result); err == nil {
			_ = s.cache.Set(ctx, cacheKey, string(data), cacheTTL)
		}
	}

	return result, nil
}

func (s *MonitorService) GetGPUTrend(ctx context.Context) ([]map[string]interface{}, error) {
	if s.promClient != nil {
		trend, err := s.promClient.GetGPUTrend(ctx, 24*time.Hour)
		if err == nil && len(trend) > 0 {
			result := make([]map[string]interface{}, len(trend))
			for i, p := range trend {
				result[i] = map[string]interface{}{
					"time":  p.Time,
					"value": p.Usage,
				}
			}
			return result, nil
		}
	}
	// 返回默认数据
	return []map[string]interface{}{
		{"time": "00:00", "value": 0},
		{"time": "04:00", "value": 0},
		{"time": "08:00", "value": 0},
		{"time": "12:00", "value": 0},
		{"time": "16:00", "value": 0},
		{"time": "20:00", "value": 0},
	}, nil
}

// getGPUUtilization 从 Prometheus 获取 GPU 平均利用率
func (s *MonitorService) getGPUUtilization(ctx context.Context) float64 {
	if s.promClient != nil {
		metrics, err := s.promClient.GetGPUUtilization(ctx)
		if err == nil {
			return metrics.AvgUtilization
		}
	}
	return 0.0
}

type DashboardService struct {
	machineService    *machine.MachineService
	customerService   *customer.CustomerService
	allocationService *allocation.AllocationService
	promClient        *prometheus.Client
}

func NewDashboardService(ms *machine.MachineService, cs *customer.CustomerService, as *allocation.AllocationService, prom *prometheus.Client) *DashboardService {
	return &DashboardService{
		machineService:    ms,
		customerService:   cs,
		allocationService: as,
		promClient:        prom,
	}
}

// GetAggregatedStats 获取仪表盘聚合统计数据
// @description 聚合机器状态、客户数量等核心指标
// @modified 2026-02-04
func (s *DashboardService) GetAggregatedStats(ctx context.Context) (map[string]any, error) {
	// 获取机器状态统计
	machineStats, err := s.machineService.GetStatusStats(ctx)
	if err != nil {
		return nil, err
	}

	// 计算机器总数
	var totalMachines int64
	for _, count := range machineStats {
		totalMachines += count
	}

	// 获取活跃客户数
	activeCustomers, err := s.customerService.CountActive(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"total_machines":     totalMachines,
		"allocated_machines": machineStats["allocated"],
		"idle_machines":      machineStats["idle"],
		"offline_machines":   machineStats["offline"],
		"active_customers":   activeCustomers,
	}, nil
}



// GetGPUTrend 获取 GPU 使用趋势数据
// @description 返回最近时间段的 GPU 利用率趋势，从 Prometheus 获取
// @modified 2026-02-05
func (s *DashboardService) GetGPUTrend(ctx context.Context) ([]map[string]any, error) {
	if s.promClient != nil {
		trend, err := s.promClient.GetGPUTrend(ctx, 24*time.Hour)
		if err == nil && len(trend) > 0 {
			result := make([]map[string]any, len(trend))
			for i, p := range trend {
				result[i] = map[string]any{
					"time":  p.Time,
					"usage": p.Usage,
				}
			}
			return result, nil
		}
	}
	// 返回默认数据
	return []map[string]any{
		{"time": "00:00", "usage": 0},
		{"time": "04:00", "usage": 0},
		{"time": "08:00", "usage": 0},
		{"time": "12:00", "usage": 0},
		{"time": "16:00", "usage": 0},
		{"time": "20:00", "usage": 0},
	}, nil
}



// GetRecentAllocations 获取最近的分配记录
// @description 用于仪表盘展示最近的机器分配情况
// @modified 2026-02-04
func (s *DashboardService) GetRecentAllocations(ctx context.Context) ([]entity.Allocation, error) {
	return s.allocationService.GetRecent(ctx)
}
