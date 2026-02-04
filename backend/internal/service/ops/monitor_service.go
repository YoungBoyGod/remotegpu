package ops

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	"github.com/YoungBoyGod/remotegpu/internal/service/customer"
	"github.com/YoungBoyGod/remotegpu/internal/service/machine"
)

type MonitorService struct {
	machineService *machine.MachineService
}

func NewMonitorService(ms *machine.MachineService) *MonitorService {
	return &MonitorService{
		machineService: ms,
	}
}

// GetGlobalSnapshot 获取全局监控快照
// @author Claude
// @description 获取系统实时监控数据，包括机器状态统计
// @reason 原实现返回Mock数据，现改为从MachineService获取真实数据
// @modified 2026-02-04
// TODO: 接入Redis缓存，设置采样频率（如30秒）避免频繁查询数据库
func (s *MonitorService) GetGlobalSnapshot(ctx context.Context) (map[string]interface{}, error) {
	// 获取机器状态统计
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

	return map[string]interface{}{
		"total_machines":     totalMachines,
		"online_machines":    onlineMachines,
		"idle_machines":      stats["idle"],
		"allocated_machines": stats["allocated"],
		"offline_machines":   stats["offline"],
		"avg_gpu_util":       0.0, // TODO: 从监控系统获取GPU利用率
	}, nil
}

func (s *MonitorService) GetGPUTrend(ctx context.Context) ([]map[string]interface{}, error) {
	// Mock trend data
	return []map[string]interface{}{
		{"time": "10:00", "value": 45},
		{"time": "11:00", "value": 60},
		{"time": "12:00", "value": 80},
	}, nil
}

type DashboardService struct {
	machineService *machine.MachineService
	customerService *customer.CustomerService
	allocationService *allocation.AllocationService
}

func NewDashboardService(ms *machine.MachineService, cs *customer.CustomerService, as *allocation.AllocationService) *DashboardService {
	return &DashboardService{
		machineService: ms,
		customerService: cs,
		allocationService: as,
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
// @description 返回最近时间段的 GPU 利用率趋势
// @modified 2026-02-04
// TODO: 接入真实监控数据源 (Prometheus/InfluxDB)
func (s *DashboardService) GetGPUTrend(ctx context.Context) ([]map[string]any, error) {
	// TODO: 从监控系统获取真实数据
	return []map[string]any{
		{"time": "00:00", "usage": 45},
		{"time": "04:00", "usage": 30},
		{"time": "08:00", "usage": 65},
		{"time": "12:00", "usage": 80},
		{"time": "16:00", "usage": 75},
		{"time": "20:00", "usage": 60},
	}, nil
}



// GetRecentAllocations 获取最近的分配记录
// @description 用于仪表盘展示最近的机器分配情况
// @modified 2026-02-04
func (s *DashboardService) GetRecentAllocations(ctx context.Context) ([]entity.Allocation, error) {
	return s.allocationService.GetRecent(ctx)
}
