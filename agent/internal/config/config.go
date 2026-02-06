package config

import (
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config Agent 配置
type Config struct {
	Port       int    `yaml:"port"`
	DBPath     string `yaml:"db_path"`
	MaxWorkers int    `yaml:"max_workers"`

	Server   ServerConfig   `yaml:"server"`
	Poll     PollConfig     `yaml:"poll"`
	Limits   LimitsConfig   `yaml:"limits"`
	Security SecurityConfig `yaml:"security"`
}

// ServerConfig Server 连接配置
type ServerConfig struct {
	URL       string        `yaml:"url"`
	AgentID   string        `yaml:"agent_id"`
	MachineID string        `yaml:"machine_id"`
	Token     string        `yaml:"token"`
	Timeout   time.Duration `yaml:"timeout"`
}

// PollConfig 轮询配置
type PollConfig struct {
	Interval  time.Duration `yaml:"interval"`
	BatchSize int           `yaml:"batch_size"`
}

// LimitsConfig 限制配置
type LimitsConfig struct {
	MaxOutputSize int `yaml:"max_output_size"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	AllowedCommands []string `yaml:"allowed_commands"`
	BlockedPatterns []string `yaml:"blocked_patterns"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Port:       8090,
		DBPath:     "/var/lib/remotegpu-agent/tasks.db",
		MaxWorkers: 4,
		Server: ServerConfig{
			Timeout: 30 * time.Second,
		},
		Poll: PollConfig{
			Interval:  5 * time.Second,
			BatchSize: 10,
		},
		Limits: LimitsConfig{
			MaxOutputSize: 1 << 20, // 1MB
		},
	}
}

// Load 加载配置：YAML 文件 → 环境变量覆盖
func Load() *Config {
	cfg := DefaultConfig()

	// 尝试加载 YAML 配置文件
	for _, path := range []string{"./agent.yaml", "/etc/remotegpu-agent/agent.yaml"} {
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
	if v := os.Getenv("AGENT_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	if v := os.Getenv("AGENT_DB_PATH"); v != "" {
		cfg.DBPath = v
	}
	if v := os.Getenv("AGENT_MAX_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MaxWorkers = n
		}
	}
	if v := os.Getenv("SERVER_URL"); v != "" {
		cfg.Server.URL = v
	}
	if v := os.Getenv("AGENT_ID"); v != "" {
		cfg.Server.AgentID = v
	}
	if v := os.Getenv("MACHINE_ID"); v != "" {
		cfg.Server.MachineID = v
	}
	if v := os.Getenv("AGENT_TOKEN"); v != "" {
		cfg.Server.Token = v
	}
}

// ServerConfigured 检查 Server 配置是否完整
func (c *Config) ServerConfigured() bool {
	return c.Server.URL != "" && c.Server.AgentID != "" && c.Server.MachineID != ""
}
