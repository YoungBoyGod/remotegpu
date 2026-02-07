package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Harbor API 客户端
type Client struct {
	endpoint   string
	username   string
	password   string
	httpClient *http.Client
}

// NewClient 创建 Harbor 客户端
func NewClient(endpoint, username, password string) *Client {
	return &Client{
		endpoint: endpoint,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Repository Harbor 仓库信息
type Repository struct {
	Name          string    `json:"name"`
	ArtifactCount int       `json:"artifact_count"`
	UpdateTime    time.Time `json:"update_time"`
	Description   string    `json:"description"`
}

// Artifact Harbor 镜像制品信息
type Artifact struct {
	Digest    string    `json:"digest"`
	Size      int64     `json:"size"`
	PushTime  time.Time `json:"push_time"`
	Tags      []Tag     `json:"tags"`
	ExtraAttrs map[string]interface{} `json:"extra_attrs"`
}

// Tag 镜像标签
type Tag struct {
	Name     string    `json:"name"`
	PushTime time.Time `json:"push_time"`
}

// doRequest 发送带认证的 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v2.0%s", c.endpoint, path)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Harbor 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Harbor API 返回错误 %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// ListRepositories 列出项目下的所有仓库
func (c *Client) ListRepositories(ctx context.Context, project string) ([]Repository, error) {
	path := fmt.Sprintf("/projects/%s/repositories?page_size=100", project)
	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var repos []Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("解析仓库列表失败: %w", err)
	}
	return repos, nil
}

// ListArtifacts 列出仓库下的所有制品（镜像标签）
func (c *Client) ListArtifacts(ctx context.Context, project, repoName string) ([]Artifact, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts?page_size=100&with_tag=true", project, repoName)
	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var artifacts []Artifact
	if err := json.Unmarshal(body, &artifacts); err != nil {
		return nil, fmt.Errorf("解析制品列表失败: %w", err)
	}
	return artifacts, nil
}

// Ping 检查 Harbor 连接是否正常
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.doRequest(ctx, http.MethodGet, "/systeminfo")
	return err
}
