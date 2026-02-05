package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type MachineEnrollmentDao struct {
	db *gorm.DB
}

func NewMachineEnrollmentDao(db *gorm.DB) *MachineEnrollmentDao {
	return &MachineEnrollmentDao{db: db}
}

func (d *MachineEnrollmentDao) Create(ctx context.Context, enrollment *entity.MachineEnrollment) error {
	return d.db.WithContext(ctx).Create(enrollment).Error
}

func (d *MachineEnrollmentDao) FindByID(ctx context.Context, id uint) (*entity.MachineEnrollment, error) {
	var enrollment entity.MachineEnrollment
	if err := d.db.WithContext(ctx).First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (d *MachineEnrollmentDao) ListByCustomer(ctx context.Context, customerID uint, page, pageSize int) ([]entity.MachineEnrollment, int64, error) {
	var list []entity.MachineEnrollment
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.MachineEnrollment{}).
		Where("customer_id = ?", customerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d *MachineEnrollmentDao) ListPending(ctx context.Context, limit int) ([]entity.MachineEnrollment, error) {
	var list []entity.MachineEnrollment
	query := d.db.WithContext(ctx).Model(&entity.MachineEnrollment{}).
		Where("status = ?", "pending").Order("created_at asc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *MachineEnrollmentDao) UpdateStatus(ctx context.Context, id uint, status, errorMessage, hostID string) error {
	updates := map[string]interface{}{
		"status":        status,
		"error_message": errorMessage,
		"host_id":       hostID,
	}
	return d.db.WithContext(ctx).Model(&entity.MachineEnrollment{}).Where("id = ?", id).Updates(updates).Error
}
