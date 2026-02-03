import request from '@/utils/request'
import type { ApiResponse } from '@/types/common'
import type { UserInfo } from '@/types/common'

/**
 * 登录请求参数
 */
export interface LoginRequest {
  username: string
  password: string
}

/**
 * 登录响应数据
 */
export interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  user: UserInfo
}

/**
 * 刷新Token响应
 */
export interface RefreshTokenResponse {
  accessToken: string
  expiresIn: number
}

/**
 * 密码重置请求
 */
export interface ResetPasswordRequest {
  email: string
}

/**
 * 确认密码重置请求
 */
export interface ConfirmPasswordResetRequest {
  email: string
  code: string
  password: string
}

/**
 * 用户登录
 */
export function login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  return request.post('/api/auth/login', data)
}

/**
 * 用户登出
 */
export function logout(): Promise<ApiResponse<void>> {
  return request.post('/api/auth/logout')
}

/**
 * 刷新访问令牌
 */
export function refreshToken(refreshToken: string): Promise<ApiResponse<RefreshTokenResponse>> {
  return request.post('/api/auth/refresh', { refreshToken })
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser(): Promise<ApiResponse<UserInfo>> {
  return request.get('/api/auth/profile')
}

/**
 * 请求密码重置
 */
export function requestPasswordReset(data: ResetPasswordRequest): Promise<ApiResponse<void>> {
  return request.post('/api/auth/password-reset/request', data)
}

/**
 * 提交新密码
 */
export function confirmPasswordReset(data: ConfirmPasswordResetRequest): Promise<ApiResponse<void>> {
  return request.post('/api/auth/password-reset/confirm', data)
}

/**
 * 验证Token有效性
 */
export function validateToken(): Promise<ApiResponse<{ valid: boolean }>> {
  return request.get('/api/auth/validate')
}
