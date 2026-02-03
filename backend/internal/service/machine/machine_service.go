package machine

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
	// 实际实现中，应该使用 gorm 的批量插入
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

	// 判断使用公网 IP 还是内网 IP（例如，如果在 NAT 后面）
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
