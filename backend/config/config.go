package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	JWT        JWTConfig        `yaml:"jwt"`
	Log        LogConfig        `yaml:"log"`
	Storage    StorageConfig    `yaml:"storage"`
	Harbor     HarborConfig     `yaml:"harbor"`
	Etcd       EtcdConfig       `yaml:"etcd"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Jumpserver JumpserverConfig `yaml:"jumpserver"`
	Nginx      NginxConfig      `yaml:"nginx"`
	UptimeKuma UptimeKumaConfig `yaml:"uptime_kuma"`
	Guacamole  GuacamoleConfig  `yaml:"guacamole"`
	K8s        K8sConfig        `yaml:"k8s"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
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
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expire_time"` // 小时
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	OutputPath string `yaml:"output_path"` // 日志输出路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件数量
	MaxAge     int    `yaml:"max_age"`     // 保留旧日志文件的最大天数
}

// StorageConfig 存储配置（RustFS）
type StorageConfig struct {
	Type      string `yaml:"type"`       // local, rustfs, s3
	LocalPath string `yaml:"local_path"` // 本地存储路径
	RustFS    struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucket"`
	} `yaml:"rustfs"`
	MaxUploadSize int64 `yaml:"max_upload_size"` // 最大上传大小(字节)
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
