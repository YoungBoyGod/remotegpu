package service

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type DatasetService struct {
	datasetDao *dao.DatasetRepo
}

func NewDatasetService(db *gorm.DB) *DatasetService {
	return &DatasetService{
		datasetDao: dao.NewDatasetRepo(db),
	}
}

func (s *DatasetService) ListDatasets(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Dataset, int64, error) {
	return s.datasetDao.ListByCustomerID(ctx, customerID, page, pageSize)
}