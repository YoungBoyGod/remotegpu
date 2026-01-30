package dao

import (
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EnvironmentDao 环境数据访问对象
type EnvironmentDao struct {
	db *gorm.DB
}

// NewEnvironmentDao 创建环境 DAO
func NewEnvironmentDao() *EnvironmentDao {
	return &EnvironmentDao{
		db: database.GetDB(),
	}
}

// Create 创建环境
func (d *EnvironmentDao) Create(env *entity.Environment) error {
	return d.db.Create(env).Error
}

// GetByID 根据ID获取环境
func (d *EnvironmentDao) GetByID(id string) (*entity.Environment, error) {
	var env entity.Environment
	err := d.db.Preload("User").Preload("Workspace").Preload("Host").Preload("PortMappings").
		Where("id = ?", id).First(&env).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// Update 更新环境
func (d *EnvironmentDao) Update(env *entity.Environment) error {
	return d.db.Save(env).Error
}

// Delete 删除环境
func (d *EnvironmentDao) Delete(id string) error {
	return d.db.Delete(&entity.Environment{}, "id = ?", id).Error
}

// GetByUserID 根据用户ID获取环境列表
func (d *EnvironmentDao) GetByUserID(userID uint) ([]*entity.Environment, error) {
	var envs []*entity.Environment
	err := d.db.Where("user_id = ?", userID).Find(&envs).Error
	if err != nil {
		return nil, err
	}
	return envs, nil
}

// GetByWorkspaceID 根据工作空间ID获取环境列表
func (d *EnvironmentDao) GetByWorkspaceID(workspaceID uint) ([]*entity.Environment, error) {
	var envs []*entity.Environment
	err := d.db.Where("workspace_id = ?", workspaceID).Find(&envs).Error
	if err != nil {
		return nil, err
	}
	return envs, nil
}

// GetByHostID 根据主机ID获取环境列表
func (d *EnvironmentDao) GetByHostID(hostID string) ([]*entity.Environment, error) {
	var envs []*entity.Environment
	err := d.db.Where("host_id = ?", hostID).Find(&envs).Error
	if err != nil {
		return nil, err
	}
	return envs, nil
}

// GetByStatus 根据状态获取环境列表
func (d *EnvironmentDao) GetByStatus(status string) ([]*entity.Environment, error) {
	var envs []*entity.Environment
	err := d.db.Where("status = ?", status).Find(&envs).Error
	if err != nil {
		return nil, err
	}
	return envs, nil
}

// List 获取环境列表（分页）
func (d *EnvironmentDao) List(page, pageSize int) ([]*entity.Environment, int64, error) {
	var envs []*entity.Environment
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&entity.Environment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Offset(offset).Limit(pageSize).Find(&envs).Error; err != nil {
		return nil, 0, err
	}

	return envs, total, nil
}

// PortMappingDao 端口映射数据访问对象
type PortMappingDao struct {
	db *gorm.DB
}

// NewPortMappingDao 创建端口映射 DAO
func NewPortMappingDao() *PortMappingDao {
	return &PortMappingDao{
		db: database.GetDB(),
	}
}

// Create 创建端口映射
func (d *PortMappingDao) Create(pm *entity.PortMapping) error {
	return d.db.Create(pm).Error
}

// GetByID 根据ID获取端口映射
func (d *PortMappingDao) GetByID(id int64) (*entity.PortMapping, error) {
	var pm entity.PortMapping
	err := d.db.Where("id = ?", id).First(&pm).Error
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

// GetByEnvironmentID 根据环境ID获取端口映射列表
func (d *PortMappingDao) GetByEnvironmentID(envID string) ([]*entity.PortMapping, error) {
	var pms []*entity.PortMapping
	err := d.db.Where("env_id = ?", envID).Find(&pms).Error
	if err != nil {
		return nil, err
	}
	return pms, nil
}

// Delete 删除端口映射
func (d *PortMappingDao) Delete(id int64) error {
	return d.db.Delete(&entity.PortMapping{}, id).Error
}

// DeleteByEnvironmentID 根据环境ID删除所有端口映射
func (d *PortMappingDao) DeleteByEnvironmentID(envID string) error {
	return d.db.Where("env_id = ?", envID).Delete(&entity.PortMapping{}).Error
}

// AllocatePort 分配端口（并发安全）
// envID: 环境ID
// serviceType: 服务类型（ssh/rdp/jupyter/custom）
// internalPort: 容器内部端口
// 返回: 分配的端口映射记录
func (d *PortMappingDao) AllocatePort(envID, serviceType string, internalPort int) (*entity.PortMapping, error) {
	const (
		minPort    = 30000
		maxPort    = 32767
		maxRetries = 10 // 最大重试次数
	)

	var pm *entity.PortMapping
	var lastErr error

	// 重试机制处理并发冲突
	for retry := 0; retry < maxRetries; retry++ {
		err := d.db.Transaction(func(tx *gorm.DB) error {
			// 查询已使用的端口（使用行锁）
			var usedPorts []int
			if err := tx.Model(&entity.PortMapping{}).
				Where("status = ?", "active").
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Pluck("external_port", &usedPorts).Error; err != nil {
				return err
			}

			// 创建已使用端口的map，便于快速查找
			usedPortMap := make(map[int]bool)
			for _, port := range usedPorts {
				usedPortMap[port] = true
			}

			// 查找第一个可用端口
			availablePort := 0
			for port := minPort; port <= maxPort; port++ {
				if !usedPortMap[port] {
					availablePort = port
					break
				}
			}

			if availablePort == 0 {
				return gorm.ErrRecordNotFound // 没有可用端口
			}

			// 创建端口映射记录
			pm = &entity.PortMapping{
				EnvID:        envID,
				ServiceType:  serviceType,
				ExternalPort: availablePort,
				InternalPort: internalPort,
				Status:       "active",
			}

			return tx.Create(pm).Error
		})

		if err == nil {
			// 成功分配端口
			return pm, nil
		}

		lastErr = err

		// 如果是唯一约束冲突，重试
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"idx_port_mappings_external_port\" (SQLSTATE 23505)" ||
			err.Error() == "UNIQUE constraint failed: port_mappings.external_port" {
			continue
		}

		// 其他错误直接返回
		return nil, err
	}

	// 重试次数用尽
	if lastErr != nil {
		return nil, lastErr
	}

	return nil, gorm.ErrRecordNotFound
}

// ReleasePort 释放端口
func (d *PortMappingDao) ReleasePort(id int64) error {
	return d.db.Model(&entity.PortMapping{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "released",
			"released_at":  gorm.Expr("NOW()"),
		}).Error
}

// GetAvailablePort 获取可用端口数量
func (d *PortMappingDao) GetAvailablePort() (int, error) {
	const (
		minPort   = 30000
		maxPort   = 32767
		totalPort = maxPort - minPort + 1
	)

	var usedCount int64
	err := d.db.Model(&entity.PortMapping{}).
		Where("status = ?", "active").
		Count(&usedCount).Error
	if err != nil {
		return 0, err
	}

	return totalPort - int(usedCount), nil
}
