package config

import (
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config Proxy 配置
type Config struct {
	Port int `yaml:"port"`

	Server    ServerConfig    `yaml:"server"`
	PortPool  PortPoolConfig  `yaml:"port_pool"`
	Heartbeat HeartbeatConfig `yaml:"heartbeat"`
	HTTPProxy HTTPProxyConfig `yaml:"http_proxy"`
	Network   NetworkConfig   `yaml:"network"`
}

// ServerConfig 后端服务器连接配置
type ServerConfig struct {
	URL     string        `yaml:"url"`
	ProxyID string        `yaml:"proxy_id"`
	Token   string        `yaml:"token"`
	Timeout time.Duration `yaml:"timeout"`
}

// PortPoolConfig 端口池配置
type PortPoolConfig struct {
	RangeStart int `yaml:"range_start"`
	RangeEnd   int `yaml:"range_end"`
}

// HeartbeatConfig 心跳配置
type HeartbeatConfig struct {
	Interval time.Duration `yaml:"interval"`
}

// HTTPProxyConfig HTTP 反向代理配置
type HTTPProxyConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	InnerIP string `yaml:"inner_ip"`
	OuterIP string `yaml:"outer_ip"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Port: 9090,
		Server: ServerConfig{
			Timeout: 30 * time.Second,
		},
		PortPool: PortPoolConfig{
			RangeStart: 20000,
			RangeEnd:   60000,
		},
		Heartbeat: HeartbeatConfig{
			Interval: 30 * time.Second,
		},
		HTTPProxy: HTTPProxyConfig{
			Enabled: false,
			Port:    9091,
		},
	}
}

// Load 加载配置：YAML 文件 → 环境变量覆盖
func Load() *Config {
	cfg := DefaultConfig()

	// 尝试加载 YAML 配置文件
	for _, path := range []string{"./proxy.yaml", "/etc/remotegpu-proxy/proxy.yaml"} {
		if data, err := os.ReadFile(path); err == nil {
			yaml.Unmarshal(data, cfg)
			break
		}
	}

	// 环境变量覆盖
	applyEnv(cfg)

	return cfg
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("PROXY_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	if v := os.Getenv("SERVER_URL"); v != "" {
		cfg.Server.URL = v
	}
	if v := os.Getenv("PROXY_ID"); v != "" {
		cfg.Server.ProxyID = v
	}
	if v := os.Getenv("PROXY_TOKEN"); v != "" {
		cfg.Server.Token = v
	}
	if v := os.Getenv("PORT_RANGE_START"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.PortPool.RangeStart = n
		}
	}
	if v := os.Getenv("PORT_RANGE_END"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.PortPool.RangeEnd = n
		}
	}
	if v := os.Getenv("INNER_IP"); v != "" {
		cfg.Network.InnerIP = v
	}
	if v := os.Getenv("OUTER_IP"); v != "" {
		cfg.Network.OuterIP = v
	}
}

// ServerConfigured 检查 Server 配置是否完整
func (c *Config) ServerConfigured() bool {
	return c.Server.URL != "" && c.Server.ProxyID != ""
}
