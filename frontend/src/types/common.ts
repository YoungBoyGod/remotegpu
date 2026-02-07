/**
 * 通用类型定义
 */

// API 响应结构
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  traceId?: string
}

// 分页请求参数
export interface PageRequest {
  page: number
  pageSize: number
  sortBy?: string
  order?: 'asc' | 'desc'
  filters?: Record<string, any>
}

// 分页响应结构
export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 用户角色（与后端 openapi.yaml 一致）
export type UserRole = 'admin' | 'customer' | 'customer_owner' | 'customer_member'

// 用户信息
export interface UserInfo {
  id: number
  username: string
  email: string
  role: UserRole
  tenantId?: number
  tenantName?: string
  avatar?: string
  must_change_password?: boolean
}

// 登录响应
export interface LoginResponse {
  accessToken: string
  refreshToken: string
  user: UserInfo
}

// 统计卡片数据
export interface StatCardData {
  label: string
  value: number | string
  detail?: string
  icon?: string
  color?: string
}
