package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client Prometheus 客户端
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// Config Prometheus 配置
type Config struct {
	Enabled  bool
	Endpoint string
}

// NewClient 创建 Prometheus 客户端
func NewClient(cfg *Config) *Client {
	if !cfg.Enabled || cfg.Endpoint == "" {
		return nil
	}
	return &Client{
		endpoint: cfg.Endpoint,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// QueryResponse Prometheus 查询响应
type QueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string        `json:"resultType"`
		Result     []QueryResult `json:"result"`
	} `json:"data"`
	Error     string `json:"error,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
}

// QueryResult 查询结果
type QueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}   `json:"value"`  // [timestamp, value]
	Values [][]interface{}   `json:"values"` // for range queries
}

// Query 执行即时查询
func (c *Client) Query(ctx context.Context, query string) (*QueryResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("prometheus client not initialized")
	}

	u, err := url.Parse(c.endpoint + "/api/v1/query")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s", result.Error)
	}

	return &result, nil
}

// QueryRange 执行范围查询
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*QueryResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("prometheus client not initialized")
	}

	u, err := url.Parse(c.endpoint + "/api/v1/query_range")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.Unix()))
	q.Set("end", fmt.Sprintf("%d", end.Unix()))
	q.Set("step", step.String())
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s", result.Error)
	}

	return &result, nil
}
