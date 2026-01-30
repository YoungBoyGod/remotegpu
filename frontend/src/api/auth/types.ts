/**
 * 用户与权限模块 - 类型定义
 */

export type BackendUserRole = 'admin' | 'internal' | 'external'

// 后端返回的用户信息
export interface BackendUserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  role: BackendUserRole
  status: number
}

// 前端使用的用户信息
export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  role: 'admin' | 'customer'
  status: number
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
  user: BackendUserInfo
}

// 注册请求
export interface RegisterRequest {
  username: string
  email: string
  password: string
  nickname?: string
}

// 更新用户信息请求
export interface UpdateUserRequest {
  nickname?: string
  avatar?: string
}
