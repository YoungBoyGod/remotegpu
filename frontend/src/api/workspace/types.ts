/**
 * 工作空间模块 - 类型定义
 */

// 工作空间角色
export type WorkspaceRole = 'owner' | 'admin' | 'member'

// 工作空间信息
export interface WorkspaceInfo {
  id: number
  name: string
  description: string
  owner_id: number
  member_count: number
  created_at: string
  updated_at: string
}

// 创建工作空间请求
export interface CreateWorkspaceRequest {
  name: string
  description?: string
}

// 更新工作空间请求
export interface UpdateWorkspaceRequest {
  name?: string
  description?: string
}

// 工作空间列表响应
export interface WorkspaceListResponse {
  items: WorkspaceInfo[]
  total: number
  page: number
  page_size: number
}

// 添加成员请求
export interface AddMemberRequest {
  user_id: number
  role: WorkspaceRole
}

// 工作空间成员信息
export interface WorkspaceMemberInfo {
  id: number
  workspace_id: number
  customer_id: number
  username: string
  email: string
  role: WorkspaceRole
  joined_at: string
}
