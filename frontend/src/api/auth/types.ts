/**
 * 用户与权限模块 - 类型定义
 */
import type { Timestamps, UUIDField } from '../common/types'

// 用户信息
export interface User extends Timestamps, UUIDField {
  id: number
  username: string
  email: string
  phone?: string
  avatar?: string
  role: 'admin' | 'customer'
  account_type: 'individual' | 'enterprise'
  status: 'active' | 'inactive' | 'suspended'
  last_login_at?: string
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
  remember_me?: boolean
}

// 登录响应
export interface LoginResponse {
  token: string
  refresh_token: string
  expires_in: number
  user: User
}

// 注册请求
export interface RegisterRequest {
  username: string
  email: string
  password: string
  phone?: string
  verification_code?: string
}

// 工作空间
export interface Workspace extends Timestamps, UUIDField {
  id: number
  name: string
  description?: string
  owner_id: number
  owner_name?: string
  type: 'personal' | 'team' | 'enterprise'
  member_count: number
  status: 'active' | 'suspended'
}

// 工作空间成员
export interface WorkspaceMember extends Timestamps {
  id: number
  workspace_id: number
  customer_id: number
  username: string
  email: string
  role: 'owner' | 'admin' | 'member' | 'viewer'
  status: 'active' | 'pending' | 'inactive'
  joined_at: string
}

// 资源配额
export interface ResourceQuota {
  customer_id: number
  workspace_id?: number
  cpu_quota: number
  memory_quota: number
  gpu_quota: number
  storage_quota: number
  environment_quota: number
}

// 配额使用情况
export interface QuotaUsage extends ResourceQuota {
  cpu_used: number
  memory_used: number
  gpu_used: number
  storage_used: number
  environment_used: number
}
