package machine

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/cache"
	"github.com/YoungBoyGod/remotegpu/pkg/crypto"
	"github.com/YoungBoyGod/remotegpu/pkg/logger"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

var (
	ErrEnrollmentAuthRequired = errors.New("ssh credential required")
	ErrEnrollmentInvalidAddr  = errors.New("host address required")
)

// SystemInfoSnapshot Agent 返回的系统信息快照
type SystemInfoSnapshot struct {
	Hostname      string
	CPUCores      int
	MemoryTotalGB int64
	DiskTotalGB   int64
	Collected     bool
}

// AgentSystemInfoProvider Agent 系统信息获取接口（避免循环依赖）
type AgentSystemInfoProvider interface {
	GetSystemInfo(ctx context.Context, hostID, address string) (*SystemInfoSnapshot, error)
}

type MachineEnrollmentService struct {
	db             *gorm.DB
	enrollmentDao  *dao.MachineEnrollmentDao
	machineService *MachineService
	agentProvider  AgentSystemInfoProvider
	redisClient    *redis.Client
	maxRetries     int
	retryDelay     time.Duration
	skipCollect    bool
}

func NewMachineEnrollmentService(db *gorm.DB, machineSvc *MachineService, agentProvider AgentSystemInfoProvider) *MachineEnrollmentService {
	maxRetries := 3
	retryDelay := 10 * time.Second
	if config.GlobalConfig != nil {
		if config.GlobalConfig.Enrollment.MaxRetries >= 0 {
			maxRetries = config.GlobalConfig.Enrollment.MaxRetries
		}
		if config.GlobalConfig.Enrollment.RetryDelay > 0 {
			retryDelay = time.Duration(config.GlobalConfig.Enrollment.RetryDelay) * time.Second
		}
	}
	return &MachineEnrollmentService{
		db:             db,
		enrollmentDao:  dao.NewMachineEnrollmentDao(db),
		machineService: machineSvc,
		agentProvider:  agentProvider,
		redisClient:    cache.GetRedis(),
		maxRetries:     maxRetries,
		retryDelay:     retryDelay,
		skipCollect:    config.GlobalConfig != nil && config.GlobalConfig.Enrollment.SkipCollect,
	}
}

const machineEnrollmentQueueKey = "machine:enrollment:queue"
const machineEnrollmentRetryKey = "machine:enrollment:retry"

type EnrollmentSpec struct {
	Hostname  string
	CPUCores  int
	MemoryGB  int64
	DiskGB    int64
	GPUCount  int
	Collected bool
}

func (s *MachineEnrollmentService) CreateEnrollment(ctx context.Context, customerID uint, req *entity.MachineEnrollment) (*entity.MachineEnrollment, error) {
	if req.Address == "" {
		return nil, ErrEnrollmentInvalidAddr
	}
	if req.SSHPassword == "" && req.SSHKey == "" {
		return nil, ErrEnrollmentAuthRequired
	}

	if req.Name == "" {
		if req.Hostname != "" {
			req.Name = req.Hostname
		} else {
			req.Name = req.Address
		}
	}

	if req.Region == "" {
		req.Region = "default"
	}

	req.CustomerID = customerID
	req.Status = "pending"

	if err := s.enrollmentDao.Create(ctx, req); err != nil {
		return nil, err
	}

	enqueueCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.enqueue(enqueueCtx, req.ID); err != nil {
		_ = s.enrollmentDao.UpdateStatus(ctx, req.ID, "failed", err.Error(), "")
		return nil, err
	}
	return req, nil
}

func (s *MachineEnrollmentService) ListByCustomer(ctx context.Context, customerID uint, page, pageSize int) ([]entity.MachineEnrollment, int64, error) {
	return s.enrollmentDao.ListByCustomer(ctx, customerID, page, pageSize)
}

func (s *MachineEnrollmentService) GetByID(ctx context.Context, id uint) (*entity.MachineEnrollment, error) {
	return s.enrollmentDao.FindByID(ctx, id)
}

func (s *MachineEnrollmentService) StartWorker(ctx context.Context) {
	if s.redisClient == nil {
		logger.GetLogger().Warn("Machine enrollment queue disabled: redis client not initialized")
		return
	}
	go s.runWorker(ctx)
}

func (s *MachineEnrollmentService) enqueue(ctx context.Context, id uint) error {
	if s.redisClient == nil {
		return errors.New("redis client not initialized")
	}
	return s.redisClient.RPush(ctx, machineEnrollmentQueueKey, strconv.FormatUint(uint64(id), 10)).Err()
}

func (s *MachineEnrollmentService) runWorker(ctx context.Context) {
	s.requeuePending(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result, err := s.redisClient.BLPop(ctx, 5*time.Second, machineEnrollmentQueueKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			logger.GetLogger().Warn(fmt.Sprintf("Machine enrollment queue error: %v", err))
			continue
		}
		if len(result) < 2 {
			continue
		}

		id, err := strconv.ParseUint(result[1], 10, 64)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Invalid enrollment id: %s", result[1]))
			continue
		}
		s.processEnrollment(uint(id))
	}
}

func (s *MachineEnrollmentService) requeuePending(ctx context.Context) {
	list, err := s.enrollmentDao.ListPending(ctx, 200)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("Failed to requeue enrollments: %v", err))
		return
	}
	for _, enrollment := range list {
		if err := s.redisClient.RPush(ctx, machineEnrollmentQueueKey, strconv.FormatUint(uint64(enrollment.ID), 10)).Err(); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to enqueue enrollment %d: %v", enrollment.ID, err))
		}
	}
}

func (s *MachineEnrollmentService) handleEnrollmentFailure(ctx context.Context, enrollmentID uint, err error) {
	if s.redisClient == nil {
		_ = s.enrollmentDao.UpdateStatus(ctx, enrollmentID, "failed", err.Error(), "")
		return
	}
	if s.maxRetries <= 0 {
		_ = s.redisClient.HDel(ctx, machineEnrollmentRetryKey, strconv.FormatUint(uint64(enrollmentID), 10)).Err()
		_ = s.enrollmentDao.UpdateStatus(ctx, enrollmentID, "failed", err.Error(), "")
		return
	}

	retryKey := strconv.FormatUint(uint64(enrollmentID), 10)
	retryCount, retryErr := s.redisClient.HIncrBy(ctx, machineEnrollmentRetryKey, retryKey, 1).Result()
	if retryErr != nil {
		_ = s.enrollmentDao.UpdateStatus(ctx, enrollmentID, "failed", err.Error(), "")
		return
	}

	if retryCount > int64(s.maxRetries) {
		_ = s.redisClient.HDel(ctx, machineEnrollmentRetryKey, retryKey).Err()
		_ = s.enrollmentDao.UpdateStatus(ctx, enrollmentID, "failed", err.Error(), "")
		return
	}

	_ = s.enrollmentDao.UpdateStatus(ctx, enrollmentID, "pending", fmt.Sprintf("retry %d/%d: %v", retryCount, s.maxRetries, err), "")
	s.scheduleRetry(enrollmentID)
}

func (s *MachineEnrollmentService) scheduleRetry(enrollmentID uint) {
	go func() {
		<-time.After(s.retryDelay)
		retryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.enqueue(retryCtx, enrollmentID); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to requeue enrollment %d: %v", enrollmentID, err))
		}
	}()
}

func (s *MachineEnrollmentService) clearRetry(ctx context.Context, enrollmentID uint) {
	if s.redisClient == nil {
		return
	}
	_ = s.redisClient.HDel(ctx, machineEnrollmentRetryKey, strconv.FormatUint(uint64(enrollmentID), 10)).Err()
}

func (s *MachineEnrollmentService) processEnrollment(id uint) {
	ctx := context.Background()
	enrollment, err := s.enrollmentDao.FindByID(ctx, id)
	if err != nil {
		return
	}

	if enrollment.Status != "pending" {
		return
	}

	if s.skipCollect {
		host := entity.Host{
			Name:         enrollment.Name,
			Hostname:     enrollment.Hostname,
			Region:       enrollment.Region,
			IPAddress:    enrollment.Address,
			SSHPort:      enrollment.SSHPort,
			SSHUsername:  enrollment.SSHUsername,
			SSHPassword:  enrollment.SSHPassword,
			SSHKey:       enrollment.SSHKey,
			Status:       "offline",
			HealthStatus: "unknown",
			NeedsCollect: true,
		}
		if host.Hostname == "" {
			host.Hostname = enrollment.Address
		}
		if host.Name == "" {
			host.Name = host.Hostname
		}

		if err := s.machineService.CreateMachine(ctx, &host); err != nil {
			_ = s.enrollmentDao.UpdateStatus(ctx, enrollment.ID, "failed", err.Error(), "")
			return
		}

		_ = s.enrollmentDao.UpdateStatus(ctx, enrollment.ID, "success", "", host.ID)
		s.clearRetry(ctx, enrollment.ID)
		return
	}

	spec, err := s.collectSpec(ctx, enrollment)
	if err != nil {
		if errors.Is(err, ErrEnrollmentAuthRequired) || errors.Is(err, ErrEnrollmentInvalidAddr) {
			_ = s.enrollmentDao.UpdateStatus(ctx, enrollment.ID, "failed", err.Error(), "")
			return
		}
		s.handleEnrollmentFailure(ctx, enrollment.ID, err)
		return
	}
	// CodeX 2026-02-04: ensure spec/hostname valid before creating host.
	if spec.Hostname == "" {
		spec.Hostname = enrollment.Hostname
		if spec.Hostname == "" {
			spec.Hostname = enrollment.Address
		}
	}
	if spec.CPUCores <= 0 || spec.MemoryGB <= 0 || spec.DiskGB <= 0 {
		_ = s.enrollmentDao.UpdateStatus(ctx, enrollment.ID, "failed", "hardware spec collection failed", "")
		return
	}

	host := entity.Host{
		Name:          enrollment.Name,
		Hostname:      spec.Hostname,
		Region:        enrollment.Region,
		IPAddress:     enrollment.Address,
		SSHPort:       enrollment.SSHPort,
		SSHUsername:   enrollment.SSHUsername,
		SSHPassword:   enrollment.SSHPassword,
		SSHKey:        enrollment.SSHKey,
		Status:        "idle",
		TotalCPU:      spec.CPUCores,
		TotalMemoryGB: spec.MemoryGB,
		TotalDiskGB:   spec.DiskGB,
		HealthStatus:  "healthy",
		NeedsCollect:  false,
	}

	if host.Name == "" {
		host.Name = host.Hostname
	}

	if err := s.machineService.CreateMachine(ctx, &host); err != nil {
		if err == ErrHostDuplicateIP || err == ErrHostDuplicateHostname || errors.Is(err, gorm.ErrDuplicatedKey) {
			errorMessage := err.Error()
			_ = s.enrollmentDao.UpdateStatus(ctx, enrollment.ID, "failed", errorMessage, "")
			return
		}
		s.handleEnrollmentFailure(ctx, enrollment.ID, err)
		return
	}

	_ = s.enrollmentDao.UpdateStatus(ctx, enrollment.ID, "success", "", host.ID)
	s.clearRetry(ctx, enrollment.ID)
}

func (s *MachineEnrollmentService) collectSpec(ctx context.Context, enrollment *entity.MachineEnrollment) (*EnrollmentSpec, error) {
	if enrollment.Address == "" {
		return nil, ErrEnrollmentInvalidAddr
	}
	if enrollment.SSHPassword == "" && enrollment.SSHKey == "" {
		return nil, ErrEnrollmentAuthRequired
	}

	if s.agentProvider != nil {
		if info, err := s.agentProvider.GetSystemInfo(ctx, enrollment.Address, enrollment.Address); err == nil {
			spec := &EnrollmentSpec{
				Hostname:  info.Hostname,
				CPUCores:  info.CPUCores,
				MemoryGB:  info.MemoryTotalGB,
				DiskGB:    info.DiskTotalGB,
				GPUCount:  0,
				Collected: info.Collected,
			}
			// CodeX 2026-02-04: only trust agent spec when data is complete.
			if spec.Collected && spec.CPUCores > 0 && spec.MemoryGB > 0 && spec.DiskGB > 0 {
				return spec, nil
			}
		}
	}

	client, err := s.connectSSH(enrollment)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	hostname, err := s.runCommand(client, "hostname")
	if err != nil {
		return nil, err
	}
	cores, err := s.runCommand(client, "nproc")
	if err != nil {
		return nil, err
	}
	mem, err := s.runCommand(client, "grep MemTotal /proc/meminfo")
	if err != nil {
		return nil, err
	}
	disk, err := s.runCommand(client, "df -k / | tail -1")
	if err != nil {
		return nil, err
	}

	return &EnrollmentSpec{
		Hostname:  strings.TrimSpace(hostname),
		CPUCores:  parseInt(cores),
		MemoryGB:  parseMemGB(mem),
		DiskGB:    parseDiskGB(disk),
		GPUCount:  0,
		Collected: true,
	}, nil
}

func (s *MachineEnrollmentService) connectSSH(enrollment *entity.MachineEnrollment) (*ssh.Client, error) {
	authMethods := []ssh.AuthMethod{}
	if enrollment.SSHPassword != "" {
		// 修复 P0 安全问题：解密 SSH 密码
		decryptedPassword, err := crypto.DecryptAES256GCM(enrollment.SSHPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt SSH password: %w", err)
		}
		authMethods = append(authMethods, ssh.Password(decryptedPassword))
	}
	if enrollment.SSHKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(enrollment.SSHKey))
		if err != nil {
			return nil, fmt.Errorf("parse ssh key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User:            enrollment.SSHUsername,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(enrollment.Address, strconv.Itoa(enrollment.SSHPort))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh dial failed: %w", err)
	}
	return client, nil
}

func (s *MachineEnrollmentService) runCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("command %s failed: %w", cmd, err)
	}
	return string(output), nil
}

func parseInt(value string) int {
	trimmed := strings.TrimSpace(value)
	number, _ := strconv.Atoi(trimmed)
	return number
}

func parseMemGB(value string) int64 {
	fields := strings.Fields(value)
	if len(fields) < 2 {
		return 0
	}
	kb, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0
	}
	return kb / 1024 / 1024
}

func parseDiskGB(value string) int64 {
	fields := strings.Fields(value)
	if len(fields) < 2 {
		return 0
	}
	kb, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0
	}
	return kb / 1024 / 1024
}
