package service

import (
	"context"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

type MachineService struct {
	machineDao *dao.MachineDao
}

func NewMachineService(db *gorm.DB) *MachineService {
	return &MachineService{
		machineDao: dao.NewMachineDao(db),
	}
}

func (s *MachineService) ListMachines(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entity.Host, int64, error) {
	return s.machineDao.List(ctx, page, pageSize, filters)
}

func (s *MachineService) CreateMachine(ctx context.Context, host *entity.Host) error {
	// TODO: Add validation logic (e.g., check IP uniqueness)
	return s.machineDao.Create(ctx, host)
}

func (s *MachineService) ImportMachines(ctx context.Context, hosts []entity.Host) error {
	// In a real implementation, this would use gorm's batch insert
	for _, host := range hosts {
		if err := s.machineDao.Create(ctx, &host); err != nil {
			return err
		}
	}
	return nil
}

func (s *MachineService) GetConnectionInfo(ctx context.Context, hostID string) (map[string]interface{}, error) {
	host, err := s.machineDao.FindByID(ctx, hostID)
	if err != nil {
		return nil, err
	}

	// Logic to determine public IP vs internal IP (e.g., if behind NAT)
	connectIP := host.PublicIP
	if connectIP == "" {
		connectIP = host.IPAddress
	}

	return map[string]interface{}{
		"ssh_command": "ssh root@" + connectIP, // Simplified
		"host":        connectIP,
		"port":        host.SSHPort,
	}, nil
}
