/**
 * 工作空间模块 - API 接口
 */
import request from '../common/request'
import type {
  WorkspaceInfo,
  CreateWorkspaceRequest,
  UpdateWorkspaceRequest,
  WorkspaceListResponse,
  AddMemberRequest,
  WorkspaceMemberInfo
} from './types'
import type { StatusResponse } from '../common/types'

// ==================== 工作空间管理 ====================

/**
 * 创建工作空间
 */
export function createWorkspace(data: CreateWorkspaceRequest) {
  return request.post<WorkspaceInfo>('/workspaces', data)
}

/**
 * 获取工作空间列表
 */
export function getWorkspaces(page: number = 1, pageSize: number = 10) {
  return request.get<WorkspaceListResponse>('/workspaces', {
    params: { page, page_size: pageSize }
  })
}

/**
 * 获取工作空间详情
 */
export function getWorkspaceById(id: number) {
  return request.get<WorkspaceInfo>(`/workspaces/${id}`)
}

/**
 * 更新工作空间
 */
export function updateWorkspace(id: number, data: UpdateWorkspaceRequest) {
  return request.put<WorkspaceInfo>(`/workspaces/${id}`, data)
}

/**
 * 删除工作空间
 */
export function deleteWorkspace(id: number) {
  return request.delete<StatusResponse>(`/workspaces/${id}`)
}

// ==================== 成员管理 ====================

/**
 * 添加成员
 */
export function addMember(workspaceId: number, data: AddMemberRequest) {
  return request.post<StatusResponse>(`/workspaces/${workspaceId}/members`, data)
}

/**
 * 移除成员
 */
export function removeMember(workspaceId: number, userId: number) {
  return request.delete<StatusResponse>(`/workspaces/${workspaceId}/members/${userId}`)
}

/**
 * 获取成员列表
 */
export function getMembers(workspaceId: number) {
  return request.get<WorkspaceMemberInfo[]>(`/workspaces/${workspaceId}/members`)
}
