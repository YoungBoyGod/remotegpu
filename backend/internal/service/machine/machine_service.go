package machine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
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

func (s *MachineService) GetHost(ctx context.Context, id string) (*entity.Host, error) {
	return s.machineDao.FindByID(ctx, id)
}

func (s *MachineService) CreateMachine(ctx context.Context, host *entity.Host) error {
	// CodeX 2026-02-04: validate unique IP/hostname before create.
	if host.ID == "" {
		host.ID = deriveHostID(host)
	}
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

func (s *MachineService) CollectHostSpec(ctx context.Context, host *entity.Host, info *SystemInfoSnapshot) error {
	if host == nil || info == nil {
		return fmt.Errorf("missing host or system info")
	}
	if info.Hostname != "" {
		host.Hostname = info.Hostname
	}
	if host.Name == "" && info.Hostname != "" {
		host.Name = info.Hostname
	}
	if info.CPUCores > 0 {
		host.TotalCPU = info.CPUCores
		host.CPUInfo = fmt.Sprintf("%d cores", info.CPUCores)
	}
	if info.MemoryTotalGB > 0 {
		host.TotalMemoryGB = info.MemoryTotalGB
	}
	if info.DiskTotalGB > 0 {
		host.TotalDiskGB = info.DiskTotalGB
	}
	if host.TotalCPU <= 0 || host.TotalMemoryGB <= 0 {
		return fmt.Errorf("invalid collected spec")
	}
	if info.Collected {
		host.Status = "idle"
		host.HealthStatus = "healthy"
	}
	host.NeedsCollect = false
	return s.machineDao.UpdateCollectFields(ctx, host)
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
		if host.ID == "" {
			host.ID = deriveHostID(&host)
		}
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

func deriveHostID(host *entity.Host) string {
	if host.Hostname != "" {
		return host.Hostname
	}
	if host.IPAddress != "" {
		return host.IPAddress
	}
	return "host-" + uuid.NewString()
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

func (s *MachineService) ListNeedCollect(ctx context.Context, limit int) ([]entity.Host, error) {
	return s.machineDao.ListNeedCollect(ctx, limit)
}

func (s *MachineService) UpdateHostSpec(ctx context.Context, host *entity.Host) error {
	if host == nil {
		return fmt.Errorf("missing host")
	}
	return s.machineDao.UpdateCollectFields(ctx, host)
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

// DeleteMachine 删除机器
func (s *MachineService) DeleteMachine(ctx context.Context, hostID string) error {
	return s.machineDao.Delete(ctx, hostID)
}

// UpdateStatus 更新机器状态
func (s *MachineService) UpdateStatus(ctx context.Context, hostID string, status string) error {
	return s.machineDao.UpdateStatus(ctx, hostID, status)
}
