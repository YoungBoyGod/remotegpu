package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HostService 主机服务
type HostService struct {
	hostDao HostDaoInterface
	gpuDao  GPUDaoInterface
}

// NewHostService 创建主机服务实例
func NewHostService() *HostService {
	return &HostService{
		hostDao: dao.NewHostDao(),
		gpuDao:  dao.NewGPUDao(),
	}
}

// Create 创建主机
func (s *HostService) Create(host *entity.Host) error {
	if host.ID == "" {
		host.ID = fmt.Sprintf("host-%s", uuid.New().String()[:8])
	}
	host.RegisteredAt = time.Now()
	return s.hostDao.Create(host)
}

// GetByID 根据ID获取主机
func (s *HostService) GetByID(id string) (*entity.Host, error) {
	return s.hostDao.GetByID(id)
}

// Update 更新主机
func (s *HostService) Update(host *entity.Host) error {
	return s.hostDao.Update(host)
}

// Delete 删除主机
func (s *HostService) Delete(id string) error {
	// 先删除关联的GPU
	if err := s.gpuDao.DeleteByHostID(id); err != nil {
		return err
	}
	return s.hostDao.Delete(id)
}

// List 获取主机列表
func (s *HostService) List(page, pageSize int) ([]*entity.Host, int64, error) {
	return s.hostDao.List(page, pageSize)
}

// UpdateStatus 更新主机状态
func (s *HostService) UpdateStatus(id, status string) error {
	return s.hostDao.UpdateStatus(id, status)
}

// Heartbeat 更新心跳
func (s *HostService) Heartbeat(id string) error {
	_, err := s.hostDao.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("主机不存在")
		}
		return err
	}
	return s.hostDao.UpdateHeartbeat(id)
}
