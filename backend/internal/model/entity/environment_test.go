package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEnvironment_TableName 测试 Environment 表名
func TestEnvironment_TableName(t *testing.T) {
	env := Environment{}
	assert.Equal(t, "environments", env.TableName())
}

// TestEnvironment_Create 测试 Environment 结构体创建
func TestEnvironment_Create(t *testing.T) {
	now := time.Now()
	workspaceID := uint(10)
	storage := int64(10737418240) // 10GB
	sshPort := 22
	rdpPort := 3389
	jupyterPort := 8888

	env := Environment{
		ID:          "env-test-001",
		UserID:  100,
		WorkspaceID: &workspaceID,
		HostID:      "host-001",
		Name:        "Test Environment",
		Description: "This is a test environment",
		Image:       "ubuntu:22.04",
		Status:      "running",
		CPU:         4,
		Memory:      8589934592, // 8GB
		GPU:         1,
		Storage:     &storage,
		SSHPort:     &sshPort,
		RDPPort:     &rdpPort,
		JupyterPort: &jupyterPort,
		ContainerID: "container-123",
		PodName:     "pod-test-001",
		CreatedAt:   now,
		UpdatedAt:   now,
		StartedAt:   &now,
	}

	assert.Equal(t, "env-test-001", env.ID)
	assert.Equal(t, uint(100), env.UserID)
	assert.Equal(t, uint(10), *env.WorkspaceID)
	assert.Equal(t, "host-001", env.HostID)
	assert.Equal(t, "Test Environment", env.Name)
	assert.Equal(t, "This is a test environment", env.Description)
	assert.Equal(t, "ubuntu:22.04", env.Image)
	assert.Equal(t, "running", env.Status)
	assert.Equal(t, 4, env.CPU)
	assert.Equal(t, int64(8589934592), env.Memory)
	assert.Equal(t, 1, env.GPU)
	assert.Equal(t, int64(10737418240), *env.Storage)
	assert.Equal(t, 22, *env.SSHPort)
	assert.Equal(t, 3389, *env.RDPPort)
	assert.Equal(t, 8888, *env.JupyterPort)
	assert.Equal(t, "container-123", env.ContainerID)
	assert.Equal(t, "pod-test-001", env.PodName)
	assert.Equal(t, now, env.CreatedAt)
	assert.Equal(t, now, env.UpdatedAt)
	assert.Equal(t, now, *env.StartedAt)
}

// TestPortMapping_TableName 测试 PortMapping 表名
func TestPortMapping_TableName(t *testing.T) {
	pm := PortMapping{}
	assert.Equal(t, "port_mappings", pm.TableName())
}

// TestPortMapping_Create 测试 PortMapping 结构体创建
func TestPortMapping_Create(t *testing.T) {
	now := time.Now()

	pm := PortMapping{
		ID:           1,
		EnvID:        "env-test-001",
		ServiceType:  "ssh",
		ExternalPort: 30022,
		InternalPort: 22,
		Status:       "active",
		AllocatedAt:  now,
	}

	assert.Equal(t, int64(1), pm.ID)
	assert.Equal(t, "env-test-001", pm.EnvID)
	assert.Equal(t, "ssh", pm.ServiceType)
	assert.Equal(t, 30022, pm.ExternalPort)
	assert.Equal(t, 22, pm.InternalPort)
	assert.Equal(t, "active", pm.Status)
	assert.Equal(t, now, pm.AllocatedAt)
	assert.Nil(t, pm.ReleasedAt)
}

// TestEnvironment_PointerFields 测试 Environment 指针字段
func TestEnvironment_PointerFields(t *testing.T) {
	// 测试所有指针字段为 nil 的情况
	env := Environment{
		ID:         "env-test-002",
		UserID: 100,
		HostID:     "host-001",
		Name:       "Test Environment",
		Image:      "ubuntu:22.04",
		CPU:        2,
		Memory:     4294967296,
	}

	assert.Nil(t, env.WorkspaceID)
	assert.Nil(t, env.Storage)
	assert.Nil(t, env.SSHPort)
	assert.Nil(t, env.RDPPort)
	assert.Nil(t, env.JupyterPort)
	assert.Nil(t, env.StartedAt)
	assert.Nil(t, env.StoppedAt)

	// 测试指针字段有值的情况
	workspaceID := uint(10)
	storage := int64(10737418240)
	sshPort := 22
	rdpPort := 3389
	jupyterPort := 8888
	now := time.Now()

	env.WorkspaceID = &workspaceID
	env.Storage = &storage
	env.SSHPort = &sshPort
	env.RDPPort = &rdpPort
	env.JupyterPort = &jupyterPort
	env.StartedAt = &now
	env.StoppedAt = &now

	assert.NotNil(t, env.WorkspaceID)
	assert.Equal(t, uint(10), *env.WorkspaceID)
	assert.NotNil(t, env.Storage)
	assert.Equal(t, int64(10737418240), *env.Storage)
	assert.NotNil(t, env.SSHPort)
	assert.Equal(t, 22, *env.SSHPort)
	assert.NotNil(t, env.RDPPort)
	assert.Equal(t, 3389, *env.RDPPort)
	assert.NotNil(t, env.JupyterPort)
	assert.Equal(t, 8888, *env.JupyterPort)
	assert.NotNil(t, env.StartedAt)
	assert.NotNil(t, env.StoppedAt)
}

// TestEnvironment_StatusValues 测试 Environment 状态枚举值
func TestEnvironment_StatusValues(t *testing.T) {
	validStatuses := []string{"creating", "running", "stopped", "error", "deleting"}

	for _, status := range validStatuses {
		env := Environment{
			ID:         "env-test-" + status,
			UserID: 100,
			HostID:     "host-001",
			Name:       "Test Environment",
			Image:      "ubuntu:22.04",
			Status:     status,
			CPU:        2,
			Memory:     4294967296,
		}

		assert.Equal(t, status, env.Status)
	}
}

// TestPortMapping_StatusValues 测试 PortMapping 状态枚举值
func TestPortMapping_StatusValues(t *testing.T) {
	validStatuses := []string{"active", "released"}

	for _, status := range validStatuses {
		pm := PortMapping{
			EnvID:        "env-test-001",
			ServiceType:  "ssh",
			ExternalPort: 30022,
			InternalPort: 22,
			Status:       status,
			AllocatedAt:  time.Now(),
		}

		assert.Equal(t, status, pm.Status)
	}
}

// TestPortMapping_ServiceTypes 测试 PortMapping 服务类型
func TestPortMapping_ServiceTypes(t *testing.T) {
	validServiceTypes := []string{"ssh", "rdp", "jupyter", "custom"}

	for _, serviceType := range validServiceTypes {
		pm := PortMapping{
			EnvID:        "env-test-001",
			ServiceType:  serviceType,
			ExternalPort: 30022,
			InternalPort: 22,
			Status:       "active",
			AllocatedAt:  time.Now(),
		}

		assert.Equal(t, serviceType, pm.ServiceType)
	}
}

// TestPortMapping_ReleasedAt 测试 PortMapping ReleasedAt 字段
func TestPortMapping_ReleasedAt(t *testing.T) {
	now := time.Now()

	// 测试未释放的端口映射
	pm1 := PortMapping{
		EnvID:        "env-test-001",
		ServiceType:  "ssh",
		ExternalPort: 30022,
		InternalPort: 22,
		Status:       "active",
		AllocatedAt:  now,
	}
	assert.Nil(t, pm1.ReleasedAt)

	// 测试已释放的端口映射
	releasedTime := now.Add(1 * time.Hour)
	pm2 := PortMapping{
		EnvID:        "env-test-001",
		ServiceType:  "ssh",
		ExternalPort: 30022,
		InternalPort: 22,
		Status:       "released",
		AllocatedAt:  now,
		ReleasedAt:   &releasedTime,
	}
	assert.NotNil(t, pm2.ReleasedAt)
	assert.Equal(t, releasedTime, *pm2.ReleasedAt)
}

// TestEnvironment_MinimalFields 测试 Environment 最小必需字段
func TestEnvironment_MinimalFields(t *testing.T) {
	env := Environment{
		ID:         "env-minimal",
		UserID: 1,
		HostID:     "host-001",
		Name:       "Minimal Env",
		Image:      "ubuntu:22.04",
		CPU:        1,
		Memory:     1073741824, // 1GB
	}

	assert.Equal(t, "env-minimal", env.ID)
	assert.Equal(t, uint(1), env.UserID)
	assert.Equal(t, "host-001", env.HostID)
	assert.Equal(t, "Minimal Env", env.Name)
	assert.Equal(t, "ubuntu:22.04", env.Image)
	assert.Equal(t, 1, env.CPU)
	assert.Equal(t, int64(1073741824), env.Memory)
	assert.Equal(t, 0, env.GPU) // 默认值
}

// TestPortMapping_MinimalFields 测试 PortMapping 最小必需字段
func TestPortMapping_MinimalFields(t *testing.T) {
	pm := PortMapping{
		EnvID:        "env-test-001",
		ServiceType:  "ssh",
		ExternalPort: 30022,
		InternalPort: 22,
	}

	assert.Equal(t, "env-test-001", pm.EnvID)
	assert.Equal(t, "ssh", pm.ServiceType)
	assert.Equal(t, 30022, pm.ExternalPort)
	assert.Equal(t, 22, pm.InternalPort)
}
