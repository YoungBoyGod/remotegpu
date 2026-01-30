/**
 * 环境管理模块 - API 接口
 */
import request from '../common/request'
import type {
  Environment,
  CreateEnvironmentRequest
} from './types'
import type { StatusResponse } from '../common/types'

// ==================== 环境管理 ====================

/**
 * 获取环境列表
 */
export function getEnvironmentList(params?: {
  workspace_id?: number
}) {
  return request.get<Environment[]>('/admin/environments', { params })
}

/**
 * 创建环境
 */
export function createEnvironment(data: CreateEnvironmentRequest) {
  return request.post<Environment>('/admin/environments', data)
}

/**
 * 获取环境详情
 */
export function getEnvironmentDetail(id: string) {
  return request.get<Environment>(`/admin/environments/${id}`)
}

/**
 * 启动环境
 */
export function startEnvironment(id: string) {
  return request.post<StatusResponse>(`/admin/environments/${id}/start`)
}

/**
 * 停止环境
 */
export function stopEnvironment(id: string) {
  return request.post<StatusResponse>(`/admin/environments/${id}/stop`)
}

/**
 * 重启环境
 */
export function deleteEnvironment(id: string) {
  return request.delete<StatusResponse>(`/admin/environments/${id}`)
}
