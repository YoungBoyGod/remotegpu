package v1

// CreateAlertRuleRequest 创建告警规则请求
type CreateAlertRuleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	MetricType  string  `json:"metric_type" binding:"required"`
	Threshold   float64 `json:"threshold" binding:"required"`
	Condition   string  `json:"condition" binding:"required"`
	Duration    int     `json:"duration"`
	Severity    string  `json:"severity"`
	Enabled     *bool   `json:"enabled"`
}

// UpdateAlertRuleRequest 更新告警规则请求
type UpdateAlertRuleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MetricType  string   `json:"metric_type"`
	Threshold   *float64 `json:"threshold"`
	Condition   string   `json:"condition"`
	Duration    *int     `json:"duration"`
	Severity    string   `json:"severity"`
	Enabled     *bool    `json:"enabled"`
}
