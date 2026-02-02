package service

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

// StorageRouterService 存储路由服务
type StorageRouterService struct {
	storageSourceDao *dao.StorageSourceDao
}

// NewStorageRouterService 创建存储路由服务
func NewStorageRouterService() *StorageRouterService {
	return &StorageRouterService{
		storageSourceDao: dao.NewStorageSourceDao(),
	}
}

// SelectStorageByUserType 根据用户类型选择存储源
// userType: "external" 外部用户, "internal" 内部用户
func (s *StorageRouterService) SelectStorageByUserType(userType string) (*entity.StorageSource, error) {
	var isPublic bool
	if userType == "external" {
		// 外部用户使用公有云存储
		isPublic = true
	} else {
		// 内部用户使用私有云存储
		isPublic = false
	}

	// 获取对应类型的存储源列表
	sources, err := s.storageSourceDao.ListByPublic(isPublic, "active")
	if err != nil {
		return nil, fmt.Errorf("获取存储源列表失败: %w", err)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("没有可用的存储源")
	}

	// 返回优先级最高的存储源(priority最小)
	return sources[0], nil
}

// SelectStorageByRegion 根据节点区域选择最近的存储源
func (s *StorageRouterService) SelectStorageByRegion(region string) (*entity.StorageSource, error) {
	// 获取所有活跃的存储源
	sources, err := s.storageSourceDao.List("active")
	if err != nil {
		return nil, fmt.Errorf("获取存储源列表失败: %w", err)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("没有可用的存储源")
	}

	// 优先选择同区域的存储源
	for _, source := range sources {
		if source.Region == region {
			return source, nil
		}
	}

	// 如果没有同区域的存储源,返回优先级最高的
	return sources[0], nil
}

// SelectStorageForDataset 为数据集选择存储源
// 综合考虑用户类型和节点区域
func (s *StorageRouterService) SelectStorageForDataset(userType, region string) (*entity.StorageSource, error) {
	// 首先根据用户类型筛选
	var isPublic bool
	if userType == "external" {
		isPublic = true
	} else {
		isPublic = false
	}

	sources, err := s.storageSourceDao.ListByPublic(isPublic, "active")
	if err != nil {
		return nil, fmt.Errorf("获取存储源列表失败: %w", err)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("没有可用的存储源")
	}

	// 在同类型的存储源中,优先选择同区域的
	if region != "" {
		for _, source := range sources {
			if source.Region == region {
				return source, nil
			}
		}
	}

	// 如果没有同区域的,返回优先级最高的
	return sources[0], nil
}

// GetBestReplicaForNode 为节点选择最佳的数据集副本
func (s *StorageRouterService) GetBestReplicaForNode(replicas []*entity.DatasetReplica, nodeRegion string) (*entity.DatasetReplica, error) {
	if len(replicas) == 0 {
		return nil, fmt.Errorf("没有可用的副本")
	}

	// 优先选择主副本
	for _, replica := range replicas {
		if replica.IsPrimary && replica.Status == "ready" {
			// 如果主副本的存储源在同区域,直接返回
			if replica.StorageSource != nil && replica.StorageSource.Region == nodeRegion {
				return replica, nil
			}
		}
	}

	// 如果主副本不在同区域,查找同区域的副本
	if nodeRegion != "" {
		for _, replica := range replicas {
			if replica.Status == "ready" && replica.StorageSource != nil {
				if replica.StorageSource.Region == nodeRegion {
					return replica, nil
				}
			}
		}
	}

	// 如果没有同区域的副本,返回第一个可用的副本
	for _, replica := range replicas {
		if replica.Status == "ready" {
			return replica, nil
		}
	}

	return nil, fmt.Errorf("没有可用的副本")
}
