package dao

import (
	"context"

	"gorm.io/gorm"
)

type BaseDao[T any] struct {
	db *gorm.DB
}

func NewBaseDao[T any](db *gorm.DB) *BaseDao[T] {
	return &BaseDao[T]{db: db}
}

func (d *BaseDao[T]) Create(ctx context.Context, entity *T) error {
	return d.db.WithContext(ctx).Create(entity).Error
}

func (d *BaseDao[T]) Update(ctx context.Context, entity *T) error {
	return d.db.WithContext(ctx).Save(entity).Error
}

func (d *BaseDao[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return d.db.WithContext(ctx).Delete(&entity, id).Error
}

func (d *BaseDao[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := d.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (d *BaseDao[T]) FindByUUID(ctx context.Context, uuid string) (*T, error) {
	var entity T
	if err := d.db.WithContext(ctx).Where("uuid = ?", uuid).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
