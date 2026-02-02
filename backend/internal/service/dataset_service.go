package service

import (
	"context"
	"fmt"
	"io"

	"github.com/YoungBoyGod/remotegpu/internal/adapter"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

// DatasetService 数据集管理服务
type DatasetService struct {
	datasetDao        *dao.DatasetDao
	replicaDao        *dao.DatasetReplicaDao
	storageSourceDao  *dao.StorageSourceDao
	storageRouter     *StorageRouterService
}

// NewDatasetService 创建数据集服务
func NewDatasetService() *DatasetService {
	return &DatasetService{
		datasetDao:       dao.NewDatasetDao(),
		replicaDao:       dao.NewDatasetReplicaDao(),
		storageSourceDao: dao.NewStorageSourceDao(),
		storageRouter:    NewStorageRouterService(),
	}
}

// CreateDataset 创建数据集
func (s *DatasetService) CreateDataset(dataset *entity.Dataset) error {
	// 验证数据集信息
	if err := s.validateDataset(dataset); err != nil {
		return fmt.Errorf("数据集验证失败: %w", err)
	}

	// 设置初始状态
	dataset.Status = "uploading"

	// 创建数据集记录
	return s.datasetDao.Create(dataset)
}

// GetDataset 获取数据集
func (s *DatasetService) GetDataset(id uint) (*entity.Dataset, error) {
	return s.datasetDao.GetByIDWithReplicas(id)
}

// ListUserDatasets 获取用户的数据集列表
func (s *DatasetService) ListUserDatasets(userID uint, status string) ([]*entity.Dataset, error) {
	return s.datasetDao.ListByUser(userID, status)
}

// ListWorkspaceDatasets 获取工作空间的数据集列表
func (s *DatasetService) ListWorkspaceDatasets(workspaceID uint, status string) ([]*entity.Dataset, error) {
	return s.datasetDao.ListByWorkspace(workspaceID, status)
}

// UploadDataset 上传数据集到存储源
func (s *DatasetService) UploadDataset(ctx context.Context, datasetID uint, userType string, reader io.Reader, size int64) error {
	// 获取数据集信息
	dataset, err := s.datasetDao.GetByID(datasetID)
	if err != nil {
		return fmt.Errorf("获取数据集失败: %w", err)
	}

	// 根据用户类型选择存储源
	storageSource, err := s.storageRouter.SelectStorageByUserType(userType)
	if err != nil {
		return fmt.Errorf("选择存储源失败: %w", err)
	}

	// 创建存储适配器
	storageAdapter, err := s.createAdapter(storageSource)
	if err != nil {
		return fmt.Errorf("创建存储适配器失败: %w", err)
	}

	// 生成对象键
	objectKey := fmt.Sprintf("datasets/%d/%s", datasetID, dataset.Name)

	// 上传文件
	if err := storageAdapter.Upload(ctx, objectKey, reader, size); err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}

	// 更新数据集状态
	dataset.Status = "ready"
	dataset.Size = size
	if err := s.datasetDao.Update(dataset); err != nil {
		return fmt.Errorf("更新数据集状态失败: %w", err)
	}

	// 创建副本记录
	replica := &entity.DatasetReplica{
		DatasetID:       datasetID,
		StorageSourceID: storageSource.ID,
		Path:            objectKey,
		Size:            size,
		Status:          "ready",
		IsPrimary:       true,
	}
	if err := s.replicaDao.Create(replica); err != nil {
		return fmt.Errorf("创建副本记录失败: %w", err)
	}

	return nil
}

// DeleteDataset 删除数据集
func (s *DatasetService) DeleteDataset(id uint) error {
	// 获取数据集及其副本
	dataset, err := s.datasetDao.GetByIDWithReplicas(id)
	if err != nil {
		return fmt.Errorf("获取数据集失败: %w", err)
	}

	// 删除所有副本
	for _, replica := range dataset.Replicas {
		if err := s.deleteReplica(replica); err != nil {
			// 记录错误但继续删除其他副本
			fmt.Printf("删除副本失败: %v\n", err)
		}
	}

	// 删除数据集记录
	return s.datasetDao.Delete(id)
}

// validateDataset 验证数据集信息
func (s *DatasetService) validateDataset(dataset *entity.Dataset) error {
	if dataset.Name == "" {
		return fmt.Errorf("数据集名称不能为空")
	}
	if dataset.UserID == 0 {
		return fmt.Errorf("用户ID不能为空")
	}
	return nil
}

// deleteReplica 删除副本
func (s *DatasetService) deleteReplica(replica *entity.DatasetReplica) error {
	// 获取存储源
	storageSource, err := s.storageSourceDao.GetByID(replica.StorageSourceID)
	if err != nil {
		return err
	}

	// 创建存储适配器
	storageAdapter, err := s.createAdapter(storageSource)
	if err != nil {
		return err
	}

	// 删除存储中的文件
	ctx := context.Background()
	if err := storageAdapter.Delete(ctx, replica.Path); err != nil {
		return err
	}

	// 删除副本记录
	return s.replicaDao.Delete(replica.ID)
}

// createAdapter 创建存储适配器
func (s *DatasetService) createAdapter(source *entity.StorageSource) (adapter.StorageAdapter, error) {
	switch source.Type {
	case "minio":
		return adapter.NewMinIOAdapter(source.Endpoint, source.AccessKey, source.SecretKey, source.Bucket, false)
	case "oss":
		return adapter.NewOSSAdapter(source.Endpoint, source.AccessKey, source.SecretKey, source.Bucket)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", source.Type)
	}
}
