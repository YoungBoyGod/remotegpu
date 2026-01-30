package service

import (
	"fmt"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
)

// GPUService GPU服务
type GPUService struct {
	gpuDao *dao.GPUDao
}

// NewGPUService 创建GPU服务实例
func NewGPUService() *GPUService {
	return &GPUService{
		gpuDao: dao.NewGPUDao(),
	}
}

// Create 创建GPU
func (s *GPUService) Create(gpu *entity.GPU) error {
	return s.gpuDao.Create(gpu)
}

// GetByID 根据ID获取GPU
func (s *GPUService) GetByID(id uint) (*entity.GPU, error) {
	return s.gpuDao.GetByID(id)
}

// GetByHostID 根据主机ID获取GPU列表
func (s *GPUService) GetByHostID(hostID string) ([]*entity.GPU, error) {
	return s.gpuDao.GetByHostID(hostID)
}

// Update 更新GPU
func (s *GPUService) Update(gpu *entity.GPU) error {
	return s.gpuDao.Update(gpu)
}

// Delete 删除GPU
func (s *GPUService) Delete(id uint) error {
	return s.gpuDao.Delete(id)
}

// UpdateStatus 更新GPU状态
func (s *GPUService) UpdateStatus(id uint, status string) error {
	return s.gpuDao.UpdateStatus(id, status)
}

// List 分页获取GPU列表
func (s *GPUService) List(page, pageSize int) ([]*entity.GPU, int64, error) {
	return s.gpuDao.List(page, pageSize)
}

// GetByStatus 根据状态获取GPU列表
func (s *GPUService) GetByStatus(status string) ([]*entity.GPU, error) {
	return s.gpuDao.GetByStatus(status)
}

// Allocate 分配GPU到环境
func (s *GPUService) Allocate(id uint, envID string) error {
	gpu, err := s.gpuDao.GetByID(id)
	if err != nil {
		return err
	}
	if gpu.Status != "available" {
		return fmt.Errorf("GPU不可用，当前状态: %s", gpu.Status)
	}
	return s.gpuDao.Allocate(id, envID)
}

// Release 释放GPU
func (s *GPUService) Release(id uint) error {
	gpu, err := s.gpuDao.GetByID(id)
	if err != nil {
		return err
	}
	if gpu.Status != "allocated" {
		return fmt.Errorf("GPU未分配，当前状态: %s", gpu.Status)
	}
	return s.gpuDao.Release(id)
}
