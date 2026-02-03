package ops

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/service/allocation"
	"github.com/YoungBoyGod/remotegpu/internal/service/customer"
	"github.com/YoungBoyGod/remotegpu/internal/service/machine"
)

type MonitorService struct {
	// redis client, influxdb client etc.
}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

func (s *MonitorService) GetGlobalSnapshot(ctx context.Context) (map[string]interface{}, error) {
	// Mock data for now
	return map[string]interface{}{
		"total_machines": 100,
		"online_machines": 95,
		"avg_gpu_util": 75.5,
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

func (s *DashboardService) GetAggregatedStats(ctx context.Context) (map[string]interface{}, error) {
	// In reality, this would call count methods on DAOs
	return map[string]interface{}{
		"total_machines": 100, // s.machineService.Count()
		"active_customers": 20, // s.customerService.Count()
	}, nil
}