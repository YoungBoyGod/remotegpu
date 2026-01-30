package service

import (
	"fmt"
	"testing"

	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	"gorm.io/gorm"
)

// MockHostDao 模拟 HostDao
type MockHostDao struct {
	mock.Mock
}

func (m *MockHostDao) ListByStatus(status string) ([]*entity.Host, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Host), args.Error(1)
}

// MockEnvironmentDao 模拟 EnvironmentDao
type MockEnvironmentDao struct {
	mock.Mock
}

func (m *MockEnvironmentDao) Create(env *entity.Environment) error {
	args := m.Called(env)
	return args.Error(0)
}

func (m *MockEnvironmentDao) Update(env *entity.Environment) error {
	args := m.Called(env)
	return args.Error(0)
}

func (m *MockEnvironmentDao) GetByID(id string) (*entity.Environment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Environment), args.Error(1)
}

func (m *MockEnvironmentDao) GetByUserID(customerID uint) ([]*entity.Environment, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Environment), args.Error(1)
}

func (m *MockEnvironmentDao) GetByWorkspaceID(workspaceID uint) ([]*entity.Environment, error) {
	args := m.Called(workspaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Environment), args.Error(1)
}

// MockPortMappingDao 模拟 PortMappingDao
type MockPortMappingDao struct {
	mock.Mock
}

func (m *MockPortMappingDao) GetByEnvironmentID(envID string) ([]*entity.PortMapping, error) {
	args := m.Called(envID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.PortMapping), args.Error(1)
}

func (m *MockPortMappingDao) DeleteByEnvironmentID(envID string) error {
	args := m.Called(envID)
	return args.Error(0)
}

// MockK8sClient 模拟 K8s 客户端
type MockK8sClient struct {
	mock.Mock
}

func (m *MockK8sClient) CreatePod(config *k8s.PodConfig) (*corev1.Pod, error) {
	args := m.Called(config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corev1.Pod), args.Error(1)
}

func (m *MockK8sClient) DeletePod(namespace, name string) error {
	args := m.Called(namespace, name)
	return args.Error(0)
}

func (m *MockK8sClient) GetPodStatus(namespace, name string) (string, error) {
	args := m.Called(namespace, name)
	return args.String(0), args.Error(1)
}

func (m *MockK8sClient) GetPodLogs(namespace, name string, opts *k8s.LogOptions) (string, error) {
	args := m.Called(namespace, name, opts)
	return args.String(0), args.Error(1)
}

// MockResourceQuotaService 模拟配额服务
type MockResourceQuotaService struct {
	mock.Mock
}

func (m *MockResourceQuotaService) CheckQuota(customerID uint, workspaceID *uint, req *ResourceRequest) (bool, error) {
	args := m.Called(customerID, workspaceID, req)
	return args.Bool(0), args.Error(1)
}

// MockDB 模拟数据库
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Transaction(fc func(tx *gorm.DB) error) error {
	args := m.Called(fc)
	return args.Error(0)
}


// TestSelectHost_Success 测试主机选择成功
func TestSelectHost_Success(t *testing.T) {
	// 准备测试数据
	hosts := []*entity.Host{
		{
			ID:          "host1",
			TotalCPU:    16,
			UsedCPU:     4,
			TotalMemory: 32000,
			UsedMemory:  8000,
			TotalGPU:    4,
			UsedGPU:     1,
		},
		{
			ID:          "host2",
			TotalCPU:    16,
			UsedCPU:     8,
			TotalMemory: 32000,
			UsedMemory:  16000,
			TotalGPU:    4,
			UsedGPU:     2,
		},
	}

	// 创建 mock
	mockHostDao := new(MockHostDao)
	mockHostDao.On("ListByStatus", "active").Return(hosts, nil)

	// 创建 service
	service := &EnvironmentService{
		hostDao: mockHostDao,
	}

	// 执行测试
	host, err := service.selectHost(4, 8000, 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, host)
	assert.Equal(t, "host1", host.ID)
	mockHostDao.AssertExpectations(t)
}

// TestSelectHost_NoAvailableHost 测试没有可用主机
func TestSelectHost_NoAvailableHost(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockHostDao.On("ListByStatus", "active").Return([]*entity.Host{}, nil)

	service := &EnvironmentService{
		hostDao: mockHostDao,
	}

	host, err := service.selectHost(4, 8000, 1)

	assert.Error(t, err)
	assert.Nil(t, host)
	assert.Contains(t, err.Error(), "没有可用的主机")
	mockHostDao.AssertExpectations(t)
}

// TestSelectHost_InsufficientResources 测试资源不足
func TestSelectHost_InsufficientResources(t *testing.T) {
	hosts := []*entity.Host{
		{
			ID:          "host1",
			TotalCPU:    16,
			UsedCPU:     14,
			TotalMemory: 32000,
			UsedMemory:  30000,
			TotalGPU:    4,
			UsedGPU:     4,
		},
	}

	mockHostDao := new(MockHostDao)
	mockHostDao.On("ListByStatus", "active").Return(hosts, nil)

	service := &EnvironmentService{
		hostDao: mockHostDao,
	}

	host, err := service.selectHost(4, 8000, 1)

	assert.Error(t, err)
	assert.Nil(t, host)
	assert.Contains(t, err.Error(), "没有满足资源要求的主机")
	mockHostDao.AssertExpectations(t)
}

// TestSelectHost_LoadBalancing 测试负载均衡
func TestSelectHost_LoadBalancing(t *testing.T) {
	hosts := []*entity.Host{
		{
			ID:          "host1",
			TotalCPU:    16,
			UsedCPU:     12, // 75% 使用率
			TotalMemory: 32000,
			UsedMemory:  24000, // 75% 使用率
			TotalGPU:    4,
			UsedGPU:     3, // 75% 使用率
		},
		{
			ID:          "host2",
			TotalCPU:    16,
			UsedCPU:     4, // 25% 使用率
			TotalMemory: 32000,
			UsedMemory:  8000, // 25% 使用率
			TotalGPU:    4,
			UsedGPU:     1, // 25% 使用率
		},
	}

	mockHostDao := new(MockHostDao)
	mockHostDao.On("ListByStatus", "active").Return(hosts, nil)

	service := &EnvironmentService{
		hostDao: mockHostDao,
	}

	host, err := service.selectHost(2, 4000, 1)

	assert.NoError(t, err)
	assert.NotNil(t, host)
	// 应该选择使用率更低的 host2
	assert.Equal(t, "host2", host.ID)
	mockHostDao.AssertExpectations(t)
}

// TestSelectHost_DatabaseError 测试数据库错误
func TestSelectHost_DatabaseError(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockHostDao.On("ListByStatus", "active").Return(nil, fmt.Errorf("database connection failed"))

	service := &EnvironmentService{
		hostDao: mockHostDao,
	}

	host, err := service.selectHost(4, 8000, 1)

	assert.Error(t, err)
	assert.Nil(t, host)
	assert.Contains(t, err.Error(), "查询主机失败")
	mockHostDao.AssertExpectations(t)
}

// TestGetEnvironment_Success 测试获取环境成功
func TestGetEnvironment_Success(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	expectedEnv := &entity.Environment{
		ID:     "env-123",
		Name:   "test-env",
		Status: "running",
	}
	mockEnvDao.On("GetByID", "env-123").Return(expectedEnv, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	env, err := service.GetEnvironment("env-123")

	assert.NoError(t, err)
	assert.NotNil(t, env)
	assert.Equal(t, "env-123", env.ID)
	assert.Equal(t, "test-env", env.Name)
	mockEnvDao.AssertExpectations(t)
}

// TestGetEnvironment_NotFound 测试环境不存在
func TestGetEnvironment_NotFound(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockEnvDao.On("GetByID", "env-999").Return(nil, fmt.Errorf("record not found"))

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	env, err := service.GetEnvironment("env-999")

	assert.Error(t, err)
	assert.Nil(t, env)
	mockEnvDao.AssertExpectations(t)
}

// TestListEnvironments_ByCustomer 测试按客户列出环境
func TestListEnvironments_ByCustomer(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	expectedEnvs := []*entity.Environment{
		{ID: "env-1", Name: "env1", UserID: 1},
		{ID: "env-2", Name: "env2", UserID: 1},
	}
	mockEnvDao.On("GetByUserID", uint(1)).Return(expectedEnvs, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	envs, err := service.ListEnvironments(1, nil)

	assert.NoError(t, err)
	assert.Len(t, envs, 2)
	mockEnvDao.AssertExpectations(t)
}

// TestListEnvironments_ByWorkspace 测试按工作空间列出环境
func TestListEnvironments_ByWorkspace(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	workspaceID := uint(10)
	expectedEnvs := []*entity.Environment{
		{ID: "env-1", Name: "env1", WorkspaceID: &workspaceID},
	}
	mockEnvDao.On("GetByWorkspaceID", uint(10)).Return(expectedEnvs, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	envs, err := service.ListEnvironments(1, &workspaceID)

	assert.NoError(t, err)
	assert.Len(t, envs, 1)
	mockEnvDao.AssertExpectations(t)
}

// TestGetAccessInfo_Success 测试获取访问信息成功
func TestGetAccessInfo_Success(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockPortMappingDao := new(MockPortMappingDao)

	env := &entity.Environment{
		ID:      "env-123",
		Status:  "running",
		PodName: "pod-123",
	}
	portMappings := []*entity.PortMapping{
		{EnvID: "env-123", ExternalPort: 8080, InternalPort: 80},
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockPortMappingDao.On("GetByEnvironmentID", "env-123").Return(portMappings, nil)

	service := &EnvironmentService{
		envDao:         mockEnvDao,
		portMappingDao: mockPortMappingDao,
	}

	accessInfo, err := service.GetAccessInfo("env-123")

	assert.NoError(t, err)
	assert.NotNil(t, accessInfo)
	assert.Equal(t, "env-123", accessInfo["environment_id"])
	assert.Equal(t, "running", accessInfo["status"])
	mockEnvDao.AssertExpectations(t)
	mockPortMappingDao.AssertExpectations(t)
}

// TestGetAccessInfo_EnvironmentNotFound 测试环境不存在
func TestGetAccessInfo_EnvironmentNotFound(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockEnvDao.On("GetByID", "env-999").Return(nil, fmt.Errorf("record not found"))

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	accessInfo, err := service.GetAccessInfo("env-999")

	assert.Error(t, err)
	assert.Nil(t, accessInfo)
	assert.Contains(t, err.Error(), "获取环境失败")
	mockEnvDao.AssertExpectations(t)
}

// TestGetStatus_WithK8s 测试获取状态（有K8s Pod）
func TestGetStatus_WithK8s(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockK8sClient := new(MockK8sClient)

	env := &entity.Environment{
		ID:      "env-123",
		Status:  "running",
		PodName: "pod-123",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockK8sClient.On("GetPodStatus", "default", "pod-123").Return("Running", nil)

	service := &EnvironmentService{
		envDao:    mockEnvDao,
		k8sClient: mockK8sClient,
	}

	status, err := service.GetStatus("env-123")

	assert.NoError(t, err)
	assert.Equal(t, "Running", status)
	mockEnvDao.AssertExpectations(t)
	mockK8sClient.AssertExpectations(t)
}

// TestGetStatus_WithoutK8s 测试获取状态（无K8s Pod）
func TestGetStatus_WithoutK8s(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:      "env-123",
		Status:  "stopped",
		PodName: "",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	status, err := service.GetStatus("env-123")

	assert.NoError(t, err)
	assert.Equal(t, "stopped", status)
	mockEnvDao.AssertExpectations(t)
}

// TestCreateEnvironment_QuotaCheckFailed 测试配额检查失败
func TestCreateEnvironment_QuotaCheckFailed(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockQuotaService := new(MockResourceQuotaService)

	req := &CreateEnvironmentRequest{
		UserID: 1,
		Name:       "test-env",
		Image:      "ubuntu:20.04",
		CPU:        4,
		Memory:     8000,
		GPU:        1,
	}

	// Mock 配额检查失败
	mockQuotaService.On("CheckQuota", uint(1), (*uint)(nil), mock.AnythingOfType("*service.ResourceRequest")).Return(false, nil)

	service := &EnvironmentService{
		hostDao:      mockHostDao,
		quotaService: mockQuotaService,
	}

	env, err := service.CreateEnvironment(req)

	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "资源配额不足")
	mockQuotaService.AssertExpectations(t)
}

// TestCreateEnvironment_HostSelectionFailed 测试主机选择失败
func TestCreateEnvironment_HostSelectionFailed(t *testing.T) {
	mockHostDao := new(MockHostDao)
	mockQuotaService := new(MockResourceQuotaService)

	req := &CreateEnvironmentRequest{
		UserID: 1,
		Name:       "test-env",
		Image:      "ubuntu:20.04",
		CPU:        4,
		Memory:     8000,
		GPU:        1,
	}

	// Mock 配额检查通过
	mockQuotaService.On("CheckQuota", uint(1), (*uint)(nil), mock.AnythingOfType("*service.ResourceRequest")).Return(true, nil)
	// Mock 主机选择失败
	mockHostDao.On("ListByStatus", "active").Return([]*entity.Host{}, nil)

	service := &EnvironmentService{
		hostDao:      mockHostDao,
		quotaService: mockQuotaService,
	}

	env, err := service.CreateEnvironment(req)

	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), "选择主机失败")
	mockQuotaService.AssertExpectations(t)
	mockHostDao.AssertExpectations(t)
}

// TestValidateCreateRequest_Success 测试验证成功
func TestValidateCreateRequest_Success(t *testing.T) {
	service := &EnvironmentService{}

	req := &CreateEnvironmentRequest{
		Name:   "test-env",
		Image:  "ubuntu:20.04",
		CPU:    4,
		Memory: 8000,
		GPU:    1,
	}

	err := service.validateCreateRequest(req)
	assert.NoError(t, err)
}

// TestValidateCreateRequest_EmptyName 测试名称为空
func TestValidateCreateRequest_EmptyName(t *testing.T) {
	service := &EnvironmentService{}

	req := &CreateEnvironmentRequest{
		Name:   "",
		Image:  "ubuntu:20.04",
		CPU:    4,
		Memory: 8000,
		GPU:    1,
	}

	err := service.validateCreateRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "环境名称不能为空")
}

// TestValidateCreateRequest_EmptyImage 测试镜像为空
func TestValidateCreateRequest_EmptyImage(t *testing.T) {
	service := &EnvironmentService{}

	req := &CreateEnvironmentRequest{
		Name:   "test-env",
		Image:  "",
		CPU:    4,
		Memory: 8000,
		GPU:    1,
	}

	err := service.validateCreateRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "镜像不能为空")
}

// TestValidateCreateRequest_InvalidCPU 测试 CPU 无效
func TestValidateCreateRequest_InvalidCPU(t *testing.T) {
	service := &EnvironmentService{}

	req := &CreateEnvironmentRequest{
		Name:   "test-env",
		Image:  "ubuntu:20.04",
		CPU:    0,
		Memory: 8000,
		GPU:    1,
	}

	err := service.validateCreateRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CPU 必须大于 0")
}

// TestValidateCreateRequest_InvalidMemory 测试内存无效
func TestValidateCreateRequest_InvalidMemory(t *testing.T) {
	service := &EnvironmentService{}

	req := &CreateEnvironmentRequest{
		Name:   "test-env",
		Image:  "ubuntu:20.04",
		CPU:    4,
		Memory: 0,
		GPU:    1,
	}

	err := service.validateCreateRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "内存必须大于 0")
}

// TestValidateCreateRequest_NegativeGPU 测试 GPU 为负数
func TestValidateCreateRequest_NegativeGPU(t *testing.T) {
	service := &EnvironmentService{}

	req := &CreateEnvironmentRequest{
		Name:   "test-env",
		Image:  "ubuntu:20.04",
		CPU:    4,
		Memory: 8000,
		GPU:    -1,
	}

	err := service.validateCreateRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GPU 不能为负数")
}

// TestGetLogs_Success 测试获取日志成功
func TestGetLogs_Success(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockK8sClient := new(MockK8sClient)

	env := &entity.Environment{
		ID:      "env-123",
		PodName: "pod-123",
	}

	expectedLogs := "2024-01-01 10:00:00 Application started\n2024-01-01 10:00:01 Processing request"

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockK8sClient.On("GetPodLogs", "default", "pod-123", mock.AnythingOfType("*k8s.LogOptions")).Return(expectedLogs, nil)

	service := &EnvironmentService{
		envDao:    mockEnvDao,
		k8sClient: mockK8sClient,
	}

	logs, err := service.GetLogs("env-123", 100)

	assert.NoError(t, err)
	assert.Equal(t, expectedLogs, logs)
	mockEnvDao.AssertExpectations(t)
	mockK8sClient.AssertExpectations(t)
}

// TestGetLogs_EnvironmentNotFound 测试环境不存在
func TestGetLogs_EnvironmentNotFound(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockEnvDao.On("GetByID", "env-999").Return(nil, fmt.Errorf("record not found"))

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	logs, err := service.GetLogs("env-999", 100)

	assert.Error(t, err)
	assert.Empty(t, logs)
	assert.Contains(t, err.Error(), "获取环境失败")
	mockEnvDao.AssertExpectations(t)
}

// TestGetLogs_NoPodName 测试环境没有关联的 Pod
func TestGetLogs_NoPodName(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:      "env-123",
		PodName: "",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	logs, err := service.GetLogs("env-123", 100)

	assert.Error(t, err)
	assert.Empty(t, logs)
	assert.Contains(t, err.Error(), "环境没有关联的 Pod")
	mockEnvDao.AssertExpectations(t)
}

// TestGetLogs_NoK8sClient 测试 K8s 客户端未初始化
func TestGetLogs_NoK8sClient(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:      "env-123",
		PodName: "pod-123",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)

	service := &EnvironmentService{
		envDao:    mockEnvDao,
		k8sClient: nil,
	}

	logs, err := service.GetLogs("env-123", 100)

	assert.Error(t, err)
	assert.Empty(t, logs)
	assert.Contains(t, err.Error(), "K8s 客户端未初始化")
	mockEnvDao.AssertExpectations(t)
}

// TestStartEnvironment_Success 测试启动环境成功
func TestStartEnvironment_Success(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:     "env-123",
		Status: "stopped",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockEnvDao.On("Update", mock.AnythingOfType("*entity.Environment")).Return(nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.StartEnvironment("env-123")

	assert.NoError(t, err)
	mockEnvDao.AssertExpectations(t)
}

// TestStartEnvironment_WrongStatus 测试启动非停止状态的环境
func TestStartEnvironment_WrongStatus(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:     "env-123",
		Status: "running",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.StartEnvironment("env-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "只能启动已停止的环境")
	mockEnvDao.AssertExpectations(t)
}

// TestStartEnvironment_NotFound 测试启动不存在的环境
func TestStartEnvironment_NotFound(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockEnvDao.On("GetByID", "env-999").Return(nil, fmt.Errorf("record not found"))

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.StartEnvironment("env-999")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "获取环境失败")
	mockEnvDao.AssertExpectations(t)
}

// TestStopEnvironment_Success 测试停止环境成功
func TestStopEnvironment_Success(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:     "env-123",
		Status: "running",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockEnvDao.On("Update", mock.AnythingOfType("*entity.Environment")).Return(nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.StopEnvironment("env-123")

	assert.NoError(t, err)
	mockEnvDao.AssertExpectations(t)
}

// TestStopEnvironment_WrongStatus 测试停止非运行状态的环境
func TestStopEnvironment_WrongStatus(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:     "env-123",
		Status: "stopped",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.StopEnvironment("env-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "只能停止运行中的环境")
	mockEnvDao.AssertExpectations(t)
}

// TestStopEnvironment_NotFound 测试停止不存在的环境
func TestStopEnvironment_NotFound(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockEnvDao.On("GetByID", "env-999").Return(nil, fmt.Errorf("record not found"))

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.StopEnvironment("env-999")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "获取环境失败")
	mockEnvDao.AssertExpectations(t)
}

// TestRestartEnvironment_Success 测试重启环境成功
func TestRestartEnvironment_Success(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	// 第一次调用 GetByID (for StopEnvironment)
	runningEnv := &entity.Environment{
		ID:     "env-123",
		Status: "running",
	}
	mockEnvDao.On("GetByID", "env-123").Return(runningEnv, nil).Once()
	mockEnvDao.On("Update", mock.AnythingOfType("*entity.Environment")).Return(nil).Once()

	// 第二次调用 GetByID (for StartEnvironment)
	stoppedEnv := &entity.Environment{
		ID:     "env-123",
		Status: "stopped",
	}
	mockEnvDao.On("GetByID", "env-123").Return(stoppedEnv, nil).Once()
	mockEnvDao.On("Update", mock.AnythingOfType("*entity.Environment")).Return(nil).Once()

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.RestartEnvironment("env-123")

	assert.NoError(t, err)
	mockEnvDao.AssertExpectations(t)
}

// TestRestartEnvironment_StopFailed 测试重启环境时停止失败
func TestRestartEnvironment_StopFailed(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)

	env := &entity.Environment{
		ID:     "env-123",
		Status: "stopped",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	err := service.RestartEnvironment("env-123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "停止环境失败")
	mockEnvDao.AssertExpectations(t)
}

// TestGetStatus_GetByIDFailed 测试获取状态时查询环境失败
func TestGetStatus_GetByIDFailed(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockEnvDao.On("GetByID", "env-999").Return(nil, fmt.Errorf("database error"))

	service := &EnvironmentService{
		envDao: mockEnvDao,
	}

	status, err := service.GetStatus("env-999")

	assert.Error(t, err)
	assert.Empty(t, status)
	assert.Contains(t, err.Error(), "获取环境失败")
	mockEnvDao.AssertExpectations(t)
}

// TestGetStatus_K8sError 测试获取状态时 K8s 查询失败
func TestGetStatus_K8sError(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockK8sClient := new(MockK8sClient)

	env := &entity.Environment{
		ID:      "env-123",
		Status:  "running",
		PodName: "pod-123",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockK8sClient.On("GetPodStatus", "default", "pod-123").Return("", fmt.Errorf("k8s error"))

	service := &EnvironmentService{
		envDao:    mockEnvDao,
		k8sClient: mockK8sClient,
	}

	status, err := service.GetStatus("env-123")

	assert.NoError(t, err)
	assert.Equal(t, "running", status)
	mockEnvDao.AssertExpectations(t)
	mockK8sClient.AssertExpectations(t)
}

// TestGetAccessInfo_PortMappingError 测试获取访问信息时端口映射查询失败
func TestGetAccessInfo_PortMappingError(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockPortMappingDao := new(MockPortMappingDao)

	env := &entity.Environment{
		ID:      "env-123",
		Status:  "running",
		PodName: "pod-123",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockPortMappingDao.On("GetByEnvironmentID", "env-123").Return(nil, fmt.Errorf("database error"))

	service := &EnvironmentService{
		envDao:         mockEnvDao,
		portMappingDao: mockPortMappingDao,
	}

	accessInfo, err := service.GetAccessInfo("env-123")

	assert.Error(t, err)
	assert.Nil(t, accessInfo)
	assert.Contains(t, err.Error(), "获取端口映射失败")
	mockEnvDao.AssertExpectations(t)
	mockPortMappingDao.AssertExpectations(t)
}

// TestGetLogs_K8sError 测试获取日志时 K8s 查询失败
func TestGetLogs_K8sError(t *testing.T) {
	mockEnvDao := new(MockEnvironmentDao)
	mockK8sClient := new(MockK8sClient)

	env := &entity.Environment{
		ID:      "env-123",
		PodName: "pod-123",
	}

	mockEnvDao.On("GetByID", "env-123").Return(env, nil)
	mockK8sClient.On("GetPodLogs", "default", "pod-123", mock.AnythingOfType("*k8s.LogOptions")).Return("", fmt.Errorf("k8s error"))

	service := &EnvironmentService{
		envDao:    mockEnvDao,
		k8sClient: mockK8sClient,
	}

	logs, err := service.GetLogs("env-123", 100)

	assert.Error(t, err)
	assert.Empty(t, logs)
	mockEnvDao.AssertExpectations(t)
	mockK8sClient.AssertExpectations(t)
}
