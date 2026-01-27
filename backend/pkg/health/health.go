package health

import (
	"context"
	"time"
)

// Status 健康状态
type Status string

const (
	StatusHealthy   Status = "healthy"   // 健康
	StatusUnhealthy Status = "unhealthy" // 不健康
	StatusUnknown   Status = "unknown"   // 未知
	StatusDisabled  Status = "disabled"  // 已禁用
)

// CheckResult 健康检查结果
type CheckResult struct {
	Service   string        `json:"service"`    // 服务名称
	Status    Status        `json:"status"`     // 健康状态
	Message   string        `json:"message"`    // 状态消息
	Latency   time.Duration `json:"latency"`    // 响应延迟
	Timestamp time.Time     `json:"timestamp"`  // 检查时间
	Details   interface{}   `json:"details"`    // 详细信息
}

// Checker 健康检查器接口
type Checker interface {
	// Check 执行健康检查
	Check(ctx context.Context) *CheckResult
	// Name 返回服务名称
	Name() string
}

// Manager 健康检查管理器
type Manager struct {
	checkers []Checker
}

// NewManager 创建健康检查管理器
func NewManager() *Manager {
	return &Manager{
		checkers: make([]Checker, 0),
	}
}

// Register 注册健康检查器
func (m *Manager) Register(checker Checker) {
	m.checkers = append(m.checkers, checker)
}

// CheckAll 检查所有服务
func (m *Manager) CheckAll(ctx context.Context) []*CheckResult {
	results := make([]*CheckResult, 0, len(m.checkers))
	for _, checker := range m.checkers {
		result := checker.Check(ctx)
		results = append(results, result)
	}
	return results
}

// CheckService 检查指定服务
func (m *Manager) CheckService(ctx context.Context, serviceName string) *CheckResult {
	for _, checker := range m.checkers {
		if checker.Name() == serviceName {
			return checker.Check(ctx)
		}
	}
	return &CheckResult{
		Service:   serviceName,
		Status:    StatusUnknown,
		Message:   "服务未注册",
		Timestamp: time.Now(),
	}
}
