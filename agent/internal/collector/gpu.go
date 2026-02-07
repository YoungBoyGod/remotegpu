package collector

import (
	"os/exec"
	"strconv"
	"strings"
)

// collectGPU 通过 nvidia-smi 采集 GPU 指标
func collectGPU(m *Metrics) {
	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=index,uuid,name,utilization.gpu,memory.used,memory.total,temperature.gpu,power.draw",
		"--format=csv,noheader,nounits",
	).Output()
	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		gm := parseGPULine(line)
		if gm != nil {
			m.GPUMetrics = append(m.GPUMetrics, *gm)
		}
	}
}

// parseGPULine 解析 nvidia-smi CSV 输出的一行
// 格式: index, uuid, name, utilization.gpu, memory.used, memory.total, temperature.gpu, power.draw
func parseGPULine(line string) *GPUMetric {
	fields := strings.Split(line, ", ")
	if len(fields) < 8 {
		return nil
	}

	idx, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return nil
	}

	gm := &GPUMetric{
		Index: idx,
		UUID:  strings.TrimSpace(fields[1]),
		Name:  strings.TrimSpace(fields[2]),
	}

	if v, err := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64); err == nil {
		gm.UtilPercent = &v
	}
	if v, err := strconv.Atoi(strings.TrimSpace(fields[4])); err == nil {
		gm.MemoryUsedMB = &v
	}
	if v, err := strconv.Atoi(strings.TrimSpace(fields[5])); err == nil {
		gm.MemoryTotalMB = &v
	}
	if v, err := strconv.Atoi(strings.TrimSpace(fields[6])); err == nil {
		gm.TemperatureC = &v
	}
	if v, err := strconv.ParseFloat(strings.TrimSpace(fields[7]), 64); err == nil {
		gm.PowerUsageW = &v
	}

	return gm
}
