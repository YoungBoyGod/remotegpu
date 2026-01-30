package service

import (
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
