package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type DocumentDao struct {
	db *gorm.DB
}

func NewDocumentDao(db *gorm.DB) *DocumentDao {
	return &DocumentDao{db: db}
}

func (d *DocumentDao) Create(ctx context.Context, doc *entity.Document) error {
	return d.db.WithContext(ctx).Create(doc).Error
}

func (d *DocumentDao) FindByID(ctx context.Context, id uint) (*entity.Document, error) {
	var doc entity.Document
	err := d.db.WithContext(ctx).Preload("Uploader").First(&doc, "id = ?", id).Error
	return &doc, err
}

func (d *DocumentDao) Update(ctx context.Context, id uint, fields map[string]any) error {
	return d.db.WithContext(ctx).Model(&entity.Document{}).Where("id = ?", id).Updates(fields).Error
}

func (d *DocumentDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&entity.Document{}, "id = ?", id).Error
}

func (d *DocumentDao) List(ctx context.Context, page, pageSize int, category, keyword string) ([]entity.Document, int64, error) {
	var docs []entity.Document
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.Document{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("title ILIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Uploader").
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&docs).Error

	return docs, total, err
}

func (d *DocumentDao) ListCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := d.db.WithContext(ctx).
		Model(&entity.Document{}).
		Distinct("category").
		Pluck("category", &categories).Error
	return categories, err
}
