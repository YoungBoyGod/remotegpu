package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestWorkspace_TableName 测试 Workspace 表名
func TestWorkspace_TableName(t *testing.T) {
	workspace := Workspace{}
	assert.Equal(t, "workspaces", workspace.TableName())
}

// TestWorkspace_Create 测试 Workspace 结构体创建
func TestWorkspace_Create(t *testing.T) {
	now := time.Now()
	testUUID := uuid.New()

	workspace := Workspace{
		ID:          1,
		UUID:        testUUID,
		OwnerID:     100,
		Name:        "Test Workspace",
		Description: "This is a test workspace",
		Type:        "personal",
		MemberCount: 1,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, uint(1), workspace.ID)
	assert.Equal(t, testUUID, workspace.UUID)
	assert.Equal(t, uint(100), workspace.OwnerID)
	assert.Equal(t, "Test Workspace", workspace.Name)
	assert.Equal(t, "This is a test workspace", workspace.Description)
	assert.Equal(t, "personal", workspace.Type)
	assert.Equal(t, 1, workspace.MemberCount)
	assert.Equal(t, "active", workspace.Status)
	assert.Equal(t, now, workspace.CreatedAt)
	assert.Equal(t, now, workspace.UpdatedAt)
}

// TestWorkspaceMember_TableName 测试 WorkspaceMember 表名
func TestWorkspaceMember_TableName(t *testing.T) {
	member := WorkspaceMember{}
	assert.Equal(t, "workspace_members", member.TableName())
}

// TestWorkspaceMember_Create 测试 WorkspaceMember 结构体创建
func TestWorkspaceMember_Create(t *testing.T) {
	now := time.Now()

	member := WorkspaceMember{
		ID:          1,
		WorkspaceID: 10,
		CustomerID:  100,
		Role:        "member",
		Status:      "active",
		JoinedAt:    now,
		CreatedAt:   now,
	}

	assert.Equal(t, uint(1), member.ID)
	assert.Equal(t, uint(10), member.WorkspaceID)
	assert.Equal(t, uint(100), member.CustomerID)
	assert.Equal(t, "member", member.Role)
	assert.Equal(t, "active", member.Status)
	assert.Equal(t, now, member.JoinedAt)
	assert.Equal(t, now, member.CreatedAt)
}
