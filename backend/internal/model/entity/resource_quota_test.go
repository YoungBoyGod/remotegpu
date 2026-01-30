package entity

import (
	"reflect"
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

// TestResourceQuota_GORMTags 测试 GORM 标签
func TestResourceQuota_GORMTags(t *testing.T) {
	quota := ResourceQuota{}
	quotaType := reflect.TypeOf(quota)

	// 测试 ID 字段的 primarykey 标签
	idField, _ := quotaType.FieldByName("ID")
	assert.Contains(t, idField.Tag.Get("gorm"), "primarykey")

	// 测试 CustomerID 字段的标签
	customerIDField, _ := quotaType.FieldByName("CustomerID")
	assert.Contains(t, customerIDField.Tag.Get("gorm"), "not null")
	assert.Contains(t, customerIDField.Tag.Get("gorm"), "uniqueIndex:idx_customer_workspace")

	// 测试 WorkspaceID 字段的标签
	workspaceIDField, _ := quotaType.FieldByName("WorkspaceID")
	assert.Contains(t, workspaceIDField.Tag.Get("gorm"), "uniqueIndex:idx_customer_workspace")

	// 测试资源字段的标签
	cpuField, _ := quotaType.FieldByName("CPU")
	assert.Contains(t, cpuField.Tag.Get("gorm"), "not null")
	assert.Contains(t, cpuField.Tag.Get("gorm"), "default:0")

	memoryField, _ := quotaType.FieldByName("Memory")
	assert.Contains(t, memoryField.Tag.Get("gorm"), "not null")
	assert.Contains(t, memoryField.Tag.Get("gorm"), "default:0")

	gpuField, _ := quotaType.FieldByName("GPU")
	assert.Contains(t, gpuField.Tag.Get("gorm"), "not null")
	assert.Contains(t, gpuField.Tag.Get("gorm"), "default:0")

	storageField, _ := quotaType.FieldByName("Storage")
	assert.Contains(t, storageField.Tag.Get("gorm"), "not null")
	assert.Contains(t, storageField.Tag.Get("gorm"), "default:0")
}

// TestResourceQuota_JSONTags 测试 JSON 标签
func TestResourceQuota_JSONTags(t *testing.T) {
	quota := ResourceQuota{}
	quotaType := reflect.TypeOf(quota)

	// 测试所有字段的 JSON 标签
	idField, _ := quotaType.FieldByName("ID")
	assert.Equal(t, "id", idField.Tag.Get("json"))

	customerIDField, _ := quotaType.FieldByName("CustomerID")
	assert.Equal(t, "customer_id", customerIDField.Tag.Get("json"))

	workspaceIDField, _ := quotaType.FieldByName("WorkspaceID")
	assert.Equal(t, "workspace_id", workspaceIDField.Tag.Get("json"))

	cpuField, _ := quotaType.FieldByName("CPU")
	assert.Equal(t, "cpu", cpuField.Tag.Get("json"))

	memoryField, _ := quotaType.FieldByName("Memory")
	assert.Equal(t, "memory", memoryField.Tag.Get("json"))

	gpuField, _ := quotaType.FieldByName("GPU")
	assert.Equal(t, "gpu", gpuField.Tag.Get("json"))

	storageField, _ := quotaType.FieldByName("Storage")
	assert.Equal(t, "storage", storageField.Tag.Get("json"))
}

// TestResourceQuota_FieldTypes 测试字段类型
func TestResourceQuota_FieldTypes(t *testing.T) {
	quota := ResourceQuota{}
	quotaType := reflect.TypeOf(quota)

	// 测试 ID 类型
	idField, _ := quotaType.FieldByName("ID")
	assert.Equal(t, reflect.TypeOf(uint(0)), idField.Type)

	// 测试 CustomerID 类型
	customerIDField, _ := quotaType.FieldByName("CustomerID")
	assert.Equal(t, reflect.TypeOf(uint(0)), customerIDField.Type)

	// 测试 WorkspaceID 类型（指针）
	workspaceIDField, _ := quotaType.FieldByName("WorkspaceID")
	assert.Equal(t, reflect.TypeOf((*uint)(nil)), workspaceIDField.Type)

	// 测试 CPU 类型
	cpuField, _ := quotaType.FieldByName("CPU")
	assert.Equal(t, reflect.TypeOf(int(0)), cpuField.Type)

	// 测试 Memory 类型
	memoryField, _ := quotaType.FieldByName("Memory")
	assert.Equal(t, reflect.TypeOf(int64(0)), memoryField.Type)

	// 测试 GPU 类型
	gpuField, _ := quotaType.FieldByName("GPU")
	assert.Equal(t, reflect.TypeOf(int(0)), gpuField.Type)

	// 测试 Storage 类型
	storageField, _ := quotaType.FieldByName("Storage")
	assert.Equal(t, reflect.TypeOf(int64(0)), storageField.Type)
}

// TestResourceQuota_ZeroValues 测试零值
func TestResourceQuota_ZeroValues(t *testing.T) {
	quota := ResourceQuota{}

	assert.Equal(t, uint(0), quota.ID)
	assert.Equal(t, uint(0), quota.CustomerID)
	assert.Nil(t, quota.WorkspaceID)
	assert.Equal(t, 0, quota.CPU)
	assert.Equal(t, int64(0), quota.Memory)
	assert.Equal(t, 0, quota.GPU)
	assert.Equal(t, int64(0), quota.Storage)
	assert.True(t, quota.CreatedAt.IsZero())
	assert.True(t, quota.UpdatedAt.IsZero())
}

// TestResourceQuota_Associations 测试关联关系
func TestResourceQuota_Associations(t *testing.T) {
	quota := ResourceQuota{}
	quotaType := reflect.TypeOf(quota)

	// 测试 Customer 关联
	customerField, _ := quotaType.FieldByName("Customer")
	assert.Contains(t, customerField.Tag.Get("gorm"), "foreignKey:CustomerID")
	assert.Contains(t, customerField.Tag.Get("json"), "customer")
	assert.Contains(t, customerField.Tag.Get("json"), "omitempty")

	// 测试 Workspace 关联
	workspaceField, _ := quotaType.FieldByName("Workspace")
	assert.Contains(t, workspaceField.Tag.Get("gorm"), "foreignKey:WorkspaceID")
	assert.Contains(t, workspaceField.Tag.Get("json"), "workspace")
	assert.Contains(t, workspaceField.Tag.Get("json"), "omitempty")
}

// TestResourceQuota_WorkspaceIDPointer 测试 WorkspaceID 指针行为
func TestResourceQuota_WorkspaceIDPointer(t *testing.T) {
	// 测试 nil 指针
	quota1 := ResourceQuota{WorkspaceID: nil}
	assert.Nil(t, quota1.WorkspaceID)

	// 测试非 nil 指针
	workspaceID := uint(100)
	quota2 := ResourceQuota{WorkspaceID: &workspaceID}
	assert.NotNil(t, quota2.WorkspaceID)
	assert.Equal(t, uint(100), *quota2.WorkspaceID)

	// 测试修改指针值
	*quota2.WorkspaceID = 200
	assert.Equal(t, uint(200), *quota2.WorkspaceID)
}

// TestResourceQuota_BoundaryValues 测试边界值
func TestResourceQuota_BoundaryValues(t *testing.T) {
	// 测试最大值
	quota := ResourceQuota{
		CPU:     int(^uint(0) >> 1), // int 最大值
		Memory:  int64(^uint64(0) >> 1), // int64 最大值
		GPU:     int(^uint(0) >> 1),
		Storage: int64(^uint64(0) >> 1),
	}

	assert.Greater(t, quota.CPU, 0)
	assert.Greater(t, quota.Memory, int64(0))
	assert.Greater(t, quota.GPU, 0)
	assert.Greater(t, quota.Storage, int64(0))

	// 测试零值
	quota2 := ResourceQuota{
		CPU:     0,
		Memory:  0,
		GPU:     0,
		Storage: 0,
	}

	assert.Equal(t, 0, quota2.CPU)
	assert.Equal(t, int64(0), quota2.Memory)
	assert.Equal(t, 0, quota2.GPU)
	assert.Equal(t, int64(0), quota2.Storage)
}
