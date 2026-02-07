package collector

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// collectSystem 采集系统指标（CPU、内存、磁盘）
func collectSystem(m *Metrics) {
	// CPU 使用率
	if percents, err := cpu.Percent(0, false); err == nil && len(percents) > 0 {
		v := percents[0]
		m.CPUUsagePercent = &v
	}

	// 内存使用率
	if vm, err := mem.VirtualMemory(); err == nil {
		pct := vm.UsedPercent
		m.MemoryUsagePercent = &pct
		usedGB := int64(vm.Used / (1024 * 1024 * 1024))
		m.MemoryUsedGB = &usedGB
	}

	// 磁盘使用率（根分区）
	if usage, err := disk.Usage("/"); err == nil {
		pct := usage.UsedPercent
		m.DiskUsagePercent = &pct
		usedGB := int64(usage.Used / (1024 * 1024 * 1024))
		m.DiskUsedGB = &usedGB
	}
}
