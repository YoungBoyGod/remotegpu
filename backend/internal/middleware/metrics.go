package middleware

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 请求计数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP 请求总数",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求延迟直方图
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP 请求延迟（秒）",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 业务指标：机器分配计数
	MachineAllocationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "machine_allocations_total",
			Help: "机器分配总数",
		},
	)

	// 业务指标：机器回收计数
	MachineReclamationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "machine_reclamations_total",
			Help: "机器回收总数",
		},
	)

	// 业务指标：任务创建计数
	TasksCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_created_total",
			Help: "任务创建总数",
		},
	)

	// 业务指标：任务完成计数
	TasksCompletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_completed_total",
			Help: "任务完成总数",
		},
	)

	// 业务指标：当前在线机器数
	MachinesOnline = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "machines_online",
			Help: "当前在线机器数",
		},
	)
)

// PrometheusMetrics Gin 中间件，记录 HTTP 请求指标
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

// dbStatsCollector 数据库连接池指标采集器
type dbStatsCollector struct {
	db *sql.DB

	openConns  *prometheus.Desc
	idleConns  *prometheus.Desc
	inUseConns *prometheus.Desc
	maxOpen    *prometheus.Desc
}

// RegisterDBMetrics 注册数据库连接池指标
func RegisterDBMetrics(db *sql.DB) {
	collector := &dbStatsCollector{
		db: db,
		openConns: prometheus.NewDesc(
			"db_open_connections",
			"当前打开的数据库连接数",
			nil, nil,
		),
		idleConns: prometheus.NewDesc(
			"db_idle_connections",
			"当前空闲的数据库连接数",
			nil, nil,
		),
		inUseConns: prometheus.NewDesc(
			"db_in_use_connections",
			"当前使用中的数据库连接数",
			nil, nil,
		),
		maxOpen: prometheus.NewDesc(
			"db_max_open_connections",
			"最大打开连接数",
			nil, nil,
		),
	}
	prometheus.MustRegister(collector)
}

func (c *dbStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.openConns
	ch <- c.idleConns
	ch <- c.inUseConns
	ch <- c.maxOpen
}

func (c *dbStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()
	ch <- prometheus.MustNewConstMetric(c.openConns, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.inUseConns, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.maxOpen, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
}
