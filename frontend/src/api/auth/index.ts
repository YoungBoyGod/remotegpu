/**
 * 用户与权限模块 - API 接口
 */
import request from '../common/request'
import type {
  User,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  Workspace,
  WorkspaceMember,
  ResourceQuota,
  QuotaUsage
} from './types'
import type { IdResponse, StatusResponse, PaginationParams, PaginationResponse } from '../common/types'

// ==================== 用户认证 ====================

/**
 * 用户登录
 */
export function login(data: LoginRequest) {
  return request.post<LoginResponse>('/auth/login', data)
}

/**
 * 用户注册
 */
export function register(data: RegisterRequest) {
  return request.post<IdResponse>('/auth/register', data)
}

/**
 * 刷新 Token
 */
export function refreshToken(refreshToken: string) {
  return request.post<LoginResponse>('/auth/refresh', { refresh_token: refreshToken })
}

/**
 * 退出登录
 */
export function logout() {
  return request.post<StatusResponse>('/auth/logout')
}

// ==================== 用户管理 ====================

/**
 * 获取当前用户信息
 */
export function getCurrentUser() {
  return request.get<User>('/users/me')
}

/**
 * 更新用户信息
 */
export function updateUserInfo(data: Partial<User>) {
  return request.put<StatusResponse>('/users/me', data)
}

/**
 * 修改密码
 */
export function changePassword(data: { old_password: string; new_password: string }) {
  return request.post<StatusResponse>('/users/me/change-password', data)
}

// ==================== 工作空间管理 ====================

/**
 * 获取工作空间列表
 */
export function getWorkspaceList(params?: PaginationParams) {
  return request.get<PaginationResponse<Workspace>>('/workspaces', { params })
}

/**
 * 创建工作空间
 */
export function createWorkspace(data: { name: string; description?: string; type: string }) {
  return request.post<IdResponse>('/workspaces', data)
}

/**
 * 获取工作空间详情
 */
export function getWorkspaceDetail(id: number) {
  return request.get<Workspace>(`/workspaces/${id}`)
}

/**
 * 更新工作空间
 */
export function updateWorkspace(id: number, data: Partial<Workspace>) {
  return request.put<StatusResponse>(`/workspaces/${id}`, data)
}

/**
 * 删除工作空间
 */
export function deleteWorkspace(id: number) {
  return request.delete<StatusResponse>(`/workspaces/${id}`)
}

// ==================== 工作空间成员管理 ====================

/**
 * 获取工作空间成员列表
 */
export function getWorkspaceMembers(workspaceId: number, params?: PaginationParams) {
  return request.get<PaginationResponse<WorkspaceMember>>(`/workspaces/${workspaceId}/members`, { params })
}

/**
 * 添加工作空间成员
 */
export function addWorkspaceMember(workspaceId: number, data: { customer_id: number; role: string }) {
  return request.post<StatusResponse>(`/workspaces/${workspaceId}/members`, data)
}

/**
 * 更新成员角色
 */
export function updateMemberRole(workspaceId: number, memberId: number, role: string) {
  return request.put<StatusResponse>(`/workspaces/${workspaceId}/members/${memberId}`, { role })
}

/**
 * 移除工作空间成员
 */
export function removeWorkspaceMember(workspaceId: number, memberId: number) {
  return request.delete<StatusResponse>(`/workspaces/${workspaceId}/members/${memberId}`)
}

// ==================== 配额管理 ====================

/**
 * 获取用户配额
 */
export function getUserQuota() {
  return request.get<ResourceQuota>('/users/me/quota')
}

/**
 * 获取配额使用情况
 */
export function getQuotaUsage() {
  return request.get<QuotaUsage>('/users/me/quota/usage')
}

/**
 * 检查权限
 */
export function checkPermission(data: { resource_type: string; resource_id: number; action: string }) {
  return request.post<{ allowed: boolean }>('/auth/check-permission', data)
}
