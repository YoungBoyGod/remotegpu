package dao

import "github.com/YoungBoyGod/remotegpu/internal/model/entity"

// ResourceQuotaDaoInterface 资源配额DAO接口
type ResourceQuotaDaoInterface interface {
	Create(quota *entity.ResourceQuota) error
	GetByID(id uint) (*entity.ResourceQuota, error)
	GetByUserID(userID uint) (*entity.ResourceQuota, error)
	Update(quota *entity.ResourceQuota) error
	Delete(id uint) error
}
