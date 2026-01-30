package service

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"gorm.io/gorm"
)

// HostDaoInterface 主机DAO接口
type HostDaoInterface interface {
	Create(host *entity.Host) error
	GetByID(id string) (*entity.Host, error)
	Update(host *entity.Host) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.Host, int64, error)
	ListByStatus(status string) ([]*entity.Host, error)
	UpdateStatus(id, status string) error
	UpdateHeartbeat(id string) error
}

// GPUDaoInterface GPU DAO接口
type GPUDaoInterface interface {
	Create(gpu *entity.GPU) error
	GetByID(id uint) (*entity.GPU, error)
	GetByHostID(hostID string) ([]*entity.GPU, error)
	Update(gpu *entity.GPU) error
	Delete(id uint) error
	DeleteByHostID(hostID string) error
	UpdateStatus(id uint, status string) error
	List(page, pageSize int) ([]*entity.GPU, int64, error)
	GetByStatus(status string) ([]*entity.GPU, error)
	Allocate(id uint, allocatedTo string) error
	Release(id uint) error
}

// EnvironmentDaoInterface 环境DAO接口
type EnvironmentDaoInterface interface {
	Create(env *entity.Environment) error
	Update(env *entity.Environment) error
	GetByID(id string) (*entity.Environment, error)
	GetByUserID(userID uint) ([]*entity.Environment, error)
	GetByWorkspaceID(workspaceID uint) ([]*entity.Environment, error)
}

// PortMappingDaoInterface 端口映射DAO接口
type PortMappingDaoInterface interface {
	GetByEnvironmentID(envID string) ([]*entity.PortMapping, error)
	DeleteByEnvironmentID(envID string) error
}

// ResourceQuotaServiceInterface 配额服务接口
type ResourceQuotaServiceInterface interface {
	CheckQuota(customerID uint, workspaceID *uint, req *ResourceRequest) (bool, error)
}

// K8sClientInterface K8s客户端接口
type K8sClientInterface interface {
	CreatePod(config *k8s.PodConfig) (*corev1.Pod, error)
	DeletePod(namespace, name string) error
	GetPodStatus(namespace, name string) (string, error)
	GetPodLogs(namespace, name string, opts *k8s.LogOptions) (string, error)
}

// ResourceQuotaServiceInterface 配额服务接口（已定义，这里重复是为了完整性）
// 实际使用时应该只定义一次

// DBInterface 数据库接口
type DBInterface interface {
	Transaction(fc func(tx *gorm.DB) error) error
}

// GormDBWrapper 包装 gorm.DB 实现 DBInterface
type GormDBWrapper struct {
	db *gorm.DB
}

// NewGormDBWrapper 创建 GormDBWrapper
func NewGormDBWrapper(db *gorm.DB) *GormDBWrapper {
	return &GormDBWrapper{db: db}
}

// Transaction 实现 DBInterface
func (w *GormDBWrapper) Transaction(fc func(tx *gorm.DB) error) error {
	return w.db.Transaction(fc)
}
