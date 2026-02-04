package machine

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

var (
	ErrHostDuplicateIP       = errors.New("host ip already exists")
	ErrHostDuplicateHostname = errors.New("host hostname already exists")
)

func (s *MachineService) ListMachines(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entity.Host, int64, error) {
	return s.machineDao.List(ctx, page, pageSize, filters)
}

func (s *MachineService) CreateMachine(ctx context.Context, host *entity.Host) error {
	// CodeX 2026-02-04: validate unique IP/hostname before create.
	if host.IPAddress != "" {
		if _, err := s.machineDao.FindByIPAddress(ctx, host.IPAddress); err == nil {
			return ErrHostDuplicateIP
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if host.Hostname != "" {
		if _, err := s.machineDao.FindByHostname(ctx, host.Hostname); err == nil {
			return ErrHostDuplicateHostname
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	return s.machineDao.Create(ctx, host)
}

func (s *MachineService) ImportMachines(ctx context.Context, hosts []entity.Host) error {
	// CodeX 2026-02-04: skip duplicates by IP/hostname during import.
	if len(hosts) == 0 {
		return nil
	}

	ips := make([]string, 0, len(hosts))
	hostnames := make([]string, 0, len(hosts))
	for _, host := range hosts {
		if host.IPAddress != "" {
			ips = append(ips, host.IPAddress)
		}
		if host.Hostname != "" {
			hostnames = append(hostnames, host.Hostname)
		}
	}

	existing, err := s.machineDao.FindExistingKeys(ctx, uniqueStrings(ips), uniqueStrings(hostnames))
	if err != nil {
		return err
	}

	for _, host := range hosts {
		key := dao.HostKey{IPAddress: host.IPAddress, Hostname: host.Hostname}
		if _, ok := existing[key]; ok {
			continue
		}
		if err := s.machineDao.Create(ctx, &host); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				continue
			}
			return fmt.Errorf("import host %s failed: %w", formatHostKey(host), err)
		}
	}

	return nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func formatHostKey(host entity.Host) string {
	if host.Hostname != "" && host.IPAddress != "" {
		return host.Hostname + " (" + host.IPAddress + ")"
	}
	if host.Hostname != "" {
		return host.Hostname
	}
	return host.IPAddress
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

// Count 获取机器总数
// @modified 2026-02-04
func (s *MachineService) Count(ctx context.Context) (int64, error) {
	return s.machineDao.Count(ctx)
}

// GetStatusStats 获取各状态机器统计
// @description 用于仪表盘展示机器状态分布
// @modified 2026-02-04
func (s *MachineService) GetStatusStats(ctx context.Context) (map[string]int64, error) {
	return s.machineDao.GetStatusStats(ctx)
}
