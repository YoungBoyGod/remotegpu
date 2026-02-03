package host

import (
	"fmt"
)

// LeastUsedSelector 选择资源使用率最低的主机
type LeastUsedSelector struct{}

// NewLeastUsedSelector 创建最低使用率选择器
func NewLeastUsedSelector() *LeastUsedSelector {
	return &LeastUsedSelector{}
}

// Select 选择主机
func (s *LeastUsedSelector) Select(hosts []*HostInfo, req *ResourceRequirement) (*HostInfo, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("没有可用的主机")
	}

	// 过滤出可以分配资源的主机
	var availableHosts []*HostInfo
	for _, host := range hosts {
		if host.Status == "active" && host.CanAllocate(req) {
			availableHosts = append(availableHosts, host)
		}
	}

	if len(availableHosts) == 0 {
		return nil, fmt.Errorf("没有满足资源要求的主机")
	}

	// 选择使用率最低的主机
	var bestHost *HostInfo
	lowestRate := 1.0

	for _, host := range availableHosts {
		rate := host.UsageRate()
		if bestHost == nil || rate < lowestRate {
			bestHost = host
			lowestRate = rate
		}
	}

	return bestHost, nil
}
