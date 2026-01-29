package health

import (
	"github.com/YoungBoyGod/remotegpu/config"
)

// InitManager 初始化健康检查管理器
func InitManager(cfg *config.Config) *Manager {
	manager := NewManager()

	// 注册数据库健康检查
	manager.Register(NewPostgreSQLChecker(cfg.Database))
	manager.Register(NewRedisChecker(cfg.Redis))

	// 注册基础设施健康检查
	manager.Register(NewEtcdChecker(cfg.Etcd))
	manager.Register(NewK8sChecker(cfg.K8s))

	// 注册HTTP服务健康检查
	manager.Register(NewHarborChecker(cfg.Harbor))
	manager.Register(NewPrometheusChecker(cfg.Prometheus))
	manager.Register(NewJumpserverChecker(cfg.Jumpserver))
	manager.Register(NewNginxChecker(cfg.Nginx))
	manager.Register(NewUptimeKumaChecker(cfg.UptimeKuma))
	manager.Register(NewGuacamoleChecker(cfg.Guacamole))
	manager.Register(NewS3Checker(cfg.Storage))

	return manager
}
