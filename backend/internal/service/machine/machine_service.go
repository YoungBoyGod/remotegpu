package machine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MachineService struct {
	machineDao    *dao.MachineDao
	hostMetricDao *dao.HostMetricDao
	allocationDao *dao.AllocationDao
	db            *gorm.DB
	statusCache   *HostStatusCache
}

func NewMachineService(db *gorm.DB) *MachineService {
	return &MachineService{
		machineDao:    dao.NewMachineDao(db),
		hostMetricDao: dao.NewHostMetricDao(db),
		allocationDao: dao.NewAllocationDao(db),
		db:            db,
	}
}

// SetStatusCache 注入设备状态缓存（可选，不影响现有调用）
func (s *MachineService) SetStatusCache(c *HostStatusCache) {
	s.statusCache = c
}

// GetStatusCache 获取设备状态缓存
func (s *MachineService) GetStatusCache() *HostStatusCache {
	return s.statusCache
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

	// 修复 P0 安全问题：加密 SSH 密码
	if host.SSHPassword != "" {
		encrypted, err := crypto.EncryptAES256GCM(host.SSHPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt SSH password: %w", err)
		}
		host.SSHPassword = encrypted
	}
	// 加密 SSH Key
	if host.SSHKey != "" {
		encrypted, err := crypto.EncryptAES256GCM(host.SSHKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt SSH key: %w", err)
		}
		host.SSHKey = encrypted
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
		host.DeviceStatus = "online"
		host.AllocationStatus = "idle"
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

		// 修复 P0 安全问题：加密 SSH 密码和 SSH Key
		if host.SSHPassword != "" {
			encrypted, err := crypto.EncryptAES256GCM(host.SSHPassword)
			if err != nil {
				return fmt.Errorf("failed to encrypt SSH password for host %s: %w", formatHostKey(host), err)
			}
			host.SSHPassword = encrypted
		}
		if host.SSHKey != "" {
			encrypted, err := crypto.EncryptAES256GCM(host.SSHKey)
			if err != nil {
				return fmt.Errorf("failed to encrypt SSH key for host %s: %w", formatHostKey(host), err)
			}
			host.SSHKey = encrypted
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

	// SSH 连接主机优先级：ssh_host > public_ip > ip_address
	connectHost := host.SSHHost
	if connectHost == "" {
		connectHost = host.PublicIP
	}
	if connectHost == "" {
		connectHost = host.IPAddress
	}

	username := host.SSHUsername
	if username == "" {
		username = "root"
	}

	port := host.SSHPort
	if port == 0 {
		port = 22
	}

	// 解密 SSH 密码
	password := ""
	if host.SSHPassword != "" {
		decrypted, err := crypto.DecryptAES256GCM(host.SSHPassword)
		if err == nil {
			password = decrypted
		}
	}

	// 生成 SSH 连接命令
	sshCommand := fmt.Sprintf("ssh -p %d %s@%s", port, username, connectHost)

	// 解密 VNC 密码
	vncPassword := ""
	if host.VNCPassword != "" {
		decrypted, err := crypto.DecryptAES256GCM(host.VNCPassword)
		if err == nil {
			vncPassword = decrypted
		}
	}

	return map[string]interface{}{
		"ssh": map[string]interface{}{
			"username": username,
			"host":     connectHost,
			"port":     port,
			"password": password,
		},
		"ssh_command": sshCommand,
		"jupyter": map[string]interface{}{
			"url":   host.JupyterURL,
			"token": host.JupyterToken,
		},
		"vnc": map[string]interface{}{
			"url":      host.VNCURL,
			"password": vncPassword,
		},
	}, nil
}

// GetMachineDetail 获取机器详情（包含 SSH 连接信息）
func (s *MachineService) GetMachineDetail(ctx context.Context, hostID string) (map[string]interface{}, error) {
	host, err := s.machineDao.FindByID(ctx, hostID)
	if err != nil {
		return nil, err
	}

	// SSH 连接主机优先级：ssh_host > public_ip > ip_address
	connectHost := host.SSHHost
	if connectHost == "" {
		connectHost = host.PublicIP
	}
	if connectHost == "" {
		connectHost = host.IPAddress
	}

	username := host.SSHUsername
	if username == "" {
		username = "root"
	}

	port := host.SSHPort
	if port == 0 {
		port = 22
	}

	sshCommand := fmt.Sprintf("ssh -p %d %s@%s", port, username, connectHost)

	// 解密 SSH 密码
	password := ""
	if host.SSHPassword != "" {
		decrypted, err := crypto.DecryptAES256GCM(host.SSHPassword)
		if err == nil {
			password = decrypted
		}
	}

	// 解密 VNC 密码
	vncPassword := ""
	if host.VNCPassword != "" {
		decrypted, err := crypto.DecryptAES256GCM(host.VNCPassword)
		if err == nil {
			vncPassword = decrypted
		}
	}

	return map[string]interface{}{
		"id":            host.ID,
		"name":          host.Name,
		"hostname":      host.Hostname,
		"region":        host.Region,
		"ip_address":    host.IPAddress,
		"public_ip":     host.PublicIP,
		"status":            host.Status,
		"device_status":     host.DeviceStatus,
		"allocation_status": host.AllocationStatus,
		"health_status":     host.HealthStatus,
		"os_type":       host.OSType,
		"os_version":    host.OSVersion,
		"cpu_info":      host.CPUInfo,
		"total_cpu":     host.TotalCPU,
		"total_memory_gb": host.TotalMemoryGB,
		"total_disk_gb":   host.TotalDiskGB,
		"ssh_host":      host.SSHHost,
		"ssh_port":      port,
		"ssh_username":  username,
		"ssh_password":  password,
		"ssh_command":   sshCommand,
		"agent_port":    host.AgentPort,
		"jupyter_url":   host.JupyterURL,
		"jupyter_token": host.JupyterToken,
		"vnc_url":       host.VNCURL,
		"vnc_password":  vncPassword,
		"external_ip":           host.ExternalIP,
		"external_ssh_port":     host.ExternalSSHPort,
		"external_jupyter_port": host.ExternalJupyterPort,
		"external_vnc_port":     host.ExternalVNCPort,
		"nginx_domain":          host.NginxDomain,
		"nginx_config_path":     host.NginxConfigPath,
		"last_heartbeat": host.LastHeartbeat,
		"created_at":    host.CreatedAt,
		"updated_at":    host.UpdatedAt,
		"gpus":          host.GPUs,
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

// UpdateMachine 更新机器信息
func (s *MachineService) UpdateMachine(ctx context.Context, hostID string, fields map[string]interface{}) error {
	// IP 唯一性校验：如果更新了 ip_address，检查是否与其他机器冲突
	if ip, ok := fields["ip_address"]; ok {
		if ipStr, _ := ip.(string); ipStr != "" {
			existing, err := s.machineDao.FindByIPAddress(ctx, ipStr)
			if err == nil && existing.ID != hostID {
				return ErrHostDuplicateIP
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
	}
	// Hostname 唯一性校验
	if hostname, ok := fields["hostname"]; ok {
		if hn, _ := hostname.(string); hn != "" {
			existing, err := s.machineDao.FindByHostname(ctx, hn)
			if err == nil && existing.ID != hostID {
				return ErrHostDuplicateHostname
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
	}

	// 如果更新了 SSH 密码，需要加密
	if password, ok := fields["ssh_password"]; ok {
		if pw, _ := password.(string); pw != "" {
			encrypted, err := crypto.EncryptAES256GCM(pw)
			if err != nil {
				return fmt.Errorf("failed to encrypt SSH password: %w", err)
			}
			fields["ssh_password"] = encrypted
		}
	}
	// 如果更新了 SSH Key，需要加密
	if key, ok := fields["ssh_key"]; ok {
		if k, _ := key.(string); k != "" {
			encrypted, err := crypto.EncryptAES256GCM(k)
			if err != nil {
				return fmt.Errorf("failed to encrypt SSH key: %w", err)
			}
			fields["ssh_key"] = encrypted
		}
	}
	return s.machineDao.UpdateFields(ctx, hostID, fields)
}

// DeleteMachine 删除机器（检查分配状态，防止误删已分配的机器）
func (s *MachineService) DeleteMachine(ctx context.Context, hostID string) error {
	host, err := s.machineDao.FindByID(ctx, hostID)
	if err != nil {
		return err
	}
	if host.AllocationStatus == "allocated" {
		return fmt.Errorf("cannot delete machine: currently allocated to a customer")
	}
	return s.machineDao.Delete(ctx, hostID)
}

// UpdateStatus 更新机器状态（兼容旧接口）
func (s *MachineService) UpdateStatus(ctx context.Context, hostID string, status string) error {
	return s.machineDao.UpdateStatus(ctx, hostID, status)
}

// UpdateDeviceStatus 更新设备在线状态
func (s *MachineService) UpdateDeviceStatus(ctx context.Context, hostID string, deviceStatus string) error {
	return s.machineDao.UpdateDeviceStatus(ctx, hostID, deviceStatus)
}

// UpdateAllocationStatus 更新分配状态
func (s *MachineService) UpdateAllocationStatus(ctx context.Context, hostID string, allocationStatus string) error {
	return s.machineDao.UpdateAllocationStatus(ctx, hostID, allocationStatus)
}

// ResolvePostMaintenanceStatus 取消维护时，根据是否有活跃分配决定恢复为 allocated 还是 idle
func (s *MachineService) ResolvePostMaintenanceStatus(ctx context.Context, hostID string) (string, error) {
	_, err := s.allocationDao.FindActiveByHostID(ctx, hostID)
	if err == nil {
		return "allocated", nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "idle", nil
	}
	return "", err
}

// BatchSetMaintenance 批量设置维护状态
// 进入维护：直接设为 maintenance
// 取消维护：检查是否有活跃分配，有则恢复为 allocated，否则恢复为 idle
func (s *MachineService) BatchSetMaintenance(ctx context.Context, hostIDs []string, maintenance bool) (int64, error) {
	if maintenance {
		return s.machineDao.BatchUpdateAllocationStatus(ctx, hostIDs, "maintenance")
	}
	// 取消维护：只将当前为 maintenance 的机器恢复为 idle
	return s.machineDao.BatchUpdateAllocationStatus(ctx, hostIDs, "idle")
}

// Heartbeat 处理 Agent 心跳上报，优先写入 Redis 缓存，减少数据库写入压力
func (s *MachineService) Heartbeat(ctx context.Context, hostID string) error {
	if s.statusCache != nil {
		status := &CachedHostStatus{
			HostID:        hostID,
			DeviceStatus:  "online",
			LastHeartbeat: time.Now(),
		}
		_ = s.statusCache.SetOnline(ctx, status)
		return nil
	}
	return s.machineDao.UpdateHeartbeat(ctx, hostID, "online")
}

// HeartbeatMetrics 心跳携带的监控指标
type HeartbeatMetrics struct {
	CPUUsagePercent    *float64
	MemoryUsagePercent *float64
	MemoryUsedGB       *int64
	DiskUsagePercent   *float64
	DiskUsedGB         *int64
	GPUMetrics         []GPUMetricData
}

// GPUMetricData 单个 GPU 的监控指标
type GPUMetricData struct {
	Index         int
	UUID          string
	Name          string
	UtilPercent   *float64
	MemoryUsedMB  *int
	MemoryTotalMB *int
	TemperatureC  *int
	PowerUsageW   *float64
}

// AgentRegistration Agent 注册信息
type AgentRegistration struct {
	AgentID    string
	MachineID  string
	Version    string
	Hostname   string
	IPAddress  string
	AgentPort  int
	MaxWorkers int
}

// HeartbeatWithMetrics 处理带监控指标的心跳上报
func (s *MachineService) HeartbeatWithMetrics(ctx context.Context, hostID string, metrics *HeartbeatMetrics) error {
	// 写入 Redis 缓存（如果可用）
	if s.statusCache != nil {
		status := &CachedHostStatus{
			HostID:        hostID,
			DeviceStatus:  "online",
			LastHeartbeat: time.Now(),
		}
		if metrics != nil {
			status.CPUUsage = metrics.CPUUsagePercent
			status.MemoryUsage = metrics.MemoryUsagePercent
			status.DiskUsage = metrics.DiskUsagePercent
			status.GPUCount = len(metrics.GPUMetrics)
		}
		_ = s.statusCache.SetOnline(ctx, status)
	} else {
		// 无缓存时直接写数据库
		if err := s.machineDao.UpdateHeartbeat(ctx, hostID, "online"); err != nil {
			return err
		}
	}

	// 监控指标仍然写入 host_metrics 表
	if metrics != nil {
		metric := &entity.HostMetric{
			HostID:             hostID,
			CPUUsagePercent:    metrics.CPUUsagePercent,
			MemoryUsagePercent: metrics.MemoryUsagePercent,
			MemoryUsedGB:       metrics.MemoryUsedGB,
			DiskUsagePercent:   metrics.DiskUsagePercent,
			DiskUsedGB:         metrics.DiskUsedGB,
			CollectedAt:        time.Now(),
		}
		if err := s.hostMetricDao.Create(ctx, metric); err != nil {
			return fmt.Errorf("保存心跳指标失败: %w", err)
		}
	}

	return nil
}

// RegisterAgent 处理 Agent 注册
func (s *MachineService) RegisterAgent(ctx context.Context, info *AgentRegistration) error {
	fields := make(map[string]interface{})
	fields["device_status"] = "online"

	if info.Hostname != "" {
		fields["hostname"] = info.Hostname
	}
	if info.AgentPort > 0 {
		fields["agent_port"] = info.AgentPort
	}

	return s.machineDao.UpdateFields(ctx, info.MachineID, fields)
}

// ListAgents 获取 Agent 列表（基于 hosts 表构建 Agent 视图）
func (s *MachineService) ListAgents(ctx context.Context) (map[string]interface{}, error) {
	hosts, err := s.machineDao.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var online, offline int64
	agents := make([]map[string]interface{}, 0, len(hosts))
	for _, h := range hosts {
		status := h.DeviceStatus
		if status == "" {
			status = "offline"
		}
		if status == "online" {
			online++
		} else {
			offline++
		}

		agent := map[string]interface{}{
			"agent_id":       h.ID,
			"machine_id":     h.ID,
			"machine_name":   h.Name,
			"ip_address":     h.IPAddress,
			"status":         status,
			"last_heartbeat": h.LastHeartbeat,
			"agent_port":     h.AgentPort,
			"region":         h.Region,
		}

		// GPU 摘要
		gpuNames := make([]string, 0, len(h.GPUs))
		for _, gpu := range h.GPUs {
			gpuNames = append(gpuNames, gpu.Name)
		}
		agent["gpu_count"] = len(h.GPUs)
		agent["gpu_models"] = gpuNames

		agents = append(agents, agent)
	}

	return map[string]interface{}{
		"total":   len(hosts),
		"online":  online,
		"offline": offline,
		"agents":  agents,
	}, nil
}

// GetMachineUsage 获取机器使用情况
func (s *MachineService) GetMachineUsage(ctx context.Context, hostID string) (map[string]interface{}, error) {
	host, err := s.machineDao.FindByID(ctx, hostID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"host_id":           host.ID,
		"status":            host.Status,
		"device_status":     host.DeviceStatus,
		"allocation_status": host.AllocationStatus,
	}

	// 从 host_metrics 获取最新监控数据
	metric, err := s.hostMetricDao.GetLatest(ctx, hostID)
	if err == nil && metric != nil {
		result["collected_at"] = metric.CollectedAt
		if metric.CPUUsagePercent != nil {
			result["cpu_usage"] = *metric.CPUUsagePercent
		}
		if metric.MemoryUsagePercent != nil {
			result["memory_usage"] = *metric.MemoryUsagePercent
		}
		if metric.DiskUsagePercent != nil {
			result["disk_usage"] = *metric.DiskUsagePercent
		}
	}

	// GPU 使用情况（从 GPU 列表构建）
	gpuUsage := make([]map[string]interface{}, 0, len(host.GPUs))
	for _, gpu := range host.GPUs {
		gpuUsage = append(gpuUsage, map[string]interface{}{
			"index":           gpu.Index,
			"name":            gpu.Name,
			"memory_total_mb": gpu.MemoryTotalMB,
			"status":          gpu.Status,
		})
	}
	result["gpu_usage"] = gpuUsage

	return result, nil
}

// GetRealtimeStatus 从 Redis 缓存读取设备实时状态，缓存未命中时回退到数据库
func (s *MachineService) GetRealtimeStatus(ctx context.Context, hostID string) (*CachedHostStatus, error) {
	if s.statusCache != nil {
		cached, err := s.statusCache.Get(ctx, hostID)
		if err == nil {
			return cached, nil
		}
	}

	// 回退到数据库
	host, err := s.machineDao.FindByID(ctx, hostID)
	if err != nil {
		return nil, err
	}

	status := &CachedHostStatus{
		HostID:       host.ID,
		DeviceStatus: host.DeviceStatus,
	}
	if host.LastHeartbeat != nil {
		status.LastHeartbeat = *host.LastHeartbeat
	}
	return status, nil
}
