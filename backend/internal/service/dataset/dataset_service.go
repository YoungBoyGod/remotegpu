package dataset

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type DatasetService struct {
	datasetDao *dao.DatasetDao
}

func NewDatasetService(db *gorm.DB) *DatasetService {
	return &DatasetService{
		datasetDao: dao.NewDatasetDao(db),
	}
}

func (s *DatasetService) ListDatasets(ctx context.Context, customerID uint, page, pageSize int) ([]entity.Dataset, int64, error) {
	return s.datasetDao.ListByCustomerID(ctx, customerID, page, pageSize)
}

// GetDataset 根据ID获取数据集
// @author Claude
// @description 获取数据集详情，用于权限校验
// @modified 2026-02-04
func (s *DatasetService) GetDataset(ctx context.Context, id uint) (*entity.Dataset, error) {
	return s.datasetDao.FindByID(ctx, id)
}

// ValidateOwnership 验证数据集所有权
// @author Claude
// @description 验证数据集是否属于指定用户
// @modified 2026-02-04
func (s *DatasetService) ValidateOwnership(ctx context.Context, datasetID uint, customerID uint) error {
	dataset, err := s.datasetDao.FindByID(ctx, datasetID)
	if err != nil {
		return err
	}
	if dataset.CustomerID != customerID {
		return entity.ErrUnauthorized
	}
	return nil
}
