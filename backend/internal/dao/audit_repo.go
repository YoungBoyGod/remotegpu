package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type AuditDao struct {
	db *gorm.DB
}

func NewAuditDao(db *gorm.DB) *AuditDao {
	return &AuditDao{db: db}
}

func (d *AuditDao) Create(ctx context.Context, log *entity.AuditLog) error {
	return d.db.WithContext(ctx).Create(log).Error
}

// ListParams 审计日志查询参数
type AuditListParams struct {
	Page         int
	PageSize     int
	Action       string
	ResourceType string
	Username     string
	StartTime    string
	EndTime      string
}

func (d *AuditDao) List(ctx context.Context, params AuditListParams) ([]entity.AuditLog, int64, error) {
	var logs []entity.AuditLog
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.AuditLog{})

	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}
	if params.ResourceType != "" {
		query = query.Where("resource_type = ?", params.ResourceType)
	}
	if params.Username != "" {
		query = query.Where("username LIKE ?", "%"+params.Username+"%")
	}
	if params.StartTime != "" {
		query = query.Where("created_at >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		query = query.Where("created_at <= ?", params.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(params.PageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
