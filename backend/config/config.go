package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	Redis         RedisConfig         `yaml:"redis"`
	JWT           JWTConfig           `yaml:"jwt"`
	Encryption    EncryptionConfig    `yaml:"encryption"`
	Log           LogConfig           `yaml:"log"`
	Storage       StorageConfig       `yaml:"storage"`
	Mail          MailConfig          `yaml:"mail"`
	Harbor        HarborConfig        `yaml:"harbor"`
	Etcd          EtcdConfig          `yaml:"etcd"`
	Prometheus    PrometheusConfig    `yaml:"prometheus"`
	Jumpserver    JumpserverConfig    `yaml:"jumpserver"`
	Nginx         NginxConfig         `yaml:"nginx"`
	UptimeKuma    UptimeKumaConfig    `yaml:"uptime_kuma"`
	Guacamole     GuacamoleConfig     `yaml:"guacamole"`
	K8s           K8sConfig           `yaml:"k8s"`
	HotReload     HotReloadConfig     `yaml:"hot_reload"`
	Graceful      GracefulConfig      `yaml:"graceful"`
	Swagger       SwaggerConfig       `yaml:"swagger"`
	Agent         AgentConfig         `yaml:"agent"`
	Enrollment    EnrollmentConfig    `yaml:"machine_enrollment"`
	MachineAction MachineActionConfig `yaml:"machine_action"`
	HeartbeatMonitor HeartbeatMonitorConfig `yaml:"heartbeat_monitor"`
	MetricsCollector MetricsCollectorConfig `yaml:"metrics_collector"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"` // 连接池大小
	Timeout  int    `yaml:"timeout"`   // 超时时间(秒)
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expire_time"` // 小时
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Key string `yaml:"key"` // AES-256 加密密钥(32字节)
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	OutputPath string `yaml:"output_path"` // 日志输出路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件数量
	MaxAge     int    `yaml:"max_age"`     // 保留旧日志文件的最大天数
}

// StorageConfig 存储配置
type StorageConfig struct {
	MaxUploadSize int64            `yaml:"max_upload_size"` // 最大上传大小(字节)
	Default       string           `yaml:"default"`         // 默认存储后端名称
	Backends      []StorageBackend `yaml:"backends"`        // 存储后端列表
}

// StorageBackend 存储后端配置
type StorageBackend struct {
	Name      string `yaml:"name"`       // 存储后端名称（唯一标识）
	Type      string `yaml:"type"`       // 类型: local, rustfs, s3
	Enabled   bool   `yaml:"enabled"`    // 是否启用
	Path      string `yaml:"path"`       // 本地存储路径（type=local时使用）
	Endpoint  string `yaml:"endpoint"`   // S3/RustFS 端点
	Region    string `yaml:"region"`     // S3 区域
	AccessKey string `yaml:"access_key"` // 访问密钥
	SecretKey string `yaml:"secret_key"` // 秘密密钥
	Bucket    string `yaml:"bucket"`     // 存储桶名称
}

// MailConfig 邮件配置
type MailConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`    // 发件人名称
	UseSSL   bool   `yaml:"use_ssl"` // 是否使用 SSL/TLS
}

// HarborConfig Harbor 镜像仓库配置
type HarborConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Project  string `yaml:"project"` // 默认项目名
}

// EtcdConfig Etcd 配置
type EtcdConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Endpoints []string `yaml:"endpoints"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	Timeout   int      `yaml:"timeout"` // 连接超时时间(秒)
}

// PrometheusConfig Prometheus 监控配置
type PrometheusConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

// JumpserverConfig Jumpserver 堡垒机配置
type JumpserverConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key"`
	OrgID    string `yaml:"org_id"`
}

// NginxConfig Nginx 配置
type NginxConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

// UptimeKumaConfig Uptime Kuma 监控配置
type UptimeKumaConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

// GuacamoleConfig Guacamole 远程桌面网关配置
type GuacamoleConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// K8sConfig Kubernetes 配置
type K8sConfig struct {
	Enabled    bool   `yaml:"enabled"`
	KubeConfig string `yaml:"kubeconfig"` // kubeconfig文件路径
	Namespace  string `yaml:"namespace"`  // 默认命名空间
	InCluster  bool   `yaml:"in_cluster"` // 是否在集群内运行
}

// HotReloadConfig 热更新配置
type HotReloadConfig struct {
	Enabled     bool     `yaml:"enabled"`      // 是否启用热更新
	WatchDirs   []string `yaml:"watch_dirs"`   // 监控的目录
	WatchExts   []string `yaml:"watch_exts"`   // 监控的文件扩展名
	ExcludeDirs []string `yaml:"exclude_dirs"` // 排除的目录
	BuildCmd    string   `yaml:"build_cmd"`    // 构建命令
	Debounce    int      `yaml:"debounce"`     // 防抖时间(秒)
}

// GracefulConfig 优雅启动配置
type GracefulConfig struct {
	ShutdownTimeout int `yaml:"shutdown_timeout"` // 关闭超时时间(秒)
	RetryInterval   int `yaml:"retry_interval"`   // 重试间隔(秒)
	MaxRetries      int `yaml:"max_retries"`      // 最大重试次数，0表示无限重试
}

// SwaggerConfig Swagger API文档配置
type SwaggerConfig struct {
	Enabled     bool     `yaml:"enabled"`      // 是否启用Swagger
	Title       string   `yaml:"title"`        // API标题
	Description string   `yaml:"description"`  // API描述
	Version     string   `yaml:"version"`      // API版本
	BasePath    string   `yaml:"base_path"`    // API基础路径
	OpenAPIPath string   `yaml:"openapi_path"` // OpenAPI 文档路径
	Schemes     []string `yaml:"schemes"`      // 支持的协议 (http, https)
}

// AgentConfig Agent 通信配置
type AgentConfig struct {
	Enabled     bool   `yaml:"enabled"`     // 是否启用 Agent 通信
	Protocol    string `yaml:"protocol"`    // 通信协议: grpc, http
	Port        int    `yaml:"port"`        // Agent 默认端口
	GRPCPort    int    `yaml:"grpc_port"`   // gRPC 端口
	HTTPPort    int    `yaml:"http_port"`   // HTTP 端口
	Timeout     int    `yaml:"timeout"`     // 请求超时时间(秒)
	RetryCount  int    `yaml:"retry_count"` // 重试次数
	RetryDelay  int    `yaml:"retry_delay"` // 重试间隔(秒)
	TLSEnabled  bool   `yaml:"tls_enabled"` // 是否启用 TLS
	TLSCertFile string `yaml:"tls_cert"`    // TLS 证书文件
	TLSKeyFile  string `yaml:"tls_key"`     // TLS 密钥文件
}

// EnrollmentConfig 用户添加机器队列配置
type EnrollmentConfig struct {
	MaxRetries  int  `yaml:"max_retries"`  // 最大重试次数
	RetryDelay  int  `yaml:"retry_delay"`  // 重试延迟(秒)
	SkipCollect bool `yaml:"skip_collect"` // 跳过采集直接入库
}

// MachineActionConfig 机器动作队列配置
type MachineActionConfig struct {
	MaxRetries int `yaml:"max_retries"` // 最大重试次数
	RetryDelay int `yaml:"retry_delay"` // 重试延迟(秒)
}

// HeartbeatMonitorConfig 心跳监控配置
type HeartbeatMonitorConfig struct {
	Enabled       bool `yaml:"enabled"`        // 是否启用
	Timeout       int  `yaml:"timeout"`        // 心跳超时时间(秒)
	CheckInterval int  `yaml:"check_interval"` // 检查间隔(秒)
}

// MetricsCollectorConfig 监控数据采集配置
type MetricsCollectorConfig struct {
	Enabled       bool `yaml:"enabled"`        // 是否启用
	Interval      int  `yaml:"interval"`       // 采集间隔(秒)
	RetentionDays int  `yaml:"retention_days"` // 数据保留天数
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	GlobalConfig = &cfg
	return nil
}
