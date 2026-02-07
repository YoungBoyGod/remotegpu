package dao

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type ImageDao struct {
	db *gorm.DB
}

func NewImageDao(db *gorm.DB) *ImageDao {
	return &ImageDao{db: db}
}

func (d *ImageDao) Create(ctx context.Context, img *entity.Image) error {
	return d.db.WithContext(ctx).Create(img).Error
}

func (d *ImageDao) FindByID(ctx context.Context, id uint) (*entity.Image, error) {
	var img entity.Image
	if err := d.db.WithContext(ctx).First(&img, id).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

type ImageListParams struct {
	Page      int
	PageSize  int
	Category  string
	Framework string
	Status    string
}

func (d *ImageDao) List(ctx context.Context, params ImageListParams) ([]entity.Image, int64, error) {
	var images []entity.Image
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.Image{})

	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.Framework != "" {
		query = query.Where("framework = ?", params.Framework)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

func (d *ImageDao) Update(ctx context.Context, img *entity.Image) error {
	return d.db.WithContext(ctx).Save(img).Error
}

func (d *ImageDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&entity.Image{}, id).Error
}

// FindByName 按镜像名称查找（用于同步去重）
func (d *ImageDao) FindByName(ctx context.Context, name string) (*entity.Image, error) {
	var img entity.Image
	if err := d.db.WithContext(ctx).Where("name = ?", name).First(&img).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

// ListExistingNames 批量查询已存在的镜像名称（用于同步去重）
func (d *ImageDao) ListExistingNames(ctx context.Context, names []string) (map[string]bool, error) {
	if len(names) == 0 {
		return map[string]bool{}, nil
	}
	var existing []struct{ Name string }
	if err := d.db.WithContext(ctx).Model(&entity.Image{}).
		Select("name").Where("name IN ?", names).
		Find(&existing).Error; err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(existing))
	for _, e := range existing {
		result[e.Name] = true
	}
	return result, nil
}
