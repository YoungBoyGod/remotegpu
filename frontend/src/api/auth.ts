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
  mustChangePassword: boolean
}

interface BackendLoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  must_change_password: boolean
}

/**
 * 刷新Token响应
 */
export interface RefreshTokenResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  mustChangePassword: boolean
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
  return (request.post('/auth/login', data) as Promise<ApiResponse<BackendLoginResponse>>).then((res) => ({
    ...res,
    data: {
      accessToken: res.data.access_token,
      refreshToken: res.data.refresh_token,
      expiresIn: res.data.expires_in,
      mustChangePassword: res.data.must_change_password,
    },
  }))
}

/**
 * 用户登出
 */
export function logout(): Promise<ApiResponse<void>> {
  return request.post('/auth/logout')
}

/**
 * 刷新访问令牌
 */
export function refreshToken(refreshToken: string): Promise<ApiResponse<RefreshTokenResponse>> {
  return (request.post('/auth/refresh', { refresh_token: refreshToken }) as Promise<ApiResponse<BackendLoginResponse>>).then((res) => ({
    ...res,
    data: {
      accessToken: res.data.access_token,
      refreshToken: res.data.refresh_token,
      expiresIn: res.data.expires_in,
      mustChangePassword: res.data.must_change_password,
    },
  }))
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export function changePassword(data: ChangePasswordRequest): Promise<ApiResponse<void>> {
  return request.post('/auth/password/change', data)
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser(): Promise<ApiResponse<UserInfo>> {
  return request.get('/auth/profile')
}

/**
 * 请求密码重置
 */
export function requestPasswordReset(data: ResetPasswordRequest): Promise<ApiResponse<void>> {
  return request.post('/auth/password-reset/request', data)
}

/**
 * 提交新密码
 */
export function confirmPasswordReset(data: ConfirmPasswordResetRequest): Promise<ApiResponse<void>> {
  return request.post('/auth/password-reset/confirm', data)
}

/**
 * 验证Token有效性
 */
export function validateToken(): Promise<ApiResponse<{ valid: boolean }>> {
  return request.get('/auth/validate')
}
