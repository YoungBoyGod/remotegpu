package dao

import (
	"context"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// HostMetricDao 机器监控指标数据访问层
type HostMetricDao struct {
	db *gorm.DB
}

// NewHostMetricDao 创建DAO实例
func NewHostMetricDao(db *gorm.DB) *HostMetricDao {
	return &HostMetricDao{db: db}
}

// Create 创建监控记录
func (d *HostMetricDao) Create(ctx context.Context, metric *entity.HostMetric) error {
	return d.db.WithContext(ctx).Create(metric).Error
}

// BatchCreate 批量创建监控记录
func (d *HostMetricDao) BatchCreate(ctx context.Context, metrics []*entity.HostMetric) error {
	if len(metrics) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Create(metrics).Error
}

// GetLatest 获取机器的最新监控数据
func (d *HostMetricDao) GetLatest(ctx context.Context, hostID string) (*entity.HostMetric, error) {
	var metric entity.HostMetric
	err := d.db.WithContext(ctx).
		Where("host_id = ?", hostID).
		Order("collected_at DESC").
		First(&metric).Error
	return &metric, err
}

// GetHistory 获取机器的历史监控数据
func (d *HostMetricDao) GetHistory(ctx context.Context, hostID string, start, end time.Time, limit int) ([]*entity.HostMetric, error) {
	var metrics []*entity.HostMetric
	query := d.db.WithContext(ctx).
		Where("host_id = ?", hostID).
		Where("collected_at BETWEEN ? AND ?", start, end).
		Order("collected_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&metrics).Error
	return metrics, err
}

// DeleteOldRecords 删除旧的监控记录
func (d *HostMetricDao) DeleteOldRecords(ctx context.Context, before time.Time) error {
	return d.db.WithContext(ctx).
		Where("collected_at < ?", before).
		Delete(&entity.HostMetric{}).Error
}
