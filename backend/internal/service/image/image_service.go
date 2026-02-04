package image

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type ImageService struct {
	imageDao *dao.ImageDao
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		imageDao: dao.NewImageDao(db),
	}
}

func (s *ImageService) List(ctx context.Context, params dao.ImageListParams) ([]entity.Image, int64, error) {
	return s.imageDao.List(ctx, params)
}

func (s *ImageService) Create(ctx context.Context, img *entity.Image) error {
	return s.imageDao.Create(ctx, img)
}

func (s *ImageService) GetByID(ctx context.Context, id uint) (*entity.Image, error) {
	return s.imageDao.FindByID(ctx, id)
}

func (s *ImageService) Update(ctx context.Context, img *entity.Image) error {
	return s.imageDao.Update(ctx, img)
}

func (s *ImageService) Delete(ctx context.Context, id uint) error {
	return s.imageDao.Delete(ctx, id)
}

func (s *ImageService) Sync(ctx context.Context) error {
	// TODO: Implement actual sync logic with Harbor or local registry
	// For now, we'll just mock it or do nothing
	return nil
}
