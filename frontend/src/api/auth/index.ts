/**
 * 用户与权限模块 - API 接口
 */
import request from '../common/request'
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  UpdateUserRequest,
  BackendUserInfo
} from './types'
import type { StatusResponse } from '../common/types'

// ==================== 用户认证 ====================

/**
 * 用户登录
 */
export function login(data: LoginRequest) {
  return request.post<LoginResponse>('/user/login', data)
}

/**
 * 用户注册
 */
export function register(data: RegisterRequest) {
  return request.post<StatusResponse>('/user/register', data)
}

// ==================== 用户管理 ====================

/**
 * 获取当前用户信息
 */
export function getCurrentUser() {
  return request.get<BackendUserInfo>('/user/info')
}

/**
 * 更新用户信息
 */
export function updateUserInfo(data: UpdateUserRequest) {
  return request.put<StatusResponse>('/user/info', data)
}
