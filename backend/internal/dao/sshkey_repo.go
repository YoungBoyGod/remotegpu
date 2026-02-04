package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type SSHKeyDao struct {
	db *gorm.DB
}

func NewSSHKeyDao(db *gorm.DB) *SSHKeyDao {
	return &SSHKeyDao{db: db}
}

func (d *SSHKeyDao) Create(ctx context.Context, key *entity.SSHKey) error {
	return d.db.WithContext(ctx).Create(key).Error
}

func (d *SSHKeyDao) FindByID(ctx context.Context, id uint) (*entity.SSHKey, error) {
	var key entity.SSHKey
	if err := d.db.WithContext(ctx).First(&key, id).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func (d *SSHKeyDao) ListByCustomerID(ctx context.Context, customerID uint) ([]entity.SSHKey, error) {
	var keys []entity.SSHKey
	if err := d.db.WithContext(ctx).Where("customer_id = ?", customerID).Order("created_at DESC").Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (d *SSHKeyDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&entity.SSHKey{}, id).Error
}

func (d *SSHKeyDao) ExistsByFingerprint(ctx context.Context, customerID uint, fingerprint string) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&entity.SSHKey{}).
		Where("customer_id = ? AND fingerprint = ?", customerID, fingerprint).
		Count(&count).Error
	return count > 0, err
}
