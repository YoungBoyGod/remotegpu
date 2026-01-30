package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestResourceQuota_TableName 测试 ResourceQuota 表名
func TestResourceQuota_TableName(t *testing.T) {
	quota := ResourceQuota{}
	assert.Equal(t, "resource_quotas", quota.TableName())
}

// TestResourceQuota_Create 测试 ResourceQuota 结构体创建
func TestResourceQuota_Create(t *testing.T) {
	now := time.Now()
	workspaceID := uint(10)

	quota := ResourceQuota{
		ID:          1,
		CustomerID:  100,
		WorkspaceID: &workspaceID,
		CPU:         8,
		Memory:      16384,
		GPU:         2,
		Storage:     500,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, uint(1), quota.ID)
	assert.Equal(t, uint(100), quota.CustomerID)
	assert.NotNil(t, quota.WorkspaceID)
	assert.Equal(t, uint(10), *quota.WorkspaceID)
	assert.Equal(t, 8, quota.CPU)
	assert.Equal(t, int64(16384), quota.Memory)
	assert.Equal(t, 2, quota.GPU)
	assert.Equal(t, int64(500), quota.Storage)
	assert.Equal(t, now, quota.CreatedAt)
	assert.Equal(t, now, quota.UpdatedAt)
}

// TestResourceQuota_CreateWithNilWorkspace 测试创建用户级别配额（WorkspaceID为空）
func TestResourceQuota_CreateWithNilWorkspace(t *testing.T) {
	now := time.Now()

	quota := ResourceQuota{
		ID:          1,
		CustomerID:  100,
		WorkspaceID: nil,
		CPU:         16,
		Memory:      32768,
		GPU:         4,
		Storage:     1000,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, uint(1), quota.ID)
	assert.Equal(t, uint(100), quota.CustomerID)
	assert.Nil(t, quota.WorkspaceID)
	assert.Equal(t, 16, quota.CPU)
	assert.Equal(t, int64(32768), quota.Memory)
	assert.Equal(t, 4, quota.GPU)
	assert.Equal(t, int64(1000), quota.Storage)
}
