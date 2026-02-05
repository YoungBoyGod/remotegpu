package prometheus

import (
	"context"
	"strconv"
	"time"
)

// GPUMetrics GPU 指标数据
type GPUMetrics struct {
	AvgUtilization float64            `json:"avg_utilization"`
	ByHost         map[string]float64 `json:"by_host,omitempty"`
}

// GPUTrendPoint GPU 趋势数据点
type GPUTrendPoint struct {
	Time  string  `json:"time"`
	Usage float64 `json:"usage"`
}

// GetGPUUtilization 获取 GPU 平均利用率
func (c *Client) GetGPUUtilization(ctx context.Context) (*GPUMetrics, error) {
	if c == nil {
		return &GPUMetrics{AvgUtilization: 0}, nil
	}

	// 查询所有 GPU 的平均利用率
	// 常见的 nvidia_gpu_exporter 指标名称
	queries := []string{
		"avg(DCGM_FI_DEV_GPU_UTIL)",           // DCGM exporter
		"avg(nvidia_gpu_duty_cycle)",           // nvidia_gpu_exporter
		"avg(nvidia_smi_utilization_gpu)",      // nvidia_smi_exporter
	}

	for _, query := range queries {
		resp, err := c.Query(ctx, query)
		if err != nil {
			continue
		}
		if len(resp.Data.Result) > 0 {
			if val, ok := parseValue(resp.Data.Result[0].Value); ok {
				return &GPUMetrics{AvgUtilization: val}, nil
			}
		}
	}

	return &GPUMetrics{AvgUtilization: 0}, nil
}

// GetGPUTrend 获取 GPU 利用率趋势
func (c *Client) GetGPUTrend(ctx context.Context, duration time.Duration) ([]GPUTrendPoint, error) {
	if c == nil {
		return defaultGPUTrend(), nil
	}

	end := time.Now()
	start := end.Add(-duration)
	step := duration / 6 // 6 个数据点

	queries := []string{
		"avg(DCGM_FI_DEV_GPU_UTIL)",
		"avg(nvidia_gpu_duty_cycle)",
		"avg(nvidia_smi_utilization_gpu)",
	}

	for _, query := range queries {
		resp, err := c.QueryRange(ctx, query, start, end, step)
		if err != nil {
			continue
		}
		if len(resp.Data.Result) > 0 {
			return parseRangeValues(resp.Data.Result[0].Values), nil
		}
	}

	return defaultGPUTrend(), nil
}

// parseValue 解析即时查询值
func parseValue(value []interface{}) (float64, bool) {
	if len(value) < 2 {
		return 0, false
	}
	if strVal, ok := value[1].(string); ok {
		if val, err := strconv.ParseFloat(strVal, 64); err == nil {
			return val, true
		}
	}
	return 0, false
}

// parseRangeValues 解析范围查询值
func parseRangeValues(values [][]interface{}) []GPUTrendPoint {
	var points []GPUTrendPoint
	for _, v := range values {
		if len(v) < 2 {
			continue
		}
		ts, ok1 := v[0].(float64)
		strVal, ok2 := v[1].(string)
		if !ok1 || !ok2 {
			continue
		}
		val, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			continue
		}
		t := time.Unix(int64(ts), 0)
		points = append(points, GPUTrendPoint{
			Time:  t.Format("15:04"),
			Usage: val,
		})
	}
	return points
}

// defaultGPUTrend 默认趋势数据
func defaultGPUTrend() []GPUTrendPoint {
	return []GPUTrendPoint{
		{Time: "00:00", Usage: 0},
		{Time: "04:00", Usage: 0},
		{Time: "08:00", Usage: 0},
		{Time: "12:00", Usage: 0},
		{Time: "16:00", Usage: 0},
		{Time: "20:00", Usage: 0},
	}
}
