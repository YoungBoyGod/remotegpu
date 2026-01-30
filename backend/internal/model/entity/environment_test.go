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
		CustomerID:  100,
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
	assert.Equal(t, uint(100), env.CustomerID)
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
